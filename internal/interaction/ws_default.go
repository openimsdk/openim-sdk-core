// Copyright © 2023 OpenIM SDK.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !js
// +build !js

package interaction

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Default struct {
	ConnType int
	conn     *websocket.Conn
	sendConn *websocket.Conn
}

func (d *Default) LocalAddr() string {
	//TODO implement me
	//panic("implement me")
	return ""
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
