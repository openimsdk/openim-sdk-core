package main

import (
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/test"
	"time"
)

func main() {

	strMyUidx := "3064833583"
	log.NewPrivateLog("", 6)
	tokenx := test.GenToken(strMyUidx)
	//tokenx := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiI3MDcwMDgxNTMiLCJQbGF0Zm9ybSI6IkFuZHJvaWQiLCJleHAiOjE5NjY0MTJ1XjJZGWj5fB3mqC7p6ytxSarvxZfsABwIjoxNjUxMDU1MDU2fQ.aWvmJ_sQxXmT5nKwiM5QsF9-tfkldzOYZtRD3nrUuko"
	test.InOutDoTest(strMyUidx, tokenx, test.WSADDR, test.APIADDR)
	test.DoSetGroupVerification()
	test.DoTestGetGroupsInfo()
	//	test.DoTestDeleteAllMsgFromLocalAndSvr()

	println("token ", tokenx)
	time.Sleep(100000 * time.Second)
	b := utils.GetCurrentTimestampBySecond()
	i := 0
	for {
		test.DoTestSendMsg2Group(strMyUidx, "a43619731c1c05eb93fc2501eab04f33", i)
		i++
		time.Sleep(100 * time.Millisecond)
		if i == 10000 {
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
