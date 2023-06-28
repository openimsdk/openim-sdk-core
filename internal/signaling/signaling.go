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
	"open_im_sdk/internal/interaction"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/sdkerrs"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
)

type LiveSignaling struct {
	*interaction.LongConnMgr
	listener    open_im_sdk_callback.OnSignalingListener
	loginUserID string
	db_interface.DataBase
	platformID int32
	isCanceled bool

	listenerForService open_im_sdk_callback.OnSignalingListener
}

func NewLiveSignaling(longConnMgr *interaction.LongConnMgr, loginUserID string, platformID int32, db db_interface.DataBase) *LiveSignaling {
	return &LiveSignaling{LongConnMgr: longConnMgr, loginUserID: loginUserID, platformID: platformID, DataBase: db}
}

func (s *LiveSignaling) setDefaultReq(req *server_api_params.InvitationInfo) {
	if req.RoomID == "" {
		req.RoomID = utils.OperationIDGenerator()
	}
	if req.Timeout == 0 {
		req.Timeout = 60 * 60
	}
}

func (s *LiveSignaling) checkInvitation(invitation *server_api_params.InvitationInfo) error {
	if invitation == nil {
		return sdkerrs.ErrArgs.Wrap("invitation is nil")
	}
	return nil
}

func (s *LiveSignaling) waitPush(ctx context.Context, req *server_api_params.SignalReq, busyLineUserList []string) {
	var invt *server_api_params.InvitationInfo
	switch payload := req.Payload.(type) {
	case *server_api_params.SignalReq_Invite:
		invt = payload.Invite.Invitation
	case *server_api_params.SignalReq_InviteInGroup:
		invt = payload.InviteInGroup.Invitation
	}
	var listenerList []open_im_sdk_callback.OnSignalingListener
	if s.listener != nil {
		listenerList = append(listenerList, s.listener)
	}
	if s.listenerForService != nil {
		listenerList = append(listenerList, s.listenerForService)
	}
	if len(listenerList) == 0 {
		log.ZError(ctx, "no listner", nil)
		return
	}
	var inviteeUserIDList []string
	for _, inviteUser := range invt.InviteeUserIDList {
		if !utils.IsContain(inviteUser, busyLineUserList) {
			inviteeUserIDList = append(inviteeUserIDList, inviteUser)
		}
	}
	for _, v := range inviteeUserIDList {
		go func(invitee string) {
			push, err := s.SignalingWaitPush(ctx, invt.InviterUserID, invitee, invt.RoomID, invt.Timeout)
			if err != nil {
				if errs.Unwrap(err) == sdkerrs.ErrNetworkTimeOut {
					log.ZError(ctx, "timeout", err, "invitee", invitee, "roomID", invt.RoomID, "timeout", invt.Timeout)
					switch payload := req.Payload.(type) {
					case *server_api_params.SignalReq_Invite:
						if !s.isCanceled {
							for _, listener := range listenerList {
								listener.OnInvitationTimeout(utils.StructToJsonString(payload.Invite))
							}
						}
					case *server_api_params.SignalReq_InviteInGroup:
						if !s.isCanceled {
							for _, listener := range listenerList {
								listener.OnInvitationTimeout(utils.StructToJsonString(payload.InviteInGroup))
							}
						}
					}
				} else {
					log.ZError(ctx, "wait push failed", err, "inviterUserID", invt.InviterUserID, "invitee", invitee, "roomID", invt.RoomID, "timeout", invt.Timeout)
				}
				return
			}
			log.ZInfo(ctx, "signalingWaitPush ", "push", push.String(), "inviterUserID", invt.InviterUserID, "invitee", invitee, "roomID", invt.RoomID, "timeout", invt.Timeout)
			s.doSignalPush(ctx, push)
		}(v)
	}
}
func (s *LiveSignaling) doSignalPush(ctx context.Context, req *server_api_params.SignalReq) {
	var listenerList []open_im_sdk_callback.OnSignalingListener
	if s.listener != nil {
		listenerList = append(listenerList, s.listener)
	}
	if s.listenerForService != nil {
		listenerList = append(listenerList, s.listenerForService)
	}
	if len(listenerList) == 0 {
		log.ZError(ctx, "no listner", nil)
		return
	}
	log.ZDebug(ctx, "doSignalPush", "req", req)
	switch payload := req.Payload.(type) {
	case *server_api_params.SignalReq_Accept:
		for _, listener := range listenerList {
			listener.OnInviteeAccepted(utils.StructToJsonString(payload.Accept))
		}

	case *server_api_params.SignalReq_Reject:
		for _, listener := range listenerList {
			listener.OnInviteeRejected(utils.StructToJsonString(payload.Reject))
		}
	default:
		log.ZError(ctx, "resp payload type failed", nil, "req", req)
	}
}

func (s *LiveSignaling) SetListener(listener open_im_sdk_callback.OnSignalingListener) {
	s.listener = listener
}

func (s *LiveSignaling) SetListenerForService(listener open_im_sdk_callback.OnSignalingListener) {
	s.listenerForService = listener
}

func (s *LiveSignaling) getSelfParticipant(ctx context.Context, groupID string) (*server_api_params.ParticipantMetaData, error) {
	p := server_api_params.ParticipantMetaData{GroupInfo: &sdkws.GroupInfo{}, GroupMemberInfo: &sdkws.GroupMemberFullInfo{}, UserInfo: &sdkws.PublicUserInfo{}}
	if groupID != "" {
		group, err := s.GetGroupInfoByGroupID(ctx, groupID)
		if err != nil {
			return nil, err
		}
		copier.Copy(p.GroupInfo, group)
		groupMemberInfo, err := s.GetGroupMemberInfoByGroupIDUserID(ctx, groupID, s.loginUserID)
		if err != nil {
			return nil, err
		}
		copier.Copy(p.GroupMemberInfo, groupMemberInfo)
	}
	user, err := s.GetLoginUser(ctx, s.loginUserID)
	if err != nil {
		return nil, err
	}
	copier.Copy(p.UserInfo, user)
	return &p, nil
}

func (s *LiveSignaling) DoNotification(ctx context.Context, msg *sdkws.MsgData) {
	var listenerList []open_im_sdk_callback.OnSignalingListener
	if s.listener != nil {
		listenerList = append(listenerList, s.listener)
	}
	if s.listenerForService != nil {
		listenerList = append(listenerList, s.listenerForService)
	}
	if len(listenerList) == 0 {
		log.ZError(ctx, "no listner", nil)
		return
	}
	switch msg.ContentType {
	case constant.SignalingNotification:
		log.ZInfo(ctx, utils.GetSelfFuncName(), "args ", msg.String())
		var resp server_api_params.SignalReq
		err := proto.Unmarshal(msg.Content, &resp)
		if err != nil {
			log.ZError(ctx, "Unmarshal failed", err, "msg", msg)
			return
		}
		switch payload := resp.Payload.(type) {
		case *server_api_params.SignalReq_Accept:
			if payload.Accept.Invitation.InviterUserID == s.loginUserID && payload.Accept.Invitation.PlatformID == s.platformID {
				var wsResp interaction.GeneralWsResp
				wsResp.ReqIdentifier = constant.SendSignalMsg
				wsResp.Data = msg.Content
				wsResp.MsgIncr = s.loginUserID + payload.Accept.UserID + payload.Accept.Invitation.RoomID
				log.ZDebug(ctx, "search msgIncr", wsResp.MsgIncr)
				if err := s.LongConnMgr.Syncer.NotifyResp(ctx, wsResp); err != nil {
					log.ZError(ctx, "notifyResp failed", err, "wsResp", wsResp)
				}
				return
			}
			if payload.Accept.OpUserPlatformID != s.platformID && payload.Accept.UserID == s.loginUserID {
				for _, listener := range listenerList {
					listener.OnInviteeAcceptedByOtherDevice(utils.StructToJsonString(payload.Accept))
					log.ZDebug(ctx, "OnInviteeAcceptedByOtherDevice", "accept", utils.StructToJsonString(payload.Accept))
				}
				return
			}
		case *server_api_params.SignalReq_Reject:
			if payload.Reject.Invitation.InviterUserID == s.loginUserID && payload.Reject.Invitation.PlatformID == s.platformID {
				var wsResp interaction.GeneralWsResp
				wsResp.ReqIdentifier = constant.SendSignalMsg
				wsResp.Data = msg.Content
				wsResp.MsgIncr = s.loginUserID + payload.Reject.UserID + payload.Reject.Invitation.RoomID
				log.ZDebug(ctx, "search msgIncr: ", wsResp.MsgIncr)
				if err := s.LongConnMgr.Syncer.NotifyResp(ctx, wsResp); err != nil {
					log.ZError(ctx, "notifyResp failed", err, "wsResp", wsResp)
				}
				return
			}
			if payload.Reject.OpUserPlatformID != s.platformID && payload.Reject.UserID == s.loginUserID {
				for _, listener := range listenerList {
					listener.OnInviteeRejectedByOtherDevice(utils.StructToJsonString(payload.Reject))
					log.ZDebug(ctx, "OnInviteeRejectedByOtherDevice", "reject", utils.StructToJsonString(payload.Reject))
				}
				return
			}

		case *server_api_params.SignalReq_HungUp:
			if s.loginUserID != payload.HungUp.UserID {
				for _, listener := range listenerList {
					listener.OnHangUp(utils.StructToJsonString(payload.HungUp))
					log.ZDebug(ctx, "OnHangUp", "hungUp", utils.StructToJsonString(payload.HungUp))
				}
			}
		case *server_api_params.SignalReq_Cancel:
			if utils.IsContain(s.loginUserID, payload.Cancel.Invitation.InviteeUserIDList) {
				for _, listener := range listenerList {
					listener.OnInvitationCancelled(utils.StructToJsonString(payload.Cancel))
					log.ZDebug(ctx, "OnInvitationCancelled", "cancel", utils.StructToJsonString(payload.Cancel))
				}
			}
		case *server_api_params.SignalReq_Invite:
			if utils.IsContain(s.loginUserID, payload.Invite.Invitation.InviteeUserIDList) {
				for _, listener := range listenerList {
					if !utils.IsContain(s.loginUserID, payload.Invite.Invitation.BusyLineUserIDList) {
						listener.OnReceiveNewInvitation(utils.StructToJsonString(payload.Invite))
						log.ZDebug(ctx, "OnReceiveNewInvitation", "invite", utils.StructToJsonString(payload.Invite))
					}
				}
			}

		case *server_api_params.SignalReq_InviteInGroup:
			if utils.IsContain(s.loginUserID, payload.InviteInGroup.Invitation.InviteeUserIDList) {
				for _, listener := range listenerList {
					if !utils.IsContain(s.loginUserID, payload.InviteInGroup.Invitation.BusyLineUserIDList) {
						listener.OnReceiveNewInvitation(utils.StructToJsonString(payload.InviteInGroup))
						log.ZDebug(ctx, "OnReceiveNewInvitation", "inviteInGroup", utils.StructToJsonString(payload.InviteInGroup))
					}
				}
			}
		default:
			log.ZError(ctx, "resp payload type failed", nil, "msg", msg)
		}
	case constant.CustomSignalNotification:
		var callback server_api_params.SignalSendCustomSignalReq
		if err := proto.Unmarshal(msg.Content, &callback); err != nil {
			log.ZError(ctx, "proto.Unmarshal failed", err, "msg", msg)
			return
		}
		for _, listener := range listenerList {
			listener.OnReceiveCustomSignal(utils.StructToJsonString(&callback))
			log.ZDebug(ctx, "SignalSendCustomSignalReq", "onReceiveCustomSignal", &callback)
		}
	case constant.StreamChangedNotification:
		var callback server_api_params.SignalOnStreamChangeReq
		if err := proto.Unmarshal(msg.Content, &callback); err != nil {
			log.ZError(ctx, "proto.Unmarshal failed", err, "msg", msg)
			return
		}
		for _, listener := range listenerList {
			listener.OnStreamChange(utils.StructToJsonString(&callback))
			log.ZDebug(ctx, "SignalOnStreamChangeReq", "onStreamChangeReq", &callback)
		}

	case constant.RoomParticipantsConnectedNotification:
		var callback server_api_params.SignalOnRoomParticipantConnectedReq
		if err := proto.Unmarshal(msg.Content, &callback); err != nil {
			log.ZError(ctx, "proto.Unmarshal failed", err, "msg", msg)
			return
		}
		for _, listener := range listenerList {
			listener.OnRoomParticipantConnected(utils.StructToJsonString(&callback))
			log.ZDebug(ctx, "SignalOnRoomParticipantConnectedReq", "onRoomParticipantConnectedReq", &callback)
		}
	case constant.RoomParticipantsDisconnectedNotification:
		var callback server_api_params.SignalOnRoomParticipantDisconnectedReq
		if err := proto.Unmarshal(msg.Content, &callback); err != nil {
			log.ZError(ctx, "proto.Unmarshal failed", err, "msg", msg)
			return
		}
		for _, listener := range listenerList {
			listener.OnRoomParticipantDisconnected(utils.StructToJsonString(&callback))
			log.ZDebug(ctx, "SignalOnRoomParticipantDisconnectedReq", "onRoomParticipantDisconnectedReq", &callback)
		}
	}

}
func (s *LiveSignaling) SendSignalingReqWaitResp(ctx context.Context, req *server_api_params.SignalReq) (*server_api_params.SignalResp, error) {
	var signalResp server_api_params.SignalResp
	err := s.LongConnMgr.SendReqWaitResp(ctx, req, constant.SendSignalMsg, &signalResp)
	if err != nil {
		return nil, err
	}
	return &signalResp, nil
}
func (s *LiveSignaling) SignalingWaitPush(ctx context.Context, inviterUserID, inviteeUserID, roomID string, timeout int32) (*server_api_params.SignalReq, error) {
	msgIncr := inviterUserID + inviteeUserID + roomID
	ch := s.LongConnMgr.Syncer.AddChByIncr(msgIncr)
	defer s.LongConnMgr.Syncer.DelCh(msgIncr)
	resp, err := s.LongConnMgr.Syncer.WaitResp(ctx, ch, int(timeout))
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var signalReq server_api_params.SignalReq
	err = proto.Unmarshal(resp.Data, &signalReq)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	return &signalReq, nil
}
