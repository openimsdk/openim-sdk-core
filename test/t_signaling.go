package test

import (
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/log"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
)

type testSignalingListener struct {
}

func (s *testSignalingListener) OnHangUp(hangUpCallback string) {
	panic("implement me")
}

func (s *testSignalingListener) OnReceiveNewInvitation(receiveNewInvitationCallback string) {
	log.Info("", utils.GetSelfFuncName(), "listener ", receiveNewInvitationCallback)
}

func (s *testSignalingListener) OnInviteeAccepted(inviteeAcceptedCallback string) {
	log.Info("", utils.GetSelfFuncName(), "listener ", inviteeAcceptedCallback)
}

func (s *testSignalingListener) OnInviteeRejected(inviteeRejectedCallback string) {
	log.Info("", utils.GetSelfFuncName(), "listener ", inviteeRejectedCallback)
}

//
func (s *testSignalingListener) OnInvitationCancelled(invitationCancelledCallback string) {
	log.Info("", utils.GetSelfFuncName(), "listener ", invitationCancelledCallback)
}

//
func (s *testSignalingListener) OnInvitationTimeout(invitationTimeoutCallback string) {
	log.Info("", utils.GetSelfFuncName(), "listener ", invitationTimeoutCallback)
}

func (s *testSignalingListener) OnInviteeAcceptedByOtherDevice(inviteeAcceptedCallback string) {
	log.Info("", utils.GetSelfFuncName(), "listener ", inviteeAcceptedCallback)
}

func (s *testSignalingListener) OnInviteeRejectedByOtherDevice(inviteeRejectedCallback string) {
	log.Info("", utils.GetSelfFuncName(), "listener ", inviteeRejectedCallback)
}

type testSingaling struct {
	baseCallback
}

var TestRoomID = "room_id_111"

func DoTestInviteInGroup() {
	t := testSingaling{baseCallback{OperationID: utils.OperationIDGenerator(), callName: utils.GetSelfFuncName()}}
	req := &api.SignalInviteInGroupReq{}
	req.Invitation = SetTestInviteInfo()
	s := utils.StructToJsonString(req)
	log.Info(t.OperationID, utils.GetSelfFuncName(), "input: ", s)
	open_im_sdk.SignalingInviteInGroup(t, t.OperationID, s)
}

func SetTestInviteInfo() *api.InvitationInfo {
	req := &api.InvitationInfo{}
	req.Timeout = 1000
	req.InviteeUserIDList = append(req.InviteeUserIDList, MemberUserID)
	req.MediaType = "video"
	req.RoomID = TestRoomID
	req.GroupID = TestgroupID
	req.SessionType = 2
	return req
}

func DoTestInvite() {
	t := testSingaling{baseCallback{OperationID: utils.OperationIDGenerator(), callName: utils.GetSelfFuncName()}}
	req := &api.SignalInviteReq{}
	req.Invitation = SetTestInviteInfo()
	req.Invitation.GroupID = ""
	req.Invitation.SessionType = 1
	req.Invitation.PlatformID = 1
	req.Invitation.Timeout = 10
	s := utils.StructToJsonString(req)
	log.Info(t.OperationID, utils.GetSelfFuncName(), "input: ", s)
	open_im_sdk.SignalingInvite(t, t.OperationID, s)
}

func DoTestAccept() {
	t := testSingaling{baseCallback{OperationID: utils.OperationIDGenerator(), callName: utils.GetSelfFuncName()}}
	req := &api.SignalAcceptReq{Invitation: &api.InvitationInfo{}, OpUserID: "18349115126"}
	req.Invitation = SetTestInviteInfo()
	req.Invitation.InviterUserID = "18666662412"
	s := utils.StructToJsonString(req)
	log.Info(t.OperationID, utils.GetSelfFuncName(), "input: ", s, req.String())
	open_im_sdk.SignalingAccept(t, t.OperationID, s)
}

func DoTestReject() {
	t := testSingaling{baseCallback{OperationID: utils.OperationIDGenerator(), callName: utils.GetSelfFuncName()}}
	req := &api.SignalRejectReq{Invitation: &api.InvitationInfo{}, OpUserID: "18349115126"}
	req.Invitation = SetTestInviteInfo()
	req.Invitation.InviterUserID = "18666662412"
	s := utils.StructToJsonString(req)
	log.Info(t.OperationID, utils.GetSelfFuncName(), "input: ", s)
	open_im_sdk.SignalingReject(t, t.OperationID, s)
}

func DoTestCancel() {
	t := testSingaling{baseCallback{OperationID: utils.OperationIDGenerator(), callName: utils.GetSelfFuncName()}}
	req := &api.SignalCancelReq{Invitation: &api.InvitationInfo{}}

	req.Invitation = SetTestInviteInfo()
	req.Invitation.GroupID = ""
	req.Invitation.SessionType = 1
	req.Invitation.PlatformID = 1
	req.Invitation.Timeout = 10
	req.Invitation.InviterUserID = "18666662412"

	req.OpUserID = "18666662412"
	s := utils.StructToJsonString(req)
	log.Info(t.OperationID, utils.GetSelfFuncName(), "input: ", s)
	open_im_sdk.SignalingCancel(t, t.OperationID, s)
}

func DoTestHungUp() {
	t := testSingaling{baseCallback{OperationID: utils.OperationIDGenerator(), callName: utils.GetSelfFuncName()}}
	req := &api.SignalHungUpReq{Invitation: &api.InvitationInfo{}}
	req.Invitation = SetTestInviteInfo()
	s := utils.StructToJsonString(req)
	log.Info(t.OperationID, utils.GetSelfFuncName(), "input: ", s)
	open_im_sdk.SignalingHungUp(t, t.OperationID, s)
}
