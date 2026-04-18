# gocoon

Minimal Go client and local OpenAI-compatible proxy for Cocoon.

## Use coding agent

### Using TONAPI.io

This example runs `aider` coding agent using `uv` python package manager:

```bash
export OPENAI_API_BASE=https://dev.tonapi.io/v2/cocoon/v1
export OPENAI_API_KEY=(API key from tonconsole.com with Cocoon capability enabled)

uvx --from aider-chat --python 3.12 aider --no-gitignore --no-stream --model openai/Qwen/Qwen3-32B --no-show-model-warnings
```

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
