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

type MsgSync struct {
	*db.DataBase
	*Ws
	loginUserID        string
	conversationCh     chan common.Cmd2Value
	PushMsgAndMaxSeqCh chan common.Cmd2Value
	seqMaxSynchronized uint32
	seqMaxNeedSync     uint32
}

func (m *MsgSync) compareSeq() {
	//todo 统计中间缺失的seq，并同步
	n, err := m.GetNormalMsgSeq()
	if err != nil {
		log.Error("", "GetNormalMsgSeq failed ", err.Error())
	}
	a, err := m.GetAbnormalMsgSeq()
	if err != nil {
		log.Error("", "GetAbnormalMsgSeq failed ", err.Error())
	}
	if n > a {
		m.seqMaxSynchronized = n
	} else {
		m.seqMaxSynchronized = a
	}
	m.seqMaxNeedSync = m.seqMaxSynchronized
	log.Info("", "load seq, normal, abnormal, ", n, a, m.seqMaxNeedSync, m.seqMaxSynchronized)
}

func (m *MsgSync) doMaxSeq(cmd common.Cmd2Value) {
	var maxSeqOnSvr = cmd.Value.(sdk_struct.CmdMaxSeqToMsgSync).MaxSeqOnSvr
	operationID := cmd.Value.(sdk_struct.CmdMaxSeqToMsgSync).OperationID
	log.Debug(operationID, "recv max seq on svr, doMaxSeq, maxSeqOnSvr, m.seqMaxSynchronized, m.seqMaxNeedSync",
		maxSeqOnSvr, m.seqMaxSynchronized, m.seqMaxNeedSync)
	if maxSeqOnSvr <= m.seqMaxNeedSync {
		return
	}
	m.seqMaxNeedSync = maxSeqOnSvr
	log.Debug(operationID, "syncMsgFromServer ", m.seqMaxSynchronized+1, m.seqMaxNeedSync)
	m.syncMsg()
}

func (m *MsgSync) doPushMsg(cmd common.Cmd2Value) {
	msg := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).Msg
	operationID := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).OperationID
	log.Debug(operationID, "recv push msg, doPushMsg ", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, m.seqMaxNeedSync, m.seqMaxSynchronized)
	if msg.Seq == 0 {
		m.TriggerCmdNewMsgCome([]*server_api_params.MsgData{msg}, operationID)
		return
	}
	if m.seqMaxNeedSync == 0 {
		return
	}

	if msg.Seq == m.seqMaxSynchronized+1 {
		log.Debug(operationID, "TriggerCmdNewMsgCome ", msg.ServerMsgID, msg.ClientMsgID, msg.Seq)
		m.TriggerCmdNewMsgCome([]*server_api_params.MsgData{msg}, operationID)
		m.seqMaxSynchronized = msg.Seq
	}
	if msg.Seq > m.seqMaxNeedSync {
		m.seqMaxNeedSync = msg.Seq
	}
	log.Debug(operationID, "syncMsgFromServer ", m.seqMaxSynchronized+1, m.seqMaxNeedSync)
	m.syncMsg()
}

func (m *MsgSync) Work(cmd common.Cmd2Value) {
	switch cmd.Cmd {
	case constant.CmdPushMsg:
		m.doPushMsg(cmd)
	case constant.CmdMaxSeq:
		m.doMaxSeq(cmd)
	default:
		log.Error("", "cmd failed ", cmd.Cmd)
	}
}

func (m *MsgSync) GetCh() chan common.Cmd2Value {
	return m.PushMsgAndMaxSeqCh
}

func NewMsgSync(dataBase *db.DataBase, ws *Ws, loginUserID string, ch chan common.Cmd2Value, pushMsgAndMaxSeqCh chan common.Cmd2Value) *MsgSync {
	p := &MsgSync{DataBase: dataBase,
		Ws: ws, loginUserID: loginUserID, conversationCh: ch, PushMsgAndMaxSeqCh: pushMsgAndMaxSeqCh}
	p.compareSeq()
	go common.DoListener(p)
	return p
}

func (m *MsgSync) syncMsg() {
	if m.seqMaxNeedSync > m.seqMaxSynchronized {
		log.Info("", "do syncMsg ", m.seqMaxSynchronized+1, m.seqMaxNeedSync)
		m.syncMsgFromServer(m.seqMaxSynchronized+1, m.seqMaxNeedSync)
		m.seqMaxSynchronized = m.seqMaxNeedSync
	}
}

func (m *MsgSync) syncMsgFromServer(beginSeq, endSeq uint32) {
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
		m.syncMsgFromServerSplit(needSyncSeqList[i*SPLIT : (i+1)*SPLIT])
	}
	m.syncMsgFromServerSplit(needSyncSeqList[SPLIT*(len(needSyncSeqList)/SPLIT):])
}

func (m *MsgSync) syncMsgFromServerSplit(needSyncSeqList []uint32) {
	var pullMsgReq server_api_params.PullMessageBySeqListReq
	pullMsgReq.SeqList = needSyncSeqList
	pullMsgReq.UserID = m.loginUserID
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

func (m *MsgSync) TriggerCmdNewMsgCome(msgList []*server_api_params.MsgData, operationID string) {

	for {
		err := common.TriggerCmdNewMsgCome(sdk_struct.CmdNewMsgComeToConversation{MsgList: msgList, OperationID: operationID}, m.conversationCh)
		if err != nil {
			log.Error(operationID, "TriggerCmdNewMsgCome failed ", err.Error(), m.loginUserID)
			continue
		}
		log.Warn(operationID, "TriggerCmdNewMsgCome ok ", m.loginUserID)
		return
	}
}
