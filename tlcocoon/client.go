package tlcocoon

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	types "github.com/tonkeeper/gocoon/tlcocoon/types"
	tl "github.com/tonkeeper/tongo/tl"
	"math"
)

type ClientAuthorizeWithProxyLongRequest struct{}

func (*ClientAuthorizeWithProxyLongRequest) CRC() uint32 {
	return uint32(0xd3474303)
}
func (t ClientAuthorizeWithProxyLongRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	return buf.Bytes(), nil
}
func ClientAuthorizeWithProxyLong(ctx context.Context, m Requester, i ClientAuthorizeWithProxyLongRequest) (types.IClientAuthorizationWithProxy, error) {
	respRaw, err := requestRaw(ctx, m, &i)
	if err != nil {
		var zero types.IClientAuthorizationWithProxy
		return zero, err
	}
	res, err := types.DecodeIClientAuthorizationWithProxy(bytes.NewReader(respRaw))
	if err != nil {
		var zero types.IClientAuthorizationWithProxy
		return zero, fmt.Errorf("response: %w", err)
	}
	return res, nil
}

type ClientAuthorizeWithProxyShortRequest struct {
	Data []byte
}

func (*ClientAuthorizeWithProxyShortRequest) CRC() uint32 {
	return uint32(0x6c276723)
}
func (t ClientAuthorizeWithProxyShortRequest) MarshalTL() ([]byte, error) {
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
func ClientAuthorizeWithProxyShort(ctx context.Context, m Requester, i ClientAuthorizeWithProxyShortRequest) (types.IClientAuthorizationWithProxy, error) {
	respRaw, err := requestRaw(ctx, m, &i)
	if err != nil {
		var zero types.IClientAuthorizationWithProxy
		return zero, err
	}
	res, err := types.DecodeIClientAuthorizationWithProxy(bytes.NewReader(respRaw))
	if err != nil {
		var zero types.IClientAuthorizationWithProxy
		return zero, fmt.Errorf("response: %w", err)
	}
	return res, nil
}

type ClientConnectToProxyRequest struct {
	Params           types.ClientParams
	MinConfigVersion int32
}

func (*ClientConnectToProxyRequest) CRC() uint32 {
	return uint32(0xff5fa0f4)
}
func (t ClientConnectToProxyRequest) MarshalTL() ([]byte, error) {
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
func ClientConnectToProxy(ctx context.Context, m Requester, i ClientConnectToProxyRequest) (types.ClientConnectedToProxy, error) {
	var res types.ClientConnectedToProxy
	return res, request(ctx, m, &i, &res)
}

type ClientGetWorkerTypesRequest struct{}

func (*ClientGetWorkerTypesRequest) CRC() uint32 {
	return uint32(0x7f062bdb)
}
func (t ClientGetWorkerTypesRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	return buf.Bytes(), nil
}
func ClientGetWorkerTypes(ctx context.Context, m Requester, i ClientGetWorkerTypesRequest) (types.ClientWorkerTypes, error) {
	var res types.ClientWorkerTypes
	return res, request(ctx, m, &i, &res)
}

type ClientGetWorkerTypesV2Request struct{}

func (*ClientGetWorkerTypesV2Request) CRC() uint32 {
	return uint32(0xb2133d72)
}
func (t ClientGetWorkerTypesV2Request) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	return buf.Bytes(), nil
}
func ClientGetWorkerTypesV2(ctx context.Context, m Requester, i ClientGetWorkerTypesV2Request) (types.ClientWorkerTypesV2, error) {
	var res types.ClientWorkerTypesV2
	return res, request(ctx, m, &i, &res)
}

type ClientRequestRefundRequest struct{}

func (*ClientRequestRefundRequest) CRC() uint32 {
	return uint32(0x238d863d)
}
func (t ClientRequestRefundRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	return buf.Bytes(), nil
}
func ClientRequestRefund(ctx context.Context, m Requester, i ClientRequestRefundRequest) (types.IClientRefund, error) {
	respRaw, err := requestRaw(ctx, m, &i)
	if err != nil {
		var zero types.IClientRefund
		return zero, err
	}
	res, err := types.DecodeIClientRefund(bytes.NewReader(respRaw))
	if err != nil {
		var zero types.IClientRefund
		return zero, fmt.Errorf("response: %w", err)
	}
	return res, nil
}

type ClientRunQueryRequest struct {
	ModelName        string
	Query            []byte
	MaxCoefficient   int32
	MaxTokens        int32
	Timeout          float64
	RequestID        [32]byte
	MinConfigVersion int32
}

func (*ClientRunQueryRequest) CRC() uint32 {
	return uint32(0xbc748f32)
}
func (t ClientRunQueryRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
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
	m20TimeoutBits := math.Float64bits(t.Timeout)
	var m20TimeoutRaw [8]byte
	binary.LittleEndian.PutUint64(m20TimeoutRaw[:], m20TimeoutBits)
	_, err = buf.Write(m20TimeoutRaw[:])
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(t.RequestID[:])
	if err != nil {
		return nil, err
	}
	_ = 32
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
func ClientRunQuery(ctx context.Context, m Requester, i ClientRunQueryRequest) (types.IClientQueryAnswer, error) {
	respRaw, err := requestRaw(ctx, m, &i)
	if err != nil {
		var zero types.IClientQueryAnswer
		return zero, err
	}
	res, err := types.DecodeIClientQueryAnswer(bytes.NewReader(respRaw))
	if err != nil {
		var zero types.IClientQueryAnswer
		return zero, fmt.Errorf("response: %w", err)
	}
	return res, nil
}

type ClientRunQueryExRequest struct {
	ModelName        string
	Query            []byte
	MaxCoefficient   int32
	MaxTokens        int32
	Timeout          float64
	RequestID        [32]byte
	MinConfigVersion int32
	Flags            uint32 `tl:"0,bitflag"`
	EnableDebug      *bool  `tl:",omitempty:Flags:0"`
}

func (*ClientRunQueryExRequest) CRC() uint32 {
	return uint32(0xf54cb74b)
}
func (t ClientRunQueryExRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	flagsVar0 := t.Flags
	if t.EnableDebug != nil {
		flagsVar0 |= uint32(0x1)
	}
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
	m22TimeoutBits := math.Float64bits(t.Timeout)
	var m22TimeoutRaw [8]byte
	binary.LittleEndian.PutUint64(m22TimeoutRaw[:], m22TimeoutBits)
	_, err = buf.Write(m22TimeoutRaw[:])
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(t.RequestID[:])
	if err != nil {
		return nil, err
	}
	_ = 32
	b, err = tl.Marshal(t.MinConfigVersion)
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
func ClientRunQueryEx(ctx context.Context, m Requester, i ClientRunQueryExRequest) (types.IClientQueryAnswerEx, error) {
	respRaw, err := requestRaw(ctx, m, &i)
	if err != nil {
		var zero types.IClientQueryAnswerEx
		return zero, err
	}
	res, err := types.DecodeIClientQueryAnswerEx(bytes.NewReader(respRaw))
	if err != nil {
		var zero types.IClientQueryAnswerEx
		return zero, fmt.Errorf("response: %w", err)
	}
	return res, nil
}

type ClientUpdatePaymentStatusRequest struct{}

func (*ClientUpdatePaymentStatusRequest) CRC() uint32 {
	return uint32(0x9ed1c697)
}
func (t ClientUpdatePaymentStatusRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	return buf.Bytes(), nil
}
func ClientUpdatePaymentStatus(ctx context.Context, m Requester, i ClientUpdatePaymentStatusRequest) (types.ClientPaymentStatus, error) {
	var res types.ClientPaymentStatus
	return res, request(ctx, m, &i, &res)
}

type ClientAPI interface {
	AuthorizeWithProxyLong(ctx context.Context, i ClientAuthorizeWithProxyLongRequest) (types.IClientAuthorizationWithProxy, error)
	AuthorizeWithProxyShort(ctx context.Context, i ClientAuthorizeWithProxyShortRequest) (types.IClientAuthorizationWithProxy, error)
	ConnectToProxy(ctx context.Context, i ClientConnectToProxyRequest) (types.ClientConnectedToProxy, error)
	GetWorkerTypes(ctx context.Context, i ClientGetWorkerTypesRequest) (types.ClientWorkerTypes, error)
	GetWorkerTypesV2(ctx context.Context, i ClientGetWorkerTypesV2Request) (types.ClientWorkerTypesV2, error)
	RequestRefund(ctx context.Context, i ClientRequestRefundRequest) (types.IClientRefund, error)
	RunQuery(ctx context.Context, i ClientRunQueryRequest) (types.IClientQueryAnswer, error)
	RunQueryEx(ctx context.Context, i ClientRunQueryExRequest) (types.IClientQueryAnswerEx, error)
	UpdatePaymentStatus(ctx context.Context, i ClientUpdatePaymentStatusRequest) (types.ClientPaymentStatus, error)
}
type Client struct {
	requester Requester
}

func NewClient(requester Requester) *Client {
	return &Client{requester: requester}
}
func (c *Client) AuthorizeWithProxyLong(ctx context.Context, i ClientAuthorizeWithProxyLongRequest) (types.IClientAuthorizationWithProxy, error) {
	return ClientAuthorizeWithProxyLong(ctx, c.requester, i)
}
func (c *Client) AuthorizeWithProxyShort(ctx context.Context, i ClientAuthorizeWithProxyShortRequest) (types.IClientAuthorizationWithProxy, error) {
	return ClientAuthorizeWithProxyShort(ctx, c.requester, i)
}
func (c *Client) ConnectToProxy(ctx context.Context, i ClientConnectToProxyRequest) (types.ClientConnectedToProxy, error) {
	return ClientConnectToProxy(ctx, c.requester, i)
}
func (c *Client) GetWorkerTypes(ctx context.Context, i ClientGetWorkerTypesRequest) (types.ClientWorkerTypes, error) {
	return ClientGetWorkerTypes(ctx, c.requester, i)
}
func (c *Client) GetWorkerTypesV2(ctx context.Context, i ClientGetWorkerTypesV2Request) (types.ClientWorkerTypesV2, error) {
	return ClientGetWorkerTypesV2(ctx, c.requester, i)
}
func (c *Client) RequestRefund(ctx context.Context, i ClientRequestRefundRequest) (types.IClientRefund, error) {
	return ClientRequestRefund(ctx, c.requester, i)
}
func (c *Client) RunQuery(ctx context.Context, i ClientRunQueryRequest) (types.IClientQueryAnswer, error) {
	return ClientRunQuery(ctx, c.requester, i)
}
func (c *Client) RunQueryEx(ctx context.Context, i ClientRunQueryExRequest) (types.IClientQueryAnswerEx, error) {
	return ClientRunQueryEx(ctx, c.requester, i)
}
func (c *Client) UpdatePaymentStatus(ctx context.Context, i ClientUpdatePaymentStatusRequest) (types.ClientPaymentStatus, error) {
	return ClientUpdatePaymentStatus(ctx, c.requester, i)
}

var _ ClientAPI = (*Client)(nil)
