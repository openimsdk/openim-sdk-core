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
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"time"

	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/log"
	utils2 "github.com/OpenIMSDK/tools/utils"
)

func (c *Conversation) Work(c2v common.Cmd2Value) {
	log.ZDebug(c2v.Ctx, "NotificationCmd start", "cmd", c2v.Cmd, "value", c2v.Value)
	defer log.ZDebug(c2v.Ctx, "NotificationCmd end", "cmd", c2v.Cmd, "value", c2v.Value)
	switch c2v.Cmd {
	case constant.CmdDeleteConversation:
		c.doDeleteConversation(c2v)
	case constant.CmdNewMsgCome:
		c.doMsgNew(c2v)
	case constant.CmdSuperGroupMsgCome:
		// c.doSuperGroupMsgNew(c2v)
	case constant.CmdUpdateConversation:
		c.doUpdateConversation(c2v)
	case constant.CmdUpdateMessage:
		c.doUpdateMessage(c2v)
	case constant.CmSyncReactionExtensions:
		// c.doSyncReactionExtensions(c2v)
	case constant.CmdNotification:
		c.doNotificationNew(c2v)
	}
}

func (c *Conversation) doDeleteConversation(c2v common.Cmd2Value) {
	node := c2v.Value.(common.DeleteConNode)
	ctx := c2v.Ctx
	// Mark messages related to this conversation for deletion
	err := c.db.UpdateMessageStatusBySourceID(context.Background(), node.SourceID, constant.MsgStatusHasDeleted, int32(node.SessionType))
	if err != nil {
		log.ZError(ctx, "setMessageStatusBySourceID", err)
		return
	}
	// Reset the session information, empty session
	err = c.db.ResetConversation(ctx, node.ConversationID)
	if err != nil {
		log.ZError(ctx, "ResetConversation err:", err)
	}
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{"", constant.TotalUnreadMessageChanged, ""}})
}

func (c *Conversation) doUpdateConversation(c2v common.Cmd2Value) {
	if c.ConversationListener == nil {
		// log.Error("internal", "not set conversationListener")
		return
	}
	ctx := c2v.Ctx
	node := c2v.Value.(common.UpdateConNode)
	switch node.Action {
	case constant.AddConOrUpLatMsg:
		var list []*model_struct.LocalConversation
		lc := node.Args.(model_struct.LocalConversation)
		oc, err := c.db.GetConversation(ctx, lc.ConversationID)
		if err == nil {
			// log.Info("this is old conversation", *oc)
			if lc.LatestMsgSendTime >= oc.LatestMsgSendTime { // The session update of asynchronous messages is subject to the latest sending time
				err := c.db.UpdateColumnsConversation(ctx, node.ConID, map[string]interface{}{"latest_msg_send_time": lc.LatestMsgSendTime, "latest_msg": lc.LatestMsg})
				if err != nil {
					// log.Error("internal", "updateConversationLatestMsgModel err: ", err)
				} else {
					oc.LatestMsgSendTime = lc.LatestMsgSendTime
					oc.LatestMsg = lc.LatestMsg
					list = append(list, oc)
					c.ConversationListener.OnConversationChanged(utils.StructToJsonString(list))
				}
			}
		} else {
			// log.Info("this is new conversation", lc)
			err4 := c.db.InsertConversation(ctx, &lc)
			if err4 != nil {
				// log.Error("internal", "insert new conversation err:", err4.Error())
			} else {
				list = append(list, &lc)
				c.ConversationListener.OnNewConversation(utils.StructToJsonString(list))
			}
		}

	case constant.UnreadCountSetZero:
		if err := c.db.UpdateColumnsConversation(ctx, node.ConID, map[string]interface{}{"unread_count": 0}); err != nil {
			log.ZError(ctx, "updateConversationUnreadCountModel err", err, "conversationID", node.ConID)
		} else {
			totalUnreadCount, err := c.db.GetTotalUnreadMsgCountDB(ctx)
			if err == nil {
				c.ConversationListener.OnTotalUnreadMessageCountChanged(totalUnreadCount)
			} else {
				log.ZError(ctx, "getTotalUnreadMsgCountDB err", err)
			}

		}
	// case ConChange:
	//	err, list := u.getAllConversationListModel()
	//	if err != nil {
	//		sdkLog("getAllConversationListModel database err:", err.Error())
	//	} else {
	//		if list == nil {
	//			u.ConversationListenerx.OnConversationChanged(structToJsonString([]ConversationStruct{}))
	//		} else {
	//			u.ConversationListenerx.OnConversationChanged(structToJsonString(list))
	//
	//		}
	//	}
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
			c.ConversationListener.OnTotalUnreadMessageCountChanged(totalUnreadCount)
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
				// log.Error("internal", "getConversationTypeByGroupID database err:", err.Error())
				return
			}
			lc.GroupID = st.SourceID
			lc.ConversationID = conversationID
			lc.ConversationType = conversationType
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
			c.ConversationListener.OnConversationChanged(utils.StructToJsonStringDefault(newCList))
		}
	case constant.NewCon:
		cidList := node.Args.([]string)
		cLists, err := c.db.GetMultipleConversationDB(ctx, cidList)
		if err != nil {
			// log.Error("internal", "getMultipleConversationModel err :", err.Error())
		} else {
			if cLists != nil {
				// log.Info("internal", "getMultipleConversationModel success :", cLists)
				c.ConversationListener.OnNewConversation(utils.StructToJsonString(cLists))
			}
		}
	case constant.ConChangeDirect:
		cidList := node.Args.(string)
		c.ConversationListener.OnConversationChanged(cidList)

	case constant.NewConDirect:
		cidList := node.Args.(string)
		// log.Debug("internal", "NewConversation", cidList)
		c.ConversationListener.OnNewConversation(cidList)

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
			c.ConversationListener.OnNewConversation(utils.StructToJsonString(result))
		}
	case constant.SyncConversation:

		c.SyncAllConversations(ctx)
		err := c.SyncConversationUnreadCount(ctx)
		if err != nil {
			// log.Error(operationID, "reconn sync conversation unread count err", err.Error())
		}
		totalUnreadCount, err := c.db.GetTotalUnreadMsgCountDB(ctx)
		if err != nil {
			// log.Error("internal", "TotalUnreadMessageChanged database err:", err.Error())
		} else {
			c.ConversationListener.OnTotalUnreadMessageCountChanged(totalUnreadCount)
		}

	}
}

func (c *Conversation) doUpdateMessage(c2v common.Cmd2Value) {
	if c.ConversationListener == nil {
		// log.Error("internal", "not set conversationListener")
		return
	}

	node := c2v.Value.(common.UpdateMessageNode)
	ctx := c2v.Ctx
	switch node.Action {
	case constant.UpdateMsgFaceUrlAndNickName:
		args := node.Args.(common.UpdateMessageInfo)
		if args.GroupID == "" {
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
		} else {
			conversationID := c.getConversationIDBySessionType(args.GroupID, constant.SuperGroupChatType)
			err := c.db.UpdateMsgSenderFaceURLAndSenderNickname(ctx, conversationID, args.UserID, args.FaceURL, args.Nickname)
			if err != nil {
				log.ZError(ctx, "UpdateMsgSenderFaceURLAndSenderNickname err", err)
			}
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
	if msg.SendTime < c.LoginTime() || c.LoginTime() == 0 {
		log.ZWarn(ctx, "ignore notification", nil, "clientMsgID", msg.ClientMsgID, "serverMsgID",
			msg.ServerMsgID, "seq", msg.Seq, "contentType", msg.ContentType,
			"sendTime", msg.SendTime, "loginTime", c.full.Group().LoginTime())
		return
	}
	if c.msgListener == nil {
		log.ZError(ctx, "msgListner is nil", nil)
		return
	}
	//var notification sdkws.ConversationChangedNotification
	tips := &sdkws.ConversationUpdateTips{}
	if err := utils.UnmarshalNotificationElem(msg.Content, tips); err != nil {
		log.ZError(ctx, "UnmarshalNotificationElem err", err, "msg", msg)
		return
	}
	go func() {
		c.SyncConversations(ctx, tips.ConversationIDList)
	}()
}

func (c *Conversation) DoConversationIsPrivateChangedNotification(ctx context.Context, msg *sdkws.MsgData) {
	if msg.SendTime < c.LoginTime() || c.LoginTime() == 0 {
		log.ZWarn(ctx, "ignore notification", nil, "clientMsgID", msg.ClientMsgID, "serverMsgID",
			msg.ServerMsgID, "seq", msg.Seq, "contentType", msg.ContentType,
			"sendTime", msg.SendTime, "loginTime", c.full.Group().LoginTime())
		return
	}
	if c.msgListener == nil {
		log.ZError(ctx, "msgListner is nil", nil)
		return
	}
	tips := &sdkws.ConversationSetPrivateTips{}
	if err := utils.UnmarshalNotificationElem(msg.Content, tips); err != nil {
		log.ZError(ctx, "UnmarshalNotificationElem err", err, "msg", msg)
		return
	}
	go func() {
		c.SyncConversations(ctx, []string{tips.ConversationID})
	}()
}

func (c *Conversation) doNotificationNew(c2v common.Cmd2Value) {
	ctx := c2v.Ctx
	allMsg := c2v.Value.(sdk_struct.CmdNewMsgComeToConversation).Msgs
	syncFlag := c2v.Value.(sdk_struct.CmdNewMsgComeToConversation).SyncFlag
	switch syncFlag {
	case constant.MsgSyncBegin:
		c.startTime = time.Now()
		c.ConversationListener.OnSyncServerStart()
		if err := c.SyncConversationHashReadSeqs(ctx); err != nil {
			log.ZError(ctx, "SyncConversationHashReadSeqs err", err)
		}
		//clear SubscriptionStatusMap
		c.cache.SubscriptionStatusMap.Range(func(key, value interface{}) bool {
			c.cache.SubscriptionStatusMap.Delete(key)
			return true
		})
		for _, syncFunc := range []func(c context.Context) error{
			c.user.SyncLoginUserInfo,
			c.friend.SyncAllBlackList, c.friend.SyncAllFriendList, c.friend.SyncAllFriendApplication, c.friend.SyncAllSelfFriendApplication,
			c.group.SyncAllJoinedGroups, c.group.SyncAllAdminGroupApplication, c.group.SyncAllSelfGroupApplication, c.group.SyncAllJoinedGroupMembers,
		} {
			go func(syncFunc func(c context.Context) error) {
				_ = syncFunc(ctx)
			}(syncFunc)
		}
	case constant.MsgSyncFailed:
		c.ConversationListener.OnSyncServerFailed()
	case constant.MsgSyncEnd:
		log.ZDebug(ctx, "MsgSyncEnd", "time", time.Since(c.startTime).Milliseconds())
		defer c.ConversationListener.OnSyncServerFinish()
		go c.SyncAllConversations(ctx)
	}

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
					c.user.DoNotification(ctx, v, c.cache.UpdateStatus)
				} else if utils2.Contain(v.ContentType, constant.GroupApplicationRejectedNotification, constant.GroupApplicationAcceptedNotification, constant.JoinGroupApplicationNotification) {
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
