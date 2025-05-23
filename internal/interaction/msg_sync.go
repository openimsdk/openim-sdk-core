// Copyright © 2023 OpenIM SDK. All rights reserved.
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
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sort_conversation"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/stringutil"

	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
)

const (
	connectPullNums       = 1
	defaultPullNums       = 10
	SplitPullMsgNum       = 100
	pullMsgGoroutineLimit = 10
	maxConversations      = 500
	synMaxConversations   = 100
)

// MsgSyncer is a central hub for message relay, responsible for sequential message gap pulling,
// handling network events, and managing app foreground and background events.
type MsgSyncer struct {
	loginUserID string                // login user ID
	longConnMgr *LongConnMgr          // long connection manager
	recvCh      chan common.Cmd2Value // channel for receiving push messages and the maximum SEQ number
	//conversationEventQueue chan common.Cmd2Value // storage and session triggering
	conversationEventQueue *common.EventQueue
	syncedMaxSeqs          map[string]int64      // map of the maximum synced SEQ numbers for all group IDs
	syncedMaxSeqsLock      sync.RWMutex          // syncedMaxSeqs map lock
	db                     db_interface.DataBase // data store
	reinstalled            bool                  //true if the app was uninstalled and reinstalled
	isSyncing              bool                  // indicates whether data is being synced
	isSyncingLock          sync.Mutex            // lock for syncing state

	isLargeDataSync bool
}

func (m *MsgSyncer) SetLoginUserID(loginUserID string) {
	m.loginUserID = loginUserID
}

func (m *MsgSyncer) SetDataBase(db db_interface.DataBase) {
	m.db = db
}

// NewMsgSyncer creates a new instance of the message synchronizer.
func NewMsgSyncer(recvCh chan common.Cmd2Value, conversationEventQueue *common.EventQueue,
	longConnMgr *LongConnMgr) *MsgSyncer {
	return &MsgSyncer{
		longConnMgr:            longConnMgr,
		recvCh:                 recvCh,
		conversationEventQueue: conversationEventQueue,
		syncedMaxSeqs:          make(map[string]int64),
	}
}

// LoadSeq seq The db reads the data to the memory,set syncedMaxSeqs
func (m *MsgSyncer) LoadSeq(ctx context.Context) error {
	conversationIDList, err := m.db.GetAllConversationIDList(ctx)
	if err != nil {
		log.ZError(ctx, "get conversation id list failed", err)
		return err
	}

	if len(conversationIDList) == 0 {
		version, err := m.db.GetAppSDKVersion(ctx)
		if err != nil && !errors.Is(err, errs.ErrRecordNotFound) {
			return err
		}
		if version == nil || !version.Installed {
			m.reinstalled = true
		}
	}

	// TODO With a large number of sessions(10w), this could potentially cause blocking and needs optimization.

	type SyncedSeq struct {
		ConversationID string
		MaxSyncedSeq   int64
		Err            error
	}

	partSize := 20
	currency := (len(conversationIDList)-1)/partSize + 1
	if len(conversationIDList) == 0 {
		currency = 0
	}
	var wg sync.WaitGroup
	resultMaps := make([]map[string]SyncedSeq, currency)

	for i := 0; i < currency; i++ {
		wg.Add(1)
		start := i * partSize
		end := start + partSize
		if i == currency-1 {
			end = len(conversationIDList)
		}

		resultMaps[i] = make(map[string]SyncedSeq)

		go func(i, start, end int) {
			defer wg.Done()
			for _, v := range conversationIDList[start:end] {
				maxSyncedSeq, err := m.db.CheckConversationNormalMsgSeq(ctx, v)
				resultMaps[i][v] = SyncedSeq{
					ConversationID: v,
					MaxSyncedSeq:   maxSyncedSeq,
					Err:            err,
				}
			}
		}(i, start, end)
	}

	wg.Wait()

	// merge map
	for _, resultMap := range resultMaps {
		for k, v := range resultMap {
			if v.Err != nil {
				log.ZError(ctx, "get group normal seq failed", errs.Wrap(v.Err), "conversationID", k)
				continue
			}
			m.syncedMaxSeqs[k] = v.MaxSyncedSeq
		}
	}
	notificationSeqs, err := m.db.GetNotificationAllSeqs(ctx)
	if err != nil {
		log.ZError(ctx, "get notification seq failed", err)
		return err
	}
	for _, notificationSeq := range notificationSeqs {
		m.syncedMaxSeqs[notificationSeq.ConversationID] = notificationSeq.Seq
	}
	log.ZDebug(ctx, "loadSeq", "syncedMaxSeqs", m.syncedMaxSeqs)
	return nil
}

// DoListener Listen to the message pipe of the message synchronizer
// and process received and pushed messages
func (m *MsgSyncer) DoListener(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("panic: %+v\n%s", r, debug.Stack())

			log.ZWarn(ctx, "DoListener panic", nil, "panic info", err)
		}
	}()
	for {
		select {
		case cmd := <-m.recvCh:
			m.handlePushMsgAndEvent(cmd)
		case <-ctx.Done():
			log.ZInfo(ctx, "msg syncer done, sdk logout.....")
			return
		}
	}
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
func (m *MsgSyncer) handlePushMsgAndEvent(cmd common.Cmd2Value) {
	switch cmd.Cmd {
	case constant.CmdConnSuccesss:
		log.ZInfo(cmd.Ctx, "recv long conn mgr connected", "cmd", cmd.Cmd, "value", cmd.Value)
		if m.startSync() {
			m.doConnected(cmd.Ctx)
		} else {
			log.ZWarn(cmd.Ctx, "syncing, ignore connected event", nil, "cmd", cmd.Cmd, "value", cmd.Value)
		}
	case constant.CmdWakeUpDataSync:
		log.ZInfo(cmd.Ctx, "app wake up, start sync msgs", "cmd", cmd.Cmd, "value", cmd.Value)
		if m.startSync() {
			m.doWakeupDataSync(cmd.Ctx)
		} else {
			log.ZWarn(cmd.Ctx, "syncing, ignore wake up event", nil, "cmd", cmd.Cmd, "value", cmd.Value)
		}
	case constant.CmdIMMessageSync:
		if conversationIDs, ok := cmd.Value.([]string); ok {
			log.ZInfo(cmd.Ctx, "manual trigger IM message synchronization", "cmd", cmd.Cmd, "value", cmd.Value)
			m.doIMMessageSync(cmd.Ctx, conversationIDs)
		} else {
			log.ZWarn(cmd.Ctx, "invalid value type for IMMessageSync", nil, "cmd", cmd.Cmd, "value", cmd.Value)
		}

	case constant.CmdPushMsg:
		m.doPushMsg(cmd.Ctx, cmd.Value.(*sdkws.PushMessages))
	}
}

func (m *MsgSyncer) getNeedSyncConversations(ctx context.Context, maxSeqToSync map[string]int64) map[string][2]int64 {
	needSyncSeqMap := make(map[string][2]int64)
	//when app reinstalled do not pull notifications messages.
	if m.reinstalled {
		notificationsSeqMap := make(map[string]int64)
		messagesSeqMap := make(map[string]int64)
		for conversationID, seq := range maxSeqToSync {
			if IsNotification(conversationID) {
				if seq != 0 { // seq is 0, no need to sync
					notificationsSeqMap[conversationID] = seq
				}
			} else {
				messagesSeqMap[conversationID] = seq
			}
		}

		var notificationSeqs []*model_struct.NotificationSeqs

		for conversationID, seq := range notificationsSeqMap {
			notificationSeqs = append(notificationSeqs, &model_struct.NotificationSeqs{
				ConversationID: conversationID,
				Seq:            seq,
			})
			m.syncedMaxSeqs[conversationID] = seq
		}

		if len(notificationSeqs) > 0 {
			err := m.db.BatchInsertNotificationSeq(ctx, notificationSeqs)
			if err != nil {
				log.ZWarn(ctx, "BatchInsertNotificationSeq err", err)
			}
		}

		for conversationID, maxSeq := range messagesSeqMap {
			if syncedMaxSeq, ok := m.syncedMaxSeqs[conversationID]; ok {
				if maxSeq > syncedMaxSeq {
					needSyncSeqMap[conversationID] = [2]int64{syncedMaxSeq + 1, maxSeq}
				}
			} else {
				needSyncSeqMap[conversationID] = [2]int64{0, maxSeq}
			}
		}
		return needSyncSeqMap
	} else {
		for conversationID, maxSeq := range maxSeqToSync {
			if syncedMaxSeq, ok := m.syncedMaxSeqs[conversationID]; ok {
				if maxSeq > syncedMaxSeq {
					needSyncSeqMap[conversationID] = [2]int64{syncedMaxSeq + 1, maxSeq}
				}
			} else {
				if maxSeq != 0 { // seq is 0, no need to sync
					needSyncSeqMap[conversationID] = [2]int64{0, maxSeq}
				}
			}
		}
		return needSyncSeqMap
	}
}

func (m *MsgSyncer) compareSeqsAndBatchSync(ctx context.Context, maxSeqToSync map[string]int64, pullNums int64) {
	needSyncSeqMap := make(map[string][2]int64)
	//when app reinstalled do not pull notifications messages.
	if m.reinstalled {
		notificationsSeqMap := make(map[string]int64)
		messagesSeqMap := make(map[string]int64)
		for conversationID, seq := range maxSeqToSync {
			if IsNotification(conversationID) {
				if seq != 0 { // seq is 0, no need to sync
					notificationsSeqMap[conversationID] = seq
				}
			} else {
				messagesSeqMap[conversationID] = seq
			}
		}

		var notificationSeqs []*model_struct.NotificationSeqs

		for conversationID, seq := range notificationsSeqMap {
			notificationSeqs = append(notificationSeqs, &model_struct.NotificationSeqs{
				ConversationID: conversationID,
				Seq:            seq,
			})
			m.syncedMaxSeqs[conversationID] = seq
		}

		if len(notificationSeqs) > 0 {
			err := m.db.BatchInsertNotificationSeq(ctx, notificationSeqs)
			if err != nil {
				log.ZWarn(ctx, "BatchInsertNotificationSeq err", err)
			}
		}

		for conversationID, maxSeq := range messagesSeqMap {
			if syncedMaxSeq, ok := m.syncedMaxSeqs[conversationID]; ok {
				if maxSeq > syncedMaxSeq {
					needSyncSeqMap[conversationID] = [2]int64{syncedMaxSeq + 1, maxSeq}
				}
			} else {
				needSyncSeqMap[conversationID] = [2]int64{0, maxSeq}
			}
		}
		defer func() {
			if err := m.db.SetAppSDKVersion(ctx, &model_struct.LocalAppSDKVersion{
				Installed: true,
			}); err != nil {
				log.ZError(ctx, "SetAppSDKVersion err", err)
			}
			m.reinstalled = false
		}()
		_ = m.syncAndTriggerReinstallMsgs(ctx, needSyncSeqMap, true, pullNums)
	} else {
		for conversationID, maxSeq := range maxSeqToSync {
			if syncedMaxSeq, ok := m.syncedMaxSeqs[conversationID]; ok {
				if maxSeq > syncedMaxSeq {
					needSyncSeqMap[conversationID] = [2]int64{syncedMaxSeq + 1, maxSeq}
				}
			} else {
				if maxSeq != 0 { // seq is 0, no need to sync
					needSyncSeqMap[conversationID] = [2]int64{0, maxSeq}
				}
			}
		}
		_ = m.syncAndTriggerMsgs(ctx, needSyncSeqMap, pullNums)
	}
}

// startSync checks if the sync is already in progress.
// If syncing is in progress, it returns false. Otherwise, it starts syncing and returns true.
func (m *MsgSyncer) startSync() bool {
	m.isSyncingLock.Lock()
	defer m.isSyncingLock.Unlock()

	if m.isSyncing {
		// If already syncing, return false
		return false
	}

	// Set syncing to true and start the sync
	m.isSyncing = true

	// Create a goroutine that waits for 5 seconds and then sets isSyncing to false
	go func() {
		time.Sleep(5 * time.Second)
		m.isSyncingLock.Lock()
		m.isSyncing = false
		m.isSyncingLock.Unlock()
	}()

	return true
}

func (m *MsgSyncer) doPushMsg(ctx context.Context, push *sdkws.PushMessages) {
	log.ZDebug(ctx, "push msgs", "push", push, "syncedMaxSeqs", m.syncedMaxSeqs)
	m.pushTriggerAndSync(ctx, push.Msgs, m.triggerConversation)
	m.pushTriggerAndSync(ctx, push.NotificationMsgs, m.triggerNotification)
}

func (m *MsgSyncer) pushTriggerAndSync(ctx context.Context, pushMessages map[string]*sdkws.PullMsgs, triggerFunc func(ctx context.Context, msgs map[string]*sdkws.PullMsgs) error) {
	if len(pushMessages) == 0 {
		return
	}
	needSyncSeqMap := make(map[string][2]int64)
	var lastSeq int64
	var storageMsgs []*sdkws.MsgData
	for conversationID, msgs := range pushMessages {
		for _, msg := range msgs.Msgs {
			if msg.Seq == 0 {
				_ = triggerFunc(ctx, map[string]*sdkws.PullMsgs{conversationID: {Msgs: []*sdkws.MsgData{msg}}})
				continue
			}
			lastSeq = msg.Seq
			storageMsgs = append(storageMsgs, msg)
		}
		if lastSeq == m.syncedMaxSeqs[conversationID]+int64(len(storageMsgs)) && lastSeq != 0 {
			log.ZDebug(ctx, "trigger msgs", "msgs", storageMsgs)
			_ = triggerFunc(ctx, map[string]*sdkws.PullMsgs{conversationID: {Msgs: storageMsgs}})
			m.syncedMaxSeqs[conversationID] = lastSeq
		} else if lastSeq != 0 && lastSeq > m.syncedMaxSeqs[conversationID] {
			//must pull message when message type is notification
			needSyncSeqMap[conversationID] = [2]int64{m.syncedMaxSeqs[conversationID] + 1, lastSeq}
		}
	}
	m.syncAndTriggerMsgs(ctx, needSyncSeqMap, defaultPullNums)
}

// Called after successful reconnection to synchronize the latest message
func (m *MsgSyncer) doConnected(ctx context.Context) {
	reinstalled := m.reinstalled
	if reinstalled {
		common.DispatchSyncFlagWithMeta(ctx, constant.AppDataSyncBegin, nil, m.conversationEventQueue)
	} else {
		common.DispatchSyncFlagWithMeta(ctx, constant.MsgSyncBegin, nil, m.conversationEventQueue)
	}
	var resp sdkws.GetMaxSeqResp
	if err := m.longConnMgr.SendReqWaitResp(ctx, &sdkws.GetMaxSeqReq{UserID: m.loginUserID}, constant.GetNewestSeq, &resp); err != nil {
		log.ZError(ctx, "get max seq error", err)
		common.DispatchSyncFlag(ctx, constant.MsgSyncFailed, m.conversationEventQueue)
		return
	} else {
		log.ZDebug(ctx, "get max seq success", "resp", resp.MaxSeqs)
	}

	// Calculate the list of conversations that need to be synchronized,
	// including the start and end sequence (seq)
	needSyncAllSeqMap := m.getNeedSyncConversations(ctx, resp.MaxSeqs)
	convCount := len(needSyncAllSeqMap)

	if convCount == 0 {
		log.ZInfo(ctx, "no conversations messages need to sync")
	}

	// In cases where there is no uninstall and reinstall,
	// the amount of conversation data to be synchronized in a single operation is too large
	if len(needSyncAllSeqMap) >= maxConversations {
		log.ZDebug(ctx, "large conversations to sync", "length", len(needSyncAllSeqMap))
		m.isLargeDataSync = true
		common.DispatchSyncFlagWithMeta(ctx, constant.LargeDataSyncBegin, nil, m.conversationEventQueue)
	}

	maxSeqs, sortConversationList, err := m.SyncAndSortConversations(ctx, reinstalled)
	if err != nil {
		log.ZError(ctx, "SyncAndSortConversations err", err)
	}
	if reinstalled {
		common.DispatchSyncFlagWithMeta(ctx, constant.AppDataSyncData, maxSeqs, m.conversationEventQueue)
	} else {
		if m.isLargeDataSync {
			log.ZWarn(ctx, "too many conversations to sync", nil, "maxConversations", maxConversations)
			common.DispatchSyncFlagWithMeta(ctx, constant.LargeDataSyncData, maxSeqs, m.conversationEventQueue)
		} else {
			common.DispatchSyncFlagWithMeta(ctx, constant.MsgSyncData, maxSeqs, m.conversationEventQueue)
		}
	}

	sort_conversation.NewConversationBatchProcessor(sortConversationList, needSyncAllSeqMap, synMaxConversations).Run(ctx, m.handleMessage)

	if reinstalled {
		defer m.markInstallDone(ctx)

	} else {
		if !m.isLargeDataSync {
			common.DispatchSyncFlag(ctx, constant.MsgSyncEnd, m.conversationEventQueue)
		}

	}
}

func (m *MsgSyncer) markInstallDone(ctx context.Context) {
	if err := m.db.SetAppSDKVersion(ctx, &model_struct.LocalAppSDKVersion{Installed: true}); err != nil {
		log.ZError(ctx, "SetAppSDKVersion failed", err)
	}
	m.reinstalled = false
}
func (m *MsgSyncer) handleMessage(ctx context.Context, batchID int, needSyncTopSeqMap map[string][2]int64, isFirst bool) {
	ctx = mcontext.WithTriggerIDContext(ctx, stringutil.IntToString(batchID))
	reinstalled := m.reinstalled
	log.ZDebug(ctx, "handle message need sync top message map", "length", len(needSyncTopSeqMap), "needSyncTopSeqMap", needSyncTopSeqMap)

	if reinstalled {
		_ = m.syncAndTriggerReinstallMsgs(ctx, needSyncTopSeqMap, isFirst, connectPullNums)
		if isFirst {
			common.DispatchSyncFlag(ctx, constant.AppDataSyncEnd, m.conversationEventQueue)
		}
	} else {
		if m.isLargeDataSync {
			log.ZDebug(ctx, "handleMessage large conversations to sync", "length", len(needSyncTopSeqMap), "isFirst", isFirst, "maxConversations", maxConversations)
			_ = m.syncAndTriggerReinstallMsgs(ctx, needSyncTopSeqMap, isFirst, connectPullNums)
			if isFirst {
				common.DispatchSyncFlag(ctx, constant.LargeDataSyncEnd, m.conversationEventQueue)
			}
		} else {
			_ = m.syncAndTriggerMsgs(ctx, needSyncTopSeqMap, connectPullNums)
		}
	}
}

func (m *MsgSyncer) doWakeupDataSync(ctx context.Context) {
	common.DispatchSyncData(ctx, m.conversationEventQueue)
	var resp sdkws.GetMaxSeqResp
	if err := m.longConnMgr.SendReqWaitResp(ctx, &sdkws.GetMaxSeqReq{UserID: m.loginUserID}, constant.GetNewestSeq, &resp); err != nil {
		log.ZError(ctx, "get max seq error", err)
		return
	} else {
		log.ZDebug(ctx, "get max seq success", "resp", resp.MaxSeqs)
	}
	m.compareSeqsAndBatchSync(ctx, resp.MaxSeqs, defaultPullNums)
}

func (m *MsgSyncer) doIMMessageSync(ctx context.Context, conversationIDs []string) {

	resp := msg.GetConversationsHasReadAndMaxSeqResp{}
	req := msg.GetConversationsHasReadAndMaxSeqReq{UserID: m.loginUserID, ConversationIDs: conversationIDs}
	err := m.longConnMgr.SendReqWaitResp(ctx, &req, constant.GetConvMaxReadSeq, &resp)
	if err != nil {
		log.ZWarn(ctx, "GetConvMaxReadSeq SendReqWaitResp err", err)
		return
	} else {
		log.ZDebug(ctx, "GetConvMaxReadSeq SendReqWaitResp success", "resp", resp.Seqs)
	}
	maxSeqMap := make(map[string]int64)
	for conversationID, seqs := range resp.Seqs {
		maxSeqMap[conversationID] = seqs.MaxSeq
	}
	m.compareSeqsAndBatchSync(ctx, maxSeqMap, defaultPullNums)
}

func IsNotification(conversationID string) bool {
	return strings.HasPrefix(conversationID, "n_")
}

func (m *MsgSyncer) syncAndTriggerMsgs(ctx context.Context, seqMap map[string][2]int64, syncMsgNum int64) error {
	if len(seqMap) == 0 {
		log.ZDebug(ctx, "nothing to sync", "syncMsgNum", syncMsgNum)
		return nil
	}

	log.ZDebug(ctx, "current sync seqMap", "seqMap", seqMap)
	var (
		tempSeqMap = make(map[string][2]int64, 50)
		msgNum     = 0
	)

	for k, v := range seqMap {
		oneConversationSyncNum := v[1] - v[0] + 1
		tempSeqMap[k] = v
		// For notification conversations, use oneConversationSyncNum directly
		if IsNotification(k) {
			msgNum += int(oneConversationSyncNum)
		} else {
			// For regular conversations, ensure msgNum is the minimum of oneConversationSyncNum and syncMsgNum
			msgNum += int(min(oneConversationSyncNum, syncMsgNum))
		}

		// If accumulated msgNum reaches SplitPullMsgNum, trigger a batch pull
		if msgNum >= SplitPullMsgNum {
			resp, err := m.pullMsgBySeqRange(ctx, tempSeqMap, syncMsgNum)
			if err != nil {
				log.ZError(ctx, "syncMsgFromServer error", err, "tempSeqMap", tempSeqMap)
				return err
			}
			_ = m.triggerConversation(ctx, resp.Msgs)
			_ = m.triggerNotification(ctx, resp.NotificationMsgs)
			for conversationID, seqs := range tempSeqMap {
				m.syncedMaxSeqs[conversationID] = seqs[1]
			}
			// Reset tempSeqMap and msgNum to handle the next batch
			tempSeqMap = make(map[string][2]int64, 50)
			msgNum = 0
		}
	}

	// Handle remaining messages to ensure all are synced
	if len(tempSeqMap) > 0 {
		resp, err := m.pullMsgBySeqRange(ctx, tempSeqMap, syncMsgNum)
		if err != nil {
			log.ZError(ctx, "syncMsgFromServer error", err, "tempSeqMap", tempSeqMap)
			return err
		}
		_ = m.triggerConversation(ctx, resp.Msgs)
		_ = m.triggerNotification(ctx, resp.NotificationMsgs)
		for conversationID, seqs := range tempSeqMap {
			m.syncedMaxSeqs[conversationID] = seqs[1]
		}
	}

	return nil
}

// Fragment synchronization message, seq refresh after successful trigger
func (m *MsgSyncer) syncAndTriggerReinstallMsgs(ctx context.Context, seqMap map[string][2]int64, isFirst bool, syncMsgNum int64) error {
	if len(seqMap) > 0 {
		log.ZDebug(ctx, "current sync seqMap", "seqMap", seqMap)
		var (
			tempSeqMap = make(map[string][2]int64, 50)
			msgNum     = 0
			total      = 0
			gr         *errgroup.Group
		)
		if isFirst {
			total = len(seqMap)
		}
		gr, _ = errgroup.WithContext(ctx)
		gr.SetLimit(pullMsgGoroutineLimit)
		for k, v := range seqMap {
			oneConversationSyncNum := min(v[1]-v[0]+1, syncMsgNum)
			tempSeqMap[k] = v
			if IsNotification(k) {
				msgNum += int(oneConversationSyncNum)
			} else {
				// For regular conversations, ensure msgNum is the minimum of oneConversationSyncNum and syncMsgNum
				msgNum += int(min(oneConversationSyncNum, syncMsgNum))
			}
			if msgNum >= SplitPullMsgNum {
				tpSeqMap := make(map[string][2]int64, len(tempSeqMap))
				for k, v := range tempSeqMap {
					tpSeqMap[k] = v
				}

				gr.Go(func() error {
					resp, err := m.pullMsgBySeqRange(ctx, tpSeqMap, syncMsgNum)
					if err != nil {
						log.ZError(ctx, "syncMsgFromServer err", err, "tempSeqMap", tpSeqMap)
						return err
					}
					m.checkMessagesAndGetLastMessage(ctx, resp.Msgs)
					_ = m.triggerReinstallConversation(ctx, resp.Msgs, total)
					_ = m.triggerNotification(ctx, resp.NotificationMsgs)
					for conversationID, seqs := range tpSeqMap {
						m.syncedMaxSeqsLock.Lock()
						m.syncedMaxSeqs[conversationID] = seqs[1]
						m.syncedMaxSeqsLock.Unlock()
					}
					return nil
				})

				tempSeqMap = make(map[string][2]int64, 50)
				msgNum = 0
			}
		}
		if len(tempSeqMap) > 0 && msgNum > 0 {
			gr.Go(func() error {
				resp, err := m.pullMsgBySeqRange(ctx, tempSeqMap, syncMsgNum)
				if err != nil {
					log.ZError(ctx, "syncMsgFromServer err", err, "seqMap", seqMap)
					return err
				}
				m.checkMessagesAndGetLastMessage(ctx, resp.Msgs)
				_ = m.triggerReinstallConversation(ctx, resp.Msgs, total)
				_ = m.triggerNotification(ctx, resp.NotificationMsgs)
				for conversationID, seqs := range tempSeqMap {
					m.syncedMaxSeqsLock.Lock()
					m.syncedMaxSeqs[conversationID] = seqs[1]
					m.syncedMaxSeqsLock.Unlock()
				}
				return nil
			})
		}
		if err := gr.Wait(); err != nil {
			return err
		}

	} else {
		log.ZDebug(ctx, "noting conversation to sync", "syncMsgNum", syncMsgNum)
	}
	return nil
}
func (m *MsgSyncer) checkMessagesAndGetLastMessage(ctx context.Context, messages map[string]*sdkws.PullMsgs) {
	var conversationIDs []string

	for conversationID, message := range messages {
		allInValid := true
		for _, data := range message.Msgs {
			if data.Status < constant.MsgStatusHasDeleted {
				allInValid = false
				break
			}
		}
		if allInValid {
			conversationIDs = append(conversationIDs, conversationID)
		}
	}
	if len(conversationIDs) > 0 {
		resp, err := m.fetchLatestValidMessages(ctx, conversationIDs)
		if err != nil {
			log.ZError(ctx, "fetchLatestValidMessages", err, "conversationIDs", conversationIDs)
			return
		}
		for conversationID, message := range resp.Msgs {
			messages[conversationID] = &sdkws.PullMsgs{Msgs: []*sdkws.MsgData{message}}
		}
	}

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

func (m *MsgSyncer) pullMsgBySeqRange(ctx context.Context, seqMap map[string][2]int64, syncMsgNum int64) (resp *sdkws.PullMessageBySeqsResp, err error) {
	log.ZDebug(ctx, "pullMsgBySeqRange", "seqMap", seqMap, "syncMsgNum", syncMsgNum)

	req := sdkws.PullMessageBySeqsReq{UserID: m.loginUserID}
	for conversationID, seqs := range seqMap {
		req.SeqRanges = append(req.SeqRanges, &sdkws.SeqRange{
			ConversationID: conversationID,
			Begin:          seqs[0],
			End:            seqs[1],
			Num:            syncMsgNum,
		})
	}
	resp = &sdkws.PullMessageBySeqsResp{}
	if err := m.longConnMgr.SendReqWaitResp(ctx, &req, constant.PullMsgByRange, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *MsgSyncer) fetchLatestValidMessages(ctx context.Context, conversationID []string) (resp *msg.GetLastMessageResp, err error) {
	log.ZDebug(ctx, "fetchLatestValidMessages", "conversationID", conversationID)

	req := msg.GetLastMessageReq{
		UserID:          m.loginUserID,
		ConversationIDs: conversationID,
	}
	resp = &msg.GetLastMessageResp{}
	if err := m.longConnMgr.SendReqWaitResp(ctx, &req, constant.PullConvLastMessage, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// synchronizes messages by SEQs.
func (m *MsgSyncer) syncMsgBySeqs(ctx context.Context, conversationID string, seqsNeedSync []int64) (allMsgs []*sdkws.MsgData, err error) {
	pullMsgReq := sdkws.PullMessageBySeqsReq{}
	pullMsgReq.UserID = m.loginUserID
	split := constant.SplitPullMsgNum
	seqsList := m.splitSeqs(split, seqsNeedSync)
	for i := 0; i < len(seqsList); {
		var pullMsgResp sdkws.PullMessageBySeqsResp
		err := m.longConnMgr.SendReqWaitResp(ctx, &pullMsgReq, constant.PullMsgByRange, &pullMsgResp)
		if err != nil {
			log.ZError(ctx, "syncMsgFromServerSplit err", err, "pullMsgReq", pullMsgReq)
			continue
		}
		i++
		allMsgs = append(allMsgs, pullMsgResp.Msgs[conversationID].Msgs...)
	}
	return allMsgs, nil
}

// triggers a conversation with a new message.
func (m *MsgSyncer) triggerConversation(ctx context.Context, msgs map[string]*sdkws.PullMsgs) error {
	if len(msgs) > 0 {
		err := common.DispatchNewMessage(ctx, sdk_struct.CmdNewMsgComeToConversation{Msgs: msgs}, m.conversationEventQueue)
		if err != nil {
			log.ZError(ctx, "triggerCmdNewMsgCome err", err, "msgs", msgs)
		}
		log.ZDebug(ctx, "triggerConversation", "msgs", msgs)
		return err
	} else {
		log.ZDebug(ctx, "triggerConversation is nil", "msgs", msgs)
	}
	return nil
}

// triggers a conversation with a new message.
func (m *MsgSyncer) triggerReinstallConversation(ctx context.Context, msgs map[string]*sdkws.PullMsgs, total int) (err error) {
	if len(msgs) > 0 {
		err = common.DispatchMsgSyncInReinstall(ctx, sdk_struct.CmdMsgSyncInReinstall{
			Msgs:  msgs,
			Total: total,
		}, m.conversationEventQueue)
		if err != nil {
			log.ZError(ctx, "DispatchMsgSyncInReinstall err", err, "msgs", msgs)
		}
		log.ZDebug(ctx, "triggerReinstallConversation", "length", len(msgs))
		return err
	} else {
		log.ZDebug(ctx, "triggerReinstallConversation is nil")
	}
	return nil
}

func (m *MsgSyncer) triggerNotification(ctx context.Context, msgs map[string]*sdkws.PullMsgs) error {
	if len(msgs) > 0 {
		common.DispatchNotification(ctx, sdk_struct.CmdNewMsgComeToConversation{Msgs: msgs}, m.conversationEventQueue)
	} else {
		log.ZDebug(ctx, "triggerNotification is nil", "notifications", msgs)
	}
	return nil

}

func (m *MsgSyncer) SyncAndSortConversations(ctx context.Context, reinstalled bool) (map[string]*msg.Seqs, *sort_conversation.SortConversationList, error) {
	startTime := time.Now()
	log.ZDebug(ctx, "start SyncConversationHashReadSeqs")

	resp := msg.GetConversationsHasReadAndMaxSeqResp{}
	req := msg.GetConversationsHasReadAndMaxSeqReq{UserID: m.loginUserID, ReturnPinned: reinstalled}
	err := m.longConnMgr.SendReqWaitResp(ctx, &req, constant.GetConvMaxReadSeq, &resp)
	if err != nil {
		log.ZWarn(ctx, "SendReqWaitResp err", err)
		return nil, nil, err
	}
	seqs := resp.Seqs
	log.ZDebug(ctx, "getServerHasReadAndMaxSeqs completed", "duration", time.Since(startTime).Seconds())

	if len(seqs) == 0 {
		return nil, nil, nil
	}

	stepStartTime := time.Now()
	conversationsOnLocal, err := m.db.GetAllConversations(ctx)
	if err != nil {
		log.ZWarn(ctx, "get all conversations err", err)
		return nil, nil, err
	}
	log.ZDebug(ctx, "GetAllConversations completed", "duration", time.Since(stepStartTime).Seconds())

	conversationsOnLocalMap := datautil.SliceToMap(conversationsOnLocal, func(e *model_struct.LocalConversation) string {
		return e.ConversationID
	})

	var (
		list                  []*sort_conversation.ConversationMetaData
		pinnedConversationIDs []string
	)
	pinnedConversationIDsMap := datautil.SliceSetAny(resp.PinnedConversationIDs, func(e string) string {
		return e
	})
	stepStartTime = time.Now()
	for conversationID, v := range seqs {
		sortConversationMetaData := &sort_conversation.ConversationMetaData{
			ConversationID:    conversationID,
			LatestMsgSendTime: v.MaxSeqTime,
		}
		if _, ok := pinnedConversationIDsMap[conversationID]; ok {
			sortConversationMetaData.IsPinned = true
			pinnedConversationIDs = append(pinnedConversationIDs, conversationID)
		}

		if conversation, ok := conversationsOnLocalMap[conversationID]; ok {
			if conversation.IsPinned {
				sortConversationMetaData.IsPinned = true
				pinnedConversationIDs = append(pinnedConversationIDs, conversationID)
			}
			if conversation.DraftTextTime > 0 {
				sortConversationMetaData.DraftTextTime = conversation.DraftTextTime
			}
		}
		list = append(list, sortConversationMetaData)
	}

	sortConversationList := sort_conversation.NewSortConversationList(list, pinnedConversationIDs)

	log.ZDebug(ctx, "Process seqs completed", "duration", time.Since(stepStartTime).Seconds())
	return resp.Seqs, sortConversationList, nil
}
