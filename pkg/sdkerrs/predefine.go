package sdkerrs

import "github.com/OpenIMSDK/Open-IM-Server/pkg/errs"

var (
	ErrArgs           = errs.NewCodeError(ArgsError, "ArgsError")
	ErrCtxDeadline    = errs.NewCodeError(CtxDeadlineExceededError, "CtxDeadlineExceededError")
	ErrSdkInternal    = errs.NewCodeError(SdkInternalError, "SdkInternalError")
	ErrNetwork        = errs.NewCodeError(NetworkError, "NetworkError")
	ErrNetworkTimeOut = errs.NewCodeError(NetworkTimeoutError, "NetworkTimeoutError")
	ErrDuplicateKey   = errs.NewCodeError(DuplicateKeyError, "DuplicateKeyError")

	ErrGroupIDNotFound = errs.NewCodeError(GroupIDNotFoundError, "GroupIDNotFoundError")
	ErrUserIDNotFound  = errs.NewCodeError(UserIDNotFoundError, "UserIDNotFoundError")

	ErrRecordNotFound = errs.NewCodeError(RecordNotFoundError, "RecordNotFoundError")
)
