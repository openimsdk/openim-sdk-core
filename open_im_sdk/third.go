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

package open_im_sdk

import (
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
)

func UpdateFcmToken(callback open_im_sdk_callback.Base, operationID, userIDList string) {
	call(callback, operationID, UserForSDK.Third().UpdateFcmToken, userIDList)
}

func SetAppBadge(callback open_im_sdk_callback.Base, operationID string, appUnreadCount int32) {
	call(callback, operationID, UserForSDK.Third().SetAppBadge, appUnreadCount)
}

func UploadLogs(callback open_im_sdk_callback.Base, operationID string, uploadlogParams string) {
	call(callback, operationID, UserForSDK.Third().UploadLogs, uploadlogParams)
}
