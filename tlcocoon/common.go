package tlcocoon

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
)

type Requester interface {
	MakeRequest(ctx context.Context, msg []byte) ([]byte, error)
}

func requestRaw(ctx context.Context, m Requester, in interface {
	CRC() uint32
	MarshalTL() ([]byte, error)
}) ([]byte, error) {
	body, err := in.MarshalTL()
	if err != nil {
		return nil, fmt.Errorf("marshaling: %w", err)
	}
	msg := make([]byte, 4+len(body))
	binary.LittleEndian.PutUint32(msg, in.CRC())
	copy(msg[4:], body)
	respRaw, err := m.MakeRequest(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("sending: %w", err)
	}
	return respRaw, nil
}

func request(ctx context.Context, m Requester, in interface {
	CRC() uint32
	MarshalTL() ([]byte, error)
}, out interface {
	CRC() uint32
	UnmarshalTL(io.Reader) error
}) error {
	respRaw, err := requestRaw(ctx, m, in)
	if err != nil {
		return err
	}
	if len(respRaw) < 4 {
		return fmt.Errorf("response: too short: %d", len(respRaw))
	}
	got := binary.LittleEndian.Uint32(respRaw)
	want := out.CRC()
	if got != want {
		return fmt.Errorf("response: invalid crc code: got 0x%08x; want 0x%08x", got, want)
	}
	err = out.UnmarshalTL(bytes.NewReader(respRaw[4:]))
	if err != nil {
		return fmt.Errorf("response: %w", err)
	}
	return nil
}
