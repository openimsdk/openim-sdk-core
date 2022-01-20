package interaction

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
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
	retryTimes := 0
	reTryInterval := 10
	operationID := utils.OperationIDGenerator()
	for {
		time.Sleep(time.Duration(heartbeatInterval) * time.Second)
		u.Lock()
		if u.LoginState() == constant.Logout {
			u.Unlock()
			return
		}
		u.Unlock()

		resp, err := u.SendReqWaitResp(&server_api_params.GetMaxAndMinSeqReq{}, constant.WSGetNewestSeq, reqTimeout, retryTimes, u.loginUserID, operationID)
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
			} else {
				log.Error(operationID, "other err ", err.Error(), " closeConn")
				u.CloseConn()
				continue
			}
		}
		var wsSeqResp server_api_params.GetMaxAndMinSeqResp
		err = proto.Unmarshal(resp.Data, &wsSeqResp)
		if err != nil {
			log.Error(operationID, "Unmarshal failed ", err.Error())
			u.CloseConn()
		} else {
			err := common.TriggerCmdMaxSeq(uint32(wsSeqResp.MaxSeq), u.PushMsgAndMaxSeqCh)
			if err != nil {
				log.Error(operationID, "TriggerMaxSeq failed ", err.Error())
			}

		}
	}
}
