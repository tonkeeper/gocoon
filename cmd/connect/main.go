package main

import (
	"context"
	"crypto/tls"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"strings"

	tongo "github.com/tonkeeper/tongo"
	abiCocoon "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/cocoon"
	"github.com/tonkeeper/tongo/liteapi"

	"github.com/tonkeeper/gococoon/pkg/cocoon"
	"github.com/tonkeeper/gococoon/pkg/tlcocoonapi"
)

const rootContractAddr = "EQCns7bYSp0igFvS1wpb5wsZjCKCV19MD5AVzI4EyxsnU73k"

func main() {
	ownerAddr := flag.String("owner", rootContractAddr, "client owner TON address (any valid address)")
	secretStr := flag.String("secret", "", "secret string for short auth")
	flag.Parse()

	// 1. Connect to TON mainnet and read registered proxies from root contract.
	client, err := liteapi.NewClientWithDefaultMainnet()
	if err != nil {
		log.Fatalf("create liteapi client: %v", err)
	}

	rootAddr := tongo.MustParseAddress(rootContractAddr)
	rootClient := abiCocoon.NewCocoonRoot(client, client).WithAccountId(rootAddr.ID)

	ctx := context.Background()
	rootStore, err := rootClient.Storage(ctx)
	if err != nil {
		log.Fatalf("get root storage: %v", err)
	}

	proxies := rootStore.Data.Value.RegisteredProxies.Values()
	if len(proxies) == 0 {
		log.Fatal("no registered proxies in root contract")
	}

	fmt.Printf("Registered proxies (%d):\n", len(proxies))
	for i, p := range proxies {
		fmt.Printf("  [%d] %s\n", i, p.Address)
	}

	// Address format: "worker_addr client_addr" (space-separated) or single addr.
	proxyAddr := clientAddr(proxies[0].Address)
	fmt.Printf("\nConnecting to proxy: %s\n", proxyAddr)

	// 2. Solve PoW and establish mutual TLS with an ephemeral Ed25519 cert.
	conn, err := cocoon.Dial(proxyAddr)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer conn.Close()

	state := conn.ConnectionState()
	fmt.Printf("Connected!\n")
	fmt.Printf("  TLS version:  %s\n", tlsVersionName(state.Version))
	fmt.Printf("  Cipher suite: %s\n", tls.CipherSuiteName(state.CipherSuite))
	fmt.Printf("  Server cert:  %d bytes (TDX quote embedded)\n", len(conn.ServerCert))

	// 3. Start TcpConnection session (sends tcp.connect framing packet).
	sess, err := cocoon.NewSession(conn)
	if err != nil {
		log.Fatalf("new session: %v", err)
	}

	// 4. Build the API client with the session as transport.
	apiClient := tlcocoonapi.NewClient(sess.Query)

	// 5. Send client.connectToProxy handshake.
	rootVersion := uint32(rootStore.Version)
	handshakeReq := tlcocoonapi.ClientConnectToProxyRequest{
		Params: tlcocoonapi.ClientParamsC{
			// flags=3: bit0=is_test present, bit1=proto versions present
			Flags:              3,
			ClientOwnerAddress: *ownerAddr,
			IsTest:             boolPtr(false),
			MinProtoVersion:    uint32Ptr(1),
			MaxProtoVersion:    uint32Ptr(1),
		},
		MinConfigVersion: rootVersion,
	}

	fmt.Printf("\nSending client.connectToProxy (root version=%d)...\n", rootVersion)
	resp, err := apiClient.ClientConnectToProxy(ctx, handshakeReq)
	if err != nil {
		log.Fatalf("connectToProxy: %v", err)
	}

	fmt.Printf("client.connectedToProxy received!\n")
	fmt.Printf("  Proxy owner:    %s\n", resp.Params.ProxyOwnerAddress)
	fmt.Printf("  Proxy SC:       %s\n", resp.Params.ProxyScAddress)
	fmt.Printf("  Client SC:      %s\n", resp.ClientScAddress)
	fmt.Printf("  Proto version:  %v\n", derefUint32(resp.Params.ProtoVersion))
	pubkey := resp.Params.ProxyPublicKey
	fmt.Printf("  Proxy pubkey:   %s\n", hex.EncodeToString(pubkey[:]))

	// 6. Determine auth type and proceed.
	auth := resp.Auth
	switch auth.SumType {
	case "ClientProxyConnectionAuthShort":
		short := auth.ClientProxyConnectionAuthShort
		fmt.Printf("\nAuth type: short (secret_hash=%s nonce=%d)\n",
			hex.EncodeToString(short.SecretHash[:]), short.Nonce)
		if *secretStr == "" {
			fmt.Println("No --secret provided — skipping short auth.")
			break
		}
		authResp, err := apiClient.ClientAuthorizeWithProxyShort(ctx,
			tlcocoonapi.ClientAuthorizeWithProxyShortRequest{Data: []byte(*secretStr)})
		if err != nil {
			log.Fatalf("authorizeWithProxyShort: %v", err)
		}
		printAuthResult(authResp)

	case "ClientProxyConnectionAuthLong":
		long := auth.ClientProxyConnectionAuthLong
		fmt.Printf("\nAuth type: long (nonce=%d)\n", long.Nonce)
		fmt.Println("Long auth requires a blockchain transaction — skipping for this demo.")
		// In production: send TON tx to proxy SC, then call ClientAuthorizeWithProxyLong.

	default:
		fmt.Printf("\nUnknown auth type: %q\n", auth.SumType)
	}
}

// clientAddr extracts the client-facing address from a RegisteredProxy address
// string. Format is "worker_addr client_addr" or just "addr" if both are same.
func clientAddr(addr string) string {
	parts := strings.Fields(addr)
	if len(parts) == 2 {
		return parts[1]
	}
	return parts[0]
}

func tlsVersionName(v uint16) string {
	switch v {
	case 0x0304:
		return "TLS 1.3"
	case 0x0303:
		return "TLS 1.2"
	default:
		return fmt.Sprintf("0x%04x", v)
	}
}

func printAuthResult(r tlcocoonapi.ClientAuthorizationWithProxy) {
	switch r.SumType {
	case "ClientAuthorizationWithProxySuccess":
		s := r.ClientAuthorizationWithProxySuccess
		fmt.Printf("Auth success! tokens_committed=%d max_tokens=%d\n",
			s.TokensCommittedToDb, s.MaxTokens)
	case "ClientAuthorizationWithProxyFailed":
		f := r.ClientAuthorizationWithProxyFailed
		fmt.Printf("Auth failed: code=%d %s\n", f.ErrorCode, f.Error)
	default:
		fmt.Printf("Unknown auth result: %q\n", r.SumType)
	}
}

func boolPtr(v bool) *bool         { return &v }
func uint32Ptr(v uint32) *uint32   { return &v }
func derefUint32(p *uint32) uint32 { if p == nil { return 0 }; return *p }
