/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/9/8 14:54).
 */
package ws_local_server

import (
	"encoding/json"
	"reflect"
)

type Req struct {
	ReqFuncName string `json:"reqFuncName" `
	OperationID string `json:"operationID"`
	Data        string `json:"data"`
	UId         string `json:"uid"`
}

func (ws *WServer) DoLogin(m Req, conn *UserConn) {
	UserRouteRwLock.RLock()
	defer UserRouteRwLock.RUnlock()
	urm, ok := UserRouteMap[m.UId]
	if !ok {
		wrapSdkLog("user first login: ", m)
		refR := GenUserRouterNoLock(m.UId)
		params := []reflect.Value{reflect.ValueOf(m.Data), reflect.ValueOf(m.OperationID)}
		vf, ok := (*refR.refName)[m.ReqFuncName]
		if ok {
			vf.Call(params)
		} else {
			wrapSdkLog("no func name: ", m.ReqFuncName, m)
			SendOneConnMessage(EventData{m.ReqFuncName, StatusBadParameter, StatusText(StatusBadParameter), "", m.OperationID}, conn)
		}

	} else {
		if urm.wsRouter.getMyLoginStatus() == 101 {
			//send ok
			SendOneConnMessage(EventData{"Login", 0, "ok", "", m.OperationID}, conn)
		} else {
			wrapSdkLog("login status pending, try after 5 second ", urm.wsRouter.getMyLoginStatus(), m.UId)
			SendOneConnMessage(EventData{"Login", StatusLoginPending, StatusText(StatusLoginPending), "", m.OperationID}, conn)
		}
	}
}

func (ws *WServer) msgParse(conn *UserConn, jsonMsg []byte) {
	m := Req{}
	if err := json.Unmarshal(jsonMsg, &m); err != nil {
		SendOneConnMessage(EventData{"error", 100, "Unmarshal failed", "", ""}, conn)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			SendOneConnMessage(EventData{m.ReqFuncName, StatusBadParameter, StatusText(StatusBadParameter), "", m.OperationID}, conn)
			wrapSdkLog("bad request, panic is ", r)
		}
	}()

	wrapSdkLog("Basic Info", "reqFuncName ", m.ReqFuncName, "data ", m.Data, "recv jsonMsg: ", string(jsonMsg))

	if m.ReqFuncName == "Login" {
		ws.DoLogin(m, conn)
		wrapSdkLog("login ", m)
		return
	}

	UserRouteRwLock.RLock()
	defer UserRouteRwLock.RUnlock()
	urm, ok := UserRouteMap[m.UId]
	if !ok {
		wrapSdkLog("user not login failed, must login first: ", m.UId)
		SendOneConnMessage(EventData{"Login", StatusNoLogin, StatusText(StatusNoLogin), "", m.OperationID}, conn)
		return
	}
	parms := []reflect.Value{reflect.ValueOf(m.Data), reflect.ValueOf(m.OperationID)}
	vf, ok := (*urm.refName)[m.ReqFuncName]
	if ok {
		vf.Call(parms)
	} else {
		wrapSdkLog("no func ", m.ReqFuncName)
		SendOneConnMessage(EventData{m.ReqFuncName, StatusBadParameter, StatusText(StatusBadParameter), "", m.OperationID}, conn)
	}

}

//func (ws *WServer) sendMsg(conn *UserConn, mReply map[string]interface{}) {
//	bMsg, _ := json.Marshal(mReply)
//	err := ws.writeMsg(conn, websocket.TextMessage, bMsg)
//	if err != nil {
//		wrapSdkLog("WS WriteMsg error", "", "userIP", conn.RemoteAddr().String(), "userUid", ws.getUserUid(conn), "error", err, "mReply", mReply)
//	}
//}
//
//func (ws *WServer) sendErrMsg(conn *UserConn, errCode int32, errMsg string) {
//	mReply := make(map[string]interface{})
//	mReply["errCode"] = errCode
//	mReply["errMsg"] = errMsg
//	ws.sendMsg(conn, mReply)
//}
