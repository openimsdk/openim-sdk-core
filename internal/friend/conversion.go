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

package friend

import (
	"open_im_sdk/pkg/db/model_struct"

	"github.com/OpenIMSDK/protocol/sdkws"
)

func ServerFriendRequestToLocalFriendRequest(info *sdkws.FriendRequest) *model_struct.LocalFriendRequest {
	return &model_struct.LocalFriendRequest{
		FromUserID:   info.FromUserID,
		FromNickname: info.FromNickname,
		FromFaceURL:  info.FromFaceURL,
		//FromGender:    info.FromGender,
		ToUserID:   info.ToUserID,
		ToNickname: info.ToNickname,
		ToFaceURL:  info.ToFaceURL,
		//ToGender:      info.ToGender,
		HandleResult:  info.HandleResult,
		ReqMsg:        info.ReqMsg,
		CreateTime:    info.CreateTime,
		HandlerUserID: info.HandlerUserID,
		HandleMsg:     info.HandleMsg,
		HandleTime:    info.HandleTime,
		Ex:            info.Ex,
		//AttachedInfo:  info.AttachedInfo,
	}
}

func ServerFriendToLocalFriend(info *sdkws.FriendInfo) *model_struct.LocalFriend {
	return &model_struct.LocalFriend{
		OwnerUserID:    info.OwnerUserID,
		FriendUserID:   info.FriendUser.UserID,
		Remark:         info.Remark,
		CreateTime:     info.CreateTime,
		AddSource:      info.AddSource,
		OperatorUserID: info.OperatorUserID,
		Nickname:       info.FriendUser.Nickname,
		FaceURL:        info.FriendUser.FaceURL,
		Ex:             info.Ex,
		//AttachedInfo:   info.FriendUser.AttachedInfo,
	}
}

func ServerBlackToLocalBlack(info *sdkws.BlackInfo) *model_struct.LocalBlack {
	return &model_struct.LocalBlack{
		OwnerUserID:    info.OwnerUserID,
		BlockUserID:    info.BlackUserInfo.UserID,
		CreateTime:     info.CreateTime,
		AddSource:      info.AddSource,
		OperatorUserID: info.OperatorUserID,
		Nickname:       info.BlackUserInfo.Nickname,
		FaceURL:        info.BlackUserInfo.FaceURL,
		Ex:             info.Ex,
		//AttachedInfo:   info.FriendUser.AttachedInfo,
	}
}
