package test

import (
	"encoding/json"
	"fmt"
	"open_im_sdk/sdk_struct"
	"strings"
	"sync"

	//"github.com/gorilla/websocket"
	//"github.com/jinzhu/copier"
	//"google.golang.org/protobuf/types/known/apipb"
	"math/rand"
	"net"
	"open_im_sdk/internal/interaction"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/network"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"os"
	"strconv"
	"time"
)

var (
	TESTIP       = "43.128.5.63"
	APIADDR      = "http://" + TESTIP + ":10000"
	WSADDR       = "ws://" + TESTIP + ":17778"
	REGISTERADDR = APIADDR + "/auth/user_register"
	TOKENADDR    = APIADDR + "/auth/user_token"
	SECRET       = "tuoyun"
	SENDINTERVAL = 20
)

func runRigister(strMyUid string) {
	for true {
		err := register(strMyUid)
		if err == nil {
			break
		} else {
			time.Sleep(time.Duration(30) * time.Second)
			continue
		}
	}
}

type GetTokenReq struct {
	Secret   string `json:"secret"`
	Platform int    `json:"platform"`
	Uid      string `json:"uid"`
}

type RegisterReq struct {
	Secret   string `json:"secret"`
	Platform int    `json:"platform"`
	Uid      string `json:"uid"`
	Name     string `json:"name"`
}

type ResToken struct {
	Data struct {
		ExpiredTime int64  `json:"expiredTime"`
		Token       string `json:"token"`
		Uid         string `json:"uid"`
	}
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

func register(uid string) error {
	url := REGISTERADDR
	var req server_api_params.UserRegisterReq
	req.OperationID = utils.OperationIDGenerator()
	req.Platform = 1
	req.UserID = uid
	req.Secret = SECRET
	req.Nickname = uid
	for {
		_, err := network.Post2Api(url, req, "")
		if err != nil && !strings.Contains(err.Error(), "status code failed") {
			log.Error(req.OperationID, "post failed ,continue ", err.Error())
			//	time.Sleep(time.Duration(1) * time.Second)
			continue
		} else {
			return nil
		}
		//status code failed
	}

	return nil
}

func getToken(uid string) string {
	url := TOKENADDR
	var req server_api_params.UserTokenReq
	req.Platform = 2
	req.UserID = uid
	req.Secret = SECRET
	req.OperationID = utils.OperationIDGenerator()
	r, err := network.Post2Api(url, req, "")
	if err != nil {
		log.Error(req.OperationID, "Post2Api failed ", err.Error(), url, req)
		return ""
	}

	var stcResp ResToken
	err = json.Unmarshal(r, &stcResp)
	if stcResp.ErrCode != 0 {
		log.Error(req.OperationID, "ErrCode failed ", stcResp.ErrMsg, stcResp.ErrMsg)
		return ""
	}
	return stcResp.Data.Token

}

func init() {
	sdk_struct.SvrConf = sdk_struct.IMConfig{Platform: 1, ApiAddr: APIADDR, WsAddr: WSADDR, DataDir: "./", LogLevel: 6, ObjectStorage: "cos"}
}

func runGetToken(strMyUid string) string {
	var token string
	for true {
		token = getToken(strMyUid)
		if token == "" {
			log.Error("test_openim: get token failed")
			time.Sleep(time.Duration(1) * time.Second)
			continue
		} else {
			log.Info("get token: ", strMyUid, token)
			break
		}
	}

	return token
}

func getMyIP() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)

		os.Exit(1)
		return ""
	}
	for _, address := range addrs {

		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//	fmt.Println(ipnet.IP.String())
				return ipnet.IP.String()
			}

		}
	}
	return ""
}

func GenUid(uid int) string {
	if getMyIP() == "" {
		fmt.Println("getMyIP() failed")
		os.Exit(1)
	}
	UidPrefix := getMyIP() + "uid"
	return UidPrefix + strconv.FormatInt(int64(uid), 10)
}

func GenToken(userID string) string {

	return runGetToken(userID)
}

func GenWs(id int) {
	//return
	userID := GenUid(id)
	userLock.Lock()
	defer userLock.Unlock()
	allUserID = append(allUserID, userID)
	//register(userID)
	//token := GenToken(userID)
	token := "testtokentokentokentokentokentokentokentokentokentoken"
	allToken = append(allToken, token)

	wsRespAsyn := interaction.NewWsRespAsyn()

	wsConn := interaction.NewWsConn(new(testInitLister), token, userID)
	cmdWsCh := make(chan common.Cmd2Value, 10)

	pushMsgAndMaxSeqCh := make(chan common.Cmd2Value, 1000)
	ws := interaction.NewWs(wsRespAsyn, wsConn, cmdWsCh, pushMsgAndMaxSeqCh)
	allWs = append(allWs, ws)

}

var userLock sync.RWMutex

var allUserID []string
var allToken []string
var allWs []*interaction.Ws
var intervalSleep int

func DoTestRun(num int, interval int, ip string) {
	TESTIP = ip
	intervalSleep = interval
	var wg sync.WaitGroup
	wg.Add(num)

	for i := 0; i < num; i++ {
		go func(t int) {

			GenWs(t)
			log.Info("genws ", t)
			wg.Done()
		}(i)

		//if allUserID[i] == "" || allToken[i] == "" || allWs[i] == nil {
		//	log.Error("", "args failed")
		//}
		//log.Debug("", "user: ", allUserID[i], "token: ", allToken[i], allWs[i])
	}

	for i := 0; i < num; i++ {
		wg.Wait()
	}

	log.Info("", "start send message...")
	time.Sleep(time.Duration(1) * time.Second)

	for i := 0; i < num; i++ {
		go testSend(i, "ok", num)
	}
}

func TestSendCostTime() {
	GenWs(0)
	sendID := allUserID[0]
	recvID := allUserID[0]
	for {
		operationID := utils.OperationIDGenerator()
		b := SendTextMessage("test", sendID, recvID, operationID, allWs[0])
		if b {
			log.Debug(operationID, sendID, recvID, "SendTextMessage success")
		} else {
			log.Error(operationID, sendID, recvID, "SendTextMessage failed")
		}
		time.Sleep(time.Duration(5) * time.Second)
		log.Debug(operationID, "//////////////////////////////////")
	}

}
func testSend(idx int, text string, uidNum int) {
	for {
		operationID := utils.OperationIDGenerator()
		sendID := allUserID[idx]
		recvID := allUserID[rand.Intn(uidNum)]
		b := SendTextMessage(text, sendID, recvID, operationID, allWs[idx])
		if b {
			log.Debug(operationID, sendID, recvID, "SendTextMessage success")
		} else {
			log.Error(operationID, sendID, recvID, "SendTextMessage failed")
		}
		time.Sleep(time.Duration(rand.Intn(intervalSleep)) * time.Second)
	}
}

func SendTextMessage(text, senderID, recvID, operationID string, ws *interaction.Ws) bool {
	var wsMsgData server_api_params.MsgData
	options := make(map[string]bool, 2)
	wsMsgData.SendID = senderID
	wsMsgData.RecvID = recvID
	wsMsgData.ClientMsgID = utils.GetMsgID(senderID)
	wsMsgData.SenderPlatformID = 1
	wsMsgData.SessionType = constant.SingleChatType
	wsMsgData.MsgFrom = constant.UserMsgType
	wsMsgData.ContentType = constant.Text
	wsMsgData.Content = []byte(text)
	wsMsgData.CreateTime = utils.GetCurrentTimestampByMill()
	wsMsgData.Options = options
	wsMsgData.OfflinePushInfo = nil
	timeout := 300
	return ws.SendReqTest(&wsMsgData, constant.WSSendMsg, timeout, senderID, operationID)
}
