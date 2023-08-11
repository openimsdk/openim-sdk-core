// Copyright © 2023 OpenIM SDK. All rights reserved.
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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/sdkerrs"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024 * 1024
)

const (
	DefaultNotConnect = iota
	Closed            = iota + 1
	Connecting
	Connected
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var (
	ErrChanClosed                = errors.New("send channel closed")
	ErrConnClosed                = errors.New("conn has closed")
	ErrNotSupportMessageProtocol = errors.New("not support message protocol")
	ErrClientClosed              = errors.New("client actively close the connection")
	ErrPanic                     = errors.New("panic error")
)

type LongConnMgr struct {
	//conn status mutex
	w          sync.Mutex
	connStatus int
	// The long connection,can be set tcp or websocket.
	conn     LongConn
	listener open_im_sdk_callback.OnConnListener
	// Buffered channel of outbound messages.
	send               chan Message
	pushMsgAndMaxSeqCh chan common.Cmd2Value
	conversationCh     chan common.Cmd2Value
	loginMgrCh         chan common.Cmd2Value
	heartbeatCh        chan common.Cmd2Value
	closedErr          error
	ctx                context.Context
	IsCompression      bool
	Syncer             *WsRespAsyn
	encoder            Encoder
	compressor         Compressor
	IsBackground       bool
	// write conn lock
	connWrite *sync.Mutex
}

type Message struct {
	Message GeneralWsReq
	Resp    chan *GeneralWsResp
}

func NewLongConnMgr(ctx context.Context, listener open_im_sdk_callback.OnConnListener, heartbeatCmdCh, pushMsgAndMaxSeqCh, loginMgrCh chan common.Cmd2Value) *LongConnMgr {
	l := &LongConnMgr{listener: listener, pushMsgAndMaxSeqCh: pushMsgAndMaxSeqCh,
		loginMgrCh: loginMgrCh, IsCompression: true,
		Syncer: NewWsRespAsyn(), encoder: NewGobEncoder(), compressor: NewGzipCompressor()}
	l.send = make(chan Message, 10)
	l.conn = NewWebSocket(WebSocket)
	l.connWrite = new(sync.Mutex)
	l.ctx = ctx
	l.heartbeatCh = heartbeatCmdCh
	return l
}
func (c *LongConnMgr) Run(ctx context.Context) {
	//fmt.Println(mcontext.GetOperationID(ctx), "login run", string(debug.Stack()))
	go c.readPump(ctx)
	go c.writePump(ctx)
	go c.heartbeat(ctx)
}

func (c *LongConnMgr) SendReqWaitResp(ctx context.Context, m proto.Message, reqIdentifier int, resp proto.Message) error {
	data, err := proto.Marshal(m)
	if err != nil {
		return sdkerrs.ErrArgs
	}
	msg := Message{
		Message: GeneralWsReq{
			ReqIdentifier: reqIdentifier,
			SendID:        ccontext.Info(ctx).UserID(),
			OperationID:   ccontext.Info(ctx).OperationID(),
			Data:          data,
		},
		Resp: make(chan *GeneralWsResp, 1),
	}
	c.send <- msg
	log.ZDebug(ctx, "send message to send channel success", "msg", m, "reqIdentifier", reqIdentifier)
	select {
	case <-ctx.Done():
		return sdkerrs.ErrCtxDeadline
	case v, ok := <-msg.Resp:
		if !ok {
			return errors.New("response channel closed")
		}
		if v.ErrCode != 0 {
			return errs.NewCodeError(v.ErrCode, v.ErrMsg)
		}
		if err := proto.Unmarshal(v.Data, resp); err != nil {
			return sdkerrs.ErrArgs
		}
		return nil
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.

func (c *LongConnMgr) readPump(ctx context.Context) {
	log.ZDebug(ctx, "readPump start", "goroutine ID:", getGoroutineID())
	defer func() {
		log.ZWarn(c.ctx, "readPump closed", c.closedErr)
	}()
	connNum := 0
	//c.conn.SetPongHandler(funcation(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		ctx = ccontext.WithOperationID(ctx, utils.OperationIDGenerator())
		needRecon, err := c.reConn(ctx, &connNum)
		if !needRecon {
			c.closedErr = err
			return
		}
		if err != nil {
			log.ZWarn(c.ctx, "reConn", err)
			time.Sleep(time.Second * 1)
			continue
		}
		c.conn.SetReadLimit(maxMessageSize)
		_ = c.conn.SetReadDeadline(pongWait)
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			//if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			//	log.Printf("error: %v", err)
			//}
			//break
			//c.closedErr = err
			log.ZError(c.ctx, "readMessage err", err, "goroutine ID:", getGoroutineID())
			_ = c.close()
			continue
		}
		switch messageType {
		case MessageBinary:
			err := c.handleMessage(message)
			if err != nil {
				c.closedErr = err
				return
			}
		case MessageText:
			c.closedErr = ErrNotSupportMessageProtocol
			return
		//case PingMessage:
		//	err := c.writePongMsg()
		//	log.ZError(c.ctx, "writePongMsg", err)
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
func (c *LongConnMgr) writePump(ctx context.Context) {
	log.ZDebug(ctx, "writePump start", "goroutine ID:", getGoroutineID())

	defer func() {
		c.close()
		close(c.send)
	}()
	for {
		select {
		case <-ctx.Done():
			c.closedErr = ctx.Err()
			return
		case message, ok := <-c.send:
			if !ok {
				// The hub closed the channel.
				_ = c.conn.SetWriteDeadline(writeWait)
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.ZError(c.ctx, "send close message error", err)
				}
				c.closedErr = ErrChanClosed
				return
			}
			log.ZDebug(c.ctx, "writePump recv message", "reqIdentifier", message.Message.ReqIdentifier,
				"operationID", message.Message.OperationID, "sendID", message.Message.SendID)
			resp, err := c.sendAndWaitResp(&message.Message)
			if err != nil {
				resp = &GeneralWsResp{
					ReqIdentifier: message.Message.ReqIdentifier,
					OperationID:   message.Message.OperationID,
					Data:          nil,
				}
				if code, ok := errs.Unwrap(err).(errs.CodeError); ok {
					resp.ErrCode = code.Code()
					resp.ErrMsg = code.Msg()
				} else {
					log.ZError(c.ctx, "writeBinaryMsgAndRetry failed", err, "wsReq", message.Message)
				}

			}
			nErr := c.Syncer.notifyCh(message.Resp, resp, 1)
			if nErr != nil {
				log.ZError(c.ctx, "TriggerCmdNewMsgCome failed", nErr, "wsResp", resp)
			}
		}
	}
}

func (c *LongConnMgr) heartbeat(ctx context.Context) {
	log.ZDebug(ctx, "heartbeat start", "goroutine ID:", getGoroutineID())
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		log.ZWarn(c.ctx, "heartbeat closed", nil, "heartbeat", "heartbeat done sdk logout.....")
	}()
	for {
		select {
		case <-ctx.Done():
			log.ZInfo(ctx, "heartbeat done sdk logout.....")
			return
		case <-c.heartbeatCh:
			c.sendPingToServer(ctx)
		case <-ticker.C:
			c.sendPingToServer(ctx)
		}
	}

}
func getGoroutineID() int64 {
	buf := make([]byte, 64)
	buf = buf[:runtime.Stack(buf, false)]
	idField := strings.Fields(strings.TrimPrefix(string(buf), "goroutine "))[0]
	id, err := strconv.ParseInt(idField, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
func (c *LongConnMgr) sendPingToServer(ctx context.Context) {
	if c.conn == nil {
		return
	}
	var m sdkws.GetMaxSeqReq
	m.UserID = ccontext.Info(ctx).UserID()
	opID := utils.OperationIDGenerator()
	sCtx := ccontext.WithOperationID(c.ctx, opID)
	log.ZInfo(sCtx, "ping and getMaxSeq start", "goroutine ID:", getGoroutineID())
	data, err := proto.Marshal(&m)
	if err != nil {
		log.ZError(sCtx, "proto.Marshal", err)
		return
	}
	req := &GeneralWsReq{
		ReqIdentifier: constant.GetNewestSeq,
		SendID:        m.UserID,
		OperationID:   opID,
		Data:          data,
	}
	resp, err := c.sendAndWaitResp(req)
	if err != nil {
		log.ZError(sCtx, "sendAndWaitResp", err)
		_ = c.close()
		time.Sleep(time.Second * 1)
		return
	} else {
		if resp.ErrCode != 0 {
			log.ZError(sCtx, "getMaxSeq failed", nil, "errCode:", resp.ErrCode, "errMsg:", resp.ErrMsg)
		}
		var wsSeqResp sdkws.GetMaxSeqResp
		err = proto.Unmarshal(resp.Data, &wsSeqResp)
		if err != nil {
			log.ZError(sCtx, "proto.Unmarshal", err)
		}
		var cmd sdk_struct.CmdMaxSeqToMsgSync
		cmd.ConversationMaxSeqOnSvr = wsSeqResp.MaxSeqs

		err := common.TriggerCmdMaxSeq(sCtx, &cmd, c.pushMsgAndMaxSeqCh)
		if err != nil {
			log.ZError(sCtx, "TriggerCmdMaxSeq failed", err)
		}
	}
}
func (c *LongConnMgr) sendAndWaitResp(msg *GeneralWsReq) (*GeneralWsResp, error) {
	tempChan, err := c.writeBinaryMsgAndRetry(msg)
	defer c.Syncer.DelCh(msg.MsgIncr)
	if err != nil {
		return nil, err
	} else {
		select {
		case resp := <-tempChan:
			return resp, nil
		case <-time.After(time.Second * 5):
			return nil, sdkerrs.ErrNetworkTimeOut
		}

	}
}

func (c *LongConnMgr) writeBinaryMsgAndRetry(msg *GeneralWsReq) (chan *GeneralWsResp, error) {
	msgIncr, tempChan := c.Syncer.AddCh(msg.SendID)
	msg.MsgIncr = msgIncr
	if c.GetConnectionStatus() != Connected && msg.ReqIdentifier == constant.GetNewestSeq {
		return tempChan, sdkerrs.ErrNetwork.Wrap("connection closed,conning...")
	}
	for i := 0; i < 60; i++ {
		err := c.writeBinaryMsg(*msg)
		if err != nil {
			log.ZError(c.ctx, "send binary message error", err, "message", msg)
			c.closedErr = err
			_ = c.close()
			time.Sleep(time.Second * 1)
			continue
		} else {
			return tempChan, nil
		}
	}
	return nil, sdkerrs.ErrNetwork.Wrap("send binary message error")
}

func (c *LongConnMgr) writeBinaryMsg(req GeneralWsReq) error {
	c.connWrite.Lock()
	defer c.connWrite.Unlock()
	encodeBuf, err := c.encoder.Encode(req)
	if err != nil {
		return err
	}
	if c.GetConnectionStatus() != Connected {
		return sdkerrs.ErrNetwork.Wrap("connection closed,re conning...")
	}
	_ = c.conn.SetWriteDeadline(writeWait)
	if c.IsCompression {
		resultBuf, compressErr := c.compressor.Compress(encodeBuf)
		if compressErr != nil {
			return compressErr
		}
		return c.conn.WriteMessage(MessageBinary, resultBuf)
	} else {
		return c.conn.WriteMessage(MessageBinary, encodeBuf)
	}
}
func (c *LongConnMgr) close() error {
	c.w.Lock()
	defer c.w.Unlock()
	if c.connStatus == Closed || c.connStatus == Connecting || c.connStatus == DefaultNotConnect {
		return nil
	}
	c.connStatus = Closed
	log.ZWarn(c.ctx, "conn closed", c.closedErr)
	return c.conn.Close()
}

func (c *LongConnMgr) handleMessage(message []byte) error {
	if c.IsCompression {
		var decompressErr error
		message, decompressErr = c.compressor.DeCompress(message)
		if decompressErr != nil {
			log.ZError(c.ctx, "DeCompress failed", decompressErr, message)
			return sdkerrs.ErrMsgDeCompression
		}
	}
	var wsResp GeneralWsResp
	err := c.encoder.Decode(message, &wsResp)
	if err != nil {
		log.ZError(c.ctx, "decodeBinaryWs err", err, "message", message)
		return sdkerrs.ErrMsgDecodeBinaryWs
	}
	ctx := context.WithValue(c.ctx, "operationID", wsResp.OperationID)
	log.ZInfo(ctx, "recv msg", "errCode", wsResp.ErrCode, "errMsg", wsResp.ErrMsg,
		"reqIdentifier", wsResp.ReqIdentifier)
	switch wsResp.ReqIdentifier {
	case constant.PushMsg:
		if err = c.doPushMsg(ctx, wsResp); err != nil {
			log.ZError(ctx, "doWSPushMsg failed", err, "wsResp", wsResp)
		}
	case constant.LogoutMsg:
		if err := c.Syncer.NotifyResp(ctx, wsResp); err != nil {
			log.ZError(ctx, "notifyResp failed", err, "wsResp", wsResp)
		}
		return sdkerrs.ErrLoginOut
	case constant.KickOnlineMsg:
		log.ZDebug(ctx, "client kicked offline")
		c.listener.OnKickedOffline()
		_ = common.TriggerCmdLogOut(ctx, c.loginMgrCh)
		return errors.New("client kicked offline")
	case constant.GetNewestSeq:
		fallthrough
	case constant.PullMsgBySeqList:
		fallthrough
	case constant.SendMsg:
		fallthrough
	case constant.SendSignalMsg:
		fallthrough
	case constant.SetBackgroundStatus:
		if err := c.Syncer.NotifyResp(ctx, wsResp); err != nil {
			log.ZError(ctx, "notifyResp failed", err, "reqIdentifier", wsResp.ReqIdentifier, "errCode",
				wsResp.ErrCode, "errMsg", wsResp.ErrMsg, "msgIncr", wsResp.MsgIncr, "operationID", wsResp.OperationID)
		}
	default:
		// log.Error(wsResp.OperationID, "type failed, ", wsResp.ReqIdentifier)
		return sdkerrs.ErrMsgBinaryTypeNotSupport
	}
	return nil
}
func (c *LongConnMgr) IsConnected() bool {
	c.w.Lock()
	defer c.w.Unlock()
	if c.connStatus == Connected {
		return true
	}
	return false

}
func (c *LongConnMgr) GetConnectionStatus() int {
	c.w.Lock()
	defer c.w.Unlock()
	return c.connStatus
}
func (c *LongConnMgr) reConn(ctx context.Context, num *int) (needRecon bool, err error) {
	if c.IsConnected() {
		return true, nil
	}
	c.connWrite.Lock()
	defer c.connWrite.Unlock()
	log.ZDebug(ctx, "conn start")
	c.listener.OnConnecting()
	c.w.Lock()
	c.connStatus = Connecting
	c.w.Unlock()
	url := fmt.Sprintf("%s?sendID=%s&token=%s&platformID=%d&operationID=%s&isBackground=%t", ccontext.Info(ctx).WsAddr(),
		ccontext.Info(ctx).UserID(), ccontext.Info(ctx).Token(), ccontext.Info(ctx).PlatformID(), ccontext.Info(ctx).OperationID(), c.IsBackground)
	if c.IsCompression {
		url += fmt.Sprintf("&compression=%s", "gzip")
	}
	resp, err := c.conn.Dial(url, nil)
	if err != nil {
		c.w.Lock()
		c.connStatus = Closed
		c.w.Unlock()
		if resp != nil {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return true, err
			}
			log.ZInfo(ctx, "reConn resp", "body", string(body))
			var apiResp struct {
				ErrCode int    `json:"errCode"`
				ErrMsg  string `json:"errMsg"`
				ErrDlt  string `json:"errDlt"`
			}
			if err := json.Unmarshal(body, &apiResp); err != nil {
				return true, err
			}
			switch apiResp.ErrCode {
			case
				errs.TokenExpiredError,
				errs.TokenInvalidError,
				errs.TokenMalformedError,
				errs.TokenNotValidYetError,
				errs.TokenUnknownError,
				errs.TokenKickedError,
				errs.TokenNotExistError:
				c.listener.OnUserTokenExpired()
				_ = common.TriggerCmdLogOut(ctx, c.loginMgrCh)
			default:
				c.listener.OnConnectFailed(int32(apiResp.ErrCode), apiResp.ErrMsg)
			}
			log.ZWarn(ctx, "long conn establish failed", sdkerrs.New(apiResp.ErrCode, apiResp.ErrMsg, apiResp.ErrDlt))
			return false, errs.NewCodeError(apiResp.ErrCode, apiResp.ErrMsg).WithDetail(apiResp.ErrDlt).Wrap()
		}
		c.listener.OnConnectFailed(sdkerrs.NetworkError, err.Error())
		return true, err
	}
	c.listener.OnConnectSuccess()
	c.ctx = newContext(c.conn.LocalAddr())
	c.ctx = context.WithValue(ctx, "ConnContext", c.ctx)
	c.w.Lock()
	c.connStatus = Connected
	c.w.Unlock()
	*num++
	log.ZInfo(c.ctx, "long conn establish success", "localAddr", c.conn.LocalAddr(), "connNum", *num)
	_ = common.TriggerCmdConnected(ctx, c.pushMsgAndMaxSeqCh)
	return true, nil
}

func (c *LongConnMgr) doPushMsg(ctx context.Context, wsResp GeneralWsResp) error {
	var msg sdkws.PushMessages
	err := proto.Unmarshal(wsResp.Data, &msg)
	if err != nil {
		return err
	}
	return common.TriggerCmdPushMsg(ctx, &msg, c.pushMsgAndMaxSeqCh)
}
func (c *LongConnMgr) Close(ctx context.Context) {
	if c.GetConnectionStatus() == Connected {
		log.ZInfo(ctx, "network change conn close")
		c.closedErr = errors.New("closed by client network change")
		_ = c.close()
	} else {
		log.ZInfo(ctx, "conn already closed")
	}

}
func (c *LongConnMgr) SetBackground(isBackground bool) {
	c.IsBackground = isBackground
}
