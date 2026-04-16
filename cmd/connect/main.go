package main

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	tonapi "github.com/tonkeeper/tonapi-go"
	"github.com/tonkeeper/tongo/ton"
	"go.uber.org/zap"

	goclient "github.com/tonkeeper/gococoon/pkg/client"
	"github.com/tonkeeper/gococoon/pkg/wallet"
)

func main() {
	secretStr := flag.String("secret", "", "secret string for short auth")
	tonapiKey := flag.String("tonapi-key", "", "tonapi.io API key (optional, increases rate limits)")
	flag.Parse()

	logger := zap.Must(zap.NewDevelopment())
	defer logger.Sync() //nolint:errcheck

	ctx := context.Background()

	ownerAddrStr := os.Getenv("OWNER_ADDRESS")
	if ownerAddrStr == "" {
		ownerAddrStr = "UQCFsVc8MnJ6PwU6Li8EFX1vplzc0f4GREt9UsyKNSdxRph4"
		//ownerAddrStr = "EQCB2xmOW5UsLnnOioX_FXDhAuFxKDtCF5G1BIJejLbAAWOs"
	}
	ownerAddr := ton.MustParseAccountID(ownerAddrStr)

	var priv ed25519.PrivateKey
	if seed := os.Getenv("PRIVATE_KEY"); seed != "" {
		raw, err := hex.DecodeString(seed)
		if err != nil || len(raw) != ed25519.SeedSize {
			logger.Fatal("PRIVATE_KEY must be a 64-char hex string (32-byte Ed25519 seed)")
		}
		priv = ed25519.NewKeyFromSeed(raw)
	} else {
		_, generated, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			logger.Fatal("generate ed25519 key", zap.Error(err))
		}
		priv = generated
		fmt.Fprintf(os.Stderr, "generated new key — set PRIVATE_KEY=%s to reuse it\n",
			hex.EncodeToString(priv.Seed()))
	}

	//walletAddrStr := os.Getenv("WALLET_ADDRESS")
	//if walletAddrStr == "" {
	//	walletAddrStr = "UQD_5KYZHQcUIhBJhPQ0n3Fpg7l2qqE6Wc5W2tMeAypMfm0C"
	//}
	//walletAddr, err := ton.ParseAccountID(walletAddrStr)
	//if err != nil {
	//	logger.Fatal("invalid WALLET_ADDRESS", zap.Error(err))
	//}

	w, err := wallet.Generate(priv, ownerAddr)
	if err != nil {
		logger.Fatal("generate wallet", zap.Error(err))
	}
	logger.Info("wallet address", zap.String("address", w.Address.ToHuman(false, false)))

	sender, err := tonapi.NewClient(tonapi.TonApiURL, tonapi.WithToken(*tonapiKey))
	if err != nil {
		logger.Fatal("create tonapi client", zap.Error(err))
	}

	cc := goclient.NewCocoonClient(w, sender, goclient.Opts{}.WithSecret(*secretStr))
	conn, err := cc.Connect(ctx, logger)
	if err != nil {
		logger.Fatal("connect", zap.Error(err))
	}
	defer conn.Close()

	const testModel = "Qwen/Qwen3-32B"
	bodyJSON, err := json.Marshal(map[string]any{
		"model":      testModel,
		"messages":   []map[string]string{{"role": "user", "content": "Tell me latest news about TON"}},
		"max_tokens": 1200,
	})
	if err != nil {
		logger.Fatal("marshal query body", zap.Error(err))
	}

	logger.Info("running query", zap.String("model", testModel))
	resp, err := conn.POST(ctx, testModel, "/v1/chat/completions", bodyJSON)
	if err != nil {
		logger.Fatal("POST", zap.Error(err))
	}

	idx := bytes.Index(resp, []byte("{"))
	if idx < 0 {
		logger.Fatal("JSON not found in response", zap.ByteString("resp", resp))
	}
	var completion struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(resp[idx:], &completion); err != nil {
		logger.Fatal("parse completion JSON", zap.Error(err))
	}
	if len(completion.Choices) > 0 {
		fmt.Println(completion.Choices[0].Message.Content)
	}
}
