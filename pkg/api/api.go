package api

import (
	"github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/third"
	"github.com/openimsdk/protocol/user"
)

var (
	ParseToken = newApi[auth.ParseTokenReq, auth.ParseTokenResp]("/auth/parse_token")
)

var (
	GetUsersInfo             = newApi[user.GetDesignateUsersReq, user.GetDesignateUsersResp]("/user/get_users_info")
	UpdateUserInfo           = newApi[user.UpdateUserInfoReq, user.UpdateUserInfoResp]("/user/update_user_info")
	UpdateUserInfoEx         = newApi[user.UpdateUserInfoExReq, user.UpdateUserInfoExResp]("/user/update_user_info_ex")
	ProcessUserCommandAdd    = newApi[user.ProcessUserCommandAddReq, user.ProcessUserCommandAddResp]("/user/process_user_command_add")
	ProcessUserCommandDelete = newApi[user.ProcessUserCommandDeleteReq, user.ProcessUserCommandDeleteResp]("/user/process_user_command_delete")
	ProcessUserCommandUpdate = newApi[user.ProcessUserCommandUpdateReq, user.ProcessUserCommandUpdateResp]("/user/process_user_command_update")
	ProcessUserCommandGet    = newApi[user.ProcessUserCommandGetReq, user.ProcessUserCommandGetResp]("/user/process_user_command_get")
	ProcessUserCommandGetAll = newApi[user.ProcessUserCommandGetAllReq, user.ProcessUserCommandGetAllResp]("/user/process_user_command_get_all")
	UserRegister             = newApi[user.UserRegisterReq, user.UserRegisterResp]("/user/user_register")
)

var (
	AddFriend                    = newApi[relation.ApplyToAddFriendReq, relation.ApplyToAddFriendResp]("/friend/add_friend")
	DeleteFriend                 = newApi[relation.DeleteFriendReq, relation.DeleteFriendResp]("/friend/delete_friend")
	GetRecvFriendApplicationList = newApi[relation.GetPaginationFriendsApplyToReq, relation.GetPaginationFriendsApplyToResp]("/friend/get_friend_apply_list")
	GetSelfFriendApplicationList = newApi[relation.GetPaginationFriendsApplyFromReq, relation.GetPaginationFriendsApplyFromResp]("/friend/get_self_friend_apply_list")
	GetSelfUnhandledApplyCount   = newApi[relation.GetSelfUnhandledApplyCountReq, relation.GetSelfUnhandledApplyCountResp]("/friend/get_self_unhandled_apply_count")
	ImportFriendList             = newApi[relation.ImportFriendReq, relation.ImportFriendResp]("/friend/import_friend")
	GetDesignatedFriendsApply    = newApi[relation.GetDesignatedFriendsApplyReq, relation.GetDesignatedFriendsApplyResp]("/friend/get_designated_friend_apply")
	GetFriendList                = newApi[relation.GetPaginationFriendsReq, relation.GetPaginationFriendsResp]("/friend/get_friend_list")
	GetDesignatedFriends         = newApi[relation.GetDesignatedFriendsReq, relation.GetDesignatedFriendsResp]("/friend/get_designated_friends")
	AddFriendResponse            = newApi[relation.RespondFriendApplyReq, relation.RespondFriendApplyResp]("/friend/add_friend_response")
	SetFriendRemark              = newApi[relation.SetFriendRemarkReq, relation.SetFriendRemarkResp]("/friend/set_friend_remark")
	UpdateFriends                = newApi[relation.UpdateFriendsReq, relation.UpdateFriendsResp]("/friend/update_friends")
	GetIncrementalFriends        = newApi[relation.GetIncrementalFriendsReq, relation.GetIncrementalFriendsResp]("/friend/get_incremental_friends")
	GetFullFriendUserIDs         = newApi[relation.GetFullFriendUserIDsReq, relation.GetFullFriendUserIDsResp]("/friend/get_full_friend_user_ids")
	AddBlack                     = newApi[relation.AddBlackReq, relation.AddBlackResp]("/friend/add_black")
	RemoveBlack                  = newApi[relation.RemoveBlackReq, relation.RemoveBlackResp]("/friend/remove_black")
	GetBlackList                 = newApi[relation.GetPaginationBlacksReq, relation.GetPaginationBlacksResp]("/friend/get_black_list")
)

var (
	ClearConversationMsg             = newApi[msg.ClearConversationsMsgReq, msg.ClearConversationsMsgResp]("/msg/clear_conversation_msg") // Clear the message of the specified conversation
	ClearAllMsg                      = newApi[msg.UserClearAllMsgReq, msg.UserClearAllMsgResp]("/msg/user_clear_all_msg")                 // Clear all messages of the current user
	DeleteMsgs                       = newApi[msg.DeleteMsgsReq, msg.DeleteMsgsResp]("/msg/delete_msgs")                                  // Delete the specified message
	RevokeMsg                        = newApi[msg.RevokeMsgReq, msg.RevokeMsgResp]("/msg/revoke_msg")
	MarkMsgsAsRead                   = newApi[msg.MarkMsgsAsReadReq, msg.MarkMsgsAsReadResp]("/msg/mark_msgs_as_read")
	GetConversationsHasReadAndMaxSeq = newApi[msg.GetConversationsHasReadAndMaxSeqReq, msg.GetConversationsHasReadAndMaxSeqResp]("/msg/get_conversations_has_read_and_max_seq")
	MarkConversationAsRead           = newApi[msg.MarkConversationAsReadReq, msg.MarkConversationAsReadResp]("/msg/mark_conversation_as_read")
	SetConversationHasReadSeq        = newApi[msg.SetConversationHasReadSeqReq, msg.SetConversationHasReadSeqResp]("/msg/set_conversation_has_read_seq")
	SendMsg                          = newApi[msg.SendMsgReq, msg.SendMsgResp]("/msg/send_msg")
	GetServerTime                    = newApi[msg.GetServerTimeReq, msg.GetServerTimeResp]("/msg/get_server_time")
)

var (
	CreateGroup                       = newApi[group.CreateGroupReq, group.CreateGroupResp]("/group/create_group")
	SetGroupInfoEx                    = newApi[group.SetGroupInfoExReq, group.SetGroupInfoExResp]("/group/set_group_info_ex")
	JoinGroup                         = newApi[group.JoinGroupReq, group.JoinGroupResp]("/group/join_group")
	QuitGroup                         = newApi[group.QuitGroupReq, group.QuitGroupResp]("/group/quit_group")
	GetGroupsInfo                     = newApi[group.GetGroupsInfoReq, group.GetGroupsInfoResp]("/group/get_groups_info")
	GetGroupMemberList                = newApi[group.GetGroupMemberListReq, group.GetGroupMemberListResp]("/group/get_group_member_list")
	GetGroupMembersInfo               = newApi[group.GetGroupMembersInfoReq, group.GetGroupMembersInfoResp]("/group/get_group_members_info")
	InviteUserToGroup                 = newApi[group.InviteUserToGroupReq, group.InviteUserToGroupResp]("/group/invite_user_to_group")
	GetJoinedGroupList                = newApi[group.GetJoinedGroupListReq, group.GetJoinedGroupListResp]("/group/get_joined_group_list")
	KickGroupMember                   = newApi[group.KickGroupMemberReq, group.KickGroupMemberResp]("/group/kick_group")
	TransferGroup                     = newApi[group.TransferGroupOwnerReq, group.TransferGroupOwnerResp]("/group/transfer_group")
	GetRecvGroupApplicationList       = newApi[group.GetGroupApplicationListReq, group.GetGroupApplicationListResp]("/group/get_recv_group_applicationList")
	GetSendGroupApplicationList       = newApi[group.GetUserReqApplicationListReq, group.GetUserReqApplicationListResp]("/group/get_user_req_group_applicationList")
	GetGroupApplicationUnhandledCount = newApi[group.GetGroupApplicationUnhandledCountReq, group.GetGroupApplicationUnhandledCountResp]("/group/get_group_application_unhandled_count")
	AcceptGroupApplication            = newApi[group.GroupApplicationResponseReq, group.GroupApplicationResponseResp]("/group/group_application_response")
	DismissGroup                      = newApi[group.DismissGroupReq, group.DismissGroupResp]("/group/dismiss_group")
	MuteGroupMember                   = newApi[group.MuteGroupMemberReq, group.MuteGroupMemberResp]("/group/mute_group_member")
	CancelMuteGroupMember             = newApi[group.CancelMuteGroupMemberReq, group.CancelMuteGroupMemberResp]("/group/cancel_mute_group_member")
	MuteGroup                         = newApi[group.MuteGroupReq, group.MuteGroupResp]("/group/mute_group")
	CancelMuteGroup                   = newApi[group.CancelMuteGroupReq, group.CancelMuteGroupResp]("/group/cancel_mute_group")
	SetGroupMemberInfo                = newApi[group.SetGroupMemberInfoReq, group.SetGroupMemberInfoResp]("/group/set_group_member_info")
	GetIncrementalJoinGroup           = newApi[group.GetIncrementalJoinGroupReq, group.GetIncrementalJoinGroupResp]("/group/get_incremental_join_groups")
	GetIncrementalGroupMemberBatch    = newApi[group.BatchGetIncrementalGroupMemberReq, group.BatchGetIncrementalGroupMemberResp]("/group/get_incremental_group_members_batch")
	GetFullJoinedGroupIDs             = newApi[group.GetFullJoinGroupIDsReq, group.GetFullJoinGroupIDsResp]("/group/get_full_join_group_ids")
	GetFullGroupMemberUserIDs         = newApi[group.GetFullGroupMemberUserIDsReq, group.GetFullGroupMemberUserIDsResp]("/group/get_full_group_member_user_ids")
)

var (
	GetConversations           = newApi[conversation.GetConversationsReq, conversation.GetConversationsResp]("/conversation/get_conversations")
	GetAllConversations        = newApi[conversation.GetAllConversationsReq, conversation.GetAllConversationsResp]("/conversation/get_all_conversations")
	SetConversations           = newApi[conversation.SetConversationsReq, conversation.SetConversationsResp]("/conversation/set_conversations")
	GetIncrementalConversation = newApi[conversation.GetIncrementalConversationReq, conversation.GetIncrementalConversationResp]("/conversation/get_incremental_conversations")
	GetFullConversationIDs     = newApi[conversation.GetFullOwnerConversationIDsReq, conversation.GetFullOwnerConversationIDsResp]("/conversation/get_full_conversation_ids")
	GetOwnerConversation       = newApi[conversation.GetOwnerConversationReq, conversation.GetOwnerConversationResp]("/conversation/get_owner_conversation")
)

var (
	GetAdminToken = newApi[auth.GetAdminTokenReq, auth.GetAdminTokenResp]("/auth/get_admin_token")
	GetUsersToken = newApi[auth.GetUserTokenReq, auth.GetUserTokenResp]("/auth/get_user_token")
)

var (
	FcmUpdateToken = newApi[third.FcmUpdateTokenReq, third.FcmUpdateTokenResp]("/third/fcm_update_token")
	SetAppBadge    = newApi[third.SetAppBadgeReq, third.SetAppBadgeResp]("/third/set_app_badge")
	UploadLogs     = newApi[third.UploadLogsReq, third.UploadLogsResp]("/third/logs/upload")
)

var (
	ObjectPartLimit               = newApi[third.PartLimitReq, third.PartLimitResp]("/object/part_limit")
	ObjectInitiateMultipartUpload = newApi[third.InitiateMultipartUploadReq, third.InitiateMultipartUploadResp]("/object/initiate_multipart_upload")
	ObjectAuthSign                = newApi[third.AuthSignReq, third.AuthSignResp]("/object/auth_sign")
	ObjectCompleteMultipartUpload = newApi[third.CompleteMultipartUploadReq, third.CompleteMultipartUploadResp]("/object/complete_multipart_upload")
	ObjectAccessURL               = newApi[third.AccessURLReq, third.AccessURLResp]("/object/access_url")
)
