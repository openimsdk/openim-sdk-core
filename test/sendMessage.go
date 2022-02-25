package test

import (
	"github.com/gorilla/websocket"
	"open_im_sdk/internal/interaction"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"time"
)

func SendTextMessage(text, senderID, recvID, operationID string, ws *interaction.Ws) bool {
	ws.
	timeout := 300
	retryTimes := 60
	var wsReq GeneralWsReq
	var connSend *websocket.Conn
	var err error
	wsReq.ReqIdentifier = constant.WSSendMsg
	wsReq.OperationID = operationID
	msgIncr, ch := w.AddCh(senderID)
	defer w.DelCh(msgIncr)
	wsReq.SendID = senderID
	wsReq.MsgIncr = msgIncr
	wsReq.Data, err = proto.Marshal(m)
	if err != nil {
		return nil, utils.Wrap(err, "proto marshal err")
	}
	for i := 0; i < retryTimes+1; i++ {
		connSend, err = w.writeBinaryMsg(wsReq)
		if err != nil {
			if !w.IsWriteTimeout(err) {
				log.Error(operationID, "Not send timeout, failed, close conn, writeBinaryMsg again ", err.Error())
				w.CloseConn()
				time.Sleep(time.Duration(1) * time.Second)
				continue
			} else {
				return nil, utils.Wrap(err, "writeBinaryMsg timeout")
			}
		}
		break
	}
	r1, r2 := w.WaitResp(ch, timeout, wsReq.OperationID, connSend)
	return r1, r2
}
