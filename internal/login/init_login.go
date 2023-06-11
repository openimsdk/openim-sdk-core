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
	"open_im_sdk/internal/signaling"
	"open_im_sdk/internal/super_group"
	"open_im_sdk/internal/third"
	"open_im_sdk/internal/user"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/sdkerrs"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/push"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
)

type LoginMgr struct {
	friend       *friend.Friend
	group        *group.Group
	superGroup   *super_group.SuperGroup
	conversation *conv.Conversation
	user         *user.User
	file         *file.File
	signaling    *signaling.LiveSignaling
	business     *business.Business

	full         *full.Full
	db           db_interface.DataBase
	longConnMgr  *interaction.LongConnMgr
	msgSyncer    *interaction.MsgSyncer
	push         *third.Push
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
	businessListener            open_im_sdk_callback.OnCustomBusinessListener

	conversationCh     chan common.Cmd2Value
	cmdWsCh            chan common.Cmd2Value
	heartbeatCmdCh     chan common.Cmd2Value
	pushMsgAndMaxSeqCh chan common.Cmd2Value

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

func (u *LoginMgr) Push() *third.Push {
	return u.push
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

func (u *LoginMgr) Signaling() *signaling.LiveSignaling {
	return u.signaling
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

func (u *LoginMgr) SetBusinessListener(listener open_im_sdk_callback.OnCustomBusinessListener) {
	if u.business != nil {
		u.business.SetListener(listener)
	} else {
		u.businessListener = listener
	}
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
		return sdkerrs.ErrSdkInternal.Wrap("init database " + err.Error())
	}
	log.ZDebug(ctx, "NewDataBase ok", "userID", userID, "dataDir", u.info.DataDir, "login cost time", time.Since(t1))
	u.conversationCh = make(chan common.Cmd2Value, 1000)
	u.heartbeatCmdCh = make(chan common.Cmd2Value, 10)
	u.pushMsgAndMaxSeqCh = make(chan common.Cmd2Value, 1000)
	u.loginTime = time.Now().UnixNano() / 1e6
	u.user = user.NewUser(u.db, u.loginUserID, u.conversationCh)
	u.user.SetListener(u.userListener)
	u.file = file.NewFile(u.db, u.loginUserID)
	u.friend = friend.NewFriend(u.loginUserID, u.db, u.user, u.conversationCh)
	u.friend.SetListener(u.friendListener)
	u.friend.SetLoginTime(u.loginTime)
	u.group = group.NewGroup(u.loginUserID, u.db, u.conversationCh)
	u.group.SetGroupListener(u.groupListener)
	u.superGroup = super_group.NewSuperGroup(u.loginUserID, u.db)
	u.cache = cache.NewCache(u.user, u.friend)
	u.full = full.NewFull(u.user, u.friend, u.group, u.conversationCh, u.cache, u.db, u.superGroup)

	u.business = business.NewBusiness(u.db)
	if u.businessListener != nil {
		u.business.SetListener(u.businessListener)
	}
	u.push = third.NewPush(u.info.PlatformID, u.loginUserID)
	log.ZDebug(ctx, "forcedSynchronization success...", "login cost time: ", time.Since(t1))
	u.longConnMgr = interaction.NewLongConnMgr(ctx, u.connListener, u.heartbeatCmdCh, u.pushMsgAndMaxSeqCh, u.conversationCh)
	u.msgSyncer, _ = interaction.NewMsgSyncer(ctx, u.conversationCh, u.pushMsgAndMaxSeqCh, u.loginUserID, u.longConnMgr, u.db, 0)
	u.conversation = conv.NewConversation(ctx, u.longConnMgr, u.db, u.conversationCh,
		u.friend, u.group, u.user, u.conversationListener, u.advancedMsgListener, u.signaling, u.business, u.cache, u.full, u.file)
	u.conversation.SetLoginTime()
	u.signaling = signaling.NewLiveSignaling(u.longConnMgr, u.loginUserID, u.info.PlatformID, u.db)
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
	ctx := ccontext.WithInfo(context.Background(), u.info)
	u.ctx, u.cancel = context.WithCancel(ctx)
	return true
}

func (u *LoginMgr) Context() context.Context {
	return u.ctx
}

func (u *LoginMgr) logout(ctx context.Context) error {
	err := u.longConnMgr.SendReqWaitResp(ctx, &push.DelUserPushTokenReq{UserID: u.info.UserID, PlatformID: u.info.PlatformID}, constant.LogoutMsg, &push.DelUserPushTokenResp{})
	if err != nil {
		return err
	}
	u.Exit()
	_ = u.db.Close(u.ctx)
	//reset global context
	resetContext := ccontext.WithInfo(context.Background(), u.info)
	u.ctx, u.cancel = context.WithCancel(resetContext)
	log.ZDebug(ctx, "TriggerCmdLogout success...")
	return nil
}

func (u *LoginMgr) setAppBackgroundStatus(ctx context.Context, isBackground bool) error {
	err := u.longConnMgr.SendReqWaitResp(ctx, &server_api_params.SetAppBackgroundStatusReq{UserID: u.loginUserID, IsBackground: isBackground}, constant.SetBackgroundStatus, nil)
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

func (u *LoginMgr) GetLoginUserID() string {
	return u.loginUserID
}

func (u *LoginMgr) GetLoginStatus() int {
	return u.longConnMgr.GetConnectionStatus()
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
