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

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/errs"
)

func (d *DataBase) GetBlackListDB(ctx context.Context) ([]*model_struct.LocalBlack, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var blackList []*model_struct.LocalBlack
	return blackList, errs.Wrap(d.conn.WithContext(ctx).Find(&blackList).Error)
}

func (d *DataBase) GetBlackListUserID(ctx context.Context) (blackListUid []string, err error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	return blackListUid, errs.WrapMsg(d.conn.WithContext(ctx).Model(&model_struct.LocalBlack{}).Select("block_user_id").Find(&blackListUid).Error, "GetBlackList failed")
}

func (d *DataBase) GetBlackInfoByBlockUserID(ctx context.Context, blockUserID string) (*model_struct.LocalBlack, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var black model_struct.LocalBlack
	return &black, errs.WrapMsg(d.conn.WithContext(ctx).Where("owner_user_id = ? AND block_user_id = ? ",
		d.loginUserID, blockUserID).Take(&black).Error, "GetBlackInfoByBlockUserID failed")
}

func (d *DataBase) GetBlackInfoList(ctx context.Context, blockUserIDList []string) ([]*model_struct.LocalBlack, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var blackList []*model_struct.LocalBlack
	if err := d.conn.WithContext(ctx).Where("block_user_id IN ? ", blockUserIDList).Find(&blackList).Error; err != nil {
		return nil, errs.WrapMsg(err, "GetBlackInfoList failed")
	}
	return blackList, nil
}

func (d *DataBase) InsertBlack(ctx context.Context, black *model_struct.LocalBlack) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Create(black).Error, "InsertBlack failed")
}

func (d *DataBase) UpdateBlack(ctx context.Context, black *model_struct.LocalBlack) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Updates(black)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.WrapMsg(t.Error, "UpdateBlack failed")
}

func (d *DataBase) DeleteBlack(ctx context.Context, blockUserID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Where("owner_user_id=? and block_user_id=?", d.loginUserID, blockUserID).Delete(&model_struct.LocalBlack{}).Error, "DeleteBlack failed")
}
