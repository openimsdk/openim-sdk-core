package conversation_msg

import (
	"errors"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

func (c *Conversation) getAllConversationList(callback common.Base, operationID string) sdk.GetAllConversationListCallback {
	conversationList, err := c.db.GetAllConversationList()
	common.CheckDBErrCallback(callback, err, operationID)
	return conversationList
}
func (c *Conversation) getConversationListSplit(callback common.Base, offset, count int, operationID string) sdk.GetConversationListSplitCallback {
	conversationList, err := c.db.GetConversationListSplit(offset, count)
	common.CheckDBErrCallback(callback, err, operationID)
	return conversationList
}

func (c *Conversation) setConversationRecvMessageOpt(callback common.Base, conversationIDList []string, opt int, operationID string) []*server_api_params.OptResult {
	apiReq := server_api_params.SetReceiveMessageOptReq{}
	apiReq.OperationID = operationID
	apiReq.FromUserID = c.loginUserID
	var temp int32
	temp = int32(opt)
	apiReq.Opt = &temp
	apiReq.ConversationIDList = conversationIDList
	var realData []*server_api_params.OptResult
	c.p.PostFatalCallback(callback, constant.SetReceiveMessageOptRouter, apiReq, realData, apiReq.OperationID)
	c.db.SetMultipleConversationRecvMsgOpt(conversationIDList, opt)
	return realData
}
func (c *Conversation) getConversationRecvMessageOpt(callback common.Base, conversationIDList []string, operationID string) []*server_api_params.OptResult {
	apiReq := server_api_params.GetReceiveMessageOptReq{}
	apiReq.OperationID = operationID
	apiReq.FromUserID = c.loginUserID
	apiReq.ConversationIDList = conversationIDList
	var realData []*server_api_params.OptResult
	c.p.PostFatalCallback(callback, constant.GetReceiveMessageOptRouter, apiReq, realData, apiReq.OperationID)
	return realData
}
func (c *Conversation) getOneConversation(callback common.Base, sourceID string, sessionType int32, operationID string) *db.LocalConversation {
	conversationID := c.GetConversationIDBySessionType(sourceID, sessionType)
	lc, err := c.db.GetConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	if lc != nil {
		return lc
	} else {
		var newConversation db.LocalConversation
		newConversation.ConversationID = conversationID
		newConversation.ConversationType = sessionType
		switch sessionType {
		case constant.SingleChatType:
			newConversation.UserID = sourceID
			//faceUrl, name, err := u.getUserNameAndFaceUrlByUid(sourceID)
			//if err != nil {
			//	callback.OnError(301, err.Error())
			//	utils.sdkLog("getUserNameAndFaceUrlByUid err:", err)
			//	return
			//}
			//c.ShowName = name
			//c.FaceURL = faceUrl
		case constant.GroupChatType:
			newConversation.GroupID = sourceID
			//faceUrl, name, err := u.getGroupNameAndFaceUrlByUid(sourceID)
			//if err != nil {
			//	callback.OnError(301, err.Error())
			//	utils.sdkLog("getGroupNameAndFaceUrlByUid err:", err)
			//}
			//c.ShowName = name
			//c.FaceURL = faceUrl

		}
		err := c.db.InsertConversation(&newConversation)
		common.CheckDBErrCallback(callback, err, operationID)
		return &newConversation
	}
}
func (c *Conversation) getMultipleConversation(callback common.Base, conversationIDList []string, operationID string) sdk.GetMultipleConversationCallback {
	conversationList, err := c.db.GetMultipleConversation(conversationIDList)
	common.CheckDBErrCallback(callback, err, operationID)
	return conversationList
}

func (c *Conversation) deleteConversation(callback common.Base, conversationID, operationID string) {
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
	err = c.UpdateMessageStatusBySourceID(sourceID, constant.MsgStatusHasDeleted, lc.ConversationType)
	common.CheckDBErrCallback(callback, err, operationID)
	//Reset the session information, empty session
	err = c.ResetConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
}
func (c *Conversation) setConversationDraft(callback common.Base, conversationID, draftText, operationID string) {
	if draftText != "" {
		err := c.db.SetConversationDraft(conversationID, draftText)
		common.CheckDBErrCallback(callback, err, operationID)
	} else {
		err := c.db.RemoveConversationDraft(conversationID, draftText)
		common.CheckDBErrCallback(callback, err, operationID)
	}
}

func (c *Conversation) pinConversation(callback common.Base, conversationID string, isPinned bool, operationID string) {
	lc := db.LocalConversation{ConversationID: conversationID}
	if isPinned {
		lc.IsPinned = constant.Pinned
		err := c.UpdateConversation(&lc)
		common.CheckDBErrCallback(callback, err, operationID)
	} else {
		lc.IsPinned = constant.NotPinned
		err := c.UnPinConversation(conversationID, constant.NotPinned)
		common.CheckDBErrCallback(callback, err, operationID)
	}
}

func (c *Conversation) getHistoryMessageList(callback common.Base, req sdk.GetHistoryMessageListParams, operationID string) sdk.GetHistoryMessageListCallback {
	var sourceID string
	var conversationID string
	var startTime int64
	var sessionType int
	if req.UserID == "" {
		sourceID = req.GroupID
		conversationID = c.GetConversationIDBySessionType(sourceID, constant.GroupChatType)
		sessionType = constant.GroupChatType
	} else {
		sourceID = req.UserID
		conversationID = c.GetConversationIDBySessionType(sourceID, constant.SingleChatType)
		sessionType = constant.SingleChatType
	}
	if req.StartMsg == nil {
		lc, err := c.db.GetConversation(conversationID)
		common.CheckDBErrCallback(callback, err, operationID)
		if lc == nil {
			startTime = 0
		} else {
			startTime = lc.LatestMsgSendTime + TimeOffset
		}

	} else {
		startTime = req.StartMsg.SendTime
	}
	log.Info(operationID, "sourceID:", sourceID, "startTime:", startTime, "count:", req.Count)
	list, err := c.db.GetMessageList(sourceID, sessionType, req.Count, startTime)
	common.CheckDBErrCallback(callback, err, operationID)
	return list

}
func (c *Conversation) revokeOneMessage(callback common.Base, req sdk.RevokeMessageParams, operationID string) {
	var recvID, groupID string
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
	case constant.GroupChatType:
		groupID = req.GroupID
	default:

		callback.OnError(200, "args err")
	}
	req.Content = message.ClientMsgID
	req.ClientMsgID = utils.GetMsgID(message.SendID)
	req.ContentType = constant.Revoke
	options := make(map[string]bool, 2)
	_ = c.internalSendMessage(callback, (*sdk_struct.MsgStruct)(&req), recvID, groupID, operationID, &server_api_params.OfflinePushInfo{}, false, options)
	//插入一条消息，以及会话最新的一条消息，触发UI的更新
	err = c.db.UpdateColumnsMessage(req.Content, map[string]interface{}{"status": constant.MsgStatusRevoked})
	common.CheckDBErrCallback(callback, err, operationID)
}
func (c *Conversation) typingStatusUpdate(callback common.Base, recvID, msgTip, operationID string) {
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Typing)
	s.Content = msgTip
	options := make(map[string]bool, 2)
	_ = c.internalSendMessage(callback, &s, recvID, "", operationID, &server_api_params.OfflinePushInfo{}, true, options)

}

func (c *Conversation) markC2CMessageAsRead(callback common.Base, msgIDList string, recvID, operationID string) {
	var list sdk.MarkC2CMessageAsReadParams
	common.JsonUnmarshalCallback(msgIDList, &list, callback, operationID)
	//conversationID := c.GetConversationIDBySessionType(recvID, constant.SingleChatType)

	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.HasReadReceipt)
	s.Content = msgIDList
	options := make(map[string]bool, 2)
	_ = c.internalSendMessage(callback, &s, recvID, "", operationID, &server_api_params.OfflinePushInfo{}, false, options)
	err := c.db.UpdateMessageHasRead(recvID, list)
	common.CheckDBErrCallback(callback, err, operationID)
	//u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{conversationID, constant.UpdateLatestMessageChange, ""}})
	//u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
}
func (c *Conversation) insertMessageToLocalStorage(callback common.Base, s *db.LocalChatLog, operationID string) string {
	err := c.db.InsertMessage(s)
	common.CheckDBErrCallback(callback, err, operationID)
	return s.ClientMsgID
}

func (c *Conversation) clearGroupHistoryMessage(callback common.Base, groupID string, operationID string) {
	conversationID := c.GetConversationIDBySessionType(groupID, constant.GroupChatType)
	err := c.db.UpdateMessageStatusBySourceID(groupID, constant.MsgStatusHasDeleted, constant.GroupChatType)
	common.CheckDBErrCallback(callback, err, operationID)
	err = c.db.ClearConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	//	u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
}

func (c *Conversation) clearC2CHistoryMessage(callback common.Base, userID string, operationID string) {
	conversationID := c.GetConversationIDBySessionType(userID, constant.SingleChatType)
	err := c.db.UpdateMessageStatusBySourceID(userID, constant.MsgStatusHasDeleted, constant.SingleChatType)
	common.CheckDBErrCallback(callback, err, operationID)
	err = c.db.ClearConversation(conversationID)
	common.CheckDBErrCallback(callback, err, operationID)
	//u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
}

func (c *Conversation) deleteMessageFromLocalStorage(callback common.Base, s *sdk_struct.MsgStruct, operationID string) {
	var conversation db.LocalConversation
	var latestMsg sdk_struct.MsgStruct
	var conversationID string
	var sourceID string
	chatLog := db.LocalChatLog{ClientMsgID: s.ClientMsgID, Status: constant.MsgStatusHasDeleted}
	err := c.db.UpdateMessage(&chatLog)
	common.CheckDBErrCallback(callback, err, operationID)

	callback.OnSuccess("")

	if s.SessionType == constant.GroupChatType {
		conversationID = c.GetConversationIDBySessionType(s.RecvID, constant.GroupChatType)
		sourceID = s.RecvID

	} else if s.SessionType == constant.SingleChatType {
		if s.SendID != c.loginUserID {
			conversationID = c.GetConversationIDBySessionType(s.SendID, constant.SingleChatType)
			sourceID = s.SendID
		} else {
			conversationID = c.GetConversationIDBySessionType(s.RecvID, constant.SingleChatType)
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
			conversation.LatestMsgSendTime = utils.GetCurrentTimestampByMill()
		} else {
			conversation.LatestMsg = utils.StructToJsonString(list[0])
			conversation.LatestMsgSendTime = list[0].SendTime
		}
		//		err = u.triggerCmdUpdateConversation(common.updateConNode{ConId: conversationID, Action: constant.AddConOrUpLatMsg, Args: conversation})

		//	u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
	}
}
