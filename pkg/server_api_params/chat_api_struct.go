package server_api_params

type DeleteMsgReq struct {
	OpUserID    string   `json:"opUserID"`
	UserID      string   `json:"userID"`
	SeqList     []uint32 `json:"seqList"`
	OperationID string   `json:"operationID"`
}

type DeleteMsgResp struct {
}

type CleanUpMsgReq struct {
	UserID      string `json:"userID"  binding:"required"`
	OperationID string `json:"operationID"  binding:"required"`
}

type CleanUpMsgResp struct {
	CommResp
}
type DelSuperGroupMsgReq struct {
	UserID      string   `json:"userID,omitempty" binding:"required"`
	GroupID     string   `json:"groupID,omitempty" binding:"required"`
	SeqList     []uint32 `json:"seqList,omitempty"`
	IsAllDelete bool     `json:"isAllDelete"`
	OperationID string   `json:"operationID,omitempty" binding:"required"`
}
type DelSuperGroupMsgResp struct {
	CommResp
}
type MsgDeleteNotificationElem struct {
	GroupID     string   `json:"groupID"`
	IsAllDelete bool     `json:"isAllDelete"`
	SeqList     []uint32 `json:"seqList"`
}
type SetMessageReactionExtensionsReq struct {
	OperationID           string               `json:"operationID" validate:"required"`
	ClientMsgID           string               `json:"clientMsgID" validate:"required"`
	SourceID              string               `json:"sourceID" validate:"required"`
	SessionType           int32                `json:"sessionType" validate:"required"`
	ReactionExtensionList map[string]*KeyValue `json:"reactionExtensionList"`
	IsReact               bool                 `json:"isReact,omitempty"`
	IsExternalExtensions  bool                 `json:"isExternalExtensions,omitempty"`
	MsgFirstModifyTime    int64                `json:"msgFirstModifyTime,omitempty"`
}
type KeyValue struct {
	TypeKey          string `json:"typeKey" validate:"required"`
	Value            string `json:"value" validate:"required"`
	LatestUpdateTime int64  `json:"latestUpdateTime"`
}
type SetMessageReactionExtensionsResp struct {
	CommResp
	ReactionExtensionListResult []*ExtensionResult
	Data                        map[string]interface{} `json:"data"`
}
type ExtensionResult struct {
	CommResp
	KeyValue
}
