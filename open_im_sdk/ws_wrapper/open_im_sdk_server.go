/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/9/8 14:35).
 */
package main

import (
	"flag"
	"fmt"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/open_im_sdk/ws_wrapper/utils"
	"open_im_sdk/open_im_sdk/ws_wrapper/ws_local_server"
	"runtime"
	"sync"
)

func main() {
	var sdkWsPort *int
	sysType := runtime.GOOS
	switch sysType {
	case "linux":
		openIMApiPort := flag.Int("openIM_api_port", 0, "openIM api listening port")
		openIMWsPort := flag.Int("openIM_ws_port", 0, "openIM ws listening port")
		sdkWsPort = flag.Int("sdk_ws_port", 7799, "openIMSDK ws listening port")
		//sdkDBDir:= flag.String("sdk_db_dir","","openIMSDK initialization path")
		flag.Parse()
		ws_local_server.InitServer(&open_im_sdk.IMConfig{IpApiAddr: "http://" + utils.ServerIP + ":" + utils.IntToString(*openIMApiPort),
			IpWsAddr: "ws://" + utils.ServerIP + ":" + utils.IntToString(*openIMWsPort), Platform: utils.WebPlatformID, DbDir: "../db/sdk/"})

	case "windows":
		sdkWsPort = flag.Int("sdk_ws_port", 7799, "openIM ws listening port")
		flag.Parse()
		ws_local_server.InitServer(&open_im_sdk.IMConfig{IpApiAddr: "https://open-im.rentsoft.cn",
			IpWsAddr: "wss://open-im.rentsoft.cn/wss", Platform: utils.WebPlatformID, DbDir: "./"})
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
