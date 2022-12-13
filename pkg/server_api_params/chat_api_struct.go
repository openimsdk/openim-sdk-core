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

type GetMessageListReactionExtensionsReq struct {
	OperationID            string                                    `json:"operationID" binding:"required"`
	SourceID               string                                    `json:"sourceID"  binding:"required"`
	SessionType            int32                                     `json:"sessionType" binding:"required"`
	MessageReactionKeyList []OperateMessageListReactionExtensionsReq `json:"messageReactionKeyList" binding:"required"`
}

type KeyValueResp struct {
	KeyValue             *KeyValue `protobuf:"bytes,1,opt,name=keyValue" json:"keyValue,omitempty"`
	ErrCode              int32     `protobuf:"varint,2,opt,name=errCode" json:"errCode,omitempty"`
	ErrMsg               string    `protobuf:"bytes,3,opt,name=errMsg" json:"errMsg,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

type ExtendMsg struct {
	ReactionExtensionList map[string]*KeyValueResp `protobuf:"bytes,1,rep,name=reactionExtensionList" json:"reactionExtensionList,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	ClientMsgID           string                   `protobuf:"bytes,2,opt,name=clientMsgID" json:"clientMsgID,omitempty"`
	MsgFirstModifyTime    int64                    `protobuf:"varint,3,opt,name=msgFirstModifyTime" json:"msgFirstModifyTime,omitempty"`
	AttachedInfo          string                   `protobuf:"bytes,4,opt,name=attachedInfo" json:"attachedInfo,omitempty"`
	Ex                    string                   `protobuf:"bytes,5,opt,name=ex" json:"ex,omitempty"`
}

type ExtendMsgResp struct {
	ExtendMsg *ExtendMsg `protobuf:"bytes,1,opt,name=extendMsg" json:"extendMsg,omitempty"`
	ErrCode   int32      `protobuf:"varint,2,opt,name=errCode" json:"errCode,omitempty"`
	ErrMsg    string     `protobuf:"bytes,3,opt,name=errMsg" json:"errMsg,omitempty"`
}

type GetMessageListReactionExtensionsResp struct {
	SuccessList []*ExtendMsgResp `json:"successList"`
	FailedList  []*ExtendMsgResp `json:"failedList"`
}

type OperateMessageListReactionExtensionsReq struct {
	ClientMsgID        string `json:"clientMsgID"`
	MsgFirstModifyTime int64  `json:"msgFirstModifyTime"`
}
