package db

import (
	"errors"
	"fmt"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) InsertGroup(groupInfo *model_struct.LocalGroup) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(groupInfo).Error, "InsertGroup failed")
}
func (d *DataBase) DeleteGroup(groupID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	localGroup := model_struct.LocalGroup{GroupID: groupID}
	return utils.Wrap(d.conn.Delete(&localGroup).Error, "DeleteGroup failed")
}
func (d *DataBase) UpdateGroup(groupInfo *model_struct.LocalGroup) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()

	t := d.conn.Model(groupInfo).Select("*").Updates(*groupInfo)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "")

}
func (d *DataBase) GetJoinedGroupList() ([]*model_struct.LocalGroup, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var groupList []model_struct.LocalGroup
	err := d.conn.Find(&groupList).Error
	var transfer []*model_struct.LocalGroup
	for _, v := range groupList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetJoinedGroupList failed ")
}
func (d *DataBase) GetGroupInfoByGroupID(groupID string) (*model_struct.LocalGroup, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var g model_struct.LocalGroup
	return &g, utils.Wrap(d.conn.Where("group_id = ?", groupID).Take(&g).Error, "GetGroupList failed")
}
func (d *DataBase) GetAllGroupInfoByGroupIDOrGroupName(keyword string, isSearchGroupID bool, isSearchGroupName bool) ([]*model_struct.LocalGroup, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var groupList []model_struct.LocalGroup
	var condition string
	if isSearchGroupID {
		if isSearchGroupName {
			condition = fmt.Sprintf("group_id like %q or name like %q", "%"+keyword+"%", "%"+keyword+"%")
		} else {
			condition = fmt.Sprintf("group_id like %q ", "%"+keyword+"%")
		}
	} else {
		condition = fmt.Sprintf("name like %q ", "%"+keyword+"%")
	}
	err := d.conn.Where(condition).Order("create_time DESC").Find(&groupList).Error
	var transfer []*model_struct.LocalGroup
	for _, v := range groupList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetAllGroupInfoByGroupIDOrGroupName failed ")
}
