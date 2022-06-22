package interaction

import (
	"github.com/golang/protobuf/proto"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/sdk_struct"
	"sync"
)

type SuperGroupMsgSync struct {
	*db.DataBase
	*Ws
	loginUserID              string
	conversationCh           chan common.Cmd2Value
	superGroupMtx            sync.Mutex
	Group2SeqMaxNeedSync     map[string]uint32
	Group2SeqMaxSynchronized map[string]uint32
	SuperGroupIDList         []string
	joinedSuperGroupCh       chan common.Cmd2Value
}

func NewSuperGroupMsgSync(dataBase *db.DataBase, ws *Ws, loginUserID string, conversationCh chan common.Cmd2Value, joinedSuperGroupCh chan common.Cmd2Value) *SuperGroupMsgSync {
	p := &SuperGroupMsgSync{DataBase: dataBase, Ws: ws, loginUserID: loginUserID, conversationCh: conversationCh, joinedSuperGroupCh: joinedSuperGroupCh}
	p.Group2SeqMaxNeedSync = make(map[string]uint32, 0)
	p.Group2SeqMaxSynchronized = make(map[string]uint32, 0)
	go p.updateJoinedSuperGroup()
	return p
}

func (m *SuperGroupMsgSync) updateJoinedSuperGroup() {
	for {
		select {
		case cmd := <-m.joinedSuperGroupCh:

			operationID := cmd.Value.(sdk_struct.CmdJoinedSuperGroup).OperationID
			log.Info(operationID, "updateJoinedSuperGroup recv cmd: ", cmd)
			g, err := m.GetJoinedSuperGroupList()
			if err == nil {
				m.superGroupMtx.Lock()
				m.SuperGroupIDList = m.SuperGroupIDList[0:0]
				for _, v := range g {
					m.SuperGroupIDList = append(m.SuperGroupIDList, v.GroupID)
				}
				m.superGroupMtx.Unlock()
				m.compareSeq(operationID)
			}
		}
	}
}

func (m *SuperGroupMsgSync) compareSeq(operationID string) {
	g, err := m.GetJoinedSuperGroupList()
	if err == nil {
		m.superGroupMtx.Lock()
		m.SuperGroupIDList = m.SuperGroupIDList[0:0]
		for _, v := range g {
			m.SuperGroupIDList = append(m.SuperGroupIDList, v.GroupID)
		}
		m.superGroupMtx.Unlock()
	}

	log.Debug(operationID, "compareSeq load groupID list ", m.SuperGroupIDList)

	m.superGroupMtx.Lock()

	defer m.superGroupMtx.Unlock()
	for _, v := range m.SuperGroupIDList {
		var seqMaxSynchronized uint32
		var seqMaxNeedSync uint32
		n, err := m.GetSuperGroupNormalMsgSeq(v)
		if err != nil {
			log.Error(operationID, "GetSuperGroupNormalMsgSeq failed ", err.Error(), v)
		}
		a, err := m.GetSuperGroupAbnormalMsgSeq(v)
		if err != nil {
			log.Error(operationID, "GetSuperGroupAbnormalMsgSeq failed ", err.Error(), v)
		}
		log.Warn(operationID, "GetSuperGroupNormalMsgSeq GetSuperGroupAbnormalMsgSeq", n, a)
		if n > a {
			seqMaxSynchronized = n
		} else {
			seqMaxSynchronized = a
		}
		seqMaxNeedSync = seqMaxSynchronized
		m.Group2SeqMaxNeedSync[v] = seqMaxNeedSync
		m.Group2SeqMaxSynchronized[v] = seqMaxSynchronized
		log.Info(operationID, "load seq, normal, abnormal, ", n, a, seqMaxNeedSync, seqMaxSynchronized)
	}
}

func (m *SuperGroupMsgSync) doMaxSeq(cmd common.Cmd2Value) {
	operationID := cmd.Value.(sdk_struct.CmdMaxSeqToMsgSync).OperationID
	m.superGroupMtx.Lock()
	for groupID, maxSeqOnSvr := range cmd.Value.(sdk_struct.CmdMaxSeqToMsgSync).GroupID2MaxSeqOnSvr {
		seqMaxNeedSync := m.Group2SeqMaxNeedSync[groupID]
		log.Debug(operationID, "super group doMaxSeq, maxSeqOnSvr, seqMaxSynchronized, seqMaxNeedSync",
			maxSeqOnSvr, m.Group2SeqMaxSynchronized[groupID], seqMaxNeedSync)
		if maxSeqOnSvr <= seqMaxNeedSync {
			continue
		}
		m.Group2SeqMaxNeedSync[groupID] = maxSeqOnSvr
	}
	m.superGroupMtx.Unlock()
	m.syncMsg(operationID)
}

func (m *SuperGroupMsgSync) doPushMsg(cmd common.Cmd2Value) {
	msg := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).Msg
	operationID := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).OperationID
	log.Debug(operationID, "recv super group push msg, doPushMsg ", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, msg.GroupID, msg.SessionType)
	if msg.Seq == 0 {
		m.TriggerCmdNewMsgCome([]*server_api_params.MsgData{msg}, operationID)
		return
	}

	//seqMaxNeedSync := m.Group2SeqMaxNeedSync[msg.GroupID]
	//	seqMaxSynchronized := m.Group2SeqMaxSynchronized[msg.GroupID]

	if m.Group2SeqMaxNeedSync[msg.GroupID] == 0 {
		return
	}
	if m.Group2SeqMaxNeedSync[msg.GroupID] == 0 {
		return
	}

	if msg.Seq == m.Group2SeqMaxSynchronized[msg.GroupID]+1 {
		log.Debug(operationID, "TriggerCmdNewMsgCome ", msg.ServerMsgID, msg.ClientMsgID, msg.Seq)
		m.TriggerCmdNewMsgCome([]*server_api_params.MsgData{msg}, operationID)
		m.Group2SeqMaxSynchronized[msg.GroupID] = msg.Seq
	}
	if msg.Seq > m.Group2SeqMaxNeedSync[msg.GroupID] {
		m.Group2SeqMaxNeedSync[msg.GroupID] = msg.Seq
	}
	log.Debug(operationID, "syncMsgFromServer ", m.Group2SeqMaxSynchronized[msg.GroupID]+1, m.Group2SeqMaxNeedSync[msg.GroupID])
	m.syncMsg(operationID)
}

func (m *SuperGroupMsgSync) syncMsg(operationID string) {
	m.superGroupMtx.Lock()
	for groupID, seqMaxNeedSync := range m.Group2SeqMaxNeedSync {
		seqMaxSynchronized := m.Group2SeqMaxSynchronized[groupID]
		if seqMaxNeedSync > seqMaxSynchronized {
			log.Info(operationID, "do syncMsg ", seqMaxSynchronized+1, seqMaxNeedSync)
			m.syncMsgFromServer(seqMaxSynchronized+1, seqMaxNeedSync, groupID, operationID)
			m.Group2SeqMaxSynchronized[groupID] = seqMaxNeedSync
		}
	}
	m.superGroupMtx.Unlock()
}

func (m *SuperGroupMsgSync) syncMsgFromServer(beginSeq, endSeq uint32, groupID, operationID string) {
	if beginSeq > endSeq {
		log.Error("", "beginSeq > endSeq", beginSeq, endSeq)
		return
	}

	var needSyncSeqList []uint32
	for i := beginSeq; i <= endSeq; i++ {
		needSyncSeqList = append(needSyncSeqList, i)
	}
	var SPLIT = 100
	for i := 0; i < len(needSyncSeqList)/SPLIT; i++ {
		m.syncMsgFromServerSplit(needSyncSeqList[i*SPLIT:(i+1)*SPLIT], groupID, operationID)
	}
	m.syncMsgFromServerSplit(needSyncSeqList[SPLIT*(len(needSyncSeqList)/SPLIT):], groupID, operationID)
}

func (m *SuperGroupMsgSync) syncMsgFromServerSplit(needSyncSeqList []uint32, groupID, operationID string) {
	var pullMsgReq server_api_params.PullMessageBySeqListReq
	//pullMsgReq.SeqList = needSyncSeqList
	pullMsgReq.UserID = m.loginUserID
	pullMsgReq.GroupSeqList = make(map[string]*server_api_params.SeqList, 0)
	pullMsgReq.GroupSeqList[groupID] = &server_api_params.SeqList{SeqList: needSyncSeqList}

	for {
		pullMsgReq.OperationID = operationID
		log.Debug(operationID, "super group pull message", groupID)
		resp, err := m.SendReqWaitResp(&pullMsgReq, constant.WSPullMsgBySeqList, 30, 2, m.loginUserID, operationID)
		if err != nil {
			log.Error(operationID, "SendReqWaitResp failed ", err.Error(), constant.WSPullMsgBySeqList, 30, 2, m.loginUserID)
			continue
		}

		var pullMsgResp server_api_params.PullMessageBySeqListResp
		err = proto.Unmarshal(resp.Data, &pullMsgResp)
		if err != nil {
			log.Error(operationID, "Unmarshal failed ", err.Error())
			return
		}
		log.Debug(operationID, "SendReqWaitResp pull msg ", pullMsgReq.String())
		log.Debug(operationID, "SendReqWaitResp pull msg result ", pullMsgResp.String())
		for _, v := range pullMsgResp.GroupMsgDataList {
			m.TriggerCmdNewMsgCome(v.MsgDataList, operationID)
		}
		return
	}
}

func (m *SuperGroupMsgSync) TriggerCmdNewMsgCome(msgList []*server_api_params.MsgData, operationID string) {
	for {
		err := common.TriggerCmdNewMsgCome(sdk_struct.CmdNewMsgComeToConversation{MsgList: msgList, OperationID: operationID}, m.conversationCh)
		if err != nil {
			log.Warn(operationID, "TriggerCmdNewMsgCome failed ", err.Error(), m.loginUserID)
			continue
		}
		return
	}
}
