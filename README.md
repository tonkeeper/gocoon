# gocoon

Minimal Go client and local OpenAI-compatible proxy for Cocoon.

## Configuration

To use this client, you need a _cocoon wallet_ with its private key.

The owner address specifies the wallet's owner

The client secret authenticates requests to the proxy. This secret is stored in a cocoon client smart contract that is
automatically deployed by the cocoon wallet.

```bash
export COCOON_WALLET_PRIVKEY=<64-char-hex-seed>
export COCOON_WALLET_OWNER=<wallet-owner-address>
export CLIENT_SECRET=<client-secret>
```

## Usage

Generate wallet key:

```bash
go run ./cmd/gocoontool wallet generate
```

Deploy wallet:

```bash
go run ./cmd/gocoontool wallet deploy
```

Run server:

```bash
go run ./cmd/server
```

Endpoints:

- `GET /v1/models`
- `POST /v1/chat/completions`

Run chat:

```bash
go run ./cmd/gocoochat
```
