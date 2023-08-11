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
	CmdFriend                     = "001"
	CmdBlackList                  = "002"
	CmdNotification               = "003"
	CmdDeleteConversation         = "004"
	CmdNewMsgCome                 = "005"
	CmdSuperGroupMsgCome          = "006"
	CmdUpdateConversation         = "007"
	CmSyncReactionExtensions      = "008"
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
	CmdUpdateMessage    = "019"

	CmdReconnect = "020"
	CmdInit      = "021"

	CmdMaxSeq       = "maxSeq"
	CmdPushMsg      = "pushMsg"
	CmdConnSuccesss = "connSuccess"
	CmdWakeUp       = "wakeUp"
	CmdLogOut       = "loginOut"
)

const (
	//ContentType
	Text                            = 101
	Picture                         = 102
	Sound                           = 103
	Video                           = 104
	File                            = 105
	AtText                          = 106
	Merger                          = 107
	Card                            = 108
	Location                        = 109
	Custom                          = 110
	Typing                          = 113
	Quote                           = 114
	Face                            = 115
	AdvancedText                    = 117
	CustomMsgNotTriggerConversation = 119
	CustomMsgOnlineOnly             = 120
	ReactionMessageModifier         = 121
	ReactionMessageDeleter          = 122

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
	FriendInfoUpdatedNotification         = 1209
	FriendNotificationEnd                 = 1299
	ConversationChangeNotification        = 1300

	UserNotificationBegin        = 1301
	UserInfoUpdatedNotification  = 1303 //SetSelfInfoTip             = 204
	UserStatusChangeNotification = 1304
	UserNotificationEnd          = 1399
	OANotification               = 1400

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

	SuperGroupNotificationBegin         = 1650
	SuperGroupUpdateNotification        = 1651
	MsgDeleteNotification               = 1652
	ReactionMessageModifierNotification = 1653
	ReactionMessageDeleteNotification   = 1654
	SuperGroupNotificationEnd           = 1699

	ConversationPrivateChatNotification = 1701
	ConversationUnreadNotification      = 1702

	WorkMomentNotificationBegin = 1900
	WorkMomentNotification      = 1901

	BusinessNotificationBegin = 2000
	BusinessNotification      = 2001
	BusinessNotificationEnd   = 2099

	RevokeNotification = 2101

	HasReadReceiptNotification      = 2150
	GroupHasReadReceiptNotification = 2155
	ClearConversationNotification   = 2101
	DeleteMsgsNotification          = 2102

	HasReadReceipt = 2200

	NotificationEnd = 5000

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
	MsgStatusDefault = 0

	MsgStatusSending     = 1
	MsgStatusSendSuccess = 2
	MsgStatusSendFailed  = 3
	MsgStatusHasDeleted  = 4
	MsgStatusFiltered    = 5

	//OptionsKey
	IsHistory                  = "history"
	IsPersistent               = "persistent"
	IsUnreadCount              = "unreadCount"
	IsConversationUpdate       = "conversationUpdate"
	IsOfflinePush              = "offlinePush"
	IsSenderSync               = "senderSync"
	IsNotPrivate               = "notPrivate"
	IsSenderConversationUpdate = "senderConversationUpdate"

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

// const (
//
//	ErrCodeInitLogin    = 1001
//	ErrCodeFriend       = 2001
//	ErrCodeConversation = 3001
//	ErrCodeUserInfo     = 4001
//	ErrCodeGroup        = 5001
//
// )
const (
	NormalGroup                       = 0
	SuperGroup                        = 1
	WorkingGroup                      = 2
	SuperGroupTableName               = "local_super_groups"
	SuperGroupErrChatLogsTableNamePre = "local_sg_err_chat_logs_"
	ChatLogsTableNamePre              = "chat_logs_"
)

const (
	SdkInit = 0

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
	AddConOrUpLatMsg                  = 2
	UnreadCountSetZero                = 3
	IncrUnread                        = 5
	TotalUnreadMessageChanged         = 6
	UpdateConFaceUrlAndNickName       = 7
	UpdateLatestMessageChange         = 8
	ConChange                         = 9
	NewCon                            = 10
	ConChangeDirect                   = 11
	NewConDirect                      = 12
	ConversationLatestMsgHasRead      = 13
	UpdateMsgFaceUrlAndNickName       = 14
	SyncConversation                  = 15
	SyncMessageListReactionExtensions = 16
	SyncMessageListTypeKeyInfo        = 17

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

// const MaxTotalMsgLen = 20480
const (
	FriendAcceptTip  = "You have successfully become friends, so start chatting"
	TransferGroupTip = "The owner of the group is transferred!"
	AcceptGroupTip   = "%s join the group"
)

const (
	GetNewestSeq        = 1001
	PullMsgBySeqList    = 1002
	SendMsg             = 1003
	SendSignalMsg       = 1004
	DelMsg              = 1005
	PushMsg             = 2001
	KickOnlineMsg       = 2002
	LogoutMsg           = 2003
	SetBackgroundStatus = 2004

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

const SuccessCallbackDefault = "" // Default value for success callback

const (
	AppOrdinaryUsers = 1 // Application user type: ordinary user
	AppAdmin         = 2 // Application user type: administrator

	GroupOwner         = 100 // Group member type: owner
	GroupAdmin         = 60  // Group member type: administrator
	GroupOrdinaryUsers = 20  // Group member type: ordinary user

	GroupFilterAll                   = 0
	GroupFilterOwner                 = 1
	GroupFilterAdmin                 = 2
	GroupFilterOrdinaryUsers         = 3
	GroupFilterAdminAndOrdinaryUsers = 4
	GroupFilterOwnerAndAdmin         = 5

	GroupResponseAgree  = 1  // Response to group application: agree
	GroupResponseRefuse = -1 // Response to group application: refuse

	FriendResponseAgree   = 1  // Response to friend request: agree
	FriendResponseRefuse  = -1 // Response to friend request: refuse
	FriendResponseDefault = 0

	Male   = 1 // Gender: male
	Female = 2 // Gender: female
)
const (
	AtAllString = "AtAllTag" // String for 'all people' mention tag
	AtNormal    = 0          // Mention mode: normal
	AtMe        = 1          // Mention mode: mention sender only
	AtAll       = 2          // Mention mode: mention all people
	AtAllAtMe   = 3          // Mention mode: mention all people and sender

)
const (
	FieldRecvMsgOpt    = 1 // Field type: message receiving options
	FieldIsPinned      = 2 // Field type: whether a message is pinned
	FieldAttachedInfo  = 3 // Field type: attached information
	FieldIsPrivateChat = 4 // Field type: whether a message is from a private chat
	FieldGroupAtType   = 5 // Field type: group mention mode
	FieldIsNotInGroup  = 6 // Field type: whether a message is not in a group
	FieldEx            = 7 // Field type: extension field
	FieldUnread        = 8 // Field type: whether a message is unread
	FieldBurnDuration  = 9 // Field type: message burn duration
)
const (
	SetMessageExtensions = 1 // Message extension operation type: set extension
	AddMessageExtensions = 2 // Message extension operation type: add extension
)
const (
	KeywordMatchOr  = 0 // Keyword match mode: match any keyword
	KeywordMatchAnd = 1 // Keyword match mode: match all keywords
)

const BigVersion = "v3"
const UpdateVersion = ".0.0"
const SdkVersion = "openim-sdk-core-"
const LogFileName = "sdk"

func GetSdkVersion() string {
	return SdkVersion + BigVersion + UpdateVersion
}

var HeartbeatInterval = 5

const (
	MsgSyncModelDefault  = 0   //SyncFlag
	MsgSyncModelLogin    = 1   //SyncFlag
	SyncOrderStartLatest = 101 //PullMsgOrder

	MsgSyncBegin      = 1001 //
	MsgSyncProcessing = 1002 //
	MsgSyncEnd        = 1003 //
	MsgSyncFailed     = 1004
)

const (
	JoinByInvitation = 2
	JoinBySearch     = 3
	JoinByQRCode     = 4
)
const (
	SplitPullMsgNum              = 100
	PullMsgNumWhenLogin          = 10000
	PullMsgNumForReadDiffusion   = 50
	NormalMsgMinNumReadDiffusion = 100
)

const SplitGetGroupMemberNum = 1000
const UseHashGroupMemberNum = 1000

const (
	Uninitialized    = -1001
	NoNetwork        = 1 //有网络->无网络
	NetworkAvailable = 2 //无网络->有网络
	NetworkVariation = 3 //有网络，但状态有变化
)
