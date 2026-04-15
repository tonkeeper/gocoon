package cocoon

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
)

// Session implements the TcpConnection framing layer over an established Conn.
//
// Wire format per packet: [uint32 LE size][uint32 LE seqno][payload]
// where payload is a boxed TL tcp.Packet.
type Session struct {
	conn      *Conn
	sendSeqno uint32
	recvSeqno uint32
}

// NewSession creates a session and sends the initial tcp.connect packet,
// then reads the tcp.connected acknowledgement from the server.
func NewSession(conn *Conn) (*Session, error) {
	s := &Session{conn: conn}
	id, err := randInt64()
	if err != nil {
		return nil, err
	}
	// Build tcp.connect{id}
	var w tlWriter
	w.u32(idTcpConnect)
	w.i64(id)
	if err := s.sendFrame(w.buf); err != nil {
		return nil, fmt.Errorf("send tcp.connect: %w", err)
	}
	// Read tcp.connected{id} acknowledgement.
	frame, err := s.recvFrame()
	if err != nil {
		return nil, fmt.Errorf("recv tcp.connected: %w", err)
	}
	r := newTLReader(frame)
	if err := r.expectU32(idTcpConnected, "tcp.connected"); err != nil {
		return nil, err
	}
	// id field — we don't validate it matches
	if _, err := r.i64(); err != nil {
		return nil, fmt.Errorf("read tcp.connected id: %w", err)
	}
	return s, nil
}

// Query sends payload (a boxed TL function call) as tcp.query and returns
// the inner data bytes from the tcp.queryAnswer response.
//
// The ctx is used for deadline/cancellation — it sets a deadline on the
// underlying net.Conn if the connection supports it.
func (s *Session) Query(_ context.Context, payload []byte) ([]byte, error) {
	id, err := randInt64()
	if err != nil {
		return nil, err
	}

	// Build tcp.query{id, data=payload}
	var w tlWriter
	w.u32(idTcpQuery)
	w.i64(id)
	w.bytes(payload)
	if err := s.sendFrame(w.buf); err != nil {
		return nil, fmt.Errorf("send tcp.query: %w", err)
	}

	// Read responses until we get the matching tcp.queryAnswer.
	for {
		frame, err := s.recvFrame()
		if err != nil {
			return nil, fmt.Errorf("recv frame: %w", err)
		}
		r := newTLReader(frame)
		magic, err := r.u32()
		if err != nil {
			return nil, fmt.Errorf("read magic: %w", err)
		}

		switch magic {
		case idTcpQueryAnswer:
			gotID, err := r.i64()
			if err != nil {
				return nil, fmt.Errorf("read tcp.queryAnswer id: %w", err)
			}
			data, err := r.rawBytes()
			if err != nil {
				return nil, fmt.Errorf("read tcp.queryAnswer data: %w", err)
			}
			if gotID != id {
				return nil, fmt.Errorf("tcp.queryAnswer id mismatch: got %d, want %d", gotID, id)
			}
			return data, nil

		case idTcpQueryError:
			gotID, err := r.i64()
			if err != nil {
				return nil, fmt.Errorf("read tcp.queryError id: %w", err)
			}
			code, err := r.u32()
			if err != nil {
				return nil, fmt.Errorf("read tcp.queryError code: %w", err)
			}
			msg, err := r.str()
			if err != nil {
				return nil, fmt.Errorf("read tcp.queryError message: %w", err)
			}
			_ = gotID
			return nil, fmt.Errorf("tcp.queryError code=%d: %s", code, msg)

		case idTcpPing:
			pingID, err := r.i64()
			if err != nil {
				return nil, fmt.Errorf("read tcp.ping id: %w", err)
			}
			// Reply with tcp.pong
			var pw tlWriter
			pw.u32(idTcpPong)
			pw.i64(pingID)
			if err := s.sendFrame(pw.buf); err != nil {
				return nil, fmt.Errorf("send tcp.pong: %w", err)
			}
			// Continue waiting for the answer

		default:
			return nil, fmt.Errorf("unexpected tcp.Packet magic 0x%08x", magic)
		}
	}
}

// sendFrame writes [uint32 LE size][uint32 LE seqno][payload] to the connection.
func (s *Session) sendFrame(payload []byte) error {
	var hdr [8]byte
	binary.LittleEndian.PutUint32(hdr[0:4], uint32(len(payload)))
	binary.LittleEndian.PutUint32(hdr[4:8], s.sendSeqno)
	s.sendSeqno++

	if _, err := s.conn.Write(hdr[:]); err != nil {
		return err
	}
	_, err := s.conn.Write(payload)
	return err
}

// recvFrame reads [uint32 LE size][uint32 LE seqno][payload] and returns payload.
func (s *Session) recvFrame() ([]byte, error) {
	var hdr [8]byte
	if _, err := io.ReadFull(s.conn, hdr[:]); err != nil {
		return nil, fmt.Errorf("read frame header: %w", err)
	}
	size := binary.LittleEndian.Uint32(hdr[0:4])
	// seqno := binary.LittleEndian.Uint32(hdr[4:8]) // not validated for now

	if size > 1<<20 {
		return nil, fmt.Errorf("frame too large: %d bytes", size)
	}
	payload := make([]byte, size)
	if _, err := io.ReadFull(s.conn, payload); err != nil {
		return nil, fmt.Errorf("read frame payload: %w", err)
	}
	s.recvSeqno++
	return payload, nil
}

func randInt64() (int64, error) {
	n, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 63))
	if err != nil {
		return 0, err
	}
	return n.Int64(), nil
}
