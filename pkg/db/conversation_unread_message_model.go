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

func (d *DataBase) BatchInsertConversationUnreadMessageList(ctx context.Context, messageList []*model_struct.LocalConversationUnreadMessage) error {
	if messageList == nil {
		return nil
	}
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Create(messageList).Error, "BatchInsertConversationUnreadMessageList failed")
}
func (d *DataBase) DeleteConversationUnreadMessageList(ctx context.Context, conversationID string, sendTime int64) int64 {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return d.conn.WithContext(ctx).Where("conversation_id = ? and send_time <= ?", conversationID, sendTime).Delete(&model_struct.LocalConversationUnreadMessage{}).RowsAffected
}
