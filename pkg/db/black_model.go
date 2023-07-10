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
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) GetBlackListDB(ctx context.Context) ([]*model_struct.LocalBlack, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	if d == nil {
		return nil, errors.New("database is not open")
	}
	var blackList []model_struct.LocalBlack

	err := d.conn.WithContext(ctx).Find(&blackList).Error
	var transfer []*model_struct.LocalBlack
	for _, v := range blackList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, err
}

func (d *DataBase) GetBlackListUserID(ctx context.Context) (blackListUid []string, err error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	return blackListUid, utils.Wrap(d.conn.WithContext(ctx).Model(&model_struct.LocalBlack{}).Select("block_user_id").Find(&blackListUid).Error, "GetBlackList failed")
}

func (d *DataBase) GetBlackInfoByBlockUserID(ctx context.Context, blockUserID string) (*model_struct.LocalBlack, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var black model_struct.LocalBlack
	return &black, utils.Wrap(d.conn.WithContext(ctx).Where("owner_user_id = ? AND block_user_id = ? ",
		d.loginUserID, blockUserID).Take(&black).Error, "GetBlackInfoByBlockUserID failed")
}

func (d *DataBase) GetBlackInfoList(ctx context.Context, blockUserIDList []string) ([]*model_struct.LocalBlack, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var blackList []model_struct.LocalBlack
	if t := d.conn.WithContext(ctx).Where("block_user_id IN ? ", blockUserIDList).Find(&blackList).Error; t != nil {
		return nil, utils.Wrap(t, "GetBlackInfoList failed")
	}

	var transfer []*model_struct.LocalBlack
	for _, v := range blackList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, nil
}

func (d *DataBase) InsertBlack(ctx context.Context, black *model_struct.LocalBlack) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Create(black).Error, "InsertBlack failed")
}

func (d *DataBase) UpdateBlack(ctx context.Context, black *model_struct.LocalBlack) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	t := d.conn.WithContext(ctx).Updates(black)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateBlack failed")
}

func (d *DataBase) DeleteBlack(ctx context.Context, blockUserID string) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Where("owner_user_id=? and block_user_id=?", d.loginUserID, blockUserID).Delete(&model_struct.LocalBlack{}).Error, "DeleteBlack failed")
}
