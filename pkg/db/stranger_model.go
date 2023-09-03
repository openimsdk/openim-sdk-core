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
	"open_im_sdk/pkg/utils"
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
	err := utils.Wrap(d.conn.Where("1 = 1").Delete(&model_struct.LocalStranger{}).Error, "Delete LocalStrangers failed")
	if err != nil {
		return err
	}
	err = utils.Wrap(d.conn.Create(localStrangerList).Error, "Creat LocalStrangers failed")
	if err != nil {
		return err
	}
	return nil
}
