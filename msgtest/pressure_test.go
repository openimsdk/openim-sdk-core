package msgtest

import (
	"context"
	"flag"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/msgtest/module"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"sync"
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

type PressureTester struct {
	friendManager  *module.TestFriendManager
	userManager    *module.TestUserManager
	groupManager   *module.TestGroupManager
	msgSender      map[string]*module.SendMsgUser
	groupMsgSender map[string]*module.SendMsgUser

	groupSenderUserIDs, friendSenderUserIDs, notfriendSenderUserIDs []string
	recvMsgUserIDs                                                  []string

	tenThousandGroupIDs, thousandGroupIDs, hundredGroupUserIDs, fiftyGroupUserIDs []string
}

func NewPressureTester() *PressureTester {
	metaManager := module.NewMetaManager(APIADDR, SECRET, MANAGERUSERID)
	return &PressureTester{friendManager: metaManager.NewFriendManager(), userManager: metaManager.NewUserManager(), groupManager: metaManager.NewGroupMananger(),
		msgSender: make(map[string]*module.SendMsgUser), groupMsgSender: make(map[string]*module.SendMsgUser)}
}

func (p *PressureTester) genUserIDs() (userIDs, fastenedUserIDs, recvMsgUserIDs []string) {
	userIDs = p.userManager.GenUserIDs(totalOnlineUserNum - fastenedUserNum)                  // 在线用户
	fastenedUserIDs = p.userManager.GenUserIDsWithPrefix(fastenedUserNum, FastenedUserPrefix) // 指定300用户
	recvMsgUserIDs = p.userManager.GenUserIDsWithPrefix(recvMsgUserNum, RecvMsgPrefix)        // 抽样用户完整SDK
	return
}

func (p *PressureTester) registerUsers(userIDs []string, fastenedUserIDs []string, recvMsgUserIDs []string) error {
	for i := 0; i < len(userIDs); i += 1000 {
		end := i + 1000
		if end > len(userIDs) {
			end = len(userIDs)
		}
		userIDsSlice := userIDs[i:end]
		if err := p.userManager.RegisterUsers(userIDsSlice...); err != nil {
			return err
		}
		if len(userIDsSlice) < 1000 {
			break
		}
	}
	if len(fastenedUserIDs) != 0 {
		if err := p.userManager.RegisterUsers(fastenedUserIDs...); err != nil {
			return err
		}
	}
	if len(recvMsgUserIDs) != 0 {
		if err := p.userManager.RegisterUsers(recvMsgUserIDs...); err != nil {
			return err
		}
	}
	return nil
}

func (p *PressureTester) initUserConns(userIDs []string, fastenedUserIDs []string) {
	for i, userID := range userIDs {
		token, err := p.userManager.GetToken(userID, int32(PLATFORMID))
		if err != nil {
			log.ZError(context.Background(), "get token failed", err, "userID", userID, "platformID", PLATFORMID)
			continue
		}
		user := module.NewUser(userID, token, sdk_struct.IMConfig{WsAddr: WSADDR, ApiAddr: APIADDR, PlatformID: int32(PLATFORMID)})
		if 0 <= i && i < friendMsgSenderNum {
			p.msgSender[userID] = user
			p.friendSenderUserIDs = append(p.friendSenderUserIDs, userID)
		} else if friendMsgSenderNum <= i && i < friendMsgSenderNum+NotFriendMsgSenderNum {
			p.msgSender[userID] = user
			p.notfriendSenderUserIDs = append(p.notfriendSenderUserIDs, userID)
		}
	}
	if len(fastenedUserIDs) != 0 {
		for _, userID := range fastenedUserIDs {
			token, err := p.userManager.GetToken(userID, int32(PLATFORMID))
			if err != nil {
				log.ZError(context.Background(), "get token failed", err, "userID", userID, "platformID", PLATFORMID)
				continue
			}
			user := module.NewUser(userID, token, sdk_struct.IMConfig{WsAddr: WSADDR, ApiAddr: APIADDR, PlatformID: int32(PLATFORMID)})
			p.msgSender[userID] = user
			p.groupSenderUserIDs = append(p.groupSenderUserIDs, userID)
		}
	}
}

func (p *PressureTester) createTestGroups(userIDs, fastenedUserIDs, recvMsgUserIDs []string) (err error) {
	// create ten thousand group
	for i := 1; i <= TenThousandGroupNum; i++ {
		groupID := p.groupManager.GenGroupID(fmt.Sprintf("tenThousandGroup_%d", i))
		err = p.groupManager.CreateGroup(groupID, "tenThousandGroup", userIDs[0], append(userIDs[(i-1)*TenThousandGroupUserNum:i*TenThousandGroupUserNum-1], fastenedUserIDs...))
		if err != nil {
			return
		}
		p.tenThousandGroupIDs = append(p.tenThousandGroupIDs, groupID)
	}
	// create two thousand group
	exclude := TenThousandGroupNum * TenThousandGroupUserNum
	for i := 1; i <= ThousandGroupNum; i++ {
		groupID := p.groupManager.GenGroupID(fmt.Sprintf("thousandGroup_%d", i))
		err = p.groupManager.CreateGroup(groupID, "thousandGroup", userIDs[0], append(userIDs[exclude+(i-1)*ThousandGroupUserNum:exclude+i*ThousandGroupUserNum-1], fastenedUserIDs...))
		if err != nil {
			return
		}
		p.thousandGroupIDs = append(p.thousandGroupIDs, groupID)
	}
	// create five hundred group
	exclude += exclude + ThousandGroupNum*ThousandGroupUserNum
	for i := 1; i <= HundredGroupNum; i++ {
		groupID := p.groupManager.GenGroupID(fmt.Sprintf("hundredGroup_%d", i))
		err = p.groupManager.CreateGroup(groupID, "hundredGroup", userIDs[0], append(fastenedUserIDs[0:80], recvMsgUserIDs...))
		if err != nil {
			return
		}
		p.hundredGroupUserIDs = append(p.hundredGroupUserIDs, groupID)
	}
	// create fifty group
	exclude += exclude + HundredGroupNum*HundredGroupUserNum
	for i := 1; i <= FiftyGroupNum; i++ {
		groupID := p.groupManager.GenGroupID(fmt.Sprintf("fiftyGroup_%d", i))
		err = p.groupManager.CreateGroup(groupID, "fiftyGroup", userIDs[0], append(fastenedUserIDs[0:30], recvMsgUserIDs...))
		if err != nil {
			return
		}
		p.fiftyGroupUserIDs = append(p.fiftyGroupUserIDs, groupID)
	}
	return
}

func (p *PressureTester) sendMsgs2Users(senderIDs, recvIDs []string, num int, duration time.Duration) {
	var wg sync.WaitGroup
	for _, senderID := range senderIDs {
		for _, recvID := range recvIDs {
			wg.Add(1)
			go func(senderID, recvID string) {
				defer wg.Done()
				for i := 0; i < num; i++ {
					if user, ok := p.msgSender[senderID]; ok {
						user.SendMsgWithContext(recvID, i)
					}
					time.Sleep(duration)
				}
			}(senderID, recvID)
		}
	}
	wg.Wait()
}

func (p *PressureTester) sendMsgs2Groups(senderIDs, groupIDs []string, num int, duration time.Duration) {
	var wg sync.WaitGroup
	for _, senderID := range senderIDs {
		for _, groupID := range groupIDs {
			wg.Add(1)
			go func(senderID, groupID string) {
				defer wg.Done()
				for i := 0; i < num; i++ {
					if user, ok := p.groupMsgSender[senderID]; ok {
						user.SendGroupMsgWithContext(groupID, i)
					}
					time.Sleep(duration)
				}
			}(senderID, groupID)
		}
	}
	wg.Wait()
}

func (p *PressureTester) pressureSendMsg() {
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		p.sendMsgs2Users(p.friendSenderUserIDs, p.recvMsgUserIDs, msgSenderNumEvreyUser, time.Second)
	}()
	go func() {
		defer wg.Done()
		p.sendMsgs2Users(p.notfriendSenderUserIDs, p.recvMsgUserIDs, msgSenderNumEvreyUser, time.Second)
	}()
	go func() {
		defer wg.Done()
		totalGroupIDs := append(append(p.tenThousandGroupIDs, p.thousandGroupIDs...), append(p.hundredGroupUserIDs, p.fiftyGroupUserIDs...)...)
		p.sendMsgs2Groups(p.groupSenderUserIDs, totalGroupIDs, msgSenderNumEvreyUser, time.Second)
	}()
	wg.Wait()
}

func (p *PressureTester) importFriends(friendSenderUserIDs, recvMsgUserIDs []string) error {
	for _, recvMsgUserID := range recvMsgUserIDs {
		if err := p.friendManager.ImportFriends(recvMsgUserID, friendSenderUserIDs); err != nil {
			return err
		}
	}
	return nil
}
