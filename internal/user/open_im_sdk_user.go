package user

import (
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/utils"
)

func (u *User) GetUsersInfo(callback common.Base, userIDList string, operationID string) {
	go func() {
		log.NewInfo(operationID, utils.RunFuncName(), "args: ", userIDList)
		var unmarshalParam sdk_params_callback.GetUsersInfoParam
		common.JsonUnmarshalAndArgsValidate(userIDList, &unmarshalParam, callback, operationID)
		result := u.GetUsersInfoFromSvr(callback, unmarshalParam, operationID)
		callback.OnSuccess(utils.StructToJsonString(utils.StructToJsonString(result)))
		log.NewInfo(operationID, utils.RunFuncName(), "callback: ", utils.StructToJsonString(result))
	}()
}

func (u *User) GetSelfUserInfo(callback common.Base, operationID string) {
	go func() {
		log.NewInfo(operationID, utils.RunFuncName(), "args: ")
		result := u.getSelfUserInfo(callback, operationID)
		log.NewInfo(operationID, utils.RunFuncName(), "callback: ", utils.StructToJsonString(result))
	}()
}

func (u *User) SetSelfInfo(callback common.Base, userInfo string, operationID string) {
	go func() {
		log.NewInfo(operationID, utils.RunFuncName(), "args: ", userInfo)
		var unmarshalParam sdk_params_callback.SetSelfUserInfoParam
		common.JsonUnmarshalAndArgsValidate(userInfo, &unmarshalParam, callback, operationID)
		u.updateSelfUserInfo(callback, unmarshalParam, operationID)
		callback.OnSuccess(utils.StructToJsonString(utils.StructToJsonString(sdk_params_callback.SetSelfUserInfoCallback)))
		log.NewInfo(operationID, utils.RunFuncName(), "callback: ")
	}()
}
