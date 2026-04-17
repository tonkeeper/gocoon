package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/ton"
	"go.uber.org/zap"

	goclient "github.com/tonkeeper/gocoon/pkg/client"
	cocoonWallet "github.com/tonkeeper/gocoon/pkg/wallet"
)

const (
	chatModel    = "Qwen/Qwen3-32B"
	maxTokens    = 1200
	systemPrompt = "Disable thinking"
)

func main() {
	logger := zap.Must(zap.NewDevelopment())
	defer logger.Sync() //nolint:errcheck
	ctx := context.Background()

	walletOwnerAddr := ton.MustParseAccountID(os.Getenv("OWNER_ADDRESS"))
	clientSecret := os.Getenv("CLIENT_SECRET")
	priv, err := privKeyFromHex(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		logger.Fatal("parse PRIVATE_KEY", zap.Error(err))
	}

	liteClient, err := liteapi.NewClientWithDefaultMainnet()
	if err != nil {
		logger.Fatal("create liteapi client", zap.Error(err))
	}

	wallet, err := cocoonWallet.New(priv, walletOwnerAddr, liteClient)
	if err != nil {
		logger.Fatal("create wallet", zap.Error(err))
	}
	logger.Info("wallet address", zap.String("address", wallet.Address().ToHuman(false, false)))

	cc, err := goclient.NewCocoonClient(ctx, wallet, clientSecret, goclient.Opts{LiteClient: liteClient})
	if err != nil {
		logger.Fatal("create cocoon client", zap.Error(err))
	}
	conn, err := cc.Connect(ctx, logger)
	if err != nil {
		logger.Fatal("connect", zap.Error(err))
	}
	defer conn.Close()

	t0 := time.Now()
	workerTypes, err := conn.GetWorkerTypes(ctx)
	if err != nil {
		logger.Fatal("GetWorkerTypes", zap.Error(err))
	}
	logger.Info("worker types fetched", zap.Duration("elapsed", time.Since(t0)))
	for _, wt := range workerTypes {
		fmt.Printf("  %s (%d workers)\n", wt.Name, len(wt.Workers))
		for _, w := range wt.Workers {
			fmt.Printf("    coefficient=%d active=%d/%d\n", w.Coefficient, w.ActiveRequests, w.MaxActiveRequests)
		}
	}

	type message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	history := []message{
		{Role: "system", Content: systemPrompt},
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("\nChat with %s (Ctrl+D to quit)\n\n", chatModel)
	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}
		userText := strings.TrimSpace(scanner.Text())
		if userText == "" {
			continue
		}
		history = append(history, message{Role: "user", Content: userText})

		bodyJSON, err := json.Marshal(map[string]any{
			"model":      chatModel,
			"messages":   history,
			"max_tokens": maxTokens,
		})
		if err != nil {
			logger.Fatal("marshal request", zap.Error(err))
		}

		t0 := time.Now()
		resp, err := conn.POST(ctx, chatModel, "/v1/chat/completions", bodyJSON)
		if err != nil {
			logger.Fatal("POST", zap.Error(err))
		}
		elapsed := time.Since(t0)

		idx := bytes.Index(resp, []byte("{"))
		if idx < 0 {
			logger.Fatal("JSON not found in response", zap.ByteString("resp", resp))
		}
		var completion struct {
			Choices []struct {
				Message message `json:"message"`
			} `json:"choices"`
		}
		if err := json.Unmarshal(resp[idx:], &completion); err != nil {
			logger.Fatal("parse completion JSON", zap.Error(err))
		}
		if len(completion.Choices) == 0 {
			logger.Warn("empty choices in response")
			continue
		}

		reply := completion.Choices[0].Message
		history = append(history, reply)
		fmt.Printf("\nAssistant (%s): %s\n\n", elapsed.Round(time.Millisecond), reply.Content)
	}
}

func privKeyFromHex(seed string) (ed25519.PrivateKey, error) {
	raw, err := hex.DecodeString(seed)
	if err != nil || len(raw) != ed25519.SeedSize {
		base64.StdEncoding.DecodeString(seed) //nolint:errcheck
		return nil, fmt.Errorf("must be a 64-char hex string (32-byte Ed25519 seed): %v", err)
	}
	return ed25519.NewKeyFromSeed(raw), nil
}
