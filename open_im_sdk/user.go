package open_im_sdk

import "open_im_sdk/open_im_sdk_callback"

func GetUsersInfo(callback open_im_sdk_callback.Base, operationID string, userIDList string) {
	call(callback, operationID, userForSDK.User().GetUsersInfo, userIDList)
}

func SetSelfInfo(callback open_im_sdk_callback.Base, operationID string, userInfo string) {
	call(callback, operationID, userForSDK.User().SetSelfInfo, userInfo)
}

func GetSelfUserInfo(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, userForSDK.User().GetSelfUserInfo)
}
