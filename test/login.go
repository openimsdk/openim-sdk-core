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
	"fmt"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"time"
)

type BaseSuccessFailed struct {
	successData string
	errCode     int
	errMsg      string
	funcName    string
	time        time.Time
}

func (b *BaseSuccessFailed) OnError(errCode int32, errMsg string) {
	b.errCode = -1
	b.errMsg = errMsg
	log.Error("login failed", errCode, errMsg)

}

func (b *BaseSuccessFailed) OnSuccess(data string) {
	b.errCode = 1
	b.successData = data
	log.Info("login success", data, time.Since(b.time))
}

func InOutDoTest(uid, tk, ws, api string) {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = api
	cf.PlatformID = constant.WindowsPlatformID
	cf.WsAddr = ws
	cf.DataDir = "./"
	cf.LogLevel = LogLevel
	cf.IsExternalExtensions = true
	cf.IsLogStandardOutput = true
	cf.LogFilePath = "./"

	b, _ := json.Marshal(cf)
	s := string(b)
	fmt.Println(s)
	var testinit testInitLister

	operationID := utils.OperationIDGenerator()
	if !open_im_sdk.InitSDK(&testinit, operationID, s) {
		fmt.Println("", "InitSDK failed")
		return
	}

	var testConversation conversationCallBack
	open_im_sdk.SetConversationListener(&testConversation)

	var testUser userCallback
	open_im_sdk.SetUserListener(testUser)

	var msgCallBack MsgListenerCallBak
	open_im_sdk.SetAdvancedMsgListener(&msgCallBack)

	var batchMsg BatchMsg
	open_im_sdk.SetBatchMsgListener(&batchMsg)

	var friendListener testFriendListener
	open_im_sdk.SetFriendListener(friendListener)

	var groupListener testGroupListener
	open_im_sdk.SetGroupListener(groupListener)

	InOutlllogin(uid, tk)
}

func InOutlllogin(uid, tk string) {
	var callback BaseSuccessFailed
	callback.time = time.Now()
	callback.funcName = utils.GetSelfFuncName()
	operationID := utils.OperationIDGenerator()
	open_im_sdk.Login(&callback, operationID, uid, tk)
	for {
		if callback.errCode == 1 {
			return
		} else if callback.errCode == -1 {
			time.Sleep(100 * time.Millisecond)
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func InOutLogout() {
	var callback BaseSuccessFailed
	callback.funcName = utils.GetSelfFuncName()
	opretaionID := utils.OperationIDGenerator()
	open_im_sdk.Logout(&callback, opretaionID)
}
