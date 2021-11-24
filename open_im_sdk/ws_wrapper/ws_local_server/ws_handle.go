package ws_local_server

import (
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type EventData struct {
	Event       string `json:"event"`
	ErrCode     int    `json:"errCode"`
	ErrMsg      string `json:"errMsg"`
	Data        string `json:"data"`
	OperationID string `json:"operationID"`
}

type BaseSuccFailed struct {
	funcName    string //e.g open_im_sdk/open_im_sdk.Login
	operationID string
	uid         string
}

//e.g open_im_sdk/open_im_sdk.Login ->Login
func cleanUpfuncName(funcName string) string {
	end := strings.LastIndex(funcName, ".")
	if end == -1 {
		wrapSdkLog("funcName not include.", funcName)
		return ""
	}
	return funcName[end+1:]
}

func (b *BaseSuccFailed) OnError(errCode int, errMsg string) {
	wrapSdkLog("!!!!!!!OnError ", b.uid, b.operationID, b.funcName)
	SendOneUserMessage(EventData{cleanUpfuncName(b.funcName), errCode, errMsg, "", b.operationID}, b.uid)
}

func (b *BaseSuccFailed) OnSuccess(data string) {
	wrapSdkLog("!!!!!!!OnSuccess ", b.uid, b.operationID, b.funcName)
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

//uid->funcname->func

type WsFuncRouter struct {
	uId string
}

func DelUserRouter(uid string) {
	wrapSdkLog("DelUserRouter ", uid)

	idx := strings.LastIndex(uid, " Web")
	if idx == -1 {
		wrapSdkLog("err uid, not Web", uid)
		return
	}

	uid = uid[:idx]

	UserRouteRwLock.Lock()
	defer UserRouteRwLock.Unlock()
	urm, ok := UserRouteMap[uid]
	if ok {
		wrapSdkLog("DelUserRouter logout, uninitsdk", uid)
		urm.wsRouter.LogoutNoCallback(uid, "0")
		urm.wsRouter.UnInitSDK()
	} else {
		wrapSdkLog("no found UserRouteMap: ", uid)
	}
	wrapSdkLog("DelUserRouter delete ", uid)
	delete(UserRouteMap, uid)
}

func GenUserRouterNoLock(uid string) *RefRouter {
	_, ok := UserRouteMap[uid]
	if ok {
		return nil
	}
	RouteMap1 := make(map[string]reflect.Value, 0)
	var wsRouter1 WsFuncRouter
	wsRouter1.uId = uid
	wsRouter1.AddAdvancedMsgListener()
	wsRouter1.SetConversationListener()
	wsRouter1.SetFriendListener()
	wsRouter1.SetGroupListener()
	vf := reflect.ValueOf(&wsRouter1)
	vft := vf.Type()

	mNum := vf.NumMethod()
	for i := 0; i < mNum; i++ {
		mName := vft.Method(i).Name
		wrapSdkLog("index:", i, " MethodName:", mName)
		RouteMap1[mName] = vf.Method(i)
	}
	wsRouter1.InitSDK(ConfigSvr, "0")
	var rr RefRouter
	rr.refName = &RouteMap1
	rr.wsRouter = &wsRouter1
	UserRouteMap[uid] = rr
	wrapSdkLog("insert UserRouteMap: ", uid)
	return &rr
}

func GenUserRouter(uid string) *map[string]reflect.Value {
	UserRouteRwLock.Lock()
	defer UserRouteRwLock.Unlock()

	_, ok := UserRouteMap[uid]
	if ok {
		return nil
	}
	RouteMap1 := make(map[string]reflect.Value, 0)
	var wsRouter1 WsFuncRouter
	wsRouter1.uId = uid
	wsRouter1.AddAdvancedMsgListener()
	wsRouter1.SetConversationListener()
	wsRouter1.SetFriendListener()
	wsRouter1.SetGroupListener()
	vf := reflect.ValueOf(&wsRouter1)
	vft := vf.Type()

	mNum := vf.NumMethod()
	for i := 0; i < mNum; i++ {
		mName := vft.Method(i).Name
		wrapSdkLog("index:", i, " MethodName:", mName)
		RouteMap1[mName] = vf.Method(i)
	}
	wsRouter1.InitSDK(ConfigSvr, "0")
	var rr RefRouter
	rr.refName = &RouteMap1
	rr.wsRouter = &wsRouter1
	UserRouteMap[uid] = rr
	wrapSdkLog("insert UserRouteMap: ", uid)
	return &RouteMap1
}
