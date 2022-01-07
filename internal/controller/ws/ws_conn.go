package ws

import (
	"bytes"
	"encoding/gob"
	"github.com/gorilla/websocket"
	"open_im_sdk/internal/open_im_sdk"
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

func (u *WsConn) decodeBinaryWs(message []byte) (*open_im_sdk.GeneralWsResp, error) {
	LogStart()
	buff := bytes.NewBuffer(message)
	dec := gob.NewDecoder(buff)
	var data GeneralWsResp
	err := dec.Decode(&data)
	if err != nil {
		LogFReturn(nil, err.Error())
		return nil, err
	}
	LogSReturn(&data, nil)
	return &data, nil
}

func (u *WsConn) WriteMsg(msg open_im_sdk.GeneralWsReq) (error, *websocket.Conn) {
	LogStart(msg.OperationID)
	LogSReturn(msg.OperationID)
	return u.writeBinaryMsg(msg)
}
