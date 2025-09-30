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
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/tools/log"
)

func (c *Conversation) SyncAllConversationHashReadSeqs(ctx context.Context) error {
	startTime := time.Now()
	log.ZDebug(ctx, "start SyncConversationHashReadSeqs")

	resp := msg.GetConversationsHasReadAndMaxSeqResp{}
	req := msg.GetConversationsHasReadAndMaxSeqReq{UserID: c.loginUserID}
	err := c.SendReqWaitResp(ctx, &req, constant.GetConvMaxReadSeq, &resp)
	if err != nil {
		log.ZWarn(ctx, "SendReqWaitResp err", err)
		return err
	}
	seqs := resp.Seqs
	log.ZDebug(ctx, "getServerHasReadAndMaxSeqs completed", "duration", time.Since(startTime).Seconds())

	if len(seqs) == 0 {
		return nil
	}
	var conversationChangedIDs []string
	var conversationIDsNeedSync []string

	stepStartTime := time.Now()
	conversationsOnLocal, err := c.db.GetAllConversations(ctx)
	if err != nil {
		log.ZWarn(ctx, "get all conversations err", err)
		return err
	}
	log.ZDebug(ctx, "GetAllConversations completed", "duration", time.Since(stepStartTime).Seconds())

	conversationsOnLocalMap := datautil.SliceToMap(conversationsOnLocal, func(e *model_struct.LocalConversation) string {
		return e.ConversationID
	})

	stepStartTime = time.Now()
	for conversationID, v := range seqs {
		var unreadCount int32
		c.maxSeqRecorder.Set(conversationID, v.MaxSeq)
		if v.MaxSeq-v.HasReadSeq < 0 {
			unreadCount = 0
			log.ZWarn(ctx, "unread count is less than 0", nil, "conversationID",
				conversationID, "maxSeq", v.MaxSeq, "hasReadSeq", v.HasReadSeq)
		} else {
			unreadCount = int32(v.MaxSeq - v.HasReadSeq)
		}
		if conversation, ok := conversationsOnLocalMap[conversationID]; ok {
			if conversation.UnreadCount != unreadCount {
				if err := c.db.UpdateColumnsConversation(ctx, conversationID, map[string]interface{}{"unread_count": unreadCount}); err != nil {
					log.ZWarn(ctx, "UpdateColumnsConversation err", err, "conversationID", conversationID)
					continue
				}
				conversationChangedIDs = append(conversationChangedIDs, conversationID)
			}
		} else {
			conversationIDsNeedSync = append(conversationIDsNeedSync, conversationID)
		}
	}
	log.ZDebug(ctx, "Process seqs completed", "duration", time.Since(stepStartTime).Seconds())

	if len(conversationIDsNeedSync) > 0 {
		stepStartTime = time.Now()
		r, err := c.getConversationsByIDsFromServer(ctx, conversationIDsNeedSync)
		if err != nil {
			log.ZWarn(ctx, "getServerConversationsByIDs err", err, "conversationIDs", conversationIDsNeedSync)
			return err
		}
		log.ZDebug(ctx, "getServerConversationsByIDs completed", "duration", time.Since(stepStartTime).Seconds())
		conversationsOnServer := datautil.Batch(ServerConversationToLocal, r.Conversations)
		stepStartTime = time.Now()
		if err := c.batchAddFaceURLAndName(ctx, conversationsOnServer...); err != nil {
			log.ZWarn(ctx, "batchAddFaceURLAndName err", err, "conversationsOnServer", conversationsOnServer)
			return err
		}
		log.ZDebug(ctx, "batchAddFaceURLAndName completed", "duration", time.Since(stepStartTime).Seconds())

		for _, conversation := range conversationsOnServer {
			var unreadCount int32
			v, ok := seqs[conversation.ConversationID]
			if !ok {
				continue
			}
			if v.MaxSeq-v.HasReadSeq < 0 {
				unreadCount = 0
				log.ZWarn(ctx, "unread count is less than 0", nil, "server seq", v, "conversation", conversation)
			} else {
				unreadCount = int32(v.MaxSeq - v.HasReadSeq)
			}
			conversation.UnreadCount = unreadCount
		}

		stepStartTime = time.Now()
		err = c.db.BatchInsertConversationList(ctx, conversationsOnServer)
		if err != nil {
			log.ZWarn(ctx, "BatchInsertConversationList err", err, "conversationsOnServer", conversationsOnServer)
		}
		log.ZDebug(ctx, "BatchInsertConversationList completed", "duration", time.Since(stepStartTime).Seconds())
	}

	log.ZDebug(ctx, "update conversations", "conversations", conversationChangedIDs)
	if len(conversationChangedIDs) > 0 {
		stepStartTime = time.Now()
		common.DispatchUpdateConversation(ctx, common.UpdateConNode{Action: constant.ConChange, Args: conversationChangedIDs}, c.ConversationEventQueue())
		common.DispatchUpdateConversation(ctx, common.UpdateConNode{Action: constant.TotalUnreadMessageChanged}, c.ConversationEventQueue())
		log.ZDebug(ctx, "TriggerCmdUpdateConversation completed", "duration", time.Since(stepStartTime).Seconds())
	}

	stepStartTime = time.Now()
	if err := c.syncAllGroupReadCursors(ctx); err != nil {
		log.ZWarn(ctx, "syncAllGroupReadCursors failed", err)
	}
	log.ZDebug(ctx, "syncAllGroupReadCursors completed", "duration", time.Since(stepStartTime).Seconds())

	log.ZDebug(ctx, "SyncAllConversationHashReadSeqs completed", "totalDuration", time.Since(startTime).Seconds())
	return nil
}

func (c *Conversation) syncAllGroupReadCursors(ctx context.Context) error {
	conversations, err := c.db.GetAllConversations(ctx)
	if err != nil {
		log.ZWarn(ctx, "GetAllConversations failed", err)
		return err
	}

	var groupConversationIDs []string
	for _, conv := range conversations {
		if conv.ConversationType == constant.ReadGroupChatType {
			groupConversationIDs = append(groupConversationIDs, conv.ConversationID)
		}
	}

	if len(groupConversationIDs) == 0 {
		log.ZDebug(ctx, "no group conversations to sync cursors")
		return nil
	}

	log.ZDebug(ctx, "found group conversations", "count", len(groupConversationIDs), "conversationIDs", groupConversationIDs)

	for _, conversationID := range groupConversationIDs {
		if _, err := c.db.GetGroupReadCursorState(ctx, conversationID); err != nil {
			if ierr := c.db.InsertGroupReadCursorState(ctx, &model_struct.LocalGroupReadCursorState{ConversationID: conversationID, CursorVersion: 1}); ierr != nil {
				log.ZWarn(ctx, "InsertGroupReadCursorState failed", ierr, "conversationID", conversationID)
			} else {
				log.ZDebug(ctx, "initialized LocalGroupReadCursorState", "conversationID", conversationID)
			}
		}
	}

	stepStartTime := time.Now()
	resp, err := c.getConversationReadCursors(ctx, groupConversationIDs)
	if err != nil {
		log.ZWarn(ctx, "getConversationReadCursorsFromServer failed", err)
		return err
	}
	log.ZDebug(ctx, "getConversationReadCursorsFromServer completed", "duration", time.Since(stepStartTime).Seconds())

	stepStartTime = time.Now()
	allCursorCount := 0
	for conversationID, cursorList := range resp.Cursors {
		curCursorCount := 0
		if cursorList == nil || len(cursorList.Cursors) == 0 {
			continue
		}

		for _, cursor := range cursorList.Cursors {
			localCursor := &model_struct.LocalGroupReadCursor{
				ConversationID: conversationID,
				UserID:         cursor.UserID,
				MaxReadSeq:     cursor.MaxReadSeq,
			}

			existingCursor, err := c.db.GetGroupReadCursor(ctx, conversationID, cursor.UserID)
			if err != nil {
				if err := c.db.InsertGroupReadCursor(ctx, localCursor); err != nil {
					log.ZWarn(ctx, "InsertGroupReadCursor failed", err, "conversationID", conversationID, "userID", cursor.UserID)
				} else {
					curCursorCount++
				}
			} else {
				if cursor.MaxReadSeq > existingCursor.MaxReadSeq {
					if err := c.db.UpdateGroupReadCursor(ctx, conversationID, cursor.UserID, cursor.MaxReadSeq); err != nil {
						log.ZWarn(ctx, "UpdateGroupReadCursor failed", err, "conversationID", conversationID, "userID", cursor.UserID)
					} else {
						curCursorCount++
					}
				}
			}
		}
		allCursorCount += curCursorCount
		if curCursorCount != 0 {
			if err := c.db.IncrementGroupReadCursorVersion(ctx, conversationID); err != nil {
				log.ZWarn(ctx, "IncrementGroupReadCursorVersion failed", err, "conversationID", conversationID)
			}
		}
	}

	log.ZDebug(ctx, "syncAllGroupReadCursors completed", "duration", time.Since(stepStartTime).Seconds(), "cursorCount", allCursorCount)
	return nil
}
