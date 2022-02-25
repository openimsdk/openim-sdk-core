package test

import (
	"encoding/json"
	"fmt"
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
	TESTIP = "43.128.5.63"
	//TESTIP       = "1.14.194.38"
	APIADDR      = "http://" + TESTIP + ":10000"
	WSADDR       = "ws://" + TESTIP + ":17778"
	REGISTERADDR = APIADDR + "/user_register"
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
	var req RegisterReq
	req.Platform = 1
	req.Uid = uid
	req.Secret = SECRET
	req.Name = uid
	r, err := network.Post2Api(url, req, "")
	if err != nil {
		fmt.Println(r, err)
		return err
	}

	return nil

}
func getToken(uid string) string {
	url := TOKENADDR
	var req GetTokenReq
	req.Platform = 2
	req.Uid = uid
	req.Secret = SECRET
	r, err := network.Post2Api(url, req, "")
	if err != nil {
		fmt.Println(r, err)
		return ""
	}

	var stcResp ResToken
	err = json.Unmarshal(r, &stcResp)
	if stcResp.ErrCode != 0 {
		fmt.Println(stcResp.ErrCode, stcResp.ErrMsg)
		return ""
	}
	return stcResp.Data.Token

}

func runGetToken(strMyUid string) string {
	var token string
	for true {
		token = getToken(strMyUid)
		if token == "" {
			fmt.Println("test_openim: get token failed")
			time.Sleep(time.Duration(30) * time.Second)
			continue
		} else {
			fmt.Println("get token: ", strMyUid, token)
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
				fmt.Println(ipnet.IP.String())
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
	UidPrefix := getMyIP() + "open_im_test_uid_"
	return UidPrefix + strconv.FormatInt(int64(uid), 10)
}

func GenToken(userID string) string {
	return runGetToken(userID)
}

func GenWs(id int) {
	userID := GenUid(id)
	allUserID = append(allUserID, userID)
	token := GenToken(userID)
	allToken = append(allToken, token)

	wsRespAsyn := interaction.NewWsRespAsyn()

	wsConn := interaction.NewWsConn(new(testInitLister), token, userID)
	cmdWsCh := make(chan common.Cmd2Value, 10)

	pushMsgAndMaxSeqCh := make(chan common.Cmd2Value, 1000)
	ws := interaction.NewWs(wsRespAsyn, wsConn, cmdWsCh, pushMsgAndMaxSeqCh)
	allWs = append(allWs, ws)

}

var allUserID []string
var allToken []string
var allWs []*interaction.Ws

func DoTestRun(num int) {
	for i := 0; i < num; i++ {
		GenWs(i)
		if allUserID[i] == "" || allToken[i] == "" || allWs[i] == nil {
			log.Error("", "args failed")
		}
		log.Debug("", "user: ", allUserID[i], "token: ", allToken[i], allWs[i])
	}

	for i := 0; i < num; i++ {
		go testSend(i, "ok", num)
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
			log.Debug(operationID, sendID, recvID, "SendTextMessage failed")
		}

		time.Sleep(time.Duration(100) * time.Second)
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
