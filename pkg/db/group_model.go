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
	localGroup := LocalGroup{GroupID: groupID}
	return utils.Wrap(d.conn.Delete(&localGroup).Error, "DeleteGroup failed")
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
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetJoinedGroupList failed ")
}
func (d *DataBase) GetGroupInfoByGroupID(groupID string) (*LocalGroup, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var g LocalGroup
	return &g, utils.Wrap(d.conn.Where("group_id = ?", groupID).Take(&g).Error, "GetGroupList failed")
}
