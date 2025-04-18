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

func GetUsersInfo(callback open_im_sdk_callback.Base, operationID string, userIDs string) {
	call(callback, operationID, IMUserContext.User().GetUsersInfo, userIDs)
}

// SetSelfInfo sets the user's own information.
func SetSelfInfo(callback open_im_sdk_callback.Base, operationID string, userInfo string) {
	call(callback, operationID, IMUserContext.User().SetSelfInfo, userInfo)
}

// GetSelfUserInfo obtains the user's own information.
func GetSelfUserInfo(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, IMUserContext.User().GetSelfUserInfo)
}

func GetUserClientConfig(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, IMUserContext.User().GetUserClientConfig)
}
