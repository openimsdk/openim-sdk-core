package open_im_sdk

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
)

func SetSignalingListener(callback open_im_sdk_callback.OnSignalingListener) {
	if callback == nil || userForSDK == nil {
		log.Error("callback or userForSDK is nil")
		return
	}
	userForSDK.SetSignalingListener(callback)
}

func SetSignalingListenerForService(callback open_im_sdk_callback.OnSignalingListener) {
	if callback == nil || userForSDK == nil {
		log.Error("callback or userForSDK is nil")
		return
	}
	userForSDK.SetSignalingListenerForService(callback)
}

func SetListenerForService(callback open_im_sdk_callback.OnListenerForService) {
	if callback == nil || userForSDK == nil {
		log.Error("callback or userForSDK is nil")
		return
	}
	userForSDK.SetListenerForService(callback)
}

func SignalingInviteInGroup(callback open_im_sdk_callback.Base, operationID string, signalInviteInGroupReq string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Signaling().InviteInGroup(callback, signalInviteInGroupReq, operationID)
}

func SignalingInvite(callback open_im_sdk_callback.Base, operationID string, signalInviteReq string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Signaling().Invite(callback, signalInviteReq, operationID)
}

func SignalingAccept(callback open_im_sdk_callback.Base, operationID string, signalAcceptReq string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Signaling().Accept(callback, signalAcceptReq, operationID)
}

func SignalingReject(callback open_im_sdk_callback.Base, operationID string, signalRejectReq string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Signaling().Reject(callback, signalRejectReq, operationID)
}

func SignalingCancel(callback open_im_sdk_callback.Base, operationID string, signalCancelReq string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Signaling().Cancel(callback, signalCancelReq, operationID)
}

func SignalingHungUp(callback open_im_sdk_callback.Base, operationID string, signalHungUpReq string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Signaling().HungUp(callback, signalHungUpReq, operationID)
}

func SignalingGetRoomByGroupID(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Signaling().SignalingGetRoomByGroupID(callback, groupID, operationID)
}

func SignalingGetTokenByRoomID(callback open_im_sdk_callback.Base, operationID, groupID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Signaling().SignalingGetTokenByRoomID(callback, groupID, operationID)
}
