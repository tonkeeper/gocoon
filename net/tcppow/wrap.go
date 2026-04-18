package tcppow

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"math/bits"
	"net"
	"sync"
	"time"
)

const (
	powSimpleMagic         uint32 = 0x418e1291
	powSimpleResponseMagic uint32 = 0x01827319
)

// Conn wraps a raw TCP connection and completes the PoW challenge/response
// before the first application read/write.
type Conn struct {
	net.Conn
	once         sync.Once
	handshakeErr error
	difficulty   int32
	timeSpent    time.Duration
}

// Wrap returns a net.Conn that transparently solves cocoon PoW on first I/O.
func Wrap(conn net.Conn) net.Conn {
	return &Conn{Conn: conn}
}

// Handshake forces PoW exchange eagerly.
func (c *Conn) Handshake() error {
	c.once.Do(func() {
		challenge, err := readChallenge(c.Conn)
		if err != nil {
			c.handshakeErr = err
			_ = c.Conn.Close()
			return
		}
		start := time.Now()
		nonce := solve(challenge)
		c.timeSpent = time.Since(start)
		c.difficulty = challenge.difficulty
		if err := sendResponse(c.Conn, nonce); err != nil {
			c.handshakeErr = fmt.Errorf("send pow response: %w", err)
			_ = c.Conn.Close()
		}
	})
	return c.handshakeErr
}

func (c *Conn) Difficulty() int32 {
	return c.difficulty
}

func (c *Conn) TimeSpent() time.Duration {
	return c.timeSpent
}

func (c *Conn) Read(b []byte) (int, error) {
	if err := c.Handshake(); err != nil {
		return 0, err
	}
	return c.Conn.Read(b)
}

func (c *Conn) Write(b []byte) (int, error) {
	if err := c.Handshake(); err != nil {
		return 0, err
	}
	return c.Conn.Write(b)
}

type challenge struct {
	difficulty int32
	salt       [16]byte
}

func readChallenge(r io.Reader) (challenge, error) {
	var buf [24]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return challenge{}, fmt.Errorf("read pow challenge: %w", err)
	}
	magic := binary.LittleEndian.Uint32(buf[0:4])
	if magic != powSimpleMagic {
		return challenge{}, fmt.Errorf("pow: unexpected magic 0x%08x (expected 0x%08x)", magic, powSimpleMagic)
	}
	var c challenge
	c.difficulty = int32(binary.LittleEndian.Uint32(buf[4:8]))
	copy(c.salt[:], buf[8:24])
	return c, nil
}

func sendResponse(w io.Writer, nonce int64) error {
	var buf [12]byte
	binary.LittleEndian.PutUint32(buf[0:4], powSimpleResponseMagic)
	binary.LittleEndian.PutUint64(buf[4:12], uint64(nonce))
	_, err := w.Write(buf[:])
	return err
}

func solve(c challenge) int64 {
	var data [24]byte
	copy(data[:16], c.salt[:])
	var nonce int64
	for {
		binary.LittleEndian.PutUint64(data[16:], uint64(nonce))
		h := sha256.Sum256(data[:])
		if leadingZeroBits(h[:]) >= int(c.difficulty) {
			return nonce
		}
		nonce++
	}
}

func leadingZeroBits(h []byte) int {
	v := binary.LittleEndian.Uint64(h[:8])
	return bits.LeadingZeros64(v)
}
