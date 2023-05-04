package login

import (
	"context"
	"open_im_sdk/internal/business"
	"open_im_sdk/internal/cache"
	conv "open_im_sdk/internal/conversation_msg"
	"open_im_sdk/internal/file"
	"open_im_sdk/internal/friend"
	"open_im_sdk/internal/full"
	"open_im_sdk/internal/group"
	"open_im_sdk/internal/interaction"
	comm2 "open_im_sdk/internal/obj_storage"
	"open_im_sdk/internal/signaling"
	"open_im_sdk/internal/super_group"
	"open_im_sdk/internal/user"
	workMoments "open_im_sdk/internal/work_moments"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/ccontext"
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
	file         *file.File
	signaling    *signaling.LiveSignaling
	//advancedFunction advanced_interface.AdvancedFunction
	workMoments *workMoments.WorkMoments
	business    *business.Business

	full         *full.Full
	db           db_interface.DataBase
	longConnMgr  *interaction.LongConnMgr
	msgSync      *interaction.MsgSync
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
	ctx                context.Context
	cancel             context.CancelFunc
	info               *ccontext.GlobalConfig
	id2MinSeq          map[string]int64
}

func (u *LoginMgr) BaseCtx() context.Context {
	return u.ctx
}

func (u *LoginMgr) Exit() {
	u.cancel()
}

func (u *LoginMgr) GetToken() string {
	return u.token
}

func (u *LoginMgr) Push() *comm2.Push {
	return u.push
}

func (u *LoginMgr) ImConfig() sdk_struct.IMConfig {
	return sdk_struct.IMConfig{
		Platform:             u.info.Platform,
		ApiAddr:              u.info.ApiAddr,
		WsAddr:               u.info.WsAddr,
		DataDir:              u.info.DataDir,
		LogLevel:             u.info.LogLevel,
		EncryptionKey:        u.info.EncryptionKey,
		IsCompression:        u.info.IsCompression,
		IsExternalExtensions: u.info.IsExternalExtensions,
	}
}

func (u *LoginMgr) Conversation() *conv.Conversation {
	return u.conversation
}

func (u *LoginMgr) User() *user.User {
	return u.user
}

func (u *LoginMgr) File() *file.File {
	return u.file
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
	u.info.UserID = userID
	u.info.Token = token
	log.ZInfo(ctx, "login start... ", "userID", userID, "token", token)
	t1 := time.Now()
	u.token = token
	u.loginUserID = userID
	var err error
	u.db, err = db.NewDataBase(ctx, userID, u.info.DataDir)
	if err != nil {
		return errs.ErrDatabase.Wrap(err.Error())
	}
	log.ZDebug(ctx, "NewDataBase ok", "userID", userID, "dataDir", u.info.DataDir, "login cost time", time.Since(t1))
	u.conversationCh = make(chan common.Cmd2Value, 1000)
	u.cmdWsCh = make(chan common.Cmd2Value, 10)

	u.heartbeatCmdCh = make(chan common.Cmd2Value, 10)
	u.pushMsgAndMaxSeqCh = make(chan common.Cmd2Value, 1000)

	u.joinedSuperGroupCh = make(chan common.Cmd2Value, 10)

	u.id2MinSeq = make(map[string]int64, 100)
	u.user = user.NewUser(u.db, u.loginUserID, u.conversationCh)
	u.user.SetListener(u.userListener)

	u.file = file.NewFile(u.db, u.loginUserID)

	u.friend = friend.NewFriend(u.loginUserID, u.db, u.user, u.conversationCh)
	u.friend.SetFriendListener(u.friendListener)

	u.group = group.NewGroup(u.loginUserID, u.db, u.joinedSuperGroupCh, u.heartbeatCmdCh, u.conversationCh)
	u.group.SetGroupListener(u.groupListener)
	u.superGroup = super_group.NewSuperGroup(u.loginUserID, u.db, u.joinedSuperGroupCh, u.heartbeatCmdCh)
	u.cache = cache.NewCache(u.user, u.friend)
	u.full = full.NewFull(u.user, u.friend, u.group, u.conversationCh, u.cache, u.db, u.superGroup)
	u.workMoments = workMoments.NewWorkMoments(u.loginUserID, u.db)
	if u.workMomentsListener != nil {
		u.workMoments.SetListener(u.workMomentsListener)
	}
	u.business = business.NewBusiness(u.db)
	if u.businessListener != nil {
		u.business.SetListener(u.businessListener)
	}
	u.push = comm2.NewPush(u.info.Platform, u.loginUserID)
	log.ZDebug(ctx, "forcedSynchronization success...", "login cost time: ", time.Since(t1))
	u.longConnMgr = interaction.NewLongConnMgr(ctx, u.connListener, u.pushMsgAndMaxSeqCh, u.conversationCh)
	//wsConn := ws.NewWsConn(u.connListener, u.token, u.loginUserID, u.imConfig.IsCompression, u.conversationCh)
	//wsRespAsyn := ws.NewWsRespAsyn()
	//u.ws = ws.NewWs(wsRespAsyn, wsConn, u.cmdWsCh, u.pushMsgAndMaxSeqCh, u.heartbeatCmdCh, u.conversationCh)
	u.msgSync = interaction.NewMsgSync(ctx, u.db, u.conversationCh, u.pushMsgAndMaxSeqCh)
	//u.heartbeat = heartbeart.NewHeartbeat(u.msgSync, u.heartbeatCmdCh, u.connListener, u.token, u.id2MinSeq, u.full)
	//var objStorage comm3.ObjectStorage
	//switch u.imConfig.ObjectStorage {
	//case "cos":
	//	objStorage = comm2.NewCOS(u.postApi)
	//case "minio":
	//	objStorage = comm2.NewMinio(u.postApi)
	//case "oss":
	//	objStorage = comm2.NewOSS(u.postApi)
	//case "aws":
	//	objStorage = comm2.NewAWS(u.postApi)
	//default:
	//	objStorage = comm2.NewCOS(u.postApi)
	//}
	u.conversation = conv.NewConversation(ctx, u.db, u.conversationCh,
		u.friend, u.group, u.user, u.conversationListener, u.advancedMsgListener, u.signaling, u.workMoments, u.business, u.cache, u.full, u.id2MinSeq)
	//var wg sync.WaitGroup
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	u.forcedSynchronization(ctx)
	//}()
	//wg.Wait()
	log.ZDebug(ctx, "forcedSynchronization success...", "login cost time: ", time.Since(t1))
	//u.ws = ws.NewWs(wsRespAsyn, wsConn, u.cmdWsCh, u.pushMsgAndMaxSeqCh, u.heartbeatCmdCh, u.conversationCh)
	//u.msgSync = ws.NewMsgSync(u.db, u.ws, u.loginUserID, u.conversationCh, u.pushMsgAndMaxSeqCh, u.joinedSuperGroupCh)
	//u.heartbeat = heartbeart.NewHeartbeat(u.msgSync, u.heartbeatCmdCh, u.connListener, u.token, u.id2MinSeq, u.full)

	u.signaling = signaling.NewLiveSignaling(u.longConnMgr, u.loginUserID, u.info.Platform, u.db)

	if u.signalingListener != nil {
		u.signaling.SetListener(u.signalingListener)
	}
	if u.signalingListenerFroService != nil {
		u.signaling.SetListenerForService(u.signalingListenerFroService)
	}

	if u.batchMsgListener != nil {
		u.conversation.SetBatchMsgListener(u.batchMsgListener)
		log.ZDebug(ctx, "SetBatchMsgListener", "batchMsgListener", u.batchMsgListener)
	}
	go common.DoListener(u.conversation)
	go u.conversation.FixVersionData(ctx)
	log.ZInfo(ctx, "login success...", "login cost time: ", time.Since(t1))
	return nil
}

func (u *LoginMgr) InitSDK(config sdk_struct.IMConfig, listener open_im_sdk_callback.OnConnListener, operationID string) bool {
	if listener == nil {
		return false
	}
	u.info = &ccontext.GlobalConfig{
		Platform: config.Platform,
		ApiAddr:  config.ApiAddr,
		WsAddr:   config.WsAddr,
		DataDir:  config.DataDir,
		LogLevel: config.LogLevel,
		//ObjectStorage:        config.ObjectStorage,
		EncryptionKey:        config.EncryptionKey,
		IsCompression:        config.IsCompression,
		IsExternalExtensions: config.IsExternalExtensions,
	}
	u.connListener = listener
	ctx := ccontext.WithInfo(context.Background(), u.info)
	u.ctx, u.cancel = context.WithCancel(ctx)
	return true
}

func (u *LoginMgr) logout(ctx context.Context) error {
	err := u.longConnMgr.SendReqWaitResp(ctx, &server_api_params.GetMaxAndMinSeqReq{}, constant.LogoutMsg, nil)
	if err != nil {
		return err
	}
	u.Exit()
	log.ZDebug(ctx, "TriggerCmdLogout success...")
	return nil
}

func (u *LoginMgr) setAppBackgroundStatus(ctx context.Context, isBackground bool) error {
	return u.longConnMgr.SendReqWaitResp(ctx, &server_api_params.SetAppBackgroundStatusReq{UserID: u.loginUserID, IsBackground: isBackground}, constant.SetBackgroundStatus, nil)

}

func (u *LoginMgr) GetLoginUser() string {
	return u.loginUserID
}

func (u *LoginMgr) GetLoginStatus() int {
	return u.longConnMgr.GetConnectionStatus()
}

func (u *LoginMgr) forcedSynchronization(ctx context.Context) {
	log.ZInfo(ctx, "sync all info begin")
	var wg sync.WaitGroup
	var errCh = make(chan error, 12)
	wg.Add(12)
	go func() {
		defer wg.Done()
		if err := u.user.SyncLoginUserInfo(ctx); err != nil {
			errCh <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := u.friend.SyncFriendList(ctx); err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := u.friend.SyncBlackList(ctx); err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := u.friend.SyncFriendApplication(ctx); err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := u.friend.SyncSelfFriendApplication(ctx); err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := u.group.SyncJoinedGroup(ctx); err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := u.group.SyncAdminGroupApplication(ctx); err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := u.group.SyncSelfGroupApplication(ctx); err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := u.group.SyncJoinedGroupMemberForFirstLogin(ctx); err != nil {
			errCh <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := u.superGroup.SyncJoinedGroupList(ctx); err != nil {
			errCh <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := u.conversation.SyncConversations(ctx); err != nil {
			errCh <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := u.conversation.SyncConversationUnreadCount(ctx); err != nil {
			errCh <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()
	for err := range errCh {
		if err != nil {
			log.ZError(ctx, "sync info failed", err)
		}
	}
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
