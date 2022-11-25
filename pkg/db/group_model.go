package db

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) InsertGroup(groupInfo *model_struct.LocalGroup) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	return utils.Wrap(d.conn.Create(groupInfo).Error, "InsertGroup failed")
}
func (d *DataBase) DeleteGroup(groupID string) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	localGroup := model_struct.LocalGroup{GroupID: groupID}
	return utils.Wrap(d.conn.Delete(&localGroup).Error, "DeleteGroup failed")
}
func (d *DataBase) UpdateGroup(groupInfo *model_struct.LocalGroup) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()

	t := d.conn.Model(groupInfo).Select("*").Updates(*groupInfo)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "")

}
func (d *DataBase) GetJoinedGroupListDB() ([]*model_struct.LocalGroup, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
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
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	var g model_struct.LocalGroup
	return &g, utils.Wrap(d.conn.Where("group_id = ?", groupID).Take(&g).Error, "GetGroupList failed")
}
func (d *DataBase) GetAllGroupInfoByGroupIDOrGroupName(keyword string, isSearchGroupID bool, isSearchGroupName bool) ([]*model_struct.LocalGroup, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()

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

func (d *DataBase) AddMemberCount(groupID string) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	group := model_struct.LocalGroup{GroupID: groupID}
	return utils.Wrap(d.conn.Model(&group).Updates(map[string]interface{}{"member_count": gorm.Expr("member_count+1")}).Error, "")
}

func (d *DataBase) SubtractMemberCount(groupID string) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	group := model_struct.LocalGroup{GroupID: groupID}
	return utils.Wrap(d.conn.Model(&group).Updates(map[string]interface{}{"member_count": gorm.Expr("member_count-1")}).Error, "")
}
