package db

import (
	"errors"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) InsertGroup(groupInfo *LocalGroup) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(groupInfo).Error, "InsertGroup failed")
}
func (d *DataBase) DeleteGroup(groupID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Where("group_id=?", groupID).Delete(&LocalGroup{}).Error, "DeleteGroup failed")
}
func (d *DataBase) UpdateGroup(groupInfo *LocalGroup) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Updates(groupInfo)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateGroup failed")
}
func (d *DataBase) GetJoinedGroupList() ([]*LocalGroup, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var groupList []LocalGroup
	err := d.conn.Find(&groupList).Error

	var transfer []*LocalGroup
	for _, v := range groupList {
		transfer = append(transfer, &v)
	}
	return transfer, utils.Wrap(err, "GetJoinedGroupList failed ")
}
func (d *DataBase) GetGroupInfoByGroupID(groupID string) (g *LocalGroup, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return g, utils.Wrap(d.conn.Where("group_id = ?", groupID).Find(g).Error, "GetGroupList failed")
}
