package sdk_params_callback

import (
	"open_im_sdk/pkg/db/model_struct"
	api "open_im_sdk/pkg/server_api_params"
)

////////////////////////////////friend////////////////////////////////////
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

//////////////////////////////user////////////////////////////////////////
type SelfInfoUpdatedCallback model_struct.LocalUser

//////////////////////////////user////////////////////////////////////////
type ConversationUpdateCallback model_struct.LocalConversation
type ConversationDeleteCallback model_struct.LocalConversation

/////////////////////////////signaling/////////////////////////////////////
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
