package group

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (g *Group) SyncGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	var members []*sdkws.GroupMemberFullInfo
	if userIDs == nil {
		req := &group.GetGroupMemberListReq{GroupID: groupID, Pagination: &sdkws.RequestPagination{}}
		fn := func(resp *group.GetGroupMemberListResp) []*sdkws.GroupMemberFullInfo { return resp.Members }
		resp, err := util.GetPageAll(ctx, constant.GetGroupAllMemberListRouter, req, fn)
		if err != nil {
			return err
		}
		members = resp
	} else {
		resp, err := util.CallApi[group.GetGroupMembersInfoResp](ctx, constant.GetGroupMembersInfoRouter, &group.GetGroupMembersInfoReq{GroupID: groupID, UserIDs: userIDs})
		if err != nil {
			return err
		}
		members = resp.Members
	}
	var serverData []*model_struct.LocalGroupMember
	for _, member := range members {
		serverData = append(serverData, &model_struct.LocalGroupMember{
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
	localData, err := g.db.GetGroupMemberListSplit(ctx, groupID, 0, 0, 0)
	if err != nil {
		return err
	}
	return g.groupMemberSyncer.Sync(ctx, serverData, localData, nil)
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
	serverData := &model_struct.LocalGroup{
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
	localData := make([]*model_struct.LocalGroup, 0, 1)
	if dbGroup, err := g.db.GetGroupInfoByGroupID(ctx, groupID); err == nil {
		localData = append(localData, dbGroup)
	}
	// todo
	if err := g.groupSyncer.Sync(ctx, []*model_struct.LocalGroup{serverData}, localData, nil); err != nil {
		return err
	}
	g.listener.OnGroupInfoChanged(utils.StructToJsonString(serverData))
	return nil
}

func (g *Group) SyncGroupAndMember(ctx context.Context, groupID string) error {
	if err := g.SyncGroup(ctx, groupID); err != nil {
		return err
	}
	return g.SyncGroupMember(ctx, groupID, nil)
}

func (g *Group) SyncSelfGroupApplication(ctx context.Context) error {
	list, err := util.GetPageAll(ctx, constant.GetSendGroupApplicationListRouter, &group.GetUserReqApplicationListReq{UserID: g.loginUserID, Pagination: &sdkws.RequestPagination{}}, func(resp *group.GetGroupApplicationListResp) []*sdkws.GroupRequest { return resp.GroupRequests })
	if err != nil {
		return err
	}
	serverData := make([]*model_struct.LocalGroupRequest, 0, len(list))
	for _, request := range list {
		serverData = append(serverData, &model_struct.LocalGroupRequest{
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
	localData, err := g.db.GetSendGroupApplication(ctx)
	if err != nil {
		return err
	}
	if err := g.groupRequestSyncer.Sync(ctx, serverData, localData, nil); err != nil {
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
