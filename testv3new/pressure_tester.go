package testv3new

import (
	"context"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/sdk_struct"
	"open_im_sdk/testv3new/testcore"
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
	p.InitSendCores([]string{sendUserID})
	p.InitRecvCores(recvUserIDs)
	sendCore := p.sendLightWeightSDKCores[sendUserID]

	for _, recvUserID := range recvUserIDs {
		for i := 0; i < msgNum; i++ {
			time.Sleep(duration)
			if err := sendCore.SendSingleMsg(ctx, recvUserID, i); err != nil {
				log.ZError(ctx, "send msg error", err, "index", i, "recvUserID", recvUserID, "sendUserID", sendUserID)
			}
		}
	}
}

// group chat send msg pressure test
func (p *PressureTester) PressureSendGroupMsgs(groupID string, msgNum int, duration time.Duration) {

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
