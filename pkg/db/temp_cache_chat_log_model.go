//go:build !js
// +build !js

package db

import "context"

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) BatchInsertTempCacheMessageList(ctx context.Context, MessageList []*model_struct.TempCacheLocalChatLog) error {
	if MessageList == nil {
		return nil
	}
	return utils.Wrap(d.conn.Create(MessageList).Error, "BatchInsertTempCacheMessageList failed")
}
func (d *DataBase) InsertTempCacheMessage(ctx context.Context, Message *model_struct.TempCacheLocalChatLog) error {

	return utils.Wrap(d.conn.Create(Message).Error, "InsertTempCacheMessage failed")

}
