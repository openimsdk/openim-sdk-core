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

package sdk_params_callback

import (
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/server_api_params"
)

type CreateGroupBaseInfoParam struct {
	GroupType int32 `json:"groupType"`
	SetGroupInfoParam
}

type CreateGroupMemberRoleParam []*server_api_params.GroupAddMemberInfo
type CreateGroupCallback map[string]interface{}

// param groupID reqMsg
const JoinGroupCallback = constant.SuccessCallbackDefault

// type QuitGroupParam // groupID
const QuitGroupCallback = constant.SuccessCallbackDefault

const DismissGroupCallback = constant.SuccessCallbackDefault

const GroupMuteChangeCallback = constant.SuccessCallbackDefault

const GroupMemberMuteChangeCallback = constant.SuccessCallbackDefault

const SetGroupMemberNicknameCallback = constant.SuccessCallbackDefault

// type GetJoinedGroupListParam null
type GetJoinedGroupListCallback []*model_struct.LocalGroup

type GetGroupsInfoParam []string
type GetGroupsInfoCallback []*model_struct.LocalGroup
type SearchGroupsParam struct {
	KeywordList       []string `json:"keywordList"`
	IsSearchGroupID   bool     `json:"isSearchGroupID"`
	IsSearchGroupName bool     `json:"isSearchGroupName"`
}
type SearchGroupsCallback []*model_struct.LocalGroup

type SearchGroupMembersParam struct {
	GroupID                string   `json:"groupID"`
	KeywordList            []string `json:"keywordList"`
	IsSearchUserID         bool     `json:"isSearchUserID"`
	IsSearchMemberNickname bool     `json:"isSearchMemberNickname"`
	//offset, count int
	Offset int `json:"offset"`
	Count  int `json:"count"`
}
type SearchGroupMembersCallback []*model_struct.LocalGroupMember

type SetGroupInfoParam struct {
	GroupName        string `json:"groupName"`
	Notification     string `json:"notification"`
	Introduction     string `json:"introduction"`
	FaceURL          string `json:"faceURL"`
	Ex               string `json:"ex"`
	NeedVerification *int32 `json:"needVerification" binding:"oneof=0 1 2"`
}

type SetGroupMemberInfoParam struct {
	GroupID string  `json:"groupID"`
	UserID  string  `json:"userID"`
	Ex      *string `json:"ex"`
}

const SetGroupMemberInfoCallback = constant.SuccessCallbackDefault

const SetGroupInfoCallback = constant.SuccessCallbackDefault

// type GetGroupMemberListParam groupID ...
type GetGroupMemberListCallback []*model_struct.LocalGroupMember

type GetGroupMembersInfoParam []string
type GetGroupMembersInfoCallback []*model_struct.LocalGroupMember

type KickGroupMemberParam []string
type KickGroupMemberCallback []*server_api_params.UserIDResult

// type TransferGroupOwnerParam
const TransferGroupOwnerCallback = constant.SuccessCallbackDefault

type InviteUserToGroupParam []string
type InviteUserToGroupCallback []*server_api_params.UserIDResult

// type GetGroupApplicationListParam
type GetGroupApplicationListCallback []*model_struct.LocalAdminGroupRequest

type GetSendGroupApplicationListCallback []*model_struct.LocalGroupRequest

// type AcceptGroupApplicationParam
const AcceptGroupApplicationCallback = constant.SuccessCallbackDefault

// type RefuseGroupApplicationParam
const RefuseGroupApplicationCallback = constant.SuccessCallbackDefault
