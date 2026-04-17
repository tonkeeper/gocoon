package tlcocoonTypes

import (
	"bytes"
	"fmt"
	tl "github.com/tonkeeper/tongo/tl"
	"io"
)

type WorkerCompareBalanceWithProxyResult struct {
	SignedPayment       IProxySignedPayment
	TokensCommittedToDb int64
	MaxTokens           int64
	ErrorCode           int32
}

func (*WorkerCompareBalanceWithProxyResult) CRC() uint32 {
	return uint32(0x9054c411)
}
func (t WorkerCompareBalanceWithProxyResult) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	m4SignedPayment := t.SignedPayment
	if m4SignedPayment == nil {
		return nil, fmt.Errorf("nil %s", "ProxySignedPayment")
	}
	b, err = tl.Marshal(m4SignedPayment.CRC())
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = m4SignedPayment.MarshalTL()
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	_ = "IProxySignedPayment"
	b, err = tl.Marshal(t.TokensCommittedToDb)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.MaxTokens)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ErrorCode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *WorkerCompareBalanceWithProxyResult) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	tmp2SignedPayment, err := decodeIProxySignedPayment(r)
	if err != nil {
		return err
	}
	t.SignedPayment = tmp2SignedPayment
	err = tl.Unmarshal(r, &t.TokensCommittedToDb)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.MaxTokens)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ErrorCode)
	if err != nil {
		return err
	}
	return nil
}

type WorkerConnectedToProxy struct {
	Params          ProxyParams
	WorkerScAddress string
}

func (*WorkerConnectedToProxy) CRC() uint32 {
	return uint32(0x66bc82ee)
}
func (t WorkerConnectedToProxy) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.Params)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.WorkerScAddress)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *WorkerConnectedToProxy) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Params)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.WorkerScAddress)
	if err != nil {
		return err
	}
	return nil
}

type WorkerEnabledDisabled struct {
	Disabled int32
}

func (*WorkerEnabledDisabled) CRC() uint32 {
	return uint32(0xb3b71039)
}
func (t WorkerEnabledDisabled) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.Disabled)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *WorkerEnabledDisabled) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Disabled)
	if err != nil {
		return err
	}
	return nil
}

type WorkerExtendedCompareBalanceWithProxyResult struct {
	ErrorCode int32
}

func (*WorkerExtendedCompareBalanceWithProxyResult) CRC() uint32 {
	return uint32(0xc10973e)
}
func (t WorkerExtendedCompareBalanceWithProxyResult) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.ErrorCode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *WorkerExtendedCompareBalanceWithProxyResult) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.ErrorCode)
	if err != nil {
		return err
	}
	return nil
}

type WorkerProxyHandshakeCompleted struct{}

func (*WorkerProxyHandshakeCompleted) CRC() uint32 {
	return uint32(0x3ad9b886)
}
func (t WorkerProxyHandshakeCompleted) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	return buf.Bytes(), nil
}
func (t *WorkerProxyHandshakeCompleted) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	return nil
}

type WorkerNewCoefficient struct {
	NewCoefficient int32
}

func (*WorkerNewCoefficient) CRC() uint32 {
	return uint32(0x6af092b0)
}
func (t WorkerNewCoefficient) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.NewCoefficient)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *WorkerNewCoefficient) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.NewCoefficient)
	if err != nil {
		return err
	}
	return nil
}

type WorkerParams struct {
	Flags              uint32 `tl:"0,bitflag"`
	WorkerOwnerAddress string
	Model              string
	Coefficient        int32
	IsTest             *bool  `tl:",omitempty:Flags:0"`
	ProxyCnt           *int32 `tl:",omitempty:Flags:0"`
	MaxActiveRequests  *int32 `tl:",omitempty:Flags:0"`
	MinProtoVersion    *int32 `tl:",omitempty:Flags:1"`
	MaxProtoVersion    *int32 `tl:",omitempty:Flags:1"`
}

func (*WorkerParams) CRC() uint32 {
	return uint32(0x869c73ed)
}
func (t WorkerParams) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	flagsVar0 := t.Flags
	if t.IsTest != nil {
		flagsVar0 |= uint32(0x1)
	}
	if t.ProxyCnt != nil {
		flagsVar0 |= uint32(0x1)
	}
	if t.MaxActiveRequests != nil {
		flagsVar0 |= uint32(0x1)
	}
	if t.MinProtoVersion != nil {
		flagsVar0 |= uint32(0x2)
	}
	if t.MaxProtoVersion != nil {
		flagsVar0 |= uint32(0x2)
	}
	b, err = tl.Marshal(flagsVar0)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.WorkerOwnerAddress)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Model)
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
	if (flagsVar0>>int32(0))&1 == 1 {
		b, err = tl.Marshal(*t.IsTest)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (flagsVar0>>int32(0))&1 == 1 {
		b, err = tl.Marshal(*t.ProxyCnt)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (flagsVar0>>int32(0))&1 == 1 {
		b, err = tl.Marshal(*t.MaxActiveRequests)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (flagsVar0>>int32(1))&1 == 1 {
		b, err = tl.Marshal(*t.MinProtoVersion)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (flagsVar0>>int32(1))&1 == 1 {
		b, err = tl.Marshal(*t.MaxProtoVersion)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
func (t *WorkerParams) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Flags)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.WorkerOwnerAddress)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Model)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Coefficient)
	if err != nil {
		return err
	}
	if (t.Flags>>int32(0))&1 == 1 {
		var tmp10IsTest bool
		err = tl.Unmarshal(r, &tmp10IsTest)
		if err != nil {
			return err
		}
		t.IsTest = &tmp10IsTest
	}
	if (t.Flags>>int32(0))&1 == 1 {
		var tmp11ProxyCnt int32
		err = tl.Unmarshal(r, &tmp11ProxyCnt)
		if err != nil {
			return err
		}
		t.ProxyCnt = &tmp11ProxyCnt
	}
	if (t.Flags>>int32(0))&1 == 1 {
		var tmp12MaxActiveRequests int32
		err = tl.Unmarshal(r, &tmp12MaxActiveRequests)
		if err != nil {
			return err
		}
		t.MaxActiveRequests = &tmp12MaxActiveRequests
	}
	if (t.Flags>>int32(1))&1 == 1 {
		var tmp13MinProtoVersion int32
		err = tl.Unmarshal(r, &tmp13MinProtoVersion)
		if err != nil {
			return err
		}
		t.MinProtoVersion = &tmp13MinProtoVersion
	}
	if (t.Flags>>int32(1))&1 == 1 {
		var tmp14MaxProtoVersion int32
		err = tl.Unmarshal(r, &tmp14MaxProtoVersion)
		if err != nil {
			return err
		}
		t.MaxProtoVersion = &tmp14MaxProtoVersion
	}
	return nil
}

type WorkerPaymentStatus struct {
	SignedPayment IProxySignedPayment
	DbTokens      int64
	MaxTokens     int64
}

func (*WorkerPaymentStatus) CRC() uint32 {
	return uint32(0xa88bd701)
}
func (t WorkerPaymentStatus) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	m4SignedPayment := t.SignedPayment
	if m4SignedPayment == nil {
		return nil, fmt.Errorf("nil %s", "ProxySignedPayment")
	}
	b, err = tl.Marshal(m4SignedPayment.CRC())
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = m4SignedPayment.MarshalTL()
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	_ = "IProxySignedPayment"
	b, err = tl.Marshal(t.DbTokens)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.MaxTokens)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *WorkerPaymentStatus) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	tmp2SignedPayment, err := decodeIProxySignedPayment(r)
	if err != nil {
		return err
	}
	t.SignedPayment = tmp2SignedPayment
	err = tl.Unmarshal(r, &t.DbTokens)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.MaxTokens)
	if err != nil {
		return err
	}
	return nil
}
