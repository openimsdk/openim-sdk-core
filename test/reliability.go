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
	wg.Add(clientNum)
	for i := 0; i < clientNum; i++ {
		go func(idx int) {
			RegisterUserReliability(idx, timeStamp)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Info("", "RegisterUserReliability finish ", clientNum)

	rand.Seed(time.Now().UnixNano())

	wg.Add(clientNum)
	for i := 0; i < clientNum; i++ {
		rdSleep := rand.Intn(randSleepMaxSecond) + 1
		isSend := rand.Intn(2)
		if isSend == 0 {
			go func(idx int) {
				ReliabilityOne(idx, rdSleep, true, intervalSleepMS)
				wg.Done()
			}(i)
			sendMsgClient++
		} else {
			go func(idx int) {
				ReliabilityOne(idx, rdSleep, false, intervalSleepMS)
				wg.Done()
			}(i)
		}
	}
	wg.Wait()
	log.Warn("CheckReliabilityResult start, send msg client number: ", sendMsgClient, "total client number: ", clientNum)

	for {
		if CheckReliabilityResult() {
			log.Warn("", "CheckReliabilityResult ok, exit")
			os.Exit(0)
			return
		} else {
			log.Warn("", "CheckReliabilityResult failed , wait.... ")
		}
		time.Sleep(time.Duration(300) * time.Second)
	}
}

func PressTest(msgNumOneClient int, intervalSleepMS int, randSleepMaxSecond int, clientNum int) {
	msgNumInOneClient = msgNumOneClient
	timeStamp := utils.Int64ToString(time.Now().Unix())

	var wg sync.WaitGroup
	wg.Add(clientNum)
	for i := 0; i < clientNum; i++ {
		go func(idx int) {
			RegisterUserReliability(idx, timeStamp)
			log.Warn("", "get user token finish ", idx)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Info("", "get all user token finish ", clientNum)

	rand.Seed(time.Now().UnixNano())

	wg.Add(clientNum)
	for i := 0; i < clientNum; i++ {
		rdSleep := rand.Intn(randSleepMaxSecond) + 1
		isSend := rand.Intn(2)
		isSend = 0
		if isSend == 0 {
			go func(idx int) {
				PressOne(idx, rdSleep, true, intervalSleepMS)
				wg.Done()
			}(i)
			sendMsgClient++
		} else {
			go func(idx int) {
				PressOne(idx, rdSleep, false, intervalSleepMS)
				wg.Done()
			}(i)
		}
	}
	wg.Wait()
	log.Warn("CheckReliabilityResult start, send msg client number: ", sendMsgClient, "total client number: ", clientNum)
}

func CheckReliabilityResult() bool {
	log.Info("", "start check map send -> map recv")
	sameNum := 0

	for ksend, vsend := range SendSuccAllMsg {
		krecv, ok := RecvAllMsg[ksend]
		if ok {
			sameNum++
			x := vsend
			y := krecv
			x = x + x
			y = y + y

		} else {
			log.Error("", "check failed  not in recv ", ksend, len(SendFailedAllMsg), len(SendSuccAllMsg), len(RecvAllMsg))
			return false
		}
	}
	log.Info("", "check map send -> map recv ok ", sameNum)
	log.Info("", "start check map recv -> map send ")
	sameNum = 0

	for k1, v1 := range RecvAllMsg {
		v2, ok := SendSuccAllMsg[k1]
		if ok {
			sameNum++
			x := v1 + v2
			x = x + x

		} else {
			log.Error("", "check failed  not in send ", k1, len(SendFailedAllMsg), len(SendSuccAllMsg), len(RecvAllMsg))
			//	return false
		}
	}

	log.Warn("", "need send msg num : ", sendMsgClient*msgNumInOneClient)
	log.Warn("", "send msg succ num ", len(SendSuccAllMsg))
	log.Warn("", "send msg failed num ", len(SendFailedAllMsg))
	log.Warn("", "recv msg succ num ", len(RecvAllMsg))
	log.Warn("", "msg in recv, and in send num ", sameNum)

	return true
}

func ReliabilityOne(index int, beforeLoginSleep int, isSendMsg bool, intervalSleepMS int) {
	time.Sleep(time.Duration(beforeLoginSleep) * time.Second)
	strMyUid := allLoginMgr[index].userID
	token := allLoginMgr[index].token
	ReliabilityInitAndLogin(index, strMyUid, token, WSADDR, APIADDR)
	log.Info("", "login ok client num: ", len(allLoginMgr))
	log.Warn("start One", index, beforeLoginSleep, isSendMsg, strMyUid, token, WSADDR, APIADDR)
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
			time.Sleep(time.Duration(intervalSleepMS) * time.Millisecond)
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
					log.Warn("", "NumGoroutine > max  ", runtime.NumGoroutine(), MaxNumGoroutine)
					continue
				} else {
					break
				}
			}

			DoTestSendMsg(index, strMyUid, recvId, idx)

		}
		//Msgwg.Done()
	}
}

func PressOne(index int, beforeLoginSleep int, isSendMsg bool, intervalSleepMS int) {
	time.Sleep(time.Duration(beforeLoginSleep) * time.Second)
	strMyUid := allLoginMgr[index].userID
	token := allLoginMgr[index].token
	ReliabilityInitAndLogin(index, strMyUid, token, WSADDR, APIADDR)
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
					log.Warn("", "NumGoroutine > max  ", runtime.NumGoroutine(), MaxNumGoroutine)
					continue
				} else {
					break
				}
			}

			DoTestSendMsgPress(index, strMyUid, recvId, idx)

		}
		//Msgwg.Done()
	}
}
