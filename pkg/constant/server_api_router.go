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

	SendMsgRouter          = "/chat/send_msg"
	PullUserMsgRouter      = "/chat/pull_msg"
	PullUserMsgBySeqRouter = "/chat/pull_msg_by_seq"
	NewestSeqRouter        = "/chat/newest_seq"
	DeleteMsgRouter        = "/chat/del_msg"

	TencentCloudStorageCredentialRouter = "/third/tencent_cloud_storage_credential"
	MinioStorageCredentialRouter        = "/third/minio_storage_credential"

	//group
	CreateGroupRouter                 = RouterGroup + "/create_group"
	SetGroupInfoRouter                = RouterGroup + "/set_group_info"
	JoinGroupRouter                   = RouterGroup + "/join_group"
	QuitGroupRouter                   = RouterGroup + "/quit_group"
	GetGroupsInfoRouter               = RouterGroup + "/get_groups_info"
	GetGroupAllMemberListRouter       = RouterGroup + "/get_group_all_member_list"
	GetGroupMembersInfoRouter         = RouterGroup + "/get_group_members_info"
	InviteUserToGroupRouter           = RouterGroup + "/invite_user_to_group"
	GetJoinedGroupListRouter          = RouterGroup + "/get_joined_group_list"
	KickGroupMemberRouter             = RouterGroup + "/kick_group"
	TransferGroupRouter               = RouterGroup + "/transfer_group"
	GetRecvGroupApplicationListRouter = RouterGroup + "/get_recv_group_applicationList"
	GetSendGroupApplicationListRouter = RouterGroup + "/get_user_req_group_applicationList"
	AcceptGroupApplicationRouter      = RouterGroup + "/group_application_response"
	RefuseGroupApplicationRouter      = RouterGroup + "/group_application_response"

	//conversation
	SetConversationOptRouter   = ConversationGroup + "/set_conversation"
	GetConversationsRouter     = ConversationGroup + "/get_conversations"
	GetAllConversationsRouter  = ConversationGroup + "/get_all_conversations"
	GetConversationRouter      = ConversationGroup + "/get_conversation"
	BatchSetConversationRouter = ConversationGroup + "/batch_set_conversation"
)
const (
	RouterGroup       = "/group"
	ConversationGroup = "/conversation"
)
