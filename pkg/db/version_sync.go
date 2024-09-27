// Copyright © 2024 OpenIM SDK. All rights reserved.
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

	"gorm.io/gorm"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/errs"
)

func (d *DataBase) GetVersionSync(ctx context.Context, tableName, entityID string) (*model_struct.LocalVersionSync, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var res model_struct.LocalVersionSync
	err := d.conn.WithContext(ctx).Where("`table_name` = ? and `entity_id` = ?", tableName, entityID).Take(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model_struct.LocalVersionSync{}, errs.ErrRecordNotFound.Wrap()
		}
		return nil, errs.Wrap(err)
	}
	return &res, errs.Wrap(err)
}

func (d *DataBase) SetVersionSync(ctx context.Context, lv *model_struct.LocalVersionSync) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()

	var existing model_struct.LocalVersionSync
	err := d.conn.WithContext(ctx).Where("`table_name` = ? AND `entity_id` = ?", lv.Table, lv.EntityID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		if createErr := d.conn.WithContext(ctx).Create(lv).Error; createErr != nil {
			return errs.Wrap(createErr)
		}
		return nil
	} else if err != nil {
		return errs.Wrap(err)
	}

	if updateErr := d.conn.WithContext(ctx).Model(&existing).Updates(lv).Error; updateErr != nil {
		return errs.Wrap(updateErr)
	}

	return nil
}

func (d *DataBase) DeleteVersionSync(ctx context.Context, tableName, entityID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	localVersionSync := model_struct.LocalVersionSync{Table: tableName, EntityID: entityID}
	return errs.WrapMsg(d.conn.WithContext(ctx).Delete(&localVersionSync).Error, "DeleteVersionSync failed")
}
