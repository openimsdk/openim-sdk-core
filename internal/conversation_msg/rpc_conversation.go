package conversation_msg

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/sdk_struct"
)

func (c *Conversation) getAllConversationList(callback common.Base, operationID string) sdk.GetAllConversationListCallback {
	conversationList, err := c.db.GetAllConversationList()
	common.CheckErr(callback, err, operationID)
	return conversationList
}
func (c *Conversation) getConversationListSplit(callback common.Base, offset, count int, operationID string) sdk.GetConversationListSplitCallback {
	conversationList, err := c.db.GetConversationListSplit(offset, count)
	common.CheckErr(callback, err, operationID)
	return conversationList
}

func (c *Conversation) setConversationRecvMessageOpt(callback common.Base, conversationIDList []string, opt int, operationID string) *server_api_params.CommDataResp {
	apiReq := server_api_params.SetReceiveMessageOptReq{}
	apiReq.OperationID = operationID
	apiReq.FromUserID = c.loginUserID
	var temp int32
	temp = int32(opt)
	apiReq.Opt = &temp
	apiReq.ConversationIDList = conversationIDList
	result := c.p.PostFatalCallback(callback, constant.SetReceiveMessageOptRouter, apiReq, operationID)
	c.db.SetMultipleConversationRecvMsgOpt(conversationIDList, opt)
	return result
}
func (c *Conversation) getConversationRecvMessageOpt(callback common.Base, conversationIDList []string, operationID string) []*server_api_params.OptResult {
	apiReq := server_api_params.GetReceiveMessageOptReq{}
	apiReq.OperationID = operationID
	apiReq.FromUserID = c.loginUserID
	apiReq.ConversationIDList = conversationIDList
	result := c.p.PostFatalCallback(callback, constant.GetReceiveMessageOptRouter, apiReq, operationID)
	var realData []*server_api_params.OptResult
	mapstructure.Decode(result.Data, realData)
	return realData
}
func (c *Conversation) getOneConversation(callback common.Base, sourceID string, sessionType int32, operationID string) *db.LocalConversation {
	conversationID := c.GetConversationIDBySessionType(sourceID, sessionType)
	lc, err := c.db.GetConversation(conversationID)
	common.CheckErr(callback, err, operationID)
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
		common.CheckErr(callback, err, operationID)
		return &newConversation
	}
}
func (c *Conversation) getMultipleConversation(callback common.Base, conversationIDList []string, operationID string) sdk.GetMultipleConversationCallback {
	conversationList, err := c.db.GetMultipleConversation(conversationIDList)
	common.CheckErr(callback, err, operationID)
	return conversationList
}

func (c *Conversation) deleteConversation(callback common.Base, conversationID, operationID string) {
	lc, err := c.db.GetConversation(conversationID)
	common.CheckErr(callback, err, operationID)
	var sourceID string
	switch lc.ConversationType {
	case constant.SingleChatType:
		sourceID = lc.UserID
	case constant.GroupChatType:
		sourceID = lc.GroupID
	}
	//Mark messages related to this conversation for deletion
	err = c.UpdateMessageStatusBySourceID(sourceID, constant.MsgStatusHasDeleted, lc.ConversationType)
	common.CheckErr(callback, err, operationID)
	//Reset the session information, empty session
	err = c.ResetConversation(conversationID)
	common.CheckErr(callback, err, operationID)
}
func (c *Conversation) setConversationDraft(callback common.Base, conversationID, draftText, operationID string) {
	if draftText != "" {
		err := c.db.SetConversationDraft(conversationID, draftText)
		common.CheckErr(callback, err, operationID)
	} else {
		err := c.db.RemoveConversationDraft(conversationID, draftText)
		common.CheckErr(callback, err, operationID)
	}
}

func (c *Conversation) pinConversation(callback common.Base, conversationID string, isPinned bool, operationID string) {
	lc := db.LocalConversation{ConversationID: conversationID}
	if isPinned {
		lc.IsPinned = constant.Pinned
		err := c.UpdateConversation(&lc)
		common.CheckErr(callback, err, operationID)
	} else {
		lc.IsPinned = constant.NotPinned
		err := c.UnPinConversation(conversationID, constant.NotPinned)
		common.CheckErr(callback, err, operationID)
	}
}

func (c *Conversation) getHistoryMessageList(callback common.Base, req sdk.GetHistoryMessageListParams, operationID string) {

	var sourceID string
	var conversationID string
	var startTime uint32
	var latestMsg sdk_struct.MsgStruct
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
		common.CheckErr(callback, err, operationID)
		if lc == nil {
			startTime = 0
		} else {
			startTime = lc.LatestMsgSendTime + TimeOffset
		}

	} else {
		startTime = req.StartMsg.SendTime
	}
	log.Info(operationID, "sourceID:", sourceID, "startTime:", startTime, "count:", req.Count)
	err, list := u.getHistoryMessage(sourceID, startTime, p.Count, sessionType)
	lc := db.LocalConversation{ConversationID: conversationID}
	if isPinned {
		lc.IsPinned = constant.Pinned
		err := c.UpdateConversation(&lc)
		common.CheckErr(callback, err, operationID)
	} else {
		lc.IsPinned = constant.NotPinned
		err := c.UnPinConversation(conversationID, constant.NotPinned)
		common.CheckErr(callback, err, operationID)
	}
}
