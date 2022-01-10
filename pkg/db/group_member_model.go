package db

import (
	"errors"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) GetGroupMemberInfoByGroupIDUserID(groupID, userID string) (*LocalGroupMember, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var groupMember LocalGroupMember
	return &groupMember, utils.Wrap(d.conn.Where("group_id = ? AND user_id = ?",
		groupID, userID).Error, "GetGroupMemberInfoByGroupIDUserID failed")
}

func (d *DataBase) GetAllGroupMemberList() ([]LocalGroupMember, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var groupMemberList []LocalGroupMember
	return groupMemberList, utils.Wrap(d.conn.Find(&groupMemberList).Error, "GetAllGroupMemberList failed")
}

func (d *DataBase) GetGroupMemberListByGroupID(groupID string) ([]*LocalGroupMember, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var groupMemberList []LocalGroupMember
	return groupMemberList, utils.Wrap(d.conn.Where("group_id = ? ", groupID).Find(&groupMemberList).Error, "GetGroupMemberListByGroupID failed")
}

func (d *DataBase) InsertGroupMember(groupMember *LocalGroupMember) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return d.conn.Create(groupMember).Error
}

func (d *DataBase) DeleteGroupMember(groupID, userID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	groupMember := LocalGroupMember{GroupID: groupID, UserID: userID}
	return d.conn.Where("group_id=? and user_id=?", groupID, userID).Delete(&groupMember).Error
}

func (d *DataBase) UpdateGroupMember(groupMember *LocalGroupMember) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Updates(groupMember)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return t.Error
}
