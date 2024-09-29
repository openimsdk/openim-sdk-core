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
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/protocol/auth"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
)

var (
	ctx context.Context
)

func init() {
	fmt.Println("------------------------>>>>>>>>>>>>>>>>>>> test init func <<<<<<<<<<<<<<<<<<<------------------------")
	rand.Seed(time.Now().UnixNano())
	listner := &OnConnListener{}
	config := getConf(APIADDR, WSADDR)
	config.DataDir = "./"
	configData, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	isInit := open_im_sdk.InitSDK(listner, "test", string(configData))
	if !isInit {
		panic("init sdk failed")
	}
	ctx = open_im_sdk.UserForSDK.Context()
	ctx = ccontext.WithOperationID(ctx, "initOperationID_"+strconv.Itoa(int(time.Now().UnixMilli())))
	token, err := GetAdminToken(ctx, UserID, Secret)
	if err != nil {
		panic(err)
	}
	if err := open_im_sdk.UserForSDK.Login(ctx, UserID, token); err != nil {
		panic(err)
	}
	open_im_sdk.UserForSDK.SetConversationListener(&onConversationListener{ctx: ctx})
	open_im_sdk.UserForSDK.SetGroupListener(&onGroupListener{ctx: ctx})
	open_im_sdk.UserForSDK.SetAdvancedMsgListener(&onAdvancedMsgListener{ctx: ctx})
	open_im_sdk.UserForSDK.SetFriendshipListener(&onFriendshipListener{ctx: ctx})
	open_im_sdk.UserForSDK.SetUserListener(&onUserListener{ctx: ctx})
}

func getConf(APIADDR, WSADDR string) sdk_struct.IMConfig {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = APIADDR
	cf.WsAddr = WSADDR
	cf.DataDir = "./"
	cf.LogLevel = 6
	cf.IsExternalExtensions = true
	cf.PlatformID = PlatformID
	cf.LogFilePath = "./"
	cf.IsLogStandardOutput = true
	return cf
}

func GetAdminToken(ctx context.Context, userID string, secret string) (string, error) {
	req := &auth.GetAdminTokenReq{
		UserID: userID,
		Secret: secret,
	}
	return api.ExtractField(ctx, api.GetAdminToken.Invoke, req, (*auth.GetAdminTokenResp).GetToken)
}
