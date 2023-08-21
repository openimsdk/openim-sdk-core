package sdk_user_simulator

import (
	"fmt"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
)

type conversationCallBack struct {
	SyncFlag int
}

func (c *conversationCallBack) OnSyncServerFailed() {
	log.Info("", utils.GetSelfFuncName())
}

func (c *conversationCallBack) OnNewConversation(conversationList string) {
	log.Info("", utils.GetSelfFuncName())
}

func (c *conversationCallBack) OnConversationChanged(conversationList string) {
	log.Info("", utils.GetSelfFuncName())
}

func (c *conversationCallBack) OnTotalUnreadMessageCountChanged(totalUnreadCount int32) {
	log.Info("", utils.GetSelfFuncName())
}

func (c *conversationCallBack) OnRecvMessageExtensionsChanged(msgID string, reactionExtensionList string) {
	log.Info("", utils.GetSelfFuncName())
}

func (c *conversationCallBack) OnRecvMessageExtensionsDeleted(msgID string, reactionExtensionKeyList string) {
	log.Info("", utils.GetSelfFuncName())
}

func (c *conversationCallBack) OnSyncServerProgress(progress int) {
	log.Info("", utils.GetSelfFuncName())
}

func (c *conversationCallBack) OnSyncServerStart() {

}

func (c *conversationCallBack) OnSyncServerFinish() {
	c.SyncFlag = 1
	log.Info("", utils.GetSelfFuncName())

}

type userCallback struct {
}

func (c userCallback) OnUserStatusChanged(statusMap string) {
	log.Info("", utils.GetSelfFuncName())
}

func (userCallback) OnSelfInfoUpdated(callbackData string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackData)
}

type MsgListenerCallBak struct {
}

func (m *MsgListenerCallBak) OnRecvNewMessage(message string) {
	log.Info("", utils.GetSelfFuncName())
}

func (m *MsgListenerCallBak) OnRecvC2CReadReceipt(msgReceiptList string) {
	log.Info("", utils.GetSelfFuncName())
}

func (m *MsgListenerCallBak) OnMsgDeleted(s string) {}

func (m *MsgListenerCallBak) OnRecvOfflineNewMessage(message string) {
	log.Info("", utils.GetSelfFuncName())
}

func (m *MsgListenerCallBak) OnRecvMessageExtensionsAdded(msgID string, reactionExtensionList string) {
	fmt.Printf("OnRecvMessageExtensionsAdded", msgID, reactionExtensionList)
	log.Info("internal", "OnRecvMessageExtensionsAdded ", msgID, reactionExtensionList)

}

func (m *MsgListenerCallBak) OnRecvGroupReadReceipt(groupMsgReceiptList string) {
	//fmt.Println("OnRecvC2CReadReceipt , ", groupMsgReceiptList)
}
func (m *MsgListenerCallBak) OnNewRecvMessageRevoked(messageRevoked string) {
	//fmt.Println("OnNewRecvMessageRevoked , ", messageRevoked)
}

func (m *MsgListenerCallBak) OnRecvMessageExtensionsChanged(msgID string, reactionExtensionList string) {
	log.Info("internal", "OnRecvMessageExtensionsChanged ", msgID, reactionExtensionList)

}
func (m *MsgListenerCallBak) OnRecvMessageExtensionsDeleted(msgID string, reactionExtensionKeyList string) {
	log.Info("internal", "OnRecvMessageExtensionsDeleted ", msgID, reactionExtensionKeyList)
}

type testFriendListener struct {
	x int
}

func (testFriendListener) OnFriendApplicationAdded(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)
}
func (testFriendListener) OnFriendApplicationDeleted(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)
}

func (testFriendListener) OnFriendApplicationAccepted(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)
}

func (testFriendListener) OnFriendApplicationRejected(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)
}

func (testFriendListener) OnFriendAdded(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)
}

func (testFriendListener) OnFriendDeleted(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)
}

func (testFriendListener) OnBlackAdded(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)
}
func (testFriendListener) OnBlackDeleted(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)
}

func (testFriendListener) OnFriendInfoChanged(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)
}

func (testFriendListener) OnSuccess() {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName())
}

func (testFriendListener) OnError(code int32, msg string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), code, msg)
}

type testGroupListener struct {
}

func (testGroupListener) OnJoinedGroupAdded(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}
func (testGroupListener) OnJoinedGroupDeleted(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}

func (testGroupListener) OnGroupMemberAdded(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}
func (testGroupListener) OnGroupMemberDeleted(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}

func (testGroupListener) OnGroupApplicationAdded(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}
func (testGroupListener) OnGroupApplicationDeleted(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}

func (testGroupListener) OnGroupInfoChanged(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}
func (testGroupListener) OnGroupMemberInfoChanged(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}

func (testGroupListener) OnGroupApplicationAccepted(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}
func (testGroupListener) OnGroupApplicationRejected(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}

func (testGroupListener) OnGroupDismissed(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}

type testConnListener struct {
}

func (t *testConnListener) OnUserTokenExpired() {
	log.Info("", utils.GetSelfFuncName())
}
func (t *testConnListener) OnConnecting() {
	log.Info("", utils.GetSelfFuncName())
}

func (t *testConnListener) OnConnectSuccess() {
	log.Info("", utils.GetSelfFuncName())
}

func (t *testConnListener) OnConnectFailed(ErrCode int32, ErrMsg string) {
	log.Info("", utils.GetSelfFuncName(), ErrCode, ErrMsg)
}

func (t *testConnListener) OnKickedOffline() {
	log.Info("", utils.GetSelfFuncName())
}

func (t *testConnListener) OnSelfInfoUpdated(info string) {
	log.Info("", utils.GetSelfFuncName())
}

func (t *testConnListener) OnSuccess() {
	log.Info("", utils.GetSelfFuncName())
}

func (t *testConnListener) OnError(code int32, msg string) {
	log.Info("", utils.GetSelfFuncName(), code, msg)
}
