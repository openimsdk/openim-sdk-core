// Copyright © 2023 OpenIM SDK. All rights reserved.
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
	"math/big"
	"math/rand"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/sdkerrs"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
)

func (s *LiveSignaling) SignalingInviteInGroup(ctx context.Context, signalInviteInGroupReq *server_api_params.SignalInviteInGroupReq) (*server_api_params.SignalInviteInGroupResp, error) {
	if err := s.checkInvitation(signalInviteInGroupReq.Invitation); err != nil {
		return nil, err
	}
	s.setDefaultReq(signalInviteInGroupReq.Invitation)
	log.ZDebug(ctx, "x", "invitation", signalInviteInGroupReq.Invitation.InviterUserID, "login", s.loginUserID)
	signalInviteInGroupReq.Invitation.InviterUserID = s.loginUserID
	signalInviteInGroupReq.UserID = s.loginUserID
	signalInviteInGroupReq.Invitation.InitiateTime = utils.GetCurrentTimestampBySecond()
	participants, err := s.getSelfParticipant(ctx, signalInviteInGroupReq.Invitation.GroupID)
	if err != nil {
		return nil, err
	}
	signalInviteInGroupReq.Participant = participants
	req := &server_api_params.SignalReq{Payload: &server_api_params.SignalReq_InviteInGroup{InviteInGroup: signalInviteInGroupReq}}
	resp, err := s.SendSignalingReqWaitResp(ctx, req)
	if err != nil {
		return nil, err
	}
	s.isCanceled = false
	reply := resp.Payload.(*server_api_params.SignalResp_InviteInGroup).InviteInGroup
	go s.waitPush(ctx, req, reply.BusyLineUserIDList)
	return reply, nil
}

func (s *LiveSignaling) SignalingInvite(ctx context.Context, signalInviteReq *server_api_params.SignalInviteReq) (*server_api_params.SignalInviteResp, error) {
	if err := s.checkInvitation(signalInviteReq.Invitation); err != nil {
		return nil, err
	}
	s.setDefaultReq(signalInviteReq.Invitation)
	signalInviteReq.Invitation.InviterUserID = s.loginUserID
	signalInviteReq.UserID = s.loginUserID
	signalInviteReq.Invitation.InitiateTime = utils.GetCurrentTimestampBySecond()
	participants, err := s.getSelfParticipant(ctx, signalInviteReq.Invitation.GroupID)
	if err != nil {
		return nil, err
	}
	signalInviteReq.Participant = participants
	req := &server_api_params.SignalReq{Payload: &server_api_params.SignalReq_Invite{Invite: signalInviteReq}}
	resp, err := s.SendSignalingReqWaitResp(ctx, req)
	if err != nil {
		return nil, err
	}
	s.isCanceled = false
	reply := resp.Payload.(*server_api_params.SignalResp_Invite).Invite
	go s.waitPush(ctx, req, reply.BusyLineUserIDList)
	return reply, nil
}

func (s *LiveSignaling) SignalingAccept(ctx context.Context, signalAcceptReq *server_api_params.SignalAcceptReq) (*server_api_params.SignalAcceptResp, error) {
	if err := s.checkInvitation(signalAcceptReq.Invitation); err != nil {
		return nil, err
	}
	s.setDefaultReq(signalAcceptReq.Invitation)
	signalAcceptReq.UserID = s.loginUserID
	signalAcceptReq.Invitation.InitiateTime = utils.GetCurrentTimestampBySecond()
	signalAcceptReq.OpUserPlatformID = s.platformID
	participants, err := s.getSelfParticipant(ctx, signalAcceptReq.Invitation.GroupID)
	if err != nil {
		return nil, err
	}
	signalAcceptReq.Participant = participants
	req := &server_api_params.SignalReq{Payload: &server_api_params.SignalReq_Accept{Accept: signalAcceptReq}}
	resp, err := s.SendSignalingReqWaitResp(ctx, req)
	if err != nil {
		return nil, err
	}
	reply := resp.Payload.(*server_api_params.SignalResp_Accept).Accept
	return reply, nil
}

func (s *LiveSignaling) SignalingReject(ctx context.Context, signalRejectReq *server_api_params.SignalRejectReq) error {
	if err := s.checkInvitation(signalRejectReq.Invitation); err != nil {
		return err
	}
	s.setDefaultReq(signalRejectReq.Invitation)
	signalRejectReq.UserID = s.loginUserID
	signalRejectReq.Invitation.InitiateTime = utils.GetCurrentTimestampBySecond()
	signalRejectReq.OpUserPlatformID = s.platformID
	participant, err := s.getSelfParticipant(ctx, signalRejectReq.Invitation.GroupID)
	if err != nil {
		return err
	}
	signalRejectReq.Participant = participant
	req := &server_api_params.SignalReq{Payload: &server_api_params.SignalReq_Reject{Reject: signalRejectReq}}
	_, err = s.SendSignalingReqWaitResp(ctx, req)
	return err
}

func (s *LiveSignaling) SignalingCancel(ctx context.Context, signalCancelReq *server_api_params.SignalCancelReq) error {
	if err := s.checkInvitation(signalCancelReq.Invitation); err != nil {
		return err
	}
	s.setDefaultReq(signalCancelReq.Invitation)
	signalCancelReq.UserID = s.loginUserID
	signalCancelReq.Invitation.InitiateTime = utils.GetCurrentTimestampBySecond()
	participant, err := s.getSelfParticipant(ctx, signalCancelReq.Invitation.GroupID)
	if err != nil {
		return err
	}
	signalCancelReq.Participant = participant
	req := &server_api_params.SignalReq{Payload: &server_api_params.SignalReq_Cancel{Cancel: signalCancelReq}}
	_, err = s.SendSignalingReqWaitResp(ctx, req)
	if err != nil {
		return err
	}
	s.isCanceled = true
	return nil
}

func (s *LiveSignaling) SignalingHungUp(ctx context.Context, signalHungUpReq *server_api_params.SignalHungUpReq) error {
	if err := s.checkInvitation(signalHungUpReq.Invitation); err != nil {
		return err
	}
	s.setDefaultReq(signalHungUpReq.Invitation)
	signalHungUpReq.UserID = s.loginUserID
	signalHungUpReq.Invitation.InitiateTime = utils.GetCurrentTimestampBySecond()
	req := &server_api_params.SignalReq{Payload: &server_api_params.SignalReq_HungUp{HungUp: signalHungUpReq}}
	_, err := s.SendSignalingReqWaitResp(ctx, req)
	return err
}

func (s *LiveSignaling) SignalingGetTokenByRoomID(ctx context.Context, groupID string) (*server_api_params.SignalGetTokenByRoomIDResp, error) {
	participant, err := s.getSelfParticipant(ctx, groupID)
	if err != nil {
		return nil, err
	}
	signalReq := &server_api_params.SignalReq{Payload: &server_api_params.SignalReq_GetTokenByRoomID{GetTokenByRoomID: &server_api_params.SignalGetTokenByRoomIDReq{
		RoomID: groupID, UserID: s.loginUserID, Participant: participant,
	}}}
	resp, err := s.SendSignalingReqWaitResp(ctx, signalReq)
	if err != nil {
		return nil, err
	}
	return resp.Payload.(*server_api_params.SignalResp_GetTokenByRoomID).GetTokenByRoomID, nil
}

func (s *LiveSignaling) SignalingGetRoomByGroupID(ctx context.Context, groupID string) (*server_api_params.SignalGetRoomByGroupIDResp, error) {
	req := &server_api_params.SignalGetRoomByGroupIDReq{GroupID: groupID}
	return util.CallApi[server_api_params.SignalGetRoomByGroupIDResp](ctx, constant.SignalGetRoomByGroupIDRouter, req)
}

func (s *LiveSignaling) GetSignalingInvitationInfoStartApp(ctx context.Context) (*server_api_params.GetSignalInvitationInfoStartAppResp, error) {
	req := &server_api_params.GetSignalInvitationInfoStartAppReq{UserID: s.loginUserID}
	return util.CallApi[server_api_params.GetSignalInvitationInfoStartAppResp](ctx, constant.GetSignalInvitationInfoStartAppRouter, req)
}

func (s *LiveSignaling) SignalingCreateMeeting(ctx context.Context, req *server_api_params.SignalCreateMeetingReq) (*server_api_params.SignalCreateMeetingResp, error) {
	participant, err := s.getSelfParticipant(ctx, "")
	if err != nil {
		return nil, err
	}
	req.MeetingHostUserID = s.loginUserID
	req.Participant = participant
	bi := big.NewInt(0)
	bi.SetString(utils.Md5(s.loginUserID + utils.Int64ToString(rand.Int63n(time.Now().UnixNano())))[0:8], 16)
	req.RoomID = bi.String()
	return util.CallApi[server_api_params.SignalCreateMeetingResp](ctx, constant.SignalCreateMeetingRouter, req)
}

func (s *LiveSignaling) SignalingJoinMeeting(ctx context.Context, req *server_api_params.SignalJoinMeetingReq) (*server_api_params.SignalJoinMeetingResp, error) {
	participant, err := s.getSelfParticipant(ctx, "")
	if err != nil {
		return nil, err
	}
	req.Participant = participant
	req.UserID = s.loginUserID
	return util.CallApi[server_api_params.SignalJoinMeetingResp](ctx, constant.SignalJoinMeetingRouter, req)
}

func (s *LiveSignaling) SignalingUpdateMeetingInfo(ctx context.Context, req *server_api_params.SignalUpdateMeetingInfoReq) error {
	if req.RoomID == "" {
		return sdkerrs.ErrRecordNotFound.Wrap("roomID is empty")
	}
	return util.ApiPost(ctx, constant.SignalUpdateMeetingInfoRouter, req, nil)
}

func (s *LiveSignaling) SignalingCloseRoom(ctx context.Context, roomID string) error {
	req := &server_api_params.SignalCloseRoomReq{RoomID: roomID}
	return util.ApiPost(ctx, constant.SignalCloseRoomRouter, req, nil)
}

func (s *LiveSignaling) SignalingGetMeetings(ctx context.Context) (resp *server_api_params.SignalGetMeetingsResp, err error) {
	req := &server_api_params.SignalGetMeetingsReq{UserID: s.loginUserID}
	return util.CallApi[server_api_params.SignalGetMeetingsResp](ctx, constant.SignalGetMeetingsRouter, req)
}

func (s *LiveSignaling) SignalingOperateStream(ctx context.Context, streamType, roomID, userID string, mute, muteAll bool) error {
	req := &server_api_params.SignalOperateStreamReq{RoomID: roomID, UserID: userID, StreamType: streamType, Mute: mute, MuteAll: muteAll}
	return util.ApiPost(ctx, constant.SignalOperateStreamRouter, req, nil)
}

func (s *LiveSignaling) SignalingSendCustomSignal(ctx context.Context, customInfo, roomID string) error {
	req := &server_api_params.SignalSendCustomSignalReq{RoomID: roomID, CustomInfo: customInfo}
	return util.ApiPost(ctx, constant.SignalSendCustomSignalRouter, req, nil)
}
