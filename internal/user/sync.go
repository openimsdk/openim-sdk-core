// Copyright © 2023 OpenIM SDK.
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
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"

	sdk "open_im_sdk/pkg/sdk_params_callback"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"gorm.io/gorm"
)

// SyncLoginUserInfo synchronizes the user's information from the server.
func (u *User) SyncLoginUserInfo(ctx context.Context) error {
	log.NewInfo(utils.GetSelfFuncName(), "args: ")
	remoteUser, err := u.GetSingleUserFromSvr(ctx, u.loginUserID)
	if err != nil {
		return err
	}
	localUser, err := u.GetLoginUser(ctx, u.loginUserID)
	if err != nil && errs.Unwrap(err) != gorm.ErrRecordNotFound {
		return err
	}
	var remoteUsers []*model_struct.LocalUser
	if err == nil {
		remoteUsers = []*model_struct.LocalUser{localUser}
	}
	log.ZDebug(ctx, "SyncLoginUserInfo", "remoteUser", remoteUser, "localUser", localUser)

	err = u.userSyncer.Sync(ctx, []*model_struct.LocalUser{remoteUser}, remoteUsers, nil)
	if err != nil {
		return err
	}
	callbackData := sdk.SelfInfoUpdatedCallback(*remoteUser)
	if u.listener == nil {
		log.Error("u.listener == nil")
		return nil
	}
	u.listener.OnSelfInfoUpdated(utils.StructToJsonString(callbackData))
	log.Info("OnSelfInfoUpdated", utils.StructToJsonString(callbackData))
	if localUser.Nickname == remoteUser.Nickname && localUser.FaceURL == remoteUser.FaceURL {
		log.NewInfo("OnSelfInfoUpdated nickname faceURL unchanged", callbackData)
		return nil
	}
	_ = common.TriggerCmdUpdateMessage(common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName, Args: common.UpdateMessageInfo{UserID: callbackData.UserID, FaceURL: callbackData.FaceURL, Nickname: callbackData.Nickname}}, u.conversationCh)

	return nil
}
