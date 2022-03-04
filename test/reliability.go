package test

import (
	"flag"
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

func TestReliability() {

	cmdfile := "./cmd.txt"
	uid := flag.Int("uid", 1, "RpcToken default listen port 10800")
	uidCount := flag.Int("uid_count", 2, "RpcToken default listen port 10800")
	messageCount := flag.Int("message_count", 1, "RpcToken default listen port 10800")
	APIADDR1 := flag.String("api_addr", "http://127.0.0.1:10000", "api addr")
	WSADDR1 := flag.String("ws_addr", "http://127.0.0.1:17778", "ws addr")
	REGISTERADDR1 := flag.String("register_addr", "http://127.0.0.1:10000/auth/user_register", "register addr")
	TOKENADDR1 := flag.String("token_addr", "http://127.0.0.1:10000/auth/user_token", "token addr")
	flag.Parse()

	APIADDR = *APIADDR1
	WSADDR = *WSADDR1
	REGISTERADDR = *REGISTERADDR1
	TOKENADDR = *TOKENADDR1

	var myUid int = *uid
	var uidNum int = *uidCount
	var msgnum int = *messageCount

	log.Info("args is ", myUid, uidNum, msgnum)
	var strMyUid string

	strMyUid = GenUid(myUid)

	runRigister(strMyUid)
	token := runGetToken(strMyUid)

	cmd := GetCmd(myUid, cmdfile)

	log.Info("getcmd value ", cmd)
	switch cmd {
	case -1:
		log.Info("GetCmd failed ")
		time.Sleep(time.Duration(1) * time.Second)
	case 5:
		log.Info("wait 2 mins, then login")
		time.Sleep(time.Duration(1*60) * time.Second)
		DoTest(strMyUid, token, WSADDR, APIADDR)
		log.Info("login do test, only login")
		log.Info("testmypid: ", os.Getpid())
	case 6:
		log.Info("wait 4 mins, then login")
		time.Sleep(time.Duration(2*60) * time.Second)
		DoTest(strMyUid, token, WSADDR, APIADDR)
		log.Info("login do test, only login")
		log.Info("testmypid: ", os.Getpid())
	case 3:
		log.Info("wait 2 mins, then login and send")
		time.Sleep(time.Duration(1*60) * time.Second)
		DoTest(strMyUid, token, WSADDR, APIADDR)
		log.Info("login do test, login and send")

		var recvId string
		var idx string
		rand.Seed(time.Now().UnixNano())
		if msgnum == 0 {
			fmt.Println("dont send,  exit")
			os.Exit(0)
		} else {
			for i := 0; i < msgnum; i++ {
				var r int
				for true {
					time.Sleep(time.Duration(SENDINTERVAL) * time.Millisecond)

					r = rand.Intn(uidNum) + 1
					fmt.Println("test rand ", myUid, uidNum, r)
					if r == myUid {
						continue
					} else {
						break
					}
				}
				recvId = GenUid(r)
				idx = strconv.FormatInt(int64(i), 10)

				DoTestSendMsg(strMyUid, recvId, idx)
			}
		}

	case 4:
		fmt.Println("wait 4 mins, then login and send")
		time.Sleep(time.Duration(2*60) * time.Second)
		DoTest(strMyUid, token, WSADDR, APIADDR)
		fmt.Println("login do test, login and send")

		var recvId string
		var idx string
		rand.Seed(time.Now().UnixNano())
		if msgnum == 0 {
			fmt.Println("dont send,  exit")
			os.Exit(0)
		} else {
			for i := 0; i < msgnum; i++ {
				var r int
				for true {
					time.Sleep(time.Duration(SENDINTERVAL) * time.Millisecond)

					r = rand.Intn(uidNum) + 1
					fmt.Println("test rand ", myUid, uidNum, r)
					if r == myUid {
						continue
					} else {
						break
					}
				}
				recvId = GenUid(r)
				idx = strconv.FormatInt(int64(i), 10)

				DoTestSendMsg(strMyUid, recvId, idx)
			}
		}

	case 1:
		fmt.Println("only login")
		DoTest(strMyUid, token, WSADDR, APIADDR)
		fmt.Println("login do test, only login...")
		fmt.Println("testmypid: ", os.Getpid())
	case 2:
		fmt.Println("login send")
		DoTest(strMyUid, token, WSADDR, APIADDR)
		fmt.Println("login do test, login and send")

		var recvId string
		var idx string
		rand.Seed(time.Now().UnixNano())
		if msgnum == 0 {
			fmt.Println("dont send,  exit")
			os.Exit(0)
		} else {
			for i := 0; i < msgnum; i++ {
				var r int
				for true {
					time.Sleep(time.Duration(SENDINTERVAL) * time.Millisecond)

					r = rand.Intn(uidNum) + 1
					fmt.Println("test rand ", myUid, uidNum, r)
					if r == myUid {
						continue
					} else {
						break
					}
				}
				recvId = GenUid(r)
				idx = strconv.FormatInt(int64(i), 10)

				DoTestSendMsg(strMyUid, recvId, idx)
			}
		}
	case 7:
		fmt.Println("random sleep and send")
		DoTest(strMyUid, token, WSADDR, APIADDR)

		var recvId string
		var idx string
		rand.Seed(time.Now().UnixNano())
		maxSleep := 60
		msgnum = 10
		if msgnum == 0 {
			fmt.Println("dont send,  exit")
			os.Exit(0)
		} else {
			for i := 0; i < msgnum; i++ {
				var r int
				for true {
					time.Sleep(time.Duration(rand.Intn(maxSleep)+1) * time.Second)
					r = rand.Intn(uidNum) + 1
					fmt.Println("test rand ", myUid, uidNum, r)
					if r == myUid {
						continue
					} else {
						break
					}
				}
				recvId = GenUid(r)
				idx = strconv.FormatInt(int64(i), 10)

				DoTestSendMsg(strMyUid, recvId, idx)
			}
		}

	}

}

//for true {
//	time.Sleep(time.Duration(60) * time.Second)
//	fmt.Println("waiting")
//}
