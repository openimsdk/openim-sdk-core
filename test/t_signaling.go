package test

import (
	"fmt"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/log"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
)

type testSignalingListener struct {
}

func (s *testSignalingListener) OnHangUp(hangUpCallback string) {
	//panic("implement me")
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

func (s *testSignalingListener) OnRoomParticipantConnected(onRoomChangeCallback string) {
	log.Info("", utils.GetSelfFuncName(), "listener ", onRoomChangeCallback)
}

func (s *testSignalingListener) OnRoomParticipantDisconnected(onRoomChangeCallback string) {
	log.Info("", utils.GetSelfFuncName(), "listener ", onRoomChangeCallback)
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
	req.InviteeUserIDList = []string{"3495023045"}
	req.MediaType = "video"
	req.RoomID = "1826384574"
	req.GroupID = "1826384574"
	req.SessionType = 2
	return req
}

func DoTestInvite(userID string) {
	t := testSingaling{baseCallback{OperationID: utils.OperationIDGenerator(), callName: utils.GetSelfFuncName()}}
	req := &api.SignalInviteReq{}
	req.OpUserID = userID
	req.Invitation = SetTestInviteInfo()
	req.Invitation.GroupID = ""
	req.Invitation.SessionType = 1
	req.Invitation.PlatformID = 1
	req.Invitation.Timeout = 30
	req.Invitation.MediaType = "video"
	req.Invitation.InviteeUserIDList = []string{"17726378428"}
	s := utils.StructToJsonString(req)
	fmt.Println(utils.GetSelfFuncName(), "input: ", s, t.OperationID)
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

func DoTestSignalGetRoomByGroupID(groupID string) {
	t := testSingaling{baseCallback{OperationID: utils.OperationIDGenerator(), callName: utils.GetSelfFuncName()}}
	open_im_sdk.SignalingGetRoomByGroupID(t, t.OperationID, groupID)
}

func DoTestSignalGetTokenByRoomID(roomID string) {
	t := testSingaling{baseCallback{OperationID: utils.OperationIDGenerator(), callName: utils.GetSelfFuncName()}}
	open_im_sdk.SignalingGetTokenByRoomID(t, t.OperationID, roomID)
}
