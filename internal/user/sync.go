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

package user

import (
	"context"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"

	sdk "open_im_sdk/pkg/sdk_params_callback"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"gorm.io/gorm"
)

func (u *User) SyncLoginUserInfo(ctx context.Context) error {
	remoteUser, err := u.GetSingleUserFromSvr(ctx, u.loginUserID)
	if err != nil {
		return err
	}
	localUser, err := u.GetLoginUser(ctx, u.loginUserID)
	if err != nil && errs.Unwrap(err) != gorm.ErrRecordNotFound {
		return err
	}
	var localUsers []*model_struct.LocalUser
	if err == nil {
		localUsers = []*model_struct.LocalUser{localUser}
	}
	log.ZDebug(ctx, "SyncLoginUserInfo", "remoteUser", remoteUser, "localUser", localUser)
	err = u.userSyncer.Sync(ctx, []*model_struct.LocalUser{remoteUser}, localUsers, nil)
	if err != nil {
		return err
	}
	callbackData := sdk.SelfInfoUpdatedCallback(*remoteUser)
	if u.listener == nil {
		return err
	}
	u.listener.OnSelfInfoUpdated(utils.StructToJsonString(callbackData))
	if localUser.Nickname == remoteUser.Nickname && localUser.FaceURL == remoteUser.FaceURL {
		return err
	}
	return nil
}
