package db

import (
	"errors"
	"open_im_sdk/pkg/utils"
)

func (u *open_im_sdk.UserRelated) _getGroupMemberInfoByGroupIDUserID(groupID, userID string) (*LocalGroupMember, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var groupMember LocalGroupMember
	return &groupMember, utils.wrap(u.imdb.Where("group_id = ? AND user_id = ?",
		groupID, userID).Error, "_getGroupMemberInfoByGroupIDUserID failed")
}

func (u *open_im_sdk.UserRelated) _getAllGroupMemberList() ([]LocalGroupMember, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var groupMemberList []LocalGroupMember
	return groupMemberList, utils.wrap(u.imdb.Find(&groupMemberList).Error, "_getAllGroupMemberList failed")
}

func (u *open_im_sdk.UserRelated) _getGroupMemberListByGroupID(groupID string) ([]LocalGroupMember, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var groupMemberList []LocalGroupMember
	return groupMemberList, utils.wrap(u.imdb.Where("group_id = ? ", groupID).Find(&groupMemberList).Error, "_getGroupMemberListByGroupID failed")
}

func (u *open_im_sdk.UserRelated) _insertGroupMember(groupMember *LocalGroupMember) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return u.imdb.Create(groupMember).Error
}

func (u *open_im_sdk.UserRelated) _deleteGroupMember(groupID, userID string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	groupMember := LocalGroupMember{GroupID: groupID, UserID: userID}
	return u.imdb.Where("group_id=? and user_id=?", groupID, userID).Delete(&groupMember).Error
}

func (u *open_im_sdk.UserRelated) _updateGroupMember(groupMember *LocalGroupMember) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(groupMember)
	if t.RowsAffected == 0 {
		return utils.wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return t.Error
}
