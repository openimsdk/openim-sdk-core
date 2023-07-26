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

type LocalGroups struct{}

func NewLocalGroups() *LocalGroups {
	return &LocalGroups{}
}

func (i *LocalGroups) InsertGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error {
	_, err := exec.Exec(utils.StructToJsonString(groupInfo))
	return err
}

func (i *LocalGroups) DeleteGroup(ctx context.Context, groupID string) error {
	_, err := exec.Exec(groupID)
	return err
}

// 该函数需要全更新
func (i *LocalGroups) UpdateGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error {
	_, err := exec.Exec(groupInfo.GroupID, utils.StructToJsonString(groupInfo))
	return err
}

func (i *LocalGroups) GetJoinedGroupListDB(ctx context.Context) (result []*model_struct.LocalGroup, err error) {
	gList, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
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

func (i *LocalGroups) GetGroups(ctx context.Context, groupIDs []string) (result []*model_struct.LocalGroup, err error) {
	gList, err := exec.Exec(utils.StructToJsonString(groupIDs))
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
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

func (i *LocalGroups) GetGroupInfoByGroupID(ctx context.Context, groupID string) (*model_struct.LocalGroup, error) {
	c, err := exec.Exec(groupID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := c.(string); ok {
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

func (i *LocalGroups) GetAllGroupInfoByGroupIDOrGroupName(ctx context.Context, keyword string, isSearchGroupID bool, isSearchGroupName bool) (result []*model_struct.LocalGroup, err error) {
	gList, err := exec.Exec(keyword, isSearchGroupID, isSearchGroupName)
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
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

func (i *LocalGroups) AddMemberCount(ctx context.Context, groupID string) error {
	_, err := exec.Exec(groupID)
	return err
}

func (i *LocalGroups) SubtractMemberCount(ctx context.Context, groupID string) error {
	_, err := exec.Exec(groupID)
	return err
}
func (i *LocalGroups) GetGroupMemberAllGroupIDs(ctx context.Context) (result []string, err error) {
	groupIDList, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := groupIDList.(string); ok {
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
