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

package db

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"open_im_sdk/pkg/db/model_struct"
	"time"
)

func (d *DataBase) GetUpload(ctx context.Context, partHash string) (*model_struct.Upload, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var upload model_struct.Upload
	err := d.conn.WithContext(ctx).Where("part_hash = ?", partHash).Take(&upload).Error
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &upload, nil
}

func (d *DataBase) InsertUpload(ctx context.Context, upload *model_struct.Upload) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.Wrap(d.conn.WithContext(ctx).Create(upload).Error)
}

func (d *DataBase) deleteUpload(ctx context.Context, partHash string) error {
	err := d.conn.WithContext(ctx).Where("part_hash = ?", partHash).Delete(&model_struct.UploadPart{}).Error
	if err != nil {
		return err
	}
	return errs.Wrap(d.conn.WithContext(ctx).Where("part_hash = ?", partHash).Delete(&model_struct.Upload{}).Error)
}

func (d *DataBase) DeleteUpload(ctx context.Context, partHash string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return d.deleteUpload(ctx, partHash)
}

func (d *DataBase) DeleteExpireUpload(ctx context.Context) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var uploads []*model_struct.Upload
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

func (d *DataBase) GetUploadPart(ctx context.Context, partHash string) ([]int32, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var index []int32
	err := d.conn.WithContext(ctx).Model(&model_struct.UploadPart{}).Where("part_hash = ?", partHash).Pluck("index", &index).Error
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return index, nil
}

func (d *DataBase) SetUploadPartPush(ctx context.Context, partHash string, index []int32) error {
	if len(index) == 0 {
		return nil
	}
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	uploadParts := make([]*model_struct.UploadPart, 0, len(index))
	for _, idx := range index {
		uploadParts = append(uploadParts, &model_struct.UploadPart{
			PartHash: partHash,
			Index:    idx,
		})
	}
	return errs.Wrap(d.conn.WithContext(ctx).Create(&uploadParts).Error)
}
