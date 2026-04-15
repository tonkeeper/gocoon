package cocoon

import (
	"encoding/binary"
	"fmt"
)

// TL type IDs (magic numbers), computed as CRC32 of the TL definition string.
const (
	idTcpConnect    uint32 = 0xa57c4261
	idTcpConnected  uint32 = 0x636d41d6
	idTcpPing       uint32 = 0xbbe9627c
	idTcpPong       uint32 = 0x0bd4302c
	idTcpQuery      uint32 = 0x3af51908
	idTcpQueryAnswer uint32 = 0xc048c311
	idTcpQueryError uint32 = 0x4cd2f602

	// client.params#40fdca64 (explicit in schema)
	idClientParams uint32 = 0x40fdca64
	// proxy.params#d5c5609f (explicit in schema)
	idProxyParams uint32 = 0xd5c5609f

	idClientConnectToProxy   uint32 = 0xff5fa0f4
	idClientConnectedToProxy uint32 = 0x95317ad1

	idClientProxyConnectionAuthShort uint32 = 0xd6ffc5af
	idClientProxyConnectionAuthLong  uint32 = 0x417bf016

	idClientAuthorizeWithProxyShort uint32 = 0x6c276723
	idClientAuthorizeWithProxyLong  uint32 = 0xd3474303

	idClientAuthorizationWithProxySuccess uint32 = 0x75d5ac34
	idClientAuthorizationWithProxyFailed  uint32 = 0x60551c96

	idProxySignedPayment      uint32 = 0x02998182
	idProxySignedPaymentEmpty uint32 = 0xb347ce64

	idBoolTrue  uint32 = 0x997275b5
	idBoolFalse uint32 = 0xbc799737
)

// ── TL writer ────────────────────────────────────────────────────────────────

type tlWriter struct{ buf []byte }

func (w *tlWriter) u32(v uint32) {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], v)
	w.buf = append(w.buf, b[:]...)
}

func (w *tlWriter) i64(v int64) {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(v))
	w.buf = append(w.buf, b[:]...)
}

// bytes encodes a TL bytes/string field.
// Short form (len < 254): 1-byte length + data + padding to 4-byte alignment.
// Long form: 0xFE + 3-byte LE length + data + padding.
func (w *tlWriter) bytes(data []byte) {
	n := len(data)
	if n < 254 {
		w.buf = append(w.buf, byte(n))
		w.buf = append(w.buf, data...)
		if pad := (4 - (1+n)%4) % 4; pad > 0 {
			w.buf = append(w.buf, make([]byte, pad)...)
		}
	} else {
		w.buf = append(w.buf, 0xFE, byte(n), byte(n>>8), byte(n>>16))
		w.buf = append(w.buf, data...)
		if pad := (4 - n%4) % 4; pad > 0 {
			w.buf = append(w.buf, make([]byte, pad)...)
		}
	}
}

func (w *tlWriter) str(s string) { w.bytes([]byte(s)) }

func (w *tlWriter) boolean(v bool) {
	if v {
		w.u32(idBoolTrue)
	} else {
		w.u32(idBoolFalse)
	}
}

func (w *tlWriter) raw(b []byte) { w.buf = append(w.buf, b...) }

// ── TL reader ────────────────────────────────────────────────────────────────

type tlReader struct {
	buf []byte
	pos int
}

func newTLReader(b []byte) *tlReader { return &tlReader{buf: b} }

func (r *tlReader) remaining() int { return len(r.buf) - r.pos }

func (r *tlReader) need(n int) error {
	if r.remaining() < n {
		return fmt.Errorf("tl: need %d bytes, have %d", n, r.remaining())
	}
	return nil
}

func (r *tlReader) u32() (uint32, error) {
	if err := r.need(4); err != nil {
		return 0, err
	}
	v := binary.LittleEndian.Uint32(r.buf[r.pos:])
	r.pos += 4
	return v, nil
}

func (r *tlReader) i64() (int64, error) {
	if err := r.need(8); err != nil {
		return 0, err
	}
	v := binary.LittleEndian.Uint64(r.buf[r.pos:])
	r.pos += 8
	return int64(v), nil
}

func (r *tlReader) expectU32(want uint32, name string) error {
	got, err := r.u32()
	if err != nil {
		return err
	}
	if got != want {
		return fmt.Errorf("tl: %s: expected magic 0x%08x, got 0x%08x", name, want, got)
	}
	return nil
}

// rawBytes reads TL bytes/string, returning a copy of the payload (no padding).
func (r *tlReader) rawBytes() ([]byte, error) {
	if err := r.need(1); err != nil {
		return nil, err
	}
	first := r.buf[r.pos]
	r.pos++
	var length int
	if first < 254 {
		length = int(first)
		// 1-byte len + data + pad, total must be multiple of 4
		total := 1 + length
		pad := (4 - total%4) % 4
		if err := r.need(length + pad); err != nil {
			return nil, err
		}
		data := make([]byte, length)
		copy(data, r.buf[r.pos:])
		r.pos += length + pad
		return data, nil
	}
	if first == 254 {
		if err := r.need(3); err != nil {
			return nil, err
		}
		length = int(r.buf[r.pos]) | int(r.buf[r.pos+1])<<8 | int(r.buf[r.pos+2])<<16
		r.pos += 3
		pad := (4 - length%4) % 4
		if err := r.need(length + pad); err != nil {
			return nil, err
		}
		data := make([]byte, length)
		copy(data, r.buf[r.pos:])
		r.pos += length + pad
		return data, nil
	}
	return nil, fmt.Errorf("tl: unsupported bytes prefix 0x%02x", first)
}

func (r *tlReader) str() (string, error) {
	b, err := r.rawBytes()
	return string(b), err
}

func (r *tlReader) boolean() (bool, error) {
	m, err := r.u32()
	if err != nil {
		return false, err
	}
	switch m {
	case idBoolTrue:
		return true, nil
	case idBoolFalse:
		return false, nil
	default:
		return false, fmt.Errorf("tl: bool: unexpected magic 0x%08x", m)
	}
}

func (r *tlReader) fixedBytes(n int) ([]byte, error) {
	if err := r.need(n); err != nil {
		return nil, err
	}
	b := make([]byte, n)
	copy(b, r.buf[r.pos:])
	r.pos += n
	return b, nil
}
