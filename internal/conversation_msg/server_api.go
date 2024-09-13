package conversation_msg

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	pbConversation "github.com/openimsdk/protocol/conversation"
	pbMsg "github.com/openimsdk/protocol/msg"
)

func (c *Conversation) markMsgAsRead2Server(ctx context.Context, conversationID string, seqs []int64) error {
	req := &pbMsg.MarkMsgsAsReadReq{UserID: c.loginUserID, ConversationID: conversationID, Seqs: seqs}
	return util.ApiPost(ctx, constant.MarkMsgsAsReadRouter, req, nil)
}

func (c *Conversation) markConversationAsReadServer(ctx context.Context, conversationID string, hasReadSeq int64, seqs []int64) error {
	req := &pbMsg.MarkConversationAsReadReq{UserID: c.loginUserID, ConversationID: conversationID, HasReadSeq: hasReadSeq, Seqs: seqs}
	return util.ApiPost(ctx, constant.MarkConversationAsRead, req, nil)
}

func (c *Conversation) setConversationHasReadSeq(ctx context.Context, conversationID string, hasReadSeq int64) error {
	req := &pbMsg.SetConversationHasReadSeqReq{UserID: c.loginUserID, ConversationID: conversationID, HasReadSeq: hasReadSeq}
	return util.ApiPost(ctx, constant.SetConversationHasReadSeq, req, nil)
}

// To delete session information, delete the server first, and then invoke the interface.
// The client receives a callback to delete all local information.
func (c *Conversation) clearConversationMsgFromServer(ctx context.Context, conversationID string) error {
	req := &pbMsg.ClearConversationsMsgReq{UserID: c.loginUserID, ConversationIDs: []string{conversationID}}
	return util.ApiPost(ctx, constant.ClearConversationMsgRouter, req, nil)
}

// Delete all server messages
func (c *Conversation) deleteAllMessageFromServer(ctx context.Context) error {
	req := &pbMsg.UserClearAllMsgReq{UserID: c.loginUserID}
	return util.ApiPost(ctx, constant.ClearAllMsgRouter, req, nil)

}

// The user deletes part of the message from the server
func (c *Conversation) deleteMessagesFromServer(ctx context.Context, conversationID string, seqs []int64) error {
	req := &pbMsg.DeleteMsgsReq{UserID: c.loginUserID, Seqs: seqs, ConversationID: conversationID}
	return util.ApiPost(ctx, constant.DeleteMsgsRouter, req, nil)

}

func (c *Conversation) revokeMessageFromServer(ctx context.Context, conversationID string, seq int64) error {
	req := &pbMsg.RevokeMsgReq{UserID: c.loginUserID, ConversationID: conversationID, Seq: seq}
	return util.ApiPost(ctx, constant.RevokeMsgRouter, req, nil)
}

func (c *Conversation) getHasReadAndMaxSeqsFromServer(ctx context.Context, conversationIDs ...string) (*pbMsg.GetConversationsHasReadAndMaxSeqResp, error) {
	resp := &pbMsg.GetConversationsHasReadAndMaxSeqResp{}
	req := pbMsg.GetConversationsHasReadAndMaxSeqReq{UserID: c.loginUserID, ConversationIDs: conversationIDs}
	return resp, util.ApiPost(ctx, constant.GetConversationsHasReadAndMaxSeqRouter, &req, resp)
}

func (c *Conversation) getConversationsByIDsFromServer(ctx context.Context, conversations []string) (*pbConversation.GetConversationsResp, error) {
	resp, err := util.CallApi[pbConversation.GetConversationsResp](ctx, constant.GetConversationsRouter,
		pbConversation.GetConversationsReq{OwnerUserID: c.loginUserID, ConversationIDs: conversations})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Conversation) getAllConversationListFromServer(ctx context.Context) (*pbConversation.GetAllConversationsResp, error) {
	resp, err := util.CallApi[pbConversation.GetAllConversationsResp](ctx, constant.GetAllConversationsRouter,
		pbConversation.GetAllConversationsReq{OwnerUserID: c.loginUserID})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Conversation) getAllConversationIDsFromServer(ctx context.Context) (*pbConversation.GetFullOwnerConversationIDsResp, error) {
	resp, err := util.CallApi[pbConversation.GetFullOwnerConversationIDsResp](ctx, constant.GetFullConversationIDs,
		pbConversation.GetFullOwnerConversationIDsReq{UserID: c.loginUserID})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Conversation) getIncrementalConversationFromServer(ctx context.Context, version uint64, versionID string) (*pbConversation.GetIncrementalConversationResp, error) {
	resp, err := util.CallApi[pbConversation.GetIncrementalConversationResp](ctx, constant.GetIncrementalConversation,
		pbConversation.GetIncrementalConversationReq{UserID: c.loginUserID, Version: version, VersionID: versionID})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
