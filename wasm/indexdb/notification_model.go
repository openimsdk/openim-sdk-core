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

type NotificationSeqs struct {
}

func NewNotificationSeqs() *NotificationSeqs {
	return &NotificationSeqs{}
}

func (i *NotificationSeqs) SetNotificationSeq(ctx context.Context, conversationID string, seq int64) error {
	_, err := exec.Exec(conversationID, seq)
	return err
}

func (i *NotificationSeqs) BatchInsertNotificationSeq(ctx context.Context, notificationSeqs []*model_struct.NotificationSeqs) error {
	_, err := exec.Exec(utils.StructToJsonString(notificationSeqs))
	return err
}

func (i *NotificationSeqs) GetNotificationAllSeqs(ctx context.Context) (result []*model_struct.NotificationSeqs, err error) {
	gList, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return result, err
		} else {
			return nil, exec.ErrType
		}
	}
}
