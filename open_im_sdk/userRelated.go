// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package open_im_sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/openimsdk/openim-sdk-core/v3/internal/third/file"
	"github.com/openimsdk/tools/errs"

	"github.com/openimsdk/openim-sdk-core/v3/internal/relation"

	conv "github.com/openimsdk/openim-sdk-core/v3/internal/conversation_msg"
	"github.com/openimsdk/openim-sdk-core/v3/internal/group"
	"github.com/openimsdk/openim-sdk-core/v3/internal/interaction"
	"github.com/openimsdk/openim-sdk-core/v3/internal/third"
	"github.com/openimsdk/openim-sdk-core/v3/internal/user"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/protocol/push"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/jsonutil"
)

const (
	LogoutStatus = iota + 1
	Logging
	Logged
)

const (
	LogoutTips = "js sdk socket close"
)

var (
	// IMUserContext is the global user context instance
	IMUserContext *UserContext
	once          sync.Once
)

func init() {
	IMUserContext = NewIMUserContext()
	IMUserContext.initResources()

}

func (u *UserContext) InitResources() {
	u.initResources()
}

func (u *UserContext) initResources() {
	ctx := ccontext.WithInfo(context.Background(), u.info)
	u.ctx, u.cancel = context.WithCancel(ctx)
	u.setFGCtx()
	u.conversationEventQueue = make(chan common.Cmd2Value, 1000)
	u.msgSyncerCh = make(chan common.Cmd2Value, 1000)
	u.loginMgrCh = make(chan common.Cmd2Value, 1)

	u.longConnMgr = interaction.NewLongConnMgr(u.ctx, u.userOnlineStatusChange, u.msgSyncerCh, u.loginMgrCh)
	u.ctx = ccontext.WithApiErrCode(u.ctx, &apiErrCallback{loginMgrCh: u.loginMgrCh, listener: u.ConnListener})
	u.setLoginStatus(LogoutStatus)
	u.user = user.NewUser(u.conversationEventQueue)
	u.file = file.NewFile()
	u.relation = relation.NewRelation(u.conversationEventQueue, u.user)
	u.group = group.NewGroup(u.conversationEventQueue)
	u.third = third.NewThird(u.file)
	u.msgSyncer = interaction.NewMsgSyncer(u.conversationEventQueue, u.msgSyncerCh, u.longConnMgr)
	u.conversation = conv.NewConversation(u.longConnMgr, u.msgSyncerCh, u.conversationEventQueue,
		u.relation, u.group, u.user, u.file)
	u.setListener(ctx)
}

// CheckResourceLoad checks the SDK is resource load status.
func CheckResourceLoad(userContext *UserContext, funcName string) error {
	if userContext.Info().IMConfig == nil {
		return sdkerrs.ErrSDKNotInit.WrapMsg(funcName)
	}
	if funcName == "" {
		return nil
	}
	parts := strings.Split(funcName, ".")
	if parts[len(parts)-1] == "Login-fm" {
		return nil
	}
	if userContext.getLoginStatus(context.Background()) != Logged {
		return sdkerrs.ErrSDKNotLogin.WrapMsg(funcName)
	}
	return nil
}

type UserContext struct {
	relation     *relation.Relation
	group        *group.Group
	conversation *conv.Conversation
	user         *user.User
	file         *file.File

	db          db_interface.DataBase
	longConnMgr *interaction.LongConnMgr
	msgSyncer   *interaction.MsgSyncer
	third       *third.Third
	token       string
	loginUserID string

	justOnceFlag bool

	w           sync.Mutex
	loginStatus int

	connListener         open_im_sdk_callback.OnConnListener
	groupListener        open_im_sdk_callback.OnGroupListener
	friendshipListener   open_im_sdk_callback.OnFriendshipListener
	conversationListener open_im_sdk_callback.OnConversationListener
	advancedMsgListener  open_im_sdk_callback.OnAdvancedMsgListener
	userListener         open_im_sdk_callback.OnUserListener
	signalingListener    open_im_sdk_callback.OnSignalingListener
	businessListener     open_im_sdk_callback.OnCustomBusinessListener
	msgKvListener        open_im_sdk_callback.OnMessageKvInfoListener

	//conversationCh chan common.Cmd2Value

	conversationEventQueue chan common.Cmd2Value
	cmdWsCh                chan common.Cmd2Value
	msgSyncerCh            chan common.Cmd2Value
	loginMgrCh             chan common.Cmd2Value

	ctx       context.Context
	cancel    context.CancelFunc
	fgCtx     context.Context
	fgCancel  context.CancelCauseFunc
	info      *ccontext.GlobalConfig
	id2MinSeq map[string]int64
}

func (u *UserContext) Info() *ccontext.GlobalConfig {
	return u.info
}

func (u *UserContext) ConnListener() open_im_sdk_callback.OnConnListener {
	return u.connListener
}

func (u *UserContext) GroupListener() open_im_sdk_callback.OnGroupListener {
	return u.groupListener
}

func (u *UserContext) FriendshipListener() open_im_sdk_callback.OnFriendshipListener {
	return u.friendshipListener
}

func (u *UserContext) ConversationListener() open_im_sdk_callback.OnConversationListener {
	return u.conversationListener
}

func (u *UserContext) AdvancedMsgListener() open_im_sdk_callback.OnAdvancedMsgListener {
	return u.advancedMsgListener
}

func (u *UserContext) UserListener() open_im_sdk_callback.OnUserListener {
	return u.userListener
}

func (u *UserContext) SignalingListener() open_im_sdk_callback.OnSignalingListener {
	return u.signalingListener
}

func (u *UserContext) BusinessListener() open_im_sdk_callback.OnCustomBusinessListener {
	return u.businessListener
}

func (u *UserContext) MsgKvListener() open_im_sdk_callback.OnMessageKvInfoListener {
	return u.msgKvListener
}

func (u *UserContext) Exit() {
	u.cancel()
}

func (u *UserContext) Third() *third.Third {
	return u.third
}

func (u *UserContext) ImConfig() sdk_struct.IMConfig {
	return sdk_struct.IMConfig{
		PlatformID: u.info.PlatformID,
		ApiAddr:    u.info.ApiAddr,
		WsAddr:     u.info.WsAddr,
		DataDir:    u.info.DataDir,
		LogLevel:   u.info.LogLevel,
	}
}

func (u *UserContext) Conversation() *conv.Conversation {
	return u.conversation
}

func (u *UserContext) User() *user.User {
	return u.user
}

func (u *UserContext) File() *file.File {
	return u.file
}

func (u *UserContext) Group() *group.Group {
	return u.group
}

func (u *UserContext) Relation() *relation.Relation {
	return u.relation
}

func (u *UserContext) SetConversationListener(conversationListener open_im_sdk_callback.OnConversationListener) {
	u.conversationListener = conversationListener
}

func (u *UserContext) SetAdvancedMsgListener(advancedMsgListener open_im_sdk_callback.OnAdvancedMsgListener) {
	u.advancedMsgListener = advancedMsgListener
}

func (u *UserContext) SetMessageKvInfoListener(messageKvInfoListener open_im_sdk_callback.OnMessageKvInfoListener) {
	u.msgKvListener = messageKvInfoListener
}

func (u *UserContext) SetFriendshipListener(friendshipListener open_im_sdk_callback.OnFriendshipListener) {
	u.friendshipListener = friendshipListener
}

func (u *UserContext) SetGroupListener(groupListener open_im_sdk_callback.OnGroupListener) {
	u.groupListener = groupListener
}

func (u *UserContext) SetUserListener(userListener open_im_sdk_callback.OnUserListener) {
	u.userListener = userListener
}

func (u *UserContext) SetCustomBusinessListener(listener open_im_sdk_callback.OnCustomBusinessListener) {
	u.businessListener = listener
}

func (u *UserContext) GetLoginUserID() string {
	return u.loginUserID
}

func (u *UserContext) logoutListener(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("panic: %+v\n%s", r, debug.Stack())

			log.ZWarn(ctx, "logoutListener panic", nil, "panic info", err)
		}
	}()

	for {
		select {
		case <-u.loginMgrCh:
			log.ZDebug(ctx, "logoutListener exit")
			err := u.logout(ctx, true)
			if err != nil {
				log.ZError(ctx, "logout error", err)
			}
		case <-ctx.Done():
			log.ZInfo(ctx, "logoutListener done sdk logout.....")
			return
		}
	}

}

func NewIMUserContext() *UserContext {
	once.Do(func() {
		IMUserContext = &UserContext{
			info: &ccontext.GlobalConfig{},
		}
	})
	return IMUserContext
}

func NewLoginMgr() *UserContext {
	return &UserContext{
		info: &ccontext.GlobalConfig{},
	}
}
func (u *UserContext) getLoginStatus(_ context.Context) int {
	u.w.Lock()
	defer u.w.Unlock()
	return u.loginStatus
}
func (u *UserContext) setLoginStatus(status int) {
	u.w.Lock()
	defer u.w.Unlock()
	u.loginStatus = status
}
func (u *UserContext) checkSendingMessage(ctx context.Context) {
	sendingMessages, err := u.db.GetAllSendingMessages(ctx)
	if err != nil {
		log.ZError(ctx, "GetAllSendingMessages failed", err)
	}
	for _, message := range sendingMessages {
		if err := u.handlerSendingMsg(ctx, message); err != nil {
			log.ZError(ctx, "handlerSendingMsg failed", err, "message", message)
		}
		if err := u.db.DeleteSendingMessage(ctx, message.ConversationID, message.ClientMsgID); err != nil {
			log.ZError(ctx, "DeleteSendingMessage failed", err, "conversationID", message.ConversationID, "clientMsgID", message.ClientMsgID)
		}
	}
}

func (u *UserContext) handlerSendingMsg(ctx context.Context, sendingMsg *model_struct.LocalSendingMessages) error {
	tableMessage, err := u.db.GetMessage(ctx, sendingMsg.ConversationID, sendingMsg.ClientMsgID)
	if err != nil {
		return err
	}
	if tableMessage.Status != constant.MsgStatusSending {
		return nil
	}
	err = u.db.UpdateMessage(ctx, sendingMsg.ConversationID, &model_struct.LocalChatLog{ClientMsgID: sendingMsg.ClientMsgID, Status: constant.MsgStatusSendFailed})
	if err != nil {
		return err
	}
	conversation, err := u.db.GetConversation(ctx, sendingMsg.ConversationID)
	if err != nil {
		return err
	}
	latestMsg := &sdk_struct.MsgStruct{}
	if err := json.Unmarshal([]byte(conversation.LatestMsg), &latestMsg); err != nil {
		return err
	}
	if latestMsg.ClientMsgID == sendingMsg.ClientMsgID {
		latestMsg.Status = constant.MsgStatusSendFailed
		conversation.LatestMsg = jsonutil.StructToJsonString(latestMsg)
		return u.db.UpdateConversation(ctx, conversation)
	}
	return nil
}

func (u *UserContext) login(ctx context.Context, userID, token string) error {
	if u.getLoginStatus(ctx) == Logged {
		return sdkerrs.ErrLoginRepeat
	}
	u.setLoginStatus(Logging)
	log.ZDebug(ctx, "login start... ", "userID", userID, "token", token)
	t1 := time.Now()

	u.info.UserID = userID
	u.info.Token = token

	if err := u.initialize(ctx, userID); err != nil {
		return err
	}

	u.run(ctx)
	u.setLoginStatus(Logged)
	log.ZDebug(ctx, "login success...", "login cost time: ", time.Since(t1))
	return nil
}

func (u *UserContext) initialize(ctx context.Context, userID string) error {
	var err error
	u.db, err = db.NewDataBase(ctx, userID, u.info.DataDir, int(u.info.LogLevel))
	if err != nil {
		return sdkerrs.ErrSdkInternal.WrapMsg("init database " + err.Error())
	}
	u.checkSendingMessage(ctx)
	u.user.SetLoginUserID(userID)
	u.user.SetDataBase(u.db)
	u.file.SetLoginUserID(userID)
	u.file.SetDataBase(u.db)
	u.relation.SetDataBase(u.db)
	u.relation.SetLoginUserID(userID)
	u.group.SetDataBase(u.db)
	u.group.SetLoginUserID(userID)
	u.third.SetPlatform(u.info.PlatformID)
	u.third.SetLoginUserID(userID)
	u.third.SetAppFramework(u.info.SystemType)
	u.third.SetLogFilePath(u.info.LogFilePath)
	u.msgSyncer.SetLoginUserID(userID)
	u.msgSyncer.SetDataBase(u.db)
	u.conversation.SetLoginUserID(userID)
	u.conversation.SetDataBase(u.db)
	u.conversation.SetPlatform(u.info.PlatformID)
	u.conversation.SetDataDir(u.info.DataDir)
	err = u.msgSyncer.LoadSeq(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserContext) setListener(ctx context.Context) {
	setListener(ctx, &u.connListener, u.ConnListener, u.longConnMgr.SetListener, nil)
	setListener(ctx, &u.userListener, u.UserListener, u.user.SetListener, newEmptyUserListener)
	setListener(ctx, &u.friendshipListener, u.FriendshipListener, u.relation.SetListener, newEmptyFriendshipListener)
	setListener(ctx, &u.groupListener, u.GroupListener, u.group.SetGroupListener, newEmptyGroupListener)
	setListener(ctx, &u.conversationListener, u.ConversationListener, u.conversation.SetConversationListener, newEmptyConversationListener)
	setListener(ctx, &u.advancedMsgListener, u.AdvancedMsgListener, u.conversation.SetMsgListener, newEmptyAdvancedMsgListener)
	setListener(ctx, &u.businessListener, u.BusinessListener, u.conversation.SetBusinessListener, newEmptyCustomBusinessListener)
}

func setListener[T any](ctx context.Context, listener *T, getter func() T, setFunc func(listener func() T), newFunc func(context.Context) T) {
	if *(*unsafe.Pointer)(unsafe.Pointer(listener)) == nil && newFunc != nil {
		*listener = newFunc(ctx)
	}
	setFunc(getter)
}

func (u *UserContext) run(ctx context.Context) {
	u.longConnMgr.Run(ctx, u.fgCtx)
	go u.msgSyncer.DoListener(ctx)
	go common.DoListener(u.ctx, u.conversation)
	go u.logoutListener(ctx)
}

func (u *UserContext) setFGCtx() {
	u.fgCtx, u.fgCancel = context.WithCancelCause(context.Background())
}

func (u *UserContext) InitSDK(config *sdk_struct.IMConfig, listener open_im_sdk_callback.OnConnListener) bool {
	if listener == nil {
		return false
	}
	u.info.IMConfig = config
	u.connListener = listener
	return true
}

func (u *UserContext) Context() context.Context {
	return u.ctx
}

func (u *UserContext) userOnlineStatusChange(users map[string][]int32) {
	u.User().UserOnlineStatusChange(users)
}

func (u *UserContext) UnInitSDK() {
	if u.getLoginStatus(context.Background()) == Logged {
		fmt.Println("sdk not logout, please logout first")
		return
	}
	u.Info().IMConfig = nil
	u.setLoginStatus(0)
}

// token error recycle recourse, kicked not recycle
func (u *UserContext) logout(ctx context.Context, isTokenValid bool) error {
	if ccontext.Info(ctx).OperationID() == LogoutTips {
		isTokenValid = true
	}
	if !isTokenValid {
		ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()
		err := u.longConnMgr.SendReqWaitResp(ctx, &push.DelUserPushTokenReq{UserID: u.info.UserID, PlatformID: u.info.PlatformID}, constant.LogoutMsg, &push.DelUserPushTokenResp{})
		if err != nil {
			log.ZWarn(ctx, "TriggerCmdLogout server recycle resources failed...", err)
		} else {
			log.ZDebug(ctx, "TriggerCmdLogout server recycle resources success...")
		}
	}
	u.Exit()
	err := u.db.Close(u.ctx)
	if err != nil {
		log.ZWarn(ctx, "TriggerCmdLogout db recycle resources failed...", err)
	}
	// user object must be rest  when user logout
	u.initResources()
	log.ZDebug(ctx, "TriggerCmdLogout client success...",
		"isTokenValid", isTokenValid)
	return nil
}

func (u *UserContext) setAppBackgroundStatus(ctx context.Context, isBackground bool) error {

	u.longConnMgr.SetBackground(isBackground)

	if !isBackground {
		if u.info.StopGoroutineOnBackground {
			u.setFGCtx()
			u.longConnMgr.ResumeForegroundTasks(u.ctx, u.fgCtx)
		}
	} else {
		if u.info.StopGoroutineOnBackground {
			u.fgCancel(errs.Wrap(fmt.Errorf("app in background")))
			u.longConnMgr.Close(ctx)
		}
	}
	var resp sdkws.SetAppBackgroundStatusResp
	err := u.longConnMgr.SendReqWaitResp(ctx, &sdkws.SetAppBackgroundStatusReq{UserID: u.loginUserID, IsBackground: isBackground}, constant.SetBackgroundStatus, &resp)
	if err != nil {
		return err
	} else {
		if !isBackground {
			_ = common.DispatchWakeUp(ctx, u.msgSyncerCh)
		}
		return nil
	}
}

func (u *UserContext) LongConnMgr() *interaction.LongConnMgr {
	return u.longConnMgr
}
