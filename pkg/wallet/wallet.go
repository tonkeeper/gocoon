// Package wallet manages the cocoon wallet: an Ed25519 key pair that identifies
// a cocoon_wallet smart contract on TON.
package wallet

import (
	"crypto/ed25519"
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/ton"
)

// Wallet holds the Ed25519 key pair and associated addresses.
type Wallet struct {
	PrivateKey   ed25519.PrivateKey
	PublicKey    ed25519.PublicKey
	OwnerAddress string        // TON address of the owner
	Address      ton.AccountID // SC address (zero-value until EnsureAddress is called)
	hasAddress   bool
}

// New creates a Wallet from an existing Ed25519 private key and owner address.
func New(priv ed25519.PrivateKey, ownerAddress string) (*Wallet, error) {
	if _, err := ton.ParseAccountID(ownerAddress); err != nil {
		return nil, fmt.Errorf("invalid owner address %q: %w", ownerAddress, err)
	}
	return &Wallet{
		PrivateKey:   priv,
		PublicKey:    priv.Public().(ed25519.PublicKey),
		OwnerAddress: ownerAddress,
	}, nil
}

// SetAddress sets the wallet SC address directly, bypassing StateInit derivation.
func (w *Wallet) SetAddress(addr ton.AccountID) {
	w.Address = addr
	w.hasAddress = true
}

// EnsureAddress derives the cocoon_wallet SC address from the given StateInit
// cell (hash of the cell = account address on workchain 0) and caches it in
// w.Address. Returns immediately if the address is already populated.
func (w *Wallet) EnsureAddress(stateInitCell *boc.Cell) error {
	if w.hasAddress {
		return nil
	}
	hash, err := stateInitCell.Hash()
	if err != nil {
		return fmt.Errorf("hash state init cell: %w", err)
	}
	var addrBits [32]byte
	copy(addrBits[:], hash)
	w.Address = ton.AccountID{Workchain: 0, Address: addrBits}
	w.hasAddress = true
	return nil
}
