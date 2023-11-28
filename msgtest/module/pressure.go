package module

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"math/rand"
	"os"
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
	log.ZWarn(context.Background(), "server time is", nil, "serverTime", serverTime, "current time",
		utils.GetCurrentTimestampByMill(), "time offset", serverTime-utils.GetCurrentTimestampByMill())

	return &PressureTester{friendManager: metaManager.NewFriendManager(), userManager: metaManager.NewUserManager(),
		groupManager: metaManager.NewGroupMananger(),
		msgSender:    make(map[string]*SendMsgUser), groupMsgSender: make(map[string]*SendMsgUser),
		timeOffset: serverTime - utils.GetCurrentTimestampByMill()}
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
	for i := 0; i < total; i += step {
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
		var receiverUserIDs []string
		for len(receiverUserIDs) < num {
			index := rand.Intn(length)
			if fastenedUserIDs[index] != userID {
				receiverUserIDs = append(receiverUserIDs, fastenedUserIDs[index])
			}
		}
		wg.Add(1)
		go func(receiverUserIDs []string) {
			defer wg.Done()
			for j, rv := range receiverUserIDs {
				if user, ok := p.msgSender[userID]; ok {
					user.SendMsgWithContext(rv, j)
				}
				time.Sleep(duration)

			}
		}(receiverUserIDs)
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

func (p *PressureTester) CheckMsg(ctx context.Context) {
	log.ZDebug(ctx, "message send finished start to check message")
	var max, min, latencySum int64
	failedMessageAllMap := make(map[string]*errorValue)
	var sampleSendLength, sampleRecvLength, failedMessageLength int
	for _, user := range p.msgSender {
		if len(user.failedMessageMap) != 0 {
			failedMessageLength += len(user.failedMessageMap)
			for s, value := range user.failedMessageMap {
				failedMessageAllMap[s] = value
			}
		}
		if len(user.sendSampleMessage) != 0 {
			sampleSendLength += len(user.sendSampleMessage)
		}
		if len(user.recvSampleMessage) != 0 {
			sampleRecvLength += len(user.recvSampleMessage)
			for _, value := range user.recvSampleMessage {
				if min == 0 && max == 0 {
					min = value.Latency
					max = value.Latency
				}
				if value.Latency < min {
					min = value.Latency
				}
				if value.Latency > max {
					max = value.Latency
				}
				latencySum += value.Latency
			}
		}
	}
	log.ZDebug(context.Background(), "check result", "failedMessageLength", failedMessageLength,
		"sampleSendLength", sampleSendLength, "sampleRecvLength", sampleRecvLength, "Average of message latency",
		utils.Int64ToString(latencySum/int64(sampleRecvLength))+" ms", "max", utils.Int64ToString(max)+" ms",
		"min", utils.Int64ToString(min)+" ms")
	if len(failedMessageAllMap) > 0 {
		err := p.saveFailedMessageToFile(failedMessageAllMap, "failedMessageAllMap")
		if err != nil {
			log.ZWarn(ctx, "save failed message to file failed", err)
		}
	}
	log.ZDebug(ctx, "message send finished start to check message")
	os.Exit(1)
}

func (p *PressureTester) saveFailedMessageToFile(m map[string]*errorValue, filename string) error {
	file, err := os.Create(filename + ".txt")
	if err != nil {
		return err
	}
	defer file.Close()

	for key, value := range m {
		line := fmt.Sprintf("Key: %s, Value: %v\n", key, value)
		_, err := file.WriteString(line)
		if err != nil {
			return err
		}
	}
	return nil
}
