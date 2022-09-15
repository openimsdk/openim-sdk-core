package db

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) BatchInsertConversationUnreadMessageList(messageList []*model_struct.LocalConversationUnreadMessage) error {
	if messageList == nil {
		return nil
	}
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(messageList).Error, "BatchInsertConversationUnreadMessageList failed")
}
func (d *DataBase) DeleteConversationUnreadMessageList(conversationID string, sendTime int64) int64 {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return d.conn.Where("conversation_id = ? and send_time <= ?", conversationID, sendTime).Delete(&model_struct.LocalConversationUnreadMessage{}).RowsAffected
}
