package constant

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

	CmdJoinedSuperGroup = "018"

	CmdMaxSeq  = "maxSeq"
	CmdPushMsg = "pushMsg"
	CmdLogout  = "Logout"
	CmdWakeUp  = "wakeUp"
)

const (
	//ContentType
	Text                = 101
	Picture             = 102
	Voice               = 103
	Video               = 104
	File                = 105
	AtText              = 106
	Merger              = 107
	Card                = 108
	Location            = 109
	Custom              = 110
	Revoke              = 111
	HasReadReceipt      = 112
	Typing              = 113
	Quote               = 114
	Face                = 115
	GroupHasReadReceipt = 116
	AdvancedText        = 117
	//////////////////////////////////////////
	NotificationBegin       = 1000
	FriendNotificationBegin = 1200

	FriendApplicationApprovedNotification = 1201 //add_friend_response
	FriendApplicationRejectedNotification = 1202 //add_friend_response
	FriendApplicationNotification         = 1203 //add_friend
	FriendAddedNotification               = 1204
	FriendDeletedNotification             = 1205 //delete_friend
	FriendRemarkSetNotification           = 1206 //set_friend_remark?
	BlackAddedNotification                = 1207 //add_black
	BlackDeletedNotification              = 1208 //remove_black
	FriendNotificationEnd                 = 1299
	ConversationChangeNotification        = 1300

	UserNotificationBegin       = 1301
	UserInfoUpdatedNotification = 1303 //SetSelfInfoTip             = 204
	UserNotificationEnd         = 1399
	OANotification              = 1400

	GroupNotificationBegin = 1500

	GroupCreatedNotification                 = 1501
	GroupInfoSetNotification                 = 1502
	JoinGroupApplicationNotification         = 1503
	MemberQuitNotification                   = 1504
	GroupApplicationAcceptedNotification     = 1505
	GroupApplicationRejectedNotification     = 1506
	GroupOwnerTransferredNotification        = 1507
	MemberKickedNotification                 = 1508
	MemberInvitedNotification                = 1509
	MemberEnterNotification                  = 1510
	GroupDismissedNotification               = 1511
	GroupMemberMutedNotification             = 1512
	GroupMemberCancelMutedNotification       = 1513
	GroupMutedNotification                   = 1514
	GroupCancelMutedNotification             = 1515
	GroupMemberInfoSetNotification           = 1516
	GroupMemberSetToAdminNotification        = 1517
	GroupMemberSetToOrdinaryUserNotification = 1518
	GroupNotificationEnd                     = 1599

	SignalingNotificationBegin = 1600
	SignalingNotification      = 1601
	SignalingNotificationEnd   = 1649

	SuperGroupNotificationBegin  = 1650
	SuperGroupUpdateNotification = 1651
	SuperGroupNotificationEnd    = 1699

	ConversationPrivateChatNotification = 1701
	OrganizationChangedNotification     = 1801

	WorkMomentNotificationBegin = 1900
	WorkMomentNotification      = 1901

	NotificationEnd = 2000

	////////////////////////////////////////

	//MsgFrom
	UserMsgType = 100
	SysMsgType  = 200

	/////////////////////////////////////
	//SessionType
	SingleChatType       = 1
	GroupChatType        = 2
	SuperGroupChatType   = 3
	NotificationChatType = 4

	//MsgStatus
	MsgStatusDefault     = 0
	MsgStatusSending     = 1
	MsgStatusSendSuccess = 2
	MsgStatusSendFailed  = 3
	MsgStatusHasDeleted  = 4
	MsgStatusRevoked     = 5
	MsgStatusFiltered    = 6

	//OptionsKey
	IsHistory                  = "history"
	IsPersistent               = "persistent"
	IsUnreadCount              = "unreadCount"
	IsConversationUpdate       = "conversationUpdate"
	IsOfflinePush              = "offlinePush"
	IsSenderSync               = "senderSync"
	IsNotPrivate               = "notPrivate"
	IsSenderConversationUpdate = "senderConversationUpdate"
	IsSenderNotificationPush   = "senderNotificationPush"

	//GroupStatus
	GroupOk              = 0
	GroupBanChat         = 1
	GroupStatusDismissed = 2
	GroupStatusMuted     = 3

	// workMoment permission
	WorkMomentPublic            = 0
	WorkMomentPrivate           = 1
	WorkMomentPermissionCanSee  = 2
	WorkMomentPermissionCantSee = 3

	// workMoment sdk notification type
	WorkMomentCommentNotification = 0
	WorkMomentLikeNotification    = 1
	WorkMomentAtUserNotification  = 2
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
	BlackRelationship  = 0
	FriendRelationship = 1
)

//const (
//	ErrCodeInitLogin    = 1001
//	ErrCodeFriend       = 2001
//	ErrCodeConversation = 3001
//	ErrCodeUserInfo     = 4001
//	ErrCodeGroup        = 5001
//)
const (
	NormalGroup                       = 0
	SuperGroup                        = 1
	SuperGroupTableName               = "local_super_groups"
	SuperGroupErrChatLogsTableNamePre = "local_sg_err_chat_logs_"
	SuperGroupChatLogsTableNamePre    = "local_sg_chat_logs_"
)

const (
	SdkInit      = 0
	LoginSuccess = 101
	Logining     = 102
	LoginFailed  = 103

	Logout = 201

	TokenFailedExpired       = 701
	TokenFailedInvalid       = 702
	TokenFailedKickedOffline = 703
)

const (
	DeFaultSuccessMsg = "ok"
)

const (
	AddConOrUpLatMsg          = 2
	UnreadCountSetZero        = 3
	IncrUnread                = 5
	TotalUnreadMessageChanged = 6
	UpdateFaceUrlAndNickName  = 7
	UpdateLatestMessageChange = 8
	ConChange                 = 9
	NewCon                    = 10

	HasRead = 1
	NotRead = 0

	IsFilter  = 1
	NotFilter = 0
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
const MaxTotalMsgLen = 51200

//const MaxTotalMsgLen = 20480
const (
	FriendAcceptTip  = "You have successfully become friends, so start chatting"
	TransferGroupTip = "The owner of the group is transferred!"
	AcceptGroupTip   = "%s join the group"
)

const (
	WSGetNewestSeq     = 1001
	WSPullMsgBySeqList = 1002
	WSSendMsg          = 1003
	WSSendSignalMsg    = 1004
	WsDelMsg           = 1005
	WSPushMsg          = 2001
	WSKickOnlineMsg    = 2002
	WsLogoutMsg        = 2003

	WSDataError = 3001
)

// conversation
const (
	//MsgReceiveOpt
	ReceiveMessage          = 0
	NotReceiveMessage       = 1
	ReceiveNotNotifyMessage = 2

	//pinned
	Pinned    = 1
	NotPinned = 0

	//privateChat
	IsPrivateChat  = true
	NotPrivateChat = false
)

const SuccessCallbackDefault = ""

const (
	AppOrdinaryUsers = 1
	AppAdmin         = 2

	GroupOrdinaryUsers = 1
	GroupOwner         = 2
	GroupAdmin         = 3

	GroupResponseAgree  = 1
	GroupResponseRefuse = -1

	FriendResponseAgree  = 1
	FriendResponseRefuse = -1

	Male   = 1
	Female = 2
)
const (
	AtAllString = "AtAllTag"
	AtNormal    = 0
	AtMe        = 1
	AtAll       = 2
	AtAllAtMe   = 3
)
const (
	FieldRecvMsgOpt    = 1
	FieldIsPinned      = 2
	FieldAttachedInfo  = 3
	FieldIsPrivateChat = 4
	FieldGroupAtType   = 5
	FieldIsNotInGroup  = 6
	FieldEx            = 7
)
const (
	KeywordMatchOr  = 0
	KeywordMatchAnd = 1
)

const BigVersion = "v2"
const UpdateVersion = ".0.0"
const SdkVersion = "Open-IM-SDK-Core-"
const LogFileName = "sdk"

var HeartbeatInterval = 5

const (
	MsgSyncModelDefault  = 0   //SyncFlag
	MsgSyncModelLogin    = 1   //SyncFlag
	SyncOrderStartLatest = 101 //PullMsgOrder
)
