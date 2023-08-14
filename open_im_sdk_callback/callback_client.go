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

package open_im_sdk_callback

type Base interface {
	OnError(errCode int32, errMsg string)
	OnSuccess(data string)
}
type SendMsgCallBack interface {
	Base
	OnProgress(progress int)
}

type OnConnListener interface {
	OnConnecting()
	OnConnectSuccess()
	OnConnectFailed(errCode int32, errMsg string)
	OnKickedOffline()
	OnUserTokenExpired()
}

type OnGroupListener interface {
	OnJoinedGroupAdded(groupInfo string)
	OnJoinedGroupDeleted(groupInfo string)
	OnGroupMemberAdded(groupMemberInfo string)
	OnGroupMemberDeleted(groupMemberInfo string)
	OnGroupApplicationAdded(groupApplication string)
	OnGroupApplicationDeleted(groupApplication string)
	OnGroupInfoChanged(groupInfo string)
	OnGroupDismissed(groupInfo string)
	OnGroupMemberInfoChanged(groupMemberInfo string)
	OnGroupApplicationAccepted(groupApplication string)
	OnGroupApplicationRejected(groupApplication string)
}
type OnFriendshipListener interface {
	OnFriendApplicationAdded(friendApplication string)
	OnFriendApplicationDeleted(friendApplication string)
	OnFriendApplicationAccepted(friendApplication string)
	OnFriendApplicationRejected(friendApplication string)
	OnFriendAdded(friendInfo string)
	OnFriendDeleted(friendInfo string)
	OnFriendInfoChanged(friendInfo string)
	OnBlackAdded(blackInfo string)
	OnBlackDeleted(blackInfo string)
}
type OnConversationListener interface {
	OnSyncServerStart()
	OnSyncServerFinish()
	//OnSyncServerProgress(progress int)
	OnSyncServerFailed()
	OnNewConversation(conversationList string)
	OnConversationChanged(conversationList string)
	OnTotalUnreadMessageCountChanged(totalUnreadCount int32)
}
type OnAdvancedMsgListener interface {
	OnRecvNewMessage(message string)
	OnRecvC2CReadReceipt(msgReceiptList string)
	OnRecvGroupReadReceipt(groupMsgReceiptList string)

	OnNewRecvMessageRevoked(messageRevoked string)
	OnRecvMessageExtensionsChanged(msgID string, reactionExtensionList string)
	OnRecvMessageExtensionsDeleted(msgID string, reactionExtensionKeyList string)
	OnRecvMessageExtensionsAdded(msgID string, reactionExtensionList string)
	OnRecvOfflineNewMessage(message string)
	OnMsgDeleted(message string)
}

type OnBatchMsgListener interface {
	OnRecvNewMessages(messageList string)
	OnRecvOfflineNewMessages(messageList string)
}

type OnUserListener interface {
	OnSelfInfoUpdated(userInfo string)
	OnUserStatusChanged(statusMap string)
}

type OnCustomBusinessListener interface {
	OnRecvCustomBusinessMessage(businessMessage string)
}
type OnMessageKvInfoListener interface {
	OnMessageKvInfoChanged(messageChangedList string)
}

type OnListenerForService interface {
	//有人申请进群
	OnGroupApplicationAdded(groupApplication string)
	//进群申请被同意
	OnGroupApplicationAccepted(groupApplication string)
	//有人申请添加你为好友
	OnFriendApplicationAdded(friendApplication string)
	//好友申请被同意
	OnFriendApplicationAccepted(groupApplication string)
	//收到新消息
	OnRecvNewMessage(message string)
}

type OnSignalingListener interface {
	OnReceiveNewInvitation(receiveNewInvitationCallback string)

	OnInviteeAccepted(inviteeAcceptedCallback string)

	OnInviteeAcceptedByOtherDevice(inviteeAcceptedCallback string)

	OnInviteeRejected(inviteeRejectedCallback string)

	OnInviteeRejectedByOtherDevice(inviteeRejectedCallback string)
	//
	OnInvitationCancelled(invitationCancelledCallback string)
	//
	OnInvitationTimeout(invitationTimeoutCallback string)
	//
	OnHangUp(hangUpCallback string)

	OnRoomParticipantConnected(onRoomParticipantConnectedCallback string)

	OnRoomParticipantDisconnected(onRoomParticipantDisconnectedCallback string)
}

type UploadFileCallback interface {
	Open(size int64)                                                    // 文件打开的大小
	PartSize(partSize int64, num int)                                   // 分片大小,数量
	HashPartProgress(index int, size int64, partHash string)            // 每块分片的hash值
	HashPartComplete(partsHash string, fileHash string)                 // 分块完成，服务端标记hash和文件最终hash
	UploadID(uploadID string)                                           // 上传ID
	UploadPartComplete(index int, partSize int64, partHash string)      // 上传分片进度
	UploadComplete(fileSize int64, streamSize int64, storageSize int64) // 整体进度
	Complete(size int64, url string, typ int)                           // 上传完成
}
