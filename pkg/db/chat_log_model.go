package db

import (
	"errors"
	"fmt"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) BatchInsertMessageList(MessageList []*LocalChatLog) error {
	if MessageList == nil {
		return nil
	}
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(MessageList).Error, "BatchInsertMessageList failed")
}
func (d *DataBase) InsertMessage(Message *LocalChatLog) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(Message).Error, "InsertMessage failed")
}
func (d *DataBase) SearchMessageByKeyword(keyword string, startTime, endTime int64, sessionType int) (result []*LocalChatLog, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var messageList []LocalChatLog
	var condition string
	switch sessionType {
	case constant.SingleChatType:
		if startTime == endTime {
			condition = fmt.Sprintf("session_type==%d And send_time<=%d AND status <=%d And content like %q", constant.SingleChatType, startTime, constant.MsgStatusSendFailed, keyword+"%%")
		}
		condition = fmt.Sprintf("session_type==%d And send_time  between %d and %d AND status <=%d And content like %q", constant.SingleChatType, startTime, endTime, constant.MsgStatusSendFailed, keyword+"%%")
	case constant.GroupChatType:
		if startTime == endTime {
			condition = fmt.Sprintf("session_type==%d And send_time<=%d AND status <=%d And content like %q", constant.GroupChatType, startTime, constant.MsgStatusSendFailed, keyword+"%%")
		}
		condition = fmt.Sprintf("session_type==%d And send_time  between %d and %d AND status <=%d And content like %q", constant.GroupChatType, startTime, endTime, constant.MsgStatusSendFailed, keyword+"%%")
	default:
		if startTime == endTime {
			condition = fmt.Sprintf(" send_time<=%d AND status <=%d And content like %q", startTime, constant.MsgStatusSendFailed, keyword+"%%")
		}
		condition = fmt.Sprintf(" send_time  between %d and %d AND status <=%d And content like %q", startTime, endTime, constant.MsgStatusSendFailed, keyword+"%%")
	}
	err = utils.Wrap(d.conn.Where(condition).Order("send_time DESC").Group("recv_id,client_msg_id").Find(&messageList).Error, "InsertMessage failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}

func (d *DataBase) BatchUpdateMessageList(MessageList []*LocalChatLog) error {
	if MessageList == nil {
		return nil
	}

	for _, v := range MessageList {
		v1 := new(LocalChatLog)
		v1.ClientMsgID = v.ClientMsgID
		v1.Seq = v.Seq
		v1.Status = v.Status
		err := d.UpdateMessage(v1)
		if err != nil {
			return utils.Wrap(err, "BatchUpdateMessageList failed")
		}

	}
	return nil
}

func (d *DataBase) MessageIfExists(ClientMsgID string) (bool, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var count int64
	t := d.conn.Model(&LocalChatLog{}).Where("client_msg_id = ?",
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
func (d *DataBase) IsExistsInErrChatLogBySeq(seq int64) bool {
	return true
}
func (d *DataBase) MessageIfExistsBySeq(seq int64) (bool, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var count int64
	t := d.conn.Model(&LocalChatLog{}).Where("seq = ?",
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
func (d *DataBase) GetMessage(ClientMsgID string) (*LocalChatLog, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var c LocalChatLog
	return &c, utils.Wrap(d.conn.Where("client_msg_id = ?",
		ClientMsgID).Take(&c).Error, "GetMessage failed")
}
func (d *DataBase) UpdateColumnsMessage(ClientMsgID string, args map[string]interface{}) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	c := LocalChatLog{ClientMsgID: ClientMsgID}
	t := d.conn.Model(&c).Updates(args)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateColumnsConversation failed")
}
func (d *DataBase) UpdateMessage(c *LocalChatLog) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Updates(c)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateMessage failed")
}
func (d *DataBase) UpdateMessageStatusBySourceID(sourceID string, status, sessionType int32) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var condition string
	if sourceID == d.loginUserID && sessionType == constant.SingleChatType {
		condition = "send_id=? And recv_id=? AND session_type=?"
	} else {
		condition = "(send_id=? or recv_id=?)AND session_type=?"
	}
	t := d.conn.Model(LocalChatLog{}).Where(condition, sourceID, sourceID, sessionType).Updates(LocalChatLog{Status: status})
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateMessageStatusBySourceID failed")
}
func (d *DataBase) UpdateMessageTimeAndStatus(clientMsgID string, serverMsgID string, sendTime int64, status int32) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Model(LocalChatLog{}).Where("client_msg_id=? And seq=?", clientMsgID, 0).Updates(LocalChatLog{Status: status, SendTime: sendTime, ServerMsgID: serverMsgID})
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateMessageStatusBySourceID failed")
}

func (d *DataBase) GetMessageList(sourceID string, sessionType, count int, startTime int64) (result []*LocalChatLog, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var messageList []LocalChatLog
	var condition string
	if sessionType == constant.SingleChatType && sourceID == d.loginUserID {
		condition = "send_id = ? And recv_id = ? AND status <=? And session_type = ? And send_time < ?"
	} else {
		condition = "(send_id = ? OR recv_id = ?) AND status <=? And session_type = ? And send_time < ?"
	}
	err = utils.Wrap(d.conn.Where(condition, sourceID, sourceID, constant.MsgStatusSendFailed, sessionType, startTime).
		Order("send_time DESC").Offset(0).Limit(count).Find(&messageList).Error, "GetMessageList failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}

func (d *DataBase) GetSendingMessageList() (result []*LocalChatLog, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var messageList []LocalChatLog
	err = utils.Wrap(d.conn.Where("status = ?", constant.MsgStatusSending).Find(&messageList).Error, "GetMessageList failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}

func (d *DataBase) UpdateMessageHasRead(sendID string, msgIDList []string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Model(LocalChatLog{}).Where("send_id=?  AND session_type=? AND client_msg_id in ?", sendID, constant.SingleChatType, msgIDList).Update("is_read", constant.HasRead)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateMessageStatusBySourceID failed")
}
func (d *DataBase) GetMultipleMessage(conversationIDList []string) (result []*LocalChatLog, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var messageList []LocalChatLog
	err = utils.Wrap(d.conn.Where("client_msg_id IN ?", conversationIDList).Find(&messageList).Error, "GetMultipleMessage failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}

func (d *DataBase) GetNormalMsgSeq() (uint32, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seq uint32
	err := d.conn.Model(LocalChatLog{}).Select("IFNULL(max(seq),0)").Find(&seq).Error
	return seq, utils.Wrap(err, "GetNormalMsgSeq")
}
func (d *DataBase) GetTestMessage(seq uint32) (*LocalChatLog, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var c LocalChatLog
	return &c, utils.Wrap(d.conn.Where("seq = ?",
		seq).Find(&c).Error, "GetTestMessage failed")
}

func (d *DataBase) UpdateMsgSenderNickname(sendID, nickname string, sType int) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Model(LocalChatLog{}).Where(
		"send_id = ? and session_type = ? and sender_nick_name != ? ", sendID, sType, nickname).Updates(
		map[string]interface{}{"sender_nick_name": nickname}).Error, utils.GetSelfFuncName()+" failed")
}

func (d *DataBase) UpdateMsgSenderFaceURL(sendID, faceURL string, sType int) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Model(LocalChatLog{}).Where(
		"send_id = ? and session_type = ? and sender_face_url != ? ", sendID, sType, faceURL).Updates(
		map[string]interface{}{"sender_face_url": faceURL}).Error, utils.GetSelfFuncName()+" failed")
}
func (d *DataBase) UpdateMsgSenderFaceURLAndSenderNickname(sendID, faceURL, nickname string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Model(LocalChatLog{}).Where(
		"send_id = ?  and (sender_face_url != ? or sender_nick_name != ?)", sendID, faceURL, nickname).Updates(
		map[string]interface{}{"sender_face_url": faceURL, "sender_nick_name": nickname}).Error, utils.GetSelfFuncName()+" failed")
}

func (d *DataBase) GetMsgSeqByClientMsgID(clientMsgID string) (uint32, error) {
	return 0, nil
}

func (d *DataBase) GetMsgSeqListByGroupID(groupID string) ([]uint32, error) {
	return nil, nil
}

func (d *DataBase) GetMsgSeqListByPeerUserID(userID string) ([]uint32, error) {
	return nil, nil
}

func (d *DataBase) GetMsgSeqListBySelfUserID(userID string) ([]uint32, error) {
	return nil, nil
}
