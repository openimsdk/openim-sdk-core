package init

import (
	ws "open_im_sdk/internal/controller/interaction"
	"open_im_sdk/internal/open_im_sdk"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"time"
)

type Heartbeat struct {
	wsConn      *ws.WsConn
	wsRespAsyn  *ws.WsRespAsyn
	token       string
	loginUserID string
}

func (u *Heartbeat) heartbeat() {
	for {
		u.wsConn.Lock()
		if u.wsConn.LoginState() == constant.LogoutCmd {
			u.wsConn.Unlock()
			return
		}
		u.wsConn.Unlock()

		msgIncr, ch := u.wsRespAsyn.AddCh(u.loginUserID)

		var wsReq ws.GeneralWsReq

		wsReq.ReqIdentifier = constant.WSGetNewestSeq
		wsReq.OperationID = utils.OperationIDGenerator()
		wsReq.SendID = u.loginUserID
		wsReq.MsgIncr = msgIncr
		//var connSend *websocket.Conn
		//	LogBegin("WriteMsg", wsReq.OperationID, wsReq.MsgIncr)
		err, connSend := u.wsConn.WriteMsg(wsReq)
		//	LogEnd("WriteMsg", wsReq.OperationID, wsReq.MsgIncr)
		if err != nil {
			utils.LogBegin("closeConn DelCh", msgIncr, wsReq.OperationID)
			u.wsConn.CloseConn()
			u.wsRespAsyn.DelCh(msgIncr)
			utils.LogEnd("closeConn DelCh continue", wsReq.OperationID)
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}

		timeout := 30
		breakFlag := 0
		for {
			if breakFlag == 1 {
				log.Info("break ", wsReq.OperationID)
				break
			}
			select {
			case r := <-ch:
				log.Info("ws ch recvMsg success: ", wsReq.OperationID)
				if r.ErrCode != 0 {
					log.Info("heartbeat response failed ", r.ErrCode, r.ErrMsg, wsReq.OperationID)
					u.wsConn.CloseConn()
				} else {
					//		sdkLog("heartbeat response success ", wsReq.OperationID)
					var wsSeqResp server_api_params.GetMaxAndMinSeqResp
					err = proto.Unmarshal(r.Data, &wsSeqResp)
					if err != nil {
						utils.sdkLog("Unmarshal failed, ", err.Error(), wsReq.OperationID)
						u.closeConn()
						//	u.DelCh(msgIncr)
						utils.LogEnd("closeConn DelCh continue")
					} else {
						needSyncSeq := u.getNeedSyncSeq(int32(wsSeqResp.MinSeq), int32(wsSeqResp.MaxSeq))
						utils.sdkLog("needSyncSeq ", wsSeqResp.MinSeq, wsSeqResp.MaxSeq, needSyncSeq)
						u.syncMsgFromServer(needSyncSeq)
					}
				}
				breakFlag = 1

			case <-time.After(time.Second * time.Duration(timeout)):
				var flag bool
				utils.sdkLog("ws ch recvMsg err: ", wsReq.OperationID)
				if connSend != u.conn {
					utils.sdkLog("old conn != current conn  ", connSend, u.conn)
					flag = false // error
				} else {
					flag = false //error
					for tr := 0; tr < 3; tr++ {
						err = u.sendPingMsg()
						if err != nil {
							utils.sdkLog("sendPingMsg failed ", wsReq.OperationID, err.Error(), tr)
							time.Sleep(time.Duration(30) * time.Second)
						} else {
							utils.sdkLog("sendPingMsg ok, break", wsReq.OperationID)
							flag = true //wait continue
							breakFlag = 1
							break
						}
					}
				}
				if breakFlag == 1 {
					utils.sdkLog("don't wait ", wsReq.OperationID)
					break
				}
				if flag == false {
					utils.sdkLog("ws ch recvMsg timeout ", timeout, "s ", wsReq.OperationID)
					utils.LogBegin("closeConn", wsReq.OperationID)
					u.closeConn()
					utils.LogEnd("closeConn", wsReq.OperationID)
					breakFlag = 1
					break
				} else {
					utils.sdkLog("wait resp continue", wsReq.OperationID)
					breakFlag = 0
					continue
				}
			}
		}

		u.DelCh(msgIncr)
		//	LogEnd("DelCh", wsReq.OperationID)
		time.Sleep(time.Duration(open_im_sdk.hearbeatInterval) * time.Second)
	}
}
