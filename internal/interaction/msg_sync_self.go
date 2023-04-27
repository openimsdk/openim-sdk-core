// Copyright © 2023 OpenIM SDK.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package interaction

import (
	"context"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"runtime"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/golang/protobuf/proto"
)

type SelfMsgSync struct {
	db_interface.DataBase
	*Ws
	loginUserID        string
	conversationCh     chan common.Cmd2Value
	seqMaxSynchronized int64
	seqMaxNeedSync     int64 //max seq in push or max seq in redis
	pushMsgCache       map[int64]*sdkws.MsgData
}

func NewSelfMsgSync(dataBase db_interface.DataBase, ws *Ws, loginUserID string, conversationCh chan common.Cmd2Value) *SelfMsgSync {
	p := &SelfMsgSync{DataBase: dataBase, Ws: ws, loginUserID: loginUserID, conversationCh: conversationCh}
	p.pushMsgCache = make(map[int64]*sdkws.MsgData, 0)
	return p
}

// 计算最大seq，初始化时调用一次
func (m *SelfMsgSync) compareSeq(ctx context.Context) {
	n, err := m.GetNormalMsgSeq(ctx)
	if err != nil {
		// log.Error(operationID, "GetNormalMsgSeq failed ", err.Error())
	}
	a, err := m.GetAbnormalMsgSeq(ctx)
	if err != nil {
		// log.Error(operationID, "GetAbnormalMsgSeq failed ", err.Error())
	}
	if n > a {
		m.seqMaxSynchronized = n
	} else {
		m.seqMaxSynchronized = a
	}
	m.seqMaxNeedSync = m.seqMaxSynchronized
	// log.Info(operationID, "load seq, normal, abnormal, ", n, a, m.seqMaxNeedSync, m.seqMaxSynchronized)
}

// 处理心跳最大seq
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
		log.Debug(operationID, "do nothing ", maxSeqOnSvr, m.seqMaxNeedSync)
		return
	}

	m.seqMaxNeedSync = maxSeqOnSvr
	m.syncMsg(operationID)
}

// 暂未启用
func (m *SelfMsgSync) doPushBatchMsg(cmd common.Cmd2Value) {
	msg := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).Msg
	operationID := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).OperationID
	log.Warn(operationID, utils.GetSelfFuncName(), " not enabled ", "msgData len: ")
	log.Debug(operationID, utils.GetSelfFuncName(), " args ", "msgData len: ", len(msg.MsgDataList))
	msgDataWrap := sdkws.MsgDataList{}
	err := proto.Unmarshal(msg.MsgDataList, &msgDataWrap)
	if err != nil {
		log.Error(operationID, "proto Unmarshal err", err.Error())
		return
	}

	if len(msgDataWrap.MsgDataList) == 1 && msgDataWrap.MsgDataList[0].Seq == 0 {
		log.Debug(operationID, utils.GetSelfFuncName(), "seq ==0 TriggerCmdNewMsgCome", msgDataWrap.MsgDataList[0].String())
		m.triggerCmdNewMsgCome([]*sdkws.MsgData{msgDataWrap.MsgDataList[0]}, operationID)
		return
	}

	//to cache
	var maxSeq int64
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
	var triggerMsgList []*sdkws.MsgData
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

// 处理消息推送，兼容单条和批量推送
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

// 处理单条消息推送
func (m *SelfMsgSync) doPushSingleMsg(cmd common.Cmd2Value) {
	msg := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).Msg
	operationID := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).OperationID
	log.Debug(operationID, utils.GetSelfFuncName(), " args: ", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, m.seqMaxNeedSync, m.seqMaxSynchronized)
	if msg.Seq == 0 {
		log.Debug(operationID, "TriggerCmdNewMsgCome ", msg.ServerMsgID, msg.ClientMsgID, msg.Seq)
		m.triggerCmdNewMsgCome([]*sdkws.MsgData{msg}, operationID)
		return
	}

	//seq连续，则消息直接触发
	if msg.Seq == m.seqMaxSynchronized+1 {
		log.Debug(operationID, "TriggerCmdNewMsgCome ", msg.ServerMsgID, msg.ClientMsgID, msg.Seq)
		m.triggerCmdNewMsgCome([]*sdkws.MsgData{msg}, operationID)
		m.seqMaxSynchronized = msg.Seq
	}
	if msg.Seq > m.seqMaxNeedSync {
		m.seqMaxNeedSync = msg.Seq
	}
	m.syncMsg(operationID)
}

// 从服务端同步消息，通过seqMaxSynchronized  seqMaxNeedSync 来同步中间的seq，同步完设置这两个值
func (m *SelfMsgSync) syncMsg(operationID string) {
	if m.seqMaxNeedSync > m.seqMaxSynchronized {
		log.Info(operationID, "do syncMsgFromServer ", m.seqMaxSynchronized+1, m.seqMaxNeedSync)
		m.syncMsgFromServer(m.seqMaxSynchronized+1, m.seqMaxNeedSync, operationID)
		m.seqMaxSynchronized = m.seqMaxNeedSync
	} else {
		log.Info(operationID, "do nothing, m.seqMaxNeedSync <= m.seqMaxSynchronized ", m.seqMaxNeedSync, m.seqMaxSynchronized)
	}
}

// 从本地缓存+服务端获取消息，内部对seq列表做了拆分
func (m *SelfMsgSync) syncMsgFromServer(beginSeq, endSeq int64, operationID string) {
	log.Debug(operationID, utils.GetSelfFuncName(), " args ", beginSeq, endSeq)
	if beginSeq > endSeq {
		log.Error(operationID, "beginSeq > endSeq", beginSeq, endSeq)
		return
	}
	var needSyncSeqList []int64
	for i := beginSeq; i <= endSeq; i++ {
		needSyncSeqList = append(needSyncSeqList, i)
	}
	var SPLIT = constant.SplitPullMsgNum
	for i := 0; i < len(needSyncSeqList)/SPLIT; i++ {
		m.syncMsgFromServerSplit(needSyncSeqList[i*SPLIT:(i+1)*SPLIT], operationID)
	}
	m.syncMsgFromServerSplit(needSyncSeqList[SPLIT*(len(needSyncSeqList)/SPLIT):], operationID)
}

// 先从本地缓存读取消息，如果不存在，再从服务端读取， 上层把seq列表拆分
func (m *SelfMsgSync) syncMsgFromCache2ServerSplit(needSyncSeqList []int64, operationID string) {
	if len(needSyncSeqList) > constant.SplitPullMsgNum {
		log.Error(operationID, "seq list too large ", len(needSyncSeqList))
		return
	}

	var msgList []*sdkws.MsgData
	var noInCache []int64
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
	var pullMsgReq sdkws.PullMessageBySeqsReq
	pullMsgReq.Seqs = noInCache
	pullMsgReq.UserID = m.loginUserID
	for {
		resp, err := m.SendReqWaitResp(context.Background(), &pullMsgReq, constant.WSPullMsgBySeqList, 60, 2, m.loginUserID)
		if err != nil && m.LoginStatus() == constant.Logout {
			log.Error(operationID, "SendReqWaitResp failed  Logout status ", err.Error(), m.LoginStatus())
			log.Warn("", "m.LoginStatus() == constant.Logout, Goexit()")
			runtime.Goexit()
			return
		}
		if err != nil {
			log.Error(operationID, "SendReqWaitResp failed ", err.Error(), constant.WSPullMsgBySeqList, 60, 2, m.loginUserID)
			continue
		}
		var pullMsgResp sdkws.PullMessageBySeqsResp
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

func (m *SelfMsgSync) syncMsgFromServerSplit(needSyncSeqList []int64, operationID string) {
	m.syncMsgFromCache2ServerSplit(needSyncSeqList, operationID)
}

// 触发新消息
func (m *SelfMsgSync) triggerCmdNewMsgCome(msgList []*sdkws.MsgData, operationID string) {
	for {
		err := common.TriggerCmdNewMsgCome(sdk_struct.CmdNewMsgComeToConversation{MsgList: msgList, OperationID: operationID}, m.conversationCh)
		if err != nil {
			log.Warn(operationID, "TriggerCmdNewMsgCome failed, try again ", err.Error(), m.loginUserID)
			continue
		}
		log.Debug(operationID, "TriggerCmdNewMsgCome ok ", m.loginUserID)
		return
	}
}
