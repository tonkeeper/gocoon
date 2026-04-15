// Code generated - DO NOT EDIT.

package tlcocoonapi

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/tonkeeper/tongo/tl"
	"io"
)

type TokensUsedC struct {
	PromptTokensUsed     uint64
	CachedTokensUsed     uint64
	CompletionTokensUsed uint64
	ReasoningTokensUsed  uint64
	TotalTokensUsed      uint64
}

func (t TokensUsedC) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.PromptTokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.CachedTokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.CompletionTokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ReasoningTokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.TotalTokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *TokensUsedC) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.PromptTokensUsed)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.CachedTokensUsed)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.CompletionTokensUsed)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ReasoningTokensUsed)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.TotalTokensUsed)
	if err != nil {
		return err
	}
	return nil
}

type ProxyParamsC struct {
	Flags             uint32
	ProxyPublicKey    tl.Int256
	ProxyOwnerAddress string
	ProxyScAddress    string
	IsTest            *bool
	ProtoVersion      *uint32
}

func (t ProxyParamsC) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Flags)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ProxyPublicKey)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ProxyOwnerAddress)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ProxyScAddress)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	if (t.Flags>>0)&1 == 1 {
		b, err = tl.Marshal(*t.IsTest)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (t.Flags>>1)&1 == 1 {
		b, err = tl.Marshal(*t.ProtoVersion)
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

func (t *ProxyParamsC) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Flags)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ProxyPublicKey)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ProxyOwnerAddress)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ProxyScAddress)
	if err != nil {
		return err
	}
	if (t.Flags>>0)&1 == 1 {
		var tempIsTest bool
		err = tl.Unmarshal(r, &tempIsTest)
		if err != nil {
			return err
		}
		t.IsTest = &tempIsTest
	}
	if (t.Flags>>1)&1 == 1 {
		var tempProtoVersion uint32
		err = tl.Unmarshal(r, &tempProtoVersion)
		if err != nil {
			return err
		}
		t.ProtoVersion = &tempProtoVersion
	}
	return nil
}

type ClientParamsC struct {
	Flags              uint32
	ClientOwnerAddress string
	IsTest             *bool
	MinProtoVersion    *uint32
	MaxProtoVersion    *uint32
}

func (t ClientParamsC) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Flags)
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
	if (t.Flags>>0)&1 == 1 {
		b, err = tl.Marshal(*t.IsTest)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (t.Flags>>1)&1 == 1 {
		b, err = tl.Marshal(*t.MinProtoVersion)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (t.Flags>>1)&1 == 1 {
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

func (t *ClientParamsC) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Flags)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ClientOwnerAddress)
	if err != nil {
		return err
	}
	if (t.Flags>>0)&1 == 1 {
		var tempIsTest bool
		err = tl.Unmarshal(r, &tempIsTest)
		if err != nil {
			return err
		}
		t.IsTest = &tempIsTest
	}
	if (t.Flags>>1)&1 == 1 {
		var tempMinProtoVersion uint32
		err = tl.Unmarshal(r, &tempMinProtoVersion)
		if err != nil {
			return err
		}
		t.MinProtoVersion = &tempMinProtoVersion
	}
	if (t.Flags>>1)&1 == 1 {
		var tempMaxProtoVersion uint32
		err = tl.Unmarshal(r, &tempMaxProtoVersion)
		if err != nil {
			return err
		}
		t.MaxProtoVersion = &tempMaxProtoVersion
	}
	return nil
}

type ProxySignedPayment struct {
	tl.SumType
	ProxySignedPayment struct {
		Data []byte
	}
	ProxySignedPaymentEmpty struct{}
}

func (t ProxySignedPayment) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	switch t.SumType {
	case "ProxySignedPayment":
		b, err = tl.Marshal(uint32(0x2998182))
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		b, err = tl.Marshal(t.ProxySignedPayment.Data)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	case "ProxySignedPaymentEmpty":
		b, err = tl.Marshal(uint32(0xb347ce64))
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
	default:
		return nil, fmt.Errorf("invalid sum type")
	}
	return buf.Bytes(), nil
}

func (t *ProxySignedPayment) UnmarshalTL(r io.Reader) error {
	var err error
	var b [4]byte
	_, err = io.ReadFull(r, b[:])
	if err != nil {
		return err
	}
	tag := int(binary.LittleEndian.Uint32(b[:]))
	switch tag {
	case 0x2998182:
		t.SumType = "ProxySignedPayment"
		err = tl.Unmarshal(r, &t.ProxySignedPayment.Data)
		if err != nil {
			return err
		}
	case 0xb347ce64:
		t.SumType = "ProxySignedPaymentEmpty"
	default:
		return fmt.Errorf("invalid tag")
	}
	return nil
}

type ClientProxyConnectionAuth struct {
	tl.SumType
	ClientProxyConnectionAuthShort struct {
		SecretHash tl.Int256
		Nonce      uint64
	}
	ClientProxyConnectionAuthLong struct {
		Nonce uint64
	}
}

func (t ClientProxyConnectionAuth) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	switch t.SumType {
	case "ClientProxyConnectionAuthShort":
		b, err = tl.Marshal(uint32(0xd6ffc5af))
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		b, err = tl.Marshal(t.ClientProxyConnectionAuthShort.SecretHash)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.ClientProxyConnectionAuthShort.Nonce)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	case "ClientProxyConnectionAuthLong":
		b, err = tl.Marshal(uint32(0x417bf016))
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		b, err = tl.Marshal(t.ClientProxyConnectionAuthLong.Nonce)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid sum type")
	}
	return buf.Bytes(), nil
}

func (t *ClientProxyConnectionAuth) UnmarshalTL(r io.Reader) error {
	var err error
	var b [4]byte
	_, err = io.ReadFull(r, b[:])
	if err != nil {
		return err
	}
	tag := int(binary.LittleEndian.Uint32(b[:]))
	switch tag {
	case 0xd6ffc5af:
		t.SumType = "ClientProxyConnectionAuthShort"
		err = tl.Unmarshal(r, &t.ClientProxyConnectionAuthShort.SecretHash)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.ClientProxyConnectionAuthShort.Nonce)
		if err != nil {
			return err
		}
	case 0x417bf016:
		t.SumType = "ClientProxyConnectionAuthLong"
		err = tl.Unmarshal(r, &t.ClientProxyConnectionAuthLong.Nonce)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid tag")
	}
	return nil
}

type ClientAuthorizationWithProxy struct {
	tl.SumType
	ClientAuthorizationWithProxySuccess struct {
		SignedPayment       ProxySignedPayment
		TokensCommittedToDb uint64
		MaxTokens           uint64
	}
	ClientAuthorizationWithProxyFailed struct {
		ErrorCode uint32
		Error     string
	}
}

func (t ClientAuthorizationWithProxy) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	switch t.SumType {
	case "ClientAuthorizationWithProxySuccess":
		b, err = tl.Marshal(uint32(0x75d5ac34))
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		b, err = tl.Marshal(t.ClientAuthorizationWithProxySuccess.SignedPayment)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.ClientAuthorizationWithProxySuccess.TokensCommittedToDb)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.ClientAuthorizationWithProxySuccess.MaxTokens)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	case "ClientAuthorizationWithProxyFailed":
		b, err = tl.Marshal(uint32(0x60551c96))
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		b, err = tl.Marshal(t.ClientAuthorizationWithProxyFailed.ErrorCode)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.ClientAuthorizationWithProxyFailed.Error)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid sum type")
	}
	return buf.Bytes(), nil
}

func (t *ClientAuthorizationWithProxy) UnmarshalTL(r io.Reader) error {
	var err error
	var b [4]byte
	_, err = io.ReadFull(r, b[:])
	if err != nil {
		return err
	}
	tag := int(binary.LittleEndian.Uint32(b[:]))
	switch tag {
	case 0x75d5ac34:
		t.SumType = "ClientAuthorizationWithProxySuccess"
		err = tl.Unmarshal(r, &t.ClientAuthorizationWithProxySuccess.SignedPayment)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.ClientAuthorizationWithProxySuccess.TokensCommittedToDb)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.ClientAuthorizationWithProxySuccess.MaxTokens)
		if err != nil {
			return err
		}
	case 0x60551c96:
		t.SumType = "ClientAuthorizationWithProxyFailed"
		err = tl.Unmarshal(r, &t.ClientAuthorizationWithProxyFailed.ErrorCode)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.ClientAuthorizationWithProxyFailed.Error)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid tag")
	}
	return nil
}

type ClientConnectedToProxyC struct {
	Params          ProxyParamsC
	ClientScAddress string
	Auth            ClientProxyConnectionAuth
	SignedPayment   ProxySignedPayment
}

func (t ClientConnectedToProxyC) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
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
	b, err = tl.Marshal(t.Auth)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.SignedPayment)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *ClientConnectedToProxyC) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Params)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ClientScAddress)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Auth)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.SignedPayment)
	if err != nil {
		return err
	}
	return nil
}

type ClientQueryAnswer struct {
	tl.SumType
	ClientQueryAnswer struct {
		Answer            []byte
		IsCompleted       bool
		RequestId         tl.Int256
		RequestTokensUsed TokensUsedC
	}
	ClientQueryAnswerError struct {
		ErrorCode         uint32
		Error             string
		RequestId         tl.Int256
		RequestTokensUsed TokensUsedC
	}
}

func (t ClientQueryAnswer) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	switch t.SumType {
	case "ClientQueryAnswer":
		b, err = tl.Marshal(uint32(0x9b943922))
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		b, err = tl.Marshal(t.ClientQueryAnswer.Answer)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.ClientQueryAnswer.IsCompleted)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.ClientQueryAnswer.RequestId)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.ClientQueryAnswer.RequestTokensUsed)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	case "ClientQueryAnswerError":
		b, err = tl.Marshal(uint32(0x6d60569a))
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		b, err = tl.Marshal(t.ClientQueryAnswerError.ErrorCode)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.ClientQueryAnswerError.Error)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.ClientQueryAnswerError.RequestId)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.ClientQueryAnswerError.RequestTokensUsed)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid sum type")
	}
	return buf.Bytes(), nil
}

func (t *ClientQueryAnswer) UnmarshalTL(r io.Reader) error {
	var err error
	var b [4]byte
	_, err = io.ReadFull(r, b[:])
	if err != nil {
		return err
	}
	tag := int(binary.LittleEndian.Uint32(b[:]))
	switch tag {
	case 0x9b943922:
		t.SumType = "ClientQueryAnswer"
		err = tl.Unmarshal(r, &t.ClientQueryAnswer.Answer)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.ClientQueryAnswer.IsCompleted)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.ClientQueryAnswer.RequestId)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.ClientQueryAnswer.RequestTokensUsed)
		if err != nil {
			return err
		}
	case 0x6d60569a:
		t.SumType = "ClientQueryAnswerError"
		err = tl.Unmarshal(r, &t.ClientQueryAnswerError.ErrorCode)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.ClientQueryAnswerError.Error)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.ClientQueryAnswerError.RequestId)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.ClientQueryAnswerError.RequestTokensUsed)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid tag")
	}
	return nil
}

type ClientQueryFinalInfoC struct {
	Flags           uint32
	TokensUsed      TokensUsedC
	WorkerDebug     *string
	ProxyDebug      *string
	ProxyStartTime  *float64
	ProxyEndTime    *float64
	WorkerStartTime *float64
	WorkerEndTime   *float64
}

func (t ClientQueryFinalInfoC) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Flags)
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
	if (t.Flags>>0)&1 == 1 {
		b, err = tl.Marshal(*t.WorkerDebug)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (t.Flags>>0)&1 == 1 {
		b, err = tl.Marshal(*t.ProxyDebug)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (t.Flags>>1)&1 == 1 {
		b, err = tl.Marshal(*t.ProxyStartTime)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (t.Flags>>1)&1 == 1 {
		b, err = tl.Marshal(*t.ProxyEndTime)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (t.Flags>>1)&1 == 1 {
		b, err = tl.Marshal(*t.WorkerStartTime)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (t.Flags>>1)&1 == 1 {
		b, err = tl.Marshal(*t.WorkerEndTime)
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

func (t *ClientQueryFinalInfoC) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Flags)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.TokensUsed)
	if err != nil {
		return err
	}
	if (t.Flags>>0)&1 == 1 {
		var tempWorkerDebug string
		err = tl.Unmarshal(r, &tempWorkerDebug)
		if err != nil {
			return err
		}
		t.WorkerDebug = &tempWorkerDebug
	}
	if (t.Flags>>0)&1 == 1 {
		var tempProxyDebug string
		err = tl.Unmarshal(r, &tempProxyDebug)
		if err != nil {
			return err
		}
		t.ProxyDebug = &tempProxyDebug
	}
	if (t.Flags>>1)&1 == 1 {
		var tempProxyStartTime float64
		err = tl.Unmarshal(r, &tempProxyStartTime)
		if err != nil {
			return err
		}
		t.ProxyStartTime = &tempProxyStartTime
	}
	if (t.Flags>>1)&1 == 1 {
		var tempProxyEndTime float64
		err = tl.Unmarshal(r, &tempProxyEndTime)
		if err != nil {
			return err
		}
		t.ProxyEndTime = &tempProxyEndTime
	}
	if (t.Flags>>1)&1 == 1 {
		var tempWorkerStartTime float64
		err = tl.Unmarshal(r, &tempWorkerStartTime)
		if err != nil {
			return err
		}
		t.WorkerStartTime = &tempWorkerStartTime
	}
	if (t.Flags>>1)&1 == 1 {
		var tempWorkerEndTime float64
		err = tl.Unmarshal(r, &tempWorkerEndTime)
		if err != nil {
			return err
		}
		t.WorkerEndTime = &tempWorkerEndTime
	}
	return nil
}

type ClientQueryAnswerEx struct {
	tl.SumType
	ClientQueryAnswerEx struct {
		RequestId tl.Int256
		Answer    []byte
		Flags     uint32
		FinalInfo *ClientQueryFinalInfoC
	}
	ClientQueryAnswerErrorEx struct {
		RequestId tl.Int256
		ErrorCode uint32
		Error     string
		Flags     uint32
		FinalInfo *ClientQueryFinalInfoC
	}
	ClientQueryAnswerPartEx struct {
		RequestId tl.Int256
		Answer    []byte
		Flags     uint32
		FinalInfo *ClientQueryFinalInfoC
	}
}

func (t ClientQueryAnswerEx) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	switch t.SumType {
	case "ClientQueryAnswerEx":
		b, err = tl.Marshal(uint32(0xcd524720))
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		b, err = tl.Marshal(t.ClientQueryAnswerEx.RequestId)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.ClientQueryAnswerEx.Answer)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.ClientQueryAnswerEx.Flags)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		if (t.ClientQueryAnswerEx.Flags>>0)&1 == 1 {
			b, err = tl.Marshal(*t.ClientQueryAnswerEx.FinalInfo)
			if err != nil {
				return nil, err
			}
			_, err = buf.Write(b)
			if err != nil {
				return nil, err
			}
		}
	case "ClientQueryAnswerErrorEx":
		b, err = tl.Marshal(uint32(0x72562a8))
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		b, err = tl.Marshal(t.ClientQueryAnswerErrorEx.RequestId)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.ClientQueryAnswerErrorEx.ErrorCode)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.ClientQueryAnswerErrorEx.Error)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.ClientQueryAnswerErrorEx.Flags)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		if (t.ClientQueryAnswerErrorEx.Flags>>0)&1 == 1 {
			b, err = tl.Marshal(*t.ClientQueryAnswerErrorEx.FinalInfo)
			if err != nil {
				return nil, err
			}
			_, err = buf.Write(b)
			if err != nil {
				return nil, err
			}
		}
	case "ClientQueryAnswerPartEx":
		b, err = tl.Marshal(uint32(0xc07bfaec))
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		b, err = tl.Marshal(t.ClientQueryAnswerPartEx.RequestId)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.ClientQueryAnswerPartEx.Answer)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.ClientQueryAnswerPartEx.Flags)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		if (t.ClientQueryAnswerPartEx.Flags>>0)&1 == 1 {
			b, err = tl.Marshal(*t.ClientQueryAnswerPartEx.FinalInfo)
			if err != nil {
				return nil, err
			}
			_, err = buf.Write(b)
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("invalid sum type")
	}
	return buf.Bytes(), nil
}

func (t *ClientQueryAnswerEx) UnmarshalTL(r io.Reader) error {
	var err error
	var b [4]byte
	_, err = io.ReadFull(r, b[:])
	if err != nil {
		return err
	}
	tag := int(binary.LittleEndian.Uint32(b[:]))
	switch tag {
	case 0xcd524720:
		t.SumType = "ClientQueryAnswerEx"
		err = tl.Unmarshal(r, &t.ClientQueryAnswerEx.RequestId)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.ClientQueryAnswerEx.Answer)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.ClientQueryAnswerEx.Flags)
		if err != nil {
			return err
		}
		if (t.ClientQueryAnswerEx.Flags>>0)&1 == 1 {
			var tempFinalInfo ClientQueryFinalInfoC
			err = tl.Unmarshal(r, &tempFinalInfo)
			if err != nil {
				return err
			}
			t.ClientQueryAnswerEx.FinalInfo = &tempFinalInfo
		}
	case 0x72562a8:
		t.SumType = "ClientQueryAnswerErrorEx"
		err = tl.Unmarshal(r, &t.ClientQueryAnswerErrorEx.RequestId)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.ClientQueryAnswerErrorEx.ErrorCode)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.ClientQueryAnswerErrorEx.Error)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.ClientQueryAnswerErrorEx.Flags)
		if err != nil {
			return err
		}
		if (t.ClientQueryAnswerErrorEx.Flags>>0)&1 == 1 {
			var tempFinalInfo ClientQueryFinalInfoC
			err = tl.Unmarshal(r, &tempFinalInfo)
			if err != nil {
				return err
			}
			t.ClientQueryAnswerErrorEx.FinalInfo = &tempFinalInfo
		}
	case 0xc07bfaec:
		t.SumType = "ClientQueryAnswerPartEx"
		err = tl.Unmarshal(r, &t.ClientQueryAnswerPartEx.RequestId)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.ClientQueryAnswerPartEx.Answer)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.ClientQueryAnswerPartEx.Flags)
		if err != nil {
			return err
		}
		if (t.ClientQueryAnswerPartEx.Flags>>0)&1 == 1 {
			var tempFinalInfo ClientQueryFinalInfoC
			err = tl.Unmarshal(r, &tempFinalInfo)
			if err != nil {
				return err
			}
			t.ClientQueryAnswerPartEx.FinalInfo = &tempFinalInfo
		}
	default:
		return fmt.Errorf("invalid tag")
	}
	return nil
}

type ClientRefund struct {
	tl.SumType
	ClientRefund struct {
		Data []byte
	}
	ClientRefundRejected struct {
		ActiveQueries uint64
	}
}

func (t ClientRefund) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	switch t.SumType {
	case "ClientRefund":
		b, err = tl.Marshal(uint32(0x83aabead))
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		b, err = tl.Marshal(t.ClientRefund.Data)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	case "ClientRefundRejected":
		b, err = tl.Marshal(uint32(0xcf9d9957))
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		b, err = tl.Marshal(t.ClientRefundRejected.ActiveQueries)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid sum type")
	}
	return buf.Bytes(), nil
}

func (t *ClientRefund) UnmarshalTL(r io.Reader) error {
	var err error
	var b [4]byte
	_, err = io.ReadFull(r, b[:])
	if err != nil {
		return err
	}
	tag := int(binary.LittleEndian.Uint32(b[:]))
	switch tag {
	case 0x83aabead:
		t.SumType = "ClientRefund"
		err = tl.Unmarshal(r, &t.ClientRefund.Data)
		if err != nil {
			return err
		}
	case 0xcf9d9957:
		t.SumType = "ClientRefundRejected"
		err = tl.Unmarshal(r, &t.ClientRefundRejected.ActiveQueries)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid tag")
	}
	return nil
}

type ClientWorkerTypeC struct {
	Name                string
	ActiveWorkers       uint32
	CoefficientMin      uint32
	CoefficientBucket50 uint32
	CoefficientMax      uint32
}

func (t ClientWorkerTypeC) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
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

func (t *ClientWorkerTypeC) UnmarshalTL(r io.Reader) error {
	var err error
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

type ClientWorkerTypesC struct {
	Types []ClientWorkerTypeC
}

func (t ClientWorkerTypesC) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
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

func (t *ClientWorkerTypesC) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Types)
	if err != nil {
		return err
	}
	return nil
}

type ClientWorkerInstanceV2C struct {
	Flags             uint32
	Coefficient       uint32
	ActiveRequests    uint32
	MaxActiveRequests uint32
}

func (t ClientWorkerInstanceV2C) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Flags)
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

func (t *ClientWorkerInstanceV2C) UnmarshalTL(r io.Reader) error {
	var err error
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

type ClientWorkerTypeV2C struct {
	Name    string
	Workers []ClientWorkerInstanceV2C
}

func (t ClientWorkerTypeV2C) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
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

func (t *ClientWorkerTypeV2C) UnmarshalTL(r io.Reader) error {
	var err error
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

type ClientWorkerTypesV2C struct {
	Types []ClientWorkerTypeV2C
}

func (t ClientWorkerTypesV2C) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
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

func (t *ClientWorkerTypesV2C) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Types)
	if err != nil {
		return err
	}
	return nil
}

type ClientPaymentStatusC struct {
	SignedPayment ProxySignedPayment
	DbTokens      uint64
	MaxTokens     uint64
}

func (t ClientPaymentStatusC) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.SignedPayment)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
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

func (t *ClientPaymentStatusC) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.SignedPayment)
	if err != nil {
		return err
	}
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

type ClientConnectToProxyRequest struct {
	Params           ClientParamsC
	MinConfigVersion uint32
}

func (t ClientConnectToProxyRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Params)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.MinConfigVersion)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Client) ClientConnectToProxy(ctx context.Context, request ClientConnectToProxyRequest) (res ClientConnectedToProxyC, err error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req ClientConnectToProxyRequest `tlSumType:"ff5fa0f4"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return res, err
	}
	resp, err := c.request(ctx, payload)
	if err != nil {
		return res, err
	}
	if len(resp) < 4 {
		return res, fmt.Errorf("not enough bytes for tag")
	}
	if binary.LittleEndian.Uint32(resp[:4]) != 0x95317ad1 {
		return res, fmt.Errorf("unexpected response tag")
	}
	err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
	return res, err
}

type ClientAuthorizeWithProxyShortRequest struct {
	Data []byte
}

func (t ClientAuthorizeWithProxyShortRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
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

func (c *Client) ClientAuthorizeWithProxyShort(ctx context.Context, request ClientAuthorizeWithProxyShortRequest) (res ClientAuthorizationWithProxy, err error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req ClientAuthorizeWithProxyShortRequest `tlSumType:"6c276723"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return res, err
	}
	resp, err := c.request(ctx, payload)
	if err != nil {
		return res, err
	}
	err = tl.Unmarshal(bytes.NewReader(resp), &res)
	return res, err
}

type ClientAuthorizeWithProxyLongRequest struct{}

func (c *Client) ClientAuthorizeWithProxyLong(ctx context.Context) (res ClientAuthorizationWithProxy, err error) {
	payload := make([]byte, 4)
	binary.LittleEndian.PutUint32(payload, 0xd3474303)
	resp, err := c.request(ctx, payload)
	if err != nil {
		return res, err
	}
	err = tl.Unmarshal(bytes.NewReader(resp), &res)
	return res, err
}

type ClientRunQueryRequest struct {
	ModelName        string
	Query            []byte
	MaxCoefficient   uint32
	MaxTokens        uint32
	Timeout          float64
	RequestId        tl.Int256
	MinConfigVersion uint32
}

func (t ClientRunQueryRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.ModelName)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Query)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.MaxCoefficient)
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
	b, err = tl.Marshal(t.Timeout)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.RequestId)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.MinConfigVersion)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Client) ClientRunQuery(ctx context.Context, request ClientRunQueryRequest) (res ClientQueryAnswer, err error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req ClientRunQueryRequest `tlSumType:"bc748f32"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return res, err
	}
	resp, err := c.request(ctx, payload)
	if err != nil {
		return res, err
	}
	err = tl.Unmarshal(bytes.NewReader(resp), &res)
	return res, err
}

type ClientRunQueryExRequest struct {
	ModelName        string
	Query            []byte
	MaxCoefficient   uint32
	MaxTokens        uint32
	Timeout          float64
	RequestId        tl.Int256
	MinConfigVersion uint32
	Flags            uint32
	EnableDebug      *bool
}

func (t ClientRunQueryExRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.ModelName)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Query)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.MaxCoefficient)
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
	b, err = tl.Marshal(t.Timeout)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.RequestId)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.MinConfigVersion)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Flags)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	if (t.Flags>>0)&1 == 1 {
		b, err = tl.Marshal(*t.EnableDebug)
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

func (c *Client) ClientRunQueryEx(ctx context.Context, request ClientRunQueryExRequest) (res ClientQueryAnswerEx, err error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req ClientRunQueryExRequest `tlSumType:"f54cb74b"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return res, err
	}
	resp, err := c.request(ctx, payload)
	if err != nil {
		return res, err
	}
	err = tl.Unmarshal(bytes.NewReader(resp), &res)
	return res, err
}

type ClientRequestRefundRequest struct{}

func (c *Client) ClientRequestRefund(ctx context.Context) (res ClientRefund, err error) {
	payload := make([]byte, 4)
	binary.LittleEndian.PutUint32(payload, 0x238d863d)
	resp, err := c.request(ctx, payload)
	if err != nil {
		return res, err
	}
	err = tl.Unmarshal(bytes.NewReader(resp), &res)
	return res, err
}

type ClientUpdatePaymentStatusRequest struct{}

func (c *Client) ClientUpdatePaymentStatus(ctx context.Context) (res ClientPaymentStatusC, err error) {
	payload := make([]byte, 4)
	binary.LittleEndian.PutUint32(payload, 0x9ed1c697)
	resp, err := c.request(ctx, payload)
	if err != nil {
		return res, err
	}
	if len(resp) < 4 {
		return res, fmt.Errorf("not enough bytes for tag")
	}
	if binary.LittleEndian.Uint32(resp[:4]) != 0xaa8a0ecc {
		return res, fmt.Errorf("unexpected response tag")
	}
	err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
	return res, err
}

type ClientGetWorkerTypesRequest struct{}

func (c *Client) ClientGetWorkerTypes(ctx context.Context) (res ClientWorkerTypesC, err error) {
	payload := make([]byte, 4)
	binary.LittleEndian.PutUint32(payload, 0x7f062bdb)
	resp, err := c.request(ctx, payload)
	if err != nil {
		return res, err
	}
	if len(resp) < 4 {
		return res, fmt.Errorf("not enough bytes for tag")
	}
	if binary.LittleEndian.Uint32(resp[:4]) != 0x65189591 {
		return res, fmt.Errorf("unexpected response tag")
	}
	err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
	return res, err
}

type ClientGetWorkerTypesV2Request struct{}

func (c *Client) ClientGetWorkerTypesV2(ctx context.Context) (res ClientWorkerTypesV2C, err error) {
	payload := make([]byte, 4)
	binary.LittleEndian.PutUint32(payload, 0xb2133d72)
	resp, err := c.request(ctx, payload)
	if err != nil {
		return res, err
	}
	if len(resp) < 4 {
		return res, fmt.Errorf("not enough bytes for tag")
	}
	if binary.LittleEndian.Uint32(resp[:4]) != 0x23507a68 {
		return res, fmt.Errorf("unexpected response tag")
	}
	err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
	return res, err
}
