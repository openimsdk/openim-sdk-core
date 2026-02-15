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

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/tools/log"
)

var ErrNotImplemented = errors.New("not set listener")

type emptyGroupListener struct {
	ctx context.Context
}

func newEmptyGroupListener(ctx context.Context) open_im_sdk_callback.OnGroupListener {
	return &emptyGroupListener{ctx}
}

func (e *emptyGroupListener) OnJoinedGroupAdded(groupInfo string) {
	log.ZWarn(e.ctx, "OnJoinedGroupAdded", nil, "groupInfo", groupInfo)
}

func (e *emptyGroupListener) OnJoinedGroupDeleted(groupInfo string) {
	log.ZWarn(e.ctx, "OnJoinedGroupDeleted", nil, "groupInfo", groupInfo)
}

func (e *emptyGroupListener) OnGroupMemberAdded(groupMemberInfo string) {
	log.ZWarn(e.ctx, "OnGroupMemberAdded", nil, "groupMemberInfo", groupMemberInfo)
}

func (e *emptyGroupListener) OnGroupMemberDeleted(groupMemberInfo string) {
	log.ZWarn(e.ctx, "OnGroupMemberDeleted", nil, "groupMemberInfo", groupMemberInfo)
}

func (e *emptyGroupListener) OnGroupApplicationAdded(groupApplication string) {
	log.ZWarn(e.ctx, "OnGroupApplicationAdded", nil, "groupApplication", groupApplication)
}

func (e *emptyGroupListener) OnGroupApplicationDeleted(groupApplication string) {
	log.ZWarn(e.ctx, "OnGroupApplicationDeleted", nil, "groupApplication", groupApplication)
}

func (e *emptyGroupListener) OnGroupInfoChanged(groupInfo string) {
	log.ZWarn(e.ctx, "OnGroupInfoChanged", nil, "groupInfo", groupInfo)
}

func (e *emptyGroupListener) OnGroupDismissed(groupInfo string) {
	log.ZWarn(e.ctx, "OnGroupDismissed", nil, "groupInfo", groupInfo)
}

func (e *emptyGroupListener) OnGroupMemberInfoChanged(groupMemberInfo string) {
	log.ZWarn(e.ctx, "OnGroupMemberInfoChanged", nil, "groupMemberInfo", groupMemberInfo)
}

func (e *emptyGroupListener) OnGroupApplicationAccepted(groupApplication string) {
	log.ZWarn(e.ctx, "OnGroupApplicationAccepted", nil, "groupApplication", groupApplication)
}

func (e *emptyGroupListener) OnGroupApplicationRejected(groupApplication string) {
	log.ZWarn(e.ctx, "OnGroupApplicationRejected", nil, "groupApplication", groupApplication)
}

type emptyFriendshipListener struct {
	ctx context.Context
}

func newEmptyFriendshipListener(ctx context.Context) open_im_sdk_callback.OnFriendshipListener {
	return &emptyFriendshipListener{ctx}
}

func (e *emptyFriendshipListener) OnFriendApplicationAdded(friendApplication string) {
	log.ZWarn(e.ctx, "OnFriendApplicationAdded", nil, "friendApplication", friendApplication)
}

func (e *emptyFriendshipListener) OnFriendApplicationDeleted(friendApplication string) {
	log.ZWarn(e.ctx, "OnFriendApplicationDeleted", nil, "friendApplication", friendApplication)
}

func (e *emptyFriendshipListener) OnFriendApplicationAccepted(friendApplication string) {
	log.ZWarn(e.ctx, "OnFriendApplicationAccepted", nil, "friendApplication", friendApplication)
}

func (e *emptyFriendshipListener) OnFriendApplicationRejected(friendApplication string) {
	log.ZWarn(e.ctx, "OnFriendApplicationRejected", nil, "friendApplication", friendApplication)
}

func (e *emptyFriendshipListener) OnFriendAdded(friendInfo string) {
	log.ZWarn(e.ctx, "OnFriendAdded", nil, "friendInfo", friendInfo)
}

func (e *emptyFriendshipListener) OnFriendDeleted(friendInfo string) {
	log.ZWarn(e.ctx, "OnFriendDeleted", nil, "friendInfo", friendInfo)
}

func (e *emptyFriendshipListener) OnFriendInfoChanged(friendInfo string) {
	log.ZWarn(e.ctx, "OnFriendInfoChanged", nil, "friendInfo", friendInfo)
}

func (e *emptyFriendshipListener) OnBlackAdded(blackInfo string) {
	log.ZWarn(e.ctx, "OnBlackAdded", nil, "blackInfo", blackInfo)
}

func (e *emptyFriendshipListener) OnBlackDeleted(blackInfo string) {
	log.ZWarn(e.ctx, "OnBlackDeleted", nil, "blackInfo", blackInfo)
}

type emptyConversationListener struct {
	ctx context.Context
}

func newEmptyConversationListener(ctx context.Context) open_im_sdk_callback.OnConversationListener {
	return &emptyConversationListener{ctx: ctx}
}

func (e *emptyConversationListener) OnSyncServerStart(reinstalled bool) {
	log.ZWarn(e.ctx, "OnSyncServerStart", nil)
}

func (e *emptyConversationListener) OnSyncServerProgress(progress int) {
	log.ZWarn(e.ctx, "OnSyncServerProgress", nil, "progress", progress)
}

func (e *emptyConversationListener) OnSyncServerFinish(reinstalled bool) {
	log.ZWarn(e.ctx, "OnSyncServerFinish", nil)
}

func (e *emptyConversationListener) OnSyncServerFailed(reinstalled bool) {
	log.ZWarn(e.ctx, "OnSyncServerFailed", nil)
}

func (e *emptyConversationListener) OnNewConversation(conversationList string) {
	log.ZWarn(e.ctx, "OnNewConversation", nil, "conversationList", conversationList)
}

func (e *emptyConversationListener) OnConversationChanged(conversationList string) {
	log.ZWarn(e.ctx, "OnConversationChanged", nil, "conversationList", conversationList)
}

func (e *emptyConversationListener) OnTotalUnreadMessageCountChanged(totalUnreadCount int32) {
	log.ZWarn(e.ctx, "OnTotalUnreadMessageCountChanged", nil, "totalUnreadCount", totalUnreadCount)
}

func (e *emptyConversationListener) OnConversationUserInputStatusChanged(change string) {
	log.ZWarn(e.ctx, "OnConversationUserInputStatusChanged", nil, "change", change)
}

type emptyAdvancedMsgListener struct {
	ctx context.Context
}

func newEmptyAdvancedMsgListener(ctx context.Context) open_im_sdk_callback.OnAdvancedMsgListener {
	return &emptyAdvancedMsgListener{ctx}
}

func (e *emptyAdvancedMsgListener) OnRecvOnlineOnlyMessage(message string) {

}

func (e *emptyAdvancedMsgListener) OnRecvNewMessage(message string) {
	log.ZWarn(e.ctx, "OnRecvNewMessage", nil, "message", message)
}

func (e *emptyAdvancedMsgListener) OnRecvC2CReadReceipt(msgReceiptList string) {
	log.ZWarn(e.ctx, "OnRecvC2CReadReceipt", nil, "msgReceiptList", msgReceiptList)
}

func (e *emptyAdvancedMsgListener) OnRecvGroupReadReceipt(groupMsgReceiptList string) {
	log.ZWarn(e.ctx, "OnRecvGroupReadReceipt", nil, "groupMsgReceiptList", groupMsgReceiptList)
}

func (e *emptyAdvancedMsgListener) OnNewRecvMessageRevoked(messageRevoked string) {
	log.ZWarn(e.ctx, "OnNewRecvMessageRevoked", nil, "messageRevoked", messageRevoked)
}

func (e *emptyAdvancedMsgListener) OnMsgEdited(msg string) {
	log.ZWarn(e.ctx, "OnMsgEdited", nil, "msg", msg)
}

func (e *emptyAdvancedMsgListener) OnRecvMessageExtensionsChanged(msgID string, reactionExtensionList string) {
	log.ZWarn(e.ctx, "OnRecvMessageExtensionsChanged", nil, "msgID", msgID,
		"reactionExtensionList", reactionExtensionList)
}

func (e *emptyAdvancedMsgListener) OnRecvMessageExtensionsDeleted(msgID string, reactionExtensionKeyList string) {
	log.ZWarn(e.ctx, "OnRecvMessageExtensionsDeleted", nil, "msgID", msgID,
		"reactionExtensionKeyList", reactionExtensionKeyList)
}

func (e *emptyAdvancedMsgListener) OnRecvMessageExtensionsAdded(msgID string, reactionExtensionList string) {
	log.ZWarn(e.ctx, "OnRecvMessageExtensionsAdded", nil, "msgID", msgID,
		"reactionExtensionList", reactionExtensionList)
}

func (e *emptyAdvancedMsgListener) OnRecvOfflineNewMessage(message string) {
	log.ZWarn(e.ctx, "OnRecvOfflineNewMessage", nil, "message", message)
}

func (e *emptyAdvancedMsgListener) OnMsgDeleted(message string) {
	log.ZWarn(e.ctx, "OnMsgDeleted", nil, "message", message)
}

type emptyUserListener struct {
	ctx context.Context
}

func newEmptyUserListener(ctx context.Context) open_im_sdk_callback.OnUserListener {
	return &emptyUserListener{ctx: ctx}
}

func (e *emptyUserListener) OnSelfInfoUpdated(userInfo string) {
	log.ZWarn(e.ctx, "OnSelfInfoUpdated", nil, "userInfo", userInfo)
}

func (e *emptyUserListener) OnUserStatusChanged(statusMap string) {
	log.ZWarn(e.ctx, "OnUserStatusChanged", nil, "statusMap", statusMap)
}

type emptyCustomBusinessListener struct {
	ctx context.Context
}

func newEmptyCustomBusinessListener(ctx context.Context) open_im_sdk_callback.OnCustomBusinessListener {
	return &emptyCustomBusinessListener{ctx: ctx}
}

func (e *emptyCustomBusinessListener) OnRecvCustomBusinessMessage(businessMessage string) {
	log.ZWarn(e.ctx, "OnRecvCustomBusinessMessage", nil, "businessMessage", businessMessage)
}
