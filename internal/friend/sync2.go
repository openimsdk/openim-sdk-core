package friend

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/protocol/friend"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"time"
)

func (f *Friend) getVersionFriendKey() string {
	return "friend:" + f.loginUserID
}

func (f *Friend) IncrSyncFriends(ctx context.Context) error {
	req := &friend.GetIncrementalFriendsReq{
		UserID: f.loginUserID,
	}
	key := f.getVersionFriendKey()
	if res, err := f.db.GetVersionSync(ctx, key); err == nil {
		req.VersionID = res.VersionID
		req.Version = res.Version
	} else {
		log.ZWarn(ctx, "get version sync failed", err, "key", key)
	}
	resp, err := util.CallApi[friend.GetIncrementalFriendsResp](ctx, constant.GetIncrementalFriends, req)
	if err != nil {
		return err
	}
	if NotChange(resp.DeleteUserIds, resp.Changes) {
		lv := model_struct.LocalVersionSync{
			Key:        key,
			VersionID:  resp.VersionID,
			Version:    resp.Version,
			CreateTime: time.Now().UnixMilli(),
		}
		return f.db.SetVersionSync(ctx, &lv)
	}
	local, err := f.db.GetAllFriendList(ctx)
	if err != nil {
		return err
	}
	b := Builder[string, *model_struct.LocalFriend]{
		Local: local,
		Key: func(v *model_struct.LocalFriend) string {
			return v.FriendUserID
		},
		Full:       resp.Full,
		DeleteKeys: resp.DeleteUserIds,
		Changes:    util.Batch(ServerFriendToLocalFriendV2, resp.Changes),
	}
	if err := f.friendSyncer.Sync(ctx, b.Build(), local, nil); err != nil {
		return err
	}
	lv := model_struct.LocalVersionSync{
		Key:        key,
		VersionID:  resp.VersionID,
		Version:    resp.Version,
		CreateTime: time.Now().UnixMilli(),
	}
	if err := f.db.SetVersionSync(ctx, &lv); err != nil {
		return err
	}
	return nil
}

type Builder[K comparable, V any] struct {
	Local      []V
	Key        func(V) K
	Full       bool
	DeleteKeys []K
	Changes    []V
}

func (b *Builder[K, V]) Build() []V {
	if b.Full {
		return b.Changes
	}
	res := make([]V, 0, len(b.Local)+len(b.Changes))
	var delSet map[K]struct{}
	if len(b.Local)+len(b.DeleteKeys) > 0 {
		delSet = datautil.SliceSet(b.DeleteKeys)
	}
	for i, v := range b.Local {
		if _, ok := delSet[b.Key(v)]; !ok {
			res = append(res, b.Local[i])
		}
	}
	return append(res, b.Changes...)
}

func NotChange[K comparable, V any](delKeys []K, changes []V) bool {
	return len(delKeys)+len(changes) == 0
}
