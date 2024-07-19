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

package conversation_msg

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
)

const (
	syncWait = iota
	asyncNoWait
	asyncWait
)

func (c *Conversation) Work(c2v common.Cmd2Value) {
	log.ZDebug(c2v.Ctx, "NotificationCmd start", "cmd", c2v.Cmd, "value", c2v.Value)
	defer log.ZDebug(c2v.Ctx, "NotificationCmd end", "cmd", c2v.Cmd, "value", c2v.Value)
	switch c2v.Cmd {
	case constant.CmdNewMsgCome:
		c.doMsgNew(c2v)
	case constant.CmdUpdateConversation:
		c.doUpdateConversation(c2v)
	case constant.CmdUpdateMessage:
		c.doUpdateMessage(c2v)
	case constant.CmSyncReactionExtensions:
	case constant.CmdNotification:
		c.doNotification(c2v)
	case constant.CmdSyncData:
		c.syncData(c2v)
	case constant.CmdSyncFlag:
		c.syncFlag(c2v)
	}
}

func (c *Conversation) syncFlag(c2v common.Cmd2Value) {
	ctx := c2v.Ctx
	syncFlag := c2v.Value.(sdk_struct.CmdNewMsgComeToConversation).SyncFlag
	switch syncFlag {
	case constant.AppDataSyncStart:
		log.ZDebug(ctx, "AppDataSyncStart")
		c.startTime = time.Now()
		c.ConversationListener().OnSyncServerStart(true)
		asyncWaitFunctions := []func(c context.Context) error{
			c.group.SyncAllJoinedGroupsAndMembers,
			c.friend.IncrSyncFriends,
		}
		runSyncFunctions(ctx, asyncWaitFunctions, asyncWait, c.ConversationListener().OnSyncServerProgress)

		syncWaitFunctions := []func(c context.Context) error{
			c.IncrSyncConversations,
			c.SyncAllConversationHashReadSeqs,
		}
		runSyncFunctions(ctx, syncWaitFunctions, syncWait, c.ConversationListener().OnSyncServerProgress)
		log.ZWarn(ctx, "core data sync over", nil, "cost time", time.Since(c.startTime).Seconds())

		asyncNoWaitFunctions := []func(c context.Context) error{
			c.user.SyncLoginUserInfoWithoutNotice,
			c.friend.SyncAllBlackListWithoutNotice,
			c.friend.SyncAllFriendApplicationWithoutNotice,
			c.friend.SyncAllSelfFriendApplicationWithoutNotice,
			c.group.SyncAllAdminGroupApplicationWithoutNotice,
			c.group.SyncAllSelfGroupApplicationWithoutNotice,
			c.user.SyncAllCommandWithoutNotice,
		}
		runSyncFunctions(ctx, asyncNoWaitFunctions, asyncNoWait, c.ConversationListener().OnSyncServerProgress)

	case constant.AppDataSyncFinish:
		log.ZDebug(ctx, "AppDataSyncFinish", "time", time.Since(c.startTime).Milliseconds())
		c.ConversationListener().OnSyncServerFinish(true)
	case constant.MsgSyncBegin:
		log.ZDebug(ctx, "MsgSyncBegin")
		c.ConversationListener().OnSyncServerStart(false)

		c.syncData(c2v)

	case constant.MsgSyncFailed:
		c.ConversationListener().OnSyncServerFailed(false)
	case constant.MsgSyncEnd:
		log.ZDebug(ctx, "MsgSyncEnd", "time", time.Since(c.startTime).Milliseconds())
		c.ConversationListener().OnSyncServerFinish(false)
	}
}

func (c *Conversation) doNotification(c2v common.Cmd2Value) {
	ctx := c2v.Ctx
	allMsg := c2v.Value.(sdk_struct.CmdNewMsgComeToConversation).Msgs

	for conversationID, msgs := range allMsg {
		log.ZDebug(ctx, "notification handling", "conversationID", conversationID, "msgs", msgs)
		if len(msgs.Msgs) != 0 {
			lastMsg := msgs.Msgs[len(msgs.Msgs)-1]
			log.ZDebug(ctx, "SetNotificationSeq", "conversationID", conversationID, "seq", lastMsg.Seq)
			if lastMsg.Seq != 0 {
				if err := c.db.SetNotificationSeq(ctx, conversationID, lastMsg.Seq); err != nil {
					log.ZError(ctx, "SetNotificationSeq err", err, "conversationID", conversationID, "lastMsg", lastMsg)
				}
			}
		}
		for _, v := range msgs.Msgs {
			switch {
			case v.ContentType == constant.ConversationChangeNotification:
				c.DoConversationChangedNotification(ctx, v)
			case v.ContentType == constant.ConversationPrivateChatNotification:
				c.DoConversationIsPrivateChangedNotification(ctx, v)
			case v.ContentType == constant.ConversationUnreadNotification:
				var tips sdkws.ConversationHasReadTips
				_ = json.Unmarshal(v.Content, &tips)
				c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: tips.ConversationID, Action: constant.UnreadCountSetZero}})
				c.db.DeleteConversationUnreadMessageList(ctx, tips.ConversationID, tips.UnreadCountTime)
				c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.ConChange, Args: []string{tips.ConversationID}}})
				continue
			case v.ContentType == constant.BusinessNotification:
				c.business.DoNotification(ctx, v)
				continue
			case v.ContentType == constant.RevokeNotification:
				c.doRevokeMsg(ctx, v)
			case v.ContentType == constant.ClearConversationNotification:
				c.doClearConversations(ctx, v)
			case v.ContentType == constant.DeleteMsgsNotification:
				c.doDeleteMsgs(ctx, v)
			case v.ContentType == constant.HasReadReceipt:
				c.doReadDrawing(ctx, v)
			}

			switch v.SessionType {
			case constant.SingleChatType:
				if v.ContentType > constant.FriendNotificationBegin && v.ContentType < constant.FriendNotificationEnd {
					c.friend.DoNotification(ctx, v)
				} else if v.ContentType > constant.UserNotificationBegin && v.ContentType < constant.UserNotificationEnd {
					c.user.DoNotification(ctx, v)
				} else if datautil.Contain(v.ContentType, constant.GroupApplicationRejectedNotification, constant.GroupApplicationAcceptedNotification, constant.JoinGroupApplicationNotification) {
					c.group.DoNotification(ctx, v)
				} else if v.ContentType > constant.SignalingNotificationBegin && v.ContentType < constant.SignalingNotificationEnd {

					continue
				}
			case constant.GroupChatType, constant.SuperGroupChatType:
				if v.ContentType > constant.GroupNotificationBegin && v.ContentType < constant.GroupNotificationEnd {
					c.group.DoNotification(ctx, v)
				} else if v.ContentType > constant.SignalingNotificationBegin && v.ContentType < constant.SignalingNotificationEnd {
					continue
				}
			}
		}
	}

}

func (c *Conversation) getConversationLatestMsgClientID(latestMsg string) string {
	msg := &sdk_struct.MsgStruct{}
	if err := json.Unmarshal([]byte(latestMsg), msg); err != nil {
		log.ZError(context.Background(), "getConversationLatestMsgClientID", err, "latestMsg", latestMsg)
	}
	return msg.ClientMsgID
}

func (c *Conversation) doUpdateConversation(c2v common.Cmd2Value) {
	ctx := c2v.Ctx
	node := c2v.Value.(common.UpdateConNode)
	switch node.Action {
	case constant.AddConOrUpLatMsg:
		var list []*model_struct.LocalConversation
		lc := node.Args.(model_struct.LocalConversation)
		oc, err := c.db.GetConversation(ctx, lc.ConversationID)
		if err == nil {
			// log.Info("this is old conversation", *oc)
			if lc.LatestMsgSendTime >= oc.LatestMsgSendTime || c.getConversationLatestMsgClientID(lc.LatestMsg) == c.getConversationLatestMsgClientID(oc.LatestMsg) { // The session update of asynchronous messages is subject to the latest sending time
				err := c.db.UpdateColumnsConversation(ctx, node.ConID, map[string]interface{}{"latest_msg_send_time": lc.LatestMsgSendTime, "latest_msg": lc.LatestMsg})
				if err != nil {
					log.ZError(ctx, "updateConversationLatestMsgModel", err, "conversationID", node.ConID)
				} else {
					oc.LatestMsgSendTime = lc.LatestMsgSendTime
					oc.LatestMsg = lc.LatestMsg
					list = append(list, oc)
					c.ConversationListener().OnConversationChanged(utils.StructToJsonString(list))
				}
			}
		} else {
			// log.Info("this is new conversation", lc)
			err4 := c.db.InsertConversation(ctx, &lc)
			if err4 != nil {
				// log.Error("internal", "insert new conversation err:", err4.Error())
			} else {
				list = append(list, &lc)
				c.ConversationListener().OnNewConversation(utils.StructToJsonString(list))
			}
		}

	case constant.UnreadCountSetZero:
		if err := c.db.UpdateColumnsConversation(ctx, node.ConID, map[string]interface{}{"unread_count": 0}); err != nil {
			log.ZError(ctx, "updateConversationUnreadCountModel err", err, "conversationID", node.ConID)
		} else {
			totalUnreadCount, err := c.db.GetTotalUnreadMsgCountDB(ctx)
			if err == nil {
				c.ConversationListener().OnTotalUnreadMessageCountChanged(totalUnreadCount)
			} else {
				log.ZError(ctx, "getTotalUnreadMsgCountDB err", err)
			}

		}
	case constant.IncrUnread:
		err := c.db.IncrConversationUnreadCount(ctx, node.ConID)
		if err != nil {
			// log.Error("internal", "incrConversationUnreadCount database err:", err.Error())
			return
		}
	case constant.TotalUnreadMessageChanged:
		totalUnreadCount, err := c.db.GetTotalUnreadMsgCountDB(ctx)
		if err != nil {
			// log.Error("internal", "TotalUnreadMessageChanged database err:", err.Error())
		} else {
			c.ConversationListener().OnTotalUnreadMessageCountChanged(totalUnreadCount)
		}
	case constant.UpdateConFaceUrlAndNickName:
		var lc model_struct.LocalConversation
		st := node.Args.(common.SourceIDAndSessionType)
		log.ZInfo(ctx, "UpdateConFaceUrlAndNickName", "st", st)
		switch st.SessionType {
		case constant.SingleChatType:
			lc.UserID = st.SourceID
			lc.ConversationID = c.getConversationIDBySessionType(st.SourceID, constant.SingleChatType)
			lc.ConversationType = constant.SingleChatType
		case constant.SuperGroupChatType:
			conversationID, conversationType, err := c.getConversationTypeByGroupID(ctx, st.SourceID)
			if err != nil {
				return
			}
			lc.GroupID = st.SourceID
			lc.ConversationID = conversationID
			lc.ConversationType = conversationType
		case constant.NotificationChatType:
			lc.UserID = st.SourceID
			lc.ConversationID = c.getConversationIDBySessionType(st.SourceID, constant.NotificationChatType)
			lc.ConversationType = constant.NotificationChatType
		default:
			log.ZError(ctx, "not support sessionType", nil, "sessionType", st.SessionType)
			return
		}
		lc.ShowName = st.Nickname
		lc.FaceURL = st.FaceURL
		err := c.db.UpdateConversation(ctx, &lc)
		if err != nil {
			// log.Error("internal", "setConversationFaceUrlAndNickName database err:", err.Error())
			return
		}
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: lc.ConversationID, Action: constant.ConChange, Args: []string{lc.ConversationID}}})

	case constant.UpdateLatestMessageChange:
		conversationID := node.ConID
		var latestMsg sdk_struct.MsgStruct
		l, err := c.db.GetConversation(ctx, conversationID)
		if err != nil {
			log.ZError(ctx, "getConversationLatestMsgModel err", err, "conversationID", conversationID)
		} else {
			err := json.Unmarshal([]byte(l.LatestMsg), &latestMsg)
			if err != nil {
				log.ZError(ctx, "latestMsg,Unmarshal err", err)
			} else {
				latestMsg.IsRead = true
				newLatestMessage := utils.StructToJsonString(latestMsg)
				err = c.db.UpdateColumnsConversation(ctx, node.ConID, map[string]interface{}{"latest_msg_send_time": latestMsg.SendTime, "latest_msg": newLatestMessage})
				if err != nil {
					log.ZError(ctx, "updateConversationLatestMsgModel err", err)
				}
			}
		}
	case constant.ConChange:
		conversationIDs := node.Args.([]string)
		conversations, err := c.db.GetMultipleConversationDB(ctx, conversationIDs)
		if err != nil {
			log.ZError(ctx, "getMultipleConversationModel err", err)
		} else {
			var newCList []*model_struct.LocalConversation
			for _, v := range conversations {
				if v.LatestMsgSendTime != 0 {
					newCList = append(newCList, v)
				}
			}
			c.ConversationListener().OnConversationChanged(utils.StructToJsonStringDefault(newCList))
		}
	case constant.NewCon:
		cidList := node.Args.([]string)
		cLists, err := c.db.GetMultipleConversationDB(ctx, cidList)
		if err != nil {
			// log.Error("internal", "getMultipleConversationModel err :", err.Error())
		} else {
			if cLists != nil {
				// log.Info("internal", "getMultipleConversationModel success :", cLists)
				c.ConversationListener().OnNewConversation(utils.StructToJsonString(cLists))
			}
		}
	case constant.ConChangeDirect:
		cidList := node.Args.(string)
		c.ConversationListener().OnConversationChanged(cidList)

	case constant.NewConDirect:
		cidList := node.Args.(string)
		// log.Debug("internal", "NewConversation", cidList)
		c.ConversationListener().OnNewConversation(cidList)

	case constant.ConversationLatestMsgHasRead:
		hasReadMsgList := node.Args.(map[string][]string)
		var result []*model_struct.LocalConversation
		var latestMsg sdk_struct.MsgStruct
		var lc model_struct.LocalConversation
		for conversationID, msgIDList := range hasReadMsgList {
			LocalConversation, err := c.db.GetConversation(ctx, conversationID)
			if err != nil {
				// log.Error("internal", "get conversation err", err.Error(), conversationID)
				continue
			}
			err = utils.JsonStringToStruct(LocalConversation.LatestMsg, &latestMsg)
			if err != nil {
				// log.Error("internal", "JsonStringToStruct err", err.Error(), conversationID)
				continue
			}
			if utils.IsContain(latestMsg.ClientMsgID, msgIDList) {
				latestMsg.IsRead = true
				lc.ConversationID = conversationID
				lc.LatestMsg = utils.StructToJsonString(latestMsg)
				LocalConversation.LatestMsg = utils.StructToJsonString(latestMsg)
				err := c.db.UpdateConversation(ctx, &lc)
				if err != nil {
					// log.Error("internal", "UpdateConversation database err:", err.Error())
					continue
				} else {
					result = append(result, LocalConversation)
				}
			}
		}
		if result != nil {
			// log.Info("internal", "getMultipleConversationModel success :", result)
			c.ConversationListener().OnNewConversation(utils.StructToJsonString(result))
		}
	case constant.SyncConversation:

	}
}

func (c *Conversation) syncData(c2v common.Cmd2Value) {
	ctx := c2v.Ctx
	c.startTime = time.Now()
	//clear SubscriptionStatusMap
	//c.user.OnlineStatusCache.DeleteAll()

	// Synchronous sync functions
	syncFuncs := []func(c context.Context) error{
		c.SyncAllConversationHashReadSeqs,
	}

	runSyncFunctions(ctx, syncFuncs, syncWait, nil)

	// Asynchronous sync functions
	asyncFuncs := []func(c context.Context) error{
		c.user.SyncLoginUserInfo,
		c.friend.SyncAllBlackList,
		c.friend.SyncAllFriendApplication,
		c.friend.SyncAllSelfFriendApplication,
		c.group.SyncAllAdminGroupApplication,
		c.group.SyncAllSelfGroupApplication,
		c.user.SyncAllCommand,
		c.group.SyncAllJoinedGroupsAndMembers,
		c.friend.IncrSyncFriends,
		c.IncrSyncConversations,
	}

	runSyncFunctions(ctx, asyncFuncs, asyncNoWait, nil)
}

func runSyncFunctions(ctx context.Context, funcs []func(c context.Context) error, mode int, progressCallback func(progress int)) {
	totalFuncs := len(funcs)
	var wg sync.WaitGroup

	for i, fn := range funcs {
		switch mode {
		case asyncWait:
			wg.Add(1)
			go executeSyncFunction(ctx, fn, i, totalFuncs, progressCallback, &wg)
		case asyncNoWait:
			go executeSyncFunction(ctx, fn, i, totalFuncs, progressCallback, nil)
		case syncWait:
			executeSyncFunction(ctx, fn, i, totalFuncs, progressCallback, nil)
		}
	}

	if mode == asyncWait {
		wg.Wait()
	}
}

func executeSyncFunction(ctx context.Context, fn func(c context.Context) error, index, total int, progressCallback func(progress int), wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	funcName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	startTime := time.Now()
	err := fn(ctx)
	duration := time.Since(startTime)
	if err != nil {
		log.ZWarn(ctx, fmt.Sprintf("%s sync error", funcName), err, "duration", duration.Seconds())
	} else {
		log.ZDebug(ctx, fmt.Sprintf("%s completed successfully", funcName), "duration", duration.Seconds())
	}
	if progressCallback != nil {
		progress := int(float64(index+1) / float64(total) * 100)
		if progress == 0 {
			progress = 1
		}
		progressCallback(progress)
	}
}

func (c *Conversation) doUpdateMessage(c2v common.Cmd2Value) {
	node := c2v.Value.(common.UpdateMessageNode)
	ctx := c2v.Ctx
	switch node.Action {
	case constant.UpdateMsgFaceUrlAndNickName:
		args := node.Args.(common.UpdateMessageInfo)
		switch args.SessionType {
		case constant.SingleChatType:
			if args.UserID == c.loginUserID {
				conversationIDList, err := c.db.GetAllSingleConversationIDList(ctx)
				if err != nil {
					log.ZError(ctx, "GetAllSingleConversationIDList err", err)
					return
				} else {
					log.ZDebug(ctx, "get single conversationID list", "conversationIDList", conversationIDList)
					for _, conversationID := range conversationIDList {
						err := c.db.UpdateMsgSenderFaceURLAndSenderNickname(ctx, conversationID, args.UserID, args.FaceURL, args.Nickname)
						if err != nil {
							log.ZError(ctx, "UpdateMsgSenderFaceURLAndSenderNickname err", err)
							continue
						}
					}

				}
			} else {
				conversationID := c.getConversationIDBySessionType(args.UserID, constant.SingleChatType)
				err := c.db.UpdateMsgSenderFaceURLAndSenderNickname(ctx, conversationID, args.UserID, args.FaceURL, args.Nickname)
				if err != nil {
					log.ZError(ctx, "UpdateMsgSenderFaceURLAndSenderNickname err", err)
				}

			}
		case constant.SuperGroupChatType:
			conversationID := c.getConversationIDBySessionType(args.GroupID, constant.SuperGroupChatType)
			err := c.db.UpdateMsgSenderFaceURLAndSenderNickname(ctx, conversationID, args.UserID, args.FaceURL, args.Nickname)
			if err != nil {
				log.ZError(ctx, "UpdateMsgSenderFaceURLAndSenderNickname err", err)
			}
		case constant.NotificationChatType:
			conversationID := c.getConversationIDBySessionType(args.UserID, constant.NotificationChatType)
			err := c.db.UpdateMsgSenderFaceURLAndSenderNickname(ctx, conversationID, args.UserID, args.FaceURL, args.Nickname)
			if err != nil {
				log.ZError(ctx, "UpdateMsgSenderFaceURLAndSenderNickname err", err)
			}
		default:
			log.ZError(ctx, "not support sessionType", nil, "args", args)
			return
		}
	}

}

// funcation (c *Conversation) doSyncReactionExtensions(c2v common.Cmd2Value) {
//	if c.ConversationListener == nil {
//		// log.Error("internal", "not set conversationListener")
//		return
//	}
//	node := c2v.Value.(common.SyncReactionExtensionsNode)
//	ctx := mcontext.NewCtx(node.OperationID)
//	switch node.Action {
//	case constant.SyncMessageListReactionExtensions:
//		args := node.Args.(syncReactionExtensionParams)
//		// log.Error(node.OperationID, "come SyncMessageListReactionExtensions", args)
//		var reqList []server_api_params.OperateMessageListReactionExtensionsReq
//		for _, v := range args.MessageList {
//			var temp server_api_params.OperateMessageListReactionExtensionsReq
//			temp.ClientMsgID = v.ClientMsgID
//			temp.MsgFirstModifyTime = v.MsgFirstModifyTime
//			reqList = append(reqList, temp)
//		}
//		var apiReq server_api_params.GetMessageListReactionExtensionsReq
//		apiReq.SourceID = args.SourceID
//		apiReq.TypeKeyList = args.TypeKeyList
//		apiReq.SessionType = args.SessionType
//		apiReq.MessageReactionKeyList = reqList
//		apiReq.IsExternalExtensions = args.IsExternalExtension
//		apiReq.OperationID = node.OperationID
//		apiResp, err := util.CallApi[server_api_params.GetMessageListReactionExtensionsResp](ctx, constant.GetMessageListReactionExtensionsRouter, &apiReq)
//		if err != nil {
//			// log.NewError(node.OperationID, utils.GetSelfFuncName(), "getMessageListReactionExtensions err:", err.Error())
//			return
//		}
//		// for _, result := range apiResp {
//		// 	log.Warn(node.OperationID, "api return reslut is:", result.ClientMsgID, result.ReactionExtensionList)
//		// }
//		onLocal := funcation(data []*model_struct.LocalChatLogReactionExtensions) []*server_api_params.SingleMessageExtensionResult {
//			var result []*server_api_params.SingleMessageExtensionResult
//			for _, v := range data {
//				temp := new(server_api_params.SingleMessageExtensionResult)
//				tempMap := make(map[string]*sdkws.KeyValue)
//				_ = json.Unmarshal(v.LocalReactionExtensions, &tempMap)
//				if len(args.TypeKeyList) != 0 {
//					for s, _ := range tempMap {
//						if !utils.IsContain(s, args.TypeKeyList) {
//							delete(tempMap, s)
//						}
//					}
//				}
//
//				temp.ReactionExtensionList = tempMap
//				temp.ClientMsgID = v.ClientMsgID
//				result = append(result, temp)
//			}
//			return result
//		}(args.ExtendMessageList)
//		var onServer []*server_api_params.SingleMessageExtensionResult
//		for _, v := range *apiResp {
//			if v.ErrCode == 0 {
//				onServer = append(onServer, v)
//			}
//		}
//		aInBNot, _, sameA, _ := common.CheckReactionExtensionsDiff(onServer, onLocal)
//		for _, v := range aInBNot {
//			// log.Error(node.OperationID, "come InsertMessageReactionExtension", args, v.ClientMsgID)
//			if len(v.ReactionExtensionList) > 0 {
//				temp := model_struct.LocalChatLogReactionExtensions{ClientMsgID: v.ClientMsgID, LocalReactionExtensions: []byte(utils.StructToJsonString(v.ReactionExtensionList))}
//				err := c.db.InsertMessageReactionExtension(ctx, &temp)
//				if err != nil {
//					// log.Error(node.OperationID, "InsertMessageReactionExtension err:", err.Error())
//					continue
//				}
//			}
//			var changedKv []*sdkws.KeyValue
//			for _, value := range v.ReactionExtensionList {
//				changedKv = append(changedKv, value)
//			}
//			if len(changedKv) > 0 {
//				c.msgListener.OnRecvMessageExtensionsChanged(v.ClientMsgID, utils.StructToJsonString(changedKv))
//			}
//		}
//		// for _, result := range sameA {
//		// log.ZWarn(ctx, "result", result.ReactionExtensionList, result.ClientMsgID)
//		// }
//		for _, v := range sameA {
//			// log.Error(node.OperationID, "come sameA", v.ClientMsgID, v.ReactionExtensionList)
//			tempMap := make(map[string]*sdkws.KeyValue)
//			for _, extensions := range args.ExtendMessageList {
//				if v.ClientMsgID == extensions.ClientMsgID {
//					_ = json.Unmarshal(extensions.LocalReactionExtensions, &tempMap)
//					break
//				}
//			}
//			if len(v.ReactionExtensionList) == 0 {
//				err := c.db.DeleteMessageReactionExtension(ctx, v.ClientMsgID)
//				if err != nil {
//					// log.Error(node.OperationID, "DeleteMessageReactionExtension err:", err.Error())
//					continue
//				}
//				var deleteKeyList []string
//				for key, _ := range tempMap {
//					deleteKeyList = append(deleteKeyList, key)
//				}
//				if len(deleteKeyList) > 0 {
//					c.msgListener.OnRecvMessageExtensionsDeleted(v.ClientMsgID, utils.StructToJsonString(deleteKeyList))
//				}
//			} else {
//				deleteKeyList, changedKv := funcation(local, server map[string]*sdkws.KeyValue) ([]string, []*sdkws.KeyValue) {
//					var deleteKeyList []string
//					var changedKv []*sdkws.KeyValue
//					for k, v := range local {
//						ia, ok := server[k]
//						if ok {
//							//服务器不同的kv
//							if ia.Value != v.Value {
//								changedKv = append(changedKv, ia)
//							}
//						} else {
//							//服务器已经没有kv
//							deleteKeyList = append(deleteKeyList, k)
//						}
//					}
//					//从服务器新增的kv
//					for k, v := range server {
//						_, ok := local[k]
//						if !ok {
//							changedKv = append(changedKv, v)
//
//						}
//					}
//					return deleteKeyList, changedKv
//				}(tempMap, v.ReactionExtensionList)
//				extendMsg := model_struct.LocalChatLogReactionExtensions{ClientMsgID: v.ClientMsgID, LocalReactionExtensions: []byte(utils.StructToJsonString(v.ReactionExtensionList))}
//				err = c.db.UpdateMessageReactionExtension(ctx, &extendMsg)
//				if err != nil {
//					// log.Error(node.OperationID, "UpdateMessageReactionExtension err:", err.Error())
//					continue
//				}
//				if len(deleteKeyList) > 0 {
//					c.msgListener.OnRecvMessageExtensionsDeleted(v.ClientMsgID, utils.StructToJsonString(deleteKeyList))
//				}
//				if len(changedKv) > 0 {
//					c.msgListener.OnRecvMessageExtensionsChanged(v.ClientMsgID, utils.StructToJsonString(changedKv))
//				}
//			}
//			//err := c.db.GetAndUpdateMessageReactionExtension(v.ClientMsgID, v.ReactionExtensionList)
//			//if err != nil {
//			//	log.Error(node.OperationID, "GetAndUpdateMessageReactionExtension err:", err.Error())
//			//	continue
//			//}
//			//var changedKv []*server_api_params.KeyValue
//			//for _, value := range v.ReactionExtensionList {
//			//	changedKv = append(changedKv, value)
//			//}
//			//if len(changedKv) > 0 {
//			//	c.msgListener.OnRecvMessageExtensionsChanged(v.ClientMsgID, utils.StructToJsonString(changedKv))
//			//}
//		}
//	case constant.SyncMessageListTypeKeyInfo:
//		messageList := node.Args.([]*sdk_struct.MsgStruct)
//		var sourceID string
//		var sessionType int32
//		var reqList []server_api_params.OperateMessageListReactionExtensionsReq
//		var temp server_api_params.OperateMessageListReactionExtensionsReq
//		for _, v := range messageList {
//			//todo syncMessage must sync
//			message, err := c.db.GetMessage(ctx, "", v.ClientMsgID)
//			if err != nil {
//				// log.Error(node.OperationID, "GetMessageController err:", err.Error(), *v)
//				continue
//			}
//			temp.ClientMsgID = message.ClientMsgID
//			temp.MsgFirstModifyTime = message.MsgFirstModifyTime
//			reqList = append(reqList, temp)
//			switch message.SessionType {
//			case constant.SingleChatType:
//				sourceID = message.SendID + message.RecvID
//			case constant.NotificationChatType:
//				sourceID = message.RecvID
//			case constant.GroupChatType, constant.SuperGroupChatType:
//				sourceID = message.RecvID
//			}
//			sessionType = message.SessionType
//		}
//		var apiReq server_api_params.GetMessageListReactionExtensionsReq
//		apiReq.SourceID = sourceID
//		apiReq.SessionType = sessionType
//		apiReq.MessageReactionKeyList = reqList
//		apiReq.OperationID = node.OperationID
//		//var apiResp server_api_params.GetMessageListReactionExtensionsResp
//
//		apiResp, err := util.CallApi[server_api_params.GetMessageListReactionExtensionsResp](ctx, constant.GetMessageListReactionExtensionsRouter, &apiReq)
//		if err != nil {
//			// log.Error(node.OperationID, "GetMessageListReactionExtensions from server err:", err.Error(), apiReq)
//			return
//		}
//		var messageChangedList []*messageKvList
//		for _, v := range *apiResp {
//			if v.ErrCode == 0 {
//				var changedKv []*sdkws.KeyValue
//				var prefixTypeKey []string
//				extendMsg, _ := c.db.GetMessageReactionExtension(ctx, v.ClientMsgID)
//				localKV := make(map[string]*sdkws.KeyValue)
//				_ = json.Unmarshal(extendMsg.LocalReactionExtensions, &localKV)
//				for typeKey, value := range v.ReactionExtensionList {
//					oldValue, ok := localKV[typeKey]
//					if ok {
//						if !cmp.Equal(value, oldValue) {
//							localKV[typeKey] = value
//							prefixTypeKey = append(prefixTypeKey, getPrefixTypeKey(typeKey))
//							changedKv = append(changedKv, value)
//						}
//					} else {
//						localKV[typeKey] = value
//						prefixTypeKey = append(prefixTypeKey, getPrefixTypeKey(typeKey))
//						changedKv = append(changedKv, value)
//
//					}
//
//				}
//				extendMsg.LocalReactionExtensions = []byte(utils.StructToJsonString(localKV))
//				_ = c.db.UpdateMessageReactionExtension(ctx, extendMsg)
//				if len(changedKv) > 0 {
//					c.msgListener.OnRecvMessageExtensionsChanged(extendMsg.ClientMsgID, utils.StructToJsonString(changedKv))
//				}
//				prefixTypeKey = utils.RemoveRepeatedStringInList(prefixTypeKey)
//				if len(prefixTypeKey) > 0 && c.msgKvListener != nil {
//					var result []*sdk.SingleTypeKeyInfoSum
//					oneMessageChanged := new(messageKvList)
//					oneMessageChanged.ClientMsgID = extendMsg.ClientMsgID
//					for _, v := range prefixTypeKey {
//						singleResult := new(sdk.SingleTypeKeyInfoSum)
//						singleResult.TypeKey = v
//						for typeKey, value := range localKV {
//							if strings.HasPrefix(typeKey, v) {
//								singleTypeKeyInfo := new(sdk.SingleTypeKeyInfo)
//								err := json.Unmarshal([]byte(value.Value), singleTypeKeyInfo)
//								if err != nil {
//									continue
//								}
//								if _, ok := singleTypeKeyInfo.InfoList[c.loginUserID]; ok {
//									singleResult.IsContainSelf = true
//								}
//								for _, info := range singleTypeKeyInfo.InfoList {
//									v := *info
//									singleResult.InfoList = append(singleResult.InfoList, &v)
//								}
//								singleResult.Counter += singleTypeKeyInfo.Counter
//							}
//						}
//						result = append(result, singleResult)
//					}
//					oneMessageChanged.ChangedKvList = result
//					messageChangedList = append(messageChangedList, oneMessageChanged)
//				}
//			}
//		}
//		if len(messageChangedList) > 0 && c.msgKvListener != nil {
//			c.msgKvListener.OnMessageKvInfoChanged(utils.StructToJsonString(messageChangedList))
//		}
//
//	}
//
// }

func (c *Conversation) DoConversationChangedNotification(ctx context.Context, msg *sdkws.MsgData) {
	//var notification sdkws.ConversationChangedNotification
	tips := &sdkws.ConversationUpdateTips{}
	if err := utils.UnmarshalNotificationElem(msg.Content, tips); err != nil {
		log.ZError(ctx, "UnmarshalNotificationElem err", err, "msg", msg)
		return
	}

	err := c.IncrSyncConversations(ctx)
	if err != nil {
		log.ZWarn(ctx, "IncrSyncConversations err", err)
	}

}

func (c *Conversation) DoConversationIsPrivateChangedNotification(ctx context.Context, msg *sdkws.MsgData) {
	tips := &sdkws.ConversationSetPrivateTips{}
	if err := utils.UnmarshalNotificationElem(msg.Content, tips); err != nil {
		log.ZError(ctx, "UnmarshalNotificationElem err", err, "msg", msg)
		return
	}

	err := c.IncrSyncConversations(ctx)
	if err != nil {
		log.ZWarn(ctx, "IncrSyncConversations err", err)
	}

}
