package conversation_msg

import (
	"context"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"

	pbMsg "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/jinzhu/copier"
)

// Delete the local and server
// Delete the local, do not change the server data
// To delete the server, you need to change the local message status to delete
func (c *Conversation) DeleteConversationFromLocalAndSvr(ctx context.Context, conversationID string) error {
	// Use conversationID to remove conversations and messages from the server first
	err := c.deleteConversationAndMsgFromSvr(ctx, conversationID)
	if err != nil {
		return err
	}
	return c.deleteConversation(ctx, conversationID)
}

// To delete session information, delete the server first, and then invoke the interface.
// The client receives a callback to delete all local information.
func (c *Conversation) deleteConversationAndMsgFromSvr(ctx context.Context, conversationID string) error {
	// Verify the existence of the session and prevent the client from deleting non-existent sessions
	_, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	// Since it is deleting a session, it is deleting the full message
	// for that session, so there is no need to pass the seqList here
	var apiReq pbMsg.ClearConversationsMsgReq
	apiReq.UserID = c.loginUserID
	apiReq.ConversationIDs = []string{conversationID}
	return util.ApiPost(ctx, constant.ClearConversationMsgRouter, &apiReq, nil)
}

// Delete all messages
func (c *Conversation) DeleteAllMessage(ctx context.Context) error {
	var apiReq pbMsg.UserClearAllMsgReq
	apiReq.UserID = c.loginUserID
	err := util.ApiPost(ctx, constant.ClearMsgRouter, &apiReq, nil)
	if err != nil {
		return err
	}

	// Delete the server first (high error rate), then delete it.
	err = c.DeleteAllMessageFromSvr(ctx)
	if err != nil {
		return err
	}

	err = c.deleteAllMsgFromLocal(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Delete all messages from the local
func (c *Conversation) deleteAllMsgFromLocal(ctx context.Context) error {
	err := c.db.DeleteAllMessage(ctx)
	if err != nil {
		return err
	}
	// getReadDiffusionGroupIDList Is to get the list of group ids that have been read to roam
	groupIDList, err := c.full.GetReadDiffusionGroupIDList(ctx)
	if err != nil {
		return err
	}
	for _, v := range groupIDList {
		err = c.db.SuperGroupDeleteAllMessage(ctx, v)
		if err != nil {
			//log.Error(operationID, "SuperGroupDeleteAllMessage err", err.Error())
			continue
		}
	}
	// TODO: GetAllConversations
	err = c.db.ClearAllConversation(ctx)
	if err != nil {
		return err
	}
	// GetAllConversationListDB Is to get a list of all sessions
	conversationList, err := c.db.GetAllConversationListDB(ctx)
	if err != nil {
		return err
	}
	var cidList []string
	for _, conversation := range conversationList {
		cidList = append(cidList, conversation.ConversationID)
	}
	_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{Action: constant.ConChange, Args: cidList}, c.GetCh())
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{"", constant.TotalUnreadMessageChanged, ""}})
	return nil

}

// Delete all server messages
func (c *Conversation) DeleteAllMessageFromSvr(ctx context.Context) error {
	var apiReq pbMsg.UserClearAllMsgReq
	apiReq.UserID = c.loginUserID
	err := util.ApiPost(ctx, constant.ClearMsgRouter, &apiReq, nil)
	if err != nil {
		return err
	}
	return nil
}

// Delete a message from the local
func (c *Conversation) deleteMessage(ctx context.Context, s *sdk_struct.MsgStruct) error {
	var conversation model_struct.LocalConversation
	var latestMsg sdk_struct.MsgStruct
	conversationID := utils.GetConversationIDByMsg(s)
	chatLog := model_struct.LocalChatLog{ClientMsgID: s.ClientMsgID, Status: constant.MsgStatusHasDeleted, SessionType: s.SessionType}
	err := c.db.UpdateMessage(ctx, conversationID, &chatLog)

	if err != nil {
		return err
	}
	LocalConversation, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	err = utils.JsonStringToStruct(LocalConversation.LatestMsg, &latestMsg)
	if err != nil {
		return err
	}

	if s.ClientMsgID == latestMsg.ClientMsgID { //If the deleted message is the latest message of the conversation, update the latest message of the conversation
		list, err := c.db.GetMessageListNoTime(ctx, conversationID, 1, false)
		if err != nil {
			return err
		}
		conversation.ConversationID = conversationID
		if list == nil {
			conversation.LatestMsg = ""
			conversation.LatestMsgSendTime = s.SendTime
		} else {
			copier.Copy(&latestMsg, list[0])
			err := c.msgConvert(&latestMsg)
			if err != nil {
				log.Error("", "Parsing data error:", err.Error(), latestMsg)
			}
			conversation.LatestMsg = utils.StructToJsonString(latestMsg)
			conversation.LatestMsgSendTime = latestMsg.SendTime
		}
		err = c.db.UpdateColumnsConversation(ctx, conversation.ConversationID, map[string]interface{}{"latest_msg_send_time": conversation.LatestMsgSendTime, "latest_msg": conversation.LatestMsg})
		if err != nil {
			log.Error("internal", "updateConversationLatestMsgModel err: ", err)
		} else {
			_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())
		}
	}
	return nil
}

// The user deletes part of the message from the server
func (c *Conversation) deleteMessageFromSvr(ctx context.Context, s *sdk_struct.MsgStruct) error {
	conversationID := utils.GetConversationIDByMsg(s)
	localMessage, err := c.db.GetMessage(ctx, conversationID, s.ClientMsgID)
	if err != nil {
		return err
	}
	var apiReq pbMsg.DeleteMsgsReq
	apiReq.UserID = c.loginUserID
	apiReq.Seqs = []int64{localMessage.Seq}
	apiReq.ConversationID = conversationID
	return util.ApiPost(ctx, constant.DeleteMsgsRouter, &apiReq, nil)

}

// Delete messages from local
func (c *Conversation) deleteMessageFromLocal(ctx context.Context, s *sdk_struct.MsgStruct) error {
	var conversationID string
	switch s.SessionType {
	case constant.GroupChatType:
		conversationID = c.getConversationIDBySessionType(s.GroupID, constant.GroupChatType)
	case constant.SingleChatType:
		if s.SendID != c.loginUserID {
			conversationID = c.getConversationIDBySessionType(s.SendID, constant.SingleChatType)
		} else {
			conversationID = c.getConversationIDBySessionType(s.RecvID, constant.SingleChatType)
		}
	case constant.SuperGroupChatType:
		conversationID = c.getConversationIDBySessionType(s.GroupID, constant.SuperGroupChatType)
	}
	err := c.db.ClearConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())
	return nil
}
