package tlcocoonTypes

import (
	"bytes"
	"fmt"
	tl "github.com/tonkeeper/tongo/tl"
	"io"
)

type ITcpPacket interface {
	CRC() uint32
	MarshalTL() ([]byte, error)
	UnmarshalTL(io.Reader) error
	_ITcpPacket()
}

var (
	_ ITcpPacket = (*TcpPing)(nil)
	_ ITcpPacket = (*TcpPong)(nil)
	_ ITcpPacket = (*TcpPacket)(nil)
	_ ITcpPacket = (*TcpQueryAnswer)(nil)
	_ ITcpPacket = (*TcpQueryError)(nil)
	_ ITcpPacket = (*TcpQuery)(nil)
	_ ITcpPacket = (*TcpConnected)(nil)
	_ ITcpPacket = (*TcpConnect)(nil)
)

func decodeITcpPacket(r io.Reader) (ITcpPacket, error) {
	var tag uint32
	err := tl.Unmarshal(r, &tag)
	if err != nil {
		return nil, err
	}
	var res ITcpPacket
	switch tag {
	case uint32(0xbbe9627c):
		res = &TcpPing{}
	case uint32(0xbd4302c):
		res = &TcpPong{}
	case uint32(0x9c17baf8):
		res = &TcpPacket{}
	case uint32(0xc048c311):
		res = &TcpQueryAnswer{}
	case uint32(0x4cd2f602):
		res = &TcpQueryError{}
	case uint32(0x3af51908):
		res = &TcpQuery{}
	case uint32(0x636d41d6):
		res = &TcpConnected{}
	case uint32(0xa57c4261):
		res = &TcpConnect{}
	default:
		return nil, fmt.Errorf("invalid crc code: got 0x%08x", tag)
	}
	err = res.UnmarshalTL(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func DecodeITcpPacket(r io.Reader) (ITcpPacket, error) {
	return decodeITcpPacket(r)
}

type TcpPing struct {
	ID int64
}

func (*TcpPing) CRC() uint32 {
	return uint32(0xbbe9627c)
}
func (t TcpPing) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.ID)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *TcpPing) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.ID)
	if err != nil {
		return err
	}
	return nil
}
func (*TcpPing) _ITcpPacket() {}

type TcpPong struct {
	ID int64
}

func (*TcpPong) CRC() uint32 {
	return uint32(0xbd4302c)
}
func (t TcpPong) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.ID)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *TcpPong) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.ID)
	if err != nil {
		return err
	}
	return nil
}
func (*TcpPong) _ITcpPacket() {}

type TcpPacket struct {
	Data []byte
}

func (*TcpPacket) CRC() uint32 {
	return uint32(0x9c17baf8)
}
func (t TcpPacket) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.Data)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *TcpPacket) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.Data)
	if err != nil {
		return err
	}
	return nil
}
func (*TcpPacket) _ITcpPacket() {}

type TcpQueryAnswer struct {
	ID   int64
	Data []byte
}

func (*TcpQueryAnswer) CRC() uint32 {
	return uint32(0xc048c311)
}
func (t TcpQueryAnswer) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.ID)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Data)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *TcpQueryAnswer) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.ID)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Data)
	if err != nil {
		return err
	}
	return nil
}
func (*TcpQueryAnswer) _ITcpPacket() {}

type TcpQueryError struct {
	ID      int64
	Code    int32
	Message string
}

func (*TcpQueryError) CRC() uint32 {
	return uint32(0x4cd2f602)
}
func (t TcpQueryError) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.ID)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Code)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Message)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *TcpQueryError) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.ID)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Code)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Message)
	if err != nil {
		return err
	}
	return nil
}
func (*TcpQueryError) _ITcpPacket() {}

type TcpQuery struct {
	ID   int64
	Data []byte
}

func (*TcpQuery) CRC() uint32 {
	return uint32(0x3af51908)
}
func (t TcpQuery) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.ID)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Data)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *TcpQuery) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.ID)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Data)
	if err != nil {
		return err
	}
	return nil
}
func (*TcpQuery) _ITcpPacket() {}

type TcpConnected struct {
	ID int64
}

func (*TcpConnected) CRC() uint32 {
	return uint32(0x636d41d6)
}
func (t TcpConnected) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.ID)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *TcpConnected) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.ID)
	if err != nil {
		return err
	}
	return nil
}
func (*TcpConnected) _ITcpPacket() {}

type TcpConnect struct {
	ID int64
}

func (*TcpConnect) CRC() uint32 {
	return uint32(0xa57c4261)
}
func (t TcpConnect) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	_ = err
	_ = b
	buf := bytes.NewBuffer(nil)
	b, err = tl.Marshal(t.ID)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *TcpConnect) UnmarshalTL(r io.Reader) error {
	var err error
	_ = err
	err = tl.Unmarshal(r, &t.ID)
	if err != nil {
		return err
	}
	return nil
}
func (*TcpConnect) _ITcpPacket() {}
