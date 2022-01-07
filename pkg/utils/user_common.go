package utils

import (
	"encoding/json"
	"open_im_sdk/internal/controller/init"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"runtime"
)

func (u *open_im_sdk.UserRelated) jsonUnmarshalAndArgsValidate(s string, args interface{}, callback init.Base, operationID string) error {
	err := json.Unmarshal([]byte(s), args)
	if err != nil {
		if callback != nil {
			log.NewError(operationID, "Unmarshal failed ", err.Error(), s)
			callback.OnError(constant.ErrArgs.ErrCode, constant.ErrArgs.ErrMsg)
			runtime.Goexit()
		} else {
			return wrap(err, "json Unmarshal failed")
		}
	}

	err = u.validate.Struct(args)
	if err != nil {
		if callback != nil {
			log.NewError(operationID, "validate failed ", err.Error(), s)
			callback.OnError(constant.ErrArgs.ErrCode, constant.ErrArgs.ErrMsg)
			runtime.Goexit()
		}
	}
	return wrap(err, "args check failed")
}

func (u *open_im_sdk.UserRelated) jsonUnmarshal(s string, args interface{}, callback init.Base, operationID string) error {
	err := json.Unmarshal([]byte(s), args)
	if err != nil {
		if callback != nil {
			log.NewError(operationID, "Unmarshal failed ", err.Error(), s)
			callback.OnError(constant.ErrArgs.ErrCode, constant.ErrArgs.ErrMsg)
			runtime.Goexit()
		} else {
			return wrap(err, "json Unmarshal failed")
		}
	}
	return nil
}
