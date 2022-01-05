package open_im_sdk

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	sLog "log"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func operationIDGenerator() string {
	return strconv.FormatInt(time.Now().UnixNano()+int64(rand.Uint32()), 10)
}
func getMsgID(sendID string) string {
	t := int64ToString(getCurrentTimestampByNano())
	return Md5(t + sendID + int64ToString(rand.Int63n(getCurrentTimestampByNano())))
}
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	cipher := h.Sum(nil)
	return hex.EncodeToString(cipher)
}

//Get the current timestamp by Second

func getCurrentTimestampBySecond() int64 {
	return time.Now().Unix()
}

//Get the current timestamp by Mill
func GetCurrentTimestampByMill() int64 {
	return time.Now().UnixNano() / 1e6
}

//Get the current timestamp by Nano
func getCurrentTimestampByNano() int64 {
	return time.Now().UnixNano()
}

//wsNotification map[string]chan GeneralWsResp

func (u *UserRelated) AddCh() (string, chan GeneralWsResp) {
	LogBegin()
	u.wsMutex.Lock()
	defer u.wsMutex.Unlock()
	msgIncr := u.GenMsgIncr()
	sdkLog("msgIncr: ", msgIncr)
	ch := make(chan GeneralWsResp, 1)
	_, ok := u.wsNotification[msgIncr]
	if ok {
		sdkLog("AddCh exist")
	}
	u.wsNotification[msgIncr] = ch
	LogSReturn(msgIncr, ch)
	LogBegin(msgIncr, ch)
	return msgIncr, ch
}

func (u *UserRelated) GetCh(msgIncr string) chan GeneralWsResp {
	LogBegin(msgIncr)
	//u.wsMutex.RLock()
	//	defer u.wsMutex.RUnlock()
	ch, ok := u.wsNotification[msgIncr]
	if ok {
		sdkLog("GetCh ok")
		LogSReturn(ch)
		return ch
	}
	sdkLog("GetCh nil")
	LogFReturn(nil)
	return nil
}

func (u *UserRelated) DelCh(msgIncr string) {
	//	LogBegin(msgIncr)
	u.wsMutex.Lock()
	defer u.wsMutex.Unlock()
	ch, ok := u.wsNotification[msgIncr]
	if ok {
		close(ch)
		delete(u.wsNotification, msgIncr)
	}
	//	LogSReturn()
}

func (u *UserRelated) sendPingMsg() error {
	u.stateMutex.Lock()
	defer u.stateMutex.Unlock()
	var ping string = "try ping"

	err := u.conn.SetWriteDeadline(time.Now().Add(8 * time.Second))
	if err != nil {
		sdkLog("SetWriteDeadline failed ", err.Error())
	}
	return u.conn.WriteMessage(websocket.PingMessage, []byte(ping))
}

func (u *UserRelated) writeBinaryMsg(msg GeneralWsReq) (error, *websocket.Conn) {
	LogStart(msg.OperationID)
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(msg)
	if err != nil {
		LogFReturn(err.Error())
		return err, nil
	}

	var connSended *websocket.Conn
	u.stateMutex.Lock()
	defer u.stateMutex.Unlock()

	if u.conn != nil {
		connSended = u.conn
		err = u.conn.SetWriteDeadline(time.Now().Add(8 * time.Second))
		if err != nil {
			sdkLog("SetWriteDeadline failed ", err.Error())
		}
		sdkLog("send ws BinaryMessage len: ", len(buff.Bytes()))
		if len(buff.Bytes()) > MaxTotalMsgLen {
			LogFReturn("msg too long", len(buff.Bytes()), MaxTotalMsgLen)
			return errors.New("msg too long"), connSended
		}
		err = u.conn.WriteMessage(websocket.BinaryMessage, buff.Bytes())
		if err != nil {
			LogFReturn(err.Error(), msg.OperationID)
		} else {
			LogSReturn(nil)
		}
		return err, connSended
	} else {
		LogFReturn("conn==nil")
		return errors.New("conn==nil"), connSended
	}
}

func (u *UserRelated) decodeBinaryWs(message []byte) (*GeneralWsResp, error) {
	LogStart()
	buff := bytes.NewBuffer(message)
	dec := gob.NewDecoder(buff)
	var data GeneralWsResp
	err := dec.Decode(&data)
	if err != nil {
		LogFReturn(nil, err.Error())
		return nil, err
	}
	LogSReturn(&data, nil)
	return &data, nil
}

func (u *UserRelated) WriteMsg(msg GeneralWsReq) (error, *websocket.Conn) {
	LogStart(msg.OperationID)
	LogSReturn(msg.OperationID)
	return u.writeBinaryMsg(msg)
}

func notifyCh(ch chan GeneralWsResp, value GeneralWsResp, timeout int64) error {
	var flag = 0
	select {
	case ch <- value:
		flag = 1
	case <-time.After(time.Second * time.Duration(timeout)):
		flag = 2
	}
	if flag == 1 {
		return nil
	} else {
		sdkLog("send cmd timeout, ", timeout, value)
		return errors.New("send cmd timeout")
	}
}

func sendCmd(ch chan cmd2Value, value cmd2Value, timeout int64) error {
	var flag = 0
	select {
	case ch <- value:
		flag = 1
	case <-time.After(time.Second * time.Duration(timeout)):
		flag = 2
	}
	if flag == 1 {
		return nil
	} else {
		sdkLog("send cmd timeout, ", timeout, value)
		return errors.New("send cmd timeout")
	}
}

func (u *UserRelated) GenMsgIncr() string {
	return u.loginUserID + "_" + int64ToString(getCurrentTimestampByNano())
}

func structToJsonString(param interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

//The incoming parameter must be a pointer
func jsonStringToStruct(s string, args interface{}) error {
	return wrap(json.Unmarshal([]byte(s), args), "json Unmarshal failed")
}

//Convert timestamp to time.Time type

func UnixSecondToTime(second int64) time.Time {
	return time.Unix(second, 0)
}
func intToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}
func int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}
func int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
func StringToInt64(i string) int64 {
	j, _ := strconv.ParseInt(i, 10, 64)
	return j
}

func stringToInt(i string) int {
	j, _ := strconv.Atoi(i)
	return j
}

func checkDiff(a []diff, b []diff) (aInBNot, bInANot, sameA, sameB []int) {
	//to map
	mapA := make(map[string]diff)
	for _, v := range a {
		mapA[v.Key()] = v
	}
	mapB := make(map[string]diff)
	for _, v := range b {
		mapB[v.Key()] = v
	}

	aInBNot = make([]int, 0)
	bInANot = make([]int, 0)
	sameA = make([]int, 0)
	sameB = make([]int, 0)

	//for a
	for i, v := range a {
		ia, ok := mapB[v.Key()]
		if !ok {
			//	sdkLog("aInBNot", i)
			aInBNot = append(aInBNot, i)
		} else {
			if ia.Value() != v.Value() {
				sameA = append(sameA, i)
			}
		}
	}

	//for b
	for i, v := range b {
		ib, ok := mapA[v.Key()]
		if !ok {
			//	sdkLog("bInANot", i)
			bInANot = append(bInANot, i)
		} else {
			if ib.Value() != v.Value() {
				sameB = append(sameB, i)
			}
		}
	}
	return aInBNot, bInANot, sameA, sameB
}

func (fr friendInfo) Key() string {
	return fr.UID
}
func (fr friendInfo) Value() interface{} {
	return fr
}

func (us userInfo) Key() string {
	return us.Uid
}
func (us userInfo) Value() interface{} {
	return us
}

func (ap applyUserInfo) Key() string {
	return ap.Uid
}

func (g groupInfo) Key() string {
	return g.GroupId
}

func (g groupInfo) Value() interface{} {
	return g
}
func (ap applyUserInfo) Value() interface{} {
	return ap
}

func (g groupMemberFullInfo) Key() string {
	return g.UserId
}

func (g groupMemberFullInfo) Value() interface{} {
	return g
}
func (g GroupReqListInfo) Key() string {
	return g.GroupID + g.FromUserID + g.ToUserID
}

func (g GroupReqListInfo) Value() interface{} {
	return g
}

func GetConversationIDBySessionType(sourceID string, sessionType int) string {
	switch sessionType {
	case SingleChatType:
		return "single_" + sourceID
	case GroupChatType:
		return "group_" + sourceID
	}
	return ""
}
func getIsRead(b bool) int {
	if b {
		return HasRead
	} else {
		return NotRead
	}
}
func getIsFilter(b bool) int {
	if b {
		return IsFilter
	} else {
		return NotFilter
	}
}
func getIsReadB(i int) bool {
	if i == HasRead {
		return true
	} else {
		return false
	}

}

func RunFuncName() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()
}

func cleanUpfuncName(funcName string) string {
	end := strings.LastIndex(funcName, ".")
	if end == -1 {
		return ""
	}
	return funcName[end+1:]
}

func LogBegin(v ...interface{}) {
	if SdkLogFlag == 1 {
		return
	}
	if logger != nil {
		NewInfo("", v...)
		return
	}
	pc, b, c, _ := runtime.Caller(1)
	fname := runtime.FuncForPC(pc).Name()
	i := strings.LastIndex(b, "/")
	if i != -1 {
		sLog.Println(" [", b[i+1:len(b)], ":", c, "]", cleanUpfuncName(fname), "call func begin, args: ", v)
	}
}

func LogEnd(v ...interface{}) {
	if SdkLogFlag == 1 {
		return
	}
	if logger != nil {
		NewInfo("", v...)
		return
	}
	pc, b, c, _ := runtime.Caller(1)
	fname := runtime.FuncForPC(pc).Name()
	i := strings.LastIndex(b, "/")
	if i != -1 {
		sLog.Println(" [", b[i+1:len(b)], ":", c, "]", cleanUpfuncName(fname), "call func end, args: ", v)
	}
}

func LogStart(v ...interface{}) {
	if SdkLogFlag == 1 {
		return
	}
	if logger != nil {
		NewInfo("", v...)
		return
	}
	pc, b, c, _ := runtime.Caller(1)
	fname := runtime.FuncForPC(pc).Name()
	i := strings.LastIndex(b, "/")
	if i != -1 {
		sLog.Println(" [", b[i+1:len(b)], ":", c, "]", cleanUpfuncName(fname), "func start, args: ", v)
	}
}

func LogFReturn(v ...interface{}) {
	if SdkLogFlag == 1 {
		return
	}
	if logger != nil {
		NewInfo("", v...)
		return
	}
	pc, b, c, _ := runtime.Caller(1)
	fname := runtime.FuncForPC(pc).Name()
	i := strings.LastIndex(b, "/")
	if i != -1 {
		sLog.Println("[", b[i+1:len(b)], ":", c, "]", cleanUpfuncName(fname), "failed return args(info): ", v)
	}
}

func LogSReturn(v ...interface{}) {
	if SdkLogFlag == 1 {
		return
	}
	if logger != nil {
		NewInfo("", v...)
		return
	}
	pc, b, c, _ := runtime.Caller(1)
	fname := runtime.FuncForPC(pc).Name()
	i := strings.LastIndex(b, "/")
	if i != -1 {
		sLog.Println("[", b[i+1:len(b)], ":", c, "]", cleanUpfuncName(fname), "success return args(info): ", v)
	}

}

func sdkLog(v ...interface{}) {
	if SdkLogFlag == 1 {
		return
	}
	if logger != nil {
		NewInfo("", v...)
		return
	}
	_, b, c, _ := runtime.Caller(1)
	i := strings.LastIndex(b, "/")
	if i != -1 {
		sLog.Println("[", b[i+1:len(b)], ":", c, "]", v)
	}

}

func sdkErrLog(err error, v ...interface{}) {
	if SdkLogFlag == 1 {
		return
	}
	if logger != nil {
		NewInfo("", v...)
		return
	}
	_, b, c, _ := runtime.Caller(1)
	i := strings.LastIndex(b, "/")
	if i != -1 {
		sLog.Println("[", b[i+1:len(b)], ":", c, "]", v)
		sLog.Print("%v", err)
	}

}

type LogInfo struct {
	Info string `json:"info"`
}

func log(info string) error {
	if SdkLogFlag == 1 {
		return nil
	}
	sdkLog(info)
	return nil
}

func get(url string) (response []byte, err error) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		sdkLog(err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		sdkLog(err.Error())
		return nil, err
	}
	return body, nil
}
func retry(url string, data interface{}, token string, attempts int, sleep time.Duration, fn func(string, interface{}, string) ([]byte, error)) ([]byte, error) {
	b, err := fn(url, data, token)
	if err != nil {
		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return retry(url, data, token, attempts, 2*sleep, fn)
		}
		return nil, err
	}
	return b, nil
}

//application/json; charset=utf-8
func Post2Api(url string, data interface{}, token string) (content []byte, err error) {
	if url == sendMsgRouter {
		return retry(url, data, token, 1, 10*time.Second, postLogic)
	} else {
		return postLogic(url, data, token)
	}
}

//application/json; charset=utf-8
func post2Api(url string, data interface{}, token string) (content []byte, err error) {
	sdkLog("call post2Api: ", url)

	if url == sendMsgRouter {
		return retry(url, data, token, 1, 10*time.Second, postLogic)
	} else {
		return postLogic(url, data, token)
	}
}

func post2ApiForRead(url string, data interface{}, token string) (content []byte, err error) {
	sdkLog("call post2Api: ", url)

	if url == sendMsgRouter {
		return retry(url, data, token, 3, 10*time.Second, postLogic)
	} else {
		return postLogic(url, data, token)
	}
}

func postLogic(url string, data interface{}, token string) (content []byte, err error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		sdkLog("marshal failed, url: ", url, "req: ", string(jsonStr))
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		sdkLog("newRequest failed, url: ", url, "req: ", string(jsonStr), err.Error())
		return nil, err
	}
	req.Close = true
	req.Header.Add("content-type", "application/json")
	req.Header.Add("token", token)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		sdkLog("client.Do failed, url: ", url, "req: ", string(jsonStr), err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		sdkLog("ioutil.ReadAll failed, url: ", url, "req: ", string(jsonStr), err.Error())
		return nil, err
	}
	return result, nil
}

func copyFile(srcName string, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}

	defer func() {
		if src != nil {
			src.Close()
		}
		if dst != nil {
			dst.Close()
		}
	}()
	return io.Copy(dst, src)
}

func fileTmpPath(fullPath string) string {
	suffix := path.Ext(fullPath)
	if len(suffix) == 0 {
		sdkLog("suffix  err:")
	}

	return SvrConf.DbDir + Md5(fullPath) + suffix //a->b
}

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//judge a string whether in the  string list
func isContain(target string, List []string) bool {

	for _, element := range List {

		if target == element {
			return true
		}
	}
	return false

}

func wrap(err error, message string) error {
	return errors.Wrap(err, "==> "+printCallerNameAndLine()+message)
}

func withMessage(err error, message string) error {
	return errors.WithMessage(err, "==> "+printCallerNameAndLine()+message)
}

func printCallerNameAndLine() string {
	pc, _, line, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name() + "()@" + strconv.Itoa(line) + ": "
}

func GetSwitchFromOptions(Options map[string]bool, key string) (result bool) {
	if flag, ok := Options[key]; !ok || flag {
		return true
	}
	return false
}
func SetSwitchFromOptions(Options map[string]bool, key string, value bool) {
	Options[key] = value
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, "==> "+printCallerNameAndLine()+message)
}

func WithMessage(err error, message string) error {
	return errors.WithMessage(err, "==> "+printCallerNameAndLine()+message)
}
