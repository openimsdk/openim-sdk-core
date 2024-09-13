package db

import (
	"context"
	"testing"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/log"
)

func Test_GetFriendListCount(t *testing.T) {
	ctx := context.Background()
	db, err := NewDataBase(ctx, "1695766238", "./", 6)
	if err != nil {
		return
	}

	info, err := db.GetFriendListCount(ctx)
	if err != nil {
		t.Fatal(err)
	}
	log.ZDebug(ctx, "info is ", "info", info)
}

func Test_BatchInsertFriend(t *testing.T) {
	ctx := context.Background()
	db, err := NewDataBase(ctx, "1695766238", "./", 6)
	if err != nil {
		return
	}
	localFriends := []*model_struct.LocalFriend{
		{
			OwnerUserID:    "1695766238",
			FriendUserID:   "1234567890",
			Remark:         "hello",
			CreateTime:     1666778999,
			AddSource:      0,
			OperatorUserID: "789",
			Nickname:       "hhhh",
			FaceURL:        "",
			Ex:             "",
			AttachedInfo:   "",
			IsPinned:       false,
		},
		{
			OwnerUserID:    "1695766238",
			FriendUserID:   "1234567891",
			Remark:         "hi",
			CreateTime:     1666779000,
			AddSource:      1,
			OperatorUserID: "790",
			Nickname:       "aaaa",
			FaceURL:        "https://example.com/face.png",
			Ex:             "example",
			AttachedInfo:   "info",
			IsPinned:       true,
		},
		{
			OwnerUserID:    "1695766238",
			FriendUserID:   "1234567892",
			Remark:         "hey",
			CreateTime:     1666779001,
			AddSource:      2,
			OperatorUserID: "791",
			Nickname:       "bbbb",
			FaceURL:        "https://example.com/face2.png",
			Ex:             "example2",
			AttachedInfo:   "info2",
			IsPinned:       false,
		},
	}

	err = db.BatchInsertFriend(ctx, localFriends)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_DeleteAllFriend(t *testing.T) {
	ctx := context.Background()
	db, err := NewDataBase(ctx, "1695766238", "./", 6)
	if err != nil {
		return
	}
	// log.ZError(ctx, "DB err test", nil, "key", "vale")

	err = db.DeleteAllFriend(ctx)
	if err != nil {
		t.Fatal(err)
	}

}

// func Test_UpdateColumnsFriend(t *testing.T) {
// 	ctx := context.Background()
// 	db, err := db.NewDataBase(ctx, "1695766238", "./", 6)
// 	if err != nil {
// 		return
// 	}
// 	// log.ZError(ctx, "DB err test", nil, "key", "vale")

// 	err = db.UpdateColumnsFriend(ctx)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// }
