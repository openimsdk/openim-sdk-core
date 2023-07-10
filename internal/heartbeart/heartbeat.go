package heartbeart

import (
	"errors"
	"open_im_sdk/internal/full"
	"open_im_sdk/internal/interaction"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"runtime"
	"time"

	"github.com/golang/protobuf/proto"
)

type Heartbeat struct {
	//	*Ws
	*interaction.MsgSync
	cmdCh             chan common.Cmd2Value //waiting logout cmd , wake up cmd
	heartbeatInterval int
	token             string
	listener          open_im_sdk_callback.OnConnListener
	//ExpireTimeSeconds uint32
	id2MinSeq          map[string]uint32
	full               *full.Full
	WsForTest          *interaction.Ws
	LoginUserIDForTest string
}

func (u *Heartbeat) SetHeartbeatInterval(heartbeatInterval int) {
	u.heartbeatInterval = heartbeatInterval
}

func NewHeartbeat(msgSync *interaction.MsgSync, cmcCh chan common.Cmd2Value, listener open_im_sdk_callback.OnConnListener, token string, id2MinSeq map[string]uint32, full *full.Full) *Heartbeat {
	p := Heartbeat{MsgSync: msgSync, cmdCh: cmcCh, full: full}
	p.heartbeatInterval = constant.HeartbeatInterval
	p.listener = listener
	p.token = token
	//p.ExpireTimeSeconds = expireTimeSeconds
	p.id2MinSeq = id2MinSeq
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

//func (u *Heartbeat) IsTokenExp(operationID string) bool {
//	if u.ExpireTimeSeconds == 0 {
//		return false
//	}
//	log.Debug(operationID, "ExpireTimeSeconds ", u.ExpireTimeSeconds, "now ", uint32(time.Now().Unix()))
//	if u.ExpireTimeSeconds < uint32(time.Now().Unix()) {
//		return true
//	} else {
//		return false
//	}
//}

func (u *Heartbeat) Run() {
	defaultTimeout := 5
	wakeUpTimeout := 1
	reqTimeout := defaultTimeout
	retryTimes := 0
	heartbeatNum := 0
	for {
		reqTimeout = defaultTimeout
		operationID := utils.OperationIDGenerator()
		if constant.OnlyForTest == 1 {
			time.Sleep(5 * time.Second)
			var groupIDList []string
			resp, err := u.WsForTest.SendReqWaitResp(&server_api_params.GetMaxAndMinSeqReq{UserID: u.LoginUserIDForTest, GroupIDList: groupIDList}, constant.WSGetNewestSeq, reqTimeout, retryTimes, u.LoginUserIDForTest, operationID)
			if err != nil {
				log.Error(operationID, "SendReqWaitResp failed ", err.Error(), constant.WSGetNewestSeq, reqTimeout, u.LoginUserIDForTest)
				if !errors.Is(err, constant.WsRecvConnSame) && !errors.Is(err, constant.WsRecvConnDiff) {
					log.Error(operationID, "other err,  close conn", err.Error())
					u.CloseConn(operationID)
				}
				continue
			}

			var wsSeqResp server_api_params.GetMaxAndMinSeqResp
			err = proto.Unmarshal(resp.Data, &wsSeqResp)
			if err != nil {
				log.Error(operationID, "Unmarshal failed, close conn", err.Error())
				u.CloseConn(operationID)
				continue
			}
			log.Debug(operationID, "heartbeat req -> resp ")
			continue
		}

		if heartbeatNum != 0 {
			select {
			case r := <-u.cmdCh:
				if r.Cmd == constant.CmdLogout {
					log.Warn(operationID, "recv logout cmd, close conn,  set logout state, Goexit...")
					u.SetLoginStatus(constant.Logout)
					u.CloseConn(operationID)
					log.Warn(operationID, "close heartbeat channel ", u.cmdCh)
					runtime.Goexit()
				}
				if r.Cmd == constant.CmdWakeUp {
					u.full.SuperGroup.SyncJoinedGroupList(operationID)
					u.full.Group().SyncJoinedGroupList(operationID)
					log.Info(operationID, "recv wake up cmd, start heartbeat ", r.Cmd)
					reqTimeout = wakeUpTimeout
					break
				}
				log.Warn(operationID, "other cmd...", r.Cmd)
			case <-time.After(time.Millisecond * time.Duration(u.heartbeatInterval*1000)):
				log.Debug(operationID, "heartbeat waiting(ms)... ", u.heartbeatInterval*1000)
			}
		}
		if u.LoginStatus() == constant.Logout {
			log.Warn(operationID, " logout state Goexit", u.cmdCh)
			runtime.Goexit()
		}
		heartbeatNum++
		log.Debug(operationID, "send heartbeat req")
		//if u.IsTokenExp(operationID) {
		//	log.Warn(operationID, "TokenExp, close heartbeat channel, call OnUserTokenExpired, set logout", u.cmdCh)
		//	u.listener.OnUserTokenExpired()
		//	u.SetLoginStatus(constant.Logout)
		//	u.CloseConn(operationID)
		//	runtime.Goexit()
		//}
		var groupIDList []string
		var err error
		if heartbeatNum == 1 {
			groupIDList, err = u.full.GetReadDiffusionGroupIDList(operationID)
			log.NewInfo(operationID, "full.GetReadDiffusionGroupIDList ", heartbeatNum)
		} else {
			groupIDList, err = u.GetReadDiffusionGroupIDList()
			log.NewInfo(operationID, "db.GetReadDiffusionGroupIDList ", heartbeatNum)
		}
		if err != nil {
			log.Error(operationID, "GetReadDiffusionGroupIDList failed ", err.Error())
		}
		log.Debug(operationID, "get GetJoinedSuperGroupIDList ", groupIDList)
		resp, err := u.SendReqWaitResp(&server_api_params.GetMaxAndMinSeqReq{UserID: u.LoginUserID, GroupIDList: groupIDList}, constant.WSGetNewestSeq, reqTimeout, retryTimes, u.LoginUserID, operationID)
		if err != nil {
			log.Error(operationID, "SendReqWaitResp failed ", err.Error(), constant.WSGetNewestSeq, reqTimeout, u.LoginUserID)
			//if !errors.Is(err, constant.WsRecvConnSame) && !errors.Is(err, constant.WsRecvConnDiff) {
			//	log.Error(operationID, "other err,  close conn", err.Error())
			u.CloseConn(operationID)
			//}
			continue
		}
		var wsSeqResp server_api_params.GetMaxAndMinSeqResp
		err = proto.Unmarshal(resp.Data, &wsSeqResp)
		if err != nil {
			log.Error(operationID, "Unmarshal failed, close conn", err.Error())
			u.CloseConn(operationID)
			continue
		}

		u.id2MinSeq[utils.GetUserIDForMinSeq(u.LoginUserID)] = wsSeqResp.MinSeq
		for g, v := range wsSeqResp.GroupMaxAndMinSeq {
			u.id2MinSeq[utils.GetGroupIDForMinSeq(g)] = v.MinSeq
		}
		if constant.OnlyForTest == 1 {
			continue
		}
		//server_api_params.MaxAndMinSeq
		log.Debug(operationID, "recv heartbeat resp,  seq on svr: ", wsSeqResp.MinSeq, wsSeqResp.MaxSeq, wsSeqResp.GroupMaxAndMinSeq)
		for {
			err = common.TriggerCmdMaxSeq(sdk_struct.CmdMaxSeqToMsgSync{OperationID: operationID, MaxSeqOnSvr: wsSeqResp.MaxSeq, MinSeqOnSvr: wsSeqResp.MinSeq, GroupID2MinMaxSeqOnSvr: wsSeqResp.GroupMaxAndMinSeq}, u.PushMsgAndMaxSeqCh)
			if err != nil {
				log.Error(operationID, "TriggerMaxSeq failed ", err.Error(), "seq ", wsSeqResp.MinSeq, wsSeqResp.MaxSeq, wsSeqResp.GroupMaxAndMinSeq)
				continue
			} else {
				log.Debug(operationID, "TriggerMaxSeq  success ", "seq ", wsSeqResp.MinSeq, wsSeqResp.MaxSeq, wsSeqResp.GroupMaxAndMinSeq)
				break
			}
		}
	}
}
