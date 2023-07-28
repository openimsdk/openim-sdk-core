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

type FriendRequest struct {
	loginUserID string
}

func NewFriendRequest(loginUserID string) *FriendRequest {
	return &FriendRequest{loginUserID: loginUserID}
}

func (i FriendRequest) InsertFriendRequest(ctx context.Context, friendRequest *model_struct.LocalFriendRequest) error {
	_, err := exec.Exec(utils.StructToJsonString(friendRequest))
	return err
}

func (i FriendRequest) DeleteFriendRequestBothUserID(ctx context.Context, fromUserID, toUserID string) error {
	_, err := exec.Exec(fromUserID, toUserID)
	return err
}

func (i FriendRequest) UpdateFriendRequest(ctx context.Context, friendRequest *model_struct.LocalFriendRequest) error {
	tempLocalFriendRequest := temp_struct.LocalFriendRequest{
		FromUserID:    friendRequest.FromUserID,
		FromNickname:  friendRequest.FromNickname,
		FromFaceURL:   friendRequest.FromFaceURL,
		ToUserID:      friendRequest.ToUserID,
		ToNickname:    friendRequest.ToNickname,
		ToFaceURL:     friendRequest.ToFaceURL,
		HandleResult:  friendRequest.HandleResult,
		ReqMsg:        friendRequest.ReqMsg,
		CreateTime:    friendRequest.CreateTime,
		HandlerUserID: friendRequest.HandlerUserID,
		HandleMsg:     friendRequest.HandleMsg,
		HandleTime:    friendRequest.HandleTime,
		Ex:            friendRequest.Ex,
		AttachedInfo:  friendRequest.AttachedInfo,
	}
	_, err := exec.Exec(utils.StructToJsonString(tempLocalFriendRequest))
	return err
}

func (i FriendRequest) GetRecvFriendApplication(ctx context.Context) (result []*model_struct.LocalFriendRequest, err error) {
	gList, err := exec.Exec(i.loginUserID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalFriendRequest
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

func (i FriendRequest) GetSendFriendApplication(ctx context.Context) (result []*model_struct.LocalFriendRequest, err error) {
	gList, err := exec.Exec(i.loginUserID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalFriendRequest
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

func (i FriendRequest) GetFriendApplicationByBothID(ctx context.Context, fromUserID, toUserID string) (*model_struct.LocalFriendRequest, error) {
	c, err := exec.Exec(fromUserID, toUserID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := c.(string); ok {
			result := model_struct.LocalFriendRequest{}
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

func (i FriendRequest) GetBothFriendReq(ctx context.Context, fromUserID, toUserID string) (result []*model_struct.LocalFriendRequest, err error) {
	gList, err := exec.Exec(fromUserID, toUserID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalFriendRequest
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
