package signaling

import (
	"context"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/log"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"strings"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
)

type LiveSignaling struct {
	*ws.Ws
	listener    open_im_sdk_callback.OnSignalingListener
	loginUserID string
	db_interface.DataBase
	platformID int32
	isCanceled bool

	listenerForService open_im_sdk_callback.OnSignalingListener
}

func NewLiveSignaling(ws *ws.Ws, loginUserID string, platformID int32, db db_interface.DataBase) *LiveSignaling {
	if ws == nil {
		log.Warn("", " ws is nil")
		return nil
	}
	return &LiveSignaling{Ws: ws, loginUserID: loginUserID, platformID: platformID, DataBase: db}
}

func (s *LiveSignaling) setDefaultReq(req *sdkws.InvitationInfo) {
	if req.RoomID == "" {
		req.RoomID = utils.OperationIDGenerator()
	}
	if req.Timeout == 0 {
		req.Timeout = 60 * 60
	}
}

func (s *LiveSignaling) waitPush(ctx context.Context, req *sdkws.SignalReq, busyLineUserList []string) {
	var invt *sdkws.InvitationInfo
	switch payload := req.Payload.(type) {
	case *sdkws.SignalReq_Invite:
		invt = payload.Invite.Invitation
	case *sdkws.SignalReq_InviteInGroup:
		invt = payload.InviteInGroup.Invitation
	}
	var listenerList []open_im_sdk_callback.OnSignalingListener
	if s.listener != nil {
		listenerList = append(listenerList, s.listener)
		// log.Info(operationID, "listenerList ", listenerList, "listener ", s.listener)
	}
	if s.listenerForService != nil {
		listenerList = append(listenerList, s.listenerForService)
		// log.Info(operationID, "listenerList ", listenerList, "listenerForService ", s.listenerForService)
	}
	if len(listenerList) == 0 {
		// log.Error(operationID, "len (listenerList) == 0 ")
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
				if strings.Contains(err.Error(), "timeout") {
					// log.Error(operationID, "wait push timeout ", err.Error(), invt.InviterUserID, v, invt.RoomID, invt.Timeout)
					switch payload := req.Payload.(type) {
					case *sdkws.SignalReq_Invite:
						if !s.isCanceled {
							for _, listener := range listenerList {
								listener.OnInvitationTimeout(utils.StructToJsonString(payload.Invite))
								// log.Info(operationID, "OnInvitationTimeout ", utils.StructToJsonString(payload.Invite), listener)
							}
						}
					case *sdkws.SignalReq_InviteInGroup:
						if !s.isCanceled {
							for _, listener := range listenerList {
								listener.OnInvitationTimeout(utils.StructToJsonString(payload.InviteInGroup))
								// log.Info(operationID, "OnInvitationTimeout ", utils.StructToJsonString(payload.InviteInGroup), listener)
							}
						}
					}

				} else {
					// log.Error(operationID, "other failed ", err.Error(), invt.InviterUserID, v, invt.RoomID, invt.Timeout)
				}
				return
			}
			// log.Info(operationID, "SignalingWaitPush ", push.String(), invt.InviterUserID, v, invt.RoomID, invt.Timeout)
			s.doSignalPush(ctx, push)
		}(v)
	}
}
func (s *LiveSignaling) doSignalPush(ctx context.Context, req *api.SignalReq) {
	var listenerList []open_im_sdk_callback.OnSignalingListener
	if s.listener != nil {
		listenerList = append(listenerList, s.listener)
		// log.Info(operationID, "listenerList ", listenerList, "listener ", s.listener)
	}
	if s.listenerForService != nil {
		listenerList = append(listenerList, s.listenerForService)
		// log.Info(operationID, "listenerList ", listenerList, "listenerForService ", s.listenerForService)
	}
	if len(listenerList) == 0 {
		// log.Error(operationID, "len (listenerList) == 0 ")
		return
	}
	switch payload := req.Payload.(type) {
	case *api.SignalReq_Accept:
		// log.Info(operationID, "recv signal push Accept ", payload.Accept.String())
		for _, listener := range listenerList {
			listener.OnInviteeAccepted(utils.StructToJsonString(payload.Accept))
			// log.Info(operationID, "OnInviteeAccepted ", utils.StructToJsonString(payload.Accept), listener)
		}

	case *api.SignalReq_Reject:
		// log.Info(operationID, "recv signal push Reject ", payload.Reject.String())
		for _, listener := range listenerList {
			listener.OnInviteeRejected(utils.StructToJsonString(payload.Reject))
			// log.Info(operationID, "OnInviteeRejected ", utils.StructToJsonString(payload.Reject), listener)
		}
	default:
		// log.Error(operationID, "payload type failed ", payload)
	}
}

func (s *LiveSignaling) SetListener(listener open_im_sdk_callback.OnSignalingListener) {
	s.listener = listener
}

func (s *LiveSignaling) SetListenerForService(listener open_im_sdk_callback.OnSignalingListener) {
	s.listenerForService = listener
}

func (s *LiveSignaling) getSelfParticipant(ctx context.Context, groupID string) (*sdkws.ParticipantMetaData, error) {
	p := sdkws.ParticipantMetaData{GroupInfo: &sdkws.GroupInfo{}, GroupMemberInfo: &sdkws.GroupMemberFullInfo{}, UserInfo: &sdkws.PublicUserInfo{}}
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
			var wsResp ws.GeneralWsResp
			wsResp.ReqIdentifier = constant.WSSendSignalMsg
			wsResp.Data = msg.Content
			wsResp.MsgIncr = s.loginUserID + payload.Accept.OpUserID + payload.Accept.Invitation.RoomID
			log.Info(operationID, "search msgIncr: ", wsResp.MsgIncr)
			s.DoWSSignal(wsResp)
			return
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
			var wsResp ws.GeneralWsResp
			wsResp.ReqIdentifier = constant.WSSendSignalMsg
			wsResp.Data = msg.Content
			wsResp.MsgIncr = s.loginUserID + payload.Reject.OpUserID + payload.Reject.Invitation.RoomID
			log.Info(operationID, "search msgIncr: ", wsResp.MsgIncr)
			s.DoWSSignal(wsResp)
			return
		}
		if payload.Reject.OpUserPlatformID != s.platformID && payload.Reject.OpUserID == s.loginUserID {
			for _, listener := range listenerList {
				listener.OnInviteeRejectedByOtherDevice(utils.StructToJsonString(payload.Reject))
				log.Info(operationID, "OnInviteeRejectedByOtherDevice ", utils.StructToJsonString(payload.Reject), listener)
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
