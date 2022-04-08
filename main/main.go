package main

import (
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/test"
	"open_im_sdk/ws_wrapper/ws_local_server"

	"time"
)

func reliabilityTest() {
	intervalSleepMs := 1
	randSleepMaxSecond := 30
	imIP := "43.128.5.63"
	oneClientSendMsgNum := 1
	testClientNum := 100
	test.ReliabilityTest(oneClientSendMsgNum, intervalSleepMs, imIP, randSleepMaxSecond, testClientNum)

	for {
		if test.CheckReliabilityResult() {
			log.Warn("", "CheckReliabilityResult ok, again")
		} else {
			log.Warn("", "CheckReliabilityResult failed , wait.... ")
		}
		time.Sleep(time.Duration(10) * time.Second)
	}
}

var (
	TESTIP       = "api.adger.me"
	TESTIP_LOCAL = "api.adger.me"
	//TESTIP       = "1.14.194.38"
	APIADDR = "http://" + TESTIP_LOCAL + ":10000"
	//APIADDR = "https://im-api.jiarenapp.com"

	WSADDR = "ws://" + TESTIP + ":17778"
	//WSADDR = "wss://im.jiarenapp.com"

	REGISTERADDR = APIADDR + "/user_register"
	TOKENADDR    = APIADDR + "/auth/user_token"
	SECRET       = "tuoyun"
	SENDINTERVAL = 20
)

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

	test.REGISTERADDR = REGISTERADDR
	test.TOKENADDR = TOKENADDR
	test.SECRET = SECRET
	test.SENDINTERVAL = SENDINTERVAL
	strMyUidx := "1509479524083957760"
	//friendID := "17726378428"
	//tokenx := test.GenToken(strMyUidx)
	tokenx := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiIxNTA5NDc5NTI0MDgzOTU3NzYwIiwiUGxhdGZvcm0iOiJBbmRyb2lkIiwiZXhwIjoxOTY0NzQ0NDY1LCJuYmYiOjE2NDkzODQ0NjUsImlhdCI6MTY0OTM4NDQ2NX0.io-mLXoL4fiBCmV1KzBbKLBY8gz_oXfGgCrgzG5lCeA"
	//tokenx := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiIxNzcyNjM3ODQyOCIsIlBsYXRmb3JtIjoiSU9TIiwiZXhwIjoxOTYzMjE2NDU1LCJuYmYiOjE2NDc4NTY0NTUsImlhdCI6MTY0Nzg1NjQ1NX0.3fOcyhw7r5lOkRTJdDy7-tG9XC4XrKj_N7ufrGHPWYM"
	test.InOutDoTest(strMyUidx, tokenx, WSADDR, APIADDR)
	//test.DoTestDeleteAllMsgFromLocalAndSvr()
	//log.Warn("", "login ok, see memory, sleep 10s")
	//time.Sleep(2 * time.Second)
	//	test.InOutLogou()
	//	log.Warn("", "logout ok, see memory, sleep 10s")
	//	time.Sleep(10 * time.Second)
	//}
	//test.DoTestSetOneConversationPrivateChat("single_17726378428", false)
	test.DoTestSetConversationPinned("single_17396220460", true)
	//test.DoTestSendMsg2(strMyUidx, test.Friend_uid)
	//test.DoTestDeleteConversationMsgFromLocalAndSvr("single_17396220460")
	//test.I
	//test.DoTestInviteInGroup()
	//test.DoTestCancel()
	//test.DoTestSendMsg2(strMyUidx, friendID)
	//test.DoTestGetAllConversation()

	//test.DoTestGetOneConversation("17726378428")
	//test.DoTestGetConversations(`["single_17726378428"]`)
	//test.DoTestGetConversationListSplit()
	//test.DoTestGetConversationRecvMessageOpt(`["single_17899999999"]`)

	//set batch
	//test.DoTestSetConversationRecvMessageOpt([]string{"single_17396220460"}, constant.ReceiveNotNotifyMessage)
	//set one
	////set batch
	//test.DoTestSetConversationRecvMessageOpt([]string{"single_17726378428"}, constant.ReceiveMessage)
	////set one

	//test.DoTestSetOneConversationRecvMessageOpt("single_17726378428", constant.NotReceiveMessage)

	//test.DoTestReject()
	//test.DoTestAccept()
	//test.DoTestMarkGroupMessageAsRead()
	//test.DoTestGetGroupHistoryMessage()
	for {
		//test.DoTestSendMsg2(strMyUidx, test.Friend_uid)
		time.Sleep(1 * time.Second)
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
