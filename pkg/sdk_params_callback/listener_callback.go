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
	"open_im_sdk/pkg/db/model_struct"
)

// //////////////////////////////friend////////////////////////////////////
type FriendApplicationAddedCallback model_struct.LocalFriendRequest
type FriendApplicationAcceptCallback model_struct.LocalFriendRequest
type FriendApplicationRejectCallback model_struct.LocalFriendRequest
type FriendApplicationDeletedCallback model_struct.LocalFriendRequest
type FriendAddedCallback model_struct.LocalFriend
type FriendDeletedCallback model_struct.LocalFriend
type FriendInfoChangedCallback model_struct.LocalFriend
type BlackAddCallback model_struct.LocalBlack
type BlackDeletedCallback model_struct.LocalBlack

////////////////////////////////group////////////////////////////////////

type JoinedGroupAddedCallback model_struct.LocalGroup
type JoinedGroupDeletedCallback model_struct.LocalGroup
type GroupMemberAddedCallback model_struct.LocalGroupMember
type GroupMemberDeletedCallback model_struct.LocalGroupMember
type GroupApplicationAddedCallback model_struct.LocalAdminGroupRequest
type GroupApplicationDeletedCallback model_struct.LocalAdminGroupRequest
type GroupApplicationAcceptCallback model_struct.LocalAdminGroupRequest
type GroupApplicationRejectCallback model_struct.LocalAdminGroupRequest
type GroupInfoChangedCallback model_struct.LocalGroup
type GroupMemberInfoChangedCallback model_struct.LocalGroupMember

// ////////////////////////////user////////////////////////////////////////
type SelfInfoUpdatedCallback model_struct.LocalUser

// ////////////////////////////user////////////////////////////////////////
type ConversationUpdateCallback model_struct.LocalConversation
type ConversationDeleteCallback model_struct.LocalConversation

// ///////////////////////////signaling/////////////////////////////////////
type InvitationInfo struct {
	InviterUserID     string
	InviteeUserIDList []string
	CustomData        string
	GroupID           string
}

//type ReceiveNewInvitationCallback sdkws.SignalInviteReq
//
//type InviteeAcceptedCallback sdkws.SignalAcceptReq
//
//type InviteeRejectedCallback sdkws.SignalRejectReq
//
//type InvitationCancelledCallback sdkws.SignalCancelReq
//
//type InvitationTimeoutCallback sdkws.SignalInviteReq
