package test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"os"
	"strconv"
	"strings"
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

var testTotalNum = 0
var intervalSleepMS = 0

//var Msgwg sync.WaitGroup
var sendMsgClient = 0

func ReliabilityTest(msgNum int, interval int, ip string, randSleepMaxSecond int, clientNum int) {
	testTotalNum = msgNum
	intervalSleepMS = interval
	TESTIP = ip
	testClientNum := clientNum
	fmt.Println("1111111")

	fmt.Println("2222222222222222")
	for i := 0; i < testClientNum; i++ {
		fmt.Println("3333333333333")
		GenWsReliability(i)
		fmt.Println("4444444444444444")
	}
	rand.Seed(time.Now().UnixNano())
	//log.NewPrivateLog("sdk", 6)
	for i := 0; i < testClientNum; i++ {
		rdSleep := rand.Intn(randSleepMaxSecond) + 1
		isSend := rand.Intn(2)
		if isSend == 0 {
			go ReliabilityOne(i, rdSleep, true)
			sendMsgClient++
		} else {
			go ReliabilityOne(i, rdSleep, false)
		}
	}

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
			log.Error("", "check failed  not in recv ", ksend)
			return false
		}
	}

	log.Info("", "start check map recv -> map send")
	sameNum = 0

	for k1, v1 := range RecvAllMsg {
		v2, ok := SendSuccAllMsg[k1]
		if ok {
			sameNum++
			x := v1 + v2
			x = x + x

		} else {
			log.Error("", "check failed  not in send ", k1)
			//	return false
		}
	}

	log.Warn("", "send msg succ num ", len(SendSuccAllMsg), "recv msg num ", len(RecvAllMsg), "same num ", sameNum)
	log.Warn("", "send msg failed num ", len(SendFailedAllMsg))
	log.Warn("", "need send msg num : ", sendMsgClient*testTotalNum)
	if len(SendSuccAllMsg) > 0 {
		return true
	}
	return false
}

func ReliabilityOne(index int, beforeLoginSleep int, isSendMsg bool) {

	time.Sleep(time.Duration(beforeLoginSleep) * time.Second)
	//	coreMgrLock.Lock()
	strMyUid := allLoginMgr[index].userID
	token := allLoginMgr[index].token
	//	coreMgrLock.Unlock()
	ReliabilityInitAndLogin(index, strMyUid, token, WSADDR, APIADDR)
	log.Info("start One", index, beforeLoginSleep, isSendMsg, strMyUid, token, WSADDR, APIADDR)
	msgnum := testTotalNum
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
			DoTestSendMsg(index, strMyUid, recvId, idx)

		}
		//Msgwg.Done()
	}
}
