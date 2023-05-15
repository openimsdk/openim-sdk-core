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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"io"
	"math/rand"
	"net/http"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/ccontext"
	"time"
)

var (
	ctx context.Context
)

func init() {
	rand.Seed(time.Now().UnixNano())
	listner := &OnConnListener{}
	config := getConf(APIADDR, WSADDR)
	configData, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	isInit := open_im_sdk.InitSDK(listner, "test", string(configData))
	if !isInit {
		panic("init sdk failed")
	}
	ctx = open_im_sdk.UserForSDK.Context()
	ctx = ccontext.WithOperationID(ctx, "initOperationID")
	token, err := GetUserToken(ctx, UserID)
	if err != nil {
		panic(err)
	}
	if err := open_im_sdk.UserForSDK.Login(ctx, UserID, token); err != nil {
		panic(err)
	}
	open_im_sdk.UserForSDK.SetListenerForService(&onListenerForService{ctx: ctx})
	open_im_sdk.UserForSDK.SetConversationListener(&onConversationListener{ctx: ctx})
}

func GetUserToken(ctx context.Context, userID string) (string, error) {
	jsonReqData, err := json.Marshal(map[string]any{
		"userID":   userID,
		"platform": 1,
		"secret":   "openIM123",
	})
	if err != nil {
		return "", err
	}
	path := APIADDR + "/auth/user_token"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(jsonReqData))
	if err != nil {
		return "", err
	}
	req.Header.Set("operationID", ctx.Value("operationID").(string))
	client := http.Client{Timeout: time.Second * 3}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	type Result struct {
		ErrCode int    `json:"errCode"`
		ErrMsg  string `json:"errMsg"`
		ErrDlt  string `json:"errDlt"`
		Data    struct {
			Token             string `json:"token"`
			ExpireTimeSeconds int    `json:"expireTimeSeconds"`
		} `json:"data"`
	}
	var result Result
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if result.ErrCode != 0 {
		return "", fmt.Errorf("errCode:%d, errMsg:%s, errDlt:%s", result.ErrCode, result.ErrMsg, result.ErrDlt)
	}
	return result.Data.Token, nil
}

type onListenerForService struct {
	ctx context.Context
}

func (o *onListenerForService) OnGroupApplicationAdded(groupApplication string) {
	log.ZInfo(o.ctx, "OnGroupApplicationAdded", "groupApplication", groupApplication)
}

func (o *onListenerForService) OnGroupApplicationAccepted(groupApplication string) {
	log.ZInfo(o.ctx, "OnGroupApplicationAccepted", "groupApplication", groupApplication)
}

func (o *onListenerForService) OnFriendApplicationAdded(friendApplication string) {
	log.ZInfo(o.ctx, "OnFriendApplicationAdded", "friendApplication", friendApplication)
}

func (o *onListenerForService) OnFriendApplicationAccepted(groupApplication string) {
	log.ZInfo(o.ctx, "OnFriendApplicationAccepted", "groupApplication", groupApplication)
}

func (o *onListenerForService) OnRecvNewMessage(message string) {
	log.ZInfo(o.ctx, "OnRecvNewMessage", "message", message)
}

type onConversationListener struct {
	ctx context.Context
}

func (o *onConversationListener) OnSyncServerStart() {
	log.ZInfo(o.ctx, "OnSyncServerStart")
}

func (o *onConversationListener) OnSyncServerFinish() {
	log.ZInfo(o.ctx, "OnSyncServerFinish")
}

func (o *onConversationListener) OnSyncServerFailed() {
	log.ZInfo(o.ctx, "OnSyncServerFailed")
}

func (o *onConversationListener) OnNewConversation(conversationList string) {
	log.ZInfo(o.ctx, "OnNewConversation", "conversationList", conversationList)
}

func (o *onConversationListener) OnConversationChanged(conversationList string) {
	log.ZInfo(o.ctx, "OnConversationChanged", "conversationList", conversationList)
}

func (o *onConversationListener) OnTotalUnreadMessageCountChanged(totalUnreadCount int32) {
	log.ZInfo(o.ctx, "OnTotalUnreadMessageCountChanged", "totalUnreadCount", totalUnreadCount)
}
