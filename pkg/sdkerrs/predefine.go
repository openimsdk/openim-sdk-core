package sdkerrs

import "github.com/OpenIMSDK/Open-IM-Server/pkg/errs"

var (
	ErrArgs           = errs.NewCodeError(ArgsError, "ArgsError")
	ErrCtxDeadline    = errs.NewCodeError(CtxDeadlineExceededError, "CtxDeadlineExceededError")
	ErrSdkInternal    = errs.NewCodeError(SdkInternalError, "SdkInternalError")
	ErrNetwork        = errs.NewCodeError(NetworkError, "NetworkError")
	ErrNetworkTimeOut = errs.NewCodeError(NetworkTimeoutError, "NetworkTimeoutError")

	ErrGroupIDNotFound = errs.NewCodeError(GroupIDNotFoundError, "GroupIDNotFoundError")
	ErrUserIDNotFound  = errs.NewCodeError(UserIDNotFoundError, "UserIDNotFoundError")

	ErrResourceLoad = errs.NewCodeError(ResourceLoadNotCompleteError, "ResourceLoadNotCompleteError")

	//消息相关
	ErrFileNotFound             = errs.NewCodeError(FileNotFoundError, "RecordNotFoundError")
	ErrMsgDecodeBinaryWs        = errs.NewCodeError(MsgDecodeBinaryWsError, "MsgDecodeBinaryWsError")
	ErrMsgDeCompression         = errs.NewCodeError(MsgDeCompressionError, "MsgDeCompressionError")
	ErrMsgBinaryTypeNotSupport  = errs.NewCodeError(MsgBinaryTypeNotSupportError, "MsgTypeNotSupportError")
	ErrMsgRepeated              = errs.NewCodeError(MsgRepeatError, "only failed message can be repeatedly send")
	ErrMsgContentTypeNotSupport = errs.NewCodeError(MsgContentTypeNotSupportError, "contentType not support currently") // msg 	// msg

	//会话相关
	ErrNotSupportOpt = errs.NewCodeError(NotSupportOptError, "super group not support this opt")
	//群组相关

	ErrGroupType = errs.NewCodeError(GroupTypeErr, "group type error")

	ErrLoginOut    = errs.NewCodeError(LoginOutError, "LoginOutError")
	ErrLoginRepeat = errs.NewCodeError(LoginRepeatError, "LoginRepeatError")
)
