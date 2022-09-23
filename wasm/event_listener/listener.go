package event_listener

import (
	"syscall/js"
)

type InitCallback struct {
	uid       string
	eventData *EventData
}

func NewInitCallback(callback *js.Value) *InitCallback {
	return &InitCallback{eventData: NewEventData(callback)}
}

func (i *InitCallback) OnConnecting() {
	i.eventData.SetSelfCallerFuncName().SendMessage()
}

func (i *InitCallback) OnConnectSuccess() {
	i.eventData.SetSelfCallerFuncName().SendMessage()

}
func (i *InitCallback) OnConnectFailed(errCode int32, errMsg string) {
	i.eventData.SetSelfCallerFuncName().SetErrCode(errCode).SetErrMsg(errMsg).SendMessage()
}

func (i *InitCallback) OnKickedOffline() {
	i.eventData.SetSelfCallerFuncName().SendMessage()
}

func (i *InitCallback) OnUserTokenExpired() {
	i.eventData.SetSelfCallerFuncName().SendMessage()
}

func (i *InitCallback) OnSelfInfoUpdated(userInfo string) {
	i.eventData.SetSelfCallerFuncName().SetData(userInfo).SendMessage()
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
