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

func (d *DataBase) GetLoginUser(ctx context.Context, userID string) (*model_struct.LocalUser, error) {
	d.userMtx.RLock()
	defer d.userMtx.RUnlock()
	var user model_struct.LocalUser
	return &user, utils.Wrap(d.conn.WithContext(ctx).Where("user_id = ? ", userID).Take(&user).Error, "GetLoginUserInfo failed")
}

func (d *DataBase) UpdateLoginUser(ctx context.Context, user *model_struct.LocalUser) error {
	d.userMtx.Lock()
	defer d.userMtx.Unlock()
	t := d.conn.WithContext(ctx).Model(user).Select("*").Updates(user)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateLoginUser failed")
}
func (d *DataBase) UpdateLoginUserByMap(ctx context.Context, user *model_struct.LocalUser, args map[string]interface{}) error {
	d.userMtx.Lock()
	defer d.userMtx.Unlock()
	t := d.conn.WithContext(ctx).Model(&user).Updates(args)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateColumnsConversation failed")
}
func (d *DataBase) InsertLoginUser(ctx context.Context, user *model_struct.LocalUser) error {
	d.userMtx.Lock()
	defer d.userMtx.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Create(user).Error, "InsertLoginUser failed")
}
