/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/9/8 14:42).
 */
package ws_local_server

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/log"
	utils2 "open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"open_im_sdk/ws_wrapper/utils"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const WriteTimeoutSeconds = 30
const POINTNUM = 10

var (
	rwLock *sync.RWMutex
	WS     WServer
)

type UserConn struct {
	*websocket.Conn
	w *sync.Mutex
}
type ChanMsg struct {
	data []byte
	uid  string
}
type WServer struct {
	wsAddr       string
	wsMaxConnNum int
	wsUpGrader   *websocket.Upgrader
	wsConnToUser map[*UserConn]map[string]string
	wsUserToConn map[string]map[string]*UserConn
	ch           chan ChanMsg
}

func (ws *WServer) OnInit(wsPort int) {
	//ip := utils.ServerIP
	ws.wsAddr = ":" + utils.IntToString(wsPort)
	ws.wsMaxConnNum = 10000
	ws.wsConnToUser = make(map[*UserConn]map[string]string)
	ws.wsUserToConn = make(map[string]map[string]*UserConn)
	ws.ch = make(chan ChanMsg, 100000)
	rwLock = new(sync.RWMutex)
	ws.wsUpGrader = &websocket.Upgrader{
		HandshakeTimeout: 10 * time.Second,
		ReadBufferSize:   4096,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}
}

func (ws *WServer) Run() {
	go ws.getMsgAndSend()
	http.HandleFunc("/", ws.wsHandler)         //Get request from client to handle by wsHandler
	err := http.ListenAndServe(ws.wsAddr, nil) //Start listening
	if err != nil {
		log.Info("", "Ws listening err", "", "err", err.Error())
	}
}

func (ws *WServer) getMsgAndSend() {
	defer func() {
		if r := recover(); r != nil {
			log.Info("", "getMsgAndSend panic", " panic is ", r, debug.Stack())
			ws.getMsgAndSend()
			log.Info("", "goroutine getMsgAndSend restart")
		}
	}()
	for {
		select {
		case r := <-ws.ch:
			go func() {
				operationID := utils2.OperationIDGenerator()
				log.Info(operationID, "getMsgAndSend channel: ", string(r.data), r.uid)

				//		conns := ws.getUserConn(r.uid + " " + "Web")
				conns := ws.getUserConn(r.uid + " " + utils.PlatformIDToName(sdk_struct.SvrConf.Platform))
				if conns == nil {
					log.Error(operationID, "uid no conn, failed ", r.uid+" "+utils.PlatformIDToName(sdk_struct.SvrConf.Platform))
					r.data = nil
				}
				log.Info(operationID, "conns  ", conns, r.uid+" "+utils.PlatformIDToName(sdk_struct.SvrConf.Platform))
				for _, conn := range conns {
					if conn != nil {
						err := WS.writeMsg(conn, websocket.TextMessage, r.data)
						if err != nil {
							log.Error(operationID, "WS WriteMsg error", "", "userIP", conn.RemoteAddr().String(), "userUid", r.uid, "error", err.Error())
						} else {
							log.Info(operationID, "writeMsg  ", conn.RemoteAddr(), string(r.data), r.uid)
						}
					} else {
						log.Error(operationID, "Conn is nil, failed")
					}
				}
				r.data = nil
			}()
		}
	}
}

func (ws *WServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	operationID := utils2.OperationIDGenerator()
	defer func() {
		if r := recover(); r != nil {
			log.Info(operationID, "wsHandler panic recover", " panic is ", r)
			buf := make([]byte, 1<<20)
			runtime.Stack(buf, true)
			log.Info(operationID, "panic", "call", string(buf))
		}
	}()
	//var mem runtime.MemStats
	//runtime.ReadMemStats(&mem)
	//if mem.Alloc > 2*1024*1024*1024 {
	//	panic("Memory leak " + int64ToString(int64(mem.Alloc)))
	//}
	//log.Info(operationID, "wsHandler ", r.URL.Query(), "js sdk svr mem: ", mem.Alloc, mem.TotalAlloc, "all: ", mem)

	if ws.headerCheck(w, r, operationID) {
		query := r.URL.Query()
		conn, err := ws.wsUpGrader.Upgrade(w, r, nil) //Conn is obtained through the upgraded escalator
		if err != nil {
			log.Info(operationID, "upgrade http conn err", "", "err", err)
			return
		} else {

			sendIDAndPlatformID := query["sendID"][0] + " " + utils.PlatformIDToName(int32(utils.StringToInt64(query["platformID"][0])))
			newConn := &UserConn{conn, new(sync.Mutex)}
			ws.addUserConn(sendIDAndPlatformID, newConn, operationID)
			go ws.readMsg(newConn, sendIDAndPlatformID)
		}
	} else {
		log.NewError(operationID, "headerCheck failed")
	}
}

func pMem() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Println("mem for test ", m)
	fmt.Println("mem for test os ", m.Sys)
	fmt.Println("mem for test HeapAlloc ", m.HeapAlloc)
}
func (ws *WServer) readMsg(conn *UserConn, sendIDAndPlatformID string) {
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Info("", "ReadMessage error", "", "userIP", conn.RemoteAddr().String(), "userUid", sendIDAndPlatformID, "error", err.Error())

			//log.Info("debug memory delUserConn begin ")
			//time.Sleep(1 * time.Second)

			ws.delUserConn(conn)
			//log.Info("debug memory delUserConn end  ")
			//time.Sleep(1 * time.Second)
			return
		} else {
			log.Info("", "ReadMessage ok ", "", "msgType", msgType, "userIP", conn.RemoteAddr().String(), "userUid", sendIDAndPlatformID)
		}
		m := Req{}
		json.Unmarshal(msg, &m)

		//log.Info("debug memory msgParse begin ", m)
		//time.Sleep(1 * time.Second)

		ws.msgParse(conn, msg)
		//log.Info("debug memory msgParse end ", m)
		//time.Sleep(1 * time.Second)
	}
}

func (ws *WServer) writeMsg(conn *UserConn, a int, msg []byte) error {
	conn.w.Lock()
	defer conn.w.Unlock()
	conn.SetWriteDeadline(time.Now().Add(time.Duration(WriteTimeoutSeconds) * time.Second))
	return conn.WriteMessage(a, msg)

}
func (ws *WServer) addUserConn(uid string, conn *UserConn, operationID string) {
	rwLock.Lock()

	var flag int32
	if oldConnMap, ok := ws.wsUserToConn[uid]; ok {
		flag = 1
		oldConnMap[conn.RemoteAddr().String()] = conn
		ws.wsUserToConn[uid] = oldConnMap
		log.Info(operationID, "this user is not first login", "", "uid", uid)
		//err := oldConn.Close()
		//delete(ws.wsConnToUser, oldConn)
		//if err != nil {
		//	log.Info("", "close err", "", "uid", uid, "conn", conn)
		//}
	} else {
		i := make(map[string]*UserConn)
		i[conn.RemoteAddr().String()] = conn
		ws.wsUserToConn[uid] = i
		log.Info(operationID, "this user is first login", "", "uid", uid)
	}
	if oldStringMap, ok := ws.wsConnToUser[conn]; ok {
		oldStringMap[conn.RemoteAddr().String()] = uid
		ws.wsConnToUser[conn] = oldStringMap
		log.Info(operationID, "find failed", "", "uid", uid)
		//err := oldConn.Close()
		//delete(ws.wsConnToUser, oldConn)
		//if err != nil {
		//	log.Info("", "close err", "", "uid", uid, "conn", conn)
		//}
	} else {
		i := make(map[string]string)
		i[conn.RemoteAddr().String()] = uid
		ws.wsConnToUser[conn] = i
		log.Info(operationID, "this user is first login", "", "uid", uid)
	}
	log.Info(operationID, "WS Add operation", "", "wsUser added", ws.wsUserToConn, "uid", uid, "online_num", len(ws.wsUserToConn))
	rwLock.Unlock()

	//log.Info("", "after add, wsConnToUser map ", ws.wsConnToUser)
	//	log.Info("", "after add, wsUserToConn  map ", ws.wsUserToConn)

	if flag == 1 {
		//	DelUserRouter(uid)
	}

}
func (ws *WServer) getConnNum(uid string) int {
	rwLock.Lock()
	defer rwLock.Unlock()
	log.Info("", "getConnNum uid: ", uid)
	if connMap, ok := ws.wsUserToConn[uid]; ok {
		log.Info("", "uid->conn ", connMap)
		return len(connMap)
	} else {
		return 0
	}

}

func (ws *WServer) delUserConn(conn *UserConn) {
	operationID := utils2.OperationIDGenerator()
	rwLock.Lock()
	var uidPlatform string
	if oldStringMap, ok := ws.wsConnToUser[conn]; ok {
		uidPlatform = oldStringMap[conn.RemoteAddr().String()]
		if oldConnMap, ok := ws.wsUserToConn[uidPlatform]; ok {

			log.Info(operationID, "old map : ", oldConnMap, "conn: ", conn.RemoteAddr().String())
			delete(oldConnMap, conn.RemoteAddr().String())

			ws.wsUserToConn[uidPlatform] = oldConnMap
			log.Info(operationID, "WS delete operation", "", "wsUser deleted", ws.wsUserToConn, "uid", uidPlatform, "online_num", len(ws.wsUserToConn))
			if len(oldConnMap) == 0 {
				log.Info(operationID, "no conn delete user router ", uidPlatform)
				log.Info(operationID, "DelUserRouter ", uidPlatform)
				DelUserRouter(uidPlatform, operationID)
				ws.wsUserToConn[uidPlatform] = make(map[string]*UserConn)
				delete(ws.wsUserToConn, uidPlatform)
			}
		} else {
			log.Info(operationID, "uid not exist", "", "wsUser deleted", ws.wsUserToConn, "uid", uidPlatform, "online_num", len(ws.wsUserToConn))
		}
		oldStringMap = make(map[string]string)
		delete(ws.wsConnToUser, conn)

	}
	err := conn.Close()
	if err != nil {
		log.Info(operationID, "close err", "", "uid", uidPlatform, "conn", conn)
	}
	rwLock.Unlock()
}

func (ws *WServer) getUserConn(uid string) (w []*UserConn) {
	rwLock.RLock()
	defer rwLock.RUnlock()
	t := ws.wsUserToConn

	if connMap, ok := t[uid]; ok {
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

func (ws *WServer) headerCheck(w http.ResponseWriter, r *http.Request, operationID string) bool {

	status := http.StatusUnauthorized
	query := r.URL.Query()
	log.Info(operationID, "headerCheck: ", query["token"], query["platformID"], query["sendID"], r.RemoteAddr)
	if len(query["token"]) != 0 && len(query["sendID"]) != 0 && len(query["platformID"]) != 0 {
		SendID := query["sendID"][0] + " " + utils.PlatformIDToName(int32(utils.StringToInt64(query["platformID"][0])))
		if ws.getConnNum(SendID) >= POINTNUM {
			log.Info(operationID, "Over quantity failed", query, ws.getConnNum(SendID), SendID)
			w.Header().Set("Sec-Websocket-Version", "13")
			http.Error(w, "Over quantity", status)
			return false
		}
		//if utils.StringToInt(query["platformID"][0]) != utils.WebPlatformID {
		//	log.Info("check platform id failed", query["sendID"][0], query["platformID"][0])
		//	w.Header().Set("Sec-Websocket-Version", "13")
		//	http.Error(w, http.StatusText(status), StatusBadRequest)
		//	return false
		//}
		checkFlag := open_im_sdk.CheckToken(query["sendID"][0], query["token"][0], operationID)
		if checkFlag != nil {
			log.Info(operationID, "check token failed", query["sendID"][0], query["token"][0], checkFlag.Error())
			w.Header().Set("Sec-Websocket-Version", "13")
			http.Error(w, http.StatusText(status), status)
			return false
		}
		log.Info(operationID, "Connection Authentication Success", "", "token", query["token"][0], "userID", query["sendID"][0], "platformID", query["platformID"][0])
		return true

	} else {
		log.Info(operationID, "Args err", "", "query", query)
		w.Header().Set("Sec-Websocket-Version", "13")
		http.Error(w, http.StatusText(status), StatusBadRequest)
		return false
	}
}
