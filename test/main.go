package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"open_im_sdk/open_im_sdk"
	"os"
	"strconv"
	"strings"
	"time"
)

//	var bb BaseSuccFailed
//	bb.OnSuccess("ddd")

//	var tk = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiI3M2IwYzYzYmY2ZWZiYjkxIiwiUGxhdGZvcm0iOiJJT1MiLCJleHAiOjE2Mjc0NzU2MTYsImlhdCI6MTYyNjg3MDgxNiwibmJmIjoxNjI2ODcwODE2fQ.oVD0-_qjNckPMdBSfNcsDBLyPlLSnyqaz1T_jU91Pxw"
//	var uid = "73b0c63bf6efbb91"

//	ws_local_server.Login(tk, uid)
//open_im_sdk.Friend_uid = ""

///func CreateVideoMessageFromFullPath(videoFullPath string, videoType string, duration int64, snapshotFullPath string) string {
//open_im_sdk.DoTest(uid, tk)
//	s := open_im_sdk.CreateSoundMessageFromFullPath("D:\\1.wav", 1)
//	fmt.Println("ssss", s)
//	open_im_sdk.DoTestSendMsg("adaa5e370d7208b2")
//open_im_sdk.ForceReConn()
//	open_im_sdk.DotestKickGroupMember()
//	open_im_sdk.DoJoinGroup()
//	open_im_sdk.DoTestCreateGroup()
//	open_im_sdk.DotestGetJoinedGroupList()
//open_im_sdk.DoJoinGroup()
//	open_im_sdk.DotesttestInviteUserToGroup()

//	open_im_sdk.DotestGetGroupMemberList()
//	open_im_sdk.DotestGetGroupMembersInfo()

//s := open_im_sdk.CreateImageMessageFromFullPath("C:\\xyz.jpg")
//open_im_sdk.SendMessage(xx, s, open_im_sdk.Friend_uid, "", false )

//
//s := open_im_sdk.CreateVideoMessageFromFullPath("D:\\22.mp4", "mp4", 58, "D:\\11.jpeg")

//	s  := open_im_sdk.CreateImageMessageFromFullPath(".//11.jpeg")
//	s := open_im_sdk.DoTestCreateImageMessage("11.jpeg")

//	time.Sleep(time.Duration(30) * time.Second)
//open_im_sdk.DoTestSendMsg(s)
//open_im_sdk.CreateImageMessage("11.jpeg")

//	open_im_sdk.DoJoinGroup()
//	open_im_sdk.DoTestSendMsg(open_im_sdk.Friend_uid)
//open_im_sdk.DoTestAcceptFriendApplicationdApplication()

//	open_im_sdk.DoTestDeleteFromFriendList()
//	open_im_sdk.DoTestRefuseFriendApplication()
//	open_im_sdk.DoTestAcceptFriendApplicationdApplication()
//	open_im_sdk.DoTestDeleteFromFriendList()
//open_im_sdk.DoTestDeleteFromFriendList()
//open_im_sdk.DoTestSendMsg(open_im_sdk.Friend_uid)
//open_im_sdk.DoTestMarkC2CMessageAsRead()
//"2021-06-23 12:25:36-7eefe8fc74afd7c6adae6d0bc76929e90074d5bc-8522589345510912161"
//	open_im_sdk.DoTestGetUsersInfo()

//open_im_sdk.DoTestGetFriendList()
//	open_im_sdk.DoTestGetHistoryMessage("c93bc8b171cce7b9d1befb389abfe52f")
//open_im_sdk.DoTestGetUsersInfo()
//open_im_sdk.DoTest(uid, tk)

//open_im_sdk.DoCreateGroup()
//open_im_sdk.DoSetGroupInfo()
//open_im_sdk.DoGetGroupsInfo()
//open_im_sdk.DoJoinGroup()
//open_im_sdk.DoQuitGroup()

//--------------------------------------
//var cc = open_im_sdk.IMConfig{
//	Platform:  1,
//	IpApiAddr: "http://47.112.160.66:10000",
//	IpWsAddr:  "47.112.160.66:7777",
//}
//b, _ := json.Marshal(cc)
//v1, v2, v3 := InitSdk{}, InitSdk{}, InitSdk{}
//open_im_sdk.InitSDK(string(b), v1)
//open_im_sdk.Login(uid, tk, v2)

// 转让群
//open_im_sdk.TransferGroupOwner("05dc84b52829e82242a710ecf999c72c", "uid_1234", v3)
//open_im_sdk.GetGroupApplicationList(v3)
//
//var sctApplication groupApplication
//sctApplication.GroupId = "05dc84b52829e82242a710ecf999c72c"
//sctApplication.FromUser = "61cd9ff3c88d64b42ff5ef930b9f007b"
//sctApplication.ToUser = "0"
//
//application, _ := json.Marshal(sctApplication)
//open_im_sdk.AcceptGroupApplication(string(application), "111", v3)
//open_im_sdk.RefuseGroupApplication(string(application), "111", v3)

//
//resp, _ := open_im_sdk.Upload("D:\\\\open-im-client-sdk\\test\\11.jpg", ss)
//
//fmt.Println(resp)
//
//resp, _ = open_im_sdk.Upload("D:\\\\open-im-client-sdk\\test\\11.jpg", ss)
//
//fmt.Println(resp)
//for {
//	time.Sleep(time.Second)
//	open_im_sdk.Login(uid, tk, v2)
//}

//open_im_sdk.upload("D:\\open-im-client-sdk\\test\\1.zip", &open_im_sdk.SelfListener{})
//open_im_sdk.Friend_uid = "355d8dcb9582b6f51b14dee7be83cc7987ab08d8"
//
//open_im_sdk.DoTest(uid, tk)
//open_im_sdk.DotestSetSelfInfo()
//open_im_sdk.DoTestGetUsersInfo()

//	time.Sleep(time.Duration(5) * time.Second)
//open_im_sdk.ForceReConn()open_im_sdk.LogBegin("")
//	myUid1 := 1
//	strMyUid1 := GenUid(myUid1)

//	runRigister(strMyUid1)
//	token1 := runGetToken(strMyUid1)
//	open_im_sdk.DoTest(strMyUid1, token1, WSADDR, APIADDR)
//	//recvId1 := GenUid(1)
//	//recvId1 := "18666662412"
//	/*
//		var i int64
//		for i = 0; i < 1; i++ {
//			time.Sleep(time.Duration(1) * time.Millisecond)
//			cont := "test data: 0->skkkkkkkkkkkkkkkkkk idx:" + strconv.FormatInt(i, 10)
//			open_im_sdk.DoTestSendMsg(strMyUid1, recvId1, cont)
//			fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~", i, "~~~~~~~~~~~~~~~~~~~~")
//		}
//	*/
//
//	//open_im_sdk.DoTestaddFriend()
//	for true {
//		time.Sleep(time.Duration(60) * time.Second)
//		fmt.Println("waiting")
//	}

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
	r, err := open_im_sdk.Post2Api(url, req, "")
	if err != nil {
		fmt.Println(r, err)
		return err
	}

	return nil

}
func getToken(uid string) string {
	url := TOKENADDR
	var req GetTokenReq
	req.Platform = 1
	req.Uid = uid
	req.Secret = SECRET
	r, err := open_im_sdk.Post2Api(url, req, "")
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

type zx struct {
}

func (z zx) txexfc(uid int) int {
	open_im_sdk.LogBegin(uid)
	if uid == 0 {
		return -1
		open_im_sdk.LogFReturn(-1)
	}
	open_im_sdk.LogSReturn(1)
	return 1

}
func GenUid(uid int) string {
	if uid > 1000 {
		return strconv.FormatInt(int64(uid), 10)
	}
	open_im_sdk.LogBegin(uid)

	if getMyIP() == "" {
		fmt.Println("getMyIP() failed")
		os.Exit(1)
	}
	UidPrefix := getMyIP() + "open_im_test_uid_"
	open_im_sdk.LogSReturn(UidPrefix + strconv.FormatInt(int64(uid), 10))
	return UidPrefix + strconv.FormatInt(int64(uid), 10)
}

func GetFileContentAsStringLines(filePath string) ([]string, error) {
	result := []string{}
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return result, err
	}
	s := string(b)
	for _, lineStr := range strings.Split(s, "\n") {
		lineStr = strings.TrimSpace(lineStr)
		if lineStr == "" {
			continue
		}
		result = append(result, lineStr)
	}
	return result, nil
}

func GetCmd(myUid int, filename string) int {
	cmd, err := GetFileContentAsStringLines("cmd.txt")
	if err != nil {
		fmt.Println("GetFileContentAsStringLines failed")
		return -1
	}
	if len(cmd) < myUid {
		fmt.Println("len failed")
		return -1
	}
	return int(open_im_sdk.StringToInt64(cmd[myUid-1]))
}

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

func runGetToken(strMyUid string) string {
	var token string
	for true {
		token = getToken(strMyUid)
		if token == "" {
			fmt.Println("test_openim: get token failed")
			time.Sleep(time.Duration(30) * time.Second)
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
				fmt.Println(ipnet.IP.String())
				return ipnet.IP.String()
			}

		}
	}
	return ""
}

var (
	APIADDR      = "http://121.37.25.71:10000"
	WSADDR       = "ws://121.37.25.71:17778"
	REGISTERADDR = "http://121.37.25.71:10000/auth/user_register"
	TOKENADDR    = "http://121.37.25.71:10000/auth/user_token"
	SECRET       = "tuoyun"
	SENDINTERVAL = 20
)

// myuid,  maxuid,  msgnum
func main() {

	myUid1 := 13333333333
	strMyUid1 := GenUid(myUid1)

	runRigister(strMyUid1)
	token1 := runGetToken(strMyUid1)
	open_im_sdk.InOutDoTest(strMyUid1, token1, WSADDR, APIADDR)
	time.Sleep(time.Duration(5) * time.Second)
	open_im_sdk.InOutDoTestSendMsg(strMyUid1, "18666662412")
	//open_im_sdk.InOutDoTest(strMyUid1, token1, WSADDR, APIADDR)
	open_im_sdk.InOutLogou()

	open_im_sdk.InOutDoTest(strMyUid1, token1, WSADDR, APIADDR)

	for true {
		time.Sleep(time.Duration(60) * time.Second)
		fmt.Println("waiting")
	}

	cmdfile := "./cmd.txt"
	uid := flag.Int("uid", 1, "RpcToken default listen port 10800")
	uidCount := flag.Int("uid_count", 2, "RpcToken default listen port 10800")
	messageCount := flag.Int("message_count", 1, "RpcToken default listen port 10800")
	APIADDR1 := flag.String("api_addr", "http://127.0.0.1:10000", "api addr")
	WSADDR1 := flag.String("ws_addr", "http://127.0.0.1:17778", "ws addr")
	REGISTERADDR1 := flag.String("register_addr", "http://127.0.0.1:10000/auth/user_register", "register addr")
	TOKENADDR1 := flag.String("token_addr", "http://127.0.0.1:10000/auth/user_token", "token addr")
	flag.Parse()

	APIADDR = *APIADDR1
	WSADDR = *WSADDR1
	REGISTERADDR = *REGISTERADDR1
	TOKENADDR = *TOKENADDR1

	var myUid int = *uid
	var uidNum int = *uidCount
	var msgnum int = *messageCount
	fmt.Println("args is ", myUid, uidNum, msgnum)
	var strMyUid string
	open_im_sdk.LogBegin()
	strMyUid = GenUid(myUid)

	runRigister(strMyUid)
	token := runGetToken(strMyUid)

	cmd := GetCmd(myUid, cmdfile)

	fmt.Println("getcmd value ", cmd)
	switch cmd {
	case -1:
		fmt.Println("GetCmd failed ")
		time.Sleep(time.Duration(1) * time.Second)
	case 5:
		fmt.Println("wait 2 mins, then login")
		time.Sleep(time.Duration(1*60) * time.Second)
		open_im_sdk.DoTest(strMyUid, token, WSADDR, APIADDR)
		fmt.Println("login do test, only login")
		fmt.Println("testmypid: ", os.Getpid())
	case 6:
		fmt.Println("wait 4 mins, then login")
		time.Sleep(time.Duration(2*60) * time.Second)
		open_im_sdk.DoTest(strMyUid, token, WSADDR, APIADDR)
		fmt.Println("login do test, only login")
		fmt.Println("testmypid: ", os.Getpid())
	case 3:
		fmt.Println("wait 2 mins, then login and send")
		time.Sleep(time.Duration(1*60) * time.Second)
		open_im_sdk.DoTest(strMyUid, token, WSADDR, APIADDR)
		fmt.Println("login do test, login and send")

		var recvId string
		var idx string
		rand.Seed(time.Now().UnixNano())
		if msgnum == 0 {
			fmt.Println("dont send,  exit")
			os.Exit(0)
		} else {
			for i := 0; i < msgnum; i++ {
				var r int
				for true {
					time.Sleep(time.Duration(SENDINTERVAL) * time.Millisecond)

					r = rand.Intn(uidNum) + 1
					fmt.Println("test rand ", myUid, uidNum, r)
					if r == myUid {
						continue
					} else {
						break
					}
				}
				recvId = GenUid(r)
				idx = strconv.FormatInt(int64(i), 10)

				open_im_sdk.DoTestSendMsg(strMyUid, recvId, idx)
			}
		}

	case 4:
		fmt.Println("wait 4 mins, then login and send")
		time.Sleep(time.Duration(2*60) * time.Second)
		open_im_sdk.DoTest(strMyUid, token, WSADDR, APIADDR)
		fmt.Println("login do test, login and send")

		var recvId string
		var idx string
		rand.Seed(time.Now().UnixNano())
		if msgnum == 0 {
			fmt.Println("dont send,  exit")
			os.Exit(0)
		} else {
			for i := 0; i < msgnum; i++ {
				var r int
				for true {
					time.Sleep(time.Duration(SENDINTERVAL) * time.Millisecond)

					r = rand.Intn(uidNum) + 1
					fmt.Println("test rand ", myUid, uidNum, r)
					if r == myUid {
						continue
					} else {
						break
					}
				}
				recvId = GenUid(r)
				idx = strconv.FormatInt(int64(i), 10)

				open_im_sdk.DoTestSendMsg(strMyUid, recvId, idx)
			}
		}

	case 1:
		fmt.Println("only login")
		open_im_sdk.DoTest(strMyUid, token, WSADDR, APIADDR)
		fmt.Println("login do test, only login...")
		fmt.Println("testmypid: ", os.Getpid())
	case 2:
		fmt.Println("login send")
		open_im_sdk.DoTest(strMyUid, token, WSADDR, APIADDR)
		fmt.Println("login do test, login and send")

		var recvId string
		var idx string
		rand.Seed(time.Now().UnixNano())
		if msgnum == 0 {
			fmt.Println("dont send,  exit")
			os.Exit(0)
		} else {
			for i := 0; i < msgnum; i++ {
				var r int
				for true {
					time.Sleep(time.Duration(SENDINTERVAL) * time.Millisecond)

					r = rand.Intn(uidNum) + 1
					fmt.Println("test rand ", myUid, uidNum, r)
					if r == myUid {
						continue
					} else {
						break
					}
				}
				recvId = GenUid(r)
				idx = strconv.FormatInt(int64(i), 10)

				open_im_sdk.DoTestSendMsg(strMyUid, recvId, idx)
			}
		}
	case 7:
		fmt.Println("random sleep and send")
		open_im_sdk.DoTest(strMyUid, token, WSADDR, APIADDR)

		var recvId string
		var idx string
		rand.Seed(time.Now().UnixNano())
		maxSleep := 60
		msgnum = 10
		if msgnum == 0 {
			fmt.Println("dont send,  exit")
			os.Exit(0)
		} else {
			for i := 0; i < msgnum; i++ {
				var r int
				for true {
					time.Sleep(time.Duration(rand.Intn(maxSleep)+1) * time.Second)
					r = rand.Intn(uidNum) + 1
					fmt.Println("test rand ", myUid, uidNum, r)
					if r == myUid {
						continue
					} else {
						break
					}
				}
				recvId = GenUid(r)
				idx = strconv.FormatInt(int64(i), 10)

				open_im_sdk.DoTestSendMsg(strMyUid, recvId, idx)
			}
		}

	}

	for true {
		time.Sleep(time.Duration(60) * time.Second)
		fmt.Println("waiting")
	}

}
