package interaction

import (
	"context"
	"errors"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"net/http"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"sync"
	"time"
)

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
	closedErr          error
	ctx                context.Context
	IsCompression      bool
	Syncer             *WsRespAsyn
	encoder            Encoder
	compressor         Compressor
}
type Message struct {
	Message GeneralWsReq
	Resp    chan GeneralWsResp
}

func NewLongConnMgr(ctx context.Context, listener open_im_sdk_callback.OnConnListener, pushMsgAndMaxSeqCh, conversationCh chan common.Cmd2Value) *LongConnMgr {
	l := &LongConnMgr{listener: listener, pushMsgAndMaxSeqCh: pushMsgAndMaxSeqCh,
		conversationCh: conversationCh, IsCompression: ccontext.Info(ctx).IsCompression(),
		Syncer: NewWsRespAsyn(), encoder: NewGobEncoder(), compressor: NewGzipCompressor()}
	l.send = make(chan Message, 10)
	l.conn = NewWebSocket(WebSocket)
	go l.readPump(ctx)
	go l.writePump(ctx)
	return l
}
func (c *LongConnMgr) SendReqWaitResp(ctx context.Context, m proto.Message, reqIdentifier int, resp proto.Message) error {
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
	c.send <- msg
	log.ZDebug(ctx, "send message to send channel success", "msg", m, "reqIdentifier", reqIdentifier)
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
func (c *LongConnMgr) readPump(ctx context.Context) {
	defer func() {
		//c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(pongWait)
	//c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		err := c.reConn(ctx)
		if err != nil {
			log.ZError(c.ctx, "reConn", err)
			time.Sleep(time.Second * 1)
			continue
		}
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			//if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			//	log.Printf("error: %v", err)
			//}
			//break
			//c.closedErr = err

		}
		switch messageType {
		case MessageBinary:
			c.handleMessage(message)
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
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		close(c.send)
	}()
	for {
		select {
		case <-ctx.Done():
			c.closedErr = ctx.Err()
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
			tempChan, err := c.writeBinaryMsgAndRetry(&message.Message)
			if err != nil {
				resp := GeneralWsResp{
					ReqIdentifier: message.Message.ReqIdentifier,
					ErrCode:       1,
					ErrMsg:        "",
					OperationID:   message.Message.OperationID,
					Data:          nil,
				}
				err := c.Syncer.notifyCh(message.Resp, resp, 1)
				if err != nil {
					//log.Warn(wsResp.OperationID, "TriggerCmdNewMsgCome failed ", err.Error(), ch, wsResp.ReqIdentifier, wsResp.MsgIncr)
					log.ZError(c.ctx, "TriggerCmdNewMsgCome failed", err, "wsResp", resp)
				}
			} else {
				go func() {
					select {
					case resp := <-tempChan:
						log.ZInfo(c.ctx, "receive response", "local address", c.conn.LocalAddr(), "message", message.Message, "response", resp)
						err := c.Syncer.notifyCh(message.Resp, resp, 1)
						if err != nil {
							//log.Warn(wsResp.OperationID, "TriggerCmdNewMsgCome failed ", err.Error(), ch, wsResp.ReqIdentifier, wsResp.MsgIncr)
							log.ZError(c.ctx, "TriggerCmdNewMsgCome failed", err, "wsResp", resp)
						}
						log.ZInfo(c.ctx, "receive response", "local address", c.conn.LocalAddr(), "message", message.Message, "response", resp)
						//_ = c.close()
					case <-time.After(time.Second * 3):
						resp := GeneralWsResp{
							ReqIdentifier: message.Message.ReqIdentifier,
							ErrCode:       0,
							ErrMsg:        "",
							OperationID:   message.Message.OperationID,
							Data:          nil,
						}
						err := c.Syncer.notifyCh(message.Resp, resp, 1)
						if err != nil {
							//log.Warn(wsResp.OperationID, "TriggerCmdNewMsgCome failed ", err.Error(), ch, wsResp.ReqIdentifier, wsResp.MsgIncr)
							log.ZError(c.ctx, "TriggerCmdNewMsgCome failed", err, "wsResp", resp)
						}
					}
					c.Syncer.DelCh(message.Message.MsgIncr)

				}()

			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(writeWait)
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *LongConnMgr) writeBinaryMsgAndRetry(msg *GeneralWsReq) (chan GeneralWsResp, error) {
	msgIncr, tempChan := c.Syncer.AddCh(msg.SendID)
	msg.MsgIncr = msgIncr
	for i := 0; i < 3; i++ {
		err := c.writeBinaryMsg(*msg)
		if err != nil {
			log.ZError(c.ctx, "send binary message error", err, "local address", c.conn.LocalAddr(), "message", msg)
			_ = c.close()
			time.Sleep(time.Second * 1)
			continue
		} else {
			return tempChan, nil
		}
	}
	return nil, errors.New("send binary message error")
}

func (c *LongConnMgr) writeBinaryMsg(req GeneralWsReq) error {
	encodeBuf, err := c.encoder.Encode(req)
	if err != nil {
		return err
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
	c.connStatus = Closed
	return c.conn.Close()

}

func (c *LongConnMgr) handleMessage(message []byte) {
	if c.IsCompression {
		var decompressErr error
		message, decompressErr = c.compressor.DeCompress(message)
		if decompressErr != nil {
			log.ZError(c.ctx, "DeCompress failed", decompressErr, message)
			return
		}
	}
	var wsResp GeneralWsResp
	err := c.encoder.Decode(message, &wsResp)
	if err != nil {
		log.ZError(c.ctx, "decodeBinaryWs err", err, "message", message)
		return
	}
	ctx := context.WithValue(c.ctx, "operationID", wsResp.OperationID)
	log.ZInfo(ctx, "recv msg", "code", wsResp.ErrCode, "reqIdentifier", wsResp.ReqIdentifier)
	switch wsResp.ReqIdentifier {
	case constant.PushMsg:
		if err = c.doPushMsg(ctx, wsResp); err != nil {
			log.ZError(ctx, "doWSPushMsg failed", err, "wsResp", wsResp)
		}
	case constant.KickOnlineMsg:
		//log.Warn(wsResp.OperationID, "kick...  logout")
		//w.kickOnline(wsResp)
		//w.Logout(ctx)
	case constant.GetNewestSeq:
		fallthrough
	case constant.PullMsgBySeqList:
		fallthrough
	case constant.SendMsg:
		fallthrough
	case constant.LogoutMsg:
		fallthrough
	case constant.SendSignalMsg:
		fallthrough
	case constant.SetBackgroundStatus:
		if err := c.Syncer.NotifyResp(ctx, wsResp); err != nil {
			log.ZError(ctx, "notifyResp failed", err, "wsResp", wsResp)
		}
	default:
		log.Error(wsResp.OperationID, "type failed, ", wsResp.ReqIdentifier)
		return
	}
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
func (c *LongConnMgr) reConn(ctx context.Context) error {
	if c.IsConnected() {
		return nil
	}
	c.listener.OnConnecting()
	c.w.Lock()
	c.connStatus = Connecting
	c.w.Unlock()
	url := fmt.Sprintf("%s?sendID=%s&token=%s&platformID=%d&operationID=%s", ccontext.Info(ctx).WsAddr(),
		ccontext.Info(ctx).UserID(), ccontext.Info(ctx).Token(), ccontext.Info(ctx).Platform(), ccontext.Info(ctx).OperationID())
	var header http.Header
	if c.IsCompression {
		header = http.Header{"compression": []string{"gzip"}}
	}
	_, err := c.conn.Dial(url, header)
	if err != nil {
		//if httpResp != nil {
		//	errMsg := httpResp.Header.Get("ws_err_msg") + " operationID " + ctx.Value("operationID").(string) + err.Error()
		//	//log.Error(operationID, "websocket.DefaultDialer.Dial failed ", errMsg, httpResp.StatusCode)
		//	u.listener.OnConnectFailed(int32(httpResp.StatusCode), errMsg)
		//	switch int32(httpResp.StatusCode) {
		//	case constant.ErrTokenExpired.ErrCode:
		//		u.listener.OnUserTokenExpired()
		//		u.tokenErrCode = constant.ErrTokenExpired.ErrCode
		//		return false, false, utils.Wrap(err, errMsg)
		//	case constant.ErrTokenInvalid.ErrCode:
		//		u.tokenErrCode = constant.ErrTokenInvalid.ErrCode
		//		return false, false, utils.Wrap(err, errMsg)
		//	case constant.ErrTokenMalformed.ErrCode:
		//		u.tokenErrCode = constant.ErrTokenMalformed.ErrCode
		//		return false, false, utils.Wrap(err, errMsg)
		//	case constant.ErrTokenNotValidYet.ErrCode:
		//		u.tokenErrCode = constant.ErrTokenNotValidYet.ErrCode
		//		return false, false, utils.Wrap(err, errMsg)
		//	case constant.ErrTokenUnknown.ErrCode:
		//		u.tokenErrCode = constant.ErrTokenUnknown.ErrCode
		//		return false, false, utils.Wrap(err, errMsg)
		//	case constant.ErrTokenDifferentPlatformID.ErrCode:
		//		u.tokenErrCode = constant.ErrTokenDifferentPlatformID.ErrCode
		//		return false, false, utils.Wrap(err, errMsg)
		//	case constant.ErrTokenDifferentUserID.ErrCode:
		//		u.tokenErrCode = constant.ErrTokenDifferentUserID.ErrCode
		//		return false, false, utils.Wrap(err, errMsg)
		//	case constant.ErrTokenKicked.ErrCode:
		//		u.tokenErrCode = constant.ErrTokenKicked.ErrCode
		//		//if u.loginStatus != constant.Logout {
		//		//	u.listener.OnKickedOffline()
		//		//	u.SetLoginStatus(constant.Logout)
		//		//}
		//
		//		return false, true, utils.Wrap(err, errMsg)
		//	default:
		//		//errMsg = err.Error() + " operationID " + operationID
		//		errMsg = err.Error() + " operationID " + ctx.Value("operationID").(string)
		//		u.listener.OnConnectFailed(1001, errMsg)
		//		return true, false, utils.Wrap(err, errMsg)
		//	}
		//} else {
		//	errMsg := err.Error() + " operationID " + ctx.Value("operationID").(string)
		//	u.listener.OnConnectFailed(1001, errMsg)
		//	if u.ConversationCh != nil {
		//		common.TriggerCmdSuperGroupMsgCome(sdk_struct.CmdNewMsgComeToConversation{MsgList: nil, OperationID: ctx.Value("operationID").(string), SyncFlag: constant.MsgSyncBegin}, u.ConversationCh)
		//		common.TriggerCmdSuperGroupMsgCome(sdk_struct.CmdNewMsgComeToConversation{MsgList: nil, OperationID: ctx.Value("operationID").(string), SyncFlag: constant.MsgSyncFailed}, u.ConversationCh)
		//	}
		//
		//	//log.Error(operationID, "websocket.DefaultDialer.Dial failed ", errMsg, "url ", url)
		//	return true, false, utils.Wrap(err, errMsg)
		//}
		c.listener.OnConnectFailed(1001, err.Error())
		c.w.Lock()
		c.connStatus = Closed
		c.w.Unlock()
		return err
	}
	c.listener.OnConnectSuccess()
	c.ctx = newContext(c.conn.LocalAddr())
	c.ctx = context.WithValue(ctx, "ConnContext", c.ctx)
	c.w.Lock()
	c.connStatus = Connected
	c.w.Unlock()
	_ = common.TriggerCmdConnected(ctx, c.pushMsgAndMaxSeqCh)
	return nil
}
func (c *LongConnMgr) doPushMsg(ctx context.Context, wsResp GeneralWsResp) error {
	var msg sdkws.MsgData
	err := proto.Unmarshal(wsResp.Data, &msg)
	if err != nil {
		return err
	}
	return common.TriggerCmdPushMsg(ctx, &msg, c.pushMsgAndMaxSeqCh)
}
