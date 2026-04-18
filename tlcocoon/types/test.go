package tlcocoonTypes

import (
	"bytes"
	tl "github.com/tonkeeper/tongo/tl"
	"io"
)

type TestAnswer struct {
	ID int64
}

func (*TestAnswer) CRC() uint32 {
	return uint32(0xd4fcdb6)
}
func (t TestAnswer) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.ID)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *TestAnswer) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.ID)
	if err != nil {
		return err
	}
	return nil
}
