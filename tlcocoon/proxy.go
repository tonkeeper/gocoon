package tlcocoon

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	types "github.com/tonkeeper/gocoon/tlcocoon/types"
	tl "github.com/tonkeeper/tongo/tl"
	"math"
)

type ProxyRunQueryRequest struct {
	Query         []byte
	SignedPayment types.IProxySignedPayment
	Coefficient   int64
	Timeout       float64
	RequestID     [32]byte
}

func (*ProxyRunQueryRequest) CRC() uint32 {
	return uint32(0x47182416)
}
func (t ProxyRunQueryRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.Query)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	m8SignedPayment := t.SignedPayment
	if m8SignedPayment == nil {
		return nil, fmt.Errorf("nil %s", "ProxySignedPayment")
	}
	b, err = tl.Marshal(m8SignedPayment.CRC())
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = m8SignedPayment.MarshalTL()
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	_ = "IProxySignedPayment"
	b, err = tl.Marshal(t.Coefficient)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	m23TimeoutBits := math.Float64bits(t.Timeout)
	var m23TimeoutRaw [8]byte
	binary.LittleEndian.PutUint64(m23TimeoutRaw[:], m23TimeoutBits)
	_, err = buf.Write(m23TimeoutRaw[:])
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(t.RequestID[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	return buf.Bytes(), nil
}
func ProxyRunQuery(ctx context.Context, m Requester, i ProxyRunQueryRequest) (types.IProxyQueryAnswer, error) {
	respRaw, err := requestRaw(ctx, m, &i)
	if err != nil {
		var zero types.IProxyQueryAnswer
		return zero, err
	}
	res, err := types.DecodeIProxyQueryAnswer(bytes.NewReader(respRaw))
	if err != nil {
		var zero types.IProxyQueryAnswer
		return zero, fmt.Errorf("response: %w", err)
	}
	return res, nil
}

type ProxyRunQueryExRequest struct {
	Query         []byte
	SignedPayment types.IProxySignedPayment
	Coefficient   int64
	Timeout       float64
	RequestID     [32]byte
	Flags         uint32 `tl:"0,bitflag"`
	EnableDebug   *bool  `tl:",omitempty:Flags:0"`
}

func (*ProxyRunQueryExRequest) CRC() uint32 {
	return uint32(0xc805e7e2)
}
func (t ProxyRunQueryExRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	flagsVar0 := t.Flags
	if t.EnableDebug != nil {
		flagsVar0 |= uint32(0x1)
	}
	b, err = tl.Marshal(t.Query)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	m10SignedPayment := t.SignedPayment
	if m10SignedPayment == nil {
		return nil, fmt.Errorf("nil %s", "ProxySignedPayment")
	}
	b, err = tl.Marshal(m10SignedPayment.CRC())
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = m10SignedPayment.MarshalTL()
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	_ = "IProxySignedPayment"
	b, err = tl.Marshal(t.Coefficient)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	m25TimeoutBits := math.Float64bits(t.Timeout)
	var m25TimeoutRaw [8]byte
	binary.LittleEndian.PutUint64(m25TimeoutRaw[:], m25TimeoutBits)
	_, err = buf.Write(m25TimeoutRaw[:])
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(t.RequestID[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	b, err = tl.Marshal(flagsVar0)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	if (flagsVar0>>int32(0))&1 == 1 {
		b, err = tl.Marshal(*t.EnableDebug)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
func ProxyRunQueryEx(ctx context.Context, m Requester, i ProxyRunQueryExRequest) (types.IProxyQueryAnswer, error) {
	respRaw, err := requestRaw(ctx, m, &i)
	if err != nil {
		var zero types.IProxyQueryAnswer
		return zero, err
	}
	res, err := types.DecodeIProxyQueryAnswer(bytes.NewReader(respRaw))
	if err != nil {
		var zero types.IProxyQueryAnswer
		return zero, fmt.Errorf("response: %w", err)
	}
	return res, nil
}

type ProxyAPI interface {
	RunQuery(ctx context.Context, i ProxyRunQueryRequest) (types.IProxyQueryAnswer, error)
	RunQueryEx(ctx context.Context, i ProxyRunQueryExRequest) (types.IProxyQueryAnswer, error)
}
type Proxy struct {
	requester Requester
}

func NewProxy(requester Requester) *Proxy {
	return &Proxy{requester: requester}
}
func (c *Proxy) RunQuery(ctx context.Context, i ProxyRunQueryRequest) (types.IProxyQueryAnswer, error) {
	return ProxyRunQuery(ctx, c.requester, i)
}
func (c *Proxy) RunQueryEx(ctx context.Context, i ProxyRunQueryExRequest) (types.IProxyQueryAnswer, error) {
	return ProxyRunQueryEx(ctx, c.requester, i)
}

var _ ProxyAPI = (*Proxy)(nil)
