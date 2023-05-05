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

	//	_ "net/http/pprof"
	_ "net/http/pprof"
	"open_im_sdk/sdk_struct"

	//"open_im_sdk/open_im_sdk"

	"open_im_sdk/ws_wrapper/utils"
	"open_im_sdk/ws_wrapper/ws_local_server"
	"runtime"
	"sync"
)

func main() {
	var sdkWsPort *int
	var openIMWsAddress, openIMApiAddress, openIMDbDir, objectStorage, encryptionKey *string

	sdkWsPort = flag.Int("sdkWsPort", 7788, "openIMSDK ws listening port")
	openIMApiAddress = flag.String("openIMApiAddress", "", "openIM api listening port")
	openIMWsAddress = flag.String("openIMWsAddress", "", "openIM ws listening port")
	openIMDbDir = flag.String("openIMDbDir", "./", "openIM db dir")
	objectStorage = flag.String("objectStorage", "cos", "openIM objectStorage")
	encryptionKey = flag.String("encryptionKey", "", "openIM encryptionKey")
	flag.Parse()

	sysType := runtime.GOOS
	open_im_sdk.SetHeartbeatInterval(5)
	switch sysType {

	case "darwin":
		ws_local_server.InitServer(&sdk_struct.IMConfig{ApiAddr: *openIMApiAddress,
			WsAddr: *openIMWsAddress, Platform: utils.OSXPlatformID, DataDir: *openIMDbDir, ObjectStorage: *objectStorage, EncryptionKey: *encryptionKey})

	case "linux":
		ws_local_server.InitServer(&sdk_struct.IMConfig{ApiAddr: *openIMApiAddress,
			WsAddr: *openIMWsAddress, Platform: utils.LinuxPlatformID, DataDir: *openIMDbDir, ObjectStorage: *objectStorage, EncryptionKey: *encryptionKey})
	case "windows":
		ws_local_server.InitServer(&sdk_struct.IMConfig{ApiAddr: *openIMApiAddress,
			WsAddr: *openIMWsAddress, Platform: utils.WindowsPlatformID, DataDir: *openIMDbDir, ObjectStorage: *objectStorage, EncryptionKey: *encryptionKey})
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
