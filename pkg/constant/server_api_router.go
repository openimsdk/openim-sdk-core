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
	SetGlobalRecvMessageOptRouter = "/user/set_global_msg_recv_opt"
	ProcessUserCommandGetAll      = "/user/process_user_command_get_all"

	UserRegister = "/user/user_register"

	AddFriendRouter        = "/friend/add_friend"
	DeleteFriendRouter     = "/friend/delete_friend"
	ImportFriendListRouter = "/friend/import_friend"

	GetFriendListRouter = "/friend/get_friend_list"
	AddFriendResponse   = "/friend/add_friend_response"
	SetFriendRemark     = "/friend/set_friend_remark"
	UpdateFriends       = "/friend/update_friends"

	AddBlackRouter    = "/friend/add_black"
	RemoveBlackRouter = "/friend/remove_black"

	PullUserMsgBySeqRouter = "/chat/pull_msg_by_seq"

	// msg
	ClearConversationMsgRouter             = RouterMsg + "/clear_conversation_msg" // Clear the message of the specified conversation
	ClearAllMsgRouter                      = RouterMsg + "/user_clear_all_msg"     // Clear all messages of the current user
	DeleteMsgsRouter                       = RouterMsg + "/delete_msgs"            // Delete the specified message
	RevokeMsgRouter                        = RouterMsg + "/revoke_msg"
	MarkMsgsAsReadRouter                   = RouterMsg + "/mark_msgs_as_read"
	GetConversationsHasReadAndMaxSeqRouter = RouterMsg + "/get_conversations_has_read_and_max_seq"

	MarkConversationAsRead    = RouterMsg + "/mark_conversation_as_read"
	SetConversationHasReadSeq = RouterMsg + "/set_conversation_has_read_seq"

	GetConversationsRouter     = ConversationGroup + "/get_conversations"
	GetAllConversationsRouter  = ConversationGroup + "/get_all_conversations"
	SetConversationsRouter     = ConversationGroup + "/set_conversations"
	GetIncrementalConversation = ConversationGroup + "/get_incremental_conversations"
	GetFullConversationIDs     = ConversationGroup + "/get_full_conversation_ids"
	GetOwnerConversationRouter = ConversationGroup + "/get_owner_conversation"

	// auth
	GetUsersToken = RouterAuth + "/user_token"
)
const (
	ConversationGroup = "/conversation"
	RouterAuth        = "/auth"
	RouterMsg         = "/msg"
)
