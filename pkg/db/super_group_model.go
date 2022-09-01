package db

import (
	"errors"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) GetJoinedSuperGroupList() ([]*model_struct.LocalGroup, error) {
	d.superGroupMtx.Lock()
	defer d.superGroupMtx.Unlock()
	var groupList []model_struct.LocalGroup
	err := d.conn.Table(constant.SuperGroupTableName).Find(&groupList).Error
	var transfer []*model_struct.LocalGroup
	for _, v := range groupList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetJoinedSuperGroupList failed ")
}

func (d *DataBase) GetJoinedSuperGroupIDList() ([]string, error) {
	groupList, err := d.GetJoinedSuperGroupList()
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var groupIDList []string
	for _, v := range groupList {
		groupIDList = append(groupIDList, v.GroupID)
	}
	return groupIDList, nil
}

func (d *DataBase) InsertSuperGroup(groupInfo *model_struct.LocalGroup) error {
	d.superGroupMtx.Lock()
	defer d.superGroupMtx.Unlock()
	return utils.Wrap(d.conn.Table(constant.SuperGroupTableName).Create(groupInfo).Error, "InsertSuperGroup failed")
}

func (d *DataBase) DeleteAllSuperGroup() error {
	d.superGroupMtx.Lock()
	defer d.superGroupMtx.Unlock()
	return utils.Wrap(d.conn.Table(constant.SuperGroupTableName).Delete(&model_struct.LocalGroup{}).Error, "DeleteAllSuperGroup failed")
}

func (d *DataBase) GetSuperGroupInfoByGroupID(groupID string) (*model_struct.LocalGroup, error) {
	d.superGroupMtx.Lock()
	defer d.superGroupMtx.Unlock()
	var g model_struct.LocalGroup
	return &g, utils.Wrap(d.conn.Table(constant.SuperGroupTableName).Where("group_id = ?", groupID).Take(&g).Error, "GetGroupList failed")
}

func (d *DataBase) UpdateSuperGroup(groupInfo *model_struct.LocalGroup) error {
	d.superGroupMtx.Lock()
	defer d.superGroupMtx.Unlock()

	t := d.conn.Table(constant.SuperGroupTableName).Select("*").Updates(*groupInfo)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "")
}

func (d *DataBase) DeleteSuperGroup(groupID string) error {
	d.superGroupMtx.Lock()
	defer d.superGroupMtx.Unlock()
	localGroup := model_struct.LocalGroup{GroupID: groupID}
	return utils.Wrap(d.conn.Table(constant.SuperGroupTableName).Delete(&localGroup).Error, "DeleteSuperGroup failed")
}

func (d *DataBase) GetReadDiffusionGroupIDList() ([]string, error) {
	g1, err1 := d.GetJoinedSuperGroupIDList()
	g2, err2 := d.GetJoinedWorkingGroupIDList()
	var groupIDList []string
	if err1 == nil {
		groupIDList = append(groupIDList, g1...)
	}
	if err2 == nil {
		groupIDList = append(groupIDList, g2...)
	}
	var err error
	if err1 != nil {
		err = err1
	}
	if err2 != nil {
		err = err2
	}
	return groupIDList, err
}
