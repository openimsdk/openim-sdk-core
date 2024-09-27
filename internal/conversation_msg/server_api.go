package conversation_msg

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	pbConversation "github.com/openimsdk/protocol/conversation"
	pbMsg "github.com/openimsdk/protocol/msg"
)

func (c *Conversation) markMsgAsRead2Server(ctx context.Context, conversationID string, seqs []int64) error {
	req := &pbMsg.MarkMsgsAsReadReq{UserID: c.loginUserID, ConversationID: conversationID, Seqs: seqs}
	return api.MarkMsgsAsRead.Execute(ctx, req)
}

func (c *Conversation) markConversationAsReadServer(ctx context.Context, conversationID string, hasReadSeq int64, seqs []int64) error {
	req := &pbMsg.MarkConversationAsReadReq{UserID: c.loginUserID, ConversationID: conversationID, HasReadSeq: hasReadSeq, Seqs: seqs}
	return api.MarkConversationAsRead.Execute(ctx, req)
}

func (c *Conversation) setConversationHasReadSeq(ctx context.Context, conversationID string, hasReadSeq int64) error {
	req := &pbMsg.SetConversationHasReadSeqReq{UserID: c.loginUserID, ConversationID: conversationID, HasReadSeq: hasReadSeq}
	return api.SetConversationHasReadSeq.Execute(ctx, req)
}

// To delete session information, delete the server first, and then invoke the interface.
// The client receives a callback to delete all local information.
func (c *Conversation) clearConversationMsgFromServer(ctx context.Context, conversationID string) error {
	req := &pbMsg.ClearConversationsMsgReq{UserID: c.loginUserID, ConversationIDs: []string{conversationID}}
	return api.ClearConversationMsg.Execute(ctx, req)
}

// Delete all server messages
func (c *Conversation) deleteAllMessageFromServer(ctx context.Context) error {
	req := &pbMsg.UserClearAllMsgReq{UserID: c.loginUserID}
	return api.ClearAllMsg.Execute(ctx, req)
}

// The user deletes part of the message from the server
func (c *Conversation) deleteMessagesFromServer(ctx context.Context, conversationID string, seqs []int64) error {
	req := &pbMsg.DeleteMsgsReq{UserID: c.loginUserID, Seqs: seqs, ConversationID: conversationID}
	return api.DeleteMsgs.Execute(ctx, req)
}

func (c *Conversation) revokeMessageFromServer(ctx context.Context, conversationID string, seq int64) error {
	req := &pbMsg.RevokeMsgReq{UserID: c.loginUserID, ConversationID: conversationID, Seq: seq}
	return api.RevokeMsg.Execute(ctx, req)
}

func (c *Conversation) getHasReadAndMaxSeqsFromServer(ctx context.Context, conversationIDs ...string) (*pbMsg.GetConversationsHasReadAndMaxSeqResp, error) {
	req := pbMsg.GetConversationsHasReadAndMaxSeqReq{UserID: c.loginUserID, ConversationIDs: conversationIDs}
	return api.GetConversationsHasReadAndMaxSeq.Invoke(ctx, &req)
}

func (c *Conversation) getConversationsByIDsFromServer(ctx context.Context, conversations []string) (*pbConversation.GetConversationsResp, error) {
	req := &pbConversation.GetConversationsReq{OwnerUserID: c.loginUserID, ConversationIDs: conversations}
	return api.GetConversations.Invoke(ctx, req)
}

func (c *Conversation) getAllConversationListFromServer(ctx context.Context) (*pbConversation.GetAllConversationsResp, error) {
	req := &pbConversation.GetAllConversationsReq{OwnerUserID: c.loginUserID}
	return api.GetAllConversations.Invoke(ctx, req)
}

func (c *Conversation) getAllConversationIDsFromServer(ctx context.Context) (*pbConversation.GetFullOwnerConversationIDsResp, error) {
	req := &pbConversation.GetFullOwnerConversationIDsReq{UserID: c.loginUserID}
	return api.GetFullConversationIDs.Invoke(ctx, req)
}

func (c *Conversation) getIncrementalConversationFromServer(ctx context.Context, version uint64, versionID string) (*pbConversation.GetIncrementalConversationResp, error) {
	req := &pbConversation.GetIncrementalConversationReq{UserID: c.loginUserID, Version: version, VersionID: versionID}
	return api.GetIncrementalConversation.Invoke(ctx, req)
}
