//go:build !js
// +build !js

package interaction

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type Default struct {
	ConnType int
	conn     *websocket.Conn
	sendConn *websocket.Conn
}

func (d *Default) LocalAddr() string {
	return d.conn.LocalAddr().String()
}

func NewWebSocket(connType int) *Default {
	return &Default{ConnType: connType}
}
func (d *Default) Close() error {
	return d.conn.Close()
}

func (d *Default) WriteMessage(messageType int, message []byte) error {
	d.setSendConn(d.conn)
	return d.conn.WriteMessage(messageType, message)
}

func (d *Default) setSendConn(sendConn *websocket.Conn) {
	d.sendConn = sendConn
}

func (d *Default) ReadMessage() (int, []byte, error) {
	return d.conn.ReadMessage()
}
func (d *Default) SetReadTimeout(timeout int) error {
	return d.conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
}

func (d *Default) SetWriteTimeout(timeout int) error {
	return d.conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
}

func (d *Default) Dial(urlStr string, requestHeader http.Header) (*http.Response, error) {
	conn, httpResp, err := websocket.DefaultDialer.Dial(urlStr, requestHeader)
	if err == nil {
		d.conn = conn
	}
	return httpResp, err

}

func (d *Default) IsNil() bool {
	if d.conn != nil {
		return false
	}
	return true
}

func (d *Default) SetConnNil() {
	d.conn = nil
}
func (d *Default) CheckSendConnDiffNow() bool {
	return d.conn == d.sendConn
}
