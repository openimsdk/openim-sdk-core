package main

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/wasm_wrapper"

	"open_im_sdk/wasm/indexdb"
	"syscall/js"
)

//type BaseSuccessFailed struct {
//	funcName    string //e.g open_im_sdk/open_im_sdk.Login
//	operationID string
//	uid         string
//}
//
//func (b *BaseSuccessFailed) OnError(errCode int32, errMsg string) {
//	fmt.Println("OnError", errCode, errMsg)
//
//}
//
//func (b *BaseSuccessFailed) OnSuccess(data string) {
//	fmt.Println("OnError", data)
//}
//
//type InitCallback struct {
//	uid string
//}
//
//func (i *InitCallback) OnConnecting() {
//	fmt.Println("OnConnecting")
//}
//
//func (i *InitCallback) OnConnectSuccess() {
//	fmt.Println("OnConnecting")
//
//}
//
//func (i *InitCallback) OnConnectFailed(ErrCode int32, ErrMsg string) {
//	fmt.Println("OnConnecting", ErrCode, ErrMsg)
//
//}
//
//func (i *InitCallback) OnKickedOffline() {
//	fmt.Println("OnConnecting")
//
//}
//
//func (i *InitCallback) OnUserTokenExpired() {
//	fmt.Println("OnConnecting")
//
//}
//
//func (i *InitCallback) OnSelfInfoUpdated(userInfo string) {
//	fmt.Println("OnSelfInfoUpdated", userInfo)
//
//}
//
//var (
//	TESTIP = "43.155.69.205"
//	//TESTIP          = "121.37.25.71"
//	APIADDR = "http://" + TESTIP + ":10002"
//	WSADDR  = "ws://" + TESTIP + ":10001"
//)

func main() {
	//config := sdk_struct.IMConfig{
	//	Platform: 1,
	//	ApiAddr:  APIADDR,
	//	WsAddr:   WSADDR,
	//	DataDir:  "./",
	//	LogLevel: 6,
	//}
	//var listener InitCallback
	//var base BaseSuccessFailed
	//operationID := utils.OperationIDGenerator()
	//open_im_sdk.InitSDK(&listener, operationID, utils.StructToJsonString(config))
	//strMyUidx := "3984071717"
	//tokenx := test.RunGetToken(strMyUidx)
	//open_im_sdk.Login(&base, operationID, strMyUidx, tokenx)
	opid := utils.OperationIDGenerator()
	db := indexdb.NewIndexDB()
	msg, err := db.GetMessage("client_msg_id_123")
	if err != nil {
		log.Error(opid, "get message err:", err.Error())
	} else {
		log.Info(opid, "get message is :", *msg, "get args is :", msg.ClientMsgID, msg.ServerMsgID)
	}
	var user model_struct.LocalUser
	user.UserID = "111"
	user.CreateTime = 1232
	err = db.InsertLoginUser(&user)
	if err != nil {
		log.Error(opid, "InsertLoginUser:", err.Error())
	} else {
		log.Info(opid, "InsertLoginUser success:")
	}
	err = db.UpdateLoginUserByMap(&user, map[string]interface{}{"1": 3})
	if err != nil {
		log.Error(opid, "UpdateLoginUserByMap:", err.Error())
	} else {
		log.Info(opid, "UpdateLoginUserByMap success:")
	}
	seq, err := db.GetNormalMsgSeq()
	if err != nil {
		log.Error(opid, "GetNormalMsgSeq:", err.Error())
	} else {
		log.Info(opid, "GetNormalMsgSeq seq  success:", seq)
	}

	registerFunc()
	<-make(chan bool)
}

func registerFunc() {
	js.Global().Set(wasm_wrapper.COMMONEVENTFUNC, js.FuncOf(wasm_wrapper.CommonEventFunc))
	js.Global().Set("initSDK", js.FuncOf(wasm_wrapper.InitSDK))
	js.Global().Set("login", js.FuncOf(wasm_wrapper.Login))

}
