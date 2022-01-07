package db

import (
	"errors"
	"open_im_sdk/pkg/utils"
)

func (u *DataBase) InsertGroup(groupInfo *LocalGroup) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.Wrap(u.conn.Create(groupInfo).Error, "InsertGroup failed")
}
func (u *DataBase) DeleteGroup(groupID string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.Wrap(u.conn.Where("group_id=?", groupID).Delete(&LocalGroup{}).Error, "DeleteGroup failed")
}
func (u *DataBase) UpdateGroup(groupInfo *LocalGroup) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.conn.Updates(groupInfo)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateGroup failed")
}
func (u *DataBase) GetGroupList() ([]LocalGroup, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var groupList []LocalGroup
	return groupList, utils.Wrap(u.conn.Find(&groupList).Error, "GetGroupList failed")
}
func (u *DataBase) GetGroupInfoByGroupID(groupID string) (g *LocalGroup, err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return g, utils.Wrap(u.conn.Where("group_id = ?", groupID).Find(g).Error, "GetGroupList failed")
}
