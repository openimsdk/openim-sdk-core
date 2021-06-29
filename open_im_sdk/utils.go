package open_im_sdk

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	sLog "log"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func operationIDGenerator() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
func getMsgID(sendID string) string {
	t := time.Now().Format("2006-01-02 15:04:05") + "client"
	return Md5(t + "-" + sendID + "-" + strconv.Itoa(rand.Int()))
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

func structToJsonString(param interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

//The incoming parameter must be a pointer
func jsonStringToStruct(s string, args interface{}) error {
	err := json.Unmarshal([]byte(s), args)
	return err
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

func (ap applyUserInfo) Value() interface{} {
	return ap
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
		return MessageHasRead
	} else {
		return MessageHasNotRead
	}
}

func sdkLog(v ...interface{}) {
	_, b, c, _ := runtime.Caller(1)
	i := strings.LastIndex(b, "/")
	if i != -1 {
		sLog.Println("[", b[i+1:len(b)], ":", c, "]", v)
	}
}

type LogInfo struct {
	Info string `json:"info"`
}

func log(info string) error {
	sdkLog(info)
	//	fmt.Println(info)
	return nil
	url := "http://47.112.160.66:6969/log/log"
	_, err := post2Api(url, LogInfo{Info: info}, token)
	if err != nil {
		return err
	}
	return nil
}

func get(url string) (response []byte, err error) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log(err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log(err.Error())
		return nil, err
	}
	return body, nil
}

//application/json; charset=utf-8
func post2Api(url string, data interface{}, token string) (content []byte, err error) {
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
	req.Header.Add("content-type", "application/json")
	req.Header.Add("token", token)
	defer req.Body.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		sdkLog("newRequest failed, url: ", url, "req: ", string(jsonStr), err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		sdkLog("newRequest failed, url: ", url, "req: ", string(jsonStr), err.Error())
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
	return SvrConf.DbDir + Md5(fullPath) //a->b
}

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

var MyIP = ""

func getMyIp() string {
	//fixme In the configuration file, ip takes precedence, if not, get the valid network card ip of the machine

	//fixme Get the ip of the local network card
	netInterfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(netInterfaces); i++ {
		//Exclude useless network cards by judging the net.flag Up flag
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			address, _ := netInterfaces[i].Addrs()
			for _, addr := range address {
				if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						MyIP = ipNet.IP.String()
						return MyIP
					}
				}
			}
		}
	}
	return MyIP
}
