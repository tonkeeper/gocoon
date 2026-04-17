package tlcocoon

import (
	"bytes"
	"context"
	types "github.com/tonkeeper/gocoon/pkg/tlcocoon/types"
	tl "github.com/tonkeeper/tongo/tl"
)

type WorkerCompareBalanceWithProxyRequest struct {
	TokensCommittedToBlockchain int64
	TokensCommittedToDb         int64
	MaxTokens                   int64
}

func (*WorkerCompareBalanceWithProxyRequest) CRC() uint32 {
	return uint32(0xdb386501)
}
func (t WorkerCompareBalanceWithProxyRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.TokensCommittedToBlockchain)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
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
func WorkerCompareBalanceWithProxy(ctx context.Context, m Requester, i WorkerCompareBalanceWithProxyRequest) (types.WorkerCompareBalanceWithProxyResult, error) {
	var res types.WorkerCompareBalanceWithProxyResult
	return res, request(ctx, m, &i, &res)
}

type WorkerConnectToProxyRequest struct {
	Params types.WorkerParams
}

func (*WorkerConnectToProxyRequest) CRC() uint32 {
	return uint32(0xcc484c5d)
}
func (t WorkerConnectToProxyRequest) MarshalTL() ([]byte, error) {
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
	return buf.Bytes(), nil
}
func WorkerConnectToProxy(ctx context.Context, m Requester, i WorkerConnectToProxyRequest) (types.WorkerConnectedToProxy, error) {
	var res types.WorkerConnectedToProxy
	return res, request(ctx, m, &i, &res)
}

type WorkerExtendedCompareBalanceWithProxyRequest struct {
	TokensCommittedToDb int64
	Other               []byte
}

func (*WorkerExtendedCompareBalanceWithProxyRequest) CRC() uint32 {
	return uint32(0xf8dd9914)
}
func (t WorkerExtendedCompareBalanceWithProxyRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.TokensCommittedToDb)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Other)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func WorkerExtendedCompareBalanceWithProxy(ctx context.Context, m Requester, i WorkerExtendedCompareBalanceWithProxyRequest) (types.WorkerExtendedCompareBalanceWithProxyResult, error) {
	var res types.WorkerExtendedCompareBalanceWithProxyResult
	return res, request(ctx, m, &i, &res)
}

type WorkerProxyHandshakeCompleteRequest struct {
	IsDisabled bool
}

func (*WorkerProxyHandshakeCompleteRequest) CRC() uint32 {
	return uint32(0xf843b066)
}
func (t WorkerProxyHandshakeCompleteRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.IsDisabled)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func WorkerProxyHandshakeComplete(ctx context.Context, m Requester, i WorkerProxyHandshakeCompleteRequest) (types.WorkerProxyHandshakeCompleted, error) {
	var res types.WorkerProxyHandshakeCompleted
	return res, request(ctx, m, &i, &res)
}

type WorkerUpdatePaymentStatusRequest struct{}

func (*WorkerUpdatePaymentStatusRequest) CRC() uint32 {
	return uint32(0xdc70a524)
}
func (t WorkerUpdatePaymentStatusRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	return buf.Bytes(), nil
}
func WorkerUpdatePaymentStatus(ctx context.Context, m Requester, i WorkerUpdatePaymentStatusRequest) (types.WorkerPaymentStatus, error) {
	var res types.WorkerPaymentStatus
	return res, request(ctx, m, &i, &res)
}

type WorkerAPI interface {
	CompareBalanceWithProxy(ctx context.Context, i WorkerCompareBalanceWithProxyRequest) (types.WorkerCompareBalanceWithProxyResult, error)
	ConnectToProxy(ctx context.Context, i WorkerConnectToProxyRequest) (types.WorkerConnectedToProxy, error)
	ExtendedCompareBalanceWithProxy(ctx context.Context, i WorkerExtendedCompareBalanceWithProxyRequest) (types.WorkerExtendedCompareBalanceWithProxyResult, error)
	ProxyHandshakeComplete(ctx context.Context, i WorkerProxyHandshakeCompleteRequest) (types.WorkerProxyHandshakeCompleted, error)
	UpdatePaymentStatus(ctx context.Context, i WorkerUpdatePaymentStatusRequest) (types.WorkerPaymentStatus, error)
}
type Worker struct {
	requester Requester
}

func NewWorker(requester Requester) *Worker {
	return &Worker{requester: requester}
}
func (c *Worker) CompareBalanceWithProxy(ctx context.Context, i WorkerCompareBalanceWithProxyRequest) (types.WorkerCompareBalanceWithProxyResult, error) {
	return WorkerCompareBalanceWithProxy(ctx, c.requester, i)
}
func (c *Worker) ConnectToProxy(ctx context.Context, i WorkerConnectToProxyRequest) (types.WorkerConnectedToProxy, error) {
	return WorkerConnectToProxy(ctx, c.requester, i)
}
func (c *Worker) ExtendedCompareBalanceWithProxy(ctx context.Context, i WorkerExtendedCompareBalanceWithProxyRequest) (types.WorkerExtendedCompareBalanceWithProxyResult, error) {
	return WorkerExtendedCompareBalanceWithProxy(ctx, c.requester, i)
}
func (c *Worker) ProxyHandshakeComplete(ctx context.Context, i WorkerProxyHandshakeCompleteRequest) (types.WorkerProxyHandshakeCompleted, error) {
	return WorkerProxyHandshakeComplete(ctx, c.requester, i)
}
func (c *Worker) UpdatePaymentStatus(ctx context.Context, i WorkerUpdatePaymentStatusRequest) (types.WorkerPaymentStatus, error) {
	return WorkerUpdatePaymentStatus(ctx, c.requester, i)
}

var _ WorkerAPI = (*Worker)(nil)
