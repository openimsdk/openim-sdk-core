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
var pullMsgNumForReadDiffusion = 100

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

//计算最大seq，初始化时调用一次
func (m *SelfMsgSync) compareSeq(operationID string) {
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

//处理心跳最大seq
func (m *SelfMsgSync) doMaxSeq(cmd common.Cmd2Value) {
	var maxSeqOnSvr = cmd.Value.(sdk_struct.CmdMaxSeqToMsgSync).MaxSeqOnSvr
	var minSeqOnSvr = cmd.Value.(sdk_struct.CmdMaxSeqToMsgSync).MinSeqOnSvr
	operationID := cmd.Value.(sdk_struct.CmdMaxSeqToMsgSync).OperationID
	log.Debug(operationID, utils.GetSelfFuncName(), " args ", " maxSeqOnSvr, minSeqOnSvr, m.seqMaxSynchronized, m.seqMaxNeedSync", maxSeqOnSvr, minSeqOnSvr, m.seqMaxSynchronized, m.seqMaxNeedSync)
	if minSeqOnSvr > maxSeqOnSvr {
		log.Error(operationID, "minSeqOnSvr > maxSeqOnSvr", minSeqOnSvr, maxSeqOnSvr)
		return
	}
	if minSeqOnSvr > m.seqMaxSynchronized {
		m.seqMaxSynchronized = minSeqOnSvr - 1
	}
	if maxSeqOnSvr <= m.seqMaxNeedSync {
		log.Warn(operationID, "do nothing ", maxSeqOnSvr, m.seqMaxNeedSync)
		return
	}

	m.seqMaxNeedSync = maxSeqOnSvr
	m.syncMsg(operationID)
}

//暂未启用
func (m *SelfMsgSync) doPushBatchMsg(cmd common.Cmd2Value) {
	msg := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).Msg
	operationID := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).OperationID
	log.Warn(operationID, utils.GetSelfFuncName(), " not enabled ", "msgData len: ")
	log.Debug(operationID, utils.GetSelfFuncName(), " args ", "msgData len: ", len(msg.MsgDataList))
	msgDataWrap := server_api_params.MsgDataList{}
	err := proto.Unmarshal(msg.MsgDataList, &msgDataWrap)
	if err != nil {
		log.Error(operationID, "proto Unmarshal err", err.Error())
		return
	}

	if len(msgDataWrap.MsgDataList) == 1 && msgDataWrap.MsgDataList[0].Seq == 0 {
		log.Debug(operationID, utils.GetSelfFuncName(), "seq ==0 TriggerCmdNewMsgCome", msgDataWrap.MsgDataList[0].String())
		m.triggerCmdNewMsgCome([]*server_api_params.MsgData{msgDataWrap.MsgDataList[0]}, operationID)
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
		m.triggerCmdNewMsgCome(triggerMsgList, operationID)
	}
	for _, v := range triggerMsgList {
		delete(m.pushMsgCache, v.Seq)
	}
	m.syncMsg(operationID)
}

//处理消息推送，兼容单条和批量推送
func (m *SelfMsgSync) doPushMsg(cmd common.Cmd2Value) {
	msg := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).Msg
	operationID := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).OperationID
	if len(msg.MsgDataList) == 0 {
		log.Debug(operationID, "single push")
		m.doPushSingleMsg(cmd)
	} else {
		log.Debug(operationID, "batch push")
		m.doPushBatchMsg(cmd)
	}
}

//处理单条消息推送
func (m *SelfMsgSync) doPushSingleMsg(cmd common.Cmd2Value) {
	msg := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).Msg
	operationID := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).OperationID
	log.Debug(operationID, utils.GetSelfFuncName(), " args: ", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, m.seqMaxNeedSync, m.seqMaxSynchronized)
	if msg.Seq == 0 {
		log.Debug(operationID, "TriggerCmdNewMsgCome ", msg.ServerMsgID, msg.ClientMsgID, msg.Seq)
		m.triggerCmdNewMsgCome([]*server_api_params.MsgData{msg}, operationID)
		return
	}

	//seq连续，则消息直接触发
	if msg.Seq == m.seqMaxSynchronized+1 {
		log.Debug(operationID, "TriggerCmdNewMsgCome ", msg.ServerMsgID, msg.ClientMsgID, msg.Seq)
		m.triggerCmdNewMsgCome([]*server_api_params.MsgData{msg}, operationID)
		m.seqMaxSynchronized = msg.Seq
	}
	if msg.Seq > m.seqMaxNeedSync {
		m.seqMaxNeedSync = msg.Seq
	}
	m.syncMsg(operationID)
}

//从服务端同步消息，通过seqMaxSynchronized  seqMaxNeedSync 来同步中间的seq，同步完设置这两个值
func (m *SelfMsgSync) syncMsg(operationID string) {
	if m.seqMaxNeedSync > m.seqMaxSynchronized {
		log.Info(operationID, "do syncMsgFromServer ", m.seqMaxSynchronized+1, m.seqMaxNeedSync)
		m.syncMsgFromServer(m.seqMaxSynchronized+1, m.seqMaxNeedSync, operationID)
		m.seqMaxSynchronized = m.seqMaxNeedSync
	} else {
		log.Info(operationID, "do nothing, m.seqMaxNeedSync <= m.seqMaxSynchronized ", m.seqMaxNeedSync, m.seqMaxSynchronized)
	}
}

//从本地缓存+服务端获取消息，内部对seq列表做了拆分
func (m *SelfMsgSync) syncMsgFromServer(beginSeq, endSeq uint32, operationID string) {
	log.Debug(operationID, utils.GetSelfFuncName(), " args ", beginSeq, endSeq)
	if beginSeq > endSeq {
		log.Error(operationID, "beginSeq > endSeq", beginSeq, endSeq)
		return
	}
	var needSyncSeqList []uint32
	for i := beginSeq; i <= endSeq; i++ {
		needSyncSeqList = append(needSyncSeqList, i)
	}
	var SPLIT = splitPullMsgNum
	for i := 0; i < len(needSyncSeqList)/SPLIT; i++ {
		m.syncMsgFromServerSplit(needSyncSeqList[i*SPLIT:(i+1)*SPLIT], operationID)
	}
	m.syncMsgFromServerSplit(needSyncSeqList[SPLIT*(len(needSyncSeqList)/SPLIT):], operationID)
}

//先从本地缓存读取消息，如果不存在，再从服务端读取， 上层把seq列表拆分
func (m *SelfMsgSync) syncMsgFromCache2ServerSplit(needSyncSeqList []uint32, operationID string) {
	if len(needSyncSeqList) > splitPullMsgNum {
		log.Error(operationID, "seq list too large ", len(needSyncSeqList))
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
		m.triggerCmdNewMsgCome(msgList, operationID)
		return
	}
	log.Debug(operationID, "seq no in cache num: ", len(noInCache), " all seq num: ", len(needSyncSeqList))
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
		m.triggerCmdNewMsgCome(msgList, operationID)
		break
	}
}

func (m *SelfMsgSync) syncMsgFromServerSplit(needSyncSeqList []uint32, operationID string) {
	m.syncMsgFromCache2ServerSplit(needSyncSeqList, operationID)
}

//触发新消息
func (m *SelfMsgSync) triggerCmdNewMsgCome(msgList []*server_api_params.MsgData, operationID string) {
	for {
		err := common.TriggerCmdNewMsgCome(sdk_struct.CmdNewMsgComeToConversation{MsgList: msgList, OperationID: operationID}, m.conversationCh)
		if err != nil {
			log.Warn(operationID, "TriggerCmdNewMsgCome failed, try again ", err.Error(), m.loginUserID)
			continue
		}
		log.Warn(operationID, "TriggerCmdNewMsgCome ok ", m.loginUserID)
		return
	}
}
