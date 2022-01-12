package init

import (
	"errors"
	"github.com/golang/protobuf/proto"
	ws "open_im_sdk/internal/controller/interaction"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"time"
)

type Heartbeat struct {
	*ws.Ws
	*MsgSync
}

func NewHeartbeat(ws *ws.Ws, msgSync *MsgSync) *Heartbeat {
	return &Heartbeat{Ws: ws, MsgSync: msgSync}
}

func (u *Heartbeat) Run() {
	heartbeatInterval := 5
	for {
		u.Lock()
		if u.LoginState() == constant.Logout {
			u.Unlock()
			return
		}
		u.Unlock()

		timeout := 30
		resp, err, operationID := u.SendReqWaitResp(nil, constant.WSGetNewestSeq, timeout, u.loginUserID)
		if err != nil {
			log.Error(operationID, "SendReqWaitResp failed ", err.Error(), constant.WSGetNewestSeq, timeout, u.loginUserID)
			if errors.Is(err, constant.WsRecvCode) {
				u.CloseConn()
				continue
			}
			if errors.Is(err, constant.WsRecvConnSame) {
				for tr := 0; tr < 3; tr++ {
					err = u.SendPingMsg()
					if err != nil {
						log.Error("sendPingMsg failed ", operationID, err.Error(), tr)
						time.Sleep(time.Duration(timeout) * time.Second)
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
