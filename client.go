package gocoon

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/tonkeeper/gocoon/contracts/client_contract"
	"github.com/tonkeeper/gocoon/internal/session"
	"github.com/tonkeeper/gocoon/net/proxyconn"
	"github.com/tonkeeper/gocoon/tlcocoon"
	tlcocoonTypes "github.com/tonkeeper/gocoon/tlcocoon/types"
	"github.com/tonkeeper/tongo"
	abiCocoon "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/cocoon"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"go.uber.org/zap"
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

// Opts holds configuration for Client.
type Opts struct {
	RootAddress ton.AccountID
	LiteClient  LiteClient
}

// Client connects to a cocoon proxy and handles the full auth flow.
type Client struct {
	wallet    Wallet
	secret    string
	opts      Opts
	rootStore abiCocoon.RootStorage
}

// New creates a client, fetching the root contract state immediately.
func New(ctx context.Context, wallet Wallet, secret string, opts Opts) (*Client, error) {
	if (opts.RootAddress == ton.AccountID{}) {
		opts.RootAddress = defaultRootAddr.ID
	}
	c := &Client{wallet: wallet, secret: secret, opts: opts}
	if err := c.UpdateRootState(ctx); err != nil {
		return nil, fmt.Errorf("fetch root state: %w", err)
	}
	return c, nil
}

// NewCocoonClient is an alias for New.
func NewCocoonClient(ctx context.Context, wallet Wallet, secret string, opts Opts) (*Client, error) {
	return New(ctx, wallet, secret, opts)
}

// UpdateRootState re-fetches the root contract state, refreshing the proxy list and params.
func (c *Client) UpdateRootState(ctx context.Context) error {
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
func (c *Client) Connect(ctx context.Context, logger *zap.Logger) (*Connection, error) {
	proxies := c.rootStore.Data.Value.RegisteredProxies.Values()
	if len(proxies) == 0 {
		return nil, fmt.Errorf("no registered proxies in root contract")
	}

	proxyAddr := clientAddr(proxies[0].Address)
	logger.Info("connecting to proxy", zap.String("address", proxyAddr))

	conn, err := proxyconn.Dial(proxyAddr, logger)
	if err != nil {
		return nil, fmt.Errorf("dial proxy: %w", err)
	}

	sess, err := session.New(conn)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("new session: %w", err)
	}

	apiClient := tlcocoon.NewClient(sess)
	rootVersion := uint32(c.rootStore.Version)

	resp, err := doProxyAuth(ctx, logger, apiClient, c, rootVersion)
	if err != nil {
		_ = sess.Close()
		_ = conn.Close()
		return nil, err
	}

	var protoVersion uint32
	if resp.Params.ProtoVersion != nil {
		protoVersion = uint32(*resp.Params.ProtoVersion)
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

func doProxyAuth(ctx context.Context, logger *zap.Logger, apiClient *tlcocoon.Client, c *Client, rootVersion uint32) (tlcocoonTypes.ClientConnectedToProxy, error) {
	resp, err := apiClient.ConnectToProxy(ctx, tlcocoon.ClientConnectToProxyRequest{
		Params: tlcocoonTypes.ClientParams{
			ClientOwnerAddress: c.wallet.Address().ToHuman(false, false),
			IsTest:             boolPtr(false),
			MinProtoVersion:    int32Ptr(1),
			MaxProtoVersion:    int32Ptr(1),
		},
		MinConfigVersion: int32(rootVersion),
	})
	if err != nil {
		return tlcocoonTypes.ClientConnectedToProxy{}, fmt.Errorf("connectToProxy: %w", err)
	}

	pubkey := resp.Params.ProxyPublicKey
	logger.Info("connected to proxy",
		zap.String("proxy_owner", resp.Params.ProxyOwnerAddress),
		zap.String("proxy_sc", resp.Params.ProxyScAddress),
		zap.String("client_sc", resp.ClientScAddress),
		zap.String("proxy_pubkey", hex.EncodeToString(pubkey[:])),
	)

	clientScAddr, err := ton.ParseAccountID(resp.ClientScAddress)
	if err != nil {
		return tlcocoonTypes.ClientConnectedToProxy{}, fmt.Errorf("parse client SC address: %w", err)
	}

	clientSC := client_contract.New(clientScAddr, c.secret, c.wallet, c.opts.LiteClient)
	deployed, err := clientSC.Sync(ctx, c.rootStore.Params.Value.MinClientStake, logger)
	if err != nil {
		return tlcocoonTypes.ClientConnectedToProxy{}, fmt.Errorf("sync client contract: %w", err)
	}
	_ = deployed

	switch auth := resp.Auth.(type) {
	case *tlcocoonTypes.ClientProxyConnectionAuthShort:
		logger.Info("auth type: short",
			zap.String("secret_hash", hex.EncodeToString(auth.SecretHash[:])),
		)
		authResp, err := apiClient.AuthorizeWithProxyShort(ctx, tlcocoon.ClientAuthorizeWithProxyShortRequest{
			Data: []byte(c.secret),
		})
		if err != nil {
			return tlcocoonTypes.ClientConnectedToProxy{}, fmt.Errorf("authorizeWithProxyShort: %w", err)
		}
		if err := checkAuth(authResp, logger); err != nil {
			return tlcocoonTypes.ClientConnectedToProxy{}, err
		}
	case *tlcocoonTypes.ClientProxyConnectionAuthLong:
		logger.Info("auth type: long", zap.Uint64("nonce", uint64(auth.Nonce)))
		proxyScAddr := ton.MustParseAccountID(resp.Params.ProxyScAddress)
		if err := clientSC.Register(ctx, uint64(auth.Nonce), logger, proxyScAddr,
			resp.Params.ProxyPublicKey,
			c.rootStore.Params.Value.MinClientStake,
			c.rootStore.Params.Value,
		); err != nil {
			return tlcocoonTypes.ClientConnectedToProxy{}, fmt.Errorf("send register tx: %w", err)
		}
		logger.Info("waiting for on-chain confirmation (up to 300s)")
		longCtx, cancel := context.WithTimeout(ctx, 300*time.Second)
		defer cancel()
		authResp, err := apiClient.AuthorizeWithProxyLong(longCtx, tlcocoon.ClientAuthorizeWithProxyLongRequest{})
		if err != nil {
			return tlcocoonTypes.ClientConnectedToProxy{}, fmt.Errorf("authorizeWithProxyLong: %w", err)
		}
		if err := checkAuth(authResp, logger); err != nil {
			return tlcocoonTypes.ClientConnectedToProxy{}, err
		}
	default:
		return tlcocoonTypes.ClientConnectedToProxy{}, fmt.Errorf("unknown auth type: %T", auth)
	}

	return resp, nil
}

func checkAuth(r tlcocoonTypes.IClientAuthorizationWithProxy, logger *zap.Logger) error {
	switch s := r.(type) {
	case *tlcocoonTypes.ClientAuthorizationWithProxySuccess:
		logger.Info("auth success",
			zap.Uint64("tokens_committed", uint64(s.TokensCommittedToDb)),
			zap.Uint64("max_tokens", uint64(s.MaxTokens)),
		)
		return nil
	case *tlcocoonTypes.ClientAuthorizationWithProxyFailed:
		return fmt.Errorf("authorization failed (code %d): %s", s.ErrorCode, s.Error)
	default:
		return fmt.Errorf("unknown auth result: %T", r)
	}
}

func clientAddr(addr string) string {
	parts := strings.Fields(addr)
	if len(parts) == 2 {
		return parts[1]
	}
	return parts[0]
}

func boolPtr(v bool) *bool    { return &v }
func int32Ptr(v int32) *int32 { return &v }
