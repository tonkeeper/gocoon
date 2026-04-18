package tlcocoon

import (
	"bytes"
	"context"
	types "github.com/tonkeeper/gocoon/tlcocoon/types"
	tl "github.com/tonkeeper/tongo/tl"
)

type HttpRequestRequest struct {
	Method      string
	URL         string
	HttpVersion string
	Headers     []types.HttpHeader
	Payload     []byte
}

func (*HttpRequestRequest) CRC() uint32 {
	return uint32(0x47492de5)
}
func (t HttpRequestRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.Method)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.URL)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.HttpVersion)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Headers)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Payload)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func HttpRequest(ctx context.Context, m Requester, i HttpRequestRequest) (types.HttpResponse, error) {
	var res types.HttpResponse
	return res, request(ctx, m, &i, &res)
}

type HttpAPI interface {
	Request(ctx context.Context, i HttpRequestRequest) (types.HttpResponse, error)
}
type Http struct {
	requester Requester
}

func NewHttp(requester Requester) *Http {
	return &Http{requester: requester}
}
func (c *Http) Request(ctx context.Context, i HttpRequestRequest) (types.HttpResponse, error) {
	return HttpRequest(ctx, c.requester, i)
}

var _ HttpAPI = (*Http)(nil)
