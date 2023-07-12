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

package server_api_params

type DeleteMsgReq struct {
	OpUserID    string  `json:"opUserID"`
	UserID      string  `json:"userID"`
	SeqList     []int64 `json:"seqList"`
	OperationID string  `json:"operationID"`
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
	UserID      string  `json:"userID,omitempty" binding:"required"`
	GroupID     string  `json:"groupID,omitempty" binding:"required"`
	SeqList     []int64 `json:"seqList,omitempty"`
	IsAllDelete bool    `json:"isAllDelete"`
	OperationID string  `json:"operationID,omitempty" binding:"required"`
}
type DelSuperGroupMsgResp struct {
	CommResp
}
type MsgDeleteNotificationElem struct {
	GroupID     string  `json:"groupID"`
	IsAllDelete bool    `json:"isAllDelete"`
	SeqList     []int64 `json:"seqList"`
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
type AddMessageReactionExtensionsReq struct {
	OperationID           string               `json:"operationID" validate:"required"`
	ClientMsgID           string               `json:"clientMsgID" validate:"required"`
	SourceID              string               `json:"sourceID" validate:"required"`
	SessionType           int32                `json:"sessionType" validate:"required"`
	ReactionExtensionList map[string]*KeyValue `json:"reactionExtensionList"`
	IsReact               bool                 `json:"isReact,omitempty"`
	IsExternalExtensions  bool                 `json:"isExternalExtensions,omitempty"`
	MsgFirstModifyTime    int64                `json:"msgFirstModifyTime,omitempty"`
	Seq                   int64                `json:"seq"`
}
type DeleteMessageReactionExtensionsReq struct {
	OperationID           string      `json:"operationID" binding:"required"`
	SourceID              string      `json:"sourceID" binding:"required"`
	SessionType           int32       `json:"sessionType" binding:"required"`
	ClientMsgID           string      `json:"clientMsgID" binding:"required"`
	IsExternalExtensions  bool        `json:"isExternalExtensions"`
	MsgFirstModifyTime    int64       `json:"msgFirstModifyTime" binding:"required"`
	ReactionExtensionList []*KeyValue `json:"reactionExtensionList" binding:"required"`
}
type DeleteMessageReactionExtensionsResp struct {
	CommResp
	Result []*ExtensionResult
	Data   map[string]interface{} `json:"data"`
}
type KeyValue struct {
	TypeKey          string `json:"typeKey" validate:"required"`
	Value            string `json:"value" validate:"required"`
	LatestUpdateTime int64  `json:"latestUpdateTime"`
}
type ApiResult struct {
	Result             []*ExtensionResult `json:"result"`
	MsgFirstModifyTime int64              `json:"msgFirstModifyTime"`
	IsReact            bool               `json:"isReact"`
}
type SetMessageReactionExtensionsResp struct {
	CommResp
	ApiResult struct {
		Result             []*ExtensionResult `json:"result"`
		MsgFirstModifyTime int64              `json:"msgFirstModifyTime"`
		IsReact            bool               `json:"isReact"`
	}
	Data map[string]interface{} `json:"data"`
}
type AddMessageReactionExtensionsResp struct {
	CommResp
	ApiResult struct {
		Result             []*ExtensionResult `json:"result"`
		MsgFirstModifyTime int64              `json:"msgFirstModifyTime"`
		IsReact            bool               `json:"isReact"`
	}
	Data map[string]interface{} `json:"data"`
}
type ExtensionResult struct {
	CommResp
	KeyValue
}

type GetMessageListReactionExtensionsReq struct {
	OperationID            string                                    `json:"operationID" binding:"required"`
	SourceID               string                                    `json:"sourceID"  binding:"required"`
	SessionType            int32                                     `json:"sessionType" binding:"required"`
	IsExternalExtensions   bool                                      `json:"isExternalExtensions"`
	TypeKeyList            []string                                  `json:"typeKeyList"`
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

//type ExtendMsg struct {
//	ReactionExtensionList map[string]*KeyValueResp `protobuf:"bytes,1,rep,name=reactionExtensionList" json:"reactionExtensionList,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
//	ClientMsgID           string                   `protobuf:"bytes,2,opt,name=clientMsgID" json:"clientMsgID,omitempty"`
//	MsgFirstModifyTime    int64                    `protobuf:"varint,3,opt,name=msgFirstModifyTime" json:"msgFirstModifyTime,omitempty"`
//	AttachedInfo          string                   `protobuf:"bytes,4,opt,name=attachedInfo" json:"attachedInfo,omitempty"`
//	Ex                    string                   `protobuf:"bytes,5,opt,name=ex" json:"ex,omitempty"`
//}
//
//type ExtendMsgResp struct {
//	ExtendMsg *ExtendMsg `protobuf:"bytes,1,opt,name=extendMsg" json:"extendMsg,omitempty"`
//	ErrCode   int32      `protobuf:"varint,2,opt,name=errCode" json:"errCode,omitempty"`
//	ErrMsg    string     `protobuf:"bytes,3,opt,name=errMsg" json:"errMsg,omitempty"`
//}

type GetMessageListReactionExtensionsResp []*SingleMessageExtensionResult

type OperateMessageListReactionExtensionsReq struct {
	ClientMsgID        string `json:"clientMsgID"`
	MsgFirstModifyTime int64  `json:"msgFirstModifyTime"`
}
type ReactionMessageModifierNotification struct {
	Operation   int    `json:"operation" binding:"required"`
	SourceID    string `json:"sourceID"  binding:"required"`
	OpUserID    string `json:"opUserID"  binding:"required"`
	SessionType int32  `json:"sessionType" binding:"required"`
	// SuccessReactionExtensionList map[string]*sdkws.KeyValue `json:"reactionExtensionList,omitempty" binding:"required"`
	ClientMsgID          string `json:"clientMsgID" binding:"required"`
	IsReact              bool   `json:"isReact"`
	IsExternalExtensions bool   `json:"isExternalExtensions"`
	MsgFirstModifyTime   int64  `json:"msgFirstModifyTime"`
	Seq                  uint32 `json:"seq"`
}
type SingleMessageExtensionResult struct {
	ErrCode int32  `protobuf:"varint,1,opt,name=errCode" json:"errCode,omitempty"`
	ErrMsg  string `protobuf:"bytes,2,opt,name=errMsg" json:"errMsg,omitempty"`
	// ReactionExtensionList map[string]*sdkws.KeyValue `protobuf:"bytes,3,rep,name=reactionExtensionList" json:"reactionExtensionList,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	ClientMsgID          string   `protobuf:"bytes,4,opt,name=clientMsgID" json:"clientMsgID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}
type ReactionMessageDeleteNotification struct {
	SourceID    string `json:"sourceID"  binding:"required"`
	OpUserID    string `json:"opUserID"  binding:"required"`
	SessionType int32  `json:"sessionType" binding:"required"`
	// SuccessReactionExtensionList map[string]*sdkws.KeyValue `json:"reactionExtensionList,omitempty" binding:"required"`
	ClientMsgID        string `json:"clientMsgID" binding:"required"`
	MsgFirstModifyTime int64  `json:"msgFirstModifyTime"`
}
