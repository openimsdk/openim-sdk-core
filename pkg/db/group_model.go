// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !js
// +build !js

package db

import (
	"context"
	"errors"
	"fmt"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"

	"gorm.io/gorm"
)

func (d *DataBase) InsertGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Create(groupInfo).Error, "InsertGroup failed")
}
func (d *DataBase) DeleteGroup(ctx context.Context, groupID string) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	localGroup := model_struct.LocalGroup{GroupID: groupID}
	return utils.Wrap(d.conn.WithContext(ctx).Delete(&localGroup).Error, "DeleteGroup failed")
}
func (d *DataBase) UpdateGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()

	t := d.conn.WithContext(ctx).Model(groupInfo).Select("*").Updates(*groupInfo)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "")

}
func (d *DataBase) GetJoinedGroupListDB(ctx context.Context) ([]*model_struct.LocalGroup, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	var groupList []model_struct.LocalGroup
	err := d.conn.WithContext(ctx).Find(&groupList).Error
	var transfer []*model_struct.LocalGroup
	for _, v := range groupList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetJoinedGroupList failed ")
}

func (d *DataBase) GetGroups(ctx context.Context, groupIDs []string) ([]*model_struct.LocalGroup, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	var groupList []model_struct.LocalGroup
	err := d.conn.WithContext(ctx).Where("group_id in (?)", groupIDs).Find(&groupList).Error
	var transfer []*model_struct.LocalGroup
	for _, v := range groupList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetGroups failed ")
}

func (d *DataBase) GetGroupInfoByGroupID(ctx context.Context, groupID string) (*model_struct.LocalGroup, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	var g model_struct.LocalGroup
	return &g, utils.Wrap(d.conn.WithContext(ctx).Where("group_id = ?", groupID).Take(&g).Error, "GetGroupList failed")
}
func (d *DataBase) GetAllGroupInfoByGroupIDOrGroupName(ctx context.Context, keyword string, isSearchGroupID bool, isSearchGroupName bool) ([]*model_struct.LocalGroup, error) {
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
	err := d.conn.WithContext(ctx).Where(condition).Order("create_time DESC").Find(&groupList).Error
	var transfer []*model_struct.LocalGroup
	for _, v := range groupList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetAllGroupInfoByGroupIDOrGroupName failed ")
}

func (d *DataBase) AddMemberCount(ctx context.Context, groupID string) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	group := model_struct.LocalGroup{GroupID: groupID}
	return utils.Wrap(d.conn.WithContext(ctx).Model(&group).Updates(map[string]interface{}{"member_count": gorm.Expr("member_count+1")}).Error, "")
}

func (d *DataBase) SubtractMemberCount(ctx context.Context, groupID string) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	group := model_struct.LocalGroup{GroupID: groupID}
	return utils.Wrap(d.conn.WithContext(ctx).Model(&group).Updates(map[string]interface{}{"member_count": gorm.Expr("member_count-1")}).Error, "")
}
