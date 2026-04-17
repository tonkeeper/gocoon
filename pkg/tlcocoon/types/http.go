package tlcocoonTypes

import (
	"bytes"
	tl "github.com/tonkeeper/tongo/tl"
	"io"
)

type HttpHeader struct {
	Name  string
	Value string
}

func (*HttpHeader) CRC() uint32 {
	return uint32(0x8e9be511)
}
func (t HttpHeader) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.Name)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Value)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *HttpHeader) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Name)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Value)
	if err != nil {
		return err
	}
	return nil
}

type HttpResponse struct {
	HttpVersion string
	StatusCode  int32
	Reason      string
	Headers     []HttpHeader
	Payload     []byte
}

func (*HttpResponse) CRC() uint32 {
	return uint32(0x1cd0c42b)
}
func (t HttpResponse) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.HttpVersion)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.StatusCode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Reason)
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
func (t *HttpResponse) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.HttpVersion)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.StatusCode)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Reason)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Headers)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Payload)
	if err != nil {
		return err
	}
	return nil
}
