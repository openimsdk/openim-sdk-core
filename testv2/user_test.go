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

package testv2

import (
	"open_im_sdk/open_im_sdk"
	"testing"
	"time"

	"github.com/OpenIMSDK/protocol/sdkws"
)

func Test_GetSelfUserInfo(t *testing.T) {
	userInfo, err := open_im_sdk.UserForSDK.User().GetSelfUserInfo(ctx)
	if err != nil {
		t.Error(err)
	}

	t.Log(userInfo)
}

func Test_GetUsersInfo(t *testing.T) {
	userInfo, err := open_im_sdk.UserForSDK.Full().GetUsersInfo(ctx, []string{"friendUserID"})
	if err != nil {
		t.Error(err)
	}
	if userInfo[0].BlackInfo != nil {
		t.Log(userInfo[0].BlackInfo)
	}
	if userInfo[0].FriendInfo != nil {
		t.Log(userInfo[0].FriendInfo)
	}
	if userInfo[0].PublicInfo != nil {
		t.Log(userInfo[0].PublicInfo)
	}
}

func Test_SetSelfInfo(t *testing.T) {
	newNickName := "test"
	//newFaceURL := "http://test.com"
	err := open_im_sdk.UserForSDK.User().SetSelfInfo(ctx, &sdkws.UserInfo{
		Nickname: newNickName,
		//FaceURL:  newFaceURL,
	})
	newFaceURL := "http://test.com"

	if err != nil {
		t.Error(err)
	}
	userInfo, err := open_im_sdk.UserForSDK.User().GetSelfUserInfo(ctx)
	if err != nil {
		t.Error(err)
	}
	if userInfo.UserID != UserID && userInfo.Nickname != newNickName && userInfo.FaceURL != newFaceURL {
		t.Error("user id not match")
	}
	t.Log(userInfo)
	time.Sleep(time.Second * 10)
}

func Test_UpdateMsgSenderInfo(t *testing.T) {
	err := open_im_sdk.UserForSDK.User().UpdateMsgSenderInfo(ctx, "test", "http://test.com")
	if err != nil {
		t.Error(err)
	}
	userInfo, err := open_im_sdk.UserForSDK.User().GetSelfUserInfo(ctx)
	if err != nil {
		t.Error(err)
	}
	t.Log(userInfo)
}

func Test_Sub(t *testing.T) {
	var users []string
	users = append(users, "2926672950")
	status, err := open_im_sdk.UserForSDK.User().SubscribeUsersStatus(ctx, "2285951027", users)
	if err != nil {
		t.Error(err)
	}
	t.Log(status)

	for i := 0; i < 20; i++ {
		status, err = open_im_sdk.UserForSDK.User().SubscribeUsersStatus(ctx, "2285951027", users)
		t.Log(status)
		time.Sleep(time.Second * 3)
	}
}

func Test_GetSubscribeUsersStatus(t *testing.T) {
	status, err := open_im_sdk.UserForSDK.User().GetSubscribeUsersStatus(ctx, "2285951027")
	if err != nil {
		return
	}
	t.Log(status)
}

func Test_GetUserStatus(t *testing.T) {
	var UserIDs []string
	UserIDs = append(UserIDs, "2926672950")
	status, err := open_im_sdk.UserForSDK.User().GetUserStatus(ctx, "2285951027", UserIDs)
	if err != nil {
		return
	}
	t.Log(status)
}

func Test_UnSub(t *testing.T) {
	var users []string
	users = append(users, "2926672950")
	err := open_im_sdk.UserForSDK.User().UnsubscribeUsersStatus(ctx, "2285951027", users)
	if err != nil {
		t.Error(err)
	}
}
