package interaction

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"runtime"
	"time"
)

type Heartbeat struct {
	//*Ws
	*MsgSync
	cmdCh             chan common.Cmd2Value //waiting logout cmd , wake up cmd
	heartbeatInterval int
}

func (u *Heartbeat) SetHeartbeatInterval(heartbeatInterval int) {
	u.heartbeatInterval = heartbeatInterval
}

func NewHeartbeat(msgSync *MsgSync, cmcCh chan common.Cmd2Value) *Heartbeat {
	p := Heartbeat{MsgSync: msgSync, cmdCh: cmcCh}
	p.heartbeatInterval = 30
	go p.Run()
	return &p
}

func (u *Heartbeat) Run() {
	//	heartbeatInterval := 30
	reqTimeout := 30
	retryTimes := 0
	heartbeatNum := 0
	for {
		operationID := utils.OperationIDGenerator()
		if heartbeatNum != 0 {
			select {
			case r := <-u.cmdCh:
				if r.Cmd == constant.CmdLogout {
					log.Warn(operationID, "recv logout cmd, close conn,  set logout state, Goexit...")
					u.SetLoginState(constant.Logout)
					u.CloseConn()
					log.Warn(operationID, "close heartbeat channel ", u.cmdCh)
					//	close(u.cmdCh)
					runtime.Goexit()
				}
				if r.Cmd == constant.CmdWakeUp {
					log.Info(operationID, "recv wake up cmd, start heartbeat ", r.Cmd)
					break
				}

				log.Warn(operationID, "other cmd...", r.Cmd)
			case <-time.After(time.Millisecond * time.Duration(u.heartbeatInterval*1000)):
				log.Debug(operationID, "heartbeat waiting(ms)... ", u.heartbeatInterval*1000)
			}
		}

		heartbeatNum++
		log.Debug(operationID, "send heartbeat req")
		resp, err := u.SendReqWaitResp(&server_api_params.GetMaxAndMinSeqReq{}, constant.WSGetNewestSeq, reqTimeout, retryTimes, u.loginUserID, operationID)
		if err != nil {
			log.Error(operationID, "SendReqWaitResp failed ", err.Error(), constant.WSGetNewestSeq, reqTimeout, u.loginUserID)
			if !errors.Is(err, constant.WsRecvConnSame) && !errors.Is(err, constant.WsRecvConnDiff) {
				log.Error(operationID, "other err,  close conn", err.Error())
				u.CloseConn()
			}
			continue
		}

		var wsSeqResp server_api_params.GetMaxAndMinSeqResp
		err = proto.Unmarshal(resp.Data, &wsSeqResp)
		if err != nil {
			log.Error(operationID, "Unmarshal failed, close conn", err.Error())
			u.CloseConn()
			continue
		}

		log.Debug(operationID, "recv heartbeat resp, max seq on svr: ", wsSeqResp.MaxSeq)

		err = common.TriggerCmdMaxSeq(sdk_struct.CmdMaxSeqToMsgSync{OperationID: operationID, MaxSeqOnSvr: wsSeqResp.MaxSeq}, u.PushMsgAndMaxSeqCh)
		if err != nil {
			log.Error(operationID, "TriggerMaxSeq failed ", err.Error())
		}
	}
}
