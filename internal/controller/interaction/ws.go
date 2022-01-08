package interaction

import (
	"errors"
	"github.com/gorilla/websocket"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"time"
)

type Ws struct {
	*WsRespAsyn
	*WsConn
}

func NewWs(wsRespAsyn *WsRespAsyn, wsConn *WsConn) *Ws {
	return &Ws{WsRespAsyn: wsRespAsyn, WsConn: wsConn}
}

func (ws *Ws) WaitResp(ch chan GeneralWsResp, timeout int, operationID string, connSend *websocket.Conn) (*GeneralWsResp, error) {
	select {
	case r := <-ch:
		log.Info(operationID, "ws ch recvMsg success: ")
		if r.ErrCode != 0 {
			return nil, constant.WsRecvCode
		} else {
			log.Info(operationID, "ws ch recvMsg success")
			return &r, nil
		}

	case <-time.After(time.Second * time.Duration(timeout)):
		log.Error(operationID, "ws ch recvMsg err, timeout")
		if connSend != ws.conn {
			return nil, constant.WsRecvConnDiff
		} else {
			return nil, constant.WsRecvConnSame
		}
	}
}

func (ws *Ws) SendReqWaitResp(buff []byte, reqIdentifier int32, timeout int, SenderID string) (*GeneralWsResp, error, string) {
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
		return nil, err, wsReq.OperationID
	}
	r1, r2 := ws.WaitResp(ch, timeout, wsReq.OperationID, connSend)
	return r1, r2, wsReq.OperationID
}
