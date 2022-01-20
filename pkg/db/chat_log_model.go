package db

import (
	"errors"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) BatchInsertMessageList(MessageList []*LocalChatLog) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(MessageList).Error, "BatchInsertMessageList failed")
}
func (d *DataBase) InsertMessage(Message *LocalChatLog) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(Message).Error, "InsertMessage failed")
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
		ClientMsgID).Find(&c).Error, "GetMessage failed")
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
	t := d.conn.Model(LocalChatLog{}).Where("(send_id=? or recv_id=?)AND session_type=?", sourceID, sourceID, sessionType).Updates(LocalChatLog{Status: status})
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateMessageStatusBySourceID failed")
}
func (d *DataBase) UpdateMessageTimeAndStatus(ClientMsgID string, sendTime int64, status int32) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Model(LocalChatLog{}).Where("client_msg_id=? And seq=?", ClientMsgID, 0).Updates(LocalChatLog{Status: status, SendTime: sendTime})
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateMessageStatusBySourceID failed")
}
func (d *DataBase) GetMessageList(sourceID string, sessionType, count int, startTime int64) (result []*LocalChatLog, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var messageList []LocalChatLog
	err = utils.Wrap(d.conn.Where("(send_id = ? OR recv_id = ?) AND status <=? And session_type = ? And send_time < ?", sourceID, sourceID, constant.MsgStatusSendFailed, sessionType, startTime).
		Order("send_time DESC").Offset(0).Limit(count).Find(&messageList).Error, "GetMessageList failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}
func (d *DataBase) GetSelfMessageList(sourceID string, sessionType, count int, startTime int64) (result []*LocalChatLog, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var messageList []LocalChatLog
	err = utils.Wrap(d.conn.Where("send_id = ? And recv_id = ? AND status <=? And session_type = ? And send_time < ?", sourceID, sourceID, constant.MsgStatusSendFailed, sessionType, startTime).
		Order("send_time DESC").Offset(0).Limit(count).Find(&messageList).Error, "GetMessageList failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}
func (d *DataBase) UpdateMessageHasRead(sendID string, msgIDList []string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Model(LocalChatLog{}).Where("send_id=? And is_read=?AND session_type=?AND client_msg_id in ?", sendID, constant.NotRead, constant.SingleChatType, msgIDList).Updates(LocalChatLog{IsRead: constant.HasRead})
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateMessageStatusBySourceID failed")
}

func (d *DataBase) GetNormalMsgSeq() (uint32, error) {
	return 0, nil
}

func (d *DataBase) GetAbnormalMsgSeq() (uint32, error) {
	return 0, nil
}
