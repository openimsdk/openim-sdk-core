package init

import (
	conv "open_im_sdk/internal/controller/conversation_msg"
	"open_im_sdk/internal/controller/friend"
	"open_im_sdk/internal/controller/group"
	ws "open_im_sdk/internal/controller/interaction"
	"open_im_sdk/internal/controller/user"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
)

type LoginMgr struct {
	friend       *friend.Friend
	group        *group.Group
	conversation *conv.Conversation
	user         *user.User

	db      *db.DataBase  //1
	ws      *ws.Ws  //2
	msgSync *MsgSync //3

	heartbeat *Heartbeat //4

	token       string
	loginUserID string
	listener    ws.ConnListener

	justOnceFlag bool

	groupListener group.OnGroupListener
	friendListener friend.OnFriendshipListener
	conversationListener conv.OnConversationListener
	advancedMsgListener conv.OnAdvancedMsgListener


	conversationCh chan common.Cmd2Value
	cmdCh chan common.Cmd2Value
}

func (u *LoginMgr) SetConversationListener(conversationListener conv.OnConversationListener) {
	u.conversationListener = conversationListener
}

func (u *LoginMgr) SetFriendListener(friendListener friend.OnFriendshipListener) {
	u.friendListener = friendListener
}

func (u *LoginMgr) SetGroupListener(groupListener group.OnGroupListener) {
	u.groupListener = groupListener
}

func (u *LoginMgr) login(userID, token string, cb common.Base) {
	log.Info("login start ", userID, token)
	if cb == nil {
		log.Info("cb == nil ", userID)
		return
	}
	if u.justOnceFlag {
		cb.OnError(constant.ErrLogin.ErrCode, constant.ErrLogin.ErrMsg)
		return
	}
	err := u.checkToken(token)
	if err != nil {
		cb.OnError(constant.ErrTokenInvalid.ErrCode, constant.ErrTokenInvalid.ErrMsg)
		return
	}
	u.justOnceFlag = true

	u.token = token
	u.loginUserID = userID

	db, err := db.NewDataBase(userID)
	if err != nil {
		cb.OnError(constant.ErrDB.ErrCode, constant.ErrDB.ErrMsg)
		log.Error("0", "NewDataBase failed ", err.Error())
		return
	}
	u.db = db

	wsRespAsyn := ws.NewWsRespAsyn()
	wsConn := ws.NewWsConn(u.listener, token, userID)
	u.conversationCh = make(chan common.Cmd2Value, 1000)
	u.cmdCh = make(chan common.Cmd2Value, 10)


	u.ws = ws.NewWs(wsRespAsyn, wsConn, u.conversationCh, u.cmdCh)
	u.msgSync = NewMsgSync(db, u.ws, userID, u.conversationCh)

	u.heartbeat = NewHeartbeat(u.ws, u.msgSync)

	p := ws.NewPostApi(token, constant.SvrConf.ApiAddr)
	u.user = user.NewUser(db, p, u.loginUserID)

	u.friend = friend.NewFriend(u.loginUserID, u.db, p)
	u.friend.SetFriendListener(u.friendListener)

	u.group = group.NewGroup(u.loginUserID, u.db, p)
	u.group.SetGroupListener(u.groupListener)

	u.conversation = conv.NewConversation(u.ws, u.db, u.conversationCh, u.loginUserID, u.friend, u.group, u.user)
	u.conversation.SetConversationListener(u.conversationListener)

	log.Info("ws, forcedSynchronization heartbeat ws coroutine run ...")
	go u.forcedSynchronization()
	go u.heartbeat.Run()
	go u.ws.ReadData()
//	u.forycedSyncReceiveMessageOpt()
	cb.OnSuccess("")

}



func (u *LoginMgr) InitSDK(config string, listener ws.ConnListener) bool {
	if listener == nil {
		return false
	}
	u.listener = listener
	return true
}

func (u *LoginMgr) UnInitSDK() {

}

func (u *LoginMgr) GetVersion() string {
	return "v1.0.5"
}


func (u *LoginMgr) logout(callback common.Base) {
	common.TriggerCmdLogout(utils.ArrMsg{}, u.cmdCh)
	timeout := 5
	resp, err, operationID := u.ws.SendReqWaitResp(nil, constant.WsLogoutMsg, timeout, u.loginUserID)
}


func (u *LoginMgr) GetLoginUser() string {
	if u.GetLoginStatus() == constant.LoginSuccess {
		return u.loginUserID
	} else {
		return ""
	}
}

func (u *LoginMgr) GetLoginStatus() int {
	return u.GetLoginStatus()
}


func (u *LoginMgr) forcedSynchronization() {
	u.friend.SyncFriendList()
	u.friend.SyncBlackList()
	u.friend.SyncFriendApplication()
	u.friend.SyncSelfFriendApplication()
	u.user.SyncLoginUserInfo()
	u.group.SyncApplyGroupRequest()
	u.group.SyncGroupRequest()
	u.group.SyncJoinedGroupInfo()
	u.group.SyncSelfGroupRequest()
}

func (u *LoginMgr) GetMinSeqSvr() int64 {
	return u.GetMinSeqSvr()
}

func (u *LoginMgr) SetMinSeqSvr(minSeqSvr int64) {
	u.SetMinSeqSvr(minSeqSvr)
}

func (u *LoginMgr)checkToken(token string) error {
//	p := ws.NewPostApi(token, constant.SvrConf.ApiAddr)
	_, err := u.user.GetSelfUserInfoFromSvr()
	return utils.Wrap(err, "GetSelfUserInfoFromSvr failed")
}


//func (u *open_im_sdk.UserRelated) kickOnline(msg utils.GeneralWsResp) {
//	utils.sdkLog("kickOnline ", msg.ReqIdentifier, msg.ErrCode, msg.ErrMsg)
//	u.logout(nil)
//	u.cb.OnKickedOffline()
//}


//
//func (u *open_im_sdk.UserRelated) forycedSyncReceiveMessageOpt() {
//	OperationID := utils.operationIDGenerator()
//	resp, err := utils.post2ApiForRead(open_im_sdk.getAllConversationMessageOptRouter, open_im_sdk.paramGetAllConversationMessageOpt{OperationID: OperationID}, u.token)
//	if err != nil {
//		utils.sdkLog("post2Api failed, ", open_im_sdk.getAllConversationMessageOptRouter, OperationID)
//		return
//	}
//	var v open_im_sdk.getReceiveMessageOptResp
//	err = json.Unmarshal(resp, &v)
//	if err != nil {
//		utils.sdkLog("Unmarshal failed ", resp, OperationID)
//		return
//	}
//	if v.ErrCode != 0 {
//		utils.sdkLog("errCode failed, ", v.ErrCode, v.ErrMsg, string(resp), OperationID)
//		return
//	}
//
//	utils.sdkLog("get receive opt ", v)
//	u.receiveMessageOptMutex.Lock()
//	for _, v := range v.Data {
//		if v.Result != 0 {
//			u.receiveMessageOpt[v.ConversationId] = v.Result
//		}
//	}
//	u.receiveMessageOptMutex.Unlock()
//}



