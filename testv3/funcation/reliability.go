package funcation

import (
	"math/rand"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func RegisterReliabilityUser(id int, timeStamp string) {
	userID := GenUid(id, "reliability_"+timeStamp)
	register(userID)
	token, _ := RunGetToken(userID)
	coreMgrLock.Lock()
	defer coreMgrLock.Unlock()
	allLoginMgr[id] = &CoreNode{token: token, userID: userID}
}

func ReliabilityTest(msgNumOneClient int, intervalSleepMS int, randSleepMaxSecond int, clientNum int) {
	msgNumInOneClient = msgNumOneClient
	var wg sync.WaitGroup
	// 注册
	wg.Add(clientNum)
	for i := 0; i < clientNum; i++ {
		go func(idx int) {
			RegisterReliabilityUser(idx, utils.Int64ToString(time.Now().Unix()))
			wg.Done()
		}(i)
	}
	wg.Wait()

	log.Warn("", "RegisterReliabilityUser finished, clientNum: ", clientNum)
	log.Warn("", " init, login, send msg, start ")
	rand.Seed(time.Now().UnixNano())

	// 用户立刻登录发消息
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
	time.Sleep(time.Duration(3000) * time.Second)
	//for {
	//	// 消息异步落库可能出现延迟，每隔五秒再检查一次
	//	if CheckReliabilityResult(msgNumOneClient, clientNum) {
	//		log.Warn("", "CheckReliabilityResult ok, exit")
	//		os.Exit(0)
	//		return
	//	} else {
	//		log.Warn("", "CheckReliabilityResult failed , wait.... ")
	//	}
	//	time.Sleep(time.Duration(5) * time.Second)
	//}
}

func ReliabilityOne(index int, beforeLoginSleep int, isSendMsg bool, intervalSleepMS int) {
	//	time.Sleep(time.Duration(beforeLoginSleep) * time.Second)
	strMyUid := allLoginMgr[index].userID
	token := allLoginMgr[index].token
	ReliabilityInitAndLogin(index, strMyUid, token)
	log.Info("", "login ok client num: ", len(allLoginMgr))
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
			//time.Sleep(time.Duration(intervalSleepMS) * time.Millisecond)
			// 互相发送消息，非自己
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
