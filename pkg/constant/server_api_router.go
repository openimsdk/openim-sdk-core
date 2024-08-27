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

package constant

const (
	GetUsersInfoRouter            = "/user/get_users_info"
	UpdateSelfUserInfoRouter      = "/user/update_user_info"
	UpdateSelfUserInfoExRouter    = "/user/update_user_info_ex"
	SetGlobalRecvMessageOptRouter = "/user/set_global_msg_recv_opt"
	ProcessUserCommandAdd         = "/user/process_user_command_add"
	ProcessUserCommandDelete      = "/user/process_user_command_delete"
	ProcessUserCommandUpdate      = "/user/process_user_command_update"
	ProcessUserCommandGet         = "/user/process_user_command_get"
	ProcessUserCommandGetAll      = "/user/process_user_command_get_all"

	GetUsersInfoFromCacheRouter   = "/user/get_users_info_from_cache"
	AccountCheck                  = "/user/account_check"
	UserRegister                  = "/user/user_register"
	SubscribeUsersStatusRouter    = "/user/subscribe_users_status"
	GetSubscribeUsersStatusRouter = "/user/get_subscribe_users_status"
	GetUserStatusRouter           = "/user/get_users_status"

	AddFriendRouter                    = "/friend/add_friend"
	DeleteFriendRouter                 = "/friend/delete_friend"
	GetFriendApplicationListRouter     = "/friend/get_friend_apply_list"      //recv
	GetSelfFriendApplicationListRouter = "/friend/get_self_friend_apply_list" //send
	ImportFriendListRouter             = "/friend/import_friend"

	GetDesignatedFriendsApplyRouter = "/friend/get_designated_friend_apply"
	GetFriendListRouter             = "/friend/get_friend_list"
	GetDesignatedFriendsRouter      = "/friend/get_designated_friends"
	AddFriendResponse               = "/friend/add_friend_response"
	SetFriendRemark                 = "/friend/set_friend_remark"
	UpdateFriends                   = "/friend/update_friends"
	GetIncrementalFriends           = "/friend/get_incremental_friends"
	GetFullFriendUserIDs            = "/friend/get_full_friend_user_ids"

	AddBlackRouter     = "/friend/add_black"
	RemoveBlackRouter  = "/friend/remove_black"
	GetBlackListRouter = "/friend/get_black_list"

	PullUserMsgRouter      = "/chat/pull_msg"
	PullUserMsgBySeqRouter = "/chat/pull_msg_by_seq"
	NewestSeqRouter        = "/chat/newest_seq"

	// msg
	ClearConversationMsgRouter             = RouterMsg + "/clear_conversation_msg" // Clear the message of the specified conversation
	ClearAllMsgRouter                      = RouterMsg + "/user_clear_all_msg"     // Clear all messages of the current user
	DeleteMsgsRouter                       = RouterMsg + "/delete_msgs"            // Delete the specified message
	RevokeMsgRouter                        = RouterMsg + "/revoke_msg"
	MarkMsgsAsReadRouter                   = RouterMsg + "/mark_msgs_as_read"
	GetConversationsHasReadAndMaxSeqRouter = RouterMsg + "/get_conversations_has_read_and_max_seq"

	MarkConversationAsRead    = RouterMsg + "/mark_conversation_as_read"
	MarkMsgsAsRead            = RouterMsg + "/mark_msgs_as_read"
	SetConversationHasReadSeq = RouterMsg + "/set_conversation_has_read_seq"
	SendMsgRouter             = RouterMsg + "/send_msg"
	GetServerTimeRouter       = RouterMsg + "/get_server_time"

	GetConversationsRouter        = ConversationGroup + "/get_conversations"
	GetAllConversationsRouter     = ConversationGroup + "/get_all_conversations"
	GetConversationRouter         = ConversationGroup + "/get_conversation"
	BatchSetConversationRouter    = ConversationGroup + "/batch_set_conversation"
	ModifyConversationFieldRouter = ConversationGroup + "/modify_conversation_field"
	SetConversationsRouter        = ConversationGroup + "/set_conversations"
	GetIncrementalConversation    = ConversationGroup + "/get_incremental_conversations"
	GetFullConversationIDs        = ConversationGroup + "/get_full_conversation_ids"
	GetOwnerConversationRouter    = ConversationGroup + "/get_owner_conversation"

	ParseTokenRouter = RouterAuth + "/parse_token"

	// auth
	GetUsersToken = RouterAuth + "/user_token"
)
const (
	RouterGroup        = "/group"
	ConversationGroup  = "/conversation"
	RouterOrganization = "/organization"
	RouterAuth         = "/auth"
	RouterSuperGroup   = "/super_group"
	RouterMsg          = "/msg"
	RouterThird        = "/third"
)
