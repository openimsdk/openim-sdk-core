// Copyright © 2023 OpenIM SDK.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package signaling

import (
	"context"
	"open_im_sdk/pkg/utils"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

func (s *LiveSignaling) SignalingInviteInGroup(ctx context.Context, signalInviteInGroupReq *sdkws.SignalInviteInGroupReq) (*sdkws.SignalInviteInGroupReply, error) {
	s.setDefaultReq(signalInviteInGroupReq.Invitation)
	signalInviteInGroupReq.Invitation.InviterUserID = s.loginUserID
	signalInviteInGroupReq.OpUserID = s.loginUserID
	signalInviteInGroupReq.Invitation.InitiateTime = int32(utils.GetCurrentTimestampBySecond())
	participants, err := s.getSelfParticipant(ctx, signalInviteInGroupReq.Invitation.GroupID)
	if err != nil {
		return nil, err
	}
	signalInviteInGroupReq.Participant = participants
	req := &sdkws.SignalReq{Payload: &sdkws.SignalReq_InviteInGroup{InviteInGroup: signalInviteInGroupReq}}
	resp, err := s.SendSignalingReqWaitResp(ctx, req)
	if err != nil {
		return nil, err
	}
	s.isCanceled = false
	reply := resp.Payload.(*sdkws.SignalResp_InviteInGroup).InviteInGroup
	go s.waitPush(ctx, req, reply.BusyLineUserIDList)
	return reply, nil
}

func (s *LiveSignaling) SignalingInvite(ctx context.Context, signalInviteReq *sdkws.SignalInviteReq) (*sdkws.SignalInviteReply, error) {
	s.setDefaultReq(signalInviteReq.Invitation)
	signalInviteReq.Invitation.InviterUserID = s.loginUserID
	signalInviteReq.OpUserID = s.loginUserID
	signalInviteReq.Invitation.InitiateTime = int32(utils.GetCurrentTimestampBySecond())
	participants, err := s.getSelfParticipant(ctx, signalInviteReq.Invitation.GroupID)
	if err != nil {
		return nil, err
	}
	signalInviteReq.Participant = participants
	req := &sdkws.SignalReq{Payload: &sdkws.SignalReq_Invite{Invite: signalInviteReq}}
	resp, err := s.SendSignalingReqWaitResp(ctx, req)
	if err != nil {
		return nil, err
	}
	s.isCanceled = false
	reply := resp.Payload.(*sdkws.SignalResp_Invite).Invite
	go s.waitPush(ctx, req, reply.BusyLineUserIDList)
	return reply, nil
}

func (s *LiveSignaling) SignalingAccept(ctx context.Context, signalAcceptReq *sdkws.SignalAcceptReq) (*sdkws.SignalAcceptReply, error) {
	s.setDefaultReq(signalAcceptReq.Invitation)
	signalAcceptReq.OpUserID = s.loginUserID
	signalAcceptReq.Invitation.InitiateTime = int32(utils.GetCurrentTimestampBySecond())
	participants, err := s.getSelfParticipant(context.Background(), signalAcceptReq.Invitation.GroupID)
	if err != nil {
		return nil, err
	}
	signalAcceptReq.Participant = participants
	req := &sdkws.SignalReq{Payload: &sdkws.SignalReq_Accept{Accept: signalAcceptReq}}
	resp, err := s.SendSignalingReqWaitResp(ctx, req)
	if err != nil {
		return nil, err
	}
	reply := resp.Payload.(*sdkws.SignalResp_Accept).Accept
	return reply, nil
}

func (s *LiveSignaling) SignalingReject(ctx context.Context, signalRejectReq *sdkws.SignalRejectReq) error {
	s.setDefaultReq(signalRejectReq.Invitation)
	signalRejectReq.OpUserID = s.loginUserID
	signalRejectReq.Invitation.InitiateTime = int32(utils.GetCurrentTimestampBySecond())
	signalRejectReq.OpUserPlatformID = s.platformID
	participant, err := s.getSelfParticipant(ctx, signalRejectReq.Invitation.GroupID)
	if err != nil {
		return err
	}
	signalRejectReq.Participant = participant
	req := &sdkws.SignalReq{Payload: &sdkws.SignalReq_Reject{Reject: signalRejectReq}}
	_, err = s.SendSignalingReqWaitResp(ctx, req)
	return err
}

func (s *LiveSignaling) SignalingCancel(ctx context.Context, signalCancelReq *sdkws.SignalCancelReq) error {
	s.setDefaultReq(signalCancelReq.Invitation)
	signalCancelReq.OpUserID = s.loginUserID
	signalCancelReq.Invitation.InitiateTime = int32(utils.GetCurrentTimestampBySecond())
	participant, err := s.getSelfParticipant(ctx, signalCancelReq.Invitation.GroupID)
	if err != nil {
		return err
	}
	signalCancelReq.Participant = participant
	req := &sdkws.SignalReq{Payload: &sdkws.SignalReq_Cancel{Cancel: signalCancelReq}}
	_, err = s.SendSignalingReqWaitResp(ctx, req)
	if err != nil {
		return err
	}
	s.isCanceled = true
	return nil
}

func (s *LiveSignaling) SignalingHungUp(ctx context.Context, signalHungUpReq *sdkws.SignalHungUpReq) error {
	s.setDefaultReq(signalHungUpReq.Invitation)
	signalHungUpReq.OpUserID = s.loginUserID
	signalHungUpReq.Invitation.InitiateTime = int32(utils.GetCurrentTimestampBySecond())
	req := &sdkws.SignalReq{Payload: &sdkws.SignalReq_HungUp{HungUp: signalHungUpReq}}
	_, err := s.SendSignalingReqWaitResp(ctx, req)
	return err
}

func (s *LiveSignaling) SignalingGetRoomByGroupID(ctx context.Context, groupID string) (*sdkws.SignalGetRoomByGroupIDReply, error) {
	req := &sdkws.SignalReq_GetRoomByGroupID{GetRoomByGroupID: &sdkws.SignalGetRoomByGroupIDReq{
		OpUserID: s.loginUserID,
		GroupID:  groupID,
	}}
	participant, err := s.getSelfParticipant(ctx, groupID)
	if err != nil {
		return nil, err
	}
	req.GetRoomByGroupID.Participant = participant
	var signalReq sdkws.SignalReq
	signalReq.Payload = req
	resp, err := s.SendSignalingReqWaitResp(ctx, &signalReq)
	if err != nil {
		return nil, err
	}
	return resp.Payload.(*sdkws.SignalResp_GetRoomByGroupID).GetRoomByGroupID, nil
}

func (s *LiveSignaling) SignalingGetTokenByRoomID(ctx context.Context, groupID string) (*sdkws.SignalGetTokenByRoomIDReply, error) {
	req := &sdkws.SignalReq_GetTokenByRoomID{GetTokenByRoomID: &sdkws.SignalGetTokenByRoomIDReq{
		OpUserID: s.loginUserID,
		RoomID:   groupID,
	}}
	participant, err := s.getSelfParticipant(ctx, groupID)
	if err != nil {
		return nil, err
	}
	req.GetTokenByRoomID.Participant = participant
	var signalReq sdkws.SignalReq
	signalReq.Payload = req
	resp, err := s.SendSignalingReqWaitResp(ctx, &signalReq)
	if err != nil {
		return nil, err
	}
	return resp.Payload.(*sdkws.SignalResp_GetTokenByRoomID).GetTokenByRoomID, nil
}
