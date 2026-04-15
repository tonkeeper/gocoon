package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"

	tongo "github.com/tonkeeper/tongo"
	abiCocoon "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/cocoon"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tl"

	"github.com/tonkeeper/gococoon/pkg/cocoon"
	"github.com/tonkeeper/gococoon/pkg/tlcocoonapi"
)

const rootContractAddr = "EQCns7bYSp0igFvS1wpb5wsZjCKCV19MD5AVzI4EyxsnU73k"

func main() {
	ownerAddr := flag.String("owner", "UQD8tvSwKYOC0TQi-H2tgry_JlQbkPl_EARYnA4ejZgaTqI9", "client owner TON address (any valid address)")
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

	// 6. Authorize the connection before sending any client requests.
	auth := resp.Auth
	switch auth.SumType {
	case "ClientProxyConnectionAuthShort":
		short := auth.ClientProxyConnectionAuthShort
		fmt.Printf("\nAuth type: short (secret_hash=%s nonce=%d)\n",
			hex.EncodeToString(short.SecretHash[:]), short.Nonce)
		authResp, err := apiClient.ClientAuthorizeWithProxyShort(ctx,
			tlcocoonapi.ClientAuthorizeWithProxyShortRequest{Data: []byte(*secretStr)})
		if err != nil {
			log.Fatalf("authorizeWithProxyShort: %v", err)
		}
		if !printAuthResult(authResp) {
			log.Fatal("authorization failed — cannot continue")
		}

	case "ClientProxyConnectionAuthLong":
		long := auth.ClientProxyConnectionAuthLong
		fmt.Printf("\nAuth type: long (nonce=%d)\n", long.Nonce)
		fmt.Println("Long auth requires a blockchain transaction — skipping for this demo.")
		// In production: send TON tx to proxy SC, then call ClientAuthorizeWithProxyLong.
		return

	default:
		log.Fatalf("unknown auth type: %q", auth.SumType)
	}

	// 7. Fetch available models from the proxy.
	protoVersion := derefUint32(resp.Params.ProtoVersion)
	fmt.Printf("\nFetching models (proto_version=%d)...\n", protoVersion)
	if protoVersion != 0 {
		models, err := apiClient.ClientGetWorkerTypesV2(ctx)
		if err != nil {
			log.Fatalf("getWorkerTypesV2: %v", err)
		}
		fmt.Printf("Models (%d):\n", len(models.Types))
		for _, wt := range models.Types {
			fmt.Printf("  %s (%d workers)\n", wt.Name, len(wt.Workers))
			for _, w := range wt.Workers {
				fmt.Printf("    coefficient=%d active=%d/%d\n",
					w.Coefficient, w.ActiveRequests, w.MaxActiveRequests)
			}
		}
	} else {
		models, err := apiClient.ClientGetWorkerTypes(ctx)
		if err != nil {
			log.Fatalf("getWorkerTypes: %v", err)
		}
		fmt.Printf("Models (%d):\n", len(models.Types))
		for _, wt := range models.Types {
			fmt.Printf("  %s active=%d coeff=[%d..%d] p50=%d\n",
				wt.Name, wt.ActiveWorkers,
				wt.CoefficientMin, wt.CoefficientMax, wt.CoefficientBucket50)
		}
	}

	// 8. Send a test query as a message (not tcp.query — proxy routes runQueryEx
	//    via receive_message, not receive_query).
	const testModel = "Qwen/Qwen3-32B"
	fmt.Printf("\nRunning query model=%s prompt=\"1+1=?\"...\n", testModel)

	bodyJSON, err := json.Marshal(map[string]any{
		"model":      testModel,
		"messages":   []map[string]string{{"role": "user", "content": "1+1=?"}},
		"max_tokens": 200,
	})
	if err != nil {
		log.Fatalf("marshal query body: %v", err)
	}
	queryBytes, err := buildQueryBytes(bodyJSON)
	if err != nil {
		log.Fatalf("build query bytes: %v", err)
	}
	reqID, err := randInt256()
	if err != nil {
		log.Fatalf("rand request id: %v", err)
	}

	msgPayload, err := tl.Marshal(struct {
		tl.SumType
		Req tlcocoonapi.ClientRunQueryExRequest `tlSumType:"f54cb74b"`
	}{SumType: "Req", Req: tlcocoonapi.ClientRunQueryExRequest{
		ModelName:        testModel,
		Query:            queryBytes,
		MaxCoefficient:   100000,
		MaxTokens:        200,
		Timeout:          30.0,
		RequestId:        reqID,
		MinConfigVersion: rootVersion,
		Flags:            0,
	}})
	if err != nil {
		log.Fatalf("marshal runQueryEx: %v", err)
	}
	if err := sess.SendMessage(msgPayload); err != nil {
		log.Fatalf("send runQueryEx: %v", err)
	}

	// Receive answer packets, correlating by request_id.
	var answerBuf []byte
loop:
	for {
		pkt, err := sess.RecvPacket()
		if err != nil {
			log.Fatalf("recv packet: %v", err)
		}
		var ans tlcocoonapi.ClientQueryAnswerEx
		if err := tl.Unmarshal(bytes.NewReader(pkt), &ans); err != nil {
			log.Printf("unrecognised packet (len=%d), skipping: %v", len(pkt), err)
			continue
		}
		switch ans.SumType {
		case "ClientQueryAnswerEx":
			if ans.ClientQueryAnswerEx.RequestId != reqID {
				continue
			}
			answerBuf = append(answerBuf, ans.ClientQueryAnswerEx.Answer...)
			if ans.ClientQueryAnswerEx.Flags&1 == 1 {
				break loop
			}
		case "ClientQueryAnswerPartEx":
			if ans.ClientQueryAnswerPartEx.RequestId != reqID {
				continue
			}
			answerBuf = append(answerBuf, ans.ClientQueryAnswerPartEx.Answer...)
			if ans.ClientQueryAnswerPartEx.Flags&1 == 1 {
				break loop
			}
		case "ClientQueryAnswerErrorEx":
			if ans.ClientQueryAnswerErrorEx.RequestId != reqID {
				continue
			}
			log.Fatalf("query error code=%d: %s", ans.ClientQueryAnswerErrorEx.ErrorCode, ans.ClientQueryAnswerErrorEx.Error)
		}
	}
	fmt.Printf("Answer:\n%s\n", answerBuf)
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

// printAuthResult prints the auth result and returns true on success.
func printAuthResult(r tlcocoonapi.ClientAuthorizationWithProxy) bool {
	switch r.SumType {
	case "ClientAuthorizationWithProxySuccess":
		s := r.ClientAuthorizationWithProxySuccess
		fmt.Printf("Auth success! tokens_committed=%d max_tokens=%d\n",
			s.TokensCommittedToDb, s.MaxTokens)
		return true
	case "ClientAuthorizationWithProxyFailed":
		f := r.ClientAuthorizationWithProxyFailed
		fmt.Printf("Auth failed: code=%d %s\n", f.ErrorCode, f.Error)
		return false
	default:
		fmt.Printf("Unknown auth result: %q\n", r.SumType)
		return false
	}
}

func boolPtr(v bool) *bool       { return &v }
func uint32Ptr(v uint32) *uint32 { return &v }
func derefUint32(p *uint32) uint32 {
	if p == nil {
		return 0
	}
	return *p
}

// buildQueryBytes returns a boxed TL-serialized http.request for POST /v1/chat/completions.
// Layout: [tag u32][method][url][http_version][vector<http.header>][payload]
// http.header elements are bare (no per-element tag) as per cocoon TL codegen.
func buildQueryBytes(body []byte) ([]byte, error) {
	buf := new(bytes.Buffer)

	// boxed http.request tag = 1195978213 (0x47492DE5)
	binary.Write(buf, binary.LittleEndian, uint32(1195978213))

	for _, s := range []string{"POST", "/v1/chat/completions", "HTTP/1.1"} {
		b, err := tl.Marshal(s)
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}

	// vector<http.header>: count + bare elements (name, value).
	// TON TL vectors have no magic prefix — just count followed by elements.
	binary.Write(buf, binary.LittleEndian, uint32(1)) // 1 header
	for _, s := range []string{"Content-Type", "application/json"} {
		b, err := tl.Marshal(s)
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}

	// payload bytes
	b, err := tl.Marshal(body)
	if err != nil {
		return nil, err
	}
	buf.Write(b)

	return buf.Bytes(), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func randInt256() (tl.Int256, error) {
	var id tl.Int256
	_, err := rand.Read(id[:])
	return id, err
}

// httpResponsePayload parses a boxed TL http.response and returns its payload bytes.
// Layout: [tag u32][http_version str][status_code i32][reason str][headers vec][payload bytes]
// TON TL vectors: [count u32][bare elements...] — no magic prefix.
func httpResponsePayload(data []byte) ([]byte, error) {
	r := bytes.NewReader(data)

	var tag uint32
	if err := binary.Read(r, binary.LittleEndian, &tag); err != nil {
		return nil, fmt.Errorf("read tag: %w", err)
	}
	if tag != 0x1cd4212b { // http.response ID = 483443755
		return nil, fmt.Errorf("unexpected tag 0x%08x", tag)
	}

	// http_version
	if _, err := tlReadString(r); err != nil {
		return nil, fmt.Errorf("http_version: %w", err)
	}
	// status_code
	if err := binary.Read(r, binary.LittleEndian, new(int32)); err != nil {
		return nil, fmt.Errorf("status_code: %w", err)
	}
	// reason
	if _, err := tlReadString(r); err != nil {
		return nil, fmt.Errorf("reason: %w", err)
	}
	// headers: count + bare http.header elements (name, value)
	var count uint32
	if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
		return nil, fmt.Errorf("headers count: %w", err)
	}
	for i := uint32(0); i < count; i++ {
		if _, err := tlReadString(r); err != nil {
			return nil, fmt.Errorf("header name: %w", err)
		}
		if _, err := tlReadString(r); err != nil {
			return nil, fmt.Errorf("header value: %w", err)
		}
	}
	// payload
	return tlReadBytes(r)
}

func tlReadString(r *bytes.Reader) (string, error) {
	b, err := tlReadBytes(r)
	return string(b), err
}

func tlReadBytes(r *bytes.Reader) ([]byte, error) {
	first, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	var length, pad int
	if first < 254 {
		length = int(first)
		pad = (4 - (1+length)%4) % 4
	} else if first == 254 {
		var b [3]byte
		if _, err := io.ReadFull(r, b[:]); err != nil {
			return nil, err
		}
		length = int(b[0]) | int(b[1])<<8 | int(b[2])<<16
		pad = (4 - length%4) % 4
	} else {
		return nil, fmt.Errorf("unsupported TL bytes prefix 0x%02x", first)
	}
	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}
	if pad > 0 {
		if _, err := io.ReadFull(r, make([]byte, pad)); err != nil {
			return nil, err
		}
	}
	return data, nil
}
