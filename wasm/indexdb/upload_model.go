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

type LocalUpload struct{}

func NewLocalUpload() *LocalUpload {
	return &LocalUpload{}
}

func (i *LocalUpload) GetUpload(ctx context.Context, partHash string) (*model_struct.LocalUpload, error) {
	c, err := exec.Exec(partHash)
	if err != nil {
		return nil, err
	} else {
		if v, ok := c.(string); ok {
			result := model_struct.LocalUpload{}
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

func (i *LocalUpload) InsertUpload(ctx context.Context, upload *model_struct.LocalUpload) error {
	_, err := exec.Exec(utils.StructToJsonString(upload))
	return err
}

func (i *LocalUpload) DeleteUpload(ctx context.Context, partHash string) error {
	_, err := exec.Exec(partHash)
	return err
}
func (i *LocalUpload) UpdateUpload(ctx context.Context, upload *model_struct.LocalUpload) error {
	_, err := exec.Exec(utils.StructToJsonString(upload))
	return err
}

func (i *LocalUpload) DeleteExpireUpload(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
