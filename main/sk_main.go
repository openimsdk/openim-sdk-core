package main

import (
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/test"
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

func main() {

	strMyUidx := "18666662345"
	log.NewPrivateLog("", 6)
	tokenx := test.GenToken(strMyUidx)
	//tokenx := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiI3MDcwMDgxNTMiLCJQbGF0Zm9ybSI6IkFuZHJvaWQiLCJleHAiOjE5NjY0MTJ1XjJZGWj5fB3mqC7p6ytxSarvxZfsABwIjoxNjUxMDU1MDU2fQ.aWvmJ_sQxXmT5nKwiM5QsF9-tfkldzOYZtRD3nrUuko"
	test.InOutDoTest(strMyUidx, tokenx, test.WSADDR, test.APIADDR)
	println("start")
	test.DoTestSendMsg2Group(strMyUidx, "0a3b3b8cc181b9fcfdfe4887e0b4ac30", 0)
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
	time.Sleep(250000 * time.Millisecond)
	b := utils.GetCurrentTimestampBySecond()
	i := 0
	for {
		test.DoTestSendMsg2Group(strMyUidx, "42c9f515cb84ee0e82b3f3ce71eb14d6", i)
		i++
		time.Sleep(250 * time.Millisecond)
		if i == 100 {
			break
		}
		log.Warn("", "10 * time.Millisecond ###################waiting... msg: ", i)
	}

	log.Warn("", "cost time: ", utils.GetCurrentTimestampBySecond()-b)
	return
	i = 0
	for {
		test.DoTestSendMsg2Group(strMyUidx, "42c9f515cb84ee0e82b3f3ce71eb14d6", i)
		i++
		time.Sleep(1000 * time.Millisecond)
		if i == 10 {
			break
		}
		log.Warn("", "1000 * time.Millisecond ###################waiting... msg: ", i)
	}

	i = 0
	for {
		test.DoTestSendMsg2Group(strMyUidx, "42c9f515cb84ee0e82b3f3ce71eb14d6", i)
		i++
		time.Sleep(10000 * time.Millisecond)
		log.Warn("", "10000 * time.Millisecond ###################waiting... msg: ", i)
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
