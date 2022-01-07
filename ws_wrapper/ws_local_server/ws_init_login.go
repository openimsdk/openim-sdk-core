package ws_local_server

import (
	"encoding/json"
	"open_im_sdk/internal/controller/init"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/utils"
)

type InitCallback struct {
	uid string
}

func (i *InitCallback) OnConnecting() {
	var ed EventData
	ed.Event = cleanUpfuncName(runFuncName())
	ed.ErrCode = 0
	SendOneUserMessage(ed, i.uid)
}

func (i *InitCallback) OnConnectSuccess() {
	var ed EventData
	ed.Event = cleanUpfuncName(runFuncName())
	ed.ErrCode = 0
	SendOneUserMessage(ed, i.uid)
}

func (i *InitCallback) OnConnectFailed(ErrCode int32, ErrMsg string) {
	var ed EventData
	ed.Event = cleanUpfuncName(runFuncName())
	ed.ErrCode = ErrCode
	ed.ErrMsg = ErrMsg
	SendOneUserMessage(ed, i.uid)
}

func (i *InitCallback) OnKickedOffline() {
	var ed EventData
	ed.Event = cleanUpfuncName(runFuncName())
	ed.ErrCode = 0
	SendOneUserMessage(ed, i.uid)
}

func (i *InitCallback) OnUserTokenExpired() {
	var ed EventData
	ed.Event = cleanUpfuncName(runFuncName())
	ed.ErrCode = 0
	SendOneUserMessage(ed, i.uid)
}

func (i *InitCallback) OnSelfInfoUpdated(userInfo string) {
	var ed EventData
	ed.Data = userInfo
	ed.Event = cleanUpfuncName(runFuncName())
	ed.ErrCode = 0
	SendOneUserMessage(ed, i.uid)
}

var ConfigSvr string

func (wsRouter *WsFuncRouter) InitSDK(config string, operationID string) {
	var initcb InitCallback
	initcb.uid = wsRouter.uId
	wrapSdkLog("Initsdk uid: ", initcb.uid)
	userWorker := init.GetUserWorker(wsRouter.uId)
	if userWorker.InitSDK(config, &initcb) {
		//	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", "", operationID})
	} else {
		//	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), open_im_sdk.ErrCodeInitLogin, "init config failed", "", operationID})
	}
}

func (wsRouter *WsFuncRouter) UnInitSDK() {
	wrapSdkLog("UnInitSDK uid: ", wsRouter.uId)
	userWorker := init.GetUserWorker(wsRouter.uId)
	userWorker.UnInitSDK()
	constant.UserSDKRwLock.Lock()
	delete(constant.UserRouterMap, wsRouter.uId)
	wrapSdkLog("delete UnInitSDK uid: ", wsRouter.uId)
	constant.UserSDKRwLock.Unlock()
}

func (wsRouter *WsFuncRouter) checkKeysIn(input, operationID, funcName string, m map[string]interface{}, keys ...string) bool {
	for _, k := range keys {
		_, ok := m[k]
		if !ok {
			wrapSdkLog("key not in", keys, input, operationID, funcName)
			wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(funcName), StatusBadParameter, "key not in", "", operationID})
			return false
		}
	}
	return true
}

func (wsRouter *WsFuncRouter) Login(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed", err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := init.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "uid", "token") {
		return
	}
	userWorker.Login(m["uid"].(string), m["token"].(string), &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) Logout(input string, operationID string) {
	//userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//userWorker.Logout(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
	//todo just send response
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", "", operationID})
}

func (wsRouter *WsFuncRouter) LogoutNoCallback(input string, operationID string) {
	userWorker := init.GetUserWorker(wsRouter.uId)
	userWorker.Logout(nil)
}

func (wsRouter *WsFuncRouter) GetLoginStatus(input string, operationID string) {
	userWorker := init.GetUserWorker(wsRouter.uId)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", int32ToString(int32(userWorker.GetLoginStatus())), operationID})
}

//1
func (wsRouter *WsFuncRouter) getMyLoginStatus() int {
	userWorker := init.GetUserWorker(wsRouter.uId)
	return userWorker.GetLoginStatus()
}

//1
func (wsRouter *WsFuncRouter) GetLoginUser(input string, operationID string) {
	userWorker := init.GetUserWorker(wsRouter.uId)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", userWorker.GetLoginUser(), operationID})
}

func InitServer(config *utils.IMConfig) {
	data, _ := json.Marshal(config)
	ConfigSvr = string(data)
	UserRouteMap = make(map[string]RefRouter, 0)
	init.InitOnce(config)
}
