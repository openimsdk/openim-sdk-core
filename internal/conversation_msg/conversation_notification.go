package conversation_msg

import (
	"context"
	"encoding/json"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	utils2 "github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/copier"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"strings"
	"time"
)

func (c *Conversation) NotificationCmd(c2v common.Cmd2Value) {
	log.ZDebug(c2v.Ctx, "NotificationCmd start", "cmd", c2v.Cmd, "value", c2v.Value)
	defer log.ZDebug(c2v.Ctx, "NotificationCmd end", "cmd", c2v.Cmd, "value", c2v.Value)
	switch c2v.Cmd {
	case constant.CmdDeleteConversation:
		c.doDeleteConversation(c2v)
	case constant.CmdNewMsgCome:
		c.doMsgNew(c2v)
		v := c2v.Value.(sdk_struct.CmdNewMsgComeToConversation)
		ctx := v.Ctx
		if c.msgListener == nil || c.ConversationListener == nil {
			for _, msg := range v.MsgList {
				if msg.ContentType > constant.SignalingNotificationBegin && msg.ContentType < constant.SignalingNotificationEnd {
					log.ZDebug(ctx, "signaling DoNotification", "signaling:", c.signaling, "msg:", msg)
					c.signaling.DoNotification(ctx, msg, c.GetCh())
				} else {
					log.ZDebug(ctx, "listener is nil, do nothing ", "msg:", msg)
				}
			}
		}
		switch v.SyncFlag {
		case constant.MsgSyncBegin:
			c.ConversationListener.OnSyncServerStart()
		case constant.MsgSyncFailed:
			c.ConversationListener.OnSyncServerFailed()
		}

	case constant.CmdSuperGroupMsgCome:
		c.doSuperGroupMsgNew(c2v)
	case constant.CmdUpdateConversation:
		c.doUpdateConversation(c2v)
	case constant.CmdUpdateMessage:
		c.doUpdateMessage(c2v)
	case constant.CmSyncReactionExtensions:
		c.doSyncReactionExtensions(c2v)
	}
}

func (c *Conversation) doDeleteConversation(c2v common.Cmd2Value) {
	node := c2v.Value.(common.DeleteConNode)
	ctx := c2v.Ctx
	//Mark messages related to this conversation for deletion
	err := c.db.UpdateMessageStatusBySourceID(context.Background(), node.SourceID, constant.MsgStatusHasDeleted, int32(node.SessionType))
	if err != nil {
		log.ZError(ctx, "setMessageStatusBySourceID", err)
		return
	}
	//Reset the session information, empty session
	err = c.db.ResetConversation(ctx, node.ConversationID)
	if err != nil {
		log.ZError(ctx, "ResetConversation err:", err)
	}
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{"", constant.TotalUnreadMessageChanged, ""}})
}

func (c *Conversation) doUpdateConversation(c2v common.Cmd2Value) {
	if c.ConversationListener == nil {
		log.Error("internal", "not set conversationListener")
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
			log.Info("this is old conversation", *oc)
			if lc.LatestMsgSendTime >= oc.LatestMsgSendTime { //The session update of asynchronous messages is subject to the latest sending time
				err := c.db.UpdateColumnsConversation(nil, node.ConID, map[string]interface{}{"latest_msg_send_time": lc.LatestMsgSendTime, "latest_msg": lc.LatestMsg})
				if err != nil {
					log.Error("internal", "updateConversationLatestMsgModel err: ", err)
				} else {
					oc.LatestMsgSendTime = lc.LatestMsgSendTime
					oc.LatestMsg = lc.LatestMsg
					list = append(list, oc)
					c.ConversationListener.OnConversationChanged(utils.StructToJsonString(list))
				}
			}
		} else {
			log.Info("this is new conversation", lc)
			err4 := c.db.InsertConversation(ctx, &lc)
			if err4 != nil {
				log.Error("internal", "insert new conversation err:", err4.Error())
			} else {
				list = append(list, &lc)
				c.ConversationListener.OnNewConversation(utils.StructToJsonString(list))
			}
		}

	case constant.UnreadCountSetZero:
		if err := c.db.UpdateColumnsConversation(ctx, node.ConID, map[string]interface{}{"unread_count": 0}); err != nil {
			log.Error("internal", "UpdateColumnsConversation err", err.Error(), node.ConID)
		} else {
			totalUnreadCount, err := c.db.GetTotalUnreadMsgCountDB(ctx)
			if err == nil {
				c.ConversationListener.OnTotalUnreadMessageCountChanged(totalUnreadCount)
			} else {
				log.Error("internal", "getTotalUnreadMsgCountModel err", err.Error(), node.ConID)
			}

		}
	//case ConChange:
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
			log.Error("internal", "incrConversationUnreadCount database err:", err.Error())
			return
		}
	case constant.TotalUnreadMessageChanged:
		totalUnreadCount, err := c.db.GetTotalUnreadMsgCountDB(ctx)
		if err != nil {
			log.Error("internal", "TotalUnreadMessageChanged database err:", err.Error())
		} else {
			c.ConversationListener.OnTotalUnreadMessageCountChanged(totalUnreadCount)
		}
	case constant.UpdateConFaceUrlAndNickName:
		var lc model_struct.LocalConversation
		st := node.Args.(common.SourceIDAndSessionType)
		switch st.SessionType {
		case constant.SingleChatType:
			lc.UserID = st.SourceID
			lc.ConversationID = utils.GetConversationIDBySessionType(st.SourceID, constant.SingleChatType)
			lc.ConversationType = constant.SingleChatType
		case constant.GroupChatType:
			conversationID, conversationType, err := c.getConversationTypeByGroupID(ctx, st.SourceID)
			if err != nil {
				log.Error("internal", "getConversationTypeByGroupID database err:", err.Error())
				return
			}
			lc.GroupID = st.SourceID
			lc.ConversationID = conversationID
			lc.ConversationType = conversationType
		}
		c.addFaceURLAndName(ctx, &lc)
		err := c.db.UpdateConversation(ctx, &lc)
		if err != nil {
			log.Error("internal", "setConversationFaceUrlAndNickName database err:", err.Error())
			return
		}
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: lc.ConversationID, Action: constant.ConChange, Args: []string{lc.ConversationID}}})

	case constant.UpdateLatestMessageChange:
		conversationID := node.ConID
		var latestMsg sdk_struct.MsgStruct
		l, err := c.db.GetConversation(ctx, conversationID)
		if err != nil {
			log.Error("internal", "getConversationLatestMsgModel err", err.Error())
		} else {
			err := json.Unmarshal([]byte(l.LatestMsg), &latestMsg)
			if err != nil {
				log.Error("internal", "latestMsg,Unmarshal err :", err.Error())
			} else {
				latestMsg.IsRead = true
				newLatestMessage := utils.StructToJsonString(latestMsg)
				err = c.db.UpdateColumnsConversation(nil, node.ConID, map[string]interface{}{"latest_msg_send_time": latestMsg.SendTime, "latest_msg": newLatestMessage})
				if err != nil {
					log.Error("internal", "updateConversationLatestMsgModel err :", err.Error())
				}
			}
		}
	case constant.ConChange:
		cidList := node.Args.([]string)
		cLists, err := c.db.GetMultipleConversationDB(ctx, cidList)
		if err != nil {
			log.Error("internal", "getMultipleConversationModel err :", err.Error())
		} else {
			var newCList []*model_struct.LocalConversation
			for _, v := range cLists {
				if v.LatestMsgSendTime != 0 {
					newCList = append(newCList, v)
				}
			}
			log.Info("internal", "getMultipleConversationModel success :", newCList)

			c.ConversationListener.OnConversationChanged(utils.StructToJsonStringDefault(newCList))
		}
	case constant.NewCon:
		cidList := node.Args.([]string)
		cLists, err := c.db.GetMultipleConversationDB(ctx, cidList)
		if err != nil {
			log.Error("internal", "getMultipleConversationModel err :", err.Error())
		} else {
			if cLists != nil {
				log.Info("internal", "getMultipleConversationModel success :", cLists)
				c.ConversationListener.OnNewConversation(utils.StructToJsonString(cLists))
			}
		}
	case constant.ConChangeDirect:
		cidList := node.Args.(string)
		c.ConversationListener.OnConversationChanged(cidList)

	case constant.NewConDirect:
		cidList := node.Args.(string)
		log.Debug("internal", "NewConversation", cidList)
		c.ConversationListener.OnNewConversation(cidList)

	case constant.ConversationLatestMsgHasRead:
		hasReadMsgList := node.Args.(map[string][]string)
		var result []*model_struct.LocalConversation
		var latestMsg sdk_struct.MsgStruct
		var lc model_struct.LocalConversation
		for conversationID, msgIDList := range hasReadMsgList {
			LocalConversation, err := c.db.GetConversation(ctx, conversationID)
			if err != nil {
				log.Error("internal", "get conversation err", err.Error(), conversationID)
				continue
			}
			err = utils.JsonStringToStruct(LocalConversation.LatestMsg, &latestMsg)
			if err != nil {
				log.Error("internal", "JsonStringToStruct err", err.Error(), conversationID)
				continue
			}
			if utils.IsContain(latestMsg.ClientMsgID, msgIDList) {
				latestMsg.IsRead = true
				lc.ConversationID = conversationID
				lc.LatestMsg = utils.StructToJsonString(latestMsg)
				LocalConversation.LatestMsg = utils.StructToJsonString(latestMsg)
				err := c.db.UpdateConversation(ctx, &lc)
				if err != nil {
					log.Error("internal", "UpdateConversation database err:", err.Error())
					continue
				} else {
					result = append(result, LocalConversation)
				}
			}
		}
		if result != nil {
			log.Info("internal", "getMultipleConversationModel success :", result)
			c.ConversationListener.OnNewConversation(utils.StructToJsonString(result))
		}
	case constant.SyncConversation:
		operationID := node.Args.(string)
		log.Debug(operationID, "reconn sync conversation start")
		c.SyncConversations(ctx)
		err := c.SyncConversationUnreadCount(ctx)
		if err != nil {
			log.Error(operationID, "reconn sync conversation unread count err", err.Error())
		}
		totalUnreadCount, err := c.db.GetTotalUnreadMsgCountDB(ctx)
		if err != nil {
			log.Error("internal", "TotalUnreadMessageChanged database err:", err.Error())
		} else {
			c.ConversationListener.OnTotalUnreadMessageCountChanged(totalUnreadCount)
		}

	}
}

func (c *Conversation) doUpdateMessage(c2v common.Cmd2Value) {
	if c.ConversationListener == nil {
		log.Error("internal", "not set conversationListener")
		return
	}

	node := c2v.Value.(common.UpdateMessageNode)
	ctx := c2v.Ctx
	switch node.Action {
	case constant.UpdateMsgFaceUrlAndNickName:
		args := node.Args.(common.UpdateMessageInfo)
		var conversationType int32
		if args.GroupID == "" {
			conversationType = constant.SingleChatType
		} else {
			var err error
			_, conversationType, err = c.getConversationTypeByGroupID(ctx, args.GroupID)
			if err != nil {
				log.Error("internal", "getConversationTypeByGroupID database err:", err.Error())
				return
			}
		}
		err := c.db.UpdateMsgSenderFaceURLAndSenderNicknameController(ctx, args.UserID, args.FaceURL, args.Nickname, int(conversationType), args.GroupID)
		if err != nil {
			log.Error("internal", "UpdateMsgSenderFaceURLAndSenderNickname err:", err.Error())
		}

	}

}

func (c *Conversation) doSyncReactionExtensions(c2v common.Cmd2Value) {
	if c.ConversationListener == nil {
		log.Error("internal", "not set conversationListener")
		return
	}
	node := c2v.Value.(common.SyncReactionExtensionsNode)
	ctx := mcontext.NewCtx(node.OperationID)
	switch node.Action {
	case constant.SyncMessageListReactionExtensions:
		args := node.Args.(syncReactionExtensionParams)
		log.Error(node.OperationID, "come SyncMessageListReactionExtensions", args)
		var reqList []server_api_params.OperateMessageListReactionExtensionsReq
		for _, v := range args.MessageList {
			var temp server_api_params.OperateMessageListReactionExtensionsReq
			temp.ClientMsgID = v.ClientMsgID
			temp.MsgFirstModifyTime = v.MsgFirstModifyTime
			reqList = append(reqList, temp)
		}
		var apiReq server_api_params.GetMessageListReactionExtensionsReq
		apiReq.SourceID = args.SourceID
		apiReq.TypeKeyList = args.TypeKeyList
		apiReq.SessionType = args.SessionType
		apiReq.MessageReactionKeyList = reqList
		apiReq.IsExternalExtensions = args.IsExternalExtension
		apiReq.OperationID = node.OperationID
		var apiResp server_api_params.GetMessageListReactionExtensionsResp
		err := c.p.PostReturn(constant.GetMessageListReactionExtensionsRouter, apiReq, &apiResp)
		if err != nil {
			log.NewError(node.OperationID, utils.GetSelfFuncName(), "getMessageListReactionExtensions err:", err.Error())
			return
		}
		// for _, result := range apiResp {
		// 	log.Warn(node.OperationID, "api return reslut is:", result.ClientMsgID, result.ReactionExtensionList)
		// }
		onLocal := func(data []*model_struct.LocalChatLogReactionExtensions) []*server_api_params.SingleMessageExtensionResult {
			var result []*server_api_params.SingleMessageExtensionResult
			for _, v := range data {
				temp := new(server_api_params.SingleMessageExtensionResult)
				tempMap := make(map[string]*sdkws.KeyValue)
				_ = json.Unmarshal(v.LocalReactionExtensions, &tempMap)
				if len(args.TypeKeyList) != 0 {
					for s, _ := range tempMap {
						if !utils.IsContain(s, args.TypeKeyList) {
							delete(tempMap, s)
						}
					}
				}

				temp.ReactionExtensionList = tempMap
				temp.ClientMsgID = v.ClientMsgID
				result = append(result, temp)
			}
			return result
		}(args.ExtendMessageList)
		var onServer []*server_api_params.SingleMessageExtensionResult
		for _, v := range apiResp {
			if v.ErrCode == 0 {
				onServer = append(onServer, v)
			}
		}
		aInBNot, _, sameA, _ := common.CheckReactionExtensionsDiff(onServer, onLocal)
		for _, v := range aInBNot {
			log.Error(node.OperationID, "come InsertMessageReactionExtension", args, v.ClientMsgID)
			if len(v.ReactionExtensionList) > 0 {
				temp := model_struct.LocalChatLogReactionExtensions{ClientMsgID: v.ClientMsgID, LocalReactionExtensions: []byte(utils.StructToJsonString(v.ReactionExtensionList))}
				err := c.db.InsertMessageReactionExtension(ctx, &temp)
				if err != nil {
					log.Error(node.OperationID, "InsertMessageReactionExtension err:", err.Error())
					continue
				}
			}
			var changedKv []*sdkws.KeyValue
			for _, value := range v.ReactionExtensionList {
				changedKv = append(changedKv, value)
			}
			if len(changedKv) > 0 {
				c.msgListener.OnRecvMessageExtensionsChanged(v.ClientMsgID, utils.StructToJsonString(changedKv))
			}
		}
		// for _, result := range sameA {
		// log.ZWarn(ctx, "result", result.ReactionExtensionList, result.ClientMsgID)
		// }
		for _, v := range sameA {
			log.Error(node.OperationID, "come sameA", v.ClientMsgID, v.ReactionExtensionList)
			tempMap := make(map[string]*sdkws.KeyValue)
			for _, extensions := range args.ExtendMessageList {
				if v.ClientMsgID == extensions.ClientMsgID {
					_ = json.Unmarshal(extensions.LocalReactionExtensions, &tempMap)
					break
				}
			}
			if len(v.ReactionExtensionList) == 0 {
				err := c.db.DeleteMessageReactionExtension(ctx, v.ClientMsgID)
				if err != nil {
					log.Error(node.OperationID, "DeleteMessageReactionExtension err:", err.Error())
					continue
				}
				var deleteKeyList []string
				for key, _ := range tempMap {
					deleteKeyList = append(deleteKeyList, key)
				}
				if len(deleteKeyList) > 0 {
					c.msgListener.OnRecvMessageExtensionsDeleted(v.ClientMsgID, utils.StructToJsonString(deleteKeyList))
				}
			} else {
				deleteKeyList, changedKv := func(local, server map[string]*sdkws.KeyValue) ([]string, []*sdkws.KeyValue) {
					var deleteKeyList []string
					var changedKv []*sdkws.KeyValue
					for k, v := range local {
						ia, ok := server[k]
						if ok {
							//服务器不同的kv
							if ia.Value != v.Value {
								changedKv = append(changedKv, ia)
							}
						} else {
							//服务器已经没有kv
							deleteKeyList = append(deleteKeyList, k)
						}
					}
					//从服务器新增的kv
					for k, v := range server {
						_, ok := local[k]
						if !ok {
							changedKv = append(changedKv, v)

						}
					}
					return deleteKeyList, changedKv
				}(tempMap, v.ReactionExtensionList)
				extendMsg := model_struct.LocalChatLogReactionExtensions{ClientMsgID: v.ClientMsgID, LocalReactionExtensions: []byte(utils.StructToJsonString(v.ReactionExtensionList))}
				err = c.db.UpdateMessageReactionExtension(ctx, &extendMsg)
				if err != nil {
					log.Error(node.OperationID, "UpdateMessageReactionExtension err:", err.Error())
					continue
				}
				if len(deleteKeyList) > 0 {
					c.msgListener.OnRecvMessageExtensionsDeleted(v.ClientMsgID, utils.StructToJsonString(deleteKeyList))
				}
				if len(changedKv) > 0 {
					c.msgListener.OnRecvMessageExtensionsChanged(v.ClientMsgID, utils.StructToJsonString(changedKv))
				}
			}
			//err := c.db.GetAndUpdateMessageReactionExtension(v.ClientMsgID, v.ReactionExtensionList)
			//if err != nil {
			//	log.Error(node.OperationID, "GetAndUpdateMessageReactionExtension err:", err.Error())
			//	continue
			//}
			//var changedKv []*server_api_params.KeyValue
			//for _, value := range v.ReactionExtensionList {
			//	changedKv = append(changedKv, value)
			//}
			//if len(changedKv) > 0 {
			//	c.msgListener.OnRecvMessageExtensionsChanged(v.ClientMsgID, utils.StructToJsonString(changedKv))
			//}
		}
	case constant.SyncMessageListTypeKeyInfo:
		messageList := node.Args.([]*sdk_struct.MsgStruct)
		var sourceID string
		var sessionType int32
		var reqList []server_api_params.OperateMessageListReactionExtensionsReq
		var temp server_api_params.OperateMessageListReactionExtensionsReq
		for _, v := range messageList {
			message, err := c.db.GetMessageController(ctx, v)
			if err != nil {
				log.Error(node.OperationID, "GetMessageController err:", err.Error(), *v)
				continue
			}
			temp.ClientMsgID = message.ClientMsgID
			temp.MsgFirstModifyTime = message.MsgFirstModifyTime
			reqList = append(reqList, temp)
			switch message.SessionType {
			case constant.SingleChatType:
				sourceID = message.SendID + message.RecvID
			case constant.NotificationChatType:
				sourceID = message.RecvID
			case constant.GroupChatType, constant.SuperGroupChatType:
				sourceID = message.RecvID
			}
			sessionType = message.SessionType
		}
		var apiReq server_api_params.GetMessageListReactionExtensionsReq
		apiReq.SourceID = sourceID
		apiReq.SessionType = sessionType
		apiReq.MessageReactionKeyList = reqList
		apiReq.OperationID = node.OperationID
		var apiResp server_api_params.GetMessageListReactionExtensionsResp
		err := c.p.PostReturnWithTimeOut(constant.GetMessageListReactionExtensionsRouter, apiReq, &apiResp, time.Second*2)
		if err != nil {
			log.Error(node.OperationID, "GetMessageListReactionExtensions from server err:", err.Error(), apiReq)
			return
		}
		var messageChangedList []*messageKvList
		for _, v := range apiResp {
			if v.ErrCode == 0 {
				var changedKv []*sdkws.KeyValue
				var prefixTypeKey []string
				extendMsg, _ := c.db.GetMessageReactionExtension(ctx, v.ClientMsgID)
				localKV := make(map[string]*sdkws.KeyValue)
				_ = json.Unmarshal(extendMsg.LocalReactionExtensions, &localKV)
				for typeKey, value := range v.ReactionExtensionList {
					oldValue, ok := localKV[typeKey]
					if ok {
						if !cmp.Equal(value, oldValue) {
							localKV[typeKey] = value
							prefixTypeKey = append(prefixTypeKey, getPrefixTypeKey(typeKey))
							changedKv = append(changedKv, value)
						}
					} else {
						localKV[typeKey] = value
						prefixTypeKey = append(prefixTypeKey, getPrefixTypeKey(typeKey))
						changedKv = append(changedKv, value)

					}

				}
				extendMsg.LocalReactionExtensions = []byte(utils.StructToJsonString(localKV))
				_ = c.db.UpdateMessageReactionExtension(ctx, extendMsg)
				if len(changedKv) > 0 {
					c.msgListener.OnRecvMessageExtensionsChanged(extendMsg.ClientMsgID, utils.StructToJsonString(changedKv))
				}
				prefixTypeKey = utils.RemoveRepeatedStringInList(prefixTypeKey)
				if len(prefixTypeKey) > 0 && c.msgKvListener != nil {
					var result []*sdk.SingleTypeKeyInfoSum
					oneMessageChanged := new(messageKvList)
					oneMessageChanged.ClientMsgID = extendMsg.ClientMsgID
					for _, v := range prefixTypeKey {
						singleResult := new(sdk.SingleTypeKeyInfoSum)
						singleResult.TypeKey = v
						for typeKey, value := range localKV {
							if strings.HasPrefix(typeKey, v) {
								singleTypeKeyInfo := new(sdk.SingleTypeKeyInfo)
								err := json.Unmarshal([]byte(value.Value), singleTypeKeyInfo)
								if err != nil {
									continue
								}
								if _, ok := singleTypeKeyInfo.InfoList[c.loginUserID]; ok {
									singleResult.IsContainSelf = true
								}
								for _, info := range singleTypeKeyInfo.InfoList {
									v := *info
									singleResult.InfoList = append(singleResult.InfoList, &v)
								}
								singleResult.Counter += singleTypeKeyInfo.Counter
							}
						}
						result = append(result, singleResult)
					}
					oneMessageChanged.ChangedKvList = result
					messageChangedList = append(messageChangedList, oneMessageChanged)
				}
			}
		}
		if len(messageChangedList) > 0 && c.msgKvListener != nil {
			c.msgKvListener.OnMessageKvInfoChanged(utils.StructToJsonString(messageChangedList))
		}

	}

}

func (c *Conversation) DoNotification(ctx context.Context, msg *sdkws.MsgData) {
	if msg.SendTime < c.full.Group().LoginTime() || c.full.Group().LoginTime() == 0 {
		log.ZWarn(ctx, "ignore notification", nil, "clientMsgID", msg.ClientMsgID, "serverMsgID", msg.ServerMsgID, "seq", msg.Seq, "contentType", msg.ContentType)
		return
	}
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg)
	if c.msgListener == nil {
		log.Error(operationID, utils.GetSelfFuncName(), "listener == nil")
		return
	}
	go func() {
		c.SyncConversations(ctx)
	}()
}

func (c *Conversation) doMsgNew1(c2v common.Cmd2Value) {
	operationID := c2v.Value.(sdk_struct.CmdNewMsgComeToConversation).OperationID
	allMsg := c2v.Value.(sdk_struct.CmdNewMsgComeToConversation).MsgList
	syncFlag := c2v.Value.(sdk_struct.CmdNewMsgComeToConversation).SyncFlag
	ctx := mcontext.NewCtx(operationID)
	if c.msgListener == nil || c.ConversationListener == nil {
		for _, v := range allMsg {
			if v.ContentType > constant.SignalingNotificationBegin && v.ContentType < constant.SignalingNotificationEnd {
				c.signaling.DoNotification(ctx, v, c.GetCh())
			}
		}
		return
	}
	if syncFlag == constant.MsgSyncBegin {
		c.ConversationListener.OnSyncServerStart()
	}
	if syncFlag == constant.MsgSyncFailed {
		c.ConversationListener.OnSyncServerFailed()
	}

	var isTriggerUnReadCount bool
	var insertMsg, updateMsg []*model_struct.LocalChatLog
	var exceptionMsg []*model_struct.LocalErrChatLog
	var unreadMessages []*model_struct.LocalConversationUnreadMessage
	var newMessages, msgReadList, groupMsgReadList, msgRevokeList, newMsgRevokeList, reactionMsgModifierList, reactionMsgDeleterList sdk_struct.NewMsgList
	var isUnreadCount, isConversationUpdate, isHistory, isNotPrivate, isSenderConversationUpdate, isSenderNotificationPush bool
	conversationChangedSet := make(map[string]*model_struct.LocalConversation)
	newConversationSet := make(map[string]*model_struct.LocalConversation)
	conversationSet := make(map[string]*model_struct.LocalConversation)
	phConversationChangedSet := make(map[string]*model_struct.LocalConversation)
	phNewConversationSet := make(map[string]*model_struct.LocalConversation)
	log.ZDebug(ctx, "do Msg come here", "len", len(allMsg), "ch len", len(c.GetCh()))
	b := time.Now()
	for _, v := range allMsg {
		log.ZDebug(ctx, "do Msg come here", "loginUserID", c.loginUserID, "msg", v)
		isHistory = utils.GetSwitchFromOptions(v.Options, constant.IsHistory)
		isUnreadCount = utils.GetSwitchFromOptions(v.Options, constant.IsUnreadCount)
		isConversationUpdate = utils.GetSwitchFromOptions(v.Options, constant.IsConversationUpdate)
		isNotPrivate = utils.GetSwitchFromOptions(v.Options, constant.IsNotPrivate)
		isSenderConversationUpdate = utils.GetSwitchFromOptions(v.Options, constant.IsSenderConversationUpdate)
		isSenderNotificationPush = utils.GetSwitchFromOptions(v.Options, constant.IsSenderNotificationPush)
		msg := new(sdk_struct.MsgStruct)
		copier.Copy(msg, v)
		if v.OfflinePushInfo != nil {
			msg.OfflinePush = *v.OfflinePushInfo
		}
		msg.Content = string(v.Content)
		//var tips sdkws.TipsComm
		//if v.ContentType >= constant.NotificationBegin && v.ContentType <= constant.NotificationEnd {
		//	_ = proto.Unmarshal(v.Content, &tips)
		//	marshaler := jsonpb.Marshaler{
		//		OrigName:     true,
		//		EnumsAsInts:  false,
		//		EmitDefaults: false,
		//	}
		//	msg.Content, _ = marshaler.MarshalToString(&tips)
		//} else {
		//	msg.Content = string(v.Content)
		//}
		//When the message has been marked and deleted by the cloud, it is directly inserted locally without any conversation and message update.
		if msg.Status == constant.MsgStatusHasDeleted {
			insertMsg = append(insertMsg, c.msgStructToLocalChatLog(msg))
			continue
		}
		msg.Status = constant.MsgStatusSendSuccess
		msg.IsRead = false
		//De-analyze data
		if err := c.msgHandleByContentType(msg); err != nil {
			log.Error(operationID, "Parsing data error:", err.Error(), *msg, "type: ", msg.ContentType)
			continue
		}
		if !isSenderNotificationPush {
			msg.AttachedInfoElem.NotSenderNotificationPush = true
			msg.AttachedInfo = utils.StructToJsonString(msg.AttachedInfoElem)
		}
		if !isNotPrivate {
			msg.AttachedInfoElem.IsPrivateChat = true
			msg.AttachedInfo = utils.StructToJsonString(msg.AttachedInfoElem)
		}
		if msg.ClientMsgID == "" {
			exceptionMsg = append(exceptionMsg, c.msgStructToLocalErrChatLog(msg))
			continue
		}
		ctx := context.Background()
		mcontext.SetOperationID(ctx, operationID)
		switch {
		case v.ContentType == constant.ConversationChangeNotification || v.ContentType == constant.ConversationPrivateChatNotification:
			c.DoNotification(ctx, v)
		case v.ContentType == constant.MsgDeleteNotification:
			c.full.SuperGroup.DoNotification(v, c.GetCh(), operationID)
		case v.ContentType == constant.SuperGroupUpdateNotification:
			c.full.SuperGroup.DoNotification(v, c.GetCh(), operationID)
			continue
		case v.ContentType == constant.ConversationUnreadNotification:
			var unreadArgs sdkws.ConversationUpdateTips
			_ = json.Unmarshal([]byte(msg.Content), &unreadArgs)
			for _, v := range unreadArgs.ConversationIDList {
				c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: v, Action: constant.UnreadCountSetZero}})
				c.db.DeleteConversationUnreadMessageList(ctx, v, unreadArgs.UpdateUnreadCountTime)
			}
			c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.ConChange, Args: unreadArgs.ConversationIDList}})
			continue
		case v.ContentType == constant.BusinessNotification:
			c.business.DoNotification(ctx, msg.Content)
			continue
		}

		switch v.SessionType {
		case constant.SingleChatType:
			if v.ContentType > constant.FriendNotificationBegin && v.ContentType < constant.FriendNotificationEnd {
				c.friend.DoNotification(ctx, v)
			} else if v.ContentType > constant.UserNotificationBegin && v.ContentType < constant.UserNotificationEnd {
				c.user.DoNotification(ctx, v)
			} else if utils2.Contain(v.ContentType, constant.GroupApplicationRejectedNotification, constant.GroupApplicationAcceptedNotification, constant.JoinGroupApplicationNotification) {
				c.group.DoNotification(ctx, v)
			} else if v.ContentType > constant.SignalingNotificationBegin && v.ContentType < constant.SignalingNotificationEnd {
				c.signaling.DoNotification(ctx, v, c.GetCh())
				continue
			} else if v.ContentType == constant.WorkMomentNotification {
				c.workMoments.DoNotification(msg.Content, operationID)
			}
		case constant.GroupChatType, constant.SuperGroupChatType:
			if v.ContentType > constant.GroupNotificationBegin && v.ContentType < constant.GroupNotificationEnd {
				c.group.DoNotification(ctx, v)
			} else if v.ContentType > constant.SignalingNotificationBegin && v.ContentType < constant.SignalingNotificationEnd {
				c.signaling.DoNotification(ctx, v, c.GetCh())
				continue
			}
		}
		if v.SendID == c.loginUserID { //seq
			// Messages sent by myself  //if  sent through  this terminal
			m, err := c.db.GetMessageController(ctx, msg)
			if err == nil {
				log.Info(operationID, "have message", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, *msg)
				if m.Seq == 0 {
					if !isConversationUpdate {
						msg.Status = constant.MsgStatusFiltered
					}
					updateMsg = append(updateMsg, c.msgStructToLocalChatLog(msg))
				} else {
					exceptionMsg = append(exceptionMsg, c.msgStructToLocalErrChatLog(msg))
				}
			} else {
				log.Info(operationID, "sync message", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, *msg)
				lc := model_struct.LocalConversation{
					ConversationType:  v.SessionType,
					LatestMsg:         utils.StructToJsonString(msg),
					LatestMsgSendTime: msg.SendTime,
				}
				switch v.SessionType {
				case constant.SingleChatType:
					lc.ConversationID = utils.GetConversationIDBySessionType(v.RecvID, constant.SingleChatType)
					lc.UserID = v.RecvID
				case constant.GroupChatType:
					lc.GroupID = v.GroupID
					lc.ConversationID = utils.GetConversationIDBySessionType(lc.GroupID, constant.GroupChatType)
				case constant.SuperGroupChatType:
					lc.GroupID = v.GroupID
					lc.ConversationID = utils.GetConversationIDBySessionType(lc.GroupID, constant.SuperGroupChatType)
				}
				if isConversationUpdate {
					if isSenderConversationUpdate {
						log.Debug(operationID, "updateConversation msg", v, lc)
						c.updateConversation(&lc, conversationSet)
					}
					newMessages = append(newMessages, msg)
				} else {
					msg.Status = constant.MsgStatusFiltered
				}
				if isHistory {
					insertMsg = append(insertMsg, c.msgStructToLocalChatLog(msg))
				}
				switch msg.ContentType {
				case constant.Revoke:
					msgRevokeList = append(msgRevokeList, msg)
				case constant.HasReadReceipt:
					msgReadList = append(msgReadList, msg)
				case constant.GroupHasReadReceipt:
					groupMsgReadList = append(groupMsgReadList, msg)
				case constant.AdvancedRevoke:
					newMsgRevokeList = append(newMsgRevokeList, msg)
					newMessages = removeElementInList(newMessages, msg)
				case constant.ReactionMessageModifier:
					reactionMsgModifierList = append(reactionMsgModifierList, msg)
				case constant.ReactionMessageDeleter:
					reactionMsgDeleterList = append(reactionMsgDeleterList, msg)
				default:
				}
			}
		} else { //Sent by others
			if _, err := c.db.GetMessageController(ctx, msg); err != nil { //Deduplication operation
				lc := model_struct.LocalConversation{
					ConversationType:  v.SessionType,
					LatestMsg:         utils.StructToJsonString(msg),
					LatestMsgSendTime: msg.SendTime,
				}
				switch v.SessionType {
				case constant.SingleChatType:
					lc.ConversationID = utils.GetConversationIDBySessionType(v.SendID, constant.SingleChatType)
					lc.UserID = v.SendID
					lc.ShowName = msg.SenderNickname
					lc.FaceURL = msg.SenderFaceURL
				case constant.GroupChatType:
					lc.GroupID = v.GroupID
					lc.ConversationID = utils.GetConversationIDBySessionType(lc.GroupID, constant.GroupChatType)
				case constant.SuperGroupChatType:
					lc.GroupID = v.GroupID
					lc.ConversationID = utils.GetConversationIDBySessionType(lc.GroupID, constant.SuperGroupChatType)
					//faceUrl, name, err := u.getGroupNameAndFaceUrlByUid(c.GroupID)
					//if err != nil {
					//	utils.sdkLog("getGroupNameAndFaceUrlByUid err:", err)
					//} else {
					//	c.ShowName = name
					//	c.FaceURL = faceUrl
					//}
				case constant.NotificationChatType:
					lc.ConversationID = utils.GetConversationIDBySessionType(v.SendID, constant.NotificationChatType)
					lc.UserID = v.SendID
				}
				if isUnreadCount {
					cacheConversation := c.cache.GetConversation(lc.ConversationID)
					if msg.SendTime > cacheConversation.UpdateUnreadCountTime {
						isTriggerUnReadCount = true
						lc.UnreadCount = 1
						tempUnreadMessages := model_struct.LocalConversationUnreadMessage{ConversationID: lc.ConversationID, ClientMsgID: msg.ClientMsgID, SendTime: msg.SendTime}
						unreadMessages = append(unreadMessages, &tempUnreadMessages)
					}
				}
				if isConversationUpdate {
					c.updateConversation(&lc, conversationSet)
					newMessages = append(newMessages, msg)
				} else {
					msg.Status = constant.MsgStatusFiltered
				}
				if isHistory {
					log.Debug(operationID, "trigger msg is ", msg.SenderNickname, msg.SenderFaceURL)
					insertMsg = append(insertMsg, c.msgStructToLocalChatLog(msg))
				}
				switch msg.ContentType {
				case constant.Revoke:
					msgRevokeList = append(msgRevokeList, msg)
				case constant.HasReadReceipt:
					msgReadList = append(msgReadList, msg)
				case constant.GroupHasReadReceipt:
					groupMsgReadList = append(groupMsgReadList, msg)
				case constant.Typing:
					newMessages = append(newMessages, msg)
				case constant.CustomMsgOnlineOnly:
					newMessages = append(newMessages, msg)
				case constant.CustomMsgNotTriggerConversation:
					newMessages = append(newMessages, msg)
				case constant.OANotification:
					if !isConversationUpdate {
						newMessages = append(newMessages, msg)
					}
				case constant.AdvancedRevoke:
					newMsgRevokeList = append(newMsgRevokeList, msg)
					newMessages = removeElementInList(newMessages, msg)
				case constant.ReactionMessageModifier:
					reactionMsgModifierList = append(reactionMsgModifierList, msg)
				case constant.ReactionMessageDeleter:
					reactionMsgDeleterList = append(reactionMsgDeleterList, msg)
				default:
				}

			} else {
				exceptionMsg = append(exceptionMsg, c.msgStructToLocalErrChatLog(msg))
				log.ZWarn(ctx, "Deduplication operation ", nil, "msg", *c.msgStructToLocalErrChatLog(msg))
			}
		}
	}
	b1 := utils.GetCurrentTimestampByMill()
	log.Info(operationID, "generate conversation map is :", conversationSet)
	log.Debug(operationID, "before insert msg cost time : ", time.Since(b))

	list, err := c.db.GetAllConversationListDB(ctx)
	if err != nil {
		log.Error(operationID, "GetAllConversationListDB", "error", err.Error())
	}
	m := make(map[string]*model_struct.LocalConversation)
	listToMap(list, m)
	log.Debug(operationID, "listToMap: ", list, conversationSet)
	c.diff(ctx, m, conversationSet, conversationChangedSet, newConversationSet)
	log.Info(operationID, "trigger map is :", "newConversations", newConversationSet, "changedConversations", conversationChangedSet)
	b2 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "listToMap diff, cost time : ", b2-b1)

	//seq sync message update
	err5 := c.db.BatchUpdateMessageList(ctx, updateMsg)
	if err5 != nil {
		log.Error(operationID, "sync seq normal message err  :", err5.Error())
	}
	b3 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "BatchUpdateMessageList, cost time : ", b3-b2)

	//Normal message storage
	err1 := c.db.BatchInsertMessageListController(ctx, insertMsg)
	if err1 != nil {
		log.Error(operationID, "insert GetMessage detail err:", err1.Error(), len(insertMsg))
		for _, v := range insertMsg {
			e := c.db.InsertMessageController(ctx, v)
			if e != nil {
				errChatLog := &model_struct.LocalErrChatLog{}
				copier.Copy(errChatLog, v)
				exceptionMsg = append(exceptionMsg, errChatLog)
				log.ZWarn(ctx, "InsertMessage operation", err, "chatErrLog", errChatLog, "chatLog", v)
			}
		}
	}
	b4 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "BatchInsertMessageListController, cost time : ", b4-b3)

	//Exception message storage
	log.ZWarn(ctx, "exceptionMsgs", nil, "msgs", exceptionMsg)
	err2 := c.db.BatchInsertExceptionMsgController(ctx, exceptionMsg)
	if err2 != nil {
		log.Error(operationID, "insert err message err  :", err2.Error())

	}
	hList, _ := c.db.GetHiddenConversationList(ctx)
	for _, v := range hList {
		if nc, ok := newConversationSet[v.ConversationID]; ok {
			phConversationChangedSet[v.ConversationID] = nc
			nc.RecvMsgOpt = v.RecvMsgOpt
			nc.GroupAtType = v.GroupAtType
			nc.IsPinned = v.IsPinned
			nc.IsPrivateChat = v.IsPrivateChat
			if nc.IsPrivateChat {
				nc.BurnDuration = v.BurnDuration
			}
			nc.IsNotInGroup = v.IsNotInGroup
			nc.AttachedInfo = v.AttachedInfo
			nc.Ex = v.Ex
		}
	}
	b5 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "GetHiddenConversationList, cost time : ", b5-b4)

	for k, v := range newConversationSet {
		if _, ok := phConversationChangedSet[v.ConversationID]; !ok {
			phNewConversationSet[k] = v
		}
	}
	//Changed conversation storage
	err3 := c.db.BatchUpdateConversationList(ctx, append(mapConversationToList(conversationChangedSet), mapConversationToList(phConversationChangedSet)...))
	if err3 != nil {
		log.Error(operationID, "insert changed conversation err :", err3.Error())
	}
	b6 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "BatchUpdateConversationList, cost time : ", b6-b5)
	//New conversation storage
	err4 := c.db.BatchInsertConversationList(ctx, mapConversationToList(phNewConversationSet))
	if err4 != nil {
		log.Error(operationID, "insert new conversation err:", err4.Error())
	}
	b7 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "BatchInsertConversationList, cost time : ", b7-b6)
	unreadMessageErr := c.db.BatchInsertConversationUnreadMessageList(ctx, unreadMessages)
	if unreadMessageErr != nil {
		log.Error(operationID, "insert BatchInsertConversationUnreadMessageList err:", unreadMessageErr.Error())
	}
	c.doMsgReadState(ctx, msgReadList)
	b8 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "doMsgReadState  cost time : ", b8-b7)

	c.DoGroupMsgReadState(ctx, groupMsgReadList)
	b9 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "DoGroupMsgReadState  cost time : ", b9-b8, "len: ", len(groupMsgReadList))

	c.revokeMessage(ctx, msgRevokeList)
	b10 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "revokeMessage  cost time : ", b10-b9)
	if c.batchMsgListener != nil {
		c.batchNewMessages(ctx, newMessages)
		b11 := utils.GetCurrentTimestampByMill()
		log.Debug(operationID, "batchNewMessages  cost time : ", b11-b10)
	} else {
		c.newMessage(newMessages)
		b12 := utils.GetCurrentTimestampByMill()
		log.Debug(operationID, "newMessage  cost time : ", b12-b10)
	}
	c.newRevokeMessage(ctx, newMsgRevokeList)
	c.doReactionMsgModifier(ctx, reactionMsgModifierList)
	c.doReactionMsgDeleter(ctx, reactionMsgDeleterList)
	//log.Info(operationID, "trigger map is :", newConversationSet, conversationChangedSet)
	if len(newConversationSet) > 0 {
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.NewConDirect, Args: utils.StructToJsonString(mapConversationToList(newConversationSet))}})

	}
	if len(conversationChangedSet) > 0 {
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.ConChangeDirect, Args: utils.StructToJsonString(mapConversationToList(conversationChangedSet))}})
	}

	if isTriggerUnReadCount {
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.TotalUnreadMessageChanged, Args: ""}})
	}
	if syncFlag == constant.MsgSyncEnd {
		c.ConversationListener.OnSyncServerFinish()
	}
	log.Debug(operationID, "insert msg, total cost time: ", time.Since(b), "len:  ", len(allMsg))
}
