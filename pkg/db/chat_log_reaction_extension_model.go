//go:build !js
// +build !js

package db

import (
	"context"
	"errors"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/errs"
)

func (d *DataBase) GetMessageReactionExtension(ctx context.Context, msgID string) (result *model_struct.LocalChatLogReactionExtensions, err error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var l model_struct.LocalChatLogReactionExtensions
	return &l, errs.WrapMsg(d.conn.WithContext(ctx).Where("client_msg_id = ?",
		msgID).Take(&l).Error, "GetMessageReactionExtension failed")
}

func (d *DataBase) InsertMessageReactionExtension(ctx context.Context, messageReactionExtension *model_struct.LocalChatLogReactionExtensions) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Create(messageReactionExtension).Error, "InsertMessageReactionExtension failed")
}

func (d *DataBase) UpdateMessageReactionExtension(ctx context.Context, c *model_struct.LocalChatLogReactionExtensions) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Updates(c)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.WrapMsg(t.Error, "UpdateConversation failed")
}

func (d *DataBase) DeleteMessageReactionExtension(ctx context.Context, msgID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	temp := model_struct.LocalChatLogReactionExtensions{ClientMsgID: msgID}
	return d.conn.WithContext(ctx).Delete(&temp).Error

}
