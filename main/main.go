package main

import (
	"fmt"
	"open_im_sdk/pkg/log"
	"open_im_sdk/test"
	"time"
)

func main() {
	//APIADDR := "https://test-web.rentsoft.cn/api"
	//WSADDR := "wss://test-web.rentsoft.cn/msg_gateway"
	APIADDR := "http://59.36.173.89:10002"
	WSADDR := "ws://59.36.173.89:10001"
	REGISTERADDR := APIADDR + "/user_register"
	ACCOUNTCHECK := APIADDR + "/manager/account_check"
	TOKENADDR := APIADDR + "/auth/user_token"
	SECRET := "openIM123"
	SENDINTERVAL := 20
	test.REGISTERADDR = REGISTERADDR
	test.TOKENADDR = TOKENADDR
	test.SECRET = SECRET
	test.SENDINTERVAL = SENDINTERVAL
	test.WSADDR = WSADDR
	test.ACCOUNTCHECK = ACCOUNTCHECK
	strMyUidx := "1010930629"

	//	var onlineNum *int          //Number of online users
	//	var senderNum *int          //Number of users sending messages
	//	var singleSenderMsgNum *int //Number of single user send messages
	//	var intervalTime *int       //Sending time interval, in millseconds
	//	onlineNum = flag.Int("on", 10000, "online num")
	//	senderNum = flag.Int("sn", 100, "sender num")
	//	singleSenderMsgNum = flag.Int("mn", 1000, "single sender msg num")
	//	intervalTime = flag.Int("t", 1000, "interval time mill second")
	//	flag.Parse()
	//strMyUidx := "13900000000"

	log.NewPrivateLog("", 6)
	tokenx := test.RunGetToken(strMyUidx)
	fmt.Println(tokenx)
	test.InOutDoTest(strMyUidx, tokenx, test.WSADDR, test.APIADDR)
	println("start")
	//test.DoTestCreateGroup()
	//test.DoTestSearchLocalMessages()
	// test.DoTestInviteInGroup()
	//time.Sleep(time.Second*6)
	//test.DoTestGetSelfUserInfo()
	//test.DoTestSetBurnDuration("single_2861383134")
	//test.DoTestInvite("123123")
	test.DoTestSetAppBackgroundStatus(true)
	for {
		//test.DotestDeleteFriend()
		//test.DoTestSendMsg2("", "1443506268")
		time.Sleep(time.Second * 10)
	}
	//test.DoTestSignalGetRoomByGroupID("1826384574")
	//test.DoTestSignalGetTokenByRoomID("1826384574")
	//test.DoTestSendImageMsg("3433303585")
	//test.DoTestGetUserInDepartment()
	//test.DoTestGetDepartmentMemberAndSubDepartment()
	//test.DoTestDeleteAllMsgFromLocalAndSvr()
	//	test.DoTestGetDepartmentMemberAndSubDepartment()
	//test.DotestUploadFile()
	//test.DotestMinio()
	//test.DotestSearchFriends()
	//if *senderNum == 0 {
	//	test.RegisterAccounts(*onlineNum)
	//	return
	//}
	//
	//test.OnlineTest(*onlineNum)
	////test.TestSendCostTime()
	//test.ReliabilityTest(*singleSenderMsgNum, *intervalTime, 10, *senderNum)
	//test.DoTestSearchLocalMessages()
	//println("start")
	//test.DoTestSendImageMsg(strMyUidx, "17726378428")
	//test.DoTestSearchGroups()
	//test.DoTestGetHistoryMessage("")
	//test.DoTestGetHistoryMessageReverse("")
	//test.DoTestInviteInGroup()
	//test.DoTestCancel()
	//test.DoTestSendMsg2(strMyUidx, friendID)
	//test.DoTestGetAllConversation()

	//test.DoTestGetOneConversation("17726378428")
	//test.DoTestGetConversations(`["single_17726378428"]`)
	//test.DoTestGetConversationListSplit()
	//test.DoTestGetConversationRecvMessageOpt(`["single_17726378428"]`)

	//set batch
	//test.DoTestSetConversationRecvMessageOpt([]string{"single_17726378428"}, constant.NotReceiveMessage)
	//set one
	////set batch
	//test.DoTestSetConversationRecvMessageOpt([]string{"single_17726378428"}, constant.ReceiveMessage)
	////set one
	//test.DoTestSetConversationPinned("single_17726378428", false)
	//test.DoTestSetOneConversationRecvMessageOpt("single_17726378428", constant.NotReceiveMessage)
	//test.DoTestSetOneConversationPrivateChat("single_17726378428", false)
	//test.DoTestReject()
	//test.DoTestAccept()
	//test.DoTestMarkGroupMessageAsRead()
	//test.DoTestGetGroupHistoryMessage()
	//test.DoTestGetHistoryMessage("17396220460")
	for {
		time.Sleep(10000 * time.Millisecond)
		log.Warn("", "10000 * time.Millisecond ###################waiting... msg: ")
	}

	//reliabilityTest()
	//	test.PressTest(testClientNum, intervalSleep, imIP)
}
