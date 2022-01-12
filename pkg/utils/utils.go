package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"

	"open_im_sdk/pkg/constant"

	"github.com/pkg/errors"
	"io"

	"reflect"

	"os"
	"path"

	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"
)


func OperationIDGenerator() string {
	return strconv.FormatInt(time.Now().UnixNano()+int64(rand.Uint32()), 10)
}
func getMsgID(sendID string) string {
	t := Int64ToString(GetCurrentTimestampByNano())
	return Md5(t + sendID + Int64ToString(rand.Int63n(GetCurrentTimestampByNano())))
}
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	cipher := h.Sum(nil)
	return hex.EncodeToString(cipher)
}

//Get the current timestamp by Second

func GetCurrentTimestampBySecond() int64 {
	return time.Now().Unix()
}

//Get the current timestamp by Mill
func GetCurrentTimestampByMill() int64 {
	return time.Now().UnixNano() / 1e6
}

//Convert nano timestamp to time.Time type
func UnixNanoSecondToTime(nanoSecond int64) time.Time {
	return time.Unix(0, nanoSecond)
}

//Get the current timestamp by Nano
func GetCurrentTimestampByNano() int64 {
	return time.Now().UnixNano()
}

func StructToJsonString(param interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

func StructToJsonStringDefault(param interface{}) string {
	if reflect.TypeOf(param).Kind() == reflect.Slice && reflect.ValueOf(param).Len() == 0 {
		return "[]"
	}
	return StructToJsonString(param)
}

//The incoming parameter must be a pointer
func JsonStringToStruct(s string, args interface{}) error {
	return Wrap(json.Unmarshal([]byte(s), args), "json Unmarshal failed")
}

//Convert timestamp to time.Time type

func UnixSecondToTime(second int64) time.Time {
	return time.Unix(second, 0)
}
func IntToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}
func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
func StringToInt64(i string) int64 {
	j, _ := strconv.ParseInt(i, 10, 64)
	return j
}

func StringToInt(i string) int {
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
func IsContain(target string, List []string) bool {

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
