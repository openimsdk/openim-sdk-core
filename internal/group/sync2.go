package group

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/internal/incrversion"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	pconstant "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
)

func (g *Group) getIncrementalGroupMemberBatch(ctx context.Context, groups []*group.GetIncrementalGroupMemberReq) (map[string]*group.GetIncrementalGroupMemberResp, error) {
	type BatchIncrementalReq struct {
		UserID string                                `json:"user_id"`
		List   []*group.GetIncrementalGroupMemberReq `json:"list"`
	}
	type BatchIncrementalResp struct {
		List map[string]*group.GetIncrementalGroupMemberResp `json:"list"`
	}
	resp, err := util.CallApi[BatchIncrementalResp](ctx, constant.GetIncrementalGroupMemberBatch, &BatchIncrementalReq{UserID: g.loginUserID, List: groups})
	if err != nil {
		return nil, err
	}
	return resp.List, nil
}

func (g *Group) groupMemberVersionKey(groupID string) string {
	return "friend:" + groupID
}

func (g *Group) IncrSyncJoinGroupMember(ctx context.Context) error {
	groups, err := g.db.GetJoinedGroupListDB(ctx)
	if err != nil {
		return err
	}
	groupIDs := datautil.Slice(groups, func(e *model_struct.LocalGroup) string {
		return e.GroupID
	})
	return g.IncrSyncGroupMember(ctx, groupIDs...)
}

func (g *Group) IncrSyncGroupMember(ctx context.Context, groupIDs ...string) error {
	if len(groupIDs) == 0 {
		return nil
	}
	const maxSyncNum = pconstant.MaxSyncPullNumber
	groupIDSet := datautil.SliceSet(groupIDs)
	var groups []*group.GetIncrementalGroupMemberReq
	if len(groupIDs) > maxSyncNum {
		groups = make([]*group.GetIncrementalGroupMemberReq, 0, maxSyncNum)
	} else {
		groups = make([]*group.GetIncrementalGroupMemberReq, 0, len(groupIDs))
	}
	for {
		if len(groupIDSet) == 0 {
			return nil
		}
		for groupID := range groupIDSet {
			if len(groups) == cap(groups) {
				break
			}
			req := group.GetIncrementalGroupMemberReq{
				GroupID: groupID,
			}
			lvs, err := g.db.GetVersionSync(ctx, g.groupMemberVersionKey(groupID))
			if err == nil {
				req.VersionID = lvs.VersionID
				req.Version = lvs.Version
			} else {
				log.ZInfo(ctx, "get group version", "groupID", groupID, "error", err)
			}
			groups = append(groups, &req)
		}
		groupVersion, err := g.getIncrementalGroupMemberBatch(ctx, groups)
		if err != nil {
			return err
		}
		groups = groups[:0]
		for groupID, resp := range groupVersion {
			if err := g.syncGroupMember(ctx, groupID, resp); err != nil {
				return err
			}
			delete(groupIDSet, groupID)
		}
		num := len(groupIDSet)
		_ = num
	}
}

func (g *Group) syncGroupMember(ctx context.Context, groupID string, resp *group.GetIncrementalGroupMemberResp) error {
	opt := incrversion.Option[*model_struct.LocalGroupMember, *group.GetIncrementalGroupMemberResp]{
		Ctx: ctx,
		DB:  g.db,
		Key: func(localFriend *model_struct.LocalGroupMember) string {
			return localFriend.UserID
		},
		SyncKey: func() string {
			return g.groupMemberVersionKey(groupID)
		},
		Local: func() ([]*model_struct.LocalGroupMember, error) {
			return g.db.GetGroupMemberListSplit(ctx, groupID, 0, 0, 9999999)
		},
		ServerVersion: func() *group.GetIncrementalGroupMemberResp {
			return resp
		},
		Full: func(resp *group.GetIncrementalGroupMemberResp) bool {
			return resp.Full
		},
		Version: func(resp *group.GetIncrementalGroupMemberResp) (string, uint64) {
			return resp.VersionID, resp.Version
		},
		DeleteIDs: func(resp *group.GetIncrementalGroupMemberResp) []string {
			return resp.DeleteUserIds
		},
		Changes: func(resp *group.GetIncrementalGroupMemberResp) []*model_struct.LocalGroupMember {
			return util.Batch(ServerGroupMemberToLocalGroupMember, resp.Changes)
		},
		Syncer: func(server, local []*model_struct.LocalGroupMember) error {
			return g.groupMemberSyncer.Sync(ctx, server, local, nil)
		},
	}
	return opt.Sync()
}

func (g *Group) groupJoinVersionKey() string {
	return "join_group:" + g.loginUserID
}

func (g *Group) IncrSyncJoinGroup(ctx context.Context) error {
	opt := incrversion.Option[*model_struct.LocalGroup, *group.GetIncrementalJoinGroupResp]{
		Ctx: ctx,
		DB:  g.db,
		Key: func(localFriend *model_struct.LocalGroup) string {
			return localFriend.GroupID
		},
		SyncKey: g.groupJoinVersionKey,
		Local: func() ([]*model_struct.LocalGroup, error) {
			return g.db.GetJoinedGroupListDB(ctx)
		},
		Server: func(version *model_struct.LocalVersionSync) (*group.GetIncrementalJoinGroupResp, error) {
			return util.CallApi[group.GetIncrementalJoinGroupResp](ctx, constant.GetIncrementalJoinGroup, &group.GetIncrementalJoinGroupReq{
				UserID:    g.loginUserID,
				Version:   version.Version,
				VersionID: version.VersionID,
			})
		},
		Full: func(resp *group.GetIncrementalJoinGroupResp) bool {
			return resp.Full
		},
		Version: func(resp *group.GetIncrementalJoinGroupResp) (string, uint64) {
			return resp.VersionID, resp.Version
		},
		DeleteIDs: func(resp *group.GetIncrementalJoinGroupResp) []string {
			return resp.DeleteGroupIds
		},
		Changes: func(resp *group.GetIncrementalJoinGroupResp) []*model_struct.LocalGroup {
			return util.Batch(ServerGroupToLocalGroup, resp.Changes)
		},
		Syncer: func(server, local []*model_struct.LocalGroup) error {
			return g.groupSyncer.Sync(ctx, server, local, nil)
		},
	}
	return opt.Sync()
}
