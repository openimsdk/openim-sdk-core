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

	"open_im_sdk/sdk_struct"
	"open_im_sdk/ws_wrapper/utils"
	"open_im_sdk/ws_wrapper/ws_local_server"
	"os"
	"runtime"
	"sync"
)

func main() {
	APIADDR := "http://172.20.10.5:10000"
	WSADDR := "ws://172.20.10.5:17778"

	sqliteDir := flag.String("sdk_db_dir", "./", "openIMSDK initialization path")
	sdkWsIp := flag.String("sdk_ws_ip", "", "openIM ws listening ip")
	sdkWsPort := flag.Int("sdk_ws_port", 7799, "openIM ws listening port")
	openIMApiAddress := flag.String("api_address", APIADDR, "openIM api address")
	openIMWsAddress := flag.String("ws_address", WSADDR, "openIM ws address")

	flag.Parse()
	sysType := runtime.GOOS
	switch sysType {
	case "darwin":

		ws_local_server.InitServer(&sdk_struct.IMConfig{ApiAddr: *openIMApiAddress,
			WsAddr: *openIMWsAddress, Platform: utils.OSXPlatformID, DataDir: *sqliteDir})
	case "linux":

		ws_local_server.InitServer(&sdk_struct.IMConfig{ApiAddr: *openIMApiAddress,
			WsAddr: *openIMWsAddress, Platform: utils.LinuxPlatformID, DataDir: *sqliteDir})
	case "windows":

		ws_local_server.InitServer(&sdk_struct.IMConfig{ApiAddr: *openIMApiAddress,
			WsAddr: *openIMWsAddress, Platform: utils.WindowsPlatformID, DataDir: *sqliteDir})
	default:
		fmt.Println("this os not support", sysType)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Println("ws server is starting")

	ws_local_server.WS.OnInit(*sdkWsPort, *sdkWsIp)
	if ws_local_server.WS.Run()!=nil {
		os.Exit(-10);
	}
	fmt.Println("run success")
	wg.Wait()

}
