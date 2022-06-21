package test

import (
	"encoding/json"
	"fmt"
	"open_im_sdk/internal/login"
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

func runRigister(strMyUid string) {
	for true {
		err := register(strMyUid)
		if err == nil {
			break
		} else {
			time.Sleep(time.Duration(1) * time.Second)
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
			continue
		} else {
			log.Info(req.OperationID, "Post2Api ok ", req)
			return nil
		}
	}
}

func getToken(uid string) string {
	url := TOKENADDR
	var req server_api_params.UserTokenReq
	req.Platform = 1
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
		log.Error(req.OperationID, "ErrCode failed ", stcResp.ErrCode, stcResp.ErrMsg, url, req)
		return ""
	}
	log.Info(req.OperationID, "get token: ", stcResp.Data.Token)
	return stcResp.Data.Token
}

func init() {
	sdk_struct.SvrConf = sdk_struct.IMConfig{Platform: 1, ApiAddr: APIADDR, WsAddr: WSADDR, DataDir: "./", LogLevel: 6, ObjectStorage: "cos"}
	allLoginMgr = make(map[int]*CoreNode)

}

func runGetToken(strMyUid string) string {
	var token string
	for true {
		token = getToken(strMyUid)
		if token == "" {

			time.Sleep(time.Duration(1) * time.Second)
			continue
		} else {
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

func GenUid(uid int, prefix string) string {
	if getMyIP() == "" {
		fmt.Println("getMyIP() failed")
		os.Exit(1)
	}
	UidPrefix := getMyIP() + "_" + prefix
	return UidPrefix + strconv.FormatInt(int64(uid), 10)
}

func GenToken(userID string) string {
	return runGetToken(userID)
}

func RegisterAccounts(number int) {
	var wg sync.WaitGroup
	wg.Add(number)
	for i := 0; i < number; i++ {
		go func(t int) {
			userID := GenUid(t, "online")
			register(userID)
			log.Info("register ", userID)
			wg.Done()
		}(i)

	}
	wg.Wait()
	log.Info("", "RegisterAccounts finish ", number)
}

func GenWsConn(id int) {
	userID := GenUid(id, "online")
	userLock.Lock()
	defer userLock.Unlock()
	allUserID = append(allUserID, userID)
	//register(userID)
	token := GenToken(userID)
	allToken = append(allToken, token)

	wsRespAsyn := interaction.NewWsRespAsyn()

	wsConn := interaction.NewWsConn(new(testInitLister), token, userID)
	cmdWsCh := make(chan common.Cmd2Value, 10)

	pushMsgAndMaxSeqCh := make(chan common.Cmd2Value, 1000)
	ws := interaction.NewWs(wsRespAsyn, wsConn, cmdWsCh, pushMsgAndMaxSeqCh, nil)
	allWs = append(allWs, ws)
}

func RegisterUserReliability(id int, timeStamp string) {
	userID := GenUid(id, "reliability"+timeStamp+"_")
	coreMgrLock.Lock()
	defer coreMgrLock.Unlock()
	register(userID)
	token := GenToken(userID)
	allLoginMgr[id] = &CoreNode{token: token, userID: userID}
}

type CoreNode struct {
	token  string
	userID string
	mgr    *login.LoginMgr
}

func addSendSuccess() {
	sendSuccessLock.Lock()
	defer sendSuccessLock.Unlock()
	sendSuccessCount++
}
func addSendFailed() {
	sendFailedLock.Lock()
	defer sendFailedLock.Unlock()
	sendFailedCount++
}

func OnlineTest(number int) {
	RegisterAccounts(number)
	var wg sync.WaitGroup
	wg.Add(number)
	for i := 0; i < number; i++ {
		go func(t int) {
			GenWsConn(t)
			log.Info("GenWsConn ", t)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Info("", "OnlineTest finish ", number)
}

func TestSendCostTime() {
	GenWsConn(0)
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
func TestSend(idx int, text string, uidNum, intervalSleep int) {
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
		time.Sleep(time.Duration(rand.Intn(intervalSleep)) * time.Millisecond)
	}
}
func testSendReliability(idx int, text string, uidNum, intervalSleep int) {
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
