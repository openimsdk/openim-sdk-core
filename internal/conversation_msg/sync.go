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
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/syncer"
	"time"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
)

func (c *Conversation) SyncConversationsAndTriggerCallback(ctx context.Context, conversationsOnServer []*model_struct.LocalConversation) error {
	conversationsOnLocal, err := c.db.GetAllConversations(ctx)
	if err != nil {
		return err
	}
	if err := c.batchAddFaceURLAndName(ctx, conversationsOnServer...); err != nil {
		return err
	}
	if err = c.conversationSyncer.Sync(ctx, conversationsOnServer, conversationsOnLocal, func(ctx context.Context, state int, server, local *model_struct.LocalConversation) error {
		if state == syncer.Update || state == syncer.Insert {
			c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: server.ConversationID, Action: constant.ConChange, Args: []string{server.ConversationID}}})
		}
		return nil
	}, true); err != nil {
		return err
	}
	conversationsOnLocal, err = c.db.GetAllConversations(ctx)
	if err != nil {
		return err
	}
	c.cache.UpdateConversations(conversationsOnLocal)
	return nil
}

func (c *Conversation) SyncConversations(ctx context.Context, conversationIDs []string) error {
	conversationsOnServer, err := c.getServerConversationsByIDs(ctx, conversationIDs)
	if err != nil {
		return err
	}
	return c.SyncConversationsAndTriggerCallback(ctx, conversationsOnServer)
}

func (c *Conversation) SyncAllConversations(ctx context.Context) error {
	ccTime := time.Now()
	conversationsOnServer, err := c.getServerConversationList(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "get server cost time", "cost time", time.Since(ccTime), "conversation on server", conversationsOnServer)
	return c.SyncConversationsAndTriggerCallback(ctx, conversationsOnServer)
}

func (c *Conversation) SyncConversationUnreadCount(ctx context.Context) error {
	var conversationChangedList []string
	allConversations := c.cache.GetAllHasUnreadMessageConversations()
	log.ZDebug(ctx, "get unread message length", "len", len(allConversations))
	for _, conversation := range allConversations {
		if deleteRows := c.db.DeleteConversationUnreadMessageList(ctx, conversation.ConversationID, conversation.UpdateUnreadCountTime); deleteRows > 0 {
			log.ZDebug(ctx, "DeleteConversationUnreadMessageList", conversation.ConversationID, conversation.UpdateUnreadCountTime, "delete rows:", deleteRows)
			if err := c.db.DecrConversationUnreadCount(ctx, conversation.ConversationID, deleteRows); err != nil {
				log.ZDebug(ctx, "DecrConversationUnreadCount", conversation.ConversationID, conversation.UpdateUnreadCountTime, "decr unread count err:", err.Error())
			} else {
				conversationChangedList = append(conversationChangedList, conversation.ConversationID)
			}
		}
	}
	if len(conversationChangedList) > 0 {
		if err := common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{Action: constant.ConChange, Args: conversationChangedList}, c.GetCh()); err != nil {
			return err
		}
	}
	return nil
}

func (c *Conversation) SyncConversationHashReadSeqs(ctx context.Context) error {
	log.ZDebug(ctx, "start SyncConversationHashReadSeqs")
	seqs, err := c.getServerHasReadAndMaxSeqs(ctx)
	if err != nil {
		return err
	}
	if len(seqs) == 0 {
		return nil
	}
	var conversations []*model_struct.LocalConversation
	var conversationIDs []string
	var conversationIDsNeedSync []string
	for conversationID, v := range seqs {
		var unreadCount int32
		c.maxSeqRecorder.Set(conversationID, v.MaxSeq)
		if v.MaxSeq-v.HasReadSeq < 0 {
			unreadCount = 0
		} else {
			unreadCount = int32(v.MaxSeq - v.HasReadSeq)
		}
		if err := c.db.UpdateColumnsConversation(ctx, conversationID, map[string]interface{}{"unread_count": unreadCount, "has_read_seq": v.HasReadSeq}); err != nil {
			if errs.Unwrap(err) == errs.ErrRecordNotFound {
				conversationIDsNeedSync = append(conversationIDsNeedSync, conversationID)
			} else {
				log.ZWarn(ctx, "UpdateColumnsConversation err", err, "conversationID", conversationID)
			}
			continue
		}
		conversationIDs = append(conversationIDs, conversationID)
	}
	if len(conversationIDsNeedSync) > 0 {
		if err := c.SyncConversations(ctx, conversationIDsNeedSync); err != nil {
			log.ZWarn(ctx, "sync new conversations failed", nil, "conversationIDs", conversationIDsNeedSync)
		} else {
			conversationIDs = append(conversationIDs, conversationIDsNeedSync...)
		}
	}

	log.ZDebug(ctx, "update conversations", "conversations", conversations)
	if len(conversationIDs) > 0 {
		common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{Action: constant.ConChange, Args: conversationIDs}, c.GetCh())
		common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{Action: constant.TotalUnreadMessageChanged}, c.GetCh())
	}
	return nil
}
