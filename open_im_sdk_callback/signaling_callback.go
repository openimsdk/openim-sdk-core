package open_im_sdk_callback

type OnSignalingListener interface {
	OnReceiveNewInvitation( )
	OnInviteeAccepted()

	OnInviteeRejected()
	//
	OnInvitationCancelled()
	//
	OnInvitationTimeout()
}
