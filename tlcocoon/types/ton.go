package tlcocoonTypes

import (
	"bytes"
	tl "github.com/tonkeeper/tongo/tl"
	"io"
)

type TonBlockIDExt struct {
	Workchain int32
	Shard     int64
	Seqno     int32
	RootHash  []byte
	FileHash  []byte
}

func (*TonBlockIDExt) CRC() uint32 {
	return uint32(0xbc3f6da5)
}
func (t TonBlockIDExt) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.Workchain)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Shard)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Seqno)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.RootHash)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.FileHash)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *TonBlockIDExt) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Workchain)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Shard)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Seqno)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.RootHash)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.FileHash)
	if err != nil {
		return err
	}
	return nil
}
