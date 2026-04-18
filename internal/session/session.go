package session

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tonkeeper/gocoon/net/proxyconn"
	"github.com/tonkeeper/gocoon/tlcocoon"
	tlcocoonTypes "github.com/tonkeeper/gocoon/tlcocoon/types"
	tl "github.com/tonkeeper/tongo/tl"
)

const (
	maxFrameSize       = 1 << 20
	idlePingInterval   = 10 * time.Second
	idlePingCheckEvery = time.Second
)

var (
	ErrSessionClosed = errors.New("session closed")
)

// Session implements the cocoon TcpConnection framing layer over an established Conn.
//
// Wire format per packet: [uint32 LE size][uint32 LE seqno][payload]
// where payload is a boxed TL tcp.Packet.
type Session struct {
	conn      *proxyconn.Conn
	sendSeqno uint32
	recvSeqno uint32

	writeCh  chan writeRequest
	packetCh chan []byte
	doneCh   chan struct{}

	pendingMu sync.Mutex
	pending   map[int64]chan queryResult

	packetPendingMu sync.Mutex
	packetPending   map[[32]byte]*packetPendingState

	errOnce sync.Once
	errMu   sync.RWMutex
	err     error

	lastActivity atomic.Int64
}

type writeRequest struct {
	packet tlObject
	ackCh  chan error
}

type tlObject interface {
	CRC() uint32
	MarshalTL() ([]byte, error)
}

type queryResult struct {
	data []byte
	err  error
}

type RunQueryExOptions struct {
	MaxCoefficient   int32
	MaxTokens        int32
	Timeout          float64
	MinConfigVersion int32
	EnableDebug      *bool
}

type packetPendingState struct {
	resCh chan queryResult
	buf   []byte
}

// New creates a session and performs the initial tcp.connect/tcp.connected handshake.
// After handshake completion it starts background read/write/keepalive loops.
func New(conn *proxyconn.Conn) (*Session, error) {
	s := &Session{
		conn:          conn,
		writeCh:       make(chan writeRequest),
		packetCh:      make(chan []byte, 16),
		doneCh:        make(chan struct{}),
		pending:       make(map[int64]chan queryResult),
		packetPending: make(map[[32]byte]*packetPendingState),
		sendSeqno:     0,
		recvSeqno:     0,
	}
	s.markActivity()

	connectID, err := randInt64()
	if err != nil {
		return nil, err
	}
	if err := s.writePacket(&tlcocoonTypes.TcpConnect{ID: connectID}); err != nil {
		return nil, fmt.Errorf("send tcp.connect: %w", err)
	}

	payload, err := s.recvFrame()
	if err != nil {
		return nil, fmt.Errorf("recv tcp.connected: %w", err)
	}
	pkt, err := decodeTCPPacket(payload)
	if err != nil {
		return nil, fmt.Errorf("decode tcp.connected: %w", err)
	}
	connected, ok := pkt.(*tlcocoonTypes.TcpConnected)
	if !ok {
		return nil, fmt.Errorf("unexpected first packet %T (want tcp.connected)", pkt)
	}
	if connected.ID != connectID {
		return nil, fmt.Errorf("tcp.connected id mismatch: got %d, want %d", connected.ID, connectID)
	}

	go s.writeLoop()
	go s.readLoop()
	go s.keepAliveLoop()

	return s, nil
}

// Close closes the underlying connection and all background loops.
func (s *Session) Close() error {
	err := s.sessionErrOr()
	s.fail(err)
	return s.conn.Close()
}

// Query sends payload (a boxed TL function call) as tcp.query and returns
// tcp.queryAnswer data for the correlated query id.
func (s *Session) Query(ctx context.Context, payload []byte) ([]byte, error) {
	return s.MakeRequest(ctx, payload)
}

// MakeRequest implements tlcocoon.Requester by dispatching requests via tcp.query
// and correlating tcp.queryAnswer/tcp.queryError by query id.
func (s *Session) MakeRequest(ctx context.Context, msg []byte) ([]byte, error) {
	queryID, err := randInt64()
	if err != nil {
		return nil, err
	}

	resCh := make(chan queryResult, 1)
	s.pendingMu.Lock()
	s.pending[queryID] = resCh
	s.pendingMu.Unlock()

	if err := s.sendPacket(ctx, &tlcocoonTypes.TcpQuery{ID: queryID, Data: msg}); err != nil {
		s.dropPending(queryID)
		return nil, err
	}

	select {
	case <-ctx.Done():
		s.dropPending(queryID)
		return nil, ctx.Err()
	case <-s.doneCh:
		s.dropPending(queryID)
		return nil, s.sessionErrOr()
	case res := <-resCh:
		return res.data, res.err
	}
}

// SendMessage sends payload as tcp.packet.
func (s *Session) SendMessage(payload []byte) error {
	return s.sendPacket(context.Background(), &tlcocoonTypes.TcpPacket{Data: payload})
}

// RunClientQueryEx builds client.runQueryEx request, sends it over tcp.packet,
// and returns fully assembled answer bytes correlated by request_id.
func (s *Session) RunClientQueryEx(ctx context.Context, model string, query []byte, opts RunQueryExOptions) ([]byte, error) {
	var reqID [32]byte
	if _, err := rand.Read(reqID[:]); err != nil {
		return nil, fmt.Errorf("generate request id: %w", err)
	}
	req := tlcocoon.ClientRunQueryExRequest{
		ModelName:        model,
		Query:            query,
		MaxCoefficient:   opts.MaxCoefficient,
		MaxTokens:        opts.MaxTokens,
		Timeout:          opts.Timeout,
		RequestID:        reqID,
		MinConfigVersion: opts.MinConfigVersion,
		EnableDebug:      opts.EnableDebug,
	}
	body, err := req.MarshalTL()
	if err != nil {
		return nil, fmt.Errorf("marshal runQueryEx: %w", err)
	}
	payload := make([]byte, 4+len(body))
	binary.LittleEndian.PutUint32(payload[:4], req.CRC())
	copy(payload[4:], body)
	return s.doPacketRequest(ctx, payload, reqID)
}

// doPacketRequest sends payload as tcp.packet and correlates response chunks
// by requestID from client.queryAnswerPartEx/client.queryAnswerEx.
func (s *Session) doPacketRequest(ctx context.Context, payload []byte, requestID [32]byte) ([]byte, error) {
	resCh := make(chan queryResult, 1)

	s.packetPendingMu.Lock()
	if _, exists := s.packetPending[requestID]; exists {
		s.packetPendingMu.Unlock()
		return nil, fmt.Errorf("request already pending")
	}
	s.packetPending[requestID] = &packetPendingState{resCh: resCh}
	s.packetPendingMu.Unlock()

	if err := s.sendPacket(ctx, &tlcocoonTypes.TcpPacket{Data: payload}); err != nil {
		s.dropPacketPending(requestID)
		return nil, err
	}

	select {
	case <-ctx.Done():
		s.dropPacketPending(requestID)
		return nil, ctx.Err()
	case <-s.doneCh:
		s.dropPacketPending(requestID)
		return nil, s.sessionErrOr()
	case res := <-resCh:
		return res.data, res.err
	}
}

// RecvPacket receives payload from inbound tcp.packet.
func (s *Session) RecvPacket() ([]byte, error) {
	select {
	case <-s.doneCh:
		return nil, s.sessionErrOr()
	case p := <-s.packetCh:
		return p, nil
	}
}

func (s *Session) sendPacket(ctx context.Context, packet tlObject) error {
	ackCh := make(chan error, 1)
	req := writeRequest{packet: packet, ackCh: ackCh}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.doneCh:
		return s.sessionErrOr()
	case s.writeCh <- req:
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.doneCh:
		return s.sessionErrOr()
	case err := <-ackCh:
		return err
	}
}

// Err returns a terminal session error when the session is closed, otherwise nil.
func (s *Session) Err() error {
	select {
	case <-s.doneCh:
		return s.sessionErrOr()
	default:
		return nil
	}
}

func (s *Session) writeLoop() {
	for {
		select {
		case <-s.doneCh:
			return
		case req := <-s.writeCh:
			payload, err := marshalBoxed(req.packet)
			if err == nil {
				err = s.sendFrame(payload)
			}
			if err != nil {
				s.fail(fmt.Errorf("write packet: %w", err))
			}
			req.ackCh <- err
			if err != nil {
				return
			}
		}
	}
}

func (s *Session) readLoop() {
	for {
		payload, err := s.recvFrame()
		if err != nil {
			s.fail(fmt.Errorf("read frame: %w", err))
			return
		}

		pkt, err := decodeTCPPacket(payload)
		if err != nil {
			s.fail(fmt.Errorf("decode tcp packet: %w", err))
			return
		}

		switch p := pkt.(type) {
		case *tlcocoonTypes.TcpPing:
			_ = s.sendPacket(context.Background(), &tlcocoonTypes.TcpPong{ID: p.ID})
		case *tlcocoonTypes.TcpPong:
			// keepalive ack
		case *tlcocoonTypes.TcpQueryAnswer:
			s.resolvePending(p.ID, queryResult{data: p.Data})
		case *tlcocoonTypes.TcpQueryError:
			s.resolvePending(p.ID, queryResult{err: fmt.Errorf("tcp.queryError code=%d: %s", p.Code, p.Message)})
		case *tlcocoonTypes.TcpQuery:
			_ = s.sendPacket(context.Background(), &tlcocoonTypes.TcpQueryError{
				ID:      p.ID,
				Code:    -1,
				Message: "incoming tcp.query is not supported",
			})
		case *tlcocoonTypes.TcpPacket:
			if s.handleClientQueryPacket(p.Data) {
				continue
			}
			select {
			case <-s.doneCh:
				return
			case s.packetCh <- p.Data:
			}
		default:
			s.fail(fmt.Errorf("unexpected tcp packet %T", pkt))
			return
		}
	}
}

func (s *Session) keepAliveLoop() {
	ticker := time.NewTicker(idlePingCheckEvery)
	defer ticker.Stop()

	for {
		select {
		case <-s.doneCh:
			return
		case <-ticker.C:
			last := time.Unix(0, s.lastActivity.Load())
			if time.Since(last) < idlePingInterval {
				continue
			}
			pingID, err := randInt64()
			if err != nil {
				s.fail(fmt.Errorf("generate ping id: %w", err))
				return
			}
			_ = s.sendPacket(context.Background(), &tlcocoonTypes.TcpPing{ID: pingID})
		}
	}
}

func (s *Session) resolvePending(queryID int64, res queryResult) {
	s.pendingMu.Lock()
	ch, ok := s.pending[queryID]
	if ok {
		delete(s.pending, queryID)
	}
	s.pendingMu.Unlock()
	if !ok {
		return
	}
	select {
	case ch <- res:
	default:
	}
}

func (s *Session) dropPending(queryID int64) {
	s.pendingMu.Lock()
	delete(s.pending, queryID)
	s.pendingMu.Unlock()
}

func (s *Session) resolvePacketPending(requestID [32]byte, res queryResult) {
	s.packetPendingMu.Lock()
	st, ok := s.packetPending[requestID]
	if ok {
		delete(s.packetPending, requestID)
	}
	s.packetPendingMu.Unlock()
	if !ok {
		return
	}
	select {
	case st.resCh <- res:
	default:
	}
}

func (s *Session) dropPacketPending(requestID [32]byte) {
	s.packetPendingMu.Lock()
	delete(s.packetPending, requestID)
	s.packetPendingMu.Unlock()
}

func (s *Session) handleClientQueryPacket(data []byte) bool {
	ans, err := tlcocoonTypes.DecodeIClientQueryAnswerEx(bytes.NewReader(data))
	if err != nil {
		return false
	}

	switch a := ans.(type) {
	case *tlcocoonTypes.ClientQueryAnswerPartEx:
		s.appendPacketAnswer(a.RequestID, a.Answer, a.FinalInfo != nil)
		return true
	case *tlcocoonTypes.ClientQueryAnswerEx:
		s.appendPacketAnswer(a.RequestID, a.Answer, a.FinalInfo != nil)
		return true
	case *tlcocoonTypes.ClientQueryAnswerErrorEx:
		s.resolvePacketPending(a.RequestID, queryResult{
			err: fmt.Errorf("query error (code %d): %s", a.ErrorCode, a.Error),
		})
		return true
	default:
		return false
	}
}

func (s *Session) appendPacketAnswer(requestID [32]byte, chunk []byte, isFinal bool) {
	var (
		state *packetPendingState
		out   []byte
	)

	s.packetPendingMu.Lock()
	state = s.packetPending[requestID]
	if state != nil {
		state.buf = append(state.buf, chunk...)
		if isFinal {
			delete(s.packetPending, requestID)
			out = append([]byte(nil), state.buf...)
		}
	}
	s.packetPendingMu.Unlock()

	if state == nil || !isFinal {
		return
	}
	select {
	case state.resCh <- queryResult{data: out}:
	default:
	}
}

func (s *Session) fail(err error) {
	s.errOnce.Do(func() {
		s.errMu.Lock()
		s.err = err
		s.errMu.Unlock()
		close(s.doneCh)

		s.pendingMu.Lock()
		for id, ch := range s.pending {
			delete(s.pending, id)
			select {
			case ch <- queryResult{err: err}:
			default:
			}
		}
		s.pendingMu.Unlock()

		s.packetPendingMu.Lock()
		for id, st := range s.packetPending {
			delete(s.packetPending, id)
			select {
			case st.resCh <- queryResult{err: err}:
			default:
			}
		}
		s.packetPendingMu.Unlock()
	})
}

func (s *Session) sessionErrOr() error {
	s.errMu.RLock()
	defer s.errMu.RUnlock()
	if s.err != nil {
		return s.err
	}
	return ErrSessionClosed
}

func (s *Session) markActivity() {
	s.lastActivity.Store(time.Now().UnixNano())
}

func decodeTCPPacket(payload []byte) (tlObject, error) {
	if len(payload) < 4 {
		return nil, fmt.Errorf("payload too short: %d", len(payload))
	}
	magic := binary.LittleEndian.Uint32(payload[:4])
	body := payload[4:]

	var pkt tlObject
	switch magic {
	case (&tlcocoonTypes.TcpPing{}).CRC():
		pkt = &tlcocoonTypes.TcpPing{}
	case (&tlcocoonTypes.TcpPong{}).CRC():
		pkt = &tlcocoonTypes.TcpPong{}
	case (&tlcocoonTypes.TcpPacket{}).CRC():
		pkt = &tlcocoonTypes.TcpPacket{}
	case (&tlcocoonTypes.TcpQueryAnswer{}).CRC():
		pkt = &tlcocoonTypes.TcpQueryAnswer{}
	case (&tlcocoonTypes.TcpQueryError{}).CRC():
		pkt = &tlcocoonTypes.TcpQueryError{}
	case (&tlcocoonTypes.TcpQuery{}).CRC():
		pkt = &tlcocoonTypes.TcpQuery{}
	case (&tlcocoonTypes.TcpConnected{}).CRC():
		pkt = &tlcocoonTypes.TcpConnected{}
	case (&tlcocoonTypes.TcpConnect{}).CRC():
		pkt = &tlcocoonTypes.TcpConnect{}
	default:
		return nil, fmt.Errorf("unknown tcp packet magic 0x%08x", magic)
	}
	if err := tl.Unmarshal(bytes.NewReader(body), pkt); err != nil {
		return nil, err
	}
	return pkt, nil
}

// sendFrame writes [uint32 LE size][uint32 LE seqno][payload] to the connection.
func (s *Session) sendFrame(payload []byte) error {
	var hdr [8]byte
	binary.LittleEndian.PutUint32(hdr[0:4], uint32(len(payload)))
	binary.LittleEndian.PutUint32(hdr[4:8], s.sendSeqno)
	s.sendSeqno++

	if _, err := s.conn.Write(hdr[:]); err != nil {
		return err
	}
	if _, err := s.conn.Write(payload); err != nil {
		return err
	}
	s.markActivity()
	return nil
}

// recvFrame reads [uint32 LE size][uint32 LE seqno][payload] and returns payload.
func (s *Session) recvFrame() ([]byte, error) {
	var hdr [8]byte
	if _, err := io.ReadFull(s.conn, hdr[:]); err != nil {
		return nil, fmt.Errorf("read frame header: %w", err)
	}
	size := binary.LittleEndian.Uint32(hdr[0:4])
	seqno := binary.LittleEndian.Uint32(hdr[4:8])
	if seqno != s.recvSeqno {
		return nil, fmt.Errorf("unexpected seqno: got %d want %d", seqno, s.recvSeqno)
	}

	if size > maxFrameSize {
		return nil, fmt.Errorf("frame too large: %d bytes", size)
	}
	payload := make([]byte, size)
	if _, err := io.ReadFull(s.conn, payload); err != nil {
		return nil, fmt.Errorf("read frame payload: %w", err)
	}
	s.recvSeqno++
	s.markActivity()
	return payload, nil
}

func (s *Session) writePacket(packet tlObject) error {
	payload, err := marshalBoxed(packet)
	if err != nil {
		return err
	}
	return s.sendFrame(payload)
}

func marshalBoxed(v tlObject) ([]byte, error) {
	body, err := v.MarshalTL()
	if err != nil {
		return nil, err
	}
	out := make([]byte, 4+len(body))
	binary.LittleEndian.PutUint32(out[:4], v.CRC())
	copy(out[4:], body)
	return out, nil
}

func randInt64() (int64, error) {
	n, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 63))
	if err != nil {
		return 0, err
	}
	return n.Int64(), nil
}
