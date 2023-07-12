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

func (d *DataBase) SetNotificationSeq(ctx context.Context, conversationID string, seq int64) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	cursor := d.conn.WithContext(ctx).Model(&model_struct.NotificationSeqs{}).Where("conversation_id = ?", conversationID).Updates(map[string]interface{}{"seq": seq})
	if cursor.Error != nil {
		return utils.Wrap(cursor.Error, "Updates failed")
	}
	if cursor.RowsAffected == 0 {
		return utils.Wrap(d.conn.WithContext(ctx).Create(&model_struct.NotificationSeqs{ConversationID: conversationID, Seq: seq}).Error, "Create failed")
	}
	return nil
}

func (d *DataBase) GetNotificationAllSeqs(ctx context.Context) ([]*model_struct.NotificationSeqs, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seqs []*model_struct.NotificationSeqs
	return seqs, utils.Wrap(d.conn.WithContext(ctx).Where("1=1").Find(&seqs).Error, "GetNotificationAllSeqs failed")
}
