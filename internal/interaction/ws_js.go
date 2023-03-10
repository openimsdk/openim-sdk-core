//go:build js && wasm
// +build js,wasm

package interaction

import (
	"context"
	"net/http"
	"nhooyr.io/websocket"
	"time"
)

type WebSocket struct {
	ConnType int
	conn     *websocket.Conn
	sendConn *websocket.Conn
}

func NewWebSocket(connType int) *WebSocket {
	return &WebSocket{ConnType: connType}
}

func (w *WebSocket) Close() error {
	return w.conn.Close(websocket.StatusGoingAway, "Actively close the conn have old conn")
}

func (w *WebSocket) WriteMessage(messageType int, message []byte) error {
	w.setSendConn(w.conn)
	return w.conn.Write(context.Background(), websocket.MessageType(messageType), message)
}
func (w *WebSocket) setSendConn(sendConn *websocket.Conn) {
	w.sendConn = sendConn
}
func (w *WebSocket) ReadMessage() (int, []byte, error) {
	messageType, b, err := w.conn.Read(context.Background())
	return int(messageType), b, err
}

func (w *WebSocket) SetReadTimeout(timeout int) error {
	return nil
}

func (w *WebSocket) SetWriteTimeout(timeout int) error {
	return nil
}

func (w *WebSocket) Dial(urlStr string, requestHeader http.Header) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	conn, httpResp, err := websocket.Dial(ctx, urlStr, nil)
	if err == nil {
		w.conn = conn
	}
	return httpResp, err
}

func (w *WebSocket) IsNil() bool {
	if w.conn != nil {
		return false
	}
	return true
}

func (w *WebSocket) SetConnNil() {
	w.conn = nil
}
func (w *WebSocket) CheckSendConnDiffNow() bool {
	return w.sendConn == w.conn
}
