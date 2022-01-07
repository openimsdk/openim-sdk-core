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
	"open_im_sdk/internal/open_im_sdk"
	utils2 "open_im_sdk/pkg/utils"
	"open_im_sdk/ws_wrapper/utils"
	"open_im_sdk/ws_wrapper/ws_local_server"
	"runtime"
	"sync"
)

func main() {
	var sdkWsPort, openIMApiPort, openIMWsPort *int
	var openIMWsAddress, openIMApiAddress *string
	APIADDR := "http://121.37.25.71:10000"
	WSADDR := "ws://121.37.25.71:17778"
	openIMApiAddress = flag.String("openIMApiAddress", "", "openIM api listening port")
	openIMWsAddress = flag.String("openIMWsAddress", "", "openIM ws listening port")
	flag.Parse()
	open_im_sdk.SetHearbeatInterval(5)
	sysType := runtime.GOOS
	switch sysType {
	case "darwin":
		ws_local_server.InitServer(&utils2.IMConfig{IpApiAddr: *openIMApiAddress,
			IpWsAddr: *openIMWsAddress, Platform: utils.OSXPlatformID, DbDir: "./"})
	case "linux":

		//sdkDBDir:= flag.String("sdk_db_dir","","openIMSDK initialization path")
		ws_local_server.InitServer(&utils2.IMConfig{IpApiAddr: "http://" + utils.ServerIP + ":" + utils.IntToString(*openIMApiPort),
			IpWsAddr: "ws://" + utils.ServerIP + ":" + utils.IntToString(*openIMWsPort), Platform: utils.WebPlatformID, DbDir: "../db/sdk/"})

	case "windows":
		sdkWsPort = flag.Int("sdk_ws_port", 30000, "openIM ws listening port")
		flag.Parse()
		ws_local_server.InitServer(&utils2.IMConfig{IpApiAddr: APIADDR,
			IpWsAddr: WSADDR, Platform: utils.WebPlatformID, DbDir: "./"})
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
