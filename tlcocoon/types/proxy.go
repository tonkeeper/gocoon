package tlcocoonTypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	tl "github.com/tonkeeper/tongo/tl"
	"io"
	"math"
)

type ProxyParams struct {
	Flags             uint32 `tl:"0,bitflag"`
	ProxyPublicKey    [32]byte
	ProxyOwnerAddress string
	ProxyScAddress    string
	IsTest            *bool  `tl:",omitempty:Flags:0"`
	ProtoVersion      *int32 `tl:",omitempty:Flags:1"`
}

func (*ProxyParams) CRC() uint32 {
	return uint32(0xd5c5609f)
}
func (t ProxyParams) MarshalTL() ([]byte, error) {
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
	if t.ProtoVersion != nil {
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
	_, err = buf.Write(t.ProxyPublicKey[:])
	if err != nil {
		return nil, err
	}
	_ = 32
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
func (t *ProxyParams) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Flags)
	if err != nil {
		return err
	}
	_, err = io.ReadFull(r, t.ProxyPublicKey[:])
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
	if (t.Flags>>int32(0))&1 == 1 {
		var tmp10IsTest bool
		err = tl.Unmarshal(r, &tmp10IsTest)
		if err != nil {
			return err
		}
		t.IsTest = &tmp10IsTest
	}
	if (t.Flags>>int32(1))&1 == 1 {
		var tmp11ProtoVersion int32
		err = tl.Unmarshal(r, &tmp11ProtoVersion)
		if err != nil {
			return err
		}
		t.ProtoVersion = &tmp11ProtoVersion
	}
	return nil
}

type IProxyQueryAnswer interface {
	CRC() uint32
	MarshalTL() ([]byte, error)
	UnmarshalTL(io.Reader) error
	_IProxyQueryAnswer()
}

var (
	_ IProxyQueryAnswer = (*ProxyQueryAnswer)(nil)
	_ IProxyQueryAnswer = (*ProxyQueryAnswerError)(nil)
)

func decodeIProxyQueryAnswer(r io.Reader) (IProxyQueryAnswer, error) {
	var tag uint32
	err := tl.Unmarshal(r, &tag)
	if err != nil {
		return nil, err
	}
	var res IProxyQueryAnswer
	switch tag {
	case uint32(0x37e5f725):
		res = &ProxyQueryAnswer{}
	case uint32(0x669d014e):
		res = &ProxyQueryAnswerError{}
	default:
		return nil, fmt.Errorf("invalid crc code: got 0x%08x", tag)
	}
	err = res.UnmarshalTL(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func DecodeIProxyQueryAnswer(r io.Reader) (IProxyQueryAnswer, error) {
	return decodeIProxyQueryAnswer(r)
}

type ProxyQueryAnswer struct {
	Answer      []byte
	IsCompleted bool
	RequestID   [32]byte
	TokensUsed  TokensUsed
}

func (*ProxyQueryAnswer) CRC() uint32 {
	return uint32(0x37e5f725)
}
func (t ProxyQueryAnswer) MarshalTL() ([]byte, error) {
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
	b, err = tl.Marshal(t.TokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ProxyQueryAnswer) UnmarshalTL(r io.Reader) error {
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
	err = tl.Unmarshal(r, &t.TokensUsed)
	if err != nil {
		return err
	}
	return nil
}
func (*ProxyQueryAnswer) _IProxyQueryAnswer() {}

type ProxyQueryAnswerError struct {
	ErrorCode  int32
	Error      string
	RequestID  [32]byte
	TokensUsed TokensUsed
}

func (*ProxyQueryAnswerError) CRC() uint32 {
	return uint32(0x669d014e)
}
func (t ProxyQueryAnswerError) MarshalTL() ([]byte, error) {
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
	b, err = tl.Marshal(t.TokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ProxyQueryAnswerError) UnmarshalTL(r io.Reader) error {
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
	err = tl.Unmarshal(r, &t.TokensUsed)
	if err != nil {
		return err
	}
	return nil
}
func (*ProxyQueryAnswerError) _IProxyQueryAnswer() {}

type IProxyQueryAnswerEx interface {
	CRC() uint32
	MarshalTL() ([]byte, error)
	UnmarshalTL(io.Reader) error
	_IProxyQueryAnswerEx()
}

var (
	_ IProxyQueryAnswerEx = (*ProxyQueryAnswerEx)(nil)
	_ IProxyQueryAnswerEx = (*ProxyQueryAnswerErrorEx)(nil)
	_ IProxyQueryAnswerEx = (*ProxyQueryAnswerPartEx)(nil)
)

func decodeIProxyQueryAnswerEx(r io.Reader) (IProxyQueryAnswerEx, error) {
	var tag uint32
	err := tl.Unmarshal(r, &tag)
	if err != nil {
		return nil, err
	}
	var res IProxyQueryAnswerEx
	switch tag {
	case uint32(0x1ce574e8):
		res = &ProxyQueryAnswerEx{}
	case uint32(0x8f7fb0a7):
		res = &ProxyQueryAnswerErrorEx{}
	case uint32(0x6d85b50f):
		res = &ProxyQueryAnswerPartEx{}
	default:
		return nil, fmt.Errorf("invalid crc code: got 0x%08x", tag)
	}
	err = res.UnmarshalTL(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func DecodeIProxyQueryAnswerEx(r io.Reader) (IProxyQueryAnswerEx, error) {
	return decodeIProxyQueryAnswerEx(r)
}

type ProxyQueryAnswerEx struct {
	RequestID [32]byte
	Answer    []byte
	Flags     uint32               `tl:"0,bitflag"`
	FinalInfo *ProxyQueryFinalInfo `tl:",omitempty:Flags:0"`
}

func (*ProxyQueryAnswerEx) CRC() uint32 {
	return uint32(0x1ce574e8)
}
func (t ProxyQueryAnswerEx) MarshalTL() ([]byte, error) {
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
func (t *ProxyQueryAnswerEx) UnmarshalTL(r io.Reader) error {
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
		var tmp8FinalInfo ProxyQueryFinalInfo
		err = tl.Unmarshal(r, &tmp8FinalInfo)
		if err != nil {
			return err
		}
		t.FinalInfo = &tmp8FinalInfo
	}
	return nil
}
func (*ProxyQueryAnswerEx) _IProxyQueryAnswerEx() {}

type ProxyQueryAnswerErrorEx struct {
	RequestID [32]byte
	ErrorCode int32
	Error     string
	Flags     uint32               `tl:"0,bitflag"`
	FinalInfo *ProxyQueryFinalInfo `tl:",omitempty:Flags:0"`
}

func (*ProxyQueryAnswerErrorEx) CRC() uint32 {
	return uint32(0x8f7fb0a7)
}
func (t ProxyQueryAnswerErrorEx) MarshalTL() ([]byte, error) {
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
func (t *ProxyQueryAnswerErrorEx) UnmarshalTL(r io.Reader) error {
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
		var tmp10FinalInfo ProxyQueryFinalInfo
		err = tl.Unmarshal(r, &tmp10FinalInfo)
		if err != nil {
			return err
		}
		t.FinalInfo = &tmp10FinalInfo
	}
	return nil
}
func (*ProxyQueryAnswerErrorEx) _IProxyQueryAnswerEx() {}

type ProxyQueryAnswerPartEx struct {
	RequestID [32]byte
	Answer    []byte
	Flags     uint32               `tl:"0,bitflag"`
	FinalInfo *ProxyQueryFinalInfo `tl:",omitempty:Flags:0"`
}

func (*ProxyQueryAnswerPartEx) CRC() uint32 {
	return uint32(0x6d85b50f)
}
func (t ProxyQueryAnswerPartEx) MarshalTL() ([]byte, error) {
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
func (t *ProxyQueryAnswerPartEx) UnmarshalTL(r io.Reader) error {
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
		var tmp8FinalInfo ProxyQueryFinalInfo
		err = tl.Unmarshal(r, &tmp8FinalInfo)
		if err != nil {
			return err
		}
		t.FinalInfo = &tmp8FinalInfo
	}
	return nil
}
func (*ProxyQueryAnswerPartEx) _IProxyQueryAnswerEx() {}

type IProxyQueryAnswerPart interface {
	CRC() uint32
	MarshalTL() ([]byte, error)
	UnmarshalTL(io.Reader) error
	_IProxyQueryAnswerPart()
}

var (
	_ IProxyQueryAnswerPart = (*ProxyQueryAnswerPart)(nil)
	_ IProxyQueryAnswerPart = (*ProxyQueryAnswerPartError)(nil)
)

func decodeIProxyQueryAnswerPart(r io.Reader) (IProxyQueryAnswerPart, error) {
	var tag uint32
	err := tl.Unmarshal(r, &tag)
	if err != nil {
		return nil, err
	}
	var res IProxyQueryAnswerPart
	switch tag {
	case uint32(0x72eb00ee):
		res = &ProxyQueryAnswerPart{}
	case uint32(0x940e1168):
		res = &ProxyQueryAnswerPartError{}
	default:
		return nil, fmt.Errorf("invalid crc code: got 0x%08x", tag)
	}
	err = res.UnmarshalTL(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func DecodeIProxyQueryAnswerPart(r io.Reader) (IProxyQueryAnswerPart, error) {
	return decodeIProxyQueryAnswerPart(r)
}

type ProxyQueryAnswerPart struct {
	Answer      []byte
	IsCompleted bool
	RequestID   [32]byte
	TokensUsed  TokensUsed
}

func (*ProxyQueryAnswerPart) CRC() uint32 {
	return uint32(0x72eb00ee)
}
func (t ProxyQueryAnswerPart) MarshalTL() ([]byte, error) {
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
	b, err = tl.Marshal(t.TokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ProxyQueryAnswerPart) UnmarshalTL(r io.Reader) error {
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
	err = tl.Unmarshal(r, &t.TokensUsed)
	if err != nil {
		return err
	}
	return nil
}
func (*ProxyQueryAnswerPart) _IProxyQueryAnswerPart() {}

type ProxyQueryAnswerPartError struct {
	ErrorCode  int32
	Error      string
	RequestID  [32]byte
	TokensUsed TokensUsed
}

func (*ProxyQueryAnswerPartError) CRC() uint32 {
	return uint32(0x940e1168)
}
func (t ProxyQueryAnswerPartError) MarshalTL() ([]byte, error) {
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
	b, err = tl.Marshal(t.TokensUsed)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ProxyQueryAnswerPartError) UnmarshalTL(r io.Reader) error {
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
	err = tl.Unmarshal(r, &t.TokensUsed)
	if err != nil {
		return err
	}
	return nil
}
func (*ProxyQueryAnswerPartError) _IProxyQueryAnswerPart() {}

type ProxyQueryFinalInfo struct {
	Flags           uint32 `tl:"0,bitflag"`
	TokensUsed      TokensUsed
	WorkerDebug     *string  `tl:",omitempty:Flags:0"`
	WorkerStartTime *float64 `tl:",omitempty:Flags:1"`
	WorkerEndTime   *float64 `tl:",omitempty:Flags:1"`
}

func (*ProxyQueryFinalInfo) CRC() uint32 {
	return uint32(0xe794e6fc)
}
func (t ProxyQueryFinalInfo) MarshalTL() ([]byte, error) {
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
	if (flagsVar0>>int32(1))&1 == 1 {
		m17WorkerStartTimeBits := math.Float64bits(*t.WorkerStartTime)
		var m17WorkerStartTimeRaw [8]byte
		binary.LittleEndian.PutUint64(m17WorkerStartTimeRaw[:], m17WorkerStartTimeBits)
		_, err = buf.Write(m17WorkerStartTimeRaw[:])
		if err != nil {
			return nil, err
		}
	}
	if (flagsVar0>>int32(1))&1 == 1 {
		m18WorkerEndTimeBits := math.Float64bits(*t.WorkerEndTime)
		var m18WorkerEndTimeRaw [8]byte
		binary.LittleEndian.PutUint64(m18WorkerEndTimeRaw[:], m18WorkerEndTimeBits)
		_, err = buf.Write(m18WorkerEndTimeRaw[:])
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
func (t *ProxyQueryFinalInfo) UnmarshalTL(r io.Reader) error {
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
	if (t.Flags>>int32(1))&1 == 1 {
		var tmp7WorkerStartTimeBits uint64
		err = tl.Unmarshal(r, &tmp7WorkerStartTimeBits)
		if err != nil {
			return err
		}
		tmp7WorkerStartTimeVal := math.Float64frombits(tmp7WorkerStartTimeBits)
		t.WorkerStartTime = &tmp7WorkerStartTimeVal
	}
	if (t.Flags>>int32(1))&1 == 1 {
		var tmp8WorkerEndTimeBits uint64
		err = tl.Unmarshal(r, &tmp8WorkerEndTimeBits)
		if err != nil {
			return err
		}
		tmp8WorkerEndTimeVal := math.Float64frombits(tmp8WorkerEndTimeBits)
		t.WorkerEndTime = &tmp8WorkerEndTimeVal
	}
	return nil
}

type IProxySignedPayment interface {
	CRC() uint32
	MarshalTL() ([]byte, error)
	UnmarshalTL(io.Reader) error
	_IProxySignedPayment()
}

var (
	_ IProxySignedPayment = (*ProxySignedPayment)(nil)
	_ IProxySignedPayment = (*ProxySignedPaymentEmpty)(nil)
)

func decodeIProxySignedPayment(r io.Reader) (IProxySignedPayment, error) {
	var tag uint32
	err := tl.Unmarshal(r, &tag)
	if err != nil {
		return nil, err
	}
	var res IProxySignedPayment
	switch tag {
	case uint32(0x2998182):
		res = &ProxySignedPayment{}
	case uint32(0xb347ce64):
		res = &ProxySignedPaymentEmpty{}
	default:
		return nil, fmt.Errorf("invalid crc code: got 0x%08x", tag)
	}
	err = res.UnmarshalTL(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func DecodeIProxySignedPayment(r io.Reader) (IProxySignedPayment, error) {
	return decodeIProxySignedPayment(r)
}

type ProxySignedPayment struct {
	Data []byte
}

func (*ProxySignedPayment) CRC() uint32 {
	return uint32(0x2998182)
}
func (t ProxySignedPayment) MarshalTL() ([]byte, error) {
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
func (t *ProxySignedPayment) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Data)
	if err != nil {
		return err
	}
	return nil
}
func (*ProxySignedPayment) _IProxySignedPayment() {}

type ProxySignedPaymentEmpty struct{}

func (*ProxySignedPaymentEmpty) CRC() uint32 {
	return uint32(0xb347ce64)
}
func (t ProxySignedPaymentEmpty) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	return buf.Bytes(), nil
}
func (t *ProxySignedPaymentEmpty) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	return nil
}
func (*ProxySignedPaymentEmpty) _IProxySignedPayment() {}

type IProxyWorkerRequestPayment interface {
	CRC() uint32
	MarshalTL() ([]byte, error)
	UnmarshalTL(io.Reader) error
	_IProxyWorkerRequestPayment()
}

var (
	_ IProxyWorkerRequestPayment = (*ProxyClientRequestPayment)(nil)
	_ IProxyWorkerRequestPayment = (*ProxyWorkerRequestPayment)(nil)
)

func decodeIProxyWorkerRequestPayment(r io.Reader) (IProxyWorkerRequestPayment, error) {
	var tag uint32
	err := tl.Unmarshal(r, &tag)
	if err != nil {
		return nil, err
	}
	var res IProxyWorkerRequestPayment
	switch tag {
	case uint32(0xc6bcf127):
		res = &ProxyClientRequestPayment{}
	case uint32(0x1436c0d5):
		res = &ProxyWorkerRequestPayment{}
	default:
		return nil, fmt.Errorf("invalid crc code: got 0x%08x", tag)
	}
	err = res.UnmarshalTL(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func DecodeIProxyWorkerRequestPayment(r io.Reader) (IProxyWorkerRequestPayment, error) {
	return decodeIProxyWorkerRequestPayment(r)
}

type ProxyClientRequestPayment struct {
	RequestID     [32]byte
	SignedPayment IProxySignedPayment
	DbTokens      int64
	MaxTokens     int64
	RequestTokens int64
}

func (*ProxyClientRequestPayment) CRC() uint32 {
	return uint32(0xc6bcf127)
}
func (t ProxyClientRequestPayment) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	_, err = buf.Write(t.RequestID[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	m7SignedPayment := t.SignedPayment
	if m7SignedPayment == nil {
		return nil, fmt.Errorf("nil %s", "ProxySignedPayment")
	}
	b, err = tl.Marshal(m7SignedPayment.CRC())
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = m7SignedPayment.MarshalTL()
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
	b, err = tl.Marshal(t.RequestTokens)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ProxyClientRequestPayment) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	_, err = io.ReadFull(r, t.RequestID[:])
	if err != nil {
		return err
	}
	tmp4SignedPayment, err := decodeIProxySignedPayment(r)
	if err != nil {
		return err
	}
	t.SignedPayment = tmp4SignedPayment
	err = tl.Unmarshal(r, &t.DbTokens)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.MaxTokens)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.RequestTokens)
	if err != nil {
		return err
	}
	return nil
}
func (*ProxyClientRequestPayment) _IProxyWorkerRequestPayment() {}

type ProxyWorkerRequestPayment struct {
	RequestID          [32]byte
	SignedPayment      IProxySignedPayment
	DbTokens           int64
	MaxTokens          int64
	RequestTokens      int64
	RequestClientOwner string
}

func (*ProxyWorkerRequestPayment) CRC() uint32 {
	return uint32(0x1436c0d5)
}
func (t ProxyWorkerRequestPayment) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	_, err = buf.Write(t.RequestID[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	m7SignedPayment := t.SignedPayment
	if m7SignedPayment == nil {
		return nil, fmt.Errorf("nil %s", "ProxySignedPayment")
	}
	b, err = tl.Marshal(m7SignedPayment.CRC())
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = m7SignedPayment.MarshalTL()
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
	b, err = tl.Marshal(t.RequestTokens)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.RequestClientOwner)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *ProxyWorkerRequestPayment) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	_, err = io.ReadFull(r, t.RequestID[:])
	if err != nil {
		return err
	}
	tmp4SignedPayment, err := decodeIProxySignedPayment(r)
	if err != nil {
		return err
	}
	t.SignedPayment = tmp4SignedPayment
	err = tl.Unmarshal(r, &t.DbTokens)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.MaxTokens)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.RequestTokens)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.RequestClientOwner)
	if err != nil {
		return err
	}
	return nil
}
func (*ProxyWorkerRequestPayment) _IProxyWorkerRequestPayment() {}
