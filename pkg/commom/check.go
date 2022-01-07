package commom

import (
	"encoding/json"
	"errors"
	"open_im_sdk/pkg/constant"
	log2 "open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"runtime"
)

func CheckErr(callback Base, err error, operationID string) {
	if err != nil {
		if callback != nil {
			log2.NewError(operationID, "checkErr ", err, constant.ErrDB.ErrCode, constant.ErrDB.ErrMsg)
			callback.OnError(constant.ErrDB.ErrCode, constant.ErrDB.ErrMsg)
			runtime.Goexit()
		}
	}
}

func checkErrAndResp(callback Base, err error, resp []byte, operationID string) *server_api_params.CommDataResp {
	CheckErr(callback, err, operationID)
	return checkResp(callback, resp, operationID)
}

func checkResp(callback Base, resp []byte, operationID string) *server_api_params.CommDataResp {
	var c server_api_params.CommDataResp
	err := json.Unmarshal(resp, &c)
	if err != nil {
		log2.NewError(operationID, "Unmarshal ", err)
		callback.OnError(constant.ErrArgs.ErrCode, constant.ErrArgs.ErrMsg)
		runtime.Goexit()
		return nil
	}

	if c.ErrCode != 0 {
		log2.NewError(operationID, "errCode ", c.ErrCode, "errMsg ", c.ErrMsg)
		callback.OnError(c.ErrCode, c.ErrMsg)
		runtime.Goexit()
		return nil
	}
	return &c
}

func checkErrAndRespReturn(err error, resp []byte, operationID string) (*server_api_params.CommDataResp, error) {
	if err != nil {
		log2.NewError(operationID, "checkErr ", err)
		return nil, err
	}
	var c server_api_params.CommDataResp
	err = json.Unmarshal(resp, &c)
	if err != nil {
		log2.NewError(operationID, "Unmarshal ", err)
		return nil, err
	}

	if c.ErrCode != 0 {
		log2.NewError(operationID, "errCode ", c.ErrCode, "errMsg ", c.ErrMsg)
		return nil, errors.New(c.ErrMsg)
	}
	return &c, nil
}
