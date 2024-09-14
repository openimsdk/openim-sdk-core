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

package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/openimsdk/protocol/user"
	"github.com/openimsdk/protocol/wrapperspb"

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

func Test_SetSelfInfoEx(t *testing.T) {
	newNickName := "test"
	//newFaceURL := "http://test.com"
	err := open_im_sdk.UserForSDK.User().SetSelfInfo(ctx, &sdkws.UserInfoWithEx{
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
