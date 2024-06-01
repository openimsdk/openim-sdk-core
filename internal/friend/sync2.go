package friend

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/protocol/friend"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"time"
)

func (f *Friend) IncrSyncFriends(ctx context.Context) error {
	opt := Option[*model_struct.LocalFriend, *friend.GetIncrementalFriendsResp]{
		Ctx: ctx,
		DB:  f.db,
		Key: func(localFriend *model_struct.LocalFriend) string {
			return localFriend.FriendUserID
		},
		SyncKey: func() string {
			return "friend:" + f.loginUserID
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
		DeleteIDs: func(resp *friend.GetIncrementalFriendsResp) []string {
			return resp.DeleteUserIds
		},
		Changes: func(resp *friend.GetIncrementalFriendsResp) []*model_struct.LocalFriend {
			return util.Batch(ServerFriendToLocalFriendV2, resp.Changes)
		},
		Syncer: func(server, local []*model_struct.LocalFriend) error {
			return f.friendSyncer.Sync(ctx, server, local, nil)
		},
	}
	return opt.Sync()
}

type Option[V, R any] struct {
	Ctx       context.Context
	DB        db_interface.VersionSyncModel
	Key       func(V) string
	SyncKey   func() string
	Local     func() ([]V, error)
	Server    func(version *model_struct.LocalVersionSync) (R, error)
	Full      func(resp R) bool
	Version   func(resp R) (string, uint64)
	DeleteIDs func(resp R) []string
	Changes   func(resp R) []V
	Syncer    func(server, local []V) error
}

func (o *Option[V, R]) getVersionInfo() *model_struct.LocalVersionSync {
	key := o.SyncKey()
	versionInfo, err := o.DB.GetVersionSync(o.Ctx, key)
	if err != nil {
		log.ZInfo(o.Ctx, "get version info", "error", err)
		return &model_struct.LocalVersionSync{
			Key: key,
		}
	}
	return versionInfo
}

func (o *Option[V, R]) Sync() error {
	versionInfo := o.getVersionInfo()
	resp, err := o.Server(versionInfo)
	if err != nil {
		return err
	}
	delIDs := o.DeleteIDs(resp)
	changes := o.Changes(resp)
	updateVersionInfo := func() error {
		lvs := &model_struct.LocalVersionSync{
			Key:        versionInfo.Key,
			CreateTime: time.Now().UnixMilli(),
		}
		lvs.VersionID, lvs.Version = o.Version(resp)
		return o.DB.SetVersionSync(o.Ctx, lvs)
	}
	if len(delIDs)+len(changes) == 0 {
		return updateVersionInfo()
	}
	local, err := o.Local()
	if err != nil {
		return err
	}
	var server []V
	if o.Full(resp) {
		server = changes
	} else {
		kv := datautil.SliceToMapAny(local, func(v V) (string, V) {
			return o.Key(v), v
		})
		for i, change := range changes {
			key := o.Key(change)
			kv[key] = changes[i]
		}
		for _, id := range delIDs {
			delete(kv, id)
		}
		server = datautil.Values(kv)
	}
	if err := o.Syncer(server, local); err != nil {
		return err
	}
	return updateVersionInfo()
}
