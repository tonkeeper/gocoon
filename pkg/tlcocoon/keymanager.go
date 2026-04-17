package tlcocoon

import (
	"bytes"
	"context"
	types "github.com/tonkeeper/gocoon/pkg/tlcocoon/types"
	tl "github.com/tonkeeper/tongo/tl"
)

type KeyManagerGetProxyPrivateKeysRequest struct {
	MinConfigVersion int32
}

func (*KeyManagerGetProxyPrivateKeysRequest) CRC() uint32 {
	return uint32(0x47a9c2de)
}
func (t KeyManagerGetProxyPrivateKeysRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
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
func KeyManagerGetProxyPrivateKeys(ctx context.Context, m Requester, i KeyManagerGetProxyPrivateKeysRequest) (types.KeyManagerPrivateKeys, error) {
	var res types.KeyManagerPrivateKeys
	return res, request(ctx, m, &i, &res)
}

type KeyManagerGetWorkerPrivateKeysRequest struct {
	MinConfigVersion int32
}

func (*KeyManagerGetWorkerPrivateKeysRequest) CRC() uint32 {
	return uint32(0x2a26eff5)
}
func (t KeyManagerGetWorkerPrivateKeysRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
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
func KeyManagerGetWorkerPrivateKeys(ctx context.Context, m Requester, i KeyManagerGetWorkerPrivateKeysRequest) (types.KeyManagerPrivateKeys, error) {
	var res types.KeyManagerPrivateKeys
	return res, request(ctx, m, &i, &res)
}

type KeyManagerAPI interface {
	GetProxyPrivateKeys(ctx context.Context, i KeyManagerGetProxyPrivateKeysRequest) (types.KeyManagerPrivateKeys, error)
	GetWorkerPrivateKeys(ctx context.Context, i KeyManagerGetWorkerPrivateKeysRequest) (types.KeyManagerPrivateKeys, error)
}
type KeyManager struct {
	requester Requester
}

func NewKeyManager(requester Requester) *KeyManager {
	return &KeyManager{requester: requester}
}
func (c *KeyManager) GetProxyPrivateKeys(ctx context.Context, i KeyManagerGetProxyPrivateKeysRequest) (types.KeyManagerPrivateKeys, error) {
	return KeyManagerGetProxyPrivateKeys(ctx, c.requester, i)
}
func (c *KeyManager) GetWorkerPrivateKeys(ctx context.Context, i KeyManagerGetWorkerPrivateKeysRequest) (types.KeyManagerPrivateKeys, error) {
	return KeyManagerGetWorkerPrivateKeys(ctx, c.requester, i)
}

var _ KeyManagerAPI = (*KeyManager)(nil)
