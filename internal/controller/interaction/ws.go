package interaction

import (
	"errors"
	"github.com/gorilla/websocket"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"time"
)

type Ws struct {
	WsRespAsyn
	WsConn
}

func (ws *Ws) WaitResp(ch chan GeneralWsResp, timeout int, operationID string, connSend *websocket.Conn) (*GeneralWsResp, error) {
	select {
	case r := <-ch:
		log.Info(operationID, "ws ch recvMsg success: ")
		if r.ErrCode != 0 {
			return nil, errors.New("errCode failed")
		} else {
			log.Info(operationID, "ws ch recvMsg success")
			return &r, nil
		}

	case <-time.After(time.Second * time.Duration(timeout)):
		log.Error(operationID, "ws ch recvMsg err, timeout")
		if connSend != ws.conn {
			return nil, errors.New("recv timeout, conn diff")
		} else {
			return nil, errors.New("recv timeout, conn same")
		}
	}
}

func (ws *Ws) SendReqWaitResp(buff []byte, reqIdentifier int32, timeout int, SenderID string) (*GeneralWsResp, error) {
	var wsReq GeneralWsReq
	wsReq.ReqIdentifier = reqIdentifier
	wsReq.OperationID = utils.OperationIDGenerator()
	ws.Lock()
	msgIncr, ch := ws.AddCh(SenderID)
	ws.Unlock()
	wsReq.SendID = SenderID
	wsReq.MsgIncr = msgIncr
	wsReq.Data = buff
	err, connSend := ws.writeBinaryMsg(wsReq)
	if err != nil {
		log.Error(wsReq.OperationID, "ws send err ", err.Error(), wsReq)
		return nil, err
	}
	return ws.WaitResp(ch, timeout, wsReq.OperationID, connSend)
}
