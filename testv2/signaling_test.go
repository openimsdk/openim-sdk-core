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

package testv2

// func Test_SignalingInviteInGroup(t *testing.T) {
// 	resp, err := open_im_sdk.UserForSDK.Signaling().SignalingInviteInGroup(ctx, &sdkws.SignalInviteInGroupReq{
// 		Invitation: &sdkws.InvitationInfo{
// 			InviterUserID:     UserID,
// 			InviteeUserIDList: []string{"targetUser"},
// 			CustomData:        "",
// 			GroupID:           "testgroup",
// 			RoomID:            "testgroup",
// 			Timeout:           30,
// 			MediaType:         "video",
// 			PlatformID:        1,
// 			SessionType:       3,
// 		},
// 	})
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	t.Log(resp)
// }

// func Test_SignalingInite(t *testing.T) {
// 	resp, err := open_im_sdk.UserForSDK.Signaling().SignalingInvite(ctx, &sdkws.SignalInviteReq{
// 		Invitation: &sdkws.InvitationInfo{
// 			InviterUserID:     UserID,
// 			InviteeUserIDList: []string{"targetUser"},
// 			CustomData:        "",
// 			GroupID:           "",
// 			RoomID:            "testroomID",
// 			Timeout:           30,
// 			MediaType:         "video",
// 			PlatformID:        1,
// 			SessionType:       1,
// 		},
// 	})
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	t.Log(resp)
// }

// func Test_SignalingAccept(t *testing.T) {
// 	resp, err := open_im_sdk.UserForSDK.Signaling().SignalingAccept(ctx, &sdkws.SignalAcceptReq{
// 		Invitation: &sdkws.InvitationInfo{
// 			InviterUserID:     UserID,
// 			InviteeUserIDList: []string{"targetUser"},
// 			CustomData:        "",
// 			GroupID:           "",
// 			RoomID:            "testroomID",
// 			Timeout:           30,
// 			MediaType:         "video",
// 			PlatformID:        1,
// 			SessionType:       1,
// 		},
// 	})
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	t.Log(resp)
// }

// func Test_SignalingReject(t *testing.T) {
// 	err := open_im_sdk.UserForSDK.Signaling().SignalingReject(ctx, &sdkws.SignalRejectReq{
// 		Invitation: &sdkws.InvitationInfo{
// 			InviterUserID:     UserID,
// 			InviteeUserIDList: []string{"targetUser"},
// 			CustomData:        "",
// 			GroupID:           "",
// 			RoomID:            "testroomID",
// 			Timeout:           30,
// 			MediaType:         "video",
// 			PlatformID:        1,
// 			SessionType:       1,
// 		},
// 	})
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func Test_SignalingCancel(t *testing.T) {
// 	err := open_im_sdk.UserForSDK.Signaling().SignalingCancel(ctx, &sdkws.SignalCancelReq{
// 		Invitation: &sdkws.InvitationInfo{
// 			InviterUserID:     UserID,
// 			InviteeUserIDList: []string{"targetUser"},
// 			CustomData:        "",
// 			GroupID:           "",
// 			RoomID:            "testroomID",
// 			Timeout:           30,
// 			MediaType:         "video",
// 			PlatformID:        1,
// 			SessionType:       1,
// 		},
// 	})
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func Test_SignalingHungUp(t *testing.T) {
// 	err := open_im_sdk.UserForSDK.Signaling().SignalingHungUp(ctx, &sdkws.SignalHungUpReq{
// 		Invitation: &sdkws.InvitationInfo{
// 			InviterUserID:     UserID,
// 			InviteeUserIDList: []string{"targetUser"},
// 			CustomData:        "",
// 			GroupID:           "",
// 			RoomID:            "testroomID",
// 			Timeout:           30,
// 			MediaType:         "video",
// 			PlatformID:        1,
// 			SessionType:       1,
// 		},
// 	})
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func Test_SignalingGetRoomByGroupID(t *testing.T) {
// 	resp, err := open_im_sdk.UserForSDK.Signaling().SignalingGetRoomByGroupID(ctx, "testgroupID")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	t.Log(resp)
// }

// func Test_SignalingGetTokenByRoomID(t *testing.T) {
// 	resp, err := open_im_sdk.UserForSDK.Signaling().SignalingGetTokenByRoomID(ctx, "testroomID")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	t.Log(resp)
// }
