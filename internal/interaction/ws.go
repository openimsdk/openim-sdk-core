package interaction

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"runtime"
	"sync"
	"time"
)

type Ws struct {
	*WsRespAsyn
	*WsConn
	seqMsg      map[int32]server_api_params.MsgData
	seqMsgMutex sync.RWMutex
	*db.DataBase
	conversationCh chan common.Cmd2Value
	cmdCh          chan common.Cmd2Value
}

func NewWs(wsRespAsyn *WsRespAsyn, wsConn *WsConn, conversationCh, cmdCh chan common.Cmd2Value) *Ws {
	p := Ws{WsRespAsyn: wsRespAsyn, WsConn: wsConn, conversationCh: conversationCh, cmdCh: cmdCh}
	go p.ReadData()
	return &p
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

func (ws *Ws) WaitResp(ch chan GeneralWsResp, timeout int, operationID string, connSend *websocket.Conn) (*GeneralWsResp, error) {
	select {
	case r := <-ch:
		log.Debug(operationID, "ws ch recvMsg success, code ", r.ErrCode, r.ErrMsg)
		if r.ErrCode != 0 {
			return nil, constant.WsRecvCode
		} else {
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

func (ws *Ws) SendReqWaitResp(m proto.Message, reqIdentifier int32, timeout, retryTimes int, senderID, operationID string) (*GeneralWsResp, error) {
	var wsReq GeneralWsReq
	var connSend *websocket.Conn
	var err error
	wsReq.ReqIdentifier = reqIdentifier
	wsReq.OperationID = operationID
	msgIncr, ch := ws.AddCh(senderID)
	defer ws.DelCh(msgIncr)
	wsReq.SendID = senderID
	wsReq.MsgIncr = msgIncr
	wsReq.Data, err = proto.Marshal(m)
	if err != nil {
		return nil, utils.Wrap(err, "proto marshal err")
	}
	for i := 0; i < retryTimes+1; i++ {
		err, connSend = ws.writeBinaryMsg(wsReq)
		if err != nil {
			if !ws.IsWriteTimeout(err) {
				newErr := connSend.Close()
				log.Error(operationID, m, "ws write Timeout", newErr, err.Error())
				time.Sleep(time.Duration(1) * time.Second)
				continue
			} else {
				return nil, utils.Wrap(err, "writeBinaryMsg err")
			}
		} else {
			break
		}
	}
	r1, r2 := ws.WaitResp(ch, timeout, wsReq.OperationID, connSend)
	return r1, r2
}

func (u *Ws) ReadData() {
	for {
		isErrorOccurred := false
		operationID := utils.OperationIDGenerator()
		if u.WsConn.conn != nil {
			//	timeout := 5
			//	u.WsConn.SetReadTimeout(timeout)
			msgType, message, err := u.WsConn.conn.ReadMessage()
			if err != nil {
				isErrorOccurred = true
				if u.WsConn.IsFatalError(err) {
					log.Error(operationID, "IsFatalError ", err.Error(), "ReConn")
					c, err := u.WsConn.ReConn()
					if err != nil {
						log.Error(operationID, "reconn failed ", c, err.Error())
						time.Sleep(time.Duration(2) * time.Second)
					}
				} else {
					log.Warn(operationID, "other err  ", err.Error())
				}
			} else {
				if msgType == websocket.CloseMessage {
					log.Error(operationID, "type websocket.CloseMessage, ReConn")
					c, err := u.WsConn.ReConn()
					if err != nil {
						log.Error(operationID, "reconn failed ", c, err.Error())
						time.Sleep(time.Duration(2) * time.Second)
					}
				} else if msgType == websocket.TextMessage {
					log.Warn(operationID, "type websocket.TextMessage")
				} else if msgType == websocket.BinaryMessage {
					u.doWsMsg(message)
				} else {
					log.Warn(operationID, "recv other type ", msgType)
				}
			}
		} else {
			log.Error(operationID, "conn == nil, ReConn")
			_, err := u.WsConn.ReConn()
			if err != nil {
				isErrorOccurred = true
				log.Error(operationID, "ReConn failed ", err.Error())
			}
		}

		if isErrorOccurred {
			select {
			case r := <-u.cmdCh:
				if r.Cmd == constant.CmdLogout {
					u.SetLoginState(constant.Logout)
					return
				} else {
					log.Warn(operationID, "other cmd ...", r.Cmd)
					break
				}
			case <-time.After(time.Microsecond * time.Duration(1000)):
				log.Warn(operationID, "timeout... ", 1000)
				break
			}
		}
	}
}

func (u *Ws) doWsMsg(message []byte) {
	wsResp, err := u.decodeBinaryWs(message)
	if err != nil {
		log.Error("decodeBinaryWs err", err.Error())
		return
	}
	switch wsResp.ReqIdentifier {
	case constant.WSGetNewestSeq:
		go u.doWSGetNewestSeq(*wsResp)
	case constant.WSPullMsgBySeqList:
		go u.doWSPullMsg(*wsResp)
	case constant.WSPushMsg:
		go u.doWSPushMsg(*wsResp)
	case constant.WSSendMsg:
		go u.doWSSendMsg(*wsResp)
	case constant.WSKickOnlineMsg:
		go u.kickOnline(*wsResp)
	case constant.WsLogoutMsg:
		log.Warn(wsResp.OperationID, "logout.. ")
		u.SetLoginState(constant.Logout)
		runtime.Goexit()
	default:
		log.Error(wsResp.OperationID, "type failed, ", wsResp.ReqIdentifier, wsResp.OperationID)
		return
	}
}

func (u *Ws) doWSGetNewestSeq(wsResp GeneralWsResp) error {
	if err := u.notifyResp(wsResp); err != nil {
		log.Error(wsResp.OperationID, "doWSGetNewestSeq failed ", err.Error())
		return err
	}
	return nil
}

func (u *Ws) doWSPullMsg(wsResp GeneralWsResp) error {
	if err := u.notifyResp(wsResp); err != nil {
		log.Error(wsResp.OperationID, "doWSPullMsg failed ", err.Error())
		return err
	}
	return nil
}

func (u *Ws) doWSSendMsg(wsResp GeneralWsResp) error {
	if err := u.notifyResp(wsResp); err != nil {
		log.Error(wsResp.OperationID, "doWSSendMsg failed ", err.Error())
		return err
	}
	return nil
}

func (u *Ws) doWSPushMsg(wsResp GeneralWsResp) error {
	if err := u.doSendMsg(wsResp); err != nil {
		log.Error(wsResp.OperationID, "doWSPushMsg failed ", err.Error())
		return err
	}
	return nil
}

func (u *Ws) kickOnline(msg GeneralWsResp) {
	u.listener.OnKickedOffline()
}

func (u *Ws) doSendMsg(wsResp GeneralWsResp) error {
	if wsResp.ErrCode != 0 {
		return utils.Wrap(errors.New("errCode"), wsResp.ErrMsg)
	}
	var msg server_api_params.MsgData
	err := proto.Unmarshal(wsResp.Data, &msg)
	if err != nil {
		return utils.Wrap(err, "Unmarshal failed")
	}

	u.seqMsgMutex.Lock()
	defer u.seqMsgMutex.Unlock()
	b1 := u.IsExistsInErrChatLogBySeq(msg.Seq)
	b2, _ := u.MessageIfExists(msg.ClientMsgID)
	_, ok := u.seqMsg[int32(msg.Seq)]
	if b1 || b2 || ok {
		log.Debug("0", "seq in : ", msg.Seq, b1, b2, ok)
		return nil
	}
	u.seqMsg[int32(msg.Seq)] = msg
	arrMsg := sdk_struct.ArrMsg{}
	common.TriggerCmdNewMsgCome(arrMsg, u.conversationCh)
	return nil
}
