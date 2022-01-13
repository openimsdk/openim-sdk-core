package common

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"runtime"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func CheckAnyErr(callback Base, errCode int32, err error, operationID string) {
	if err != nil {
		if callback != nil {
			errInfo := "operationID[" + operationID + "], " + "info[" + err.Error() + "]"
			log.NewError(operationID, "checkErr ", errInfo)
			callback.OnError(errCode, errInfo)
			runtime.Goexit()
		}
	}
}

func CheckDBErr(callback Base, err error, operationID string) {
	if err != nil {
		if callback != nil {
			errInfo := operationID + err.Error() + constant.ErrDB.ErrMsg
			log.NewError(operationID, "checkErr ", errInfo)
			callback.OnError(constant.ErrDB.ErrCode, errInfo)
			runtime.Goexit()
		}
	}
}

func CheckDataErr(callback Base, err error, operationID string) {
	if err != nil {
		if callback != nil {
			log.NewError(operationID, "checkErr ", err, constant.ErrData.ErrCode, constant.ErrData.ErrMsg)
			callback.OnError(constant.ErrData.ErrCode, constant.ErrData.ErrMsg)
			runtime.Goexit()
		}
	}
}

func CheckErr(callback Base, err error, operationID string) {
	if err != nil {
		if callback != nil {
			log.NewError(operationID, "checkErr ", err, constant.ErrDB.ErrCode, constant.ErrDB.ErrMsg)
			callback.OnError(constant.ErrDB.ErrCode, constant.ErrDB.ErrMsg)
			runtime.Goexit()
		}
	}
}

func CheckErrAndResp(callback Base, err error, resp []byte, operationID string) *server_api_params.CommDataResp {
	CheckErr(callback, err, operationID)
	return CheckResp(callback, resp, operationID)
}

func CheckResp(callback Base, resp []byte, operationID string) *server_api_params.CommDataResp {
	var c server_api_params.CommDataResp
	err := json.Unmarshal(resp, &c)
	if err != nil {
		log.NewError(operationID, "Unmarshal ", err)
		callback.OnError(constant.ErrArgs.ErrCode, constant.ErrArgs.ErrMsg)
		runtime.Goexit()
		return nil
	}
	if c.ErrCode != 0 {
		log.NewError(operationID, "errCode ", c.ErrCode, "errMsg ", c.ErrMsg)
		callback.OnError(c.ErrCode, c.ErrMsg)
		runtime.Goexit()
		return nil
	}
	return &c
}

func CheckErrAndRespReturn(err error, resp []byte) (*server_api_params.CommDataResp, error) {
	if err != nil {
		return nil, utils.Wrap(err, "resp failed")
	}
	var c server_api_params.CommDataResp
	err = json.Unmarshal(resp, &c)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}

	if c.ErrCode != 0 {
		return nil, utils.Wrap(errors.New(c.ErrMsg), "")
	}
	return &c, nil
}

func JsonUnmarshalAndArgsValidate(s string, args interface{}, callback Base, operationID string) error {
	err := json.Unmarshal([]byte(s), args)
	if err != nil {
		if callback != nil {
			log.NewError(operationID, "Unmarshal failed ", err.Error(), s)
			callback.OnError(constant.ErrArgs.ErrCode, constant.ErrArgs.ErrMsg)
			runtime.Goexit()
		} else {
			return utils.Wrap(err, "json Unmarshal failed")
		}
	}
	err = validate.Struct(args)
	if err != nil {
		if callback != nil {
			log.NewError(operationID, "validate failed ", err.Error(), s)
			callback.OnError(constant.ErrArgs.ErrCode, constant.ErrArgs.ErrMsg)
			runtime.Goexit()
		}
	}
	return utils.Wrap(err, "args check failed")
}

func JsonUnmarshal(s string, args interface{}, callback Base, operationID string) error {
	err := json.Unmarshal([]byte(s), args)
	if err != nil {
		if callback != nil {
			log.NewError(operationID, "Unmarshal failed ", err.Error(), s)
			callback.OnError(constant.ErrArgs.ErrCode, constant.ErrArgs.ErrMsg)
			runtime.Goexit()
		} else {
			return utils.Wrap(err, "json Unmarshal failed")
		}
	}
	return nil
}
