package conversation_msg

import (
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

func (c *Conversation) insertMessageToLocalStorage(callback common.Base, s *db.LocalChatLog, operationID string) string {
	err := c.db.InsertMessage(s)
	common.CheckDBErr(callback, err, operationID)
	return s.ClientMsgID
}

func (c *Conversation) clearGroupHistoryMessage(callback common.Base, groupID string, operationID string) {
	conversationID := c.GetConversationIDBySessionType(groupID, constant.GroupChatType)
	err := c.db.UpdateMessageStatusBySourceID(groupID, constant.MsgStatusHasDeleted, constant.GroupChatType)
	common.CheckDBErr(callback, err, operationID)
	err = c.db.ClearConversation(conversationID)
	common.CheckDBErr(callback, err, operationID)
	//	u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
}

func (c *Conversation) clearC2CHistoryMessage(callback common.Base, userID string, operationID string) {
	conversationID := c.GetConversationIDBySessionType(userID, constant.SingleChatType)
	err := c.db.UpdateMessageStatusBySourceID(userID, constant.MsgStatusHasDeleted, constant.SingleChatType)
	common.CheckDBErr(callback, err, operationID)
	err = c.db.ClearConversation(conversationID)
	common.CheckDBErr(callback, err, operationID)
	//u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
}

func (c *Conversation) deleteMessageFromLocalStorage(callback common.Base, s *sdk_struct.MsgStruct, operationID string) {
	var conversation db.LocalConversation
	var latestMsg sdk_struct.MsgStruct
	var conversationID string
	var sourceID string
	chatLog := db.LocalChatLog{ClientMsgID: s.ClientMsgID, Status: constant.MsgStatusHasDeleted}
	err := c.db.UpdateMessage(&chatLog)
	common.CheckDBErr(callback, err, operationID)

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
	common.CheckDBErr(callback, err, operationID)
	common.JsonUnmarshal(LocalConversation.LatestMsg, &latestMsg, callback, operationID)

	if s.ClientMsgID == latestMsg.ClientMsgID { //If the deleted message is the latest message of the conversation, update the latest message of the conversation
		list, err := c.db.GetMessageList(sourceID, int(s.SessionType), 1, s.SendTime+TimeOffset)
		common.CheckDBErr(callback, err, operationID)

		conversation.ConversationID = conversationID
		if list == nil {
			conversation.LatestMsg = ""
			conversation.LatestMsgSendTime = uint32(utils.GetCurrentTimestampByNano())
		} else {
			conversation.LatestMsg = utils.StructToJsonString(list[0])
			conversation.LatestMsgSendTime = list[0].SendTime
		}
		//		err = u.triggerCmdUpdateConversation(common.updateConNode{ConId: conversationID, Action: constant.AddConOrUpLatMsg, Args: conversation})

		//	u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
	}
}
