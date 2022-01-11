package interaction

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"

	"time"

	"open_im_sdk/pkg/utils"

	"errors"
	"open_im_sdk/pkg/constant"
	"sync"
)

type ConnListener interface {
	OnConnecting()
	OnConnectSuccess()
	OnConnectFailed(ErrCode int32, ErrMsg string)
	OnKickedOffline()
	OnUserTokenExpired()
	OnSelfInfoUpdated(userInfo string)
}

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
	p.conn, _, _ = p.ReConn()
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
		return errors.New("conn == nil")
	}
	var ping string = "try ping"
	err := u.conn.SetWriteDeadline(time.Now().Add(8 * time.Second))
	if err != nil {

	}
	return u.conn.WriteMessage(websocket.PingMessage, []byte(ping))
}

func (u *WsConn) SetWriteTimeout(timeout uint32)error{
	return u.conn.SetWriteDeadline(time.Now().Add(timeout * time.Second))
}


func (u *WsConn) SetReadTimeout(timeout int)error{
	return u.conn.SetReadDeadline(time.Now().Add(timeout * time.Second))
}


func (u *WsConn) writeBinaryMsg(msg GeneralWsReq) (error, *websocket.Conn) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(msg)
	if err != nil {
		return err, nil
	}

	var connSended *websocket.Conn
	u.stateMutex.Lock()
	defer u.stateMutex.Unlock()

	if u.conn != nil {
		connSended = u.conn

		err := u.SetWriteTimeout(8)
		if err != nil {
		}
		if len(buff.Bytes()) > constant.MaxTotalMsgLen {
			return errors.New("msg too long"), connSended
		}
		err = u.conn.WriteMessage(websocket.BinaryMessage, buff.Bytes())
		if err != nil {
		} else {
		}
		return err, connSended
	} else {
		return errors.New("conn==nil"), connSended
	}
}

func (u *WsConn) decodeBinaryWs(message []byte) (*GeneralWsResp, error) {

	buff := bytes.NewBuffer(message)
	dec := gob.NewDecoder(buff)
	var data GeneralWsResp
	err := dec.Decode(&data)
	if err != nil {

		return nil, err
	}

	return &data, nil
}

func (u *WsConn) WriteMsg(msg GeneralWsReq) (error, *websocket.Conn) {
	return u.writeBinaryMsg(msg)
}

func (u *WsConn) IsReadTimeout() bool{
	return false
}

func (u *WsConn) IsWriteTimeout() bool{
	return false
}

func (u *WsConn) IsFatalError() bool{
	return false
}


func (u *WsConn) ReConn() (*websocket.Conn, *http.Response, error) {
	u.stateMutex.Lock()
	defer u.stateMutex.Unlock()
	if u.conn != nil {
		u.conn.Close()
		u.conn = nil
	}
	if u.loginState == constant.TokenFailedKickedOffline || u.loginState == constant.TokenFailedExpired || u.loginState == constant.TokenFailedInvalid {
		return nil, nil, errors.New("don't reconn")
	}

	u.listener.OnConnecting()
	url := fmt.Sprintf("%s?sendID=%s&token=%s&platformID=%d", constant.SvrConf.WsAddr, u.loginUserID, u.token, constant.SvrConf.Platform)
	conn, httpResp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		if httpResp != nil {
			u.listener.OnConnectFailed(int32(httpResp.StatusCode), err.Error())
		} else {
			u.listener.OnConnectFailed(1001, err.Error())
		}

		utils.LogFReturn(nil, err.Error(), url)
		return nil, httpResp, err
	}
	u.listener.OnConnectSuccess()
	u.loginState = constant.LoginSuccess

	return conn, httpResp, nil
}
