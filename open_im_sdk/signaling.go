// Copyright © 2023 OpenIM SDK.
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

// SignalingInviteInGroup sends an invitation to join a signaling group.
func SignalingInviteInGroup(callback open_im_sdk_callback.Base, operationID string, signalInviteInGroupReq string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingInviteInGroup, signalInviteInGroupReq)
}

// SignalingInvite sends a signaling invitation to a user or a group of users.
func SignalingInvite(callback open_im_sdk_callback.Base, operationID string, signalInviteReq string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingInvite, signalInviteReq)
}

// SignalingAccept accepts a signaling invitation.
func SignalingAccept(callback open_im_sdk_callback.Base, operationID string, signalAcceptReq string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingAccept, signalAcceptReq)
}

// SignalingReject rejects a signaling invitation.
func SignalingReject(callback open_im_sdk_callback.Base, operationID string, signalRejectReq string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingReject, signalRejectReq)
}

// SignalingCancel cancels a signaling invitation or a signaling group.
func SignalingCancel(callback open_im_sdk_callback.Base, operationID string, signalCancelReq string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingCancel, signalCancelReq)
}

// SignalingHungUp terminates a signaling connection.
func SignalingHungUp(callback open_im_sdk_callback.Base, operationID string, signalHungUpReq string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingHungUp, signalHungUpReq)
}

// SignalingGetRoomByGroupID gets the signaling room information by group ID.
func SignalingGetRoomByGroupID(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingGetRoomByGroupID, groupID)
}

// SignalingGetTokenByRoomID gets the signaling token by room ID.
func SignalingGetTokenByRoomID(callback open_im_sdk_callback.Base, operationID string, roomID string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingGetTokenByRoomID, roomID)
}
