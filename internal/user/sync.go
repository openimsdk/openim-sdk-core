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
	"errors"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
)

func (u *User) SyncLoginUserInfo(ctx context.Context) error {
	remoteUser, err := u.GetSingleUserFromSvr(ctx, u.loginUserID)
	if err != nil {
		return err
	}
	localUser, err := u.GetLoginUser(ctx, u.loginUserID)
	if err != nil && errors.Is(errs.Unwrap(err), errs.ErrRecordNotFound) {
		return err
	}
	var localUsers []*model_struct.LocalUser
	if err == nil {
		localUsers = []*model_struct.LocalUser{localUser}
	}
	log.ZDebug(ctx, "SyncLoginUserInfo", "remoteUser", remoteUser, "localUser", localUser)
	return u.userSyncer.Sync(ctx, []*model_struct.LocalUser{remoteUser}, localUsers, nil)
}
func (u *User) SyncLoginUserInfoWithoutNotice(ctx context.Context) error {
	remoteUser, err := u.GetSingleUserFromSvr(ctx, u.loginUserID)
	if err != nil {
		return err
	}
	localUser, err := u.GetLoginUser(ctx, u.loginUserID)
	if err != nil && errors.Is(errs.Unwrap(err), errs.ErrRecordNotFound) {
		log.ZError(ctx, "SyncLoginUserInfo", err)
	}
	var localUsers []*model_struct.LocalUser
	if err == nil {
		localUsers = []*model_struct.LocalUser{localUser}
	}
	log.ZDebug(ctx, "SyncLoginUserInfo", "remoteUser", remoteUser, "localUser", localUser)
	return u.userSyncer.Sync(ctx, []*model_struct.LocalUser{remoteUser}, localUsers, nil, false, true)
}

func (u *User) SyncAllCommand(ctx context.Context) error {
	return u.syncAllCommand(ctx, true)
}

func (u *User) SyncAllCommandWithoutNotice(ctx context.Context) error {
	return u.syncAllCommand(ctx, false)
}

func (u *User) syncAllCommand(ctx context.Context, withNotice bool) error {
	resp, err := u.processUserCommandGetAll(ctx)
	if err != nil {
		return err
	}
	localData, err := u.DataBase.ProcessUserCommandGetAll(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "sync command", "data from server", resp, "data from local", localData)
	if withNotice {
		return u.commandSyncer.Sync(ctx, datautil.Batch(ServerCommandToLocalCommand, resp), localData, nil)
	} else {
		return u.commandSyncer.Sync(ctx, datautil.Batch(ServerCommandToLocalCommand, resp), localData, nil, false, true)
	}
}
