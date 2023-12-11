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
	"errors"
	"fmt"
	"github.com/OpenIMSDK/protocol/push"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/openim-sdk-core/v3/internal/business"
	conv "github.com/openimsdk/openim-sdk-core/v3/internal/conversation_msg"
	"github.com/openimsdk/openim-sdk-core/v3/internal/file"
	"github.com/openimsdk/openim-sdk-core/v3/internal/friend"
	"github.com/openimsdk/openim-sdk-core/v3/internal/full"
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
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"strings"
	"sync"
	"time"
)

const (
	LogoutStatus = iota + 1
	Logging
	Logged
)

var (
	// Client-independent user class
	UserForSDK *LoginMgr
)

// CheckResourceLoad checks the SDK is resource load status.
func CheckResourceLoad(uSDK *LoginMgr, funcName string) error {
	if uSDK == nil {
		return utils.Wrap(errors.New("CheckResourceLoad failed uSDK == nil "), "")
	}
	if funcName == "" {
		return nil
	}
	parts := strings.Split(funcName, ".")
	if parts[len(parts)-1] == "Login-fm" {
		return nil
	}
	if uSDK.Friend() == nil || uSDK.User() == nil || uSDK.Group() == nil || uSDK.Conversation() == nil ||
		uSDK.Full() == nil {
		return utils.Wrap(errors.New("CheckResourceLoad failed, resource nil "), "")
	}
	return nil
}

type LoginMgr struct {
	friend       *friend.Friend
	group        *group.Group
	conversation *conv.Conversation
	user         *user.User
	file         *file.File
	business     *business.Business

	full         *full.Full
	db           db_interface.DataBase
	longConnMgr  *interaction.LongConnMgr
	msgSyncer    *interaction.MsgSyncer
	third        *third.Third
	token        string
	loginUserID  string
	connListener open_im_sdk_callback.OnConnListener

	loginTime int64

	justOnceFlag bool

	w           sync.Mutex
	loginStatus int

	groupListener               open_im_sdk_callback.OnGroupListener
	friendListener              open_im_sdk_callback.OnFriendshipListener
	conversationListener        open_im_sdk_callback.OnConversationListener
	advancedMsgListener         open_im_sdk_callback.OnAdvancedMsgListener
	batchMsgListener            open_im_sdk_callback.OnBatchMsgListener
	userListener                open_im_sdk_callback.OnUserListener
	signalingListener           open_im_sdk_callback.OnSignalingListener
	signalingListenerFroService open_im_sdk_callback.OnSignalingListener
	businessListener            open_im_sdk_callback.OnCustomBusinessListener

	conversationCh     chan common.Cmd2Value
	cmdWsCh            chan common.Cmd2Value
	heartbeatCmdCh     chan common.Cmd2Value
	pushMsgAndMaxSeqCh chan common.Cmd2Value
	loginMgrCh         chan common.Cmd2Value

	ctx       context.Context
	cancel    context.CancelFunc
	info      *ccontext.GlobalConfig
	id2MinSeq map[string]int64
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

func (u *LoginMgr) Third() *third.Third {
	return u.third
}

func (u *LoginMgr) ImConfig() sdk_struct.IMConfig {
	return sdk_struct.IMConfig{
		PlatformID:           u.info.PlatformID,
		ApiAddr:              u.info.ApiAddr,
		WsAddr:               u.info.WsAddr,
		DataDir:              u.info.DataDir,
		LogLevel:             u.info.LogLevel,
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
		u.friend.SetListener(friendListener)
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

func (u *LoginMgr) SetListenerForService(listener open_im_sdk_callback.OnListenerForService) {
	if u.friend == nil || u.group == nil || u.conversation == nil {
		return
	}
	u.friend.SetListenerForService(listener)
	u.group.SetListenerForService(listener)
	u.conversation.SetListenerForService(listener)
}

func (u *LoginMgr) SetBusinessListener(listener open_im_sdk_callback.OnCustomBusinessListener) {
	if u.business != nil {
		u.business.SetListener(listener)
	} else {
		u.businessListener = listener
	}
}
func (u *LoginMgr) GetLoginUserID() string {
	return u.loginUserID
}
func (u *LoginMgr) logoutListener(ctx context.Context) {
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

func NewLoginMgr() *LoginMgr {
	return &LoginMgr{
		info: &ccontext.GlobalConfig{}, // 分配内存空间
	}
}
func (u *LoginMgr) getLoginStatus(_ context.Context) int {
	u.w.Lock()
	defer u.w.Unlock()
	return u.loginStatus
}
func (u *LoginMgr) setLoginStatus(status int) {
	u.w.Lock()
	defer u.w.Unlock()
	u.loginStatus = status
}
func (u *LoginMgr) checkSendingMessage(ctx context.Context) {
	sendingMessages, err := u.db.GetAllSendingMessages(ctx)
	if err != nil {
		log.ZError(ctx, "GetAllSendingMessages failed", err)
	}
	for _, message := range sendingMessages {
		tableMessage, err := u.db.GetMessage(ctx, message.ConversationID, message.ClientMsgID)
		if err != nil {
			log.ZError(ctx, "GetMessage failed", err, "message", message)
			continue
		}
		if tableMessage.Status == constant.MsgStatusSending {
			err := u.db.UpdateMessage(ctx, message.ConversationID, &model_struct.LocalChatLog{ClientMsgID: message.ClientMsgID, Status: constant.MsgStatusSendFailed})
			if err != nil {
				log.ZError(ctx, "UpdateMessage failed", err, "tableMessage", tableMessage)
			} else {
				err := u.db.DeleteSendingMessage(ctx, message.ConversationID, message.ClientMsgID)
				if err != nil {
					log.ZError(ctx, "DeleteSendingMessage failed", err, "tableMessage", tableMessage)
				}
			}

		}
	}
}

func (u *LoginMgr) login(ctx context.Context, userID, token string) error {
	if u.getLoginStatus(ctx) == Logged {
		return sdkerrs.ErrLoginRepeat
	}
	u.setLoginStatus(Logging)
	u.info.UserID = userID
	u.info.Token = token
	log.ZInfo(ctx, "login start... ", "userID", userID, "token", token)
	t1 := time.Now()
	u.token = token
	u.loginUserID = userID
	var err error
	u.db, err = db.NewDataBase(ctx, userID, u.info.DataDir, int(u.info.LogLevel))
	if err != nil {
		return sdkerrs.ErrSdkInternal.Wrap("init database " + err.Error())
	}
	u.checkSendingMessage(ctx)
	log.ZDebug(ctx, "NewDataBase ok", "userID", userID, "dataDir", u.info.DataDir, "login cost time", time.Since(t1))
	u.loginTime = time.Now().UnixNano() / 1e6
	u.user = user.NewUser(u.db, u.loginUserID, u.conversationCh)
	u.user.SetListener(u.userListener)
	u.file = file.NewFile(u.db, u.loginUserID)
	u.friend = friend.NewFriend(u.loginUserID, u.db, u.user, u.conversationCh)
	u.friend.SetListener(u.friendListener)
	u.friend.SetLoginTime(u.loginTime)
	u.group = group.NewGroup(u.loginUserID, u.db, u.conversationCh)
	u.group.SetGroupListener(u.groupListener)
	u.full = full.NewFull(u.user, u.friend, u.group, u.conversationCh, u.db)
	u.business = business.NewBusiness(u.db)
	if u.businessListener != nil {
		u.business.SetListener(u.businessListener)
	}
	u.third = third.NewThird(u.info.PlatformID, u.loginUserID, constant.SdkVersion, u.info.LogFilePath, u.file)
	log.ZDebug(ctx, "forcedSynchronization success...", "login cost time: ", time.Since(t1))

	u.longConnMgr.Run(ctx)
	u.msgSyncer, _ = interaction.NewMsgSyncer(ctx, u.conversationCh, u.pushMsgAndMaxSeqCh, u.loginUserID, u.longConnMgr, u.db, 0)
	u.conversation = conv.NewConversation(ctx, u.longConnMgr, u.db, u.conversationCh,
		u.friend, u.group, u.user, u.conversationListener, u.advancedMsgListener, u.business, u.full, u.file)
	u.conversation.SetLoginTime()
	if u.batchMsgListener != nil {
		u.conversation.SetBatchMsgListener(u.batchMsgListener)
		log.ZDebug(ctx, "SetBatchMsgListener", "batchMsgListener", u.batchMsgListener)
	}
	go common.DoListener(u.conversation, u.ctx)
	go func() {
		memberGroupIDs, err := u.db.GetGroupMemberAllGroupIDs(ctx)
		if err != nil {
			log.ZError(ctx, "GetGroupMemberAllGroupIDs failed", err)
			return
		}
		if len(memberGroupIDs) > 0 {
			groups, err := u.db.GetJoinedGroupListDB(ctx)
			if err != nil {
				log.ZError(ctx, "GetJoinedGroupListDB failed", err)
				return
			}
			memberGroupIDMap := make(map[string]struct{})
			for _, groupID := range memberGroupIDs {
				memberGroupIDMap[groupID] = struct{}{}
			}
			for _, info := range groups {
				delete(memberGroupIDMap, info.GroupID)
			}
			for groupID := range memberGroupIDMap {
				if err := u.db.DeleteGroupAllMembers(ctx, groupID); err != nil {
					log.ZError(ctx, "DeleteGroupAllMembers failed", err, "groupID", groupID)
				}
			}
		}
	}()

	go u.logoutListener(ctx)
	u.setLoginStatus(Logged)
	log.ZInfo(ctx, "login success...", "login cost time: ", time.Since(t1))
	return nil
}

func (u *LoginMgr) InitSDK(config sdk_struct.IMConfig, listener open_im_sdk_callback.OnConnListener) bool {
	if listener == nil {
		return false
	}
	u.info = &ccontext.GlobalConfig{}
	u.info.IMConfig = config
	u.connListener = listener
	u.initResources()
	return true
}

func (u *LoginMgr) Context() context.Context {
	return u.ctx
}

func (u *LoginMgr) initResources() {
	ctx := ccontext.WithInfo(context.Background(), u.info)
	u.ctx, u.cancel = context.WithCancel(ctx)
	u.conversationCh = make(chan common.Cmd2Value, 1000)
	u.heartbeatCmdCh = make(chan common.Cmd2Value, 10)
	u.pushMsgAndMaxSeqCh = make(chan common.Cmd2Value, 1000)
	u.loginMgrCh = make(chan common.Cmd2Value, 1)
	u.longConnMgr = interaction.NewLongConnMgr(u.ctx, u.connListener, u.heartbeatCmdCh, u.pushMsgAndMaxSeqCh, u.loginMgrCh)
	u.setLoginStatus(LogoutStatus)
}

func (u *LoginMgr) UnInitSDK() {
	if u.getLoginStatus(context.Background()) == Logged {
		fmt.Println("sdk not logout, please logout first")
		return
	}
	u.info = nil
	u.setLoginStatus(0)
}

// token error recycle recourse, kicked not recycle
func (u *LoginMgr) logout(ctx context.Context, isTokenValid bool) error {
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
	_ = u.db.Close(u.ctx)
	// user object must be rest  when user logout
	u.initResources()
	log.ZDebug(ctx, "TriggerCmdLogout client success...", "isTokenValid", isTokenValid)
	return nil
}

func (u *LoginMgr) setAppBackgroundStatus(ctx context.Context, isBackground bool) error {
	if u.longConnMgr.GetConnectionStatus() == 0 {
		u.longConnMgr.SetBackground(isBackground)
		return nil
	}
	var resp sdkws.SetAppBackgroundStatusResp
	err := u.longConnMgr.SendReqWaitResp(ctx, &sdkws.SetAppBackgroundStatusReq{UserID: u.loginUserID, IsBackground: isBackground}, constant.SetBackgroundStatus, &resp)
	if err != nil {
		return err
	} else {
		u.longConnMgr.SetBackground(isBackground)
		if isBackground == false {
			_ = common.TriggerCmdWakeUp(u.heartbeatCmdCh)
		}
		return nil
	}
}
