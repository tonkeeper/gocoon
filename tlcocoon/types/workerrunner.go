package tlcocoonTypes

import (
	"bytes"
	tl "github.com/tonkeeper/tongo/tl"
	"io"
)

type WorkerRunnerConfig struct {
	IsTest              bool
	IsTestnet           bool
	HttpPort            int32
	OwnerAddress        string
	ProxyConnections    int32
	ConnectToProxyVia   string
	ModelName           string
	Coefficient         int32
	RootContractAddress string
	NodeWalletKey       [32]byte
	ForwardRequestsTo   string
	TonConfigFilename   string
	CheckProxyHashes    bool
	MaxActiveRequests   int32
	ImageHash           [32]byte
	HttpAccessHash      string
}

func (*WorkerRunnerConfig) CRC() uint32 {
	return uint32(0x870db56d)
}
func (t WorkerRunnerConfig) MarshalTL() ([]byte, error) {
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
	b, err = tl.Marshal(t.HttpPort)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.OwnerAddress)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ProxyConnections)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ConnectToProxyVia)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ModelName)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Coefficient)
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
	_, err = buf.Write(t.NodeWalletKey[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	b, err = tl.Marshal(t.ForwardRequestsTo)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.TonConfigFilename)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.CheckProxyHashes)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.MaxActiveRequests)
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
func (t *WorkerRunnerConfig) UnmarshalTL(r io.Reader) error {
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
	err = tl.Unmarshal(r, &t.HttpPort)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.OwnerAddress)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ProxyConnections)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ConnectToProxyVia)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ModelName)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Coefficient)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.RootContractAddress)
	if err != nil {
		return err
	}
	_, err = io.ReadFull(r, t.NodeWalletKey[:])
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ForwardRequestsTo)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.TonConfigFilename)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.CheckProxyHashes)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.MaxActiveRequests)
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
