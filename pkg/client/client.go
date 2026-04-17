package client

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand/v2"
	"strings"
	"time"

	tonapi "github.com/tonkeeper/tonapi-go"
	tongo "github.com/tonkeeper/tongo"
	abiCocoon "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/cocoon"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"go.uber.org/zap"

	"github.com/tonkeeper/gocoon/pkg/cocoon"
	"github.com/tonkeeper/gocoon/pkg/tlcocoonapi"
	"github.com/tonkeeper/gocoon/pkg/wallet"
)

const defaultRootAddr = "EQCns7bYSp0igFvS1wpb5wsZjCKCV19MD5AVzI4EyxsnU73k"

// BlockchainSender can broadcast a signed BOC to the TON network.
// *tonapi.Client satisfies this interface.
type BlockchainSender interface {
	SendBlockchainMessage(ctx context.Context, request *tonapi.SendBlockchainMessageReq) error
}

// Opts holds optional configuration for CocoonClient.
type Opts struct {
	rootAddress ton.AccountID
	secret      string
}

func (o Opts) WithRootAddress(addr ton.AccountID) Opts {
	o.rootAddress = addr
	return o
}

func (o Opts) WithSecret(secret string) Opts {
	o.secret = secret
	return o
}

// CocoonClient connects to a cocoon proxy and handles the full auth flow.
type CocoonClient struct {
	w      *wallet.Wallet
	sender BlockchainSender
	opts   Opts
}

// NewCocoonClient creates a client for the given wallet identity.
// opts is optional — call Opts{}.WithSecret(...) etc. to configure.
// Root address defaults to the production cocoon root contract.
func NewCocoonClient(w *wallet.Wallet, sender BlockchainSender, opts Opts) *CocoonClient {
	if (opts.rootAddress == ton.AccountID{}) {
		opts.rootAddress = tongo.MustParseAddress(defaultRootAddr).ID
	}
	return &CocoonClient{w: w, sender: sender, opts: opts}
}

// Connect dials the proxy, performs the handshake and authorization, and
// returns a ready-to-use Connection.
func (c *CocoonClient) Connect(ctx context.Context, logger *zap.Logger) (*Connection, error) {
	// Lite API for reading on-chain state and submitting seqno queries.
	lc, err := liteapi.NewClientWithDefaultMainnet()
	if err != nil {
		return nil, fmt.Errorf("create liteapi client: %w", err)
	}

	// Read registered proxies from the root contract.
	rootClient := abiCocoon.NewCocoonRoot(lc, lc).WithAccountId(c.opts.rootAddress)
	_, rootStore, err := rootClient.AccountState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get root storage: %w", err)
	}
	proxies := rootStore.Data.Value.RegisteredProxies.Values()
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
	rootVersion := uint32(rootStore.Version)

	resp, err := apiClient.ClientConnectToProxy(ctx, tlcocoonapi.ClientConnectToProxyRequest{
		Params: tlcocoonapi.ClientParamsC{
			Flags:              3, // bit0=is_test, bit1=proto versions
			ClientOwnerAddress: c.w.Address.ToHuman(false, false),
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

	clientContr := abiCocoon.NewCocoonClient(lc, lc).WithAccountId(clientScAddr)
	accState, _, err := clientContr.AccountState(ctx)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("get client storage: %w", err)
	}
	if accState.Account.SumType == "Account" {
		clientData, err := clientContr.GetCocoonClientData(ctx)
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("get cocoon client data: %w", err)
		}

		if c.opts.secret == "" {
			conn.Close()
			return nil, fmt.Errorf("secret is required for long auth")
		}
		secretHash := sha256hash(c.opts.secret)
		if diff, ok := clientData.SecretHash.Compare(secretHash); ok && diff != 0 {
			secHashJson, _ := secretHash.MarshalJSON()
			expectedHashJson, _ := clientData.SecretHash.MarshalJSON()
			logger.Warn("secret hash mismatch",
				zap.String("actual", string(secHashJson)),
				zap.String("expected", string(expectedHashJson)),
			)
			msg, err := abiCocoon.OwnerClientChangeSecretHash{
				QueryId:        tlb.Uint64(rand.Uint64()),
				NewSecretHash:  secretHash,
				SendExcessesTo: c.w.Address.ToInternal(),
			}.ToInternal(clientScAddr.ToInternal(), toncents(20), false, nil)
			if err != nil {
				conn.Close()
				return nil, fmt.Errorf("building change secret hash: %w", err)
			}
			if err := c.w.ForwardMessage(ctx, lc, msg, c.sender); err != nil {
				return nil, fmt.Errorf("forward change secret request failed: %w", err)
			}
			logger.Info("change secret request sent")
			time.Sleep(7 * time.Second)
		}

		if clientData.Balance <= rootStore.Params.Value.MinClientStake {
			topUpAmount := rootStore.Params.Value.MinClientStake - accState.Account.Account.Storage.Balance.Grams + toncents(50)
			logger.Info("amount is low, topping up",
				zap.Uint64("amount", uint64(topUpAmount)),
				zap.Uint64("balance", uint64(accState.Account.Account.Storage.Balance.Grams)),
				zap.Uint64("client_state_balance", uint64(clientData.Balance)),
				zap.Uint64("min_stake", uint64(rootStore.Params.Value.MinClientStake)),
			)
			msg, err := abiCocoon.ExtClientTopUp{
				QueryId:        tlb.Uint64(rand.Uint64()),
				TopUpAmount:    topUpAmount,
				SendExcessesTo: c.w.Address.ToInternal(),
			}.ToInternal(clientScAddr.ToInternal(), topUpAmount+toncents(20), false, nil)
			if err != nil {
				conn.Close()
				return nil, fmt.Errorf("building top-up message failed: %w", err)
			}
			if err := c.w.ForwardMessage(ctx, lc, msg, c.sender); err != nil {
				return nil, fmt.Errorf("forward top-up request failed: %w", err)
			}
			logger.Info("top-up request sent")
			time.Sleep(7 * time.Second)
		}
	}

	// Authorization.
	auth := resp.Auth
	switch auth.SumType {
	case "ClientProxyConnectionAuthShort":
		short := auth.ClientProxyConnectionAuthShort
		logger.Info("auth type: short",
			zap.String("secret_hash", hex.EncodeToString(short.SecretHash[:])))
		authResp, err := apiClient.ClientAuthorizeWithProxyShort(ctx,
			tlcocoonapi.ClientAuthorizeWithProxyShortRequest{Data: []byte(c.opts.secret)})
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
		clientScAddr, err := ton.ParseAccountID(resp.ClientScAddress)
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("parse client SC address: %w", err)
		}

		proxtAddr := ton.MustParseAccountID(resp.Params.ProxyScAddress)

		if err := wallet.SendRegisterTx(ctx, lc, c.w, clientScAddr, long.Nonce, c.sender,
			rootStore.Params.Value.ClientScCode, c.opts.secret, logger, proxtAddr,
			resp.Params.ProxyPublicKey,
			rootStore.Params.Value.MinClientStake,
			rootStore.Params.Value,
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

	return &Connection{
		conn:        conn,
		sess:        sess,
		apiClient:   apiClient,
		rootVersion: rootVersion,
		logger:      logger,
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

// clientAddr extracts the client-facing address from a RegisteredProxy address.
// Format is "worker_addr client_addr" (space-separated) or a single address.
func clientAddr(addr string) string {
	parts := strings.Fields(addr)
	if len(parts) == 2 {
		return parts[1]
	}
	return parts[0]
}

func boolPtr(v bool) *bool       { return &v }
func uint32Ptr(v uint32) *uint32 { return &v }

func sha256hash(data string) tlb.Uint256 {
	secHash := sha256.New()
	secHash.Write([]byte(data))
	secHashBytes := big.NewInt(0).SetBytes(secHash.Sum(nil))
	return tlb.Uint256(*secHashBytes)
}

func toncents(cents uint64) tlb.Coins {
	return tlb.Coins(cents * 1_000_000_000 / 100)
}
