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
	userPb "github.com/OpenIMSDK/protocol/user"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"

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

func (u *User) SyncUserStatus(ctx context.Context, fromUserID string, status int32, platformID int32) {
	userOnlineStatus := userPb.OnlineStatus{
		UserID:      fromUserID,
		Status:      status,
		PlatformIDs: []int32{platformID},
	}
	if v, ok := u.OnlineStatusCache.Load(fromUserID); ok {
		if status == constant.Online {
			v.PlatformIDs = utils.RemoveRepeatedElementsInList(append(v.PlatformIDs, platformID))
			u.OnlineStatusCache.Store(fromUserID, v)
		} else {
			v.PlatformIDs = utils.RemoveOneInList(v.PlatformIDs, platformID)
			if len(v.PlatformIDs) == 0 {
				v.Status = constant.Offline
				v.PlatformIDs = []int32{}
				u.OnlineStatusCache.Delete(fromUserID)
			}
		}
		u.listener.OnUserStatusChanged(utils.StructToJsonString(v))
	} else {
		if status == constant.Online {
			u.OnlineStatusCache.Store(fromUserID, &userOnlineStatus)
			u.listener.OnUserStatusChanged(utils.StructToJsonString(userOnlineStatus))
		} else {
			log.ZWarn(ctx, "exception", errors.New("user not exist"), "fromUserID", fromUserID,
				"status", status, "platformID", platformID)
		}
	}
}
