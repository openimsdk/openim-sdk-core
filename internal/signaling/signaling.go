package signaling

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
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
	db_interface.DataBase
	platformID int32
	isCanceled bool

	listenerForService open_im_sdk_callback.OnSignalingListener

	acceptRoomIDs map[string]interface{}
}

func NewLiveSignaling(ws *ws.Ws, loginUserID string, platformID int32, db db_interface.DataBase) *LiveSignaling {
	if ws == nil {
		log.Warn("", " ws is nil")
		return nil
	}
	m := make(map[string]interface{}, 0)
	return &LiveSignaling{Ws: ws, loginUserID: loginUserID, platformID: platformID, DataBase: db, acceptRoomIDs: m}
}

func (s *LiveSignaling) waitPush(req *api.SignalReq, busyLineUserList []string, operationID string) {
	var invt *api.InvitationInfo
	switch payload := req.Payload.(type) {
	case *api.SignalReq_Invite:
		invt = payload.Invite.Invitation
	case *api.SignalReq_InviteInGroup:
		invt = payload.InviteInGroup.Invitation
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
		log.Error(operationID, "len (listenerList) == 0 ")
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
			push, err := s.SignalingWaitPush(invt.InviterUserID, invitee, invt.RoomID, invt.Timeout, operationID)
			if err != nil {
				if strings.Contains(err.Error(), "timeout") {
					log.Error(operationID, "wait push timeout ", err.Error(), invt.InviterUserID, v, invt.RoomID, invt.Timeout)
					switch payload := req.Payload.(type) {
					case *api.SignalReq_Invite:
						if !s.isCanceled {
							for _, listener := range listenerList {
								listener.OnInvitationTimeout(utils.StructToJsonString(payload.Invite))
								log.Info(operationID, "OnInvitationTimeout ", utils.StructToJsonString(payload.Invite), listener)
							}
						}
					case *api.SignalReq_InviteInGroup:
						if !s.isCanceled {
							for _, listener := range listenerList {
								listener.OnInvitationTimeout(utils.StructToJsonString(payload.InviteInGroup))
								log.Info(operationID, "OnInvitationTimeout ", utils.StructToJsonString(payload.InviteInGroup), listener)
							}
						}
					}

				} else {
					log.Error(operationID, "other failed ", err.Error(), invt.InviterUserID, v, invt.RoomID, invt.Timeout)
				}
				return
			}
			log.Info(operationID, "SignalingWaitPush ", push.String(), invt.InviterUserID, v, invt.RoomID, invt.Timeout)
			s.doSignalPush(push, operationID)
		}(v)
	}
}
func (s *LiveSignaling) doSignalPush(req *api.SignalReq, operationID string) {
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
		log.Error(operationID, "len (listenerList) == 0 ")
		return
	}
	switch payload := req.Payload.(type) {
	//case *api.SignalReq_Invite:
	//	log.Info(operationID, "recv signal push ", payload.Invite.String())
	//	s.listener.OnReceiveNewInvitation(utils.StructToJsonString(payload.Invite))
	case *api.SignalReq_Accept:
		log.Info(operationID, "recv signal push Accept ", payload.Accept.String())
		for _, listener := range listenerList {
			listener.OnInviteeAccepted(utils.StructToJsonString(payload.Accept))
			log.Info(operationID, "OnInviteeAccepted ", utils.StructToJsonString(payload.Accept), listener)
		}

	case *api.SignalReq_Reject:
		log.Info(operationID, "recv signal push Reject ", payload.Reject.String())
		for _, listener := range listenerList {
			listener.OnInviteeRejected(utils.StructToJsonString(payload.Reject))
			log.Info(operationID, "OnInviteeRejected ", utils.StructToJsonString(payload.Reject), listener)
		}

	//case *api.SignalReq_HungUp:
	//	log.Info(operationID, "recv signal push HungUp ", payload.HungUp.String())
	//	s.listener.OnHangUp(utils.StructToJsonString(payload.HungUp))
	//case *api.SignalReq_Cancel:
	//	log.Info(operationID, "recv signal push ", payload.Cancel.String())
	//	s.listener.OnInvitationCancelled(utils.StructToJsonString(payload.Cancel))
	//case *api.SignalReq_InviteInGroup:
	//	log.Info(operationID, "recv signal push ", payload.InviteInGroup.String())
	default:
		log.Error(operationID, "payload type failed ", payload)
	}
}

func (s *LiveSignaling) SetListener(listener open_im_sdk_callback.OnSignalingListener) {
	log.Info("", utils.GetSelfFuncName(), "args ", listener)
	s.listener = listener
}

func (s *LiveSignaling) SetListenerForService(listener open_im_sdk_callback.OnSignalingListener) {
	log.Info("", utils.GetSelfFuncName(), "args ", listener)
	s.listenerForService = listener
}

func (s *LiveSignaling) getSelfParticipant(groupID string, callback open_im_sdk_callback.Base, operationID string) *api.ParticipantMetaData {
	log.Info(operationID, utils.GetSelfFuncName(), "args ", groupID)
	p := api.ParticipantMetaData{GroupInfo: &api.GroupInfo{}, GroupMemberInfo: &api.GroupMemberFullInfo{}, UserInfo: &api.PublicUserInfo{}}
	if groupID != "" {
		g, err := s.GetGroupInfoByGroupID(groupID)
		if err != nil {
			log.NewError(operationID, "GetGroupInfoByGroupID failed", err.Error())
			if !strings.Contains(err.Error(), "record not found") {
				common.CheckDBErrCallback(callback, err, operationID)
			}
		} else {
			copier.Copy(p.GroupInfo, g)
			mInfo, err := s.GetGroupMemberInfoByGroupIDUserID(groupID, s.loginUserID)
			common.CheckDBErrCallback(callback, err, operationID)
			copier.Copy(p.GroupMemberInfo, mInfo)
		}
	}

	sf, err := s.GetLoginUser(s.loginUserID)
	common.CheckDBErrCallback(callback, err, operationID)
	copier.Copy(p.UserInfo, sf)
	log.Info(operationID, utils.GetSelfFuncName(), "return ", p)
	return &p
}

func (s *LiveSignaling) DoNotification(msg *api.MsgData, conversationCh chan common.Cmd2Value, operationID string) {
	log.Info(operationID, utils.GetSelfFuncName(), "args ", msg.String())
	var resp api.SignalReq
	err := proto.Unmarshal(msg.Content, &resp)
	if err != nil {
		log.Error(operationID, "Unmarshal failed")
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
		log.Error(operationID, "len (listenerList) == 0 ")
		return
	}
	switch payload := resp.Payload.(type) {
	case *api.SignalReq_Accept:
		log.Info(operationID, "signaling response ", payload.Accept.String())
		if payload.Accept.Invitation.InviterUserID == s.loginUserID && payload.Accept.Invitation.PlatformID == s.platformID {
			if _, ok := s.acceptRoomIDs[payload.Accept.Invitation.RoomID]; ok {
				var wsResp ws.GeneralWsResp
				wsResp.ReqIdentifier = constant.WSSendSignalMsg
				wsResp.Data = msg.Content
				wsResp.MsgIncr = s.loginUserID + payload.Accept.OpUserID + payload.Accept.Invitation.RoomID
				log.Info(operationID, "search msgIncr: ", wsResp.MsgIncr)
				s.DoWSSignal(wsResp)
				return
			} else {
				for _, listener := range listenerList {
					listener.OnInviteeAcceptedByOtherDevice(utils.StructToJsonString(payload.Accept))
					log.Info(operationID, "OnInviteeAcceptedByOtherDevice same platform", utils.StructToJsonString(payload.Accept), listener)
				}
				return
			}
		}
		if payload.Accept.OpUserPlatformID != s.platformID && payload.Accept.OpUserID == s.loginUserID {
			for _, listener := range listenerList {
				listener.OnInviteeAcceptedByOtherDevice(utils.StructToJsonString(payload.Accept))
				log.Info(operationID, "OnInviteeAcceptedByOtherDevice ", utils.StructToJsonString(payload.Accept), listener)
			}
			return
		}
	case *api.SignalReq_Reject:
		log.Info(operationID, "signaling response ", payload.Reject.String())
		if payload.Reject.Invitation.InviterUserID == s.loginUserID && payload.Reject.Invitation.PlatformID == s.platformID {
			if _, ok := s.acceptRoomIDs[payload.Reject.Invitation.RoomID]; ok {
				var wsResp ws.GeneralWsResp
				wsResp.ReqIdentifier = constant.WSSendSignalMsg
				wsResp.Data = msg.Content
				wsResp.MsgIncr = s.loginUserID + payload.Reject.OpUserID + payload.Reject.Invitation.RoomID
				log.Info(operationID, "search msgIncr: ", wsResp.MsgIncr)
				s.DoWSSignal(wsResp)
				return
			} else {
				for _, listener := range listenerList {
					listener.OnInviteeRejectedByOtherDevice(utils.StructToJsonString(payload.Reject))
					log.Info(operationID, "OnInviteeRejectedByOtherDevice ", utils.StructToJsonString(payload.Reject), listener)
				}
				return
			}
		}
		if payload.Reject.OpUserPlatformID != s.platformID && payload.Reject.OpUserID == s.loginUserID {
			for _, listener := range listenerList {
				listener.OnInviteeRejectedByOtherDevice(utils.StructToJsonString(payload.Reject))
				log.Info(operationID, "OnInviteeRejectedByOtherDevice same platform", utils.StructToJsonString(payload.Reject), listener)
			}
			return
		}

	case *api.SignalReq_HungUp:
		log.Info(operationID, "signaling response HungUp", payload.HungUp.String())
		if s.loginUserID != payload.HungUp.OpUserID {
			for _, listener := range listenerList {
				listener.OnHangUp(utils.StructToJsonString(payload.HungUp))
				log.Info(operationID, "OnHangUp ", utils.StructToJsonString(payload.HungUp), listener)
			}
		}
	case *api.SignalReq_Cancel:
		log.Info(operationID, "signaling response ", payload.Cancel.String())
		if utils.IsContain(s.loginUserID, payload.Cancel.Invitation.InviteeUserIDList) {
			for _, listener := range listenerList {
				listener.OnInvitationCancelled(utils.StructToJsonString(payload.Cancel))
				log.Info(operationID, "OnInvitationCancelled ", utils.StructToJsonString(payload.Cancel), listener)
			}
		}
	case *api.SignalReq_Invite:
		log.Info(operationID, "signaling response ", payload.Invite.String())
		if utils.IsContain(s.loginUserID, payload.Invite.Invitation.InviteeUserIDList) {
			for _, listener := range listenerList {
				if !utils.IsContain(s.loginUserID, payload.Invite.Invitation.BusyLineUserIDList) {
					listener.OnReceiveNewInvitation(utils.StructToJsonString(payload.Invite))
					log.Info(operationID, "OnReceiveNewInvitation ", utils.StructToJsonString(payload.Invite), listener)
				}
			}
		}

	case *api.SignalReq_InviteInGroup:
		log.Info(operationID, "signaling response ", payload.InviteInGroup.String())
		if utils.IsContain(s.loginUserID, payload.InviteInGroup.Invitation.InviteeUserIDList) {
			for _, listener := range listenerList {
				if !utils.IsContain(s.loginUserID, payload.InviteInGroup.Invitation.BusyLineUserIDList) {
					listener.OnReceiveNewInvitation(utils.StructToJsonString(payload.InviteInGroup))
					log.Info(operationID, "OnReceiveNewInvitation ", utils.StructToJsonString(payload.InviteInGroup), listener)
				}
			}
		}
	case *api.SignalReq_OnRoomParticipantConnectedReq:
		log.Info(operationID, "signaling response ", payload.OnRoomParticipantConnectedReq.String())
		for _, listener := range listenerList {
			listener.OnRoomParticipantConnected(utils.StructToJsonString(payload.OnRoomParticipantConnectedReq))
			log.Info(operationID, "SignalOnRoomParticipantConnectedReq", utils.StructToJsonString(payload.OnRoomParticipantConnectedReq), listener)
		}
	case *api.SignalReq_OnRoomParticipantDisconnectedReq:
		log.Info(operationID, "signaling response ", payload.OnRoomParticipantDisconnectedReq.String())
		for _, listener := range listenerList {
			listener.OnRoomParticipantDisconnected(utils.StructToJsonString(payload.OnRoomParticipantDisconnectedReq))
			log.Info(operationID, "SignalOnRoomParticipantDisconnectedReq", utils.StructToJsonString(payload.OnRoomParticipantDisconnectedReq), listener)
		}

	default:
		log.Error(operationID, "resp payload type failed ", payload)
	}
}

func (s *LiveSignaling) handleSignaling(req *api.SignalReq, callback open_im_sdk_callback.Base, operationID string) {
	log.Info(operationID, utils.GetSelfFuncName(), "args ", req.String())
	resp, err := s.SendSignalingReqWaitResp(req, operationID)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "SendSignalingReqWaitResp error", err.Error())
		common.CheckAnyErrCallback(callback, 3003, err, operationID)
	}
	common.CheckAnyErrCallback(callback, 3001, err, operationID)
	var busyLineUserIDList []string
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
		s.isCanceled = true
		log.Info(operationID, "signaling response ", payload.Cancel.String())
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.CancelCallback(payload.Cancel)))
	case *api.SignalResp_Invite:
		s.isCanceled = false
		busyLineUserIDList = payload.Invite.BusyLineUserIDList
		log.Info(operationID, "signaling response ", payload.Invite.String())
		s.acceptRoomIDs[payload.Invite.RoomID] = nil
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.InviteCallback(payload.Invite)))
	case *api.SignalResp_InviteInGroup:
		s.isCanceled = false
		busyLineUserIDList = payload.InviteInGroup.BusyLineUserIDList
		log.Info(operationID, "signaling response ", payload.InviteInGroup.String())
		s.acceptRoomIDs[payload.InviteInGroup.RoomID] = nil
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.InviteInGroupCallback(payload.InviteInGroup)))
	case *api.SignalResp_GetRoomByGroupID:
		log.Info(operationID, "signaling response ", payload.GetRoomByGroupID.String())
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.GetRoomByGroupIDCallback(payload.GetRoomByGroupID)))
	case *api.SignalResp_GetTokenByRoomID:
		log.Info(operationID, "signaling response", payload.GetTokenByRoomID.String())
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.GetTokenByRoomID(payload.GetTokenByRoomID)))
	default:
		log.Error(operationID, "resp payload type failed ", payload)
		common.CheckAnyErrCallback(callback, 3002, errors.New("resp payload type failed"), operationID)
	}
	switch req.Payload.(type) {
	case *api.SignalReq_Invite:
		log.Info(operationID, "wait push ", req.String())
		s.waitPush(req, busyLineUserIDList, operationID)
	case *api.SignalReq_InviteInGroup:
		log.Info(operationID, "wait push ", req.String())
		s.waitPush(req, busyLineUserIDList, operationID)
	}
}
