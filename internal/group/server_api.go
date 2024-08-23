package group

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
)

func (g *Group) getFullGroupMemberUserIDs(ctx context.Context, req *group.GetFullGroupMemberUserIDsReq) (*group.GetFullGroupMemberUserIDsResp, error) {
	return api.GetFullGroupMemberUserIDs(ctx, req)
}

func (g *Group) getIncrementalJoinGroup(ctx context.Context, req *group.GetIncrementalJoinGroupReq) (*group.GetIncrementalJoinGroupResp, error) {
	return api.GetIncrementalJoinGroup(ctx, req)
}

func (g *Group) getFullJoinGroupIDs(ctx context.Context, req *group.GetFullJoinGroupIDsReq) (*group.GetFullJoinGroupIDsResp, error) {
	return api.GetFullJoinedGroupIDs(ctx, req)
}

func (g *Group) getIncrementalGroupMemberBatch(ctx context.Context, reqs []*group.GetIncrementalGroupMemberReq) (map[string]*group.GetIncrementalGroupMemberResp, error) {
	resp, err := api.GetIncrementalGroupMemberBatch(ctx, &group.BatchGetIncrementalGroupMemberReq{UserID: g.loginUserID, ReqList: reqs})
	if err != nil {
		return nil, err
	}
	return resp.RespList, nil
}

func (g *Group) createGroup(ctx context.Context, req *group.CreateGroupReq) (*group.CreateGroupResp, error) {
	return api.CreateGroup(ctx, req)
}

func (g *Group) joinGroup(ctx context.Context, req *group.JoinGroupReq) error {
	_, err := api.JoinGroup(ctx, req)
	return err
}

func (g *Group) quitGroup(ctx context.Context, groupID string) error {
	_, err := api.QuitGroup(ctx, &group.QuitGroupReq{GroupID: groupID, UserID: g.loginUserID})
	return err
}

func (g *Group) dismissGroup(ctx context.Context, groupID string) error {
	_, err := api.DismissGroup(ctx, &group.DismissGroupReq{GroupID: groupID})
	return err
}

func (g *Group) setGroupInfo(ctx context.Context, req *group.SetGroupInfoReq) error {
	_, err := api.SetGroupInfo(ctx, req)
	return err
}

func (g *Group) setGroupMemberInfo(ctx context.Context, req *group.SetGroupMemberInfoReq) error {
	_, err := api.SetGroupMemberInfo(ctx, req)
	return err
}

func (g *Group) kickGroupMember(ctx context.Context, req *group.KickGroupMemberReq) error {
	_, err := api.KickGroupMember(ctx, req)
	return err
}

func (g *Group) transferGroup(ctx context.Context, req *group.TransferGroupOwnerReq) error {
	_, err := api.TransferGroup(ctx, req)
	return err
}

func (g *Group) cancelMuteGroupMember(ctx context.Context, req *group.CancelMuteGroupMemberReq) error {
	_, err := api.CancelMuteGroupMember(ctx, req)
	return err
}

func (g *Group) muteGroupMember(ctx context.Context, req *group.MuteGroupMemberReq) error {
	_, err := api.MuteGroupMember(ctx, req)
	return err
}

func (g *Group) muteGroup(ctx context.Context, groupID string) error {
	_, err := api.MuteGroup(ctx, &group.MuteGroupReq{GroupID: groupID})
	return err
}

func (g *Group) cancelMuteGroup(ctx context.Context, groupID string) error {
	_, err := api.CancelMuteGroup(ctx, &group.CancelMuteGroupReq{GroupID: groupID})
	return err
}

func (g *Group) getDesignatedGroupMembers(ctx context.Context, groupID string, userIDs []string) ([]*sdkws.GroupMemberFullInfo, error) {
	resp, err := api.GetGroupMembersInfo(ctx, &group.GetGroupMembersInfoReq{GroupID: groupID, UserIDs: userIDs})
	if err != nil {
		return nil, err
	}
	return resp.Members, nil
}

func (g *Group) getServerSelfGroupApplication(ctx context.Context) ([]*sdkws.GroupRequest, error) {
	req := &group.GetUserReqApplicationListReq{UserID: g.loginUserID, Pagination: &sdkws.RequestPagination{}}
	return util.PageNext(ctx, req, api.GetSendGroupApplicationList, func(resp *group.GetUserReqApplicationListResp) []*sdkws.GroupRequest { return resp.GroupRequests })
}

func (g *Group) getServerJoinGroup(ctx context.Context) ([]*sdkws.GroupInfo, error) {
	req := &group.GetJoinedGroupListReq{FromUserID: g.loginUserID, Pagination: &sdkws.RequestPagination{}}
	return util.PageNext(ctx, req, api.GetJoinedGroupList, func(resp *group.GetJoinedGroupListResp) []*sdkws.GroupInfo { return resp.Groups })
}

func (g *Group) getServerAdminGroupApplicationList(ctx context.Context) ([]*sdkws.GroupRequest, error) {
	req := &group.GetGroupApplicationListReq{FromUserID: g.loginUserID, Pagination: &sdkws.RequestPagination{}}
	return util.PageNext(ctx, req, api.GetRecvGroupApplicationList, func(resp *group.GetGroupApplicationListResp) []*sdkws.GroupRequest { return resp.GroupRequests })
}

func (g *Group) getGroupsInfoFromSvr(ctx context.Context, groupIDs []string) ([]*sdkws.GroupInfo, error) {
	resp, err := api.GetGroupsInfo(ctx, &group.GetGroupsInfoReq{GroupIDs: groupIDs})
	if err != nil {
		return nil, err
	}
	return resp.GroupInfos, nil
}
