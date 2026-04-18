package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/tonkeeper/gocoon/contracts/wallet_contract"
	abiCocoon "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/cocoon"
	"github.com/tonkeeper/tongo/ton"
	"go.uber.org/zap"
)

func cmdWalletGenerate() {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stderr, "generated new key — set COCOON_WALLET_PRIVKEY=%s to reuse it\n",
		hex.EncodeToString(priv.Seed()))
}

func cmdWalletDeploy() {
	logger := zap.Must(zap.NewDevelopment())
	defer logger.Sync() //nolint:errcheck
	ctx := context.Background()

	priv, err := privKeyFromHex(os.Getenv("COCOON_WALLET_PRIVKEY"))
	if err != nil {
		logger.Fatal("parse COCOON_WALLET_PRIVKEY", zap.Error(err))
	}
	ownerAddr := ton.MustParseAccountID(os.Getenv("COCOON_WALLET_OWNER"))

	lc := mustLiteClient()
	w, err := wallet_contract.New(priv, ownerAddr, lc)
	if err != nil {
		logger.Fatal("create wallet", zap.Error(err))
	}
	logger.Info("wallet address", zap.String("address", w.Address().ToHuman(false, false)))

	if err := w.Deploy(ctx); err != nil {
		logger.Fatal("deploy wallet", zap.Error(err))
	}
	logger.Info("deploy tx sent — wallet will appear on-chain after the next block")
}

func cmdWalletState() {
	addr := mustArg(3, "wallet state <address>")
	ctx := context.Background()
	lc := mustLiteClient()

	wc := abiCocoon.NewCocoonWallet(lc, lc).WithAccountId(ton.MustParseAccountID(addr))

	seqno, err := wc.GetSeqno(ctx)
	if err != nil {
		fatalf("get seqno: %v", err)
	}
	pubkey, err := wc.GetPublicKey(ctx)
	if err != nil {
		fatalf("get public key: %v", err)
	}
	ownerAddr, err := wc.GetOwnerAddress(ctx)
	if err != nil {
		fatalf("get owner address: %v", err)
	}
	fmt.Printf("WalletState (%s):\n", addr)
	fmt.Printf("  Seqno:        %v\n", seqno)
	fmt.Printf("  PublicKey:    %v\n", pubkey)
	fmt.Printf("  OwnerAddress: %v\n", ownerAddr)
}

func privKeyFromHex(seed string) (ed25519.PrivateKey, error) {
	raw, err := hex.DecodeString(seed)
	if err != nil || len(raw) != ed25519.SeedSize {
		return nil, fmt.Errorf("must be a 64-char hex string (32-byte Ed25519 seed): %v", err)
	}
	return ed25519.NewKeyFromSeed(raw), nil
}
