package login

import (
	"errors"
	comm2 "open_im_sdk/internal/common"
	conv "open_im_sdk/internal/conversation_msg"
	"open_im_sdk/internal/friend"
	"open_im_sdk/internal/full"
	"open_im_sdk/internal/group"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/internal/user"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

type LoginMgr struct {
	friend       *friend.Friend
	group        *group.Group
	conversation *conv.Conversation
	user         *user.User
	full         *full.Full
	db           *db.DataBase
	ws           *ws.Ws
	msgSync      *ws.MsgSync

	heartbeat *ws.Heartbeat

	token        string
	loginUserID  string
	connListener open_im_sdk_callback.OnConnListener

	justOnceFlag bool

	groupListener        open_im_sdk_callback.OnGroupListener
	friendListener       open_im_sdk_callback.OnFriendshipListener
	conversationListener open_im_sdk_callback.OnConversationListener
	advancedMsgListener  open_im_sdk_callback.OnAdvancedMsgListener
	userListener         open_im_sdk_callback.OnUserListener

	conversationCh chan common.Cmd2Value
	cmdWsCh        chan common.Cmd2Value
	heartbeatCmdCh chan common.Cmd2Value
	imConfig       sdk_struct.IMConfig
}

func (u *LoginMgr) ImConfig() sdk_struct.IMConfig {
	return u.imConfig
}

func (u *LoginMgr) Conversation() *conv.Conversation {
	return u.conversation
}

func (u *LoginMgr) User() *user.User {
	return u.user
}

func (u *LoginMgr) Full() *full.Full {
	return u.full
}

func (u *LoginMgr) Group() *group.Group {
	return u.group
}

func (u *LoginMgr) Friend() *friend.Friend {
	return u.friend
}

func (u *LoginMgr) SetConversationListener(conversationListener open_im_sdk_callback.OnConversationListener) {
	u.conversationListener = conversationListener
}

func (u *LoginMgr) SetAdvancedMsgListener(advancedMsgListener open_im_sdk_callback.OnAdvancedMsgListener) {
	u.advancedMsgListener = advancedMsgListener
}

func (u *LoginMgr) SetFriendListener(friendListener open_im_sdk_callback.OnFriendshipListener) {
	u.friendListener = friendListener
}

func (u *LoginMgr) SetGroupListener(groupListener open_im_sdk_callback.OnGroupListener) {
	u.groupListener = groupListener
}

func (u *LoginMgr) SetUserListener(userListener open_im_sdk_callback.OnUserListener) {
	u.userListener = userListener
}

func (u *LoginMgr) login(userID, token string, cb open_im_sdk_callback.Base, operationID string) {
	log.Info(operationID, "login start... ", userID, token)
	if u.justOnceFlag {
		cb.OnError(constant.ErrLogin.ErrCode, constant.ErrLogin.ErrMsg)
		return
	}
	err := CheckToken(userID, token, operationID)
	common.CheckTokenErrCallback(cb, err, operationID)
	log.Info(operationID, "checkToken ok ", userID, token)
	u.justOnceFlag = true

	u.token = token
	u.loginUserID = userID

	db, err := db.NewDataBase(userID, sdk_struct.SvrConf.DataDir)
	if err != nil {
		cb.OnError(constant.ErrDB.ErrCode, constant.ErrDB.ErrMsg)
		log.Error(operationID, "NewDataBase failed ", err.Error())
		return
	}
	u.db = db
	log.Info(operationID, "NewDataBase ok ", userID, sdk_struct.SvrConf.DataDir)
	wsRespAsyn := ws.NewWsRespAsyn()
	wsConn := ws.NewWsConn(u.connListener, token, userID)
	u.conversationCh = make(chan common.Cmd2Value, 1000)
	u.cmdWsCh = make(chan common.Cmd2Value, 10)

	pushMsgAndMaxSeqCh := make(chan common.Cmd2Value, 1000)
	u.ws = ws.NewWs(wsRespAsyn, wsConn, u.cmdWsCh, pushMsgAndMaxSeqCh)
	u.msgSync = ws.NewMsgSync(db, u.ws, userID, u.conversationCh, pushMsgAndMaxSeqCh)

	u.heartbeatCmdCh = make(chan common.Cmd2Value, 10)
	u.heartbeat = ws.NewHeartbeat(u.msgSync, u.heartbeatCmdCh)

	p := ws.NewPostApi(token, sdk_struct.SvrConf.ApiAddr)
	u.user = user.NewUser(db, p, u.loginUserID)
	u.user.SetListener(u.userListener)

	u.friend = friend.NewFriend(u.loginUserID, u.db, u.user, p)
	u.friend.SetFriendListener(u.friendListener)

	u.group = group.NewGroup(u.loginUserID, u.db, p)
	u.group.SetGroupListener(u.groupListener)

	u.full = full.NewFull(u.user, u.friend, u.group)
	if u.imConfig.ObjectStorage != "cos" && u.imConfig.ObjectStorage != "" {
		err = errors.New("u.imConfig.ObjectStorage failed ")
		common.CheckConfigErrCallback(cb, err, operationID)
	}
	u.forcedSynchronization()
	log.Info(operationID, "forcedSynchronization success...")
	objStorage := comm2.NewCOS(p)
	u.conversation = conv.NewConversation(u.ws, u.db, p, u.conversationCh,
		u.loginUserID, u.imConfig.Platform, u.imConfig.DataDir,
		u.friend, u.group, u.user, objStorage)
	u.conversation.SetConversationListener(u.conversationListener)
	u.conversation.SetMsgListener(u.advancedMsgListener)
	log.Info(operationID, "login success...")
	//u.forycedSyncReceiveMessageOpt()
	cb.OnSuccess("")

}

func (u *LoginMgr) InitSDK(config sdk_struct.IMConfig, listener open_im_sdk_callback.OnConnListener, operationID string) bool {
	u.imConfig = config
	log.NewInfo(operationID, utils.GetSelfFuncName(), config)
	if listener == nil {
		return false
	}
	u.connListener = listener
	return true
}

func (u *LoginMgr) logout(callback open_im_sdk_callback.Base, operationID string) {
	log.Info(operationID, "TriggerCmdLogout ws...")

	if u.friend == nil || u.conversation == nil || u.user == nil || u.full == nil ||
		u.db == nil || u.ws == nil || u.msgSync == nil || u.heartbeat == nil {
		log.Info(operationID, "nil, no TriggerCmdLogout ", *u)
		return
	}

	err := common.TriggerCmdLogout(u.cmdWsCh)
	if err != nil {
		log.Error(operationID, "TriggerCmdLogout failed ", err.Error())
	}
	log.Info(operationID, "TriggerCmdLogout heartbeat...")
	err = common.TriggerCmdLogout(u.heartbeatCmdCh)
	if err != nil {
		log.Error(operationID, "TriggerCmdLogout failed ", err.Error())
	}

	timeout := 5
	retryTimes := 0
	log.Info(operationID, "send to svr logout ...", u.loginUserID)
	resp, err := u.ws.SendReqWaitResp(&server_api_params.GetMaxAndMinSeqReq{}, constant.WsLogoutMsg, timeout, retryTimes, u.loginUserID, operationID)
	if err != nil {
		log.Warn(operationID, "SendReqWaitResp failed ", err.Error(), constant.WsLogoutMsg, timeout, u.loginUserID, resp)
		//if callback != nil {
		//	callback.OnError(constant.ErrArgs.ErrCode, constant.ErrArgs.ErrMsg)
		//} else {
		//	return
		//}
	}
	if callback != nil {
		callback.OnSuccess("")
	}
	u.justOnceFlag = false
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
	u.group.SyncAdminGroupApplication(operationID)
	u.group.SyncSelfGroupApplication(operationID)
	u.group.SyncJoinedGroupMember(operationID)
}

func (u *LoginMgr) GetMinSeqSvr() int64 {
	return u.GetMinSeqSvr()
}

func (u *LoginMgr) SetMinSeqSvr(minSeqSvr int64) {
	u.SetMinSeqSvr(minSeqSvr)
}

func CheckToken(userID, token string, operationID string) error {
	if operationID == "" {
		operationID = utils.OperationIDGenerator()
	}

	log.Debug(operationID, utils.GetSelfFuncName(), userID, token)
	p := ws.NewPostApi(token, sdk_struct.SvrConf.ApiAddr)
	_, err := user.NewUser(nil, p, userID).GetSelfUserInfoFromSvr(operationID)
	return utils.Wrap(err, "GetSelfUserInfoFromSvr failed "+operationID)
}

func (u *LoginMgr) uploadImage(callback open_im_sdk_callback.Base, filePath string, token, obj string, operationID string) string {
	if obj == "cos" {
		p := ws.NewPostApi(token, u.ImConfig().ApiAddr)
		o := comm2.NewCOS(p)
		url, _, err := o.UploadImage(filePath, func(progress int) {
			if progress == 100 {
				callback.OnSuccess("")
			}
		})
		if err != nil {
			log.Error(operationID, "UploadImage failed ", err.Error(), filePath)
			return ""
		}
		return url
	} else {
		return ""
	}
}
