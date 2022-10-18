package event_listener

import (
	"syscall/js"
)

type CallbackWriter interface {
	SendMessage()
	SetEvent(event string) CallbackWriter
	SetData(data interface{}) CallbackWriter
	SetErrCode(errCode int32) CallbackWriter
	SetOperationID(operationID string) CallbackWriter
	SetErrMsg(errMsg string) CallbackWriter
	GetOperationID() string
	HandlerFunc() interface{}
}
type Caller interface {
	NewCaller(funcName interface{}, callback CallbackWriter, arguments *[]js.Value) Caller
	AsyncCallWithCallback() (result []interface{})
	AsyncCallWithOutCallback() (fn func() interface{})
}
