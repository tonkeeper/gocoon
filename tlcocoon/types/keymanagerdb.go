package tlcocoonTypes

import (
	"bytes"
	"fmt"
	tl "github.com/tonkeeper/tongo/tl"
	"io"
)

type IKeyManagerDbConfig interface {
	CRC() uint32
	MarshalTL() ([]byte, error)
	UnmarshalTL(io.Reader) error
	_IKeyManagerDbConfig()
}

var (
	_ IKeyManagerDbConfig = (*KeyManagerDbConfigEmpty)(nil)
	_ IKeyManagerDbConfig = (*KeyManagerDbConfigV1)(nil)
)

func decodeIKeyManagerDbConfig(r io.Reader) (IKeyManagerDbConfig, error) {
	var tag uint32
	err := tl.Unmarshal(r, &tag)
	if err != nil {
		return nil, err
	}
	var res IKeyManagerDbConfig
	switch tag {
	case uint32(0xd3fc2385):
		res = &KeyManagerDbConfigEmpty{}
	case uint32(0x6d2817ac):
		res = &KeyManagerDbConfigV1{}
	default:
		return nil, fmt.Errorf("invalid crc code: got 0x%08x", tag)
	}
	err = res.UnmarshalTL(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func DecodeIKeyManagerDbConfig(r io.Reader) (IKeyManagerDbConfig, error) {
	return decodeIKeyManagerDbConfig(r)
}

type KeyManagerDbConfigEmpty struct{}

func (*KeyManagerDbConfigEmpty) CRC() uint32 {
	return uint32(0xd3fc2385)
}
func (t KeyManagerDbConfigEmpty) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	return buf.Bytes(), nil
}
func (t *KeyManagerDbConfigEmpty) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	return nil
}
func (*KeyManagerDbConfigEmpty) _IKeyManagerDbConfig() {}

type KeyManagerDbConfigV1 struct {
	RootContractVersion int32
}

func (*KeyManagerDbConfigV1) CRC() uint32 {
	return uint32(0x6d2817ac)
}
func (t KeyManagerDbConfigV1) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.RootContractVersion)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *KeyManagerDbConfigV1) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.RootContractVersion)
	if err != nil {
		return err
	}
	return nil
}
func (*KeyManagerDbConfigV1) _IKeyManagerDbConfig() {}

type KeyManagerDbKey struct {
	PrivateKey              [32]byte
	ForProxies              bool
	ForWorkers              bool
	ValidSinceConfigVersion int32
	ValidSinceUtime         int32
	ValidUntilUtime         int32
}

func (*KeyManagerDbKey) CRC() uint32 {
	return uint32(0x4d0a1e17)
}
func (t KeyManagerDbKey) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	_, err = buf.Write(t.PrivateKey[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	b, err = tl.Marshal(t.ForProxies)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ForWorkers)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ValidSinceConfigVersion)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ValidSinceUtime)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ValidUntilUtime)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *KeyManagerDbKey) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	_, err = io.ReadFull(r, t.PrivateKey[:])
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ForProxies)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ForWorkers)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ValidSinceConfigVersion)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ValidSinceUtime)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ValidUntilUtime)
	if err != nil {
		return err
	}
	return nil
}
