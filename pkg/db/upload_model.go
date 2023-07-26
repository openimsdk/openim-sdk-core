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
	"open_im_sdk/pkg/db/model_struct"
	"time"

	"github.com/OpenIMSDK/tools/errs"
)

func (d *DataBase) GetUpload(ctx context.Context, partHash string) (*model_struct.LocalUpload, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var upload model_struct.LocalUpload
	err := d.conn.WithContext(ctx).Where("part_hash = ?", partHash).Take(&upload).Error
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &upload, nil
}

func (d *DataBase) InsertUpload(ctx context.Context, upload *model_struct.LocalUpload) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.Wrap(d.conn.WithContext(ctx).Create(upload).Error)
}

func (d *DataBase) deleteUpload(ctx context.Context, partHash string) error {
	return errs.Wrap(d.conn.WithContext(ctx).Where("part_hash = ?", partHash).Delete(&model_struct.LocalUpload{}).Error)
}

func (d *DataBase) UpdateUpload(ctx context.Context, upload *model_struct.LocalUpload) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	return errs.Wrap(d.conn.WithContext(ctx).Updates(upload).Error)
}

func (d *DataBase) DeleteUpload(ctx context.Context, partHash string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return d.deleteUpload(ctx, partHash)
}

func (d *DataBase) DeleteExpireUpload(ctx context.Context) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var uploads []*model_struct.LocalUpload
	err := d.conn.WithContext(ctx).Where("expire_time <= ?", time.Now().UnixMilli()).Find(&uploads).Error
	if err != nil {
		return errs.Wrap(err)
	}
	for _, upload := range uploads {
		if err := d.deleteUpload(ctx, upload.PartHash); err != nil {
			return err
		}
	}
	return nil
}
