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

var splitPullMsgNum = 1000
var pullMsgNumWhenLogin = 10000

type SelfMsgSync struct {
	*db.DataBase
	*Ws
	loginUserID        string
	conversationCh     chan common.Cmd2Value
	seqMaxSynchronized uint32
	seqMaxNeedSync     uint32 //max seq in push or max seq in redis
	pushMsgCache       map[uint32]*server_api_params.MsgData
}

func NewSelfMsgSync(dataBase *db.DataBase, ws *Ws, loginUserID string, conversationCh chan common.Cmd2Value) *SelfMsgSync {
	p := &SelfMsgSync{DataBase: dataBase, Ws: ws, loginUserID: loginUserID, conversationCh: conversationCh}
	p.pushMsgCache = make(map[uint32]*server_api_params.MsgData, 0)
	return p
}

func (m *SelfMsgSync) GetNormalMsgMaxSeq() (uint32, error) {
	return 0, nil
}

func (m *SelfMsgSync) GetLostMsgSeqList(minSeqInSvr uint32, maxSeqInSvr uint32) ([]uint32, error) {
	return nil, nil
}

func (m *SelfMsgSync) compareSeq(operationID string) {
	//todo 统计中间缺失的seq，并同步
	n, err := m.GetNormalMsgSeq()
	if err != nil {
		log.Error(operationID, "GetNormalMsgSeq failed ", err.Error())
	}
	a, err := m.GetAbnormalMsgSeq()
	if err != nil {
		log.Error(operationID, "GetAbnormalMsgSeq failed ", err.Error())
	}
	if n > a {
		m.seqMaxSynchronized = n
	} else {
		m.seqMaxSynchronized = a
	}
	m.seqMaxNeedSync = m.seqMaxSynchronized
	log.Info(operationID, "load seq, normal, abnormal, ", n, a, m.seqMaxNeedSync, m.seqMaxSynchronized)

}

func (m *SelfMsgSync) doMaxSeq(cmd common.Cmd2Value) {
	var maxSeqOnSvr = cmd.Value.(sdk_struct.CmdMaxSeqToMsgSync).MaxSeqOnSvr
	operationID := cmd.Value.(sdk_struct.CmdMaxSeqToMsgSync).OperationID
	log.Debug(operationID, "recv max seq on svr, doMaxSeq, maxSeqOnSvr, m.seqMaxSynchronized, m.seqMaxNeedSync",
		maxSeqOnSvr, m.seqMaxSynchronized, m.seqMaxNeedSync)
	if maxSeqOnSvr < m.seqMaxNeedSync {
		return
	}
	m.seqMaxNeedSync = maxSeqOnSvr
	m.syncMsg(operationID)
}

func (m *SelfMsgSync) doPushBatchMsg(cmd common.Cmd2Value) {
	msg := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).Msg
	operationID := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).OperationID
	log.Debug(operationID, utils.GetSelfFuncName(), "recv push msg, doPushBatchMsg ", "msgData len: ", len(msg.MsgDataList))
	msgDataWrap := server_api_params.MsgDataList{}
	err := proto.Unmarshal(msg.MsgDataList, &msgDataWrap)
	if err != nil {
		log.Error(operationID, "proto Unmarshal err", err.Error())
		return
	}

	if len(msgDataWrap.MsgDataList) == 1 && msgDataWrap.MsgDataList[0].Seq == 0 {
		log.Debug(operationID, utils.GetSelfFuncName(), "seq ==0 TriggerCmdNewMsgCome", msgDataWrap.MsgDataList[0].String())
		m.TriggerCmdNewMsgCome([]*server_api_params.MsgData{msgDataWrap.MsgDataList[0]}, operationID)
		return
	}

	//to cache
	var maxSeq uint32
	for _, v := range msgDataWrap.MsgDataList {
		if v.Seq > m.seqMaxSynchronized {
			m.pushMsgCache[v.Seq] = v
			log.Debug(operationID, "doPushBatchMsg insert cache v.Seq > m.seqMaxSynchronized", v.Seq, m.seqMaxSynchronized)
		} else {
			log.Debug(operationID, "doPushBatchMsg don't insert cache v.Seq <= m.seqMaxSynchronized", v.Seq, m.seqMaxSynchronized)
		}
		if v.Seq > maxSeq {
			maxSeq = v.Seq
		}
	}

	//update m.seqMaxNeedSync
	log.Debug(operationID, "max Seq in push batch msg, m.seqMaxNeedSync ", maxSeq, m.seqMaxNeedSync)
	if maxSeq > m.seqMaxNeedSync {
		m.seqMaxNeedSync = maxSeq
	}

	seqMaxSynchronizedBegin := m.seqMaxSynchronized
	var triggerMsgList []*server_api_params.MsgData
	for {
		seqMaxSynchronizedBegin++
		cacheMsg, ok := m.pushMsgCache[seqMaxSynchronizedBegin]
		if !ok {
			break
		}
		log.Debug(operationID, "TriggerCmdNewMsgCome, node seq ", cacheMsg.Seq)
		triggerMsgList = append(triggerMsgList, cacheMsg)
		m.seqMaxSynchronized = seqMaxSynchronizedBegin
	}

	log.Debug(operationID, "TriggerCmdNewMsgCome, len:  ", len(triggerMsgList))
	if len(triggerMsgList) != 0 {
		m.TriggerCmdNewMsgCome(triggerMsgList, operationID)
	}
	for _, v := range triggerMsgList {
		delete(m.pushMsgCache, v.Seq)
	}
	m.syncMsg(operationID)
}

func (m *SelfMsgSync) doPushMsg(cmd common.Cmd2Value) {
	msg := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).Msg
	operationID := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).OperationID
	if len(msg.MsgDataList) == 0 {
		log.Debug(operationID, "no batch push")
		m.doPushSingleMsg(cmd)
	} else {
		log.Debug(operationID, "batch push")
		m.doPushBatchMsg(cmd)
	}
}

func (m *SelfMsgSync) doPushSingleMsg(cmd common.Cmd2Value) {
	msg := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).Msg
	operationID := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).OperationID
	log.Debug(operationID, utils.GetSelfFuncName(), "recv normal push msg, doPushMsg ", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, m.seqMaxNeedSync, m.seqMaxSynchronized)
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
	m.syncMsg(operationID)
}

func (m *SelfMsgSync) syncMsg(operationID string) {
	if m.seqMaxNeedSync > m.seqMaxSynchronized {
		log.Info(operationID, "do syncMsg ", m.seqMaxSynchronized+1, m.seqMaxNeedSync)
		m.syncMsgFromServer(m.seqMaxSynchronized+1, m.seqMaxNeedSync)
		m.seqMaxSynchronized = m.seqMaxNeedSync
	} else {
		log.Info(operationID, "syncMsg do nothing, m.seqMaxNeedSync <= m.seqMaxSynchronized ",
			m.seqMaxNeedSync, m.seqMaxSynchronized)
	}
}

func (m *SelfMsgSync) syncMsgFromServer(beginSeq, endSeq uint32) {
	if beginSeq > endSeq {
		log.Error("", "beginSeq > endSeq", beginSeq, endSeq)
		return
	}

	var needSyncSeqList []uint32
	for i := beginSeq; i <= endSeq; i++ {
		needSyncSeqList = append(needSyncSeqList, i)
	}
	var SPLIT = splitPullMsgNum
	for i := 0; i < len(needSyncSeqList)/SPLIT; i++ {
		m.syncMsgFromServerSplit(needSyncSeqList[i*SPLIT : (i+1)*SPLIT])
	}
	m.syncMsgFromServerSplit(needSyncSeqList[SPLIT*(len(needSyncSeqList)/SPLIT):])
}

func (m *SelfMsgSync) syncMsgFromCache2ServerSplit(needSyncSeqList []uint32) {
	var msgList []*server_api_params.MsgData
	var noInCache []uint32
	for _, v := range needSyncSeqList {
		cacheMsg, ok := m.pushMsgCache[v]
		if !ok {
			noInCache = append(noInCache, v)
		} else {
			msgList = append(msgList, cacheMsg)
			delete(m.pushMsgCache, v)
		}
	}
	operationID := utils.OperationIDGenerator()
	if len(noInCache) == 0 {
		m.TriggerCmdNewMsgCome(msgList, operationID)
		return
	}

	var pullMsgReq server_api_params.PullMessageBySeqListReq
	pullMsgReq.SeqList = noInCache
	pullMsgReq.UserID = m.loginUserID
	for {
		operationID = utils.OperationIDGenerator()
		pullMsgReq.OperationID = operationID
		resp, err := m.SendReqWaitResp(&pullMsgReq, constant.WSPullMsgBySeqList, 60, 2, m.loginUserID, operationID)
		if err != nil {
			log.Error(operationID, "SendReqWaitResp failed ", err.Error(), constant.WSPullMsgBySeqList, 60, 2, m.loginUserID)
			continue
		}
		var pullMsgResp server_api_params.PullMessageBySeqListResp
		err = proto.Unmarshal(resp.Data, &pullMsgResp)
		if err != nil {
			log.Error(operationID, "Unmarshal failed ", err.Error())
			return

		}
		msgList = append(msgList, pullMsgResp.List...)
		m.TriggerCmdNewMsgCome(msgList, operationID)
		break
	}
}

func (m *SelfMsgSync) syncMsgFromServerSplit(needSyncSeqList []uint32) {
	m.syncMsgFromCache2ServerSplit(needSyncSeqList)
}

func (m *SelfMsgSync) TriggerCmdNewMsgCome(msgList []*server_api_params.MsgData, operationID string) {
	for {
		err := common.TriggerCmdNewMsgCome(sdk_struct.CmdNewMsgComeToConversation{MsgList: msgList, OperationID: operationID}, m.conversationCh)
		if err != nil {
			log.Warn(operationID, "TriggerCmdNewMsgCome failed ", err.Error(), m.loginUserID)
			continue
		}
		return
	}
}
