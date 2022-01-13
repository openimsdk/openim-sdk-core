package interaction

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"strings"
	"sync"
	"time"
)

type ConnListener interface {
	OnConnecting()
	OnConnectSuccess()
	OnConnectFailed(ErrCode int32, ErrMsg string)
	OnKickedOffline()
	OnUserTokenExpired()
	OnSelfInfoUpdated(userInfo string)
}

const writeTimeoutSeconds = 5

type WsConn struct {
	stateMutex  sync.Mutex
	conn        *websocket.Conn
	loginState  int32
	listener    ConnListener
	token       string
	loginUserID string
}

func NewWsConn(listener ConnListener, token string, loginUserID string) *WsConn {
	p := WsConn{listener: listener, token: token, loginUserID: loginUserID}
	p.conn, _ = p.ReConn()
	return &p
}

func (u *WsConn) CloseConn() error {
	u.Lock()
	defer u.Unlock()
	if u.conn != nil {
		return u.conn.Close()
	}
	return nil
}

func (u *WsConn) LoginState() int32 {
	return u.loginState
}

func (u *WsConn) SetLoginState(loginState int32) {
	u.loginState = loginState
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

func (u *WsConn) writeBinaryMsg(msg GeneralWsReq) (error, *websocket.Conn) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(msg)
	if err != nil {
		return utils.Wrap(err, ""), nil
	}

	var connSended *websocket.Conn
	u.stateMutex.Lock()
	defer u.stateMutex.Unlock()

	if u.conn != nil {
		connSended = u.conn
		err := u.SetWriteTimeout(writeTimeoutSeconds)
		if err != nil {
			return utils.Wrap(err, "SetWriteTimeout"), nil
		}
		if len(buff.Bytes()) > constant.MaxTotalMsgLen {
			return utils.Wrap(errors.New("msg too long"), ""), nil
		}
		err = u.conn.WriteMessage(websocket.BinaryMessage, buff.Bytes())
		return utils.Wrap(err, "WriteMessage failed"), connSended
	} else {
		return utils.Wrap(errors.New("conn==nil"), ""), connSended
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

func (u *WsConn) ReConn() (*websocket.Conn, error) {
	u.stateMutex.Lock()
	defer u.stateMutex.Unlock()
	if u.conn != nil {
		u.conn.Close()
		u.conn = nil
	}
	if u.loginState == constant.TokenFailedKickedOffline {
		return nil, utils.Wrap(errors.New("don't re conn"), "TokenFailedKickedOffline")
	}

	u.listener.OnConnecting()
	url := fmt.Sprintf("%s?sendID=%s&token=%s&platformID=%d", sdk_struct.SvrConf.WsAddr, u.loginUserID, u.token, sdk_struct.SvrConf.Platform)
	conn, httpResp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		u.loginState = constant.LoginFailed
		if httpResp != nil {
			u.listener.OnConnectFailed(int32(httpResp.StatusCode), err.Error())
		} else {
			u.listener.OnConnectFailed(1001, err.Error())
		}
		return nil, err
	}
	u.listener.OnConnectSuccess()
	u.loginState = constant.LoginSuccess
	return conn, nil
}
