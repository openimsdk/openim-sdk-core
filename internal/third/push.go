// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package third

import (
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"

	"github.com/openimsdk/openim-sdk-core/v3/internal/file"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/log"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdk_params_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/server_api_params"
)

type Third struct {
	platformID   int32
	loginUserID  string
	version      string
	LogFilePath  string
	fileUploader *file.File
}

func NewThird(platformID int32, loginUserID, version, LogFilePath string, fileUploader *file.File) *Third {
	return &Third{platformID: platformID, loginUserID: loginUserID, version: version, LogFilePath: LogFilePath, fileUploader: fileUploader}
}

func (c *Third) UpdateFcmToken(callback open_im_sdk_callback.Base, fcmToken, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "UpdateFcmToken args: ", fcmToken)
		c.fmcUpdateToken(callback, fcmToken, operationID)
		callback.OnSuccess(sdk_params_callback.UpdateFcmTokenCallback)
		log.NewInfo(operationID, "UpdateFcmToken callback: ", sdk_params_callback.UpdateFcmTokenCallback)
	}()

}

func (c *Third) fmcUpdateToken(callback open_im_sdk_callback.Base, fcmToken, operationID string) {
	apiReq := server_api_params.FcmUpdateTokenReq{}
	apiReq.OperationID = operationID
	apiReq.Platform = int(c.platformID)
	apiReq.FcmToken = fcmToken
	//c.p.PostFatalCallback(callback, constant.FcmUpdateTokenRouter, apiReq, nil, apiReq.OperationID)
}
func (c *Third) SetAppBadge(callback open_im_sdk_callback.Base, appUnreadCount int32, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "SetAppBadge args: ", appUnreadCount)
		c.setAppBadge(callback, appUnreadCount, operationID)
		callback.OnSuccess(sdk_params_callback.SetAppBadgeCallback)
		log.NewInfo(operationID, "SetAppBadge callback: ", sdk_params_callback.SetAppBadgeCallback)
	}()
}
func (c *Third) setAppBadge(callback open_im_sdk_callback.Base, appUnreadCount int32, operationID string) {
	apiReq := server_api_params.SetAppBadgeReq{}
	apiReq.OperationID = operationID
	apiReq.FromUserID = c.loginUserID
	apiReq.AppUnreadCount = appUnreadCount
	//c.p.PostFatalCallback(callback, constant.SetAppBadgeRouter, apiReq, nil, apiReq.OperationID)
}
