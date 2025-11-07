//go:build js && wasm
// +build js,wasm

package interaction

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/coder/websocket"

	"github.com/openimsdk/tools/log"
)

const StatusNetWorkChanged websocket.StatusCode = 3001

const (
	TextPing = "ping"
	TextPong = "pong"
)

type TextMessage struct {
	Type string          `json:"type"`
	Body json.RawMessage `json:"body"`
}

type JSWebSocket struct {
	ctx           context.Context
	pingHandler   PingPongHandler
	pongHandler   PingPongHandler
	readDeadline  time.Time
	writeDeadline time.Time
	ConnType      int
	conn          *websocket.Conn
}

func (w *JSWebSocket) SetReadDeadline(timeout time.Duration) error {
	w.readDeadline = time.Now().Add(timeout)
	return nil
}

func (w *JSWebSocket) SetWriteDeadline(timeout time.Duration) error {
	w.writeDeadline = time.Now().Add(timeout)
	return nil
}

func (w *JSWebSocket) SetReadLimit(limit int64) {
	w.conn.SetReadLimit(limit)
}

func (w *JSWebSocket) SetPingHandler(handler PingPongHandler) {
	w.pingHandler = handler
}

func (w *JSWebSocket) SetPongHandler(handler PingPongHandler) {
	w.pongHandler = handler
}

func (w *JSWebSocket) LocalAddr() string {
	return ""
}

func NewWebSocket(connType int) *JSWebSocket {
	return &JSWebSocket{
		ctx:      context.Background(),
		ConnType: connType,
	}
}

func (w *JSWebSocket) Close() error {
	return w.conn.Close(StatusNetWorkChanged, "Actively close the conn have old conn")
}

func (w *JSWebSocket) sendText(typ string, msg string) error {
	jsonStr, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	data, err := json.Marshal(TextMessage{
		Type: typ,
		Body: jsonStr,
	})
	if err != nil {
		return err
	}
	return w.Write(websocket.MessageText, data)
}

func (w *JSWebSocket) WriteMessage(messageType int, message []byte) error {
	switch messageType {
	case PingMessage:
		return w.sendText(TextPing, string(message))
	case PongMessage:
		return w.sendText(TextPong, string(message))
	default:
		return w.Write(websocket.MessageType(messageType), message)
	}
}

func (w *JSWebSocket) Read() (websocket.MessageType, []byte, error) {
	ctx := w.ctx
	if deadline := w.readDeadline; !deadline.IsZero() {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, deadline)
		defer cancel()
	}
	return w.conn.Read(ctx)
}

func (w *JSWebSocket) Write(typ websocket.MessageType, p []byte) error {
	ctx := w.ctx
	if deadline := w.writeDeadline; !deadline.IsZero() {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, deadline)
		defer cancel()
	}
	return w.conn.Write(ctx, typ, p)
}

func (w *JSWebSocket) handlerText(b []byte) error {
	var msg TextMessage
	if err := json.Unmarshal(b, &msg); err != nil {
		return err
	}
	var handler PingPongHandler
	switch msg.Type {
	case TextPing:
		handler = w.pingHandler
	case TextPong:
		handler = w.pongHandler
	default:
		return fmt.Errorf("wasm ws read text message %s", string(b))
	}
	var str string
	if err := json.Unmarshal(msg.Body, &str); err != nil {
		return err
	}
	if handler != nil {
		return handler(str)
	}
	return nil
}

func (w *JSWebSocket) ReadMessage() (int, []byte, error) {
	for {
		messageType, b, err := w.Read()
		if err != nil {
			return 0, nil, err
		}
		switch messageType {
		case websocket.MessageText:
			if err := w.handlerText(b); err != nil {
				return 0, nil, err
			}
			continue
		case websocket.MessageBinary:
			return int(messageType), b, nil
		default:
			return 0, nil, fmt.Errorf("wasm ws read type %d msg %v", messageType, b)
		}
	}
}

func (w *JSWebSocket) dial(ctx context.Context, urlStr string) (*websocket.Conn, *http.Response, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, nil, err
	}
	query := u.Query()
	query.Set("isMsgResp", "true")
	u.RawQuery = query.Encode()
	conn, httpResp, err := websocket.Dial(ctx, u.String(), nil)
	if err != nil {
		return nil, nil, err
	}
	if httpResp == nil {
		httpResp = &http.Response{
			StatusCode: http.StatusSwitchingProtocols,
		}
	}
	_, data, err := conn.Read(ctx)
	if err != nil {
		_ = conn.CloseNow()
		return nil, nil, fmt.Errorf("read response error %w", err)
	}
	var apiResp struct {
		ErrCode int    `json:"errCode"`
		ErrMsg  string `json:"errMsg"`
		ErrDlt  string `json:"errDlt"`
	}
	if err := json.Unmarshal(data, &apiResp); err != nil {
		return nil, nil, fmt.Errorf("unmarshal response error %w", err)
	}
	if apiResp.ErrCode == 0 {
		return conn, httpResp, nil
	}
	log.ZDebug(ctx, "ws msg read resp", "data", string(data))
	httpResp.Body = io.NopCloser(bytes.NewReader(data))
	return conn, httpResp, fmt.Errorf("read response error %d %s %s",
		apiResp.ErrCode, apiResp.ErrMsg, apiResp.ErrDlt)
}

func (w *JSWebSocket) Dial(urlStr string, _ http.Header) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	conn, httpResp, err := w.dial(ctx, urlStr)
	if err == nil {
		w.conn = conn
	}
	return httpResp, err
}

func (w *JSWebSocket) IsNil() bool {
	if w.conn != nil {
		return false
	}
	return true
}
