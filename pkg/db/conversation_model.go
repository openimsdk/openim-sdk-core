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

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"

	"gorm.io/gorm"

	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

const (
	batchSize = 200
)

func (d *DataBase) GetConversationByUserID(ctx context.Context, userID string) (*model_struct.LocalConversation, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var conversation model_struct.LocalConversation
	err := errs.WrapMsg(d.conn.WithContext(ctx).Where("user_id=?", userID).Find(&conversation).Error, "GetConversationByUserID error")
	return &conversation, err
}

func (d *DataBase) GetAllConversationListDB(ctx context.Context) ([]*model_struct.LocalConversation, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var conversationList []*model_struct.LocalConversation
	err := errs.WrapMsg(d.conn.WithContext(ctx).Where("latest_msg_send_time > ?", 0).Order("case when is_pinned=1 then 0 else 1 end,max(latest_msg_send_time,draft_text_time) DESC").Find(&conversationList).Error,
		"GetAllConversationList failed")
	if err != nil {
		return nil, err
	}
	return conversationList, err
}

func (d *DataBase) FindAllConversationConversationID(ctx context.Context) (conversationIDs []string, err error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	return conversationIDs, errs.WrapMsg(d.conn.WithContext(ctx).Model(&model_struct.LocalConversation{}).Where("latest_msg_send_time > ?", 0).Pluck("conversation_id", &conversationIDs).Error, "")
}

func (d *DataBase) FindAllUnreadConversationConversationID(ctx context.Context) (conversationIDs []string, err error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	return conversationIDs, errs.WrapMsg(d.conn.WithContext(ctx).Model(&model_struct.LocalConversation{}).Where("unread_count > ?", 0).Pluck("conversation_id", &conversationIDs).Error, "")
}

func (d *DataBase) GetHiddenConversationList(ctx context.Context) ([]*model_struct.LocalConversation, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var conversationList []*model_struct.LocalConversation
	return conversationList, errs.WrapMsg(d.conn.WithContext(ctx).Where("latest_msg_send_time = ?", 0).Find(&conversationList).Error,
		"GetHiddenConversationList failed")
}

func (d *DataBase) GetAllConversations(ctx context.Context) ([]*model_struct.LocalConversation, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var conversationList []*model_struct.LocalConversation
	return conversationList, errs.WrapMsg(d.conn.WithContext(ctx).Find(&conversationList).Error, "GetAllConversations failed")
}

func (d *DataBase) GetAllConversationIDList(ctx context.Context) (result []string, err error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var c model_struct.LocalConversation
	return result, errs.WrapMsg(d.conn.WithContext(ctx).Model(&c).Pluck("conversation_id", &result).Error, "GetAllConversationIDList failed ")
}

func (d *DataBase) GetAllSingleConversationIDList(ctx context.Context) (result []string, err error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var c model_struct.LocalConversation
	return result, errs.WrapMsg(d.conn.WithContext(ctx).Model(&c).Where("conversation_type = ?", constant.SingleChatType).Pluck("conversation_id", &result).Error, "GetAllConversationIDList failed ")
}

func (d *DataBase) GetConversationListSplitDB(ctx context.Context, offset, count int) ([]*model_struct.LocalConversation, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var conversationList []*model_struct.LocalConversation
	return conversationList, errs.Wrap(d.conn.WithContext(ctx).Where("latest_msg_send_time > ?", 0).Order("case when is_pinned=1 then 0 else 1 end,max(latest_msg_send_time,draft_text_time) DESC").Offset(offset).Limit(count).Find(&conversationList).Error)
}

func (d *DataBase) BatchInsertConversationList(ctx context.Context, conversationList []*model_struct.LocalConversation) error {
	if conversationList == nil {
		return nil
	}

	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()

	for i := 0; i < len(conversationList); i += batchSize {
		end := i + batchSize
		if end > len(conversationList) {
			end = len(conversationList)
		}

		batch := conversationList[i:end]
		if err := d.conn.WithContext(ctx).Create(batch).Error; err != nil {
			return errs.WrapMsg(err, "BatchInsertConversationList failed")
		}
	}

	return nil
}

func (d *DataBase) UpdateOrCreateConversations(ctx context.Context, conversationList []*model_struct.LocalConversation) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var conversationIDs []string
	if err := d.conn.WithContext(ctx).Model(&model_struct.LocalConversation{}).Pluck("conversation_id", &conversationIDs).Error; err != nil {
		return err
	}
	var notExistConversations []*model_struct.LocalConversation
	var existConversations []*model_struct.LocalConversation
	for i, v := range conversationList {
		if utils.IsContain(v.ConversationID, conversationIDs) {
			existConversations = append(existConversations, v)
			continue
		} else {
			notExistConversations = append(notExistConversations, conversationList[i])
		}
	}
	if len(notExistConversations) > 0 {
		if err := d.conn.WithContext(ctx).Create(notExistConversations).Error; err != nil {
			return err
		}
	}
	for _, v := range existConversations {
		if err := d.conn.WithContext(ctx).Model(&model_struct.LocalConversation{}).Where("conversation_id = ?", v.ConversationID).Updates(map[string]interface{}{"unread_count": v.UnreadCount}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d *DataBase) InsertConversation(ctx context.Context, conversationList *model_struct.LocalConversation) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Create(conversationList).Error, "InsertConversation failed")
}

func (d *DataBase) DeleteConversation(ctx context.Context, conversationID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Where("conversation_id = ?", conversationID).Delete(&model_struct.LocalConversation{}).Error, "DeleteConversation failed")
}

func (d *DataBase) DeleteAllConversation(ctx context.Context) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model_struct.LocalConversation{}).Error, "DeleteAllConversation failed")
}

func (d *DataBase) GetConversation(ctx context.Context, conversationID string) (*model_struct.LocalConversation, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var c model_struct.LocalConversation
	return &c, errs.WrapMsg(d.conn.WithContext(ctx).Where("conversation_id = ?",
		conversationID).Take(&c).Error, "GetConversation failed, conversationID: "+conversationID)
}

func (d *DataBase) UpdateConversation(ctx context.Context, c *model_struct.LocalConversation) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	d.conn.WithContext(ctx).Logger.LogMode(6)
	t := d.conn.WithContext(ctx).Updates(c)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.WrapMsg(t.Error, "UpdateConversation failed")
}

func (d *DataBase) UpdateConversationForSync(ctx context.Context, c *model_struct.LocalConversation) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Model(&model_struct.LocalConversation{}).Where("conversation_id = ?", c.ConversationID).
		Updates(map[string]interface{}{"recv_msg_opt": c.RecvMsgOpt, "is_pinned": c.IsPinned, "is_private_chat": c.IsPrivateChat,
			"group_at_type": c.GroupAtType, "is_not_in_group": c.IsNotInGroup, "update_unread_count_time": c.UpdateUnreadCountTime, "ex": c.Ex, "attached_info": c.AttachedInfo,
			"burn_duration": c.BurnDuration, "msg_destruct_time": c.MsgDestructTime, "is_msg_destruct": c.IsMsgDestruct})
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.WrapMsg(t.Error, "UpdateConversation failed")
}

func (d *DataBase) BatchUpdateConversationList(ctx context.Context, conversationList []*model_struct.LocalConversation) error {
	for _, v := range conversationList {
		err := d.UpdateConversation(ctx, v)
		if err != nil {
			return errs.WrapMsg(err, "BatchUpdateConversationList failed")
		}

	}
	return nil
}

func (d *DataBase) ConversationIfExists(ctx context.Context, conversationID string) (bool, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var count int64
	t := d.conn.WithContext(ctx).Model(&model_struct.LocalConversation{}).Where("conversation_id = ?",
		conversationID).Count(&count)
	if t.Error != nil {
		return false, errs.WrapMsg(t.Error, "ConversationIfExists get failed")
	}
	if count != 1 {
		return false, nil
	} else {
		return true, nil
	}
}

// Reset the conversation is equivalent to deleting the conversation,
// and the GetAllConversation or GetConversationListSplit interface will no longer be obtained.
func (d *DataBase) ResetConversation(ctx context.Context, conversationID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	c := model_struct.LocalConversation{ConversationID: conversationID, UnreadCount: 0, LatestMsg: "", LatestMsgSendTime: 0, DraftText: "", DraftTextTime: 0}
	t := d.conn.WithContext(ctx).Select("unread_count", "latest_msg", "latest_msg_send_time", "draft_text", "draft_text_time").Updates(c)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.WrapMsg(t.Error, "ResetConversation failed")
}

// ResetAllConversation Reset ALL conversation is equivalent to deleting the conversation,
// and the GetAllConversation or GetConversationListSplit interface will no longer be obtained.
func (d *DataBase) ResetAllConversation(ctx context.Context) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	c := model_struct.LocalConversation{UnreadCount: 0, LatestMsg: "", LatestMsgSendTime: 0, DraftText: "", DraftTextTime: 0}
	t := d.conn.WithContext(ctx).Session(&gorm.Session{AllowGlobalUpdate: true}).Select("unread_count", "latest_msg", "latest_msg_send_time", "draft_text", "draft_text_time").Updates(c)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.WrapMsg(t.Error, "ResetConversation failed")
}

// Clear the conversation, which is used to delete the conversation history message and clear the conversation at the same time.
// The GetAllConversation or GetConversationListSplit interface can still be obtained,
// but there is no latest message.
func (d *DataBase) ClearConversation(ctx context.Context, conversationID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	c := model_struct.LocalConversation{ConversationID: conversationID, UnreadCount: 0, LatestMsg: "", DraftText: "", DraftTextTime: 0}
	t := d.conn.WithContext(ctx).Select("unread_count", "latest_msg", "draft_text", "draft_text_time").Updates(c)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.WrapMsg(t.Error, "ClearConversation failed")
}

func (d *DataBase) SetConversationDraftDB(ctx context.Context, conversationID, draftText string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	nowTime := utils.GetCurrentTimestampByMill()
	t := d.conn.WithContext(ctx).Exec("update local_conversations set draft_text=?,draft_text_time=?,latest_msg_send_time=case when latest_msg_send_time=? then ? else latest_msg_send_time  end where conversation_id=?",
		draftText, nowTime, 0, nowTime, conversationID)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.WrapMsg(t.Error, "SetConversationDraft failed")
}

func (d *DataBase) RemoveConversationDraft(ctx context.Context, conversationID, draftText string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	c := model_struct.LocalConversation{ConversationID: conversationID, DraftText: draftText, DraftTextTime: 0}
	t := d.conn.WithContext(ctx).Select("draft_text", "draft_text_time").Updates(c)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.WrapMsg(t.Error, "RemoveConversationDraft failed")
}

func (d *DataBase) UnPinConversation(ctx context.Context, conversationID string, isPinned int) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Exec("update local_conversations set is_pinned=?,draft_text_time=case when draft_text=? then ? else draft_text_time  end where conversation_id=?",
		isPinned, "", 0, conversationID)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.WrapMsg(t.Error, "UnPinConversation failed")
}

func (d *DataBase) UpdateColumnsConversation(ctx context.Context, conversationID string, args map[string]interface{}) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Model(model_struct.LocalConversation{ConversationID: conversationID}).Updates(args)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errs.ErrRecordNotFound, "no update")
	}
	return errs.WrapMsg(t.Error, "UpdateColumnsConversation failed")
}

func (d *DataBase) UpdateAllConversation(ctx context.Context, conversation *model_struct.LocalConversation) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	if conversation.ConversationID != "" {
		return errs.WrapMsg(errors.New("not update all conversation"), "UpdateAllConversation failed")
	}
	t := d.conn.WithContext(ctx).Model(conversation).Updates(conversation)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.WrapMsg(t.Error, "UpdateColumnsConversation failed")
}

func (d *DataBase) IncrConversationUnreadCount(ctx context.Context, conversationID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	c := model_struct.LocalConversation{ConversationID: conversationID}
	t := d.conn.WithContext(ctx).Model(&c).Update("unread_count", gorm.Expr("unread_count+?", 1))
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.WrapMsg(t.Error, "IncrConversationUnreadCount failed")
}

func (d *DataBase) GetTotalUnreadMsgCountDB(ctx context.Context) (totalUnreadCount int32, err error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var result []int64
	err = d.conn.WithContext(ctx).Model(&model_struct.LocalConversation{}).Where("recv_msg_opt < ? and latest_msg_send_time > ?", constant.ReceiveNotNotifyMessage, 0).Pluck("unread_count", &result).Error
	if err != nil {
		return totalUnreadCount, errs.WrapMsg(errors.New("GetTotalUnreadMsgCount err"), "GetTotalUnreadMsgCount err")
	}
	for _, v := range result {
		totalUnreadCount += int32(v)
	}
	return totalUnreadCount, nil
}

func (d *DataBase) SetMultipleConversationRecvMsgOpt(ctx context.Context, conversationIDList []string, opt int) (err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Model(&model_struct.LocalConversation{}).Where("conversation_id IN ?", conversationIDList).Updates(map[string]interface{}{"recv_msg_opt": opt})
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.WrapMsg(t.Error, "SetMultipleConversationRecvMsgOpt failed")
}

func (d *DataBase) GetMultipleConversationDB(ctx context.Context, conversationIDList []string) (result []*model_struct.LocalConversation, err error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var conversationList []model_struct.LocalConversation
	err = errs.WrapMsg(d.conn.WithContext(ctx).Where("conversation_id IN ?", conversationIDList).Find(&conversationList).Error, "GetMultipleConversation failed")
	for _, v := range conversationList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}

func (d *DataBase) DecrConversationUnreadCount(ctx context.Context, conversationID string, count int64) (err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	tx := d.conn.WithContext(ctx).Begin()
	c := model_struct.LocalConversation{ConversationID: conversationID}
	t := tx.Model(&c).Update("unread_count", gorm.Expr("unread_count-?", count))
	if t.Error != nil {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	if err := tx.Where("conversation_id = ?",
		conversationID).Take(&c).Error; err != nil {
		tx.Rollback()
		return errs.WrapMsg(errors.New("get conversation err"), "")
	}
	if c.UnreadCount < 0 {
		log.ZWarn(ctx, "decr unread count < 0", nil, "conversationID", conversationID, "count", count)
		if t = tx.Model(&c).Update("unread_count", 0); t.Error != nil {
			tx.Rollback()
			return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
		}
	}
	tx.Commit()
	return nil
}

func (d *DataBase) SearchConversations(ctx context.Context, searchParam string) ([]*model_struct.LocalConversation, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	// Define the search condition based on the searchParam
	condition := fmt.Sprintf("show_name like %q ", "%"+searchParam+"%")
	var conversationList []*model_struct.LocalConversation
	return conversationList, errs.WrapMsg(d.conn.WithContext(ctx).Where(condition).Order("latest_msg_send_time DESC").Find(&conversationList).Error, "SearchConversation failed ")
}
