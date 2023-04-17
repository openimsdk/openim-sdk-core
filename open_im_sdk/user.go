package open_im_sdk

import "open_im_sdk/open_im_sdk_callback"

func GetUsersInfo(callback open_im_sdk_callback.Base, operationID string, userIDList string) {
	call(callback, operationID, UserForSDK.User().GetUsersInfo, userIDList)
}

func SetSelfInfo(callback open_im_sdk_callback.Base, operationID string, userInfo string) {
	call(callback, operationID, UserForSDK.User().SetSelfInfo, userInfo)
}

func GetSelfUserInfo(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.User().GetSelfUserInfo)
}

func UpdateMsgSenderInfo(callback open_im_sdk_callback.Base, operationID string, nickname, faceURL string) {
	call(callback, operationID, UserForSDK.User().UpdateMsgSenderInfo, nickname, faceURL)
}
