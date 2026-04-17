package tlcocoonTypes

import (
	"bytes"
	tl "github.com/tonkeeper/tongo/tl"
	"io"
)

type KeyStorageRunnerConfig struct {
	IsTest                    bool
	IsTestnet                 bool
	RpcPort                   int32
	HttpPort                  int32
	RootContractAddress       string
	MachineSpecificPrivateKey [32]byte
	TonConfigFilename         string
	DbPath                    string
	CheckHashes               bool
	ImageHash                 [32]byte
	HttpAccessHash            string
}

func (*KeyStorageRunnerConfig) CRC() uint32 {
	return uint32(0xd7874e40)
}
func (t KeyStorageRunnerConfig) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.IsTest)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.IsTestnet)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.RpcPort)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.HttpPort)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.RootContractAddress)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(t.MachineSpecificPrivateKey[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	b, err = tl.Marshal(t.TonConfigFilename)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.DbPath)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.CheckHashes)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(t.ImageHash[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	b, err = tl.Marshal(t.HttpAccessHash)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *KeyStorageRunnerConfig) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.IsTest)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.IsTestnet)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.RpcPort)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.HttpPort)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.RootContractAddress)
	if err != nil {
		return err
	}
	_, err = io.ReadFull(r, t.MachineSpecificPrivateKey[:])
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.TonConfigFilename)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.DbPath)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.CheckHashes)
	if err != nil {
		return err
	}
	_, err = io.ReadFull(r, t.ImageHash[:])
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.HttpAccessHash)
	if err != nil {
		return err
	}
	return nil
}
