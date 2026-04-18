package proxyconn

import (
	"fmt"
	"net"
	"time"

	"github.com/tonkeeper/gocoon/net/ratls"
	"github.com/tonkeeper/gocoon/net/tcppow"
	"go.uber.org/zap"
)

// Conn is an established cocoon proxy connection (post-PoW, post-TLS).
type Conn struct {
	net.Conn
	// ServerCert is the server's raw DER certificate (contains TDX quote).
	ServerCert []byte
}

// Dial connects to a cocoon proxy at addr (e.g. "91.108.4.11:8888"),
// solves the PoW challenge, and completes mutual TLS with an ephemeral
// Ed25519 self-signed certificate.
func Dial(addr string, logger *zap.Logger) (*Conn, error) {

	logger.Debug("connecting to cocoon proxy", zap.String("addr", addr))
	tcp, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return nil, fmt.Errorf("tcp dial %s: %w", addr, err)
	}

	// Step 1: solve PoW
	powConn, ok := tcppow.Wrap(tcp).(*tcppow.Conn)
	if !ok {
		tcp.Close()
		return nil, fmt.Errorf("pow wrap returned unexpected type")
	}
	logger.Debug("solving PoW response")
	if err := powConn.Handshake(); err != nil {
		return nil, err
	}
	logger.Debug("PoW solved",
		zap.Int32("difficulty", powConn.Difficulty()),
		zap.Int64("time_spent_ms", powConn.TimeSpent().Milliseconds()),
	)

	// Step 2: mutual TLS 1.3 (RA-TLS)
	tlsConn := ratls.Wrap(powConn)
	if err := tlsConn.Handshake(); err != nil {
		return nil, err
	}

	// connection is ready for TL :-)
	return &Conn{Conn: tlsConn, ServerCert: tlsConn.ServerCert()}, nil
}
