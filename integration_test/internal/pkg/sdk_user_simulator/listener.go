package sdk_user_simulator

import (
	"context"
	"math/rand"

	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

type conversationCallBack struct {
}

func (c *conversationCallBack) OnSyncServerFailed(reinstalled bool) {
}

func (c *conversationCallBack) OnNewConversation(conversationList string) {
}

func (c *conversationCallBack) OnConversationChanged(conversationList string) {
}

func (c *conversationCallBack) OnTotalUnreadMessageCountChanged(totalUnreadCount int32) {
}

func (c *conversationCallBack) OnRecvMessageExtensionsChanged(msgID string, reactionExtensionList string) {
}

func (c *conversationCallBack) OnRecvMessageExtensionsDeleted(msgID string, reactionExtensionKeyList string) {
}

func (c *conversationCallBack) OnSyncServerProgress(progress int) {
}

func (c *conversationCallBack) OnSyncServerStart(reinstalled bool) {

}

func (c *conversationCallBack) OnSyncServerFinish(reinstalled bool) {

}

func (c *conversationCallBack) OnConversationUserInputStatusChanged(change string) {

}

type userCallback struct {
}

func (c userCallback) OnUserStatusChanged(statusMap string) {

}

func (userCallback) OnSelfInfoUpdated(callbackData string) {

}
func (userCallback) OnUserCommandAdd(callbackData string) {

}
func (userCallback) OnUserCommandUpdate(callbackData string) {

}
func (userCallback) OnUserCommandDelete(callbackData string) {

}

type SingleMessage struct {
	SendID      string
	ClientMsgID string
	Delay       int64
}
type MsgListenerCallBak struct {
	userID      string
	GroupDelay  map[string][]*SingleMessage
	SingleDelay map[string][]*SingleMessage
}

func NewMsgListenerCallBak(userID string) *MsgListenerCallBak {
	return &MsgListenerCallBak{userID: userID,
		GroupDelay:  make(map[string][]*SingleMessage),
		SingleDelay: make(map[string][]*SingleMessage)}
}

func (m *MsgListenerCallBak) OnRecvNewMessage(message string) {
	var sm sdk_struct.MsgStruct
	_ = utils.JsonStringToStruct(message, &sm)

	if rand.Float64() < config.CheckMsgRate && sm.ContentType == constant.Text {
		rev := utils.GetCurrentTimestampByMill()
		stm := &vars.StatMsg{
			CostTime:    rev - sm.CreateTime,
			ReceiveTime: rev,
			Msg:         &sm,
		}
		select {
		case vars.RecvMsgConsuming <- stm:
		default:
		}
	}
	switch sm.SessionType {
	case constant.SingleChatType:
		m.SingleDelay[sm.SendID] =
			append(m.SingleDelay[sm.SendID], &SingleMessage{SendID: sm.SendID, ClientMsgID: sm.ClientMsgID, Delay: GetRelativeServerTime() - sm.SendTime})
	case constant.ReadGroupChatType:
		m.GroupDelay[sm.GroupID] =
			append(m.GroupDelay[sm.GroupID], &SingleMessage{SendID: sm.SendID, ClientMsgID: sm.ClientMsgID, Delay: GetRelativeServerTime() - sm.SendTime})
	default:
	}

}

func (m *MsgListenerCallBak) OnRecvC2CReadReceipt(msgReceiptList string) {
}

func (m *MsgListenerCallBak) OnMsgDeleted(s string) {}

func (m *MsgListenerCallBak) OnRecvOfflineNewMessage(message string) {
}

func (m *MsgListenerCallBak) OnRecvMessageExtensionsAdded(msgID string, reactionExtensionList string) {

}

func (m *MsgListenerCallBak) OnRecvGroupReadReceipt(groupMsgReceiptList string) {
}
func (m *MsgListenerCallBak) OnNewRecvMessageRevoked(messageRevoked string) {
}

func (m *MsgListenerCallBak) OnRecvMessageExtensionsChanged(msgID string, reactionExtensionList string) {

}
func (m *MsgListenerCallBak) OnRecvMessageExtensionsDeleted(msgID string, reactionExtensionKeyList string) {
}

func (m *MsgListenerCallBak) OnRecvOnlineOnlyMessage(message string) {

}

type testFriendshipListener struct {
}

func (testFriendshipListener) OnFriendApplicationAdded(callbackInfo string) {

}
func (testFriendshipListener) OnFriendApplicationDeleted(callbackInfo string) {

}

func (testFriendshipListener) OnFriendApplicationAccepted(callbackInfo string) {

}

func (testFriendshipListener) OnFriendApplicationRejected(callbackInfo string) {

}

func (testFriendshipListener) OnFriendAdded(callbackInfo string) {
}

func (testFriendshipListener) OnFriendDeleted(callbackInfo string) {

}

func (testFriendshipListener) OnBlackAdded(callbackInfo string) {

}
func (testFriendshipListener) OnBlackDeleted(callbackInfo string) {

}

func (testFriendshipListener) OnFriendInfoChanged(callbackInfo string) {

}

func (testFriendshipListener) OnSuccess() {

}

func (testFriendshipListener) OnError(code int32, msg string) {

}

type testGroupListener struct {
}

func (testGroupListener) OnJoinedGroupAdded(callbackInfo string) {

}
func (testGroupListener) OnJoinedGroupDeleted(callbackInfo string) {

}

func (testGroupListener) OnGroupMemberAdded(callbackInfo string) {

}
func (testGroupListener) OnGroupMemberDeleted(callbackInfo string) {

}

func (testGroupListener) OnGroupApplicationAdded(callbackInfo string) {

}
func (testGroupListener) OnGroupApplicationDeleted(callbackInfo string) {

}

func (testGroupListener) OnGroupInfoChanged(callbackInfo string) {

}
func (testGroupListener) OnGroupMemberInfoChanged(callbackInfo string) {

}

func (testGroupListener) OnGroupApplicationAccepted(callbackInfo string) {

}
func (testGroupListener) OnGroupApplicationRejected(callbackInfo string) {

}

func (testGroupListener) OnGroupDismissed(callbackInfo string) {

}

type testConnListener struct {
	UserID string
}

func (t *testConnListener) OnUserTokenInvalid(errMsg string) {
	log.ZError(context.TODO(), "user token invalid", errs.New("user token invalid").Wrap(), "userID", t.UserID)
}

func (t *testConnListener) OnUserTokenExpired() {
	log.ZError(context.TODO(), "user token expired", errs.New("user token expired").Wrap(), "userID", t.UserID)
}
func (t *testConnListener) OnConnecting() {

}

func (t *testConnListener) OnConnectSuccess() {
	vars.NowLoginNum.Add(1)
}

func (t *testConnListener) OnConnectFailed(ErrCode int32, ErrMsg string) {
	log.ZError(context.TODO(), "connect failed", errs.NewCodeError(int(ErrCode), ErrMsg), "userID", t.UserID)
}

func (t *testConnListener) OnKickedOffline() {
	log.ZError(context.TODO(), "kicked offline", errs.New("kicked offline").Wrap(), "userID", t.UserID)
}

func (t *testConnListener) OnSelfInfoUpdated(info string) {

}
func (t *testConnListener) OnUserCommandAdd(info string) {

}
func (t *testConnListener) OnUserCommandUpdate(info string) {

}
func (t *testConnListener) OnUserCommandDelete(info string) {

}
func (t *testConnListener) OnSuccess() {

}

func (t *testConnListener) OnError(code int32, msg string) {

}
