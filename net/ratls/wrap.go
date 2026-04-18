package ratls

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
	"sync"
	"time"
)

// Conn wraps a net.Conn and upgrades it to mutual TLS 1.3 on first I/O.
type Conn struct {
	net.Conn
	once         sync.Once
	handshakeErr error
	tlsConn      *tls.Conn
	serverCert   []byte
	clientCert   *tls.Certificate
}

// Wrap returns a net.Conn that transparently completes RA-TLS on first I/O.
func Wrap(conn net.Conn) *Conn {
	return &Conn{Conn: conn}
}

// WrapWithCert is like Wrap, but uses a caller-provided client certificate.
func WrapWithCert(conn net.Conn, cert tls.Certificate) net.Conn {
	return &Conn{Conn: conn, clientCert: &cert}
}

// Handshake forces TLS handshake eagerly.
func (c *Conn) Handshake() error {
	c.once.Do(func() {
		var clientCert tls.Certificate
		if c.clientCert != nil {
			clientCert = *c.clientCert
		} else {
			var err error
			clientCert, err = generateEphemeralCert()
			if err != nil {
				c.handshakeErr = fmt.Errorf("generate cert: %w", err)
				_ = c.Conn.Close()
				return
			}
		}

		tlsCfg := &tls.Config{
			MinVersion:             tls.VersionTLS13,
			InsecureSkipVerify:     true,
			Certificates:           []tls.Certificate{clientCert},
			SessionTicketsDisabled: true,
		}

		tlsConn := tls.Client(c.Conn, tlsCfg)
		if err := tlsConn.Handshake(); err != nil {
			c.handshakeErr = fmt.Errorf("tls handshake: %w", err)
			_ = c.Conn.Close()
			return
		}

		c.tlsConn = tlsConn
		if certs := tlsConn.ConnectionState().PeerCertificates; len(certs) > 0 {
			c.serverCert = certs[0].Raw
		}
	})
	return c.handshakeErr
}

func (c *Conn) Read(b []byte) (int, error) {
	if err := c.Handshake(); err != nil {
		return 0, err
	}
	return c.tlsConn.Read(b)
}

func (c *Conn) Write(b []byte) (int, error) {
	if err := c.Handshake(); err != nil {
		return 0, err
	}
	return c.tlsConn.Write(b)
}

func (c *Conn) Close() error {
	if c.tlsConn != nil {
		return c.tlsConn.Close()
	}
	return c.Conn.Close()
}

// ServerCert returns server raw DER certificate bytes after successful handshake.
func (c *Conn) ServerCert() []byte {
	if len(c.serverCert) == 0 {
		return nil
	}
	out := make([]byte, len(c.serverCert))
	copy(out, c.serverCert)
	return out
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
		Subject:      pkix.Name{CommonName: "tonkeeper/gocoon"},
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
