package open_im_sdk

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/log"
)

func SetGroupListener(callback open_im_sdk_callback.OnGroupListener) {
	if callback == nil || userForSDK == nil {
		log.Error("callback or userForSDK is nil")
		return
	}
	userForSDK.SetGroupListener(callback)
}

func SetOrganizationListener(callback open_im_sdk_callback.OnOrganizationListener) {
	if callback == nil || userForSDK == nil {
		log.Error("callback or userForSDK is nil")
		return
	}
	userForSDK.SetOrganizationListener(callback)
}
func SetConversationListener(listener open_im_sdk_callback.OnConversationListener) {
	userForSDK.SetConversationListener(listener)
}
func SetAdvancedMsgListener(listener open_im_sdk_callback.OnAdvancedMsgListener) {
	userForSDK.SetAdvancedMsgListener(listener)
}
func SetBatchMsgListener(listener open_im_sdk_callback.OnBatchMsgListener) {
	userForSDK.SetBatchMsgListener(listener)
}

func SetUserListener(listener open_im_sdk_callback.OnUserListener) {
	userForSDK.SetUserListener(listener)
}

func SetWorkMomentsListener(listener open_im_sdk_callback.OnWorkMomentsListener) {
	userForSDK.SetWorkMomentsListener(listener)
}
