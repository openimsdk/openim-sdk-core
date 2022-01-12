package common

import "open_im_sdk/pkg/utils"

type Base interface {
	OnError(errCode int32, errMsg string)
	OnSuccess(data string)
}
