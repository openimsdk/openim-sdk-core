package open_im_sdk_callback

type OnSignalingListener interface {
	OnReceiveNewInvitation(receiveNewInvitationCallback string)

	OnInviteeAccepted(inviteeAcceptedCallback string)

	OnInviteeRejected(inviteeRejectedCallback string)
	//
	OnInvitationCancelled(invitationCancelledCallback string)
	//
	OnInvitationTimeout(invitationTimeoutCallback string)
}
