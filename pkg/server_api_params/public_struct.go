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

import (
	"open_im_sdk/pkg/db/model_struct"
)

type ApiUserInfo struct {
	UserID           string `json:"userID" binding:"required,min=1,max=64"`
	Nickname         string `json:"nickname" binding:"omitempty,min=1,max=64"`
	FaceURL          string `json:"faceURL" binding:"omitempty,max=1024"`
	GlobalRecvMsgOpt int32  `json:"globalRecvMsgOpt" binding:"omitempty,oneof=0 1 2"`
	Ex               string `json:"ex" binding:"omitempty,max=1024"`
}

type GroupAddMemberInfo struct {
	UserID    string `json:"userID" validate:"required"`
	RoleLevel int32  `json:"roleLevel" validate:"required"`
}

type PublicUser struct {
	UserID     string `json:"userID"`
	Nickname   string `json:"nickname"`
	FaceURL    string `json:"faceURL"`
	Ex         string `json:"ex"`
	CreateTime int64  `json:"createTime"`
}

type FullUserInfo struct {
	PublicInfo *PublicUser               `json:"publicInfo"`
	FriendInfo *model_struct.LocalFriend `json:"friendInfo"`
	BlackInfo  *model_struct.LocalBlack  `json:"blackInfo"`
}

//GroupName    string                `json:"groupName"`
//	Introduction string                `json:"introduction"`
//	NotificationCmd string                `json:"notification"`
//	FaceUrl      string                `json:"faceUrl"`
//	OperationID  string                `json:"operationID" binding:"required"`
//	GroupType    int32                 `json:"groupType"`
//	Ex           string                `json:"ex"`

//type GroupInfo struct {
//	GroupID       string `json:"groupID"`
//	GroupName     string `json:"groupName"`
//	NotificationCmd  string `json:"notification"`
//	Introduction  string `json:"introduction"`
//	FaceUrl       string `json:"faceUrl"`
//	OwnerUserID   string `json:"ownerUserID"`
//	Ex            string `json:"ex"`
//	GroupType     int32  `json:"groupType"`
//}

//type GroupMemberFullInfo struct {
//	GroupID        string `json:"groupID"`
//	UserID         string `json:"userID"`
//	RoleLevel      int32  `json:"roleLevel"`
//	JoinTime       uint64 `json:"joinTime"`
//	Nickname       string `json:"nickname"`
//	FaceUrl        string `json:"faceUrl"`
//	FriendRemark   string `json:"friendRemark"`
//	AppMangerLevel int32  `json:"appMangerLevel"`
//	JoinSource     int32  `json:"joinSource"`
//	OperatorUserID string `json:"operatorUserID"`
//	Ex             string `json:"ex"`
//}
//
//type PublicUserInfo struct {
//	UserID   string `json:"userID"`
//	Nickname string `json:"nickname"`
//	FaceUrl  string `json:"faceUrl"`
//	Gender   int32  `json:"gender"`
//}
//
//type UserInfo struct {
//	UserID   string `json:"userID"`
//	Nickname string `json:"nickname"`
//	FaceUrl  string `json:"faceUrl"`
//	Gender   int32  `json:"gender"`
//	Mobile   string `json:"mobile"`
//	Birth    string `json:"birth"`
//	Email    string `json:"email"`
//	Ex       string `json:"ex"`
//}
//
//type FriendInfo struct {
//	OwnerUserID    string   `json:"ownerUserID"`
//	Remark         string   `json:"remark"`
//	CreateTime     int64    `json:"createTime"`
//	FriendUser     UserInfo `json:"friendUser"`
//	AddSource      int32    `json:"addSource"`
//	OperatorUserID string   `json:"operatorUserID"`
//	Ex             string   `json:"ex"`
//}
//
//type BlackInfo struct {
//	OwnerUserID    string         `json:"ownerUserID"`
//	CreateTime     int64          `json:"createTime"`
//	BlackUser      PublicUserInfo `json:"friendUser"`
//	AddSource      int32          `json:"addSource"`
//	OperatorUserID string         `json:"operatorUserID"`
//	Ex             string         `json:"ex"`
//}
//
//type GroupRequest struct {
//	UserID       string `json:"userID"`
//	GroupID      string `json:"groupID"`
//	HandleResult string `json:"handleResult"`
//	ReqMsg       string `json:"reqMsg"`
//	HandleMsg    string `json:"handleMsg"`
//	ReqTime      int64  `json:"reqTime"`
//	HandleUserID string `json:"handleUserID"`
//	HandleTime   int64  `json:"handleTime"`
//	Ex           string `json:"ex"`
//}
//
//type FriendRequest struct {
//	FromUserID    string `json:"fromUserID"`
//	ToUserID      string `json:"toUserID"`
//	HandleResult  int32  `json:"handleResult"`
//	ReqMessage    string `json:"reqMessage"`
//	CreateTime    int64  `json:"createTime"`
//	HandlerUserID string `json:"handlerUserID"`
//	HandleMsg     string `json:"handleMsg"`
//	HandleTime    int64  `json:"handleTime"`
//	Ex            string `json:"ex"`
//}
//
//
//
