package interaction

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"time"
)

type Heartbeat struct {
	//*Ws
	*MsgSync
}

func NewHeartbeat(msgSync *MsgSync) *Heartbeat {
	p := Heartbeat{MsgSync: msgSync}
	go p.Run()
	return &p
}

func (u *Heartbeat) Run() {
	heartbeatInterval := 5
	reqTimeout := 30
	reTryInterval := 10
	for {
		u.Lock()
		if u.LoginState() == constant.Logout {
			u.Unlock()
			return
		}
		u.Unlock()

		resp, err, operationID := u.SendReqWaitResp(nil, constant.WSGetNewestSeq, reqTimeout, u.loginUserID)
		if err != nil {
			log.Error(operationID, "SendReqWaitResp failed ", err.Error(), constant.WSGetNewestSeq, reqTimeout, u.loginUserID)
			//	if  u.IsWriteTimeout(err)
			if errors.Is(err, constant.WsRecvCode) {
				log.Error(operationID, "is WsRecvCode, CloseConn")
				u.CloseConn()
				time.Sleep(time.Duration(reTryInterval) * time.Second)
				continue
			}
			if errors.Is(err, constant.WsRecvConnSame) {
				for tr := 0; tr < 3; tr++ {
					err = u.SendPingMsg()
					if err != nil {
						log.Error("sendPingMsg failed ", operationID, err.Error(), tr)
						time.Sleep(time.Duration(reTryInterval) * time.Second)
					} else {
						break
					}
				}
				continue
			}
			if errors.Is(err, constant.WsRecvConnDiff) {
				continue
			}
		}
		var wsSeqResp server_api_params.GetMaxAndMinSeqResp
		err = proto.Unmarshal(resp.Data, &wsSeqResp)
		if err != nil {
			log.Error(operationID, "Unmarshal failed ", err.Error())
			u.CloseConn()
		} else {
			needSyncSeq := u.getNeedSyncSeq(int32(wsSeqResp.MinSeq), int32(wsSeqResp.MaxSeq))
			log.Info("needSyncSeq ", wsSeqResp.MinSeq, wsSeqResp.MaxSeq, needSyncSeq)
			u.syncMsgFromServer(needSyncSeq)
		}
		time.Sleep(time.Duration(heartbeatInterval) * time.Second)
	}
}
