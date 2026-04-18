package gocoon

import (
	"context"
	"fmt"

	"github.com/tonkeeper/gocoon/internal/session"
	"github.com/tonkeeper/gocoon/net/proxyconn"
	"github.com/tonkeeper/gocoon/tlcocoon"
	tlcocoonTypes "github.com/tonkeeper/gocoon/tlcocoon/types"
	"go.uber.org/zap"
)

const (
	defaultMaxCoefficient = 100_000
	defaultMaxTokens      = 200
	defaultTimeout        = 30.0
)

// WorkerInstance describes a single inference worker.
type WorkerInstance struct {
	Coefficient       uint32
	ActiveRequests    uint32
	MaxActiveRequests uint32
}

// WorkerType describes a model type and its available workers.
type WorkerType struct {
	Name    string
	Workers []WorkerInstance
}

// Connection is a live, authorized connection to a cocoon proxy.
type Connection struct {
	conn         *proxyconn.Conn
	sess         *session.Session
	apiClient    *tlcocoon.Client
	rootVersion  uint32
	protoVersion uint32
	logger       *zap.Logger
}

// Close tears down the session and the underlying TLS connection.
func (c *Connection) Close() error {
	_ = c.sess.Close()
	return c.conn.Close()
}

// POST sends an HTTP POST request through the cocoon proxy and returns the response body bytes.
func (c *Connection) POST(ctx context.Context, model, path string, body []byte) ([]byte, error) {
	query, err := tlcocoon.EncodeHTTPRequest(
		"POST",
		path,
		"HTTP/1.1",
		[]tlcocoonTypes.HttpHeader{{Name: "Content-Type", Value: "application/json"}},
		body,
	)
	if err != nil {
		return nil, err
	}
	raw, err := c.runQuery(ctx, model, query)
	if err != nil {
		return nil, err
	}
	return tlcocoon.DecodeHTTPResponsePayload(raw)
}

// GET sends an HTTP GET request through the cocoon proxy and returns the response body bytes.
func (c *Connection) GET(ctx context.Context, model, path string) ([]byte, error) {
	query, err := tlcocoon.EncodeHTTPRequest("GET", path, "HTTP/1.1", nil, nil)
	if err != nil {
		return nil, err
	}
	raw, err := c.runQuery(ctx, model, query)
	if err != nil {
		return nil, err
	}
	return tlcocoon.DecodeHTTPResponsePayload(raw)
}

// GetWorkerTypes returns the list of worker types available on the proxy,
// normalized to a unified type regardless of the negotiated proto version.
func (c *Connection) GetWorkerTypes(ctx context.Context) ([]WorkerType, error) {
	if c.protoVersion != 0 {
		res, err := c.apiClient.GetWorkerTypesV2(ctx, tlcocoon.ClientGetWorkerTypesV2Request{})
		if err != nil {
			return nil, fmt.Errorf("getWorkerTypesV2: %w", err)
		}
		out := make([]WorkerType, len(res.Types))
		for i, wt := range res.Types {
			workers := make([]WorkerInstance, len(wt.Workers))
			for j, w := range wt.Workers {
				workers[j] = WorkerInstance{
					Coefficient:       uint32(w.Coefficient),
					ActiveRequests:    uint32(w.ActiveRequests),
					MaxActiveRequests: uint32(w.MaxActiveRequests),
				}
			}
			out[i] = WorkerType{Name: wt.Name, Workers: workers}
		}
		return out, nil
	}

	res, err := c.apiClient.GetWorkerTypes(ctx, tlcocoon.ClientGetWorkerTypesRequest{})
	if err != nil {
		return nil, fmt.Errorf("getWorkerTypes: %w", err)
	}
	out := make([]WorkerType, len(res.Types))
	for i, wt := range res.Types {
		out[i] = WorkerType{
			Name: wt.Name,
			Workers: []WorkerInstance{{
				Coefficient:       uint32(wt.CoefficientBucket50),
				ActiveRequests:    uint32(wt.ActiveWorkers),
				MaxActiveRequests: 0,
			}},
		}
	}
	return out, nil
}

// runQuery marshals a RunQueryEx message, sends it, collects streamed answer
// packets, and returns the assembled response body.
func (c *Connection) runQuery(ctx context.Context, model string, queryBytes []byte) ([]byte, error) {
	raw, err := c.sess.RunClientQueryEx(ctx, model, queryBytes, session.RunQueryExOptions{
		MaxCoefficient:   int32(defaultMaxCoefficient),
		MaxTokens:        int32(defaultMaxTokens),
		Timeout:          defaultTimeout,
		MinConfigVersion: int32(c.rootVersion),
	})
	if err != nil {
		return nil, fmt.Errorf("runQueryEx: %w", err)
	}
	return raw, nil
}
