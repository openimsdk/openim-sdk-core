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

type LocalGroupRequest struct {
}

func NewLocalGroupRequest() *LocalGroupRequest {
	return &LocalGroupRequest{}
}

func (i *LocalGroupRequest) InsertGroupRequest(ctx context.Context, groupRequest *model_struct.LocalGroupRequest) error {
	_, err := exec.Exec(utils.StructToJsonString(groupRequest))
	return err
}

func (i *LocalGroupRequest) DeleteGroupRequest(ctx context.Context, groupID, userID string) error {
	_, err := exec.Exec(groupID, userID)
	return err
}

func (i *LocalGroupRequest) UpdateGroupRequest(ctx context.Context, groupRequest *model_struct.LocalGroupRequest) error {
	_, err := exec.Exec(utils.StructToJsonString(groupRequest))
	return err
}

func (i *LocalGroupRequest) GetSendGroupApplication(ctx context.Context) ([]*model_struct.LocalGroupRequest, error) {
	result, err := exec.Exec()
	if err != nil {
		return nil, err
	}
	if v, ok := result.(string); ok {
		var request []*model_struct.LocalGroupRequest
		if err := utils.JsonStringToStruct(v, &request); err != nil {
			return nil, err
		}
		return request, nil
	} else {
		return nil, exec.ErrType
	}
}

func (i *LocalGroupRequest) InsertAdminGroupRequest(ctx context.Context, groupRequest *model_struct.LocalAdminGroupRequest) error {
	_, err := exec.Exec(utils.StructToJsonString(groupRequest))
	return err
}

func (i *LocalGroupRequest) DeleteAdminGroupRequest(ctx context.Context, groupID, userID string) error {
	_, err := exec.Exec(groupID, userID)
	return err
}

func (i *LocalGroupRequest) UpdateAdminGroupRequest(ctx context.Context, groupRequest *model_struct.LocalAdminGroupRequest) error {
	_, err := exec.Exec(utils.StructToJsonString(groupRequest))
	return err
}

func (i *LocalGroupRequest) GetAdminGroupApplication(ctx context.Context) ([]*model_struct.LocalAdminGroupRequest, error) {
	result, err := exec.Exec()
	if err != nil {
		return nil, err
	}
	if v, ok := result.(string); ok {
		var request []*model_struct.LocalAdminGroupRequest
		if err := utils.JsonStringToStruct(v, &request); err != nil {
			return nil, err
		}
		return request, nil
	} else {
		return nil, exec.ErrType
	}
}
