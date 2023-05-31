package conversation_msg

import (
	"context"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"

	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	pbMsg "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/jinzhu/copier"
)

// Delete the local and server
// Delete the local, do not change the server data
// To delete the server, you need to change the local message status to delete
func (c *Conversation) clearConversationFromLocalAndSvr(ctx context.Context, conversationID string) error {
	_, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	// Use conversationID to remove conversations and messages from the server first
	err = c.clearConversationMsgFromSvr(ctx, conversationID)
	if err != nil {
		return err
	}
	if err := c.clearConversationAndDeleteAllMsg(ctx, conversationID); err != nil {
		return err
	}
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.ConChange, Args: []string{conversationID}}})
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.TotalUnreadMessageChanged}})
	return nil
}

func (c *Conversation) clearConversationAndDeleteAllMsg(ctx context.Context, conversationID string) error {
	err := c.db.DeleteConversationAllMessages(ctx, conversationID)
	if err != nil {
		return err
	}
	err = c.db.ResetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	return nil
}

// To delete session information, delete the server first, and then invoke the interface.
// The client receives a callback to delete all local information.
func (c *Conversation) clearConversationMsgFromSvr(ctx context.Context, conversationID string) error {
	var apiReq pbMsg.ClearConversationsMsgReq
	apiReq.UserID = c.loginUserID
	apiReq.ConversationIDs = []string{conversationID}
	return util.ApiPost(ctx, constant.ClearConversationMsgRouter, &apiReq, nil)
}

// Delete all messages
func (c *Conversation) deleteAllMessage(ctx context.Context) error {
	// Delete the server first (high error rate), then delete it.
	err := c.deleteAllMessageFromSvr(ctx)
	if err != nil {
		return err
	}
	err = c.deleteAllMsgFromLocal(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Delete all server messages
func (c *Conversation) deleteAllMessageFromSvr(ctx context.Context) error {
	var apiReq pbMsg.UserClearAllMsgReq
	apiReq.UserID = c.loginUserID
	err := util.ApiPost(ctx, constant.ClearAllMsgRouter, &apiReq, nil)
	if err != nil {
		return err
	}
	return nil
}

// Delete all messages from the local
func (c *Conversation) deleteAllMsgFromLocal(ctx context.Context) error {
	conversationList, err := c.db.GetAllConversationListDB(ctx)
	if err != nil {
		return err
	}
	var cidList []string
	for _, conversation := range conversationList {
		cidList = append(cidList, conversation.ConversationID)
	}
	for _, v := range conversationList {
		c.clearConversationAndDeleteAllMsg(ctx, v.ConversationID)
	}
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.ConChange, Args: cidList}})
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.TotalUnreadMessageChanged}})
	return nil

}

// Delete a message from the local
func (c *Conversation) deleteMessage(ctx context.Context, s *sdk_struct.MsgStruct) error {
	if err := c.deleteMessageFromSvr(ctx, s); err != nil {
		return err
	}
	return c.deleteMessageFromLocal(ctx, s)
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
	conversationID := utils.GetConversationIDByMsg(s)
	if err := c.db.DeleteConversationMsgs(ctx, conversationID, []string{s.ClientMsgID}); err != nil {
		return err
	}
	conversation, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	var latestMsg sdk_struct.MsgStruct
	utils.JsonStringToStruct(conversation.LatestMsg, &latestMsg)
	if latestMsg.ClientMsgID == s.ClientMsgID {
		log.ZDebug(ctx, "latesetMsg deleted", "seq", latestMsg.Seq, "clientMsgID", latestMsg.ClientMsgID)
		msgs, err := c.db.GetMessageListNoTime(ctx, conversationID, 1, false)
		if err != nil {
			return err
		}
		latestMsgSendTime := latestMsg.SendTime
		latestMsgStr := ""
		if len(msgs) > 0 {
			copier.Copy(&latestMsg, msgs[0])
			err := c.msgConvert(&latestMsg)
			if err != nil {
				log.ZError(ctx, "parsing data error", err, latestMsg)
			}
			latestMsgStr = utils.StructToJsonString(latestMsg)
			latestMsgSendTime = latestMsg.SendTime
		}
		if err := c.db.UpdateColumnsConversation(ctx, conversationID, map[string]interface{}{"latest_msg": latestMsgStr, "latest_msg_send_time": latestMsgSendTime}); err != nil {
			return err
		}
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.ConChange, Args: []string{conversationID}}})
	}
	c.msgListener.OnMsgDeleted(utils.StructToJsonString(s))
	return nil
}

func (c *Conversation) doDeleteMsgs(ctx context.Context, msg *sdkws.MsgData) {
	tips := sdkws.DeleteMsgsTips{}
	utils.UnmarshalNotificationElem(msg.Content, &tips)
	for _, v := range tips.Seqs {
		msg, err := c.db.GetMessageBySeq(ctx, tips.ConversationID, v)
		if err != nil {
			log.ZError(ctx, "GetMessageBySeq err", err, "conversationID", tips.ConversationID, "seq", v)
			continue
		}
		var s sdk_struct.MsgStruct
		copier.Copy(&s, msg)
		err = c.msgConvert(&s)
		if err != nil {
			log.ZError(ctx, "parsing data error", err, "msg", msg)
		}
		if err := c.deleteMessageFromLocal(ctx, &s); err != nil {
			log.ZError(ctx, "deleteMessageFromLocal err", err, "conversationID", tips.ConversationID, "seq", v)
		}
	}
}

func (c *Conversation) doClearConversations(ctx context.Context, msg *sdkws.MsgData) {
	tips := sdkws.ClearConversationTips{}
	utils.UnmarshalNotificationElem(msg.Content, &tips)
	log.ZDebug(ctx, "doClearConversations", "tips", tips)
	for _, v := range tips.ConversationIDs {
		if err := c.clearConversationAndDeleteAllMsg(ctx, v); err != nil {
			log.ZError(ctx, "clearConversation err", err, "conversationID", v)
		}
	}
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.ConChange, Args: tips.ConversationIDs}})
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.TotalUnreadMessageChanged}})
}
