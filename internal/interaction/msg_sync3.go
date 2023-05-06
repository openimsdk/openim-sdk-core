package interaction

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/golang/protobuf/proto"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/sdk_struct"
)

type LongConnMgr struct {
}

func (c *LongConnMgr) SendReqWaitResp(ctx context.Context, m proto.Message, reqIdentifier int, resp proto.Message) error {
	return nil
}

// type MsgSyncer struct {
// 	loginUserID       string
// 	longConnMgr       *LongConnMgr
// 	pushMsgAndEventCh chan common.Cmd2Value
// 	conversationCh    chan common.Cmd2Value
// 	ctx               context.Context
// 	syncedMaxSeqs     map[string]uint64
// 	db                db_interface.DataBase
// }

// func NewMsgSyncer(ctx context.Context, conversationCh, pushMsgAndMaxSeqCh, recvSeqch chan common.Cmd2Value, loginUserID string, longConnMgr *LongConnMgr, db db_interface.DataBase) *MsgSyncer {
// 	p := &MsgSyncer{
// 		loginUserID:       loginUserID,
// 		longConnMgr:       longConnMgr,
// 		pushMsgAndEventCh: pushMsgAndMaxSeqCh,
// 		conversationCh:    conversationCh,
// 		ctx:               ctx,
// 		syncedMaxSeqs:     make(map[string]uint64),
// 		db:                db,
// 	}
// 	p.loadSeq()
// 	go p.DoListener()
// 	return p
// }

// set syncedMaxSeqs
func (m *MsgSyncer) loadSeq() error {
	return nil
}

func (m *MsgSyncer) DoListener() {
	for {
		select {
		case cmd := <-m.pushMsgAndEventCh:
			m.handlePushMsgAndEvent(cmd)
		case <-m.ctx.Done():
			log.ZInfo(m.ctx, "msg syncer done, sdk logout.....")
			return
		}
	}
}

func (m *MsgSyncer) handlePushMsgAndEvent(cmd common.Cmd2Value) {
	switch cmd.Cmd {
	case constant.CmdConnSuccesss:
		m.doConnected()
	case constant.CmdMaxSeq:
		m.doMaxSeq(cmd.Value.(*sdk_struct.CmdMaxSeqToMsgSync))
	case constant.CmdPushMsg:
		m.doPushMsg(cmd.Value.(*sdkws.PushMessages))
	}
}

// finishes a synchronization.
func (m *MsgSyncer) connected() error {
	req := m.GenReqForConnected()
	m.pullMsg(req)
}

func (m *MsgSyncer) doMaxSeq(seq *sdk_struct.CmdMaxSeqToMsgSync) error {
	req := m.GenReqForMaxSeq(seq)
	m.pullMsg(req)
}

func (m *MsgSyncer) doPushMsg(push *sdkws.PushMessages) error {
	req := m.GenReqForPushMsg(push)
	m.pullMsg(req)
}

// 重连成功后调用，同步最新的1条消息
func (m *MsgSyncer) doConnected() {
	//同步开始
	common.TriggerCmdNotification(m.ctx, sdk_struct.CmdNewMsgComeToConversation{SyncFlag: constant.MsgSyncBegin}, m.conversationCh)
	//同步
	req := m.GenReqForConnected()
	m.pullMsg(req)
	//同步结束
	common.TriggerCmdNotification(m.ctx, sdk_struct.CmdNewMsgComeToConversation{SyncFlag: constant.MsgSyncEnd}, m.conversationCh)

}

func (m *MsgSyncer) GenReqForConnected() *sdkws.PullMessageBySeqsReq {
	return nil

}

func (m *MsgSyncer) GenReqForMaxSeq(seq *sdk_struct.CmdMaxSeqToMsgSync) *sdkws.PullMessageBySeqsReq {
	return nil
}

func (m *MsgSyncer) GenReqForPushMsg(push *sdkws.PushMessages) *sdkws.PullMessageBySeqsReq {
	return nil
}

// 1 拉取消息  2 丢给conversation ch
func (m *MsgSyncer) pullMsg(req *sdkws.PullMessageBySeqsReq) error {
	var resp sdkws.PullMessageBySeqsResp
	m.longConnMgr.SendReqWaitResp(m.ctx, req, constant.PullMsgBySeqList, &resp)

	//通知->ch
	common.TriggerCmdNotification()
	//消息->ch
	common.TriggerCmdNewMsgCome()
	return nil
}
