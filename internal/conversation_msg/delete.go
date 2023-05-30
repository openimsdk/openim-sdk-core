package conversation_msg

import (
	"context"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"

	pbMsg "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
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
	err = c.deleteMessageFromSvr(ctx)
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

func (c *Conversation) deleteMessageFromSvr(ctx context.Context) error {
	var apiReq pbMsg.DeleteMsgsReq
	// GetAllConversations

	if err != nil {
		return err
	}
	apiReq.UserID = c.loginUserID
	err = util.ApiPost(ctx, constant.ClearMsgRouter, &apiReq, nil)
	if err != nil {
		return err
	}

	// getReadDiffusionGroupIDList Is to get the list of group ids that have been read to roam
	groupIDList, err := c.full.GetReadDiffusionGroupIDList(ctx)
	if err != nil {
		return err
	}
	// Delete the roaming message of the group
	var superGroupApiReq pbMsg.DelSuperGroupMsgReq
	superGroupApiReq.UserID = c.loginUserID
	for _, v := range groupIDList {
		superGroupApiReq.GroupID = v
		err := util.ApiPost(ctx, constant.DeleteSuperGroupMsgRouter, &superGroupApiReq, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
