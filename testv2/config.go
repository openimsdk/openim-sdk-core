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
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/protocol/constant"
)

const (
	APIADDR = "http://172.16.8.48:10002"
	WSADDR  = "ws://172.16.8.48:10001"

	//APIADDR = "http://127.0.0.1:10002"
	//WSADDR  = "ws://127.0.0.1:10001"

	//APIADDR = "http://127.0.0.1:10002"
	//WSADDR  = "ws://127.0.0.1:10001"

	UserID       = "7327731536"
	friendUserID = "3281432310"
)

func getConf(APIADDR, WSADDR string) sdk_struct.IMConfig {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = APIADDR
	cf.WsAddr = WSADDR
	cf.DataDir = "../"
	cf.LogLevel = 6
	cf.IsExternalExtensions = true
	cf.PlatformID = constant.LinuxPlatformID
	cf.LogFilePath = ""
	cf.IsLogStandardOutput = true
	return cf
}
