package tlcocoonapi

import "context"

// Client is the cocoon API client. Generated methods are attached to this type.
// Embed it or compose it; implement request() to wire up the actual transport.
type Client struct {
	transport func(ctx context.Context, payload []byte) ([]byte, error)
}

func NewClient(transport func(ctx context.Context, payload []byte) ([]byte, error)) *Client {
	return &Client{transport: transport}
}

func (c *Client) request(ctx context.Context, payload []byte) ([]byte, error) {
	return c.transport(ctx, payload)
}
