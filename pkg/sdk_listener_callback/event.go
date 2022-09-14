package sdk_listener_callback

import (
	"open_im_sdk/pkg/utils"
	"runtime"
	"syscall/js"
)

type EventData struct {
	Event       string      `json:"event"`
	ErrCode     int32       `json:"errCode"`
	ErrMsg      string      `json:"errMsg"`
	Data        interface{} `json:"data"`
	OperationID string      `json:"operationID"`
	callback    js.Value
}

func NewEventData(callback js.Value) *EventData {
	return &EventData{callback: callback}
}
func (e *EventData) SendMessage() {
	e.callback.Invoke(utils.StructToJsonString(e))
}
func (e *EventData) SetEvent(event string) *EventData {
	e.Event = event
	return e
}
func (e *EventData) SetData(data interface{}) *EventData {
	e.Data = data
	return e
}
func (e *EventData) SetErrCode(errCode int32) *EventData {
	e.ErrCode = errCode
	return e
}
func (e *EventData) SetErrMsg(errMsg string) *EventData {
	e.ErrMsg = errMsg
	return e
}
func (e *EventData) SetSelfCallerFuncName() *EventData {
	pc, _, _, _ := runtime.Caller(1)
	e.Event = utils.CleanUpfuncName(runtime.FuncForPC(pc).Name())
	return e
}
