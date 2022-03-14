package server_api_params

type Conversation struct {
	OwnerUserID      string `json:"ownerUserID" binding:"required"`
	ConversationID   string `json:"conversationID"`
	ConversationType int32  `json:"conversationType"`
	UserID           string `json:"userID"`
	GroupID          string `json:"groupID"`
	RecvMsgOpt       int32  `json:"recvMsgOpt"`
	UnreadCount      int32  `json:"unreadCount"`
	DraftTextTime    int64  `json:"draftTextTime"`
	IsPinned         bool   `json:"isPinned"`
	IsPrivateChat    bool   `json:"isPrivateChat"`
	AttachedInfo     string `json:"attachedInfo"`
	Ex               string `json:"ex"`
}

type SetConversationReq struct {
	Conversation
	OperationID string `json:"operationID" binding:"required"`
}

type SetConversationResp struct {
}

type BatchSetConversationsReq struct {
	Conversations []Conversation `json:"conversations" binding:"required"`
	OwnerUserID   string         `json:"ownerUserID" binding:"required"`
	OperationID   string         `json:"operationID" binding:"required"`
}

type BatchSetConversationsResp struct {
	Success []string `json:"success"`
	Failed  []string `json:"failed"`
}

type GetConversationReq struct {
	ConversationID string `json:"conversationID" binding:"required"`
	OwnerUserID    string `json:"ownerUserID" binding:"required"`
	OperationID    string `json:"operationID" binding:"required"`
}

type GetConversationResp struct {
	Conversation Conversation `json:"data"`
}

type GetAllConversationsReq struct {
	OwnerUserID string `json:"ownerUserID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}

type GetAllConversationsResp struct {
	Conversations []Conversation `json:"data"`
}

type GetConversationsReq struct {
	ConversationIDs []string `json:"conversationIDs" binding:"required"`
	OwnerUserID     string   `json:"ownerUserID" binding:"required"`
	OperationID     string   `json:"operationID" binding:"required"`
}

type GetConversationsResp struct {
	Conversations []Conversation `json:"data"`
}
