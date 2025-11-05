//go:build !js
// +build !js

package db

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/errs"
)

func (d *DataBase) InsertFriend(ctx context.Context, friend *model_struct.LocalFriend) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Create(friend).Error, "InsertFriend failed")
}

func (d *DataBase) DeleteFriendDB(ctx context.Context, friendUserID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Where("owner_user_id=? and friend_user_id=?", d.loginUserID, friendUserID).Delete(&model_struct.LocalFriend{}).Error, "DeleteFriend failed")
}

func (d *DataBase) GetFriendListCount(ctx context.Context) (int64, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var count int64
	return count, errs.WrapMsg(d.conn.WithContext(ctx).Model(&model_struct.LocalFriend{}).Count(&count).Error, "GetFriendListCount failed")
}

func (d *DataBase) UpdateFriend(ctx context.Context, friend *model_struct.LocalFriend) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()

	t := d.conn.WithContext(ctx).Model(friend).Select("*").Updates(*friend)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.Wrap(t.Error)

}
func (d *DataBase) GetAllFriendList(ctx context.Context) ([]*model_struct.LocalFriend, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var friendList []*model_struct.LocalFriend
	return friendList, errs.WrapMsg(d.conn.WithContext(ctx).Where("owner_user_id = ?", d.loginUserID).Find(&friendList).Error,
		"GetFriendList failed")
}

func (d *DataBase) GetPageFriendList(ctx context.Context, offset, count int) ([]*model_struct.LocalFriend, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var friendList []*model_struct.LocalFriend
	err := errs.WrapMsg(d.conn.WithContext(ctx).Where("owner_user_id = ?", d.loginUserID).Offset(offset).Limit(count).Order("name").Find(&friendList).Error,
		"GetFriendList failed")
	return friendList, err
}

func (d *DataBase) BatchInsertFriend(ctx context.Context, friendList []*model_struct.LocalFriend) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	if friendList == nil {
		return errs.New("nil").Wrap()
	}
	return errs.WrapMsg(d.conn.WithContext(ctx).Create(friendList).Error, "BatchInsertFriendList failed")
}

func (d *DataBase) DeleteAllFriend(ctx context.Context) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model_struct.LocalFriend{}).Error, "DeleteAllFriend failed")
}

func (d *DataBase) SearchFriendList(ctx context.Context, keyword string, isSearchUserID, isSearchNickname, isSearchRemark bool) ([]*model_struct.LocalFriend, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()

	var friendList []*model_struct.LocalFriend
	query := d.conn.WithContext(ctx)

	var conditions []string
	var args []any

	if isSearchUserID {
		conditions = append(conditions, "friend_user_id LIKE ?")
		args = append(args, "%"+keyword+"%")
	}

	if isSearchNickname {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+keyword+"%")
	}

	if isSearchRemark {
		conditions = append(conditions, "remark LIKE ?")
		args = append(args, "%"+keyword+"%")
	}

	if len(conditions) > 0 {
		query = query.Where(
			"("+strings.Join(conditions, " OR ")+")",
			args...,
		)
	}

	err := query.Order("create_time DESC").Find(&friendList).Error
	return friendList, errs.WrapMsg(err, "SearchFriendList failed")
}

func (d *DataBase) GetFriendInfoByFriendUserID(ctx context.Context, FriendUserID string) (*model_struct.LocalFriend, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var friend model_struct.LocalFriend
	return &friend, errs.WrapMsg(d.conn.WithContext(ctx).Where("owner_user_id = ? AND friend_user_id = ?",
		d.loginUserID, FriendUserID).Take(&friend).Error, "GetFriendInfoByFriendUserID failed")
}

func (d *DataBase) GetFriendInfoList(ctx context.Context, friendUserIDList []string) ([]*model_struct.LocalFriend, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var friendList []*model_struct.LocalFriend
	err := errs.WrapMsg(d.conn.WithContext(ctx).Where("friend_user_id IN ?", friendUserIDList).Find(&friendList).Error, "GetFriendInfoListByFriendUserID failed")
	return friendList, err
}
func (d *DataBase) UpdateColumnsFriend(ctx context.Context, friendIDs []string, args map[string]any) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Model(&model_struct.LocalFriend{}).Where("friend_user_id IN ?", friendIDs).Updates(args).Error, "UpdateColumnsFriend failed")
}
