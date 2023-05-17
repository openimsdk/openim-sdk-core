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

//	func (d *DataBase) BatchInsertMessageList(ctx context.Context, MessageList []*model_struct.LocalChatLog) error {
//		if MessageList == nil {
//			return nil
//		}
//		d.mRWMutex.Lock()
//		defer d.mRWMutex.Unlock()
//		return utils.Wrap(d.conn.WithContext(ctx).Create(MessageList).Error, "BatchInsertMessageList failed")
//	}
//func (d *DataBase) BatchInsertMessageListController(ctx context.Context, MessageList []*model_struct.LocalChatLog) error {
//	if len(MessageList) == 0 {
//		return nil
//	}
//	switch MessageList[len(MessageList)-1].SessionType {
//	case constant.SuperGroupChatType:
//		return d.SuperGroupBatchInsertMessageList(ctx, MessageList, MessageList[len(MessageList)-1].RecvID)
//	default:
//		return d.BatchInsertMessageList(ctx, MessageList)
//	}
//}

//	func (d *DataBase) InsertMessage(ctx context.Context, Message *model_struct.LocalChatLog) error {
//		d.mRWMutex.Lock()
//		defer d.mRWMutex.Unlock()
//		return utils.Wrap(d.conn.WithContext(ctx).Create(Message).Error, "InsertMessage failed")
//	}
//
//	func (d *DataBase) InsertMessageController(ctx context.Context, message *model_struct.LocalChatLog) error {
//		switch message.SessionType {
//		case constant.SuperGroupChatType:
//			return d.SuperGroupInsertMessage(ctx, message, message.RecvID)
//		default:
//			return d.InsertMessage(ctx, message)
//		}
//	}
func (d *DataBase) SearchMessageByKeyword(ctx context.Context, contentType []int, keywordList []string, keywordListMatchType int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
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
	switch sessionType {
	case constant.SingleChatType, constant.NotificationChatType:
		condition = fmt.Sprintf("session_type==%d And (send_id==%q OR recv_id==%q) And send_time  between %d and %d AND status <=%d  And content_type IN ? ", constant.SingleChatType, sourceID, sourceID, startTime, endTime, constant.MsgStatusSendFailed)
	case constant.GroupChatType:
		condition = fmt.Sprintf("session_type==%d And recv_id==%q And send_time between %d and %d AND status <=%d  And content_type IN ? ", constant.GroupChatType, sourceID, startTime, endTime, constant.MsgStatusSendFailed)
	default:
		//condition = fmt.Sprintf("(send_id==%q OR recv_id==%q) And send_time between %d and %d AND status <=%d And content_type == %d And content like %q", sourceID, sourceID, startTime, endTime, constant.MsgStatusSendFailed, constant.Text, "%"+keyword+"%")
		return nil, err
	}
	condition += subCondition
	err = utils.Wrap(d.conn.WithContext(ctx).Where(condition, contentType).Order("send_time DESC").Offset(offset).Limit(count).Find(&messageList).Error, "InsertMessage failed")

	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}
func (d *DataBase) SearchMessageByKeywordController(ctx context.Context, contentType []int, keywordList []string, keywordListMatchType int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	switch sessionType {
	case constant.SuperGroupChatType:
		return d.SuperGroupSearchMessageByKeyword(ctx, contentType, keywordList, keywordListMatchType, sourceID, startTime, endTime, sessionType, offset, count)
	default:
		return d.SearchMessageByKeyword(ctx, contentType, keywordList, keywordListMatchType, sourceID, startTime, endTime, sessionType, offset, count)
	}
}

func (d *DataBase) SearchMessageByContentType(ctx context.Context, contentType []int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	var messageList []model_struct.LocalChatLog
	var condition string
	switch sessionType {
	case constant.SingleChatType, constant.NotificationChatType:
		condition = fmt.Sprintf("session_type==%d And (send_id==%q OR recv_id==%q) And send_time between %d and %d AND status <=%d And content_type IN ?", constant.SingleChatType, sourceID, sourceID, startTime, endTime, constant.MsgStatusSendFailed)
	case constant.GroupChatType:
		condition = fmt.Sprintf("session_type==%d And recv_id==%q And send_time between %d and %d AND status <=%d And content_type IN ?", constant.GroupChatType, sourceID, startTime, endTime, constant.MsgStatusSendFailed)
	default:
		return nil, err
	}
	err = utils.Wrap(d.conn.WithContext(ctx).Where(condition, contentType).Order("send_time DESC").Offset(offset).Limit(count).Find(&messageList).Error, "SearchMessage failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}

func (d *DataBase) SearchAllMessageByContentType(ctx context.Context, contentType int) (result []*model_struct.LocalChatLog, err error) {
	err = d.conn.WithContext(ctx).Model(&model_struct.LocalChatLog{}).Where("content_type = ?", contentType).Find(&result).Error
	return result, err
}

func (d *DataBase) SearchMessageByContentTypeController(ctx context.Context, contentType []int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	switch sessionType {
	case constant.SuperGroupChatType:
		return d.SuperGroupSearchMessageByContentType(ctx, contentType, sourceID, startTime, endTime, sessionType, offset, count)
	default:
		return d.SearchMessageByContentType(ctx, contentType, sourceID, startTime, endTime, sessionType, offset, count)
	}
}

func (d *DataBase) SearchMessageByContentTypeAndKeyword(ctx context.Context, contentType []int, keywordList []string, keywordListMatchType int, startTime, endTime int64) (result []*model_struct.LocalChatLog, err error) {
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
	err = utils.Wrap(d.conn.WithContext(ctx).Where(condition, contentType).Order("send_time DESC").Find(&messageList).Error, "SearchMessage failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}
func (d *DataBase) SearchMessageByContentTypeAndKeywordController(ctx context.Context, contentType []int, keywordList []string, keywordListMatchType int, startTime, endTime int64) (result []*model_struct.LocalChatLog, err error) {
	list, err := d.SearchMessageByContentTypeAndKeyword(ctx, contentType, keywordList, keywordListMatchType, startTime, endTime)
	if err != nil {
		return nil, err
	}
	superGroupIDList, err := d.GetJoinedSuperGroupIDList(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range superGroupIDList {
		sList, err := d.SuperGroupSearchMessageByContentTypeAndKeyword(ctx, contentType, keywordList, keywordListMatchType, startTime, endTime, v)
		if err != nil {
			// TODO: log.Error(operationID, "search message in group err", err.Error(), v)
			continue
		}
		list = append(list, sList...)
	}
	workingGroupIDList, err := d.GetJoinedWorkingGroupIDList(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range workingGroupIDList {
		sList, err := d.SuperGroupSearchMessageByContentTypeAndKeyword(ctx, contentType, keywordList, keywordListMatchType, startTime, endTime, v)
		if err != nil {
			// TODO: log.Error(operationID, "search message in group err", err.Error(), v)
			continue
		}
		list = append(list, sList...)
	}
	return list, nil
}

//func (d *DataBase) BatchUpdateMessageList(ctx context.Context, MessageList []*model_struct.LocalChatLog) error {
//	if MessageList == nil {
//		return nil
//	}
//
//	for _, v := range MessageList {
//		v1 := new(model_struct.LocalChatLog)
//		v1.ClientMsgID = v.ClientMsgID
//		v1.Seq = v.Seq
//		v1.Status = v.Status
//		v1.RecvID = v.RecvID
//		v1.SessionType = v.SessionType
//		v1.ServerMsgID = v.ServerMsgID
//		err := d.UpdateMessageController(ctx, v1)
//		if err != nil {
//			return utils.Wrap(err, "BatchUpdateMessageList failed")
//		}
//
//	}
//	return nil
//}

//func (d *DataBase) BatchSpecialUpdateMessageList(ctx context.Context, MessageList []*model_struct.LocalChatLog) error {
//	if MessageList == nil {
//		return nil
//	}
//
//	for _, v := range MessageList {
//		v1 := new(model_struct.LocalChatLog)
//		v1.ClientMsgID = v.ClientMsgID
//		v1.ServerMsgID = v.ServerMsgID
//		v1.SendID = v.SendID
//		v1.RecvID = v.RecvID
//		v1.SenderPlatformID = v.SenderPlatformID
//		v1.SenderNickname = v.SenderNickname
//		v1.SenderFaceURL = v.SenderFaceURL
//		v1.SessionType = v.SessionType
//		v1.MsgFrom = v.MsgFrom
//		v1.ContentType = v.ContentType
//		v1.Content = v.Content
//		v1.Seq = v.Seq
//		v1.SendTime = v.SendTime
//		v1.CreateTime = v.CreateTime
//		v1.AttachedInfo = v.AttachedInfo
//		v1.Ex = v.Ex
//		err := d.UpdateMessageController(ctx, v1)
//		if err != nil {
//			log.Error("", "update single message failed", *v)
//			return utils.Wrap(err, "BatchUpdateMessageList failed")
//		}
//
//	}
//	return nil
//}

func (d *DataBase) MessageIfExists(ctx context.Context, ClientMsgID string) (bool, error) {
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
func (d *DataBase) IsExistsInErrChatLogBySeq(ctx context.Context, seq int64) bool {
	return true
}
func (d *DataBase) MessageIfExistsBySeq(ctx context.Context, seq int64) (bool, error) {
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

//	func (d *DataBase) GetMessage(ctx context.Context, ClientMsgID string) (*model_struct.LocalChatLog, error) {
//		var c model_struct.LocalChatLog
//		return &c, utils.Wrap(d.conn.WithContext(ctx).Where("client_msg_id = ?",
//			ClientMsgID).Take(&c).Error, "GetMessage failed")
//	}
//func (d *DataBase) GetMessageController(ctx context.Context, msg *sdk_struct.MsgStruct) (*model_struct.LocalChatLog, error) {
//	switch msg.SessionType {
//	case constant.SuperGroupChatType:
//		return d.SuperGroupGetMessage(ctx, msg)
//	default:
//		return d.GetMessage(ctx, msg.ClientMsgID)
//	}
//}

func (d *DataBase) GetAllUnDeleteMessageSeqList(ctx context.Context) ([]uint32, error) {
	var seqList []uint32
	return seqList, utils.Wrap(d.conn.WithContext(ctx).Model(&model_struct.LocalChatLog{}).Where("status != ?", constant.MsgStatusHasDeleted).Select("seq").Find(&seqList).Error, "")
}
func (d *DataBase) UpdateColumnsMessageList(ctx context.Context, clientMsgIDList []string, args map[string]interface{}) error {
	c := model_struct.LocalChatLog{}
	t := d.conn.WithContext(ctx).Model(&c).Where("client_msg_id IN", clientMsgIDList).Updates(args)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateColumnsConversation failed")
}
func (d *DataBase) UpdateColumnsMessage(ctx context.Context, ClientMsgID string, args map[string]interface{}) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	c := model_struct.LocalChatLog{ClientMsgID: ClientMsgID}
	t := d.conn.WithContext(ctx).Model(&c).Updates(args)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateColumnsConversation failed")
}
func (d *DataBase) UpdateColumnsMessageController(ctx context.Context, ClientMsgID string, groupID string, sessionType int32, args map[string]interface{}) error {
	switch sessionType {
	case constant.SuperGroupChatType:
		return utils.Wrap(d.SuperGroupUpdateColumnsMessage(ctx, ClientMsgID, groupID, args), "")
	default:
		return utils.Wrap(d.UpdateColumnsMessage(ctx, ClientMsgID, args), "")
	}
}

//	func (d *DataBase) UpdateMessage(ctx context.Context, c *model_struct.LocalChatLog) error {
//		t := d.conn.WithContext(ctx).Updates(c)
//		if t.RowsAffected == 0 {
//			return utils.Wrap(errors.New("RowsAffected == 0"), "no update ")
//		}
//		return utils.Wrap(t.Error, "UpdateMessage failed")
//	}
//func (d *DataBase) UpdateMessageController(ctx context.Context, c *model_struct.LocalChatLog) error {
//	switch c.SessionType {
//	case constant.SuperGroupChatType:
//		return utils.Wrap(d.SuperGroupUpdateMessage(ctx, c), "")
//	default:
//		return utils.Wrap(d.UpdateMessage(ctx, c), "")
//	}
//}

func (d *DataBase) DeleteAllMessage(ctx context.Context) error {
	m := model_struct.LocalChatLog{Status: constant.MsgStatusHasDeleted, Content: ""}
	err := d.conn.WithContext(ctx).Session(&gorm.Session{AllowGlobalUpdate: true}).Select("status", "content").Updates(m).Error
	return utils.Wrap(err, "delete all message error")
}
func (d *DataBase) UpdateMessageStatusBySourceID(ctx context.Context, sourceID string, status, sessionType int32) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var condition string
	if sourceID == d.loginUserID && sessionType == constant.SingleChatType {
		condition = "send_id=? And recv_id=? AND session_type=?"
	} else {
		condition = "(send_id=? or recv_id=?)AND session_type=?"
	}
	t := d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Where(condition, sourceID, sourceID, sessionType).Updates(model_struct.LocalChatLog{Status: status})
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateMessageStatusBySourceID failed")
}
func (d *DataBase) UpdateMessageStatusBySourceIDController(ctx context.Context, sourceID string, status, sessionType int32) error {
	switch sessionType {
	case constant.SuperGroupChatType:
		return d.SuperGroupUpdateMessageStatusBySourceID(ctx, sourceID, status, sessionType)
	default:
		return d.UpdateMessageStatusBySourceID(ctx, sourceID, status, sessionType)
	}
}

//	func (d *DataBase) UpdateMessageTimeAndStatus(ctx context.Context, clientMsgID string, serverMsgID string, sendTime int64, status int32) error {
//		d.mRWMutex.Lock()
//		defer d.mRWMutex.Unlock()
//		return utils.Wrap(d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Where("client_msg_id=? And seq=?", clientMsgID, 0).
//			Updates(model_struct.LocalChatLog{Status: status, SendTime: sendTime, ServerMsgID: serverMsgID}).Error, "UpdateMessageStatusBySourceID failed")
//
// }
//func (d *DataBase) UpdateMessageTimeAndStatusController(ctx context.Context, msg *sdk_struct.MsgStruct) error {
//	switch msg.SessionType {
//	case constant.SuperGroupChatType:
//		return d.SuperGroupUpdateMessageTimeAndStatus(ctx, msg)
//	default:
//		return d.UpdateMessageTimeAndStatus(ctx, msg.ClientMsgID, msg.ServerMsgID, msg.SendTime, msg.Status)
//	}
//}

//func (d *DataBase) UpdateMessageAttachedInfo(ctx context.Context, msg *sdk_struct.MsgStruct) error {
//	info, err := json.Marshal(msg.AttachedInfoElem)
//	if err != nil {
//		return err
//	}
//	d.mRWMutex.Lock()
//	defer d.mRWMutex.Unlock()
//	switch msg.SessionType {
//	case constant.SuperGroupChatType:
//		t := d.conn.WithContext(ctx).Table(utils.GetSuperGroupTableName(msg.GroupID)).Where("client_msg_id=?", msg.ClientMsgID).Updates(map[string]any{"attached_info": string(info)})
//		if t.RowsAffected == 0 {
//			return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
//		}
//		return utils.Wrap(t.Error, "SuperGroupUpdateMessageTimeAndStatus failed")
//	default:
//		return utils.Wrap(d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Where("client_msg_id=?", msg.ClientMsgID).Updates(map[string]any{"attached_info": string(info)}).Error, "")
//	}
//}

func (d *DataBase) UpdateMessageByClientMsgID(ctx context.Context, clientMsgID string, data map[string]any) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Where("client_msg_id=?", clientMsgID).Updates(data).Error, "")
}

// group ,index_recv_id and index_send_time only one can be used,when index_recv_id be used,temp B tree use for order by,Query speed decrease
//func (d *DataBase) GetMessageList(ctx context.Context, sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
//	d.mRWMutex.Lock()
//	defer d.mRWMutex.Unlock()
//	var messageList []model_struct.LocalChatLog
//	var condition, timeOrder, timeSymbol string
//	if isReverse {
//		timeOrder = "send_time ASC"
//		timeSymbol = ">"
//	} else {
//		timeOrder = "send_time DESC"
//		timeSymbol = "<"
//	}
//	if sessionType == constant.SingleChatType && sourceID == d.loginUserID {
//		condition = "send_id = ? And recv_id = ? AND status <=? And session_type = ? And send_time " + timeSymbol + " ?"
//	} else {
//		condition = "(send_id = ? OR recv_id = ?) AND status <=? And session_type = ? And send_time " + timeSymbol + " ?"
//	}
//	err = utils.Wrap(d.conn.WithContext(ctx).Where(condition, sourceID, sourceID, constant.MsgStatusSendFailed, sessionType, startTime).
//		Order(timeOrder).Offset(0).Limit(count).Find(&messageList).Error, "GetMessageList failed")
//	for _, v := range messageList {
//		v1 := v
//		result = append(result, &v1)
//	}
//	return result, err
//}

func (d *DataBase) GetAllMessageForTest(ctx context.Context) (result []*model_struct.LocalChatLog, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var messageList []model_struct.LocalChatLog

	err = utils.Wrap(d.conn.WithContext(ctx).Find(&messageList).Error, "GetMessageList failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}

//func (d *DataBase) GetMessageListController(ctx context.Context, sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
//	switch sessionType {
//	case constant.SuperGroupChatType:
//		return d.SuperGroupGetMessageList(ctx, sourceID, sessionType, count, startTime, isReverse)
//	default:
//		return d.GetMessageList(ctx, sourceID, sessionType, count, startTime, isReverse)
//	}
//}

//func (d *DataBase) GetMessageListNoTime(ctx context.Context, sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
//	d.mRWMutex.Lock()
//	defer d.mRWMutex.Unlock()
//	var messageList []model_struct.LocalChatLog
//	var condition, timeOrder string
//	if isReverse {
//		timeOrder = "send_time ASC"
//	} else {
//		timeOrder = "send_time DESC"
//	}
//	switch sessionType {
//	case constant.SingleChatType:
//		if sourceID == d.loginUserID {
//			condition = "send_id = ? And recv_id = ? AND status <=? And session_type = ?"
//		} else {
//			condition = "(send_id = ? OR recv_id = ?) AND status <=? And session_type = ? "
//		}
//		err = utils.Wrap(d.conn.WithContext(ctx).Where(condition, sourceID, sourceID, constant.MsgStatusSendFailed, sessionType).
//			Order(timeOrder).Offset(0).Limit(count).Find(&messageList).Error, "GetMessageList failed")
//	case constant.GroupChatType:
//		condition = " recv_id = ? AND status <=? And session_type = ? "
//		err = utils.Wrap(d.conn.WithContext(ctx).Where(condition, sourceID, constant.MsgStatusSendFailed, sessionType).
//			Order(timeOrder).Offset(0).Limit(count).Find(&messageList).Error, "GetMessageList failed")
//	default:
//		condition = "(send_id = ? OR recv_id = ?) AND status <=? And session_type = ? "
//		err = utils.Wrap(d.conn.WithContext(ctx).Where(condition, sourceID, sourceID, constant.MsgStatusSendFailed, sessionType).
//			Order(timeOrder).Offset(0).Limit(count).Find(&messageList).Error, "GetMessageList failed")
//	}
//
//	for _, v := range messageList {
//		v1 := v
//		result = append(result, &v1)
//	}
//	return result, err
//}

//func (d *DataBase) GetMessageListNoTimeController(ctx context.Context, sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
//	switch sessionType {
//	case constant.SuperGroupChatType:
//		return d.SuperGroupGetMessageListNoTime(ctx, sourceID, sessionType, count, isReverse)
//	default:
//		return d.GetMessageListNoTime(ctx, sourceID, sessionType, count, isReverse)
//	}
//}

func (d *DataBase) GetSendingMessageList(ctx context.Context) (result []*model_struct.LocalChatLog, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var messageList []model_struct.LocalChatLog
	err = utils.Wrap(d.conn.WithContext(ctx).Where("status = ?", constant.MsgStatusSending).Find(&messageList).Error, "GetMessageList failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}
func (d *DataBase) UpdateSingleMessageHasRead(ctx context.Context, sendID string, msgIDList []string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Where("send_id=?  AND session_type=? AND client_msg_id in ?", sendID, constant.SingleChatType, msgIDList).Update("is_read", constant.HasRead)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateMessageStatusBySourceID failed")
}
func (d *DataBase) UpdateGroupMessageHasRead(ctx context.Context, msgIDList []string, sessionType int32) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Where("session_type=? AND client_msg_id in ?", sessionType, msgIDList).Update("is_read", constant.HasRead)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateMessageStatusBySourceID failed")
}
func (d *DataBase) UpdateGroupMessageHasReadController(ctx context.Context, msgIDList []string, groupID string, sessionType int32) error {
	switch sessionType {
	case constant.SuperGroupChatType:
		return d.SuperGroupUpdateGroupMessageHasRead(ctx, msgIDList, groupID)
	default:
		return d.UpdateGroupMessageHasRead(ctx, msgIDList, sessionType)
	}
}
func (d *DataBase) GetMultipleMessage(ctx context.Context, msgIDList []string) (result []*model_struct.LocalChatLog, err error) {
	var messageList []model_struct.LocalChatLog
	err = utils.Wrap(d.conn.WithContext(ctx).Where("client_msg_id IN ?", msgIDList).Order("send_time DESC").Find(&messageList).Error, "GetMultipleMessage failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}
func (d *DataBase) GetMultipleMessageController(ctx context.Context, msgIDList []string, groupID string, sessionType int32) (result []*model_struct.LocalChatLog, err error) {
	switch sessionType {
	case constant.SuperGroupChatType:
		return d.SuperGroupGetMultipleMessage(ctx, msgIDList, groupID)
	default:
		return d.GetMultipleMessage(ctx, msgIDList)
	}
}

func (d *DataBase) GetNormalMsgSeq(ctx context.Context) (int64, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seq int64
	err := d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Select("IFNULL(max(seq),0)").Find(&seq).Error
	return seq, utils.Wrap(err, "GetNormalMsgSeq")
}
func (d *DataBase) GetLostMsgSeqList(ctx context.Context, minSeqInSvr int64) ([]int64, error) {
	var hasSeqList []int64
	err := d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Select("seq").Order("seq ASC").Find(&hasSeqList).Error
	if err != nil {
		return nil, err
	}
	seqLength := len(hasSeqList)
	if seqLength == 0 {
		return nil, nil
	}
	var normalSeqList []int64

	var i int64
	for i = 1; i <= hasSeqList[seqLength-1]; i++ {
		normalSeqList = append(normalSeqList, i)
	}
	lostSeqList := utils.DifferenceSubset(normalSeqList, hasSeqList)
	if len(lostSeqList) == 0 {
		return nil, nil
	}
	abnormalSeqList, err := d.GetAbnormalMsgSeqList(ctx)
	if err != nil {
		return nil, err
	}
	if len(abnormalSeqList) == 0 {
		return lostSeqList, nil
	}
	return utils.DifferenceSubset(lostSeqList, utils.Intersect(lostSeqList, abnormalSeqList)), nil

}
func (d *DataBase) GetConversationNormalMsgSeq(ctx context.Context, conversationID string) (int64, error) {
	var seq int64
	err := d.conn.WithContext(ctx).Table(utils.GetSuperGroupTableName(conversationID)).Select("IFNULL(max(seq),0)").Find(&seq).Error
	return seq, utils.Wrap(err, "GetConversationNormalMsgSeq")
}

func (d *DataBase) GetTestMessage(ctx context.Context, seq uint32) (*model_struct.LocalChatLog, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var c model_struct.LocalChatLog
	return &c, utils.Wrap(d.conn.WithContext(ctx).Where("seq = ?",
		seq).Find(&c).Error, "GetTestMessage failed")
}

func (d *DataBase) UpdateMsgSenderNickname(ctx context.Context, sendID, nickname string, sType int) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Where(
		"send_id = ? and session_type = ? and sender_nick_name != ? ", sendID, sType, nickname).Updates(
		map[string]interface{}{"sender_nick_name": nickname}).Error, utils.GetSelfFuncName()+" failed")
}

func (d *DataBase) UpdateMsgSenderFaceURL(ctx context.Context, sendID, faceURL string, sType int) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Where(
		"send_id = ? and session_type = ? and sender_face_url != ? ", sendID, sType, faceURL).Updates(
		map[string]interface{}{"sender_face_url": faceURL}).Error, utils.GetSelfFuncName()+" failed")
}
func (d *DataBase) UpdateMsgSenderFaceURLAndSenderNicknameController(ctx context.Context, sendID, faceURL, nickname string, sessionType int, groupID string) error {
	switch sessionType {
	case constant.SuperGroupChatType:
		return d.SuperGroupUpdateMsgSenderFaceURLAndSenderNickname(ctx, sendID, faceURL, nickname, sessionType, groupID)
	default:
		return d.UpdateMsgSenderFaceURLAndSenderNickname(ctx, sendID, faceURL, nickname, sessionType)
	}
}
func (d *DataBase) UpdateMsgSenderFaceURLAndSenderNickname(ctx context.Context, sendID, faceURL, nickname string, sessionType int) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Where(
		"send_id = ? and session_type = ?", sendID, sessionType).Updates(
		map[string]interface{}{"sender_face_url": faceURL, "sender_nick_name": nickname}).Error, utils.GetSelfFuncName()+" failed")
}

func (d *DataBase) GetMsgSeqByClientMsgID(ctx context.Context, clientMsgID string) (uint32, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seq uint32
	err := utils.Wrap(d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Select("seq").Where("client_msg_id=?", clientMsgID).First(&seq).Error, utils.GetSelfFuncName()+" failed")
	return seq, err
}

func (d *DataBase) GetMsgSeqByClientMsgIDController(ctx context.Context, m *sdk_struct.MsgStruct) (uint32, error) {
	switch m.SessionType {
	case constant.SuperGroupChatType:
		return d.SuperGroupGetMsgSeqByClientMsgID(ctx, m.ClientMsgID, m.GroupID)
	default:
		return d.GetMsgSeqByClientMsgID(ctx, m.ClientMsgID)
	}
}

func (d *DataBase) GetMsgSeqListByGroupID(ctx context.Context, groupID string) ([]uint32, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seqList []uint32
	err := utils.Wrap(d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Select("seq").Where("recv_id=?", groupID).Find(&seqList).Error, utils.GetSelfFuncName()+" failed")
	return seqList, err
}

func (d *DataBase) GetMsgSeqListByPeerUserID(ctx context.Context, userID string) ([]uint32, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seqList []uint32
	err := utils.Wrap(d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Select("seq").Where("recv_id=? or send_id=?", userID, userID).Find(&seqList).Error, utils.GetSelfFuncName()+" failed")
	return seqList, err
}

func (d *DataBase) GetMsgSeqListBySelfUserID(ctx context.Context, userID string) ([]uint32, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seqList []uint32
	err := utils.Wrap(d.conn.WithContext(ctx).Model(model_struct.LocalChatLog{}).Select("seq").Where("recv_id=? and send_id=?", userID, userID).Find(&seqList).Error, utils.GetSelfFuncName()+" failed")
	return seqList, err
}
