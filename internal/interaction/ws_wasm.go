// +build js,wasm

package interaction

import (
	"context"
	"net/http"
	"nhooyr.io/websocket"
	"time"
)

type WebSocketJS struct {
	ConnType int
	conn     *websocket.Conn
	sendConn *websocket.Conn
}

func NewWebSocketJS(connType int) *WebSocketJS {
	return &WebSocketJS{ConnType: connType}
}

func (w *WebSocketJS) Close() error {
	return w.conn.Close(websocket.StatusGoingAway, "Actively close the conn have old conn")
}

func (w *WebSocketJS) WriteMessage(messageType int, message []byte) error {
	w.setSendConn(w.conn)
	return w.conn.Write(context.Background(), websocket.MessageType(messageType), message)
}
func (w *WebSocketJS) setSendConn(sendConn *websocket.Conn) {
	w.sendConn = sendConn
}
func (w *WebSocketJS) ReadMessage() (int, []byte, error) {
	messageType, b, err := w.conn.Read(context.Background())
	return int(messageType), b, err
}

func (w *WebSocketJS) SetReadTimeout(timeout int) error {
	return nil
}

func (w *WebSocketJS) SetWriteTimeout(timeout int) error {
	return nil
}

func (w *WebSocketJS) Dial(urlStr string, requestHeader http.Header) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	conn, httpResp, err := websocket.Dial(ctx, urlStr, nil)
	if err == nil {
		w.conn = conn
	}
	return httpResp, err
}

func (w *WebSocketJS) IsNil() bool {
	if w.conn != nil {
		return false
	}
	return true
}

func (w *WebSocketJS) SetConnNil() {
	w.conn = nil
}
func (w *WebSocketJS) CheckSendConnDiffNow() bool {
	return w.sendConn == w.conn
}
