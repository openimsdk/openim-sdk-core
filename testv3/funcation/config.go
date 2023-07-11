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

package funcation

import (
	"context"
	"fmt"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"sync"

	imLog "github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
)

var cf = sdk_struct.IMConfig{
	ApiAddr:             APIADDR,
	WsAddr:              WSADDR,
	PlatformID:          PlatformID,
	DataDir:             "./../",
	LogLevel:            LogLevel,
	IsLogStandardOutput: true,
}

// log and token
var (
	rotateCount         = uint(0)
	LogLevel            = uint32(6)
	PlatformID          = int32(1)
	Secret              = "tuoyun"
	IsLogStandardOutput = true
	isLogJson           = false
	LogName             = ""
	LogFilePath         = ""
	DataDir             = "./../"

	AdminToken = ""
)

// ctx and it's config
var (
	ctx    context.Context
	config ccontext.GlobalConfig
)

func init() {
	AdminToken, _ = getToken("openIM123456")
	// init log
	if err := imLog.InitFromConfig(
		"open-im-sdk-core", LogName, int(LogLevel), IsLogStandardOutput, isLogJson, LogFilePath, rotateCount); err != nil {
		fmt.Println(utils.OperationIDGenerator(), "log init failed ", err.Error())
	}
}

// system
var (
	// TESTIP       = "59.36.173.89"
	TESTIP       = "203.56.175.233"
	APIADDR      = "http://" + TESTIP + ":10002"
	WSADDR       = "ws://" + TESTIP + ":10001"
	REGISTERADDR = APIADDR + "/auth/user_register"
	TOKENADDR    = APIADDR + "/auth/user_token"
)

var coreMgrLock sync.RWMutex
var AllLoginMgr map[string]*CoreNode
var userLock sync.RWMutex
var AllUserID []string

// var allWs []*interaction.Ws
var sendSuccessCount, sendFailedCount int
var sendSuccessLock sync.RWMutex
var sendFailedLock sync.RWMutex
var SendMsgMapLock sync.RWMutex

var msgNumInOneClient = 0
var MaxNumGoroutine = 100000

// var Msgwg sync.WaitGroup
var sendMsgClient = 0

// Listener
var (
	testConversation conversationCallBack
)
