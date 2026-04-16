package client

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"

	"github.com/tonkeeper/tongo/tl"
	"go.uber.org/zap"

	"github.com/tonkeeper/gococoon/pkg/cocoon"
	"github.com/tonkeeper/gococoon/pkg/tlcocoonapi"
)

const (
	defaultMaxCoefficient = 100_000
	defaultMaxTokens      = 200
	defaultTimeout        = 30.0
)

// Connection is a live, authorized connection to a cocoon proxy.
type Connection struct {
	conn        *cocoon.Conn
	sess        *cocoon.Session
	apiClient   *tlcocoonapi.Client
	rootVersion uint32
	logger      *zap.Logger
}

// Close tears down the underlying TLS connection.
func (c *Connection) Close() error {
	return c.conn.Close()
}

// POST sends an HTTP POST request through the cocoon proxy and returns the
// response body bytes. Content-Type: application/json is set automatically.
func (c *Connection) POST(ctx context.Context, model, path string, body []byte) ([]byte, error) {
	query, err := buildHTTPRequest("POST", path, [][2]string{{"Content-Type", "application/json"}}, body)
	if err != nil {
		return nil, err
	}
	return c.runQuery(ctx, model, query)
}

// GET sends an HTTP GET request through the cocoon proxy and returns the
// response body bytes.
func (c *Connection) GET(ctx context.Context, model, path string) ([]byte, error) {
	query, err := buildHTTPRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}
	return c.runQuery(ctx, model, query)
}

// runQuery marshals a RunQueryEx message, sends it, collects streamed answer
// packets, and returns the assembled response body.
func (c *Connection) runQuery(ctx context.Context, model string, queryBytes []byte) ([]byte, error) {
	var reqID tl.Int256
	if _, err := rand.Read(reqID[:]); err != nil {
		return nil, fmt.Errorf("generate request id: %w", err)
	}

	payload, err := tl.Marshal(struct {
		tl.SumType
		Req tlcocoonapi.ClientRunQueryExRequest `tlSumType:"f54cb74b"`
	}{SumType: "Req", Req: tlcocoonapi.ClientRunQueryExRequest{
		ModelName:        model,
		Query:            queryBytes,
		MaxCoefficient:   defaultMaxCoefficient,
		MaxTokens:        defaultMaxTokens,
		Timeout:          defaultTimeout,
		RequestId:        reqID,
		MinConfigVersion: c.rootVersion,
		Flags:            0,
	}})
	if err != nil {
		return nil, fmt.Errorf("marshal runQueryEx: %w", err)
	}
	if err := c.sess.SendMessage(payload); err != nil {
		return nil, fmt.Errorf("send runQueryEx: %w", err)
	}

	var buf []byte
	for {
		pkt, err := c.sess.RecvPacket()
		if err != nil {
			return nil, fmt.Errorf("recv packet: %w", err)
		}
		var ans tlcocoonapi.ClientQueryAnswerEx
		if err := tl.Unmarshal(bytes.NewReader(pkt), &ans); err != nil {
			c.logger.Warn("unrecognised packet", zap.Int("len", len(pkt)), zap.Error(err))
			continue
		}
		switch ans.SumType {
		case "ClientQueryAnswerEx":
			a := ans.ClientQueryAnswerEx
			if a.RequestId != reqID {
				continue
			}
			buf = append(buf, a.Answer...)
			if a.Flags&1 == 1 {
				return buf, nil
			}
		case "ClientQueryAnswerPartEx":
			a := ans.ClientQueryAnswerPartEx
			if a.RequestId != reqID {
				continue
			}
			buf = append(buf, a.Answer...)
			if a.Flags&1 == 1 {
				return buf, nil
			}
		case "ClientQueryAnswerErrorEx":
			a := ans.ClientQueryAnswerErrorEx
			if a.RequestId != reqID {
				continue
			}
			return nil, fmt.Errorf("query error (code %d): %s", a.ErrorCode, a.Error)
		}
	}
}

// buildHTTPRequest serializes a boxed TL http.request.
// Layout: [tag u32][method][url][http_version][vector<http.header>][payload]
func buildHTTPRequest(method, path string, headers [][2]string, body []byte) ([]byte, error) {
	buf := new(bytes.Buffer)

	// boxed http.request tag = 0x47492DE5
	binary.Write(buf, binary.LittleEndian, uint32(0x47492DE5))

	for _, s := range []string{method, path, "HTTP/1.1"} {
		b, err := tl.Marshal(s)
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}

	binary.Write(buf, binary.LittleEndian, uint32(len(headers)))
	for _, h := range headers {
		for _, s := range h {
			b, err := tl.Marshal(s)
			if err != nil {
				return nil, err
			}
			buf.Write(b)
		}
	}

	b, err := tl.Marshal(body)
	if err != nil {
		return nil, err
	}
	buf.Write(b)

	return buf.Bytes(), nil
}
