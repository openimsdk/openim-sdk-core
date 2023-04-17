package open_im_sdk

import "open_im_sdk/open_im_sdk_callback"

func GetWorkMomentsUnReadCount(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.WorkMoments().GetWorkMomentsUnReadCount)
}

func GetWorkMomentsNotification(callback open_im_sdk_callback.Base, operationID string, offset int, count int) {
	call(callback, operationID, UserForSDK.WorkMoments().GetWorkMomentsNotification, offset, count)
}

func ClearWorkMomentsNotification(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.WorkMoments().ClearWorkMomentsNotification)
}
