package interaction

import (
	"github.com/golang/protobuf/proto"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

type SuperGroupMsgSync struct {
	*db.DataBase
	*Ws
	loginUserID    string
	conversationCh chan common.Cmd2Value
	//  PushMsgAndMaxSeqCh chan common.Cmd2Value
	//seqMaxSynchronized uint32
	//seqMaxNeedSync     uint32
	Group2SeqMaxNeedSync     map[string]uint32
	Group2SeqMaxSynchronized map[string]uint32
	GroupIDList              []string
}

func (m *SuperGroupMsgSync) compareSeq() {
	for _, v := range m.GroupIDList {
		var seqMaxSynchronized uint32
		var seqMaxNeedSync uint32
		n, err := m.GetSuperGroupNormalMsgSeq(v)
		if err != nil {
			log.Error("", "GetSuperGroupNormalMsgSeq failed ", err.Error(), v)
		}
		a, err := m.GetSuperGroupAbnormalMsgSeq(v)
		if err != nil {
			log.Error("", "GetSuperGroupAbnormalMsgSeq failed ", err.Error(), v)
		}
		if n > a {
			seqMaxSynchronized = n
		} else {
			seqMaxSynchronized = a
		}
		seqMaxNeedSync = seqMaxSynchronized
		m.Group2SeqMaxNeedSync[v] = seqMaxNeedSync
		m.Group2SeqMaxSynchronized[v] = seqMaxSynchronized
		log.Info("", "load seq, normal, abnormal, ", n, a, seqMaxNeedSync, seqMaxSynchronized)
	}
}

func (m *SuperGroupMsgSync) doMaxSeq(cmd common.Cmd2Value) {
	operationID := cmd.Value.(sdk_struct.CmdMaxSeqToMsgSync).OperationID
	for groupID, maxSeqOnSvr := range cmd.Value.(sdk_struct.CmdMaxSeqToMsgSync).GroupID2MaxSeqOnSvr {
		seqMaxNeedSync := m.Group2SeqMaxNeedSync[groupID]
		log.Debug(operationID, "super group doMaxSeq, maxSeqOnSvr, seqMaxSynchronized, seqMaxNeedSync",
			maxSeqOnSvr, m.Group2SeqMaxSynchronized[groupID], seqMaxNeedSync)
		if maxSeqOnSvr <= seqMaxNeedSync {
			continue
		}
		m.Group2SeqMaxNeedSync[groupID] = maxSeqOnSvr
	}
	m.syncMsg()
}

func (m *SuperGroupMsgSync) doPushMsg(cmd common.Cmd2Value) {
	msg := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).Msg
	operationID := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).OperationID
	log.Debug(operationID, "recv push msg, doPushMsg ", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, msg.GroupID, msg.SessionType)
	if msg.Seq == 0 {
		m.TriggerCmdNewMsgCome([]*server_api_params.MsgData{msg}, operationID)
		return
	}

	seqMaxNeedSync := m.Group2SeqMaxNeedSync[msg.GroupID]
	seqMaxSynchronized := m.Group2SeqMaxSynchronized[msg.GroupID]

	if m.Group2SeqMaxNeedSync[msg.GroupID] == 0 {
		return
	}
	if seqMaxNeedSync == 0 {
		return
	}

	if msg.Seq == seqMaxSynchronized+1 {
		log.Debug(operationID, "TriggerCmdNewMsgCome ", msg.ServerMsgID, msg.ClientMsgID, msg.Seq)
		m.TriggerCmdNewMsgCome([]*server_api_params.MsgData{msg}, operationID)
		seqMaxSynchronized = msg.Seq
	}
	if msg.Seq > seqMaxNeedSync {
		seqMaxNeedSync = msg.Seq
	}
	log.Debug(operationID, "syncMsgFromServer ", seqMaxSynchronized+1, seqMaxNeedSync)
	m.syncMsg()
}

func (m *SuperGroupMsgSync) syncMsg() {
	for groupID, seqMaxNeedSync := range m.Group2SeqMaxNeedSync {
		seqMaxSynchronized := m.Group2SeqMaxSynchronized[groupID]
		if seqMaxNeedSync > seqMaxSynchronized {
			log.Info("", "do syncMsg ", seqMaxSynchronized+1, seqMaxNeedSync)
			m.syncMsgFromServer(seqMaxSynchronized+1, seqMaxNeedSync, groupID)
			m.Group2SeqMaxSynchronized[groupID] = seqMaxNeedSync
		}
	}

}

func (m *SuperGroupMsgSync) syncMsgFromServer(beginSeq, endSeq uint32, groupID string) {
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
		m.syncMsgFromServerSplit(needSyncSeqList[i*SPLIT:(i+1)*SPLIT], groupID)
	}
	m.syncMsgFromServerSplit(needSyncSeqList[SPLIT*(len(needSyncSeqList)/SPLIT):], groupID)
}

func (m *SuperGroupMsgSync) syncMsgFromServerSplit(needSyncSeqList []uint32, groupID string) {
	var pullMsgReq server_api_params.PullMessageBySeqListReq
	pullMsgReq.SeqList = needSyncSeqList
	pullMsgReq.UserID = m.loginUserID
	pullMsgReq.GroupSeqList[groupID] = &server_api_params.SeqList{SeqList: needSyncSeqList}

	for {
		operationID := utils.OperationIDGenerator()
		pullMsgReq.OperationID = operationID
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
		m.TriggerCmdNewMsgCome(pullMsgResp.List, operationID)
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
