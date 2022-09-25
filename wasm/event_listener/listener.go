package event_listener

import (
	"syscall/js"
)

type ConnCallback struct {
	uid string
	CallbackWriter
}

func NewConnCallback(callback *js.Value) *ConnCallback {
	return &ConnCallback{CallbackWriter: NewEventData(callback)}
}

func (i *ConnCallback) OnConnecting() {
	i.CallbackWriter.SetSelfCallerFuncName().SendMessage()
}

func (i *ConnCallback) OnConnectSuccess() {
	i.CallbackWriter.SetSelfCallerFuncName().SendMessage()

}
func (i *ConnCallback) OnConnectFailed(errCode int32, errMsg string) {
	i.CallbackWriter.SetSelfCallerFuncName().SetErrCode(errCode).SetErrMsg(errMsg).SendMessage()
}

func (i *ConnCallback) OnKickedOffline() {
	i.CallbackWriter.SetSelfCallerFuncName().SendMessage()
}

func (i *ConnCallback) OnUserTokenExpired() {
	i.CallbackWriter.SetSelfCallerFuncName().SendMessage()
}

func (i *ConnCallback) OnSelfInfoUpdated(userInfo string) {
	i.CallbackWriter.SetSelfCallerFuncName().SetData(userInfo).SendMessage()
}

type ConversationCallback struct {
	uid string
	CallbackWriter
}

func NewConversationCallback(callback *js.Value) *ConversationCallback {
	return &ConversationCallback{CallbackWriter: NewEventData(callback)}
}
func (c ConversationCallback) OnSyncServerStart() {
	c.CallbackWriter.SetSelfCallerFuncName().SendMessage()
}

func (c ConversationCallback) OnSyncServerFinish() {
	c.CallbackWriter.SetSelfCallerFuncName().SendMessage()
}

func (c ConversationCallback) OnSyncServerFailed() {
	c.CallbackWriter.SetSelfCallerFuncName().SendMessage()

}

func (c ConversationCallback) OnNewConversation(conversationList string) {
	c.CallbackWriter.SetSelfCallerFuncName().SetData(conversationList).SendMessage()

}

func (c ConversationCallback) OnConversationChanged(conversationList string) {
	c.CallbackWriter.SetSelfCallerFuncName().SetData(conversationList).SendMessage()

}

func (c ConversationCallback) OnTotalUnreadMessageCountChanged(totalUnreadCount int32) {
	c.CallbackWriter.SetSelfCallerFuncName().SetData(totalUnreadCount).SendMessage()
}

type BaseCallback struct {
	CallbackWriter
}

func (b *BaseCallback) EventData() CallbackWriter {
	return b.CallbackWriter
}

func NewBaseCallback(funcName string, callback *js.Value) *BaseCallback {
	return &BaseCallback{CallbackWriter: NewEventData(callback).SetEvent(funcName)}
}

func (b *BaseCallback) OnError(errCode int32, errMsg string) {
	b.CallbackWriter.SetErrCode(errCode).SetErrMsg(errMsg).SendMessage()
}
func (b *BaseCallback) OnSuccess(data string) {
	b.CallbackWriter.SetData(data).SendMessage()
}

type SendMessageCallback struct {
	BaseCallback
	clientMsgID string
}

func (s *SendMessageCallback) SetClientMsgID(clientMsgID string) {
	s.clientMsgID = clientMsgID
}
func NewSendMessageCallback(funcName string, callback *js.Value) *SendMessageCallback {
	return &SendMessageCallback{BaseCallback: BaseCallback{CallbackWriter: NewEventData(callback).SetEvent(funcName)}}
}

func (s SendMessageCallback) OnProgress(progress int) {
	panic("implement me")
}
