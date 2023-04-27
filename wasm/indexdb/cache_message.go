// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

//go:build js && wasm
// +build js,wasm

package indexdb

import "context"

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

type LocalCacheMessage struct {
}

func NewLocalCacheMessage() *LocalCacheMessage {
	return &LocalCacheMessage{}
}

func (i *LocalCacheMessage) BatchInsertTempCacheMessageList(ctx context.Context, MessageList []*model_struct.TempCacheLocalChatLog) error {
	_, err := Exec(utils.StructToJsonString(MessageList))
	return err
}

func (i *LocalCacheMessage) InsertTempCacheMessage(ctx context.Context, Message *model_struct.TempCacheLocalChatLog) error {
	_, err := Exec(utils.StructToJsonString(Message))
	return err
}
