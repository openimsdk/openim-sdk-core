package event_listener

import (
	"syscall/js"
)

type ConnCallback struct {
	uid       string
	eventData *EventData
}

func NewConnCallback(callback *js.Value) *ConnCallback {
	return &ConnCallback{eventData: NewEventData(callback)}
}

func (i *ConnCallback) OnConnecting() {
	i.eventData.SetSelfCallerFuncName().SendMessage()
}

func (i *ConnCallback) OnConnectSuccess() {
	i.eventData.SetSelfCallerFuncName().SendMessage()

}
func (i *ConnCallback) OnConnectFailed(errCode int32, errMsg string) {
	i.eventData.SetSelfCallerFuncName().SetErrCode(errCode).SetErrMsg(errMsg).SendMessage()
}

func (i *ConnCallback) OnKickedOffline() {
	i.eventData.SetSelfCallerFuncName().SendMessage()
}

func (i *ConnCallback) OnUserTokenExpired() {
	i.eventData.SetSelfCallerFuncName().SendMessage()
}

func (i *ConnCallback) OnSelfInfoUpdated(userInfo string) {
	i.eventData.SetSelfCallerFuncName().SetData(userInfo).SendMessage()
}

type ConversationCallback struct {
	uid       string
	eventData *EventData
}

func NewConversationCallback(callback *js.Value) *ConversationCallback {
	return &ConversationCallback{eventData: NewEventData(callback)}
}
func (c ConversationCallback) OnSyncServerStart() {
	c.eventData.SetSelfCallerFuncName().SendMessage()
}

func (c ConversationCallback) OnSyncServerFinish() {
	c.eventData.SetSelfCallerFuncName().SendMessage()
}

func (c ConversationCallback) OnSyncServerFailed() {
	c.eventData.SetSelfCallerFuncName().SendMessage()

}

func (c ConversationCallback) OnNewConversation(conversationList string) {
	c.eventData.SetSelfCallerFuncName().SetData(conversationList).SendMessage()

}

func (c ConversationCallback) OnConversationChanged(conversationList string) {
	c.eventData.SetSelfCallerFuncName().SetData(conversationList).SendMessage()

}

func (c ConversationCallback) OnTotalUnreadMessageCountChanged(totalUnreadCount int32) {
	c.eventData.SetSelfCallerFuncName().SetData(totalUnreadCount).SendMessage()
}

type BaseCallback struct {
	funcName  string
	eventData *EventData
}

func (b *BaseCallback) EventData() *EventData {
	return b.eventData
}

func NewBaseCallback(funcName string, callback *js.Value) *BaseCallback {
	return &BaseCallback{funcName: funcName, eventData: NewEventData(callback)}
}

func (b *BaseCallback) OnError(errCode int32, errMsg string) {
	b.eventData.SetEvent(b.funcName).SetErrCode(errCode).SetErrMsg(errMsg).SendMessage()
}
func (b *BaseCallback) OnSuccess(data string) {
	b.eventData.SetEvent(b.funcName).SetData(data).SendMessage()
}

type SendMessageCallback struct {
	BaseCallback
	clientMsgID string
}

func (s *SendMessageCallback) SetClientMsgID(clientMsgID string) {
	s.clientMsgID = clientMsgID
}
func NewSendMessageCallback(funcName string, callback *js.Value) *SendMessageCallback {
	return &SendMessageCallback{BaseCallback: BaseCallback{funcName: funcName, eventData: NewEventData(callback)}}
}

func (s SendMessageCallback) OnProgress(progress int) {
	panic("implement me")
}
