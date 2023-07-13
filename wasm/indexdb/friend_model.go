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
	"open_im_sdk/wasm/indexdb/temp_struct"
)

type Friend struct {
	loginUserID string
}

func NewFriend(loginUserID string) *Friend {
	return &Friend{loginUserID: loginUserID}
}

func (i *Friend) InsertFriend(ctx context.Context, friend *model_struct.LocalFriend) error {
	_, err := exec.Exec(utils.StructToJsonString(friend))
	return err
}

func (i *Friend) DeleteFriendDB(ctx context.Context, friendUserID string) error {
	_, err := exec.Exec(friendUserID, i.loginUserID)
	return err
}

func (i *Friend) UpdateFriend(ctx context.Context, friend *model_struct.LocalFriend) error {
	tempLocalFriend := temp_struct.LocalFriend{
		OwnerUserID:    friend.OwnerUserID,
		FriendUserID:   friend.FriendUserID,
		Remark:         friend.Remark,
		CreateTime:     friend.CreateTime,
		AddSource:      friend.AddSource,
		OperatorUserID: friend.OperatorUserID,
		Nickname:       friend.Nickname,
		FaceURL:        friend.FaceURL,
		Ex:             friend.Ex,
		AttachedInfo:   friend.AttachedInfo,
	}
	_, err := exec.Exec(utils.StructToJsonString(tempLocalFriend))
	return err
}

func (i *Friend) GetAllFriendList(ctx context.Context) (result []*model_struct.LocalFriend, err error) {
	gList, err := exec.Exec(i.loginUserID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalFriend
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

func (i *Friend) SearchFriendList(ctx context.Context, keyword string, isSearchUserID, isSearchNickname, isSearchRemark bool) (result []*model_struct.LocalFriend, err error) {
	gList, err := exec.Exec(keyword, isSearchUserID, isSearchNickname, isSearchRemark)
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalFriend
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

func (i *Friend) GetFriendInfoByFriendUserID(ctx context.Context, FriendUserID string) (*model_struct.LocalFriend, error) {
	c, err := exec.Exec(FriendUserID, i.loginUserID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := c.(string); ok {
			result := model_struct.LocalFriend{}
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

func (i *Friend) GetFriendInfoList(ctx context.Context, friendUserIDList []string) (result []*model_struct.LocalFriend, err error) {
	gList, err := exec.Exec(utils.StructToJsonString(friendUserIDList))
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalFriend
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
func (i *Friend) GetPageFriendList(ctx context.Context, offset, count int) (result []*model_struct.LocalFriend, err error) {
	gList, err := exec.Exec(offset, count, i.loginUserID)
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
