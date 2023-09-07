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

type LocalStrangers struct {
}

func NewLocalStrangers() *LocalStrangers {
	return &LocalStrangers{}
}
func (l *LocalStrangers) GetStrangerInfo(ctx context.Context, userIDs []string) (result []*model_struct.LocalStranger, err error) {
	gList, err := exec.Exec(utils.StructToJsonString(userIDs))
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalStranger
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

func (l *LocalStrangers) SetStrangerInfo(ctx context.Context, strangerList []*model_struct.LocalStranger) error {
	_, err := exec.Exec(utils.StructToJsonString(strangerList))
	return err
}
