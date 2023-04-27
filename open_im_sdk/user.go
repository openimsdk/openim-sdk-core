package open_im_sdk

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
)

// GetUsersInfo obtains the information about multiple users.
func GetUsersInfo(callback open_im_sdk_callback.Base, operationID string, userIDList string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Full().GetUsersInfo(callback, userIDList, operationID)
}

// SetSelfInfo sets the user's own information.
func SetSelfInfo(callback open_im_sdk_callback.Base, operationID string, userInfo string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.User().SetSelfInfo(callback, userInfo, operationID)
}

// GetSelfUserInfo obtains the user's own information.
func GetSelfUserInfo(callback open_im_sdk_callback.Base, operationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.User().GetSelfUserInfo(callback, operationID)
}
