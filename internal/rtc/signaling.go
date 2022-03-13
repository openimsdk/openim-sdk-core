package rtc

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	api "open_im_sdk/pkg/server_api_params"
)

type Signaling interface {

	//invitee 被邀请者
	Invite(inviteeUserID, customData string, offlinePushInfo *api.OfflinePushInfo, timeout uint32, callback open_im_sdk_callback.Base, operationID string) error

	InviteInGroup(groupID string, inviteeUserIDList []string, customData string, offlinePushInfo *api.OfflinePushInfo, timeout uint32, callback open_im_sdk_callback.Base, operationID string) error

	Cancel(inviteeUserID, customData string, callback open_im_sdk_callback.Base, operationID string) error

	Accept(inviteUserID, customData string, callback open_im_sdk_callback.Base, operationID string) error

	Reject(inviteUserID, customData string, callback open_im_sdk_callback.Base, operationID string) error

	HungUp(peerUserID, customData string, callback open_im_sdk_callback.Base, operationID string) error

	SetListener(listener open_im_sdk_callback.OnSignalingListener, operationID string) error

	DoNotification(msg *api.MsgData, conversationCh chan common.Cmd2Value)
}
