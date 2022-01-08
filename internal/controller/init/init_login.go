package init

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net/http"
	conv "open_im_sdk/internal/controller/conversation_msg"
	"open_im_sdk/internal/controller/friend"
	"open_im_sdk/internal/controller/group"
	ws "open_im_sdk/internal/controller/interaction"
	"open_im_sdk/internal/open_im_sdk"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"

	"time"

	"github.com/jinzhu/copier"
	"github.com/mitchellh/mapstructure"
	"open_im_sdk/internal/controller/user"
)

type LoginMgr struct {
	friend       *friend.Friend
	group        *group.Group
	conversation *conv.Conversation
	user         *user.User

	db      *db.DataBase
	ws      *ws.Ws
	msgSync *MsgSync

	hearbeat *Heartbeat

	token       string
	loginUserID string
	listener    ws.ConnListener
}

func (u *LoginMgr) login(userID, token string, cb common.Base) {
	if cb == nil {
		return
	}
	log.Info("login start ", userID, token)
	u.token = token
	u.loginUserID = userID

	db, err := db.NewDataBase(userID)
	if err != nil {
		cb.OnError(constant.ErrDB.ErrCode, constant.ErrDB.ErrMsg)
		log.Error("", "NewDataBase failed ", err.Error())
		return
	}
	u.db = db
	u.ws = ws.NewWs(ws.NewWsRespAsyn(), ws.NewWsConn(u.listener, token, userID))
	u.msgSync = NewMsgSync(db, u.ws, userID)

	u.hearbeat = NewHeartbeat(u.ws, u.msgSync)

	log.Info("ws, forcedSynchronization heartbeat coroutine run ...")
	go u.forcedSynchronization()
	go u.hearbeat.heartbeat()
	go u.run()
	//		go u.timedCloseDB()
	u.forycedSyncReceiveMessageOpt()
	cb.OnSuccess("")

}

//func (u *open_im_sdk.UserRelated) closeListenerCh() {
//	if u.ConversationCh != nil {
//		close(u.ConversationCh)
//		u.ConversationCh = nil
//	}
//}

func (u *LoginMgr) initSDK(config string, cb *ws.ConnListener) bool {
	if cb == nil {
		log.Error("callback == nil")
		return false
	}

	log.Info("initSDK LoginState ")

	u.listener = cb
	u.initListenerCh()
	utils.sdkLog("init success, ", config)

	go open_im_sdk.doListener(u)
	return true
}

func (u *open_im_sdk.UserRelated) unInitSDK() {
	u.unInitAll()
	u.closeListenerCh()
}

func (im *utils.open_im_sdk) getVersion() string {
	return "v1.0.5"
}

func (im *utils.open_im_sdk) getServerTime() int64 {
	return 0
}

func (u *open_im_sdk.UserRelated) logout(cb Base) {
	go func() {
		u.stateMutex.Lock()
		defer u.stateMutex.Unlock()

		u.LoginState = constant.LogoutCmd

		utils.sdkLog("set LoginState ", u.LoginState)

		err := u.closeConn()
		if err != nil {
			if cb != nil {
				cb.OnError(constant.ErrCodeInitLogin, err.Error())
			}
			return
		}
		utils.sdkLog("closeConn ok")

		//err = u.closeDB()
		if err != nil {
			if cb != nil {
				cb.OnError(constant.ErrCodeInitLogin, err.Error())
			}
			return
		}
		utils.sdkLog("close db ok")

		u.loginUserID = ""
		u.token = ""
		time.Sleep(time.Duration(6) * time.Second)
		if cb != nil {
			cb.OnSuccess("")
		}
		utils.sdkLog("logout return")
	}()
}

//
//func (u *open_im_sdk.UserRelated) timedCloseDB() {
//	timeTicker := time.NewTicker(time.Second * 5)
//	num := 0
//	for {
//		<-timeTicker.C
//		u.stateMutex.Lock()
//		if u.LoginState == constant.LogoutCmd {
//			utils.sdkLog("logout timedCloseDB return", constant.LogoutCmd)
//			u.stateMutex.Unlock()
//			return
//		}
//		u.stateMutex.Unlock()
//		num++
//		if num%60 == 0 {
//			utils.sdkLog("closeDBSetNil begin")
//			//u.closeDBSetNil()
//			utils.sdkLog("closeDBSetNil end")
//		}
//	}
//}
//
//func (u *open_im_sdk.UserRelated) closeConn() error {
//	utils.LogBegin()
//	if u.conn != nil {
//		err := u.conn.Close()
//		if err != nil {
//			utils.LogFReturn(err.Error())
//			return err
//		}
//	}
//	utils.LogSReturn(nil)
//	return nil
//}

func (u *open_im_sdk.UserRelated) getLoginUser() string {
	if u.LoginState == constant.LoginSuccess {
		return u.loginUserID
	} else {
		return ""
	}
}

func (im *utils.open_im_sdk) getLoginStatus() int {
	return im.LoginState
}

func (u *open_im_sdk.UserRelated) forycedSyncReceiveMessageOpt() {
	OperationID := utils.operationIDGenerator()
	resp, err := utils.post2ApiForRead(open_im_sdk.getAllConversationMessageOptRouter, open_im_sdk.paramGetAllConversationMessageOpt{OperationID: OperationID}, u.token)
	if err != nil {
		utils.sdkLog("post2Api failed, ", open_im_sdk.getAllConversationMessageOptRouter, OperationID)
		return
	}
	var v open_im_sdk.getReceiveMessageOptResp
	err = json.Unmarshal(resp, &v)
	if err != nil {
		utils.sdkLog("Unmarshal failed ", resp, OperationID)
		return
	}
	if v.ErrCode != 0 {
		utils.sdkLog("errCode failed, ", v.ErrCode, v.ErrMsg, string(resp), OperationID)
		return
	}

	utils.sdkLog("get receive opt ", v)
	u.receiveMessageOptMutex.Lock()
	for _, v := range v.Data {
		if v.Result != 0 {
			u.receiveMessageOpt[v.ConversationId] = v.Result
		}
	}
	u.receiveMessageOptMutex.Unlock()
}

func (u *open_im_sdk.UserRelated) forcedSynchronization() {
	u.ForceSyncFriend()
	u.ForceSyncBlackList()
	u.ForceSyncFriendApplication()
	//	u.ForceSyncSelfFriendApplication()
	u.ForceSyncLoginUserInfo()
	//	u.ForceSyncJoinedGroup()
	//	u.ForceSyncGroupRequest()
	//	u.ForceSyncSelfGroupRequest()
	//	u.ForceSyncJoinedGroupMember()
}

func (u *open_im_sdk.UserRelated) doWsMsg(message []byte) {
	utils.LogBegin()
	utils.LogBegin("decodeBinaryWs")
	wsResp, err := u.decodeBinaryWs(message)
	if err != nil {
		utils.LogFReturn("decodeBinaryWs err", err.Error())
		return
	}
	utils.LogEnd("decodeBinaryWs ", wsResp.OperationID, wsResp.ReqIdentifier)

	switch wsResp.ReqIdentifier {
	case constant.WSGetNewestSeq:
		u.doWSGetNewestSeq(*wsResp)
	case constant.WSPullMsgBySeqList:
		u.doWSPullMsg(*wsResp)
	case constant.WSPushMsg:
		u.doWSPushMsg(*wsResp)
	case constant.WSSendMsg:
		u.doWSSendMsg(*wsResp)
	case constant.WSKickOnlineMsg:
		u.kickOnline(*wsResp)
	default:
		utils.LogFReturn("type failed, ", wsResp.ReqIdentifier, wsResp.OperationID, wsResp.ErrCode, wsResp.ErrMsg)
		return
	}
	utils.LogSReturn()
	return
}

func (u *open_im_sdk.UserRelated) doWSGetNewestSeq(wsResp utils.GeneralWsResp) {
	utils.LogBegin(wsResp.OperationID)
	u.notifyResp(wsResp)
	utils.LogSReturn(wsResp.OperationID)
}

func (u *open_im_sdk.UserRelated) doWSPullMsg(wsResp utils.GeneralWsResp) {
	utils.LogBegin(wsResp.OperationID)
	u.notifyResp(wsResp)
	utils.LogSReturn(wsResp.OperationID)
}

func (u *open_im_sdk.UserRelated) doWSSendMsg(wsResp utils.GeneralWsResp) {
	utils.LogBegin(wsResp.OperationID)
	u.notifyResp(wsResp)
	utils.LogSReturn(wsResp.OperationID)
}

func (u *open_im_sdk.UserRelated) doWSPushMsg(wsResp utils.GeneralWsResp) {
	utils.LogBegin()
	u.doMsg(wsResp)
	utils.LogSReturn()
}

func (u *open_im_sdk.UserRelated) doMsg(wsResp utils.GeneralWsResp) {
	utils.LogBegin(wsResp.OperationID)
	var msg server_api_params.MsgData
	if wsResp.ErrCode != 0 {
		utils.sdkLog("errcode: ", wsResp.ErrCode, " errmsg: ", wsResp.ErrMsg)
		utils.LogFReturn()
		return
	}
	err := proto.Unmarshal(wsResp.Data, &msg)
	if err != nil {
		utils.sdkLog("Unmarshal failed", err.Error())
		utils.LogFReturn()
		return
	}

	utils.sdkLog("openim ws  recv push msg do push seq in : ", msg.Seq)
	u.seqMsgMutex.Lock()
	b1 := u.isExistsInErrChatLogBySeq(msg.Seq)
	b2 := u.judgeMessageIfExists(msg.ClientMsgID)
	_, ok := u.seqMsg[int32(msg.Seq)]
	if b1 || b2 || ok {
		utils.sdkLog("seq in : ", msg.Seq, b1, b2, ok)
		u.seqMsgMutex.Unlock()
		return
	}

	u.seqMsg[int32(msg.Seq)] = &msg
	u.seqMsgMutex.Unlock()

	arrMsg := utils.ArrMsg{}
	u.triggerCmdNewMsgCome(arrMsg)
}

func (u *open_im_sdk.UserRelated) GetMinSeqSvr() int64 {
	u.minSeqSvrRWMutex.RLock()
	min := u.minSeqSvr
	u.minSeqSvrRWMutex.RUnlock()
	return min
}

func (u *open_im_sdk.UserRelated) SetMinSeqSvr(minSeqSvr int64) {

	u.minSeqSvrRWMutex.Lock()
	if minSeqSvr > u.minSeqSvr {
		u.minSeqSvr = minSeqSvr
	}
	u.minSeqSvrRWMutex.Unlock()

}

func (u *open_im_sdk.UserRelated) syncSeq2Msg() error {
	svrMaxSeq, svrMinSeq, err := u.getUserNewestSeq()
	if err != nil {
		utils.sdkLog("getUserNewestSeq failed ", err.Error())
		return err
	}

	needSyncSeq := u.getNeedSyncSeq(int32(svrMinSeq), int32(svrMaxSeq))

	err = u.syncMsgFromServer(needSyncSeq)
	return err
}

func (u *open_im_sdk.UserRelated) syncLoginUserInfo() error {
	userSvr, err := u.getServerUserInfo()
	if err != nil {
		log.NewError("0", "getServerUserInfo failed , user: ", err.Error())
		return err
	}

	log.NewInfo("0", "getServerUserInfo ", userSvr)

	userLocal, err := u._getLoginUser()
	needInsert := 0
	if err != nil {
		log.NewError("0", "_getLoginUser failed  ", err.Error())
		needInsert = 1
	}

	if utils.CompFields(&userLocal, &userSvr) {
		return nil
	}

	var updateLocalUser db.LocalUser
	copier.Copy(&updateLocalUser, userSvr)
	log.NewInfo("0", "copy: ", updateLocalUser)
	if needInsert == 1 {
		err = u._insertLoginUser(&updateLocalUser)
		if err != nil {
			log.NewError("0 ", "_insertLoginUser failed ", err.Error())
		}
		return err
	}
	err = u._updateLoginUser(&updateLocalUser)
	if err != nil {
		log.NewError("0 ", "_updateLoginUser failed ", err.Error())
	}
	return err

	//if err != nil {
	//	return err
	//}
	//sdkLog("getServerUserInfo ok, user: ", *userSvr)
	//
	//userLocal, err := u.getLoginUserInfoFromLocal()
	//userLocal, err := u._getLoginUser()
	//if err != nil {
	//	return err
	//}
	//sdkLog("getLoginUserInfoFromLocal ok, user: ", userLocal)
	//
	//if userSvr.Uid != userLocal.Uid ||
	//	userSvr.Name != userLocal.Name ||
	//	userSvr.Icon != userLocal.Icon ||
	//	userSvr.Gender != userLocal.Gender ||
	//	userSvr.Mobile != userLocal.Mobile ||
	//	userSvr.Birth != userLocal.Birth ||
	//	userSvr.Email != userLocal.Email ||
	//	userSvr.Ex != userLocal.Ex {
	//	bUserInfo, err := json.Marshal(userSvr)
	//	if err != nil {
	//		sdkLog("marshal failed, ", err.Error())
	//		return err
	//	}
	//
	//	copier.Copy(a, b)
	//	err = u._updateLoginUser(userSvr)
	//	if err != nil {
	//		u.cb.OnSelfInfoUpdated(string(bUserInfo))
	//	}
	//}
}

func (u *LoginMgr) run() {
	for {
		utils.LogStart()
		if u.conn == nil {
			utils.LogBegin("reConn", nil)
			re, _, _ := u.reConn(nil)
			utils.LogEnd("reConn", re)
			u.conn = re
		}
		if u.conn != nil {
			msgType, message, err := u.conn.ReadMessage()
			utils.sdkLog("ReadMessage message ", msgType, err)
			if err != nil {
				u.stateMutex.Lock()
				utils.sdkLog("ws read message failed ", err.Error(), u.LoginState)
				if u.LoginState == constant.LogoutCmd {
					utils.sdkLog("logout, ws close, return ", constant.LogoutCmd, err)
					u.conn = nil
					u.stateMutex.Unlock()
					return
				}
				u.stateMutex.Unlock()
				time.Sleep(time.Duration(5) * time.Second)
				utils.sdkLog("ws  ReadMessage failed, sleep 5s, reconn, ", err)
				utils.LogBegin("reConn", u.conn)
				u.conn, _, err = u.reConn(u.conn)
				utils.LogEnd("reConn", u.conn)
			} else {
				if msgType == websocket.CloseMessage {
					u.conn, _, _ = u.reConn(u.conn)
				} else if msgType == websocket.TextMessage {
					utils.sdkLog("type failed, recv websocket.TextMessage ", string(message))
				} else if msgType == websocket.BinaryMessage {
					go u.doWsMsg(message)
				} else {
					utils.sdkLog("recv other msg: type ", msgType)
				}
			}
		} else {
			u.stateMutex.Lock()
			if u.LoginState == constant.LogoutCmd {
				utils.sdkLog("logout, ws close, return ", constant.LogoutCmd)
				u.stateMutex.Unlock()
				return
			}
			u.stateMutex.Unlock()
			utils.sdkLog("ws failed, sleep 5s, reconn... ")
			time.Sleep(time.Duration(5) * time.Second)
		}
	}
}

/*
func (u *UserRelated) pullBySplit(beginSeq int64, endSeq int64) error {
	LogBegin(beginSeq, endSeq)
	if beginSeq > endSeq {
		LogFReturn("beginSeq > endSeq")
		return nil
	}
	var SPLIT int64 = 100
	var bSeq, eSeq int64
	if endSeq-beginSeq > SPLIT {
		bSeq = beginSeq
		//1, 118 117/10 = 11  i: 0- 10  1-> 11 12->22 23->33  34->44   111->121
		for i := 0; int64(i) < (endSeq-beginSeq)/SPLIT; i++ {
			eSeq = bSeq + SPLIT - 1
			sdkLog("pull args: ", i, bSeq, eSeq, (endSeq-beginSeq)/SPLIT)
			err := u.pullOldMsgAndMergeNewMsg(bSeq, eSeq)
			if err != nil {
				LogFReturn(err.Error())
				return err
			}
			bSeq = eSeq + 1
		}
		if bSeq <= endSeq {
			sdkLog("pull remainder args: ", bSeq, endSeq)
			err := u.pullOldMsgAndMergeNewMsg(bSeq, endSeq)
			if err != nil {
				LogFReturn(err.Error())
				return err
			}
		}
	} else {
		err := u.pullOldMsgAndMergeNewMsg(beginSeq, endSeq)
		if err != nil {
			LogFReturn(err.Error())
			return err
		}
	}
	return nil
}

*/

/*
func (u *UserRelated) pullOldMsgAndMergeNewMsgByWs(beginSeq int64, endSeq int64) (err error) {
	LogBegin(beginSeq, endSeq)
	if beginSeq > endSeq {
		LogSReturn(nil)
		return nil
	}
	LogBegin("AddCh")
	msgIncr, ch := u.AddCh()

	var wsReq GeneralWsReq
	wsReq.ReqIdentifier = WSPullMsgBySeqList
	wsReq.OperationID = operationIDGenerator()
	wsReq.SendID = u.LoginUid
	//wsReq.Token = u.token
	wsReq.MsgIncr = msgIncr

	var pullMsgReq PullMessageBySeqListReq
	LogBegin("getNoInSeq ", beginSeq, endSeq)
	pullMsgReq.SeqList = u.getNotInSeq(beginSeq, endSeq)
	LogEnd("getNoInSeq ", pullMsgReq.SeqList)

	wsReq.Data, err = proto.Marshal(&pullMsgReq)
	if err != nil {
		sdkLog("Marshl failed ")
		LogFReturn(err.Error())
		u.DelCh(msgIncr)
		return err
	}
	LogBegin("WriteMsg ", wsReq.OperationID)
	err, _ = u.WriteMsg(wsReq)
	LogEnd("WriteMsg ", wsReq.OperationID, err)
	if err != nil {
		sdkLog("close conn, WriteMsg failed ", err.Error())
		u.DelCh(msgIncr)
		return err
	}

	timeout := 10
	select {
	case r := <-ch:
		sdkLog("ws ch recvMsg success: ", wsReq.OperationID)
		if r.ErrCode != 0 {
			sdkLog("pull msg failed ", r.ErrCode, r.ErrMsg, wsReq.OperationID)
			u.DelCh(msgIncr)
			return errors.New(r.ErrMsg)
		} else {
			sdkLog("pull msg success ", wsReq.OperationID)
			var pullMsg PullUserMsgResp

			pullMsg.ErrCode = 0

			var pullMsgResp PullMessageBySeqListResp
			err := proto.Unmarshal(r.Data, &pullMsgResp)
			if err != nil {
				sdkLog("Unmarshal failed ", err.Error())
				LogFReturn(err.Error())
				return err
			}
			pullMsg.Data.Group = pullMsgResp.GroupUserMsg
			pullMsg.Data.Single = pullMsgResp.SingleUserMsg
			pullMsg.Data.MaxSeq = pullMsgResp.MaxSeq
			pullMsg.Data.MinSeq = pullMsgResp.MinSeq

			u.seqMsgMutex.Lock()

			arrMsg := ArrMsg{}
			isInmap := false
			for i := 0; i < len(pullMsg.Data.Single); i++ {
				for j := 0; j < len(pullMsg.Data.Single[i].List); j++ {
					sdkLog("open_im pull one msg: |", pullMsg.Data.Single[i].List[j].ClientMsgID, "|")
					singleMsg := MsgData{
						SendID:           pullMsg.Data.Single[i].List[j].SendID,
						RecvID:           pullMsg.Data.Single[i].List[j].RecvID,
						SessionType:      SingleChatType,
						MsgFrom:          pullMsg.Data.Single[i].List[j].MsgFrom,
						ContentType:      pullMsg.Data.Single[i].List[j].ContentType,
						ServerMsgID:      pullMsg.Data.Single[i].List[j].ServerMsgID,
						Content:          pullMsg.Data.Single[i].List[j].Content,
						SendTime:         pullMsg.Data.Single[i].List[j].SendTime,
						Seq:              pullMsg.Data.Single[i].List[j].Seq,
						SenderNickName:   pullMsg.Data.Single[i].List[j].SenderNickName,
						SenderFaceURL:    pullMsg.Data.Single[i].List[j].SenderFaceURL,
						ClientMsgID:      pullMsg.Data.Single[i].List[j].ClientMsgID,
						SenderPlatformID: pullMsg.Data.Single[i].List[j].SenderPlatformID,
					}
					//	arrMsg.SingleData = append(arrMsg.SingleData, singleMsg)
					u.seqMsg[pullMsg.Data.Single[i].List[j].Seq] = singleMsg
					sdkLog("into map, seq: ", pullMsg.Data.Single[i].List[j].Seq, pullMsg.Data.Single[i].List[j].ClientMsgID, pullMsg.Data.Single[i].List[j].ServerMsgID)
				}
			}



			for i := 0; i < len(pullMsg.Data.Group); i++ {
				for j := 0; j < len(pullMsg.Data.Group[i].List); j++ {
					groupMsg := MsgData{
						SendID:           pullMsg.Data.Group[i].List[j].SendID,
						RecvID:           pullMsg.Data.Group[i].List[j].RecvID,
						SessionType:      GroupChatType,
						MsgFrom:          pullMsg.Data.Group[i].List[j].MsgFrom,
						ContentType:      pullMsg.Data.Group[i].List[j].ContentType,
						ServerMsgID:      pullMsg.Data.Group[i].List[j].ServerMsgID,
						Content:          pullMsg.Data.Group[i].List[j].Content,
						SendTime:         pullMsg.Data.Group[i].List[j].SendTime,
						Seq:              pullMsg.Data.Group[i].List[j].Seq,
						SenderNickName:   pullMsg.Data.Group[i].List[j].SenderNickName,
						SenderFaceURL:    pullMsg.Data.Group[i].List[j].SenderFaceURL,
						ClientMsgID:      pullMsg.Data.Group[i].List[j].ClientMsgID,
						SenderPlatformID: pullMsg.Data.Group[i].List[j].SenderPlatformID,
					}
					//	arrMsg.GroupData = append(arrMsg.GroupData, groupMsg)
					u.seqMsg[pullMsg.Data.Group[i].List[j].Seq] = groupMsg
					sdkLog("into map, seq: ", pullMsg.Data.Group[i].List[j].Seq, pullMsg.Data.Group[i].List[j].ClientMsgID, pullMsg.Data.Group[i].List[j].ServerMsgID)
				}
			}
			u.seqMsgMutex.Unlock()

			u.seqMsgMutex.RLock()
			for i := beginSeq; i <= endSeq; i++ {
				v, ok := u.seqMsg[i]
				if ok {
					if v.SessionType == SingleChatType {
						arrMsg.SingleData = append(arrMsg.SingleData, v)
						sdkLog("pull seq: ", v.Seq, v)
						if v.ContentType > SingleTipBegin && v.ContentType < SingleTipEnd {
							var msgRecv MsgData
							msgRecv.ContentType = v.ContentType
							msgRecv.Content = v.Content
							msgRecv.SendID = v.SendID
							msgRecv.RecvID = v.RecvID
							LogBegin("doFriendMsg ", msgRecv)
							u.doFriendMsg(msgRecv)
							LogEnd("doFriendMsg ", msgRecv)
						}
					} else if v.SessionType == GroupChatType {
						sdkLog("pull seq: ", v.Seq, v)
						arrMsg.GroupData = append(arrMsg.GroupData, v)
						if v.ContentType > GroupTipBegin && v.ContentType < GroupTipEnd {
							LogBegin("doGroupMsg ", v)
							u.doGroupMsg(v)
							LogEnd("doGroupMsg ", v)
						}
					} else {
						sdkLog("type failed, ", v.SessionType, v)
					}
				} else {
					sdkLog("seq no in map, failed, seq: ", i)
				}
			}
			u.seqMsgMutex.RUnlock()

			sdkLog("triggerCmdNewMsgCome len: ", len(arrMsg.SingleData), len(arrMsg.GroupData))
			err = u.triggerCmdNewMsgCome(arrMsg)
			if err != nil {
				sdkLog("triggerCmdNewMsgCome failed, ", err.Error())
			}
			u.DelCh(msgIncr)
		}
	case <-time.After(time.Second * time.Duration(timeout)):
		sdkLog("ws ch recvMsg timeout,", wsReq.OperationID)
		u.DelCh(msgIncr)
	}
	return nil
}

*/
func CheckToken(uId, token string) int {
	_, err := utils.post2Api(open_im_sdk.newestSeqRouter, open_im_sdk.paramsNewestSeqReq{ReqIdentifier: 1001, OperationID: utils.operationIDGenerator(), SendID: uId, MsgIncr: 1}, token)
	if err != nil {
		return -1
	}
	return 0
}

func (u *open_im_sdk.UserRelated) getUserNewestSeq() (int64, int64, error) {
	utils.LogBegin()
	resp, err := utils.post2Api(open_im_sdk.newestSeqRouter, open_im_sdk.paramsNewestSeqReq{ReqIdentifier: 1001, OperationID: utils.operationIDGenerator(), SendID: u.loginUserID, MsgIncr: 1}, u.token)
	if err != nil {
		utils.LogFReturn(0, err.Error())
		return 0, 0, err
	}
	var seqResp open_im_sdk.paramsNewestSeqResp
	err = json.Unmarshal(resp, &seqResp)
	if err != nil {
		utils.sdkLog("UnMarshal failed, ", err.Error())
		utils.LogFReturn(0, err.Error())
		return 0, 0, err
	}

	if seqResp.ErrCode != 0 {
		utils.sdkLog("errcode: ", seqResp.ErrCode, "errmsg: ", seqResp.ErrMsg)
		utils.LogFReturn(0, seqResp.ErrMsg)
		return 0, 0, errors.New(seqResp.ErrMsg)
	}
	utils.LogSReturn(seqResp.Data.Seq, nil)
	return seqResp.Data.Seq, seqResp.Data.MinSeq, nil
}

func (u *open_im_sdk.UserRelated) getServerUserInfo() (*server_api_params.UserInfo, error) {
	apiReq := server_api_params.GetUserInfoReq{OperationID: utils.operationIDGenerator(), UserIDList: []string{u.loginUserID}}
	resp, err := utils.post2Api(open_im_sdk.getUserInfoRouter, apiReq, u.token)
	commData, err := utils.checkErrAndRespReturn(err, resp, apiReq.OperationID)
	if err != nil {
		return nil, utils.wrap(err, apiReq.OperationID)
	}
	realData := server_api_params.GetUserInfoResp{}
	err = mapstructure.Decode(commData.Data, &realData.UserInfoList)
	if err != nil {
		log.NewError(apiReq.OperationID, "Decode failed ", err.Error())
		return nil, err
	}
	log.NewInfo(apiReq.OperationID, "realData.UserInfoList", realData.UserInfoList, commData.Data)
	if len(realData.UserInfoList) == 0 {
		log.NewInfo(apiReq.OperationID, "failed, no user : ", u.loginUserID)
		return nil, errors.New("no login user")
	}
	log.NewInfo(apiReq.OperationID, "realData.UserInfoList[0]", realData.UserInfoList[0])
	return realData.UserInfoList[0], nil
}

func (u *open_im_sdk.UserRelated) getUserNameAndFaceUrlByUid(uid string) (faceUrl, name string, err error) {
	friendInfo, err := u._getFriendInfoByFriendUserID(uid)
	if err != nil {
		return "", "", err
	}
	if friendInfo.FriendUserID == "" {
		userInfo, err := u.getUserInfoByUid(uid)
		if err != nil {
			return "", "", err
		} else {
			return userInfo.Icon, userInfo.Name, nil
		}
	} else {
		if friendInfo.Remark != "" {
			return friendInfo.FaceUrl, friendInfo.Remark, nil
		} else {
			return friendInfo.FaceUrl, friendInfo.Nickname, nil
		}
	}
}
func (u *open_im_sdk.UserRelated) getUserInfoByUid(uid string) (*open_im_sdk.userInfo, error) {
	var uidList []string
	uidList = append(uidList, uid)
	resp, err := utils.post2Api(open_im_sdk.getUserInfoRouter, open_im_sdk.paramsGetUserInfo{OperationID: utils.operationIDGenerator(), UidList: uidList}, u.token)
	if err != nil {
		utils.sdkLog("post2Api failed, ", open_im_sdk.getUserInfoRouter, uidList, err.Error())
		return nil, err
	}
	utils.sdkLog("post api: ", open_im_sdk.getUserInfoRouter, open_im_sdk.paramsGetUserInfo{OperationID: utils.operationIDGenerator(), UidList: uidList}, "uid ", uid)
	var userResp open_im_sdk.getUserInfoResp
	err = json.Unmarshal(resp, &userResp)
	if err != nil {
		utils.sdkLog("Unmarshal failed, ", resp, err.Error())
		return nil, err
	}

	if userResp.ErrCode != 0 {
		utils.sdkLog("errcode: ", userResp.ErrCode, "errmsg:", userResp.ErrMsg)
		return nil, errors.New(userResp.ErrMsg)
	}

	if len(userResp.Data) == 0 {
		utils.sdkLog("failed, no user :", uid)
		return nil, errors.New("no user")
	}
	return &userResp.Data[0], nil
}

func (u *open_im_sdk.UserRelated) doFriendMsg(msg *server_api_params.MsgData) {
	utils.sdkLog("doFriendMsg ", msg)
	if u.cb == nil || u.friendListener == nil {
		utils.sdkLog("listener is null")
		return
	}

	if msg.SendID == u.loginUserID && msg.SenderPlatformID == u.SvrConf.Platform {
		utils.sdkLog("sync msg ", msg.ContentType)
		return
	}

	go func() {
		switch msg.ContentType {
		case constant.AddFriendTip:
			utils.sdkLog("addFriendNew ", msg)
			u.addFriendNew(msg) //
		case constant.AcceptFriendApplicationTip:
			utils.sdkLog("acceptFriendApplicationNew ", msg)
			u.acceptFriendApplicationNew(msg)
		case constant.RefuseFriendApplicationTip:
			utils.sdkLog("refuseFriendApplicationNew ", msg)
			u.refuseFriendApplicationNew(msg)
		case constant.SetSelfInfoTip:
			utils.sdkLog("setSelfInfo ", msg)
			u.setSelfInfo(msg)
			//	case KickOnlineTip:
			//		sdkLog("kickOnline ", msg)
			//		u.kickOnline(&msg)
		default:
			utils.sdkLog("type failed, ", msg)
		}
	}()
}

func (u *open_im_sdk.UserRelated) acceptFriendApplicationNew(msg *server_api_params.MsgData) {
	utils.LogBegin(msg.ContentType, msg.ServerMsgID, msg.ClientMsgID)
	u.syncFriendList()
	utils.sdkLog(msg.SendID, msg.RecvID)
	utils.sdkLog("acceptFriendApplicationNew", msg.ServerMsgID, msg)

	fInfoList, err := u.getServerFriendList()
	if err != nil {
		return
	}
	utils.sdkLog("fInfoList", fInfoList)

	//for _, fInfo := range fInfoList {
	//	if fInfo.UID == msg.SendID {
	//		jData, err := json.Marshal(fInfo)
	//		if err != nil {
	//			sdkLog("err: ", err.Error())
	//			return
	//		}
	//		u.friendListener.OnFriendListAdded(string(jData))
	//		u.friendListener.OnFriendApplicationListAccept(string(jData))
	//		return
	//	}
	//}
}

func (u *open_im_sdk.UserRelated) refuseFriendApplicationNew(msg *server_api_params.MsgData) {
	utils.sdkLog(msg.SendID, msg.RecvID)
	applyList, err := u.getServerSelfApplication()

	if err != nil {
		return
	}
	for _, v := range applyList {
		if v.Uid == msg.SendID {
			jData, err := json.Marshal(v)
			if err != nil {
				utils.sdkLog("err: ", err.Error())
				return
			}
			u.friendListener.OnFriendApplicationListReject(string(jData))
			return
		}
	}

}

func (u *open_im_sdk.UserRelated) addFriendNew(msg *server_api_params.MsgData) {
	utils.sdkLog("addFriend start ")
	u.syncFriendApplication()

	var ui2GetUserInfo open_im_sdk.ui2ClientCommonReq
	ui2GetUserInfo.UidList = append(ui2GetUserInfo.UidList, msg.SendID)
	resp, err := utils.post2Api(open_im_sdk.getUserInfoRouter, open_im_sdk.paramsGetUserInfo{UidList: ui2GetUserInfo.UidList, OperationID: utils.operationIDGenerator()}, u.token)
	if err != nil {
		utils.sdkLog("getUserInfo failed", err)
		return
	}
	var vgetUserInfoResp open_im_sdk.getUserInfoResp
	err = json.Unmarshal(resp, &vgetUserInfoResp)
	if err != nil {
		utils.sdkLog("Unmarshal failed, ", err.Error())
		return
	}
	if vgetUserInfoResp.ErrCode != 0 {
		utils.sdkLog(vgetUserInfoResp.ErrCode, vgetUserInfoResp.ErrMsg)
		return
	}
	if len(vgetUserInfoResp.Data) == 0 {
		utils.sdkLog(vgetUserInfoResp.ErrCode, vgetUserInfoResp.ErrMsg, msg)
		return
	}
	var appUserNode open_im_sdk.applyUserInfo
	appUserNode.Uid = vgetUserInfoResp.Data[0].Uid
	appUserNode.Name = vgetUserInfoResp.Data[0].Name
	appUserNode.Icon = vgetUserInfoResp.Data[0].Icon
	appUserNode.Gender = vgetUserInfoResp.Data[0].Gender
	appUserNode.Mobile = vgetUserInfoResp.Data[0].Mobile
	appUserNode.Birth = vgetUserInfoResp.Data[0].Birth
	appUserNode.Email = vgetUserInfoResp.Data[0].Email
	appUserNode.Ex = vgetUserInfoResp.Data[0].Ex
	appUserNode.Flag = 0

	jsonInfo, err := json.Marshal(appUserNode)
	if err != nil {
		utils.sdkLog("  marshal failed", err.Error())
		return
	}
	u.friendListener.OnFriendApplicationListAdded(string(jsonInfo))
}

func (u *open_im_sdk.UserRelated) kickOnline(msg utils.GeneralWsResp) {
	utils.sdkLog("kickOnline ", msg.ReqIdentifier, msg.ErrCode, msg.ErrMsg)
	u.logout(nil)
	u.cb.OnKickedOffline()
}

func (u *open_im_sdk.UserRelated) setSelfInfo(msg *server_api_params.MsgData) {
	var uidList []string
	uidList = append(uidList, msg.SendID)
	resp, err := utils.post2Api(open_im_sdk.getUserInfoRouter, open_im_sdk.paramsGetUserInfo{OperationID: utils.operationIDGenerator(), UidList: uidList}, u.token)
	if err != nil {
		utils.sdkLog("post2Api failed, ", open_im_sdk.getUserInfoRouter, uidList, err.Error())
		return
	}
	var userResp open_im_sdk.getUserInfoResp
	err = json.Unmarshal(resp, &userResp)
	if err != nil {
		utils.sdkLog("Unmarshal failed, ", resp, err.Error())
		return
	}

	if userResp.ErrCode != 0 {
		utils.sdkLog("errcode: ", userResp.ErrCode, "errmsg:", userResp.ErrMsg)
		return
	}

	if len(userResp.Data) == 0 {
		utils.sdkLog("failed, no user : ", u.loginUserID)
		return
	}

	err = u.updateFriendInfo(userResp.Data[0].Uid, userResp.Data[0].Name, userResp.Data[0].Icon, userResp.Data[0].Gender, userResp.Data[0].Mobile, userResp.Data[0].Birth, userResp.Data[0].Email, userResp.Data[0].Ex)
	if err != nil {
		utils.sdkLog("  db change failed", err.Error())
		return
	}

	jsonInfo, err := json.Marshal(userResp.Data[0])
	if err != nil {
		utils.sdkLog("  marshal failed", err.Error())
		return
	}

	u.friendListener.OnFriendInfoChanged(string(jsonInfo))
}
