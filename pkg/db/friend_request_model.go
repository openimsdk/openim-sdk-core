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

func (d *DataBase) InsertFriendRequest(ctx context.Context, friendRequest *model_struct.LocalFriendRequest) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Create(friendRequest).Error, "InsertFriendRequest failed")
}

func (d *DataBase) DeleteFriendRequestBothUserID(ctx context.Context, fromUserID, toUserID string) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Where("from_user_id=? and to_user_id=?", fromUserID, toUserID).Delete(&model_struct.LocalFriendRequest{}).Error, "DeleteFriendRequestBothUserID failed")
}

func (d *DataBase) UpdateFriendRequest(ctx context.Context, friendRequest *model_struct.LocalFriendRequest) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	t := d.conn.WithContext(ctx).Model(friendRequest).Select("*").Updates(*friendRequest)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "")
}

func (d *DataBase) GetRecvFriendApplication(ctx context.Context) ([]*model_struct.LocalFriendRequest, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friendRequestList []model_struct.LocalFriendRequest
	err := utils.Wrap(d.conn.WithContext(ctx).Where("to_user_id = ?", d.loginUserID).Order("create_time DESC").Find(&friendRequestList).Error, "GetRecvFriendApplication failed")

	var transfer []*model_struct.LocalFriendRequest
	for _, v := range friendRequestList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetRecvFriendApplication failed")
}

func (d *DataBase) GetSendFriendApplication(ctx context.Context) ([]*model_struct.LocalFriendRequest, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friendRequestList []model_struct.LocalFriendRequest
	err := utils.Wrap(d.conn.WithContext(ctx).Where("from_user_id = ?", d.loginUserID).Order("create_time DESC").Find(&friendRequestList).Error, "GetSendFriendApplication failed")

	var transfer []*model_struct.LocalFriendRequest
	for _, v := range friendRequestList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetSendFriendApplication failed")
}

func (d *DataBase) GetFriendApplicationByBothID(ctx context.Context, fromUserID, toUserID string) (*model_struct.LocalFriendRequest, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friendRequest model_struct.LocalFriendRequest
	err := utils.Wrap(d.conn.WithContext(ctx).Where("from_user_id = ? AND to_user_id = ?", fromUserID, toUserID).Take(&friendRequest).Error, "GetFriendApplicationByBothID failed")
	return &friendRequest, utils.Wrap(err, "GetFriendApplicationByBothID failed")
}

func (d *DataBase) GetBothFriendReq(ctx context.Context, fromUserID, toUserID string) (friendRequests []*model_struct.LocalFriendRequest, err error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	err = utils.Wrap(d.conn.WithContext(ctx).Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)", fromUserID, toUserID, toUserID, fromUserID).Find(&friendRequests).Error, "GetFriendApplicationByBothID failed")
	return friendRequests, utils.Wrap(err, "GetFriendApplicationByBothID failed")
}
