package tlcocoonTypes

import (
	"bytes"
	tl "github.com/tonkeeper/tongo/tl"
	"io"
)

type KeyManagerPrivateKey struct {
	ValidUntilUtime int32
	PrivateKey      [32]byte
}

func (*KeyManagerPrivateKey) CRC() uint32 {
	return uint32(0x60f4f8e7)
}
func (t KeyManagerPrivateKey) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.ValidUntilUtime)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(t.PrivateKey[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	return buf.Bytes(), nil
}
func (t *KeyManagerPrivateKey) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.ValidUntilUtime)
	if err != nil {
		return err
	}
	_, err = io.ReadFull(r, t.PrivateKey[:])
	if err != nil {
		return err
	}
	return nil
}

type KeyManagerPrivateKeys struct {
	Keys []KeyManagerPrivateKey
}

func (*KeyManagerPrivateKeys) CRC() uint32 {
	return uint32(0xed66ce5)
}
func (t KeyManagerPrivateKeys) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.Keys)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *KeyManagerPrivateKeys) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Keys)
	if err != nil {
		return err
	}
	return nil
}
