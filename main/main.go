package main

import (
	"C"
	"open_im_sdk/pkg/log"
	"open_im_sdk/test"
	"time"
)

//export reliabilityTest
func reliabilityTest() {
	intervalSleepMs := 1
	randSleepMaxSecond := 30
	imIP := "43.128.5.63"
	oneClientSendMsgNum := 5000
	testClientNum := 50
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

func main() {
	reliabilityTest()
}

//
//func main() {
//	testClientNum := 100
//	intervalSleep := 2
//	imIP := "43.128.5.63"
//	test.PressTest(testClientNum, intervalSleep, imIP)
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
