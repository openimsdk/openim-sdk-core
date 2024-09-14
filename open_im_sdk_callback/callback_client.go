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
	OnUserTokenInvalid(errMsg string)
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
	OnSyncServerStart(reinstalled bool)
	OnSyncServerFinish(reinstalled bool)
	OnSyncServerProgress(progress int)
	OnSyncServerFailed(reinstalled bool)
	OnNewConversation(conversationList string)
	OnConversationChanged(conversationList string)
	OnTotalUnreadMessageCountChanged(totalUnreadCount int32)
	OnConversationUserInputStatusChanged(change string)
}

type OnAdvancedMsgListener interface {
	OnRecvNewMessage(message string)
	OnRecvC2CReadReceipt(msgReceiptList string)
	//OnRecvGroupReadReceipt(groupMsgReceiptList string)

	OnNewRecvMessageRevoked(messageRevoked string)
	//OnRecvMessageExtensionsChanged(msgID string, reactionExtensionList string)
	//OnRecvMessageExtensionsDeleted(msgID string, reactionExtensionKeyList string)
	//OnRecvMessageExtensionsAdded(msgID string, reactionExtensionList string)
	OnRecvOfflineNewMessage(message string)
	OnMsgDeleted(message string)
	OnRecvOnlineOnlyMessage(message string)
}

type OnBatchMsgListener interface {
	OnRecvNewMessages(messageList string)
	OnRecvOfflineNewMessages(messageList string)
}

type OnUserListener interface {
	OnSelfInfoUpdated(userInfo string)
	OnUserStatusChanged(userOnlineStatus string)
	OnUserCommandAdd(userCommand string)
	OnUserCommandDelete(userCommand string)
	OnUserCommandUpdate(userCommand string)
}

type OnCustomBusinessListener interface {
	OnRecvCustomBusinessMessage(businessMessage string)
}
type OnMessageKvInfoListener interface {
	OnMessageKvInfoChanged(messageChangedList string)
}

type OnListenerForService interface {
	// OnGroupApplicationAdded someone has requested to join the group.
	OnGroupApplicationAdded(groupApplication string)
	// OnGroupApplicationAccepted the group join request has been approved.
	OnGroupApplicationAccepted(groupApplication string)
	// OnFriendApplicationAdded someone has requested to add you as a friend.
	OnFriendApplicationAdded(friendApplication string)
	// OnFriendApplicationAccepted the friend request has been accepted.
	OnFriendApplicationAccepted(groupApplication string)
	// OnRecvNewMessage new message received
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
	Open(size int64)                                                    // Open a file with the specified size
	PartSize(partSize int64, num int)                                   // Set the size and number of parts for chunking
	HashPartProgress(index int, size int64, partHash string)            // Track the hash value of each part
	HashPartComplete(partsHash string, fileHash string)                 // Mark chunking complete with the final hash values
	UploadID(uploadID string)                                           // Assign an upload ID
	UploadPartComplete(index int, partSize int64, partHash string)      // Track the progress of each uploaded part
	UploadComplete(fileSize int64, streamSize int64, storageSize int64) // Track the overall upload progress
	Complete(size int64, url string, typ int)                           // Mark the upload as complete with final details
}

type UploadLogProgress interface {
	OnProgress(current int64, size int64)
}
