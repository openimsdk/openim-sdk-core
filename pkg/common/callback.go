package common

type Base interface {
	OnError(errCode int32, errMsg string)
	OnSuccess(data string)
}
type SendMsgCallBack interface {
	Base
	OnProgress(progress int)
}
