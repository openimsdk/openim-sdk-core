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
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"

	"github.com/OpenIMSDK/protocol/sdkws"
)

func (u *User) GetUsersInfo(ctx context.Context, userIDs []string) ([]*model_struct.LocalUser, error) {
	return u.GetUsersInfoFromSvr(ctx, userIDs)
}

func (u *User) GetSelfUserInfo(ctx context.Context) (*model_struct.LocalUser, error) {
	return u.getSelfUserInfo(ctx)
}

func (u *User) SetSelfInfo(ctx context.Context, userInfo *sdkws.UserInfo) error {
	return u.updateSelfUserInfo(ctx, userInfo)
}

func (u *User) UpdateMsgSenderInfo(ctx context.Context, nickname, faceURL string) (err error) {
	if nickname != "" {
		if err = u.DataBase.UpdateMsgSenderNickname(ctx, u.loginUserID, nickname, constant.SingleChatType); err != nil {
			return err
		}
	}
	if faceURL != "" {
		if err = u.DataBase.UpdateMsgSenderFaceURL(ctx, u.loginUserID, faceURL, constant.SingleChatType); err != nil {
			return err
		}
	}
	return nil
}

func (u *User) SubscribeUsersStatus(ctx context.Context, userID string, userIDs []string) ([]*userPb.OnlineStatus, error) {
	return u.subscribeUsersStatus(ctx, userID, userIDs)
}

func (u *User) UnsubscribeUsersStatus(ctx context.Context, userID string, userIDs []string) error {
	return u.unsubscribeUsersStatus(ctx, userID, userIDs)
}

func (u *User) GetSubscribeUsersStatus(ctx context.Context, userID string) ([]*userPb.OnlineStatus, error) {
	return u.getSubscribeUsersStatus(ctx, userID)
}

func (u *User) GetUserStatus(ctx context.Context, userID string, userIDs []string) ([]*userPb.OnlineStatus, error) {
	return u.getUserStatus(ctx, userID, userIDs)
}
