package constant

import "errors"

var (
	ErrServer = ErrInfo{500, "server error"}

	ErrTokenDifferentPlatformID = ErrInfo{707, TokenDifferentPlatformIDMsg.Error()}
	ErrTokenDifferentUserID     = ErrInfo{708, TokenDifferentUserIDMsg.Error()}

	ErrStatus                = ErrInfo{ErrCode: 804, ErrMsg: StatusMsg.Error()}
	ErrCallback              = ErrInfo{ErrCode: 809, ErrMsg: CallBackMsg.Error()}
	ErrSendLimit             = ErrInfo{ErrCode: 810, ErrMsg: "send msg limit, to many request, try again later"}
	ErrMessageHasReadDisable = ErrInfo{ErrCode: 811, ErrMsg: "message has read disable"}
	ErrInternal              = ErrInfo{ErrCode: 812, ErrMsg: "internal error"}
)

var (
	TokenDifferentPlatformIDMsg = errors.New("different platformID")
	TokenDifferentUserIDMsg     = errors.New("different userID")

	StatusMsg = errors.New("status is abnormal")

	CallBackMsg = errors.New("callback failed")
)

const (
	NoError              = 0
	FormattingError      = 10001
	HasRegistered        = 10002
	NotRegistered        = 10003
	PasswordErr          = 10004
	GetIMTokenErr        = 10005
	RepeatSendCode       = 10006
	MailSendCodeErr      = 10007
	SmsSendCodeErr       = 10008
	CodeInvalidOrExpired = 10009
	RegisterFailed       = 10010
	ResetPasswordFailed  = 10011
	DatabaseError        = 10002
	ServerError          = 10004
	HttpError            = 10005
	IoError              = 10006
	IntentionalError     = 10007
)

func (e *ErrInfo) Code() int32 {
	return e.ErrCode
}
