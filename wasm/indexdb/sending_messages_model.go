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
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
)

type LocalSendingMessages struct {
}

func NewLocalSendingMessages() *LocalSendingMessages {
	return &LocalSendingMessages{}
}
func (i *LocalSendingMessages) InsertSendingMessage(ctx context.Context, message *model_struct.LocalSendingMessages) error {
	_, err := exec.Exec(utils.StructToJsonString(message))
	return err
}

func (i *LocalSendingMessages) DeleteSendingMessage(ctx context.Context, conversationID, clientMsgID string) error {
	_, err := exec.Exec(conversationID, clientMsgID)
	return err
}
func (i *LocalSendingMessages) GetAllSendingMessages(ctx context.Context) (result []*model_struct.LocalSendingMessages, err error) {
	gList, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalSendingMessages
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, exec.ErrType
		}
	}
}
