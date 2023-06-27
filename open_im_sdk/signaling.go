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

func SignalingInviteInGroup(callback open_im_sdk_callback.Base, operationID string, signalInviteInGroupReq string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingInviteInGroup, signalInviteInGroupReq)
}

func SignalingInvite(callback open_im_sdk_callback.Base, operationID string, signalInviteReq string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingInvite, signalInviteReq)
}

func SignalingAccept(callback open_im_sdk_callback.Base, operationID string, signalAcceptReq string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingAccept, signalAcceptReq)
}

func SignalingReject(callback open_im_sdk_callback.Base, operationID string, signalRejectReq string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingReject, signalRejectReq)
}

func SignalingCancel(callback open_im_sdk_callback.Base, operationID string, signalCancelReq string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingCancel, signalCancelReq)
}

func SignalingHungUp(callback open_im_sdk_callback.Base, operationID string, signalHungUpReq string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingHungUp, signalHungUpReq)
}

func SignalingGetRoomByGroupID(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingGetRoomByGroupID, groupID)
}

func SignalingGetTokenByRoomID(callback open_im_sdk_callback.Base, operationID string, roomID string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingGetTokenByRoomID, roomID)
}

func GetSignalingInvitationInfoStartApp(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Signaling().GetSignalingInvitationInfoStartApp)
}

func SignalingCreateMeeting(callback open_im_sdk_callback.Base, operationID string, signalingCreateMeetingReq string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingCreateMeeting, signalingCreateMeetingReq)
}

func SignalingJoinMeeting(callback open_im_sdk_callback.Base, operationID string, signalingJoinMeetingReq string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingJoinMeeting, signalingJoinMeetingReq)
}

func SignalingUpdateMeetingInfo(callback open_im_sdk_callback.Base, operationID string, signalingUpdateMeetingInfoReq string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingUpdateMeetingInfo, signalingUpdateMeetingInfoReq)
}

func SignalingCloseRoom(callback open_im_sdk_callback.Base, operationID string, roomID string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingCloseRoom, roomID)
}

func SignalingGetMeetings(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingGetMeetings)
}

func SignalingOperateStream(callback open_im_sdk_callback.Base, operationID string, streamType, roomID, userID string, mute, muteAll bool) {
	call(callback, operationID, UserForSDK.Signaling().SignalingOperateStream)
}

func SignalingSendCustomSignal(callback open_im_sdk_callback.Base, operationID string, customInfo, roomID string) {
	call(callback, operationID, UserForSDK.Signaling().SignalingSendCustomSignal, customInfo, roomID)
}
