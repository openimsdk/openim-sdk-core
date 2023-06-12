package conversation_msg

import (
	"context"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	pbMsg "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

func (c *Conversation) markMsgAsRead2Svr(ctx context.Context, conversationID string, seqs []int64) error {
	req := &pbMsg.MarkMsgsAsReadReq{UserID: c.loginUserID, ConversationID: conversationID, Seqs: seqs}
	return util.ApiPost(ctx, constant.MarkMsgsAsReadRouter, req, nil)
}

func (c *Conversation) markConversationAsReadSvr(ctx context.Context, conversationID string, hasReadSeq int64) error {
	req := &pbMsg.MarkConversationAsReadReq{UserID: c.loginUserID, ConversationID: conversationID, HasReadSeq: hasReadSeq}
	return util.ApiPost(ctx, constant.MarkConversationAsRead, req, nil)
}

// mark a conversation's all message as read
func (c *Conversation) markConversationMessageAsRead(ctx context.Context, conversationID string) error {
	_, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	peerUserMaxSeq, err := c.db.GetConversationPeerNormalMsgSeq(ctx, conversationID)
	if err != nil {
		return err
	}
	maxSeq, err := c.db.GetConversationNormalMsgSeq(ctx, conversationID)
	if err != nil {
		return err
	}
	msgs, err := c.db.GetUnreadMessage(ctx, conversationID)
	if err != nil {
		return err
	}
	msgIDs, _ := c.getAsReadMsgMapAndList(ctx, msgs)
	if err := c.markConversationAsReadSvr(ctx, conversationID, maxSeq); err != nil {
		return err
	}
	_, err = c.db.MarkConversationMessageAsRead(ctx, conversationID, msgIDs)
	if err != nil {
		return err
	}
	if err := c.db.UpdateColumnsConversation(ctx, conversationID, map[string]interface{}{"unread_count": 0}); err != nil {
		log.ZError(ctx, "UpdateColumnsConversation err", err, "conversationID", conversationID)
	}
	c.unreadChangeTrigger(ctx, conversationID, peerUserMaxSeq == maxSeq)
	return nil
}

// mark a conversation's message as read by seqs
func (c *Conversation) markConversationMessageAsReadByMsgID(ctx context.Context, conversationID string, msgIDs []string) error {
	_, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	msgs, err := c.db.GetMessagesByClientMsgIDs(ctx, conversationID, msgIDs)
	if err != nil {
		return err
	}
	if len(msgs) == 0 {
		return nil
	}
	var hasReadSeq = msgs[0].Seq
	maxSeq, err := c.db.GetConversationNormalMsgSeq(ctx, conversationID)
	if err != nil {
		return err
	}
	markAsReadMsgIDs, seqs := c.getAsReadMsgMapAndList(ctx, msgs)
	if err := c.markMsgAsRead2Svr(ctx, conversationID, seqs); err != nil {
		return err
	}
	decrCount, err := c.db.MarkConversationMessageAsRead(ctx, conversationID, markAsReadMsgIDs)
	if err != nil {
		return err
	}
	if err := c.db.DecrConversationUnreadCount(ctx, conversationID, decrCount); err != nil {
		log.ZError(ctx, "decrConversationUnreadCount err", err, "conversationID", conversationID, "decrCount", decrCount)
	}
	c.unreadChangeTrigger(ctx, conversationID, hasReadSeq == maxSeq && msgs[0].SendID != c.loginUserID)
	return nil
}

func (c *Conversation) getAsReadMsgMapAndList(ctx context.Context, msgs []*model_struct.LocalChatLog) (asReadMsgIDs []string, seqs []int64) {
	for _, msg := range msgs {
		if !msg.IsRead && msg.ContentType < constant.NotificationBegin && msg.SendID != c.loginUserID {
			asReadMsgIDs = append(asReadMsgIDs, msg.ClientMsgID)
			seqs = append(seqs, msg.Seq)
		} else {
			log.ZWarn(ctx, "msg can't marked as read", nil, "msg", msg)
		}
	}
	return
}

func (c *Conversation) unreadChangeTrigger(ctx context.Context, conversationID string, latestMsgIsRead bool) {
	if latestMsgIsRead {
		_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{ConID: conversationID, Action: constant.UpdateLatestMessageChange}, c.GetCh())
	}
	_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())
	_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{ConID: conversationID, Action: constant.TotalUnreadMessageChanged}, c.GetCh())
}

func (c *Conversation) doUnreadCount(ctx context.Context, conversationID string, hasReadSeq int64) {
	conversation, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		log.ZError(ctx, "GetConversation err", err, "conversationID", conversationID)
		return
	}
	var seqs []int64
	if hasReadSeq > conversation.HasReadSeq {
		for i := conversation.HasReadSeq + 1; i <= hasReadSeq; i++ {
			seqs = append(seqs, i)
		}
		_, err := c.db.MarkConversationMessageAsReadBySeqs(ctx, conversationID, seqs)
		if err != nil {
			log.ZError(ctx, "MarkConversationMessageAsReadBySeqs err", err, "conversationID", conversationID, "seqs", seqs)
			return
		}
		if err := c.db.DecrConversationUnreadCount(ctx, conversationID, int64(len(seqs))); err != nil {
			log.ZError(ctx, "decrConversationUnreadCount err", err, "conversationID", conversationID, "decrCount", int64(len(seqs)))
		}
		if err := c.db.UpdateColumnsConversation(ctx, conversationID, map[string]interface{}{"has_read_seq": hasReadSeq}); err != nil {
			log.ZError(ctx, "UpdateColumnsConversation err", err, "conversationID", conversationID)
		}
	} else {
		log.ZWarn(ctx, "hasReadSeq <= conversation.HasReadSeq", nil, "hasReadSeq", hasReadSeq, "conversation.HasReadSeq", conversation.HasReadSeq)
	}
}

func (c *Conversation) doReadDrawing(ctx context.Context, msg *sdkws.MsgData) {
	tips := &sdkws.MarkAsReadTips{}
	utils.UnmarshalNotificationElem(msg.Content, tips)
	if tips.MarkAsReadUserID != c.loginUserID {
		log.ZDebug(ctx, "do readDrawing", "tips", tips)
		conversation, err := c.db.GetConversation(ctx, tips.ConversationID)
		if err != nil {
			log.ZError(ctx, "GetConversation err", err, "conversationID", tips.ConversationID)
			return
		}
		messages, err := c.db.GetMessagesBySeqs(ctx, tips.ConversationID, tips.Seqs)
		if err != nil {
			log.ZError(ctx, "GetMessagesBySeqs err", err, "conversationID", tips.ConversationID, "seqs", tips.Seqs)
			return
		}
		if conversation.ConversationType == constant.SingleChatType {
			var successMsgIDs []string
			for _, message := range messages {
				attachInfo := sdk_struct.AttachedInfoElem{}
				_ = utils.JsonStringToStruct(message.AttachedInfo, &attachInfo)
				attachInfo.HasReadTime = msg.SendTime
				message.AttachedInfo = utils.StructToJsonString(attachInfo)
				message.IsRead = true
				if err = c.db.UpdateMessage(ctx, tips.ConversationID, message); err != nil {
					log.ZError(ctx, "UpdateMessage err", err, "conversationID", tips.ConversationID, "message", message)
				} else {
					successMsgIDs = append(successMsgIDs, message.ClientMsgID)
				}
			}
			var messageReceiptResp = []*sdk_struct.MessageReceipt{{UserID: tips.MarkAsReadUserID, MsgIDList: successMsgIDs,
				SessionType: conversation.ConversationType, ReadTime: msg.SendTime}}
			c.msgListener.OnRecvC2CReadReceipt(utils.StructToJsonString(messageReceiptResp))
		} else if conversation.ConversationType == constant.SuperGroupChatType {
			var successMsgIDs []string
			for _, message := range messages {
				attachInfo := sdk_struct.AttachedInfoElem{}
				_ = utils.JsonStringToStruct(message.AttachedInfo, &attachInfo)
				attachInfo.HasReadTime = msg.SendTime
				attachInfo.GroupHasReadInfo.HasReadUserIDList = utils.RemoveRepeatedStringInList(append(attachInfo.GroupHasReadInfo.HasReadUserIDList, tips.MarkAsReadUserID))
				attachInfo.GroupHasReadInfo.HasReadCount = int32(len(attachInfo.GroupHasReadInfo.HasReadUserIDList))
				message.AttachedInfo = utils.StructToJsonString(attachInfo)
				if err = c.db.UpdateMessage(ctx, tips.ConversationID, message); err != nil {
					log.ZError(ctx, "UpdateMessage err", err, "conversationID", tips.ConversationID, "message", message)
				} else {
					successMsgIDs = append(successMsgIDs, message.ClientMsgID)
				}
			}
			var messageReceiptResp = []*sdk_struct.MessageReceipt{{GroupID: conversation.GroupID, MsgIDList: successMsgIDs,
				SessionType: conversation.ConversationType, ReadTime: msg.SendTime}}
			c.msgListener.OnRecvGroupReadReceipt(utils.StructToJsonString(messageReceiptResp))
		}
	} else {
		log.ZDebug(ctx, "do unread count", "tips", tips)
		c.doUnreadCount(ctx, tips.ConversationID, tips.HasReadSeq)
	}
}
