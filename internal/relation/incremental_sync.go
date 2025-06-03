package relation

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/syncer"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/tools/utils/datautil"
)

func (r *Relation) IncrSyncFriends(ctx context.Context) error {
	friendSyncer := syncer.VersionSynchronizer[*model_struct.LocalFriend, *relation.GetIncrementalFriendsResp]{
		Ctx:       ctx,
		DB:        r.db,
		TableName: r.friendListTableName(),
		EntityID:  r.loginUserID,
		Key: func(localFriend *model_struct.LocalFriend) string {
			return localFriend.FriendUserID
		},
		Local: func() ([]*model_struct.LocalFriend, error) {
			return r.db.GetAllFriendList(ctx)
		},
		Server: func(version *model_struct.LocalVersionSync) (*relation.GetIncrementalFriendsResp, error) {
			return r.getIncrementalFriends(ctx, &relation.GetIncrementalFriendsReq{
				UserID:    r.loginUserID,
				Version:   version.Version,
				VersionID: version.VersionID,
			})
		},
		Full: func(resp *relation.GetIncrementalFriendsResp) bool {
			return resp.Full
		},
		Version: func(resp *relation.GetIncrementalFriendsResp) (string, uint64) {
			return resp.VersionID, resp.Version
		},
		Delete: func(resp *relation.GetIncrementalFriendsResp) []string {
			return resp.Delete
		},
		Update: func(resp *relation.GetIncrementalFriendsResp) []*model_struct.LocalFriend {
			return datautil.Batch(ServerFriendToLocalFriend, resp.Update)
		},
		Insert: func(resp *relation.GetIncrementalFriendsResp) []*model_struct.LocalFriend {
			return datautil.Batch(ServerFriendToLocalFriend, resp.Insert)
		},
		Syncer: func(server, local []*model_struct.LocalFriend) error {
			return r.friendSyncer.Sync(ctx, server, local, nil)
		},
		FullSyncer: func(ctx context.Context) error {
			return r.friendSyncer.FullSync(ctx, r.loginUserID)
		},
		FullID: func(ctx context.Context) ([]string, error) {
			resp, err := r.getFullFriendUserIDs(ctx, &relation.GetFullFriendUserIDsReq{
				UserID: r.loginUserID,
			})
			if err != nil {
				return nil, err
			}
			return resp.UserIDs, nil
		},
		IDOrderChanged: func(resp *relation.GetIncrementalFriendsResp) bool {
			return resp.SortVersion > 0
		},
	}
	return friendSyncer.IncrementalSync()
}

func (r *Relation) friendListTableName() string {
	return model_struct.LocalFriend{}.TableName()
}
func (r *Relation) IncrSyncFriendsWithLock(ctx context.Context) error {
	r.relationSyncMutex.Lock()
	defer r.relationSyncMutex.Unlock()
	return r.IncrSyncFriends(ctx)
}
