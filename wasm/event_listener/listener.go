package event_listener

import (
	"open_im_sdk/pkg/utils"
	"syscall/js"
)

type ConnCallback struct {
	uid string
	CallbackWriter
}

func NewConnCallback(funcName string, callback *js.Value) *ConnCallback {
	return &ConnCallback{CallbackWriter: NewEventData(callback).SetEvent(funcName)}
}

func (i *ConnCallback) OnConnecting() {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SendMessage()
}

func (i *ConnCallback) OnConnectSuccess() {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SendMessage()

}
func (i *ConnCallback) OnConnectFailed(errCode int32, errMsg string) {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetErrCode(errCode).SetErrMsg(errMsg).SendMessage()
}

func (i *ConnCallback) OnKickedOffline() {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SendMessage()
}

func (i *ConnCallback) OnUserTokenExpired() {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SendMessage()
}

func (i *ConnCallback) OnSelfInfoUpdated(userInfo string) {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(userInfo).SendMessage()
}

type ConversationCallback struct {
	uid string
	CallbackWriter
}

func NewConversationCallback(callback *js.Value) *ConversationCallback {
	return &ConversationCallback{CallbackWriter: NewEventData(callback)}
}
func (c ConversationCallback) OnSyncServerStart() {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SendMessage()
}

func (c ConversationCallback) OnSyncServerFinish() {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SendMessage()
}

func (c ConversationCallback) OnSyncServerFailed() {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SendMessage()

}

func (c ConversationCallback) OnNewConversation(conversationList string) {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(conversationList).SendMessage()

}

func (c ConversationCallback) OnConversationChanged(conversationList string) {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(conversationList).SendMessage()

}

func (c ConversationCallback) OnTotalUnreadMessageCountChanged(totalUnreadCount int32) {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(totalUnreadCount).SendMessage()
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
