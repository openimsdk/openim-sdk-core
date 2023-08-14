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
	userPb "github.com/OpenIMSDK/protocol/user"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"gorm.io/gorm"
)

func (u *User) SyncLoginUserInfo(ctx context.Context) error {
	remoteUser, err := u.GetSingleUserFromSvr(ctx, u.loginUserID)
	if err != nil {
		return err
	}
	localUser, err := u.GetLoginUser(ctx, u.loginUserID)
	if err != nil && errs.Unwrap(err) != gorm.ErrRecordNotFound {
		log.ZError(ctx, "SyncLoginUserInfo", err)
	}
	var localUsers []*model_struct.LocalUser
	if err == nil {
		localUsers = []*model_struct.LocalUser{localUser}
	}
	log.ZDebug(ctx, "SyncLoginUserInfo", "remoteUser", remoteUser, "localUser", localUser)
	return u.userSyncer.Sync(ctx, []*model_struct.LocalUser{remoteUser}, localUsers, nil)
}

func (u *User) SyncUserStatus(ctx context.Context, fromId string, toUserID string, status int32, platformID int32, c func(userID string, statusMap *userPb.OnlineStatus)) {
	statusMap := userPb.OnlineStatus{
		UserID:     fromId,
		Status:     status,
		PlatformID: platformID,
	}
	c(fromId, &statusMap)
	u.listener.OnUserStatusChanged(utils.StructToJsonString(statusMap))
}
