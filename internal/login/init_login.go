package login

import (
	"open_im_sdk/internal/business"
	"open_im_sdk/internal/cache"
	comm3 "open_im_sdk/internal/common"
	conv "open_im_sdk/internal/conversation_msg"
	"open_im_sdk/internal/friend"
	"open_im_sdk/internal/full"
	"open_im_sdk/internal/group"
	"open_im_sdk/internal/heartbeart"
	ws "open_im_sdk/internal/interaction"
	comm2 "open_im_sdk/internal/obj_storage"
	"open_im_sdk/internal/organization"
	"open_im_sdk/internal/signaling"
	"open_im_sdk/internal/super_group"
	"open_im_sdk/internal/user"
	workMoments "open_im_sdk/internal/work_moments"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"sync"
	"time"
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
	workMoments *workMoments.WorkMoments
	business    *business.Business

	full         *full.Full
	db           db_interface.DataBase
	ws           *ws.Ws
	msgSync      *ws.MsgSync
	heartbeat    *heartbeart.Heartbeat
	push         *comm2.Push
	cache        *cache.Cache
	token        string
	loginUserID  string
	connListener open_im_sdk_callback.OnConnListener

	loginTime int64

	justOnceFlag bool

	groupListener               open_im_sdk_callback.OnGroupListener
	friendListener              open_im_sdk_callback.OnFriendshipListener
	conversationListener        open_im_sdk_callback.OnConversationListener
	advancedMsgListener         open_im_sdk_callback.OnAdvancedMsgListener
	batchMsgListener            open_im_sdk_callback.OnBatchMsgListener
	userListener                open_im_sdk_callback.OnUserListener
	signalingListener           open_im_sdk_callback.OnSignalingListener
	signalingListenerFroService open_im_sdk_callback.OnSignalingListener
	organizationListener        open_im_sdk_callback.OnOrganizationListener
	workMomentsListener         open_im_sdk_callback.OnWorkMomentsListener
	businessListener            open_im_sdk_callback.OnCustomBusinessListener

	conversationCh     chan common.Cmd2Value
	cmdWsCh            chan common.Cmd2Value
	heartbeatCmdCh     chan common.Cmd2Value
	pushMsgAndMaxSeqCh chan common.Cmd2Value
	joinedSuperGroupCh chan common.Cmd2Value
	imConfig           sdk_struct.IMConfig

	id2MinSeq map[string]uint32
	postApi   *ws.PostApi
}

func (u *LoginMgr) Push() *comm2.Push {
	return u.push
}

func (u *LoginMgr) Organization() *organization.Organization {
	return u.organization
}

func (u *LoginMgr) Heartbeat() *heartbeart.Heartbeat {
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
	if u.conversation != nil {
		u.conversation.SetConversationListener(conversationListener)
	} else {
		u.conversationListener = conversationListener
	}
}

func (u *LoginMgr) SetAdvancedMsgListener(advancedMsgListener open_im_sdk_callback.OnAdvancedMsgListener) {
	if u.conversation != nil {
		u.conversation.SetMsgListener(advancedMsgListener)
	} else {
		u.advancedMsgListener = advancedMsgListener
	}
}
func (u *LoginMgr) SetMessageKvInfoListener(messageKvInfoListener open_im_sdk_callback.OnMessageKvInfoListener) {
	if u.conversation != nil {
		u.conversation.SetMsgKvListener(messageKvInfoListener)
	}
}
func (u *LoginMgr) SetBatchMsgListener(batchMsgListener open_im_sdk_callback.OnBatchMsgListener) {
	if u.conversation != nil {
		u.conversation.SetBatchMsgListener(batchMsgListener)
	} else {
		u.batchMsgListener = batchMsgListener
	}
}
func (u *LoginMgr) SetFriendListener(friendListener open_im_sdk_callback.OnFriendshipListener) {
	if u.friend != nil {
		u.friend.SetFriendListener(friendListener)
	} else {
		u.friendListener = friendListener
	}
}

func (u *LoginMgr) SetGroupListener(groupListener open_im_sdk_callback.OnGroupListener) {
	if u.group != nil {
		u.group.SetGroupListener(groupListener)
	} else {
		u.groupListener = groupListener
	}
}

func (u *LoginMgr) SetOrganizationListener(listener open_im_sdk_callback.OnOrganizationListener) {
	if u.organization != nil {
		u.organization.SetListener(listener)
	} else {
		u.organizationListener = listener
	}
}

func (u *LoginMgr) SetUserListener(userListener open_im_sdk_callback.OnUserListener) {
	//if u.signaling != nil {
	//		u.signaling.SetListener(listener)
	//	} else {
	//		u.signalingListener = listener
	//	}

	if u.user != nil {
		u.user.SetListener(userListener)
	} else {
		u.userListener = userListener
	}
}

func (u *LoginMgr) SetSignalingListener(listener open_im_sdk_callback.OnSignalingListener) {
	if u.signaling != nil {
		u.signaling.SetListener(listener)
	} else {
		u.signalingListener = listener
	}
}

func (u *LoginMgr) SetSignalingListenerForService(listener open_im_sdk_callback.OnSignalingListener) {
	if u.signaling != nil {
		u.signaling.SetListenerForService(listener)
	} else {
		u.signalingListenerFroService = listener
	}
}

func (u *LoginMgr) SetListenerForService(listener open_im_sdk_callback.OnListenerForService) {
	if u.friend == nil || u.group == nil || u.conversation == nil {
		log.Error("", "is nil ", u.friend, u.group, u.conversation)
		return
	}
	u.friend.SetListenerForService(listener)
	u.group.SetListenerForService(listener)
	u.conversation.SetListenerForService(listener)
}

func (u *LoginMgr) SetWorkMomentsListener(listener open_im_sdk_callback.OnWorkMomentsListener) {
	if u.workMoments != nil {
		u.workMoments.SetListener(listener)
	} else {
		u.workMomentsListener = listener
	}
}

func (u *LoginMgr) SetBusinessListener(listener open_im_sdk_callback.OnCustomBusinessListener) {
	if u.business != nil {
		u.business.SetListener(listener)
	} else {
		u.businessListener = listener
	}
}

func (u *LoginMgr) wakeUp(cb open_im_sdk_callback.Base, operationID string) {
	log.Info(operationID, utils.GetSelfFuncName(), "args ")
	err := common.TriggerCmdWakeUp(u.heartbeatCmdCh)
	common.CheckAnyErrCallback(cb, 2001, err, operationID)
	cb.OnSuccess("")
}

func (u *LoginMgr) login(userID, token string, cb open_im_sdk_callback.Base, operationID string) {
	log.Info(operationID, "login start... ", userID, token, sdk_struct.SvrConf)
	t1 := time.Now()
	u.token = token
	u.loginUserID = userID
	var sqliteConn *db.DataBase
	var err error
	if constant.OnlyForTest == 1 {
		wsConn := ws.NewWsConn(u.connListener, u.token, u.loginUserID, u.imConfig.IsCompression, u.conversationCh)
		wsRespAsyn := ws.NewWsRespAsyn()
		u.ws = ws.NewWs(wsRespAsyn, wsConn, u.cmdWsCh, u.pushMsgAndMaxSeqCh, u.heartbeatCmdCh, u.conversationCh)
		u.heartbeat = heartbeart.NewHeartbeat(u.msgSync, u.heartbeatCmdCh, u.connListener, u.token, u.id2MinSeq, u.full)
		u.heartbeat.WsForTest = u.ws
		u.heartbeat.LoginUserIDForTest = u.loginUserID
		cb.OnSuccess("")
		return
	}

	sqliteConn, err = db.NewDataBase(userID, sdk_struct.SvrConf.DataDir, operationID)
	if err != nil {
		cb.OnError(constant.ErrDB.ErrCode, err.Error())
		log.Error(operationID, "NewDataBase failed ", err.Error())
		return
	}

	u.db = sqliteConn
	log.Info(operationID, "NewDataBase ok ", userID, sdk_struct.SvrConf.DataDir, "login cost time: ", time.Since(t1))

	u.conversationCh = make(chan common.Cmd2Value, 1000)
	u.cmdWsCh = make(chan common.Cmd2Value, 10)

	u.heartbeatCmdCh = make(chan common.Cmd2Value, 10)
	u.pushMsgAndMaxSeqCh = make(chan common.Cmd2Value, 1000)

	u.joinedSuperGroupCh = make(chan common.Cmd2Value, 10)

	u.id2MinSeq = make(map[string]uint32, 100)
	p := ws.NewPostApi(token, sdk_struct.SvrConf.ApiAddr)
	u.postApi = p
	u.user = user.NewUser(sqliteConn, p, u.loginUserID, u.conversationCh)
	u.user.SetListener(u.userListener)

	u.friend = friend.NewFriend(u.loginUserID, u.db, u.user, p, u.conversationCh)
	u.friend.SetFriendListener(u.friendListener)

	u.group = group.NewGroup(u.loginUserID, u.db, p, u.joinedSuperGroupCh, u.heartbeatCmdCh, u.conversationCh)
	u.group.SetGroupListener(u.groupListener)
	u.superGroup = super_group.NewSuperGroup(u.loginUserID, u.db, p, u.joinedSuperGroupCh, u.heartbeatCmdCh)
	u.organization = organization.NewOrganization(u.loginUserID, u.db, p)
	u.organization.SetListener(u.organizationListener)
	u.cache = cache.NewCache(u.user, u.friend)
	u.full = full.NewFull(u.user, u.friend, u.group, u.conversationCh, u.cache, u.db, u.superGroup)
	u.workMoments = workMoments.NewWorkMoments(u.loginUserID, u.db, p)
	if u.workMomentsListener != nil {
		u.workMoments.SetListener(u.workMomentsListener)
	}
	u.business = business.NewBusiness(u.db)
	if u.businessListener != nil {
		u.business.SetListener(u.businessListener)
	}
	log.NewInfo(operationID, u.imConfig.ObjectStorage, "new obj login cost time: ", time.Since(t1))
	log.NewInfo(operationID, u.imConfig.ObjectStorage, "SyncLoginUserInfo login cost time: ", time.Since(t1))
	u.push = comm2.NewPush(p, u.imConfig.Platform, u.loginUserID)
	go u.forcedSynchronization()

	log.Info(operationID, "forcedSynchronization success...", "login cost time: ", time.Since(t1))
	log.Info(operationID, "all channel ", u.pushMsgAndMaxSeqCh, u.conversationCh, u.heartbeatCmdCh, u.cmdWsCh)

	wsConn := ws.NewWsConn(u.connListener, u.token, u.loginUserID, u.imConfig.IsCompression, u.conversationCh)
	wsRespAsyn := ws.NewWsRespAsyn()
	u.ws = ws.NewWs(wsRespAsyn, wsConn, u.cmdWsCh, u.pushMsgAndMaxSeqCh, u.heartbeatCmdCh, u.conversationCh)
	u.msgSync = ws.NewMsgSync(u.db, u.ws, u.loginUserID, u.conversationCh, u.pushMsgAndMaxSeqCh, u.joinedSuperGroupCh)
	u.heartbeat = heartbeart.NewHeartbeat(u.msgSync, u.heartbeatCmdCh, u.connListener, u.token, u.id2MinSeq, u.full)
	log.NewInfo(operationID, u.imConfig.ObjectStorage)

	var objStorage comm3.ObjectStorage
	switch u.imConfig.ObjectStorage {
	case "cos":
		objStorage = comm2.NewCOS(u.postApi)
	case "minio":
		objStorage = comm2.NewMinio(u.postApi)
	case "oss":
		objStorage = comm2.NewOSS(u.postApi)
	case "aws":
		objStorage = comm2.NewAWS(u.postApi)
	default:
		objStorage = comm2.NewCOS(u.postApi)
	}
	u.signaling = signaling.NewLiveSignaling(u.ws, u.loginUserID, u.imConfig.Platform, u.db)
	if u.signalingListener != nil {
		u.signaling.SetListener(u.signalingListener)
	}
	if u.signalingListenerFroService != nil {
		u.signaling.SetListenerForService(u.signalingListenerFroService)
	}
	u.conversation = conv.NewConversation(u.ws, u.db, u.postApi, u.conversationCh,
		u.loginUserID, u.imConfig.Platform, u.imConfig.DataDir, u.imConfig.EncryptionKey,
		u.friend, u.group, u.user, objStorage, u.conversationListener, u.advancedMsgListener,
		u.organization, u.signaling, u.workMoments, u.business, u.cache, u.full, u.id2MinSeq, u.imConfig.IsExternalExtensions)
	if u.batchMsgListener != nil {
		u.conversation.SetBatchMsgListener(u.batchMsgListener)
		log.Info(operationID, "SetBatchMsgListener ", u.batchMsgListener)
	}
	log.Debug(operationID, "SyncConversations begin ")
	u.conversation.SyncConversations(operationID, time.Second*2)
	go u.conversation.SyncConversationUnreadCount(operationID)
	go common.DoListener(u.conversation)
	log.Debug(operationID, "SyncConversations end ")
	go u.conversation.FixVersionData()
	log.Info(operationID, "ws heartbeat end ")

	log.Info(operationID, "login success...", "login cost time: ", time.Since(t1))
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

	timeout := 2
	retryTimes := 0
	log.Info(operationID, "send to svr logout ...", u.loginUserID)
	resp, err := u.ws.SendReqWaitResp(&server_api_params.GetMaxAndMinSeqReq{}, constant.WsLogoutMsg, timeout, retryTimes, u.loginUserID, operationID)
	if err != nil {
		log.Warn(operationID, "SendReqWaitResp failed ", err.Error(), constant.WsLogoutMsg, timeout, u.loginUserID, resp)
		if !u.ws.IsInterruptReconnection() {
			callback.OnError(100, err.Error())
			return
		} else {
			log.Warn(operationID, "SendReqWaitResp failed, but interrupt reconnection ", err.Error(), constant.WsLogoutMsg, timeout, u.loginUserID, resp)
		}
	}

	err = common.TriggerCmdLogout(u.cmdWsCh)
	if err != nil {
		log.Error(operationID, "TriggerCmdLogout failed ", err.Error())
	}
	log.Info(operationID, "TriggerCmdLogout heartbeat...")
	err = common.TriggerCmdLogout(u.heartbeatCmdCh)
	if err != nil {
		log.Error(operationID, "TriggerCmdLogout failed ", err.Error())
	}
	err = common.TriggerCmdLogout(u.joinedSuperGroupCh)
	if err != nil {
		log.Error(operationID, "TriggerCmdLogout  joinedSuperGroupCh failed ", err.Error())
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

	log.Info(operationID, "close db ")
	u.db.Close()

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

func (u *LoginMgr) setAppBackgroundStatus(callback open_im_sdk_callback.Base, isBackground bool, operationID string) {
	timeout := 5
	retryTimes := 2
	log.Info(operationID, "send to svr WsSetBackgroundStatus ...", u.loginUserID)
	resp, err := u.ws.SendReqWaitResp(&server_api_params.SetAppBackgroundStatusReq{UserID: u.loginUserID, IsBackground: isBackground}, constant.WsSetBackgroundStatus, timeout, retryTimes, u.loginUserID, operationID)
	if err != nil {
		log.Error(operationID, "SendReqWaitResp failed ", err.Error(), constant.WsSetBackgroundStatus, timeout, u.loginUserID, resp)
	}
	common.CheckAnyErrCallback(callback, constant.ErrInternal.ErrCode, err, operationID)
	callback.OnSuccess("")
}

func (u *LoginMgr) GetLoginUser() string {
	return u.loginUserID
}

func (u *LoginMgr) GetLoginStatus() int32 {
	return u.ws.LoginStatus()
}

func (u *LoginMgr) forcedSynchronization() {
	operationID := utils.OperationIDGenerator()

	log.Info(operationID, "sync all info begin")
	var wg sync.WaitGroup
	wg.Add(10)
	go func() {
		u.user.SyncLoginUserInfo(operationID)
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
	if u.organizationListener != nil {
		go func() {
			u.organization.SyncOrganization(operationID)
			wg.Done()
		}()
	} else {
		wg.Done()
	}
	go func() {
		u.superGroup.SyncJoinedGroupList(operationID)
		wg.Done()
	}()
	wg.Wait()

	u.loginTime = utils.GetCurrentTimestampByMill()
	u.user.SetLoginTime(u.loginTime)
	u.friend.SetLoginTime(u.loginTime)
	u.group.SetLoginTime(u.loginTime)
	u.superGroup.SetLoginTime(u.loginTime)
	u.organization.SetLoginTime(u.loginTime)
	log.Info(operationID, "login init sync finished")
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
	user := user.NewUser(nil, p, userID, nil)
	//_, err := user.GetSelfUserInfoFromSvr(operationID)
	//if err != nil {
	//	return utils.Wrap(err, "GetSelfUserInfoFromSvr failed "+operationID), 0
	//}
	exp, err := user.ParseTokenFromSvr(operationID)
	return err, exp
}

func (u *LoginMgr) uploadImage(callback open_im_sdk_callback.Base, filePath string, token, obj string, operationID string) string {
	p := ws.NewPostApi(token, u.ImConfig().ApiAddr)
	var o comm3.ObjectStorage
	switch obj {
	case "cos":
		o = comm2.NewCOS(p)
	case "minio":
		o = comm2.NewMinio(p)
	case "aws":
		o = comm2.NewAWS(p)
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
