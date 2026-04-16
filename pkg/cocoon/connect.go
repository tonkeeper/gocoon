// Package cocoon implements the client-side cocoon proxy connection protocol:
// TCP → PoW challenge/response → mutual TLS 1.3 (Ed25519 self-signed cert).
package cocoon

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"

	"go.uber.org/zap"
)

// Conn is an established cocoon proxy connection (post-PoW, post-TLS).
type Conn struct {
	*tls.Conn
	// ServerCert is the server's raw DER certificate (contains TDX quote).
	ServerCert []byte
}

// Dial connects to a cocoon proxy at addr (e.g. "91.108.4.11:8888"),
// solves the PoW challenge, and completes mutual TLS with an ephemeral
// Ed25519 self-signed certificate.
//
// The server uses policy "any" for client certificates, so no TDX
// attestation is required — any valid Ed25519 self-signed cert works.
func Dial(addr string, logger *zap.Logger) (*Conn, error) {
	clientCert, err := generateEphemeralCert()
	if err != nil {
		return nil, fmt.Errorf("generate cert: %w", err)
	}
	return DialWithCert(addr, clientCert, logger)
}

// DialWithCert is like Dial but uses a pre-generated TLS certificate.
func DialWithCert(addr string, clientCert tls.Certificate, logger *zap.Logger) (*Conn, error) {
	logger.Debug("connecting to cocoon proxy", zap.String("addr", addr))
	tcp, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return nil, fmt.Errorf("tcp dial %s: %w", addr, err)
	}

	// Step 1: solve PoW
	logger.Debug("reading PoW challenge")
	challenge, err := readPowChallenge(tcp)
	if err != nil {
		tcp.Close()
		return nil, err
	}
	logger.Debug("got PoW challenge", zap.Any("challenge", challenge))
	nonce, timeSpent := solvePow(challenge)
	logger.Debug("solving PoW response",
		zap.Duration("spent", timeSpent),
		zap.Int32("complexity", challenge.Difficulty),
	)
	if err := sendPowResponse(tcp, nonce); err != nil {
		tcp.Close()
		return nil, fmt.Errorf("send pow response: %w", err)
	}

	// Step 2: mutual TLS 1.3
	// InsecureSkipVerify disables the CA chain check (cert is self-signed).
	// The server's identity is verified via TDX attestation embedded in the
	// X.509 extensions — handled separately by VerifyServerCert if needed.
	tlsCfg := &tls.Config{
		MinVersion:             tls.VersionTLS13,
		InsecureSkipVerify:     true,
		Certificates:           []tls.Certificate{clientCert},
		SessionTicketsDisabled: true,
	}
	tlsConn := tls.Client(tcp, tlsCfg)
	if err := tlsConn.Handshake(); err != nil {
		tcp.Close()
		return nil, fmt.Errorf("tls handshake: %w", err)
	}

	var serverCert []byte
	if certs := tlsConn.ConnectionState().PeerCertificates; len(certs) > 0 {
		serverCert = certs[0].Raw
	}

	return &Conn{Conn: tlsConn, ServerCert: serverCert}, nil
}

// generateEphemeralCert creates a throwaway Ed25519 self-signed cert.
// The server's verifier requires Ed25519 — RSA/ECDSA are rejected.
func generateEphemeralCert() (tls.Certificate, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return tls.Certificate{}, err
	}

	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: "cocoon-client"},
		NotBefore:    time.Now().Add(-time.Minute),
		NotAfter:     time.Now().Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, pub, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return tls.Certificate{}, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})

	return tls.X509KeyPair(certPEM, keyPEM)
}
