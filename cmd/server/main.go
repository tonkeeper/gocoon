package main

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/tonkeeper/gocoon"
	"github.com/tonkeeper/gocoon/contracts/wallet_contract"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/ton"
	"go.uber.org/zap"
)

const defaultListenAddr = "127.0.0.1:8080"

type server struct {
	conn   *gocoon.Connection
	logger *zap.Logger
}

type modelObject struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

type modelsResponse struct {
	Object string        `json:"object"`
	Data   []modelObject `json:"data"`
}

type errorBody struct {
	Error struct {
		Message string      `json:"message"`
		Type    string      `json:"type"`
		Param   interface{} `json:"param"`
		Code    interface{} `json:"code"`
	} `json:"error"`
}

type modelRequest struct {
	Model string `json:"model"`
}

func main() {
	logger := zap.Must(zap.NewDevelopment())
	defer logger.Sync() //nolint:errcheck

	ctx := context.Background()
	conn, err := connectCocoon(ctx, logger)
	if err != nil {
		logger.Fatal("connect cocoon", zap.Error(err))
	}
	defer conn.Close()

	s := &server{conn: conn, logger: logger}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleProxy)

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = defaultListenAddr
	}

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info("server listening", zap.String("addr", addr))
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("http server", zap.Error(err))
	}
}

func connectCocoon(ctx context.Context, logger *zap.Logger) (*gocoon.Connection, error) {
	walletOwnerAddr := ton.MustParseAccountID(os.Getenv("COCOON_WALLET_OWNER"))
	clientSecret := os.Getenv("CLIENT_SECRET")
	priv, err := privKeyFromHex(os.Getenv("COCOON_WALLET_PRIVKEY"))
	if err != nil {
		return nil, fmt.Errorf("parse COCOON_WALLET_PRIVKEY: %w", err)
	}

	liteClient, err := liteapi.NewClientWithDefaultMainnet()
	if err != nil {
		return nil, fmt.Errorf("create liteapi client: %w", err)
	}

	wallet, err := wallet_contract.New(priv, walletOwnerAddr, liteClient)
	if err != nil {
		return nil, fmt.Errorf("create wallet: %w", err)
	}
	logger.Info("wallet address", zap.String("address", wallet.Address().ToHuman(false, false)))

	cc, err := gocoon.New(ctx, wallet, clientSecret, gocoon.Opts{LiteClient: liteClient})
	if err != nil {
		return nil, fmt.Errorf("create cocoon client: %w", err)
	}
	conn, err := cc.Connect(ctx, logger)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	return conn, nil
}

func (s *server) handleProxy(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		resp, err := s.conn.GET(r.Context(), "Qwen/Qwen3-32B", r.URL.RequestURI())
		if err != nil {
			s.writeOpenAIError(w, http.StatusBadGateway, "server_error", fmt.Sprintf("upstream error: %v", err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resp)
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.writeOpenAIError(w, http.StatusBadRequest, "invalid_request_error", "failed to read request body")
			return
		}

		model := extractModelFromRequest(r, body)
		if model == "" {
			s.writeOpenAIError(w, http.StatusBadRequest, "invalid_request_error", "model is required (query param `model`, header `X-Model`, or JSON body field `model`)")
			return
		}

		resp, err := s.conn.POST(r.Context(), model, r.URL.RequestURI(), body)
		if err != nil {
			s.writeOpenAIError(w, http.StatusBadGateway, "server_error", fmt.Sprintf("upstream error: %v", err))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resp)
	default:
		s.writeOpenAIError(w, http.StatusMethodNotAllowed, "invalid_request_error", "method not allowed")
	}
}

func extractModelFromRequest(r *http.Request, body []byte) string {
	if q := r.URL.Query().Get("model"); q != "" {
		return q
	}
	if h := r.Header.Get("X-Model"); h != "" {
		return h
	}
	if len(body) == 0 {
		return ""
	}
	var req modelRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return ""
	}
	return req.Model
}

func (s *server) writeOpenAIError(w http.ResponseWriter, status int, typ, msg string) {
	var body errorBody
	body.Error.Message = msg
	body.Error.Type = typ
	writeJSON(w, status, body)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func privKeyFromHex(seed string) (ed25519.PrivateKey, error) {
	raw, err := hex.DecodeString(seed)
	if err != nil || len(raw) != ed25519.SeedSize {
		base64.StdEncoding.DecodeString(seed) //nolint:errcheck
		return nil, fmt.Errorf("must be a 64-char hex string (32-byte Ed25519 seed): %v", err)
	}
	return ed25519.NewKeyFromSeed(raw), nil
}
