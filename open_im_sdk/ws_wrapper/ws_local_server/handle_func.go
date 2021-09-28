/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/9/8 14:54).
 */
package ws_local_server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"reflect"
)

type Req struct {
	ReqFuncName string `json:"reqFuncName" `
	OperationID string `json:"operationID"`
	Data        string `json:"data"`
	UId         string `json:"uid"`
}

func (ws *WServer) msgParse(conn *websocket.Conn, jsonMsg []byte) {
	m := Req{}
	if err := json.Unmarshal(jsonMsg, &m); err != nil {
		//GlobalSendMessage(EventData{m.ReqFuncName, -1, "ws json Unmarshal err ", "", "0"})
		return
	}
	if m.OperationID == "" {
		//	GlobalSendMessage(EventData{m.ReqFuncName, -2, "no OperationID", "", "0"})
		return
	}
	wrapSdkLog("Basic Info Authentication Success", "reqFuncName ", m.ReqFuncName, "data ", m.Data, "recv jsonMsg: ", string(jsonMsg))

	if m.ReqFuncName == "Login" {
		wrapSdkLog("login ", m.UId)
		GenUserRouter(m.UId)
	}
	UserRouteRwLock.RLock()
	defer UserRouteRwLock.RUnlock()
	urm, ok := UserRouteMap[m.UId]
	if !ok {
		wrapSdkLog("user not login error: ", m.UId)
		return
	}
	parms := []reflect.Value{reflect.ValueOf(m.Data), reflect.ValueOf(m.OperationID)}
	vf, ok := (*urm.refName)[m.ReqFuncName]
	if ok {
		vf.Call(parms)
	} else {
		//	GlobalSendMessage(EventData{m.ReqFuncName, -1, "no func ", "", m.OperationID})
	}
	defer func() {
		if r := recover(); r != nil {
			wrapSdkLog("panic is ", r)
			//	GlobalSendMessage(EventData{m.ReqFuncName, -3, "panic ", "", "0"})
		}
	}()
}

func (ws *WServer) sendMsg(conn *websocket.Conn, mReply map[string]interface{}) {
	bMsg, _ := json.Marshal(mReply)
	err := ws.writeMsg(conn, websocket.TextMessage, bMsg)
	if err != nil {
		wrapSdkLog("WS WriteMsg error", "", "userIP", conn.RemoteAddr().String(), "userUid", ws.getUserUid(conn), "error", err, "mReply", mReply)
	}
}

func (ws *WServer) sendErrMsg(conn *websocket.Conn, errCode int32, errMsg string) {
	mReply := make(map[string]interface{})
	mReply["errCode"] = errCode
	mReply["errMsg"] = errMsg
	ws.sendMsg(conn, mReply)
}
