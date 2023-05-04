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

package interaction

import (
	"context"
	"errors"
	"net/http"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/utils"
	"runtime"
	"sync"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

// 1.协程消息的收发模块一直运行，直到收到关闭信号（一般是用户退出登陆）
// 2.消息调用模块传递channle，收到消息后，调用channel，将消息传递给协程消息收发模块
// 3.消息通过ws发送后
// 4.
const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 51200
)
const (
	Closed = iota + 1
	Connecting
	Connected
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)
var ErrChanClosed = errors.New("send channel closed")
var ErrConnClosed = errors.New("conn has closed")
var ErrNotSupportMessageProtocol = errors.New("not support message protocol")
var ErrClientClosed = errors.New("client actively close the connection")
var ErrPanic = errors.New("panic error")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	w sync.Mutex
	// The long connection,can be set tcp or websocket.
	conn LongConn

	// Buffered channel of outbound messages.
	send chan Message
	//
	closedErr  error
	ctx        *ConnContext
	isCompress bool
	connStatus int
	syncer     *WsRespAsyn
	encoder    Encoder
	compressor Compressor
}
type Message struct {
	Message GeneralWsReq
	Resp    chan GeneralWsResp
}

func (c *Client) SendReqWaitResp(ctx context.Context, m proto.Message, reqIdentifier int32, resp proto.Message) error {
	data, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	msg := Message{
		Message: GeneralWsReq{
			ReqIdentifier: reqIdentifier,
			Token:         "",
			SendID:        mcontext.GetOpUserID(ctx),
			OperationID:   mcontext.GetOperationID(ctx),
			MsgIncr:       "",
			Data:          data,
		},
		Resp: make(chan GeneralWsResp, 1),
	}
	select {
	case <-ctx.Done():
		close(msg.Resp)
		return errors.New("send message timeout")
	case c.send <- msg:
	}
	select {
	case <-ctx.Done():
		return errors.New("wait response timeout")
	case v, ok := <-msg.Resp:
		if !ok {
			return errors.New("response channel closed")
		}
		if v.ErrCode != 0 {
			return errs.NewCodeError(v.ErrCode, v.ErrMsg)
		}
		if err := proto.Unmarshal(v.Data, resp); err != nil {
			return err
		}
		return nil
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump(ctx context.Context) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(pongWait)
	//c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
			c.closedErr = err
		}
		switch messageType {
		case MessageBinary:
			parseDataErr := c.handleMessage(message)
			if parseDataErr != nil {
				c.closedErr = parseDataErr
				return
			}
		case MessageText:
			c.closedErr = ErrNotSupportMessageProtocol
			return
		case PingMessage:
			err := c.writePongMsg()
			log.ZError(c.ctx, "writePongMsg", err)
		case CloseMessage:
			c.closedErr = ErrClientClosed
			return
		default:
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case <-c.ctx.Done():
			c.closedErr = c.ctx.Err()
			return
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(writeWait)
			if !ok {
				// The hub closed the channel.
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.ZError(c.ctx, "send close message error", err)
				}
				c.closedErr = ErrChanClosed
				return
			}
			msgIncr, tempChan := c.syncer.AddCh(message.message.SendID)

			for i := 0; i < 3; i++ {
				err := c.writeBinaryMsg(message.message)
				if err != nil {
					log.ZError(c.ctx, "send binary message error", err, "local address", c.conn.LocalAddr(), "message", message.message)
					_ = c.close()
					continue
				} else {
					break
				}
			}
			go func() {
				select {
				case <-time.After(time.Second * 3):
					log.ZError(c.ctx, "send message timeout", "local address", c.conn.LocalAddr(), "message", message.message)
					_ = c.close()
				case resp := <-message.resp:
					log.ZInfo(c.ctx, "receive response", "local address", c.conn.LocalAddr(), "message", message.message, "response", resp)
					_ = c.close()
				}
			}()
			c.syncer.DelCh(msgIncr)

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(writeWait)
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) writeBinaryMsg(req GeneralWsReq) error {
	encodeBuf, err := c.encoder.Encode(req)
	if err != nil {
		return err
	}
	_ = c.conn.SetWriteDeadline(writeWait)
	if c.isCompress {
		resultBuf, compressErr := c.compressor.Compress(encodeBuf)
		if compressErr != nil {
			return compressErr
		}
		return c.conn.WriteMessage(MessageBinary, resultBuf)
	} else {
		return c.conn.WriteMessage(MessageBinary, encodeBuf)
	}
}
func (c *Client) close() error {
	c.w.Lock()
	defer c.w.Unlock()
	c.connStatus = Closed
	return c.conn.Close()

}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
func (c *Client) handleMessage(message []byte) {
	var wsResp GeneralWsResp
	err := c.encoder.Decode(message, &wsResp)
	if err != nil {
		log.Error("decodeBinaryWs err", err.Error())
		return
	}
	ctx := context.WithValue(context.Background(), "operationID", wsResp.OperationID)
	log.Debug(wsResp.OperationID, "ws recv msg, code: ", wsResp.ErrCode, wsResp.ReqIdentifier)
	switch wsResp.ReqIdentifier {
	case constant.WSGetNewestSeq:
		if err := c.notifyResp(wsResp); err != nil {
			return utils.Wrap(err, "")
		}
		if err = w.doWSGetNewestSeq(wsResp); err != nil {
			log.Error(wsResp.OperationID, "doWSGetNewestSeq failed ", err.Error(), wsResp.ReqIdentifier, wsResp.MsgIncr)
		}
	case constant.WSPullMsgBySeqList:
		if err = w.doWSPullMsg(wsResp); err != nil {
			log.Error(wsResp.OperationID, "doWSPullMsg failed ", err.Error())
		}
	case constant.WSPushMsg:
		// todo
		//if constant.OnlyForTest == 1 {
		//	return
		//}
		if err = w.doWSPushMsg(wsResp); err != nil {
			log.Error(wsResp.OperationID, "doWSPushMsg failed ", err.Error())
		}
		//if err = w.doWSPushMsgForTest(*wsResp); err != nil {
		//	log.Error(wsResp.OperationID, "doWSPushMsgForTest failed ", err.Error())
		//}

	case constant.WSSendMsg:
		if err = w.doWSSendMsg(wsResp); err != nil {
			log.Error(wsResp.OperationID, "doWSSendMsg failed ", err.Error(), wsResp.ReqIdentifier, wsResp.MsgIncr)
		}
	case constant.WSKickOnlineMsg:
		log.Warn(wsResp.OperationID, "kick...  logout")
		w.kickOnline(wsResp)
		w.Logout(ctx)

	case constant.WsLogoutMsg:
		log.Warn(wsResp.OperationID, "WsLogoutMsg... Ws goroutine exit")
		if err = w.doWSLogoutMsg(wsResp); err != nil {
			log.Error(wsResp.OperationID, "doWSLogoutMsg failed ", err.Error())
		}
		runtime.Goexit()
	case constant.WSSendSignalMsg:
		log.Info(wsResp.OperationID, "signaling...")
		w.DoWSSignal(wsResp)
	case constant.WsSetBackgroundStatus:
		log.Info(wsResp.OperationID, "WsSetBackgroundStatus...")
		if err = w.setAppBackgroundStatus(wsResp); err != nil {
			log.Error(wsResp.OperationID, "WsSetBackgroundStatus failed ", err.Error(), wsResp.ReqIdentifier, wsResp.MsgIncr)
		}
		log.NewDebug(wsResp.OperationID, wsResp)
	default:
		log.Error(wsResp.OperationID, "type failed, ", wsResp.ReqIdentifier)
		return
	}
}
