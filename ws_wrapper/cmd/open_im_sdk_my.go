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
	"runtime"
	"sync"
)

func main() {
	var sdkWsPort, openIMApiPort, openIMWsPort *int
	//var openIMWsAddress, openIMApiAddress *string
	APIADDR := "http://43.155.69.205:10001"
	WSADDR := "ws://43.155.69.205:10002"
	sdkWsPort = flag.Int("sdk_ws_port", 30000, "openIM ws listening port")
	//openIMApiAddress = flag.String("openIMApiAddress", "", "openIM api listening port")
	//openIMWsAddress = flag.String("openIMWsAddress", "", "openIM ws listening port")
	flag.Parse()
	sysType := runtime.GOOS
	fmt.Println(sysType)
	switch sysType {
	case "darwin":
		ws_local_server.InitServer(&sdk_struct.IMConfig{ApiAddr: APIADDR,
			WsAddr: WSADDR, Platform: utils.WebPlatformID, DataDir: "./", LogLevel: 6, ObjectStorage: "minio"})
	case "linux":
		//sdkDBDir:= flag.String("sdk_db_dir","","openIMSDK initialization path")
		ws_local_server.InitServer(&sdk_struct.IMConfig{ApiAddr: "http://" + utils.ServerIP + ":" + utils.IntToString(*openIMApiPort),
			WsAddr: "ws://" + utils.ServerIP + ":" + utils.IntToString(*openIMWsPort), Platform: utils.WebPlatformID, DataDir: "../db/sdk/"})

	case "windows":
		ws_local_server.InitServer(&sdk_struct.IMConfig{ApiAddr: APIADDR,
			WsAddr: WSADDR, Platform: utils.WindowsPlatformID, DataDir: "./", LogLevel: 6})
	default:
		fmt.Println("this os not support", sysType)

	}
	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Println("ws server is starting")
	ws_local_server.WS.OnInit(*sdkWsPort)
	//go funcation() {
	//	log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	//}()

	ws_local_server.WS.Run()

	fmt.Println("run success")
	wg.Wait()

}
