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
	OnGroupMemberInfoChanged(groupMemberInfo string)
	OnGroupApplicationAccepted(groupApplication string)
	OnGroupApplicationRejected(groupApplication string)
}
type OnFriendshipListener interface {
	OnFriendApplicationAdded(friendApplication string)
	OnFriendApplicationDeleted(friendApplication string)
	OnFriendApplicationAccepted(groupApplication string)
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
	//deprecated
	OnRecvMessageRevoked(msgID string)
	OnNewRecvMessageRevoked(messageRevoked string)
	OnRecvMessageExtensionsChanged(msgID string, reactionExtensionList string)
	OnRecvMessageExtensionsDeleted(msgID string, reactionExtensionKeyList string)
	OnRecvMessageExtensionsAdded(msgID string, reactionExtensionList string)
}

type OnBatchMsgListener interface {
	OnRecvNewMessages(messageList string)
}

type OnUserListener interface {
	OnSelfInfoUpdated(userInfo string)
}

type OnOrganizationListener interface {
	OnOrganizationUpdated()
}

type OnWorkMomentsListener interface {
	OnRecvNewNotification()
}

type OnCustomBusinessListener interface {
	OnRecvCustomBusinessMessage(businessMessage string)
}
type OnMessageKvInfoListener interface {
	OnMessageKvInfoChanged(messageChangedList string)
}
