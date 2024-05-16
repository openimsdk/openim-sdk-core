package friend

import (
	"crypto/md5"
	"encoding/binary"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/utils/datautil"
	"strconv"
	"strings"
)

func (f *Friend) CalculateHash(friends []*model_struct.LocalFriend) uint64 {
	datautil.SortAny(friends, func(a, b *model_struct.LocalFriend) bool {
		return a.CreateTime > b.CreateTime
	})
	if len(friends) > constant.MaxSyncPullNumber {
		friends = friends[:constant.MaxSyncPullNumber]
	}
	hashStr := strings.Join(datautil.Slice(friends, func(f *model_struct.LocalFriend) string {
		return strings.Join([]string{
			f.FriendUserID,
			f.Remark,
			strconv.FormatInt(f.CreateTime, 10),
			strconv.Itoa(int(f.AddSource)),
			f.OperatorUserID,
			f.Ex,
			strconv.FormatBool(f.IsPinned),
		}, ",")
	}), ";")
	sum := md5.Sum([]byte(hashStr))
	return binary.BigEndian.Uint64(sum[:])
}
