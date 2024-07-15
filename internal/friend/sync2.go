package friend

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/internal/incrversion"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	friend "github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/tools/utils/datautil"
)

func (f *Friend) IncrSyncFriends(ctx context.Context) error {
	friendSyncer := incrversion.VersionSynchronizer[*model_struct.LocalFriend, *friend.GetIncrementalFriendsResp]{
		Ctx:       ctx,
		DB:        f.db,
		TableName: f.friendListTableName(),
		EntityID:  f.loginUserID,
		Key: func(localFriend *model_struct.LocalFriend) string {
			return localFriend.FriendUserID
		},
		Local: func() ([]*model_struct.LocalFriend, error) {
			return f.db.GetAllFriendList(ctx)
		},
		Server: func(version *model_struct.LocalVersionSync) (*friend.GetIncrementalFriendsResp, error) {
			return util.CallApi[friend.GetIncrementalFriendsResp](ctx, constant.GetIncrementalFriends, &friend.GetIncrementalFriendsReq{
				UserID:    f.loginUserID,
				Version:   version.Version,
				VersionID: version.VersionID,
			})
		},
		Full: func(resp *friend.GetIncrementalFriendsResp) bool {
			return resp.Full
		},
		Version: func(resp *friend.GetIncrementalFriendsResp) (string, uint64) {
			return resp.VersionID, resp.Version
		},
		Delete: func(resp *friend.GetIncrementalFriendsResp) []string {
			return resp.Delete
		},
		Update: func(resp *friend.GetIncrementalFriendsResp) []*model_struct.LocalFriend {
			return datautil.Batch(ServerFriendToLocalFriend, resp.Update)
		},
		Insert: func(resp *friend.GetIncrementalFriendsResp) []*model_struct.LocalFriend {
			return datautil.Batch(ServerFriendToLocalFriend, resp.Insert)
		},
		Syncer: func(server, local []*model_struct.LocalFriend) error {
			return f.friendSyncer.Sync(ctx, server, local, nil)
		},
		FullSyncer: func(ctx context.Context) error {
			return f.friendSyncer.FullSync(ctx, f.loginUserID)
		},
		FullID: func(ctx context.Context) ([]string, error) {
			resp, err := util.CallApi[friend.GetFullFriendUserIDsResp](ctx, constant.GetFullFriendUserIDs, &friend.GetFullFriendUserIDsReq{
				UserID: f.loginUserID,
			})
			if err != nil {
				return nil, err
			}
			return resp.UserIDs, nil
		},
	}
	return friendSyncer.Sync()
}

func (f *Friend) friendListTableName() string {
	return model_struct.LocalFriend{}.TableName()
}
