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
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"

	"gorm.io/gorm"
)

func (d *DataBase) initSuperLocalChatLog(ctx context.Context, groupID string) {
	if !d.conn.WithContext(ctx).Migrator().HasTable(utils.GetConversationTableName(groupID)) {
		d.conn.WithContext(ctx).Table(utils.GetConversationTableName(groupID)).AutoMigrate(&model_struct.LocalChatLog{})
	}
}
func (d *DataBase) SuperGroupBatchInsertMessageList(ctx context.Context, MessageList []*model_struct.LocalChatLog, groupID string) error {
	if MessageList == nil {
		return nil
	}
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetConversationTableName(groupID)).Create(MessageList).Error, "SuperGroupBatchInsertMessageList failed")
}
func (d *DataBase) SuperGroupInsertMessage(ctx context.Context, Message *model_struct.LocalChatLog, groupID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetConversationTableName(groupID)).Create(Message).Error, "SuperGroupInsertMessage failed")
}
func (d *DataBase) SuperGroupDeleteAllMessage(ctx context.Context, groupID string) error {
	return utils.Wrap(d.conn.WithContext(ctx).Session(&gorm.Session{AllowGlobalUpdate: true}).Table(utils.GetConversationTableName(groupID)).Delete(&model_struct.LocalChatLog{}).Error, "SuperGroupDeleteAllMessage failed")
}
func (d *DataBase) SuperGroupSearchMessageByKeyword(ctx context.Context, contentType []int, keywordList []string, keywordListMatchType int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var messageList []model_struct.LocalChatLog
	var condition string
	var subCondition string
	if keywordListMatchType == constant.KeywordMatchOr {
		for i := 0; i < len(keywordList); i++ {
			if i == 0 {
				subCondition += "And ("
			}
			if i+1 >= len(keywordList) {
				subCondition += "content like " + "'%" + keywordList[i] + "%') "
			} else {
				subCondition += "content like " + "'%" + keywordList[i] + "%' " + "or "

			}
		}
	} else {
		for i := 0; i < len(keywordList); i++ {
			if i == 0 {
				subCondition += "And ("
			}
			if i+1 >= len(keywordList) {
				subCondition += "content like " + "'%" + keywordList[i] + "%') "
			} else {
				subCondition += "content like " + "'%" + keywordList[i] + "%' " + "and "
			}
		}
	}
	condition = fmt.Sprintf("recv_id=%q And send_time between %d and %d AND status <=%d  And content_type IN ? ", sourceID, startTime, endTime, constant.MsgStatusSendFailed)

	condition += subCondition
	err = utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetConversationTableName(sourceID)).Where(condition, contentType).Order("send_time DESC").Offset(offset).Limit(count).Find(&messageList).Error, "InsertMessage failed")

	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}

func (d *DataBase) SuperGroupSearchAllMessageByContentType(ctx context.Context, groupID string, contentType int32) (result []*model_struct.LocalChatLog, err error) {
	err = d.conn.WithContext(ctx).Table(utils.GetConversationTableName(groupID)).Where("content_type = ?", contentType).Find(&result).Error
	return result, err
}

func (d *DataBase) SuperGroupSearchMessageByContentType(ctx context.Context, contentType []int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	var messageList []model_struct.LocalChatLog
	var condition string
	condition = fmt.Sprintf("session_type=%d And recv_id==%q And send_time between %d and %d AND status <=%d And content_type IN ?", sessionType, sourceID, startTime, endTime, constant.MsgStatusSendFailed)

	err = utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetConversationTableName(sourceID)).Where(condition, contentType).Order("send_time DESC").Offset(offset).Limit(count).Find(&messageList).Error, "SearchMessage failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}

func (d *DataBase) SuperGroupSearchMessageByContentTypeAndKeyword(ctx context.Context, contentType []int, keywordList []string, keywordListMatchType int, startTime, endTime int64, groupID string) (result []*model_struct.LocalChatLog, err error) {
	var messageList []model_struct.LocalChatLog
	var condition string
	var subCondition string
	if keywordListMatchType == constant.KeywordMatchOr {
		for i := 0; i < len(keywordList); i++ {
			if i == 0 {
				subCondition += "And ("
			}
			if i+1 >= len(keywordList) {
				subCondition += "content like " + "'%" + keywordList[i] + "%') "
			} else {
				subCondition += "content like " + "'%" + keywordList[i] + "%' " + "or "

			}
		}
	} else {
		for i := 0; i < len(keywordList); i++ {
			if i == 0 {
				subCondition += "And ("
			}
			if i+1 >= len(keywordList) {
				subCondition += "content like " + "'%" + keywordList[i] + "%') "
			} else {
				subCondition += "content like " + "'%" + keywordList[i] + "%' " + "and "
			}
		}
	}
	condition = fmt.Sprintf("send_time between %d and %d AND status <=%d  And content_type IN ? ", startTime, endTime, constant.MsgStatusSendFailed)
	condition += subCondition
	log.Info("key owrd", condition)
	err = utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetConversationTableName(groupID)).Where(condition, contentType).Order("send_time DESC").Find(&messageList).Error, "SearchMessage failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}
func (d *DataBase) SuperGroupBatchUpdateMessageList(ctx context.Context, MessageList []*model_struct.LocalChatLog) error {
	if MessageList == nil {
		return nil
	}

	for _, v := range MessageList {
		v1 := new(model_struct.LocalChatLog)
		v1.ClientMsgID = v.ClientMsgID
		v1.Seq = v.Seq
		v1.Status = v.Status
		err := d.SuperGroupUpdateMessage(ctx, v1)
		if err != nil {
			return utils.Wrap(err, "BatchUpdateMessageList failed")
		}

	}
	return nil
}

func (d *DataBase) SuperGroupMessageIfExists(ctx context.Context, ClientMsgID string) (bool, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var count int64
	t := d.conn.WithContext(ctx).Model(&model_struct.LocalChatLog{}).Where("client_msg_id = ?",
		ClientMsgID).Count(&count)
	if t.Error != nil {
		return false, utils.Wrap(t.Error, "MessageIfExists get failed")
	}
	if count != 1 {
		return false, nil
	} else {
		return true, nil
	}
}
func (d *DataBase) SuperGroupIsExistsInErrChatLogBySeq(ctx context.Context, seq int64) bool {
	return true
}
func (d *DataBase) SuperGroupMessageIfExistsBySeq(ctx context.Context, seq int64) (bool, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var count int64
	t := d.conn.WithContext(ctx).Model(&model_struct.LocalChatLog{}).Where("seq = ?",
		seq).Count(&count)
	if t.Error != nil {
		return false, utils.Wrap(t.Error, "MessageIfExistsBySeq get failed")
	}
	if count != 1 {
		return false, nil
	} else {
		return true, nil
	}
}
func (d *DataBase) SuperGroupGetMessage(ctx context.Context, msg *sdk_struct.MsgStruct) (*model_struct.LocalChatLog, error) {
	d.initSuperLocalChatLog(ctx, msg.GroupID)
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var c model_struct.LocalChatLog
	return &c, utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetConversationTableName(msg.GroupID)).Where("client_msg_id = ?",
		msg.ClientMsgID).Take(&c).Error, "GetMessage failed")
}

func (d *DataBase) SuperGroupGetAllUnDeleteMessageSeqList(ctx context.Context) ([]uint32, error) {
	var seqList []uint32
	return seqList, utils.Wrap(d.conn.WithContext(ctx).Model(&model_struct.LocalChatLog{}).Where("status != 4").Select("seq").Find(&seqList).Error, "")
}

func (d *DataBase) SuperGroupUpdateColumnsMessage(ctx context.Context, ClientMsgID, groupID string, args map[string]interface{}) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Table(utils.GetConversationTableName(groupID)).Where(
		"client_msg_id = ? ", ClientMsgID).Updates(args)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateColumnsConversation failed")
}
func (d *DataBase) SuperGroupUpdateMessage(ctx context.Context, c *model_struct.LocalChatLog) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Table(utils.GetConversationTableName(c.RecvID)).Updates(c)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateMessage failed")
}
func (d *DataBase) SuperGroupUpdateSpecificContentTypeMessage(ctx context.Context, contentType int, groupID string, args map[string]interface{}) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Table(utils.GetConversationTableName(groupID)).Where("content_type = ?", contentType).Updates(args)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateMessage failed")
}

//funcation (d *DataBase) SuperGroupDeleteAllMessage(ctx context.Context, ) error {
//	d.mRWMutex.Lock()
//	defer d.mRWMutex.Unlock()
//	err := d.conn.WithContext(ctx).Model(&model_struct.LocalChatLog{}).Exec("update local_chat_logs set status = ?,content = ? ", constant.MsgStatusHasDeleted, "").Error
//	return utils.Wrap(err, "delete all message error")
//}

func (d *DataBase) SuperGroupUpdateMessageStatusBySourceID(ctx context.Context, sourceID string, status, sessionType int32) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var condition string
	if sourceID == d.loginUserID && sessionType == constant.SingleChatType {
		condition = "send_id=? And recv_id=? AND session_type=?"
	} else {
		condition = "(send_id=? or recv_id=?)AND session_type=?"
	}
	t := d.conn.WithContext(ctx).Table(utils.GetConversationTableName(sourceID)).Where(condition, sourceID, sourceID, sessionType).Updates(model_struct.LocalChatLog{Status: status})
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateMessageStatusBySourceID failed")
}
func (d *DataBase) SuperGroupUpdateMessageTimeAndStatus(ctx context.Context, msg *sdk_struct.MsgStruct) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Table(utils.GetConversationTableName(msg.GroupID)).Where("client_msg_id=? And seq=?", msg.ClientMsgID, 0).Updates(model_struct.LocalChatLog{Status: msg.Status, SendTime: msg.SendTime, ServerMsgID: msg.ServerMsgID})
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "SuperGroupUpdateMessageTimeAndStatus failed")
}

func (d *DataBase) SuperGroupGetMessageList(ctx context.Context, sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var messageList []model_struct.LocalChatLog
	var condition, timeOrder, timeSymbol string
	if isReverse {
		timeOrder = "send_time ASC"
		timeSymbol = ">"
	} else {
		timeOrder = "send_time DESC"
		timeSymbol = "<"
	}
	condition = " recv_id = ? AND status <=? And session_type = ? And send_time " + timeSymbol + " ?"

	err = utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetConversationTableName(sourceID)).Where(condition, sourceID, constant.MsgStatusSendFailed, sessionType, startTime).
		Order(timeOrder).Offset(0).Limit(count).Find(&messageList).Error, "GetMessageList failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}
func (d *DataBase) SuperGroupGetMessageListNoTime(ctx context.Context, sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	d.initSuperLocalChatLog(ctx, sourceID)
	var messageList []model_struct.LocalChatLog
	var condition, timeOrder string
	if isReverse {
		timeOrder = "send_time ASC"
	} else {
		timeOrder = "send_time DESC"
	}
	condition = "recv_id = ? AND status <=? And session_type = ? "

	err = utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetConversationTableName(sourceID)).Where(condition, sourceID, constant.MsgStatusSendFailed, sessionType).
		Order(timeOrder).Offset(0).Limit(count).Find(&messageList).Error, "GetMessageList failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}

func (d *DataBase) SuperGroupGetSendingMessageList(ctx context.Context, groupID string) (result []*model_struct.LocalChatLog, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var messageList []model_struct.LocalChatLog
	err = utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetConversationTableName(groupID)).Where("status = ?", constant.MsgStatusSending).Find(&messageList).Error, "GetMessageList failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}

func (d *DataBase) SuperGroupUpdateGroupMessageHasRead(ctx context.Context, msgIDList []string, groupID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Table(utils.GetConversationTableName(groupID)).Where(" client_msg_id in ?", msgIDList).Update("is_read", constant.HasRead)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateMessageStatusBySourceID failed")
}
func (d *DataBase) SuperGroupUpdateGroupMessageFields(ctx context.Context, msgIDList []string, groupID string, args map[string]interface{}) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Table(utils.GetConversationTableName(groupID)).Where(" client_msg_id in ?", msgIDList).Updates(args)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateMessageStatusBySourceID failed")
}

func (d *DataBase) SuperGroupGetNormalMsgSeq(ctx context.Context) (int64, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seq int64
	err := d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Select("IFNULL(max(seq),0)").Find(&seq).Error
	return seq, utils.Wrap(err, "GetNormalMsgSeq")
}
func (d *DataBase) SuperGroupGetNormalMinSeq(ctx context.Context, groupID string) (int64, error) {
	var seq int64
	err := d.conn.WithContext(ctx).Table(utils.GetConversationTableName(groupID)).Select("IFNULL(min(seq),0)").Where("seq >?", 0).Find(&seq).Error
	return seq, utils.Wrap(err, "SuperGroupGetNormalMinSeq")
}
func (d *DataBase) SuperGroupGetTestMessage(ctx context.Context, seq int64) (*model_struct.LocalChatLog, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var c model_struct.LocalChatLog
	return &c, utils.Wrap(d.conn.WithContext(ctx).Where("seq = ?",
		seq).Find(&c).Error, "GetTestMessage failed")
}

func (d *DataBase) SuperGroupUpdateMsgSenderNickname(ctx context.Context, sendID, nickname string, sType int) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Where(
		"send_id = ? and session_type = ? and sender_nick_name != ? ", sendID, sType, nickname).Updates(
		map[string]interface{}{"sender_nick_name": nickname}).Error, utils.GetSelfFuncName()+" failed")
}

func (d *DataBase) SuperGroupUpdateMsgSenderFaceURL(ctx context.Context, sendID, faceURL string, sType int) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Where(
		"send_id = ? and session_type = ? and sender_face_url != ? ", sendID, sType, faceURL).Updates(
		map[string]interface{}{"sender_face_url": faceURL}).Error, utils.GetSelfFuncName()+" failed")
}
func (d *DataBase) SuperGroupUpdateMsgSenderFaceURLAndSenderNickname(ctx context.Context, sendID, faceURL, nickname string, sessionType int, groupID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetConversationTableName(groupID)).Where(
		"send_id = ? and session_type = ?", sendID, sessionType).Updates(
		map[string]interface{}{"sender_face_url": faceURL, "sender_nick_name": nickname}).Error, utils.GetSelfFuncName()+" failed")
}

func (d *DataBase) SuperGroupGetMsgSeqByClientMsgID(ctx context.Context, clientMsgID string, groupID string) (uint32, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seq uint32
	err := utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetConversationTableName(groupID)).Select("seq").Where("client_msg_id=?", clientMsgID).Take(&seq).Error, utils.GetSelfFuncName()+" failed")
	return seq, err
}

func (d *DataBase) SuperGroupGetMsgSeqListByGroupID(ctx context.Context, groupID string) ([]uint32, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seqList []uint32
	err := utils.Wrap(d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Select("seq").Where("recv_id=?", groupID).Find(&seqList).Error, utils.GetSelfFuncName()+" failed")
	return seqList, err
}

func (d *DataBase) SuperGroupGetMsgSeqListByPeerUserID(ctx context.Context, userID string) ([]uint32, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seqList []uint32
	err := utils.Wrap(d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Select("seq").Where("recv_id=? or send_id=?", userID, userID).Find(&seqList).Error, utils.GetSelfFuncName()+" failed")
	return seqList, err
}

func (d *DataBase) SuperGroupGetMsgSeqListBySelfUserID(ctx context.Context, userID string) ([]uint32, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seqList []uint32
	err := utils.Wrap(d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Select("seq").Where("recv_id=? and send_id=?", userID, userID).Find(&seqList).Error, utils.GetSelfFuncName()+" failed")
	return seqList, err
}
