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
	"open_im_sdk/pkg/constant"
	"sync"
)

var LogLevel uint32 = 6
var PlatformID = int32(3)
var LogName = ""
var IsLogStandardOutput = true
var LogFilePath = ""

var ReliabilityUserA = 1234567
var ReliabilityUserB = 1234567
var (
	//TESTIP = "121.5.182.23"
	TESTIP = "59.36.173.89"
	//TESTIP              = "121.37.25.71"

	//TESTIP  = "open-im-test.rentsoft.cn"
	APIADDR = "http://" + TESTIP + ":10002"

	WSADDR              = "ws://" + TESTIP + ":10001"
	REGISTERADDR        = APIADDR + "/auth/user_register"
	TOKENADDR           = APIADDR + "/auth/user_token"
	SECRET              = "openIM123"
	SENDINTERVAL        = 20
	GETSELFUSERINFO     = APIADDR + "/user/get_self_user_info"
	CREATEGROUP         = APIADDR + constant.CreateGroupRouter
	ACCOUNTCHECK        = APIADDR + "/user/account_check"
	GETGROUPSINFOROUTER = APIADDR + constant.GetGroupsInfoRouter
)

var coreMgrLock sync.RWMutex

var allLoginMgr map[int]*CoreNode

var allLoginMgrtmp []*CoreNode

var userLock sync.RWMutex

var allUserID []string
var allToken []string

// var allWs []*interaction.Ws
var sendSuccessCount, sendFailedCount int
var sendSuccessLock sync.RWMutex
var sendFailedLock sync.RWMutex

var msgNumInOneClient = 0

// var Msgwg sync.WaitGroup
var sendMsgClient = 0

var MaxNumGoroutine = 100000

// 常量
var (
	RELIABILITY = "reliability_"
)
