// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package open_im_sdk

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
)

func UpdateFcmToken(callback open_im_sdk_callback.Base, operationID, fmcToken string) {
	if err := CheckResourceLoad(UserForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	UserForSDK.Push().UpdateFcmToken(callback, fmcToken, operationID)
}
func SetAppBadge(callback open_im_sdk_callback.Base, operationID string, appUnreadCount int32) {
	if err := CheckResourceLoad(UserForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	UserForSDK.Push().SetAppBadge(callback, appUnreadCount, operationID)
}
