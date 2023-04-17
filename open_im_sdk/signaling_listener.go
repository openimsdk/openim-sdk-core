package open_im_sdk

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/log"
)

func SetSignalingListener(callback open_im_sdk_callback.OnSignalingListener) {
	if callback == nil || UserForSDK == nil {
		log.Error("callback or UserForSDK is nil")
		return
	}
	UserForSDK.SetSignalingListener(callback)
}

func SetSignalingListenerForService(callback open_im_sdk_callback.OnSignalingListener) {
	if callback == nil || UserForSDK == nil {
		log.Error("callback or UserForSDK is nil")
		return
	}
	UserForSDK.SetSignalingListenerForService(callback)
}

func SetListenerForService(callback open_im_sdk_callback.OnListenerForService) {
	if callback == nil || UserForSDK == nil {
		log.Error("callback or UserForSDK is nil")
		return
	}
	UserForSDK.SetListenerForService(callback)
}
