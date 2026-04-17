package tlcocoonTypes

import (
	"bytes"
	"fmt"
	tl "github.com/tonkeeper/tongo/tl"
	"io"
)

type IRootConfigConfig interface {
	CRC() uint32
	MarshalTL() ([]byte, error)
	UnmarshalTL(io.Reader) error
	_IRootConfigConfig()
}

var (
	_ IRootConfigConfig = (*RootConfigPseudo)(nil)
	_ IRootConfigConfig = (*RootConfigConfigV5)(nil)
)

func decodeIRootConfigConfig(r io.Reader) (IRootConfigConfig, error) {
	var tag uint32
	err := tl.Unmarshal(r, &tag)
	if err != nil {
		return nil, err
	}
	var res IRootConfigConfig
	switch tag {
	case uint32(0x87a976cc):
		res = &RootConfigPseudo{}
	case uint32(0x5e84869b):
		res = &RootConfigConfigV5{}
	default:
		return nil, fmt.Errorf("invalid crc code: got 0x%08x", tag)
	}
	err = res.UnmarshalTL(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func DecodeIRootConfigConfig(r io.Reader) (IRootConfigConfig, error) {
	return decodeIRootConfigConfig(r)
}

type RootConfigPseudo struct {
	ProxyHashes       [][32]byte
	WorkerHashes      [][32]byte
	ModelHashes       [][32]byte
	ModelTypes        []string
	RegisteredProxies []RootConfigRegisteredProxy
	PricePerToken     int32
	WorkerFeePerToken int32
	LastProxySeqno    int32
	Version           int32
	ParamsVersion     int32
	ProxyScCode       string
	WorkerScCode      string
	ClientScCode      string
	RootOwnerAddress  string
}

func (*RootConfigPseudo) CRC() uint32 {
	return uint32(0x87a976cc)
}
func (t RootConfigPseudo) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.ProxyHashes)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.WorkerHashes)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ModelHashes)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ModelTypes)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.RegisteredProxies)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.PricePerToken)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.WorkerFeePerToken)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.LastProxySeqno)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Version)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ParamsVersion)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ProxyScCode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.WorkerScCode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ClientScCode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.RootOwnerAddress)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *RootConfigPseudo) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.ProxyHashes)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.WorkerHashes)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ModelHashes)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ModelTypes)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.RegisteredProxies)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.PricePerToken)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.WorkerFeePerToken)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.LastProxySeqno)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Version)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ParamsVersion)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ProxyScCode)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.WorkerScCode)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ClientScCode)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.RootOwnerAddress)
	if err != nil {
		return err
	}
	return nil
}
func (*RootConfigPseudo) _IRootConfigConfig() {}

type RootConfigConfigV5 struct {
	RootOwnerAddress                string
	ProxyHashes                     [][32]byte
	RegisteredProxies               []RootConfigRegisteredProxy
	LastProxySeqno                  int32
	WorkerHashes                    [][32]byte
	ModelHashes                     [][32]byte
	Version                         int32
	StructVersion                   int32
	ParamsVersion                   int32
	UniqueID                        int32
	IsTest                          int32
	PricePerToken                   int64
	WorkerFeePerToken               int64
	PromptTokensPriceMultiplier     int32
	CachedTokensPriceMultiplier     int32
	CompletionTokensPriceMultiplier int32
	ReasoningTokensPriceMultiplier  int32
	ProxyDelayBeforeClose           int32
	ClientDelayBeforeClose          int32
	MinProxyStake                   int64
	MinClientStake                  int64
	ProxyScCode                     string
	WorkerScCode                    string
	ClientScCode                    string
}

func (*RootConfigConfigV5) CRC() uint32 {
	return uint32(0x5e84869b)
}
func (t RootConfigConfigV5) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.RootOwnerAddress)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ProxyHashes)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.RegisteredProxies)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.LastProxySeqno)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.WorkerHashes)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ModelHashes)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Version)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.StructVersion)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ParamsVersion)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.UniqueID)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.IsTest)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.PricePerToken)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.WorkerFeePerToken)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.PromptTokensPriceMultiplier)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.CachedTokensPriceMultiplier)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.CompletionTokensPriceMultiplier)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ReasoningTokensPriceMultiplier)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ProxyDelayBeforeClose)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ClientDelayBeforeClose)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.MinProxyStake)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.MinClientStake)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ProxyScCode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.WorkerScCode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ClientScCode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *RootConfigConfigV5) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.RootOwnerAddress)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ProxyHashes)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.RegisteredProxies)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.LastProxySeqno)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.WorkerHashes)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ModelHashes)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Version)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.StructVersion)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ParamsVersion)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.UniqueID)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.IsTest)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.PricePerToken)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.WorkerFeePerToken)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.PromptTokensPriceMultiplier)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.CachedTokensPriceMultiplier)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.CompletionTokensPriceMultiplier)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ReasoningTokensPriceMultiplier)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ProxyDelayBeforeClose)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ClientDelayBeforeClose)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.MinProxyStake)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.MinClientStake)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ProxyScCode)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.WorkerScCode)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ClientScCode)
	if err != nil {
		return err
	}
	return nil
}
func (*RootConfigConfigV5) _IRootConfigConfig() {}

type RootConfigRegisteredProxy struct {
	Seqno   int32
	Address string
}

func (*RootConfigRegisteredProxy) CRC() uint32 {
	return uint32(0x9f4e446a)
}
func (t RootConfigRegisteredProxy) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.Seqno)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Address)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *RootConfigRegisteredProxy) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Seqno)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Address)
	if err != nil {
		return err
	}
	return nil
}
