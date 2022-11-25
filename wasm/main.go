package main

import (
	"fmt"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"open_im_sdk/test"
	"syscall/js"
)

type BaseSuccessFailed struct {
	funcName    string //e.g open_im_sdk/open_im_sdk.Login
	operationID string
	uid         string
}

func (b *BaseSuccessFailed) OnError(errCode int32, errMsg string) {
	fmt.Println("OnError", errCode, errMsg)

}

func (b *BaseSuccessFailed) OnSuccess(data string) {
	fmt.Println("OnError", data)
}

type InitCallback struct {
	uid string
}

func (i *InitCallback) OnConnecting() {
	fmt.Println("OnConnecting")
}

func (i *InitCallback) OnConnectSuccess() {
	fmt.Println("OnConnecting")

}

func (i *InitCallback) OnConnectFailed(ErrCode int32, ErrMsg string) {
	fmt.Println("OnConnecting", ErrCode, ErrMsg)

}

func (i *InitCallback) OnKickedOffline() {
	fmt.Println("OnConnecting")

}

func (i *InitCallback) OnUserTokenExpired() {
	fmt.Println("OnConnecting")

}

func (i *InitCallback) OnSelfInfoUpdated(userInfo string) {
	fmt.Println("OnSelfInfoUpdated", userInfo)

}

var (
	TESTIP = "43.155.69.205"
	//TESTIP          = "121.37.25.71"
	APIADDR = "http://" + TESTIP + ":10002"
	WSADDR  = "ws://" + TESTIP + ":10001"
)

func readUserAgent() string {
	return js.Global().Get("navigator").Get("userAgent").String()
}

func alert(str string) {
	js.Global().Get("alert")
	//js.Get("alert").Invoke(str)
}
func main() {
	config := sdk_struct.IMConfig{
		Platform:      1,
		ApiAddr:       APIADDR,
		WsAddr:        WSADDR,
		DataDir:       "./",
		LogLevel:      6,
		IsCompression: true,
	}
	var listener InitCallback
	var base BaseSuccessFailed
	operationID := utils.OperationIDGenerator()
	open_im_sdk.InitSDK(&listener, operationID, utils.StructToJsonString(config))
	strMyUidx := "3984071717"
	tokenx := test.RunGetToken(strMyUidx)
	open_im_sdk.Login(&base, operationID, strMyUidx, tokenx)
	userAgent := readUserAgent()
	alert(userAgent)
	<-make(chan bool)
}
