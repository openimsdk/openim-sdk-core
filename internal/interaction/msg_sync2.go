package interaction

import (
	"context"
	"fmt"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"sync"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
)

type MsgSyncer struct {
	conversationCh     chan common.Cmd2Value
	pushMsgAndMaxSeqCh chan common.Cmd2Value
	ctx                context.Context
	minSeq             int64
	maxSeq             int64
	seqMap             map[string]int64

	lock sync.Mutex
}

func NewMsgSyncer(conversationCh, pushMsgAndMaxSeqCh chan common.Cmd2Value) *MsgSyncer {
	return &MsgSyncer{
		conversationCh:     conversationCh,
		pushMsgAndMaxSeqCh: pushMsgAndMaxSeqCh,
	}
}

func (m *MsgSyncer) GetCh() chan common.Cmd2Value {
	return m.pushMsgAndMaxSeqCh
}

// from pushMsgAndMaxSeqCh
func (m *MsgSyncer) Work(cmd common.Cmd2Value) {
	switch cmd.Cmd {
	case constant.CmdPushMsg:
		m.doPushMsg(cmd)
	case constant.CmdMaxSeq:
		m.doMaxSeq(cmd)
	default:
		log.ZError(cmd.Ctx, "inviald cmd error", fmt.Errorf("inviald cmd %s", cmd.Cmd))
	}
}

func (m *MsgSyncer) DoListener() {
	for {
		select {
		case cmd := <-m.GetCh():
			m.Work(cmd)
		case <-m.ctx.Done():
			return
		}
	}
}

func (m *MsgSyncer) doPushMsg(ctx context.Context) {

}

func (m *MsgSyncer) doMaxSeq(ctx context.Context) {

}

// 内部状态保持正确
func (m *MsgSyncer) SyncMsg(ctx context.Context, sourceID string, sessionType int32) error {
	return nil
}

func (m *MsgSyncer) triggerCmdNewMsgCome(ctx context.Context) error {
	return nil
}
