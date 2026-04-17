package tlcocoon

import (
	"bytes"
	"context"
	types "github.com/tonkeeper/gocoon/pkg/tlcocoon/types"
	tl "github.com/tonkeeper/tongo/tl"
)

type TestQueryRequest struct {
	ID int64
}

func (*TestQueryRequest) CRC() uint32 {
	return uint32(0x394e9143)
}
func (t TestQueryRequest) MarshalTL() ([]byte, error) {
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
func TestQuery(ctx context.Context, m Requester, i TestQueryRequest) (types.TestAnswer, error) {
	var res types.TestAnswer
	return res, request(ctx, m, &i, &res)
}

type TestAPI interface {
	Query(ctx context.Context, i TestQueryRequest) (types.TestAnswer, error)
}
type Test struct {
	requester Requester
}

func NewTest(requester Requester) *Test {
	return &Test{requester: requester}
}
func (c *Test) Query(ctx context.Context, i TestQueryRequest) (types.TestAnswer, error) {
	return TestQuery(ctx, c.requester, i)
}

var _ TestAPI = (*Test)(nil)
