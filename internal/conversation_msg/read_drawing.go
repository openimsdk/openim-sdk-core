// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conversation_msg

import (
	"context"
	"encoding/json"
	"errors"
	utils2 "github.com/OpenIMSDK/tools/utils"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"

	pbMsg "github.com/OpenIMSDK/protocol/msg"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/log"
)

func (c *Conversation) markMsgAsRead2Svr(ctx context.Context, conversationID string, seqs []int64) error {
	req := &pbMsg.MarkMsgsAsReadReq{UserID: c.loginUserID, ConversationID: conversationID, Seqs: seqs}
	return util.ApiPost(ctx, constant.MarkMsgsAsReadRouter, req, nil)
}

func (c *Conversation) markConversationAsReadSvr(ctx context.Context, conversationID string, hasReadSeq int64, seqs []int64) error {
	req := &pbMsg.MarkConversationAsReadReq{UserID: c.loginUserID, ConversationID: conversationID, HasReadSeq: hasReadSeq, Seqs: seqs}
	return util.ApiPost(ctx, constant.MarkConversationAsRead, req, nil)
}

func (c *Conversation) setConversationHasReadSeq(ctx context.Context, conversationID string, hasReadSeq int64) error {
	req := &pbMsg.SetConversationHasReadSeqReq{UserID: c.loginUserID, ConversationID: conversationID, HasReadSeq: hasReadSeq}
	return util.ApiPost(ctx, constant.SetConversationHasReadSeq, req, nil)
}

func (c *Conversation) getConversationMaxSeqAndSetHasRead(ctx context.Context, conversationID string) error {
	maxSeq, err := c.db.GetConversationNormalMsgSeq(ctx, conversationID)
	if err != nil {
		return err
	}
	if maxSeq == 0 {
		return nil
	}
	if err := c.setConversationHasReadSeq(ctx, conversationID, maxSeq); err != nil {
		return err
	}
	if err := c.db.UpdateColumnsConversation(ctx, conversationID, map[string]interface{}{"has_read_seq": maxSeq}); err != nil {
		return err
	}
	return nil
}

// mark a conversation's all message as read
func (c *Conversation) markConversationMessageAsRead(ctx context.Context, conversationID string) error {
	_, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	// get the maximum sequence number of messages in the table that are not sent by oneself
	peerUserMaxSeq, err := c.db.GetConversationPeerNormalMsgSeq(ctx, conversationID)
	if err != nil {
		return err
	}
	// get the maximum sequence number of messages in the table
	maxSeq, err := c.db.GetConversationNormalMsgSeq(ctx, conversationID)
	if err != nil {
		return err
	}
	msgs, err := c.db.GetUnreadMessage(ctx, conversationID)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "get unread message", "msgs", len(msgs))
	msgIDs, seqs := c.getAsReadMsgMapAndList(ctx, msgs)
	if len(seqs) == 0 {
		log.ZWarn(ctx, "seqs is empty", nil, "conversationID", conversationID)
		return nil
	}
	log.ZDebug(ctx, "markConversationMessageAsRead", "conversationID", conversationID, "seqs",
		seqs, "peerUserMaxSeq", peerUserMaxSeq, "maxSeq", maxSeq)
	if err := c.markConversationAsReadSvr(ctx, conversationID, maxSeq, seqs); err != nil {
		return err
	}
	_, err = c.db.MarkConversationMessageAsReadDB(ctx, conversationID, msgIDs)
	if err != nil {
		log.ZWarn(ctx, "MarkConversationMessageAsRead err", err, "conversationID", conversationID, "msgIDs", msgIDs)
	}
	if err := c.db.UpdateColumnsConversation(ctx, conversationID, map[string]interface{}{"unread_count": 0}); err != nil {
		log.ZError(ctx, "UpdateColumnsConversation err", err, "conversationID", conversationID)
	}
	log.ZDebug(ctx, "update columns sucess")
	c.unreadChangeTrigger(ctx, conversationID, peerUserMaxSeq == maxSeq)
	return nil
}

// mark a conversation's message as read by seqs
func (c *Conversation) markMessagesAsReadByMsgID(ctx context.Context, conversationID string, msgIDs []string) error {
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
	log.ZDebug(ctx, "msgs len", "markAsReadMsgIDs", len(markAsReadMsgIDs), "seqs", seqs)
	if len(seqs) == 0 {
		log.ZWarn(ctx, "seqs is empty", nil, "conversationID", conversationID)
		return nil
	}
	if err := c.markMsgAsRead2Svr(ctx, conversationID, seqs); err != nil {
		return err
	}
	decrCount, err := c.db.MarkConversationMessageAsReadDB(ctx, conversationID, markAsReadMsgIDs)
	if err != nil {
		return err
	}
	if err := c.db.DecrConversationUnreadCount(ctx, conversationID, decrCount); err != nil {
		log.ZError(ctx, "decrConversationUnreadCount err", err, "conversationID", conversationID,
			"decrCount", decrCount)
	}
	c.unreadChangeTrigger(ctx, conversationID, hasReadSeq == maxSeq && msgs[0].SendID != c.loginUserID)
	return nil
}

func (c *Conversation) getAsReadMsgMapAndList(ctx context.Context,
	msgs []*model_struct.LocalChatLog) (asReadMsgIDs []string, seqs []int64) {
	for _, msg := range msgs {
		if !msg.IsRead && msg.SendID != c.loginUserID {
			if msg.Seq == 0 {
				log.ZWarn(ctx, "exception seq", errors.New("exception message "), "msg", msg)
			} else {
				asReadMsgIDs = append(asReadMsgIDs, msg.ClientMsgID)
				seqs = append(seqs, msg.Seq)
			}
		} else {
			log.ZWarn(ctx, "msg can't marked as read", nil, "msg", msg)
		}
	}
	return
}

func (c *Conversation) unreadChangeTrigger(ctx context.Context, conversationID string, latestMsgIsRead bool) {
	if latestMsgIsRead {
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: conversationID,
			Action: constant.UpdateLatestMessageChange, Args: []string{conversationID}}, Ctx: ctx})
	}
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: conversationID,
		Action: constant.ConChange, Args: []string{conversationID}}, Ctx: ctx})
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.TotalUnreadMessageChanged},
		Ctx: ctx})
}

func (c *Conversation) doUnreadCount(ctx context.Context, conversation *model_struct.LocalConversation, hasReadSeq int64, seqs []int64) {
	if conversation.ConversationType == constant.SingleChatType {
		if len(seqs) != 0 {
			_, err := c.db.MarkConversationMessageAsReadBySeqs(ctx, conversation.ConversationID, seqs)
			if err != nil {
				log.ZWarn(ctx, "MarkConversationMessageAsReadBySeqs err", err, "conversationID", conversation.ConversationID, "seqs", seqs)
			}
		} else {
			log.ZWarn(ctx, "seqs is empty", nil, "conversationID", conversation.ConversationID, "hasReadSeq", hasReadSeq)
		}
		if hasReadSeq > conversation.HasReadSeq {
			decrUnreadCount := hasReadSeq - conversation.HasReadSeq
			if err := c.db.DecrConversationUnreadCount(ctx, conversation.ConversationID, decrUnreadCount); err != nil {
				log.ZError(ctx, "DecrConversationUnreadCount err", err, "conversationID", conversation.ConversationID, "decrUnreadCount", decrUnreadCount)
			}
			if err := c.db.UpdateColumnsConversation(ctx, conversation.ConversationID, map[string]interface{}{"has_read_seq": hasReadSeq}); err != nil {
				log.ZError(ctx, "UpdateColumnsConversation err", err, "conversationID", conversation.ConversationID)
			}
		}
		var latestMsg sdkws.MsgData
		if err := json.Unmarshal([]byte(conversation.LatestMsg), &latestMsg); err != nil {
			log.ZError(ctx, "Unmarshal err", err, "conversationID", conversation.ConversationID, "latestMsg", conversation.LatestMsg)
		}
		if (!latestMsg.IsRead) && utils2.Contain(latestMsg.Seq, seqs...) {
			latestMsg.IsRead = true
			conversation.LatestMsg = utils.StructToJsonString(&latestMsg)
			_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{ConID: conversation.ConversationID, Action: constant.AddConOrUpLatMsg, Args: *conversation}, c.GetCh())
		}
	} else {
		if err := c.db.UpdateColumnsConversation(ctx, conversation.ConversationID, map[string]interface{}{"unread_count": 0}); err != nil {
			log.ZError(ctx, "UpdateColumnsConversation err", err, "conversationID", conversation.ConversationID)
		}
	}

	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: conversation.ConversationID, Action: constant.ConChange, Args: []string{conversation.ConversationID}}})
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.TotalUnreadMessageChanged}})

}

func (c *Conversation) doReadDrawing(ctx context.Context, msg *sdkws.MsgData) {
	tips := &sdkws.MarkAsReadTips{}
	err := utils.UnmarshalNotificationElem(msg.Content, tips)
	if err != nil {
		log.ZWarn(ctx, "UnmarshalNotificationElem err", err, "msg", msg)
		return
	}
	log.ZDebug(ctx, "do readDrawing", "tips", tips)
	conversation, err := c.db.GetConversation(ctx, tips.ConversationID)
	if err != nil {
		log.ZError(ctx, "GetConversation err", err, "conversationID", tips.ConversationID)
		return
	}
	if tips.MarkAsReadUserID != c.loginUserID {
		if len(tips.Seqs) == 0 {
			return
		}
		messages, err := c.db.GetMessagesBySeqs(ctx, tips.ConversationID, tips.Seqs)
		if err != nil {
			log.ZError(ctx, "GetMessagesBySeqs err", err, "conversationID", tips.ConversationID, "seqs", tips.Seqs)
			return
		}
		if conversation.ConversationType == constant.SingleChatType {
			var latestMsg sdkws.MsgData
			if err := json.Unmarshal([]byte(conversation.LatestMsg), &latestMsg); err != nil {
				log.ZError(ctx, "Unmarshal err", err, "conversationID", tips.ConversationID, "latestMsg", conversation.LatestMsg)
			}
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
					if latestMsg.ClientMsgID == message.ClientMsgID {
						if latestMsg.IsRead != message.IsRead {
							conversation.LatestMsg = utils.StructToJsonString(message)
							_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{ConID: conversation.ConversationID, Action: constant.AddConOrUpLatMsg, Args: *conversation}, c.GetCh())
						}
					}
					successMsgIDs = append(successMsgIDs, message.ClientMsgID)
				}
			}
			var messageReceiptResp = []*sdk_struct.MessageReceipt{{UserID: tips.MarkAsReadUserID, MsgIDList: successMsgIDs,
				SessionType: conversation.ConversationType, ReadTime: msg.SendTime}}
			c.msgListener().OnRecvC2CReadReceipt(utils.StructToJsonString(messageReceiptResp))
		}
		//else if conversation.ConversationType == constant.SuperGroupChatType {
		//	var successMsgIDs []string
		//	for _, message := range messages {
		//		attachInfo := sdk_struct.AttachedInfoElem{}
		//		_ = utils.JsonStringToStruct(message.AttachedInfo, &attachInfo)
		//		attachInfo.HasReadTime = msg.SendTime
		//		attachInfo.GroupHasReadInfo.HasReadUserIDList = utils.RemoveRepeatedStringInList(append(attachInfo.GroupHasReadInfo.HasReadUserIDList, tips.MarkAsReadUserID))
		//		attachInfo.GroupHasReadInfo.HasReadCount = int32(len(attachInfo.GroupHasReadInfo.HasReadUserIDList))
		//		message.AttachedInfo = utils.StructToJsonString(attachInfo)
		//		if err = c.db.UpdateMessage(ctx, tips.ConversationID, message); err != nil {
		//			log.ZError(ctx, "UpdateMessage err", err, "conversationID", tips.ConversationID, "message", message)
		//		} else {
		//			successMsgIDs = append(successMsgIDs, message.ClientMsgID)
		//		}
		//	}
		//	var messageReceiptResp = []*sdk_struct.MessageReceipt{{GroupID: conversation.GroupID, MsgIDList: successMsgIDs,
		//		SessionType: conversation.ConversationType, ReadTime: msg.SendTime}}
		//	c.msgListener.OnRecvGroupReadReceipt(utils.StructToJsonString(messageReceiptResp))
		//}
	} else {
		c.doUnreadCount(ctx, conversation, tips.HasReadSeq, tips.Seqs)
	}
}
