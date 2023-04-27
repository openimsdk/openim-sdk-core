// Copyright © 2023 OpenIM SDK.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package interaction

import (
	"context"
	"errors"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"runtime"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

type Ws struct {
	*WsRespAsyn
	*WsConn
	//conversationCh chan common.Cmd2Value
	cmdCh              chan common.Cmd2Value //waiting logout cmd
	pushMsgAndMaxSeqCh chan common.Cmd2Value //recv push msg  -> channel
	cmdHeartbeatCh     chan common.Cmd2Value //
	conversationCH     chan common.Cmd2Value
	JustOnceFlag       bool
	IsBackground       bool
}

func NewWs(wsRespAsyn *WsRespAsyn, wsConn *WsConn, cmdCh chan common.Cmd2Value, pushMsgAndMaxSeqCh chan common.Cmd2Value, cmdHeartbeatCh, conversationCH chan common.Cmd2Value) *Ws {
	p := Ws{WsRespAsyn: wsRespAsyn, WsConn: wsConn, cmdCh: cmdCh, pushMsgAndMaxSeqCh: pushMsgAndMaxSeqCh, cmdHeartbeatCh: cmdHeartbeatCh, conversationCH: conversationCH}
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

func (w *Ws) WaitResp(ch chan GeneralWsResp, timeout int, operationID string) (*GeneralWsResp, error) {
	select {
	case r := <-ch:
		log.Debug(operationID, "ws ch recvMsg success, code ", r.ErrCode)
		if r.ErrCode != 0 {
			log.Error(operationID, "ws ch recvMsg failed, code, err msg: ", r.ErrCode, r.ErrMsg)
			switch r.ErrCode {
			case int(constant.ErrInBlackList.ErrCode):
				return nil, &constant.ErrInBlackList
			case int(constant.ErrNotFriend.ErrCode):
				return nil, &constant.ErrNotFriend
			}
			return nil, errors.New(utils.IntToString(r.ErrCode) + ":" + r.ErrMsg)
		} else {
			return &r, nil
		}

	case <-time.After(time.Second * time.Duration(timeout)):
		log.Error(operationID, "ws ch recvMsg err, timeout")
		if w.conn.IsNil() {
			return nil, errors.New("ws ch recvMsg err, timeout,conn is nil")
		}
		if w.conn.CheckSendConnDiffNow() {
			return nil, constant.WsRecvConnDiff
		} else {
			return nil, constant.WsRecvConnSame
		}
	}
}

func (w *Ws) SendReqWaitResp(ctx context.Context, m proto.Message, reqIdentifier int32, timeout, retryTimes int, senderID string) (*GeneralWsResp, error) {
	switch reqIdentifier {
	case constant.WsSetBackgroundStatus:
		if v, ok := m.(*server_api_params.SetAppBackgroundStatusReq); ok {
			w.IsBackground = v.IsBackground
		}
	}
	var wsReq GeneralWsReq
	var err error
	wsReq.ReqIdentifier = reqIdentifier
	wsReq.OperationID = ctx.Value("operationID").(string)
	msgIncr, ch := w.AddCh(senderID)
	log.Debug(wsReq.OperationID, "SendReqWaitResp AddCh msgIncr:", msgIncr, reqIdentifier)
	defer w.DelCh(msgIncr)
	defer log.Debug(wsReq.OperationID, "SendReqWaitResp DelCh msgIncr:", msgIncr, reqIdentifier)
	wsReq.SendID = senderID
	wsReq.MsgIncr = msgIncr
	wsReq.Data, err = proto.Marshal(m)
	if err != nil {
		return nil, utils.Wrap(err, "proto marshal err")
	}
	flag := 0
	for i := 0; i < retryTimes+1; i++ {
		err = w.writeBinaryMsg(wsReq)
		if err != nil {
			if !w.IsWriteTimeout(err) {
				//log.Error(operationID, "Not send timeout, failed, close conn, writeBinaryMsg again ", err.Error(), w.conn, reqIdentifier)
				w.CloseConn(ctx)
				time.Sleep(time.Duration(1) * time.Second)
				continue
			} else {
				return nil, utils.Wrap(err, "writeBinaryMsg timeout")
			}
		}
		flag = 1
		break
	}
	if flag == 1 {
		//log.Debug(operationID, "send ok wait resp")
		r1, r2 := w.WaitResp(ch, timeout, wsReq.OperationID)
		return r1, r2
	} else {
		//log.Error(operationID, "send failed")
		err := errors.New("send failed")
		return nil, utils.Wrap(err, "SendReqWaitResp failed")
	}
}
func (w *Ws) SendReqTest(m proto.Message, reqIdentifier int32, timeout int, senderID, operationID string) bool {
	var wsReq GeneralWsReq

	var err error
	wsReq.ReqIdentifier = reqIdentifier
	wsReq.OperationID = operationID
	msgIncr, ch := w.AddCh(senderID)
	defer w.DelCh(msgIncr)
	wsReq.SendID = senderID
	wsReq.MsgIncr = msgIncr
	wsReq.Data, err = proto.Marshal(m)
	if err != nil {
		return false
	}
	err = w.writeBinaryMsg(wsReq)
	if err != nil {
		log.Error(operationID, "writeBinaryMsg timeout", m.String(), senderID, err.Error())
		return false
	} else {
		log.Debug(operationID, "writeBinaryMsg success", m.String(), senderID)
	}
	startTime := time.Now()
	result := w.WaitTest(ch, timeout, wsReq.OperationID, m, senderID)
	log.Debug(operationID, "ws Response time：", time.Since(startTime), m.String(), senderID, result)
	return result
}
func (w *Ws) WaitTest(ch chan GeneralWsResp, timeout int, operationID string, m proto.Message, senderID string) bool {
	select {
	case r := <-ch:
		if r.ErrCode != 0 {
			log.Error(operationID, "ws ch recvMsg success, code ", r.ErrCode, r.ErrMsg, m.String(), senderID)
			return false
		} else {
			log.Debug(operationID, "ws ch recvMsg send success, code ", m.String(), senderID)
			return true
		}

	case <-time.After(time.Second * time.Duration(timeout)):
		log.Error(operationID, "ws ch recvMsg err, timeout ", m.String(), senderID)

		return false
	}
}
func (w *Ws) reConnSleep(ctx context.Context, sleep int32) (error, bool) {
	isNeedReConn, isKicked, err := w.WsConn.ReConn(ctx)
	if err != nil {
		if isKicked {
			//log.Warn(operationID, "kicked, when re conn ")
			w.kickOnline(GeneralWsResp{})
			w.Logout(ctx)
		}
		//log.Error(operationID, "ReConn failed ", err.Error(), "is need re connect ", isNeedReConn)
		time.Sleep(time.Duration(sleep) * time.Second)
	} else {
		resp, err := w.SendReqWaitResp(ctx, &server_api_params.SetAppBackgroundStatusReq{UserID: w.loginUserID, IsBackground: w.IsBackground}, constant.WsSetBackgroundStatus, 5, 2, w.loginUserID)
		if err != nil {
			_ = resp
			//log.Error(operationID, "SendReqWaitResp failed ", err.Error(), constant.WsSetBackgroundStatus, 5, w.loginUserID, resp)
		}
		_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{Action: constant.SyncConversation, Args: ctx.Value("operationID")}, w.conversationCH)
	}
	return err, isNeedReConn
}

func (w *Ws) ReadData(ctx context.Context) {
	isErrorOccurred := false
	for {
		operationID := utils.OperationIDGenerator()
		ctx := context.WithValue(context.Background(), "operationID", operationID)
		if isErrorOccurred {
			select {
			case r := <-w.cmdCh:
				if r.Cmd == constant.CmdLogout {
					log.Info(operationID, "recv CmdLogout, return, close conn")
					log.Warn(operationID, "close ws read channel ", w.cmdCh)
					//		close(w.cmdCh)
					w.SetLoginStatus(constant.Logout)
					return
				}
				log.Warn(operationID, "other cmd ...", r.Cmd)
			case <-time.After(time.Millisecond * time.Duration(100)):
				log.Info(operationID, "timeout(ms)... ", 100)
			}
		}
		isErrorOccurred = false
		if w.WsConn.conn.IsNil() {
			isErrorOccurred = true
			log.Warn(operationID, "conn == nil, ReConn ")
			err, isNeedReConnect := w.reConnSleep(ctx, 1)
			if err != nil && isNeedReConnect == false {
				log.Warn(operationID, "token failed, don't connect again")
				return
			}
			continue
		}

		//	timeout := 5
		//	u.WsConn.SetReadTimeout(timeout)
		msgType, message, err := w.WsConn.conn.ReadMessage()
		if err != nil {
			isErrorOccurred = true
			if w.loginStatus == constant.Logout {
				log.Warn(operationID, "loginState == logout ")
				log.Warn(operationID, "close ws read channel ", w.cmdCh)
				//	close(w.cmdCh)
				return
			}
			if w.WsConn.IsFatalError(err) {
				log.Error(operationID, "IsFatalError ", err.Error(), "ReConn", w.WsConn.conn.LocalAddr())
				//sleep 500 millisecond,waiting for network reconn,when network switch
				time.Sleep(time.Millisecond * 500)
				err, isNeedReConnect := w.reConnSleep(ctx, 5)
				if err != nil && isNeedReConnect == false {
					log.Warn(operationID, "token failed, don't connect again ")
					return
				}
			} else {
				log.Warn(operationID, "timeout failed ", err.Error())
			}
			continue
		}
		if msgType == websocket.CloseMessage {
			log.Error(operationID, "type websocket.CloseMessage, ReConn")
			err, isNeedReConnect := w.reConnSleep(ctx, 1)

			if err != nil && isNeedReConnect == false {
				log.Warn(operationID, "token failed, don't connect again")
				return
			}
			continue
		} else if msgType == websocket.TextMessage {
			log.Warn(operationID, "type websocket.TextMessage")
		} else if msgType == websocket.BinaryMessage {
			if w.IsCompression {
				var decompressErr error
				message, decompressErr = w.compressor.DeCompress(message)
				if decompressErr != nil {
					log.NewWarn(operationID, "decompress failed", decompressErr.Error())
					continue
				}
			}
			w.doWsMsg(message)
		} else {
			log.Warn(operationID, "recv other type ", msgType)
		}
	}
}

func (w *Ws) doWsMsg(message []byte) {
	var wsResp GeneralWsResp
	err := w.encoder.Decode(message, &wsResp)
	if err != nil {
		log.Error("decodeBinaryWs err", err.Error())
		return
	}
	ctx := context.WithValue(context.Background(), "operationID", wsResp.OperationID)
	log.Debug(wsResp.OperationID, "ws recv msg, code: ", wsResp.ErrCode, wsResp.ReqIdentifier)
	switch wsResp.ReqIdentifier {
	case constant.WSGetNewestSeq:
		if err = w.doWSGetNewestSeq(wsResp); err != nil {
			log.Error(wsResp.OperationID, "doWSGetNewestSeq failed ", err.Error(), wsResp.ReqIdentifier, wsResp.MsgIncr)
		}
	case constant.WSPullMsgBySeqList:
		if err = w.doWSPullMsg(wsResp); err != nil {
			log.Error(wsResp.OperationID, "doWSPullMsg failed ", err.Error())
		}
	case constant.WSPushMsg:
		// todo
		//if constant.OnlyForTest == 1 {
		//	return
		//}
		if err = w.doWSPushMsg(wsResp); err != nil {
			log.Error(wsResp.OperationID, "doWSPushMsg failed ", err.Error())
		}
		//if err = w.doWSPushMsgForTest(*wsResp); err != nil {
		//	log.Error(wsResp.OperationID, "doWSPushMsgForTest failed ", err.Error())
		//}

	case constant.WSSendMsg:
		if err = w.doWSSendMsg(wsResp); err != nil {
			log.Error(wsResp.OperationID, "doWSSendMsg failed ", err.Error(), wsResp.ReqIdentifier, wsResp.MsgIncr)
		}
	case constant.WSKickOnlineMsg:
		log.Warn(wsResp.OperationID, "kick...  logout")
		w.kickOnline(wsResp)
		w.Logout(ctx)

	case constant.WsLogoutMsg:
		log.Warn(wsResp.OperationID, "WsLogoutMsg... Ws goroutine exit")
		if err = w.doWSLogoutMsg(wsResp); err != nil {
			log.Error(wsResp.OperationID, "doWSLogoutMsg failed ", err.Error())
		}
		runtime.Goexit()
	case constant.WSSendSignalMsg:
		log.Info(wsResp.OperationID, "signaling...")
		w.DoWSSignal(wsResp)
	case constant.WsSetBackgroundStatus:
		log.Info(wsResp.OperationID, "WsSetBackgroundStatus...")
		if err = w.setAppBackgroundStatus(wsResp); err != nil {
			log.Error(wsResp.OperationID, "WsSetBackgroundStatus failed ", err.Error(), wsResp.ReqIdentifier, wsResp.MsgIncr)
		}
		log.NewDebug(wsResp.OperationID, wsResp)
	default:
		log.Error(wsResp.OperationID, "type failed, ", wsResp.ReqIdentifier)
		return
	}
}

func (w *Ws) Logout(ctx context.Context) error {
	w.SetLoginStatus(constant.Logout)
	w.CloseConn(ctx)
	//log.Warn(operationID, "TriggerCmdLogout ws...", w.conn)
	err := common.TriggerCmdLogout(w.cmdCh)
	if err != nil {
		return err
		//log.Error(operationID, "TriggerCmdLogout failed ", err.Error())
	}
	//log.Info(operationID, "TriggerCmdLogout heartbeat...")
	err = common.TriggerCmdLogout(w.cmdHeartbeatCh)
	if err != nil {
		return err
		//log.Error(operationID, "TriggerCmdLogout failed ", err.Error())
	}
	return nil
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

func (w *Ws) DoWSSignal(wsResp GeneralWsResp) error {
	if err := w.notifyResp(wsResp); err != nil {
		return utils.Wrap(err, "")
	}
	return nil
}
func (w *Ws) doWSLogoutMsg(wsResp GeneralWsResp) error {
	if err := w.notifyResp(wsResp); err != nil {
		return utils.Wrap(err, "")
	}
	return nil
}
func (w *Ws) doWSPushMsg(wsResp GeneralWsResp) error {
	if wsResp.ErrCode != 0 {
		return utils.Wrap(errors.New("errCode"), wsResp.ErrMsg)
	}
	var msg sdkws.MsgData
	err := proto.Unmarshal(wsResp.Data, &msg)
	if err != nil {
		return utils.Wrap(err, "Unmarshal failed")
	}
	return utils.Wrap(common.TriggerCmdPushMsg(sdk_struct.CmdPushMsgToMsgSync{Msg: &msg, OperationID: wsResp.OperationID}, w.pushMsgAndMaxSeqCh), "")
}

func (w *Ws) setAppBackgroundStatus(wsResp GeneralWsResp) error {
	if err := w.notifyResp(wsResp); err != nil {
		return utils.Wrap(err, "")
	}
	return nil
}

func (w *Ws) doWSPushMsgForTest(wsResp GeneralWsResp) error {
	if wsResp.ErrCode != 0 {
		return utils.Wrap(errors.New("errCode"), wsResp.ErrMsg)
	}
	var msg server_api_params.MsgData
	err := proto.Unmarshal(wsResp.Data, &msg)
	if err != nil {
		return utils.Wrap(err, "Unmarshal failed")
	}
	log.Debug(wsResp.OperationID, "recv push doWSPushMsgForTest")
	return nil
	//	return utils.Wrap(common.TriggerCmdPushMsg(sdk_struct.CmdPushMsgToMsgSync{Msg: &msg, OperationID: wsResp.OperationID}, w.pushMsgAndMaxSeqCh), "")
}

func (w *Ws) kickOnline(msg GeneralWsResp) {
	w.listener.OnKickedOffline()
}

func (w *Ws) SendSignalingReqWaitResp(ctx context.Context, req *sdkws.SignalReq) (*sdkws.SignalResp, error) {
	resp, err := w.SendReqWaitResp(ctx, req, constant.WSSendSignalMsg, 10, 12, w.loginUserID)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var signalResp sdkws.SignalResp
	err = proto.Unmarshal(resp.Data, &signalResp)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	return &signalResp, nil
}

func (w *Ws) SignalingWaitPush(ctx context.Context, inviterUserID, inviteeUserID, roomID string, timeout int32) (*sdkws.SignalReq, error) {
	msgIncr := inviterUserID + inviteeUserID + roomID
	ch := w.AddChByIncr(msgIncr)
	defer w.DelCh(msgIncr)
	resp, err := w.WaitResp(ch, int(timeout), mcontext.GetOperationID(ctx))
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var signalReq sdkws.SignalReq
	err = proto.Unmarshal(resp.Data, &signalReq)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}

	return &signalReq, nil
}
