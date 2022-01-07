package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"open_im_sdk/internal/controller/init"
	"open_im_sdk/internal/open_im_sdk"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	log2 "open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"

	"github.com/gorilla/websocket"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"

	"os"
	"path"

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

func (u *open_im_sdk.UserRelated) AddCh() (string, chan open_im_sdk.GeneralWsResp) {
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

func (u *open_im_sdk.UserRelated) GetCh(msgIncr string) chan open_im_sdk.GeneralWsResp {
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

func (u *open_im_sdk.UserRelated) DelCh(msgIncr string) {
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

func (u *open_im_sdk.UserRelated) sendPingMsg() error {
	u.stateMutex.Lock()
	defer u.stateMutex.Unlock()
	var ping string = "try ping"

	err := u.conn.SetWriteDeadline(time.Now().Add(8 * time.Second))
	if err != nil {
		sdkLog("SetWriteDeadline failed ", err.Error())
	}
	return u.conn.WriteMessage(websocket.PingMessage, []byte(ping))
}

func (u *open_im_sdk.UserRelated) writeBinaryMsg(msg open_im_sdk.GeneralWsReq) (error, *websocket.Conn) {
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
		if len(buff.Bytes()) > constant.MaxTotalMsgLen {
			LogFReturn("msg too long", len(buff.Bytes()), constant.MaxTotalMsgLen)
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

func (u *open_im_sdk.UserRelated) decodeBinaryWs(message []byte) (*open_im_sdk.GeneralWsResp, error) {
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

func (u *open_im_sdk.UserRelated) WriteMsg(msg open_im_sdk.GeneralWsReq) (error, *websocket.Conn) {
	LogStart(msg.OperationID)
	LogSReturn(msg.OperationID)
	return u.writeBinaryMsg(msg)
}

func notifyCh(ch chan open_im_sdk.GeneralWsResp, value open_im_sdk.GeneralWsResp, timeout int64) error {
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

func sendCmd(ch chan open_im_sdk.cmd2Value, value open_im_sdk.cmd2Value, timeout int64) error {
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

func GenMsgIncr(userID string) string {
	return userID + "_" + int64ToString(getCurrentTimestampByNano())
}

func structToJsonString(param interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

func structToJsonStringDefault(param interface{}) string {
	if reflect.TypeOf(param).Kind() == reflect.Slice && reflect.ValueOf(param).Len() == 0 {
		return "[]"
	}
	return structToJsonString(param)
}

//The incoming parameter must be a pointer
func jsonStringToStruct(s string, args interface{}) error {
	return Wrap(json.Unmarshal([]byte(s), args), "json Unmarshal failed")
}

//Convert timestamp to time.Time type

func UnixSecondToTime(second int64) time.Time {
	return time.Unix(second, 0)
}
func IntToString(i int) string {
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

func GetConversationIDBySessionType(sourceID string, sessionType int) string {
	switch sessionType {
	case constant.SingleChatType:
		return "single_" + sourceID
	case constant.GroupChatType:
		return "group_" + sourceID
	}
	return ""
}
func getIsRead(b bool) int {
	if b {
		return constant.HasRead
	} else {
		return constant.NotRead
	}
}
func getIsFilter(b bool) int {
	if b {
		return constant.IsFilter
	} else {
		return constant.NotFilter
	}
}
func getIsReadB(i int) bool {
	if i == constant.HasRead {
		return true
	} else {
		return false
	}

}

func RunFuncName() string {
	pc, _, _, _ := runtime.Caller(2)
	return cleanUpfuncName(runtime.FuncForPC(pc).Name())
}

func cleanUpfuncName(funcName string) string {
	end := strings.LastIndex(funcName, ".")
	if end == -1 {
		return ""
	}
	return funcName[end+1:]
}

func LogBegin(v ...interface{}) {
	//if constant.SdkLogFlag == 1 {
	//	return
	//}
	//if open_im_sdk.logger != nil {
	//	log2.NewInfo("", v...)
	//	return
	//}
	//pc, b, c, _ := runtime.Caller(1)
	//fname := runtime.FuncForPC(pc).Name()
	//i := strings.LastIndex(b, "/")
	//if i != -1 {
	//	sLog.Println(" [", b[i+1:len(b)], ":", c, "]", cleanUpfuncName(fname), "call func begin, args: ", v)
	//}
}

func LogEnd(v ...interface{}) {
	//if constant.SdkLogFlag == 1 {
	//	return
	//}
	//if open_im_sdk.logger != nil {
	//	log2.NewInfo("", v...)
	//	return
	//}
	//pc, b, c, _ := runtime.Caller(1)
	//fname := runtime.FuncForPC(pc).Name()
	//i := strings.LastIndex(b, "/")
	//if i != -1 {
	//	sLog.Println(" [", b[i+1:len(b)], ":", c, "]", cleanUpfuncName(fname), "call func end, args: ", v)
	//}
}

func LogStart(v ...interface{}) {
	//if constant.SdkLogFlag == 1 {
	//	return
	//}
	//if open_im_sdk.logger != nil {
	//	log2.NewInfo("", v...)
	//	return
	//}
	//pc, b, c, _ := runtime.Caller(1)
	//fname := runtime.FuncForPC(pc).Name()
	//i := strings.LastIndex(b, "/")
	//if i != -1 {
	//	sLog.Println(" [", b[i+1:len(b)], ":", c, "]", cleanUpfuncName(fname), "func start, args: ", v)
	//}
}

func LogFReturn(v ...interface{}) {
	//if constant.SdkLogFlag == 1 {
	//	return
	//}
	//if open_im_sdk.logger != nil {
	//	log2.NewInfo("", v...)
	//	return
	//}
	//pc, b, c, _ := runtime.Caller(1)
	//fname := runtime.FuncForPC(pc).Name()
	//i := strings.LastIndex(b, "/")
	//if i != -1 {
	//	sLog.Println("[", b[i+1:len(b)], ":", c, "]", cleanUpfuncName(fname), "failed return args(info): ", v)
	//}
}

func LogSReturn(v ...interface{}) {
	//if constant.SdkLogFlag == 1 {
	//	return
	//}
	//if open_im_sdk.logger != nil {
	//	log2.NewInfo("", v...)
	//	return
	//}
	//pc, b, c, _ := runtime.Caller(1)
	//fname := runtime.FuncForPC(pc).Name()
	//i := strings.LastIndex(b, "/")
	//if i != -1 {
	//	sLog.Println("[", b[i+1:len(b)], ":", c, "]", cleanUpfuncName(fname), "success return args(info): ", v)
	//}

}

func sdkLog(v ...interface{}) {
	//if constant.SdkLogFlag == 1 {
	//	return
	//}
	//if open_im_sdk.logger != nil {
	//	log2.NewInfo("", v...)
	//	return
	//}
	//_, b, c, _ := runtime.Caller(1)
	//i := strings.LastIndex(b, "/")
	//if i != -1 {
	//	sLog.Println("[", b[i+1:len(b)], ":", c, "]", v)
	//}

}

type LogInfo struct {
	Info string `json:"info"`
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

	return constant.SvrConf.DbDir + Md5(fullPath) + suffix //a->b
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

func GetSelfFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func printCallerNameAndLine() string {
	pc, _, line, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name() + "()@" + strconv.Itoa(line) + ": "
}
