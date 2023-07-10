package db

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

func (d *DataBase) initSuperLocalChatLog(groupID string) {
	if !d.conn.Migrator().HasTable(utils.GetSuperGroupTableName(groupID)) {
		d.conn.Table(utils.GetSuperGroupTableName(groupID)).AutoMigrate(&model_struct.LocalChatLog{})
	}
}
func (d *DataBase) SuperGroupBatchInsertMessageList(MessageList []*model_struct.LocalChatLog, groupID string) error {
	if MessageList == nil {
		return nil
	}
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Table(utils.GetSuperGroupTableName(groupID)).Create(MessageList).Error, "SuperGroupBatchInsertMessageList failed")
}
func (d *DataBase) SuperGroupInsertMessage(Message *model_struct.LocalChatLog, groupID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Table(utils.GetSuperGroupTableName(groupID)).Create(Message).Error, "SuperGroupInsertMessage failed")
}
func (d *DataBase) SuperGroupDeleteAllMessage(groupID string) error {
	return utils.Wrap(d.conn.Session(&gorm.Session{AllowGlobalUpdate: true}).Table(utils.GetSuperGroupTableName(groupID)).Delete(&model_struct.LocalChatLog{}).Error, "SuperGroupDeleteAllMessage failed")
}
func (d *DataBase) SuperGroupSearchMessageByKeyword(contentType []int, keywordList []string, keywordListMatchType int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
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
	err = utils.Wrap(d.conn.Table(utils.GetSuperGroupTableName(sourceID)).Where(condition, contentType).Order("send_time DESC").Offset(offset).Limit(count).Find(&messageList).Error, "InsertMessage failed")

	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}

func (d *DataBase) SuperGroupSearchAllMessageByContentType(groupID string, contentType int32) (result []*model_struct.LocalChatLog, err error) {
	err = d.conn.Table(utils.GetSuperGroupTableName(groupID)).Where("content_type = ?", contentType).Find(&result).Error
	return result, err
}

func (d *DataBase) SuperGroupSearchMessageByContentType(contentType []int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	var messageList []model_struct.LocalChatLog
	var condition string
	condition = fmt.Sprintf("session_type=%d And recv_id==%q And send_time between %d and %d AND status <=%d And content_type IN ?", sessionType, sourceID, startTime, endTime, constant.MsgStatusSendFailed)

	err = utils.Wrap(d.conn.Table(utils.GetSuperGroupTableName(sourceID)).Where(condition, contentType).Order("send_time DESC").Offset(offset).Limit(count).Find(&messageList).Error, "SearchMessage failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}

func (d *DataBase) SuperGroupSearchMessageByContentTypeAndKeyword(contentType []int, keywordList []string, keywordListMatchType int, startTime, endTime int64, groupID string) (result []*model_struct.LocalChatLog, err error) {
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
	err = utils.Wrap(d.conn.Table(utils.GetSuperGroupTableName(groupID)).Where(condition, contentType).Order("send_time DESC").Find(&messageList).Error, "SearchMessage failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}
func (d *DataBase) SuperGroupBatchUpdateMessageList(MessageList []*model_struct.LocalChatLog) error {
	if MessageList == nil {
		return nil
	}

	for _, v := range MessageList {
		v1 := new(model_struct.LocalChatLog)
		v1.ClientMsgID = v.ClientMsgID
		v1.Seq = v.Seq
		v1.Status = v.Status
		err := d.SuperGroupUpdateMessage(v1)
		if err != nil {
			return utils.Wrap(err, "BatchUpdateMessageList failed")
		}

	}
	return nil
}

func (d *DataBase) SuperGroupMessageIfExists(ClientMsgID string) (bool, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var count int64
	t := d.conn.Model(&model_struct.LocalChatLog{}).Where("client_msg_id = ?",
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
func (d *DataBase) SuperGroupIsExistsInErrChatLogBySeq(seq int64) bool {
	return true
}
func (d *DataBase) SuperGroupMessageIfExistsBySeq(seq int64) (bool, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var count int64
	t := d.conn.Model(&model_struct.LocalChatLog{}).Where("seq = ?",
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
func (d *DataBase) SuperGroupGetMessage(msg *sdk_struct.MsgStruct) (*model_struct.LocalChatLog, error) {
	d.initSuperLocalChatLog(msg.GroupID)
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var c model_struct.LocalChatLog
	return &c, utils.Wrap(d.conn.Table(utils.GetSuperGroupTableName(msg.GroupID)).Where("client_msg_id = ?",
		msg.ClientMsgID).Take(&c).Error, "GetMessage failed")
}

func (d *DataBase) SuperGroupGetAllUnDeleteMessageSeqList() ([]uint32, error) {
	var seqList []uint32
	return seqList, utils.Wrap(d.conn.Model(&model_struct.LocalChatLog{}).Where("status != 4").Select("seq").Find(&seqList).Error, "")
}

func (d *DataBase) SuperGroupUpdateColumnsMessage(ClientMsgID, groupID string, args map[string]interface{}) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Table(utils.GetSuperGroupTableName(groupID)).Where(
		"client_msg_id = ? ", ClientMsgID).Updates(args)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateColumnsConversation failed")
}
func (d *DataBase) SuperGroupUpdateMessage(c *model_struct.LocalChatLog) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Table(utils.GetSuperGroupTableName(c.RecvID)).Updates(c)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateMessage failed")
}
func (d *DataBase) SuperGroupUpdateSpecificContentTypeMessage(contentType int, groupID string, args map[string]interface{}) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Table(utils.GetSuperGroupTableName(groupID)).Where("content_type = ?", contentType).Updates(args)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateMessage failed")
}

//func (d *DataBase) SuperGroupDeleteAllMessage() error {
//	d.mRWMutex.Lock()
//	defer d.mRWMutex.Unlock()
//	err := d.conn.Model(&model_struct.LocalChatLog{}).Exec("update local_chat_logs set status = ?,content = ? ", constant.MsgStatusHasDeleted, "").Error
//	return utils.Wrap(err, "delete all message error")
//}

func (d *DataBase) SuperGroupUpdateMessageStatusBySourceID(sourceID string, status, sessionType int32) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var condition string
	if sourceID == d.loginUserID && sessionType == constant.SingleChatType {
		condition = "send_id=? And recv_id=? AND session_type=?"
	} else {
		condition = "(send_id=? or recv_id=?)AND session_type=?"
	}
	t := d.conn.Table(utils.GetSuperGroupTableName(sourceID)).Where(condition, sourceID, sourceID, sessionType).Updates(model_struct.LocalChatLog{Status: status})
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateMessageStatusBySourceID failed")
}
func (d *DataBase) SuperGroupUpdateMessageTimeAndStatus(msg *sdk_struct.MsgStruct) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Table(utils.GetSuperGroupTableName(msg.GroupID)).Where("client_msg_id=? And seq=?", msg.ClientMsgID, 0).Updates(model_struct.LocalChatLog{Status: msg.Status, SendTime: msg.SendTime, ServerMsgID: msg.ServerMsgID})
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "SuperGroupUpdateMessageTimeAndStatus failed")
}

func (d *DataBase) SuperGroupGetMessageList(sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
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

	err = utils.Wrap(d.conn.Table(utils.GetSuperGroupTableName(sourceID)).Where(condition, sourceID, constant.MsgStatusSendFailed, sessionType, startTime).
		Order(timeOrder).Offset(0).Limit(count).Find(&messageList).Error, "GetMessageList failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}
func (d *DataBase) SuperGroupGetMessageListNoTime(sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	d.initSuperLocalChatLog(sourceID)
	var messageList []model_struct.LocalChatLog
	var condition, timeOrder string
	if isReverse {
		timeOrder = "send_time ASC"
	} else {
		timeOrder = "send_time DESC"
	}
	condition = "recv_id = ? AND status <=? And session_type = ? "

	err = utils.Wrap(d.conn.Table(utils.GetSuperGroupTableName(sourceID)).Where(condition, sourceID, constant.MsgStatusSendFailed, sessionType).
		Order(timeOrder).Offset(0).Limit(count).Find(&messageList).Error, "GetMessageList failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}

func (d *DataBase) SuperGroupGetSendingMessageList(groupID string) (result []*model_struct.LocalChatLog, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var messageList []model_struct.LocalChatLog
	err = utils.Wrap(d.conn.Table(utils.GetSuperGroupTableName(groupID)).Where("status = ?", constant.MsgStatusSending).Find(&messageList).Error, "GetMessageList failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}

func (d *DataBase) SuperGroupUpdateGroupMessageHasRead(msgIDList []string, groupID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Table(utils.GetSuperGroupTableName(groupID)).Where(" client_msg_id in ?", msgIDList).Update("is_read", constant.HasRead)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateMessageStatusBySourceID failed")
}
func (d *DataBase) SuperGroupUpdateGroupMessageFields(msgIDList []string, groupID string, args map[string]interface{}) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Table(utils.GetSuperGroupTableName(groupID)).Where(" client_msg_id in ?", msgIDList).Updates(args)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateMessageStatusBySourceID failed")
}
func (d *DataBase) SuperGroupGetMultipleMessage(conversationIDList []string, groupID string) (result []*model_struct.LocalChatLog, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var messageList []model_struct.LocalChatLog
	err = utils.Wrap(d.conn.Table(utils.GetSuperGroupTableName(groupID)).Where("client_msg_id IN ?", conversationIDList).Order("send_time DESC").Find(&messageList).Error, "GetMultipleMessage failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}

func (d *DataBase) SuperGroupGetNormalMsgSeq() (uint32, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seq uint32
	err := d.conn.Model(model_struct.LocalChatLog{}).Select("IFNULL(max(seq),0)").Find(&seq).Error
	return seq, utils.Wrap(err, "GetNormalMsgSeq")
}
func (d *DataBase) SuperGroupGetNormalMinSeq(groupID string) (uint32, error) {
	var seq uint32
	err := d.conn.Table(utils.GetSuperGroupTableName(groupID)).Select("IFNULL(min(seq),0)").Where("seq >?", 0).Find(&seq).Error
	return seq, utils.Wrap(err, "SuperGroupGetNormalMinSeq")
}
func (d *DataBase) SuperGroupGetTestMessage(seq uint32) (*model_struct.LocalChatLog, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var c model_struct.LocalChatLog
	return &c, utils.Wrap(d.conn.Where("seq = ?",
		seq).Find(&c).Error, "GetTestMessage failed")
}

func (d *DataBase) SuperGroupUpdateMsgSenderNickname(sendID, nickname string, sType int) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Model(model_struct.LocalChatLog{}).Where(
		"send_id = ? and session_type = ? and sender_nick_name != ? ", sendID, sType, nickname).Updates(
		map[string]interface{}{"sender_nick_name": nickname}).Error, utils.GetSelfFuncName()+" failed")
}

func (d *DataBase) SuperGroupUpdateMsgSenderFaceURL(sendID, faceURL string, sType int) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Model(model_struct.LocalChatLog{}).Where(
		"send_id = ? and session_type = ? and sender_face_url != ? ", sendID, sType, faceURL).Updates(
		map[string]interface{}{"sender_face_url": faceURL}).Error, utils.GetSelfFuncName()+" failed")
}
func (d *DataBase) SuperGroupUpdateMsgSenderFaceURLAndSenderNickname(sendID, faceURL, nickname string, sessionType int, groupID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Table(utils.GetSuperGroupTableName(groupID)).Where(
		"send_id = ? and session_type = ?", sendID, sessionType).Updates(
		map[string]interface{}{"sender_face_url": faceURL, "sender_nick_name": nickname}).Error, utils.GetSelfFuncName()+" failed")
}

func (d *DataBase) SuperGroupGetMsgSeqByClientMsgID(clientMsgID string, groupID string) (uint32, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seq uint32
	err := utils.Wrap(d.conn.Table(utils.GetSuperGroupTableName(groupID)).Select("seq").Where("client_msg_id=?", clientMsgID).Take(&seq).Error, utils.GetSelfFuncName()+" failed")
	return seq, err
}

func (d *DataBase) SuperGroupGetMsgSeqListByGroupID(groupID string) ([]uint32, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seqList []uint32
	err := utils.Wrap(d.conn.Model(model_struct.LocalChatLog{}).Select("seq").Where("recv_id=?", groupID).Find(&seqList).Error, utils.GetSelfFuncName()+" failed")
	return seqList, err
}

func (d *DataBase) SuperGroupGetMsgSeqListByPeerUserID(userID string) ([]uint32, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seqList []uint32
	err := utils.Wrap(d.conn.Model(model_struct.LocalChatLog{}).Select("seq").Where("recv_id=? or send_id=?", userID, userID).Find(&seqList).Error, utils.GetSelfFuncName()+" failed")
	return seqList, err
}

func (d *DataBase) SuperGroupGetMsgSeqListBySelfUserID(userID string) ([]uint32, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seqList []uint32
	err := utils.Wrap(d.conn.Model(model_struct.LocalChatLog{}).Select("seq").Where("recv_id=? and send_id=?", userID, userID).Find(&seqList).Error, utils.GetSelfFuncName()+" failed")
	return seqList, err
}
func (d *DataBase) SuperGroupGetAlreadyExistSeqList(groupID string, lostSeqList []uint32) (seqList []uint32, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	err = utils.Wrap(d.conn.Table(utils.GetSuperGroupTableName(groupID)).Where("seq IN ?", lostSeqList).Pluck("seq", &seqList).Error, utils.GetSelfFuncName()+" failed")
	if err != nil {
		return nil, err
	}
	return seqList, nil
}
