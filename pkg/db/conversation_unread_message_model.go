//go:build !js
// +build !js

package db

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/errs"
)

func (d *DataBase) BatchInsertConversationUnreadMessageList(ctx context.Context, messageList []*model_struct.LocalConversationUnreadMessage) error {
	if messageList == nil {
		return nil
	}
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Create(messageList).Error, "BatchInsertConversationUnreadMessageList failed")
}
func (d *DataBase) DeleteConversationUnreadMessageList(ctx context.Context, conversationID string, sendTime int64) int64 {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return d.conn.WithContext(ctx).Where("conversation_id = ? and send_time <= ?", conversationID, sendTime).Delete(&model_struct.LocalConversationUnreadMessage{}).RowsAffected
}
