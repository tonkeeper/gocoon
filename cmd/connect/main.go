package main

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/ton"
	"go.uber.org/zap"

	goclient "github.com/tonkeeper/gocoon/pkg/client"
	cocoonWallet "github.com/tonkeeper/gocoon/pkg/wallet"
)

func main() {
	logger := zap.Must(zap.NewDevelopment())
	defer logger.Sync() //nolint:errcheck
	var err error
	ctx := context.Background()

	walletOwnerAddr := ton.MustParseAccountID(os.Getenv("OWNER_ADDRESS"))
	clientSecret := os.Getenv("CLIENT_SECRET")
	var priv ed25519.PrivateKey
	priv, err = privKeyFromHex(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		logger.Fatal("parse PRIVATE_KEY", zap.Error(err))
	}

	lc, err := liteapi.NewClientWithDefaultMainnet()
	if err != nil {
		logger.Fatal("create liteapi client", zap.Error(err))
	}

	w, err := cocoonWallet.New(priv, walletOwnerAddr, lc)
	if err != nil {
		logger.Fatal("create wallet", zap.Error(err))
	}
	logger.Info("wallet address", zap.String("address", w.Address().ToHuman(false, false)))

	cc := goclient.NewCocoonClient(w, goclient.Opts{}.WithSecret(clientSecret))
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
