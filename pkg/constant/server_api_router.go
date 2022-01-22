package constant

const (
	GetSelfUserInfoRouter    = "/user/get_self_user_info"
	GetUsersInfoRouter       = "/user/get_users_info"
	UpdateSelfUserInfoRouter = "/user/update_user_info"

	AddFriendRouter                    = "/friend/add_friend"
	DeleteFriendRouter                 = "/friend/delete_friend"
	GetFriendApplicationListRouter     = "/friend/get_friend_apply_list"      //recv
	GetSelfFriendApplicationListRouter = "/friend/get_self_friend_apply_list" //send
	GetFriendListRouter                = "/friend/get_friend_list"
	AddFriendResponse                  = "/friend/add_friend_response"
	SetFriendRemark                    = "/friend/set_friend_remark"

	AddBlackRouter     = "/friend/add_black"
	RemoveBlackRouter  = "/friend/remove_black"
	GetBlackListRouter = "/friend/get_black_list"

	SendMsgRouter                       = "/chat/send_msg"
	PullUserMsgRouter                   = "/chat/pull_msg"
	PullUserMsgBySeqRouter              = "/chat/pull_msg_by_seq"
	NewestSeqRouter                     = "/chat/newest_seq"
	TencentCloudStorageCredentialRouter = "/third/tencent_cloud_storage_credential"

	//group
	CreateGroupRouter                  = "/group/create_group"
	SetGroupInfoRouter                 = "/group/set_group_info"
	JoinGroupRouter                    = "/group/join_group"
	QuitGroupRouter                    = "/group/quit_group"
	GetGroupsInfoRouter                = "/group/get_groups_info"
	GetGroupAllMemberListRouter        = "/group/get_group_all_member_list"
	GetGroupMembersInfoRouter          = "/group/get_group_members_info"
	InviteUserToGroupRouter            = "/group/invite_user_to_group"
	GetJoinedGroupListRouter           = "/group/get_joined_group_list"
	KickGroupMemberRouter              = "/group/kick_group"
	TransferGroupRouter                = "/group/transfer_group"
	GetRecvGroupApplicationListRouter  = "/group/get_recv_group_applicationList"
	AcceptGroupApplicationRouter       = "/group/group_application_response"
	RefuseGroupApplicationRouter       = "/group/group_application_response"
	SetReceiveMessageOptRouter         = "/conversation/set_receive_message_opt"
	GetReceiveMessageOptRouter         = "/conversation/get_receive_message_opt"
	GetAllConversationMessageOptRouter = "/conversation/get_all_conversation_message_opt"
)
