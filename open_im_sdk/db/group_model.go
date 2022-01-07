package db

import (
	"errors"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/open_im_sdk/utils"
)

func (u *open_im_sdk.UserRelated) _insertGroup(groupInfo *open_im_sdk.LocalGroup) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.Wrap(u.imdb.Create(groupInfo).Error, "_insertGroup failed")
}
func (u *open_im_sdk.UserRelated) _deleteGroup(groupID string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.Wrap(u.imdb.Where("group_id=?", groupID).Delete(&open_im_sdk.LocalGroup{}).Error, "_deleteGroup failed")
}
func (u *open_im_sdk.UserRelated) _updateGroup(groupInfo *open_im_sdk.LocalGroup) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(groupInfo)
	if t.RowsAffected == 0 {
		return open_im_sdk.wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return open_im_sdk.wrap(t.Error, "_updateGroup failed")
}
func (u *open_im_sdk.UserRelated) _getGroupList() ([]open_im_sdk.LocalGroup, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var groupList []open_im_sdk.LocalGroup
	return groupList, open_im_sdk.wrap(u.imdb.Find(&groupList).Error, "_getGroupList failed")
}
func (u *open_im_sdk.UserRelated) _getGroupInfoByGroupID(groupID string) (g *open_im_sdk.LocalGroup, err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return g, open_im_sdk.wrap(u.imdb.Where("group_id = ?", groupID).Find(g).Error, "_getGroupList failed")
}
