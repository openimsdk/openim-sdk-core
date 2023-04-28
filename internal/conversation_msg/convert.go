// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package conversation_msg

import (
	"open_im_sdk/pkg/db/model_struct"

	pbConversation "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
)

func ServerConversationToLocal(conversation *pbConversation.Conversation) *model_struct.LocalConversation {
	return &model_struct.LocalConversation{
		ConversationID:        conversation.ConversationID,
		ConversationType:      conversation.ConversationType,
		UserID:                conversation.UserID,
		GroupID:               conversation.GroupID,
		RecvMsgOpt:            conversation.RecvMsgOpt,
		UnreadCount:           conversation.UnreadCount,
		GroupAtType:           conversation.GroupAtType,
		IsPinned:              conversation.IsPinned,
		DraftTextTime:         conversation.DraftTextTime,
		IsNotInGroup:          conversation.IsNotInGroup,
		BurnDuration:          conversation.BurnDuration,
		IsPrivateChat:         conversation.IsPrivateChat,
		UpdateUnreadCountTime: conversation.UpdateUnreadCountTime,
		AttachedInfo:          conversation.AttachedInfo,
		Ex:                    conversation.Ex,
	}
}

func LocalConversationToServer(conversation *model_struct.LocalConversation) *pbConversation.Conversation {
	return &pbConversation.Conversation{
		ConversationID:        conversation.ConversationID,
		ConversationType:      conversation.ConversationType,
		UserID:                conversation.UserID,
		GroupID:               conversation.GroupID,
		RecvMsgOpt:            conversation.RecvMsgOpt,
		UnreadCount:           conversation.UnreadCount,
		GroupAtType:           conversation.GroupAtType,
		IsPinned:              conversation.IsPinned,
		DraftTextTime:         conversation.DraftTextTime,
		IsNotInGroup:          conversation.IsNotInGroup,
		BurnDuration:          conversation.BurnDuration,
		IsPrivateChat:         conversation.IsPrivateChat,
		UpdateUnreadCountTime: conversation.UpdateUnreadCountTime,
		AttachedInfo:          conversation.AttachedInfo,
		Ex:                    conversation.Ex,
	}
}
