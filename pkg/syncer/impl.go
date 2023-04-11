package syncer

import (
	"context"
	"fmt"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	"reflect"
)

var _ DBInterface = (*DB)(nil)

func Inserts[T any](ctx context.Context, ts []T, fn func(ctx context.Context, t *T) error) error {
	for i := range ts {
		if err := fn(ctx, &ts[i]); err != nil {
			return err
		}
	}
	return nil
}

func getOnlyWhere[T any](where *Where, column string) (T, error) {
	var t T
	if len(where.Columns) != 1 {
		return t, fmt.Errorf("where columns length must be 1")
	}
	c := where.Columns[0]
	if c.Name != column {
		return t, fmt.Errorf("where column name must be %s", column)
	}
	return c.Value.(T), nil
}

func getKeyWhere[T any](where *Where, column string) (T, error) {
	for _, c := range where.Columns {
		if c.Name == column {
			return c.Value.(T), nil
		}
	}
	var t T
	return t, fmt.Errorf("where column name must be %s", column)
}

func updateToModel(data map[string]any, model interface{}) {
	v := reflect.ValueOf(model).Elem()
	for k, v1 := range data {
		v.FieldByName(k).Set(reflect.ValueOf(v1))
	}
}

}


type DB struct {
	db db_interface.DataBase
}

func (a *DB) Insert(ctx context.Context, ms any) error {
	switch v := ms.(type) {
	case []model_struct.LocalFriend:
		return Inserts(ctx, v, a.db.InsertFriend)
	case []model_struct.LocalFriendRequest:
		return Inserts(ctx, v, a.db.InsertFriendRequest)
	case []model_struct.LocalGroup:
		return Inserts(ctx, v, a.db.InsertGroup)
	case []model_struct.LocalGroupMember:
		return Inserts(ctx, v, a.db.InsertGroupMember)
	case []model_struct.LocalGroupRequest:
		return Inserts(ctx, v, a.db.InsertGroupRequest)
	case []model_struct.LocalUser:
		return Inserts(ctx, v, a.db.InsertLoginUser)
	case []model_struct.LocalBlack:
		return Inserts(ctx, v, a.db.InsertBlack)
	case []model_struct.LocalSeqData:
		//return Inserts(ctx, v, a.db.InsertSeqData)
	case []model_struct.LocalSeq:
		//return Inserts(ctx, v, a.db.InsertSeq)
	case []model_struct.LocalChatLog:
		//return Inserts(ctx, v, a.db.InsertChatLog)
	case []model_struct.LocalErrChatLog:
		//return Inserts(ctx, v, a.db.InsertErrChatLog)
	case []model_struct.TempCacheLocalChatLog:
		//return Inserts(ctx, v, a.db.InsertTempCacheChatLog)
	case []model_struct.LocalConversation:
		return Inserts(ctx, v, a.db.InsertConversation)
	case []model_struct.LocalConversationUnreadMessage:
		//return Inserts(ctx, v, a.db.InsertConversationUnreadMessage)
	case []model_struct.LocalAdminGroupRequest:
		return Inserts(ctx, v, a.db.InsertAdminGroupRequest)
	case []model_struct.LocalDepartment:
		return Inserts(ctx, v, a.db.InsertDepartment)
	case []model_struct.LocalDepartmentMember:
		return Inserts(ctx, v, a.db.InsertDepartmentMember)
	case []model_struct.SearchDepartmentMemberResult:
		//return Inserts(ctx, v, a.db.InsertSearchDepartmentMemberResult)
	case []model_struct.LocalChatLogReactionExtensions:
		//return Inserts(ctx, v, a.db.InsertChatLogReactionExtensions)
	case []model_struct.LocalWorkMomentsNotification:
		//return Inserts(ctx, v, a.db.InsertWorkMomentsNotification)
	case []model_struct.LocalWorkMomentsNotificationUnreadCount:
		//return Inserts(ctx, v, a.db.InsertWorkMomentsNotificationUnreadCount)
	}
	return fmt.Errorf("untreated insert type %T", ms)
}

func (a *DB) delete(ctx context.Context, m reflect.Type, where *Where) error {
	switch reflect.Zero(m).Interface().(type) {
	case model_struct.LocalFriend:
		val, err := getOnlyWhere[string](where, "friend_user_id")
		if err != nil {
			return err
		}
		return a.db.DeleteFriendDB(ctx, val)
	case model_struct.LocalFriendRequest:
		friendUserId, err := getKeyWhere[string](where, "friend_user_id")
		if err != nil {
			return err
		}
		toUserId, err := getKeyWhere[string](where, "to_user_id")
		if err == nil {
			return a.db.DeleteFriendRequestBothUserID(ctx, friendUserId, toUserId)
		}
		return a.db.DeleteFriendDB(ctx, friendUserId)
	case model_struct.LocalGroup:
		val, err := getOnlyWhere[string](where, "group_id")
		if err != nil {
			return err
		}
		return a.db.DeleteGroup(ctx, val)
	case model_struct.LocalGroupMember:
		groupID, err := getKeyWhere[string](where, "group_id")
		if err != nil {
			return err
		}
		userId, err := getKeyWhere[string](where, "user_id")
		if err != nil {
			return err
		}
		return a.db.DeleteGroupMember(ctx, groupID, userId)
	case model_struct.LocalGroupRequest:
		groupID, err := getKeyWhere[string](where, "group_id")
		if err != nil {
			return err
		}
		userId, err := getKeyWhere[string](where, "user_id")
		if err != nil {
			return err
		}
		return a.db.DeleteGroupRequest(ctx, groupID, userId)
	case model_struct.LocalUser:
		//userId, err := getOnlyWhere[string](where, "user_id")
		//if err == nil {
		//	return err
		//}

	case model_struct.LocalBlack:
		userId, err := getKeyWhere[string](where, "block_user_id")
		if err == nil {
			return err
		}
		return a.db.DeleteBlack(ctx, userId)
	case model_struct.LocalSeqData:

	case model_struct.LocalSeq:

	case model_struct.LocalChatLog:

	case model_struct.LocalErrChatLog:

	case model_struct.TempCacheLocalChatLog:

	case model_struct.LocalConversation:

	case model_struct.LocalConversationUnreadMessage:

	case model_struct.LocalAdminGroupRequest:

	case model_struct.LocalDepartment:

	case model_struct.LocalDepartmentMember:

	case model_struct.SearchDepartmentMemberResult:

	case model_struct.LocalChatLogReactionExtensions:

	case model_struct.LocalWorkMomentsNotification:

	case model_struct.LocalWorkMomentsNotificationUnreadCount:

	}
	return fmt.Errorf("untreated delete type %s", m.String())
}

func (a *DB) Delete(ctx context.Context, m reflect.Type, where []*Where) error {
	for i := range where {
		if err := a.delete(ctx, m, where[i]); err != nil {
			return err
		}
	}
	return nil
}

func (a *DB) Update(ctx context.Context, m reflect.Type, where *Where, data map[string]any) error {
	switch reflect.Zero(m).Interface().(type) {
	case model_struct.LocalFriend:
		//getOnlyWhere(where, "id")
		//a.db.UpdateFriend(ctx)

	case model_struct.LocalFriendRequest:
		//a.db.UpdateFriendRequest(ctx)
	case model_struct.LocalGroup:

	case model_struct.LocalGroupMember:

	case model_struct.LocalGroupRequest:

	case model_struct.LocalUser:

	case model_struct.LocalBlack:

	case model_struct.LocalSeqData:

	case model_struct.LocalSeq:

	case model_struct.LocalChatLog:

	case model_struct.LocalErrChatLog:

	case model_struct.TempCacheLocalChatLog:

	case model_struct.LocalConversation:

	case model_struct.LocalConversationUnreadMessage:

	case model_struct.LocalAdminGroupRequest:

	case model_struct.LocalDepartment:

	case model_struct.LocalDepartmentMember:

	case model_struct.SearchDepartmentMemberResult:

	case model_struct.LocalChatLogReactionExtensions:

	case model_struct.LocalWorkMomentsNotification:

	case model_struct.LocalWorkMomentsNotificationUnreadCount:

	}
	return fmt.Errorf("untreated update type %s", m.String())
}

func (a *DB) FindOffset(ctx context.Context, m reflect.Type, where *Where, offset int, limit int) (any, error) {
	switch reflect.Zero(m).Interface().(type) {
	case model_struct.LocalFriend:
		a.db.GetAllFriendList(ctx)

	case model_struct.LocalFriendRequest:
	case model_struct.LocalGroup:
	case model_struct.LocalGroupMember:
	case model_struct.LocalGroupRequest:
	case model_struct.LocalUser:
	case model_struct.LocalBlack:
	case model_struct.LocalSeqData:
	case model_struct.LocalSeq:
	case model_struct.LocalChatLog:
	case model_struct.LocalErrChatLog:
	case model_struct.TempCacheLocalChatLog:
	case model_struct.LocalConversation:
	case model_struct.LocalConversationUnreadMessage:
	case model_struct.LocalAdminGroupRequest:
	case model_struct.LocalDepartment:
	case model_struct.LocalDepartmentMember:
	case model_struct.SearchDepartmentMemberResult:
	case model_struct.LocalChatLogReactionExtensions:
	case model_struct.LocalWorkMomentsNotification:
	case model_struct.LocalWorkMomentsNotificationUnreadCount:
	}
	return nil, fmt.Errorf("untreated findoffset type %s", m.String())
}
