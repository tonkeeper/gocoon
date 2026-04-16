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
	PrivateKey   ed25519.PrivateKey // nil when only the SC address is known
	PublicKey    ed25519.PublicKey
	OwnerAddress ton.AccountID // TON address of the owner
	Address      ton.AccountID // wallet SC address
}

// New creates a Wallet for a known, already-deployed wallet SC address.
// priv may be nil when only short auth is needed (no on-chain signing).
func New(addr, ownerAddr ton.AccountID, priv ed25519.PrivateKey) *Wallet {
	w := &Wallet{
		OwnerAddress: ownerAddr,
		Address:      addr,
	}
	if priv != nil {
		w.PrivateKey = priv
		w.PublicKey = priv.Public().(ed25519.PublicKey)
	}
	return w
}

// Generate creates a Wallet from a private key and owner address, deriving the
// wallet SC address from the contract StateInit. Use this when the wallet SC
// has not yet been deployed or the address is not known in advance.
func Generate(priv ed25519.PrivateKey, ownerAddr ton.AccountID) (*Wallet, error) {
	w := &Wallet{
		PrivateKey:   priv,
		PublicKey:    priv.Public().(ed25519.PublicKey),
		OwnerAddress: ownerAddr,
	}
	si, err := BuildStateInit(w)
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
	w.Address = ton.AccountID{Workchain: 0, Address: addrBits}
	return w, nil
}

// EnsureAddress derives the wallet SC address from the given StateInit cell
// and sets w.Address. No-op if Address is already non-zero.
func (w *Wallet) EnsureAddress(stateInitCell *boc.Cell) error {
	if w.Address != (ton.AccountID{}) {
		return nil
	}
	hash, err := stateInitCell.Hash()
	if err != nil {
		return fmt.Errorf("hash state init cell: %w", err)
	}
	var addrBits [32]byte
	copy(addrBits[:], hash)
	w.Address = ton.AccountID{Workchain: 0, Address: addrBits}
	return nil
}
