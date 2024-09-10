package api

import (
	"github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/third"
	"github.com/openimsdk/protocol/user"
)

var (
	ParseToken = api[auth.ParseTokenReq, auth.ParseTokenResp]("/auth/parse_token")
)

var (
	GetUsersInfo             = api[user.GetDesignateUsersReq, user.GetDesignateUsersResp]("/user/get_users_info")
	UpdateUserInfo           = api[user.UpdateUserInfoReq, user.UpdateUserInfoResp]("/user/update_user_info")
	UpdateUserInfoEx         = api[user.UpdateUserInfoExReq, user.UpdateUserInfoExResp]("/user/update_user_info_ex")
	SetGlobalRecvMessageOpt  = api[user.SetGlobalRecvMessageOptReq, user.SetGlobalRecvMessageOptResp]("/user/set_global_msg_recv_opt")
	ProcessUserCommandAdd    = api[user.ProcessUserCommandAddReq, user.ProcessUserCommandAddResp]("/user/process_user_command_add")
	ProcessUserCommandDelete = api[user.ProcessUserCommandDeleteReq, user.ProcessUserCommandDeleteResp]("/user/process_user_command_delete")
	ProcessUserCommandUpdate = api[user.ProcessUserCommandUpdateReq, user.ProcessUserCommandUpdateResp]("/user/process_user_command_update")
	ProcessUserCommandGet    = api[user.ProcessUserCommandGetReq, user.ProcessUserCommandGetResp]("/user/process_user_command_get")
	ProcessUserCommandGetAll = api[user.ProcessUserCommandGetAllReq, user.ProcessUserCommandGetAllResp]("/user/process_user_command_get_all")
	UserRegister             = api[user.UserRegisterReq, user.UserRegisterResp]("/user/user_register")
)

var (
	AddFriend                    = api[relation.ApplyToAddFriendReq, relation.ApplyToAddFriendResp]("/friend/add_friend")
	DeleteFriend                 = api[relation.DeleteFriendReq, relation.DeleteFriendResp]("/friend/delete_friend")
	GetFriendApplicationList     = api[relation.GetPaginationFriendsApplyToReq, relation.GetPaginationFriendsApplyToResp]("/friend/get_friend_apply_list")
	GetSelfFriendApplicationList = api[relation.GetPaginationFriendsApplyFromReq, relation.GetPaginationFriendsApplyFromResp]("/friend/get_self_friend_apply_list")
	ImportFriendList             = api[relation.ImportFriendReq, relation.ImportFriendResp]("/friend/import_friend")
	GetDesignatedFriendsApply    = api[relation.GetDesignatedFriendsApplyReq, relation.GetDesignatedFriendsApplyResp]("/friend/get_designated_friend_apply")
	GetFriendList                = api[relation.GetPaginationFriendsReq, relation.GetPaginationFriendsResp]("/friend/get_friend_list")
	GetDesignatedFriends         = api[relation.GetDesignatedFriendsReq, relation.GetDesignatedFriendsResp]("/friend/get_designated_friends")
	AddFriendResponse            = api[relation.RespondFriendApplyReq, relation.RespondFriendApplyResp]("/friend/add_friend_response")
	SetFriendRemark              = api[relation.SetFriendRemarkReq, relation.SetFriendRemarkResp]("/friend/set_friend_remark")
	UpdateFriends                = api[relation.UpdateFriendsReq, relation.UpdateFriendsResp]("/friend/update_friends")
	GetIncrementalFriends        = api[relation.GetIncrementalFriendsReq, relation.GetIncrementalFriendsResp]("/friend/get_incremental_friends")
	GetFullFriendUserIDs         = api[relation.GetFullFriendUserIDsReq, relation.GetFullFriendUserIDsResp]("/friend/get_full_friend_user_ids")
	AddBlack                     = api[relation.AddBlackReq, relation.AddBlackResp]("/friend/add_black")
	RemoveBlack                  = api[relation.RemoveBlackReq, relation.RemoveBlackResp]("/friend/remove_black")
	GetBlackList                 = api[relation.GetPaginationBlacksReq, relation.GetPaginationBlacksResp]("/friend/get_black_list")
)

var (
	PullUserMsgBySeq = api[sdkws.PullMessageBySeqsReq, sdkws.PullMessageBySeqsResp]("/chat/pull_msg_by_seq")
)

var (
	ClearConversationMsg             = api[msg.ClearConversationsMsgReq, msg.ClearConversationsMsgResp]("/msg/clear_conversation_msg") // Clear the message of the specified conversation
	ClearAllMsg                      = api[msg.UserClearAllMsgReq, msg.UserClearAllMsgResp]("/msg/user_clear_all_msg")                 // Clear all messages of the current user
	DeleteMsgs                       = api[msg.DeleteMsgsReq, msg.DeleteMsgsResp]("/msg/delete_msgs")                                  // Delete the specified message
	RevokeMsg                        = api[msg.RevokeMsgReq, msg.RevokeMsgResp]("/msg/revoke_msg")
	MarkMsgsAsRead                   = api[msg.MarkMsgsAsReadReq, msg.MarkMsgsAsReadResp]("/msg/mark_msgs_as_read")
	GetConversationsHasReadAndMaxSeq = api[msg.GetConversationsHasReadAndMaxSeqReq, msg.GetConversationsHasReadAndMaxSeqResp]("/msg/get_conversations_has_read_and_max_seq")
	MarkConversationAsRead           = api[msg.MarkConversationAsReadReq, msg.MarkConversationAsReadResp]("/msg/mark_conversation_as_read")
	SetConversationHasReadSeq        = api[msg.SetConversationHasReadSeqReq, msg.SetConversationHasReadSeqResp]("/msg/set_conversation_has_read_seq")
	SendMsg                          = api[msg.SendMsgReq, msg.SendMsgResp]("/msg/send_msg")
	GetServerTime                    = api[msg.GetServerTimeReq, msg.GetServerTimeResp]("/msg/get_server_time")
)

var (
	CreateGroup                    = api[group.CreateGroupReq, group.CreateGroupResp]("/group/create_group")
	SetGroupInfo                   = api[group.SetGroupInfoReq, group.SetGroupInfoResp]("/group/set_group_info")
	JoinGroup                      = api[group.JoinGroupReq, group.JoinGroupResp]("/group/join_group")
	QuitGroup                      = api[group.QuitGroupReq, group.QuitGroupResp]("/group/quit_group")
	GetGroupsInfo                  = api[group.GetGroupsInfoReq, group.GetGroupsInfoResp]("/group/get_groups_info")
	GetGroupMemberList             = api[group.GetGroupMemberListReq, group.GetGroupMemberListResp]("/group/get_group_member_list")
	GetGroupMembersInfo            = api[group.GetGroupMembersInfoReq, group.GetGroupMembersInfoResp]("/group/get_group_members_info")
	InviteUserToGroup              = api[group.InviteUserToGroupReq, group.InviteUserToGroupResp]("/group/invite_user_to_group")
	GetJoinedGroupList             = api[group.GetJoinedGroupListReq, group.GetJoinedGroupListResp]("/group/get_joined_group_list")
	KickGroupMember                = api[group.KickGroupMemberReq, group.KickGroupMemberResp]("/group/kick_group")
	TransferGroup                  = api[group.TransferGroupOwnerReq, group.TransferGroupOwnerResp]("/group/transfer_group")
	GetRecvGroupApplicationList    = api[group.GetGroupApplicationListReq, group.GetGroupApplicationListResp]("/group/get_recv_group_applicationList")
	GetSendGroupApplicationList    = api[group.GetUserReqApplicationListReq, group.GetUserReqApplicationListResp]("/group/get_user_req_group_applicationList")
	AcceptGroupApplication         = api[group.GroupApplicationResponseReq, group.GroupApplicationResponseResp]("/group/group_application_response")
	DismissGroup                   = api[group.DismissGroupReq, group.DismissGroupResp]("/group/dismiss_group")
	MuteGroupMember                = api[group.MuteGroupMemberReq, group.MuteGroupMemberResp]("/group/mute_group_member")
	CancelMuteGroupMember          = api[group.CancelMuteGroupMemberReq, group.CancelMuteGroupMemberResp]("/group/cancel_mute_group_member")
	MuteGroup                      = api[group.MuteGroupReq, group.MuteGroupResp]("/group/mute_group")
	CancelMuteGroup                = api[group.CancelMuteGroupReq, group.CancelMuteGroupResp]("/group/cancel_mute_group")
	SetGroupMemberInfo             = api[group.SetGroupMemberInfoReq, group.SetGroupMemberInfoResp]("/group/set_group_member_info")
	GetIncrementalJoinGroup        = api[group.GetIncrementalJoinGroupReq, group.GetIncrementalJoinGroupResp]("/group/get_incremental_join_groups")
	GetIncrementalGroupMemberBatch = api[group.BatchGetIncrementalGroupMemberReq, group.BatchGetIncrementalGroupMemberResp]("/group/get_incremental_group_members_batch")
	GetFullJoinedGroupIDs          = api[group.GetFullJoinGroupIDsReq, group.GetFullJoinGroupIDsResp]("/group/get_full_join_group_ids")
	GetFullGroupMemberUserIDs      = api[group.GetFullGroupMemberUserIDsReq, group.GetFullGroupMemberUserIDsResp]("/group/get_full_group_member_user_ids")
)

var (
	GetConversations           = api[conversation.GetConversationsReq, conversation.GetConversationsResp]("/conversation/get_conversations")
	GetAllConversations        = api[conversation.GetAllConversationsReq, conversation.GetAllConversationsResp]("/conversation/get_all_conversations")
	SetConversations           = api[conversation.SetConversationsReq, conversation.SetConversationsResp]("/conversation/set_conversations")
	GetIncrementalConversation = api[conversation.GetIncrementalConversationReq, conversation.GetIncrementalConversationResp]("/conversation/get_incremental_conversations")
	GetFullConversationIDs     = api[conversation.GetFullOwnerConversationIDsReq, conversation.GetFullOwnerConversationIDsResp]("/conversation/get_full_conversation_ids")
	GetOwnerConversation       = api[conversation.GetOwnerConversationReq, conversation.GetOwnerConversationResp]("/conversation/get_owner_conversation")
)

var (
	GetUsersToken = api[auth.UserTokenReq, auth.UserTokenResp]("/auth/user_token")
)

var (
	FcmUpdateToken = api[third.FcmUpdateTokenReq, third.FcmUpdateTokenResp]("/third/fcm_update_token")
	SetAppBadge    = api[third.SetAppBadgeReq, third.SetAppBadgeResp]("/third/set_app_badge")
	UploadLogs     = api[third.UploadLogsReq, third.UploadLogsResp]("/third/logs/upload")
)

var (
	ObjectPartLimit               = api[third.PartLimitReq, third.PartLimitResp]("/object/part_limit")
	ObjectInitiateMultipartUpload = api[third.InitiateMultipartUploadReq, third.InitiateMultipartUploadResp]("/object/initiate_multipart_upload")
	ObjectAuthSign                = api[third.AuthSignReq, third.AuthSignResp]("/object/auth_sign")
	ObjectCompleteMultipartUpload = api[third.CompleteMultipartUploadReq, third.CompleteMultipartUploadResp]("/object/complete_multipart_upload")
	ObjectAccessURL               = api[third.AccessURLReq, third.AccessURLResp]("/object/access_url")
)
