package module

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"math/rand"
	"sync"
	"time"
)

var (
	TESTIP        = "14.29.168.56"
	APIADDR       = fmt.Sprintf("http://%v:10002", TESTIP)
	WSADDR        = fmt.Sprintf("ws://%v:10001", TESTIP)
	SECRET        = "openIM123"
	MANAGERUSERID = "openIMAdmin"

	PLATFORMID = constant.WindowsPlatformID
	LogLevel   = uint32(5)

	REGISTERADDR = APIADDR + constant.UserRegister
	TOKENADDR    = APIADDR + constant.GetUsersToken
)

var (
	totalOnlineUserNum    = 200000 // 总在线用户数
	friendMsgSenderNum    = 200    // 好友消息发送者数
	NotFriendMsgSenderNum = 200    // 非好友消息发送者数
	groupMsgSenderNum     = 200    // 群消息发送者数
	msgSenderNumEvreyUser = 100    // 每个用户的消息数
	fastenedUserNum       = 600    // 固定用户数

	recvMsgUserNum = 20 // 消息接收者数, 抽样账号
	SampleUserList []string
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

	FastenedUserPrefix  = "fastened_user_prefix"
	RecvMsgPrefix       = "recv_msg_prefix"
	singleMsgRecvPrefix = "single_msg_recv_prefix"
)

type PressureTester struct {
	friendManager  *TestFriendManager
	userManager    *TestUserManager
	groupManager   *TestGroupManager
	msgSender      map[string]*SendMsgUser
	groupMsgSender map[string]*SendMsgUser
	timeOffset     int64

	groupSenderUserIDs, friendSenderUserIDs, notfriendSenderUserIDs []string
	recvMsgUserIDs                                                  []string

	tenThousandGroupIDs, thousandGroupIDs, hundredGroupUserIDs, fiftyGroupUserIDs []string
}

func NewPressureTester() *PressureTester {
	metaManager := NewMetaManager(APIADDR, SECRET, MANAGERUSERID)
	metaManager.initToken()
	serverTime, err := metaManager.GetServerTime()
	if err != nil {
		panic(err)
	}
	return &PressureTester{friendManager: metaManager.NewFriendManager(), userManager: metaManager.NewUserManager(), groupManager: metaManager.NewGroupMananger(),
		msgSender: make(map[string]*SendMsgUser), groupMsgSender: make(map[string]*SendMsgUser), timeOffset: serverTime - utils.GetCurrentTimestampByMill()}
}

func (p *PressureTester) genUserIDs() (userIDs, fastenedUserIDs, recvMsgUserIDs []string) {
	userIDs = p.userManager.GenUserIDs(totalOnlineUserNum - fastenedUserNum)                  // 在线用户
	fastenedUserIDs = p.userManager.GenUserIDsWithPrefix(fastenedUserNum, FastenedUserPrefix) // 指定发消息的固定用户
	recvMsgUserIDs = p.userManager.GenUserIDsWithPrefix(recvMsgUserNum, RecvMsgPrefix)        // 抽样用户完整SDK
	return
}

// selectSample
func (p *PressureTester) SelectSample(total int, percentage float64) (fastenedUserIDs []string,
	sampleReceiver []string, err error) {
	if percentage < 0 || percentage > 1 {
		return nil, nil, fmt.Errorf("percentage must be between 0 and 1")
	}
	fastenedUserIDs = p.userManager.GenUserIDsWithPrefix(total, FastenedUserPrefix)
	step := int(1.0 / percentage)
	for i := 0; i <= total; i += step {
		sampleReceiver = append(sampleReceiver, fmt.Sprintf("%s_testv3new_%d", FastenedUserPrefix, i))
	}
	SampleUserList = sampleReceiver
	return fastenedUserIDs, sampleReceiver, nil

}

func (p *PressureTester) RegisterUsers(userIDs []string, fastenedUserIDs []string, recvMsgUserIDs []string) error {
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

func (p *PressureTester) InitUserConns(userIDs []string) {
	for _, userID := range userIDs {
		token, err := p.userManager.GetToken(userID, int32(PLATFORMID))
		if err != nil {
			log.ZError(context.Background(), "get token failed", err, "userID", userID, "platformID", PLATFORMID)
			continue
		}
		user := NewUser(userID, token, p.timeOffset, sdk_struct.IMConfig{WsAddr: WSADDR, ApiAddr: APIADDR, PlatformID: int32(PLATFORMID)})
		p.msgSender[userID] = user

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

func (p *PressureTester) SendSingleMessages(fastenedUserIDs []string, num int, duration time.Duration) {
	var wg sync.WaitGroup
	length := len(fastenedUserIDs)
	rand.Seed(time.Now().UnixNano())
	for _, userID := range fastenedUserIDs {
		receiverUserIDs := make([]string, 100)
		for len(receiverUserIDs) < num {
			index := rand.Intn(length)
			if fastenedUserIDs[index] != userID {
				receiverUserIDs = append(receiverUserIDs, fastenedUserIDs[index])
			}
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j, rv := range receiverUserIDs {
				if user, ok := p.msgSender[userID]; ok {
					user.SendMsgWithContext(rv, j)
				}
				time.Sleep(duration)

			}
		}()
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

func (p *PressureTester) CheckMsg() {
	var sampleSendLength, sampleRecvLength, failedMessageLength int
	for _, user := range p.msgSender {
		if len(user.failedMessageMap) != 0 {
			failedMessageLength += len(user.failedMessageMap)
		}
		if len(user.sendSampleMessage) != 0 {
			sampleSendLength += len(user.sendSampleMessage)
		}
		if len(user.recvSampleMessage) != 0 {
			sampleRecvLength += len(user.recvSampleMessage)
		}
	}
	log.ZDebug(context.Background(), "check result", "failedMessageLength", failedMessageLength,
		"sampleSendLength", sampleSendLength, "sampleRecvLength", sampleRecvLength)
}
