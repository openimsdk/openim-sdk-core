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

package main

import (
	"flag"
	"fmt"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"

	//	_ "net/http/pprof"
	"net/http"
	_ "net/http/pprof"
	"open_im_sdk/sdk_struct"

	//"open_im_sdk/open_im_sdk"

	log1 "log"
	"open_im_sdk/ws_wrapper/utils"
	"open_im_sdk/ws_wrapper/ws_local_server"
	"runtime"
	"sync"
)

func main() {
	go func() {

		log1.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()
	var sdkWsPort, logLevel *int
	var openIMWsAddress, openIMApiAddress, openIMDbDir *string

	// 	log1.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	// }()
	// var sdkWsPort, logLevel *int
	// var openIMWsAddress, openIMDbDir, openIMApiAddress *string

	openIMApiAddress = flag.String("openIM_api_address", "http://127.0.0.1:10002", "openIM api listening address")
	openIMWsAddress = flag.String("openIM_ws_address", "ws://127.0.0.1:10001", "openIM ws listening address")
	sdkWsPort = flag.Int("sdk_ws_port", 10003, "openIMSDK ws listening port")
	logLevel = flag.Int("openIM_log_level", 6, "control log output level")
	openIMDbDir = flag.String("openIMDbDir", "../db/sdk/", "openIM db dir")
	flag.Parse()
	fmt.Println("sdk server init args is :", "apiAddress:", *openIMApiAddress, "wsAddress:", *openIMWsAddress, *sdkWsPort, *logLevel)
	log.NewPrivateLog(constant.LogFileName, uint32(*logLevel))

	sysType := runtime.GOOS
	open_im_sdk.SetHeartbeatInterval(5)
	switch sysType {
	case "darwin":
		fallthrough
	case "linux":
		fallthrough
	case "windows":
		ws_local_server.InitServer(&sdk_struct.IMConfig{ApiAddr: *openIMApiAddress,
			WsAddr: *openIMWsAddress, Platform: utils.WebPlatformID, DataDir: *openIMDbDir, LogLevel: uint32(*logLevel)})
	default:
		fmt.Println("this os not support", sysType)

	}
	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Println("ws server is starting")
	ws_local_server.WS.OnInit(*sdkWsPort)
	ws_local_server.WS.Run()
	wg.Wait()

}
