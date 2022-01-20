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
	cmdCh chan common.Cmd2Value //waiting logout cmd
}

func NewHeartbeat(msgSync *MsgSync, cmcCh chan common.Cmd2Value) *Heartbeat {
	p := Heartbeat{MsgSync: msgSync, cmdCh: cmcCh}
	go p.Run()
	return &p
}

func (u *Heartbeat) Run() {
	heartbeatInterval := 5
	reqTimeout := 30
	retryTimes := 0

	for {
		operationID := utils.OperationIDGenerator()

		select {
		case r := <-u.cmdCh:
			if r.Cmd == constant.CmdLogout {
				return
			}
			log.Warn(operationID, "other cmd ...", r.Cmd)
		case <-time.After(time.Millisecond * time.Duration(heartbeatInterval*1000)):
			log.Debug(operationID, "heartbeat waiting(ms)... ", heartbeatInterval*1000)
		}

		resp, err := u.SendReqWaitResp(&server_api_params.GetMaxAndMinSeqReq{}, constant.WSGetNewestSeq, reqTimeout, retryTimes, u.loginUserID, operationID)
		if err != nil {
			log.Error(operationID, "SendReqWaitResp failed ", err.Error(), constant.WSGetNewestSeq, reqTimeout, u.loginUserID)
			if !errors.Is(err, constant.WsRecvConnSame) && !errors.Is(err, constant.WsRecvConnDiff) {
				u.CloseConn()
			}
			continue
		}

		var wsSeqResp server_api_params.GetMaxAndMinSeqResp
		err = proto.Unmarshal(resp.Data, &wsSeqResp)
		if err != nil {
			log.Error(operationID, "Unmarshal failed ", err.Error())
			u.CloseConn()
			continue
		}

		err = common.TriggerCmdMaxSeq(uint32(wsSeqResp.MaxSeq), u.PushMsgAndMaxSeqCh)
		if err != nil {
			log.Error(operationID, "TriggerMaxSeq failed ", err.Error())
		}
	}
}
