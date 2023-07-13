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
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/exec"
)

type LocalChatLogReactionExtensions struct {
	ExtKey  string `json:"ext_key"`
	ExtVal  string `json:"ext_val"`
	ExtKey2 string `json:"ext_key2"`
	ExtVal2 string `json:"ext_val2"`
}

func NewLocalChatLogReactionExtensions() *LocalChatLogReactionExtensions {
	return &LocalChatLogReactionExtensions{}
}

func (i *LocalChatLogReactionExtensions) GetMessageReactionExtension(ctx context.Context, clientMsgID string) (result *model_struct.LocalChatLogReactionExtensions, err error) {
	msg, err := exec.Exec(clientMsgID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msg.(string); ok {
			result := model_struct.LocalChatLogReactionExtensions{}
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return &result, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalChatLogReactionExtensions) InsertMessageReactionExtension(ctx context.Context, messageReactionExtension *model_struct.LocalChatLogReactionExtensions) error {
	_, err := exec.Exec(utils.StructToJsonString(messageReactionExtension))
	return err
}

//	func (i *LocalChatLogReactionExtensions) GetAndUpdateMessageReactionExtension(ctx context.Context, clientMsgID string, m map[string]*sdkws.KeyValue) error {
//		_, err := exec.Exec(clientMsgID, utils.StructToJsonString(m))
//		return err
//	}
//
//	func (i *LocalChatLogReactionExtensions) DeleteAndUpdateMessageReactionExtension(ctx context.Context, clientMsgID string, m map[string]*sdkws.KeyValue) error {
//		_, err := exec.Exec(clientMsgID, utils.StructToJsonString(m))
//		return err
//	}
func (i *LocalChatLogReactionExtensions) GetMultipleMessageReactionExtension(ctx context.Context, msgIDList []string) (result []*model_struct.LocalChatLogReactionExtensions, err error) {
	msgReactionExtensionList, err := exec.Exec(utils.StructToJsonString(msgIDList))
	if err != nil {
		return nil, err
	} else {
		if v, ok := msgReactionExtensionList.(string); ok {
			var temp []model_struct.LocalChatLogReactionExtensions
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
func (i *LocalChatLogReactionExtensions) DeleteMessageReactionExtension(ctx context.Context, msgID string) error {
	_, err := exec.Exec(msgID)
	return err
}
func (i *LocalChatLogReactionExtensions) UpdateMessageReactionExtension(ctx context.Context, c *model_struct.LocalChatLogReactionExtensions) error {
	_, err := exec.Exec(c.ClientMsgID, utils.StructToJsonString(c))
	return err
}
