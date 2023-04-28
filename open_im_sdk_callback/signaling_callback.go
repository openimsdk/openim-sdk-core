// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package open_im_sdk_callback

type OnSignalingListener interface {
	OnReceiveNewInvitation(receiveNewInvitationCallback string)

	OnInviteeAccepted(inviteeAcceptedCallback string)

	OnInviteeAcceptedByOtherDevice(inviteeAcceptedCallback string)

	OnInviteeRejected(inviteeRejectedCallback string)

	OnInviteeRejectedByOtherDevice(inviteeRejectedCallback string)
	//
	OnInvitationCancelled(invitationCancelledCallback string)
	//
	OnInvitationTimeout(invitationTimeoutCallback string)
	//
	OnHangUp(hangUpCallback string)

	OnRoomParticipantConnected(onRoomParticipantConnectedCallback string)

	OnRoomParticipantDisconnected(onRoomParticipantDisconnectedCallback string)
}
