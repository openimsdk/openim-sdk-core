package rtc

import (
	"errors"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/sdk_params_callback"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"strings"
)

type LiveSignaling struct {
	*ws.Ws
	listener    open_im_sdk_callback.OnSignalingListener
	loginUserID string
}

func NewLiveSignaling(ws *ws.Ws, listener open_im_sdk_callback.OnSignalingListener, loginUserID string) *LiveSignaling {
	if ws == nil || listener == nil {
		log.Error("", "ws or listener is nil")
	}
	return &LiveSignaling{Ws: ws, listener: listener, loginUserID: loginUserID}
}

func (s *LiveSignaling) waitPush(req *api.SignalReq, operationID string) {
	var invt *api.InvitationInfo
	switch payload := req.Payload.(type) {
	case *api.SignalReq_Invite:
		invt = payload.Invite.Invitation
	case *api.SignalReq_InviteInGroup:
		invt = payload.InviteInGroup.Invitation
	}

	for _, v := range invt.InviteeUserIDList {
		go func() {
			push, err := s.SignalingWaitPush(invt.InviterUserID, v, invt.RoomID, invt.Timeout, operationID)
			if err != nil {
				if strings.Contains(err.Error(), "timeout") {
					log.Error(operationID, "wait push timeout ", err.Error(), invt.InviterUserID, v, invt.RoomID, invt.Timeout)

				} else {
					log.Error(operationID, "other failed ", err.Error(), invt.InviterUserID, v, invt.RoomID, invt.Timeout)
				}
				return
			}
			log.Info(operationID, "SignalingWaitPush ", push.String(), invt.InviterUserID, v, invt.RoomID, invt.Timeout)
			s.doSignalPush(push, operationID)
		}()
	}
}

func (s *LiveSignaling) doSignalPush(req *api.SignalReq, operationID string) {
	switch payload := req.Payload.(type) {
	case *api.SignalReq_Invite:
		log.Info(operationID, "recv signal push ", payload.Invite.String())
		s.listener.OnReceiveNewInvitation(utils.StructToJsonString(payload.Invite))
	case *api.SignalReq_Accept:
		log.Info(operationID, "recv signal push ", payload.Accept.String())
		s.listener.OnInviteeAccepted(utils.StructToJsonString(payload.Accept))
	case *api.SignalReq_Reject:
		log.Info(operationID, "recv signal push ", payload.Reject.String())
		s.listener.OnInviteeRejected(utils.StructToJsonString(payload.Reject))
	case *api.SignalReq_Cancel:
		log.Info(operationID, "recv signal push ", payload.Cancel.String())
		s.listener.OnInvitationCancelled(utils.StructToJsonString(payload.Cancel))
	default:
		log.Error(operationID, "payload type failed ")
	}
}

func (s *LiveSignaling) SetListener(listener open_im_sdk_callback.OnSignalingListener, operationID string) {
	s.listener = listener
}

func (s *LiveSignaling) handleSignaling(req *api.SignalReq, callback open_im_sdk_callback.Base, operationID string) {
	log.Info(operationID, utils.GetSelfFuncName(), "args ", req.String())
	resp, err := s.SendSignalingReqWaitResp(req, operationID)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "SendSignalingReqWaitResp error", err.Error())
		common.CheckAnyErrCallback(callback, 3003, errors.New("timeout"), operationID)
	}
	common.CheckAnyErrCallback(callback, 3001, err, operationID)
	switch payload := resp.Payload.(type) {
	case *api.SignalResp_Accept:
		log.Info(operationID, "signaling response ", payload.Accept.String())
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.AcceptCallback(payload.Accept)))
	case *api.SignalResp_Reject:
		log.Info(operationID, "signaling response ", payload.Reject.String())
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.RejectCallback(payload.Reject)))
	case *api.SignalResp_HungUp:
		log.Info(operationID, "signaling response ", payload.HungUp.String())
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.HungUpCallback(payload.HungUp)))
	case *api.SignalResp_Cancel:
		log.Info(operationID, "signaling response ", payload.Cancel.String())
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.CancelCallback(payload.Cancel)))
	case *api.SignalResp_Invite:
		log.Info(operationID, "signaling response ", payload.Invite.String())
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.InviteCallback(payload.Invite)))
	default:
		log.Error(operationID, "resp payload type failed ", payload)
		common.CheckAnyErrCallback(callback, 3002, errors.New("resp payload type failed"), operationID)
	}
	switch req.Payload.(type) {
	case *api.SignalReq_Invite:
		log.Info(operationID, "wait push ", req.String())
		s.waitPush(req, operationID)
	case *api.SignalReq_InviteInGroup:
		log.Info(operationID, "wait push ", req.String())
		s.waitPush(req, operationID)
	}
}
