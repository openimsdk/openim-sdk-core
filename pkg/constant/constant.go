package constant

import (
	"errors"
	"open_im_sdk/pkg/utils"
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

	CmdLogout = "Logout"
)

const (
	//ContentType
	Text           = 101
	Picture        = 102
	Voice          = 103
	Video          = 104
	File           = 105
	AtText         = 106
	Merger         = 107
	Card           = 108
	Location       = 109
	Custom         = 110
	Revoke         = 111
	HasReadReceipt = 112
	Typing         = 113
	Quote          = 114
	//////////////////////////////////////////
	SingleTipBegin             = 200
	AcceptFriendApplicationTip = 201
	AddFriendTip               = 202
	RefuseFriendApplicationTip = 203
	SetSelfInfoTip             = 204

	SingleTipEnd = 399
	/////////////////////////////////////////
	GroupTipBegin             = 500
	TransferGroupOwnerTip     = 501
	CreateGroupTip            = 502
	JoinGroupTip              = 504
	QuitGroupTip              = 505
	SetGroupInfoTip           = 506
	AcceptGroupApplicationTip = 507
	RefuseGroupApplicationTip = 508
	KickGroupMemberTip        = 509
	InviteUserToGroupTip      = 510

	GroupTipEnd = 599
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
	MsgStatusRevoked     = 5
	MsgStatusFiltered    = 6

	//OptionsKey
	IsHistory            = "history"
	IsPersistent         = "persistent"
	IsUnreadCount        = "unreadCount"
	IsConversationUpdate = "conversationUpdate"
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
	SdkInit      = 0
	LoginSuccess = 101
	Logining     = 102
	LoginFailed  = 103

	LogoutCmd = 201

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
	NewConChange              = 9
	NewCon                    = 10

	HasRead = 1
	NotRead = 0

	IsFilter  = 1
	NotFilter = 0

	Pinned    = 1
	NotPinned = 0
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
const MaxTotalMsgLen = 2048
const (
	FriendAcceptTip  = "You have successfully become friends, so start chatting"
	TransferGroupTip = "The owner of the group is transferred!"
	AcceptGroupTip   = "%s join the group"
)

const (
	WSGetNewestSeq     = 1001
	WSPullMsg          = 1002
	WSSendMsg          = 1003
	WSPullMsgBySeqList = 1004
	WSPushMsg          = 2001
	WSKickOnlineMsg    = 2002
	WsLogoutMsg = 2003
	WSDataError        = 3001

)

const (
	//MsgReceiveOpt
	ReceiveMessage          = 0
	NotReceiveMessage       = 1
	ReceiveNotNotifyMessage = 2
)

// key = errCode, string = errMsg
type ErrInfo struct {
	ErrCode int32
	ErrMsg  string
}

var (
	OK = ErrInfo{0, ""}

	ErrParseToken = ErrInfo{200, ParseTokenMsg.Error()}

	ErrTencentCredential = ErrInfo{400, ThirdPartyMsg.Error()}

	ErrTokenExpired     = ErrInfo{701, TokenExpiredMsg.Error()}
	ErrTokenInvalid     = ErrInfo{702, TokenInvalidMsg.Error()}
	ErrTokenMalformed   = ErrInfo{703, TokenMalformedMsg.Error()}
	ErrTokenNotValidYet = ErrInfo{704, TokenNotValidYetMsg.Error()}
	ErrTokenUnknown     = ErrInfo{705, TokenUnknownMsg.Error()}

	ErrAccess = ErrInfo{ErrCode: 801, ErrMsg: AccessMsg.Error()}
	ErrDB     = ErrInfo{ErrCode: 802, ErrMsg: DBMsg.Error()}
	ErrArgs   = ErrInfo{ErrCode: 803, ErrMsg: ArgsMsg.Error()}
	ErrApi    = ErrInfo{ErrCode: 804, ErrMsg: ApiMsg.Error()}
	ErrData = ErrInfo{ErrCode: 805, ErrMsg: DataMsg.Error()}
	ErrLogin = ErrInfo{ErrCode: 806, ErrMsg: LoginMsg.Error()}

	ErrWsRecvConnDiff = ErrInfo{ErrCode: 901, ErrMsg: WsRecvConnDiff.Error()}
	ErrWsRecvConnSame = ErrInfo{ErrCode: 902, ErrMsg: WsRecvConnSame.Error()}
	ErrWsRecvCode     = ErrInfo{ErrCode: 903, ErrMsg: WsRecvCode.Error()}
)

var (
	ParseTokenMsg       = errors.New("parse token failed")
	TokenExpiredMsg     = errors.New("token is timed out, please log in again")
	TokenInvalidMsg     = errors.New("token has been invalidated")
	TokenNotValidYetMsg = errors.New("token not active yet")
	TokenMalformedMsg   = errors.New("that's not even a token")
	TokenUnknownMsg     = errors.New("couldn't handle this token")

	AccessMsg = errors.New("no permission")
	DBMsg     = errors.New("db failed")
	ArgsMsg   = errors.New("args failed")
	ApiMsg    = errors.New("api failed")
	DataMsg = errors.New("data failed ")
	LoginMsg = errors.New("you can only login once")

	ThirdPartyMsg = errors.New("third party error")

	WsRecvConnDiff = errors.New("recv timeout, conn diff")
	WsRecvConnSame = errors.New("recv timeout, conn same")
	WsRecvCode     = errors.New("recv code err")
)

func (e *ErrInfo) Error() string {
	return e.ErrMsg
}

const SuccessCallbackDefault = ""

var SvrConf utils.IMConfig