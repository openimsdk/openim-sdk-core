package db

import (
	"errors"
	"open_im_sdk/pkg/utils"
)

func (u *DataBase) GetGroupMemberInfoByGroupIDUserID(groupID, userID string) (*LocalGroupMember, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var groupMember LocalGroupMember
	return &groupMember, utils.Wrap(u.conn.Where("group_id = ? AND user_id = ?",
		groupID, userID).Error, "GetGroupMemberInfoByGroupIDUserID failed")
}

func (u *DataBase) GetAllGroupMemberList() ([]LocalGroupMember, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var groupMemberList []LocalGroupMember
	return groupMemberList, utils.Wrap(u.conn.Find(&groupMemberList).Error, "GetAllGroupMemberList failed")
}

func (u *DataBase) GetGroupMemberListByGroupID(groupID string) ([]LocalGroupMember, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var groupMemberList []LocalGroupMember
	return groupMemberList, utils.Wrap(u.conn.Where("group_id = ? ", groupID).Find(&groupMemberList).Error, "GetGroupMemberListByGroupID failed")
}

func (u *DataBase) InsertGroupMember(groupMember *LocalGroupMember) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return u.conn.Create(groupMember).Error
}

func (u *DataBase) DeleteGroupMember(groupID, userID string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	groupMember := LocalGroupMember{GroupID: groupID, UserID: userID}
	return u.conn.Where("group_id=? and user_id=?", groupID, userID).Delete(&groupMember).Error
}

func (u *DataBase) UpdateGroupMember(groupMember *LocalGroupMember) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.conn.Updates(groupMember)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return t.Error
}
