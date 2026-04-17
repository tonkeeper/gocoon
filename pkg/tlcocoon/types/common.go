package tlcocoonTypes

import (
	"bytes"
	tl "github.com/tonkeeper/tongo/tl"
	"io"
)

type Function struct{}

func (*Function) CRC() uint32 {
	return uint32(0x7acbc197)
}
func (t Function) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	return buf.Bytes(), nil
}
func (t *Function) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	return nil
}

type Object struct{}

func (*Object) CRC() uint32 {
	return uint32(0x29704ca0)
}
func (t Object) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	return buf.Bytes(), nil
}
func (t *Object) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	return nil
}

type TokensUsed struct {
	PromptTokensUsed     int64
	CachedTokensUsed     int64
	CompletionTokensUsed int64
	ReasoningTokensUsed  int64
	TotalTokensUsed      int64
}

func (*TokensUsed) CRC() uint32 {
	return uint32(0x70c5b15c)
}
func (t TokensUsed) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.PromptTokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.CachedTokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.CompletionTokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ReasoningTokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.TotalTokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *TokensUsed) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.PromptTokensUsed)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.CachedTokensUsed)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.CompletionTokensUsed)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ReasoningTokensUsed)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.TotalTokensUsed)
	if err != nil {
		return err
	}
	return nil
}
