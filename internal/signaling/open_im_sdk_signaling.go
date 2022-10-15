package signaling

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/log"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
)

func (s *LiveSignaling) setDefaultReq(req *api.InvitationInfo) {
	if req.RoomID == "" {
		req.RoomID = utils.OperationIDGenerator()
	}
	if req.Timeout == 0 {
		req.Timeout = 60 * 60
	}
}

func (s *LiveSignaling) InviteInGroup(callback open_im_sdk_callback.Base, signalInviteInGroupReq string, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	var listenerList []open_im_sdk_callback.OnSignalingListener
	if s.listener != nil {
		listenerList = append(listenerList, s.listener)
		log.Info(operationID, "listenerList ", listenerList, "listener ", s.listener)
	}
	if s.listenerForService != nil {
		listenerList = append(listenerList, s.listenerForService)
		log.Info(operationID, "listenerList ", listenerList, "listenerForService ", s.listenerForService)
	}

	if len(listenerList) == 0 {
		log.Error(operationID, "not set signaling listener ")
		callback.OnError(3004, "not set signaling listener ")
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", signalInviteInGroupReq)
		req := &api.SignalReq_InviteInGroup{InviteInGroup: &api.SignalInviteInGroupReq{Invitation: &api.InvitationInfo{}, OfflinePushInfo: &api.OfflinePushInfo{}}}
		var signalReq api.SignalReq
		common.JsonUnmarshalCallback(signalInviteInGroupReq, req.InviteInGroup, callback, operationID)
		s.setDefaultReq(req.InviteInGroup.Invitation)
		req.InviteInGroup.Invitation.InviterUserID = s.loginUserID
		req.InviteInGroup.OpUserID = s.loginUserID
		req.InviteInGroup.Invitation.InitiateTime = int32(utils.GetCurrentTimestampBySecond())
		signalReq.Payload = req
		req.InviteInGroup.Participant = s.getSelfParticipant(req.InviteInGroup.Invitation.GroupID, callback, operationID)
		s.handleSignaling(&signalReq, callback, operationID)
		log.NewInfo(operationID, fName, " callback: finished")
	}()
}

func (s *LiveSignaling) Invite(callback open_im_sdk_callback.Base, signalInviteReq string, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}

	var listenerList []open_im_sdk_callback.OnSignalingListener
	if s.listener != nil {
		listenerList = append(listenerList, s.listener)
		log.Info(operationID, "listenerList ", listenerList, "listener ", s.listener)
	}
	if s.listenerForService != nil {
		listenerList = append(listenerList, s.listenerForService)
		log.Info(operationID, "listenerList ", listenerList, "listenerForService ", s.listenerForService)
	}

	if len(listenerList) == 0 {
		log.Error(operationID, "listener is nil")
		callback.OnError(3004, "listener is nil")
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", signalInviteReq)
		req := &api.SignalReq_Invite{Invite: &api.SignalInviteReq{Invitation: &api.InvitationInfo{}, OfflinePushInfo: &api.OfflinePushInfo{}}}
		var signalReq api.SignalReq
		common.JsonUnmarshalCallback(signalInviteReq, req.Invite, callback, operationID)
		s.setDefaultReq(req.Invite.Invitation)
		req.Invite.Invitation.InviterUserID = s.loginUserID
		req.Invite.OpUserID = s.loginUserID
		req.Invite.Invitation.InitiateTime = int32(utils.GetCurrentTimestampBySecond())
		signalReq.Payload = req
		req.Invite.Participant = s.getSelfParticipant(req.Invite.Invitation.GroupID, callback, operationID)
		s.handleSignaling(&signalReq, callback, operationID)
		log.NewInfo(operationID, fName, " callback: finished")
	}()
}

func (s *LiveSignaling) Accept(callback open_im_sdk_callback.Base, signalAcceptReq string, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	var listenerList []open_im_sdk_callback.OnSignalingListener
	if s.listener != nil {
		listenerList = append(listenerList, s.listener)
		log.Info(operationID, "listenerList ", listenerList, "listener ", s.listener)
	}
	if s.listenerForService != nil {
		listenerList = append(listenerList, s.listenerForService)
		log.Info(operationID, "listenerList ", listenerList, "listenerForService ", s.listenerForService)
	}

	if len(listenerList) == 0 {
		log.Error(operationID, "not set signaling listener")
		callback.OnError(3004, "not set signaling listener")
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", signalAcceptReq)
		req := &api.SignalReq_Accept{Accept: &api.SignalAcceptReq{Invitation: &api.InvitationInfo{}, OfflinePushInfo: &api.OfflinePushInfo{}}}
		var signalReq api.SignalReq
		common.JsonUnmarshalCallback(signalAcceptReq, req.Accept, callback, operationID)
		s.setDefaultReq(req.Accept.Invitation)
		req.Accept.OpUserID = s.loginUserID
		req.Accept.Invitation.InitiateTime = int32(utils.GetCurrentTimestampBySecond())
		signalReq.Payload = req
		req.Accept.Participant = s.getSelfParticipant(req.Accept.Invitation.GroupID, callback, operationID)
		req.Accept.OpUserPlatformID = s.platformID
		s.handleSignaling(&signalReq, callback, operationID)
		log.NewInfo(operationID, fName, " callback finished")
	}()
}

func (s *LiveSignaling) Reject(callback open_im_sdk_callback.Base, signalRejectReq string, operationID string) {
	if callback == nil {
		log.NewError(operationID, "callback is nil")
		return
	}
	var listenerList []open_im_sdk_callback.OnSignalingListener
	if s.listener != nil {
		listenerList = append(listenerList, s.listener)
		log.Info(operationID, "listenerList ", listenerList, "listener ", s.listener)
	}
	if s.listenerForService != nil {
		listenerList = append(listenerList, s.listenerForService)
		log.Info(operationID, "listenerList ", listenerList, "listenerForService ", s.listenerForService)
	}

	if len(listenerList) == 0 {
		log.Error(operationID, "not set signaling listener")
		callback.OnError(3004, "not set signaling listener")
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", signalRejectReq)
		req := &api.SignalReq_Reject{Reject: &api.SignalRejectReq{Invitation: &api.InvitationInfo{}, OfflinePushInfo: &api.OfflinePushInfo{}}}
		var signalReq api.SignalReq
		common.JsonUnmarshalCallback(signalRejectReq, req.Reject, callback, operationID)
		s.setDefaultReq(req.Reject.Invitation)
		req.Reject.OpUserID = s.loginUserID
		req.Reject.Invitation.InitiateTime = int32(utils.GetCurrentTimestampBySecond())
		signalReq.Payload = req
		req.Reject.OpUserPlatformID = s.platformID
		req.Reject.Participant = s.getSelfParticipant(req.Reject.Invitation.GroupID, callback, operationID)
		s.handleSignaling(&signalReq, callback, operationID)
		log.NewInfo(operationID, fName, " callback finished")
	}()
}

func (s *LiveSignaling) Cancel(callback open_im_sdk_callback.Base, signalCancelReq string, operationID string) {
	if callback == nil {
		log.NewError(operationID, "callback is nil")
	}
	var listenerList []open_im_sdk_callback.OnSignalingListener
	if s.listener != nil {
		listenerList = append(listenerList, s.listener)
		log.Info(operationID, "listenerList ", listenerList, "listener ", s.listener)
	}
	if s.listenerForService != nil {
		listenerList = append(listenerList, s.listenerForService)
		log.Info(operationID, "listenerList ", listenerList, "listenerForService ", s.listenerForService)
	}

	if len(listenerList) == 0 {
		log.Error(operationID, "not set signaling listener")
		callback.OnError(3004, "not set signaling listener")
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", signalCancelReq)
		req := &api.SignalReq_Cancel{Cancel: &api.SignalCancelReq{Invitation: &api.InvitationInfo{}, OfflinePushInfo: &api.OfflinePushInfo{}}}
		var signalReq api.SignalReq
		common.JsonUnmarshalCallback(signalCancelReq, req.Cancel, callback, operationID)
		s.setDefaultReq(req.Cancel.Invitation)
		req.Cancel.OpUserID = s.loginUserID
		req.Cancel.Invitation.InitiateTime = int32(utils.GetCurrentTimestampBySecond())
		signalReq.Payload = req
		s.handleSignaling(&signalReq, callback, operationID)
		log.NewInfo(operationID, fName, " callback finished")
	}()
}

func (s *LiveSignaling) HungUp(callback open_im_sdk_callback.Base, signalHungUpReq string, operationID string) {
	if callback == nil {
		log.NewError(operationID, "callback is nil")
	}
	var listenerList []open_im_sdk_callback.OnSignalingListener
	if s.listener != nil {
		listenerList = append(listenerList, s.listener)
		log.Info(operationID, "listenerList ", listenerList, "listener ", s.listener)
	}
	if s.listenerForService != nil {
		listenerList = append(listenerList, s.listenerForService)
		log.Info(operationID, "listenerList ", listenerList, "listenerForService ", s.listenerForService)
	}

	if len(listenerList) == 0 {
		log.Error(operationID, "not set signaling listener")
		callback.OnError(3004, "not set signaling listener")
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", signalHungUpReq)
		req := &api.SignalReq_HungUp{HungUp: &api.SignalHungUpReq{Invitation: &api.InvitationInfo{}, OfflinePushInfo: &api.OfflinePushInfo{}}}
		var signalReq api.SignalReq
		common.JsonUnmarshalCallback(signalHungUpReq, req.HungUp, callback, operationID)
		s.setDefaultReq(req.HungUp.Invitation)
		req.HungUp.OpUserID = s.loginUserID
		req.HungUp.Invitation.InitiateTime = int32(utils.GetCurrentTimestampBySecond())
		signalReq.Payload = req
		s.handleSignaling(&signalReq, callback, operationID)
		log.NewInfo(operationID, fName, " callback finished")
	}()
}

func (s *LiveSignaling) SignalingGetRoomByGroupID(callback open_im_sdk_callback.Base, groupID, operationID string) {
	if callback == nil {
		log.NewError(operationID, "callback is nil")
	}
	if s.listener == nil {
		log.Error(operationID, "listener is nil")
		callback.OnError(3004, "listener is nil")
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args groupID:", groupID)
		req := &api.SignalReq_GetRoomByGroupID{GetRoomByGroupID: &api.SignalGetRoomByGroupIDReq{
			OpUserID: s.loginUserID,
			GroupID:  groupID,
		}}
		req.GetRoomByGroupID.Participant = s.getSelfParticipant(req.GetRoomByGroupID.GroupID, callback, operationID)
		var signalReq api.SignalReq
		signalReq.Payload = req
		log.NewDebug(operationID, "SignalGetRoomByGroupID", req.GetRoomByGroupID.String(), signalReq.Payload)
		s.handleSignaling(&signalReq, callback, operationID)
		log.NewInfo(operationID, fName, " callback finished")
	}()
}

func (s *LiveSignaling) SignalingGetTokenByRoomID(callback open_im_sdk_callback.Base, groupID, operationID string) {
	if callback == nil {
		log.NewError(operationID, "callback is nil")
	}
	if s.listener == nil {
		log.Error(operationID, "listener is nil")
		callback.OnError(3004, "listener is nil")
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args groupID:", groupID)
		req := &api.SignalReq_GetTokenByRoomID{GetTokenByRoomID: &api.SignalGetTokenByRoomIDReq{
			OpUserID: s.loginUserID,
			RoomID:   groupID,
		}}
		req.GetTokenByRoomID.Participant = s.getSelfParticipant(req.GetTokenByRoomID.RoomID, callback, operationID)
		var signalReq api.SignalReq
		signalReq.Payload = req
		log.NewDebug(operationID, "SignalGetRoomByGroupID", req.GetTokenByRoomID.String(), signalReq.Payload)
		s.handleSignaling(&signalReq, callback, operationID)
		log.NewInfo(operationID, fName, " callback finished")
	}()
}
