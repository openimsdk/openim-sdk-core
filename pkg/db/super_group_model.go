//go:build !js
// +build !js

package db

import (
	"context"
	"errors"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) GetJoinedSuperGroupList(ctx context.Context) ([]*model_struct.LocalGroup, error) {
	d.superGroupMtx.Lock()
	defer d.superGroupMtx.Unlock()
	var groupList []model_struct.LocalGroup
	err := d.conn.WithContext(ctx).Table(constant.SuperGroupTableName).Find(&groupList).Error
	var transfer []*model_struct.LocalGroup
	for _, v := range groupList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetJoinedSuperGroupList failed ")
}

func (d *DataBase) GetJoinedSuperGroupIDList(ctx context.Context) ([]string, error) {
	groupList, err := d.GetJoinedSuperGroupList(ctx)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var groupIDList []string
	for _, v := range groupList {
		groupIDList = append(groupIDList, v.GroupID)
	}
	return groupIDList, nil
}

func (d *DataBase) InsertSuperGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error {
	d.superGroupMtx.Lock()
	defer d.superGroupMtx.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Table(constant.SuperGroupTableName).Create(groupInfo).Error, "InsertSuperGroup failed")
}

func (d *DataBase) DeleteAllSuperGroup(ctx context.Context) error {
	d.superGroupMtx.Lock()
	defer d.superGroupMtx.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Table(constant.SuperGroupTableName).Delete(&model_struct.LocalGroup{}).Error, "DeleteAllSuperGroup failed")
}

func (d *DataBase) GetSuperGroupInfoByGroupID(ctx context.Context, groupID string) (*model_struct.LocalGroup, error) {
	d.superGroupMtx.Lock()
	defer d.superGroupMtx.Unlock()
	var g model_struct.LocalGroup
	return &g, utils.Wrap(d.conn.WithContext(ctx).Table(constant.SuperGroupTableName).Where("group_id = ?", groupID).Take(&g).Error, "GetGroupList failed")
}

func (d *DataBase) UpdateSuperGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error {
	d.superGroupMtx.Lock()
	defer d.superGroupMtx.Unlock()

	t := d.conn.WithContext(ctx).Table(constant.SuperGroupTableName).Select("*").Updates(*groupInfo)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "")
}

func (d *DataBase) DeleteSuperGroup(ctx context.Context, groupID string) error {
	d.superGroupMtx.Lock()
	defer d.superGroupMtx.Unlock()
	localGroup := model_struct.LocalGroup{GroupID: groupID}
	return utils.Wrap(d.conn.WithContext(ctx).Table(constant.SuperGroupTableName).Delete(&localGroup).Error, "DeleteSuperGroup failed")
}

func (d *DataBase) GetReadDiffusionGroupIDList(ctx context.Context) ([]string, error) {
	sg, err := d.GetJoinedSuperGroupIDList(ctx)
	if err != nil {
		return nil, err
	}
	wg, err := d.GetJoinedWorkingGroupIDList(ctx)
	if err != nil {
		return nil, err
	}
	return append(sg, wg...), err
}
