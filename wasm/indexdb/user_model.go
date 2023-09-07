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
	"time"
)

type LocalUsers struct {
}

func NewLocalUsers() *LocalUsers {
	return &LocalUsers{}
}

func (l *LocalUsers) GetLoginUser(ctx context.Context, userID string) (*model_struct.LocalUser, error) {
	user, err := exec.Exec(userID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := user.(string); ok {
			result := model_struct.LocalUser{}
			temp := temp_struct.LocalUser{}
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			result.UserID = temp.UserID
			result.Nickname = temp.Nickname
			result.FaceURL = temp.FaceURL
			result.CreateTime = temp.CreateTime
			result.AppMangerLevel = temp.AppMangerLevel
			result.Ex = temp.Ex
			result.AttachedInfo = temp.Ex
			result.GlobalRecvMsgOpt = temp.GlobalRecvMsgOpt
			return &result, err
		} else {
			return nil, exec.ErrType
		}

	}
}

func (l *LocalUsers) UpdateLoginUser(ctx context.Context, user *model_struct.LocalUser) error {
	_, err := exec.Exec(utils.StructToJsonString(user))
	return err

}
func (l *LocalUsers) UpdateLoginUserByMap(ctx context.Context, user *model_struct.LocalUser, args map[string]interface{}) error {
	if v, ok := args["birth_time"]; ok {
		if t, ok := v.(time.Time); ok {
			args["birth_time"] = utils.TimeToString(t)
		} else {
			return exec.ErrType
		}
	}
	_, err := exec.Exec(user.UserID, args)
	return err
}

func (l *LocalUsers) InsertLoginUser(ctx context.Context, user *model_struct.LocalUser) error {
	temp := temp_struct.LocalUser{}
	temp.UserID = user.UserID
	temp.Nickname = user.Nickname
	temp.FaceURL = user.FaceURL
	temp.CreateTime = user.CreateTime
	temp.AppMangerLevel = user.AppMangerLevel
	temp.Ex = user.Ex
	temp.AttachedInfo = user.Ex
	temp.GlobalRecvMsgOpt = user.GlobalRecvMsgOpt
	_, err := exec.Exec(utils.StructToJsonString(temp))
	return err
}
