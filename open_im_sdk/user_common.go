package open_im_sdk

import (
	"encoding/json"
	"open_im_sdk/open_im_sdk/log"
	"open_im_sdk/open_im_sdk/utils"
	"runtime"
)

func (u *UserRelated) jsonUnmarshalAndArgsValidate(s string, args interface{}, callback Base, operationID string) error {
	err := json.Unmarshal([]byte(s), args)
	if err != nil {
		if callback != nil {
			log.NewError(operationID, "Unmarshal failed ", err.Error(), s)
			callback.OnError(ErrArgs.ErrCode, ErrArgs.ErrMsg)
			runtime.Goexit()
		} else {
			return utils.wrap(err, "json Unmarshal failed")
		}
	}

	err = u.validate.Struct(args)
	if err != nil {
		if callback != nil {
			log.NewError(operationID, "validate failed ", err.Error(), s)
			callback.OnError(ErrArgs.ErrCode, ErrArgs.ErrMsg)
			runtime.Goexit()
		}
	}
	return utils.wrap(err, "args check failed")
}

func (u *UserRelated) jsonUnmarshal(s string, args interface{}, callback Base, operationID string) error {
	err := json.Unmarshal([]byte(s), args)
	if err != nil {
		if callback != nil {
			log.NewError(operationID, "Unmarshal failed ", err.Error(), s)
			callback.OnError(ErrArgs.ErrCode, ErrArgs.ErrMsg)
			runtime.Goexit()
		} else {
			return utils.wrap(err, "json Unmarshal failed")
		}
	}
	return nil
}
