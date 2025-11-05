//go:build !js
// +build !js

package db

import (
	"context"
	"errors"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/errs"
)

func (d *DataBase) GetGroupMemberInfoByGroupIDUserID(ctx context.Context, groupID, userID string) (*model_struct.LocalGroupMember, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var groupMember model_struct.LocalGroupMember
	return &groupMember, errs.WrapMsg(d.conn.WithContext(ctx).Where("group_id = ? AND user_id = ?",
		groupID, userID).Take(&groupMember).Error, "GetGroupMemberInfoByGroupIDUserID failed")
}

func (d *DataBase) GetAllGroupMemberList(ctx context.Context) ([]model_struct.LocalGroupMember, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var groupMemberList []model_struct.LocalGroupMember
	return groupMemberList, errs.WrapMsg(d.conn.WithContext(ctx).Find(&groupMemberList).Error, "GetAllGroupMemberList failed")
}

func (d *DataBase) GetGroupMemberCount(ctx context.Context, groupID string) (int32, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var count int64
	err := d.conn.WithContext(ctx).Model(&model_struct.LocalGroupMember{}).Where("group_id = ? ", groupID).Count(&count).Error
	return int32(count), errs.WrapMsg(err, "GetGroupMemberCount failed")
}

func (d *DataBase) GetGroupSomeMemberInfo(ctx context.Context, groupID string, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var groupMemberList []*model_struct.LocalGroupMember
	err := d.conn.WithContext(ctx).Where("group_id = ? AND user_id IN ? ", groupID, userIDList).Find(&groupMemberList).Error
	return groupMemberList, errs.WrapMsg(err, "GetGroupMemberListByGroupID failed ")
}

func (d *DataBase) GetGroupMemberListByGroupID(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var groupMemberList []*model_struct.LocalGroupMember
	err := d.conn.WithContext(ctx).Where("group_id = ? ", groupID).Find(&groupMemberList).Error
	return groupMemberList, errs.WrapMsg(err, "GetGroupMemberListByGroupID failed ")
}

func (d *DataBase) GetGroupMemberListByUserIDs(ctx context.Context, groupID string, filter int32, userIDs []string) ([]*model_struct.LocalGroupMember, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var groupMemberList []*model_struct.LocalGroupMember
	var err error
	switch filter {
	case constant.GroupFilterAll:
		err = d.conn.WithContext(ctx).Where("group_id = ? AND user_id IN ?", groupID, userIDs).Find(&groupMemberList).Error
	case constant.GroupFilterOwner:
		err = d.conn.WithContext(ctx).Where("group_id = ? AND role_level = ? AND user_id IN ?", groupID, constant.GroupOwner, userIDs).Find(&groupMemberList).Error
	case constant.GroupFilterAdmin:
		err = d.conn.WithContext(ctx).Where("group_id = ? AND role_level = ? AND user_id IN ?", groupID, constant.GroupAdmin, userIDs).Find(&groupMemberList).Error
	case constant.GroupFilterOrdinaryUsers:
		err = d.conn.WithContext(ctx).Where("group_id = ? AND role_level = ? AND user_id IN ?", groupID, constant.GroupOrdinaryUsers, userIDs).Find(&groupMemberList).Error
	case constant.GroupFilterAdminAndOrdinaryUsers:
		err = d.conn.WithContext(ctx).Where("group_id = ? AND (role_level = ? OR role_level = ?) AND user_id IN ?", groupID, constant.GroupAdmin, constant.GroupOrdinaryUsers, userIDs).Find(&groupMemberList).Error
	case constant.GroupFilterOwnerAndAdmin:
		err = d.conn.WithContext(ctx).Where("group_id = ? AND (role_level = ? OR role_level = ?) AND user_id IN ?", groupID, constant.GroupOwner, constant.GroupAdmin, userIDs).Find(&groupMemberList).Error
	default:
		return nil, errs.New("filter args failed.", "filter", filter).Wrap()
	}

	return groupMemberList, errs.Wrap(err)

}

func (d *DataBase) GetGroupMemberListSplit(ctx context.Context, groupID string, filter int32, offset, count int) ([]*model_struct.LocalGroupMember, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var groupMemberList []*model_struct.LocalGroupMember
	var err error
	switch filter {
	case constant.GroupFilterAll:
		err = d.conn.WithContext(ctx).Where("group_id = ?", groupID).Order("role_level DESC,join_time ASC").Offset(offset).Limit(count).Find(&groupMemberList).Error
	case constant.GroupFilterOwner:
		err = d.conn.WithContext(ctx).Where("group_id = ? And role_level = ?", groupID, constant.GroupOwner).Offset(offset).Limit(count).Find(&groupMemberList).Error
	case constant.GroupFilterAdmin:
		err = d.conn.WithContext(ctx).Where("group_id = ? And role_level = ?", groupID, constant.GroupAdmin).Order("join_time ASC").Offset(offset).Limit(count).Find(&groupMemberList).Error
	case constant.GroupFilterOrdinaryUsers:
		err = d.conn.WithContext(ctx).Where("group_id = ? And role_level = ?", groupID, constant.GroupOrdinaryUsers).Order("join_time ASC").Offset(offset).Limit(count).Find(&groupMemberList).Error
	case constant.GroupFilterAdminAndOrdinaryUsers:
		err = d.conn.WithContext(ctx).Where("group_id = ? And (role_level = ? or role_level = ?)", groupID, constant.GroupAdmin, constant.GroupOrdinaryUsers).Order("role_level DESC,join_time ASC").Offset(offset).Limit(count).Find(&groupMemberList).Error
	case constant.GroupFilterOwnerAndAdmin:
		err = d.conn.WithContext(ctx).Where("group_id = ? And (role_level = ? or role_level = ?)", groupID, constant.GroupOwner, constant.GroupAdmin).Order("role_level DESC,join_time ASC").Offset(offset).Limit(count).Find(&groupMemberList).Error
	default:
		return nil, errs.New("filter args failed", "filter", filter).Wrap()
	}

	return groupMemberList, errs.WrapMsg(err, "GetGroupMemberListSplit failed ")
}

func (d *DataBase) GetGroupMemberOwnerAndAdminDB(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var groupMemberList []*model_struct.LocalGroupMember
	err := d.conn.WithContext(ctx).Where("group_id = ? And (role_level = ? OR role_level = ?)", groupID, constant.GroupOwner, constant.GroupAdmin).Order("join_time DESC").Find(&groupMemberList).Error

	return groupMemberList, errs.WrapMsg(err, "GetGroupMemberListSplit failed ")
}

func (d *DataBase) GetGroupMemberListSplitByJoinTimeFilter(ctx context.Context, groupID string, offset, count int, joinTimeBegin, joinTimeEnd int64, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var groupMemberList []*model_struct.LocalGroupMember
	var err error
	if len(userIDList) == 0 {
		err = d.conn.WithContext(ctx).Where("group_id = ? And join_time  between ? and ? ", groupID, joinTimeBegin, joinTimeEnd).Order("join_time DESC").Offset(offset).Limit(count).Find(&groupMemberList).Error
	} else {
		err = d.conn.WithContext(ctx).Where("group_id = ? And join_time  between ? and ? And user_id NOT IN ?", groupID, joinTimeBegin, joinTimeEnd, userIDList).Order("join_time DESC").Offset(offset).Limit(count).Find(&groupMemberList).Error
	}
	return groupMemberList, errs.WrapMsg(err, "GetGroupMemberListSplitByJoinTimeFilter failed ")
}

func (d *DataBase) InsertGroupMember(ctx context.Context, groupMember *model_struct.LocalGroupMember) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Create(groupMember).Error, "")
}

func (d *DataBase) BatchInsertGroupMember(ctx context.Context, groupMemberList []*model_struct.LocalGroupMember) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Create(groupMemberList).Error, "BatchInsertGroupMember failed")
}

func (d *DataBase) DeleteGroupMember(ctx context.Context, groupID, userID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var groupMember model_struct.LocalGroupMember
	return d.conn.WithContext(ctx).Where("group_id=? and user_id=?", groupID, userID).Delete(&groupMember).Error
}

func (d *DataBase) DeleteGroupAllMembers(ctx context.Context, groupID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var groupMember model_struct.LocalGroupMember
	return d.conn.WithContext(ctx).Where("group_id=? ", groupID).Delete(&groupMember).Error
}

func (d *DataBase) UpdateGroupMember(ctx context.Context, groupMember *model_struct.LocalGroupMember) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Model(groupMember).Select("*").Updates(*groupMember)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.WrapMsg(t.Error, "")
}

func (d *DataBase) SearchGroupMembersDB(ctx context.Context, keyword string, groupID string, isSearchMemberNickname, isSearchUserID bool, offset, count int) (result []*model_struct.LocalGroupMember, err error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	if !isSearchMemberNickname && !isSearchUserID {
		return nil, errors.New("args failed")
	}

	var groupMemberList []*model_struct.LocalGroupMember
	query := d.conn.WithContext(ctx)

	if isSearchUserID && isSearchMemberNickname {
		query = query.Where("user_id LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	} else if isSearchUserID {
		query = query.Where("user_id LIKE ?", "%"+keyword+"%")
	} else if isSearchMemberNickname {
		query = query.Where("nickname LIKE ?", "%"+keyword+"%")
	}

	if groupID != "" {
		query = query.Where("group_id IN ?", []string{groupID})
	}

	err = query.Order("role_level DESC,join_time ASC").
		Offset(offset).
		Limit(count).
		Find(&groupMemberList).Error

	return groupMemberList, err
}
