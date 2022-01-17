package login

import (
	conv "open_im_sdk/internal/conversation_msg"
	"open_im_sdk/internal/friend"
	"open_im_sdk/internal/group"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/internal/user"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

type LoginMgr struct {
	friend       *friend.Friend
	group        *group.Group
	conversation *conv.Conversation
	user         *user.User

	db      *db.DataBase //1
	ws      *ws.Ws       //2
	msgSync *ws.MsgSync  //3

	heartbeat *ws.Heartbeat //4

	token        string
	loginUserID  string
	connListener ws.ConnListener

	justOnceFlag bool

	groupListener        group.OnGroupListener
	friendListener       friend.OnFriendshipListener
	conversationListener conv.OnConversationListener
	advancedMsgListener  conv.OnAdvancedMsgListener
	userListener         user.OnUserListener

	conversationCh chan common.Cmd2Value
	cmdCh          chan common.Cmd2Value

	imConfig sdk_struct.IMConfig
}

func (u *LoginMgr) Conversation() *conv.Conversation {
	return u.conversation
}

func (u *LoginMgr) User() *user.User {
	return u.user
}

func (u *LoginMgr) Group() *group.Group {
	return u.group
}

func (u *LoginMgr) Friend() *friend.Friend {
	return u.friend
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

func (u *LoginMgr) SetUserListener(userListener user.OnUserListener) {
	u.userListener = userListener
}

func (u *LoginMgr) login(userID, token string, cb common.Base, operationID string) {
	log.Info("login start ", userID, token)
	if u.justOnceFlag {
		cb.OnError(constant.ErrLogin.ErrCode, constant.ErrLogin.ErrMsg)
		return
	}
	err := CheckToken(userID, token)
	common.CheckAnyErr(cb, 1111, err, operationID)
	log.Info("0", "checkToken ok ", userID, token)
	u.justOnceFlag = true

	u.token = token
	u.loginUserID = userID

	db, err := db.NewDataBase(userID, sdk_struct.SvrConf.DbDir)
	if err != nil {
		cb.OnError(constant.ErrDB.ErrCode, constant.ErrDB.ErrMsg)
		log.Error("0", "NewDataBase failed ", err.Error())
		return
	}
	u.db = db
	log.Info("0", "NewDataBase ok ", userID, sdk_struct.SvrConf.DbDir)
	wsRespAsyn := ws.NewWsRespAsyn()
	wsConn := ws.NewWsConn(u.connListener, token, userID)
	u.conversationCh = make(chan common.Cmd2Value, 1000)
	u.cmdCh = make(chan common.Cmd2Value, 10)

	u.ws = ws.NewWs(wsRespAsyn, wsConn, u.conversationCh, u.cmdCh)
	u.msgSync = ws.NewMsgSync(db, u.ws, userID, u.conversationCh)

	u.heartbeat = ws.NewHeartbeat(u.msgSync)

	p := ws.NewPostApi(token, sdk_struct.SvrConf.ApiAddr)
	u.user = user.NewUser(db, p, u.loginUserID)
	u.SetUserListener(u.userListener)

	u.friend = friend.NewFriend(u.loginUserID, u.db, p)
	u.friend.SetFriendListener(u.friendListener)

	u.group = group.NewGroup(u.loginUserID, u.db, p)
	u.group.SetGroupListener(u.groupListener)

	u.conversation = conv.NewConversation(u.ws, u.db, u.conversationCh, u.loginUserID, u.friend, u.group, u.user)
	u.conversation.SetConversationListener(u.conversationListener)

	log.Info("forcedSynchronization run ...")
	go u.forcedSynchronization()
	//	u.forycedSyncReceiveMessageOpt()
	cb.OnSuccess("")

}

func (u *LoginMgr) InitSDK(config sdk_struct.IMConfig, listener ws.ConnListener, operationID string) bool {
	u.imConfig = config
	log.NewInfo(operationID, utils.GetSelfFuncName(), config)
	if listener == nil {
		return false
	}
	u.connListener = listener
	return true
}

func (u *LoginMgr) logout(callback common.Base, operationID string) {
	err := common.TriggerCmdLogout(sdk_struct.ArrMsg{}, u.cmdCh)
	if err != nil {
		log.Error(operationID, "TriggerCmdLogout failed ", err.Error())
		return
	}
	timeout := 5
	retryTimes := 0
	resp, err := u.ws.SendReqWaitResp(nil, constant.WsLogoutMsg, timeout, retryTimes, u.loginUserID, operationID)
	if err != nil {
		log.Error(operationID, "SendReqWaitResp failed ", err.Error(), constant.WsLogoutMsg, timeout, u.loginUserID, resp)
		if callback != nil {
			callback.OnError(constant.ErrArgs.ErrCode, constant.ErrArgs.ErrMsg)
		} else {
			return
		}
	}
	if callback != nil {
		callback.OnSuccess("")
	}
}

func (u *LoginMgr) GetLoginUser() string {
	if u.GetLoginStatus() == constant.LoginSuccess {
		return u.loginUserID
	} else {
		return ""
	}
}

func (u *LoginMgr) GetLoginStatus() int32 {
	return u.ws.LoginState()
}

func (u *LoginMgr) forcedSynchronization() {
	operationID := utils.OperationIDGenerator()
	u.friend.SyncFriendList(operationID)
	u.friend.SyncBlackList(operationID)
	u.friend.SyncFriendApplication(operationID)
	u.friend.SyncSelfFriendApplication(operationID)

	u.user.SyncLoginUserInfo(operationID)

	u.group.SyncJoinedGroupList(operationID)
	u.group.SyncGroupApplication(operationID)
	u.group.SyncSelfGroupApplication(operationID)
}

func (u *LoginMgr) GetMinSeqSvr() int64 {
	return u.GetMinSeqSvr()
}

func (u *LoginMgr) SetMinSeqSvr(minSeqSvr int64) {
	u.SetMinSeqSvr(minSeqSvr)
}

func CheckToken(userID, token string) error {
	operationID := utils.OperationIDGenerator()
	log.Debug(operationID, utils.GetSelfFuncName(), userID, token)
	p := ws.NewPostApi(token, sdk_struct.SvrConf.ApiAddr)
	_, err := user.NewUser(nil, p, userID).GetSelfUserInfoFromSvr(operationID)
	return utils.Wrap(err, "GetSelfUserInfoFromSvr failed "+operationID)
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
