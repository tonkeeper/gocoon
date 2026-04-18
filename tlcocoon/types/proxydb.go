package tlcocoonTypes

import (
	"bytes"
	"fmt"
	tl "github.com/tonkeeper/tongo/tl"
	"io"
)

type IProxyDbClientInfo interface {
	CRC() uint32
	MarshalTL() ([]byte, error)
	UnmarshalTL(io.Reader) error
	_IProxyDbClientInfo()
}

var (
	_ IProxyDbClientInfo = (*ProxyDbClientInfo)(nil)
	_ IProxyDbClientInfo = (*ProxyDbClientInfoV2)(nil)
)

func decodeIProxyDbClientInfo(r io.Reader) (IProxyDbClientInfo, error) {
	var tag uint32
	err := tl.Unmarshal(r, &tag)
	if err != nil {
		return nil, err
	}
	var res IProxyDbClientInfo
	switch tag {
	case uint32(0x6cfeaa21):
		res = &ProxyDbClientInfo{}
	case uint32(0x3078c9bf):
		res = &ProxyDbClientInfoV2{}
	default:
		return nil, fmt.Errorf("invalid crc code: got 0x%08x", tag)
	}
	err = res.UnmarshalTL(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func DecodeIProxyDbClientInfo(r io.Reader) (IProxyDbClientInfo, error) {
	return decodeIProxyDbClientInfo(r)
}

type ProxyDbClientInfo struct {
	OwnerAddress  string
	ScAddress     string
	Status        int32
	Balance       int64
	ScTokensUsed  int64
	TokensUsed    int64
	SecretHash    [32]byte
	LastRequestAt int32
}

func (*ProxyDbClientInfo) CRC() uint32 {
	return uint32(0x6cfeaa21)
}
func (t ProxyDbClientInfo) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.OwnerAddress)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ScAddress)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Status)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Balance)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ScTokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.TokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(t.SecretHash[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	b, err = tl.Marshal(t.LastRequestAt)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ProxyDbClientInfo) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.OwnerAddress)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ScAddress)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Status)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Balance)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ScTokensUsed)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.TokensUsed)
	if err != nil {
		return err
	}
	_, err = io.ReadFull(r, t.SecretHash[:])
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.LastRequestAt)
	if err != nil {
		return err
	}
	return nil
}
func (*ProxyDbClientInfo) _IProxyDbClientInfo() {}

type ProxyDbClientInfoV2 struct {
	OwnerAddress  string
	ScAddress     string
	Status        int32
	Balance       int64
	Stake         int64
	ScTokensUsed  int64
	TokensUsed    int64
	SecretHash    [32]byte
	LastRequestAt int32
}

func (*ProxyDbClientInfoV2) CRC() uint32 {
	return uint32(0x3078c9bf)
}
func (t ProxyDbClientInfoV2) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.OwnerAddress)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ScAddress)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Status)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Balance)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Stake)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ScTokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.TokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(t.SecretHash[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	b, err = tl.Marshal(t.LastRequestAt)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ProxyDbClientInfoV2) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.OwnerAddress)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ScAddress)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Status)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Balance)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Stake)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ScTokensUsed)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.TokensUsed)
	if err != nil {
		return err
	}
	_, err = io.ReadFull(r, t.SecretHash[:])
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.LastRequestAt)
	if err != nil {
		return err
	}
	return nil
}
func (*ProxyDbClientInfoV2) _IProxyDbClientInfo() {}

type IProxyDbConfig interface {
	CRC() uint32
	MarshalTL() ([]byte, error)
	UnmarshalTL(io.Reader) error
	_IProxyDbConfig()
}

var (
	_ IProxyDbConfig = (*ProxyDbConfigEmpty)(nil)
	_ IProxyDbConfig = (*ProxyDbConfigV4)(nil)
	_ IProxyDbConfig = (*ProxyDbConfigV4Disabled)(nil)
)

func decodeIProxyDbConfig(r io.Reader) (IProxyDbConfig, error) {
	var tag uint32
	err := tl.Unmarshal(r, &tag)
	if err != nil {
		return nil, err
	}
	var res IProxyDbConfig
	switch tag {
	case uint32(0x43dc40b4):
		res = &ProxyDbConfigEmpty{}
	case uint32(0xbde43703):
		res = &ProxyDbConfigV4{}
	case uint32(0xad2cfb4a):
		res = &ProxyDbConfigV4Disabled{}
	default:
		return nil, fmt.Errorf("invalid crc code: got 0x%08x", tag)
	}
	err = res.UnmarshalTL(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func DecodeIProxyDbConfig(r io.Reader) (IProxyDbConfig, error) {
	return decodeIProxyDbConfig(r)
}

type ProxyDbConfigEmpty struct{}

func (*ProxyDbConfigEmpty) CRC() uint32 {
	return uint32(0x43dc40b4)
}
func (t ProxyDbConfigEmpty) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	return buf.Bytes(), nil
}
func (t *ProxyDbConfigEmpty) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	return nil
}
func (*ProxyDbConfigEmpty) _IProxyDbConfig() {}

type ProxyDbConfigV4 struct {
	RootContractAddress            string
	IsTestnet                      int32
	RootContractState              IRootConfigConfig
	RootContractStateBlockTs       int32
	ScBlockID                      TonBlockIDExt
	LastSeqnoCommittedToBlockchain int32
	PendingBlockchainSeqnoCommits  []ProxyDbPendingBlockchainSeqnoCommit
}

func (*ProxyDbConfigV4) CRC() uint32 {
	return uint32(0xbde43703)
}
func (t ProxyDbConfigV4) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.RootContractAddress)
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
	m12RootContractState := t.RootContractState
	if m12RootContractState == nil {
		return nil, fmt.Errorf("nil %s", "RootConfigConfig")
	}
	b, err = tl.Marshal(m12RootContractState.CRC())
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = m12RootContractState.MarshalTL()
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	_ = "IRootConfigConfig"
	b, err = tl.Marshal(t.RootContractStateBlockTs)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ScBlockID)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.LastSeqnoCommittedToBlockchain)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.PendingBlockchainSeqnoCommits)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ProxyDbConfigV4) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.RootContractAddress)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.IsTestnet)
	if err != nil {
		return err
	}
	tmp6RootContractState, err := decodeIRootConfigConfig(r)
	if err != nil {
		return err
	}
	t.RootContractState = tmp6RootContractState
	err = tl.Unmarshal(r, &t.RootContractStateBlockTs)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ScBlockID)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.LastSeqnoCommittedToBlockchain)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.PendingBlockchainSeqnoCommits)
	if err != nil {
		return err
	}
	return nil
}
func (*ProxyDbConfigV4) _IProxyDbConfig() {}

type ProxyDbConfigV4Disabled struct {
	RootContractAddress      string
	IsTestnet                int32
	RootContractState        IRootConfigConfig
	RootContractStateBlockTs int32
	DisabledUntilVersion     int64
}

func (*ProxyDbConfigV4Disabled) CRC() uint32 {
	return uint32(0xad2cfb4a)
}
func (t ProxyDbConfigV4Disabled) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.RootContractAddress)
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
	m12RootContractState := t.RootContractState
	if m12RootContractState == nil {
		return nil, fmt.Errorf("nil %s", "RootConfigConfig")
	}
	b, err = tl.Marshal(m12RootContractState.CRC())
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = m12RootContractState.MarshalTL()
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	_ = "IRootConfigConfig"
	b, err = tl.Marshal(t.RootContractStateBlockTs)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.DisabledUntilVersion)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ProxyDbConfigV4Disabled) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.RootContractAddress)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.IsTestnet)
	if err != nil {
		return err
	}
	tmp6RootContractState, err := decodeIRootConfigConfig(r)
	if err != nil {
		return err
	}
	t.RootContractState = tmp6RootContractState
	err = tl.Unmarshal(r, &t.RootContractStateBlockTs)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.DisabledUntilVersion)
	if err != nil {
		return err
	}
	return nil
}
func (*ProxyDbConfigV4Disabled) _IProxyDbConfig() {}

type ProxyDbOldClient struct {
	OwnerAddress string
	Tokens       int64
	NextClient   string
}

func (*ProxyDbOldClient) CRC() uint32 {
	return uint32(0x2fae7b42)
}
func (t ProxyDbOldClient) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.OwnerAddress)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Tokens)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.NextClient)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ProxyDbOldClient) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.OwnerAddress)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Tokens)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.NextClient)
	if err != nil {
		return err
	}
	return nil
}

type ProxyDbOldInstance struct {
	ContractAddress   string
	ClosingState      int32
	CloseAt           int32
	NextClient        string
	NextWorker        string
	RootContractState IRootConfigConfig
}

func (*ProxyDbOldInstance) CRC() uint32 {
	return uint32(0xa8c9e178)
}
func (t ProxyDbOldInstance) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.ContractAddress)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ClosingState)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.CloseAt)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.NextClient)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.NextWorker)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	m24RootContractState := t.RootContractState
	if m24RootContractState == nil {
		return nil, fmt.Errorf("nil %s", "RootConfigConfig")
	}
	b, err = tl.Marshal(m24RootContractState.CRC())
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = m24RootContractState.MarshalTL()
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	_ = "IRootConfigConfig"
	return buf.Bytes(), nil
}
func (t *ProxyDbOldInstance) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.ContractAddress)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ClosingState)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.CloseAt)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.NextClient)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.NextWorker)
	if err != nil {
		return err
	}
	tmp12RootContractState, err := decodeIRootConfigConfig(r)
	if err != nil {
		return err
	}
	t.RootContractState = tmp12RootContractState
	return nil
}

type ProxyDbOldWorker struct {
	OwnerAddress string
	Tokens       int64
	NextWorker   string
}

func (*ProxyDbOldWorker) CRC() uint32 {
	return uint32(0xf04156a0)
}
func (t ProxyDbOldWorker) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.OwnerAddress)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Tokens)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.NextWorker)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ProxyDbOldWorker) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.OwnerAddress)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Tokens)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.NextWorker)
	if err != nil {
		return err
	}
	return nil
}

type ProxyDbPendingBlockchainSeqnoCommit struct {
	Seqno       int32
	SessionHash [32]byte
}

func (*ProxyDbPendingBlockchainSeqnoCommit) CRC() uint32 {
	return uint32(0x525b14a9)
}
func (t ProxyDbPendingBlockchainSeqnoCommit) MarshalTL() ([]byte, error) {
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
	_, err = buf.Write(t.SessionHash[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	return buf.Bytes(), nil
}
func (t *ProxyDbPendingBlockchainSeqnoCommit) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Seqno)
	if err != nil {
		return err
	}
	_, err = io.ReadFull(r, t.SessionHash[:])
	if err != nil {
		return err
	}
	return nil
}

type ProxyDbWorkerInfo struct {
	OwnerAddress  string
	ScAddress     string
	ScTokens      int64
	Tokens        int64
	LastRequestAt int32
}

func (*ProxyDbWorkerInfo) CRC() uint32 {
	return uint32(0x2517b38f)
}
func (t ProxyDbWorkerInfo) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.OwnerAddress)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ScAddress)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ScTokens)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Tokens)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.LastRequestAt)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ProxyDbWorkerInfo) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.OwnerAddress)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ScAddress)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ScTokens)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Tokens)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.LastRequestAt)
	if err != nil {
		return err
	}
	return nil
}
