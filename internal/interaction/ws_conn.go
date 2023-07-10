package interaction

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"strings"
	"sync"
	"time"
)

const writeTimeoutSeconds = 30

type WsConn struct {
	stateMutex     sync.Mutex
	conn           *websocket.Conn
	loginStatus    int32
	listener       open_im_sdk_callback.OnConnListener
	token          string
	loginUserID    string
	IsCompression  bool
	ConversationCh chan common.Cmd2Value
	tokenErrCode   int32
}

func (u *WsConn) IsInterruptReconnection() bool {
	if u.tokenErrCode != 0 {
		return true
	}
	return false
}

func NewWsConn(listener open_im_sdk_callback.OnConnListener, token string, loginUserID string, isCompression bool, conversationCh chan common.Cmd2Value) *WsConn {
	p := WsConn{listener: listener, token: token, loginUserID: loginUserID, IsCompression: isCompression, ConversationCh: conversationCh}
	//	go func() {
	p.conn, _, _, _ = p.ReConn("init:" + utils.OperationIDGenerator())
	//	}()
	return &p
}

func (u *WsConn) CloseConn(operationID string) error {
	u.Lock()
	defer u.Unlock()
	if u.conn != nil {
		err := u.conn.Close()
		if err != nil {
			log.NewWarn(operationID, "close conn, ", u.conn, err.Error())
		}
		//	u.conn = nil
		return utils.Wrap(err, "")
	}
	return nil
}

func (u *WsConn) LoginStatus() int32 {
	return u.loginStatus
}

func (u *WsConn) SetLoginStatus(loginState int32) {
	u.loginStatus = loginState
}

func (u *WsConn) Lock() {
	u.stateMutex.Lock()
}

func (u *WsConn) Unlock() {
	u.stateMutex.Unlock()
}

func (u *WsConn) SendPingMsg() error {
	u.stateMutex.Lock()
	defer u.stateMutex.Unlock()
	if u.conn == nil {
		return utils.Wrap(errors.New("conn == nil"), "")
	}
	ping := "try ping"
	err := u.SetWriteTimeout(writeTimeoutSeconds)
	if err != nil {
		return utils.Wrap(err, "SetWriteDeadline failed")
	}
	err = u.conn.WriteMessage(websocket.PingMessage, []byte(ping))
	if err != nil {
		return utils.Wrap(err, "WriteMessage failed")
	}
	return nil
}

func (u *WsConn) SetWriteTimeout(timeout int) error {
	return u.conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
}

func (u *WsConn) SetReadTimeout(timeout int) error {
	return u.conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
}

func (u *WsConn) writeBinaryMsg(msg GeneralWsReq) (*websocket.Conn, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(msg)
	if err != nil {
		return nil, utils.Wrap(err, "Encode error")
	}

	u.stateMutex.Lock()
	defer u.stateMutex.Unlock()
	if u.conn != nil {
		err := u.SetWriteTimeout(writeTimeoutSeconds)
		if err != nil {
			return nil, utils.Wrap(err, "SetWriteTimeout")
		}
		log.Debug("this msg length is :", float32(len(buff.Bytes()))/float32(1024), "kb")
		if len(buff.Bytes()) > constant.MaxTotalMsgLen {
			return nil, utils.Wrap(errors.New("msg too long"), utils.IntToString(len(buff.Bytes())))
		}
		var data []byte
		if u.IsCompression {
			var gzipBuffer bytes.Buffer
			gz := gzip.NewWriter(&gzipBuffer)
			if _, err := gz.Write(buff.Bytes()); err != nil {
				return nil, utils.Wrap(err, "")
			}
			if err := gz.Close(); err != nil {
				return nil, utils.Wrap(err, "")
			}
			data = gzipBuffer.Bytes()
		} else {
			data = buff.Bytes()
		}
		return u.conn, utils.Wrap(u.conn.WriteMessage(websocket.BinaryMessage, data), "")
	} else {
		return nil, utils.Wrap(errors.New("conn==nil"), "")
	}
}

func (u *WsConn) decodeBinaryWs(message []byte) (*GeneralWsResp, error) {
	buff := bytes.NewBuffer(message)
	dec := gob.NewDecoder(buff)
	var data GeneralWsResp
	err := dec.Decode(&data)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	return &data, nil
}

func (u *WsConn) IsReadTimeout(err error) bool {
	if strings.Contains(err.Error(), "timeout") {
		return true
	}
	return false
}

func (u *WsConn) IsWriteTimeout(err error) bool {
	if strings.Contains(err.Error(), "timeout") {
		return true
	}
	return false
}

func (u *WsConn) IsFatalError(err error) bool {
	if strings.Contains(err.Error(), "timeout") {
		return false
	}
	return true
}

func (u *WsConn) ReConn(operationID string) (*websocket.Conn, error, bool, bool) {
	u.stateMutex.Lock()
	u.tokenErrCode = 0
	defer u.stateMutex.Unlock()
	if u.conn != nil {
		log.NewInfo(operationID, "close old conn", u.conn.LocalAddr())
		err := u.conn.Close()
		if err != nil {
			log.NewWarn(operationID, "close old conn", u.conn.LocalAddr(), err.Error())
		}
		u.conn = nil
	}
	if u.loginStatus == constant.TokenFailedKickedOffline {
		return nil, utils.Wrap(errors.New("don't re conn"), "TokenFailedKickedOffline"), false, false
	}
	u.listener.OnConnecting()

	url := fmt.Sprintf("%s?sendID=%s&token=%s&platformID=%d&operationID=%s", sdk_struct.SvrConf.WsAddr, u.loginUserID, u.token, sdk_struct.SvrConf.Platform, operationID)
	log.Info(operationID, "ws connect begin, dail: ", url)
	var header http.Header
	if u.IsCompression {
		header = http.Header{"compression": []string{"gzip"}}
	}
	conn, httpResp, err := websocket.DefaultDialer.Dial(url, header)
	log.Info(operationID, "ws connect end, dail : ", url)
	if err != nil {
		log.Error(operationID, "ws connect failed ", url, err.Error())
		u.loginStatus = constant.LoginFailed
		if httpResp != nil {
			errMsg := httpResp.Header.Get("ws_err_msg") + " operationID " + operationID + err.Error()
			log.Error(operationID, "websocket.DefaultDialer.Dial failed ", errMsg, httpResp.StatusCode)
			u.listener.OnConnectFailed(int32(httpResp.StatusCode), errMsg)
			switch int32(httpResp.StatusCode) {
			case constant.ErrTokenExpired.ErrCode:
				u.listener.OnUserTokenExpired()
				u.tokenErrCode = constant.ErrTokenExpired.ErrCode
				return nil, utils.Wrap(err, errMsg), false, false
			case constant.ErrTokenInvalid.ErrCode:
				u.tokenErrCode = constant.ErrTokenInvalid.ErrCode
				return nil, utils.Wrap(err, errMsg), false, false
			case constant.ErrTokenMalformed.ErrCode:
				u.tokenErrCode = constant.ErrTokenMalformed.ErrCode
				return nil, utils.Wrap(err, errMsg), false, false
			case constant.ErrTokenNotValidYet.ErrCode:
				u.tokenErrCode = constant.ErrTokenNotValidYet.ErrCode
				return nil, utils.Wrap(err, errMsg), false, false
			case constant.ErrTokenUnknown.ErrCode:
				u.tokenErrCode = constant.ErrTokenUnknown.ErrCode
				return nil, utils.Wrap(err, errMsg), false, false
			case constant.ErrTokenDifferentPlatformID.ErrCode:
				u.tokenErrCode = constant.ErrTokenDifferentPlatformID.ErrCode
				return nil, utils.Wrap(err, errMsg), false, false
			case constant.ErrTokenDifferentUserID.ErrCode:
				u.tokenErrCode = constant.ErrTokenDifferentUserID.ErrCode
				return nil, utils.Wrap(err, errMsg), false, false
			case constant.ErrTokenKicked.ErrCode:
				u.tokenErrCode = constant.ErrTokenKicked.ErrCode
				//if u.loginStatus != constant.Logout {
				//	u.listener.OnKickedOffline()
				//	u.SetLoginStatus(constant.Logout)
				//}

				return nil, utils.Wrap(err, errMsg), false, true
			default:
				errMsg = err.Error() + " operationID " + operationID
				u.listener.OnConnectFailed(1001, errMsg)
				return nil, utils.Wrap(err, errMsg), true, false
			}
		} else {
			errMsg := err.Error() + " operationID " + operationID
			u.listener.OnConnectFailed(1001, errMsg)
			if u.ConversationCh != nil {
				common.TriggerCmdSuperGroupMsgCome(sdk_struct.CmdNewMsgComeToConversation{MsgList: nil, OperationID: operationID, SyncFlag: constant.MsgSyncBegin}, u.ConversationCh)
				common.TriggerCmdSuperGroupMsgCome(sdk_struct.CmdNewMsgComeToConversation{MsgList: nil, OperationID: operationID, SyncFlag: constant.MsgSyncFailed}, u.ConversationCh)
			}

			log.Error(operationID, "websocket.DefaultDialer.Dial failed ", errMsg, "url ", url)
			return nil, utils.Wrap(err, errMsg), true, false
		}
	}
	u.listener.OnConnectSuccess()
	u.loginStatus = constant.LoginSuccess
	u.conn = conn
	return conn, nil, true, false
}
