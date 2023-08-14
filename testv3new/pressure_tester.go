package testv3new

import (
	"context"
	"fmt"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"open_im_sdk/testv3new/testcore"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mcontext"
)

type PressureTester struct {
	cores map[string]*testcore.BaseCore

	testUserMananger TestUserManager
	platformID       int32
	apiAddr          string
	wsAddr           string

	adminUserID string
}

func NewPressureTester(apiAddr, wsAddr string, secret, adminUserID string) *PressureTester {
	return &PressureTester{
		cores:            map[string]*testcore.BaseCore{},
		testUserMananger: *NewTestUserManager(secret),
		apiAddr:          apiAddr,
		wsAddr:           wsAddr,
		adminUserID:      adminUserID,
		platformID:       int32(PLATFORMID),
	}
}

func InitContext(userID string, platformID int32) context.Context {
	config := ccontext.GlobalConfig{
		UserID: userID, Token: "",
		IMConfig: sdk_struct.IMConfig{
			PlatformID: platformID,
			ApiAddr:    APIADDR,
			WsAddr:     WSADDR,
		}}
	ctx := ccontext.WithInfo(context.Background(), &config)
	ctx = ccontext.WithOperationID(ctx, utils.OperationIDGenerator())
	return ctx
}

func (p *PressureTester) NewAdminCtx() context.Context {
	ctx := p.testUserMananger.NewCtx()
	ctx = mcontext.SetOperationID(ctx, utils.OperationIDGenerator())
	token, err := p.testUserMananger.GetToken(ctx, p.adminUserID, constant.WindowsPlatformID)
	if err != nil {
		panic(err)
	}
	return ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
		UserID: p.adminUserID,
		Token:  token,
		IMConfig: sdk_struct.IMConfig{
			PlatformID: p.platformID,
			ApiAddr:    p.apiAddr,
			WsAddr:     p.wsAddr,
		}})
}

func (p *PressureTester) NewCtx(userID, token string) context.Context {
	return ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
		UserID: userID,
		Token:  token,
		IMConfig: sdk_struct.IMConfig{
			PlatformID: p.platformID,
			ApiAddr:    p.apiAddr,
			WsAddr:     p.wsAddr,
		}})
}

func (p *PressureTester) InitCores(userIDs []string) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	for _, userID := range userIDs {
		wg.Add(1)
		go func(userID string) {
			defer wg.Done()
			ctx := p.NewAdminCtx()

			token, err := p.testUserMananger.GetToken(ctx, userID, p.platformID)
			if err != nil {
				log.ZError(context.Background(), "get token error", err, "userID", userID)
				return
			}
			mutex.Lock()
			p.cores[userID] = testcore.NewBaseCore(p.NewCtx(userID, token), userID, p.platformID)
			mutex.Unlock()
		}(userID)
	}
	wg.Wait()
}

// PressureSendMsgs user single chat send msg pressure test
func (p *PressureTester) PressureSendMsgs(ctx context.Context, sendUserID string, recvUserIDs []string, msgNum int, duration time.Duration) {
	var wg sync.WaitGroup
	wg.Add(len(recvUserIDs))
	for _, recvUserID := range recvUserIDs {
		go func(recvUserID string) {
			defer wg.Done() // Mark this goroutine as done when finished

			// Create a new context for each goroutine to avoid shared state
			// ctx, _ := InitContext(sendUserID)

			// Send messages concurrently
			var sendWG sync.WaitGroup
			sendWG.Add(msgNum)
			for i := 1; i <= msgNum; i++ {
				go func(i int) {
					defer sendWG.Done()
					// p.WithTimer(sendCore.SendSingleMsg)(ctx, recvUserID, i)
					// if err := sendCore.SendSingleMsg(ctx, recvUserID, i); err != nil {
					// 	log.ZError(ctx, "send msg error", err, "index", i, "recvUserID", recvUserID, "sendUserID", sendUserID)
					// }
				}(i)
			}
			sendWG.Wait()

			// Delay before querying the received messages
			time.Sleep(100 * time.Millisecond)

			// Query the received messages
			recvCore := p.cores[recvUserID]
			if recvCore != nil {
				// recvMap := recvCore.GetRecvMap()
				// if recvMap != nil {
				// 	count := recvMap[sendUserID+"_"+recvUserID]
				// 	fmt.Println(fmt.Sprintf("recvUserID: %v ==> recv msg num: %d %v", recvUserID, count, count == msgNum))
				// 	log.ZInfo(ctx, "recv msg", "recv num", count, "recvUserID", recvUserID, "recv status", count == msgNum)
				// }
			}
		}(recvUserID)
	}
	wg.Wait()
}

func (p *PressureTester) CreateConversations(ctx context.Context, conversationNum int, recvUserID string) error {
	userIDs := p.testUserMananger.GenUserIDs(conversationNum)
	if err := p.testUserMananger.RegisterUsers(ctx, userIDs...); err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, userID := range userIDs {
		time.Sleep(time.Millisecond * 100)
		token, _ := p.testUserMananger.GetToken(ctx, userID, p.platformID)
		ctx2 := NewUserCtx(userID, token)
		baseCore := testcore.NewBaseCore(ctx2, userID, p.platformID)
		ctx2 = mcontext.SetOperationID(ctx2, utils.OperationIDGenerator())
		if err := baseCore.SendSingleMsg(ctx2, recvUserID, 0); err != nil {
			log.ZError(ctx2, "send msg error", err, "sendUserID", userID)
		}
	}
	wg.Add(1)
	wg.Wait()
	return nil
}
func (p *PressureTester) CreateConversationsAndBatchSendMsg(ctx context.Context, conversationNum int, onePeopleMessageNum int,
	recvUserID string, fixedUserIDs []string) error {
	var userIDs []string
	if len(fixedUserIDs) <= 0 {
		userIDs = p.testUserMananger.GenUserIDs(conversationNum)
		if err := p.testUserMananger.RegisterUsers(ctx, userIDs...); err != nil {
			return err
		}
	} else {
		userIDs = fixedUserIDs
	}

	var wg sync.WaitGroup
	for _, userID := range userIDs {
		go func() {
			time.Sleep(time.Millisecond * 100)
			token, _ := p.testUserMananger.GetToken(ctx, userID, p.platformID)
			ctx2 := NewUserCtx(userID, token)
			baseCore := testcore.NewBaseCore(ctx2, userID, p.platformID)
			ctx2 = mcontext.SetOperationID(ctx2, utils.OperationIDGenerator())
			for i := 0; i < onePeopleMessageNum; i++ {
				if err := baseCore.BatchSendSingleMsg(ctx2, recvUserID, i); err != nil {
					log.ZError(ctx2, "send msg error", err, "sendUserID", userID)
				}
				time.Sleep(time.Millisecond * 5000)
			}
		}()

	}
	wg.Add(1)
	wg.Wait()
	return nil
}
func (p *PressureTester) CreateConversationsAndBatchSendGroupMsg(ctx context.Context, conversationNum int, onePeopleMessageNum int,
	groupID string, fixedUserIDs []string) error {
	var userIDs []string
	if len(fixedUserIDs) <= 0 {
		userIDs = p.testUserMananger.GenUserIDs(conversationNum)
		if err := p.testUserMananger.RegisterUsers(ctx, userIDs...); err != nil {
			return err
		}
	} else {
		userIDs = fixedUserIDs
	}

	var wg sync.WaitGroup
	for _, userID := range userIDs {
		go func(u string) {
			log.ZDebug(ctx, "start send msg", "userID", u)
			time.Sleep(time.Millisecond * 100)
			token, _ := p.testUserMananger.GetToken(ctx, u, p.platformID)
			ctx2 := NewUserCtx(u, token)
			baseCore := testcore.NewBaseCore(ctx2, u, p.platformID)
			ctx2 = mcontext.SetOperationID(ctx2, utils.OperationIDGenerator())
			for i := 0; i < onePeopleMessageNum; i++ {
				if err := baseCore.BatchSendGroupMsg(ctx2, groupID, i); err != nil {
					log.ZError(ctx2, "send msg error", err, "sendUserID", u)
				}
				time.Sleep(time.Millisecond * 500)
			}
		}(userID)

	}
	wg.Add(1)
	wg.Wait()
	return nil
}

// PressureSendMsgs2 user single chat send msg pressure test
func (p *PressureTester) PressureSendMsgs2(ctx context.Context, sendUserIDs []string, recvUserIDs []string, msgNum int, duration time.Duration) {
	var wg sync.WaitGroup
	for _, sendUserID := range sendUserIDs {
		// ctx, _ := InitContext(sendUserID)
		sendCore := p.cores[sendUserID]
		if sendCore == nil {
			log.ZInfo(ctx, "sendCore is nil", "sendUserID", sendUserID)
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			// for _, recvUserID := range recvUserIDs {
			// 	for i := 1; i <= msgNum; i++ {
			// 		p.sendCore.SendSingleMsg(ctx, recvUserID, i)
			// 		time.Sleep(duration)
			// 	}
			// }
		}()
	}
	wg.Wait()
	//close(msgChan)
}

// PressureSendGroupMsgs group chat send msg pressure test
func (p *PressureTester) PressureSendGroupMsgs(ctx context.Context, sendUserIDs []string, groupID string, msgNum int, duration time.Duration) {
	if resp, err := p.testUserMananger.GetGroupMembersInfo(ctx, groupID, sendUserIDs); err != nil {
		log.ZError(context.Background(), "get group members info failed", err)
		return
	} else if resp.Members != nil {
		log.ZError(context.Background(), "get group members info failed", err, "userIDs", sendUserIDs)
		return
	}
	startTime := time.Now().UnixNano()
	p.InitCores(sendUserIDs)
	endTime := time.Now().UnixNano()
	fmt.Println("bantanger init send cores time:", float64(endTime-startTime))
	// 管理员邀请进群
	err := p.testUserMananger.InviteUserToGroup(ctx, groupID, sendUserIDs)
	if err != nil {
		return
	}

	for _, sendUserID := range sendUserIDs {
		// ctx, _ := InitContext(sendUserID)
		sendCore := p.cores[sendUserID]
		for i := 1; i <= msgNum; i++ {
			time.Sleep(duration)
			if err := sendCore.SendGroupMsg(ctx, groupID, i); err != nil {
				log.ZError(ctx, "send msg error", err, "index", i, "recvUserID", groupID, "sendUserID", sendUserID)
			}
		}
	}
}

// PressureSendGroupMsgs group chat send msg pressure test
func (p *PressureTester) PressureSendGroupMsgs2(ctx context.Context, sendUserIDs []string, groupIDs []string, msgNum int, duration time.Duration) {
	for _, groupID := range groupIDs {
		if resp, err := p.testUserMananger.GetGroupMembersInfo(ctx, groupID, sendUserIDs); err != nil {
			log.ZError(context.Background(), "get group members info failed", err)
			return
		} else if resp.Members != nil {
			log.ZError(context.Background(), "get group members info failed", err, "userIDs", sendUserIDs)
			return
		}

		startTime := time.Now().UnixNano()
		p.InitCores(sendUserIDs)
		endTime := time.Now().UnixNano()
		fmt.Println("bantanger init send cores time:", float64(endTime-startTime))

		// 管理员邀请进群
		err := p.testUserMananger.InviteUserToGroup(ctx, groupID, sendUserIDs)
		if err != nil {
			return
		}

		for _, sendUserID := range sendUserIDs {
			// ctx, _ := InitContext(sendUserID)
			sendCore := p.cores[sendUserID]
			for i := 1; i <= msgNum; i++ {
				time.Sleep(duration)
				if err := sendCore.SendGroupMsg(ctx, groupID, i); err != nil {
					log.ZError(ctx, "send msg error", err, "index", i, "recvUserID", groupID, "sendUserID", sendUserID)
				}
			}
		}
	}
}

// OrderingSendMsg msg ordering test
func (p *PressureTester) OrderingSendMsg(groupID string, msgNum int) {

}

// MsgReliabilityTest reliability test
func (p *PressureTester) MsgReliabilityTest(ctx context.Context, recvUserID, sendUserID string, msgNum int, duration time.Duration) {
	// ctx, _ := InitContext(sendUserID)
	sendCore := p.cores[sendUserID]

	for i := 0; i < msgNum; i++ {
		if err := sendCore.SendSingleMsg(ctx, recvUserID, i); err != nil {
			log.ZError(ctx, "send msg error", err, "index", i, "recvUserID", recvUserID, "sendUserID", sendUserID)
		}
	}
}

func (p *PressureTester) WithTimer(f interface{}) func(...interface{}) interface{} {
	return func(args ...interface{}) interface{} {
		start := time.Now().UnixNano()
		v := reflect.ValueOf(f)
		if v.Kind() != reflect.Func {
			log.ZError(context.Background(), "pass parameter is not a function", nil,
				"actual parameter", v.Kind(), "expected parameter", reflect.Func)
			return nil
		}
		funcName := runtime.FuncForPC(v.Pointer()).Name()
		var in []reflect.Value
		for _, arg := range args {
			in = append(in, reflect.ValueOf(arg))
		}
		call := v.Call(in) // Execute the original function
		end := time.Now().UnixNano()
		fmt.Printf("Execute Function: %s\nExecute Function spent time: %v\n", funcName, float64(end-start))

		// Get and return the first return value from the original function
		if len(call) > 0 {
			return call[0].Interface()
		}
		return nil
	}
}
