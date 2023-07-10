// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func GetFileContentAsStringLines(filePath string) ([]string, error) {
	result := []string{}
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return result, err
	}
	s := string(b)
	for _, lineStr := range strings.Split(s, "\n") {
		lineStr = strings.TrimSpace(lineStr)
		if lineStr == "" {
			continue
		}
		result = append(result, lineStr)
	}
	return result, nil
}

func GetCmd(myUid int, filename string) int {
	cmd, err := GetFileContentAsStringLines("cmd.txt")
	if err != nil {
		fmt.Println("GetFileContentAsStringLines failed")
		return -1
	}
	if len(cmd) < myUid {
		fmt.Println("len failed")
		return -1
	}
	return int(utils.StringToInt64(cmd[myUid-1]))
}

func ReliabilityTest(msgNumOneClient int, intervalSleepMS int, randSleepMaxSecond int, clientNum int) {
	msgNumInOneClient = msgNumOneClient
	timeStamp := utils.Int64ToString(time.Now().Unix())

	var wg sync.WaitGroup
	// 注册
	wg.Add(clientNum)
	for i := 0; i < clientNum; i++ {
		go func(idx int) {
			RegisterReliabilityUser(idx, timeStamp)
			wg.Done()
		}(i)
	}
	wg.Wait()

	log.Warn("", "RegisterReliabilityUser finished, clientNum: ", clientNum)
	log.Warn("", " init, login, send msg, start ")
	rand.Seed(time.Now().UnixNano())

	// 一半用户立刻登录发消息
	wg.Add(clientNum)
	for i := 0; i < clientNum; i++ {
		rdSleep := rand.Intn(randSleepMaxSecond) + 1
		isSend := 0 // 消息是否成功发送控制量
		if isSend == 0 {
			go func(idx int) {
				log.Warn("", " send msg flag true ", idx)
				ReliabilityOne(idx, rdSleep, true, intervalSleepMS)
				wg.Done()
			}(i)
			sendMsgClient++
		} else {
			go func(idx int) {
				log.Warn("", " send msg flag false ", idx)
				ReliabilityOne(idx, rdSleep, false, intervalSleepMS)
				wg.Done()
			}(i)
		}
	}
	wg.Wait()
	log.Warn("send msg finish,  CheckReliabilityResult")

	for {
		// 消息异步落库可能出现延迟，每隔五秒再检查一次
		if CheckReliabilityResult(msgNumOneClient, clientNum) {
			log.Warn("", "CheckReliabilityResult ok, exit")
			os.Exit(0)
			return
		} else {
			log.Warn("", "CheckReliabilityResult failed , wait.... ")
		}
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func WorkGroupReliabilityTest(msgNumOneClient int, intervalSleepMS int, randSleepMaxSecond int, clientNum int, groupID string) {
	msgNumInOneClient = msgNumOneClient
	//timeStamp := utils.Int64ToString(time.Now().Unix())

	var wg sync.WaitGroup
	wg.Add(clientNum)
	for i := 0; i < clientNum; i++ {
		go func(idx int) {
			WorkGroupRegisterReliabilityUser(idx)
			wg.Done()
		}(i)
	}
	wg.Wait()

	log.Warn("", "RegisterReliabilityUser finished, clientNum: ", clientNum)
	log.Warn("", " init, login, send msg, start ")
	rand.Seed(time.Now().UnixNano())

	wg.Add(clientNum)
	for i := 0; i < clientNum; i++ {
		rdSleep := rand.Intn(randSleepMaxSecond) + 1
		isSend := 0
		if isSend == 0 {
			go func(idx int) {
				log.Warn("", " send msg flag true ", idx)
				WorkGroupReliabilityOne(idx, rdSleep, true, intervalSleepMS, groupID)
				wg.Done()
			}(i)
			sendMsgClient++
		} else {
			go func(idx int) {
				log.Warn("", " send msg flag false ", idx)
				ReliabilityOne(idx, rdSleep, false, intervalSleepMS)
				wg.Done()
			}(i)
		}
	}
	wg.Wait()
	log.Warn("send msg finish,  CheckReliabilityResult")

	for {
		if CheckReliabilityResult(msgNumOneClient, clientNum) {
			log.Warn("", "CheckReliabilityResult ok, exit")
			os.Exit(0)
			return
		} else {
			log.Warn("", "CheckReliabilityResult failed , wait.... ")
		}
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func WorkGroupMsgDelayTest(msgNumOneClient int, intervalSleepMS int, randSleepMaxSecond int, clientBegin int, clientEnd int, groupID string) {
	msgNumInOneClient = msgNumOneClient

	var wg sync.WaitGroup

	wg.Add(clientEnd - clientBegin + 1)
	for i := clientBegin; i <= clientEnd; i++ {
		go func(idx int) {
			WorkGroupRegisterReliabilityUser(idx)
			wg.Done()
		}(i)
	}
	wg.Wait()

	log.Warn("", "RegisterReliabilityUser finished, client: ", clientBegin, clientEnd)
	log.Warn("", " init, login, send msg, start ")
	rand.Seed(time.Now().UnixNano())

	wg.Add(clientEnd - clientBegin + 1)
	for i := clientBegin; i <= clientEnd; i++ {
		rdSleep := rand.Intn(randSleepMaxSecond) + 1
		isSend := 0
		if isSend == 0 {
			go func(idx int) {
				log.Warn("", " send msg flag true ", idx)
				WorkGroupReliabilityOne(idx, rdSleep, true, intervalSleepMS, groupID)
				wg.Done()
			}(i)
			sendMsgClient++
		} else {
			go func(idx int) {
				log.Warn("", " send msg flag false ", idx)
				WorkGroupReliabilityOne(idx, rdSleep, false, intervalSleepMS, groupID)
				wg.Done()
			}(i)
		}
	}
	wg.Wait()
	log.Warn("send msg finish,  CheckReliabilityResult")

	for {
		if CheckReliabilityResult(msgNumOneClient, clientEnd-clientBegin+1) {
			log.Warn("", "CheckReliabilityResult ok, exit")
			os.Exit(0)
			return
		} else {
			log.Warn("", "CheckReliabilityResult failed , wait.... ")
		}
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func PressTest(msgNumOneClient int, intervalSleepMS int, clientNum int) {
	msgNumInOneClient = msgNumOneClient
	//timeStamp := utils.Int64ToString(time.Now().Unix())
	t1 := time.Now()
	var wg sync.WaitGroup
	wg.Add(clientNum)
	for i := 0; i < clientNum; i++ {
		go func(idx int) {
			RegisterPressUser(idx)
			log.Info("", "get user token finish ", idx)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Warn("", "get all user token finish ", clientNum, " cost time: ", time.Since(t1))

	log.Warn("", "init and login begin ")
	t1 = time.Now()
	wg.Add(clientNum)
	for i := 0; i < clientNum; i++ {
		go func(idx int) {
			strMyUid := allLoginMgr[idx].userID
			token := allLoginMgr[idx].token
			PressInitAndLogin(idx, strMyUid, token, WSADDR, APIADDR)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Warn("", "init and login end ", " cost time: ", time.Since(t1))

	log.Warn("", "send msg begin ")
	t1 = time.Now()
	wg.Add(clientNum)
	for i := 0; i < clientNum; i++ {
		go func(idx int) {
			PressOne(idx, 0, true, intervalSleepMS)
			log.Warn("", "press finished  ", idx)
			wg.Done()
		}(i)
	}
	wg.Wait()
	sendMsgTotalSuccessNum := uint32(0)
	sendMsgTotalFailedNum := uint32(0)
	for _, v := range allLoginMgr {
		sendMsgTotalSuccessNum += v.sendMsgSuccessNum
		sendMsgTotalFailedNum += v.sendMsgFailedNum
	}
	log.Warn("send msg end  ", "number of messages expected to be sent: ", clientNum*msgNumOneClient, " sendMsgTotalSuccessNum: ", sendMsgTotalSuccessNum, " sendMsgTotalFailedNum: ", sendMsgTotalFailedNum, "cost time: ", time.Since(t1))
}

func WorkGroupPressTest(msgNumOneClient int, intervalSleepMS int, clientNum int, groupID string) {
	msgNumInOneClient = msgNumOneClient
	t1 := time.Now()
	var wg sync.WaitGroup
	wg.Add(clientNum)
	for i := 0; i < clientNum; i++ {
		go func(idx int) {
			WorkGroupRegisterReliabilityUser(idx)
			log.Info("", "get user token finish ", idx)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Warn("", "get all user token finish ", clientNum, " cost time: ", time.Since(t1))

	log.Warn("", "init and login begin ")
	t1 = time.Now()
	wg.Add(clientNum)
	for i := 0; i < clientNum; i++ {
		go func(idx int) {
			strMyUid := allLoginMgr[idx].userID
			token := allLoginMgr[idx].token
			ReliabilityInitAndLogin(idx, strMyUid, token, WSADDR, APIADDR)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Warn("", "init and login end ", " cost time: ", time.Since(t1))

	log.Warn("", "send msg begin ")
	t1 = time.Now()
	wg.Add(clientNum)
	for i := 0; i < clientNum; i++ {
		go func(idx int) {
			WorkGroupPressOne(idx, 0, true, intervalSleepMS, groupID)
			wg.Done()
		}(i)
	}
	wg.Wait()
	sendMsgTotalSuccessNum := uint32(0)
	sendMsgTotalFailedNum := uint32(0)
	for _, v := range allLoginMgr {
		sendMsgTotalSuccessNum += v.sendMsgSuccessNum
		sendMsgTotalFailedNum += v.sendMsgFailedNum
	}
	log.Warn("send msg end  ", "number of messages expected to be sent: ", clientNum*msgNumOneClient, " sendMsgTotalSuccessNum: ", sendMsgTotalSuccessNum, " sendMsgTotalFailedNum: ", sendMsgTotalFailedNum, "cost time: ", time.Since(t1))
}

func CheckReliabilityResult(msgNumOneClient int, clientNum int) bool {
	log.Info("", "start check map send -> map recv")
	sameNum := 0

	// 消息数量不一致说明出现丢失
	if len(SendSuccAllMsg)+len(SendFailedAllMsg) != msgNumOneClient*clientNum {
		log.Warn("", utils.GetSelfFuncName(), " send msg success number: ", len(SendSuccAllMsg),
			" send msg failed number: ", len(SendFailedAllMsg), " all: ", msgNumOneClient*clientNum)
		return false
	}

	for ksend, _ := range SendSuccAllMsg {
		_, ok := RecvAllMsg[ksend] // RecvAllMsg 的初始化何时？
		if ok {
			sameNum++
		} else {
			// 埋点日志，第 ksend 个消息数据 本地和服务器不一致
			log.Error("", "check failed not in recv ", ksend)
			log.Error("", "send failed num: ", len(SendFailedAllMsg),
				" send success num: ", len(SendSuccAllMsg), " recv num: ", len(RecvAllMsg))
			return false
		}
	}
	log.Info("", "check map send -> map recv ok ", sameNum)
	//log.Info("", "start check map recv -> map send ")
	//sameNum = 0

	//for k1, _ := range RecvAllMsg {
	//	_, ok := SendSuccAllMsg[k1]
	//	if ok {
	//		sameNum++
	//		//x := v1 + v2
	//		//x = x + x
	//
	//	} else {
	//		log.Error("", "check failed  not in send ", k1, len(SendFailedAllMsg), len(SendSuccAllMsg), len(RecvAllMsg))
	//		//	return false
	//	}
	//}
	maxCostMsgID := ""
	minCostTime := int64(1000000)
	maxCostTime := int64(0)
	totalCostTime := int64(0)
	for ksend, vsend := range SendSuccAllMsg {
		krecv, ok := RecvAllMsg[ksend]
		if ok {
			sameNum++
			costTime := krecv.RecvTime - vsend.SendTime
			totalCostTime += costTime
			if costTime > maxCostTime {
				maxCostMsgID = ksend
				maxCostTime = costTime
			}
			if minCostTime > costTime {
				minCostTime = costTime
			}
		}
	}

	log.Warn("", "need send msg num : ", sendMsgClient*msgNumInOneClient)
	log.Warn("", "send msg succ num ", len(SendSuccAllMsg))
	log.Warn("", "send msg failed num ", len(SendFailedAllMsg))
	log.Warn("", "recv msg succ num ", len(RecvAllMsg))
	log.Warn("", "minCostTime: ", minCostTime, "ms, maxCostTime: ", maxCostTime, "ms, average cost time: ",
		totalCostTime/(int64(sendMsgClient*msgNumInOneClient)), "ms", " maxCostMsgID: ", maxCostMsgID)

	return true
}

func ReliabilityOne(index int, beforeLoginSleep int, isSendMsg bool, intervalSleepMS int) {
	//	time.Sleep(time.Duration(beforeLoginSleep) * time.Second)
	strMyUid := allLoginMgr[index].userID
	token := allLoginMgr[index].token
	log.Info("", "login ok client num: ", len(allLoginMgr), "userID ", strMyUid, "token: ", token, " index: ", index)
	ReliabilityInitAndLogin(index, strMyUid, token, WSADDR, APIADDR)

	log.Warn("start One", index, beforeLoginSleep, isSendMsg, strMyUid, token, WSADDR, APIADDR)

	msgnum := msgNumInOneClient
	uidNum := len(allLoginMgr)
	rand.Seed(time.Now().UnixNano())
	if msgnum == 0 {
		os.Exit(0)
	}
	if !isSendMsg {
		//	Msgwg.Done()
	} else {
		for i := 0; i < msgnum; i++ {
			var r int
			time.Sleep(time.Duration(intervalSleepMS) * time.Millisecond)
			for {
				r = rand.Intn(uidNum)
				if r == index {
					continue
				} else {
					break
				}
			}
			recvId := allLoginMgr[r].userID
			idx := strconv.FormatInt(int64(i), 10)
			for {
				if runtime.NumGoroutine() > MaxNumGoroutine {
					time.Sleep(time.Duration(intervalSleepMS) * time.Millisecond)
					log.Warn("", "NumGoroutine > max  ", runtime.NumGoroutine(), MaxNumGoroutine)
					continue
				} else {
					break
				}
			}
			DoTestSendMsg(index, strMyUid, recvId, "", idx)
		}
		//Msgwg.Done()
	}
}

func WorkGroupReliabilityOne(index int, beforeLoginSleep int, isSendMsg bool, intervalSleepMS int, groupID string) {
	//	time.Sleep(time.Duration(beforeLoginSleep) * time.Second)
	strMyUid := allLoginMgr[index].userID
	token := allLoginMgr[index].token
	ReliabilityInitAndLogin(index, strMyUid, token, WSADDR, APIADDR)
	log.Info("", "login ok client num: ", len(allLoginMgr))
	log.Warn("start One", index, beforeLoginSleep, isSendMsg, strMyUid, token, WSADDR, APIADDR)
	msgnum := msgNumInOneClient
	uidNum := len(allLoginMgr)
	var idx string
	rand.Seed(time.Now().UnixNano())
	if msgnum == 0 {
		os.Exit(0)
	}
	if !isSendMsg {
		//	Msgwg.Done()
	} else {
		for i := 0; i < msgnum; i++ {
			var r int
			time.Sleep(time.Duration(intervalSleepMS) * time.Millisecond)
			for {
				r = rand.Intn(uidNum)
				if r == index {
					continue
				} else {

					break
				}

			}

			idx = strconv.FormatInt(int64(i), 10)
			for {
				if runtime.NumGoroutine() > MaxNumGoroutine {
					time.Sleep(time.Duration(intervalSleepMS) * time.Millisecond)
					log.Warn("", "NumGoroutine > max  ", runtime.NumGoroutine(), MaxNumGoroutine)
					continue
				} else {
					break
				}
			}

			DoTestSendMsg(index, strMyUid, "", groupID, idx)

		}
		//Msgwg.Done()
	}
}

func WorkGroupMsgDelayOne(index int, beforeLoginSleep int, isSendMsg bool, intervalSleepMS int, groupID string) {
	//	time.Sleep(time.Duration(beforeLoginSleep) * time.Second)
	strMyUid := allLoginMgr[index].userID
	token := allLoginMgr[index].token
	ReliabilityInitAndLogin(index, strMyUid, token, WSADDR, APIADDR)
	log.Info("", "login ok client num: ", len(allLoginMgr))
	log.Warn("start One", index, beforeLoginSleep, isSendMsg, strMyUid, token, WSADDR, APIADDR)
	msgnum := msgNumInOneClient
	uidNum := len(allLoginMgr)
	var idx string
	rand.Seed(time.Now().UnixNano())
	if msgnum == 0 {
		os.Exit(0)
	}
	if !isSendMsg {
		//	Msgwg.Done()
	} else {
		for i := 0; i < msgnum; i++ {
			var r int
			time.Sleep(time.Duration(intervalSleepMS) * time.Millisecond)
			for {
				r = rand.Intn(uidNum)
				if r == index {
					continue
				} else {

					break
				}

			}

			idx = strconv.FormatInt(int64(i), 10)
			for {
				if runtime.NumGoroutine() > MaxNumGoroutine {
					time.Sleep(time.Duration(intervalSleepMS) * time.Millisecond)
					log.Warn("", "NumGoroutine > max  ", runtime.NumGoroutine(), MaxNumGoroutine)
					continue
				} else {
					break
				}
			}

			DoTestSendMsg(index, strMyUid, "", groupID, idx)

		}
		//Msgwg.Done()
	}
}

//
//funcation WorkGroupMsgDelayOne(sendID1 string, beforeLoginSleep int, isSendMsg bool, intervalSleepMS int, groupID string) {
//	//	time.Sleep(time.Duration(beforeLoginSleep) * time.Second)
//	strMyUid := allLoginMgr[index].userID
//	token := allLoginMgr[index].token
//	ReliabilityInitAndLogin(index, strMyUid, token, WSADDR, APIADDR)
//	log.Info("", "login ok client num: ", len(allLoginMgr))
//	log.Warn("start One", index, beforeLoginSleep, isSendMsg, strMyUid, token, WSADDR, APIADDR)
//	msgnum := msgNumInOneClient
//	uidNum := len(allLoginMgr)
//	var idx string
//	rand.Seed(time.Now().UnixNano())
//	if msgnum == 0 {
//		os.Exit(0)
//	}
//	if !isSendMsg {
//		//	Msgwg.Done()
//	} else {
//		for i := 0; i < msgnum; i++ {
//			var r int
//			time.Sleep(time.Duration(intervalSleepMS) * time.Millisecond)
//			for {
//				r = rand.Intn(uidNum)
//				if r == index {
//					continue
//				} else {
//
//					break
//				}
//
//			}
//
//			idx = strconv.FormatInt(int64(i), 10)
//			for {
//				if runtime.NumGoroutine() > MaxNumGoroutine {
//					time.Sleep(time.Duration(intervalSleepMS) * time.Millisecond)
//					log.Warn("", "NumGoroutine > max  ", runtime.NumGoroutine(), MaxNumGoroutine)
//					continue
//				} else {
//					break
//				}
//			}
//
//			DoTestSendMsg(index, strMyUid, "", groupID, idx)
//
//		}
//		//Msgwg.Done()
//	}
//}

func PressOne(index int, beforeLoginSleep int, isSendMsg bool, intervalSleepMS int) {
	if beforeLoginSleep != 0 {
		time.Sleep(time.Duration(beforeLoginSleep) * time.Millisecond)
	}

	strMyUid := allLoginMgr[index].userID
	token := allLoginMgr[index].token
	//	ReliabilityInitAndLogin(index, strMyUid, token, WSADDR, APIADDR)
	log.Info("", "login ok client num: ", len(allLoginMgr))
	log.Info("start One", index, beforeLoginSleep, isSendMsg, strMyUid, token, WSADDR, APIADDR)
	msgnum := msgNumInOneClient
	uidNum := len(allLoginMgr)
	var recvId string
	var idx string
	rand.Seed(time.Now().UnixNano())
	if msgnum == 0 {
		os.Exit(0)
	}
	if !isSendMsg {
		//	Msgwg.Done()
	} else {
		for i := 0; i < msgnum; i++ {
			var r int
			//	time.Sleep(time.Duration(intervalSleepMS) * time.Millisecond)
			for {
				r = rand.Intn(uidNum)
				if r == index {
					continue
				} else {

					break
				}

			}

			recvId = allLoginMgr[r].userID
			idx = strconv.FormatInt(int64(i), 10)
			for {
				if runtime.NumGoroutine() > MaxNumGoroutine {
					time.Sleep(time.Duration(intervalSleepMS) * time.Millisecond)
					log.Warn("", " NumGoroutine > max ", runtime.NumGoroutine(), MaxNumGoroutine)
					continue
				} else {
					break
				}
			}
			time.Sleep(time.Duration(intervalSleepMS) * time.Millisecond)
			//DoTestSendMsg(index, strMyUid, recvId, idx)
			if sendPressMsg(index, strMyUid, recvId, "", idx) {
				allLoginMgr[index].sendMsgSuccessNum++
			} else {
				allLoginMgr[index].sendMsgFailedNum++
			}
		}
		//Msgwg.Done()
	}
}

func WorkGroupPressOne(index int, beforeLoginSleep int, isSendMsg bool, intervalSleepMS int, groupID string) {
	if beforeLoginSleep != 0 {
		time.Sleep(time.Duration(beforeLoginSleep) * time.Millisecond)
	}
	strMyUid := allLoginMgr[index].userID
	token := allLoginMgr[index].token
	//ReliabilityInitAndLogin(index, strMyUid, token, WSADDR, APIADDR)
	log.Info("", "login ok, client num: ", len(allLoginMgr))
	log.Info("start One ", index, beforeLoginSleep, isSendMsg, strMyUid, token, WSADDR, APIADDR)
	msgnum := msgNumInOneClient
	var idx string
	rand.Seed(time.Now().UnixNano())
	if msgnum == 0 {
		os.Exit(0)
	}
	if !isSendMsg {
	} else {
		for i := 0; i < msgnum; i++ {
			idx = strconv.FormatInt(int64(i), 10)

			for {
				if runtime.NumGoroutine() > MaxNumGoroutine {
					time.Sleep(time.Duration(intervalSleepMS) * time.Millisecond)
					log.Warn("", " NumGoroutine > max ", runtime.NumGoroutine(), MaxNumGoroutine)
					continue
				} else {
					break
				}
			}
			log.Info("sendPressMsg begin", index, strMyUid, groupID)
			if sendPressMsg(index, strMyUid, "", groupID, idx) {
				allLoginMgr[index].sendMsgSuccessNum++
			} else {
				allLoginMgr[index].sendMsgFailedNum++
			}
			log.Info("sendPressMsg end")
			time.Sleep(time.Duration(intervalSleepMS) * time.Millisecond)
		}
	}
}
