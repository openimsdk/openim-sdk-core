package interaction

import (
	"context"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/sdk_struct"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/golang/protobuf/proto"
)

const (
	timeout         = 60
	retryTimes      = 2
	defaultPullNums = 1
)

type Seq struct {
	maxSeq      int64
	minSeq      int64
	sessionType int32
}

type SyncedSeq struct {
	maxSeqSynced int64
	sessionType  int32
}

// 回调同步开始 结束， 重连结束
type MsgSyncer struct {
	loginUserID string
	// listen ch
	ws        *Ws
	recvSeqch chan common.Cmd2Value // recv max seq from heartbeat map[string][2]int64
	recvMsgCh chan common.Cmd2Value // recv msg from ws, sync or push []*sdkws.MsgData

	conversationCh chan common.Cmd2Value // trigger conversation

	ctx  context.Context
	seqs map[string]SyncedSeq
	// msgCache  map[string]map[int64]*sdkws.MsgData
	db        db_interface.DataBase
	syncTimes int
}

func NewMsgSyncer(ctx context.Context, conversationCh, recvMsgCh, recvSeqch chan common.Cmd2Value,
	loginUserID string, ws *Ws) (*MsgSyncer, error) {
	m := &MsgSyncer{
		recvSeqch: recvSeqch,
		recvMsgCh: recvMsgCh,

		conversationCh: conversationCh,
		ws:             ws,
		loginUserID:    loginUserID,
		ctx:            ctx,
		// msgCache:       make(map[string]map[int64]*sdkws.MsgData),
	}
	err := m.loadSeq(ctx)
	return m, err
}

// seq db读取到内存
func (m *MsgSyncer) loadSeq(ctx context.Context) error {
	m.seqs = make(map[string]SyncedSeq)
	groupIDs, err := m.db.GetReadDiffusionGroupIDList(ctx)
	if err != nil {
		log.ZError(ctx, "get group id list failed", err)
		return err
	}
	for _, groupID := range groupIDs {
		nMaxSeq, err := m.db.GetSuperGroupNormalMsgSeq(ctx, groupID)
		if err != nil {
			log.ZError(ctx, "get group normal seq failed", err, "groupID", groupID)
			return err
		}
		aMaxSeq, err := m.db.GetSuperGroupAbnormalMsgSeq(ctx, groupID)
		if err != nil {
			log.ZError(ctx, "get group abnormal seq failed", err, "groupID", groupID)
			return err
		}
		var maxSeqSynced int64
		maxSeqSynced = nMaxSeq
		if aMaxSeq > nMaxSeq {
			maxSeqSynced = aMaxSeq
		}

		m.seqs[groupID] = SyncedSeq{
			maxSeqSynced: maxSeqSynced,
			sessionType:  constant.SuperGroupChatType,
		}
	}
	return nil
}

func (m *MsgSyncer) DoListener() {
	for {
		select {
		case cmd := <-m.recvSeqch:
			m.compareSeqsAndTrigger(cmd.Ctx, cmd.Value.(map[string]Seq), cmd.Cmd)
		case cmd := <-m.recvMsgCh:
			m.handleRecvMsgAndSyncSeqs(cmd.Ctx, cmd.Value.(*sdkws.MsgData))
		case <-m.ctx.Done():
			log.ZInfo(m.ctx, "msg syncer done, sdk logout.....")
			return
		}
	}
}

// init, reconnect, sync by heartbeat
func (m *MsgSyncer) compareSeqsAndTrigger(ctx context.Context, newSeqMap map[string]Seq, cmd string) {
	// sync callback to conversation
	switch cmd {
	case constant.CmdInit:
		m.triggerSync()
		defer m.triggerSyncFinished()
	case constant.CmdReconnect:
		m.triggerReconnect()
		defer m.triggerReconnectFinished()
	}
	for sourceID, newSeq := range newSeqMap {
		if syncedSeq, ok := m.seqs[sourceID]; ok {
			if newSeq.maxSeq > syncedSeq.maxSeqSynced {
				_ = m.sync(ctx, sourceID, newSeq.sessionType, syncedSeq.maxSeqSynced, newSeq.maxSeq)
			}
		} else {
			// new conversation
			_ = m.sync(ctx, sourceID, newSeq.sessionType, 0, newSeq.maxSeq)
		}
	}
	m.syncTimes++
}

func (m *MsgSyncer) sync(ctx context.Context, sourceID string, sessionType int32, syncedMaxSeq, maxSeq int64) (err error) {
	if err = m.syncAndTriggerMsgs(ctx, sourceID, sessionType, syncedMaxSeq, maxSeq); err != nil {
		log.ZError(ctx, "sync msgs failed", err, "sourceID", sourceID)
		return err
	}
	m.seqs[sourceID] = SyncedSeq{
		maxSeqSynced: maxSeq,
		sessionType:  sessionType,
	}
	return nil
}

// get seqs need sync interval
func (m *MsgSyncer) getSeqsNeedSync(syncedMaxSeq, maxSeq int64) []int64 {
	var seqs []int64
	for i := syncedMaxSeq + 1; i <= maxSeq; i++ {
		seqs = append(seqs, i)
	}
	return seqs
}

// recv msg from
func (m *MsgSyncer) handleRecvMsgAndSyncSeqs(ctx context.Context, msg *sdkws.MsgData) {
	// online msg
	if msg.Seq == 0 {
		_ = m.triggerConversation(ctx, []*sdkws.MsgData{msg})
		return
	}
	// 连续直接触发并且刷新seq
	if msg.Seq == m.seqs[msg.GroupID].maxSeqSynced+1 {
		_ = m.triggerConversation(ctx, []*sdkws.MsgData{msg})
		oldSeq := m.seqs[msg.GroupID]
		oldSeq.maxSeqSynced = msg.Seq
		m.seqs[msg.GroupID] = oldSeq
	} else {
		// m.msgCache[msg.GroupID][msg.Seq] = msg
		m.sync(ctx, msg.GroupID, msg.SessionType, m.seqs[msg.GroupID].maxSeqSynced, msg.Seq)
	}
}

// 分片同步消息，触发成功后刷新seq
func (m *MsgSyncer) syncAndTriggerMsgs(ctx context.Context, sourceID string, sessionType int32, syncedMaxSeq, maxSeq int64) error {
	msgs, err := m.syncMsgBySeqsInterval(ctx, sourceID, sessionType, syncedMaxSeq, maxSeq)
	if err != nil {
		log.ZError(ctx, "syncMsgFromSvr err", err, "sourceID", sourceID, "sessionType", sessionType, "syncedMaxSeq", syncedMaxSeq, "maxSeq", maxSeq)
		return err
	}
	_ = m.triggerConversation(ctx, msgs)
	return err
}

func (m *MsgSyncer) splitSeqs(split int, seqsNeedSync []int64) (splitSeqs [][]int64) {
	if len(seqsNeedSync) <= split {
		splitSeqs = append(splitSeqs, seqsNeedSync)
		return
	}
	for i := 0; i < len(seqsNeedSync); i += split {
		end := i + split
		if end > len(seqsNeedSync) {
			end = len(seqsNeedSync)
		}
		splitSeqs = append(splitSeqs, seqsNeedSync[i:end])
	}
	return
}

// cached的不拉取
func (m *MsgSyncer) syncMsgBySeqsInterval(ctx context.Context, sourceID string, sesstionType int32, syncedMaxSeq, syncedMinSeq int64) (partMsgs []*sdkws.MsgData, err error) {
	return partMsgs, nil
}

func (m *MsgSyncer) syncMsgBySeqs(ctx context.Context, sourceID string, sessionType int32, seqsNeedSync []int64) (allMsgs []*sdkws.MsgData, err error) {
	var pullMsgReq sdkws.PullMessageBySeqsReq
	pullMsgReq.UserID = m.loginUserID
	pullMsgReq.GroupSeqs = make(map[string]*sdkws.Seqs, 0)
	split := constant.SplitPullMsgNum
	seqsList := m.splitSeqs(split, seqsNeedSync)
	for i := 0; i < len(seqsList); {
		pullMsgReq.GroupSeqs[sourceID] = &sdkws.Seqs{Seqs: seqsList[i]}
		resp, err := m.ws.SendReqWaitResp(ctx, &pullMsgReq, constant.WSPullMsgBySeqList, timeout, retryTimes, m.loginUserID)
		if err != nil {
			log.ZError(ctx, "syncMsgFromSvrSplit err", err, "pullMsgReq", pullMsgReq)
			continue
		}
		i++
		var pullMsgResp sdkws.PullMessageBySeqsResp
		err = proto.Unmarshal(resp.Data, &pullMsgResp)
		if err != nil {
			log.ZError(ctx, "Unmarshal failed", err, "resp", resp)
			continue
		}
		allMsgs = append(allMsgs, pullMsgResp.List...)
	}
	return allMsgs, nil
}

func (m *MsgSyncer) triggerConversation(ctx context.Context, msgs []*sdkws.MsgData) error {
	err := common.TriggerCmdNewMsgCome(sdk_struct.CmdNewMsgComeToConversation{Ctx: ctx, MsgList: msgs}, m.conversationCh)
	if err != nil {
		log.ZError(ctx, "triggerCmdNewMsgCome err", err, "msgs", msgs)
	}
	return err
}

func (m *MsgSyncer) triggerReconnect() {

}

func (m *MsgSyncer) triggerReconnectFinished() {

}

func (m *MsgSyncer) triggerSync() {

}

func (m *MsgSyncer) triggerSyncFinished() {

}
