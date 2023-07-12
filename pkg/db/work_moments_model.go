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
	"gorm.io/gorm"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
	"time"
)

func (d *DataBase) InsertWorkMomentsNotification(ctx context.Context, jsonDetail string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	workMomentsNotification := model_struct.LocalWorkMomentsNotification{
		JsonDetail: jsonDetail,
		CreateTime: time.Now().Unix(),
	}
	return utils.Wrap(d.conn.WithContext(ctx).Create(workMomentsNotification).Error, "")
}

func (d *DataBase) GetWorkMomentsNotification(ctx context.Context, offset, count int) (WorkMomentsNotifications []*model_struct.LocalWorkMomentsNotification, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	WorkMomentsNotifications = []*model_struct.LocalWorkMomentsNotification{}
	err = utils.Wrap(d.conn.WithContext(ctx).Table("local_work_moments_notification").Order("create_time DESC").Offset(offset).Limit(count).Find(&WorkMomentsNotifications).Error, "")
	return WorkMomentsNotifications, err
}

func (d *DataBase) GetWorkMomentsNotificationLimit(ctx context.Context, pageNumber, showNumber int) (WorkMomentsNotifications []*model_struct.LocalWorkMomentsNotification, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	WorkMomentsNotifications = []*model_struct.LocalWorkMomentsNotification{}
	err = utils.Wrap(d.conn.WithContext(ctx).Table("local_work_moments_notification").Select("json_detail").Find(WorkMomentsNotifications).Error, "")
	return WorkMomentsNotifications, err
}

func (d *DataBase) InitWorkMomentsNotificationUnreadCount(ctx context.Context) error {
	var n int64
	err := utils.Wrap(d.conn.WithContext(ctx).Model(&model_struct.LocalWorkMomentsNotificationUnreadCount{}).Count(&n).Error, "")
	if err == nil {
		if n == 0 {
			c := model_struct.LocalWorkMomentsNotificationUnreadCount{UnreadCount: 0}
			return utils.Wrap(d.conn.WithContext(ctx).Model(&model_struct.LocalWorkMomentsNotificationUnreadCount{}).Create(c).Error, "IncrConversationUnreadCount failed")
		}
	}
	return err
}

func (d *DataBase) IncrWorkMomentsNotificationUnreadCount(ctx context.Context) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	c := model_struct.LocalWorkMomentsNotificationUnreadCount{}
	t := d.conn.WithContext(ctx).Model(&c).Where("1=1").Update("unread_count", gorm.Expr("unread_count+?", 1))
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "IncrConversationUnreadCount failed")
}

func (d *DataBase) MarkAllWorkMomentsNotificationAsRead(ctx context.Context) (err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Model(&model_struct.LocalWorkMomentsNotificationUnreadCount{}).Where("1=1").Updates(map[string]interface{}{"unread_count": 0}).Error, "")
}

func (d *DataBase) GetWorkMomentsUnReadCount(ctx context.Context) (workMomentsNotificationUnReadCount model_struct.LocalWorkMomentsNotificationUnreadCount, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	workMomentsNotificationUnReadCount = model_struct.LocalWorkMomentsNotificationUnreadCount{}
	err = utils.Wrap(d.conn.WithContext(ctx).Model(&model_struct.LocalWorkMomentsNotificationUnreadCount{}).First(&workMomentsNotificationUnReadCount).Error, "")
	return workMomentsNotificationUnReadCount, err
}

func (d *DataBase) ClearWorkMomentsNotification(ctx context.Context) (err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Table("local_work_moments_notification").Where("1=1").Delete(&model_struct.LocalWorkMomentsNotification{}).Error, "")
}
