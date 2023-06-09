package conversation_msg

import (
	"context"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/syncer"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
)

func (c *Conversation) SyncConversations(ctx context.Context) error {
	ccTime := time.Now()
	conversationsOnServer, err := c.getServerConversationList(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "get server cost time", "cost time", time.Since(ccTime), "conversation on server", conversationsOnServer)
	conversationsOnLocal, err := c.db.GetAllConversations(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "get local cost time", "cost time", time.Since(ccTime), "conversation on local", conversationsOnLocal)
	for _, v := range conversationsOnServer {
		c.addFaceURLAndName(ctx, v)
	}
	if err = c.conversationSyncer.Sync(ctx, conversationsOnServer, conversationsOnLocal, func(ctx context.Context, state int, server, local *model_struct.LocalConversation) error {
		if state == syncer.Update {
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
	//seqs, err := c.getServerHasReadAndMaxSeqs(ctx)
	//if err != nil {
	//	return err
	//}
	//if len(seqs) == 0 {
	//	return nil
	//}
	//var conversations []*model_struct.LocalConversation
	//for conversationID, v := range seqs {
	//	var unreadCount int32
	//	c.maxSeqRecorder.Set(conversationID, v.MaxSeq)
	//	if v.MaxSeq-v.HasReadSeq < 0 {
	//		unreadCount = 0
	//	} else {
	//		unreadCount = int32(v.MaxSeq - v.HasReadSeq)
	//	}
	//	conversations = append(conversations, &model_struct.LocalConversation{
	//		ConversationID: conversationID,
	//		UnreadCount:    unreadCount,
	//	})
	//}
	//return c.db.UpdateOrCreateConversations(ctx, conversations)
	return nil
}
