package friend

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
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
	local, err := f.db.GetAllFriendList(ctx)
	if err != nil {
		return err
	}
	userIDs := datautil.Slice(local, func(v *model_struct.LocalFriend) string {
		return v.FriendUserID
	})
	req := &friend.GetIncrementalFriendsReq{
		UserID: f.loginUserID,
		IdHash: IDHash(userIDs),
	}
	key := f.getVersionFriendKey()
	if res, err := f.db.GetVersionSync(ctx, key); err == nil {
		req.VersionID = res.VersionID
		req.Version = res.Version
		req.IdHash = res.IDHash
	} else {
		log.ZWarn(ctx, "get version sync failed", err, "key", key)
	}
	resp, err := util.CallApi[friend.GetIncrementalFriendsResp](ctx, constant.GetIncrementalFriends, req)
	if err != nil {
		return err
	}
	if len(resp.DeleteUserIds)+len(resp.Changes)+len(resp.SortUserIds) == 0 {
		return nil
	}
	if NotChange(resp.DeleteUserIds, resp.SortUserIds, resp.Changes) {
		return nil
	}
	b := Builder[string, *model_struct.LocalFriend]{
		Local: local,
		Key: func(v *model_struct.LocalFriend) string {
			return v.FriendUserID
		},
		Sort: func(v *model_struct.LocalFriend, index int32) *model_struct.LocalFriend {
			v.SortValue = index
			return v
		},
		Full:       resp.Full,
		SortKeys:   resp.SortUserIds,
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
		IDHash:     0,
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
	Sort       func(v V, index int32) V
	Full       bool
	SortKeys   []K
	DeleteKeys []K
	Changes    []V
}

func (b *Builder[K, V]) Build() []V {
	if b.Full {
		for i := range b.Changes {
			b.Changes[i] = b.Sort(b.Changes[i], int32(i))
		}
		return b.Changes
	}
	tmp := make(map[K]V)
	for i, e := range b.Local {
		tmp[b.Key(e)] = b.Local[i]
	}
	for _, key := range b.DeleteKeys {
		delete(tmp, key)
	}
	for i, e := range b.Changes {
		tmp[b.Key(e)] = b.Changes[i]
	}
	res := make([]V, 0, len(tmp))
	var incr int32
	if len(b.SortKeys) > 0 {
		for _, key := range b.SortKeys {
			if v, ok := tmp[key]; ok {
				incr++
				res = append(res, b.Sort(v, incr))
			}
		}
	} else {
		for _, t := range b.Local {
			if v, ok := tmp[b.Key(t)]; ok {
				incr++
				res = append(res, b.Sort(v, incr))
			}
		}
	}
	return res
}

func NotChange[K comparable, V any](delKeys []K, sortKeys []K, changes []V) bool {
	return len(delKeys)+len(sortKeys)+len(changes) > 0
}

func IDHash(ids []string) uint64 {
	if len(ids) == 0 {
		return 0
	}
	data, _ := json.Marshal(ids)
	sum := md5.Sum(data)
	return binary.BigEndian.Uint64(sum[:])
}
