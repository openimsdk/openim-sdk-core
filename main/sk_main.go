package main

import (
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/test"
	"time"
)

var allDB []*db.DataBase

//func TestDB(loginUserID string) {
//	operationID := utils.OperationIDGenerator()
//	dbUser, err := db.NewDataBase(loginUserID, "/data/test/Open-IM-Server/db/sdk/", operationID)
//	if err != nil {
//		log.Error(operationID, "NewDataBase failed ", err.Error(), loginUserID)
//		return
//	}
//	conversationList, err := dbUser.GetAllConversationList()
//	if err != nil {
//		log.Error(operationID, "GetAllConversationList failed ", err.Error())
//	}
//	log.Info(operationID, "GetAllConversationList len: ", len(conversationList))
//
//	groupIDList, err := dbUser.GetJoinedGroupList()
//	if err != nil {
//		log.Error(operationID, "GetJoinedGroupList failed ", err.Error())
//	}
//	log.Info(operationID, "GetJoinedGroupList len: ", len(groupIDList))
//
//	groupMemberList, err := dbUser.GetAllGroupMemberList()
//	if err != nil {
//		log.Error(operationID, "GetAllGroupMemberList failed ", err.Error())
//	}
//	log.Info(operationID, "GetAllGroupMemberList len: ", len(groupMemberList))
//	//GetAllMessageForTest
//	msgList, err := dbUser.GetAllMessageForTest()
//	if err != nil {
//		log.Error(operationID, "GetAllMessageForTest failed ", err.Error())
//	}
//	log.Info(operationID, "GetAllMessageForTest len: ", len(msgList))
//	allDB = append(allDB, dbUser)
//
//	dbUser.CloseDB(operationID)
//	log.Info(operationID, "close db finished ")
//
//}

func main() {
	//var userIDList []string
	//f, err := os.Open("/data/test/Open-IM-Server/db/sdk")
	//if err != nil {
	//	log.Error("", "open failed ", err.Error())
	//	return
	//}
	//files, err := f.Readdir(-1)
	//f.Close()
	//if err != nil {
	//	log.Error("", "Readdir failed ", err.Error())
	//	return
	//}
	//
	//for _, file := range files {
	//	begin := strings.Index(file.Name(), "OpenIM_v2_")
	//	end := strings.Index(file.Name(), ".db")
	//	userID := file.Name()[begin+len("OpenIM_v2_") : end]
	//	// OpenIM_v2_3380999461.db
	//	log.Info("", "file name: ", file.Name(), userID)
	//	TestDB(userID)
	//}
	//log.Info("", "files: ", len(allDB))
	////for _, v := range allDB {
	////	v.CloseDB("aa")
	////}
	//
	//log.Info("", "gc begin ")
	//runtime.GC()
	//log.Info("", "gc end ")
	//time.Sleep(100000 * time.Second)
	//return
	strMyUidx := "3370431052"
	tokenx := test.RunGetToken(strMyUidx)
	//tokenx := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiI3MDcwMDgxNTMiLCJQbGF0Zm9ybSI6IkFuZHJvaWQiLCJleHAiOjE5NjY0MTJ1XjJZGWj5fB3mqC7p6ytxSarvxZfsABwIjoxNjUxMDU1MDU2fQ.aWvmJ_sQxXmT5nKwiM5QsF9-tfkldzOYZtRD3nrUuko"
	//go func() {
	//	time.Sleep(2 * time.Second)
	//	test.InOutLogou()
	//}()

	test.InOutDoTest(strMyUidx, tokenx, test.WSADDR, test.APIADDR)
	//	test.InOutDoTest(strMyUidx, tokenx, test.WSADDR, test.APIADDR)

	//	time.Sleep(5 * time.Second)
	//	test.SetListenerAndLogin(strMyUidx, tokenx)
	//test.DoTestSetGroupMemberInfo("1104164664", "3188816039", "set ex")

	//	test.DotestGetGroupMemberList()
	//time.Sleep(100000 * time.Second)
	//	test.DoTestCreateGroup()

	//	test.DoTestJoinGroup()
	//	test.DoTestGetGroupsInfo()
	//	test.DoTestDeleteAllMsgFromLocalAndSvr()

	//	println("token ", tokenx)
	time.Sleep(100000 * time.Second)
	b := utils.GetCurrentTimestampBySecond()
	i := 0
	for {
		test.DoTestSendMsg2c2c(strMyUidx, "3380999461", i)
		i++
		time.Sleep(100 * time.Millisecond)
		if i == 10000 {
			break
		}
		log.Warn("", "10 * time.Millisecond ###################waiting... msg: ", i)
	}

	log.Warn("", "cost time: ", utils.GetCurrentTimestampBySecond()-b)
	time.Sleep(100000 * time.Second)
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
