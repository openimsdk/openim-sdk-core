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

type SelfMsgSyncLatestModel struct {
	*db.DataBase
	*Ws
	loginUserID        string
	conversationCh     chan common.Cmd2Value
	seqMaxSynchronized uint32
	seqMaxNeedSync     uint32 //max seq in push or max seq in redis
	pushMsgCache       map[uint32]*server_api_params.MsgData

	lostMsgSeqList  []uint32
	syncMsgFinished bool
	maxSeqOnLocal   uint32
}

func NewSelfMsgSyncLatestModel(dataBase *db.DataBase, ws *Ws, loginUserID string, conversationCh chan common.Cmd2Value) *SelfMsgSyncLatestModel {
	p := &SelfMsgSyncLatestModel{DataBase: dataBase, Ws: ws, loginUserID: loginUserID, conversationCh: conversationCh}
	p.pushMsgCache = make(map[uint32]*server_api_params.MsgData, 0)
	return p
}

//1
func (m *SelfMsgSyncLatestModel) GetLocalNormalMsgMaxSeq() (uint32, error) {
	return m.GetNormalMsgSeq()
}

//1
func (m *SelfMsgSyncLatestModel) GetLocalLostMsgSeqList(minSeqInSvr uint32) ([]uint32, error) {
	return m.GetLostMsgSeqList(minSeqInSvr)
}

//1
func (m *SelfMsgSyncLatestModel) compareSeq(operationID string, minSeqInSvr uint32) {
	n, err := m.GetLocalNormalMsgMaxSeq()
	if err != nil {
		log.Error(operationID, "GetNormalMsgSeq failed ", err.Error())
	} else {
		m.seqMaxSynchronized = n
		m.maxSeqOnLocal = n
	}
	lostSeqList, err := m.GetLocalLostMsgSeqList(minSeqInSvr)
	if err != nil {
		log.Error(operationID, "GetLostMsgSeqList failed ", err.Error(), minSeqInSvr)
	} else {
		m.lostMsgSeqList = lostSeqList
	}
	log.Info(operationID, "load,  seqMaxSynchronized, lostMsgSeqList maxSeqOnLocal ", m.seqMaxSynchronized, m.lostMsgSeqList, m.maxSeqOnLocal)
}

//1
func (m *SelfMsgSyncLatestModel) doMaxSeq(cmd common.Cmd2Value) {
	var maxSeqOnSvr = cmd.Value.(sdk_struct.CmdMaxSeqToMsgSync).MaxSeqOnSvr
	var minSeqOnSvr = cmd.Value.(sdk_struct.CmdMaxSeqToMsgSync).MinSeqOnSvr
	operationID := cmd.Value.(sdk_struct.CmdMaxSeqToMsgSync).OperationID
	log.Debug(operationID, "doMaxSeq, minSeqOnSvr, maxSeqOnSvr, m.seqMaxSynchronized, m.seqMaxNeedSync, m.maxSeqOnLocal",
		minSeqOnSvr, maxSeqOnSvr, m.seqMaxSynchronized, m.seqMaxNeedSync, m.maxSeqOnLocal)
	if !m.syncMsgFinished {
		log.Info(operationID, " syncMsgWhenLogin start ", minSeqOnSvr, maxSeqOnSvr)
		m.syncMsgWhenLogin(minSeqOnSvr, maxSeqOnSvr, operationID)
	}
	if maxSeqOnSvr <= m.seqMaxNeedSync {
		return
	}
	m.seqMaxNeedSync = maxSeqOnSvr
	m.syncMsg(operationID, constant.MsgSyncModelDefault)
}

//1
func (m *SelfMsgSyncLatestModel) syncMsgWhenLogin(minSeqOnSvr, maxSeqOnSvr uint32, operationID string) {
	log.Debug(operationID, utils.GetSelfFuncName(), "args: ", minSeqOnSvr, maxSeqOnSvr, m.syncMsgFinished)
	if m.syncMsgFinished {
		return
	}
	m.compareSeq(operationID, minSeqOnSvr)
	if maxSeqOnSvr > m.seqMaxNeedSync {
		m.seqMaxNeedSync = maxSeqOnSvr
	}
	m.syncMsg(operationID, constant.MsgSyncModelLogin)
	go m.syncLostMsg(operationID)
	m.syncMsgFinished = true
}

//
//func (m *SelfMsgSyncLatestModel) doPushBatchMsg(cmd common.Cmd2Value) {
//	msg := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).Msg
//	operationID := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).OperationID
//	log.Debug(operationID, utils.GetSelfFuncName(), "recv push msg, doPushBatchMsg ", "msgData len: ", len(msg.MsgDataList))
//	msgDataWrap := server_api_params.MsgDataList{}
//	err := proto.Unmarshal(msg.MsgDataList, &msgDataWrap)
//	if err != nil {
//		log.Error(operationID, "proto Unmarshal err", err.Error())
//		return
//	}
//
//	if len(msgDataWrap.MsgDataList) == 1 && msgDataWrap.MsgDataList[0].Seq == 0 {
//		log.Debug(operationID, utils.GetSelfFuncName(), "seq ==0 TriggerCmdNewMsgCome", msgDataWrap.MsgDataList[0].String())
//		m.TriggerCmdNewMsgCome([]*server_api_params.MsgData{msgDataWrap.MsgDataList[0]}, operationID)
//		return
//	}
//
//	//to cache
//	var maxSeq uint32
//	for _, v := range msgDataWrap.MsgDataList {
//		if v.Seq > m.seqMaxSynchronized {
//			m.pushMsgCache[v.Seq] = v
//			log.Debug(operationID, "doPushBatchMsg insert cache v.Seq > m.seqMaxSynchronized", v.Seq, m.seqMaxSynchronized)
//		} else {
//			log.Debug(operationID, "doPushBatchMsg don't insert cache v.Seq <= m.seqMaxSynchronized", v.Seq, m.seqMaxSynchronized)
//		}
//		if v.Seq > maxSeq {
//			maxSeq = v.Seq
//		}
//	}
//
//	//update m.seqMaxNeedSync
//	log.Debug(operationID, "max Seq in push batch msg, m.seqMaxNeedSync ", maxSeq, m.seqMaxNeedSync)
//	if maxSeq > m.seqMaxNeedSync {
//		m.seqMaxNeedSync = maxSeq
//	}
//
//	seqMaxSynchronizedBegin := m.seqMaxSynchronized
//	var triggerMsgList []*server_api_params.MsgData
//	for {
//		seqMaxSynchronizedBegin++
//		cacheMsg, ok := m.pushMsgCache[seqMaxSynchronizedBegin]
//		if !ok {
//			break
//		}
//		log.Debug(operationID, "TriggerCmdNewMsgCome, node seq ", cacheMsg.Seq)
//		triggerMsgList = append(triggerMsgList, cacheMsg)
//		m.seqMaxSynchronized = seqMaxSynchronizedBegin
//	}
//
//	log.Debug(operationID, "TriggerCmdNewMsgCome, len:  ", len(triggerMsgList))
//	if len(triggerMsgList) != 0 {
//		m.TriggerCmdNewMsgCome(triggerMsgList, operationID)
//	}
//	for _, v := range triggerMsgList {
//		delete(m.pushMsgCache, v.Seq)
//	}
//	m.syncMsg(operationID)
//}

func (m *SelfMsgSyncLatestModel) doPushMsg(cmd common.Cmd2Value) {
	msg := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).Msg
	operationID := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).OperationID
	if len(msg.MsgDataList) == 0 {
		log.Debug(operationID, "no batch push")
		m.doPushSingleMsg(cmd)
	} else {
		log.NewWarn(operationID, "batch push ")
		//	m.doPushBatchMsg(cmd)
	}
}

func (m *SelfMsgSyncLatestModel) doPushSingleMsg(cmd common.Cmd2Value) {
	msg := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).Msg
	operationID := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).OperationID
	log.Debug(operationID, utils.GetSelfFuncName(), "doPushSingleMsg ", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, m.seqMaxNeedSync, m.seqMaxSynchronized)
	if msg.Seq == 0 {
		m.TriggerCmdNewMsgCome([]*server_api_params.MsgData{msg}, operationID, constant.MsgSyncModelDefault, 0)
		return
	}
	if m.seqMaxNeedSync == 0 {
		return
	}

	if msg.Seq == m.seqMaxSynchronized+1 {
		log.Debug(operationID, "TriggerCmdNewMsgCome ", msg.ServerMsgID, msg.ClientMsgID, msg.Seq)
		m.TriggerCmdNewMsgCome([]*server_api_params.MsgData{msg}, operationID, constant.MsgSyncModelDefault, 0)
		m.seqMaxSynchronized = msg.Seq
	}
	if msg.Seq > m.seqMaxNeedSync {
		m.seqMaxNeedSync = msg.Seq
	}
	log.Debug(operationID, "syncMsgFromServer ", m.seqMaxSynchronized+1, m.seqMaxNeedSync)
	m.syncMsg(operationID, constant.MsgSyncModelDefault)
}

//1
func (m *SelfMsgSyncLatestModel) syncMsg(operationID string, syncFlag int) {
	if m.seqMaxNeedSync > m.seqMaxSynchronized {
		log.Info(operationID, "do syncMsg ", m.seqMaxSynchronized+1, m.seqMaxNeedSync)
		if syncFlag == constant.MsgSyncModelDefault {
			m.syncMsgFromServer(m.seqMaxSynchronized+1, m.seqMaxNeedSync, syncFlag, operationID)
			m.seqMaxSynchronized = m.seqMaxNeedSync
			return
		}

		if m.seqMaxNeedSync-m.seqMaxSynchronized < uint32(pullMsgNumWhenLogin) {
			m.syncMsgFromServer(m.seqMaxSynchronized+1, m.seqMaxNeedSync, syncFlag, operationID)
			m.seqMaxSynchronized = m.seqMaxNeedSync
		} else { //50000-20000
			m.syncMsgFromServer(m.seqMaxNeedSync-uint32(pullMsgNumWhenLogin)+1, m.seqMaxNeedSync, constant.MsgSyncModelDefault, operationID) //40000+1,50000
			m.seqMaxSynchronized = m.seqMaxNeedSync
			go m.syncMsgFromServer(m.seqMaxSynchronized+1, m.seqMaxNeedSync-uint32(pullMsgNumWhenLogin), constant.MsgSyncModelDefault, operationID) //20000+1 , 40000
		}
	} else {
		log.Info(operationID, "syncMsg do nothing, m.seqMaxNeedSync <= m.seqMaxSynchronized ",
			m.seqMaxNeedSync, m.seqMaxSynchronized)
	}
}

//1
func (m *SelfMsgSyncLatestModel) syncLostMsg(operationID string) {
	if len(m.lostMsgSeqList) == 0 {
		return
	}
	needSyncSeqList := m.lostMsgSeqList
	var SPLIT = splitPullMsgNum
	for i := 0; i < len(needSyncSeqList)/SPLIT; i++ {
		m.syncMsgFromServerSplit(needSyncSeqList[i*SPLIT:(i+1)*SPLIT], constant.MsgSyncModelDefault, operationID)
	}
	m.syncMsgFromServerSplit(needSyncSeqList[SPLIT*(len(needSyncSeqList)/SPLIT):], constant.MsgSyncModelDefault, operationID)
	needSyncSeqList = nil
}

//1
func (m *SelfMsgSyncLatestModel) syncMsgFromServer(beginSeq, endSeq uint32, syncFlag int, operationID string) {
	log.Info(operationID, utils.GetSelfFuncName(), "args: ", beginSeq, endSeq, syncFlag)
	if beginSeq > endSeq {
		log.Error(operationID, "syncMsgFromServer beginSeq > endSeq", beginSeq, endSeq)
		return
	}
	var needSyncSeqList []uint32
	for i := endSeq; i >= beginSeq; i-- {
		needSyncSeqList = append(needSyncSeqList, i)
	}
	var SPLIT = splitPullMsgNum
	for i := 0; i < len(needSyncSeqList)/SPLIT; i++ {
		log.Debug(operationID, "syncMsgFromServerSplit", syncFlag, len(needSyncSeqList[i*SPLIT:(i+1)*SPLIT]))
		m.syncMsgFromServerSplit(needSyncSeqList[i*SPLIT:(i+1)*SPLIT], syncFlag, operationID)
	}
	log.Debug(operationID, "syncMsgFromServerSplit", syncFlag, len(needSyncSeqList[SPLIT*(len(needSyncSeqList)/SPLIT):]))
	m.syncMsgFromServerSplit(needSyncSeqList[SPLIT*(len(needSyncSeqList)/SPLIT):], syncFlag, operationID)
}

//1
func (m *SelfMsgSyncLatestModel) syncMsgFromCache2ServerSplit(needSyncSeqList []uint32, syncFlag int, operationID string) {
	log.Info(operationID, utils.GetSelfFuncName(), "args: ", len(needSyncSeqList), syncFlag)
	if len(needSyncSeqList) == 0 {
		return
	}
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
	if len(noInCache) == 0 {
		m.TriggerCmdNewMsgCome(msgList, operationID, syncFlag, needSyncSeqList[0])
		return
	}

	var pullMsgReq server_api_params.PullMessageBySeqListReq
	pullMsgReq.SeqList = noInCache
	pullMsgReq.UserID = m.loginUserID
	for {
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
		m.TriggerCmdNewMsgCome(msgList, operationID, syncFlag, needSyncSeqList[0])
		return
	}
}

//1
func (m *SelfMsgSyncLatestModel) syncMsgFromServerSplit(needSyncSeqList []uint32, syncFlag int, operationID string) {
	m.syncMsgFromCache2ServerSplit(needSyncSeqList, syncFlag, operationID)
}

//1
func (m *SelfMsgSyncLatestModel) TriggerCmdNewMsgCome(msgList []*server_api_params.MsgData, operationID string, syncFlag int, currentMaxSeq uint32) {
	for {
		if syncFlag == constant.MsgSyncModelLogin {
			err := common.TriggerCmdNewMsgCome(sdk_struct.CmdNewMsgComeToConversation{MsgList: msgList, OperationID: operationID,
				SyncFlag: syncFlag, MaxSeqOnSvr: m.seqMaxNeedSync, MaxSeqOnLocal: m.maxSeqOnLocal, CurrentMaxSeq: currentMaxSeq, PullMsgOrder: constant.SyncOrderStartLatest}, m.conversationCh)
			if err != nil {
				log.Warn(operationID, "TriggerCmdNewMsgCome failed ", err.Error(), m.loginUserID)
				continue
			}
			return
		} else {
			err := common.TriggerCmdNewMsgCome(sdk_struct.CmdNewMsgComeToConversation{MsgList: msgList, OperationID: operationID, SyncFlag: syncFlag}, m.conversationCh)
			if err != nil {
				log.Warn(operationID, "TriggerCmdNewMsgCome failed ", err.Error(), m.loginUserID)
				continue
			}
			return
		}
	}
}
