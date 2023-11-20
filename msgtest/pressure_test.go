package msgtest

import (
	"flag"
	"testing"
	"time"

	"github.com/OpenIMSDK/tools/log"
)

const (
	TenThousandGroupUserNum = 10000
	ThousandGroupUserNum    = 1000
	HundredGroupUserNum     = 100
	FiftyGroupUserNum       = 50

	TenThousandGroupNum = 2
	ThousandGroupNum    = 5
	HundredGroupNum     = 50
	FiftyGroupNum       = 100

	FastenedUserPrefix = "fastened_user_prefix"
	RecvMsgPrefix      = "recv_msg_prefix"
)

var (
	totalOnlineUserNum    int // 总在线用户数
	friendMsgSenderNum    int // 好友消息发送者数
	NotFriendMsgSenderNum int // 非好友消息发送者数
	groupMsgSenderNum     int // 群消息发送者数
	msgSenderNumEvreyUser int // 每个用户的消息数
	fastenedUserNum       int // 固定用户数

	recvMsgUserNum int // 消息接收者数, 抽样账号

)

func InitWithFlag() {
	flag.IntVar(&totalOnlineUserNum, "t", 100000, "total online user num")
	flag.IntVar(&friendMsgSenderNum, "f", 100, "friend msg sender num")
	flag.IntVar(&NotFriendMsgSenderNum, "n", 100, "not friend msg sender num")
	flag.IntVar(&groupMsgSenderNum, "g", 100, "group msg sender num")
	flag.IntVar(&msgSenderNumEvreyUser, "m", 100, "msg sender num evrey user")

	flag.IntVar(&recvMsgUserNum, "r", 20, "recv msg user num")
	flag.IntVar(&fastenedUserNum, "u", 300, "fastened user num")
}

func init() {

	InitWithFlag()

	if err := log.InitFromConfig("sdk.log", "sdk", 4,
		true, false, "./chat_log", 2, 24); err != nil {
		panic(err)
	}
}

func Test_PressureFull(t *testing.T) {
	flag.Parse()
	if friendMsgSenderNum+NotFriendMsgSenderNum+groupMsgSenderNum > totalOnlineUserNum {
		t.Fatal("sender num > total online user num")
	}

	p := NewPressureTester()
	// gen userIDs
	userIDs, fastenedUserIDs, recvMsgUserIDs := p.genUserIDs()

	//// register
	//if err := p.registerUsers(userIDs, fastenedUserIDs, recvMsgUserIDs); err != nil {
	//	t.Fatalf("register users failed, err: %v", err)
	//}
	// init users
	p.initUserConns(userIDs, fastenedUserIDs)

	// create groups
	err := p.createTestGroups(userIDs, fastenedUserIDs, recvMsgUserIDs)
	if err != nil {
		t.Fatal(err)
	}

	// import friends
	if err := p.importFriends(p.friendSenderUserIDs, fastenedUserIDs); err != nil {
		t.Fatal(err)
	}

	p.pressureSendMsg()
	// send msg test
}

func Test_InitUserConn(t *testing.T) {
	flag.Parse()
	p := NewPressureTester()
	userNum := 50000
	// gen userIDs
	userIDs := p.userManager.GenUserIDs(userNum)
	// register
	//if err := p.registerUsers(userIDs, nil, nil); err != nil {
	//	t.Fatalf("register users failed, err: %v", err)
	//}
	// init users
	p.initUserConns(userIDs, nil)
	time.Sleep(time.Hour * 60)
}
