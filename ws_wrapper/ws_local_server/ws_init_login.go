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

package ws_local_server

import (
	"encoding/json"
	"open_im_sdk/internal/login"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"

	//	"open_im_sdk/pkg/constant"
	//	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
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
	log.Info(operationID, "Initsdk uid: ", initcb.uid, config)
	c := sdk_struct.IMConfig{}
	json.Unmarshal([]byte(config), &c)
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if userWorker.InitSDK(c, &initcb) {
		//	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", "", operationID})
	} else {
		//	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), open_im_sdk.ErrCodeInitLogin, "init config failed", "", operationID})
	}
}

func (wsRouter *WsFuncRouter) UnInitSDK() {
	log.Info("", "UnInitSDK uid: ", wsRouter.uId)
	open_im_sdk.UserSDKRwLock.Lock()
	_, ok := open_im_sdk.UserRouterMap[wsRouter.uId]
	if ok {
		log.Info("", "delete UserRouterMap uid: ", wsRouter.uId)
		//*v = login.LoginMgr{}
	}
	delete(open_im_sdk.UserRouterMap, wsRouter.uId)
	log.Info("", "delete UnInitSDK uid: ", wsRouter.uId)
	open_im_sdk.UserSDKRwLock.Unlock()
}

func (wsRouter *WsFuncRouter) checkResourceLoadingAndKeysIn(mgr *login.LoginMgr, input, operationID, funcName string, m map[string]interface{}, keys ...string) bool {
	for _, k := range keys {
		_, ok := m[k]
		if !ok {
			log.Info(operationID, "key not in", keys, input, operationID, funcName)
			wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(funcName), StatusBadParameter, "key not in", "", operationID})
			return false
		}
	}

	if err := open_im_sdk.CheckResourceLoad(mgr); err != nil {
		log.Info(operationID, "Resource Loading ", mgr, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(funcName), StatusResourceNotCompleted, "resource loading is not completed", "", operationID})
		return false
	}
	return true
}

func (wsRouter *WsFuncRouter) checkKeysIn(input, operationID, funcName string, m map[string]interface{}, keys ...string) bool {
	for _, k := range keys {
		_, ok := m[k]
		if !ok {
			log.Info(operationID, "key not in", keys, input, funcName, m)
			wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(funcName), StatusBadParameter, "key not in", "", operationID})
			return false
		}
	}
	return true
}

func (wsRouter *WsFuncRouter) Login(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "userID", "token") {
		return
	}
	userWorker.Login(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["userID"].(string), m["token"].(string), operationID)
}

func (wsRouter *WsFuncRouter) Logout(input string, operationID string) {
	//userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//userWorker.Logout(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId})
	//todo just send response
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", "", operationID})
}

func (wsRouter *WsFuncRouter) LogoutNoCallback(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Logout(nil, operationID)
}

func (wsRouter *WsFuncRouter) GetLoginStatus(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", int32ToString(int32(userWorker.GetLoginStatus())), operationID})
}

// 1
func (wsRouter *WsFuncRouter) getMyLoginStatus() int32 {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, "", "", runFuncName(), nil) {
		return constant.SdkInit
	}
	return userWorker.GetLoginStatus()
}

// 1
func (wsRouter *WsFuncRouter) GetLoginUser(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", userWorker.GetLoginUser(), operationID})
}

func InitServer(config *sdk_struct.IMConfig) {
	data, _ := json.Marshal(config)
	ConfigSvr = string(data)
	UserRouteMap = make(map[string]RefRouter, 0)
	open_im_sdk.InitOnce(config)
	log.Info("", "InitServer ", ConfigSvr)
}
