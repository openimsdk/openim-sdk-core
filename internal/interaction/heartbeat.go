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
	//*Ws
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

	//
	//
	//b := strings.IndexAny(u.token, ".")
	//e := strings.LastIndex(u.token, ".")
	//if b == -1 || e == -1 || b >= e {
	//	return false
	//}
	//log.Debug(operationID, "sub token ", u.token[b+1:e])
	//decodeBytes, err := base64.StdEncoding.DecodeString(u.token[b+1 : e])
	//if err != nil {
	//	//	log.Error(operationID, "DecodeString failed ", err.Error(), u.token[b+1:e])
	//	return false
	//}
	//log.Debug(operationID, "decodeBytes ", string(decodeBytes))
	//parseToken := ParseToken{}
	//err = json.Unmarshal(decodeBytes, &parseToken)
	//if err != nil {
	//	log.Error(operationID, "Unmarshal failed ", err.Error())
	//	return false
	//}
	//log.Debug(operationID, "exp ", parseToken.Exp, "now ", time.Now().Unix())
	//if parseToken.Exp < int(time.Now().Unix()) {
	//	return true
	//}
	//return false
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
