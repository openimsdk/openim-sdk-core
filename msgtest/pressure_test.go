package msgtest

import (
	"context"
	"flag"
	"open_im_sdk/msgtest/module"
	"open_im_sdk/sdk_struct"
	"sync"
	"testing"
	"time"

	"github.com/OpenIMSDK/tools/log"
)

func InitWithFlag() {
	flag.IntVar(&totalOnlineUserNum, "t", 100000, "total online user num")
	flag.IntVar(&friendMsgSenderNum, "f", 100, "friend msg sender num")
	flag.IntVar(&NotFriendMsgSenderNum, "n", 100, "not friend msg sender num")
	flag.IntVar(&groupMsgSenderNum, "g", 100, "group msg sender num")
	flag.IntVar(&msgSenderNumEvreyUser, "m", 100, "msg sender num evrey user")

	flag.IntVar(&recvMsgUserNum, "r", 20, "recv msg user num")

}

const (
	TenThousandGroupUserNum = 10000
	ThousandGroupUserNum    = 1000
	HundredGroupUserNum     = 100
	FiftyGroupUserNum       = 50

	TenThousandGroupNum = 2
	ThousandGroupNum    = 5
	HundredGroupNum     = 50
	FiftyGroupNum       = 100
)

var (
	totalOnlineUserNum    int // 总在线用户数
	friendMsgSenderNum    int // 好友消息发送者数
	NotFriendMsgSenderNum int // 非好友消息发送者数
	groupMsgSenderNum     int // 群消息发送者数
	msgSenderNumEvreyUser int // 每个用户的消息数

	recvMsgUserNum int // 消息接收者数, 抽样账号
)

func init() {
	InitWithFlag()
	flag.Parse()
	if err := log.InitFromConfig("sdk.log", "sdk", 4,
		true, false, "./chat_log", 2, 24); err != nil {
		panic(err)
	}
}

func Test_Pressure(t *testing.T) {
	if friendMsgSenderNum+NotFriendMsgSenderNum+groupMsgSenderNum > totalOnlineUserNum {
		t.Fatal("sender num > total online user num")
	}
	p := NewPressureTester()
	// sample recv msg user
	recvMsgUserIDs := p.userManager.GenUserIDs(recvMsgUserNum)
	userIDs := p.userManager.GenUserIDs(totalOnlineUserNum)
	var groupSenderUserIDs, friendSenderUserIDs, notfriendSenderUserIDs []string
	if err := p.userManager.RegisterUsers(userIDs...); err != nil {
		t.Fatal(err)
	}
	for i, userID := range userIDs {
		token, err := p.userManager.GetToken(userID, int32(PLATFORMID))
		if err != nil {
			log.ZError(context.Background(), "get token failed", err)
			continue
		}
		user := module.NewUser(userID, token, sdk_struct.IMConfig{WsAddr: WSADDR, ApiAddr: APIADDR, PlatformID: int32(PLATFORMID)})
		if 0 <= i && i < friendMsgSenderNum {
			p.msgSender[userID] = user
			friendSenderUserIDs = append(friendSenderUserIDs, userID)
		} else if friendMsgSenderNum <= i && i < friendMsgSenderNum+NotFriendMsgSenderNum {
			p.msgSender[userID] = user
			notfriendSenderUserIDs = append(notfriendSenderUserIDs, userID)
		} else if friendMsgSenderNum+NotFriendMsgSenderNum <= i && i < friendMsgSenderNum+NotFriendMsgSenderNum+groupMsgSenderNum {
			p.groupMsgSender[userID] = user
			groupSenderUserIDs = append(groupSenderUserIDs, userID)
		}
	}
	tenThousandGroupIDs, thousandGroupIDs, hundredGroupUserIDs, fiftyGroupUserIDs, err := p.createTestGroups(userIDs)
	if err != nil {
		t.Fatal(err)
	}
	totalGroupIDs := append(append(append(tenThousandGroupIDs, thousandGroupIDs...), hundredGroupUserIDs...), fiftyGroupUserIDs...)
	// import friends
	for _, recvMsgUserID := range recvMsgUserIDs {
		p.friendManager.ImportFriends(recvMsgUserID, friendSenderUserIDs)
	}
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		p.sendMsgs2Users(friendSenderUserIDs, recvMsgUserIDs, msgSenderNumEvreyUser, time.Second)
	}()
	go func() {
		defer wg.Done()
		p.sendMsgs2Users(notfriendSenderUserIDs, recvMsgUserIDs, msgSenderNumEvreyUser, time.Second)
	}()
	go func() {
		defer wg.Done()
		p.sendMsgs2Groups(groupSenderUserIDs, totalGroupIDs, msgSenderNumEvreyUser, time.Second)
	}()
	wg.Wait()
}

type PressureTester struct {
	friendManager  *module.TestFriendManager
	userManager    *module.TestUserManager
	groupManager   *module.TestGroupManager
	msgSender      map[string]*module.SendMsgUser
	groupMsgSender map[string]*module.SendMsgUser
}

func NewPressureTester() *PressureTester {
	metaManager := module.NewMetaManager(APIADDR, SECRET, MANAGERUSERID)
	return &PressureTester{friendManager: metaManager.NewFriendManager(), userManager: metaManager.NewUserManager(), groupManager: metaManager.NewGroupMananger(),
		msgSender: make(map[string]*module.SendMsgUser), groupMsgSender: make(map[string]*module.SendMsgUser)}
}

func (p *PressureTester) createTestGroups(userIDs []string) (tenThousandGroupIDs, thousandGroupIDs, hundredGroupUserIDs, fiftyGroupUserIDs []string, err error) {
	p.groupManager.CreateGroup("", "", userIDs[0], userIDs[:TenThousandGroupUserNum])
	return
}

func (p *PressureTester) sendMsgs2Users(senderIDs, recvIDs []string, num int, duration time.Duration) {
	var wg sync.WaitGroup
	for _, senderID := range senderIDs {
		for _, recvID := range recvIDs {
			wg.Add(1)
			go func(recvID string) {
				defer wg.Done()
				for i := 0; i < num; i++ {
					if user, ok := p.msgSender[senderID]; ok {
						user.SendMsgWithContext(recvID, i)
					}
					time.Sleep(duration)
				}
			}(recvID)
		}
	}
	wg.Wait()
}

func (p *PressureTester) sendMsgs2Groups(senderIDs, groupIDs []string, num int, duration time.Duration) {
	var wg sync.WaitGroup
	for _, senderID := range senderIDs {
		for _, groupID := range groupIDs {
			wg.Add(1)
			go func(groupID string) {
				for i := 0; i < num; i++ {
					if user, ok := p.groupMsgSender[senderID]; ok {
						user.SendGroupMsgWithContext(groupID, i)
					}
					time.Sleep(duration)
				}
			}(groupID)
		}
	}
	wg.Wait()
}
