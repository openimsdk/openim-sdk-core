package main

import (
	"flag"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/test"
	"open_im_sdk/ws_wrapper/ws_local_server"

	"time"
)

//func reliabilityTest() {
//	intervalSleepMs := 1
//	randSleepMaxSecond := 30
//	imIP := "43.128.5.63"
//	oneClientSendMsgNum := 1
//	testClientNum := 100
//	test.ReliabilityTest(oneClientSendMsgNum, intervalSleepMs, imIP, randSleepMaxSecond, testClientNum)
//
//	for {
//		if test.CheckReliabilityResult() {
//			log.Warn("", "CheckReliabilityResult ok, again")
//		} else {
//			log.Warn("", "CheckReliabilityResult failed , wait.... ")
//		}
//		time.Sleep(time.Duration(10) * time.Second)
//	}
//}

//var (
//	TESTIP = "43.128.5.63"
//	APIADDR      = "http://" + TESTIP + ":10000"
//	WSADDR       = "ws://" + TESTIP + ":17778"
//	REGISTERADDR = APIADDR + "/user_register"
//	ACCOUNTCHECK = APIADDR+"/manager/account_check"
//	TOKENADDR    = APIADDR + "/auth/user_token"
//	SECRET       = "tuoyun"
//	SENDINTERVAL = 20
//)

type ChanMsg struct {
	data []byte
	uid  string
}

func testMem() {
	s := server_api_params.MsgData{}
	s.RecvID = "11111111sdfaaaaaaaaaaaaaaaaa11111"
	s.RecvID = "222222222afsddddddddddddddddddddddd22"
	s.ClientMsgID = "aaaaaaaaaaaadfsaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	s.SenderNickname = "asdfafdassssssssssssssssssssssfds"
	s.SenderFaceURL = "bbbbbbbbbbbbbbbbsfdaaaaaaaaaaaaaaaaaaaaaaaaa"

	ws_local_server.SendOneUserMessageForTest(s, "aaaa")
}

func main() {
	var onlineNum *int          //Number of online users
	var senderNum *int          //Number of users sending messages
	var singleSenderMsgNum *int //Number of single user send messages
	var intervalTime *int       //Sending time interval, in millseconds
	onlineNum = flag.Int("on", 10000, "online num")
	senderNum = flag.Int("sn", 100, "sender num")
	singleSenderMsgNum = flag.Int("mn", 1000, "single sender msg num")
	intervalTime = flag.Int("t", 1000, "interval time mill second")
	flag.Parse()
	//strMyUidx := "18666662412"
	//friendID := "17726378428"
	//tokenx := test.GenToken(strMyUidx)
	//tokenx := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiI3MDcwMDgxNTMiLCJQbGF0Zm9ybSI6IkFuZHJvaWQiLCJleHAiOjE5NjY0MTJ1XjJZGWj5fB3mqC7p6ytxSarvxZfsABwIjoxNjUxMDU1MDU2fQ.aWvmJ_sQxXmT5nKwiM5QsF9-tfkldzOYZtRD3nrUuko"
	//test.InOutDoTest(strMyUidx, tokenx, test.WSADDR, test.APIADDR)
	//test.DotestSearchFriends()
	if *senderNum == 0 {
		test.RegisterAccounts(*onlineNum)
		return
	}

	test.OnlineTest(*onlineNum)
	//test.TestSendCostTime()
	test.ReliabilityTest(*singleSenderMsgNum, *intervalTime, 10, *senderNum)
	//test.DoTestSearchLocalMessages()
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
		time.Sleep(30 * time.Second)
		log.Info("", "waiting...")
	}
	//reliabilityTest()
	//	test.PressTest(testClientNum, intervalSleep, imIP)
}

//
//func main() {
//	testClientNum := 100
//	intervalSleep := 2
//	imIP := "43.128.5.63"

//
//	msgNum := 1000
//	test.ReliabilityTest(msgNum, intervalSleep, imIP)
//	for i := 0; i < 6; i++ {
//		test.Msgwg.Wait()
//	}
//
//	for {
//
//		if test.CheckReliabilityResult() {
//			log.Warn("CheckReliabilityResult ok, again")
//
//		} else {
//			log.Warn("CheckReliabilityResult failed , wait.... ")
//		}
//
//		time.Sleep(time.Duration(10) * time.Second)
//	}
//
//}

//func printCallerNameAndLine() string {
//	pc, _, line, _ := runtime.Caller(2)
//	return runtime.FuncForPC(pc).Name() + "()@" + strconv.Itoa(line) + ": "
//}

// myuid,  maxuid,  msgnum
