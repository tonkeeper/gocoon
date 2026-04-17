package tlcocoonTypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	tl "github.com/tonkeeper/tongo/tl"
	"io"
	"math"
)

type IClientAuthorizationWithProxy interface {
	CRC() uint32
	MarshalTL() ([]byte, error)
	UnmarshalTL(io.Reader) error
	_IClientAuthorizationWithProxy()
}

var (
	_ IClientAuthorizationWithProxy = (*ClientAuthorizationWithProxySuccess)(nil)
	_ IClientAuthorizationWithProxy = (*ClientAuthorizationWithProxyFailed)(nil)
)

func decodeIClientAuthorizationWithProxy(r io.Reader) (IClientAuthorizationWithProxy, error) {
	var tag uint32
	err := tl.Unmarshal(r, &tag)
	if err != nil {
		return nil, err
	}
	var res IClientAuthorizationWithProxy
	switch tag {
	case uint32(0x75d5ac34):
		res = &ClientAuthorizationWithProxySuccess{}
	case uint32(0x60551c96):
		res = &ClientAuthorizationWithProxyFailed{}
	default:
		return nil, fmt.Errorf("invalid crc code: got 0x%08x", tag)
	}
	err = res.UnmarshalTL(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func DecodeIClientAuthorizationWithProxy(r io.Reader) (IClientAuthorizationWithProxy, error) {
	return decodeIClientAuthorizationWithProxy(r)
}

type ClientAuthorizationWithProxySuccess struct {
	SignedPayment       IProxySignedPayment
	TokensCommittedToDb int64
	MaxTokens           int64
}

func (*ClientAuthorizationWithProxySuccess) CRC() uint32 {
	return uint32(0x75d5ac34)
}
func (t ClientAuthorizationWithProxySuccess) MarshalTL() ([]byte, error) {
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
	return buf.Bytes(), nil
}
func (t *ClientAuthorizationWithProxySuccess) UnmarshalTL(r io.Reader) error {
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
	return nil
}
func (*ClientAuthorizationWithProxySuccess) _IClientAuthorizationWithProxy() {}

type ClientAuthorizationWithProxyFailed struct {
	ErrorCode int32
	Error     string
}

func (*ClientAuthorizationWithProxyFailed) CRC() uint32 {
	return uint32(0x60551c96)
}
func (t ClientAuthorizationWithProxyFailed) MarshalTL() ([]byte, error) {
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
	b, err = tl.Marshal(t.Error)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ClientAuthorizationWithProxyFailed) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.ErrorCode)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Error)
	if err != nil {
		return err
	}
	return nil
}
func (*ClientAuthorizationWithProxyFailed) _IClientAuthorizationWithProxy() {}

type ClientConnectedToProxy struct {
	Params          ProxyParams
	ClientScAddress string
	Auth            IClientProxyConnectionAuth
	SignedPayment   IProxySignedPayment
}

func (*ClientConnectedToProxy) CRC() uint32 {
	return uint32(0x95317ad1)
}
func (t ClientConnectedToProxy) MarshalTL() ([]byte, error) {
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
	b, err = tl.Marshal(t.ClientScAddress)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	m12Auth := t.Auth
	if m12Auth == nil {
		return nil, fmt.Errorf("nil %s", "ClientProxyConnectionAuth")
	}
	b, err = tl.Marshal(m12Auth.CRC())
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = m12Auth.MarshalTL()
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	_ = "IClientProxyConnectionAuth"
	m23SignedPayment := t.SignedPayment
	if m23SignedPayment == nil {
		return nil, fmt.Errorf("nil %s", "ProxySignedPayment")
	}
	b, err = tl.Marshal(m23SignedPayment.CRC())
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = m23SignedPayment.MarshalTL()
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	_ = "IProxySignedPayment"
	return buf.Bytes(), nil
}
func (t *ClientConnectedToProxy) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Params)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ClientScAddress)
	if err != nil {
		return err
	}
	tmp6Auth, err := decodeIClientProxyConnectionAuth(r)
	if err != nil {
		return err
	}
	t.Auth = tmp6Auth
	tmp9SignedPayment, err := decodeIProxySignedPayment(r)
	if err != nil {
		return err
	}
	t.SignedPayment = tmp9SignedPayment
	return nil
}

type ClientParams struct {
	Flags              uint32 `tl:"0,bitflag"`
	ClientOwnerAddress string
	IsTest             *bool  `tl:",omitempty:Flags:0"`
	MinProtoVersion    *int32 `tl:",omitempty:Flags:1"`
	MaxProtoVersion    *int32 `tl:",omitempty:Flags:1"`
}

func (*ClientParams) CRC() uint32 {
	return uint32(0x40fdca64)
}
func (t ClientParams) MarshalTL() ([]byte, error) {
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
	b, err = tl.Marshal(t.ClientOwnerAddress)
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
func (t *ClientParams) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Flags)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ClientOwnerAddress)
	if err != nil {
		return err
	}
	if (t.Flags>>int32(0))&1 == 1 {
		var tmp6IsTest bool
		err = tl.Unmarshal(r, &tmp6IsTest)
		if err != nil {
			return err
		}
		t.IsTest = &tmp6IsTest
	}
	if (t.Flags>>int32(1))&1 == 1 {
		var tmp7MinProtoVersion int32
		err = tl.Unmarshal(r, &tmp7MinProtoVersion)
		if err != nil {
			return err
		}
		t.MinProtoVersion = &tmp7MinProtoVersion
	}
	if (t.Flags>>int32(1))&1 == 1 {
		var tmp8MaxProtoVersion int32
		err = tl.Unmarshal(r, &tmp8MaxProtoVersion)
		if err != nil {
			return err
		}
		t.MaxProtoVersion = &tmp8MaxProtoVersion
	}
	return nil
}

type ClientPaymentStatus struct {
	SignedPayment IProxySignedPayment
	DbTokens      int64
	MaxTokens     int64
}

func (*ClientPaymentStatus) CRC() uint32 {
	return uint32(0xaa8a0ecc)
}
func (t ClientPaymentStatus) MarshalTL() ([]byte, error) {
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
func (t *ClientPaymentStatus) UnmarshalTL(r io.Reader) error {
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

type IClientProxyConnectionAuth interface {
	CRC() uint32
	MarshalTL() ([]byte, error)
	UnmarshalTL(io.Reader) error
	_IClientProxyConnectionAuth()
}

var (
	_ IClientProxyConnectionAuth = (*ClientProxyConnectionAuthShort)(nil)
	_ IClientProxyConnectionAuth = (*ClientProxyConnectionAuthLong)(nil)
)

func decodeIClientProxyConnectionAuth(r io.Reader) (IClientProxyConnectionAuth, error) {
	var tag uint32
	err := tl.Unmarshal(r, &tag)
	if err != nil {
		return nil, err
	}
	var res IClientProxyConnectionAuth
	switch tag {
	case uint32(0xd6ffc5af):
		res = &ClientProxyConnectionAuthShort{}
	case uint32(0x417bf016):
		res = &ClientProxyConnectionAuthLong{}
	default:
		return nil, fmt.Errorf("invalid crc code: got 0x%08x", tag)
	}
	err = res.UnmarshalTL(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func DecodeIClientProxyConnectionAuth(r io.Reader) (IClientProxyConnectionAuth, error) {
	return decodeIClientProxyConnectionAuth(r)
}

type ClientProxyConnectionAuthShort struct {
	SecretHash [32]byte
	Nonce      int64
}

func (*ClientProxyConnectionAuthShort) CRC() uint32 {
	return uint32(0xd6ffc5af)
}
func (t ClientProxyConnectionAuthShort) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	_, err = buf.Write(t.SecretHash[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	b, err = tl.Marshal(t.Nonce)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ClientProxyConnectionAuthShort) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	_, err = io.ReadFull(r, t.SecretHash[:])
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Nonce)
	if err != nil {
		return err
	}
	return nil
}
func (*ClientProxyConnectionAuthShort) _IClientProxyConnectionAuth() {}

type ClientProxyConnectionAuthLong struct {
	Nonce int64
}

func (*ClientProxyConnectionAuthLong) CRC() uint32 {
	return uint32(0x417bf016)
}
func (t ClientProxyConnectionAuthLong) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.Nonce)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ClientProxyConnectionAuthLong) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Nonce)
	if err != nil {
		return err
	}
	return nil
}
func (*ClientProxyConnectionAuthLong) _IClientProxyConnectionAuth() {}

type IClientQueryAnswer interface {
	CRC() uint32
	MarshalTL() ([]byte, error)
	UnmarshalTL(io.Reader) error
	_IClientQueryAnswer()
}

var (
	_ IClientQueryAnswer = (*ClientQueryAnswer)(nil)
	_ IClientQueryAnswer = (*ClientQueryAnswerError)(nil)
)

func decodeIClientQueryAnswer(r io.Reader) (IClientQueryAnswer, error) {
	var tag uint32
	err := tl.Unmarshal(r, &tag)
	if err != nil {
		return nil, err
	}
	var res IClientQueryAnswer
	switch tag {
	case uint32(0x9b943922):
		res = &ClientQueryAnswer{}
	case uint32(0x6d60569a):
		res = &ClientQueryAnswerError{}
	default:
		return nil, fmt.Errorf("invalid crc code: got 0x%08x", tag)
	}
	err = res.UnmarshalTL(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func DecodeIClientQueryAnswer(r io.Reader) (IClientQueryAnswer, error) {
	return decodeIClientQueryAnswer(r)
}

type ClientQueryAnswer struct {
	Answer            []byte
	IsCompleted       bool
	RequestID         [32]byte
	RequestTokensUsed TokensUsed
}

func (*ClientQueryAnswer) CRC() uint32 {
	return uint32(0x9b943922)
}
func (t ClientQueryAnswer) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.Answer)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.IsCompleted)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(t.RequestID[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	b, err = tl.Marshal(t.RequestTokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ClientQueryAnswer) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Answer)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.IsCompleted)
	if err != nil {
		return err
	}
	_, err = io.ReadFull(r, t.RequestID[:])
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.RequestTokensUsed)
	if err != nil {
		return err
	}
	return nil
}
func (*ClientQueryAnswer) _IClientQueryAnswer() {}

type ClientQueryAnswerError struct {
	ErrorCode         int32
	Error             string
	RequestID         [32]byte
	RequestTokensUsed TokensUsed
}

func (*ClientQueryAnswerError) CRC() uint32 {
	return uint32(0x6d60569a)
}
func (t ClientQueryAnswerError) MarshalTL() ([]byte, error) {
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
	b, err = tl.Marshal(t.Error)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(t.RequestID[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	b, err = tl.Marshal(t.RequestTokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ClientQueryAnswerError) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.ErrorCode)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Error)
	if err != nil {
		return err
	}
	_, err = io.ReadFull(r, t.RequestID[:])
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.RequestTokensUsed)
	if err != nil {
		return err
	}
	return nil
}
func (*ClientQueryAnswerError) _IClientQueryAnswer() {}

type IClientQueryAnswerEx interface {
	CRC() uint32
	MarshalTL() ([]byte, error)
	UnmarshalTL(io.Reader) error
	_IClientQueryAnswerEx()
}

var (
	_ IClientQueryAnswerEx = (*ClientQueryAnswerEx)(nil)
	_ IClientQueryAnswerEx = (*ClientQueryAnswerErrorEx)(nil)
	_ IClientQueryAnswerEx = (*ClientQueryAnswerPartEx)(nil)
)

func decodeIClientQueryAnswerEx(r io.Reader) (IClientQueryAnswerEx, error) {
	var tag uint32
	err := tl.Unmarshal(r, &tag)
	if err != nil {
		return nil, err
	}
	var res IClientQueryAnswerEx
	switch tag {
	case uint32(0xcd524720):
		res = &ClientQueryAnswerEx{}
	case uint32(0x72562a8):
		res = &ClientQueryAnswerErrorEx{}
	case uint32(0xc07bfaec):
		res = &ClientQueryAnswerPartEx{}
	default:
		return nil, fmt.Errorf("invalid crc code: got 0x%08x", tag)
	}
	err = res.UnmarshalTL(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func DecodeIClientQueryAnswerEx(r io.Reader) (IClientQueryAnswerEx, error) {
	return decodeIClientQueryAnswerEx(r)
}

type ClientQueryAnswerEx struct {
	RequestID [32]byte
	Answer    []byte
	Flags     uint32                `tl:"0,bitflag"`
	FinalInfo *ClientQueryFinalInfo `tl:",omitempty:Flags:0"`
}

func (*ClientQueryAnswerEx) CRC() uint32 {
	return uint32(0xcd524720)
}
func (t ClientQueryAnswerEx) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	flagsVar0 := t.Flags
	if t.FinalInfo != nil {
		flagsVar0 |= uint32(0x1)
	}
	_, err = buf.Write(t.RequestID[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	b, err = tl.Marshal(t.Answer)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(flagsVar0)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	if (flagsVar0>>int32(0))&1 == 1 {
		b, err = tl.Marshal(*t.FinalInfo)
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
func (t *ClientQueryAnswerEx) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	_, err = io.ReadFull(r, t.RequestID[:])
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Answer)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Flags)
	if err != nil {
		return err
	}
	if (t.Flags>>int32(0))&1 == 1 {
		var tmp8FinalInfo ClientQueryFinalInfo
		err = tl.Unmarshal(r, &tmp8FinalInfo)
		if err != nil {
			return err
		}
		t.FinalInfo = &tmp8FinalInfo
	}
	return nil
}
func (*ClientQueryAnswerEx) _IClientQueryAnswerEx() {}

type ClientQueryAnswerErrorEx struct {
	RequestID [32]byte
	ErrorCode int32
	Error     string
	Flags     uint32                `tl:"0,bitflag"`
	FinalInfo *ClientQueryFinalInfo `tl:",omitempty:Flags:0"`
}

func (*ClientQueryAnswerErrorEx) CRC() uint32 {
	return uint32(0x72562a8)
}
func (t ClientQueryAnswerErrorEx) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	flagsVar0 := t.Flags
	if t.FinalInfo != nil {
		flagsVar0 |= uint32(0x1)
	}
	_, err = buf.Write(t.RequestID[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	b, err = tl.Marshal(t.ErrorCode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Error)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(flagsVar0)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	if (flagsVar0>>int32(0))&1 == 1 {
		b, err = tl.Marshal(*t.FinalInfo)
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
func (t *ClientQueryAnswerErrorEx) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	_, err = io.ReadFull(r, t.RequestID[:])
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ErrorCode)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Error)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Flags)
	if err != nil {
		return err
	}
	if (t.Flags>>int32(0))&1 == 1 {
		var tmp10FinalInfo ClientQueryFinalInfo
		err = tl.Unmarshal(r, &tmp10FinalInfo)
		if err != nil {
			return err
		}
		t.FinalInfo = &tmp10FinalInfo
	}
	return nil
}
func (*ClientQueryAnswerErrorEx) _IClientQueryAnswerEx() {}

type ClientQueryAnswerPartEx struct {
	RequestID [32]byte
	Answer    []byte
	Flags     uint32                `tl:"0,bitflag"`
	FinalInfo *ClientQueryFinalInfo `tl:",omitempty:Flags:0"`
}

func (*ClientQueryAnswerPartEx) CRC() uint32 {
	return uint32(0xc07bfaec)
}
func (t ClientQueryAnswerPartEx) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	flagsVar0 := t.Flags
	if t.FinalInfo != nil {
		flagsVar0 |= uint32(0x1)
	}
	_, err = buf.Write(t.RequestID[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	b, err = tl.Marshal(t.Answer)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(flagsVar0)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	if (flagsVar0>>int32(0))&1 == 1 {
		b, err = tl.Marshal(*t.FinalInfo)
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
func (t *ClientQueryAnswerPartEx) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	_, err = io.ReadFull(r, t.RequestID[:])
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Answer)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Flags)
	if err != nil {
		return err
	}
	if (t.Flags>>int32(0))&1 == 1 {
		var tmp8FinalInfo ClientQueryFinalInfo
		err = tl.Unmarshal(r, &tmp8FinalInfo)
		if err != nil {
			return err
		}
		t.FinalInfo = &tmp8FinalInfo
	}
	return nil
}
func (*ClientQueryAnswerPartEx) _IClientQueryAnswerEx() {}

type IClientQueryAnswerPart interface {
	CRC() uint32
	MarshalTL() ([]byte, error)
	UnmarshalTL(io.Reader) error
	_IClientQueryAnswerPart()
}

var (
	_ IClientQueryAnswerPart = (*ClientQueryAnswerPart)(nil)
	_ IClientQueryAnswerPart = (*ClientQueryAnswerPartError)(nil)
)

func decodeIClientQueryAnswerPart(r io.Reader) (IClientQueryAnswerPart, error) {
	var tag uint32
	err := tl.Unmarshal(r, &tag)
	if err != nil {
		return nil, err
	}
	var res IClientQueryAnswerPart
	switch tag {
	case uint32(0xb765de4c):
		res = &ClientQueryAnswerPart{}
	case uint32(0xd790a022):
		res = &ClientQueryAnswerPartError{}
	default:
		return nil, fmt.Errorf("invalid crc code: got 0x%08x", tag)
	}
	err = res.UnmarshalTL(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func DecodeIClientQueryAnswerPart(r io.Reader) (IClientQueryAnswerPart, error) {
	return decodeIClientQueryAnswerPart(r)
}

type ClientQueryAnswerPart struct {
	Answer            []byte
	IsCompleted       bool
	RequestID         [32]byte
	RequestTokensUsed TokensUsed
}

func (*ClientQueryAnswerPart) CRC() uint32 {
	return uint32(0xb765de4c)
}
func (t ClientQueryAnswerPart) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.Answer)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.IsCompleted)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(t.RequestID[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	b, err = tl.Marshal(t.RequestTokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ClientQueryAnswerPart) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Answer)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.IsCompleted)
	if err != nil {
		return err
	}
	_, err = io.ReadFull(r, t.RequestID[:])
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.RequestTokensUsed)
	if err != nil {
		return err
	}
	return nil
}
func (*ClientQueryAnswerPart) _IClientQueryAnswerPart() {}

type ClientQueryAnswerPartError struct {
	ErrorCode         int32
	Error             string
	RequestID         [32]byte
	RequestTokensUsed TokensUsed
}

func (*ClientQueryAnswerPartError) CRC() uint32 {
	return uint32(0xd790a022)
}
func (t ClientQueryAnswerPartError) MarshalTL() ([]byte, error) {
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
	b, err = tl.Marshal(t.Error)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(t.RequestID[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	b, err = tl.Marshal(t.RequestTokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ClientQueryAnswerPartError) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.ErrorCode)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Error)
	if err != nil {
		return err
	}
	_, err = io.ReadFull(r, t.RequestID[:])
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.RequestTokensUsed)
	if err != nil {
		return err
	}
	return nil
}
func (*ClientQueryAnswerPartError) _IClientQueryAnswerPart() {}

type ClientQueryFinalInfo struct {
	Flags           uint32 `tl:"0,bitflag"`
	TokensUsed      TokensUsed
	WorkerDebug     *string  `tl:",omitempty:Flags:0"`
	ProxyDebug      *string  `tl:",omitempty:Flags:0"`
	ProxyStartTime  *float64 `tl:",omitempty:Flags:1"`
	ProxyEndTime    *float64 `tl:",omitempty:Flags:1"`
	WorkerStartTime *float64 `tl:",omitempty:Flags:1"`
	WorkerEndTime   *float64 `tl:",omitempty:Flags:1"`
}

func (*ClientQueryFinalInfo) CRC() uint32 {
	return uint32(0x69a452f0)
}
func (t ClientQueryFinalInfo) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	flagsVar0 := t.Flags
	if t.WorkerDebug != nil {
		flagsVar0 |= uint32(0x1)
	}
	if t.ProxyDebug != nil {
		flagsVar0 |= uint32(0x1)
	}
	if t.ProxyStartTime != nil {
		flagsVar0 |= uint32(0x2)
	}
	if t.ProxyEndTime != nil {
		flagsVar0 |= uint32(0x2)
	}
	if t.WorkerStartTime != nil {
		flagsVar0 |= uint32(0x2)
	}
	if t.WorkerEndTime != nil {
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
	b, err = tl.Marshal(t.TokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	if (flagsVar0>>int32(0))&1 == 1 {
		b, err = tl.Marshal(*t.WorkerDebug)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (flagsVar0>>int32(0))&1 == 1 {
		b, err = tl.Marshal(*t.ProxyDebug)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (flagsVar0>>int32(1))&1 == 1 {
		m21ProxyStartTimeBits := math.Float64bits(*t.ProxyStartTime)
		var m21ProxyStartTimeRaw [8]byte
		binary.LittleEndian.PutUint64(m21ProxyStartTimeRaw[:], m21ProxyStartTimeBits)
		_, err = buf.Write(m21ProxyStartTimeRaw[:])
		if err != nil {
			return nil, err
		}
	}
	if (flagsVar0>>int32(1))&1 == 1 {
		m22ProxyEndTimeBits := math.Float64bits(*t.ProxyEndTime)
		var m22ProxyEndTimeRaw [8]byte
		binary.LittleEndian.PutUint64(m22ProxyEndTimeRaw[:], m22ProxyEndTimeBits)
		_, err = buf.Write(m22ProxyEndTimeRaw[:])
		if err != nil {
			return nil, err
		}
	}
	if (flagsVar0>>int32(1))&1 == 1 {
		m23WorkerStartTimeBits := math.Float64bits(*t.WorkerStartTime)
		var m23WorkerStartTimeRaw [8]byte
		binary.LittleEndian.PutUint64(m23WorkerStartTimeRaw[:], m23WorkerStartTimeBits)
		_, err = buf.Write(m23WorkerStartTimeRaw[:])
		if err != nil {
			return nil, err
		}
	}
	if (flagsVar0>>int32(1))&1 == 1 {
		m24WorkerEndTimeBits := math.Float64bits(*t.WorkerEndTime)
		var m24WorkerEndTimeRaw [8]byte
		binary.LittleEndian.PutUint64(m24WorkerEndTimeRaw[:], m24WorkerEndTimeBits)
		_, err = buf.Write(m24WorkerEndTimeRaw[:])
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
func (t *ClientQueryFinalInfo) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Flags)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.TokensUsed)
	if err != nil {
		return err
	}
	if (t.Flags>>int32(0))&1 == 1 {
		var tmp6WorkerDebug string
		err = tl.Unmarshal(r, &tmp6WorkerDebug)
		if err != nil {
			return err
		}
		t.WorkerDebug = &tmp6WorkerDebug
	}
	if (t.Flags>>int32(0))&1 == 1 {
		var tmp7ProxyDebug string
		err = tl.Unmarshal(r, &tmp7ProxyDebug)
		if err != nil {
			return err
		}
		t.ProxyDebug = &tmp7ProxyDebug
	}
	if (t.Flags>>int32(1))&1 == 1 {
		var tmp8ProxyStartTimeBits uint64
		err = tl.Unmarshal(r, &tmp8ProxyStartTimeBits)
		if err != nil {
			return err
		}
		tmp8ProxyStartTimeVal := math.Float64frombits(tmp8ProxyStartTimeBits)
		t.ProxyStartTime = &tmp8ProxyStartTimeVal
	}
	if (t.Flags>>int32(1))&1 == 1 {
		var tmp9ProxyEndTimeBits uint64
		err = tl.Unmarshal(r, &tmp9ProxyEndTimeBits)
		if err != nil {
			return err
		}
		tmp9ProxyEndTimeVal := math.Float64frombits(tmp9ProxyEndTimeBits)
		t.ProxyEndTime = &tmp9ProxyEndTimeVal
	}
	if (t.Flags>>int32(1))&1 == 1 {
		var tmp10WorkerStartTimeBits uint64
		err = tl.Unmarshal(r, &tmp10WorkerStartTimeBits)
		if err != nil {
			return err
		}
		tmp10WorkerStartTimeVal := math.Float64frombits(tmp10WorkerStartTimeBits)
		t.WorkerStartTime = &tmp10WorkerStartTimeVal
	}
	if (t.Flags>>int32(1))&1 == 1 {
		var tmp11WorkerEndTimeBits uint64
		err = tl.Unmarshal(r, &tmp11WorkerEndTimeBits)
		if err != nil {
			return err
		}
		tmp11WorkerEndTimeVal := math.Float64frombits(tmp11WorkerEndTimeBits)
		t.WorkerEndTime = &tmp11WorkerEndTimeVal
	}
	return nil
}

type IClientRefund interface {
	CRC() uint32
	MarshalTL() ([]byte, error)
	UnmarshalTL(io.Reader) error
	_IClientRefund()
}

var (
	_ IClientRefund = (*ClientRefund)(nil)
	_ IClientRefund = (*ClientRefundRejected)(nil)
)

func decodeIClientRefund(r io.Reader) (IClientRefund, error) {
	var tag uint32
	err := tl.Unmarshal(r, &tag)
	if err != nil {
		return nil, err
	}
	var res IClientRefund
	switch tag {
	case uint32(0x83aabead):
		res = &ClientRefund{}
	case uint32(0xcf9d9957):
		res = &ClientRefundRejected{}
	default:
		return nil, fmt.Errorf("invalid crc code: got 0x%08x", tag)
	}
	err = res.UnmarshalTL(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func DecodeIClientRefund(r io.Reader) (IClientRefund, error) {
	return decodeIClientRefund(r)
}

type ClientRefund struct {
	Data []byte
}

func (*ClientRefund) CRC() uint32 {
	return uint32(0x83aabead)
}
func (t ClientRefund) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.Data)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ClientRefund) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Data)
	if err != nil {
		return err
	}
	return nil
}
func (*ClientRefund) _IClientRefund() {}

type ClientRefundRejected struct {
	ActiveQueries int64
}

func (*ClientRefundRejected) CRC() uint32 {
	return uint32(0xcf9d9957)
}
func (t ClientRefundRejected) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.ActiveQueries)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ClientRefundRejected) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.ActiveQueries)
	if err != nil {
		return err
	}
	return nil
}
func (*ClientRefundRejected) _IClientRefund() {}

type ClientWorkerInstanceV2 struct {
	Flags             uint32 `tl:"0,bitflag"`
	Coefficient       int32
	ActiveRequests    int32
	MaxActiveRequests int32
}

func (*ClientWorkerInstanceV2) CRC() uint32 {
	return uint32(0x3ea93d00)
}
func (t ClientWorkerInstanceV2) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	flagsVar0 := t.Flags
	b, err = tl.Marshal(flagsVar0)
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
	b, err = tl.Marshal(t.ActiveRequests)
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
	return buf.Bytes(), nil
}
func (t *ClientWorkerInstanceV2) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Flags)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Coefficient)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ActiveRequests)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.MaxActiveRequests)
	if err != nil {
		return err
	}
	return nil
}

type ClientWorkerType struct {
	Name                string
	ActiveWorkers       int32
	CoefficientMin      int32
	CoefficientBucket50 int32
	CoefficientMax      int32
}

func (*ClientWorkerType) CRC() uint32 {
	return uint32(0x8210e21c)
}
func (t ClientWorkerType) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.Name)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ActiveWorkers)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.CoefficientMin)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.CoefficientBucket50)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.CoefficientMax)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ClientWorkerType) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Name)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ActiveWorkers)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.CoefficientMin)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.CoefficientBucket50)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.CoefficientMax)
	if err != nil {
		return err
	}
	return nil
}

type ClientWorkerTypeV2 struct {
	Name    string
	Workers []ClientWorkerInstanceV2
}

func (*ClientWorkerTypeV2) CRC() uint32 {
	return uint32(0xb27d8197)
}
func (t ClientWorkerTypeV2) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.Name)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Workers)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ClientWorkerTypeV2) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Name)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Workers)
	if err != nil {
		return err
	}
	return nil
}

type ClientWorkerTypes struct {
	Types []ClientWorkerType
}

func (*ClientWorkerTypes) CRC() uint32 {
	return uint32(0xf64e0e01)
}
func (t ClientWorkerTypes) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.Types)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ClientWorkerTypes) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Types)
	if err != nil {
		return err
	}
	return nil
}

type ClientWorkerTypesV2 struct {
	Types []ClientWorkerTypeV2
}

func (*ClientWorkerTypesV2) CRC() uint32 {
	return uint32(0xcf0dc67)
}
func (t ClientWorkerTypesV2) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.Types)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ClientWorkerTypesV2) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Types)
	if err != nil {
		return err
	}
	return nil
}
