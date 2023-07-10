package funcation

import (
	"math/rand"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/utils"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func RegisterReliabilityUser(id int, timeStamp string) {
	userID := GenUid(id, "reliability_"+timeStamp)
	//Register(userID)
	token, _ := RunGetToken(userID)
	coreMgrLock.Lock()
	defer coreMgrLock.Unlock()
	AllLoginMgr[id] = &CoreNode{token: token, userID: userID}
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

	// 所有成员拉取自己的会话是否有更新
	for i := 0; i < clientNum; i++ {
		var params sdk_params_callback.GetAdvancedHistoryMessageListParams
		params.UserID = AllLoginMgr[i].userID
		//params.ConversationID = "si_7788_7789"
		//params.StartClientMsgID = "83ca933d559d0374258550dd656a661c"
		params.Count = 20
		open_im_sdk.GetAdvancedHistoryMessageList(&testConversation, utils.OperationIDGenerator(), utils.StructToJsonString(params))
	}

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
	strMyUid := AllLoginMgr[index].userID
	token := AllLoginMgr[index].token
	ReliabilityInitAndLogin(index, strMyUid, token)
	log.Info("", "login ok client num: ", len(AllLoginMgr))
	log.Warn("start One", index, beforeLoginSleep, isSendMsg, strMyUid, token, WSADDR, APIADDR)

	msgnum := msgNumInOneClient
	uidNum := len(AllLoginMgr)
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
			recvId := AllLoginMgr[r].userID
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
