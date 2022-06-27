package login

import (
	"open_im_sdk/internal/cache"
	comm2 "open_im_sdk/internal/common"
	conv "open_im_sdk/internal/conversation_msg"
	"open_im_sdk/internal/friend"
	"open_im_sdk/internal/full"
	"open_im_sdk/internal/group"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/internal/organization"
	"open_im_sdk/internal/signaling"
	"open_im_sdk/internal/super_group"
	"open_im_sdk/internal/user"
	workMoments "open_im_sdk/internal/work_moments"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"sync"
)

type LoginMgr struct {
	organization *organization.Organization
	friend       *friend.Friend
	group        *group.Group
	superGroup   *super_group.SuperGroup
	conversation *conv.Conversation
	user         *user.User
	signaling    *signaling.LiveSignaling
	//advancedFunction advanced_interface.AdvancedFunction
	workMoments  *workMoments.WorkMoments
	full         *full.Full
	db           *db.DataBase
	ws           *ws.Ws
	msgSync      *ws.MsgSync
	heartbeat    *ws.Heartbeat
	cache        *cache.Cache
	token        string
	loginUserID  string
	connListener open_im_sdk_callback.OnConnListener

	loginTime int64

	justOnceFlag bool

	groupListener        open_im_sdk_callback.OnGroupListener
	friendListener       open_im_sdk_callback.OnFriendshipListener
	conversationListener open_im_sdk_callback.OnConversationListener
	advancedMsgListener  open_im_sdk_callback.OnAdvancedMsgListener
	batchMsgListener     open_im_sdk_callback.OnBatchMsgListener
	userListener         open_im_sdk_callback.OnUserListener
	signalingListener    open_im_sdk_callback.OnSignalingListener
	organizationListener open_im_sdk_callback.OnOrganizationListener
	workMomentsListener  open_im_sdk_callback.OnWorkMomentsListener

	conversationCh     chan common.Cmd2Value
	cmdWsCh            chan common.Cmd2Value
	heartbeatCmdCh     chan common.Cmd2Value
	pushMsgAndMaxSeqCh chan common.Cmd2Value
	joinedSuperGroupCh chan common.Cmd2Value
	imConfig           sdk_struct.IMConfig
}

func (u *LoginMgr) Organization() *organization.Organization {
	return u.organization
}

func (u *LoginMgr) Heartbeat() *ws.Heartbeat {
	return u.heartbeat
}

func (u *LoginMgr) Ws() *ws.Ws {
	return u.ws
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

func (u *LoginMgr) Signaling() *signaling.LiveSignaling {
	return u.signaling
}

func (u *LoginMgr) WorkMoments() *workMoments.WorkMoments {
	return u.workMoments
}

func (u *LoginMgr) SetConversationListener(conversationListener open_im_sdk_callback.OnConversationListener) {
	u.conversationListener = conversationListener
}

func (u *LoginMgr) SetAdvancedMsgListener(advancedMsgListener open_im_sdk_callback.OnAdvancedMsgListener) {
	u.advancedMsgListener = advancedMsgListener
}

func (u *LoginMgr) SetBatchMsgListener(batchMsgListener open_im_sdk_callback.OnBatchMsgListener) {
	u.batchMsgListener = batchMsgListener
}
func (u *LoginMgr) SetFriendListener(friendListener open_im_sdk_callback.OnFriendshipListener) {
	u.friendListener = friendListener
}

func (u *LoginMgr) SetGroupListener(groupListener open_im_sdk_callback.OnGroupListener) {
	u.groupListener = groupListener
}

func (u *LoginMgr) SetOrganizationListener(listener open_im_sdk_callback.OnOrganizationListener) {
	u.organizationListener = listener
}

func (u *LoginMgr) SetUserListener(userListener open_im_sdk_callback.OnUserListener) {
	u.userListener = userListener
}

func (u *LoginMgr) SetSignalingListener(listener open_im_sdk_callback.OnSignalingListener) {
	u.signalingListener = listener
}

func (u *LoginMgr) SetWorkMomentsListener(listener open_im_sdk_callback.OnWorkMomentsListener) {
	u.workMomentsListener = listener
}

func (u *LoginMgr) wakeUp(cb open_im_sdk_callback.Base, operationID string) {
	log.Info(operationID, utils.GetSelfFuncName(), "args ")
	err := common.TriggerCmdWakeUp(u.heartbeatCmdCh)
	common.CheckAnyErrCallback(cb, 2001, err, operationID)
	cb.OnSuccess("")
}

func (u *LoginMgr) login(userID, token string, cb open_im_sdk_callback.Base, operationID string) {
	log.Info(operationID, "login start... ", userID, token, sdk_struct.SvrConf)
	err, exp := CheckToken(userID, token, operationID)
	common.CheckTokenErrCallback(cb, err, operationID)
	log.Info(operationID, "checkToken ok ", userID, token, exp)
	u.token = token
	u.loginUserID = userID

	db, err := db.NewDataBase(userID, sdk_struct.SvrConf.DataDir)
	if err != nil {
		cb.OnError(constant.ErrDB.ErrCode, err.Error())
		log.Error(operationID, "NewDataBase failed ", err.Error())
		return
	}
	u.db = db
	log.Info(operationID, "NewDataBase ok ", userID, sdk_struct.SvrConf.DataDir)

	wsRespAsyn := ws.NewWsRespAsyn()
	wsConn := ws.NewWsConn(u.connListener, token, userID)
	u.conversationCh = make(chan common.Cmd2Value, 1000)
	u.cmdWsCh = make(chan common.Cmd2Value, 10)

	u.heartbeatCmdCh = make(chan common.Cmd2Value, 10)

	u.pushMsgAndMaxSeqCh = make(chan common.Cmd2Value, 1000)
	u.ws = ws.NewWs(wsRespAsyn, wsConn, u.cmdWsCh, u.pushMsgAndMaxSeqCh, u.heartbeatCmdCh)
	u.joinedSuperGroupCh = make(chan common.Cmd2Value, 10)
	u.msgSync = ws.NewMsgSync(db, u.ws, userID, u.conversationCh, u.pushMsgAndMaxSeqCh, u.joinedSuperGroupCh)

	u.heartbeat = ws.NewHeartbeat(u.msgSync, u.heartbeatCmdCh, u.connListener, token, exp)

	p := ws.NewPostApi(token, sdk_struct.SvrConf.ApiAddr)

	u.user = user.NewUser(db, p, u.loginUserID)
	u.user.SetListener(u.userListener)

	u.friend = friend.NewFriend(u.loginUserID, u.db, u.user, p)
	u.friend.SetFriendListener(u.friendListener)

	u.group = group.NewGroup(u.loginUserID, u.db, p)
	u.group.SetGroupListener(u.groupListener)
	u.superGroup = super_group.NewSuperGroup(u.loginUserID, u.db, p, u.joinedSuperGroupCh)
	u.organization = organization.NewOrganization(u.loginUserID, u.db, p)
	u.organization.SetListener(u.organizationListener)
	u.cache = cache.NewCache(u.user, u.friend)
	u.full = full.NewFull(u.user, u.friend, u.group, u.conversationCh, u.cache, u.db, u.superGroup)
	u.workMoments = workMoments.NewWorkMoments(u.loginUserID, u.db, p)
	u.workMoments.SetListener(u.workMomentsListener)
	log.NewInfo(operationID, u.imConfig.ObjectStorage)
	u.user.SyncLoginUserInfo(operationID)
	u.loginTime = utils.GetCurrentTimestampByMill()
	u.user.SetLoginTime(u.loginTime)
	u.friend.SetLoginTime(u.loginTime)
	u.group.SetLoginTime(u.loginTime)
	u.superGroup.SetLoginTime(u.loginTime)
	u.organization.SetLoginTime(u.loginTime)
	go u.forcedSynchronization()
	log.Info(operationID, "forcedSynchronization success...")
	log.Info(operationID, "all channel ", u.pushMsgAndMaxSeqCh, u.conversationCh, u.heartbeatCmdCh, u.cmdWsCh)
	log.NewInfo(operationID, u.imConfig.ObjectStorage)
	var objStorage comm2.ObjectStorage
	switch u.imConfig.ObjectStorage {
	case "cos":
		objStorage = comm2.NewCOS(p)
	case "minio":
		objStorage = comm2.NewMinio(p)
	case "oss":
		objStorage = comm2.NewOSS(p)
	default:
		objStorage = comm2.NewCOS(p)
	}
	u.signaling = signaling.NewLiveSignaling(u.ws, u.signalingListener, u.loginUserID, u.imConfig.Platform, u.db)

	u.conversation = conv.NewConversation(u.ws, u.db, p, u.conversationCh,
		u.loginUserID, u.imConfig.Platform, u.imConfig.DataDir,
		u.friend, u.group, u.user, objStorage, u.conversationListener, u.advancedMsgListener,
		u.organization, u.signaling, u.workMoments, u.cache, u.full)
	if u.batchMsgListener != nil {
		u.conversation.SetBatchMsgListener(u.batchMsgListener)
		log.Info(operationID, "SetBatchMsgListener ", u.batchMsgListener)
	}

	u.conversation.SyncConversations(operationID)
	go common.DoListener(u.conversation)
	log.Info(operationID, "login success...")
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
	log.Info(operationID, "TriggerCmd conversationCh UnInit...")
	common.UnInitAll(u.conversationCh)
	if err != nil {
		log.Error(operationID, "TriggerCmd UnInit conversation failed ", err.Error())
	}

	log.Info(operationID, "TriggerCmd pushMsgAndMaxSeqCh UnInit...")
	common.UnInitAll(u.pushMsgAndMaxSeqCh)
	if err != nil {
		log.Error(operationID, "TriggerCmd UnInit pushMsgAndMaxSeqCh failed ", err.Error())
	}

	timeout := 2
	retryTimes := 0
	log.Info(operationID, "send to svr logout ...", u.loginUserID)
	resp, err := u.ws.SendReqWaitResp(&server_api_params.GetMaxAndMinSeqReq{}, constant.WsLogoutMsg, timeout, retryTimes, u.loginUserID, operationID)
	if err != nil {
		log.Warn(operationID, "SendReqWaitResp failed ", err.Error(), constant.WsLogoutMsg, timeout, u.loginUserID, resp)
	}
	if callback != nil {
		callback.OnSuccess("")
	}
	u.justOnceFlag = false

	//go func(mgr *LoginMgr) {
	//	time.Sleep(5 * time.Second)
	//	if mgr == nil {
	//		log.Warn(operationID, "login mgr == nil")
	//		return
	//	}
	//	log.Warn(operationID, "logout close   channel ", mgr.heartbeatCmdCh, mgr.cmdWsCh, mgr.pushMsgAndMaxSeqCh, mgr.conversationCh, mgr.loginUserID)
	//	close(mgr.heartbeatCmdCh)
	//	close(mgr.cmdWsCh)
	//	close(mgr.pushMsgAndMaxSeqCh)
	//	close(mgr.conversationCh)
	//	mgr = nil
	//}(u)
}

func (u *LoginMgr) GetLoginUser() string {
	return u.loginUserID
}

func (u *LoginMgr) GetLoginStatus() int32 {
	return u.ws.LoginState()
}

func (u *LoginMgr) forcedSynchronization() {
	operationID := utils.OperationIDGenerator()
	var wg sync.WaitGroup
	wg.Add(10)

	go func() {
		u.friend.SyncFriendList(operationID)
		wg.Done()
	}()

	go func() {
		u.friend.SyncBlackList(operationID)
		wg.Done()
	}()

	go func() {
		u.friend.SyncFriendApplication(operationID)
		wg.Done()
	}()

	go func() {
		u.friend.SyncSelfFriendApplication(operationID)
		wg.Done()
	}()

	go func() {
		u.group.SyncJoinedGroupList(operationID)
		wg.Done()
	}()

	go func() {
		u.group.SyncAdminGroupApplication(operationID)
		wg.Done()
	}()

	go func() {
		u.group.SyncSelfGroupApplication(operationID)
		wg.Done()
	}()

	go func() {
		u.group.SyncJoinedGroupMemberForFirstLogin(operationID)
		wg.Done()
	}()

	go func() {
		u.organization.SyncOrganization(operationID)
		wg.Done()
	}()

	go func() {
		u.superGroup.SyncJoinedGroupList(operationID)
		wg.Done()
	}()
	wg.Wait()
}

func (u *LoginMgr) GetMinSeqSvr() int64 {
	return u.GetMinSeqSvr()
}

func (u *LoginMgr) SetMinSeqSvr(minSeqSvr int64) {
	u.SetMinSeqSvr(minSeqSvr)
}

func CheckToken(userID, token string, operationID string) (error, uint32) {
	if operationID == "" {
		operationID = utils.OperationIDGenerator()
	}

	log.Debug(operationID, utils.GetSelfFuncName(), userID, token)
	p := ws.NewPostApi(token, sdk_struct.SvrConf.ApiAddr)
	user := user.NewUser(nil, p, userID)
	_, err := user.GetSelfUserInfoFromSvr(operationID)
	if err != nil {
		return utils.Wrap(err, "GetSelfUserInfoFromSvr failed "+operationID), 0
	}
	exp, _ := user.ParseTokenFromSvr(operationID)
	return nil, exp
}

func (u *LoginMgr) uploadImage(callback open_im_sdk_callback.Base, filePath string, token, obj string, operationID string) string {
	p := ws.NewPostApi(token, u.ImConfig().ApiAddr)
	var o comm2.ObjectStorage
	switch obj {
	case "cos":
		o = comm2.NewCOS(p)
	case "minio":
		o = comm2.NewMinio(p)
	default:
		o = comm2.NewCOS(p)
	}
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
}

func (u LoginMgr) uploadFile(callback open_im_sdk_callback.SendMsgCallBack, filePath, operationID string) {
	url, _, err := u.conversation.UploadFile(filePath, callback.OnProgress)
	log.NewInfo(operationID, utils.GetSelfFuncName(), url)
	if err != nil {
		log.Error(operationID, "UploadImage failed ", err.Error(), filePath)
		callback.OnError(constant.ErrApi.ErrCode, err.Error())
	}
	callback.OnSuccess(url)
}
