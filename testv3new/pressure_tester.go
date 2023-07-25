package testv3new

import (
	"context"
	"fmt"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/sdk_struct"
	"open_im_sdk/testv3new/testcore"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
)

type PressureTester struct {
	sendLightWeightSDKCores map[string]*testcore.BaseCore
	recvLightWeightSDKCores map[string]*testcore.BaseCore

	registerManager RegisterManager
	apiAddr         string
	wsAddr          string
}

func NewPressureTester(apiAddr, wsAddr string) *PressureTester {
	return &PressureTester{
		sendLightWeightSDKCores: map[string]*testcore.BaseCore{},
		recvLightWeightSDKCores: map[string]*testcore.BaseCore{},
		registerManager:         *NewRegisterManager(),
		apiAddr:                 apiAddr,
		wsAddr:                  wsAddr,
	}
}

func NewCtx(apiAddr, wsAddr, userID, token string) context.Context {
	return ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
		UserID: userID, Token: token,
		IMConfig: sdk_struct.IMConfig{
			PlatformID:          constant.AndroidPlatformID,
			ApiAddr:             apiAddr,
			WsAddr:              wsAddr,
			LogLevel:            2,
			IsLogStandardOutput: true,
			LogFilePath:         "./",
		}})
}

func (p *PressureTester) initCores(m *map[string]*testcore.BaseCore, userIDs []string) {
	for _, userID := range userIDs {
		token, err := p.registerManager.GetToken(userID)
		if err != nil {
			log.ZError(context.Background(), "get token error", err, "userID", userID)
			continue
		}
		mV := *m
		mV[userID] = testcore.NewBaseCore(NewCtx(p.apiAddr, p.wsAddr, userID, token), userID)
	}
}

func (p *PressureTester) InitSendCores(userIDs []string) {
	p.initCores(&p.sendLightWeightSDKCores, userIDs)
}

func (p *PressureTester) InitRecvCores(userIDs []string) {
	p.initCores(&p.recvLightWeightSDKCores, userIDs)
}

// user single chat send msg pressure test
func (p *PressureTester) PressureSendMsgs(sendUserID string, recvUserIDs []string, msgNum int, duration time.Duration) {
	// 每秒发送多少条消息
	ctx, _ := InitContext(sendUserID)
	p.WithTimer(p.InitSendCores)([]string{sendUserID})
	p.WithTimer(p.InitRecvCores)(recvUserIDs)

	sendCore := p.sendLightWeightSDKCores[sendUserID]

	fmt.Println("\nmsgNum ==> ", msgNum)
	var wg sync.WaitGroup
	wg.Add(len(recvUserIDs))
	for idx, recvUserID := range recvUserIDs {
		go func(idx int) {
			for i := 1; i <= msgNum; i++ {
				p.WithTimer(sendCore.SendSingleMsg)(ctx, recvUserID, i)
				// if err := sendCore.SendSingleMsg(ctx, recvUserID, i); err != nil {
				// 	log.ZError(ctx, "send msg error", err, "index", i, "recvUserID", recvUserID, "sendUserID", sendUserID)
				// }
				time.Sleep(duration)
			}
			recvCore := p.recvLightWeightSDKCores[recvUserID]
			if recvCore != nil {
				time.Sleep(100 * time.Millisecond)
				recvMap := recvCore.GetRecvMap()
				if recvMap != nil {
					count := recvMap[sendUserID+"_"+recvUserID]
					fmt.Println(fmt.Sprintf("recvUserID: %v ==> recv msg num: %d %v", recvUserID, count, count == msgNum))
					log.ZInfo(ctx, "recv msg", "recv num", count, "recvUserID", recvUserID, "recv status", count == msgNum)
				}
			}
			wg.Done()
		}(idx)
	}
	wg.Wait()
}

// group chat send msg pressure test
func (p *PressureTester) PressureSendGroupMsgs(sendUserIDs []string, groupID string, msgNum int, duration time.Duration) {
	if resp, err := p.GetGroupMembersInfo(groupID, sendUserIDs); err != nil {
		log.ZError(context.Background(), "get group members info failed", err)
		return
	} else if resp.Members != nil {
		log.ZError(context.Background(), "get group members info failed", err, "userIDs", sendUserIDs)
		return
	}

	startTime := time.Now().UnixNano()
	p.InitSendCores(sendUserIDs)
	endTime := time.Now().UnixNano()
	fmt.Println("bantanger init send cores time:", float64(endTime-startTime))
	// p.InitRecvCores([]string{groupID})
	// 管理员邀请进群
	err := p.InviteUserToGroup(groupID, sendUserIDs)
	if err != nil {
		return
	}

	for _, sendUserID := range sendUserIDs {
		ctx, _ := InitContext(sendUserID)
		sendCore := p.sendLightWeightSDKCores[sendUserID]
		for i := 1; i <= msgNum; i++ {
			time.Sleep(duration)
			if err := sendCore.SendGroupMsg(ctx, groupID, i); err != nil {
				log.ZError(ctx, "send msg error", err, "index", i, "recvUserID", groupID, "sendUserID", sendUserID)
			}
		}
	}
}

// msg ordering test
func (p *PressureTester) OrderingSendMsg(groupID string, msgNum int) {

}

// reliability test
func (p *PressureTester) MsgReliabilityTest(sendUserID, recvUserID string, msgNum int, duration time.Duration) {
	ctx, _ := InitContext(sendUserID)
	sendCore := p.sendLightWeightSDKCores[sendUserID]

	for i := 0; i < msgNum; i++ {
		if err := sendCore.SendSingleMsg(ctx, recvUserID, i); err != nil {
			log.ZError(ctx, "send msg error", err, "index", i, "recvUserID", recvUserID, "sendUserID", sendUserID)
		}
	}
	// recvCore := p.recvLightWeightSDKCores[recvUserID]
	// log.ZInfo(context.Background(), "send msg done", "reliability", recvCore.GetRecvMsgNum() == msgNum)
}

// WithTimer Decorative function, accept a function as a parameter, and return a packaging function
func (p *PressureTester) WithTimer(f interface{}) func(...interface{}) {
	return func(args ...interface{}) {
		start := time.Now()
		v := reflect.ValueOf(f)
		if v.Kind() != reflect.Func {
			log.ZError(context.Background(), "pass parameter is not a function", nil,
				"actual", v.Kind(), "expected", reflect.Func)
			return
		}
		funcName := runtime.FuncForPC(v.Pointer()).Name()
		var in []reflect.Value
		for _, arg := range args {
			in = append(in, reflect.ValueOf(arg))
		}
		v.Call(in) // 执行原始函数
		duration := time.Since(start)
		fmt.Printf("Run Funcation: %s\nRun Funcation spent time: %s\n", funcName, duration)
	}
}
