package group

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
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
	dbGroup, err := g.db.GetGroupInfoByGroupID(ctx, groupID)
	if err != nil {
		// todo error
	}
	// todo
	if err := g.groupSyncer.Sync(ctx, util.AppendNotNil(ServerGroupToLocalGroup(resp.GroupInfos[0])), util.AppendNotNil(dbGroup), nil); err != nil {
		return err
	}
	//g.listener.OnGroupInfoChanged(utils.StructToJsonString(serverData))
	return nil
}

func (g *Group) SyncGroupAndMember(ctx context.Context, groupID string) error {
	if err := g.SyncGroup(ctx, groupID); err != nil {
		return err
	}
	return g.SyncGroupMember(ctx, groupID, nil)
}

func (g *Group) SyncSelfGroupApplication(ctx context.Context) error {
	fn := func(resp *group.GetGroupApplicationListResp) []*sdkws.GroupRequest { return resp.GroupRequests }
	req := &group.GetUserReqApplicationListReq{UserID: g.loginUserID, Pagination: &sdkws.RequestPagination{}}
	list, err := util.GetPageAll(ctx, constant.GetSendGroupApplicationListRouter, req, fn)
	if err != nil {
		return err
	}
	localData, err := g.db.GetSendGroupApplication(ctx)
	if err != nil {
		return err
	}
	if err := g.groupRequestSyncer.Sync(ctx, util.Batch(ServerGroupRequestToLocalGroupRequest, list), localData, nil); err != nil {
		return err
	}
	// todo
	return nil
}
