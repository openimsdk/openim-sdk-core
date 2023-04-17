package login

import (
	"context"
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
	"open_im_sdk/internal/signaling"
	"open_im_sdk/internal/super_group"
	"open_im_sdk/internal/user"
	workMoments "open_im_sdk/internal/work_moments"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"sync"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
)

type LoginMgr struct {
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

func (u *LoginMgr) GetToken() string {
	return u.token
}

func (u *LoginMgr) GetConfig() sdk_struct.IMConfig {
	return u.imConfig
}

func (u *LoginMgr) Push() *comm2.Push {
	return u.push
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

func (u *LoginMgr) SetUserListener(userListener open_im_sdk_callback.OnUserListener) {
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

func (u *LoginMgr) wakeUp(ctx context.Context) error {
	return common.TriggerCmdWakeUp(u.heartbeatCmdCh)
}

func (u *LoginMgr) login(ctx context.Context, userID, token string) error {
	log.ZInfo(ctx, "login start... ", "userID", userID, "token", token, "config", sdk_struct.SvrConf)
	t1 := time.Now()
	u.token = token
	u.loginUserID = userID
	var err error
	u.db, err = db.NewDataBase(ctx, userID, sdk_struct.SvrConf.DataDir)
	if err != nil {
		return errs.ErrDatabase.Wrap(err.Error())
	}
	log.ZInfo(ctx, "NewDataBase ok ", "userID", userID, "dataDir", sdk_struct.SvrConf.DataDir, "login cost time: ", time.Since(t1))

	u.conversationCh = make(chan common.Cmd2Value, 1000)
	u.cmdWsCh = make(chan common.Cmd2Value, 10)

	u.heartbeatCmdCh = make(chan common.Cmd2Value, 10)
	u.pushMsgAndMaxSeqCh = make(chan common.Cmd2Value, 1000)

	u.joinedSuperGroupCh = make(chan common.Cmd2Value, 10)

	u.id2MinSeq = make(map[string]uint32, 100)
	p := ws.NewPostApi(token, sdk_struct.SvrConf.ApiAddr)
	u.postApi = p
	u.user = user.NewUser(u.db, u.loginUserID, u.conversationCh)
	u.user.SetListener(u.userListener)

	u.friend = friend.NewFriend(u.loginUserID, u.db, u.user, p, u.conversationCh)
	u.friend.SetFriendListener(u.friendListener)

	u.group = group.NewGroup(u.loginUserID, u.db, p, u.joinedSuperGroupCh, u.heartbeatCmdCh, u.conversationCh)
	u.group.SetGroupListener(u.groupListener)
	u.superGroup = super_group.NewSuperGroup(u.loginUserID, u.db, p, u.joinedSuperGroupCh, u.heartbeatCmdCh)
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
	log.ZInfo(ctx, u.imConfig.ObjectStorage, "SyncLoginUserInfo login cost time: ", time.Since(t1))
	u.push = comm2.NewPush(p, u.imConfig.Platform, u.loginUserID)
	go u.forcedSynchronization()
	log.ZInfo(ctx, "forcedSynchronization success...", "login cost time: ", time.Since(t1))
	log.ZInfo(ctx, "all channel ", "pushMsgAndMaxSeqCh", u.pushMsgAndMaxSeqCh, "conversationCh", u.conversationCh, "heartbeatCmdCh", u.heartbeatCmdCh, "cmdWsCh", u.cmdWsCh)
	wsConn := ws.NewWsConn(u.connListener, u.token, u.loginUserID, u.imConfig.IsCompression, u.conversationCh)
	wsRespAsyn := ws.NewWsRespAsyn()
	u.ws = ws.NewWs(wsRespAsyn, wsConn, u.cmdWsCh, u.pushMsgAndMaxSeqCh, u.heartbeatCmdCh, u.conversationCh)
	u.msgSync = ws.NewMsgSync(u.db, u.ws, u.loginUserID, u.conversationCh, u.pushMsgAndMaxSeqCh, u.joinedSuperGroupCh)
	u.heartbeat = heartbeart.NewHeartbeat(u.msgSync, u.heartbeatCmdCh, u.connListener, u.token, u.id2MinSeq, u.full)
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
		u.friend, u.group, u.user, objStorage, u.conversationListener, u.advancedMsgListener, u.signaling, u.workMoments, u.business, u.cache, u.full, u.id2MinSeq, u.imConfig.IsExternalExtensions)
	if u.batchMsgListener != nil {
		u.conversation.SetBatchMsgListener(u.batchMsgListener)
		log.ZInfo(ctx, "SetBatchMsgListener", "batchMsgListener", u.batchMsgListener)
	}
	log.ZDebug(ctx, "SyncConversations begin ")
	u.conversation.SyncConversations(ctx, time.Second*2)
	go u.conversation.SyncConversationUnreadCount(mcontext.GetOperationID(ctx))
	go common.DoListener(u.conversation)
	log.ZDebug(ctx, "SyncConversations end ")
	go u.conversation.FixVersionData(ctx)
	log.ZInfo(ctx, "login success...", "login cost time: ", time.Since(t1))
	return nil
}

func (u *LoginMgr) InitSDK(config sdk_struct.IMConfig, listener open_im_sdk_callback.OnConnListener, operationID string) bool {
	u.imConfig = config
	if listener == nil {
		return false
	}
	u.connListener = listener
	return true
}

func (u *LoginMgr) logout(ctx context.Context) error {
	log.ZInfo(ctx, "TriggerCmdLogout ws...")
	if u.friend == nil || u.conversation == nil || u.user == nil || u.full == nil ||
		u.db == nil || u.ws == nil || u.msgSync == nil || u.heartbeat == nil {
		log.ZInfo(ctx, "nil, no TriggerCmdLogout ", "LoginMgr", *u)
		return nil
	}
	operationID := mcontext.GetOperationID(ctx)
	timeout := 2
	retryTimes := 0
	log.ZInfo(ctx, "send to svr logout ...", u.loginUserID)
	resp, err := u.ws.SendReqWaitResp(ctx, &server_api_params.GetMaxAndMinSeqReq{}, constant.WsLogoutMsg, timeout, retryTimes, u.loginUserID)
	if err != nil {
		log.ZWarn(ctx, "SendReqWaitResp failed ", err, "timeout", timeout, "loginUserID", u.loginUserID, "resp", resp)
		if !u.ws.IsInterruptReconnection() {
			return err
		} else {
			log.ZWarn(ctx, "SendReqWaitResp failed, but interrupt reconnection ", err, "timeout", timeout, "loginUserID", u.loginUserID, "resp", resp)
		}
	}
	err = common.TriggerCmdLogout(u.cmdWsCh)
	if err != nil {
		log.ZError(ctx, "TriggerCmdLogout u.cmdWsCh failed ", err)
	}
	log.Info(operationID, "TriggerCmdLogout heartbeat...")
	err = common.TriggerCmdLogout(u.heartbeatCmdCh)
	if err != nil {
		log.ZError(ctx, "TriggerCmdLogout u.heartbeatCmdCh failed ", err)
	}
	err = common.TriggerCmdLogout(u.joinedSuperGroupCh)
	if err != nil {
		log.ZError(ctx, "TriggerCmdLogout u.joinedSuperGroupCh failed ", err)
	}
	log.Info(operationID, "TriggerCmd conversationCh UnInit...")
	common.UnInitAll(u.conversationCh)
	if err != nil {
		log.Error(operationID, "TriggerCmd UnInit conversationCh failed ", err.Error())
	}
	common.UnInitAll(u.pushMsgAndMaxSeqCh)
	if err != nil {
		log.Error(operationID, "TriggerCmd UnInit pushMsgAndMaxSeqCh failed ", err.Error())
	}
	u.db.Close(ctx)
	u.justOnceFlag = false
	return nil
}

func (u *LoginMgr) setAppBackgroundStatus(ctx context.Context, isBackground bool) error {
	timeout := 5
	retryTimes := 2
	_, err := u.ws.SendReqWaitResp(ctx, &server_api_params.SetAppBackgroundStatusReq{UserID: u.loginUserID, IsBackground: isBackground}, constant.WsSetBackgroundStatus, timeout, retryTimes, u.loginUserID)
	return err
}

func (u *LoginMgr) GetLoginUser() string {
	return u.loginUserID
}

func (u *LoginMgr) GetLoginStatus() int32 {
	return u.ws.LoginStatus()
}

func (u *LoginMgr) forcedSynchronization() {
	operationID := utils.OperationIDGenerator()
	ctx := mcontext.NewCtx(operationID)
	log.ZInfo(ctx, "sync all info begin")
	var wg sync.WaitGroup
	wg.Add(9)
	go func() {
		u.user.SyncLoginUserInfo(ctx)
		u.friend.SyncFriendList(ctx)
		wg.Done()
	}()

	go func() {
		u.friend.SyncBlackList(ctx)
		wg.Done()
	}()

	go func() {
		u.friend.SyncFriendApplication(ctx)
		wg.Done()
	}()

	go func() {
		u.friend.SyncSelfFriendApplication(ctx)
		wg.Done()
	}()

	go func() {
		u.group.SyncJoinedGroupList(ctx)
		wg.Done()
	}()

	go func() {
		u.group.SyncAdminGroupApplication(ctx)
		wg.Done()
	}()

	go func() {
		u.group.SyncSelfGroupApplication(ctx)
		wg.Done()
	}()

	go func() {
		u.group.SyncJoinedGroupMemberForFirstLogin(ctx)
		wg.Done()
	}()
	go func() {
		u.superGroup.SyncJoinedGroupList(ctx)
		wg.Done()
	}()
	wg.Wait()

	u.loginTime = utils.GetCurrentTimestampByMill()
	u.user.SetLoginTime(u.loginTime)
	u.friend.SetLoginTime(u.loginTime)
	u.group.SetLoginTime(u.loginTime)
	u.superGroup.SetLoginTime(u.loginTime)
	log.ZInfo(ctx, "login init sync finished")
}

func CheckToken(userID, token string, operationID string) (int64, error) {
	if operationID == "" {
		operationID = utils.OperationIDGenerator()
	}
	ctx := mcontext.NewCtx(operationID)
	log.ZDebug(ctx, utils.GetSelfFuncName(), "userID", userID, "token", token)
	user := user.NewUser(nil, userID, nil)
	exp, err := user.ParseTokenFromSvr(ctx)
	return exp, err
}

func (u *LoginMgr) uploadImage(ctx context.Context, filePath string, token, obj string) (string, error) {
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
	ch := make(chan struct{}, 1)
	f := func(progress int) {
		if progress == 100 {
			ch <- struct{}{}
		}
	}
	url, _, err := o.UploadImage(filePath, f)
	if err != nil {
		return "", err
	}
	for {
		<-ch
		break
	}
	return url, nil
}

func (u LoginMgr) uploadFile(ctx context.Context, filePath string) (string, error) {
	// url, _, err := u.conversation.UploadFile(filePath, callback.OnProgress)
	// // log.NewInfo(operationID, utils.GetSelfFuncName(), url)
	// if err != nil {
	// 	log.Error(operationID, "UploadImage failed ", err.Error(), filePath)
	// 	callback.OnError(constant.ErrApi.ErrCode, err.Error())
	// }
	// callback.OnSuccess(url)
	return "", nil
}
