package interaction

import (
	"github.com/golang/protobuf/proto"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
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
}

func (m *MsgSync) doMaxSeq(cmd common.Cmd2Value) {
	var cmdSeq = cmd.Value.(uint32)
	if cmdSeq <= m.seqMaxNeedSync {
		return
	}
	m.seqMaxNeedSync = cmdSeq
	m.syncMsgFromServer(m.seqMaxSynchronized+1, m.seqMaxNeedSync)
	m.seqMaxSynchronized = m.seqMaxNeedSync
}

func (m *MsgSync) doPushMsg(cmd common.Cmd2Value) {
	if m.seqMaxNeedSync == 0 {
		return
	}
	msg := cmd.Value.(*server_api_params.MsgData)
	if uint32(msg.Seq)+1 == m.seqMaxNeedSync && m.seqMaxNeedSync == m.seqMaxSynchronized {
		m.TriggerCmdNewMsgCome([]*server_api_params.MsgData{msg})
		m.seqMaxNeedSync = uint32(msg.Seq) + 1
		m.seqMaxSynchronized = uint32(msg.Seq) + 1
		return
	}
	if uint32(msg.Seq) > m.seqMaxNeedSync {
		m.seqMaxNeedSync = uint32(msg.Seq)
		m.syncMsgFromServer(m.seqMaxSynchronized+1, m.seqMaxNeedSync)
		m.seqMaxSynchronized = m.seqMaxNeedSync
		return
	}
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
	return m.cmdCh
}

func NewMsgSync(dataBase *db.DataBase, ws *Ws, loginUserID string, ch chan common.Cmd2Value, pushMsgAndMaxSeqCh chan common.Cmd2Value) *MsgSync {
	p := &MsgSync{DataBase: dataBase,
		Ws: ws, loginUserID: loginUserID, conversationCh: ch, PushMsgAndMaxSeqCh: pushMsgAndMaxSeqCh}
	p.compareSeq()
	go common.DoListener(p)
	return p
}

func (m *MsgSync) syncMsgFromServer(beginSeq, endSeq uint32) {
	var needSyncSeqList []uint32
	for i := beginSeq; i <= endSeq; i++ {
		needSyncSeqList = append(needSyncSeqList, i)
	}
	var SPLIT = 100
	for i := 0; i < len(needSyncSeqList)/SPLIT; i++ {
		//0-99 100-199
		m.syncMsgFromServerSplit(needSyncSeqList[i*SPLIT : (i+1)*SPLIT])
	}
	m.syncMsgFromServerSplit(needSyncSeqList[SPLIT*(len(needSyncSeqList)/SPLIT):])
}

func (m *MsgSync) syncMsgFromServerSplit(needSyncSeqList []uint32) {
	operationID := utils.OperationIDGenerator()
	var pullMsgReq server_api_params.PullMessageBySeqListReq
	pullMsgReq.SeqList = needSyncSeqList
	pullMsgReq.UserID = m.loginUserID
	pullMsgReq.OperationID = operationID
	for {
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
		m.TriggerCmdNewMsgCome(pullMsgResp.List)
		return
	}
}

func (m *MsgSync) TriggerCmdNewMsgCome(msgList []*server_api_params.MsgData) {
	for {
		err := common.TriggerCmdNewMsgCome(msgList, m.conversationCh)
		if err != nil {
			continue
		}
		return
	}
}
