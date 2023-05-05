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
	"encoding/json"
	"net/url"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"open_im_sdk/ws_wrapper/test/client"
	"open_im_sdk/ws_wrapper/ws_local_server"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var totalConnNum int
var lock sync.Mutex
var TotalSendMsgNum int

func StartSimulationJSClient(api, jssdkURL, userID string, i int, userIDList []string) {
	// 模拟登录 认证 ws连接初始化
	user := client.NewIMClient("", userID, api, jssdkURL, 5)
	var err error
	user.Token, err = user.GetToken()
	if err != nil {
		log.NewError("", "generate token failed", userID, api, err.Error())
		user.Token, err = user.GetToken()
		if err != nil {
			log.NewError("", "generate token failed", userID, api, err.Error())
			return
		}
	}
	v := url.Values{}
	v.Set("sendID", userID)
	v.Set("token", user.Token)
	v.Set("platformID", utils.IntToString(5))
	c, _, err := websocket.DefaultDialer.Dial(jssdkURL+"?"+v.Encode(), nil)
	if err != nil {
		log.NewInfo("", "dial:", err.Error(), "userID", userID, "i: ", i)
		return
	}
	lock.Lock()
	totalConnNum += 1
	log.NewInfo("", "connect success", userID, "total conn num", totalConnNum)
	lock.Unlock()
	user.Conn = c
	// user.WsLogout()
	user.WsLogin()
	time.Sleep(time.Second * 2)

	// 模拟登录同步
	go user.GetSelfUserInfo()
	go user.GetAllConversationList()
	go user.GetBlackList()
	go user.GetFriendList()
	go user.GetRecvFriendApplicationList()
	go user.GetRecvGroupApplicationList()
	go user.GetSendFriendApplicationList()
	go user.GetSendGroupApplicationList()

	// 模拟监听回调
	go func() {
		for {
			resp := ws_local_server.EventData{}
			_, message, err := c.ReadMessage()
			if err != nil {
				log.NewError("", "read:", err, "error an connet failed", userID, i)
				return
			}
			// log.Printf("recv: %s", message)
			_ = json.Unmarshal(message, &resp)
			if resp.Event == "CreateTextMessage" {
				msg := sdk_struct.MsgStruct{}
				_ = json.Unmarshal([]byte(resp.Data), &msg)
				type Data struct {
					RecvID          string `json:"recvID"`
					GroupID         string `json:"groupID"`
					OfflinePushInfo string `json:"offlinePushInfo"`
					Message         string `json:"message"`
				}
				offlinePushBytes, _ := json.Marshal(server_api_params.OfflinePushInfo{Title: "push offline"})
				messageBytes, _ := json.Marshal(msg)
				data := Data{RecvID: userID, OfflinePushInfo: string(offlinePushBytes), Message: string(messageBytes)}
				err = user.SendMsg(userID, data)
				//fmt.Println(msg)
				lock.Lock()
				TotalSendMsgNum += 1
				lock.Unlock()
			}
		}
	}()

	// 模拟给随机用户发消息
	go func() {
		for {
			err = user.CreateTextMessage(userID)
			if err != nil {
				log.NewError("", err, i, userID)
			}
			time.Sleep(time.Second * 2)
		}
	}()

	// 模拟获取登陆状态
	go func() {
		for {
			if err = user.GetLoginStatus(); err != nil {
				log.NewError("", err, i, userID)
			}
			time.Sleep(time.Second * 10)
		}
	}()
}
