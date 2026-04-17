// Package contract manages interactions with the cocoon_client smart contract.
package contract

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"
	"math/rand/v2"
	"time"

	abiCocoon "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/cocoon"
	"github.com/tonkeeper/tongo/tl"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"go.uber.org/zap"
)

// LiteClient is the minimal liteapi interface the contract package needs.
type LiteClient interface {
	RunSmcMethodByID(ctx context.Context, accountID ton.AccountID, methodID int, params tlb.VmStack) (uint32, tlb.VmStack, error)
	GetAccountState(ctx context.Context, accountID ton.AccountID) (tlb.ShardAccount, error)
}

// Wallet is the signing interface required by the contract.
type Wallet interface {
	Address() ton.AccountID
	ForwardMessage(ctx context.Context, msg tlb.Message) error
}

// ClientContract manages interactions with a cocoon_client SC on TON.
type ClientContract struct {
	addr   ton.AccountID
	secret string
	wallet Wallet
	lc     LiteClient
}

// New creates a ClientContract for the given client SC address.
func New(addr ton.AccountID, secret string, wallet Wallet, lc LiteClient) *ClientContract {
	return &ClientContract{addr: addr, secret: secret, wallet: wallet, lc: lc}
}

// Sync loads the client SC state. Returns false (not deployed) with no error when the
// account does not exist on-chain yet. When deployed, it sends any needed maintenance
// txs: secret-hash update and balance top-up.
func (c *ClientContract) Sync(ctx context.Context, minStake tlb.Coins, logger *zap.Logger) (deployed bool, err error) {
	sc := abiCocoon.NewCocoonClient(c.lc, c.lc).WithAccountId(c.addr)
	accState, _, err := sc.AccountState(ctx)
	if err != nil {
		return false, fmt.Errorf("get account state: %w", err)
	}
	if accState.Account.SumType != "Account" {
		return false, nil
	}

	clientData, err := sc.GetCocoonClientData(ctx)
	if err != nil {
		return true, fmt.Errorf("get client data: %w", err)
	}
	if err := c.syncSecretHash(ctx, clientData, logger); err != nil {
		return true, err
	}
	if err := c.syncBalance(ctx, clientData, accState, minStake, logger); err != nil {
		return true, err
	}
	return true, nil
}

// Register sends the OwnerClientRegister internal message to the client SC, triggering
// the on-chain registration that long auth waits for.
func (c *ClientContract) Register(
	ctx context.Context, nonce uint64, logger *zap.Logger,
	proxyAddr ton.AccountID, proxyPubKey tl.Int256,
	minStake tlb.Coins, params *abiCocoon.CocoonParams,
) error {
	if !params.ClientScCode.Exists {
		panic("client SC code must exist in cocoon params")
	}

	proxyPubKey256 := tlb.Uint256(*new(big.Int).SetBytes(proxyPubKey[:]))

	paramsCopy := *params
	paramsCopy.ClientScCode.Exists = false
	paramsCopy.WorkerScCode.Exists = false
	paramsCopy.ProxyScCode.Exists = false

	clientInitState := tlb.StateInitT[*abiCocoon.ClientStorage]{
		Code: tlb.JustRef(params.ClientScCode.Value),
		Data: tlb.JustRef(&abiCocoon.ClientStorage{
			Stake:      minStake,
			SecretHash: tlb.Uint256(*big.NewInt(0)),
			ConstDataRef: tlb.RefT[*abiCocoon.ClientConstData]{
				Value: &abiCocoon.ClientConstData{
					OwnerAddress:   c.wallet.Address().ToInternal(),
					ProxyAddress:   proxyAddr.ToInternal(),
					ProxyPublicKey: proxyPubKey256,
				},
			},
			Params: tlb.RefT[*abiCocoon.CocoonParams]{Value: &paramsCopy},
		}),
	}

	cliInStCell, _ := clientInitState.ToCell()
	cliInStCellHash, _ := cliInStCell.HashString()
	fmt.Println("cliInStCellHash", cliInStCellHash)

	registerMsg, err := abiCocoon.OwnerClientRegister{
		QueryId:        0,
		Nonce:          tlb.Uint64(nonce),
		SendExcessesTo: c.wallet.Address().ToInternal(),
	}.ToInternal(
		c.addr.ToInternal(),
		tlb.Grams(50_000_000)+minStake,
		true,
		&clientInitState,
	)
	if err != nil {
		return fmt.Errorf("build OwnerClientRegister: %w", err)
	}
	if err := c.wallet.ForwardMessage(ctx, registerMsg); err != nil {
		return fmt.Errorf("forward register message: %w", err)
	}
	logger.Info("register tx sent")
	return nil
}

func (c *ClientContract) syncSecretHash(ctx context.Context, data abiCocoon.CocoonClientData, logger *zap.Logger) error {
	secretHash := sha256hash(c.secret)
	diff, ok := data.SecretHash.Compare(secretHash)
	if !ok || diff == 0 {
		return nil
	}
	secHashJSON, _ := secretHash.MarshalJSON()
	expectedJSON, _ := data.SecretHash.MarshalJSON()
	logger.Warn("secret hash mismatch",
		zap.String("actual", string(secHashJSON)),
		zap.String("expected", string(expectedJSON)),
	)
	msg, err := abiCocoon.OwnerClientChangeSecretHash{
		QueryId:        tlb.Uint64(rand.Uint64()),
		NewSecretHash:  secretHash,
		SendExcessesTo: c.wallet.Address().ToInternal(),
	}.ToInternal(c.addr.ToInternal(), toncents(20), false, nil)
	if err != nil {
		return fmt.Errorf("build ChangeSecretHash: %w", err)
	}
	if err := c.wallet.ForwardMessage(ctx, msg); err != nil {
		return fmt.Errorf("forward ChangeSecretHash: %w", err)
	}
	logger.Info("change secret request sent")
	time.Sleep(7 * time.Second)
	return nil
}

func (c *ClientContract) syncBalance(ctx context.Context, data abiCocoon.CocoonClientData, accState tlb.ShardAccount, minStake tlb.Coins, logger *zap.Logger) error {
	if data.Balance > minStake {
		return nil
	}
	onChainBalance := accState.Account.Account.Storage.Balance.Grams
	topUpAmount := minStake - onChainBalance + toncents(50)
	logger.Info("balance low, topping up",
		zap.Uint64("amount", uint64(topUpAmount)),
		zap.Uint64("on_chain_balance", uint64(onChainBalance)),
		zap.Uint64("client_balance", uint64(data.Balance)),
		zap.Uint64("min_stake", uint64(minStake)),
	)
	msg, err := abiCocoon.ExtClientTopUp{
		QueryId:        tlb.Uint64(rand.Uint64()),
		TopUpAmount:    topUpAmount,
		SendExcessesTo: c.wallet.Address().ToInternal(),
	}.ToInternal(c.addr.ToInternal(), topUpAmount+toncents(20), false, nil)
	if err != nil {
		return fmt.Errorf("build TopUp: %w", err)
	}
	if err := c.wallet.ForwardMessage(ctx, msg); err != nil {
		return fmt.Errorf("forward TopUp: %w", err)
	}
	logger.Info("top-up request sent")
	time.Sleep(7 * time.Second)
	return nil
}

func sha256hash(data string) tlb.Uint256 {
	h := sha256.New()
	h.Write([]byte(data))
	return tlb.Uint256(*new(big.Int).SetBytes(h.Sum(nil)))
}

func toncents(cents uint64) tlb.Coins {
	return tlb.Coins(cents * 1_000_000_000 / 100)
}
