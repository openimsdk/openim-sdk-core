package rtc

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/log"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
)

func (s *LiveSignaling) Invite(signalInviteReq string, callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", signalInviteReq)
		req := &api.SignalInviteReq{}
		var signalReq api.SignalReq
		common.JsonUnmarshalCallback(signalInviteReq, req, callback, operationID)
		*signalReq.GetInvite() = *req
		s.handleSignaling(&signalReq, callback, operationID)
		log.NewInfo(operationID, fName, " callback: finished")
	}()
}

func (s *LiveSignaling) Accept(signalAcceptReq string, callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", signalAcceptReq)
		req := &api.SignalAcceptReq{}
		var signalReq api.SignalReq
		utils.JsonStringToStruct(signalAcceptReq, req)
		*signalReq.GetAccept() = *req
		s.handleSignaling(&signalReq, callback, operationID)
		log.NewInfo(operationID, fName, " callback finished")
	}()
}

func (s *LiveSignaling) Reject(signalRejectReq string, callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		log.NewError(operationID, "callback is nil")
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", signalRejectReq)
		req := &api.SignalRejectReq{}
		var signalReq api.SignalReq
		utils.JsonStringToStruct(signalRejectReq, req)
		*signalReq.GetReject() = *req
		s.handleSignaling(&signalReq, callback, operationID)
		log.NewInfo(operationID, fName, " callback finished")
	}()
}

func (s *LiveSignaling) Cancel(signalCancelReq string, callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		log.NewError(operationID, "callback is nil")
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", signalCancelReq)
		req := &api.SignalCancelReq{}
		var signalReq api.SignalReq
		utils.JsonStringToStruct(signalCancelReq, req)
		*signalReq.GetCancel() = *req
		s.handleSignaling(&signalReq, callback, operationID)
		log.NewInfo(operationID, fName, " callback finished")
	}()
}

func (s *LiveSignaling) HungUp(signalHungUpReq string, callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		log.NewError(operationID, "callback is nil")
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", signalHungUpReq)
		req := &api.SignalHungUpReq{}
		var signalReq api.SignalReq
		utils.JsonStringToStruct(signalHungUpReq, req)
		*signalReq.GetHungUp() = *req
		s.handleSignaling(&signalReq, callback, operationID)
		log.NewInfo(operationID, fName, " callback finished")
	}()
}
