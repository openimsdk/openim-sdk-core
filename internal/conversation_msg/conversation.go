package conversation_msg

import (
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
	_ "open_im_sdk/internal/common"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"sort"
	"strings"
	"time"
)

func (c *Conversation) getAllConversationList(callback open_im_sdk_callback.Base, operationID string) sdk.GetAllConversationListCallback {
	conversationList, err := c.db.GetAllConversationListDB()
	common.CheckDBErrCallback(callback, err, operationID)
	return conversationList
}
func (c *Conversation) hideConversation(callback open_im_sdk_callback.Base, conversationID string, operationID string) {
	err := c.db.UpdateColumnsConversation(conversationID, map[string]interface{}{"latest_msg_send_time": 0})
	common.CheckDBErrCallback(callback, err, operationID)
}

func (c *Conversation) getConversationListSplit(callback open_im_sdk_callback.Base, offset, count int, operationID string) sdk.GetConversationListSplitCallback {
	conversationList, err := c.db.GetConversationListSplitDB(offset, count)
	common.CheckDBErrCallback(callback, err, operationID)
	return conversationList
}

func (c *Conversation) setConversationRecvMessageOpt(callback open_im_sdk_callback.Base, conversationIDList []string, opt int, operationID string) {
	apiReq := server_api_params.BatchSetConversationsReq{}
	apiResp := server_api_params.BatchSetConversationsResp{}
	apiReq.OperationID = operationID
	apiReq.OwnerUserID = c.loginUserID
	apiReq.NotificationType = constant.ConversationChangeNotification
	var conversations []server_api_params.Conversation
	for _, conversationID := range conversationIDList {
		localConversation, err := c.db.GetConversation(conversationID)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), "GetConversation failed", err.Error())
			continue
		}
		if localConversation.ConversationType == constant.SuperGroupChatType && opt == constant.NotReceiveMessage {
			common.CheckAnyErrCallback(callback, 100, errors.New("super group not support this opt"), operationID)
		}
		conversations = append(conversations, server_api_params.Conversation{
			OwnerUserID:      c.loginUserID,
			ConversationID:   conversationID,
			ConversationType: localConversation.ConversationType,
			UserID:           localConversation.UserID,
			GroupID:          localConversation.GroupID,
			RecvMsgOpt:       int32(opt),
			IsPinned:         localConversation.IsPinned,
			IsPrivateChat:    localConversation.IsPrivateChat,
			AttachedInfo:     localConversation.AttachedInfo,
			Ex:               localConversation.Ex,
		})
	}
	apiReq.Conversations = conversations
	c.p.PostFatalCallback(callback, constant.BatchSetConversationRouter, apiReq, &apiResp, apiReq.OperationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "output: ", apiResp)
	c.SyncConversations(operationID, 0)
}

func (c *Conversation) setConversation(callback open_im_sdk_callback.Base, apiReq *server_api_params.ModifyConversationFieldReq, conversationID string, localConversation *model_struct.LocalConversation, operationID string) {
	apiResp := server_api_params.ModifyConversationFieldResp{}
	apiReq.OwnerUserID = c.loginUserID
	apiReq.OperationID = operationID
	apiReq.ConversationID = conversationID
	apiReq.ConversationType = localConversation.ConversationType
	apiReq.UserID = localConversation.UserID
	apiReq.GroupID = localConversation.GroupID
	apiReq.UserIDList = []string{c.loginUserID}
	c.p.PostFatalCallback(callback, constant.ModifyConversationFieldRouter, apiReq, nil, apiReq.OperationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "request success, output: ", apiResp)
}

func (c *Conversation) setGlobalRecvMessageOpt(callback open_im_sdk_callback.Base, opt int32, operationID string) {
	apiReq := server_api_params.SetGlobalRecvMessageOptReq{}
	apiReq.OperationID = operationID
	apiReq.GlobalRecvMsgOpt = &opt
	c.p.PostFatalCallback(callback, constant.SetGlobalRecvMessageOptRouter, apiReq, nil, apiReq.OperationID)
	c.user.SyncLoginUserInfo(operationID)
}
func (c *Conversation) setOneConversationRecvMessageOpt(callback open_im_sdk_callback.Base, conversationID string, opt int, operationID string) {
	apiReq := &server_api_params.ModifyConversationFieldReq{}
	localConversation, err := c.db.GetConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	apiReq.RecvMsgOpt = int32(opt)
	apiReq.FieldType = constant.FieldRecvMsgOpt
	c.setConversation(callback, apiReq, conversationID, localConversation, operationID)
	c.SyncConversations(operationID, 0)
}
func (c *Conversation) setOneConversationUnread(callback open_im_sdk_callback.Base, conversationID string, unreadCount int, operationID string) {
	apiReq := &server_api_params.ModifyConversationFieldReq{}
	localConversation, err := c.db.GetConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	if localConversation.UnreadCount == 0 {
		return
	}
	apiReq.UpdateUnreadCountTime = localConversation.LatestMsgSendTime
	apiReq.UnreadCount = int32(unreadCount)
	apiReq.FieldType = constant.FieldUnread
	c.setConversation(callback, apiReq, conversationID, localConversation, operationID)
	deleteRows := c.db.DeleteConversationUnreadMessageList(localConversation.ConversationID, localConversation.LatestMsgSendTime)
	if deleteRows == 0 {
		log.Error(operationID, "DeleteConversationUnreadMessageList err", localConversation.ConversationID, localConversation.LatestMsgSendTime)
	}
}

func (c *Conversation) setOneConversationPrivateChat(callback open_im_sdk_callback.Base, conversationID string, isPrivate bool, operationID string) {
	apiReq := &server_api_params.ModifyConversationFieldReq{}
	localConversation, err := c.db.GetConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	apiReq.IsPrivateChat = isPrivate
	apiReq.FieldType = constant.FieldIsPrivateChat
	c.setConversation(callback, apiReq, conversationID, localConversation, operationID)
	c.SyncConversations(operationID, 0)
}

func (c *Conversation) setOneConversationBurnDuration(callback open_im_sdk_callback.Base, conversationID string, burnDuration int32, operationID string) {
	apiReq := &server_api_params.ModifyConversationFieldReq{}
	localConversation, err := c.db.GetConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	apiReq.BurnDuration = burnDuration
	apiReq.FieldType = constant.FieldBurnDuration
	c.setConversation(callback, apiReq, conversationID, localConversation, operationID)
	c.SyncConversations(operationID, 0)
}

func (c *Conversation) setOneConversationPinned(callback open_im_sdk_callback.Base, conversationID string, isPinned bool, operationID string) {
	apiReq := &server_api_params.ModifyConversationFieldReq{}
	localConversation, err := c.db.GetConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	apiReq.IsPinned = isPinned
	apiReq.FieldType = constant.FieldIsPinned
	c.setConversation(callback, apiReq, conversationID, localConversation, operationID)
	c.SyncConversations(operationID, 0)
}

func (c *Conversation) setOneConversationGroupAtType(callback open_im_sdk_callback.Base, conversationID, operationID string) {
	lc, err := c.db.GetConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	if lc.GroupAtType == constant.AtNormal || lc.ConversationType != constant.GroupChatType {
		common.CheckAnyErrCallback(callback, 201, errors.New("conversation don't need to reset"), operationID)
	}
	apiReq := &server_api_params.ModifyConversationFieldReq{}
	localConversation, err := c.db.GetConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	apiReq.GroupAtType = constant.AtNormal
	apiReq.FieldType = constant.FieldGroupAtType
	c.setConversation(callback, apiReq, conversationID, localConversation, operationID)
	c.SyncConversations(operationID, 0)
}
func (c *Conversation) getConversationRecvMessageOpt(callback open_im_sdk_callback.Base, conversationIDList []string, operationID string) []server_api_params.GetConversationRecvMessageOptResp {
	apiReq := server_api_params.GetConversationsReq{}
	apiReq.OperationID = operationID
	apiReq.OwnerUserID = c.loginUserID
	apiReq.ConversationIDs = conversationIDList
	var resp []server_api_params.GetConversationRecvMessageOptResp
	conversations := c.getMultipleConversation(callback, conversationIDList, operationID)
	for _, conversation := range conversations {
		resp = append(resp, server_api_params.GetConversationRecvMessageOptResp{
			ConversationID: conversation.ConversationID,
			Result:         &conversation.RecvMsgOpt,
		})
	}
	return resp
}

func (c *Conversation) getOneConversation(callback open_im_sdk_callback.Base, sourceID string, sessionType int32, operationID string) *model_struct.LocalConversation {
	conversationID := utils.GetConversationIDBySessionType(sourceID, int(sessionType))
	lc, err := c.db.GetConversation(conversationID)
	if err == nil {
		return lc
	} else {
		var newConversation model_struct.LocalConversation
		newConversation.ConversationID = conversationID
		newConversation.ConversationType = sessionType
		switch sessionType {
		case constant.SingleChatType:
			newConversation.UserID = sourceID
			faceUrl, name, err, isFromSvr := c.friend.GetUserNameAndFaceUrlByUid(sourceID, operationID)
			//	faceUrl, name, err := c.cache.GetUserNameAndFaceURL(sourceID, operationID)
			common.CheckDBErrCallback(callback, err, operationID)
			if isFromSvr {
				c.cache.Update(sourceID, faceUrl, name)
			}
			newConversation.ShowName = name
			newConversation.FaceURL = faceUrl
		case constant.GroupChatType, constant.SuperGroupChatType:
			newConversation.GroupID = sourceID
			g, err := c.full.GetGroupInfoFromLocal2Svr(sourceID, sessionType)
			//g, err := c.db.GetGroupInfoByGroupID(sourceID)
			common.CheckDBErrCallback(callback, err, operationID)
			newConversation.ShowName = g.GroupName
			newConversation.FaceURL = g.FaceURL
		}
		lc, errTemp := c.db.GetConversation(conversationID)
		if errTemp == nil {
			return lc
		}
		err := c.db.InsertConversation(&newConversation)
		common.CheckDBErrCallback(callback, err, operationID)
		return &newConversation
	}
}
func (c *Conversation) getMultipleConversation(callback open_im_sdk_callback.Base, conversationIDList []string, operationID string) sdk.GetMultipleConversationCallback {
	conversationList, err := c.db.GetMultipleConversationDB(conversationIDList)
	common.CheckDBErrCallback(callback, err, operationID)
	return conversationList
}

func (c *Conversation) deleteConversation(callback open_im_sdk_callback.Base, conversationID, operationID string) {
	lc, err := c.db.GetConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	var sourceID string
	switch lc.ConversationType {
	case constant.SingleChatType, constant.NotificationChatType:
		sourceID = lc.UserID
	case constant.GroupChatType, constant.SuperGroupChatType:
		sourceID = lc.GroupID
	}
	if lc.ConversationType == constant.SuperGroupChatType {
		err = c.db.SuperGroupDeleteAllMessage(lc.GroupID)
		common.CheckDBErrCallback(callback, err, operationID)
	} else {
		//Mark messages related to this conversation for deletion
		err = c.db.UpdateMessageStatusBySourceIDController(sourceID, constant.MsgStatusHasDeleted, lc.ConversationType)
		common.CheckDBErrCallback(callback, err, operationID)
	}
	//Reset the session information, empty session
	err = c.db.ResetConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{"", constant.TotalUnreadMessageChanged, ""}})

}
func (c *Conversation) setConversationDraft(callback open_im_sdk_callback.Base, conversationID, draftText, operationID string) {
	if draftText != "" {
		err := c.db.SetConversationDraft(conversationID, draftText)
		common.CheckDBErrCallback(callback, err, operationID)
	} else {
		err := c.db.RemoveConversationDraft(conversationID, draftText)
		common.CheckDBErrCallback(callback, err, operationID)
	}
}
func (c *Conversation) pinConversation(callback open_im_sdk_callback.Base, conversationID string, isPinned bool, operationID string) {
	//lc := db.LocalConversation{ConversationID: conversationID, IsPinned: isPinned}
	//if isPinned {
	c.setOneConversationPinned(callback, conversationID, isPinned, operationID)
	//err := c.db.UpdateConversation(&lc)
	//common.CheckDBErrCallback(callback, err, operationID)
	//} else {
	//	err := c.db.UnPinConversation(conversationID, constant.NotPinned)
	//	common.CheckDBErrCallback(callback, err, operationID)
	//}
}
func (c *Conversation) getServerConversationList(operationID string, timeout time.Duration) (server_api_params.GetAllConversationsResp, error) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	var req server_api_params.GetAllConversationsReq
	var resp server_api_params.GetAllConversationsResp
	req.OwnerUserID = c.loginUserID
	req.OperationID = operationID
	if timeout == 0 {
		err := c.p.PostReturn(constant.GetAllConversationsRouter, req, &resp.Conversations)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
			return resp, err
		}
	} else {
		err := c.p.PostReturnWithTimeOut(constant.GetAllConversationsRouter, req, &resp.Conversations, timeout)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
			return resp, err
		}
	}

	return resp, nil
}
func (c *Conversation) SyncConversations(operationID string, timeout time.Duration) {
	//log.Error(operationID,"SyncConversations start")
	var newConversationList []*model_struct.LocalConversation
	ccTime := time.Now()
	log.NewInfo(operationID, utils.GetSelfFuncName())
	conversationsOnServer, err := c.getServerConversationList(operationID, timeout)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
		return
	}
	log.Info(operationID, "get server cost time", time.Since(ccTime))
	cTime := time.Now()
	conversationsOnLocal, err := c.db.GetAllConversationListToSync()
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
	}
	log.Info(operationID, "get local cost time", time.Since(cTime))
	cTime = time.Now()
	conversationsOnLocalTempFormat := common.LocalTransferToTempConversation(conversationsOnLocal)
	conversationsOnServerTempFormat := common.ServerTransferToTempConversation(conversationsOnServer)
	conversationsOnServerLocalFormat := common.TransferToLocalConversation(conversationsOnServer)

	aInBNot, bInANot, sameA, sameB := common.CheckConversationListDiff(conversationsOnServerTempFormat, conversationsOnLocalTempFormat)
	log.Info(operationID, "diff server cost time", time.Since(cTime))

	log.NewInfo(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	log.NewInfo(operationID, "server have", len(aInBNot), "local have", len(bInANot))
	// server有 local没有
	// 可能是其他点开一下生成会话设置免打扰 插入到本地 不回调..
	cTime = time.Now()
	for _, index := range aInBNot {
		conversation := conversationsOnServerLocalFormat[index]
		var newConversation model_struct.LocalConversation
		newConversation.ConversationID = conversation.ConversationID
		newConversation.ConversationType = conversation.ConversationType
		newConversation.UserID = conversation.UserID
		newConversation.GroupID = conversation.GroupID
		newConversation.RecvMsgOpt = conversation.RecvMsgOpt
		newConversation.IsPinned = conversation.IsPinned
		newConversation.IsPrivateChat = conversation.IsPrivateChat
		newConversation.BurnDuration = conversation.BurnDuration
		newConversation.GroupAtType = conversation.GroupAtType
		newConversation.IsNotInGroup = conversation.IsNotInGroup
		newConversation.Ex = conversation.Ex
		newConversation.AttachedInfo = conversation.AttachedInfo
		newConversation.AttachedInfo = conversation.AttachedInfo
		newConversation.UpdateUnreadCountTime = conversation.UpdateUnreadCountTime
		//newConversation.UnreadCount = conversation.UnreadCount
		newConversationList = append(newConversationList, &newConversation)
		c.addFaceURLAndName(&newConversation)
		//err := c.db.InsertConversation(&newConversation)
		//if err != nil {
		//	log.NewError(operationID, utils.GetSelfFuncName(), "InsertConversation error", err.Error(), conversation)
		//	continue
		//}
	}
	log.Info(operationID, "Assemble a new conversations cost time", time.Since(cTime))
	//New conversation storage
	cTime = time.Now()
	err2 := c.db.BatchInsertConversationList(newConversationList)
	if err2 != nil {
		log.Error(operationID, "insert new conversation err:", err2.Error(), newConversationList)
	}
	log.Info(operationID, "batch insert cost time", time.Since(cTime))
	// 本地服务器有的会话 以服务器为准更新
	cTime = time.Now()
	var conversationChangedList []string
	for _, index := range sameA {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "server and client both have", *conversationsOnServerLocalFormat[index])
		err := c.db.UpdateConversationForSync(conversationsOnServerLocalFormat[index])
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), "UpdateConversation failed ", err.Error(), *conversationsOnServerLocalFormat[index])
			continue
		}
		conversationChangedList = append(conversationChangedList, conversationsOnServerLocalFormat[index].ConversationID)
	}
	// callback
	if len(conversationChangedList) > 0 {
		if err = common.TriggerCmdUpdateConversation(common.UpdateConNode{Action: constant.ConChange, Args: conversationChangedList}, c.GetCh()); err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
		}
	}
	log.Info(operationID, "batch update cost time", time.Since(cTime))

	// local有 server没有 代表没有修改公共字段
	for _, index := range bInANot {
		log.NewDebug(operationID, utils.GetSelfFuncName(), index, conversationsOnLocal[index].ConversationID,
			conversationsOnLocal[index].RecvMsgOpt, conversationsOnLocal[index].IsPinned, conversationsOnLocal[index].IsPrivateChat)
	}
	cTime = time.Now()
	conversationsOnLocal, err = c.db.GetAllConversationListToSync()
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
	}
	c.cache.UpdateConversations(conversationsOnLocal)
	log.Info(operationID, "cache update cost time", time.Since(cTime))
	log.Info(operationID, utils.GetSelfFuncName(), "all  cost time", time.Since(ccTime))
}
func (c *Conversation) SyncConversationUnreadCount(operationID string) {
	var conversationChangedList []string
	allConversations := c.cache.GetAllHasUnreadMessageConversations()
	log.Debug(operationID, "get unread message length is ", len(allConversations))
	for _, conversation := range allConversations {
		log.Debug(operationID, "has unread message conversation is:", *conversation)
		if deleteRows := c.db.DeleteConversationUnreadMessageList(conversation.ConversationID, conversation.UpdateUnreadCountTime); deleteRows > 0 {
			log.Debug(operationID, conversation.ConversationID, conversation.UpdateUnreadCountTime, "delete rows:", deleteRows)
			if err := c.db.DecrConversationUnreadCount(conversation.ConversationID, deleteRows); err != nil {
				log.Debug(operationID, conversation.ConversationID, conversation.UpdateUnreadCountTime, "decr unread count err:", err.Error())
			} else {
				conversationChangedList = append(conversationChangedList, conversation.ConversationID)
			}
		}
	}
	if len(conversationChangedList) > 0 {
		if err := common.TriggerCmdUpdateConversation(common.UpdateConNode{Action: constant.ConChange, Args: conversationChangedList}, c.GetCh()); err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
		}
	}

}
func (c *Conversation) FixVersionData() {
	switch constant.SdkVersion + constant.BigVersion + constant.UpdateVersion {
	case "v2.0.0":
		t := time.Now()
		groupIDList, err := c.db.GetReadDiffusionGroupIDList()
		if err != nil {
			log.Error("", "GetReadDiffusionGroupIDList failed ", err.Error())
			return
		}
		log.Info("", "fix version data start", groupIDList)
		for _, v := range groupIDList {
			err := c.db.SuperGroupUpdateSpecificContentTypeMessage(constant.ReactionMessageModifier, v, map[string]interface{}{"status": constant.MsgStatusFiltered})
			if err != nil {
				log.Error("", "SuperGroupUpdateSpecificContentTypeMessage failed ", err.Error())
				continue
			}
			msgList, err := c.db.SuperGroupSearchAllMessageByContentType(v, constant.ReactionMessageModifier)
			if err != nil {
				log.NewError("internal", "SuperGroupSearchMessageByContentTypeNotOffset failed", v, err.Error())
				continue
			}
			var reactionMsgIDList []string
			for _, value := range msgList {
				var n server_api_params.ReactionMessageModifierNotification
				err := json.Unmarshal([]byte(value.Content), &n)
				if err != nil {
					log.Error("internal", "unmarshal failed err:", err.Error(), *value)
					continue
				}
				reactionMsgIDList = append(reactionMsgIDList, n.ClientMsgID)
			}
			if len(reactionMsgIDList) > 0 {
				err := c.db.SuperGroupUpdateGroupMessageFields(reactionMsgIDList, v, map[string]interface{}{"is_react": true})
				if err != nil {
					log.Error("internal", "unmarshal failed err:", err.Error(), reactionMsgIDList, v)
					continue
				}
			}

		}
		log.Info("", "fix version data end", groupIDList, "cost time:", time.Since(t))

	default:

	}

}

func (c *Conversation) SyncOneConversation(conversationID, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "conversationID: ", conversationID)
	// todo
}

func (c *Conversation) findMessageList(req sdk.FindMessageListParams, operationID string) (r sdk.FindMessageListCallback) {
	type tempConversationAndMessageList struct {
		conversation *model_struct.LocalConversation
		msgIDList    []string
	}
	var s []*tempConversationAndMessageList
	for _, conversationsArgs := range req {
		localConversation, err := c.db.GetConversation(conversationsArgs.ConversationID)
		if err == nil {
			t := new(tempConversationAndMessageList)
			t.conversation = localConversation
			t.msgIDList = conversationsArgs.ClientMsgIDList
			s = append(s, t)
		} else {
			log.Error(operationID, "GetConversation err:", err.Error(), conversationsArgs.ConversationID)
		}
	}
	for _, v := range s {
		messages, err := c.db.GetMultipleMessageController(v.msgIDList, v.conversation.GroupID, v.conversation.ConversationType)
		if err == nil {
			var tempMessageList []*sdk_struct.MsgStruct
			for _, message := range messages {
				temp := sdk_struct.MsgStruct{}
				temp.ClientMsgID = message.ClientMsgID
				temp.ServerMsgID = message.ServerMsgID
				temp.CreateTime = message.CreateTime
				temp.SendTime = message.SendTime
				temp.SessionType = message.SessionType
				temp.SendID = message.SendID
				temp.RecvID = message.RecvID
				temp.MsgFrom = message.MsgFrom
				temp.ContentType = message.ContentType
				temp.SenderPlatformID = message.SenderPlatformID
				temp.SenderNickname = message.SenderNickname
				temp.SenderFaceURL = message.SenderFaceURL
				temp.Content = message.Content
				temp.Seq = message.Seq
				temp.IsRead = message.IsRead
				temp.Status = message.Status
				temp.AttachedInfo = message.AttachedInfo
				temp.Ex = message.Ex
				err := c.msgHandleByContentType(&temp)
				if err != nil {
					log.Error(operationID, "Parsing data error:", err.Error(), temp)
					continue
				}
				switch message.SessionType {
				case constant.GroupChatType:
					fallthrough
				case constant.SuperGroupChatType:
					temp.GroupID = temp.RecvID
					temp.RecvID = c.loginUserID
				}
				tempMessageList = append(tempMessageList, &temp)
			}
			findResultItem := sdk.SearchByConversationResult{}
			findResultItem.ConversationID = v.conversation.ConversationID
			findResultItem.FaceURL = v.conversation.FaceURL
			findResultItem.ShowName = v.conversation.ShowName
			findResultItem.ConversationType = v.conversation.ConversationType
			findResultItem.MessageList = tempMessageList
			findResultItem.MessageCount = len(findResultItem.MessageList)
			r.FindResultItems = append(r.FindResultItems, &findResultItem)
			r.TotalCount += findResultItem.MessageCount
		} else {
			log.Error(operationID, "GetMultipleMessageController err:", err.Error(), v)
		}
	}
	return r
}

func (c *Conversation) getHistoryMessageList(callback open_im_sdk_callback.Base, req sdk.GetHistoryMessageListParams, operationID string, isReverse bool) sdk.GetHistoryMessageListCallback {
	t := time.Now()
	var sourceID string
	var conversationID string
	var startTime int64
	var sessionType int
	var list []*model_struct.LocalChatLog
	var err error
	var messageList sdk_struct.NewMsgList
	var msg sdk_struct.MsgStruct
	var notStartTime bool
	if req.ConversationID != "" {
		conversationID = req.ConversationID
		lc, err := c.db.GetConversation(conversationID)
		if err != nil {
			return nil
		}
		switch lc.ConversationType {
		case constant.SingleChatType, constant.NotificationChatType:
			sourceID = lc.UserID
		case constant.GroupChatType, constant.SuperGroupChatType:
			sourceID = lc.GroupID
			msg.GroupID = lc.GroupID
		}
		sessionType = int(lc.ConversationType)
		if req.StartClientMsgID == "" {
			//startTime = lc.LatestMsgSendTime + TimeOffset
			////startTime = utils.GetCurrentTimestampByMill()
			notStartTime = true
		} else {
			msg.SessionType = lc.ConversationType
			msg.ClientMsgID = req.StartClientMsgID
			m, err := c.db.GetMessageController(&msg)
			common.CheckDBErrCallback(callback, err, operationID)
			startTime = m.SendTime
		}
	} else {
		if req.UserID == "" {
			newConversationID, newSessionType, err := c.getConversationTypeByGroupID(req.GroupID)
			common.CheckDBErrCallback(callback, err, operationID)
			sourceID = req.GroupID
			sessionType = int(newSessionType)
			conversationID = newConversationID
			msg.GroupID = req.GroupID
			msg.SessionType = newSessionType
		} else {
			sourceID = req.UserID
			conversationID = utils.GetConversationIDBySessionType(sourceID, constant.SingleChatType)
			sessionType = constant.SingleChatType
		}
		if req.StartClientMsgID == "" {
			//lc, err := c.db.GetConversation(conversationID)
			//if err != nil {
			//	return nil
			//}
			//startTime = lc.LatestMsgSendTime + TimeOffset
			//startTime = utils.GetCurrentTimestampByMill()
			notStartTime = true
		} else {
			msg.ClientMsgID = req.StartClientMsgID
			m, err := c.db.GetMessageController(&msg)
			common.CheckDBErrCallback(callback, err, operationID)
			startTime = m.SendTime
		}
	}
	log.Debug(operationID, "Assembly parameters cost time", time.Since(t))
	t = time.Now()
	log.Info(operationID, "sourceID:", sourceID, "startTime:", startTime, "count:", req.Count, "not start_time", notStartTime)
	if notStartTime {
		list, err = c.db.GetMessageListNoTimeController(sourceID, sessionType, req.Count, isReverse)
	} else {
		list, err = c.db.GetMessageListController(sourceID, sessionType, req.Count, startTime, isReverse)
	}
	log.Debug(operationID, "db cost time", time.Since(t))
	common.CheckDBErrCallback(callback, err, operationID)
	t = time.Now()
	for _, v := range list {
		temp := sdk_struct.MsgStruct{}
		tt := time.Now()
		temp.ClientMsgID = v.ClientMsgID
		temp.ServerMsgID = v.ServerMsgID
		temp.CreateTime = v.CreateTime
		temp.SendTime = v.SendTime
		temp.SessionType = v.SessionType
		temp.SendID = v.SendID
		temp.RecvID = v.RecvID
		temp.MsgFrom = v.MsgFrom
		temp.ContentType = v.ContentType
		temp.SenderPlatformID = v.SenderPlatformID
		temp.SenderNickname = v.SenderNickname
		temp.SenderFaceURL = v.SenderFaceURL
		temp.Content = v.Content
		temp.Seq = v.Seq
		temp.IsRead = v.IsRead
		temp.Status = v.Status
		temp.AttachedInfo = v.AttachedInfo
		temp.Ex = v.Ex
		temp.IsReact = v.IsReact
		temp.IsExternalExtensions = v.IsExternalExtensions
		err := c.msgHandleByContentType(&temp)
		if err != nil {
			log.Error(operationID, "Parsing data error:", err.Error(), temp)
			continue
		}
		log.Debug(operationID, "internal unmarshal cost time", time.Since(tt))

		switch sessionType {
		case constant.GroupChatType:
			fallthrough
		case constant.SuperGroupChatType:
			temp.GroupID = temp.RecvID
			temp.RecvID = c.loginUserID
		}
		messageList = append(messageList, &temp)
	}
	log.Debug(operationID, "unmarshal cost time", time.Since(t))
	t = time.Now()
	if !isReverse {
		sort.Sort(messageList)
	}
	log.Debug(operationID, "sort cost time", time.Since(t))
	return sdk.GetHistoryMessageListCallback(messageList)
}
func (c *Conversation) getAdvancedHistoryMessageList(callback open_im_sdk_callback.Base, req sdk.GetAdvancedHistoryMessageListParams, operationID string, isReverse bool) sdk.GetAdvancedHistoryMessageListCallback {
	t := time.Now()
	var messageListCallback sdk.GetAdvancedHistoryMessageListCallback
	var sourceID string
	var conversationID string
	var startTime int64

	var sessionType int
	var list []*model_struct.LocalChatLog
	var err error
	var messageList sdk_struct.NewMsgList
	var msg sdk_struct.MsgStruct
	var notStartTime bool
	if req.ConversationID != "" {
		conversationID = req.ConversationID
		lc, err := c.db.GetConversation(conversationID)
		if err != nil {
			messageListCallback.ErrCode = 100
			messageListCallback.ErrMsg = "conversation get err"
			return messageListCallback
		}
		switch lc.ConversationType {
		case constant.SingleChatType, constant.NotificationChatType:
			sourceID = lc.UserID
		case constant.GroupChatType, constant.SuperGroupChatType:
			sourceID = lc.GroupID
			msg.GroupID = lc.GroupID
		}
		sessionType = int(lc.ConversationType)
		if req.StartClientMsgID == "" {
			//startTime = lc.LatestMsgSendTime + TimeOffset
			////startTime = utils.GetCurrentTimestampByMill()
			notStartTime = true
		} else {
			msg.SessionType = lc.ConversationType
			msg.ClientMsgID = req.StartClientMsgID
			m, err := c.db.GetMessageController(&msg)
			common.CheckDBErrCallback(callback, err, operationID)
			startTime = m.SendTime
		}
	} else {
		if req.UserID == "" {
			newConversationID, newSessionType, err := c.getConversationTypeByGroupID(req.GroupID)
			common.CheckDBErrCallback(callback, err, operationID)
			sourceID = req.GroupID
			sessionType = int(newSessionType)
			conversationID = newConversationID
			msg.GroupID = req.GroupID
			msg.SessionType = newSessionType
		} else {
			sourceID = req.UserID
			conversationID = utils.GetConversationIDBySessionType(sourceID, constant.SingleChatType)
			sessionType = constant.SingleChatType
		}
		if req.StartClientMsgID == "" {
			//lc, err := c.db.GetConversation(conversationID)
			//if err != nil {
			//	return nil
			//}
			//startTime = lc.LatestMsgSendTime + TimeOffset
			//startTime = utils.GetCurrentTimestampByMill()
			notStartTime = true
		} else {
			msg.ClientMsgID = req.StartClientMsgID
			m, err := c.db.GetMessageController(&msg)
			common.CheckDBErrCallback(callback, err, operationID)
			startTime = m.SendTime
		}
	}
	log.Debug(operationID, "Assembly parameters cost time", time.Since(t))
	t = time.Now()
	log.Info(operationID, "sourceID:", sourceID, "startTime:", startTime, "count:", req.Count, "not start_time", notStartTime)
	if notStartTime {
		list, err = c.db.GetMessageListNoTimeController(sourceID, sessionType, req.Count, isReverse)
	} else {
		list, err = c.db.GetMessageListController(sourceID, sessionType, req.Count, startTime, isReverse)
	}
	log.Error(operationID, "db cost time", time.Since(t), len(list), err, sourceID)
	t = time.Now()
	common.CheckDBErrCallback(callback, err, operationID)
	if isReverse {
		if len(list) < req.Count {
			messageListCallback.IsEnd = true
		}
	} else {
		switch sessionType {
		case constant.SuperGroupChatType:
			if len(list) < req.Count {
				var minSeq uint32
				var maxSeq uint32
				resp, err := c.SendReqWaitResp(&server_api_params.GetMaxAndMinSeqReq{UserID: c.loginUserID, GroupIDList: []string{sourceID}}, constant.WSGetNewestSeq, 1, 2, c.loginUserID, operationID)
				if err != nil {
					log.Error(operationID, "SendReqWaitResp failed ", err.Error(), constant.WSGetNewestSeq, 30, c.loginUserID)
				} else {
					var wsSeqResp server_api_params.GetMaxAndMinSeqResp
					err = proto.Unmarshal(resp.Data, &wsSeqResp)
					if err != nil {
						log.Error(operationID, "Unmarshal failed", err.Error())
					} else if wsSeqResp.ErrCode != 0 {
						log.Error(operationID, "GetMaxAndMinSeqReq failed ", wsSeqResp.ErrCode, wsSeqResp.ErrMsg)
					} else {
						if value, ok := wsSeqResp.GroupMaxAndMinSeq[sourceID]; ok {
							minSeq = value.MinSeq
							if value.MinSeq == 0 {
								minSeq = 1
							}
							maxSeq = value.MaxSeq
						}
					}
				}
				log.Error(operationID, "from server min seq is", minSeq, maxSeq)
				seq, err := c.db.SuperGroupGetNormalMinSeq(sourceID)
				if err != nil {
					log.Error(operationID, "SuperGroupGetNormalMinSeq err:", err.Error())
				}
				log.Error(operationID, sourceID+":table min seq is ", seq)
				if seq != 0 {
					if seq <= minSeq {
						messageListCallback.IsEnd = true
					} else {
						seqList := func(seq uint32) (seqList []uint32) {
							startSeq := int64(seq) - constant.PullMsgNumForReadDiffusion
							if startSeq <= 0 {
								startSeq = 1
							}
							log.Debug(operationID, "pull start is ", startSeq)
							if startSeq < int64(minSeq) {
								startSeq = int64(minSeq)
							}
							for i := startSeq; i < int64(seq); i++ {
								seqList = append(seqList, uint32(i))
							}
							log.Debug(operationID, "pull seqList is ", seqList)
							return seqList
						}(seq)
						log.Debug(operationID, "pull seqList is ", seqList, len(seqList))
						if len(seqList) > 0 {
							c.pullMessageAndReGetHistoryMessages(sourceID, seqList, notStartTime, isReverse, req.Count, sessionType, startTime, &list, &messageListCallback, operationID)
						}
					}
				} else {
					//local don't have messages,本地无消息，但是服务器最大消息不为0
					if int64(maxSeq)-int64(minSeq) > 0 {
						messageListCallback.IsEnd = false
					} else {
						messageListCallback.IsEnd = true
					}

				}
			} else if len(list) == req.Count {
				maxSeq, minSeq, haveSeqList := func(messages []*model_struct.LocalChatLog) (max, min uint32, seqList []uint32) {
					for _, message := range messages {
						if message.Seq != 0 {
							max = message.Seq
							min = message.Seq
							break
						}
					}
					for i := 0; i < len(messages); i++ {
						if messages[i].Seq != 0 {
							seqList = append(seqList, messages[i].Seq)
						}
						if messages[i].Seq > max {
							max = messages[i].Seq

						}
						if messages[i].Seq < min {
							min = messages[i].Seq
						}
					}
					return max, min, seqList
				}(list)
				log.Debug(operationID, "get message from local db max seq:", maxSeq, "minSeq:", minSeq, "haveSeqList:", haveSeqList, "length:", len(haveSeqList))
				if maxSeq != 0 && minSeq != 0 {
					successiveSeqList := func(max, min uint32) (seqList []uint32) {
						for i := min; i <= max; i++ {
							seqList = append(seqList, i)
						}
						return seqList
					}(maxSeq, minSeq)
					lostSeqList := utils.DifferenceSubset(successiveSeqList, haveSeqList)
					lostSeqListLength := len(lostSeqList)
					log.Debug(operationID, "get lost seqList is :", lostSeqList, "length:", lostSeqListLength)
					if lostSeqListLength > 0 {
						var pullSeqList []uint32
						if lostSeqListLength <= constant.PullMsgNumForReadDiffusion {
							pullSeqList = lostSeqList
						} else {
							pullSeqList = lostSeqList[lostSeqListLength-constant.PullMsgNumForReadDiffusion : lostSeqListLength]
						}
						c.pullMessageAndReGetHistoryMessages(sourceID, pullSeqList, notStartTime, isReverse, req.Count, sessionType, startTime, &list, &messageListCallback, operationID)
					} else {
						if req.LastMinSeq != 0 {
							var thisMaxSeq uint32
							for i := 0; i < len(list); i++ {
								if list[i].Seq != 0 && thisMaxSeq == 0 {
									thisMaxSeq = list[i].Seq
								}
								if list[i].Seq > thisMaxSeq {
									thisMaxSeq = list[i].Seq
								}
							}
							log.Debug(operationID, "get lost LastMinSeq is :", req.LastMinSeq, "thisMaxSeq is :", thisMaxSeq)
							if thisMaxSeq != 0 {
								if thisMaxSeq+1 != req.LastMinSeq {
									startSeq := int64(req.LastMinSeq) - constant.PullMsgNumForReadDiffusion
									if startSeq <= int64(thisMaxSeq) {
										startSeq = int64(thisMaxSeq) + 1
									}
									successiveSeqList := func(max, min uint32) (seqList []uint32) {
										for i := min; i <= max; i++ {
											seqList = append(seqList, i)
										}
										return seqList
									}(req.LastMinSeq-1, uint32(startSeq))
									log.Debug(operationID, "get lost successiveSeqList is :", successiveSeqList, len(successiveSeqList))
									if len(successiveSeqList) > 0 {
										c.pullMessageAndReGetHistoryMessages(sourceID, successiveSeqList, notStartTime, isReverse, req.Count, sessionType, startTime, &list, &messageListCallback, operationID)
									}
								}

							}

						}
					}

				}
			}
		default:
			if len(list) < req.Count {
				messageListCallback.IsEnd = true
			}
		}

	}

	log.Debug(operationID, "pull cost time", time.Since(t))
	t = time.Now()
	var thisMinSeq uint32
	for _, v := range list {
		if v.Seq != 0 && thisMinSeq == 0 {
			thisMinSeq = v.Seq
		}
		if v.Seq < thisMinSeq {
			thisMinSeq = v.Seq
		}
		temp := sdk_struct.MsgStruct{}
		tt := time.Now()
		temp.ClientMsgID = v.ClientMsgID
		temp.ServerMsgID = v.ServerMsgID
		temp.CreateTime = v.CreateTime
		temp.SendTime = v.SendTime
		temp.SessionType = v.SessionType
		temp.SendID = v.SendID
		temp.RecvID = v.RecvID
		temp.MsgFrom = v.MsgFrom
		temp.ContentType = v.ContentType
		temp.SenderPlatformID = v.SenderPlatformID
		temp.SenderNickname = v.SenderNickname
		temp.SenderFaceURL = v.SenderFaceURL
		temp.Content = v.Content
		temp.Seq = v.Seq
		temp.IsRead = v.IsRead
		temp.Status = v.Status
		temp.AttachedInfo = v.AttachedInfo
		temp.Ex = v.Ex
		temp.IsReact = v.IsReact
		temp.IsExternalExtensions = v.IsExternalExtensions
		err := c.msgHandleByContentType(&temp)
		if err != nil {
			log.Error(operationID, "Parsing data error:", err.Error(), temp)
			continue
		}
		log.Debug(operationID, "internal unmarshal cost time", time.Since(tt))

		switch sessionType {
		case constant.GroupChatType:
			fallthrough
		case constant.SuperGroupChatType:
			temp.GroupID = temp.RecvID
			temp.RecvID = c.loginUserID
		}
		messageList = append(messageList, &temp)
	}
	log.Debug(operationID, "unmarshal cost time", time.Since(t))
	t = time.Now()
	if !isReverse {
		sort.Sort(messageList)
	}
	log.Debug(operationID, "sort cost time", time.Since(t))
	messageListCallback.MessageList = messageList
	messageListCallback.LastMinSeq = thisMinSeq
	return messageListCallback
}

func (c *Conversation) getAdvancedHistoryMessageList2(callback open_im_sdk_callback.Base, req sdk.GetAdvancedHistoryMessageListParams, operationID string, isReverse bool) sdk.GetAdvancedHistoryMessageListCallback {
	t := time.Now()
	var messageListCallback sdk.GetAdvancedHistoryMessageListCallback
	var sourceID string
	var conversationID string
	var startTime int64

	var sessionType int
	var list []*model_struct.LocalChatLog
	var err error
	var messageList sdk_struct.NewMsgList
	var msg sdk_struct.MsgStruct
	var notStartTime bool
	if req.ConversationID != "" {
		conversationID = req.ConversationID
		lc, err := c.db.GetConversation(conversationID)
		if err != nil {
			messageListCallback.ErrCode = 100
			messageListCallback.ErrMsg = "conversation get err"
			return messageListCallback
		}
		switch lc.ConversationType {
		case constant.SingleChatType, constant.NotificationChatType:
			sourceID = lc.UserID
		case constant.GroupChatType, constant.SuperGroupChatType:
			sourceID = lc.GroupID
			msg.GroupID = lc.GroupID
		}
		sessionType = int(lc.ConversationType)
		if req.StartClientMsgID == "" {
			//startTime = lc.LatestMsgSendTime + TimeOffset
			////startTime = utils.GetCurrentTimestampByMill()
			notStartTime = true
		} else {
			msg.SessionType = lc.ConversationType
			msg.ClientMsgID = req.StartClientMsgID
			m, err := c.db.GetMessageController(&msg)
			common.CheckDBErrCallback(callback, err, operationID)
			startTime = m.SendTime
		}
	} else {
		if req.UserID == "" {
			newConversationID, newSessionType, err := c.getConversationTypeByGroupID(req.GroupID)
			common.CheckDBErrCallback(callback, err, operationID)
			sourceID = req.GroupID
			sessionType = int(newSessionType)
			conversationID = newConversationID
			msg.GroupID = req.GroupID
			msg.SessionType = newSessionType
		} else {
			sourceID = req.UserID
			conversationID = utils.GetConversationIDBySessionType(sourceID, constant.SingleChatType)
			sessionType = constant.SingleChatType
		}
		if req.StartClientMsgID == "" {
			//lc, err := c.db.GetConversation(conversationID)
			//if err != nil {
			//	return nil
			//}
			//startTime = lc.LatestMsgSendTime + TimeOffset
			//startTime = utils.GetCurrentTimestampByMill()
			notStartTime = true
		} else {
			msg.ClientMsgID = req.StartClientMsgID
			m, err := c.db.GetMessageController(&msg)
			common.CheckDBErrCallback(callback, err, operationID)
			startTime = m.SendTime
		}
	}
	log.Debug(operationID, "Assembly parameters cost time", time.Since(t))
	t = time.Now()
	log.Info(operationID, "sourceID:", sourceID, "startTime:", startTime, "count:", req.Count, "not start_time", notStartTime)
	if notStartTime {
		list, err = c.db.GetMessageListNoTimeController(sourceID, sessionType, req.Count, isReverse)
	} else {
		list, err = c.db.GetMessageListController(sourceID, sessionType, req.Count, startTime, isReverse)
	}
	log.Error(operationID, "db cost time", time.Since(t), len(list), err, sourceID)
	t = time.Now()
	common.CheckDBErrCallback(callback, err, operationID)
	if sessionType == constant.SuperGroupChatType {
		rawMessageLength := len(list)
		if rawMessageLength < req.Count {
			maxSeq, minSeq, lostSeqListLength := c.messageBlocksInternalContinuityCheck(sourceID, notStartTime, isReverse, req.Count, sessionType, startTime, &list, &messageListCallback, operationID)
			_ = c.messageBlocksBetweenContinuityCheck(req.LastMinSeq, maxSeq, sourceID, notStartTime, isReverse, req.Count, sessionType, startTime, &list, &messageListCallback, operationID)
			if minSeq == 1 && lostSeqListLength == 0 {
				messageListCallback.IsEnd = true
			} else {
				c.messageBlocksEndContinuityCheck(minSeq, sourceID, notStartTime, isReverse, req.Count, sessionType, startTime, &list, &messageListCallback, operationID)
			}
		} else {
			maxSeq, _, _ := c.messageBlocksInternalContinuityCheck(sourceID, notStartTime, isReverse, req.Count, sessionType, startTime, &list, &messageListCallback, operationID)
			c.messageBlocksBetweenContinuityCheck(req.LastMinSeq, maxSeq, sourceID, notStartTime, isReverse, req.Count, sessionType, startTime, &list, &messageListCallback, operationID)

		}

	}
	//if len(list) < req.Count && sessionType == constant.SuperGroupChatType {
	//
	//} else if len(list) == req.Count && sessionType == constant.SuperGroupChatType {
	//	if maxSeq != 0 && minSeq != 0 {
	//		successiveSeqList := func(max, min uint32) (seqList []uint32) {
	//			for i := min; i <= max; i++ {
	//				seqList = append(seqList, i)
	//			}
	//			return seqList
	//		}(maxSeq, minSeq)
	//		lostSeqList := utils.DifferenceSubset(successiveSeqList, haveSeqList)
	//		lostSeqListLength := len(lostSeqList)
	//		log.Debug(operationID, "get lost seqList is :", lostSeqList, "length:", lostSeqListLength)
	//		if lostSeqListLength > 0 {
	//			var pullSeqList []uint32
	//			if lostSeqListLength <= constant.PullMsgNumForReadDiffusion {
	//				pullSeqList = lostSeqList
	//			} else {
	//				pullSeqList = lostSeqList[lostSeqListLength-constant.PullMsgNumForReadDiffusion : lostSeqListLength]
	//			}
	//			c.pullMessageAndReGetHistoryMessages(sourceID, pullSeqList, notStartTime, isReverse, req.Count, sessionType, startTime, &list, &messageListCallback, operationID)
	//		} else {
	//			if req.LastMinSeq != 0 {
	//				var thisMaxSeq uint32
	//				for i := 0; i < len(list); i++ {
	//					if list[i].Seq != 0 && thisMaxSeq == 0 {
	//						thisMaxSeq = list[i].Seq
	//					}
	//					if list[i].Seq > thisMaxSeq {
	//						thisMaxSeq = list[i].Seq
	//					}
	//				}
	//				log.Debug(operationID, "get lost LastMinSeq is :", req.LastMinSeq, "thisMaxSeq is :", thisMaxSeq)
	//				if thisMaxSeq != 0 {
	//					if thisMaxSeq+1 != req.LastMinSeq {
	//						startSeq := int64(req.LastMinSeq) - constant.PullMsgNumForReadDiffusion
	//						if startSeq <= int64(thisMaxSeq) {
	//							startSeq = int64(thisMaxSeq) + 1
	//						}
	//						successiveSeqList := func(max, min uint32) (seqList []uint32) {
	//							for i := min; i <= max; i++ {
	//								seqList = append(seqList, i)
	//							}
	//							return seqList
	//						}(req.LastMinSeq-1, uint32(startSeq))
	//						log.Debug(operationID, "get lost successiveSeqList is :", successiveSeqList, len(successiveSeqList))
	//						if len(successiveSeqList) > 0 {
	//							c.pullMessageAndReGetHistoryMessages(sourceID, successiveSeqList, notStartTime, isReverse, req.Count, sessionType, startTime, &list, &messageListCallback, operationID)
	//						}
	//					}
	//
	//				}
	//
	//			}
	//		}
	//
	//	}
	//}
	log.Debug(operationID, "pull cost time", time.Since(t))
	t = time.Now()
	var thisMinSeq uint32
	for _, v := range list {
		if v.Seq != 0 && thisMinSeq == 0 {
			thisMinSeq = v.Seq
		}
		if v.Seq < thisMinSeq {
			thisMinSeq = v.Seq
		}
		temp := sdk_struct.MsgStruct{}
		tt := time.Now()
		temp.ClientMsgID = v.ClientMsgID
		temp.ServerMsgID = v.ServerMsgID
		temp.CreateTime = v.CreateTime
		temp.SendTime = v.SendTime
		temp.SessionType = v.SessionType
		temp.SendID = v.SendID
		temp.RecvID = v.RecvID
		temp.MsgFrom = v.MsgFrom
		temp.ContentType = v.ContentType
		temp.SenderPlatformID = v.SenderPlatformID
		temp.SenderNickname = v.SenderNickname
		temp.SenderFaceURL = v.SenderFaceURL
		temp.Content = v.Content
		temp.Seq = v.Seq
		temp.IsRead = v.IsRead
		temp.Status = v.Status
		temp.AttachedInfo = v.AttachedInfo
		temp.Ex = v.Ex
		err := c.msgHandleByContentType(&temp)
		if err != nil {
			log.Error(operationID, "Parsing data error:", err.Error(), temp)
			continue
		}
		log.Debug(operationID, "internal unmarshal cost time", time.Since(tt))

		switch sessionType {
		case constant.GroupChatType:
			fallthrough
		case constant.SuperGroupChatType:
			temp.GroupID = temp.RecvID
			temp.RecvID = c.loginUserID
		}
		messageList = append(messageList, &temp)
	}
	log.Debug(operationID, "unmarshal cost time", time.Since(t))
	t = time.Now()
	if !isReverse {
		sort.Sort(messageList)
	}
	log.Debug(operationID, "sort cost time", time.Since(t))
	messageListCallback.MessageList = messageList
	if thisMinSeq == 0 {
		thisMinSeq = req.LastMinSeq
	}
	messageListCallback.LastMinSeq = thisMinSeq
	return messageListCallback
}

func (c *Conversation) revokeOneMessage(callback open_im_sdk_callback.Base, req sdk.RevokeMessageParams, operationID string) {
	var recvID, groupID string
	var localMessage model_struct.LocalChatLog
	var lc model_struct.LocalConversation
	var conversationID string
	message, err := c.db.GetMessageController((*sdk_struct.MsgStruct)(&req))
	common.CheckDBErrCallback(callback, err, operationID)
	if message.Status != constant.MsgStatusSendSuccess {
		common.CheckAnyErrCallback(callback, 201, errors.New("only send success message can be revoked"), operationID)
	}
	if message.SendID != c.loginUserID {
		common.CheckAnyErrCallback(callback, 201, errors.New("only you send message can be revoked"), operationID)
	}
	//Send message internally
	switch req.SessionType {
	case constant.SingleChatType:
		recvID = req.RecvID
		conversationID = utils.GetConversationIDBySessionType(recvID, constant.SingleChatType)
	case constant.GroupChatType:
		groupID = req.GroupID
		conversationID = utils.GetConversationIDBySessionType(groupID, constant.GroupChatType)
	case constant.SuperGroupChatType:
		groupID = req.GroupID
		conversationID = utils.GetConversationIDBySessionType(groupID, constant.SuperGroupChatType)
	default:
		common.CheckAnyErrCallback(callback, 201, errors.New("SessionType err"), operationID)
	}
	req.Content = message.ClientMsgID
	req.ClientMsgID = utils.GetMsgID(message.SendID)
	req.ContentType = constant.Revoke
	req.SendTime = 0
	req.CreateTime = utils.GetCurrentTimestampByMill()
	options := make(map[string]bool, 5)
	utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
	utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
	resp, _ := c.InternalSendMessage(callback, (*sdk_struct.MsgStruct)(&req), recvID, groupID, operationID, &server_api_params.OfflinePushInfo{}, false, options)
	req.ServerMsgID = resp.ServerMsgID
	req.SendTime = resp.SendTime
	req.Status = constant.MsgStatusSendSuccess
	msgStructToLocalChatLog(&localMessage, (*sdk_struct.MsgStruct)(&req))
	err = c.db.InsertMessageController(&localMessage)
	if err != nil {
		log.Error(operationID, "inset into chat log err", localMessage, req)
	}
	err = c.db.UpdateColumnsMessageController(req.Content, groupID, req.SessionType, map[string]interface{}{"status": constant.MsgStatusRevoked})
	if err != nil {
		log.Error(operationID, "update revoke message err", localMessage, req)
	}
	lc.LatestMsg = utils.StructToJsonString(req)
	lc.LatestMsgSendTime = req.SendTime
	lc.ConversationID = conversationID
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: lc.ConversationID, Action: constant.AddConOrUpLatMsg, Args: lc}, c.GetCh())
}
func (c *Conversation) newRevokeOneMessage(callback open_im_sdk_callback.Base, req sdk.RevokeMessageParams, operationID string) {
	var recvID, groupID string
	var localMessage model_struct.LocalChatLog
	var revokeMessage sdk_struct.MessageRevoked
	var lc model_struct.LocalConversation
	var conversationID string
	message, err := c.db.GetMessageController((*sdk_struct.MsgStruct)(&req))
	common.CheckDBErrCallback(callback, err, operationID)
	if message.Status != constant.MsgStatusSendSuccess {
		common.CheckAnyErrCallback(callback, 201, errors.New("only send success message can be revoked"), operationID)
	}
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.AdvancedRevoke, operationID)
	revokeMessage.ClientMsgID = message.ClientMsgID
	revokeMessage.RevokerID = c.loginUserID
	revokeMessage.RevokeTime = utils.GetCurrentTimestampBySecond()
	revokeMessage.RevokerNickname = s.SenderNickname
	revokeMessage.SourceMessageSendTime = message.SendTime
	revokeMessage.SessionType = message.SessionType
	revokeMessage.SourceMessageSendID = message.SendID
	revokeMessage.SourceMessageSenderNickname = message.SenderNickname
	revokeMessage.Seq = message.Seq
	revokeMessage.Ex = message.Ex
	//Send message internally
	switch message.SessionType {
	case constant.SingleChatType:
		if message.SendID != c.loginUserID {
			common.CheckAnyErrCallback(callback, 201, errors.New("only you send message can be revoked"), operationID)
		}
		recvID = message.RecvID
		conversationID = utils.GetConversationIDBySessionType(recvID, constant.SingleChatType)
	case constant.GroupChatType:
		if message.SendID != c.loginUserID {
			ownerID, adminIDList, err := c.group.GetGroupOwnerIDAndAdminIDList(message.RecvID, operationID)
			common.CheckDBErrCallback(callback, err, operationID)
			if c.loginUserID == ownerID {
				revokeMessage.RevokerRole = constant.GroupOwner
			} else if utils.IsContain(c.loginUserID, adminIDList) {
				if utils.IsContain(message.SendID, adminIDList) || message.SendID == ownerID {
					common.CheckAnyErrCallback(callback, 201, errors.New("you do not have this permission"), operationID)
				} else {
					revokeMessage.RevokerRole = constant.GroupAdmin
				}

			} else {
				common.CheckAnyErrCallback(callback, 201, errors.New("you do not have this permission"), operationID)
			}
		}
		groupID = message.RecvID
		conversationID = utils.GetConversationIDBySessionType(groupID, constant.GroupChatType)
	case constant.SuperGroupChatType:
		if message.SendID != c.loginUserID {
			ownerID, adminIDList, err := c.group.GetGroupOwnerIDAndAdminIDList(message.RecvID, operationID)
			common.CheckDBErrCallback(callback, err, operationID)
			if c.loginUserID == ownerID {
				revokeMessage.RevokerRole = constant.GroupOwner
			} else if utils.IsContain(c.loginUserID, adminIDList) {
				if utils.IsContain(message.SendID, adminIDList) || message.SendID == ownerID {
					common.CheckAnyErrCallback(callback, 201, errors.New("you do not have this permission"), operationID)
				} else {
					revokeMessage.RevokerRole = constant.GroupAdmin
				}

			} else {
				common.CheckAnyErrCallback(callback, 201, errors.New("you do not have this permission"), operationID)
			}
		}
		groupID = message.RecvID
		conversationID = utils.GetConversationIDBySessionType(groupID, constant.SuperGroupChatType)
	default:
		common.CheckAnyErrCallback(callback, 201, errors.New("SessionType err"), operationID)
	}
	s.Content = utils.StructToJsonString(revokeMessage)
	options := make(map[string]bool, 5)
	utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
	utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
	resp, _ := c.InternalSendMessage(callback, &s, recvID, groupID, operationID, &server_api_params.OfflinePushInfo{}, false, options)
	s.ServerMsgID = resp.ServerMsgID
	s.SendTime = message.SendTime //New message takes the old place
	s.Status = constant.MsgStatusSendSuccess
	msgStructToLocalChatLog(&localMessage, &s)
	err = c.db.InsertMessageController(&localMessage)
	if err != nil {
		log.Error(operationID, "inset into chat log err", localMessage, s)
	}
	err = c.db.UpdateColumnsMessageController(message.ClientMsgID, groupID, message.SessionType, map[string]interface{}{"status": constant.MsgStatusRevoked})
	if err != nil {
		log.Error(operationID, "update revoke message err", localMessage, message, err.Error())
	}
	s.SendTime = resp.SendTime
	lc.LatestMsg = utils.StructToJsonString(s)
	lc.LatestMsgSendTime = s.SendTime
	lc.ConversationID = conversationID
	s.GroupID = groupID
	s.RecvID = recvID
	c.newRevokeMessage([]*sdk_struct.MsgStruct{&s})
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: lc.ConversationID, Action: constant.AddConOrUpLatMsg, Args: lc}, c.GetCh())
}

func (c *Conversation) typingStatusUpdate(callback open_im_sdk_callback.Base, recvID, msgTip, operationID string) {
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Typing, operationID)
	s.Content = msgTip
	options := make(map[string]bool, 6)
	utils.SetSwitchFromOptions(options, constant.IsHistory, false)
	utils.SetSwitchFromOptions(options, constant.IsPersistent, false)
	utils.SetSwitchFromOptions(options, constant.IsSenderSync, false)
	utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsSenderConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
	utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
	c.InternalSendMessage(callback, &s, recvID, "", operationID, &server_api_params.OfflinePushInfo{}, true, options)

}

func (c *Conversation) markC2CMessageAsRead(callback open_im_sdk_callback.Base, msgIDList sdk.MarkC2CMessageAsReadParams, userID, operationID string) {
	var localMessage model_struct.LocalChatLog
	var newMessageIDList []string
	messages, err := c.db.GetMultipleMessage(msgIDList)
	common.CheckDBErrCallback(callback, err, operationID)
	for _, v := range messages {
		if v.IsRead == false && v.ContentType < constant.NotificationBegin && v.SendID != c.loginUserID {
			newMessageIDList = append(newMessageIDList, v.ClientMsgID)
		}
	}
	if len(newMessageIDList) == 0 {
		common.CheckAnyErrCallback(callback, 201, errors.New("message has been marked read or sender is yourself or notification message not support"), operationID)
	}
	conversationID := utils.GetConversationIDBySessionType(userID, constant.SingleChatType)
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.HasReadReceipt, operationID)
	s.Content = utils.StructToJsonString(newMessageIDList)
	options := make(map[string]bool, 5)
	utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsSenderConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
	utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
	//If there is an error, the coroutine ends, so judgment is not  required
	resp, _ := c.InternalSendMessage(callback, &s, userID, "", operationID, &server_api_params.OfflinePushInfo{}, false, options)
	s.ServerMsgID = resp.ServerMsgID
	s.SendTime = resp.SendTime
	s.Status = constant.MsgStatusFiltered
	msgStructToLocalChatLog(&localMessage, &s)
	err = c.db.InsertMessage(&localMessage)
	if err != nil {
		log.Error(operationID, "inset into chat log err", localMessage, s, err.Error())
	}

	err2 := c.db.UpdateSingleMessageHasRead(userID, newMessageIDList)
	if err2 != nil {
		log.Error(operationID, "update message has read error", newMessageIDList, userID, err2.Error())
	}
	newMessages, err3 := c.db.GetMultipleMessage(newMessageIDList)
	if err3 != nil {
		log.Error(operationID, "get messages error", newMessageIDList, userID, err3.Error())
	}
	for _, v := range newMessages {
		attachInfo := sdk_struct.AttachedInfoElem{}
		_ = utils.JsonStringToStruct(v.AttachedInfo, &attachInfo)
		attachInfo.HasReadTime = s.SendTime
		v.AttachedInfo = utils.StructToJsonString(attachInfo)
		err = c.db.UpdateMessage(v)
		if err != nil {
			log.Error("internal", "setMessageHasReadByMsgID err:", err, "ClientMsgID", v)
			continue
		}
	}
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.UpdateLatestMessageChange}, c.GetCh())
	//_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.ch)
}
func (c *Conversation) markGroupMessageAsRead(callback open_im_sdk_callback.Base, msgIDList sdk.MarkGroupMessageAsReadParams, groupID, operationID string) {
	conversationID, conversationType, err := c.getConversationTypeByGroupID(groupID)
	common.CheckAnyErrCallback(callback, 202, err, operationID)
	if len(msgIDList) == 0 {
		c.setOneConversationUnread(callback, conversationID, 0, operationID)
		_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.UnreadCountSetZero}, c.GetCh())
		_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())
		return
	}
	var localMessage model_struct.LocalChatLog
	allUserMessage := make(map[string][]string, 3)
	messages, err := c.db.GetMultipleMessageController(msgIDList, groupID, conversationType)
	common.CheckDBErrCallback(callback, err, operationID)
	for _, v := range messages {
		log.Debug(operationID, "get group info is test2", v.ClientMsgID, v.SessionType)
		if v.IsRead == false && v.ContentType < constant.NotificationBegin && v.SendID != c.loginUserID {
			if msgIDList, ok := allUserMessage[v.SendID]; ok {
				msgIDList = append(msgIDList, v.ClientMsgID)
				allUserMessage[v.SendID] = msgIDList
			} else {
				allUserMessage[v.SendID] = []string{v.ClientMsgID}
			}
		}
	}
	if len(allUserMessage) == 0 {
		common.CheckAnyErrCallback(callback, 201, errors.New("message has been marked read or sender is yourself or notification message not support"), operationID)
	}

	for userID, list := range allUserMessage {
		s := sdk_struct.MsgStruct{}
		s.GroupID = groupID
		c.initBasicInfo(&s, constant.UserMsgType, constant.GroupHasReadReceipt, operationID)
		s.Content = utils.StructToJsonString(list)
		options := make(map[string]bool, 5)
		utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
		utils.SetSwitchFromOptions(options, constant.IsSenderConversationUpdate, false)
		utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
		utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
		//If there is an error, the coroutine ends, so judgment is not  required
		resp, _ := c.InternalSendMessage(callback, &s, userID, "", operationID, &server_api_params.OfflinePushInfo{}, false, options)
		s.ServerMsgID = resp.ServerMsgID
		s.SendTime = resp.SendTime
		s.Status = constant.MsgStatusFiltered
		msgStructToLocalChatLog(&localMessage, &s)
		err = c.db.InsertMessageController(&localMessage)
		if err != nil {
			log.Error(operationID, "inset into chat log err", localMessage, s, err.Error())
		}
		log.Debug(operationID, "get group info is test3", list, conversationType)
		err2 := c.db.UpdateGroupMessageHasReadController(list, groupID, conversationType)
		if err2 != nil {
			log.Error(operationID, "update message has read err", list, userID, err2.Error())
		}
	}
}

//	func (c *Conversation) markMessageAsReadByConID(callback open_im_sdk_callback.Base, msgIDList sdk.MarkMessageAsReadByConIDParams, conversationID, operationID string) {
//		var localMessage db.LocalChatLog
//		var newMessageIDList []string
//		messages, err := c.db.GetMultipleMessage(msgIDList)
//		common.CheckDBErrCallback(callback, err, operationID)
//		for _, v := range messages {
//			if v.IsRead == false && v.ContentType < constant.NotificationBegin && v.SendID != c.loginUserID {
//				newMessageIDList = append(newMessageIDList, v.ClientMsgID)
//			}
//		}
//		if len(newMessageIDList) == 0 {
//			common.CheckAnyErrCallback(callback, 201, errors.New("message has been marked read or sender is yourself"), operationID)
//		}
//		conversationID := utils.GetConversationIDBySessionType(userID, constant.SingleChatType)
//		s := sdk_struct.MsgStruct{}
//		c.initBasicInfo(&s, constant.UserMsgType, constant.HasReadReceipt, operationID)
//		s.Content = utils.StructToJsonString(newMessageIDList)
//		options := make(map[string]bool, 5)
//		utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
//		utils.SetSwitchFromOptions(options, constant.IsSenderConversationUpdate, false)
//		utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
//		utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
//		//If there is an error, the coroutine ends, so judgment is not  required
//		resp, _ := c.InternalSendMessage(callback, &s, userID, "", operationID, &server_api_params.OfflinePushInfo{}, false, options)
//		s.ServerMsgID = resp.ServerMsgID
//		s.SendTime = resp.SendTime
//		s.Status = constant.MsgStatusFiltered
//		msgStructToLocalChatLog(&localMessage, &s)
//		err = c.db.InsertMessage(&localMessage)
//		if err != nil {
//			log.Error(operationID, "inset into chat log err", localMessage, s, err.Error())
//		}
//		err2 := c.db.UpdateMessageHasRead(userID, newMessageIDList, constant.SingleChatType)
//		if err2 != nil {
//			log.Error(operationID, "update message has read error", newMessageIDList, userID, err2.Error())
//		}
//		_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.UpdateLatestMessageChange}, c.ch)
//		//_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.ch)
//	}
func (c *Conversation) insertMessageToLocalStorage(callback open_im_sdk_callback.Base, s *model_struct.LocalChatLog, operationID string) string {
	err := c.db.InsertMessageController(s)
	common.CheckDBErrCallback(callback, err, operationID)
	return s.ClientMsgID
}

func (c *Conversation) clearGroupHistoryMessage(callback open_im_sdk_callback.Base, groupID string, operationID string) {
	_, sessionType, err := c.getConversationTypeByGroupID(groupID)
	common.CheckAnyErrCallback(callback, 202, err, operationID)

	conversationID := utils.GetConversationIDBySessionType(groupID, int(sessionType))
	switch sessionType {
	case constant.SuperGroupChatType:
		err = c.db.SuperGroupDeleteAllMessage(groupID)
		common.CheckDBErrCallback(callback, err, operationID)
	default:
		err = c.db.UpdateMessageStatusBySourceIDController(groupID, constant.MsgStatusHasDeleted, sessionType)
		common.CheckDBErrCallback(callback, err, operationID)
	}

	err = c.db.ClearConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())

}

func (c *Conversation) clearC2CHistoryMessage(callback open_im_sdk_callback.Base, userID string, operationID string) {
	conversationID := utils.GetConversationIDBySessionType(userID, constant.SingleChatType)
	err := c.db.UpdateMessageStatusBySourceID(userID, constant.MsgStatusHasDeleted, constant.SingleChatType)
	common.CheckDBErrCallback(callback, err, operationID)
	err = c.db.ClearConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())
}

func (c *Conversation) deleteMessageFromSvr(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, operationID string) {
	seq, err := c.db.GetMsgSeqByClientMsgIDController(s)
	common.CheckDBErrCallback(callback, err, operationID)
	switch s.SessionType {
	case constant.SingleChatType, constant.GroupChatType:
		var apiReq server_api_params.DeleteMsgReq
		apiReq.SeqList = []uint32{seq}
		apiReq.OpUserID = c.loginUserID
		apiReq.UserID = c.loginUserID
		apiReq.OperationID = operationID
		c.p.PostFatalCallback(callback, constant.DeleteMsgRouter, apiReq, nil, apiReq.OperationID)
	case constant.SuperGroupChatType:
		var apiReq server_api_params.DelSuperGroupMsgReq
		apiReq.UserID = c.loginUserID
		apiReq.IsAllDelete = false
		apiReq.GroupID = s.GroupID
		apiReq.OperationID = operationID
		apiReq.SeqList = []uint32{seq}
		c.p.PostFatalCallback(callback, constant.DeleteSuperGroupMsgRouter, apiReq, nil, apiReq.OperationID)
		return
	}

}

func (c *Conversation) clearMessageFromSvr(callback open_im_sdk_callback.Base, operationID string) {
	var apiReq server_api_params.CleanUpMsgReq
	apiReq.UserID = c.loginUserID
	apiReq.OperationID = operationID
	c.p.PostFatalCallback(callback, constant.ClearMsgRouter, apiReq, nil, apiReq.OperationID)
	groupIDList, err := c.full.GetReadDiffusionGroupIDList(operationID)
	common.CheckDBErrCallback(callback, err, operationID)
	var superGroupApiReq server_api_params.DelSuperGroupMsgReq
	superGroupApiReq.UserID = c.loginUserID
	superGroupApiReq.IsAllDelete = true
	for _, v := range groupIDList {
		superGroupApiReq.GroupID = v
		superGroupApiReq.OperationID = operationID
		c.p.PostFatalCallback(callback, constant.DeleteSuperGroupMsgRouter, superGroupApiReq, nil, apiReq.OperationID)
	}
}

func (c *Conversation) deleteMessageFromLocalStorage(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, operationID string) {
	var conversation model_struct.LocalConversation
	var latestMsg sdk_struct.MsgStruct
	var conversationID string
	var sourceID string
	chatLog := model_struct.LocalChatLog{ClientMsgID: s.ClientMsgID, Status: constant.MsgStatusHasDeleted, SessionType: s.SessionType}

	switch s.SessionType {
	case constant.GroupChatType:
		conversationID = utils.GetConversationIDBySessionType(s.GroupID, constant.GroupChatType)
		sourceID = s.GroupID
	case constant.SingleChatType:
		if s.SendID != c.loginUserID {
			conversationID = utils.GetConversationIDBySessionType(s.SendID, constant.SingleChatType)
			sourceID = s.SendID
		} else {
			conversationID = utils.GetConversationIDBySessionType(s.RecvID, constant.SingleChatType)
			sourceID = s.RecvID
		}
	case constant.SuperGroupChatType:
		conversationID = utils.GetConversationIDBySessionType(s.GroupID, constant.SuperGroupChatType)
		sourceID = s.GroupID
		chatLog.RecvID = s.GroupID
	}
	err := c.db.UpdateMessageController(&chatLog)
	common.CheckDBErrCallback(callback, err, operationID)
	LocalConversation, err := c.db.GetConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	common.JsonUnmarshalCallback(LocalConversation.LatestMsg, &latestMsg, callback, operationID)

	if s.ClientMsgID == latestMsg.ClientMsgID { //If the deleted message is the latest message of the conversation, update the latest message of the conversation
		list, err := c.db.GetMessageListNoTimeController(sourceID, int(s.SessionType), 1, false)
		common.CheckDBErrCallback(callback, err, operationID)

		conversation.ConversationID = conversationID
		if list == nil {
			conversation.LatestMsg = ""
			conversation.LatestMsgSendTime = s.SendTime
		} else {
			copier.Copy(&latestMsg, list[0])
			err := c.msgConvert(&latestMsg)
			if err != nil {
				log.Error(operationID, "Parsing data error:", err.Error(), latestMsg)
			}
			conversation.LatestMsg = utils.StructToJsonString(latestMsg)
			conversation.LatestMsgSendTime = latestMsg.SendTime
		}
		err = c.db.UpdateColumnsConversation(conversation.ConversationID, map[string]interface{}{"latest_msg_send_time": conversation.LatestMsgSendTime, "latest_msg": conversation.LatestMsg})
		if err != nil {
			log.Error("internal", "updateConversationLatestMsgModel err: ", err)
		} else {
			_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())
		}
	}
}
func (c *Conversation) judgeMultipleSubString(keywordList []string, main string, keywordListMatchType int) bool {
	if len(keywordList) == 0 {
		return true
	}
	if keywordListMatchType == constant.KeywordMatchOr {
		for _, v := range keywordList {
			if utils.KMP(main, v) {
				return true
			}
		}
		return false
	} else {
		for _, v := range keywordList {
			if !utils.KMP(main, v) {
				return false
			}
		}
	}
	return true
}

func (c *Conversation) searchLocalMessages(callback open_im_sdk_callback.Base, searchParam sdk.SearchLocalMessagesParams, operationID string) (r sdk.SearchLocalMessagesCallback) {

	var conversationID, sourceID string
	var startTime, endTime int64
	var list []*model_struct.LocalChatLog
	conversationMap := make(map[string]*sdk.SearchByConversationResult, 10)
	var err error

	if searchParam.SearchTimePosition == 0 {
		endTime = utils.GetCurrentTimestampBySecond()
	} else {
		endTime = searchParam.SearchTimePosition
	}
	if searchParam.SearchTimePeriod != 0 {
		startTime = endTime - searchParam.SearchTimePeriod
	}
	startTime = utils.UnixSecondToTime(startTime).UnixNano() / 1e6
	endTime = utils.UnixSecondToTime(endTime).UnixNano() / 1e6
	if len(searchParam.KeywordList) == 0 && len(searchParam.MessageTypeList) == 0 {
		common.CheckAnyErrCallback(callback, 201, errors.New("keywordlist and messageTypelist all null"), operationID)
	}
	if searchParam.ConversationID != "" {
		if searchParam.PageIndex < 1 || searchParam.Count < 1 {
			common.CheckAnyErrCallback(callback, 201, errors.New("page or count is null"), operationID)
		}
		offset := (searchParam.PageIndex - 1) * searchParam.Count
		localConversation, err := c.db.GetConversation(searchParam.ConversationID)
		common.CheckDBErrCallback(callback, err, operationID)
		switch localConversation.ConversationType {
		case constant.SingleChatType:
			sourceID = localConversation.UserID
		case constant.GroupChatType:
			sourceID = localConversation.GroupID
		case constant.SuperGroupChatType:
			sourceID = localConversation.GroupID
		}
		if len(searchParam.MessageTypeList) != 0 && len(searchParam.KeywordList) == 0 {
			list, err = c.db.SearchMessageByContentTypeController(searchParam.MessageTypeList, sourceID, startTime, endTime, int(localConversation.ConversationType), offset, searchParam.Count)
		} else {
			newContentTypeList := func(list []int) (result []int) {
				for _, v := range list {
					if utils.IsContainInt(v, SearchContentType) {
						result = append(result, v)
					}
				}
				return result
			}(searchParam.MessageTypeList)
			if len(newContentTypeList) == 0 {
				newContentTypeList = SearchContentType
			}
			list, err = c.db.SearchMessageByKeywordController(newContentTypeList, searchParam.KeywordList, searchParam.KeywordListMatchType, sourceID, startTime, endTime, int(localConversation.ConversationType), offset, searchParam.Count)
		}
	} else {
		//Comprehensive search, search all
		if len(searchParam.MessageTypeList) == 0 {
			searchParam.MessageTypeList = SearchContentType
		}
		list, err = c.db.SearchMessageByContentTypeAndKeywordController(searchParam.MessageTypeList, searchParam.KeywordList, searchParam.KeywordListMatchType, startTime, endTime, operationID)
	}
	common.CheckDBErrCallback(callback, err, operationID)
	//localChatLogToMsgStruct(&messageList, list)

	//log.Debug("hahh",utils.KMP("SSSsdf3434","s"))
	//log.Debug("hahh",utils.KMP("SSSsdf3434","g"))
	//log.Debug("hahh",utils.KMP("SSSsdf3434","3434"))
	//log.Debug("hahh",utils.KMP("SSSsdf3434","F3434"))
	//log.Debug("hahh",utils.KMP("SSSsdf3434","SDF3"))
	log.Debug(operationID, "get raw data length is", len(list))
	for _, v := range list {
		temp := sdk_struct.MsgStruct{}
		temp.ClientMsgID = v.ClientMsgID
		temp.ServerMsgID = v.ServerMsgID
		temp.CreateTime = v.CreateTime
		temp.SendTime = v.SendTime
		temp.SessionType = v.SessionType
		temp.SendID = v.SendID
		temp.RecvID = v.RecvID
		temp.MsgFrom = v.MsgFrom
		temp.ContentType = v.ContentType
		temp.SenderPlatformID = v.SenderPlatformID
		temp.SenderNickname = v.SenderNickname
		temp.SenderFaceURL = v.SenderFaceURL
		temp.Content = v.Content
		temp.Seq = v.Seq
		temp.IsRead = v.IsRead
		temp.Status = v.Status
		temp.AttachedInfo = v.AttachedInfo
		temp.Ex = v.Ex
		err := c.msgHandleByContentType(&temp)
		if err != nil {
			log.Error(operationID, "Parsing data error:", err.Error(), temp)
			continue
		}
		if temp.ContentType == constant.File && !c.judgeMultipleSubString(searchParam.KeywordList, temp.FileElem.FileName, searchParam.KeywordListMatchType) {
			continue
		}
		if temp.ContentType == constant.AtText && !c.judgeMultipleSubString(searchParam.KeywordList, temp.AtElem.Text, searchParam.KeywordListMatchType) {
			continue
		}
		switch temp.SessionType {
		case constant.SingleChatType:
			if temp.SendID == c.loginUserID {
				conversationID = utils.GetConversationIDBySessionType(temp.RecvID, constant.SingleChatType)
			} else {
				conversationID = utils.GetConversationIDBySessionType(temp.SendID, constant.SingleChatType)
			}
		case constant.GroupChatType:
			temp.GroupID = temp.RecvID
			temp.RecvID = c.loginUserID
			conversationID = utils.GetConversationIDBySessionType(temp.GroupID, constant.GroupChatType)
		case constant.SuperGroupChatType:
			temp.GroupID = temp.RecvID
			temp.RecvID = c.loginUserID
			conversationID = utils.GetConversationIDBySessionType(temp.GroupID, constant.SuperGroupChatType)
		}
		if oldItem, ok := conversationMap[conversationID]; !ok {
			searchResultItem := sdk.SearchByConversationResult{}
			localConversation, err := c.db.GetConversation(conversationID)
			if err != nil {
				log.Error(operationID, "get conversation err ", err.Error(), conversationID)
				continue
			}
			searchResultItem.ConversationID = conversationID
			searchResultItem.FaceURL = localConversation.FaceURL
			searchResultItem.ShowName = localConversation.ShowName
			searchResultItem.ConversationType = localConversation.ConversationType
			searchResultItem.MessageList = append(searchResultItem.MessageList, &temp)
			searchResultItem.MessageCount++
			conversationMap[conversationID] = &searchResultItem
		} else {
			oldItem.MessageCount++
			oldItem.MessageList = append(oldItem.MessageList, &temp)
			conversationMap[conversationID] = oldItem
		}
	}
	for _, v := range conversationMap {
		r.SearchResultItems = append(r.SearchResultItems, v)
		r.TotalCount += v.MessageCount

	}
	return r
}

func (c *Conversation) setConversationNotification(msg *server_api_params.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	c.SyncConversations(operationID, 0)
}

func (c *Conversation) DoNotification(msg *server_api_params.MsgData) {
	if msg.SendTime < c.full.Group().LoginTime() || c.full.Group().LoginTime() == 0 {
		log.Warn("", "ignore notification ", msg.ClientMsgID, msg.ServerMsgID, msg.Seq, msg.ContentType)
		return
	}
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg)
	if c.msgListener == nil {
		log.Error(operationID, utils.GetSelfFuncName(), "listener == nil")
		return
	}
	go func() {
		c.setConversationNotification(msg, operationID)
	}()
}

func (c *Conversation) delMsgBySeq(seqList []uint32) error {
	var SPLIT = 1000
	for i := 0; i < len(seqList)/SPLIT; i++ {
		if err := c.delMsgBySeqSplit(seqList[i*SPLIT : (i+1)*SPLIT]); err != nil {
			return utils.Wrap(err, "")
		}
	}
	return nil
}

func (c *Conversation) delMsgBySeqSplit(seqList []uint32) error {
	var req server_api_params.DelMsgListReq
	req.SeqList = seqList
	req.OperationID = utils.OperationIDGenerator()
	req.OpUserID = c.loginUserID
	req.UserID = c.loginUserID
	operationID := req.OperationID

	resp, err := c.Ws.SendReqWaitResp(&req, constant.WsDelMsg, 30, 5, c.loginUserID, req.OperationID)
	if err != nil {
		return utils.Wrap(err, "SendReqWaitResp failed")
	}
	var delResp server_api_params.DelMsgListResp
	err = proto.Unmarshal(resp.Data, &delResp)
	if err != nil {
		log.Error(operationID, "Unmarshal failed ", err.Error())
		return utils.Wrap(err, "Unmarshal failed")
	}
	return nil
}

// old WS method
//func (c *Conversation) deleteMessageFromSvr(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, operationID string) {
//	seq, err := c.db.GetMsgSeqByClientMsgID(s.ClientMsgID)
//	common.CheckDBErrCallback(callback, err, operationID)
//	if seq == 0 {
//		err = errors.New("seq == 0 ")
//		common.CheckArgsErrCallback(callback, err, operationID)
//	}
//	seqList := []uint32{seq}
//	err = c.delMsgBySeq(seqList)
//	common.CheckArgsErrCallback(callback, err, operationID)
//}

func (c *Conversation) deleteConversationAndMsgFromSvr(callback open_im_sdk_callback.Base, conversationID, operationID string) {
	local, err := c.db.GetConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	log.Debug(operationID, utils.GetSelfFuncName(), *local)
	var seqList []uint32
	switch local.ConversationType {
	case constant.SingleChatType, constant.NotificationChatType:
		peerUserID := local.UserID
		if peerUserID != c.loginUserID {
			seqList, err = c.db.GetMsgSeqListByPeerUserID(peerUserID)
		} else {
			seqList, err = c.db.GetMsgSeqListBySelfUserID(c.loginUserID)
		}
		log.NewDebug(operationID, utils.GetSelfFuncName(), "seqList: ", seqList)
		common.CheckDBErrCallback(callback, err, operationID)
	case constant.GroupChatType:
		groupID := local.GroupID
		seqList, err = c.db.GetMsgSeqListByGroupID(groupID)
		log.NewDebug(operationID, utils.GetSelfFuncName(), "seqList: ", seqList)
		common.CheckDBErrCallback(callback, err, operationID)
	case constant.SuperGroupChatType:
		var apiReq server_api_params.DelSuperGroupMsgReq
		apiReq.UserID = c.loginUserID
		apiReq.IsAllDelete = true
		apiReq.GroupID = local.GroupID
		apiReq.OperationID = operationID
		c.p.PostFatalCallback(callback, constant.DeleteSuperGroupMsgRouter, apiReq, nil, apiReq.OperationID)
		return

	}
	var apiReq server_api_params.DeleteMsgReq
	apiReq.OpUserID = c.loginUserID
	apiReq.UserID = c.loginUserID
	apiReq.OperationID = operationID
	apiReq.SeqList = seqList
	c.p.PostFatalCallback(callback, constant.DeleteMsgRouter, apiReq, nil, apiReq.OperationID)
	common.CheckArgsErrCallback(callback, err, operationID)
}

func (c *Conversation) deleteAllMsgFromLocal(callback open_im_sdk_callback.Base, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	err := c.db.DeleteAllMessage()
	common.CheckDBErrCallback(callback, err, operationID)
	groupIDList, err := c.full.GetReadDiffusionGroupIDList(operationID)
	common.CheckDBErrCallback(callback, err, operationID)
	for _, v := range groupIDList {
		err = c.db.SuperGroupDeleteAllMessage(v)
		if err != nil {
			log.Error(operationID, "SuperGroupDeleteAllMessage err", err.Error())
			continue
		}
	}
	err = c.db.ClearAllConversation()
	common.CheckDBErrCallback(callback, err, operationID)
	conversationList, err := c.db.GetAllConversationListDB()
	common.CheckDBErrCallback(callback, err, operationID)
	var cidList []string
	for _, conversation := range conversationList {
		cidList = append(cidList, conversation.ConversationID)
	}
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{Action: constant.ConChange, Args: cidList}, c.GetCh())
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{"", constant.TotalUnreadMessageChanged, ""}})

}

func (c *Conversation) deleteAllMsgFromSvr(callback open_im_sdk_callback.Base, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	seqList, err := c.db.GetAllUnDeleteMessageSeqList()
	log.NewInfo(operationID, utils.GetSelfFuncName(), seqList)
	common.CheckDBErrCallback(callback, err, operationID)
	var apiReq server_api_params.DeleteMsgReq
	apiReq.OpUserID = c.loginUserID
	apiReq.UserID = c.loginUserID
	apiReq.OperationID = operationID
	apiReq.SeqList = seqList
	c.p.PostFatalCallback(callback, constant.DeleteMsgRouter, apiReq, nil, apiReq.OperationID)
}
func (c *Conversation) setMessageReactionExtensions(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, req sdk.SetMessageReactionExtensionsParams, operationID string) []*server_api_params.ExtensionResult {
	message, err := c.db.GetMessageController(s)
	common.CheckDBErrCallback(callback, err, operationID)
	if message.Status != constant.MsgStatusSendSuccess {
		common.CheckAnyErrCallback(callback, 201, errors.New("only send success message can modify reaction extensions"), operationID)
	}
	if message.SessionType != constant.SuperGroupChatType {
		common.CheckAnyErrCallback(callback, 202, errors.New("currently only support super group message"), operationID)

	}
	extendMsg, _ := c.db.GetMessageReactionExtension(message.ClientMsgID)
	temp := make(map[string]*server_api_params.KeyValue)
	_ = json.Unmarshal(extendMsg.LocalReactionExtensions, &temp)
	reqTemp := make(map[string]*server_api_params.KeyValue)
	for _, v := range req {
		if value, ok := temp[v.TypeKey]; ok {
			v.LatestUpdateTime = value.LatestUpdateTime
		}
		reqTemp[v.TypeKey] = v
	}
	var sourceID string
	switch message.SessionType {
	case constant.SingleChatType:
		if message.SendID == c.loginUserID {
			sourceID = message.RecvID
		} else {
			sourceID = message.SendID
		}
	case constant.NotificationChatType:
		sourceID = message.RecvID
	case constant.GroupChatType, constant.SuperGroupChatType:
		sourceID = message.RecvID
	}
	var apiReq server_api_params.SetMessageReactionExtensionsReq
	apiReq.IsReact = message.IsReact
	apiReq.ClientMsgID = message.ClientMsgID
	apiReq.SourceID = sourceID
	apiReq.SessionType = message.SessionType
	apiReq.IsExternalExtensions = message.IsExternalExtensions
	apiReq.ReactionExtensionList = reqTemp
	apiReq.OperationID = operationID
	apiReq.MsgFirstModifyTime = message.MsgFirstModifyTime
	var apiResp server_api_params.SetMessageReactionExtensionsResp
	c.p.PostFatalCallback(callback, constant.SetMessageReactionExtensionsRouter, apiReq, &apiResp.ApiResult, apiReq.OperationID)
	var msg model_struct.LocalChatLogReactionExtensions
	msg.ClientMsgID = message.ClientMsgID
	resultKeyMap := make(map[string]*server_api_params.KeyValue)
	for _, v := range apiResp.ApiResult.Result {
		if v.ErrCode == 0 {
			temp := new(server_api_params.KeyValue)
			temp.TypeKey = v.TypeKey
			temp.Value = v.Value
			temp.LatestUpdateTime = v.LatestUpdateTime
			resultKeyMap[v.TypeKey] = temp
		}
	}
	err = c.db.GetAndUpdateMessageReactionExtension(message.ClientMsgID, resultKeyMap)
	if err != nil {
		log.Error(operationID, "GetAndUpdateMessageReactionExtension err:", err.Error())
	}
	if !message.IsReact {
		message.IsReact = apiResp.ApiResult.IsReact
		message.MsgFirstModifyTime = apiResp.ApiResult.MsgFirstModifyTime
		err = c.db.UpdateMessageController(message)
		if err != nil {
			log.Error(operationID, "UpdateMessageController err:", err.Error(), message)

		}
	}
	return apiResp.ApiResult.Result
}
func (c *Conversation) addMessageReactionExtensions(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, req sdk.AddMessageReactionExtensionsParams, operationID string) []*server_api_params.ExtensionResult {
	message, err := c.db.GetMessageController(s)
	common.CheckDBErrCallback(callback, err, operationID)
	if message.Status != constant.MsgStatusSendSuccess || message.Seq == 0 {
		common.CheckAnyErrCallback(callback, 201, errors.New("only send success message can modify reaction extensions"), operationID)
	}
	//if !message.IsExternalExtensions {
	//	common.CheckAnyErrCallback(callback, 202, errors.New(" only externalExtensions message can use this interface"), operationID)
	//
	//}
	reqTemp := make(map[string]*server_api_params.KeyValue)
	extendMsg, err := c.db.GetMessageReactionExtension(message.ClientMsgID)
	if err == nil && extendMsg != nil {
		temp := make(map[string]*server_api_params.KeyValue)
		_ = json.Unmarshal(extendMsg.LocalReactionExtensions, &temp)
		for _, v := range req {
			if value, ok := temp[v.TypeKey]; ok {
				v.LatestUpdateTime = value.LatestUpdateTime
			}
			reqTemp[v.TypeKey] = v
		}
	} else {
		for _, v := range req {
			reqTemp[v.TypeKey] = v
		}
	}
	var sourceID string
	switch message.SessionType {
	case constant.SingleChatType:
		if message.SendID == c.loginUserID {
			sourceID = message.RecvID
		} else {
			sourceID = message.SendID
		}
	case constant.NotificationChatType:
		sourceID = message.RecvID
	case constant.GroupChatType, constant.SuperGroupChatType:
		sourceID = message.RecvID
	}
	var apiReq server_api_params.AddMessageReactionExtensionsReq
	apiReq.IsReact = message.IsReact
	apiReq.ClientMsgID = message.ClientMsgID
	apiReq.SourceID = sourceID
	apiReq.SessionType = message.SessionType
	apiReq.IsExternalExtensions = message.IsExternalExtensions
	apiReq.ReactionExtensionList = reqTemp
	apiReq.OperationID = operationID
	apiReq.MsgFirstModifyTime = message.MsgFirstModifyTime
	apiReq.Seq = message.Seq
	var apiResp server_api_params.AddMessageReactionExtensionsResp
	c.p.PostFatalCallbackPenetrate(callback, constant.AddMessageReactionExtensionsRouter, apiReq, &apiResp.ApiResult, apiReq.OperationID)
	log.Debug(operationID, "api return:", message.IsReact, apiResp.ApiResult)
	if !message.IsReact {
		message.IsReact = apiResp.ApiResult.IsReact
		message.MsgFirstModifyTime = apiResp.ApiResult.MsgFirstModifyTime
		err = c.db.UpdateMessageController(message)
		if err != nil {
			log.Error(operationID, "UpdateMessageController err:", err.Error(), message)
		}
	}
	return apiResp.ApiResult.Result
}

func (c *Conversation) deleteMessageReactionExtensions(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, req sdk.DeleteMessageReactionExtensionsParams, operationID string) []*server_api_params.ExtensionResult {
	message, err := c.db.GetMessageController(s)
	common.CheckDBErrCallback(callback, err, operationID)
	if message.Status != constant.MsgStatusSendSuccess {
		common.CheckAnyErrCallback(callback, 201, errors.New("only send success message can modify reaction extensions"), operationID)
	}
	if message.SessionType != constant.SuperGroupChatType {
		common.CheckAnyErrCallback(callback, 202, errors.New("currently only support super group message"), operationID)

	}
	extendMsg, _ := c.db.GetMessageReactionExtension(message.ClientMsgID)
	temp := make(map[string]*server_api_params.KeyValue)
	_ = json.Unmarshal(extendMsg.LocalReactionExtensions, &temp)
	var reqTemp []*server_api_params.KeyValue
	for _, v := range req {
		if value, ok := temp[v]; ok {
			var tt server_api_params.KeyValue
			tt.LatestUpdateTime = value.LatestUpdateTime
			tt.TypeKey = v
			reqTemp = append(reqTemp, &tt)
		}
	}
	var sourceID string
	switch message.SessionType {
	case constant.SingleChatType:
		if message.SendID == c.loginUserID {
			sourceID = message.RecvID
		} else {
			sourceID = message.SendID
		}
	case constant.NotificationChatType:
		sourceID = message.RecvID
	case constant.GroupChatType, constant.SuperGroupChatType:
		sourceID = message.RecvID
	}
	var apiReq server_api_params.DeleteMessageReactionExtensionsReq
	apiReq.ClientMsgID = message.ClientMsgID
	apiReq.SourceID = sourceID
	apiReq.SessionType = message.SessionType
	apiReq.ReactionExtensionList = reqTemp
	apiReq.OperationID = operationID
	apiReq.IsExternalExtensions = message.IsExternalExtensions
	apiReq.MsgFirstModifyTime = message.MsgFirstModifyTime
	var apiResp server_api_params.DeleteMessageReactionExtensionsResp
	c.p.PostFatalCallback(callback, constant.DeleteMessageReactionExtensionsRouter, apiReq, &apiResp.Result, apiReq.OperationID)
	var msg model_struct.LocalChatLogReactionExtensions
	msg.ClientMsgID = message.ClientMsgID
	resultKeyMap := make(map[string]*server_api_params.KeyValue)
	for _, v := range apiResp.Result {
		if v.ErrCode == 0 {
			temp := new(server_api_params.KeyValue)
			temp.TypeKey = v.TypeKey
			resultKeyMap[v.TypeKey] = temp
		}
	}
	err = c.db.DeleteAndUpdateMessageReactionExtension(message.ClientMsgID, resultKeyMap)
	if err != nil {
		log.Error(operationID, "GetAndUpdateMessageReactionExtension err:", err.Error())
	}
	return apiResp.Result
}

type syncReactionExtensionParams struct {
	MessageList         []*model_struct.LocalChatLog
	SessionType         int32
	SourceID            string
	IsExternalExtension bool
	ExtendMessageList   []*model_struct.LocalChatLogReactionExtensions
	TypeKeyList         []string
}

func (c *Conversation) getMessageListReactionExtensions(callback open_im_sdk_callback.Base, messageList []*sdk_struct.MsgStruct, operationID string) server_api_params.GetMessageListReactionExtensionsResp {
	if len(messageList) == 0 {
		common.CheckAnyErrCallback(callback, 201, errors.New("message list is null"), operationID)
	}
	var msgIDList []string
	var sourceID string
	var sessionType int32
	var isExternalExtension bool
	for _, msgStruct := range messageList {
		switch msgStruct.SessionType {
		case constant.SingleChatType:
			if msgStruct.SendID == c.loginUserID {
				sourceID = msgStruct.RecvID
			} else {
				sourceID = msgStruct.SendID
			}
		case constant.NotificationChatType:
			sourceID = msgStruct.RecvID
		case constant.GroupChatType, constant.SuperGroupChatType:
			sourceID = msgStruct.GroupID
		}
		sessionType = msgStruct.SessionType
		msgIDList = append(msgIDList, msgStruct.ClientMsgID)
	}
	isExternalExtension = c.IsExternalExtensions
	localMessageList, err := c.db.GetMultipleMessageController(msgIDList, sourceID, sessionType)
	common.CheckDBErrCallback(callback, err, operationID)
	for _, v := range localMessageList {
		if v.IsReact != true {
			common.CheckAnyErrCallback(callback, 208, errors.New("have not reaction message in message list:"+v.ClientMsgID), operationID)
		}
	}
	var result server_api_params.GetMessageListReactionExtensionsResp
	extendMessage, _ := c.db.GetMultipleMessageReactionExtension(msgIDList)
	for _, v := range extendMessage {
		var singleResult server_api_params.SingleMessageExtensionResult
		temp := make(map[string]*server_api_params.KeyValue)
		_ = json.Unmarshal(v.LocalReactionExtensions, &temp)
		singleResult.ClientMsgID = v.ClientMsgID
		singleResult.ReactionExtensionList = temp
		result = append(result, &singleResult)
	}
	args := syncReactionExtensionParams{}
	args.MessageList = localMessageList
	args.SourceID = sourceID
	args.SessionType = sessionType
	args.ExtendMessageList = extendMessage
	args.IsExternalExtension = isExternalExtension
	_ = common.TriggerCmdSyncReactionExtensions(common.SyncReactionExtensionsNode{
		OperationID: operationID,
		Action:      constant.SyncMessageListReactionExtensions,
		Args:        args,
	}, c.GetCh())
	return result

}

func (c *Conversation) getMessageListSomeReactionExtensions(callback open_im_sdk_callback.Base, messageList []*sdk_struct.MsgStruct, keyList []string, operationID string) server_api_params.GetMessageListReactionExtensionsResp {
	if len(messageList) == 0 {
		common.CheckAnyErrCallback(callback, 201, errors.New("message list is null"), operationID)
	}
	var msgIDList []string
	var sourceID string
	var sessionType int32
	var isExternalExtension bool
	for _, msgStruct := range messageList {
		switch msgStruct.SessionType {
		case constant.SingleChatType:
			if msgStruct.SendID == c.loginUserID {
				sourceID = msgStruct.RecvID
			} else {
				sourceID = msgStruct.SendID
			}
		case constant.NotificationChatType:
			sourceID = msgStruct.RecvID
		case constant.GroupChatType, constant.SuperGroupChatType:
			sourceID = msgStruct.GroupID
		}
		sessionType = msgStruct.SessionType
		isExternalExtension = msgStruct.IsExternalExtensions
		msgIDList = append(msgIDList, msgStruct.ClientMsgID)
	}
	localMessageList, err := c.db.GetMultipleMessageController(msgIDList, sourceID, sessionType)
	common.CheckDBErrCallback(callback, err, operationID)
	var result server_api_params.GetMessageListReactionExtensionsResp
	extendMsgs, _ := c.db.GetMultipleMessageReactionExtension(msgIDList)
	for _, v := range extendMsgs {
		var singleResult server_api_params.SingleMessageExtensionResult
		temp := make(map[string]*server_api_params.KeyValue)
		_ = json.Unmarshal(v.LocalReactionExtensions, &temp)
		for s, _ := range temp {
			if !utils.IsContain(s, keyList) {
				delete(temp, s)
			}
		}
		singleResult.ClientMsgID = v.ClientMsgID
		singleResult.ReactionExtensionList = temp
		result = append(result, &singleResult)
	}
	args := syncReactionExtensionParams{}
	args.MessageList = localMessageList
	args.SourceID = sourceID
	args.TypeKeyList = keyList
	args.SessionType = sessionType
	args.ExtendMessageList = extendMsgs
	args.IsExternalExtension = isExternalExtension
	_ = common.TriggerCmdSyncReactionExtensions(common.SyncReactionExtensionsNode{
		OperationID: operationID,
		Action:      constant.SyncMessageListReactionExtensions,
		Args:        args,
	}, c.GetCh())
	return result
}

func (c *Conversation) setTypeKeyInfo(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, typeKey, ex string, isCanRepeat bool, operationID string) []*server_api_params.ExtensionResult {
	message, err := c.db.GetMessageController(s)
	common.CheckDBErrCallback(callback, err, operationID)
	if message.Status != constant.MsgStatusSendSuccess {
		common.CheckAnyErrCallback(callback, 201, errors.New("only send success message can modify reaction extensions"), operationID)
	}
	extendMsg, _ := c.db.GetMessageReactionExtension(message.ClientMsgID)
	temp := make(map[string]*server_api_params.KeyValue)
	_ = json.Unmarshal(extendMsg.LocalReactionExtensions, &temp)
	var flag bool
	var isContainSelfK string
	var dbIsCanRepeat bool
	var deletedKeyValue server_api_params.KeyValue
	var maxTypeKey string
	var maxTypeKeyValue server_api_params.KeyValue
	reqTemp := make(map[string]*server_api_params.KeyValue)
	for k, v := range temp {
		if strings.HasPrefix(k, typeKey) {
			flag = true
			singleTypeKeyInfo := new(sdk.SingleTypeKeyInfo)
			_ = json.Unmarshal([]byte(v.Value), singleTypeKeyInfo)
			if _, ok := singleTypeKeyInfo.InfoList[c.loginUserID]; ok {
				isContainSelfK = k
				dbIsCanRepeat = singleTypeKeyInfo.IsCanRepeat
				delete(singleTypeKeyInfo.InfoList, c.loginUserID)
				singleTypeKeyInfo.Counter--
				deletedKeyValue.TypeKey = v.TypeKey
				deletedKeyValue.Value = utils.StructToJsonString(singleTypeKeyInfo)
				deletedKeyValue.LatestUpdateTime = v.LatestUpdateTime
			}
			if k > maxTypeKey {
				maxTypeKey = k
				maxTypeKeyValue = *v
			}
		}
	}
	if !flag {
		if len(temp) >= 300 {
			common.CheckAnyErrCallback(callback, 202, errors.New("number of keys only can support 300"), operationID)
		}
		singleTypeKeyInfo := new(sdk.SingleTypeKeyInfo)
		singleTypeKeyInfo.TypeKey = getIndexTypeKey(typeKey, 0)
		singleTypeKeyInfo.Counter = 1
		singleTypeKeyInfo.IsCanRepeat = isCanRepeat
		singleTypeKeyInfo.Index = 0
		userInfo := new(sdk.Info)
		userInfo.UserID = c.loginUserID
		userInfo.Ex = ex
		singleTypeKeyInfo.InfoList[c.loginUserID] = userInfo
		keyValue := new(server_api_params.KeyValue)
		keyValue.TypeKey = singleTypeKeyInfo.TypeKey
		keyValue.Value = utils.StructToJsonString(singleTypeKeyInfo)
		reqTemp[singleTypeKeyInfo.TypeKey] = keyValue
	} else {
		if isContainSelfK != "" && !dbIsCanRepeat {
			//删除操作
			reqTemp[isContainSelfK] = &deletedKeyValue
		} else {
			singleTypeKeyInfo := new(sdk.SingleTypeKeyInfo)
			_ = json.Unmarshal([]byte(maxTypeKeyValue.Value), singleTypeKeyInfo)
			userInfo := new(sdk.Info)
			userInfo.UserID = c.loginUserID
			userInfo.Ex = ex
			singleTypeKeyInfo.Counter++
			singleTypeKeyInfo.InfoList[c.loginUserID] = userInfo
			maxTypeKeyValue.Value = utils.StructToJsonString(singleTypeKeyInfo)
			data, _ := json.Marshal(maxTypeKeyValue)
			if len(data) > 1000 { //单key超过了1kb
				if len(temp) >= 300 {
					common.CheckAnyErrCallback(callback, 202, errors.New("number of keys only can support 300"), operationID)
				}
				newSingleTypeKeyInfo := new(sdk.SingleTypeKeyInfo)
				newSingleTypeKeyInfo.TypeKey = getIndexTypeKey(typeKey, singleTypeKeyInfo.Index+1)
				newSingleTypeKeyInfo.Counter = 1
				newSingleTypeKeyInfo.IsCanRepeat = singleTypeKeyInfo.IsCanRepeat
				newSingleTypeKeyInfo.Index = singleTypeKeyInfo.Index + 1
				userInfo := new(sdk.Info)
				userInfo.UserID = c.loginUserID
				userInfo.Ex = ex
				newSingleTypeKeyInfo.InfoList[c.loginUserID] = userInfo
				keyValue := new(server_api_params.KeyValue)
				keyValue.TypeKey = newSingleTypeKeyInfo.TypeKey
				keyValue.Value = utils.StructToJsonString(newSingleTypeKeyInfo)
				reqTemp[singleTypeKeyInfo.TypeKey] = keyValue
			} else {
				reqTemp[maxTypeKey] = &maxTypeKeyValue
			}

		}
	}
	var sourceID string
	switch message.SessionType {
	case constant.SingleChatType:
		sourceID = message.SendID + message.RecvID
	case constant.NotificationChatType:
		sourceID = message.RecvID
	case constant.GroupChatType, constant.SuperGroupChatType:
		sourceID = message.RecvID
	}
	var apiReq server_api_params.SetMessageReactionExtensionsReq
	apiReq.IsReact = message.IsReact
	apiReq.ClientMsgID = message.ClientMsgID
	apiReq.SourceID = sourceID
	apiReq.SessionType = message.SessionType
	apiReq.IsExternalExtensions = message.IsExternalExtensions
	apiReq.ReactionExtensionList = reqTemp
	apiReq.OperationID = operationID
	apiReq.MsgFirstModifyTime = message.MsgFirstModifyTime
	var apiResp server_api_params.SetMessageReactionExtensionsResp
	c.p.PostFatalCallback(callback, constant.SetMessageReactionExtensionsRouter, apiReq, &apiResp.ApiResult, apiReq.OperationID)
	var msg model_struct.LocalChatLogReactionExtensions
	msg.ClientMsgID = message.ClientMsgID
	resultKeyMap := make(map[string]*server_api_params.KeyValue)
	for _, v := range apiResp.ApiResult.Result {
		if v.ErrCode == 0 {
			temp := new(server_api_params.KeyValue)
			temp.TypeKey = v.TypeKey
			temp.Value = v.Value
			temp.LatestUpdateTime = v.LatestUpdateTime
			resultKeyMap[v.TypeKey] = temp
		}
	}
	err = c.db.GetAndUpdateMessageReactionExtension(message.ClientMsgID, resultKeyMap)
	if err != nil {
		log.Error(operationID, "GetAndUpdateMessageReactionExtension err:", err.Error())
	}
	if !message.IsReact {
		message.IsReact = apiResp.ApiResult.IsReact
		message.MsgFirstModifyTime = apiResp.ApiResult.MsgFirstModifyTime
		err = c.db.UpdateMessageController(message)
		if err != nil {
			log.Error(operationID, "UpdateMessageController err:", err.Error(), message)

		}
	}
	return apiResp.ApiResult.Result
}
func getIndexTypeKey(typeKey string, index int) string {
	return typeKey + "$" + utils.IntToString(index)
}
func getPrefixTypeKey(typeKey string) string {
	list := strings.Split(typeKey, "$")
	if len(list) > 0 {
		return list[0]
	}
	return ""
}
func (c *Conversation) getTypeKeyListInfo(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, keyList []string, operationID string) (result []*sdk.SingleTypeKeyInfoSum) {
	message, err := c.db.GetMessageController(s)
	common.CheckDBErrCallback(callback, err, operationID)
	if message.Status != constant.MsgStatusSendSuccess {
		common.CheckAnyErrCallback(callback, 201, errors.New("only send success message can modify reaction extensions"), operationID)
	}
	if !message.IsReact {
		common.CheckAnyErrCallback(callback, 202, errors.New("can get message reaction ex"), operationID)
	}
	extendMsg, _ := c.db.GetMessageReactionExtension(message.ClientMsgID)
	temp := make(map[string]*server_api_params.KeyValue)
	_ = json.Unmarshal(extendMsg.LocalReactionExtensions, &temp)
	for _, v := range keyList {
		singleResult := new(sdk.SingleTypeKeyInfoSum)
		singleResult.TypeKey = v
		for typeKey, value := range temp {
			if strings.HasPrefix(typeKey, v) {
				singleTypeKeyInfo := new(sdk.SingleTypeKeyInfo)
				_ = json.Unmarshal([]byte(value.Value), singleTypeKeyInfo)
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
	messageList := []*sdk_struct.MsgStruct{s}
	_ = common.TriggerCmdSyncReactionExtensions(common.SyncReactionExtensionsNode{
		OperationID: operationID,
		Action:      constant.SyncMessageListTypeKeyInfo,
		Args:        messageList,
	}, c.GetCh())

	return result
}

func (c *Conversation) getAllTypeKeyInfo(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, operationID string) (result []*sdk.SingleTypeKeyInfoSum) {
	message, err := c.db.GetMessageController(s)
	common.CheckDBErrCallback(callback, err, operationID)
	if message.Status != constant.MsgStatusSendSuccess {
		common.CheckAnyErrCallback(callback, 201, errors.New("only send success message can modify reaction extensions"), operationID)
	}
	if !message.IsReact {
		common.CheckAnyErrCallback(callback, 202, errors.New("can get message reaction ex"), operationID)
	}
	extendMsg, _ := c.db.GetMessageReactionExtension(message.ClientMsgID)
	temp := make(map[string]*server_api_params.KeyValue)
	_ = json.Unmarshal(extendMsg.LocalReactionExtensions, &temp)
	mapResult := make(map[string]*sdk.SingleTypeKeyInfoSum)
	for typeKey, value := range temp {
		singleTypeKeyInfo := new(sdk.SingleTypeKeyInfo)
		err := json.Unmarshal([]byte(value.Value), singleTypeKeyInfo)
		if err != nil {
			log.Warn(operationID, "not this type ", value.Value)
			continue
		}
		prefixKey := getPrefixTypeKey(typeKey)
		if v, ok := mapResult[prefixKey]; ok {
			for _, info := range singleTypeKeyInfo.InfoList {
				t := *info
				v.InfoList = append(v.InfoList, &t)
			}
			if _, ok := singleTypeKeyInfo.InfoList[c.loginUserID]; ok {
				v.IsContainSelf = true
			}
			v.Counter += singleTypeKeyInfo.Counter
		} else {
			v := new(sdk.SingleTypeKeyInfoSum)
			v.TypeKey = prefixKey
			v.Counter = singleTypeKeyInfo.Counter
			for _, info := range singleTypeKeyInfo.InfoList {
				t := *info
				v.InfoList = append(v.InfoList, &t)
			}
			if _, ok := singleTypeKeyInfo.InfoList[c.loginUserID]; ok {
				v.IsContainSelf = true
			}
			mapResult[prefixKey] = v
		}
	}
	for _, v := range mapResult {
		result = append(result, v)

	}
	return result
}
