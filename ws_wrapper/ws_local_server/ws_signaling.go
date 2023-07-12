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

package ws_local_server

import (
	"open_im_sdk/open_im_sdk"
)

type SignalingCallback struct {
	uid string
}

func (s *SignalingCallback) OnReceiveNewInvitation(receiveNewInvitation string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", receiveNewInvitation, "0"}, s.uid)
}

func (s *SignalingCallback) OnInviteeAccepted(inviteeAccepted string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", inviteeAccepted, "0"}, s.uid)
}

func (s *SignalingCallback) OnInviteeRejected(inviteeRejected string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", inviteeRejected, "0"}, s.uid)
}

func (s *SignalingCallback) OnInvitationCancelled(invitationCancelled string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", invitationCancelled, "0"}, s.uid)
}

func (s *SignalingCallback) OnInvitationTimeout(invitationTimeout string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", invitationTimeout, "0"}, s.uid)
}

func (s *SignalingCallback) OnInviteeAcceptedByOtherDevice(inviteeAcceptedCallback string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", inviteeAcceptedCallback, "0"}, s.uid)
}

func (s *SignalingCallback) OnInviteeRejectedByOtherDevice(inviteeRejectedCallback string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", inviteeRejectedCallback, "0"}, s.uid)
}

func (s *SignalingCallback) OnHangUp(hangUpCallback string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", hangUpCallback, "0"}, s.uid)
}

func (s *SignalingCallback) OnRoomParticipantConnected(onRoomParticipantConnected string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", onRoomParticipantConnected, "0"}, s.uid)
}

func (s *SignalingCallback) OnRoomParticipantDisconnected(onRoomParticipantDisconnected string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", onRoomParticipantDisconnected, "0"}, s.uid)
}

func (wsRouter *WsFuncRouter) SetSignalingListener() {
	var sr SignalingCallback
	sr.uid = wsRouter.uId
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.SetSignalingListener(&sr)
}

func (wsRouter *WsFuncRouter) SignalingInvite(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Signaling().Invite(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, input, operationID)
}

func (wsRouter *WsFuncRouter) SignalingInviteInGroup(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Signaling().InviteInGroup(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, input, operationID)
}

func (wsRouter *WsFuncRouter) SignalingAccept(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Signaling().Accept(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, input, operationID)
}

func (wsRouter *WsFuncRouter) SignalingReject(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Signaling().Reject(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, input, operationID)
}

func (wsRouter *WsFuncRouter) SignalingCancel(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Signaling().Cancel(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, input, operationID)
}

func (wsRouter *WsFuncRouter) SignalingHungUp(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Signaling().HungUp(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, input, operationID)
}
