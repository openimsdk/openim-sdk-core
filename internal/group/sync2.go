package group

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/internal/incrversion"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	pconstant "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"sync"
)

type BatchIncrementalReq struct {
	UserID string                                `json:"user_id"`
	List   []*group.GetIncrementalGroupMemberReq `json:"list"`
}
type BatchIncrementalResp struct {
	List map[string]*group.GetIncrementalGroupMemberResp `json:"list"`
}

func (g *Group) getIncrementalGroupMemberBatch(ctx context.Context, groups []*group.GetIncrementalGroupMemberReq) (map[string]*group.GetIncrementalGroupMemberResp, error) {
	resp, err := util.CallApi[BatchIncrementalResp](ctx, constant.GetIncrementalGroupMemberBatch, &BatchIncrementalReq{UserID: g.loginUserID, List: groups})
	if err != nil {
		return nil, err
	}
	return resp.List, nil
}

func (g *Group) groupMemberTableName() string {
	return model_struct.LocalGroupMember{}.TableName()
}

func (g *Group) groupTableName() string {
	return model_struct.LocalGroup{}.TableName()
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
	var wg sync.WaitGroup
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
			lvs, err := g.db.GetVersionSync(ctx, g.groupMemberTableName(), groupID)
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
			tempResp := resp
			tempGroupID := groupID
			wg.Add(1)
			go func() error {
				if err := g.syncGroupMember(ctx, tempGroupID, tempResp); err != nil {
					return err
				}
				wg.Done()
				return nil
			}()
			delete(groupIDSet, tempGroupID)
		}
		wg.Wait()
		num := len(groupIDSet)
		_ = num
	}
}

func (g *Group) syncGroupMember(ctx context.Context, groupID string, resp *group.GetIncrementalGroupMemberResp) error {
	groupMemberSyncer := incrversion.VersionSynchronizer[*model_struct.LocalGroupMember, *group.GetIncrementalGroupMemberResp]{
		Ctx:      ctx,
		DB:       g.db,
		TabName:  g.groupMemberTableName(),
		EntityID: groupID,
		Key: func(localGroupMember *model_struct.LocalGroupMember) string {
			return localGroupMember.UserID
		},
		Local: func() ([]*model_struct.LocalGroupMember, error) {
			return g.db.GetGroupMemberListByGroupID(ctx, groupID)
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
		Delete: func(resp *group.GetIncrementalGroupMemberResp) []string {
			return resp.Delete
		},
		Update: func(resp *group.GetIncrementalGroupMemberResp) []*model_struct.LocalGroupMember {
			return datautil.Batch(ServerGroupMemberToLocalGroupMember, resp.Update)
		},
		Insert: func(resp *group.GetIncrementalGroupMemberResp) []*model_struct.LocalGroupMember {
			return datautil.Batch(ServerGroupMemberToLocalGroupMember, resp.Insert)
		},
		Syncer: func(server, local []*model_struct.LocalGroupMember) error {
			return g.groupMemberSyncer.Sync(ctx, server, local, nil)
		},
		FullSyncer: func(ctx context.Context) error {
			return g.groupMemberSyncer.FullSync(ctx, groupID)
		},
		FullID: func(ctx context.Context) ([]string, error) {
			resp, err := util.CallApi[group.GetIncrementalGroupMemberUserIDsResp](ctx, constant.GetGroupMemberAllIDs, &group.GetIncrementalGroupMemberUserIDsReq{
				GroupID: groupID,
			})
			if err != nil {
				return nil, err
			}
			return resp.UserIDs, nil
		},
	}
	return groupMemberSyncer.Sync()
}

func (g *Group) onlineSyncGroupMember(ctx context.Context, groupID string, delete, update, insert []*sdkws.GroupMemberFullInfo, version uint64) error {
	groupMemberSyncer := incrversion.VersionSynchronizer[*model_struct.LocalGroupMember, *group.GetIncrementalGroupMemberResp]{
		Ctx:      ctx,
		DB:       g.db,
		TabName:  g.groupMemberTableName(),
		EntityID: groupID,
		Key: func(localGroupMember *model_struct.LocalGroupMember) string {
			return localGroupMember.UserID
		},
		Local: func() ([]*model_struct.LocalGroupMember, error) {
			return g.db.GetGroupMemberListByGroupID(ctx, groupID)
		},
		ServerVersion: func() *group.GetIncrementalGroupMemberResp {
			return &group.GetIncrementalGroupMemberResp{
				Version:   version,
				VersionID: "",
				Full:      false,
				Delete: datautil.Slice(delete, func(e *sdkws.GroupMemberFullInfo) string {
					return e.UserID
				}),
				Insert: insert,
				Update: update,
			}
		},
		Server: func(version *model_struct.LocalVersionSync) (*group.GetIncrementalGroupMemberResp, error) {
			singleGroupReq := &group.GetIncrementalGroupMemberReq{
				GroupID:   groupID,
				VersionID: version.VersionID,
				Version:   version.Version,
			}
			resp, err := util.CallApi[BatchIncrementalResp](ctx, constant.GetIncrementalGroupMemberBatch,
				&BatchIncrementalReq{UserID: g.loginUserID, List: []*group.GetIncrementalGroupMemberReq{singleGroupReq}})
			if err != nil {
				return nil, err
			}
			if resp.List != nil {
				if singleGroupResp, ok := resp.List[groupID]; ok {
					return singleGroupResp, nil
				}
			}
			return nil, errs.New("group member version record not found")

		},
		Full: func(resp *group.GetIncrementalGroupMemberResp) bool {
			return resp.Full
		},
		Version: func(resp *group.GetIncrementalGroupMemberResp) (string, uint64) {
			return resp.VersionID, resp.Version
		},
		Delete: func(resp *group.GetIncrementalGroupMemberResp) []string {
			return resp.Delete
		},
		Update: func(resp *group.GetIncrementalGroupMemberResp) []*model_struct.LocalGroupMember {
			return datautil.Batch(ServerGroupMemberToLocalGroupMember, resp.Update)
		},
		Insert: func(resp *group.GetIncrementalGroupMemberResp) []*model_struct.LocalGroupMember {
			return datautil.Batch(ServerGroupMemberToLocalGroupMember, resp.Insert)
		},
		Syncer: func(server, local []*model_struct.LocalGroupMember) error {
			return g.groupMemberSyncer.Sync(ctx, server, local, nil)
		},
		FullSyncer: func(ctx context.Context) error {
			return g.groupMemberSyncer.FullSync(ctx, groupID)
		},
		FullID: func(ctx context.Context) ([]string, error) {
			resp, err := util.CallApi[group.GetIncrementalGroupMemberUserIDsResp](ctx, constant.GetGroupMemberAllIDs, &group.GetIncrementalGroupMemberUserIDsReq{
				GroupID: groupID,
			})
			if err != nil {
				return nil, err
			}
			return resp.UserIDs, nil
		},
	}
	return groupMemberSyncer.CheckVersionSync()
}

func (g *Group) IncrSyncJoinGroup(ctx context.Context) error {
	opt := incrversion.VersionSynchronizer[*model_struct.LocalGroup, *group.GetIncrementalJoinGroupResp]{
		Ctx:      ctx,
		DB:       g.db,
		TabName:  g.groupTableName(),
		EntityID: g.loginUserID,
		Key: func(LocalGroup *model_struct.LocalGroup) string {
			return LocalGroup.GroupID
		},
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
		Delete: func(resp *group.GetIncrementalJoinGroupResp) []string {
			return resp.Delete
		},
		Update: func(resp *group.GetIncrementalJoinGroupResp) []*model_struct.LocalGroup {
			return util.Batch(ServerGroupToLocalGroup, resp.Update)
		},
		Insert: func(resp *group.GetIncrementalJoinGroupResp) []*model_struct.LocalGroup {
			return datautil.Batch(ServerGroupToLocalGroup, resp.Insert)
		},
		Syncer: func(server, local []*model_struct.LocalGroup) error {
			return g.groupSyncer.Sync(ctx, server, local, nil)
		},
		FullSyncer: func(ctx context.Context) error {
			return g.groupSyncer.FullSync(ctx, g.loginUserID)
		},
		FullID: func(ctx context.Context) ([]string, error) {
			resp, err := util.CallApi[group.GetIncrementalJoinGroupIDsResp](ctx, constant.GetJoinedGroupAllIDs, &group.GetIncrementalJoinGroupIDsReq{
				UserID: g.loginUserID,
			})
			if err != nil {
				return nil, err
			}
			return resp.GroupIDs, nil

		},
	}
	return opt.Sync()
}
