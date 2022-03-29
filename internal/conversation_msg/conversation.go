package conversation_msg

import (
	"errors"
	"github.com/golang/protobuf/proto"
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
			OwnerUserID:      c.loginUserID,
			ConversationID:   conversationID,
			ConversationType: localConversation.ConversationType,
			UserID:           localConversation.UserID,
			GroupID:          localConversation.GroupID,
			RecvMsgOpt:       int32(opt),
			IsPinned:         localConversation.IsPinned,
			IsPrivateChat:    localConversation.IsPrivateChat,
			UnreadCount:      localConversation.UnreadCount,
			DraftTextTime:    localConversation.DraftTextTime,
			AttachedInfo:     localConversation.AttachedInfo,
			Ex:               localConversation.Ex,
		})
	}
	apiReq.Conversations = conversations
	c.p.PostFatalCallback(callback, constant.BatchSetConversationRouter, apiReq, &apiResp, apiReq.OperationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "output: ", apiResp)
	c.SyncConversations(operationID)
}

func (c *Conversation) setConversation(callback open_im_sdk_callback.Base, apiReq *server_api_params.SetConversationReq, conversationID string, operationID string) {
	localConversation, err := c.db.GetConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	apiResp := server_api_params.SetConversationResp{}
	apiReq.OwnerUserID = c.loginUserID
	apiReq.OperationID = operationID
	apiReq.ConversationID = conversationID
	apiReq.ConversationType = localConversation.ConversationType
	apiReq.UserID = localConversation.UserID
	apiReq.GroupID = localConversation.GroupID
	apiReq.Ex = localConversation.Ex
	apiReq.AttachedInfo = localConversation.AttachedInfo
	apiReq.DraftTextTime = localConversation.DraftTextTime
	apiReq.UnreadCount = localConversation.UnreadCount
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

func (c *Conversation) SyncConversations(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	conversationsOnServer, err := c.getServerConversationList(operationID)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
		return
	}
	conversationsOnLocal, err := c.db.GetAllConversationListToSync()
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
	}

	conversationsOnLocalTempFormat := common.LocalTransferToTempConversation(conversationsOnLocal)
	conversationsOnServerTempFormat := common.ServerTransferToTempConversation(conversationsOnServer)
	conversationsOnServerLocalFormat := common.TransferToLocalConversation(conversationsOnServer)

	aInBNot, bInANot, sameA, sameB := common.CheckConversationListDiff(conversationsOnServerTempFormat, conversationsOnLocalTempFormat)
	log.NewInfo(operationID, "diff ", aInBNot, bInANot, sameA, sameB)

	// server有 local没有
	// 可能是其他点开一下生成会话设置免打扰 插入到本地 不回调
	for _, index := range aInBNot {
		conversation := conversationsOnServerLocalFormat[index]
		var newConversation db.LocalConversation
		newConversation.ConversationID = conversation.ConversationID
		newConversation.ConversationType = conversation.ConversationType
		switch conversation.ConversationType {
		case constant.SingleChatType:
			newConversation.UserID = conversation.UserID
			faceUrl, name, err := c.friend.GetUserNameAndFaceUrlByUid(&tmpCallback{}, conversation.UserID, operationID)
			if err != nil {
				log.NewError(operationID, utils.GetSelfFuncName(), "GetUserNameAndFaceUrlByUid error", err.Error())
				continue
			}
			newConversation.ShowName = name
			newConversation.FaceURL = faceUrl
		case constant.GroupChatType:
			newConversation.GroupID = conversation.GroupID
			g, err := c.group.GetGroupInfoFromLocal2Svr(conversation.GroupID)
			if err != nil {
				log.NewError(operationID, utils.GetSelfFuncName(), "GetGroupInfoFromLocal2Svr error", err.Error())
				continue
			}
			newConversation.ShowName = g.GroupName
			newConversation.FaceURL = g.FaceURL
		}
		err := c.db.InsertConversation(&newConversation)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), "InsertConversation error", err.Error())
			continue
		}
		err = c.db.UpdateConversationForSync(conversation)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), "InsertConversation failed ", err.Error(), conversation)
			continue
		}
	}
	// 本地服务器有的会话 以服务器为准更新
	var conversationChangedList []string
	for _, index := range sameA {
		log.NewInfo("", *conversationsOnServerLocalFormat[index])
		err := c.db.UpdateConversationForSync(conversationsOnServerLocalFormat[index])
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), "UpdateConversation failed ", err.Error(), *conversationsOnServerLocalFormat[index])
			continue
		}
		conversationChangedList = append(conversationChangedList, conversationsOnServerLocalFormat[index].ConversationID)
	}
	// callback
	if len(conversationChangedList) > 0 {
		if err = common.TriggerCmdUpdateConversation(common.UpdateConNode{Action: constant.ConChange, Args: conversationChangedList}, c.ch); err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
		}
	}

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
	resp, _ := c.InternalSendMessage(callback, (*sdk_struct.MsgStruct)(&req), recvID, groupID, operationID, &server_api_params.OfflinePushInfo{}, false, options)
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
	var localMessage db.LocalChatLog
	var newMessageIDList []string
	messages, err := c.db.GetMultipleMessage(msgIDList)
	common.CheckDBErrCallback(callback, err, operationID)
	for _, v := range messages {
		if v.IsRead == false && v.ContentType < constant.NotificationBegin && v.SendID != c.loginUserID {
			newMessageIDList = append(newMessageIDList, v.ClientMsgID)
		}
	}
	if len(newMessageIDList) == 0 {
		common.CheckAnyErrCallback(callback, 201, errors.New("message has been marked read or sender is yourself"), operationID)
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
	err2 := c.db.UpdateMessageHasRead(userID, newMessageIDList, constant.SingleChatType)
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

func (c *Conversation) deleteMessageFromSvr(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, operationID string) {
	var apiReq server_api_params.DeleteMsgReq
	seq, err := c.db.GetMsgSeqByClientMsgID(s.ClientMsgID)
	common.CheckDBErrCallback(callback, err, operationID)
	apiReq.SeqList = []uint32{seq}
	apiReq.OpUserID = c.loginUserID
	apiReq.UserID = c.loginUserID
	apiReq.OperationID = operationID
	c.p.PostFatalCallback(callback, constant.DeleteMsgRouter, apiReq, nil, apiReq.OperationID)
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
		err = c.db.UpdateColumnsConversation(conversation.ConversationID, map[string]interface{}{"latest_msg_send_time": conversation.LatestMsgSendTime, "latest_msg": conversation.LatestMsg})
		if err != nil {
			log.Error("internal", "updateConversationLatestMsgModel err: ", err)
		} else {
			_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.ch)
		}
	}
}
func (c *Conversation) searchLocalMessages(callback open_im_sdk_callback.Base, searchParam sdk.SearchLocalMessagesParams, operationID string) (r sdk.SearchLocalMessagesCallback) {
	var conversationID string
	var startTime, endTime int64
	var searchResultItem sdk.SearchByConversationResult
	var messageList sdk_struct.NewMsgList
	var list []*db.LocalChatLog
	var err error
	if searchParam.PageIndex < 1 || searchParam.Count < 1 {
		common.CheckAnyErrCallback(callback, 201, errors.New("page or count is null"), operationID)
	}
	offset := (searchParam.PageIndex - 1) * searchParam.Count
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
	if searchParam.SearchTimePosition == 0 {
		endTime = 0
	} else {
		endTime = startTime - searchParam.SearchTimePeriod
	}
	if (len(searchParam.KeywordList) == 0 || searchParam.KeywordList[0] == "") && len(searchParam.MessageTypeList) == 0 {
		common.CheckAnyErrCallback(callback, 201, errors.New("keyword is null"), operationID)
	}
	if len(searchParam.MessageTypeList) != 0 && len(searchParam.KeywordList) == 0 {
		list, err = c.db.SearchMessageByContentType(searchParam.MessageTypeList, searchParam.SourceID, utils.UnixSecondToTime(endTime).UnixNano()/1e6, utils.UnixSecondToTime(startTime).UnixNano()/1e6, searchParam.SessionType, offset, searchParam.Count)
	} else {
		list, err = c.db.SearchMessageByKeyword(searchParam.KeywordList[0], searchParam.SourceID, utils.UnixSecondToTime(endTime).UnixNano()/1e6, utils.UnixSecondToTime(startTime).UnixNano()/1e6, searchParam.SessionType, offset, searchParam.Count)
	}

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
	log.Debug(operationID, utils.GetSelfFuncName(), *local)
	common.CheckDBErrCallback(callback, err, operationID)
	var seqList []uint32
	switch local.ConversationType {
	case constant.SingleChatType:
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
	}
	var apiReq server_api_params.DeleteMsgReq
	apiReq.OpUserID = c.loginUserID
	apiReq.UserID = c.loginUserID
	apiReq.OperationID = operationID
	apiReq.SeqList = seqList
	c.p.PostFatalCallback(callback, constant.DeleteMsgRouter, apiReq, nil, apiReq.OperationID)
	common.CheckArgsErrCallback(callback, err, operationID)
}
