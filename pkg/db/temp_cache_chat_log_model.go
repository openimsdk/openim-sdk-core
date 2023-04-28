// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

//go:build !js
// +build !js

package db

import (
	"context"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) BatchInsertTempCacheMessageList(ctx context.Context, MessageList []*model_struct.TempCacheLocalChatLog) error {
	if MessageList == nil {
		return nil
	}
	return utils.Wrap(d.conn.WithContext(ctx).Create(MessageList).Error, "BatchInsertTempCacheMessageList failed")
}
func (d *DataBase) InsertTempCacheMessage(ctx context.Context, Message *model_struct.TempCacheLocalChatLog) error {

	return utils.Wrap(d.conn.WithContext(ctx).Create(Message).Error, "InsertTempCacheMessage failed")

}
