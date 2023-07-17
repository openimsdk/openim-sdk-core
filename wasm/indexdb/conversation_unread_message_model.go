// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build js && wasm
// +build js,wasm

package indexdb

import (
	"context"
	"open_im_sdk/wasm/exec"
)

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

type LocalConversationUnreadMessages struct {
}

func NewLocalConversationUnreadMessages() *LocalConversationUnreadMessages {
	return &LocalConversationUnreadMessages{}
}

func (i *LocalConversationUnreadMessages) BatchInsertConversationUnreadMessageList(ctx context.Context, messageList []*model_struct.LocalConversationUnreadMessage) error {
	if messageList == nil {
		return nil
	}
	_, err := exec.Exec(utils.StructToJsonString(messageList))
	return err
}

func (i *LocalConversationUnreadMessages) DeleteConversationUnreadMessageList(ctx context.Context, conversationID string, sendTime int64) int64 {
	deleteRows, err := exec.Exec(conversationID, sendTime)
	if err != nil {
		return 0
	} else {
		if v, ok := deleteRows.(float64); ok {
			var result int64
			result = int64(v)
			return result
		} else {
			return 0
		}
	}
}
