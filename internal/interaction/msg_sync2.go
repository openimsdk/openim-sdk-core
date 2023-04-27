package interaction

import (
	"context"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/sdk_struct"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/golang/protobuf/proto"
)

const (
	timeout    = 60
	retryTimes = 2
)

type Seq struct {
	maxSeq      int64
	minSeq      int64
	sessionType int32
}

type MsgSyncer struct {
	loginUserID string
	// listen ch
	ws        *Ws
	recvSeqch chan common.Cmd2Value // recv max seq from heartbeat map[string][2]int64
	recvMsgCh chan common.Cmd2Value // recv msg from ws, sync or push []*sdkws.MsgData

	conversationCh chan common.Cmd2Value // trigger conversation

	ctx      context.Context
	seqs     map[string]Seq
	msgCache map[string]map[int64]*sdkws.MsgData
}

func NewMsgSyncer(ctx context.Context, conversationCh, recvMsgCh, recvSeqch chan common.Cmd2Value,
	loginUserID string, ws *Ws) *MsgSyncer {
	return &MsgSyncer{
		recvSeqch: recvSeqch,
		recvMsgCh: recvMsgCh,

		conversationCh: conversationCh,
		ws:             ws,
		loginUserID:    loginUserID,
		ctx:            ctx,
		seqs:           make(map[string]Seq),
		msgCache:       make(map[string]map[int64]*sdkws.MsgData),
	}
}

func (m *MsgSyncer) DoListener() {
	for {
		select {
		case cmd := <-m.recvSeqch:
			m.compareSeqsAndTrigger(cmd.Ctx, cmd.Value.(map[string]Seq))
		case cmd := <-m.recvMsgCh:
			m.handleRecvMsgAndSyncSeqs(cmd.Ctx, cmd.Value.(*sdkws.MsgData))
		case <-m.ctx.Done():
			log.ZInfo(m.ctx, "msg syncer done, sdk logout.....")
			return
		}
	}
}

func (m *MsgSyncer) compareSeqsAndTrigger(ctx context.Context, newSeqMap map[string]Seq) {
	for sourceID, newSeq := range newSeqMap {
		if oldSeq, ok := m.seqs[sourceID]; ok {
			if newSeq.maxSeq > oldSeq.maxSeq {
				if err := m.syncAndTriggerMsgs(ctx, sourceID, newSeq.sessionType, m.getSeqsNeedSync(oldSeq.maxSeq, newSeq.maxSeq)); err == nil {
					m.seqs[sourceID] = newSeq
				}
			}
		} else {
			// 新的会话
			if err := m.syncAndTriggerMsgs(ctx, sourceID, newSeq.sessionType, m.getSeqsNeedSync(0, newSeq.maxSeq)); err == nil {
				m.seqs[sourceID] = newSeq
			}
		}
	}
}

func (m *MsgSyncer) getSeqsNeedSync(oldMaxSeq, newMaxSeq int64) []int64 {
	var seqs []int64
	for i := oldMaxSeq + 1; i <= newMaxSeq; i++ {
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
	if msg.Seq == m.seqs[msg.GroupID].maxSeq+1 {
		_ = m.triggerConversation(ctx, []*sdkws.MsgData{msg})
		oldSeq := m.seqs[msg.GroupID]
		oldSeq.maxSeq = msg.Seq
		m.seqs[msg.GroupID] = oldSeq
	} else {
		m.msgCache[msg.GroupID][msg.Seq] = msg
		m.syncAndTriggerMsgs(ctx, msg.GroupID, msg.SessionType, m.getSeqsNeedSync(m.seqs[msg.GroupID].maxSeq, msg.Seq))
	}
}

// 分片同步消息，触发成功后刷新seq
func (m *MsgSyncer) syncAndTriggerMsgs(ctx context.Context, sourceID string, sessionType int32, seqsNeedSync []int64) error {
	var seqsNeedSyncExcludeCached []int64
	for _, seq := range seqsNeedSync {
		if _, ok := m.msgCache[sourceID][seq]; !ok {
			seqsNeedSyncExcludeCached = append(seqsNeedSyncExcludeCached, seq)
		}
	}
	msgs, err := m.syncMsgFromSvr(ctx, sourceID, sessionType, seqsNeedSyncExcludeCached)
	if err != nil {
		log.ZError(ctx, "syncMsgFromSvr err", err, "sourceID", sourceID, "sessionType", sessionType, "seqsNeedSync", seqsNeedSyncExcludeCached)
		return err
	}
	delete(m.msgCache, sourceID)
	_ = m.triggerConversation(ctx, msgs)
	return nil
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

func (m *MsgSyncer) syncMsgFromSvr(ctx context.Context, sourceID string, sessionType int32, seqsNeedSync []int64) (allMsgs []*sdkws.MsgData, err error) {
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
