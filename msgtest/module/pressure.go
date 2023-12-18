package module

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
)

var (
	TESTIP        = "39.108.141.92"
	APIADDR       = fmt.Sprintf("http://%v:20002", TESTIP)
	WSADDR        = fmt.Sprintf("ws://%v:20001", TESTIP)
	SECRET        = "openIM123"
	MANAGERUSERID = "openIMAdmin"

	PLATFORMID = constant.AndroidPlatformID
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

	FastenedUserPrefix  = "f"
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
	sendNum        atomic.Int64

	groupSenderUserIDs, friendSenderUserIDs, notfriendSenderUserIDs []string
	recvMsgUserIDs                                                  []string

	tenThousandGroupIDs, thousandGroupIDs, hundredGroupUserIDs, fiftyGroupUserIDs []string
}

func (p *PressureTester) GetSendNum() int64 {
	return p.sendNum.Load()
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
		sampleReceiver = append(sampleReceiver, fmt.Sprintf("%s_testv3_%d", FastenedUserPrefix, i))
	}
	SampleUserList = sampleReceiver
	return fastenedUserIDs, sampleReceiver, nil

}
func (p *PressureTester) SelectStartAndEnd(start, end int) (fastenedUserIDs []string) {
	return p.userManager.GenSEUserIDsWithPrefix(start, end, FastenedUserPrefix)
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

func (p *PressureTester) CreateTestGroups(userIDs, fastenedUserIDs, recvMsgUserIDs []string) (err error) {
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

// func (p *PressureTester) SendSingleMessages(fastenedUserIDs []string, num int, duration time.Duration) {
// 	var wg sync.WaitGroup
// 	length := len(fastenedUserIDs)
// 	rand.Seed(time.Now().UnixNano())
// 	for i, userID := range fastenedUserIDs {
// 		counter:=0
// 		for counter < num {
// 			index := rand.Intn(length)
// 			if index != i {
// 				counter++
// 				wg.Add(1)
// 				go func(reciver string,sender string,counter int) {
// 					defer wg.Done()
// 					if user, ok := p.msgSender[sender]; ok {
// 						user.SendMsgWithContext(reciver, counter)
// 					}
// 					time.Sleep(duration)
// 				}(fastenedUserIDs[index],userID,counter)
// 			}
// 		}

// 	}
// 	wg.Wait()

// }
func (p *PressureTester) SendSingleMessages(fastenedUserIDs []string, randomSender, num int, duration time.Duration) {
	var wg sync.WaitGroup
	length := len(fastenedUserIDs)
	rand.Seed(time.Now().UnixNano())
	for i, userID := range fastenedUserIDs {
		counter := 0
		var receiverUserIDs []string
		for counter < num {
			index := rand.Intn(length)
			if index != i {
				counter++
				receiverUserIDs = append(receiverUserIDs, fastenedUserIDs[index])
			}
		}
		wg.Add(1)
		go func(receiverUserIDs []string, u string) {
			//log.ZError(context.Background(), "SendSingleMessages", nil, "length", len(receiverUserIDs))
			defer wg.Done()
			user, _ := p.msgSender[u]
			for j, rv := range receiverUserIDs {
				user.SendMsgWithContext(rv, j)
				p.sendNum.Add(1)

				time.Sleep(duration)

			}
		}(receiverUserIDs, userID)
	}
	wg.Wait()

}
func (p *PressureTester) SendSingleMessages2(fastenedUserIDs []string, randomSender []string, randomReceiver, num int, duration time.Duration) {
	var wg sync.WaitGroup
	length := len(fastenedUserIDs)
	rand.Seed(time.Now().UnixNano())
	for _, userID := range randomSender {
		counter := 0
		var receiverUserIDs []string
		for counter < randomReceiver {
			index := rand.Intn(length)
			if fastenedUserIDs[index] != userID {
				counter++
				receiverUserIDs = append(receiverUserIDs, fastenedUserIDs[index])
			}
		}
		wg.Add(1)
		go func(receiverUserIDs []string, u string) {
			//log.ZError(context.Background(), "SendSingleMessages", nil, "length", len(receiverUserIDs))
			defer wg.Done()
			user, _ := p.msgSender[u]
			for _, rv := range receiverUserIDs {
				for x := 0; x < num; x++ {
					user.SendMsgWithContext(rv, x)
					p.sendNum.Add(1)

					time.Sleep(duration)
				}

			}
		}(receiverUserIDs, userID)
	}
	wg.Wait()

}
func (p *PressureTester) Shuffle(fastenedUserIDs []string, needNum int) []string {
	// 使用洗牌算法对 fastenedUserIDs 进行随机排序
	rand.Shuffle(len(fastenedUserIDs), func(i, j int) {
		fastenedUserIDs[i], fastenedUserIDs[j] = fastenedUserIDs[j], fastenedUserIDs[i]
	})

	// 选取前100个不重复的 userID
	selectedUserIDs := make([]string, 0, needNum)
	seen := make(map[string]bool)

	for _, userID := range fastenedUserIDs {
		if len(selectedUserIDs) == needNum {
			break
		}

		if !seen[userID] {
			selectedUserIDs = append(selectedUserIDs, userID)
			seen[userID] = true
		}
	}
	return selectedUserIDs
}

func (p *PressureTester) SendSingleMessagesTo(fastenedUserIDs []string, num int, duration time.Duration) {
	var wg sync.WaitGroup
	//length := len(fastenedUserIDs)
	rand.Seed(time.Now().UnixNano())
	for i, userID := range fastenedUserIDs {
		//counter := 0
		//var receiverUserIDs []string
		//for counter < num {
		//	index := rand.Intn(length)
		//	if index != i {
		//		counter++
		//		receiverUserIDs = append(receiverUserIDs, fastenedUserIDs[index])
		//	}
		//}
		var receiverUserIDs []string
		for i < num {
			receiverUserIDs = append(receiverUserIDs, utils.IntToString(i))
			i++
		}
		wg.Add(1)
		go func(receiverUserIDs []string, u string) {
			defer wg.Done()
			user, _ := p.msgSender[u]
			for j, rv := range receiverUserIDs {
				user.SendMsgWithContext(rv, j)
				time.Sleep(duration)

			}
		}(receiverUserIDs, userID)
	}
	wg.Wait()

}

func (p *PressureTester) SendMessages(sendID, recvID string, msgNum int) {
	var i = 0
	var ws sync.WaitGroup
	user, _ := p.msgSender[sendID]
	for i < msgNum {
		ws.Add(1)
		i++
		go func() {
			defer ws.Done()
			user.SendMsgWithContext(recvID, i)
		}()

	}
	ws.Wait()

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
	log.ZWarn(ctx, "message send finished start to check message", nil)
	var max, min, latencySum int64
	samepleReceiverFailedMap := make(map[string]*errorValue)
	failedMessageAllMap := make(map[string]*errorValue)
	sendSampleMessageAllMap := make(map[string]*msgValue)
	recvSampleMessageAllMap := make(map[string]*msgValue)
	var sampleSendLength, sampleRecvLength, failedMessageLength int
	for _, user := range p.msgSender {
		if len(user.failedMessageMap) != 0 {
			failedMessageLength += len(user.failedMessageMap)
			for s, value := range user.failedMessageMap {
				failedMessageAllMap[s] = value
				if utils.IsContain(value.RecvID, SampleUserList) {
					samepleReceiverFailedMap[s] = value
				}
			}
		}
		if len(user.sendSampleMessage) != 0 {
			sampleSendLength += len(user.sendSampleMessage)
			for s, value := range user.sendSampleMessage {
				sendSampleMessageAllMap[s] = value
			}
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
			for s, value := range user.recvSampleMessage {
				recvSampleMessageAllMap[s] = value
			}
		}
	}
	log.ZError(context.Background(), "check result", nil, "failedMessageLength", failedMessageLength,
		"sampleSendLength", sampleSendLength, "sampleRecvLength", sampleRecvLength, "Average of message latency",
		utils.Int64ToString(latencySum/int64(sampleRecvLength))+" ms", "max", utils.Int64ToString(max)+" ms",
		"min", utils.Int64ToString(min)+" ms")
	if len(failedMessageAllMap) > 0 {
		err := p.saveFailedMessageToFile(failedMessageAllMap, "failedMessageAllMap")
		if err != nil {
			log.ZWarn(ctx, "save failed message to file failed", err)
		}
	}

	if len(samepleReceiverFailedMap) > 0 {
		err := p.saveFailedMessageToFile(failedMessageAllMap, "sampleReceiverFailedMap")
		if err != nil {
			log.ZWarn(ctx, "save sampleReceiverFailedMap message to file failed", err)
		}
	}
	if sampleSendLength != sampleRecvLength {
		recvEx, sendEx := findMapIntersection(sendSampleMessageAllMap, recvSampleMessageAllMap)
		if len(recvEx) != 0 {
			p.saveSucessMessageToFile(recvEx, "recvAdditional")
		}
		if len(sendEx) != 0 {
			p.saveSucessMessageToFile(recvEx, "sendAdditional")
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
func (p *PressureTester) saveSucessMessageToFile(m map[string]*msgValue, filename string) error {
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

func findMapIntersection(map1, map2 map[string]*msgValue) (map[string]*msgValue, map[string]*msgValue) {
	InMap1NotInMap2 := make(map[string]*msgValue)
	InMap2NotInMap1 := make(map[string]*msgValue)
	for key := range map1 {
		if _, ok := map2[key]; !ok {
			InMap1NotInMap2[key] = map1[key]
		}
	}
	for key := range map2 {
		if _, ok := map1[key]; !ok {
			InMap2NotInMap1[key] = map2[key]
		}
	}

	return InMap1NotInMap2, InMap2NotInMap1
}
