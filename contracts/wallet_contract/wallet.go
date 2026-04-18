// Package cocoonWallet manages the cocoon wallet: an Ed25519 key pair that identifies
// a cocoon_wallet smart contract on TON.
package wallet_contract

import (
	"context"
	"crypto/ed25519"
	"fmt"

	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type BlockchainClient interface {
	RunSmcMethodByID(ctx context.Context, accountID ton.AccountID, methodID int, params tlb.VmStack) (uint32, tlb.VmStack, error)
	SendMessage(ctx context.Context, payload []byte) (uint32, error)
}

type Wallet struct {
	privateKey   ed25519.PrivateKey
	publicKey    ed25519.PublicKey
	ownerAddress ton.AccountID
	address      ton.AccountID
	lc           BlockchainClient
}

// New creates a Wallet from a private key and owner address, deriving the
// wallet SC address from the contract StateInit.
func New(priv ed25519.PrivateKey, ownerAddr ton.AccountID, lc BlockchainClient) (*Wallet, error) {
	w := &Wallet{
		privateKey:   priv,
		publicKey:    priv.Public().(ed25519.PublicKey),
		ownerAddress: ownerAddr,
		lc:           lc,
	}
	si, err := w.buildStateInit()
	if err != nil {
		return nil, fmt.Errorf("build state init: %w", err)
	}
	siCell, err := si.ToCell()
	if err != nil {
		return nil, fmt.Errorf("serialize state init: %w", err)
	}
	hash, err := siCell.Hash()
	if err != nil {
		return nil, fmt.Errorf("hash state init: %w", err)
	}
	var addrBits [32]byte
	copy(addrBits[:], hash)
	w.address = ton.AccountID{Workchain: 0, Address: addrBits}
	return w, nil
}

// Address returns the wallet SC address on TON.
func (w *Wallet) Address() ton.AccountID {
	return w.address
}
