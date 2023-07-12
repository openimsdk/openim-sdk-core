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
	"fmt"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) InsertFriend(ctx context.Context, friend *model_struct.LocalFriend) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Create(friend).Error, "InsertFriend failed")
}

func (d *DataBase) DeleteFriendDB(ctx context.Context, friendUserID string) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Where("owner_user_id=? and friend_user_id=?", d.loginUserID, friendUserID).Delete(&model_struct.LocalFriend{}).Error, "DeleteFriend failed")
}

func (d *DataBase) UpdateFriend(ctx context.Context, friend *model_struct.LocalFriend) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()

	t := d.conn.WithContext(ctx).Model(friend).Select("*").Updates(*friend)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "")

}

func (d *DataBase) GetAllFriendList(ctx context.Context) ([]*model_struct.LocalFriend, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friendList []model_struct.LocalFriend
	err := utils.Wrap(d.conn.WithContext(ctx).Where("owner_user_id = ?", d.loginUserID).Find(&friendList).Error,
		"GetFriendList failed")
	var transfer []*model_struct.LocalFriend
	for _, v := range friendList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, err
}

func (d *DataBase) GetPageFriendList(ctx context.Context, offset, count int) ([]*model_struct.LocalFriend, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friendList []*model_struct.LocalFriend
	err := utils.Wrap(d.conn.WithContext(ctx).Where("owner_user_id = ?", d.loginUserID).Offset(offset).Limit(count).Order("name").Find(&friendList).Error,
		"GetFriendList failed")
	return friendList, err
}

func (d *DataBase) SearchFriendList(ctx context.Context, keyword string, isSearchUserID, isSearchNickname, isSearchRemark bool) ([]*model_struct.LocalFriend, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var count int
	var friendList []model_struct.LocalFriend
	var condition string
	if isSearchUserID {
		condition = fmt.Sprintf("friend_user_id like %q ", "%"+keyword+"%")
		count++
	}
	if isSearchNickname {
		if count > 0 {
			condition += "or "
		}
		condition += fmt.Sprintf("name like %q ", "%"+keyword+"%")
		count++
	}
	if isSearchRemark {
		if count > 0 {
			condition += "or "
		}
		condition += fmt.Sprintf("remark like %q ", "%"+keyword+"%")
	}
	err := d.conn.WithContext(ctx).Where(condition).Order("create_time DESC").Find(&friendList).Error
	var transfer []*model_struct.LocalFriend
	for _, v := range friendList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "SearchFriendList failed ")

}

func (d *DataBase) GetFriendInfoByFriendUserID(ctx context.Context, FriendUserID string) (*model_struct.LocalFriend, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friend model_struct.LocalFriend
	return &friend, utils.Wrap(d.conn.WithContext(ctx).Where("owner_user_id = ? AND friend_user_id = ?",
		d.loginUserID, FriendUserID).Take(&friend).Error, "GetFriendInfoByFriendUserID failed")
}

func (d *DataBase) GetFriendInfoList(ctx context.Context, friendUserIDList []string) ([]*model_struct.LocalFriend, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friendList []model_struct.LocalFriend
	err := utils.Wrap(d.conn.WithContext(ctx).Where("friend_user_id IN ?", friendUserIDList).Find(&friendList).Error, "GetFriendInfoListByFriendUserID failed")
	var transfer []*model_struct.LocalFriend
	for _, v := range friendList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, err
}
