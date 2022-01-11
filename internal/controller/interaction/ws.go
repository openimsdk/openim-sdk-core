package interaction

import (
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"sync"
	"time"
)

type Ws struct {
	*WsRespAsyn
	*WsConn
	seqMsg      map[int32]server_api_params.MsgData
	seqMsgMutex *sync.RWMutex
	*db.DataBase
	conversationCh chan common.Cmd2Value
	cmdCh          chan common.Cmd2Value
}

func (ws *Ws) SeqMsg() map[int32]server_api_params.MsgData {
	ws.seqMsgMutex.RLock()
	defer ws.seqMsgMutex.RUnlock()
	return ws.seqMsg
}

func (ws *Ws) SetSeqMsg(seqMsg map[int32]server_api_params.MsgData) {
	ws.seqMsgMutex.Lock()
	defer ws.seqMsgMutex.Unlock()
	ws.seqMsg = seqMsg
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
		if connSend != ws.WsConn.conn {
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
		isErrorOccurred := false
		if u.WsConn.conn != nil {
			timeout := 5
			u.WsConn.SetReadTimeout(timeout)
			msgType, message, err := u.WsConn.conn.ReadMessage()
			if err != nil {
				isErrorOccurred = true
				if u.WsConn.IsFatalError() {
					log.Error("0", "fatal error, failed ", err.Error())
					u.WsConn.ReConn()
				} else {
					log.Warn("0", "other err  ", err.Error())
				}
			} else {
				if msgType == websocket.CloseMessage {
					u.WsConn.ReConn()
				} else if msgType == websocket.TextMessage {
					log.Warn("recv websocket.TextMessage type", string(message))
				} else if msgType == websocket.BinaryMessage {
					go u.doWsMsg(message)
				} else {
					log.Warn("recv other type", string(message), msgType)
				}
			}
		} else {
			_, _, err := u.WsConn.ReConn()
			if err != nil {
				isErrorOccurred = true
			}
		}

		if isErrorOccurred {
			select {
			case r := <-u.cmdCh:
				if r.Cmd == constant.CmdLogout {
					return
				} else {
					log.Warn("0", "other cmd ...", r.Cmd)
				}
			case <-time.After(time.Microsecond * time.Duration(1000)):
				log.Warn("0", "timeout... ", 1000)
			}
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

	u.seqMsgMutex.Lock()
	b1 := u.IsExistsInErrChatLogBySeq(msg.Seq)
	b2 := u.JudgeMessageIfExists(msg.ClientMsgID)
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
