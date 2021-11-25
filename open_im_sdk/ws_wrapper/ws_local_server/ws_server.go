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
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/open_im_sdk/ws_wrapper/utils"
	"runtime"
	"strings"
	"sync"
	"time"
)

const POINTNUM = 1

var (
	rwLock *sync.RWMutex
	WS     WServer
)

type UserConn struct {
	*websocket.Conn
	w *sync.Mutex
}

type WServer struct {
	wsAddr       string
	wsMaxConnNum int
	wsUpGrader   *websocket.Upgrader
	wsConnToUser map[*UserConn]map[string]string
	wsUserToConn map[string]map[string]*UserConn
}

func (ws *WServer) OnInit(wsPort int) {
	ip := utils.ServerIP
	ws.wsAddr = ip + ":" + utils.IntToString(wsPort)
	ws.wsMaxConnNum = 10000
	ws.wsConnToUser = make(map[*UserConn]map[string]string)
	ws.wsUserToConn = make(map[string]map[string]*UserConn)
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
	wrapSdkLog("wsHandler ", r.URL.Query())
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
			newConn := &UserConn{conn, new(sync.Mutex)}
			ws.addUserConn(SendID, newConn)
			go ws.readMsg(newConn)
		}
	}
}

func (ws *WServer) readMsg(conn *UserConn) {
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			wrapSdkLog("ReadMessage error", "", "userIP", conn.RemoteAddr().String(), "userUid", ws.getUserUid(conn), "error", err)
			ws.delUserConn(conn)
			return
		} else {
			wrapSdkLog("ReadMessage ok ", "", "msgType", msgType, "userIP", conn.RemoteAddr().String(), "userUid", ws.getUserUid(conn))
		}
		ws.msgParse(conn, msg)
	}
}

func (ws *WServer) writeMsg(conn *UserConn, a int, msg []byte) error {
	conn.w.Lock()
	defer conn.w.Unlock()
	return conn.WriteMessage(a, msg)

}
func (ws *WServer) addUserConn(uid string, conn *UserConn) {
	wrapSdkLog("addUserConn", uid)
	rwLock.Lock()
	wrapSdkLog("addUserConn lock", uid)
	var flag int32
	if oldConnMap, ok := ws.wsUserToConn[uid]; ok {
		flag = 1
		oldConnMap[conn.RemoteAddr().String()] = conn
		ws.wsUserToConn[uid] = oldConnMap
		wrapSdkLog("this user is not first login", "", "uid", uid)
		//err := oldConn.Close()
		//delete(ws.wsConnToUser, oldConn)
		//if err != nil {
		//	wrapSdkLog("close err", "", "uid", uid, "conn", conn)
		//}
	} else {
		i := make(map[string]*UserConn)
		i[conn.RemoteAddr().String()] = conn
		ws.wsUserToConn[uid] = i
		wrapSdkLog("this user is first login", "", "uid", uid)
	}
	if oldStringMap, ok := ws.wsConnToUser[conn]; ok {
		oldStringMap[conn.RemoteAddr().String()] = uid
		ws.wsConnToUser[conn] = oldStringMap
		wrapSdkLog("find failed", "", "uid", uid)
		//err := oldConn.Close()
		//delete(ws.wsConnToUser, oldConn)
		//if err != nil {
		//	wrapSdkLog("close err", "", "uid", uid, "conn", conn)
		//}
	} else {
		i := make(map[string]string)
		i[conn.RemoteAddr().String()] = uid
		ws.wsConnToUser[conn] = i
		wrapSdkLog("this user is first login", "", "uid", uid)
	}
	wrapSdkLog("WS Add operation", "", "wsUser added", ws.wsUserToConn, "uid", uid, "online_num", len(ws.wsUserToConn))
	rwLock.Unlock()

	if flag == 1 {
		//	DelUserRouter(uid)
	}

}
func (ws *WServer) getConnNum(uid string) int {
	rwLock.Lock()
	defer rwLock.Unlock()
	if connMap, ok := ws.wsUserToConn[uid]; ok {
		wrapSdkLog("uid->conn ", connMap)
		return len(connMap)
	} else {
		return 0
	}

}

func (ws *WServer) delUserConn(conn *UserConn) {
	rwLock.Lock()

	var uidPlatform string
	if oldStringMap, ok := ws.wsConnToUser[conn]; ok {
		uidPlatform = oldStringMap[conn.RemoteAddr().String()]
		if oldConnMap, ok := ws.wsUserToConn[uidPlatform]; ok {

			wrapSdkLog("old map : ", oldConnMap, "conn: ", conn.RemoteAddr().String())
			delete(oldConnMap, conn.RemoteAddr().String())

			ws.wsUserToConn[uidPlatform] = oldConnMap
			wrapSdkLog("WS delete operation", "", "wsUser deleted", ws.wsUserToConn, "uid", uidPlatform, "online_num", len(ws.wsUserToConn))
			if len(oldConnMap) == 0 {
				wrapSdkLog("no conn delete user router ", uidPlatform)
				wrapSdkLog("DelUserRouter ", uidPlatform)
				DelUserRouter(uidPlatform)
				delete(ws.wsUserToConn, uidPlatform)
			}
		} else {
			wrapSdkLog("uid not exist", "", "wsUser deleted", ws.wsUserToConn, "uid", uidPlatform, "online_num", len(ws.wsUserToConn))
		}
		delete(ws.wsConnToUser, conn)

	}
	err := conn.Close()
	if err != nil {
		wrapSdkLog("close err", "", "uid", uidPlatform, "conn", conn)
	}

	rwLock.Unlock()
}

func (ws *WServer) getUserConn(uid string) (w []*UserConn) {
	rwLock.RLock()
	defer rwLock.RUnlock()
	if connMap, ok := ws.wsUserToConn[uid]; ok {
		for _, v := range connMap {
			w = append(w, v)
		}
		return w
	}
	return nil
}

func (ws *WServer) getUserUid(conn *UserConn) string {
	return "getUserUid"
}

func (ws *WServer) headerCheck(w http.ResponseWriter, r *http.Request) bool {

	status := http.StatusUnauthorized
	query := r.URL.Query()
	wrapSdkLog("headerCheck: ", query["token"], query["platformID"], query["sendID"])
	if len(query["token"]) != 0 && len(query["sendID"]) != 0 && len(query["platformID"]) != 0 {
		SendID := query["sendID"][0] + " " + utils.PlatformIDToName(int32(utils.StringToInt64(query["platformID"][0])))
		if ws.getConnNum(SendID) >= POINTNUM {
			wrapSdkLog("Over quantity failed", query, ws.getConnNum(SendID))
			w.Header().Set("Sec-Websocket-Version", "13")
			http.Error(w, "Over quantity", status)
			return false
		}
		if utils.StringToInt(query["platformID"][0]) != utils.WebPlatformID {
			wrapSdkLog("check platform id failed", query["sendID"][0], query["platformID"][0])
			w.Header().Set("Sec-Websocket-Version", "13")
			http.Error(w, http.StatusText(status), StatusBadRequest)
			return false
		}
		checkFlag := open_im_sdk.CheckToken(query["sendID"][0], query["token"][0])
		if checkFlag != 0 {
			wrapSdkLog("check token failed", query["sendID"][0], query["token"][0])
			w.Header().Set("Sec-Websocket-Version", "13")
			http.Error(w, http.StatusText(status), status)
			return false
		}
		wrapSdkLog("Connection Authentication Success", "", "token", query["token"][0], "userID", query["sendID"][0], "platformID", query["platformID"][0])
		return true

	} else {
		wrapSdkLog("Args err", "", "query", query)
		w.Header().Set("Sec-Websocket-Version", "13")
		http.Error(w, http.StatusText(status), StatusBadRequest)
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
