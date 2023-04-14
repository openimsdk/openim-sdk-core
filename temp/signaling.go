package open_im_sdk


func SignalingInviteInGroup(callback open_im_sdk_callback.Base, operationID string, signalInviteInGroupReq string) {
	call(callback, operationID, userForSDK.Signaling().SignalingInviteInGroup, signalInviteInGroupReq)
}


func SignalingInvite(callback open_im_sdk_callback.Base, operationID string, signalInviteReq string) {
	call(callback, operationID, userForSDK.Signaling().SignalingInvite, signalInviteReq)
}


func SignalingAccept(callback open_im_sdk_callback.Base, operationID string, signalAcceptReq string) {
	call(callback, operationID, userForSDK.Signaling().SignalingAccept, signalAcceptReq)
}


func SignalingReject(callback open_im_sdk_callback.Base, operationID string, signalRejectReq string) {
	call(callback, operationID, userForSDK.Signaling().SignalingReject, signalRejectReq)
}


func SignalingCancel(callback open_im_sdk_callback.Base, operationID string, signalCancelReq string) {
	call(callback, operationID, userForSDK.Signaling().SignalingCancel, signalCancelReq)
}


func SignalingHungUp(callback open_im_sdk_callback.Base, operationID string, signalHungUpReq string) {
	call(callback, operationID, userForSDK.Signaling().SignalingHungUp, signalHungUpReq)
}


func SignalingGetRoomByGroupID(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, userForSDK.Signaling().SignalingGetRoomByGroupID, groupID)
}

