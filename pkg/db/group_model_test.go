package db

import (
	"context"
	"testing"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
)

func Test_GetGroupMemberListByUserIDs(t *testing.T) {
	ctx := context.Background()
	db, err := NewDataBase(ctx, "1695766238", "./", 6)
	if err != nil {
		return
	}
	// log.ZError(ctx, "DB err test", nil, "key", "vale")

	_, err = db.GetGroupMemberListByUserIDs(ctx, "337845818", 0, []string{"123124123", "33333333"})
	if err != nil {
		t.Fatal(err)
	}

}
func Test_BatchInsertGroup(t *testing.T) {

	ctx := context.Background()
	db, err := NewDataBase(ctx, "1695766238", "./", 6)
	if err != nil {
		return
	}
	// log.ZError(ctx, "DB err test", nil, "key", "vale")

	localGroups := []*model_struct.LocalGroup{
		{
			GroupID:                "1234567",
			GroupName:              "test1234",
			Notification:           "",
			Introduction:           "",
			FaceURL:                "",
			CreateTime:             1666777417,
			Status:                 0,
			CreatorUserID:          "",
			GroupType:              0,
			OwnerUserID:            "",
			MemberCount:            0,
			Ex:                     "",
			AttachedInfo:           "",
			NeedVerification:       0,
			LookMemberInfo:         0,
			ApplyMemberFriend:      0,
			NotificationUpdateTime: 0,
			NotificationUserID:     "",
		},
		{
			GroupID:                "1234568",
			GroupName:              "test5678",
			Notification:           "New Notification",
			Introduction:           "This is a test group",
			FaceURL:                "https://example.com/face.png",
			CreateTime:             1666777420,
			Status:                 1,
			CreatorUserID:          "user123",
			GroupType:              1,
			OwnerUserID:            "user456",
			MemberCount:            10,
			Ex:                     "ex",
			AttachedInfo:           "Attach",
			NeedVerification:       1,
			LookMemberInfo:         1,
			ApplyMemberFriend:      1,
			NotificationUpdateTime: 1666777425,
			NotificationUserID:     "user789",
		},
	}

	err = db.BatchInsertGroup(ctx, localGroups)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_DeleteAllGroup(t *testing.T) {
	ctx := context.Background()
	db, err := NewDataBase(ctx, "1695766238", "./", 6)
	if err != nil {
		return
	}
	// log.ZError(ctx, "DB err test", nil, "key", "vale")

	err = db.DeleteAllGroup(ctx)
	if err != nil {
		t.Fatal(err)
	}

}
