package server_api_params

type GetLocalKeyReq struct {
	SessionType int32  `json:"sessionType" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
	FriendID    string `json:"friendID"`
	GroupID     string `json:"groupID"`
	OperationID string `json:"operationID" binding:"required"`
}
type GetLocalKeyResp struct {
	SessionInfo string `json:"sessionInfo"`
}
type GetAllLocalKeyBySessionIDReq struct {
	UserID      string `json:"userID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}
type GetAllLocalKeyBySessionIDResp struct {
	Keys []Key `json:"keys"`
}
type Key struct {
	SessionID   string ` json:"sessionID"`
	SessionType int32  `json:"sessionType" `
	SessionKey  string `json:"sessionKey" `
}
type WsMarkGroupMessageAsReadReq struct {
	GroupID   string `json:"groupID"`
	MsgIDList string `json:"msgIDList"`
}
