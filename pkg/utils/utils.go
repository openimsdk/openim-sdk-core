// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/sdk_struct"
	"open_im_sdk/ws_wrapper/utils"
	"sort"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"

	"reflect"

	"github.com/pkg/errors"

	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func OperationIDGenerator() string {
	return strconv.FormatInt(time.Now().UnixNano()+int64(rand.Uint32()), 10)
}
func GetMsgID(sendID string) string {
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

// Get the current timestamp by Mill
func GetCurrentTimestampByMill() int64 {
	return time.Now().UnixNano() / 1e6
}

// Convert nano timestamp to time.Time type
func UnixNanoSecondToTime(nanoSecond int64) time.Time {
	return time.Unix(0, nanoSecond)
}

// Get the current timestamp by Nano
func GetCurrentTimestampByNano() int64 {
	return time.Now().UnixNano()
}

func StructToJsonString(param interface{}) string {
	dataType, err := json.Marshal(param)
	if err != nil {
		panic(err)
	}
	dataString := string(dataType)
	return dataString
}

func StructToJsonStringDefault(param interface{}) string {
	if reflect.TypeOf(param).Kind() == reflect.Slice && reflect.ValueOf(param).Len() == 0 {
		return "[]"
	}
	return StructToJsonString(param)
}

// The incoming parameter must be a pointer
func JsonStringToStruct(s string, args interface{}) error {
	return Wrap(json.Unmarshal([]byte(s), args), "json Unmarshal failed")
}
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
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

func RunFuncName() string {
	pc, _, _, _ := runtime.Caller(2)
	return CleanUpfuncName(runtime.FuncForPC(pc).Name())
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
	//	sLog.Println(" [", b[i+1:len(b)], ":", c, "]", cleanUpfuncName(fname), "call funcation begin, args: ", v)
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
	//	sLog.Println(" [", b[i+1:len(b)], ":", c, "]", cleanUpfuncName(fname), "call funcation end, args: ", v)
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
	//	sLog.Println(" [", b[i+1:len(b)], ":", c, "]", cleanUpfuncName(fname), "funcation start, args: ", v)
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

// judge a string whether in the  string list
func IsContain(target string, List []string) bool {

	for _, element := range List {

		if target == element {
			return true
		}
	}
	return false

}
func IsContainInt(target int, List []int) bool {

	for _, element := range List {

		if target == element {
			return true
		}
	}
	return false

}
func IsContainUInt32(target uint32, List []uint32) bool {

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
func Unwrap(err error) error {
	for err != nil {
		unwrap, ok := err.(interface {
			Unwrap() error
		})
		if !ok {
			break
		}
		err = unwrap.Unwrap()
	}
	return err
}

func WithMessage(err error, message string) error {
	return errors.WithMessage(err, "==> "+printCallerNameAndLine()+message)
}

func GetSelfFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return CleanUpfuncName(runtime.FuncForPC(pc).Name())
}
func CleanUpfuncName(funcName string) string {
	end := strings.LastIndex(funcName, ".")
	if end == -1 {
		return ""
	}
	return funcName[end+1:]
}

func printCallerNameAndLine() string {
	pc, _, line, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name() + "()@" + strconv.Itoa(line) + ": "
}
func StructToMap(user interface{}) map[string]interface{} {
	data, _ := json.Marshal(user)
	m := make(map[string]interface{})
	json.Unmarshal(data, &m)
	return m
}

//	funcation GetConversationIDBySessionType(sourceID string, sessionType int) string {
//		switch sessionType {
//		case constant.SingleChatType:
//			return "single_" + sourceID
//		case constant.GroupChatType:
//			return "group_" + sourceID
//		case constant.SuperGroupChatType:
//			return "super_group_" + sourceID
//		case constant.NotificationChatType:
//			return "notification_" + sourceID
//		}
//		return ""
//	}
func GetConversationIDByMsg(msg *sdk_struct.MsgStruct) string {
	switch msg.SessionType {
	case constant.SingleChatType:
		l := []string{msg.SendID, msg.RecvID}
		sort.Strings(l)
		return "si_" + strings.Join(l, "_") // single chat
	case constant.GroupChatType:
		return "g_" + msg.GroupID // group chat
	case constant.SuperGroupChatType:
		return "sg_" + msg.GroupID // super group chat
	case constant.NotificationChatType:
		return "sn_" + msg.SendID + "_" + msg.RecvID // server notification chat
	}
	return ""
}

func GetConversationIDByGroupID(groupID string) string {
	return "sg_" + groupID
}

func GetConversationTableName(conversationID string) string {
	return constant.ChatLogsTableNamePre + conversationID
}
func GetTableName(conversationID string) string {
	return constant.ChatLogsTableNamePre + conversationID
}
func GetErrTableName(conversationID string) string {
	return constant.SuperGroupErrChatLogsTableNamePre + conversationID
}
func RemoveRepeatedStringInList(slc []string) []string {
	var result []string
	tempMap := map[string]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}

/*
*
KMP
*
*/
func KMP(rMainString string, rSubString string) (isInMainString bool) {
	mainString := strings.ToLower(rMainString)
	subString := strings.ToLower(rSubString)
	mainIdx := 0
	subIdx := 0
	mainLen := len(mainString)
	subLen := len(subString)
	next := computeNextArray(subString)
	for {
		if mainIdx >= mainLen || subIdx >= subLen {
			break
		}

		if mainString[mainIdx] == subString[subIdx] {
			mainIdx++
			subIdx++
		} else {
			if subIdx != 0 {
				subIdx = next[subIdx-1]
			} else {
				mainIdx++
			}

		}
	}
	if subIdx >= subLen {
		if mainIdx-subLen >= 0 {
			return true
		}
	}
	return false

}

func computeNextArray(subString string) []int {
	next := make([]int, len(subString))
	index := 0
	i := 1
	for i < len(subString) {
		if subString[i] == subString[index] {
			next[i] = index + 1
			i++
			index++
		} else {
			if index != 0 {
				index = next[index-1]
			} else {
				i++
			}
		}
	}
	return next
}

func TrimStringList(list []string) (result []string) {
	for _, v := range list {
		if len(strings.Trim(v, " ")) != 0 {
			result = append(result, v)
		}
	}
	return result

}

// Get the intersection of two slices
func Intersect(slice1, slice2 []int64) []int64 {
	m := make(map[int64]bool)
	n := make([]int64, 0)
	for _, v := range slice1 {
		m[v] = true
	}
	for _, v := range slice2 {
		flag, _ := m[v]
		if flag {
			n = append(n, v)
		}
	}
	return n
}

// Get the diff of two slices
func DifferenceSubset(mainSlice, subSlice []int64) []int64 {
	m := make(map[int64]bool)
	n := make([]int64, 0)
	for _, v := range subSlice {
		m[v] = true
	}
	for _, v := range mainSlice {
		if !m[v] {
			n = append(n, v)
		}
	}
	return n
}
func DifferenceSubsetString(mainSlice, subSlice []string) []string {
	m := make(map[string]bool)
	n := make([]string, 0)
	for _, v := range subSlice {
		m[v] = true
	}
	for _, v := range mainSlice {
		if !m[v] {
			n = append(n, v)
		}
	}
	return n
}
func JsonDataOne(pb proto.Message) map[string]interface{} {
	return ProtoToMap(pb, false)
}

func ProtoToMap(pb proto.Message, idFix bool) map[string]interface{} {
	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: true,
	}

	s, _ := marshaler.MarshalToString(pb)
	out := make(map[string]interface{})
	json.Unmarshal([]byte(s), &out)
	if idFix {
		if _, ok := out["id"]; ok {
			out["_id"] = out["id"]
			delete(out, "id")
		}
	}
	return out
}
func GetUserIDForMinSeq(userID string) string {
	return "u_" + userID
}

func GetGroupIDForMinSeq(groupID string) string {
	return "g_" + groupID
}

func TimeStringToTime(timeString string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", timeString)
	return t, err
}

func TimeToString(t time.Time) string {
	return t.Format("2006-01-02")
}
func Uint32ListConvert(list []uint32) []int64 {
	var result []int64
	for _, v := range list {
		result = append(result, int64(v))
	}
	return result
}

func UnmarshalNotificationElem(bytes []byte, t interface{}) error {
	var n sdk_struct.NotificationElem
	err := utils.JsonStringToStruct(string(bytes), &n)
	if err != nil {
		return err
	}
	err = utils.JsonStringToStruct(n.Detail, t)
	if err != nil {
		return err
	}
	return nil
}
