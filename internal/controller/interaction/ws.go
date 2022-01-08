package interaction

import (
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/utils"
	"github.com/gorilla/websocket"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"sync"
	"time"
	"google.golang.org/protobuf/proto"
)

type Ws struct {
	*WsRespAsyn
	*WsConn
	seqMsg      map[int32]server_api_params.MsgData
	seqMsgMutex *sync.RWMutex
	*db.DataBase
	conversationCh chan common.Cmd2Value
}

func NewWs(wsRespAsyn *WsRespAsyn, wsConn *WsConn, lock *sync.RWMutex, ch chan common.Cmd2Value) *Ws {
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

func (u *Ws) run() {
	for {
		utils.LogStart()
		if u.conn == nil {
			utils.LogBegin("reConn", nil)
			re, _, _ := u.ReConn(nil)
			utils.LogEnd("reConn", re)
			u.conn = re
		}
		if u.conn != nil {
			msgType, message, err := u.conn.ReadMessage()
			log.Error("ReadMessage message ", msgType, err)
			if err != nil {
				u.stateMutex.Lock()
				utils.sdkLog("ws read message failed ", err.Error(), u.LoginState)
				if u.LoginState == constant.LogoutCmd {
					utils.sdkLog("logout, ws close, return ", constant.LogoutCmd, err)
					u.conn = nil
					u.stateMutex.Unlock()
					return
				}
				u.stateMutex.Unlock()
				time.Sleep(time.Duration(5) * time.Second)
				utils.sdkLog("ws  ReadMessage failed, sleep 5s, reconn, ", err)
				utils.LogBegin("reConn", u.conn)
				u.conn, _, err = u.ReConn(u.conn)
				utils.LogEnd("reConn", u.conn)
			} else {
				if msgType == websocket.CloseMessage {
					u.conn, _, _ = u.ReConn(u.conn)
				} else if msgType == websocket.TextMessage {
					utils.sdkLog("type failed, recv websocket.TextMessage ", string(message))
				} else if msgType == websocket.BinaryMessage {
					go u.doWsMsg(message)
				} else {
					utils.sdkLog("recv other msg: type ", msgType)
				}
			}
		} else {
			u.stateMutex.Lock()
			if u.LoginState == constant.LogoutCmd {
				utils.sdkLog("logout, ws close, return ", constant.LogoutCmd)
				u.stateMutex.Unlock()
				return
			}
			u.stateMutex.Unlock()
			utils.sdkLog("ws failed, sleep 5s, reconn... ")
			time.Sleep(time.Duration(5) * time.Second)
		}
	}
}

func (u *Ws) doWsMsg(message []byte) {
	utils.LogBegin()
	utils.LogBegin("decodeBinaryWs")
	wsResp, err := u.decodeBinaryWs(message)
	if err != nil {
		utils.LogFReturn("decodeBinaryWs err", err.Error())
		return
	}
	utils.LogEnd("decodeBinaryWs ", wsResp.OperationID, wsResp.ReqIdentifier)

	switch wsResp.ReqIdentifier {
	case constant.WSGetNewestSeq:
		u.doWSGetNewestSeq(*wsResp)
	case constant.WSPullMsgBySeqList:
		u.doWSPullMsg(*wsResp)
	case constant.WSPushMsg:
		u.doWSPushMsg(*wsResp)
	case constant.WSSendMsg:
		u.doWSSendMsg(*wsResp)
	case constant.WSKickOnlineMsg:
		u.kickOnline(*wsResp)
	default:
		utils.LogFReturn("type failed, ", wsResp.ReqIdentifier, wsResp.OperationID, wsResp.ErrCode, wsResp.ErrMsg)
		return
	}
	utils.LogSReturn()
	return
}

func (u *Ws) doWSGetNewestSeq(wsResp GeneralWsResp) {
	utils.LogBegin(wsResp.OperationID)
	u.notifyResp(wsResp)
	utils.LogSReturn(wsResp.OperationID)
}

func (u *Ws) doWSPullMsg(wsResp GeneralWsResp) {
	utils.LogBegin(wsResp.OperationID)
	u.notifyResp(wsResp)
	utils.LogSReturn(wsResp.OperationID)
}

func (u *Ws) doWSSendMsg(wsResp GeneralWsResp) {
	utils.LogBegin(wsResp.OperationID)
	u.notifyResp(wsResp)
	utils.LogSReturn(wsResp.OperationID)
}

func (u *Ws) doWSPushMsg(wsResp GeneralWsResp) {
	utils.LogBegin()
	u.doMsg(wsResp)
	utils.LogSReturn()
}

func (u *Ws) kickOnline(msg GeneralWsResp) {
	u.listener.OnKickedOffline()
}

func (u *Ws) doMsg(wsResp GeneralWsResp) {
	var msg server_api_params.MsgData
	if wsResp.ErrCode != 0 {
		utils.sdkLog("errcode: ", wsResp.ErrCode, " errmsg: ", wsResp.ErrMsg)
		utils.LogFReturn()
		return
	}
	err := proto.Unmarshal(wsResp.Data, &msg)
	if err != nil {
		utils.sdkLog("Unmarshal failed", err.Error())
		utils.LogFReturn()
		return
	}

	utils.sdkLog("openim ws  recv push msg do push seq in : ", msg.Seq)
	u.seqMsgMutex.Lock()
	b1 := u.isExistsInErrChatLogBySeq(msg.Seq)
	b2 := u.judgeMessageIfExists(msg.ClientMsgID)
	_, ok := u.seqMsg[int32(msg.Seq)]
	if b1 || b2 || ok {
		utils.sdkLog("seq in : ", msg.Seq, b1, b2, ok)
		u.seqMsgMutex.Unlock()
		return
	}

	u.seqMsg[int32(msg.Seq)] = &msg
	u.seqMsgMutex.Unlock()

	arrMsg := utils.ArrMsg{}
	common.TriggerCmdNewMsgCome(arrMsg, u.conversationCh)
}
