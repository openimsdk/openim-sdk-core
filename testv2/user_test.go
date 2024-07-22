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
	"fmt"
	"github.com/openimsdk/protocol/user"
	"github.com/openimsdk/protocol/wrapperspb"
	"testing"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"

	"github.com/openimsdk/protocol/sdkws"
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
func Test_GetUsersInfoWithCache(t *testing.T) {
	userInfo, err := open_im_sdk.UserForSDK.Full().GetUsersInfoWithCache(ctx, []string{"1"}, "")
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
func Test_SetSelfInfoEx(t *testing.T) {
	newNickName := "test"
	//newFaceURL := "http://test.com"
	err := open_im_sdk.UserForSDK.User().SetSelfInfoEx(ctx, &sdkws.UserInfoWithEx{
		Nickname: &wrapperspb.StringValue{
			Value: newNickName,
		},
		//FaceURL:  newFaceURL,
		Ex: &wrapperspb.StringValue{
			Value: "ASD",
		},
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

func Test_SetSetGlobalRecvMessageOpt(t *testing.T) {
	err := open_im_sdk.UserForSDK.User().SetGlobalRecvMessageOpt(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
}

//func Test_Sub(t *testing.T) {
//	var users []string
//	users = append(users, "2926672950")
//	status, err := open_im_sdk.UserForSDK.User().SubscribeUsersStatus(ctx, users)
//	if err != nil {
//		t.Error(err)
//	}
//	t.Log(status)
//
//	for i := 0; i < 20; i++ {
//		status, err = open_im_sdk.UserForSDK.User().SubscribeUsersStatus(ctx, users)
//		t.Log(status)
//		time.Sleep(time.Second * 3)
//	}
//}
//
//func Test_GetSubscribeUsersStatus(t *testing.T) {
//	status, err := open_im_sdk.UserForSDK.User().GetSubscribeUsersStatus(ctx)
//	if err != nil {
//		return
//	}
//	t.Log(status)
//}

//func Test_GetUserStatus(t *testing.T) {
//	var UserIDs []string
//	UserIDs = append(UserIDs, "2926672950")
//	status, err := open_im_sdk.UserForSDK.User().GetUserStatus(ctx, UserIDs)
//	if err != nil {
//		return
//	}
//	t.Log(status)
//}
//
//func Test_UnSub(t *testing.T) {
//	var users []string
//	users = append(users, "2926672950")
//	err := open_im_sdk.UserForSDK.User().UnsubscribeUsersStatus(ctx, users)
//	if err != nil {
//		t.Error(err)
//	}
//}

func Test_UserCommandAdd(t *testing.T) {
	// Creating a request with a pointer
	req := &user.ProcessUserCommandAddReq{
		UserID: "3",
		Type:   8,
		Uuid:   "1",
		Value: &wrapperspb.StringValue{
			Value: "ASD",
		},
		Ex: &wrapperspb.StringValue{
			Value: "ASD",
		},
	}

	// Passing the pointer to the function
	err := open_im_sdk.UserForSDK.User().ProcessUserCommandAdd(ctx, req)
	if err != nil {
		// Handle the error
		t.Errorf("Failed to add favorite: %v", err)
	}
}
func Test_UserCommandGet(t *testing.T) {
	// Creating a request with a pointer

	// Passing the pointer to the function
	result, err := open_im_sdk.UserForSDK.User().ProcessUserCommandGetAll(ctx)
	if err != nil {
		// Handle the error
		t.Errorf("Failed to add favorite: %v", err)
	}
	fmt.Printf("%v\n", result)
}
func Test_UserCommandDelete(t *testing.T) {
	// Creating a request with a pointer
	req := &user.ProcessUserCommandDeleteReq{
		UserID: "3",
		Type:   8,
		Uuid:   "1",
	}

	// Passing the pointer to the function
	err := open_im_sdk.UserForSDK.User().ProcessUserCommandDelete(ctx, req)
	if err != nil {
		// Handle the error
		t.Errorf("Failed to add favorite: %v", err)
	}
}
