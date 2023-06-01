package db

import (
	"context"
	"errors"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) initChatLog(ctx context.Context, conversationID string) {
	if !d.conn.Migrator().HasTable(utils.GetTableName(conversationID)) {
		d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).AutoMigrate(&model_struct.LocalChatLog{})
	}
}
func (d *DataBase) UpdateMessage(ctx context.Context, conversationID string, c *model_struct.LocalChatLog) error {
	t := d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Updates(c)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update ")
	}
	return utils.Wrap(t.Error, "UpdateMessage failed")
}

func (d *DataBase) UpdateMessageBySeq(ctx context.Context, conversationID string, c *model_struct.LocalChatLog) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Where("seq=?", c.Seq).Updates(c).Error, "UpdateMessage failed")
}

func (d *DataBase) BatchInsertMessageList(ctx context.Context, conversationID string, MessageList []*model_struct.LocalChatLog) error {
	if MessageList == nil {
		return nil
	}
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Create(MessageList).Error, "BatchInsertMessageList failed")
}

func (d *DataBase) InsertMessage(ctx context.Context, conversationID string, Message *model_struct.LocalChatLog) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Create(Message).Error, "InsertMessage failed")
}
func (d *DataBase) GetMessage(ctx context.Context, conversationID string, clientMsgID string) (*model_struct.LocalChatLog, error) {
	d.initChatLog(ctx, conversationID)
	var c model_struct.LocalChatLog
	return &c, utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Where("client_msg_id = ?",
		clientMsgID).Take(&c).Error, "GetMessage failed")
}

func (d *DataBase) GetMessageBySeq(ctx context.Context, conversationID string, seq int64) (*model_struct.LocalChatLog, error) {
	d.initChatLog(ctx, conversationID)
	var c model_struct.LocalChatLog
	return &c, utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Where("seq = ?",
		seq).Take(&c).Error, "GetMessage failed")
}

func (d *DataBase) UpdateMessageTimeAndStatus(ctx context.Context, conversationID, clientMsgID string, serverMsgID string, sendTime int64, status int32) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Model(model_struct.LocalChatLog{}).Where("client_msg_id=? And seq=?", clientMsgID, 0).
		Updates(model_struct.LocalChatLog{Status: status, SendTime: sendTime, ServerMsgID: serverMsgID}).Error, "UpdateMessageStatusBySourceID failed")
}
func (d *DataBase) GetMessageListNoTime(ctx context.Context, conversationID string,
	count int, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	d.initChatLog(ctx, conversationID)
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var timeOrder string
	if isReverse {
		timeOrder = "send_time ASC"
	} else {
		timeOrder = "send_time DESC"
	}
	err = utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Order(timeOrder).Offset(0).Limit(count).Find(&result).Error, "GetMessageList failed")
	if err != nil {
		return nil, err
	}
	return result, err
}
func (d *DataBase) GetMessageList(ctx context.Context, conversationID string, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var condition, timeOrder, timeSymbol string
	if isReverse {
		timeOrder = "send_time ASC"
		timeSymbol = ">"
	} else {
		timeOrder = "send_time DESC"
		timeSymbol = "<"
	}
	condition = "send_time " + timeSymbol + " ?"

	err = utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Where(condition, startTime).
		Order(timeOrder).Offset(0).Limit(count).Find(&result).Error, "GetMessageList failed")
	if err != nil {
		return nil, err
	}
	return result, err
}

func (d *DataBase) DeleteConversationAllMessages(ctx context.Context, conversationID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Where("1 = 1").Delete(model_struct.LocalChatLog{}).Error, "DeleteConversationAllMessages failed")
}

func (d *DataBase) MarkDeleteConversationAllMessages(ctx context.Context, conversationID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Where("1 = 1").Updates(model_struct.LocalChatLog{Status: constant.MsgStatusHasDeleted}).Error, "DeleteConversationAllMessages failed")
}

func (d *DataBase) DeleteConversationMsgs(ctx context.Context, conversationID string, msgIDs []string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Where("client_msg_id IN ?", msgIDs).Delete(model_struct.LocalChatLog{}).Error, "DeleteConversationMsgs failed")
}

func (d *DataBase) DeleteConversationMsgsBySeqs(ctx context.Context, conversationID string, seqs []int64) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Where("seq IN ?", seqs).Delete(model_struct.LocalChatLog{}).Error, "DeleteConversationMsgs failed")
}
