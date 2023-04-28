// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package sdk_params_callback

import (
	"open_im_sdk/pkg/db/model_struct"
	api "open_im_sdk/pkg/server_api_params"
)

////////////////////////////////friend////////////////////////////////////.
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

type (
	JoinedGroupAddedCallback        model_struct.LocalGroup
	JoinedGroupDeletedCallback      model_struct.LocalGroup
	GroupMemberAddedCallback        model_struct.LocalGroupMember
	GroupMemberDeletedCallback      model_struct.LocalGroupMember
	GroupApplicationAddedCallback   model_struct.LocalAdminGroupRequest
	GroupApplicationDeletedCallback model_struct.LocalAdminGroupRequest
	GroupApplicationAcceptCallback  model_struct.LocalAdminGroupRequest
	GroupApplicationRejectCallback  model_struct.LocalAdminGroupRequest
	GroupInfoChangedCallback        model_struct.LocalGroup
	GroupMemberInfoChangedCallback  model_struct.LocalGroupMember
)

//////////////////////////////user////////////////////////////////////////.
type SelfInfoUpdatedCallback model_struct.LocalUser

//////////////////////////////user////////////////////////////////////////.
type ConversationUpdateCallback model_struct.LocalConversation
type ConversationDeleteCallback model_struct.LocalConversation

/////////////////////////////signaling/////////////////////////////////////.
type InvitationInfo struct {
	InviterUserID     string
	InviteeUserIDList []string
	CustomData        string
	GroupID           string
}

type ReceiveNewInvitationCallback api.SignalInviteReq

type InviteeAcceptedCallback api.SignalAcceptReq

type InviteeRejectedCallback api.SignalRejectReq

type InvitationCancelledCallback api.SignalCancelReq

type InvitationTimeoutCallback api.SignalInviteReq
