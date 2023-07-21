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
	"open_im_sdk/wasm/indexdb/temp_struct"
)

type Black struct {
	loginUserID string
}

func NewBlack(loginUserID string) *Black {
	return &Black{loginUserID: loginUserID}
}

// GetBlackListDB gets the blacklist list from the database
func (i Black) GetBlackListDB(ctx context.Context) (result []*model_struct.LocalBlack, err error) {
	gList, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalBlack
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

// GetBlackListUserID gets the list of blocked user IDs
func (i Black) GetBlackListUserID(ctx context.Context) (result []string, err error) {
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

// GetBlackInfoByBlockUserID gets the information of a blocked user by their user ID
func (i Black) GetBlackInfoByBlockUserID(ctx context.Context, blockUserID string) (result *model_struct.LocalBlack, err error) {
	gList, err := exec.Exec(blockUserID, i.loginUserID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp model_struct.LocalBlack
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return &temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

// GetBlackInfoList gets the information of multiple blocked users by their user IDs
func (i Black) GetBlackInfoList(ctx context.Context, blockUserIDList []string) (result []*model_struct.LocalBlack, err error) {
	gList, err := exec.Exec(utils.StructToJsonString(blockUserIDList))
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalBlack
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

// InsertBlack inserts a new blocked user into the database
func (i Black) InsertBlack(ctx context.Context, black *model_struct.LocalBlack) error {
	_, err := exec.Exec(utils.StructToJsonString(black))
	return err
}

// UpdateBlack updates the information of a blocked user in the database
func (i Black) UpdateBlack(ctx context.Context, black *model_struct.LocalBlack) error {
	tempLocalBlack := temp_struct.LocalBlack{
		Nickname:       black.Nickname,
		FaceURL:        black.FaceURL,
		CreateTime:     black.CreateTime,
		AddSource:      black.AddSource,
		OperatorUserID: black.OperatorUserID,
		Ex:             black.Ex,
		AttachedInfo:   black.AttachedInfo,
	}
	_, err := exec.Exec(black.OwnerUserID, black.BlockUserID, utils.StructToJsonString(tempLocalBlack))
	return err
}

// DeleteBlack removes a blocked user from the database
func (i Black) DeleteBlack(ctx context.Context, blockUserID string) error {
	_, err := exec.Exec(blockUserID, i.loginUserID)
	return err
}
