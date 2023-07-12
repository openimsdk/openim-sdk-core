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
	"errors"
	"github.com/gorilla/websocket"
	"open_im_sdk/pkg/log"
	"open_im_sdk/sdk_struct"
	"open_im_sdk/ws_wrapper/utils"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type EventData struct {
	Event       string `json:"event"`
	ErrCode     int32  `json:"errCode"`
	ErrMsg      string `json:"errMsg"`
	Data        string `json:"data"`
	OperationID string `json:"operationID"`
}

type BaseSuccessFailed struct {
	funcName    string //e.g open_im_sdk/open_im_sdk.Login
	operationID string
	uid         string
}

// e.g open_im_sdk/open_im_sdk.Login ->Login
func cleanUpfuncName(funcName string) string {
	end := strings.LastIndex(funcName, ".")
	if end == -1 {
		log.Info("", "funcName not include.", funcName)
		return ""
	}
	return funcName[end+1:]
}

func (b *BaseSuccessFailed) OnError(errCode int32, errMsg string) {
	log.Info("", "!!!!!!!OnError ", b.uid, b.operationID, b.funcName)
	SendOneUserMessage(EventData{cleanUpfuncName(b.funcName), errCode, errMsg, "", b.operationID}, b.uid)
}

func (b *BaseSuccessFailed) OnSuccess(data string) {
	log.Info("", "!!!!!!!OnSuccess ", b.uid, b.operationID, b.funcName)
	SendOneUserMessage(EventData{cleanUpfuncName(b.funcName), 0, "", data, b.operationID}, b.uid)
}

func runFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

func int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

//uid->funcname->funcation

type WsFuncRouter struct {
	uId string
	//conn *UserConn
}

func DelUserRouter(uid string, operationID string) {
	log.Info(operationID, "DelUserRouter ", uid)
	sub := " " + utils.PlatformIDToName(sdk_struct.SvrConf.Platform)
	idx := strings.LastIndex(uid, sub)
	if idx == -1 {
		log.Info(operationID, "err uid, not Web", uid, sub)
		return
	}

	uid = uid[:idx]

	UserRouteRwLock.Lock()
	defer UserRouteRwLock.Unlock()
	urm, ok := UserRouteMap[uid]
	if ok {
		log.Info(operationID, "DelUserRouter logout, UnInitSDK ", uid, operationID)
		urm.wsRouter.LogoutNoCallback(uid, operationID)
		urm.wsRouter.UnInitSDK()
	} else {
		log.Info(operationID, "no found UserRouteMap: ", uid)
	}
	log.Info(operationID, "DelUserRouter delete ", uid)
	t, ok := UserRouteMap[uid]
	if ok {
		t.refName = make(map[string]reflect.Value)
	}

	delete(UserRouteMap, uid)
}

func GenUserRouterNoLock(uid string, batchMsg int, operationID string) *RefRouter {
	_, ok := UserRouteMap[uid]
	if ok {
		return nil
	}
	RouteMap1 := make(map[string]reflect.Value, 0)
	var wsRouter1 WsFuncRouter
	wsRouter1.uId = uid

	vf := reflect.ValueOf(&wsRouter1)
	vft := vf.Type()

	mNum := vf.NumMethod()
	for i := 0; i < mNum; i++ {
		mName := vft.Method(i).Name
		log.Info(operationID, "index:", i, " MethodName:", mName)
		RouteMap1[mName] = vf.Method(i)
	}
	wsRouter1.InitSDK(ConfigSvr, operationID)
	log.Info(operationID, "SetAdvancedMsgListener() ", uid)
	wsRouter1.SetAdvancedMsgListener()
	if batchMsg == 1 {
		log.Info(operationID, "SetBatchMsgListener() ", uid)
		wsRouter1.SetBatchMsgListener()
	}
	wsRouter1.SetConversationListener()
	log.Info(operationID, "SetFriendListener() ", uid)
	wsRouter1.SetFriendListener()
	log.Info(operationID, "SetGroupListener() ", uid)
	wsRouter1.SetGroupListener()
	log.Info(operationID, "SetUserListener() ", uid)
	wsRouter1.SetUserListener()
	log.Info(operationID, "SetSignalingListener() ", uid)
	wsRouter1.SetSignalingListener()
	log.Info(operationID, "setWorkMomentsListener()", uid)
	wsRouter1.SetWorkMomentsListener()
	log.Info(operationID, "SetOrganizationListener()", uid)
	wsRouter1.SetOrganizationListener()
	var rr RefRouter
	rr.refName = RouteMap1
	rr.wsRouter = &wsRouter1
	UserRouteMap[uid] = rr
	log.Info(operationID, "insert UserRouteMap: ", uid)
	return &rr
}

func (wsRouter *WsFuncRouter) GlobalSendMessage(data interface{}) {
	SendOneUserMessage(data, wsRouter.uId)
}

// listener
func SendOneUserMessage(data interface{}, uid string) {
	var chMsg ChanMsg
	chMsg.data, _ = json.Marshal(data)
	chMsg.uid = uid
	err := send2Ch(WS.ch, &chMsg, 2)
	if err != nil {
		log.Info("", "send2ch failed, ", err, string(chMsg.data), uid)
		return
	}
	log.Info("", "send response to web: ", string(chMsg.data))
}

func SendOneUserMessageForTest(data interface{}, uid string) {
	d, err := json.Marshal(data)
	log.Info("", "Marshal ", string(d))
	var chMsg ChanMsg
	chMsg.data = d
	chMsg.uid = uid
	err = send2ChForTest(WS.ch, chMsg, 2)
	if err != nil {
		log.Info("", "send2ch failed, ", err, string(chMsg.data), uid)
		return
	}
	log.Info("", "send response to web: ", string(chMsg.data))
}

func SendOneConnMessage(data interface{}, conn *UserConn) {
	bMsg, _ := json.Marshal(data)
	err := WS.writeMsg(conn, websocket.TextMessage, bMsg)
	log.Info("", "send response to web: ", string(bMsg), "userUid", WS.getUserUid(conn))
	if err != nil {
		log.Info("", "WS WriteMsg error", "", "userIP", conn.RemoteAddr().String(), "userUid", WS.getUserUid(conn), "error", err, "data", data)
	} else {
		log.Info("", "WS WriteMsg ok", "data", data, "userUid", WS.getUserUid(conn))
	}
}

func send2ChForTest(ch chan ChanMsg, value ChanMsg, timeout int64) error {
	var t ChanMsg
	t = value
	log.Info("", "test uid ", t.uid)
	return nil
}

func send2Ch(ch chan ChanMsg, value *ChanMsg, timeout int64) error {
	var flag = 0
	select {
	case ch <- *value:
		flag = 1
	case <-time.After(time.Second * time.Duration(timeout)):
		flag = 2
	}
	if flag == 1 {
		return nil
	} else {
		log.Info("", "send cmd timeout, ", timeout, value)
		return errors.New("send cmd timeout")
	}
}
