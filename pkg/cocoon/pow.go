package cocoon

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"math/bits"
	"time"
)

const (
	powSimpleMagic         uint32 = 0x418e1291
	powSimpleResponseMagic uint32 = 0x01827319
)

// PowChallenge is the 24-byte challenge sent by the server.
type PowChallenge struct {
	Difficulty int32
	Salt       [16]byte
}

// readPowChallenge reads and parses the 24-byte PoW challenge from the server.
func readPowChallenge(r io.Reader) (PowChallenge, error) {
	var buf [24]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return PowChallenge{}, fmt.Errorf("read pow challenge: %w", err)
	}
	magic := binary.LittleEndian.Uint32(buf[0:4])
	if magic != powSimpleMagic {
		return PowChallenge{}, fmt.Errorf("pow: unexpected magic 0x%08x (expected 0x%08x)", magic, powSimpleMagic)
	}
	var c PowChallenge
	c.Difficulty = int32(binary.LittleEndian.Uint32(buf[4:8]))
	copy(c.Salt[:], buf[8:24])
	return c, nil
}

// sendPowResponse writes the 12-byte PoW response (magic + nonce).
func sendPowResponse(w io.Writer, nonce int64) error {
	var buf [12]byte
	binary.LittleEndian.PutUint32(buf[0:4], powSimpleResponseMagic)
	binary.LittleEndian.PutUint64(buf[4:12], uint64(nonce))
	_, err := w.Write(buf[:])
	return err
}

// leadingZeroBits counts leading zero bits exactly as the C++ server does.
//
// The server reads the first 8 bytes of the SHA256 hash as a little-endian
// uint64, then calls count_leading_zeroes64 (i.e. counts from the MSB of
// that integer). On x86 (little-endian), this means bit 63 is hash[7]'s MSB.
func leadingZeroBits(h []byte) int {
	v := binary.LittleEndian.Uint64(h[:8])
	return bits.LeadingZeros64(v)
}

// solvePow finds a nonce such that SHA256(salt || nonce_le64) has at least
// difficulty leading zero bits (using the server's little-endian convention).
func solvePow(c PowChallenge) (int64, time.Duration) {
	start := time.Now()
	var data [24]byte
	copy(data[:16], c.Salt[:])
	var nonce int64
	for {
		binary.LittleEndian.PutUint64(data[16:], uint64(nonce))
		h := sha256.Sum256(data[:])
		if leadingZeroBits(h[:]) >= int(c.Difficulty) {
			return nonce, time.Since(start)
		}
		nonce++
	}
}
