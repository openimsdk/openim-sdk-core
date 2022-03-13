package rtc

import (
	"errors"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/sdk_params_callback"
	api "open_im_sdk/pkg/server_api_params"
)

type LiveSignaling struct {
	*ws.Ws
	listener open_im_sdk_callback.OnSignalingListener
	loginUserID          string
}


//invitee 被邀请者
func (s *LiveSignaling) invite(req *api.SignalInviteReq, callback open_im_sdk_callback.Base, operationID string) sdk_params_callback.InviteCallback{
	var signalReq api.SignalReq
	*signalReq.GetInvite() = *req
	resp, err := s.SendSignalingReqWaitResp(&signalReq, 0, operationID)
	common.CheckAnyErrCallback(callback, 3001, err, operationID)
	switch payload := resp.Payload.(type){
	case *api.SignalResp_Invite:
		go s.do(req.Invitation.InviterUserID, req.Invitation.InviteeUserIDList[0], "invite", 100, operationID)
		return sdk_params_callback.InviteCallback(payload.Invite)
	default:
		log.Error(operationID, "resp payload type failed ", payload)
		common.CheckAnyErrCallback(callback, 3002, errors.New("resp payload type failed"), operationID)
		return nil
	}
}

func (s *LiveSignaling) do(inviterUserID, inviteeUserID, event string, timeout int, operationID string){
	req, err := s.SignalingWaitPush(inviterUserID, inviteeUserID, "invite", 100, operationID)
	s.f(req)
}

func (s *LiveSignaling) dopush(req *api.SignalReq){
	//payload.Accept
	switch payload := req.Payload.(type) {
	case *api.SignalReq_Invite:
		s.listener.OnReceiveNewInvitation()
	case *api.SignalReq_Accept:
		s.listener.OnInviteeAccepted()
	case *api.SignalReq_Reject:
		s.listener.OnInviteeRejected()
	case *api.SignalReq_Cancel:
		s.listener.OnInvitationCancelled()
	default:
		log.Error("", "payload type failed ")
	}

}



func (s *LiveSignaling) inviteInGroup(groupID string, inviteeUserIDList[] string, customData string, offlinePushInfo *api.OfflinePushInfo, timeout uint32 , callback open_im_sdk_callback.Base, operationID string) sdk_params_callback.InviteInGroupCallback{
	return nil
}

func (s *LiveSignaling) cancel(inviteeUserID, customData string, callback open_im_sdk_callback.Base, operationID string) sdk_params_callback.CancelCallback{
	return nil
}

func (s *LiveSignaling) accept(inviteUserID, customData string, callback open_im_sdk_callback.Base, operationID string) sdk_params_callback.AcceptCallback{
	return nil
}

func (s *LiveSignaling) hungUp(peerUserID, customData string, callback open_im_sdk_callback.Base, operationID string) sdk_params_callback.HungUpCallback{
	return nil
}

func (s *LiveSignaling) reject(inviteUserID, customData string, callback open_im_sdk_callback.Base, operationID string) sdk_params_callback.RejectCallback{
	return nil
}

func (s *LiveSignaling) SetListener(listener open_im_sdk_callback.OnSignalingListener, operationID string) {
	return nil
}

