package main

import (

	"fmt"
	"open_im_sdk/open_im_sdk"

	"time"
)

//import (
//	"fmt"
//	_ "github.com/mattn/go-sqlite3"
//	"net/http"
//)

//	OnConnecting()
//	OnConnectSuccess()
//	OnConnectFailed(ErrCode int, ErrMsg string)
//	OnKickedOffline()
//	OnUserTokenExpired()
//	OnSelfInfoUpdated(info userInfo)

type InitSdk struct{}

func (i InitSdk) OnConnecting()                              {}
func (i InitSdk) OnConnectSuccess()                          {}
func (i InitSdk) OnConnectFailed(ErrCode int, ErrMsg string) {}
func (i InitSdk) OnKickedOffline()                           {}
func (i InitSdk) OnUserTokenExpired()                        {}
func (i InitSdk) OnSelfInfoUpdated(userInfo string)          {}
func (i InitSdk) OnError(errCode int, errMsg string)         {}
func (i InitSdk) OnSuccess(str string)                       {}

type Base interface {
	OnError(errCode int, errMsg string)
	OnSuccess(data string)
}

type X struct {


}
func (X)OnError(errCode int, errMsg string) {
	fmt.Println("OnError", errCode, errMsg)
}
func (X)OnSuccess(data string){
	fmt.Println("OnSuccess, ", data)
}

type SendMsgCallBack interface {
	Base
	OnProgress(progress int)
}

func (i InitSdk) OnProgress(progress int) {
	fmt.Printf("上传 %d / 100 \n", progress)
}

type groupApplication struct {
	GroupId          string `json:"groupID"`
	FromUser         string `json:"fromUserID"`
	FromUserNickName string `json:"fromUserNickName"`
	FromUserFaceUrl  string `json:"fromUserFaceUrl"`
	ToUser           string `json:"toUserID"`
	AddTime          int    `json:"addTime"`
	RequestMsg       string `json:"requestMsg"`
	HandledMsg       string `json:"handledMsg"`
	Type             int    `json:"type"`
	HandleStatus     int    `json:"handleStatus"`
	HandleResult     int    `json:"handleResult"`
}


func main() {

	var tk = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiI4MDFkMDk5NmE5MTUwNTJhIiwiUGxhdGZvcm0iOiJJT1MiLCJleHAiOjE2MjU0ODc1MzcsImlhdCI6MTYyNDg4MjczNywibmJmIjoxNjI0ODgyNzM3fQ.NTSI74nQDPrvkL8F8uhTpqtxLHCuCQpDHBbkjoe_cLQ"
	var uid = "801d0996a915052a1"

	open_im_sdk.Friend_uid = "c45baae51ab4a5d5"

	//open_im_sdk.DoTest(uid, tk)
	//open_im_sdk.DoTestDeleteFromFriendList()

	///func CreateVideoMessageFromFullPath(videoFullPath string, videoType string, duration int64, snapshotFullPath string) string {
	open_im_sdk.DoTest(uid, tk)
//	open_im_sdk.DoTestaddFriend()

//	open_im_sdk.DoTestSendMsg(open_im_sdk.Friend_uid)
	var xxx X
	open_im_sdk.Logout(xxx)
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
	for true {
		//	open_im_sdk.DoTestaddFriend()
		//	open_im_sdk.DoTestGetFriendList()
		//	open_im_sdk.DoTestDeleteFromFriendList()
		//	open_im_sdk.DoTestSetFriendInfo()
		//	open_im_sdk.DoTestCheckFriend()
		//	open_im_sdk.DoTestGetBlackList()
		//	open_im_sdk.DoTestDeleteFromBlackList()
		//	open_im_sdk.DoTestAddToBlackList()
		//	open_im_sdk.DoTestGetFriendsInfo()
		//	open_im_sdk.DoTestGetUsersInfo()
		//	open_im_sdk.DotestSetSelfInfo()
		//	open_im_sdk.DotestGetFriendApplicationList()
		//

		fmt.Println("running.................")
		time.Sleep(time.Duration(30) * time.Second)

	}
}
