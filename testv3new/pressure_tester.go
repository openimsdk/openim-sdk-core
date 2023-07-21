package testv3new

import (
	"context"
	"open_im_sdk/testv3new/testcore"
	"sync"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
)

type PressureTester struct {
	sendLightWeightSDKCores map[string]*testcore.BaseCore
	recvLightWeightSDKCores map[string]*testcore.BaseCore

	registerManager RegisterManager
}

func NewPressureTester() *PressureTester {
	return &PressureTester{sendLightWeightSDKCores: map[string]*testcore.BaseCore{}, recvLightWeightSDKCores: map[string]*testcore.BaseCore{}}
}

func (p *PressureTester) InitSendCores(userIDs []string) {
	for _, userID := range userIDs {
		p.sendLightWeightSDKCores[userID] = testcore.NewBaseCore(userID)
	}
}

func (p *PressureTester) InitRecvCores(userIDs []string) {
	for _, userID := range userIDs {
		p.recvLightWeightSDKCores[userID] = testcore.NewBaseCore(userID)
	}
}

func (p *PressureTester) PressureSendMsgs(recvUserIDs []string, msgNum int, duration time.Duration) {
}

func (p *PressureTester) PressureSendGroupMsgs(groupID string, msgNum int, duration time.Duration) {

}

func (p *PressureTester) MsgReliabilityTest(sendUserID, recvUserID string, msgNum int, duration time.Duration) {
	sendCore := p.sendLightWeightSDKCores[sendUserID]
	recvCore := p.recvLightWeightSDKCores[recvUserID]
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i := 0; i < msgNum; i++ {
			sendCore.SendMsg(i)
		}
		wg.Done()
	}()
	wg.Wait()
	log.ZInfo(context.Background(), "send msg done", "reliability", recvCore.GetRecvMsgNum() == msgNum)

}
