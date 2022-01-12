package common

type Base interface {
	OnError(errCode int32, errMsg string)
	OnSuccess(data string)
}
