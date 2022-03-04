package main

import (
	"github.com/pkg/errors"
	"open_im_sdk/test"

	//"open_im_sdk/pkg/constant"

	"fmt"
	//"open_im_sdk/pkg/log"

	"open_im_sdk/pkg/utils"

	"time"
)

func main() {

	//	strMyUidx := "18381415165"
	//	friendID := "17726378428"
	//	tokenx := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiIxODM4MTQxNTE2NSIsIlBsYXRmb3JtIjoiSU9TIiwiZXhwIjoxOTYxMTI5NTQyLCJuYmYiOjE2NDU3Njk1NDIsImlhdCI6MTY0NTc2OTU0Mn0.hwLlECGDdaJscqGFLx-Avx6lbj3cNHSHq1QhdO8-zHg"
	//	test.InOutDoTest(strMyUidx, tokenx, WSADDR, APIADDR)
	//	test.DoTestCreateGroup()
	//test.SetTestFriendID(friendUserID)
	//fmt.Println("logout ........... ")
	//test.InOutLogou()
	//test.DoTestSetConversationRecvMessageOpt("17396220460", `["id1","id2"]`, constant.ConversationNotNotification)
	//	test.DoTestSetConversationRecvMessageOpt("openIM123456", `["single_openIM101", "id2"]`, constant.ConversationNotification)
	//	test.DoTestGetConversationRecvMessageOpt(`["single_openIM101"]`)
	//test.DoTestAddToBlackList()
	//test.DoTestGetDesignatedFriendsInfo()
	//	test.DoTestGetUsersInfo()
	//	test.DoTestCreateGroup()

	//	test.DoSetGroupInfo()
	//test.DotestGetGroupMemberList()

	//test.DotestMinio()
	//test.DotestKickGroupMember()
	//	test.DotestInviteUserToGroup()
	//test.DotestGetGroupApplicationList()

	//test.DotestAcceptGroupApplication("")
	//test.DoTestGetUserReqGroupApplicationList()
	//test.DoTestSetConversationRecvMessageOpt(strMyUidx, []string{"s", "s2"})
	//test.DoTestSetConversationStatus(strMyUidx, 2)
	//
	//test.DoTestGetRecvGroupApplicationList()
	///////////////friend///////////////////////////////////
	//	test.DoTestGetFriendApplicationList()
	//test.DoTestAcceptFriendApplication()

	//	test.DotestSetFriendRemark()
	//	test.DotestGetFriendList()

	//test.DotestDeleteFriend()
	//	test.DoTestDeleteFromBlackList()
	//test.DoTestAddFriend()

	//test.DotestSetFriendRemark()

	//open_im_sdk.DoTestGetFriendList()
	//	open_im_sdk.DoTestAddToBlackList()
	//	open_im_sdk.DoTestGetBlackList()
	//	open_im_sdk.DoTestDeleteFromBlackList()
	//DoTestGetDesignatedFriendsInfo()
	//test.DoTestSendMsg(strMyUidx, test.Friend_uid)
	//test.DoTestSendImageMsg("", test.Friend_uid)
	//for true {
	//	time.Sleep(time.Duration(100) * time.Second)
	//	//	test.DoTestSendMsg(strMyUidx, test.Friend_uid)
	//	fmt.Println("waiting")
	//}
	//	test.DoTestSendImageMsg("", test.Friend_uid)
	//	test.TestSendCostTime()

	testClientNum := 100
	intervalSleep := 1
	imIP := "43.128.5.63"
	test.DoTestRun(testClientNum, intervalSleep, imIP)

	i := 0
	for true {

		i++
		fmt.Println("DoTestSendMsg count: ", i)
		fmt.Println("waiting")
		time.Sleep(time.Duration(1000) * time.Second)
	}

}

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

type zx struct {
}

func (z zx) txexfc(uid int) int {
	utils.LogBegin(uid)
	if uid == 0 {
		return -1
		utils.LogFReturn(-1)
	}
	utils.LogSReturn(1)
	return 1

}

func authenticate(a int) error {
	if a == 0 {
		return errors.New("test error")
	}
	return nil
}

// Annotate error
//func AuthenticateRequest(a int) error {
//	err := authenticate(a)
//	if err != nil {
//		var v GetTokenReq
//		v.Platform = 100
//		//	return fmt.Errorf("authenticate failed: %v", err, v)
//		return fmt.Errorf("open file error: %w", err)
//	}
//	return nil
//}

// Better
func f3() error {
	return utils.Wrap(errors.New("first error"), " wrap")

}

func f2() error {
	err := f3()
	if err != nil {
		return utils.WithMessage(err, "f3 err")
	}
	return nil
}

func f1() error {
	err := f2()
	if err != nil {
		return utils.WithMessage(err, "f2 err")
	}
	return nil
}

//
//func Wrap(err error, message string) error {
//	return errors.Wrap(err, "==> "+printCallerNameAndLine()+message)
//}
//
//func WithMessage(err error, message string) error {
//	return errors.WithMessage(err, "==> "+printCallerNameAndLine()+message)
//}
//
//func printCallerNameAndLine() string {
//	pc, _, line, _ := runtime.Caller(2)
//	return runtime.FuncForPC(pc).Name() + "()@" + strconv.Itoa(line) + ": "
//}

// myuid,  maxuid,  msgnum
