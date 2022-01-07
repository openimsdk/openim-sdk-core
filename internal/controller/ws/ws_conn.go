package ws

import (
	"bytes"
	"encoding/gob"
	"github.com/gorilla/websocket"

	"errors"
	"open_im_sdk/pkg/constant"
	"sync"
)

type WsConn struct {
	stateMutex sync.Mutex
	conn       *websocket.Conn
}

func (u *WsConn) sendPingMsg() error {
	u.stateMutex.Lock()
	defer u.stateMutex.Unlock()
	var ping string = "try ping"
	err := u.conn.SetWriteDeadline(time.Now().Add(8 * time.Second))
	if err != nil {
	}
	return u.conn.WriteMessage(websocket.PingMessage, []byte(ping))
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
		err = u.conn.SetWriteDeadline(time.Now().Add(8 * time.Second))
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
