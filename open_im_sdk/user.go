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
	"open_im_sdk/open_im_sdk_callback"
)

func GetUsersInfo(callback open_im_sdk_callback.Base, operationID string, userIDs string) {
	call(callback, operationID, UserForSDK.Full().GetUsersInfo, userIDs)
}

// GetUsersInfoFromSrv obtains the information about multiple users.
func GetUsersInfoFromSrv(callback open_im_sdk_callback.Base, operationID string, userIDs string) {
	call(callback, operationID, UserForSDK.User().GetUsersInfo, userIDs)
}

// SetSelfInfo sets the user's own information.
func SetSelfInfo(callback open_im_sdk_callback.Base, operationID string, userInfo string) {
	call(callback, operationID, UserForSDK.User().SetSelfInfo, userInfo)
}

// GetSelfUserInfo obtains the user's own information.
func GetSelfUserInfo(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.User().GetSelfUserInfo)
}

// UpdateMsgSenderInfo updates the message sender's nickname and face URL.
func UpdateMsgSenderInfo(callback open_im_sdk_callback.Base, operationID string, nickname, faceURL string) {
	call(callback, operationID, UserForSDK.User().UpdateMsgSenderInfo, nickname, faceURL)
}

// SubscribeUsersStatus Presence status of subscribed users.
func SubscribeUsersStatus(callback open_im_sdk_callback.Base, operationID string, userID string, userIDs []string) {
	call(callback, operationID, UserForSDK.User().SubscribeUsersStatus, userID, userIDs)
}

// UnsubscribeUsersStatus Unsubscribe a user's presence.
func UnsubscribeUsersStatus(callback open_im_sdk_callback.Base, operationID string, userID string, userIDs []string) {
	call(callback, operationID, UserForSDK.User().SubscribeUsersStatus, userID, userIDs)
}

// GetSubscribeUsersStatus Get the online status of subscribers.
func GetSubscribeUsersStatus(callback open_im_sdk_callback.Base, operationID string, userID string) {
	call(callback, operationID, UserForSDK.User().GetSubscribeUsersStatus, userID)
}

// GetUserStatus Get the online status of users.
func GetUserStatus(callback open_im_sdk_callback.Base, operationID string, userID string, userIDs []string) {
	call(callback, operationID, UserForSDK.User().GetUserStatus, userID, userIDs)
}
