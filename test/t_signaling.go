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

package test

import (
	"open_im_sdk/pkg/utils"

	"github.com/OpenIMSDK/tools/log"
	"golang.org/x/net/context"
)

type testSignalingListener struct {
	ctx context.Context
}

func (s *testSignalingListener) OnHangUp(hangUpCallback string) {
	log.ZInfo(s.ctx, utils.GetSelfFuncName(), "OnHangUp ", hangUpCallback)
}

func (s *testSignalingListener) OnReceiveNewInvitation(receiveNewInvitationCallback string) {
	log.ZInfo(s.ctx, utils.GetSelfFuncName(), "OnReceiveNewInvitation ", receiveNewInvitationCallback)
}

func (s *testSignalingListener) OnInviteeAccepted(inviteeAcceptedCallback string) {
	log.ZInfo(s.ctx, utils.GetSelfFuncName(), "OnInviteeAccepted ", inviteeAcceptedCallback)
}

func (s *testSignalingListener) OnInviteeRejected(inviteeRejectedCallback string) {
	log.ZInfo(s.ctx, utils.GetSelfFuncName(), "OnInviteeRejected ", inviteeRejectedCallback)
}

func (s *testSignalingListener) OnInvitationCancelled(invitationCancelledCallback string) {
	log.ZInfo(s.ctx, utils.GetSelfFuncName(), "OnInvitationCancelled ", invitationCancelledCallback)
}

func (s *testSignalingListener) OnInvitationTimeout(invitationTimeoutCallback string) {
	log.ZInfo(s.ctx, utils.GetSelfFuncName(), "OnInvitationTimeout ", invitationTimeoutCallback)
}

func (s *testSignalingListener) OnInviteeAcceptedByOtherDevice(inviteeAcceptedCallback string) {
	log.ZInfo(s.ctx, utils.GetSelfFuncName(), "OnInviteeAcceptedByOtherDevice ", inviteeAcceptedCallback)
}

func (s *testSignalingListener) OnInviteeRejectedByOtherDevice(inviteeRejectedCallback string) {
	log.ZInfo(s.ctx, utils.GetSelfFuncName(), "OnInviteeRejectedByOtherDevice ", inviteeRejectedCallback)
}

func (s *testSignalingListener) OnRoomParticipantConnected(onRoomChangeCallback string) {
	log.ZInfo(s.ctx, utils.GetSelfFuncName(), "onRoomChangeCallback", onRoomChangeCallback)
}

func (s *testSignalingListener) OnRoomParticipantDisconnected(onRoomChangeCallback string) {
	log.ZInfo(s.ctx, utils.GetSelfFuncName(), "onRoomChangeCallback", onRoomChangeCallback)
}

//type testSingaling struct {
//	baseCallback
//}
//
//funcation DoTestInviteInGroup() {
//	t := testSingaling{baseCallback{OperationID: utils.OperationIDGenerator(), callName: utils.GetSelfFuncName()}}
//	req := &sdkws.SignalInviteInGroupReq{}
//	req.Invitation = SetTestInviteInfo()
//	s := utils.StructToJsonString(req)
//	// log.Info(t.OperationID, utils.GetSelfFuncName(), "input: ", s)
//	open_im_sdk.SignalingInviteInGroup(t, t.OperationID, s)
//}
//
//funcation SetTestInviteInfo() *sdkws.InvitationInfo {
//	req := &sdkws.InvitationInfo{}
//	req.Timeout = 1000
//	req.InviteeUserIDList = []string{"3495023045"}
//	req.MediaType = "video"
//	req.RoomID = "1826384574"
//	req.GroupID = "1826384574"
//	req.SessionType = 2
//	return req
//}
//
//funcation DoTestInvite(userID string) {
//	t := testSingaling{baseCallback{OperationID: utils.OperationIDGenerator(), callName: utils.GetSelfFuncName()}}
//	req := &sdkws.SignalInviteReq{}
//	req.OpUserID = userID
//	req.Invitation = SetTestInviteInfo()
//	req.Invitation.GroupID = ""
//	req.Invitation.SessionType = 1
//	req.Invitation.PlatformID = 1
//	req.Invitation.Timeout = 30
//	req.Invitation.MediaType = "video"
//	req.Invitation.InviteeUserIDList = []string{"17726378428"}
//	s := utils.StructToJsonString(req)
//	fmt.Println(utils.GetSelfFuncName(), "input: ", s, t.OperationID)
//	open_im_sdk.SignalingInvite(t, t.OperationID, s)
//}
//
//funcation DoTestAccept() {
//	t := testSingaling{baseCallback{OperationID: utils.OperationIDGenerator(), callName: utils.GetSelfFuncName()}}
//	req := &sdkws.SignalAcceptReq{Invitation: &sdkws.InvitationInfo{}, OpUserID: "18349115126"}
//	req.Invitation = SetTestInviteInfo()
//	req.Invitation.InviterUserID = "18666662412"
//	s := utils.StructToJsonString(req)
//	// log.Info(t.OperationID, utils.GetSelfFuncName(), "input: ", s, req.String())
//	open_im_sdk.SignalingAccept(t, t.OperationID, s)
//}
//
//funcation DoTestReject() {
//	t := testSingaling{baseCallback{OperationID: utils.OperationIDGenerator(), callName: utils.GetSelfFuncName()}}
//	req := &sdkws.SignalRejectReq{Invitation: &sdkws.InvitationInfo{}, OpUserID: "18349115126"}
//	req.Invitation = SetTestInviteInfo()
//	req.Invitation.InviterUserID = "18666662412"
//	s := utils.StructToJsonString(req)
//	// log.Info(t.OperationID, utils.GetSelfFuncName(), "input: ", s)
//	open_im_sdk.SignalingReject(t, t.OperationID, s)
//}
//
//funcation DoTestCancel() {
//	t := testSingaling{baseCallback{OperationID: utils.OperationIDGenerator(), callName: utils.GetSelfFuncName()}}
//	req := &sdkws.SignalCancelReq{Invitation: &sdkws.InvitationInfo{}}
//	req.Invitation = SetTestInviteInfo()
//	req.Invitation.GroupID = ""
//	req.Invitation.SessionType = 1
//	req.Invitation.PlatformID = 1
//	req.Invitation.Timeout = 10
//	req.Invitation.InviterUserID = "18666662412"
//	req.OpUserID = "18666662412"
//	s := utils.StructToJsonString(req)
//	// log.Info(t.OperationID, utils.GetSelfFuncName(), "input: ", s)
//	open_im_sdk.SignalingCancel(t, t.OperationID, s)
//}
//
//funcation DoTestHungUp() {
//	t := testSingaling{baseCallback{OperationID: utils.OperationIDGenerator(), callName: utils.GetSelfFuncName()}}
//	req := &sdkws.SignalHungUpReq{Invitation: &sdkws.InvitationInfo{}}
//	req.Invitation = SetTestInviteInfo()
//	s := utils.StructToJsonString(req)
//	// log.Info(t.OperationID, utils.GetSelfFuncName(), "input: ", s)
//	open_im_sdk.SignalingHungUp(t, t.OperationID, s)
//}
//
//funcation DoTestSignalGetRoomByGroupID(groupID string) {
//	t := testSingaling{baseCallback{OperationID: utils.OperationIDGenerator(), callName: utils.GetSelfFuncName()}}
//	open_im_sdk.SignalingGetRoomByGroupID(t, t.OperationID, groupID)
//}
//
//funcation DoTestSignalGetTokenByRoomID(roomID string) {
//	t := testSingaling{baseCallback{OperationID: utils.OperationIDGenerator(), callName: utils.GetSelfFuncName()}}
//	open_im_sdk.SignalingGetTokenByRoomID(t, t.OperationID, roomID)
//}
