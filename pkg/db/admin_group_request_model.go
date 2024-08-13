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

func (d *DataBase) InsertAdminGroupRequest(ctx context.Context, groupRequest *model_struct.LocalAdminGroupRequest) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Create(groupRequest).Error, "InsertAdminGroupRequest failed")
}

func (d *DataBase) DeleteAdminGroupRequest(ctx context.Context, groupID, userID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Where("group_id=? and user_id=?", groupID, userID).Delete(&model_struct.LocalAdminGroupRequest{}).Error, "DeleteAdminGroupRequest failed")
}

func (d *DataBase) UpdateAdminGroupRequest(ctx context.Context, groupRequest *model_struct.LocalAdminGroupRequest) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Model(groupRequest).Select("*").Updates(*groupRequest)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.Wrap(t.Error)
}

func (d *DataBase) GetAdminGroupApplication(ctx context.Context) ([]*model_struct.LocalAdminGroupRequest, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var groupRequestList []*model_struct.LocalAdminGroupRequest
	err := errs.Wrap(d.conn.WithContext(ctx).Order("create_time DESC").Find(&groupRequestList).Error)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return groupRequestList, nil
}
