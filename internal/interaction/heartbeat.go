package interaction

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"open_im_sdk/open_im_sdk_callback"
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
	//	*Ws
	*MsgSync
	cmdCh             chan common.Cmd2Value //waiting logout cmd , wake up cmd
	heartbeatInterval int
	token             string
	listener          open_im_sdk_callback.OnConnListener
	ExpireTimeSeconds uint32
}

func (u *Heartbeat) SetHeartbeatInterval(heartbeatInterval int) {
	u.heartbeatInterval = heartbeatInterval
}

func NewHeartbeat(msgSync *MsgSync, cmcCh chan common.Cmd2Value, listener open_im_sdk_callback.OnConnListener, token string, expireTimeSeconds uint32) *Heartbeat {
	p := Heartbeat{MsgSync: msgSync, cmdCh: cmcCh}
	p.heartbeatInterval = constant.HeartbeatInterval
	p.listener = listener
	p.token = token
	p.ExpireTimeSeconds = expireTimeSeconds
	go p.Run()
	return &p
}

type ParseToken struct {
	UID      string `json:"UID"`
	Platform string `json:"Platform"`
	Exp      int    `json:"exp"`
	Nbf      int    `json:"nbf"`
	Iat      int    `json:"iat"`
}

func (u *Heartbeat) IsTokenExp(operationID string) bool {
	if u.ExpireTimeSeconds == 0 {
		return false
	}
	log.Debug(operationID, "ExpireTimeSeconds ", u.ExpireTimeSeconds, "now ", uint32(time.Now().Unix()))
	if u.ExpireTimeSeconds < uint32(time.Now().Unix()) {
		return true
	} else {
		return false
	}
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
		if u.IsTokenExp(operationID) {
			log.Warn(operationID, "TokenExp, close heartbeat channel, call OnUserTokenExpired ,set logout", u.cmdCh)
			u.listener.OnUserTokenExpired()
			u.SetLoginState(constant.Logout)
			u.CloseConn()
			runtime.Goexit()
		}
		groupIDList, err := u.GetJoinedSuperGroupIDList()
		if err != nil {
			log.Error(operationID, "GetJoinedSuperGroupIDList failed ", err.Error())
		}
		log.Debug(operationID, "GetJoinedSuperGroupIDList ", groupIDList)
		resp, err := u.SendReqWaitResp(&server_api_params.GetMaxAndMinSeqReq{UserID: u.loginUserID, GroupIDList: groupIDList}, constant.WSGetNewestSeq, reqTimeout, retryTimes, u.loginUserID, operationID)
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

		log.Debug(operationID, "recv heartbeat resp, max seq on svr: ", wsSeqResp.MaxSeq, wsSeqResp.GroupMaxAndMinSeq)
		groupID2MaxSeqOnSvr := make(map[string]uint32, 0)
		for groupID, seq := range wsSeqResp.GroupMaxAndMinSeq {
			groupID2MaxSeqOnSvr[groupID] = seq.MaxSeq
		}
		for {
			err = common.TriggerCmdMaxSeq(sdk_struct.CmdMaxSeqToMsgSync{OperationID: operationID, MaxSeqOnSvr: wsSeqResp.MaxSeq, GroupID2MaxSeqOnSvr: groupID2MaxSeqOnSvr}, u.PushMsgAndMaxSeqCh)
			if err != nil {
				log.Error(operationID, "TriggerMaxSeq failed ", err.Error(), " MaxSeq ", wsSeqResp.MaxSeq)
				continue
			} else {
				log.Debug(operationID, "TriggerMaxSeq  success ", " MaxSeq ", wsSeqResp.MaxSeq)
				break
			}
		}

	}
}
