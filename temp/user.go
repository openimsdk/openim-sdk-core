package open_im_sdk


func GetUsersInfo(callback open_im_sdk_callback.Base, operationID string, userIDList string) {
	call(callback, operationID, userForSDK.User().GetUsersInfo, userIDList)
}


func SetSelfInfo(callback open_im_sdk_callback.Base, operationID string, userInfo string) {
	call(callback, operationID, userForSDK.User().SetSelfInfo, userInfo)
}


func GetSelfUserInfo(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, userForSDK.User().GetSelfUserInfo)
}


func UpdateMsgSenderInfo(callback open_im_sdk_callback.Base, operationID string, nickname string, faceURL string) {
	call(callback, operationID, userForSDK.User().UpdateMsgSenderInfo, nickname,faceURL)
}

