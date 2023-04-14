package open_im_sdk

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/log"
)

func SetSignalingListener(callback open_im_sdk_callback.OnSignalingListener) {
	if callback == nil || userForSDK == nil {
		log.Error("callback or userForSDK is nil")
		return
	}
	userForSDK.SetSignalingListener(callback)
}

func SetSignalingListenerForService(callback open_im_sdk_callback.OnSignalingListener) {
	if callback == nil || userForSDK == nil {
		log.Error("callback or userForSDK is nil")
		return
	}
	userForSDK.SetSignalingListenerForService(callback)
}

func SetListenerForService(callback open_im_sdk_callback.OnListenerForService) {
	if callback == nil || userForSDK == nil {
		log.Error("callback or userForSDK is nil")
		return
	}
	userForSDK.SetListenerForService(callback)
}
