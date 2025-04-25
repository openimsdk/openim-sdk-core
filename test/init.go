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

var ctx context.Context

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
	ctx = open_im_sdk.IMUserContext.Context()
	ctx = ccontext.WithOperationID(ctx, "initOperationID_"+strconv.Itoa(int(time.Now().UnixMilli())))
	token, err := GetUserToken(ctx, UserID, PlatformID, Secret, config)
	if err != nil {
		panic(err)
	}
	if err := open_im_sdk.IMUserContext.Login(ctx, UserID, token); err != nil {
		panic(err)
	}
	ch := make(chan error)
	open_im_sdk.IMUserContext.SetConversationListener(&onConversationListener{ctx: ctx, ch: ch})
	open_im_sdk.IMUserContext.SetAdvancedMsgListener(&onAdvancedMsgListener{ctx: ctx})
	open_im_sdk.IMUserContext.SetMessageKvInfoListener(&onMessageKvInfoListener{})
	open_im_sdk.IMUserContext.SetFriendshipListener(&onFriendshipListener{ctx: ctx})
	open_im_sdk.IMUserContext.SetGroupListener(&onGroupListener{ctx: ctx})
	open_im_sdk.IMUserContext.SetCustomBusinessListener(&onCustomBusinessListener{ctx: ctx})

	if err := <-ch; err != nil {
		panic(err)
	}
}

func getConf(APIADDR, WSADDR string) sdk_struct.IMConfig {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = APIADDR
	cf.WsAddr = WSADDR
	cf.DataDir = "./"
	cf.LogLevel = 5
	cf.PlatformID = PlatformID
	cf.LogFilePath = "./"
	cf.IsLogStandardOutput = true
	return cf
}

func GetUserToken(ctx context.Context, userID string, platformID int32, secret string, imConf sdk_struct.IMConfig) (string, error) {
	adminReq := &auth.GetAdminTokenReq{
		UserID: AdminUserID,
		Secret: secret,
	}
	adminToken, err := api.ExtractField(ctx, api.GetAdminToken.Invoke, adminReq, (*auth.GetAdminTokenResp).GetToken)
	if err != nil {
		return "", err
	}
	userReq := &auth.GetUserTokenReq{
		UserID:     userID,
		PlatformID: platformID,
	}
	ctx = ccontext.WithInfo(ctx, &ccontext.GlobalConfig{
		UserID:   userID,
		Token:    adminToken,
		IMConfig: &imConf,
	})
	return api.ExtractField(ctx, api.GetUsersToken.Invoke, userReq, (*auth.GetUserTokenResp).GetToken)
}
