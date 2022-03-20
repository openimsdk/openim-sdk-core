package advanced_interface

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	api "open_im_sdk/pkg/server_api_params"
)

type Signaling interface {
	Invite(callback open_im_sdk_callback.Base, signalInviteReq string, operationID string)

	InviteInGroup(callback open_im_sdk_callback.Base, signalInviteInGroupReq string, operationID string)

	Cancel(callback open_im_sdk_callback.Base, signalCancelReq string, operationID string)

	Accept(callback open_im_sdk_callback.Base, signalAcceptReq string, operationID string)

	Reject(callback open_im_sdk_callback.Base, signalRejectReq string, operationID string)

	HungUp(callback open_im_sdk_callback.Base, signalHungUpReq string, operationID string)

	SetListener(listener open_im_sdk_callback.OnSignalingListener, operationID string)

	DoNotification(msg *api.MsgData, conversationCh chan common.Cmd2Value, operationID string)
}
