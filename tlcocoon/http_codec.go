package tlcocoon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	types "github.com/tonkeeper/gocoon/tlcocoon/types"
)

// EncodeHTTPRequest builds boxed TL http.request payload.
func EncodeHTTPRequest(method, url, httpVersion string, headers []types.HttpHeader, payload []byte) ([]byte, error) {
	req := HttpRequestRequest{
		Method:      method,
		URL:         url,
		HttpVersion: httpVersion,
		Headers:     headers,
		Payload:     payload,
	}
	raw, err := req.MarshalTL()
	if err != nil {
		return nil, err
	}
	out := make([]byte, 4+len(raw))
	binary.LittleEndian.PutUint32(out[:4], req.CRC())
	copy(out[4:], raw)
	return out, nil
}

// DecodeHTTPResponsePayload decodes boxed TL http.response and returns body payload.
func DecodeHTTPResponsePayload(raw []byte) ([]byte, error) {
	if len(raw) < 4 {
		return nil, fmt.Errorf("http response too short: %d", len(raw))
	}
	if binary.LittleEndian.Uint32(raw[:4]) != (&types.HttpResponse{}).CRC() {
		return nil, fmt.Errorf("unexpected response crc: 0x%08x", binary.LittleEndian.Uint32(raw[:4]))
	}
	var resp types.HttpResponse
	r := bytes.NewReader(raw[4:])
	if err := resp.UnmarshalTL(r); err != nil {
		return nil, err
	}
	if len(resp.Payload) > 0 {
		return resp.Payload, nil
	}
	// Some responses may carry chunk bytes after parsed fields.
	if tail, err := io.ReadAll(r); err == nil && len(tail) > 0 {
		return tail, nil
	}
	return resp.Payload, nil
}
