package ws_local_server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"open_im_sdk/open_im_sdk"
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

func (i *InitCallback) OnConnectFailed(ErrCode int, ErrMsg string) {
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
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if userWorker.InitSDK(config, &initcb) {
		//	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", "", operationID})
	} else {
		//	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), open_im_sdk.ErrCodeInitLogin, "init config failed", "", operationID})
	}
}

func (wsRouter *WsFuncRouter) UnInitSDK() {
	wrapSdkLog("UnInitSDK uid: ", wsRouter.uId)
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.UnInitSDK()
}

func (wsRouter *WsFuncRouter) Login(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed", err.Error())
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Login(m["uid"].(string), m["token"].(string), &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) Logout(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Logout(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

//1
func (wsRouter *WsFuncRouter) GetLoginStatus(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", int32ToString(int32(userWorker.GetLoginStatus())), operationID})
}

//1
func (wsRouter *WsFuncRouter) GetLoginUser(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", userWorker.GetLoginUser(), operationID})
}

func InitServer(config *open_im_sdk.IMConfig) {
	data, _ := json.Marshal(config)
	ConfigSvr = string(data)
	UserRouteMap = make(map[string]RefRouter, 0)
	open_im_sdk.InitOnce(config)
}

func (wsRouter *WsFuncRouter) GlobalSendMessage(data interface{}) {
	conn := WS.getUserConn(wsRouter.uId + " " + "Web")
	if conn == nil {
		wrapSdkLog("uid no conn ", wsRouter.uId)
	}
	bMsg, _ := json.Marshal(data)
	if conn != nil {
		err := WS.writeMsg(conn, websocket.TextMessage, bMsg)
		wrapSdkLog("sendmsg:", string(bMsg))
		if err != nil {
			wrapSdkLog("WS WriteMsg error", "", "userIP", conn.RemoteAddr().String(), "userUid", WS.getUserUid(conn), "error", err, "data", data)
		}
	} else {
		wrapSdkLog("Conn is nil", "data", data)
	}

}

func SendOneUserMessage(data interface{}, uid string) {
	conn := WS.getUserConn(uid + " " + "Web")
	if conn == nil {
		wrapSdkLog("uid no conn ", uid)
	}
	bMsg, _ := json.Marshal(data)
	if conn != nil {
		err := WS.writeMsg(conn, websocket.TextMessage, bMsg)
		wrapSdkLog("sendmsg:", string(bMsg), uid)
		if err != nil {
			wrapSdkLog("WS WriteMsg error", "", "userIP", conn.RemoteAddr().String(), "userUid", WS.getUserUid(conn), "error", err, "data", data)
		}
	} else {
		wrapSdkLog("Conn is nil", "data", data, uid)
	}

}
