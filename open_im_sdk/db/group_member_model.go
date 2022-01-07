package db

import (
	"errors"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/open_im_sdk/utils"
)

func (u *open_im_sdk.UserRelated) _getGroupMemberInfoByGroupIDUserID(groupID, userID string) (*open_im_sdk.LocalGroupMember, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var groupMember open_im_sdk.LocalGroupMember
	return &groupMember, utils.wrap(u.imdb.Where("group_id = ? AND user_id = ?",
		groupID, userID).Error, "_getGroupMemberInfoByGroupIDUserID failed")
}

func (u *open_im_sdk.UserRelated) _getAllGroupMemberList() ([]open_im_sdk.LocalGroupMember, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var groupMemberList []open_im_sdk.LocalGroupMember
	return groupMemberList, utils.wrap(u.imdb.Find(&groupMemberList).Error, "_getAllGroupMemberList failed")
}

func (u *open_im_sdk.UserRelated) _getGroupMemberListByGroupID(groupID string) ([]open_im_sdk.LocalGroupMember, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var groupMemberList []open_im_sdk.LocalGroupMember
	return groupMemberList, utils.wrap(u.imdb.Where("group_id = ? ", groupID).Find(&groupMemberList).Error, "_getGroupMemberListByGroupID failed")
}

func (u *open_im_sdk.UserRelated) _insertGroupMember(groupMember *open_im_sdk.LocalGroupMember) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return u.imdb.Create(groupMember).Error
}

func (u *open_im_sdk.UserRelated) _deleteGroupMember(groupID, userID string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	groupMember := open_im_sdk.LocalGroupMember{GroupID: groupID, UserID: userID}
	return u.imdb.Where("group_id=? and user_id=?", groupID, userID).Delete(&groupMember).Error
}

func (u *open_im_sdk.UserRelated) _updateGroupMember(groupMember *open_im_sdk.LocalGroupMember) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(groupMember)
	if t.RowsAffected == 0 {
		return utils.wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return t.Error
}
