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

type LocalSuperGroup struct{}

func NewLocalSuperGroup() *LocalSuperGroup {
	return &LocalSuperGroup{}
}

func (i *LocalSuperGroup) GetJoinedSuperGroupList(ctx context.Context) (result []*model_struct.LocalGroup, err error) {
	groupList, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := groupList.(string); ok {
			var temp []model_struct.LocalGroup
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

func (i *LocalSuperGroup) InsertSuperGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error {
	_, err := exec.Exec(utils.StructToJsonString(groupInfo))
	return err
}
func (i *LocalSuperGroup) UpdateSuperGroup(ctx context.Context, g *model_struct.LocalGroup) error {
	_, err := exec.Exec(g.GroupID, utils.StructToJsonString(g))
	return err
}

func (i *LocalSuperGroup) DeleteSuperGroup(ctx context.Context, groupID string) error {
	_, err := exec.Exec(groupID)
	return err
}

func (i *LocalSuperGroup) DeleteAllSuperGroup(ctx context.Context) error {
	_, err := exec.Exec()
	return err
}

func (i *LocalSuperGroup) GetSuperGroupInfoByGroupID(ctx context.Context, groupID string) (*model_struct.LocalGroup, error) {
	groupInfo, err := exec.Exec(groupID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := groupInfo.(string); ok {
			result := model_struct.LocalGroup{}
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

func (i *LocalSuperGroup) GetJoinedWorkingGroupIDList(ctx context.Context) ([]string, error) {
	IDList, err := exec.Exec()
	if err != nil {
		return nil, err
	}
	if v, ok := IDList.(string); ok {
		var temp []string
		err := utils.JsonStringToStruct(v, &temp)
		if err != nil {
			return nil, err
		}
		return temp, nil
	}
	return nil, exec.ErrType
}

func (i *LocalSuperGroup) GetJoinedWorkingGroupList(ctx context.Context) (result []*model_struct.LocalGroup, err error) {
	groupList, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := groupList.(string); ok {
			var temp []model_struct.LocalGroup
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
