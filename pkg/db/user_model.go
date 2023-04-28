// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

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
	t := d.conn.WithContext(ctx).Updates(user)
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
