package main

import (
	"fmt"
	"open_im_sdk/open_im_sdk"
	"time"
)

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

func main() {

	var tk = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiI3M2IwYzYzYmY2ZWZiYjkxIiwiUGxhdGZvcm0iOiJJT1MiLCJleHAiOjE2Mjc0NzU2MTYsImlhdCI6MTYyNjg3MDgxNiwibmJmIjoxNjI2ODcwODE2fQ.oVD0-_qjNckPMdBSfNcsDBLyPlLSnyqaz1T_jU91Pxw"
	var uid = "73b0c63bf6efbb91"

	//open_im_sdk.Friend_uid = ""

	///func CreateVideoMessageFromFullPath(videoFullPath string, videoType string, duration int64, snapshotFullPath string) string {
	open_im_sdk.DoTest(uid, tk)
	//	s := open_im_sdk.CreateSoundMessageFromFullPath("D:\\1.wav", 1)
	//	fmt.Println("ssss", s)
//	open_im_sdk.DoTestSendMsg("adaa5e370d7208b2")
	open_im_sdk.ForceReConn()
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

	time.Sleep(time.Duration(5) * time.Second)
	open_im_sdk.ForceReConn()

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

		//	open_im_sdk.DoTestDeleteFromFriendList()
		//	var xxx X
		//	open_im_sdk.Logout(xxx)
		//	fmt.Println("logouttttttttttttttttttttttt")
		//	open_im_sdk.Login(uid, tk, open_im_sdk.XtestLogin)
		fmt.Println("running.................")
		//		open_im_sdk.DotestGetJoinedGroupList()

		time.Sleep(time.Duration(30) * time.Second)
	}
}
