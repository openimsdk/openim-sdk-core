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
	"testing"
	"time"

	"github.com/openimsdk/protocol/wrapperspb"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"

	"github.com/openimsdk/protocol/sdkws"
)

func Test_GetSelfUserInfo(t *testing.T) {
	userInfo, err := open_im_sdk.IMUserContext.User().GetSelfUserInfo(ctx)
	if err != nil {
		t.Error(err)
	}

	t.Log(userInfo)
}

func Test_SetSelfInfoEx(t *testing.T) {
	newNickName := "test"
	//newFaceURL := "http://test.com"
	err := open_im_sdk.IMUserContext.User().SetSelfInfo(ctx, &sdkws.UserInfoWithEx{
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
	userInfo, err := open_im_sdk.IMUserContext.User().GetSelfUserInfo(ctx)
	if err != nil {
		t.Error(err)
	}
	if userInfo.UserID != UserID && userInfo.Nickname != newNickName && userInfo.FaceURL != newFaceURL {
		t.Error("user id not match")
	}
	t.Log(userInfo)
	time.Sleep(time.Second * 10)
}
