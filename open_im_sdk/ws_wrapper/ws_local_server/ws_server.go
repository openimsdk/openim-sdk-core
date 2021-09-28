/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/9/8 14:42).
 */
package ws_local_server

import (
	"github.com/gorilla/websocket"
	sLog "log"
	"net/http"
	"open_im_sdk/open_im_sdk/ws_wrapper/utils"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	rwLock *sync.RWMutex
	WS     WServer
)

type WServer struct {
	wsAddr       string
	wsMaxConnNum int
	wsUpGrader   *websocket.Upgrader
	wsConnToUser map[*websocket.Conn]string
	wsUserToConn map[string]*websocket.Conn
}

func (ws *WServer) OnInit(wsPort int) {
	ip := utils.ServerIP
	ws.wsAddr = ip + ":" + utils.IntToString(wsPort)
	ws.wsMaxConnNum = 10000
	ws.wsConnToUser = make(map[*websocket.Conn]string)
	ws.wsUserToConn = make(map[string]*websocket.Conn)
	rwLock = new(sync.RWMutex)
	ws.wsUpGrader = &websocket.Upgrader{
		HandshakeTimeout: 10 * time.Second,
		ReadBufferSize:   4096,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}
}

func (ws *WServer) Run() {
	http.HandleFunc("/", ws.wsHandler)         //Get request from client to handle by wsHandler
	err := http.ListenAndServe(ws.wsAddr, nil) //Start listening
	if err != nil {
		wrapSdkLog("Ws listening err", "", "err", err.Error())
	}
}

func (ws *WServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	if ws.headerCheck(w, r) {
		query := r.URL.Query()
		conn, err := ws.wsUpGrader.Upgrade(w, r, nil) //Conn is obtained through the upgraded escalator
		if err != nil {
			wrapSdkLog("upgrade http conn err", "", "err", err)
			return
		} else {
			//Connection mapping relationship,
			//userID+" "+platformID->conn
			SendID := query["sendID"][0] + " " + utils.PlatformIDToName(int32(utils.StringToInt64(query["platformID"][0])))

			ws.addUserConn(SendID, conn)
			go ws.readMsg(conn)
		}
	}
}

func (ws *WServer) readMsg(conn *websocket.Conn) {
	wrapSdkLog("readMs: ")
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			wrapSdkLog("WS ReadMsg error", "", "userIP", conn.RemoteAddr().String(), "userUid", ws.getUserUid(conn), "error", err)
			ws.delUserConn(conn)
			return
		} else {
			wrapSdkLog("test", "", "msgType", msgType, "userIP", conn.RemoteAddr().String(), "userUid", ws.getUserUid(conn))
		}
		ws.msgParse(conn, msg)
		//ws.writeMsg(conn, 1, chat)
	}

}

func (ws *WServer) writeMsg(conn *websocket.Conn, a int, msg []byte) error {
	rwLock.Lock()
	defer rwLock.Unlock()
	return conn.WriteMessage(a, msg)

}
func (ws *WServer) addUserConn(uid string, conn *websocket.Conn) {
	wrapSdkLog("addUserConn", uid)
	rwLock.Lock()
	wrapSdkLog("addUserConn lock", uid)
	var flag int32
	if oldConn, ok := ws.wsUserToConn[uid]; ok {
		flag = 1
		err := oldConn.Close()
		delete(ws.wsConnToUser, oldConn)
		if err != nil {
			wrapSdkLog("close err", "", "uid", uid, "conn", conn)
		}
	} else {
		wrapSdkLog("this user is first login", "", "uid", uid)
	}
	ws.wsConnToUser[conn] = uid
	ws.wsUserToConn[uid] = conn
	wrapSdkLog("WS Add operation", "", "wsUser added", ws.wsUserToConn, "uid", uid, "online_num", len(ws.wsUserToConn))
	rwLock.Unlock()

	if flag == 1 {
		DelUserRouter(uid)
	}

}

func (ws *WServer) delUserConn(conn *websocket.Conn) {
	rwLock.Lock()

	var uidPlatform string
	if uid, ok := ws.wsConnToUser[conn]; ok {
		uidPlatform = uid
		if _, ok = ws.wsUserToConn[uid]; ok {
			delete(ws.wsUserToConn, uid)
			wrapSdkLog("WS delete operation", "", "wsUser deleted", ws.wsUserToConn, "uid", uid, "online_num", len(ws.wsUserToConn))
		} else {
			wrapSdkLog("uid not exist", "", "wsUser deleted", ws.wsUserToConn, "uid", uid, "online_num", len(ws.wsUserToConn))
		}
		delete(ws.wsConnToUser, conn)
	}
	err := conn.Close()
	if err != nil {
		wrapSdkLog("close err", "", "uid", uidPlatform, "conn", conn)
	}
	rwLock.Unlock()
	DelUserRouter(uidPlatform)

}

func (ws *WServer) getUserConn(uid string) *websocket.Conn {
	rwLock.RLock()
	defer rwLock.RUnlock()
	if conn, ok := ws.wsUserToConn[uid]; ok {
		return conn
	}
	return nil
}
func (ws *WServer) getUserUid(conn *websocket.Conn) string {
	rwLock.RLock()
	defer rwLock.RUnlock()

	if uid, ok := ws.wsConnToUser[conn]; ok {
		return uid
	}
	return ""
}
func (ws *WServer) headerCheck(w http.ResponseWriter, r *http.Request) bool {
	status := http.StatusUnauthorized
	query := r.URL.Query()
	if len(query["token"]) != 0 && len(query["sendID"]) != 0 && len(query["platformID"]) != 0 {
		if utils.StringToInt(query["platformID"][0]) != utils.WebPlatformID {

		}
		wrapSdkLog("Connection Authentication Success", "", "token", query["token"][0], "userID", query["sendID"][0])
		return true

	} else {
		wrapSdkLog("Args err", "", "query", query)
		w.Header().Set("Sec-Websocket-Version", "13")
		http.Error(w, http.StatusText(status), status)
		return false
	}

}
func wrapSdkLog(v ...interface{}) {
	_, b, c, _ := runtime.Caller(1)
	i := strings.LastIndex(b, "/")
	if i != -1 {
		sLog.Println("[", b[i+1:len(b)], ":", c, "]", v)
	}
}
