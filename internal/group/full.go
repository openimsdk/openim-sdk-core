package group

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/utils/datautil"
)

func (g *Group) GetGroupInfoFromLocal2Svr(ctx context.Context, groupID string) (*model_struct.LocalGroup, error) {
	localGroup, err := g.db.GetGroupInfoByGroupID(ctx, groupID)
	if err == nil {
		return localGroup, nil
	}
	svrGroup, err := g.getGroupsInfoFromSvr(ctx, []string{groupID})
	if err != nil {
		return nil, err
	}
	if len(svrGroup) == 0 {
		return nil, sdkerrs.ErrGroupIDNotFound.WrapMsg("server not this group")
	}
	return ServerGroupToLocalGroup(svrGroup[0]), nil
}

func (g *Group) GetGroupsInfoFromLocal2Svr(ctx context.Context, groupIDs ...string) (map[string]*model_struct.LocalGroup, error) {
	groupMap := make(map[string]*model_struct.LocalGroup)
	if len(groupIDs) == 0 {
		return groupMap, nil
	}
	groups, err := g.db.GetGroups(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	var groupIDsNeedSync []string
	localGroupIDs := datautil.Slice(groups, func(group *model_struct.LocalGroup) string {
		return group.GroupID
	})
	for _, groupID := range groupIDs {
		if !datautil.Contain(groupID, localGroupIDs...) {
			groupIDsNeedSync = append(groupIDsNeedSync, groupID)
		}
	}

	if len(groupIDsNeedSync) > 0 {
		svrGroups, err := g.getGroupsInfoFromSvr(ctx, groupIDsNeedSync)
		if err != nil {
			return nil, err
		}
		for _, svrGroup := range svrGroups {
			groups = append(groups, ServerGroupToLocalGroup(svrGroup))
		}
	}
	for _, g := range groups {
		groupMap[g.GroupID] = g
	}
	return groupMap, nil
}

func (g *Group) getGroupsInfoFromSvr(ctx context.Context, groupIDs []string) ([]*sdkws.GroupInfo, error) {
	resp, err := util.CallApi[group.GetGroupsInfoResp](ctx, constant.GetGroupsInfoRouter, &group.GetGroupsInfoReq{GroupIDs: groupIDs})
	if err != nil {
		return nil, err
	}
	return resp.GroupInfos, nil
}

func (g *Group) GetJoinedDiffusionGroupIDListFromSvr(ctx context.Context) ([]string, error) {
	groups, err := g.GetServerJoinGroup(ctx)
	if err != nil {
		return nil, err
	}
	var groupIDs []string
	for _, g := range groups {
		if g.GroupType == constant.WorkingGroup {
			groupIDs = append(groupIDs, g.GroupID)
		}
	}
	return groupIDs, nil
}
