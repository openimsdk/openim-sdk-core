package interaction

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"runtime"

	"time"
)

type Ws struct {
	*WsRespAsyn
	*WsConn
	//*db.DataBase
	//conversationCh chan common.Cmd2Value
	cmdCh              chan common.Cmd2Value //waiting logout cmd
	pushMsgAndMaxSeqCh chan common.Cmd2Value //recv push msg  -> channel
}

func NewWs(wsRespAsyn *WsRespAsyn, wsConn *WsConn, cmdCh chan common.Cmd2Value, pushMsgAndMaxSeqCh chan common.Cmd2Value) *Ws {
	p := Ws{WsRespAsyn: wsRespAsyn, WsConn: wsConn, cmdCh: cmdCh, pushMsgAndMaxSeqCh: pushMsgAndMaxSeqCh}
	go p.ReadData()
	return &p
}

//func (w *Ws) SeqMsg() map[int32]server_api_params.MsgData {
//	w.seqMsgMutex.RLock()
//	defer w.seqMsgMutex.RUnlock()
//	return w.seqMsg
//}
//
//func (w *Ws) SetSeqMsg(seqMsg map[int32]server_api_params.MsgData) {
//	w.seqMsgMutex.Lock()
//	defer w.seqMsgMutex.Unlock()
//	w.seqMsg = seqMsg
//}

func (w *Ws) WaitResp(ch chan GeneralWsResp, timeout int, operationID string, connSend *websocket.Conn) (*GeneralWsResp, error) {
	select {
	case r := <-ch:
		log.Debug(operationID, "ws ch recvMsg success, code ", r.ErrCode)
		if r.ErrCode != 0 {
			return nil, constant.WsRecvCode
		} else {
			return &r, nil
		}

	case <-time.After(time.Second * time.Duration(timeout)):
		log.Error(operationID, "ws ch recvMsg err, timeout")
		if connSend != w.WsConn.conn {
			return nil, constant.WsRecvConnDiff
		} else {
			return nil, constant.WsRecvConnSame
		}
	}
}

func (w *Ws) SendReqWaitResp(m proto.Message, reqIdentifier int32, timeout, retryTimes int, senderID, operationID string) (*GeneralWsResp, error) {
	var wsReq GeneralWsReq
	var connSend *websocket.Conn
	var err error
	wsReq.ReqIdentifier = reqIdentifier
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

func (w *Ws) reConnSleep(operationID string, sleep int32) {
	_, err := w.WsConn.ReConn()
	if err != nil {
		log.Error(operationID, "ReConn failed ", err.Error())
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}

func (w *Ws) ReadData() {
	isErrorOccurred := false
	for {
		operationID := utils.OperationIDGenerator()
		if isErrorOccurred {
			select {
			case r := <-w.cmdCh:
				if r.Cmd == constant.CmdLogout {
					log.Info(operationID, "recv CmdLogout, return, close conn")
					w.SetLoginState(constant.Logout)
					//		w.CloseConn()
					return
				}
				log.Warn(operationID, "other cmd ...", r.Cmd)
			case <-time.After(time.Microsecond * time.Duration(100)):
				log.Warn(operationID, "timeout(ms)... ", 100)
			}
		}
		isErrorOccurred = false
		if w.WsConn.conn == nil {
			log.Error(operationID, "conn == nil, ReConn")
			w.reConnSleep(operationID, 1)
			continue
		}

		//	timeout := 5
		//	u.WsConn.SetReadTimeout(timeout)
		msgType, message, err := w.WsConn.conn.ReadMessage()
		if err != nil {
			isErrorOccurred = true
			if w.loginState == constant.Logout {
				log.Warn(operationID, "loginState == logout ")
				continue
			}
			if w.WsConn.IsFatalError(err) {
				log.Error(operationID, "IsFatalError ", err.Error(), "ReConn")
				w.reConnSleep(operationID, 5)
			} else {
				log.Warn(operationID, "timeout failed ", err.Error())
			}
			continue
		}
		if msgType == websocket.CloseMessage {
			log.Error(operationID, "type websocket.CloseMessage, ReConn")
			w.reConnSleep(operationID, 1)
			continue
		} else if msgType == websocket.TextMessage {
			log.Warn(operationID, "type websocket.TextMessage")
		} else if msgType == websocket.BinaryMessage {
			go w.doWsMsg(message)
		} else {
			log.Warn(operationID, "recv other type ", msgType)
		}
	}
}

func (w *Ws) doWsMsg(message []byte) {
	wsResp, err := w.decodeBinaryWs(message)
	if err != nil {
		log.Error("decodeBinaryWs err", err.Error())
		return
	}
	log.Debug(wsResp.OperationID, "ws recv msg, code: ", wsResp.ErrCode, wsResp.ReqIdentifier)
	switch wsResp.ReqIdentifier {
	case constant.WSGetNewestSeq:
		if err = w.doWSGetNewestSeq(*wsResp); err != nil {
			log.Error(wsResp.OperationID, "doWSGetNewestSeq failed ", err.Error())
		}
	case constant.WSPullMsgBySeqList:
		if err = w.doWSPullMsg(*wsResp); err != nil {
			log.Error(wsResp.OperationID, "doWSPullMsg failed ", err.Error())
		}
	case constant.WSPushMsg:
		if err = w.doWSPushMsg(*wsResp); err != nil {
			log.Error(wsResp.OperationID, "doWSPushMsg failed ", err.Error())
		}
	case constant.WSSendMsg:
		if err = w.doWSSendMsg(*wsResp); err != nil {
			log.Error(wsResp.OperationID, "doWSSendMsg failed ", err.Error())
		}
	case constant.WSKickOnlineMsg:
		log.Warn(wsResp.OperationID, "kick... ")
		w.kickOnline(*wsResp)
	case constant.WsLogoutMsg:
		log.Warn(wsResp.OperationID, "logout... ")
		w.SetLoginState(constant.Logout)
		w.CloseConn()
		runtime.Goexit()
	default:
		log.Error(wsResp.OperationID, "type failed, ", wsResp.ReqIdentifier)
		return
	}
}

func (w *Ws) doWSGetNewestSeq(wsResp GeneralWsResp) error {
	if err := w.notifyResp(wsResp); err != nil {
		return utils.Wrap(err, "")
	}
	return nil
}

func (w *Ws) doWSPullMsg(wsResp GeneralWsResp) error {
	if err := w.notifyResp(wsResp); err != nil {
		return utils.Wrap(err, "")
	}
	return nil
}

func (w *Ws) doWSSendMsg(wsResp GeneralWsResp) error {
	if err := w.notifyResp(wsResp); err != nil {
		return utils.Wrap(err, "")
	}
	return nil
}

func (w *Ws) doWSPushMsg(wsResp GeneralWsResp) error {
	if wsResp.ErrCode != 0 {
		return utils.Wrap(errors.New("errCode"), wsResp.ErrMsg)
	}
	var msg server_api_params.MsgData
	err := proto.Unmarshal(wsResp.Data, &msg)
	if err != nil {
		return utils.Wrap(err, "Unmarshal failed")
	}
	return utils.Wrap(common.TriggerCmdPushMsg(sdk_struct.CmdPushMsgToMsgSync{Msg: &msg, OperationID: wsResp.OperationID}, w.pushMsgAndMaxSeqCh), "")
}

func (w *Ws) kickOnline(msg GeneralWsResp) {
	w.listener.OnKickedOffline()
}

//func (u *Ws) doSendMsg(wsResp GeneralWsResp) error {
//
//}
