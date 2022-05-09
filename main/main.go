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
	TESTIP = "43.128.5.63"

	APIADDR      = "http://" + TESTIP + ":10002"
	WSADDR       = "ws://" + TESTIP + ":10001"
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
	strMyUidx := "17726378428"
	//friendID := "17726378428"
	tokenx := test.GenToken(strMyUidx)
	//tokenx :=    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiI3MDcwMDgxNTUiLCJQbGF0Zm9ybSI6IkFuZHJvaWQiLCJleHAiOjE5NjYzMTJ1XjJZGWj5fB3mqC7p6ytxSarvxZfsABwIjoxNjUwOTU1MDc5fQ.eLwd0meauHV8sBtR-MnZLkhVB9dFzU_g41Z5HI7U7YM"
	//	tokenx = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiIxNzcyNjM3ODQyOCIsIlBsYXRmb3JtIjoiSU9TIiwiZXhwIjoxOTYzMjE2NDU1LCJuYmYiOjE2NDc4NTY0NTUsImlhdCI6MTY0Nzg1NjQ1NX0.3fOcyhw7r5lOkRTJdDy7-tG9XC4XrKj_N7ufrGHPWYM"
	test.InOutDoTest(strMyUidx, tokenx, WSADDR, APIADDR)

	log.Info("", "DotestSetGroupMemberNickname start...")

	//test.TestGetWorkMomentsUnReadCount()
	//test.TestGetWorkMomentsNotification()
	//test.TestClearWorkMomentsNotification()
	test.DoTestSendImageMsg(strMyUidx, "13312341234")
	log.Info("", "test start...")
	//test.DoTestGetSubDepartment()
	//test.DoTestGetDepartmentMember()
	//test.DoTestGetUserInDepartment()
	//test.DoTestGetDepartmentMemberAndSubDepartment()
	//test.DotestSetGroupMemberNickname(strMyUidx)
	//test.DoTestDeleteAllMsgFromLocalAndSvr()
	//log.Warn("", "login ok, see memory, sleep 10s")
	//time.Sleep(2 * time.Second)
	//	test.InOutLogou()
	//	log.Warn("", "logout ok, see memory, sleep 10s")
	//	time.Sleep(10 * time.Second)
	//}

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
		//	test.DoTestSendMsg2(strMyUidx, "100908")
		time.Sleep(10 * time.Second)
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
