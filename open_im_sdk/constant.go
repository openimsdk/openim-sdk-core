package open_im_sdk

func initAddr() {
	ginAddress = SvrConf.IpApiAddr

	getUserInfoRouter = ginAddress + "/user/get_user_info"
	updateUserInfoRouter = ginAddress + "/user/update_user_info"
	addFriendRouter = ginAddress + "/friend/add_friend"
	getFriendApplicationListRouter = ginAddress + "/friend/get_friend_apply_list"
	getSelfApplicationListRouter = ginAddress + "/friend/get_self_apply_list"
	deleteFriendRouter = ginAddress + "/friend/delete_friend"
	getFriendInfoRouter = ginAddress + "/friend/get_friends_info"
	getFriendListRouter = ginAddress + "/friend/get_friend_list"
	sendMsgRouter = ginAddress + "/chat/send_msg"
	getBlackListRouter = ginAddress + "/friend/get_blacklist"
	addFriendResponse = ginAddress + "/friend/add_friend_response"
	addBlackListRouter = ginAddress + "/friend/add_blacklist"
	removeBlackListRouter = ginAddress + "/friend/remove_blacklist"
	//getFriendApplyListRouter = ginAddress + "/friend/get_friend_apply_list"
	pullUserMsgRouter = ginAddress + "/chat/pull_msg"
	newestSeqRouter = ginAddress + "/chat/newest_seq"
	setFriendComment = ginAddress + "/friend/set_friend_comment"
	tencentCloudStorageCredentialRouter = ginAddress + "/third/tencent_cloud_storage_credential"

	getGroupMemberListRouter = ginAddress + "/group/get_group_member_list"
	getGroupMembersInfoRouter = ginAddress + "/group/get_group_members_info"
	inviteUserToGroupRouter = ginAddress + "/group/invite_user_to_group"
	getJoinedGroupListRouter = ginAddress + "/group/get_joined_group_list"
	kickGroupMemberRouter = ginAddress + "/group/kick_group"
}

var (
	ginAddress = "http://47.112.160.66:10000"

	getUserInfoRouter              = ginAddress + "/user/get_user_info"
	updateUserInfoRouter           = ginAddress + "/user/update_user_info"
	addFriendRouter                = ginAddress + "/friend/add_friend"
	getFriendInfoRouter            = ginAddress + "/friend/get_friends_info"
	getFriendApplicationListRouter = ginAddress + "/friend/get_friend_apply_list"
	getSelfApplicationListRouter   = ginAddress + "/friend/get_self_apply_list"
	deleteFriendRouter             = ginAddress + "/friend/delete_friend"
	getFriendListRouter            = ginAddress + "/friend/get_friend_list"
	sendMsgRouter                  = ginAddress + "/chat/send_msg"
	getBlackListRouter             = ginAddress + "/friend/get_blacklist"
	addFriendResponse              = ginAddress + "/friend/add_friend_response"
	addBlackListRouter             = ginAddress + "/friend/add_blacklist"
	removeBlackListRouter          = ginAddress + "/friend/remove_blacklist"
	//	getFriendApplyListRouter            = ginAddress + "/friend/get_friend_apply_list"
	setFriendComment                    = ginAddress + "/friend/set_friend_comment"
	pullUserMsgRouter                   = ginAddress + "/chat/pull_msg"
	newestSeqRouter                     = ginAddress + "/chat/newest_seq"
	tencentCloudStorageCredentialRouter = ginAddress + "/third/tencent_cloud_storage_credential"

	//group
	createGroupRouter             = ginAddress + "/group/create_group"
	setGroupInfoRouter            = ginAddress + "/group/set_group_info"
	joinGroupRouter               = ginAddress + "/group/join_group"
	quitGroupRouter               = ginAddress + "/group/quit_group"
	getGroupsInfoRouter           = ginAddress + "/group/get_groups_info"
	getGroupMemberListRouter      = ginAddress + ""
	getGroupMembersInfoRouter     = ginAddress + ""
	inviteUserToGroupRouter       = ginAddress + ""
	getJoinedGroupListRouter      = ginAddress + ""
	kickGroupMemberRouter         = ginAddress + ""
	transferGroupRouter           = ginAddress + "/group/transfer_group"
	getGroupApplicationListRouter = ginAddress + "/group/get_group_applicationList"
	acceptGroupApplicationRouter  = ginAddress + "/group/group_application_response"
	refuseGroupApplicationRouter  = ginAddress + "/group/group_application_response"
)

func initListenerCh() {
	ConversationCh = make(chan cmd2Value, 100)
	InitCh = make(chan cmd2Value, 50)

	ConListener.ch = ConversationCh
	SdkInitManager.ch = InitCh
}

var (
	//       chan cmd2Value //cmd：
	ConversationCh chan cmd2Value //cmd：
	InitCh         chan cmd2Value //cmd：
	groupCh        chan cmd2Value //group channel

	SvrConf  IMConfig
	WsState  int32 //100 stop，  0 init
	token    string
	LoginUid string

	SdkInitManager IMManager
	FriendObj      Friend
	ConListener    ConversationListener
	groupManager groupListener
)

const (
	CmdFriend                     = "001"
	CmdBlackList                  = "002"
	CmdFriendApplication          = "003"
	CmdDeleteConversation         = "004"
	CmdNewMsgCome                 = "005"
	CmdGeyLoginUserInfo           = "006"
	CmdUpdateConversation         = "007"
	CmdForceSyncFriend            = "008"
	CmdFroceSyncBlackList         = "009"
	CmdForceSyncFriendApplication = "010"
	CmdForceSyncMsg               = "011"
	CmdForceSyncLoginUerInfo      = "012"
	CmdReLogin                    = "013"
	CmdUnInit                     = "014"
	CmdAcceptFriend               = "015"
	CmdRefuseFriend               = "016"
	CmdAddFriend                  = "017"
)

const (
	MessageHasNotRead = 0
	MessageHasRead    = 1
)
const (
	//ContentType
	Text    = 101
	Picture = 102
	Sound   = 103
	Video   = 104
	File    = 105
	Merger  = 106

	SyncSenderMsg              = 110
	AcceptFriendApplicationTip = 201
	AddFriendTip               = 202
	RefuseFriendApplicationTip = 203
	SetSelfInfoTip             = 204
	RevokeMessageTip           = 205
	C2CMessageAsRead           = 206

	KickOnlineTip = 303

	TransferGroupOwnerTip           = 501
	CreateGroupTip                  = 502
	GroupApplicationResponseTip     = 503
	JoinGroupTip                    = 504
	QuitGroupTip                    = 505
	SetGroupInfoTip                 = 506
	AcceptGroupApplicationTip       = 507
	RefuseGroupApplicationTip       = 508
	KickGroupMemberTip              = 509
	InviteUserToGroupTip            = 510
	AcceptGroupApplicationResultTip = 511
	RefuseGroupApplicationResultTip = 512
	////////////////////////////////////////
	//MsgFrom
	UserMsgType = 100
	SysMsgType  = 200

	/////////////////////////////////////
	//SessionType
	SingleChatType = 1
	GroupChatType  = 2

	//MsgStatus
	MsgStatusSending     = 1
	MsgStatusSendSuccess = 2
	MsgStatusSendFailed  = 3
	MsgStatusHasDeleted  = 4
)

const (
	ckWsInitConnection  string = "ws-init-connection"
	ckWsLoginConnection string = "ws-login-connection"
	ckWsClose           string = "ws-close"
	ckWsKickOffLine     string = "ws-kick-off-line"
	ckTokenExpired      string = "token-expired"
	ckSelfInfoUpdate    string = "self-info-update"
)

const (
	ErrCodeInitLogin    = 1001
	ErrCodeFriend       = 2001
	ErrCodeConversation = 3001
	ErrCodeUserInfo     = 4001
	ErrCodeGroup        = 5001
)

const (
	LoginSuccess = 101
	Logining     = 102
	LoginFailed  = 103
)

const (
	DeFaultSuccessMsg = "ok"
)

const (
	ConAndUnreadChange = 1
	AddConOrUpLatMsg   = 2
	UnreadCountSetZero = 3
	ConChange          = 4
	IncrUnread         = 5

	HasRead = 1
	NotRead = 0
)

const (
	GroupActionCreateGroup            = 1
	GroupActionApplyJoinGroup         = 2
	GroupActionQuitGroup              = 3
	GroupActionSetGroupInfo           = 4
	GroupActionKickGroupMember        = 5
	GroupActionTransferGroupOwner     = 6
	GroupActionInviteUserToGroup      = 7
	GroupActionAcceptGroupApplication = 8
	GroupActionRefuseGroupApplication = 9
)
const ZoomScale = "200"
