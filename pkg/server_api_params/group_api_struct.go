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

import "github.com/OpenIMSDK/protocol/sdkws"

type CommResp struct {
	ErrCode int32  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
}
type CommDataResp struct {
	CommResp
	Data []map[string]interface{} `json:"data"`
}
type CommDataRespOne struct {
	CommResp
	Data map[string]interface{} `json:"data"`
}

type KickGroupMemberReq struct {
	GroupID          string   `json:"groupID" binding:"required"`
	KickedUserIDList []string `json:"kickedUserIDList" binding:"required"`
	Reason           string   `json:"reason"`
	OperationID      string   `json:"operationID" binding:"required"`
}
type KickGroupMemberResp struct {
	CommResp
	UserIDResultList []*UserIDResult `json:"data"`
}

type GetGroupMembersInfoReq struct {
	GroupID     string   `json:"groupID" binding:"required"`
	MemberList  []string `json:"memberList" binding:"required"`
	OperationID string   `json:"operationID" binding:"required"`
}
type GetGroupMembersInfoResp struct {
	CommResp
	MemberList []*sdkws.GroupMemberFullInfo `json:"-"`
	Data       []map[string]interface{}     `json:"data"`
}

type InviteUserToGroupReq struct {
	GroupID           string   `json:"groupID" binding:"required"`
	InvitedUserIDList []string `json:"invitedUserIDList" binding:"required"`
	Reason            string   `json:"reason"`
	OperationID       string   `json:"operationID" binding:"required"`
}
type InviteUserToGroupResp struct {
	CommResp
	UserIDResultList []*UserIDResult `json:"data"`
}

type GetJoinedGroupListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}
type GetJoinedGroupListResp struct {
	CommResp
	GroupInfoList []*sdkws.GroupInfo
	Data          []map[string]interface{} `json:"data"`
}

type GetGroupMemberListReq struct {
	GroupID     string `json:"groupID"`
	Filter      int32  `json:"filter"`
	NextSeq     int32  `json:"nextSeq"`
	OperationID string `json:"operationID"`
}
type GetGroupMemberListResp struct {
	CommResp
	NextSeq    int32 `json:"nextSeq"`
	MemberList []*sdkws.GroupMemberFullInfo
	Data       []map[string]interface{} `json:"data"`
}

type GetGroupAllMemberReq struct {
	GroupID     string `json:"groupID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
	Offset      int32  `json:"offset"`
	Count       int32  `json:"count"`
}
type GetGroupAllMemberResp struct {
	CommResp
	MemberList []*sdkws.GroupMemberFullInfo `json:"-"`
	Data       []map[string]interface{}     `json:"data"`
}

type CreateGroupReq struct {
	MemberList   []*GroupAddMemberInfo `json:"memberList"  binding:"required"`
	OwnerUserID  string                `json:"ownerUserID" binding:"required"`
	GroupType    int32                 `json:"groupType"`
	GroupName    string                `json:"groupName"`
	Notification string                `json:"notification"`
	Introduction string                `json:"introduction"`
	FaceURL      string                `json:"faceURL"`
	Ex           string                `json:"ex"`
	OperationID  string                `json:"operationID" binding:"required"`
}

type CreateGroupResp struct {
	CommResp
	GroupInfo sdkws.GroupInfo
	Data      map[string]interface{} `json:"data"`
}

type GetGroupApplicationListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"` //作为管理员或群主收到的 进群申请
}
type GetGroupApplicationListResp struct {
	CommResp
	GroupRequestList []*sdkws.GroupRequest
	Data             []map[string]interface{} `json:"data"`
}

type GetUserReqGroupApplicationListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
}

type GetUserRespGroupApplicationResp struct {
	CommResp
	GroupRequestList []*sdkws.GroupRequest `json:"-"`
}

type GetGroupInfoReq struct {
	GroupIDList []string `json:"groupIDList" binding:"required"`
	OperationID string   `json:"operationID" binding:"required"`
}
type GetGroupInfoResp struct {
	CommResp
	GroupInfoList []*sdkws.GroupInfo       `json:"-"`
	Data          []map[string]interface{} `json:"data"`
}

//type GroupInfoAlias struct {
//	GroupID          string `protobuf:"bytes,1,opt,name=groupID" json:"groupID,omitempty"`
//	GroupName        string `protobuf:"bytes,2,opt,name=groupName" json:"groupName,omitempty"`
//	NotificationCmd     string `protobuf:"bytes,3,opt,name=notification" json:"notification,omitempty"`
//	Introduction     string `protobuf:"bytes,4,opt,name=introduction" json:"introduction,omitempty"`
//	FaceURL          string `protobuf:"bytes,5,opt,name=faceURL" json:"faceURL,omitempty"`
//	OwnerUserID      string `protobuf:"bytes,6,opt,name=ownerUserID" json:"ownerUserID,omitempty"`
//	CreateTime       uint32 `protobuf:"varint,7,opt,name=createTime" json:"createTime,omitempty"`
//	MemberCount      uint32 `protobuf:"varint,8,opt,name=memberCount" json:"memberCount,omitempty"`
//	Ex               string `protobuf:"bytes,9,opt,name=ex" json:"ex,omitempty"`
//	Status           int32  `protobuf:"varint,10,opt,name=status" json:"status,omitempty"`
//	CreatorUserID    string `protobuf:"bytes,11,opt,name=creatorUserID" json:"creatorUserID,omitempty"`
//	GroupType        int32  `protobuf:"varint,12,opt,name=groupType" json:"groupType,omitempty"`
//	NeedVerification int32  `protobuf:"bytes,13,opt,name=needVerification" json:"needVerification,omitempty"`
//}
//type GroupInfoAlias struct {
//	GroupInfo
//	NeedVerification int32 `protobuf:"bytes,13,opt,name=needVerification" json:"needVerification,omitempty"`
//}

type ApplicationGroupResponseReq struct {
	OperationID  string `json:"operationID" binding:"required"`
	GroupID      string `json:"groupID" binding:"required"`
	FromUserID   string `json:"fromUserID" binding:"required"` //application from FromUserID
	HandledMsg   string `json:"handledMsg"`
	HandleResult int32  `json:"handleResult" binding:"required,oneof=-1 1"`
}
type ApplicationGroupResponseResp struct {
	CommResp
}

type JoinGroupReq struct {
	GroupID       string `json:"groupID" binding:"required"`
	ReqMessage    string `json:"reqMessage"`
	OperationID   string `json:"operationID" binding:"required"`
	JoinSource    int32  `json:"joinSource"`
	InviterUserID string `json:"inviterUserID"`
}

type JoinGroupResp struct {
	CommResp
}

type QuitGroupReq struct {
	GroupID     string `json:"groupID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}
type QuitGroupResp struct {
	CommResp
}

type SetGroupInfoReq struct {
	GroupID           string `json:"groupID" binding:"required"`
	GroupName         string `json:"groupName"`
	Notification      string `json:"notification"`
	Introduction      string `json:"introduction"`
	FaceURL           string `json:"faceURL"`
	Ex                string `json:"ex"`
	OperationID       string `json:"operationID" binding:"required"`
	NeedVerification  *int32 `json:"needVerification" binding:"oneof=0 1 2"`
	LookMemberInfo    *int32 `json:"lookMemberInfo" binding:"oneof=0 1"`
	ApplyMemberFriend *int32 `json:"applyMemberFriend" binding:"oneof=0 1"`
}

type SetGroupInfoResp struct {
	CommResp
}

type TransferGroupOwnerReq struct {
	GroupID        string `json:"groupID" binding:"required"`
	OldOwnerUserID string `json:"oldOwnerUserID" binding:"required"`
	NewOwnerUserID string `json:"newOwnerUserID" binding:"required"`
	OperationID    string `json:"operationID" binding:"required"`
}
type TransferGroupOwnerResp struct {
	CommResp
}

type DismissGroupReq struct {
	GroupID     string `json:"groupID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}
type DismissGroupResp struct {
	CommResp
}

type MuteGroupMemberReq struct {
	OperationID  string `json:"operationID" binding:"required"`
	GroupID      string `json:"groupID" binding:"required"`
	UserID       string `json:"userID" binding:"required"`
	MutedSeconds uint32 `json:"mutedSeconds" binding:"required"`
}
type MuteGroupMemberResp struct {
	CommResp
}

type CancelMuteGroupMemberReq struct {
	OperationID string `json:"operationID" binding:"required"`
	GroupID     string `json:"groupID" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
}
type CancelMuteGroupMemberResp struct {
	CommResp
}

type MuteGroupReq struct {
	OperationID string `json:"operationID" binding:"required"`
	GroupID     string `json:"groupID" binding:"required"`
}
type MuteGroupResp struct {
	CommResp
}

type CancelMuteGroupReq struct {
	OperationID string `json:"operationID" binding:"required"`
	GroupID     string `json:"groupID" binding:"required"`
}
type CancelMuteGroupResp struct {
	CommResp
}

type SetGroupMemberNicknameReq struct {
	OperationID string `json:"operationID" binding:"required"`
	GroupID     string `json:"groupID" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
	Nickname    string `json:"nickname"`
}

type SetGroupMemberNicknameResp struct {
	CommResp
}

type SetGroupMemberBaseInfoReq struct {
	OperationID string `json:"operationID" binding:"required"`
	GroupID     string `json:"groupID" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
}

type SetGroupMemberInfoReq struct {
	OperationID string  `json:"operationID" binding:"required"`
	GroupID     string  `json:"groupID" binding:"required"`
	UserID      string  `json:"userID" binding:"required"`
	Nickname    *string `json:"nickname"`
	FaceURL     *string `json:"userGroupFaceUrl"`
	RoleLevel   *int32  `json:"roleLevel" validate:"gte=1,lte=3"`
	Ex          *string `json:"ex"`
}

type SetGroupMemberRoleLevelReq struct {
	SetGroupMemberBaseInfoReq
	RoleLevel int `json:"roleLevel"`
}

type SetGroupMemberRoleLevelResp struct {
	CommResp
}

type GetGroupAbstractInfoReq struct {
	OperationID string `json:"operationID" binding:"required"`
	GroupID     string `json:"groupID" binding:"required"`
}

type GetGroupAbstractInfoResp struct {
	CommResp
	GroupMemberNumber   int32  `json:"groupMemberNumber"`
	GroupMemberListHash uint64 `json:"groupMemberListHash"`
}
