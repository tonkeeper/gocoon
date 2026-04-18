# gocoon

Minimal Go client and local OpenAI-compatible proxy for Cocoon.

## Setup wallet to pay Cocoon network

To use Cocoon network, you need a _cocoon wallet_.
It is used to deploy and refill client contracts that carry
authorization secret and serve as payment channel.

1. Generate a new cocoon wallet:
  
      ```bash
      go run github.com/tonkeeper/gocoon/cmd/gocoontool@latest wallet generate
      
      COCOON_WALLET_PRIVKEY=2f690a0cd018bc005b92bd7a5a8a791999c461d49348481e0aae2d35a3dba2c1
      Cocoon wallet address is UQCQxnv9WlZERYgy7AXKY0SNxD1UGxD1UhSx_u-l16yLHloD
      ```

2. Refill generated wallet with 20 TON

3. Deploy wallet:

    ```
   COCOON_WALLET_PRIVKEY=x COCOON_WALLET_OWNER=y go run ./cmd/gocoontool wallet deploy
   ```

## Use coding agent

### Using TONAPI.io

This example runs `aider` coding agent using `uv` python package manager:

```bash
export OPENAI_API_BASE=https://dev.tonapi.io/v2/cocoon/v1
export OPENAI_API_KEY=(API key from tonconsole.com with Cocoon capability enabled)

uvx --from aider-chat --python 3.12 aider --no-gitignore --no-stream --model openai/Qwen/Qwen3-32B --no-show-model-warnings
```

### Using a local proxy with total confidentiality

Client secret is a password to connect to the proxy.

```bash
export COCOON_WALLET_PRIVKEY=<64-char-hex-seed>
export COCOON_WALLET_OWNER=<wallet-owner-address>
export CLIENT_SECRET=<client-secret>

go run github.com/tonkeeper/gocoon/cmd/server@latest
```

You will see this output, OpenAI compatible API is available at `http://localhost:8080/v1/chat/completions`:

```
INFO    server/main.go:104      wallet address  {"address": "UQDSMc2wDpO5kgPwOHO-wQuTc7Fjyvf5OdUWAn_-G0oaqon7"}
INFO    gocoon@v1.0.3/client.go:87      connecting to proxy     {"address": "91.108.4.11:8888"}
DEBUG   proxyconn/connect.go:25 connecting to cocoon proxy      {"addr": "91.108.4.11:8888"}
DEBUG   proxyconn/connect.go:37 solving PoW response
DEBUG   proxyconn/connect.go:41 PoW solved      {"difficulty": 20, "time_spent_ms": 171}
INFO    gocoon@v1.0.3/client.go:140     connected to proxy      {"proxy_owner": "UQDnlslXI2RtI1WhLmtelkb4CVQGxr8E_xSIjl0Hg79jNk6V", "proxy_sc": "UQDTC4XYwcRwan5FEMmfQze48J6YaK8Ao3QiOsvHx83q0r4a", "client_sc": "UQBHk_MRFeU-_ss9vH5r7CfyWY963C3Da9e3uuG6Wu2rWebH", "proxy_pubkey": "1cddc344e4ea001a119b6d27976e7056a29d6cc262d17bec92c8586cd3ab6fce"}
INFO    gocoon@v1.0.3/client.go:161     auth type: short        {"secret_hash": "65b20dea71f753d544a1fe91b3cbd56d53d2cf9432d3ded5d1ec444d3458dc44"}
INFO    gocoon@v1.0.3/client.go:203     auth success    {"tokens_committed": 590619, "max_tokens": 590619}
INFO    server/main.go:81       server listening        {"addr": "127.0.0.1:8080"}
```


### Use simple chat

```bash
go run ./cmd/gocoochat
```
