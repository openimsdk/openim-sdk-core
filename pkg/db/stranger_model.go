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
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"gorm.io/gorm"
)

func (d *DataBase) GetStrangerInfo(ctx context.Context, userIDs []string) ([]*model_struct.LocalStranger, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friendList []model_struct.LocalStranger
	err := utils.Wrap(d.conn.WithContext(ctx).Where("user_id IN ?", userIDs).Find(&friendList).Error, "GetFriendInfoListByFriendUserID failed")
	var transfer []*model_struct.LocalStranger
	for _, v := range friendList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, err
}

func (d *DataBase) SetStrangerInfo(ctx context.Context, localStrangerList []*model_struct.LocalStranger) error {
	//TODO Can be optimized into two chan batch update or insert operations
	for _, existingData := range localStrangerList {
		result := d.conn.First(&existingData, "user_id = ?", existingData.UserID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Data does not exist, perform insert operation
			err := d.conn.Create(&existingData).Error
			return err
		} else if result.Error != nil {
			return result.Error
		}
		err := d.conn.Save(&existingData).Error
		if err != nil {
			return err
		}
	}
	return nil
}
