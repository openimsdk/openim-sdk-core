package event_listener

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
