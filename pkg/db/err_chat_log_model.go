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
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) initSuperLocalErrChatLog(ctx context.Context, groupID string) {
	if !d.conn.WithContext(ctx).Migrator().HasTable(utils.GetErrTableName(groupID)) {
		d.conn.WithContext(ctx).Table(utils.GetErrTableName(groupID)).AutoMigrate(&model_struct.LocalErrChatLog{})
	}
}
func (d *DataBase) SuperBatchInsertExceptionMsg(ctx context.Context, MessageList []*model_struct.LocalErrChatLog, groupID string) error {
	if MessageList == nil {
		return nil
	}
	d.initSuperLocalErrChatLog(ctx, groupID)
	return utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetConversationTableName(groupID)).Create(MessageList).Error, "BatchInsertMessageList failed")
}
func (d *DataBase) GetAbnormalMsgSeq(ctx context.Context) (int64, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seq int64
	err := d.conn.WithContext(ctx).Model(model_struct.LocalErrChatLog{}).Select("IFNULL(max(seq),0)").Find(&seq).Error
	return seq, utils.Wrap(err, "GetAbnormalMsgSeq")
}
func (d *DataBase) GetAbnormalMsgSeqList(ctx context.Context) ([]int64, error) {
	var seqList []int64
	err := d.conn.WithContext(ctx).Model(model_struct.LocalErrChatLog{}).Select("seq").Find(&seqList).Error
	return seqList, utils.Wrap(err, "GetAbnormalMsgSeqList")
}
func (d *DataBase) BatchInsertExceptionMsg(ctx context.Context, messageList []*model_struct.LocalErrChatLog) error {
	if messageList == nil {
		return nil
	}
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Create(messageList).Error, "BatchInsertMessageList failed")
}
func (d *DataBase) BatchInsertExceptionMsgController(ctx context.Context, messageList []*model_struct.LocalErrChatLog) error {
	if len(messageList) == 0 {
		return nil
	}
	switch messageList[len(messageList)-1].SessionType {
	case constant.SuperGroupChatType:
		return d.SuperBatchInsertExceptionMsg(ctx, messageList, messageList[len(messageList)-1].RecvID)
	default:
		return d.BatchInsertExceptionMsg(ctx, messageList)
	}
}
func (d *DataBase) GetConversationAbnormalMsgSeq(ctx context.Context, conversationID string) (int64, error) {
	var seq int64
	if !d.conn.WithContext(ctx).Migrator().HasTable(utils.GetErrTableName(conversationID)) {
		return 0, nil
	}
	err := d.conn.WithContext(ctx).Table(utils.GetErrTableName(conversationID)).Select("IFNULL(max(seq),0)").Find(&seq).Error
	return seq, utils.Wrap(err, "GetConversationNormalMsgSeq")
}
