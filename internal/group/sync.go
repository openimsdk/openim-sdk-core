package group

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
)

func (g *Group) SyncGroupMember(ctx context.Context, groupID string) error {
	members, err := g.GetServerGroupMembers(ctx, groupID)
	if err != nil {
		return err
	}
	localData, err := g.db.GetGroupMemberListSplit(ctx, groupID, 0, 0, 9999999)
	if err != nil {
		return err
	}
	return g.groupMemberSyncer.Sync(ctx, util.Batch(ServerGroupMemberToLocalGroupMember, members), localData, nil)
}

func (g *Group) SyncJoinedGroup(ctx context.Context) error {
	_, err := g.syncJoinedGroup(ctx)
	if err != nil {
		return err
	}
	return err
}

func (g *Group) syncJoinedGroup(ctx context.Context) ([]*sdkws.GroupInfo, error) {
	groups, err := g.GetServerJoinGroup(ctx)
	if err != nil {
		return nil, err
	}
	localData, err := g.db.GetJoinedGroupListDB(ctx)
	if err != nil {
		return nil, err
	}
	if err := g.groupSyncer.Sync(ctx, util.Batch(ServerGroupToLocalGroup, groups), localData, nil); err != nil {
		return nil, err
	}
	return groups, nil
}

func (g *Group) SyncSelfGroupApplication(ctx context.Context) error {
	list, err := g.GetServerSelfGroupApplication(ctx)
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

func (g *Group) SyncJoinedGroupList(ctx context.Context) {
	if err := g.SyncJoinedGroup(ctx); err != nil {
		// tood log
	}
}

func (g *Group) SyncAdminGroupApplication(ctx context.Context) error {
	requests, err := g.GetServerAdminGroupApplicationList(ctx)
	if err != nil {
		return err
	}
	localData, err := g.db.GetAdminGroupApplication(ctx)
	if err != nil {
		return err
	}
	return g.groupAdminRequestSyncer.Sync(ctx, util.Batch(ServerGroupRequestToLocalAdminGroupRequest, requests), localData, nil)
}

func (g *Group) GetServerJoinGroup(ctx context.Context) ([]*sdkws.GroupInfo, error) {
	fn := func(resp *group.GetJoinedGroupListResp) []*sdkws.GroupInfo { return resp.Groups }
	req := &group.GetJoinedGroupListReq{FromUserID: g.loginUserID, Pagination: &sdkws.RequestPagination{}}
	return util.GetPageAll(ctx, constant.GetJoinedGroupListRouter, req, fn)
}

func (g *Group) GetServerAdminGroupApplicationList(ctx context.Context) ([]*sdkws.GroupRequest, error) {
	fn := func(resp *group.GetGroupApplicationListResp) []*sdkws.GroupRequest { return resp.GroupRequests }
	req := &group.GetGroupApplicationListReq{FromUserID: g.loginUserID, Pagination: &sdkws.RequestPagination{}}
	return util.GetPageAll(ctx, constant.GetRecvGroupApplicationListRouter, req, fn)
}

func (g *Group) GetServerSelfGroupApplication(ctx context.Context) ([]*sdkws.GroupRequest, error) {
	fn := func(resp *group.GetGroupApplicationListResp) []*sdkws.GroupRequest { return resp.GroupRequests }
	req := &group.GetUserReqApplicationListReq{UserID: g.loginUserID, Pagination: &sdkws.RequestPagination{}}
	return util.GetPageAll(ctx, constant.GetSendGroupApplicationListRouter, req, fn)
}

func (g *Group) GetServerGroupMembers(ctx context.Context, groupID string) ([]*sdkws.GroupMemberFullInfo, error) {
	req := &group.GetGroupMemberListReq{GroupID: groupID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *group.GetGroupMemberListResp) []*sdkws.GroupMemberFullInfo { return resp.Members }
	return util.GetPageAll(ctx, constant.GetGroupMemberListRouter, req, fn)
}
