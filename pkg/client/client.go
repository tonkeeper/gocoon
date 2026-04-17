package client

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/tonkeeper/tongo"
	abiCocoon "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/cocoon"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"go.uber.org/zap"

	"github.com/tonkeeper/gocoon/pkg/client/contract"
	"github.com/tonkeeper/gocoon/pkg/cocoon"
	"github.com/tonkeeper/gocoon/pkg/tlcocoonapi"
)

var defaultRootAddr = tongo.MustParseAddress("EQCns7bYSp0igFvS1wpb5wsZjCKCV19MD5AVzI4EyxsnU73k")

// LiteClient is the minimal TON liteapi interface the client needs.
type LiteClient interface {
	RunSmcMethodByID(ctx context.Context, accountID ton.AccountID, methodID int, params tlb.VmStack) (uint32, tlb.VmStack, error)
	GetAccountState(ctx context.Context, accountID ton.AccountID) (tlb.ShardAccount, error)
}

// Wallet is the subset of cocoonWallet.Wallet used by the client.
type Wallet interface {
	Address() ton.AccountID
	ForwardMessage(ctx context.Context, msg tlb.Message) error
}

// Opts holds configuration for CocoonClient.
type Opts struct {
	RootAddress ton.AccountID
	LiteClient  LiteClient
}

// CocoonClient connects to a cocoon proxy and handles the full auth flow.
type CocoonClient struct {
	wallet    Wallet
	secret    string
	opts      Opts
	rootStore abiCocoon.RootStorage
}

// NewCocoonClient creates a client, fetching the root contract state immediately.
func NewCocoonClient(ctx context.Context, wallet Wallet, secret string, opts Opts) (*CocoonClient, error) {
	if (opts.RootAddress == ton.AccountID{}) {
		opts.RootAddress = defaultRootAddr.ID
	}
	c := &CocoonClient{wallet: wallet, secret: secret, opts: opts}
	if err := c.UpdateRootState(ctx); err != nil {
		return nil, fmt.Errorf("fetch root state: %w", err)
	}
	return c, nil
}

// UpdateRootState re-fetches the root contract state, refreshing the proxy list and params.
func (c *CocoonClient) UpdateRootState(ctx context.Context) error {
	rootClient := abiCocoon.NewCocoonRoot(c.opts.LiteClient, c.opts.LiteClient).WithAccountId(c.opts.RootAddress)
	_, rootStore, err := rootClient.AccountState(ctx)
	if err != nil {
		return fmt.Errorf("get root storage: %w", err)
	}
	c.rootStore = rootStore
	return nil
}

// Connect dials the proxy, performs the handshake and authorization, and
// returns a ready-to-use Connection.
func (c *CocoonClient) Connect(ctx context.Context, logger *zap.Logger) (*Connection, error) {
	proxies := c.rootStore.Data.Value.RegisteredProxies.Values()
	if len(proxies) == 0 {
		return nil, fmt.Errorf("no registered proxies in root contract")
	}

	proxyAddr := clientAddr(proxies[0].Address)
	logger.Info("connecting to proxy", zap.String("address", proxyAddr))

	conn, err := cocoon.Dial(proxyAddr, logger)
	if err != nil {
		return nil, fmt.Errorf("dial proxy: %w", err)
	}

	sess, err := cocoon.NewSession(conn)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("new session: %w", err)
	}

	apiClient := tlcocoonapi.NewClient(sess.Query)
	rootVersion := uint32(c.rootStore.Version)

	resp, err := apiClient.ClientConnectToProxy(ctx, tlcocoonapi.ClientConnectToProxyRequest{
		Params: tlcocoonapi.ClientParamsC{
			Flags:              3,
			ClientOwnerAddress: c.wallet.Address().ToHuman(false, false),
			IsTest:             boolPtr(false),
			MinProtoVersion:    uint32Ptr(1),
			MaxProtoVersion:    uint32Ptr(1),
		},
		MinConfigVersion: rootVersion,
	})
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("connectToProxy: %w", err)
	}

	pubkey := resp.Params.ProxyPublicKey
	logger.Info("connected to proxy",
		zap.String("proxy_owner", resp.Params.ProxyOwnerAddress),
		zap.String("proxy_sc", resp.Params.ProxyScAddress),
		zap.String("client_sc", resp.ClientScAddress),
		zap.String("proxy_pubkey", hex.EncodeToString(pubkey[:])))

	clientScAddr, err := ton.ParseAccountID(resp.ClientScAddress)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("parse client SC address: %w", err)
	}

	clientSC := contract.New(clientScAddr, c.secret, c.wallet, c.opts.LiteClient)
	deployed, err := clientSC.Sync(ctx, c.rootStore.Params.Value.MinClientStake, logger)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("sync client contract: %w", err)
	}
	_ = deployed

	auth := resp.Auth
	switch auth.SumType {
	case "ClientProxyConnectionAuthShort":
		short := auth.ClientProxyConnectionAuthShort
		logger.Info("auth type: short",
			zap.String("secret_hash", hex.EncodeToString(short.SecretHash[:])))
		authResp, err := apiClient.ClientAuthorizeWithProxyShort(ctx,
			tlcocoonapi.ClientAuthorizeWithProxyShortRequest{Data: []byte(c.secret)})
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("authorizeWithProxyShort: %w", err)
		}
		if err := checkAuth(authResp, logger); err != nil {
			conn.Close()
			return nil, err
		}

	case "ClientProxyConnectionAuthLong":
		long := auth.ClientProxyConnectionAuthLong
		logger.Info("auth type: long", zap.Uint64("nonce", long.Nonce))
		proxyScAddr := ton.MustParseAccountID(resp.Params.ProxyScAddress)
		if err := clientSC.Register(ctx, long.Nonce, logger, proxyScAddr,
			resp.Params.ProxyPublicKey,
			c.rootStore.Params.Value.MinClientStake,
			c.rootStore.Params.Value,
		); err != nil {
			conn.Close()
			return nil, fmt.Errorf("send register tx: %w", err)
		}
		logger.Info("waiting for on-chain confirmation (up to 300s)")
		longCtx, cancel := context.WithTimeout(ctx, 300*time.Second)
		defer cancel()
		authResp, err := apiClient.ClientAuthorizeWithProxyLong(longCtx)
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("authorizeWithProxyLong: %w", err)
		}
		if err := checkAuth(authResp, logger); err != nil {
			conn.Close()
			return nil, err
		}

	default:
		conn.Close()
		return nil, fmt.Errorf("unknown auth type: %s", auth.SumType)
	}

	var protoVersion uint32
	if resp.Params.ProtoVersion != nil {
		protoVersion = *resp.Params.ProtoVersion
	}

	return &Connection{
		conn:         conn,
		sess:         sess,
		apiClient:    apiClient,
		rootVersion:  rootVersion,
		protoVersion: protoVersion,
		logger:       logger,
	}, nil
}

func checkAuth(r tlcocoonapi.ClientAuthorizationWithProxy, logger *zap.Logger) error {
	switch r.SumType {
	case "ClientAuthorizationWithProxySuccess":
		s := r.ClientAuthorizationWithProxySuccess
		logger.Info("auth success",
			zap.Uint64("tokens_committed", s.TokensCommittedToDb),
			zap.Uint64("max_tokens", s.MaxTokens))
		return nil
	case "ClientAuthorizationWithProxyFailed":
		f := r.ClientAuthorizationWithProxyFailed
		return fmt.Errorf("authorization failed (code %d): %s", f.ErrorCode, f.Error)
	default:
		return fmt.Errorf("unknown auth result: %s", r.SumType)
	}
}

func clientAddr(addr string) string {
	parts := strings.Fields(addr)
	if len(parts) == 2 {
		return parts[1]
	}
	return parts[0]
}

func boolPtr(v bool) *bool       { return &v }
func uint32Ptr(v uint32) *uint32 { return &v }
