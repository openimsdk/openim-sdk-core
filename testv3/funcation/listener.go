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

package funcation

import (
	"fmt"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"os"
	"time"
)

func (t *initLister) OnUserTokenExpired() {
	log.Info("", utils.GetSelfFuncName())
}
func (t *initLister) OnConnecting() {
	log.Info("", utils.GetSelfFuncName())
}

func (t *initLister) OnConnectSuccess() {
	log.Info("", utils.GetSelfFuncName())
}

func (t *initLister) OnConnectFailed(ErrCode int32, ErrMsg string) {
	log.Info("", utils.GetSelfFuncName(), ErrCode, ErrMsg)
}

func (t *initLister) OnKickedOffline() {
	log.Error("", utils.GetSelfFuncName(), "kick!!!!!!!!!!!!!!!!!!!!")
	os.Exit(-1)
}

func (t *initLister) OnSelfInfoUpdated(info string) {
	log.Info("", utils.GetSelfFuncName())
}

func (t *initLister) OnSuccess() {
	log.Info("", utils.GetSelfFuncName())
}

func (t *initLister) OnError(code int32, msg string) {
	log.Info("", utils.GetSelfFuncName(), code, msg)
}

func (b *BaseSuccessFailed) OnError(errCode int32, errMsg string) {
	b.errCode = -1
	b.errMsg = errMsg
	log.Error("login failed", errCode, errMsg)

}

func (b *BaseSuccessFailed) OnSuccess(data string) {
	b.errCode = 1
	b.successData = data
	log.Info("login success", data, time.Since(b.time))
}

func (c *conversationCallBack) OnError(errCode int32, errMsg string) {
	// TODO implement me
	// panic("implement me")
}

func (c *conversationCallBack) OnSuccess(data string) {
	// TODO implement me
	// panic("implement me")
}

func (c *conversationCallBack) OnNewConversation(conversationList string) {
	log.Info("", "OnNewConversation returnList is ", conversationList)
}

func (c *conversationCallBack) OnConversationChanged(conversationList string) {
	log.Info("", "OnConversationChanged returnList is", conversationList)
}

func (m *MsgListenerCallBak) OnRecvOfflineNewMessage(message string) {
	// TODO implement me
	// panic("implement me")
}

func (m *MsgListenerCallBak) OnRecvNewMessage(message string) {
	log.Info("", message)
}

func (m *MsgListenerCallBak) OnRecvC2CReadReceipt(msgReceiptList string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName())
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

func (m *MsgListenerCallBak) OnMsgDeleted(s string) {}

func (m *MsgListenerCallBak) OnRecvOfflineNewMessages(messageList string) {
	// TODO implement me
	// panic("implement me")
}

func (m *MsgListenerCallBak) OnRecvMessageExtensionsAdded(msgID string, reactionExtensionList string) {
	fmt.Printf("OnRecvMessageExtensionsAdded", msgID, reactionExtensionList)
	log.Info("internal", "OnRecvMessageExtensionsAdded ", msgID, reactionExtensionList)

}

func (m *MsgListenerCallBak) OnRecvGroupReadReceipt(groupMsgReceiptList string) {
	// fmt.Println("OnRecvC2CReadReceipt , ", groupMsgReceiptList)
}
func (m *MsgListenerCallBak) OnNewRecvMessageRevoked(messageRevoked string) {
	// fmt.Println("OnNewRecvMessageRevoked , ", messageRevoked)
}

func (m *MsgListenerCallBak) OnRecvMessageExtensionsChanged(msgID string, reactionExtensionList string) {
	log.Info("internal", "OnRecvMessageExtensionsChanged ", msgID, reactionExtensionList)

}
func (m *MsgListenerCallBak) OnRecvMessageExtensionsDeleted(msgID string, reactionExtensionKeyList string) {
	log.Info("internal", "OnRecvMessageExtensionsDeleted ", msgID, reactionExtensionKeyList)
}

func (userCallback) OnSelfInfoUpdated(callbackData string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackData)
}

func (c *conversationCallBack) OnRecvMessageExtensionsChanged(msgID string, reactionExtensionList string) {
	// panic("implement me")
}

func (c *conversationCallBack) OnRecvMessageExtensionsDeleted(msgID string, reactionExtensionKeyList string) {
	// panic("implement me")
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

func (c *conversationCallBack) OnSyncServerFailed() {
	log.Info("", utils.GetSelfFuncName())
}

func (c *conversationCallBack) OnTotalUnreadMessageCountChanged(totalUnreadCount int32) {
	log.Info("", "OnTotalUnreadMessageCountChanged returnTotalUnreadCount is ", totalUnreadCount)
}

func (t TestSendMsgCallBack) OnError(errCode int32, errMsg string) {
	log.Info("", "OnError is ", errCode, errMsg)
}

func (t TestSendMsgCallBack) OnSuccess(data string) {
	log.Info("", "OnSuccess is ", data)
}

func (t TestSendMsgCallBack) OnProgress(progress int) {
	log.Info("", "OnProgress is ", progress)
}
