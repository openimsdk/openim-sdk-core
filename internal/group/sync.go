package group

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/syncer"
	"open_im_sdk/pkg/utils"
)

func (g *Group) SyncGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	resp, err := util.CallApi[group.GetGroupMembersInfoResp](ctx, constant.GetGroupMembersInfoRouter, &group.GetGroupMembersInfoReq{GroupID: groupID, UserIDs: userIDs})
	if err != nil {
		return err
	}
	var members []any
	for _, member := range resp.Members {
		members = append(members, &model_struct.LocalGroupMember{
			GroupID:        member.GroupID,
			UserID:         member.UserID,
			Nickname:       member.Nickname,
			FaceURL:        member.FaceURL,
			RoleLevel:      member.RoleLevel,
			JoinTime:       member.JoinTime,
			JoinSource:     member.JoinSource,
			InviterUserID:  member.InviterUserID,
			MuteEndTime:    member.MuteEndTime,
			OperatorUserID: member.OperatorUserID,
			Ex:             member.Ex,
			//AttachedInfo:   member.AttachedInfo, // todo
		})
	}
	return syncer.New(nil).AddLocally([]any{members}).Start()
}

func (g *Group) SyncGroup(ctx context.Context, groupID string) error {
	resp, err := util.CallApi[group.GetGroupsInfoResp](ctx, constant.GetGroupsInfoRouter, &group.GetGroupsInfoReq{GroupIDs: []string{groupID}})
	if err != nil {
		return err
	}
	if len(resp.GroupInfos) == 0 {
		return errs.ErrGroupIDNotFound.Wrap(groupID)
	}
	groupInfo := resp.GroupInfos[0]
	groupModel := &model_struct.LocalGroup{
		GroupID:                groupInfo.GroupID,
		GroupName:              groupInfo.GroupName,
		Notification:           groupInfo.Notification,
		Introduction:           groupInfo.Introduction,
		FaceURL:                groupInfo.FaceURL,
		CreateTime:             groupInfo.CreateTime,
		Status:                 groupInfo.Status,
		CreatorUserID:          groupInfo.CreatorUserID,
		GroupType:              groupInfo.GroupType,
		OwnerUserID:            groupInfo.OwnerUserID,
		MemberCount:            int32(groupInfo.MemberCount),
		Ex:                     groupInfo.Ex,
		NeedVerification:       groupInfo.NeedVerification,
		LookMemberInfo:         groupInfo.LookMemberInfo,
		ApplyMemberFriend:      groupInfo.ApplyMemberFriend,
		NotificationUpdateTime: groupInfo.NotificationUpdateTime,
		NotificationUserID:     groupInfo.NotificationUserID,
		//AttachedInfo:           groupInfo.AttachedInfo, // TODO
	}
	if err := syncer.New(nil).AddLocally([]any{groupModel}).Start(); err != nil {
		return err
	}
	g.listener.OnGroupInfoChanged(utils.StructToJsonString(groupModel))
	return nil
}

func (g *Group) SyncGroupAndMember(ctx context.Context, groupID string) error {
	groupResp, err := util.CallApi[group.GetGroupsInfoResp](ctx, constant.GetGroupsInfoRouter, &group.GetGroupsInfoReq{GroupIDs: []string{groupID}})
	if err != nil {
		return err
	}
	if len(groupResp.GroupInfos) == 0 {
		return errs.ErrGroupIDNotFound.Wrap(groupID)
	}
	groupInfo := groupResp.GroupInfos[0]
	req := &group.GetGroupMemberListReq{GroupID: groupInfo.GroupID, Pagination: &sdkws.RequestPagination{PageNumber: 0, ShowNumber: 20}}
	members, err := util.GetPageAll(ctx, constant.GetGroupAllMemberListRouter, req, func(resp *group.GetGroupMemberListResp) []*sdkws.GroupMemberFullInfo { return resp.Members })
	if err != nil {
		return err
	}
	groupModel := &model_struct.LocalGroup{
		GroupID:       groupInfo.GroupID,
		GroupName:     groupInfo.GroupName,
		Notification:  groupInfo.Notification,
		Introduction:  groupInfo.Introduction,
		FaceURL:       groupInfo.FaceURL,
		CreateTime:    groupInfo.CreateTime,
		Status:        groupInfo.Status,
		CreatorUserID: groupInfo.CreatorUserID,
		GroupType:     groupInfo.GroupType,
		OwnerUserID:   groupInfo.OwnerUserID,
		MemberCount:   int32(groupInfo.MemberCount),
		Ex:            groupInfo.Ex,
		//AttachedInfo:           groupInfo.AttachedInfo, // TODO
		NeedVerification:       groupInfo.NeedVerification,
		LookMemberInfo:         groupInfo.LookMemberInfo,
		ApplyMemberFriend:      groupInfo.ApplyMemberFriend,
		NotificationUpdateTime: groupInfo.NotificationUpdateTime,
		NotificationUserID:     groupInfo.NotificationUserID,
	}
	s := syncer.New(nil).AddLocally([]any{groupModel}).AddGlobal(map[string]any{"group_id": groupID}, members)
	if err := s.Start(); err != nil {
		return err
	}
	g.listener.OnGroupInfoChanged(utils.StructToJsonString(groupModel))
	for _, member := range members {
		g.listener.OnGroupMemberInfoChanged(utils.StructToJsonString(member))
	}
	return nil
}

func (g *Group) SyncSelfGroupApplication(ctx context.Context) error {
	list, err := util.GetPageAll(ctx, constant.GetSendGroupApplicationListRouter, &group.GetUserReqApplicationListReq{UserID: g.loginUserID, Pagination: &sdkws.RequestPagination{}}, func(resp *group.GetGroupApplicationListResp) []*sdkws.GroupRequest { return resp.GroupRequests })
	if err != nil {
		return err
	}
	ms := make([]*model_struct.LocalGroupRequest, 0, len(list))
	for _, request := range list {
		ms = append(ms, &model_struct.LocalGroupRequest{
			GroupID:       request.GroupInfo.GroupID,
			GroupName:     request.GroupInfo.GroupName,
			Notification:  request.GroupInfo.Notification,
			Introduction:  request.GroupInfo.Introduction,
			GroupFaceURL:  request.GroupInfo.FaceURL,
			CreateTime:    request.GroupInfo.CreateTime,
			Status:        request.GroupInfo.Status,
			CreatorUserID: request.GroupInfo.CreatorUserID,
			GroupType:     request.GroupInfo.GroupType,
			OwnerUserID:   request.GroupInfo.OwnerUserID,
			MemberCount:   int32(request.GroupInfo.MemberCount),
			UserID:        request.UserInfo.UserID,
			Nickname:      request.UserInfo.Nickname,
			UserFaceURL:   request.UserInfo.FaceURL,
			Gender:        request.UserInfo.Gender,
			HandleResult:  request.HandleResult,
			ReqMsg:        request.ReqMsg,
			HandledMsg:    request.HandleMsg,
			ReqTime:       request.ReqTime,
			HandleUserID:  request.HandleUserID,
			HandledTime:   request.HandleTime,
			Ex:            request.Ex,
			//AttachedInfo:  request.AttachedInfo,
			JoinSource:    request.JoinSource,
			InviterUserID: request.InviterUserID,
		})
	}
	s := syncer.New(nil).AddGlobal(map[string]any{"user_id": g.loginUserID}, ms)
	if err := s.Start(); err != nil {
		return err
	}
	// todo
	return nil

	//svrList, err := g.getSendGroupApplicationListFromSvr(operationID)
	//if err != nil {
	//	log.NewError(operationID, "getSendGroupApplicationListFromSvr failed ", err.Error())
	//	return
	//}
	//onServer := common.TransferToLocalSendGroupRequest(svrList)
	//onLocal, err := g.db.GetSendGroupApplication()
	//if err != nil {
	//	log.NewError(operationID, "GetSendGroupApplication failed ", err.Error())
	//	return
	//}
	//
	//log.NewInfo(operationID, "svrList onServer onLocal ", svrList, onServer, onLocal)
	//aInBNot, bInANot, sameA, sameB := common.CheckGroupRequestDiff(onServer, onLocal)
	//log.Info(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	//for _, index := range aInBNot {
	//	err := g.db.InsertGroupRequest(onServer[index])
	//	if err != nil {
	//		log.NewError(operationID, "InsertGroupRequest failed ", err.Error(), *onServer[index])
	//		continue
	//	}
	//	callbackData := *onServer[index]
	//	if g.listener != nil {
	//		g.listener.OnGroupApplicationAdded(utils.StructToJsonString(callbackData))
	//		log.Info(operationID, "OnGroupApplicationAdded ", utils.StructToJsonString(callbackData))
	//	}
	//}
	//for _, index := range sameA {
	//	err := g.db.UpdateGroupRequest(onServer[index])
	//	if err != nil {
	//		log.NewError(operationID, "UpdateGroupRequest failed ", err.Error())
	//		continue
	//	}
	//	if onServer[index].HandleResult == constant.GroupResponseRefuse {
	//		callbackData := *onServer[index]
	//		if g.listener != nil {
	//			g.listener.OnGroupApplicationRejected(utils.StructToJsonString(callbackData))
	//			log.Info(operationID, "OnGroupApplicationRejected", utils.StructToJsonString(callbackData))
	//		}
	//
	//	} else if onServer[index].HandleResult == constant.GroupResponseAgree {
	//		callbackData := *onServer[index]
	//		if g.listener != nil {
	//			g.listener.OnGroupApplicationAccepted(utils.StructToJsonString(callbackData))
	//			log.Info(operationID, "OnGroupApplicationAccepted", utils.StructToJsonString(callbackData))
	//		}
	//		if g.listenerForService != nil {
	//			g.listenerForService.OnGroupApplicationAccepted(utils.StructToJsonString(callbackData))
	//			log.Info(operationID, "OnGroupApplicationAccepted", utils.StructToJsonString(callbackData))
	//		}
	//	} else {
	//		callbackData := *onServer[index]
	//		if g.listener != nil {
	//			g.listener.OnGroupApplicationAdded(utils.StructToJsonString(callbackData))
	//			log.Info(operationID, "OnGroupApplicationAdded", utils.StructToJsonString(callbackData))
	//		}
	//	}
	//}
	//for _, index := range bInANot {
	//	err := g.db.DeleteGroupRequest(onLocal[index].GroupID, onLocal[index].UserID)
	//	if err != nil {
	//		log.NewError(operationID, "DeleteGroupRequest failed ", err.Error())
	//		continue
	//	}
	//	callbackData := *onLocal[index]
	//	if g.listener != nil {
	//		g.listener.OnGroupApplicationDeleted(utils.StructToJsonString(callbackData))
	//	}
	//	log.Info(operationID, "OnGroupApplicationDeleted", utils.StructToJsonString(callbackData))
	//}
}
