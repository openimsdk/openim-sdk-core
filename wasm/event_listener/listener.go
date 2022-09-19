package event_listener

import (
	"open_im_sdk/pkg/sdk_listener_callback"
	"syscall/js"
)

type InitCallback struct {
	uid       string
	eventData *sdk_listener_callback.EventData
}

func NewInitCallback(callback js.Value) *InitCallback {
	return &InitCallback{eventData: sdk_listener_callback.NewEventData(callback)}
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
	funcName    string
	operationID string
	eventData   *sdk_listener_callback.EventData
}

func NewBaseCallback(funcName string, operationID string, callback js.Value) *BaseCallback {
	return &BaseCallback{funcName: funcName, operationID: operationID, eventData: sdk_listener_callback.NewEventData(callback)}
}

func (b *BaseCallback) OnError(errCode int32, errMsg string) {
	b.eventData.SetEvent(b.funcName).SetErrCode(errCode).SetErrMsg(errMsg).SendMessage()
}
func (b *BaseCallback) OnSuccess(data string) {
	b.eventData.SetEvent(b.funcName).SetData(data).SendMessage()
}
