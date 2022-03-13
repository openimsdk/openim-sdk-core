package conversation_msg

import (
	"errors"
	_ "open_im_sdk/internal/common"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"sort"
	"time"
)

func (c *Conversation) getAllConversationList(callback open_im_sdk_callback.Base, operationID string) sdk.GetAllConversationListCallback {
	conversationList, err := c.db.GetAllConversationList()
	common.CheckDBErrCallback(callback, err, operationID)
	return conversationList
}
func (c *Conversation) getConversationListSplit(callback open_im_sdk_callback.Base, offset, count int, operationID string) sdk.GetConversationListSplitCallback {
	conversationList, err := c.db.GetConversationListSplit(offset, count)
	common.CheckDBErrCallback(callback, err, operationID)
	return conversationList
}

func (c *Conversation) setConversationRecvMessageOpt(callback open_im_sdk_callback.Base, conversationIDList []string, opt int, operationID string) {
	apiReq := server_api_params.BatchSetConversationsReq{}
	apiResp := server_api_params.BatchSetConversationsResp{}
	apiReq.OperationID = operationID
	apiReq.OwnerUserID = c.loginUserID
	var conversations []server_api_params.Conversation
	for _, conversationID := range conversationIDList {
		localConversation, err := c.db.GetConversation(conversationID)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), "GetConversation failed", err.Error())
			continue
		}
		conversations = append(conversations, server_api_params.Conversation{
			OwnerUserID:    c.loginUserID,
			ConversationID: conversationID,
			RecvMsgOpt:     int32(opt),
			IsPinned:       localConversation.IsPinned,
			IsPrivateChat:  localConversation.IsPinned,
		})
	}
	apiReq.Conversations = conversations
	c.p.PostFatalCallback(callback, constant.BatchSetConversationRouter, apiReq, &apiResp, apiReq.OperationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "output: ", apiResp)
	c.SyncConversations(operationID)
}

func (c *Conversation) setConversation(callback open_im_sdk_callback.Base, apiReq *server_api_params.SetConversationReq, conversationID string, operationID string) {
	apiResp := server_api_params.SetConversationResp{}
	apiReq.OwnerUserID = c.loginUserID
	apiReq.OperationID = operationID
	apiReq.ConversationID = conversationID
	c.p.PostFatalCallback(callback, constant.SetConversationOptRouter, apiReq, nil, apiReq.OperationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "request success, output: ", apiResp)
}

func (c *Conversation) setOneConversationRecvMessageOpt(callback open_im_sdk_callback.Base, conversationID string, opt int, operationID string) {
	apiReq := &server_api_params.SetConversationReq{}
	localConversation, err := c.db.GetConversation(conversationID)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "GetConversation failed", err.Error())
		callback.OnError(constant.ErrDB.ErrCode, constant.ErrDB.ErrMsg)
		return
	}
	apiReq.RecvMsgOpt = int32(opt)
	apiReq.IsPinned = localConversation.IsPinned
	apiReq.IsPrivateChat = localConversation.IsPrivateChat
	c.setConversation(callback, apiReq, conversationID, operationID)
	c.SyncConversations(operationID)
}

func (c *Conversation) setOneConversationPrivateChat(callback open_im_sdk_callback.Base, conversationID string, isPrivate bool, operationID string) {
	apiReq := &server_api_params.SetConversationReq{}
	localConversation, err := c.db.GetConversation(conversationID)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "GetConversation failed", err.Error())
		callback.OnError(constant.ErrDB.ErrCode, constant.ErrDB.ErrMsg)
		return
	}
	apiReq.RecvMsgOpt = localConversation.RecvMsgOpt
	apiReq.IsPinned = localConversation.IsPinned
	apiReq.IsPrivateChat = isPrivate
	c.setConversation(callback, apiReq, conversationID, operationID)
	c.SyncConversations(operationID)
}

func (c *Conversation) setOneConversationPinned(callback open_im_sdk_callback.Base, conversationID string, isPinned bool, operationID string) {
	apiReq := &server_api_params.SetConversationReq{}
	localConversation, err := c.db.GetConversation(conversationID)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "GetConversation failed", err.Error())
		callback.OnError(constant.ErrDB.ErrCode, constant.ErrDB.ErrMsg)
		return
	}
	apiReq.RecvMsgOpt = localConversation.RecvMsgOpt
	apiReq.IsPinned = isPinned
	apiReq.IsPrivateChat = localConversation.IsPrivateChat
	c.setConversation(callback, apiReq, conversationID, operationID)
}

func (c *Conversation) getConversationRecvMessageOpt(callback open_im_sdk_callback.Base, conversationIDList []string, operationID string) server_api_params.GetConversationResp {
	apiReq := server_api_params.GetConversationsReq{}
	apiReq.OperationID = operationID
	apiReq.OwnerUserID = c.loginUserID
	apiReq.ConversationIDs = conversationIDList
	var realData server_api_params.GetConversationResp
	c.p.PostFatalCallback(callback, constant.GetConversationRouter, apiReq, &realData, apiReq.OperationID)
	return realData
}
func (c *Conversation) getOneConversation(callback open_im_sdk_callback.Base, sourceID string, sessionType int32, operationID string) *db.LocalConversation {
	conversationID := utils.GetConversationIDBySessionType(sourceID, int(sessionType))
	lc, err := c.db.GetConversation(conversationID)
	if err == nil {
		return lc
	} else {
		var newConversation db.LocalConversation
		newConversation.ConversationID = conversationID
		newConversation.ConversationType = sessionType
		switch sessionType {
		case constant.SingleChatType:
			newConversation.UserID = sourceID
			faceUrl, name, err := c.friend.GetUserNameAndFaceUrlByUid(callback, sourceID, operationID)
			common.CheckDBErrCallback(callback, err, operationID)
			newConversation.ShowName = name
			newConversation.FaceURL = faceUrl
		case constant.GroupChatType:
			newConversation.GroupID = sourceID
			g, err := c.group.GetGroupInfoFromLocal2Svr(sourceID)
			//g, err := c.db.GetGroupInfoByGroupID(sourceID)
			common.CheckDBErrCallback(callback, err, operationID)
			newConversation.ShowName = g.GroupName
			newConversation.FaceURL = g.FaceURL
		}
		err := c.db.InsertConversation(&newConversation)
		common.CheckDBErrCallback(callback, err, operationID)
		return &newConversation
	}
}
func (c *Conversation) getMultipleConversation(callback open_im_sdk_callback.Base, conversationIDList []string, operationID string) sdk.GetMultipleConversationCallback {
	conversationList, err := c.db.GetMultipleConversation(conversationIDList)
	common.CheckDBErrCallback(callback, err, operationID)
	return conversationList
}

func (c *Conversation) deleteConversation(callback open_im_sdk_callback.Base, conversationID, operationID string) {
	lc, err := c.db.GetConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	var sourceID string
	switch lc.ConversationType {
	case constant.SingleChatType:
		sourceID = lc.UserID
	case constant.GroupChatType:
		sourceID = lc.GroupID
	}
	//Mark messages related to this conversation for deletion
	err = c.db.UpdateMessageStatusBySourceID(sourceID, constant.MsgStatusHasDeleted, lc.ConversationType)
	common.CheckDBErrCallback(callback, err, operationID)
	//Reset the session information, empty session
	err = c.db.DeleteConversation(conversationID)
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
	c.SyncConversations(operationID)
}

func (c *Conversation) pinConversation(callback open_im_sdk_callback.Base, conversationID string, isPinned bool, operationID string) {
	lc := db.LocalConversation{ConversationID: conversationID, IsPinned: isPinned}
	if isPinned {
		err := c.db.UpdateConversation(&lc)
		common.CheckDBErrCallback(callback, err, operationID)
	} else {
		err := c.db.UnPinConversation(conversationID, constant.NotPinned)
		common.CheckDBErrCallback(callback, err, operationID)
	}
	c.setOneConversationPinned(callback, conversationID, isPinned, operationID)
	c.SyncConversations(operationID)
}

func (c *Conversation) getServerConversationList(operationID string) (server_api_params.GetAllConversationsResp, error) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	var req server_api_params.GetAllConversationsReq
	var resp server_api_params.GetAllConversationsResp
	req.OwnerUserID = c.loginUserID
	req.OperationID = operationID
	err := c.p.PostReturn(constant.GetAllConversationsRouter, req, &resp.Conversations)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
		return resp, err
	}
	return resp, nil
}

func (c *Conversation) getServerConversation(conversationID, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "conversationID: ", conversationID)
}

func (c *Conversation) SyncConversations(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	svrList, err := c.getServerConversationList(operationID)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
		return
	}
	// 判断不出本地有没有 重写
	conversationsOnLocal, err := c.db.GetAllConversationListToSync()
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
	}
	for _, v := range svrList.Conversations {
		log.Debug(operationID, v.IsPrivateChat)
	}
	conversationsOnServer := common.TransferToLocalConversation(svrList)
	aInBNot, bInANot, sameA, sameB := common.CheckConversationListDiff(conversationsOnServer, conversationsOnLocal)
	log.NewInfo(operationID, "diff ", aInBNot, bInANot, sameA, sameB)

	// server有 local没有
	// 可能是其他点开一下生成会话设置免打扰 插入到本地 不回调
	for _, index := range aInBNot {
		conversation := conversationsOnServer[index]
		conversation.LatestMsgSendTime = time.Now().Unix()
		err := c.db.InsertConversation(conversation)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), "InsertConversation failed ", err.Error(), conversation)
			continue
		}
	}
	// 本地服务器有的会话 以服务器为准更新 触发回调
	var conversationChangedList []string
	for _, index := range sameA {
		log.NewInfo("", *conversationsOnServer[index])
		err := c.db.UpdateConversationForSync(conversationsOnServer[index])
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), "UpdateConversation failed ", err.Error(), *conversationsOnServer[index])
			continue
		}
	}
	c.ConversationListener.OnConversationChanged(utils.StructToJsonString(conversationChangedList))
	// local有 server没有 代表没有修改公共字段
	for _, index := range bInANot {
		log.NewDebug(operationID, utils.GetSelfFuncName(), index, conversationsOnLocal[index].ConversationID,
			conversationsOnLocal[index].RecvMsgOpt, conversationsOnLocal[index].IsPinned, conversationsOnLocal[index].IsPrivateChat)
	}
}

func (c *Conversation) SyncOneConversation(conversationID, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "conversationID: ", conversationID)
	// todo
}

func (c *Conversation) getHistoryMessageList(callback open_im_sdk_callback.Base, req sdk.GetHistoryMessageListParams, operationID string) sdk.GetHistoryMessageListCallback {
	var sourceID string
	var conversationID string
	var startTime int64
	var sessionType int
	var messageList sdk_struct.NewMsgList
	if req.UserID == "" {
		sourceID = req.GroupID
		conversationID = utils.GetConversationIDBySessionType(sourceID, constant.GroupChatType)
		sessionType = constant.GroupChatType
	} else {
		sourceID = req.UserID
		conversationID = utils.GetConversationIDBySessionType(sourceID, constant.SingleChatType)
		sessionType = constant.SingleChatType
	}
	if req.StartClientMsgID == "" {
		lc, err := c.db.GetConversation(conversationID)
		if err != nil {
			return nil
		}
		startTime = lc.LatestMsgSendTime + TimeOffset

	} else {
		m, err := c.db.GetMessage(req.StartClientMsgID)
		common.CheckDBErrCallback(callback, err, operationID)
		startTime = m.SendTime
	}
	log.Info(operationID, "sourceID:", sourceID, "startTime:", startTime, "count:", req.Count)
	list, err := c.db.GetMessageList(sourceID, sessionType, req.Count, startTime)
	common.CheckDBErrCallback(callback, err, operationID)
	localChatLogToMsgStruct(&messageList, list)
	if req.UserID == "" {
		for _, v := range messageList {
			err := c.msgHandleByContentType(v)
			if err != nil {
				log.Error(operationID, "Parsing data error:", err.Error(), v)
				continue
			}
			v.GroupID = v.RecvID
			v.RecvID = c.loginUserID
		}
	} else {
		for _, v := range messageList {
			err := c.msgHandleByContentType(v)
			if err != nil {
				log.Error(operationID, "Parsing data error:", err.Error(), v)
				continue
			}
		}
	}
	sort.Sort(messageList)
	return sdk.GetHistoryMessageListCallback(messageList)
}
func (c *Conversation) revokeOneMessage(callback open_im_sdk_callback.Base, req sdk.RevokeMessageParams, operationID string) {
	var recvID, groupID string
	var localMessage db.LocalChatLog
	var lc db.LocalConversation
	var conversationID string
	message, err := c.db.GetMessage(req.ClientMsgID)
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
	resp, _ := c.internalSendMessage(callback, (*sdk_struct.MsgStruct)(&req), recvID, groupID, operationID, &server_api_params.OfflinePushInfo{}, false, options)
	req.ServerMsgID = resp.ServerMsgID
	req.SendTime = resp.SendTime
	req.Status = constant.MsgStatusSendSuccess
	msgStructToLocalChatLog(&localMessage, (*sdk_struct.MsgStruct)(&req))
	err = c.db.InsertMessage(&localMessage)
	if err != nil {
		log.Error(operationID, "inset into chat log err", localMessage, req)
	}
	err = c.db.UpdateColumnsMessage(req.Content, map[string]interface{}{"status": constant.MsgStatusRevoked})
	if err != nil {
		log.Error(operationID, "update revoke message err", localMessage, req)
	}
	lc.LatestMsg = utils.StructToJsonString(req)
	lc.LatestMsgSendTime = req.SendTime
	lc.ConversationID = conversationID
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: lc.ConversationID, Action: constant.AddConOrUpLatMsg, Args: lc}, c.ch)
}
func (c *Conversation) typingStatusUpdate(callback open_im_sdk_callback.Base, recvID, msgTip, operationID string) {
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Typing, operationID)
	s.Content = msgTip
	options := make(map[string]bool, 5)
	utils.SetSwitchFromOptions(options, constant.IsHistory, false)
	utils.SetSwitchFromOptions(options, constant.IsPersistent, false)
	utils.SetSwitchFromOptions(options, constant.IsSenderSync, false)
	utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
	utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
	c.internalSendMessage(callback, &s, recvID, "", operationID, &server_api_params.OfflinePushInfo{}, true, options)

}

func (c *Conversation) markC2CMessageAsRead(callback open_im_sdk_callback.Base, msgIDList sdk.MarkC2CMessageAsReadParams, userID, operationID string) {
	var localMessage db.LocalChatLog
	var newMessageIDList []string
	messages, err := c.db.GetMultipleMessage(msgIDList)
	common.CheckDBErrCallback(callback, err, operationID)
	for _, v := range messages {
		if v.IsRead == false && v.ContentType < constant.NotificationBegin {
			newMessageIDList = append(newMessageIDList, v.ClientMsgID)
		}
	}
	if len(newMessageIDList) == 0 {
		common.CheckAnyErrCallback(callback, 201, errors.New("message has been marked read"), operationID)
	}
	conversationID := utils.GetConversationIDBySessionType(userID, constant.SingleChatType)
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.HasReadReceipt, operationID)
	s.Content = utils.StructToJsonString(newMessageIDList)
	options := make(map[string]bool, 5)
	utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
	utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
	//If there is an error, the coroutine ends, so judgment is not  required
	resp, _ := c.internalSendMessage(callback, &s, userID, "", operationID, &server_api_params.OfflinePushInfo{}, false, options)
	s.ServerMsgID = resp.ServerMsgID
	s.SendTime = resp.SendTime
	s.Status = constant.MsgStatusFiltered
	msgStructToLocalChatLog(&localMessage, &s)
	err = c.db.InsertMessage(&localMessage)
	if err != nil {
		log.Error(operationID, "inset into chat log err", localMessage, s, err.Error())
	}
	err2 := c.db.UpdateMessageHasRead(userID, newMessageIDList)
	if err2 != nil {
		log.Error(operationID, "update message has read error", newMessageIDList, userID, err2.Error())
	}
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.UpdateLatestMessageChange}, c.ch)
	//_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.ch)
}
func (c *Conversation) insertMessageToLocalStorage(callback open_im_sdk_callback.Base, s *db.LocalChatLog, operationID string) string {
	err := c.db.InsertMessage(s)
	common.CheckDBErrCallback(callback, err, operationID)
	return s.ClientMsgID
}

func (c *Conversation) clearGroupHistoryMessage(callback open_im_sdk_callback.Base, groupID string, operationID string) {
	conversationID := utils.GetConversationIDBySessionType(groupID, constant.GroupChatType)
	err := c.db.UpdateMessageStatusBySourceID(groupID, constant.MsgStatusHasDeleted, constant.GroupChatType)
	common.CheckDBErrCallback(callback, err, operationID)
	err = c.db.ClearConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.ch)

}

func (c *Conversation) clearC2CHistoryMessage(callback open_im_sdk_callback.Base, userID string, operationID string) {
	conversationID := utils.GetConversationIDBySessionType(userID, constant.SingleChatType)
	err := c.db.UpdateMessageStatusBySourceID(userID, constant.MsgStatusHasDeleted, constant.SingleChatType)
	common.CheckDBErrCallback(callback, err, operationID)
	err = c.db.ClearConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.ch)
}

func (c *Conversation) deleteMessageFromLocalStorage(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, operationID string) {
	var conversation db.LocalConversation
	var latestMsg sdk_struct.MsgStruct
	var conversationID string
	var sourceID string
	chatLog := db.LocalChatLog{ClientMsgID: s.ClientMsgID, Status: constant.MsgStatusHasDeleted}
	err := c.db.UpdateMessage(&chatLog)
	common.CheckDBErrCallback(callback, err, operationID)

	if s.SessionType == constant.GroupChatType {
		conversationID = utils.GetConversationIDBySessionType(s.GroupID, constant.GroupChatType)
		sourceID = s.GroupID

	} else if s.SessionType == constant.SingleChatType {
		if s.SendID != c.loginUserID {
			conversationID = utils.GetConversationIDBySessionType(s.SendID, constant.SingleChatType)
			sourceID = s.SendID
		} else {
			conversationID = utils.GetConversationIDBySessionType(s.RecvID, constant.SingleChatType)
			sourceID = s.RecvID
		}
	}
	LocalConversation, err := c.db.GetConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	common.JsonUnmarshalCallback(LocalConversation.LatestMsg, &latestMsg, callback, operationID)

	if s.ClientMsgID == latestMsg.ClientMsgID { //If the deleted message is the latest message of the conversation, update the latest message of the conversation
		list, err := c.db.GetMessageList(sourceID, int(s.SessionType), 1, s.SendTime+TimeOffset)
		common.CheckDBErrCallback(callback, err, operationID)

		conversation.ConversationID = conversationID
		if list == nil {
			conversation.LatestMsg = ""
			conversation.LatestMsgSendTime = s.SendTime
		} else {
			conversation.LatestMsg = utils.StructToJsonString(list[0])
			conversation.LatestMsgSendTime = list[0].SendTime
		}
		_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversation.ConversationID, Action: constant.AddConOrUpLatMsg, Args: conversation}, c.ch)

	}
}
func (c *Conversation) searchLocalMessages(callback open_im_sdk_callback.Base, searchParam sdk.SearchLocalMessagesParams, operationID string) (r sdk.SearchLocalMessagesCallback) {
	var conversationID string
	var startTime, endTime int64
	//var searchResultItems []sdk.SearchByConversationResult
	var searchResultItem sdk.SearchByConversationResult
	var messageList sdk_struct.NewMsgList
	switch searchParam.SessionType {
	case constant.SingleChatType:
		conversationID = utils.GetConversationIDBySessionType(searchParam.SourceID, constant.SingleChatType)
	case constant.GroupChatType:
		conversationID = utils.GetConversationIDBySessionType(searchParam.SourceID, constant.GroupChatType)
	default:
	}
	if searchParam.SearchTimePosition == 0 {
		startTime = utils.GetCurrentTimestampBySecond()
	} else {
		startTime = searchParam.SearchTimePosition
	}
	endTime = startTime - searchParam.SearchTimePeriod
	if len(searchParam.KeywordList) == 0 {
		common.CheckAnyErrCallback(callback, 201, errors.New("keyword is null"), operationID)
	}
	list, err := c.db.SearchMessageByKeyword(searchParam.KeywordList[0], utils.UnixSecondToTime(endTime).UnixNano()/1e6, utils.UnixSecondToTime(startTime).UnixNano()/1e6, searchParam.SessionType)
	common.CheckDBErrCallback(callback, err, operationID)
	r.TotalCount = len(list)
	localChatLogToMsgStruct(&messageList, list)
	switch searchParam.SessionType {
	case constant.SingleChatType:
		for _, v := range messageList {
			err := c.msgHandleByContentType(v)
			if err != nil {
				log.Error(operationID, "Parsing data error:", err.Error(), v)
				continue
			}
		}
		sort.Sort(messageList)
		searchResultItem.ConversationID = conversationID
		searchResultItem.MessageCount = r.TotalCount
		searchResultItem.MessageList = messageList
		r.SearchResultItems = append(r.SearchResultItems, &searchResultItem)
	case constant.GroupChatType:

		for _, v := range messageList {
			err := c.msgHandleByContentType(v)
			if err != nil {
				log.Error(operationID, "Parsing data error:", err.Error(), v)
				continue
			}
			v.GroupID = v.RecvID
			v.RecvID = c.loginUserID
		}
		sort.Sort(messageList)
		searchResultItem.ConversationID = conversationID
		searchResultItem.MessageCount = r.TotalCount
		searchResultItem.MessageList = messageList
		r.SearchResultItems = append(r.SearchResultItems, &searchResultItem)
	default:

	}
	return r
}

func (c *Conversation) setConversationNotification(msg *server_api_params.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	c.SyncConversations(operationID)
}

func (c *Conversation) DoNotification(msg *server_api_params.MsgData) {
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
