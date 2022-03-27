package ws_local_server

import (
	"encoding/json"
	"errors"
	utils2 "open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"open_im_sdk/ws_wrapper/utils"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
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

//e.g open_im_sdk/open_im_sdk.Login ->Login
func cleanUpfuncName(funcName string) string {
	end := strings.LastIndex(funcName, ".")
	if end == -1 {
		wrapSdkLog("", "funcName not include.", funcName)
		return ""
	}
	return funcName[end+1:]
}

func (b *BaseSuccessFailed) OnError(errCode int32, errMsg string) {
	wrapSdkLog("", "!!!!!!!OnError ", b.uid, b.operationID, b.funcName)
	SendOneUserMessage(EventData{cleanUpfuncName(b.funcName), errCode, errMsg, "", b.operationID}, b.uid)
}

func (b *BaseSuccessFailed) OnSuccess(data string) {
	wrapSdkLog("", "!!!!!!!OnSuccess ", b.uid, b.operationID, b.funcName)
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
	//conn *UserConn
}

func DelUserRouter(uid string) {
	wrapSdkLog("", "DelUserRouter ", uid)
	sub := " " + utils.PlatformIDToName(sdk_struct.SvrConf.Platform)
	idx := strings.LastIndex(uid, sub)
	if idx == -1 {
		wrapSdkLog("", "err uid, not Web", uid)
		return
	}

	uid = uid[:idx]

	UserRouteRwLock.Lock()
	defer UserRouteRwLock.Unlock()
	urm, ok := UserRouteMap[uid]
	if ok {
		operationID := utils2.OperationIDGenerator()
		wrapSdkLog("", "DelUserRouter logout, UnInitSDK ", uid, operationID)

		urm.wsRouter.LogoutNoCallback(uid, operationID)
		urm.wsRouter.UnInitSDK()
	} else {
		wrapSdkLog("", "no found UserRouteMap: ", uid)
	}
	wrapSdkLog("", "DelUserRouter delete ", uid)
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
	//	wsRouter1.conn = conn

	vf := reflect.ValueOf(&wsRouter1)
	vft := vf.Type()

	mNum := vf.NumMethod()
	for i := 0; i < mNum; i++ {
		mName := vft.Method(i).Name
		wrapSdkLog("", "index:", i, " MethodName:", mName)
		RouteMap1[mName] = vf.Method(i)
	}
	wsRouter1.InitSDK(ConfigSvr, "0")
	wsRouter1.SetAdvancedMsgListener()
	wsRouter1.SetConversationListener()
	wrapSdkLog("", "SetFriendListener() ", uid)
	wsRouter1.SetFriendListener()
	wrapSdkLog("", "SetGroupListener() ", uid)
	wsRouter1.SetGroupListener()
	wrapSdkLog("", "SetUserListener() ", uid)
	wsRouter1.SetUserListener()
	wrapSdkLog("", "SetSignalingListener() ", uid)
	wsRouter1.SetSignalingListener()

	var rr RefRouter
	rr.refName = &RouteMap1
	rr.wsRouter = &wsRouter1
	UserRouteMap[uid] = rr
	wrapSdkLog("", "insert UserRouteMap: ", uid)
	return &rr
}

func (wsRouter *WsFuncRouter) GlobalSendMessage(data interface{}) {
	SendOneUserMessage(data, wsRouter.uId)
}

//listener
func SendOneUserMessage(data interface{}, uid string) {
	var chMsg ChanMsg
	chMsg.data, _ = json.Marshal(data)
	chMsg.uid = uid
	err := send2Ch(WS.ch, &chMsg, 2)
	if err != nil {
		wrapSdkLog("", "send2ch failed, ", err, string(chMsg.data), uid)
		return
	}
	wrapSdkLog("", "send response to web: ", string(chMsg.data))
}

func SendOneUserMessageForTest(data interface{}, uid string) {
	d, err := json.Marshal(data)
	wrapSdkLog("", "Marshal ", string(d))
	var chMsg ChanMsg
	chMsg.data = d
	chMsg.uid = uid
	err = send2ChForTest(WS.ch, chMsg, 2)
	if err != nil {
		wrapSdkLog("", "send2ch failed, ", err, string(chMsg.data), uid)
		return
	}
	wrapSdkLog("", "send response to web: ", string(chMsg.data))
}

func SendOneConnMessage(data interface{}, conn *UserConn) {
	bMsg, _ := json.Marshal(data)
	err := WS.writeMsg(conn, websocket.TextMessage, bMsg)
	wrapSdkLog("", "send response to web: ", string(bMsg), "userUid", WS.getUserUid(conn))
	if err != nil {
		wrapSdkLog("", "WS WriteMsg error", "", "userIP", conn.RemoteAddr().String(), "userUid", WS.getUserUid(conn), "error", err, "data", data)
	} else {
		wrapSdkLog("", "WS WriteMsg ok", "data", data, "userUid", WS.getUserUid(conn))
	}
}

func send2ChForTest(ch chan ChanMsg, value ChanMsg, timeout int64) error {
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
		wrapSdkLog("", "send cmd timeout, ", timeout, value)
		return errors.New("send cmd timeout")
	}
}
