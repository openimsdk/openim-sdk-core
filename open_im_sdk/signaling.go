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
