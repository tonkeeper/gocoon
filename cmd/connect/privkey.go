package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
)

func privKeyFromHex(seed string) (ed25519.PrivateKey, error) {
	raw, err := hex.DecodeString(seed)
	if err != nil || len(raw) != ed25519.SeedSize {
		return nil, fmt.Errorf("must be a 64-char hex string (32-byte Ed25519 seed): %v", err)
	}
	return ed25519.NewKeyFromSeed(raw), nil
}
