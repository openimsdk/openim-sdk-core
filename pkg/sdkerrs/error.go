package sdkerrs

import "github.com/OpenIMSDK/Open-IM-Server/pkg/errs"

func New(code int, msg string, dtl string) errs.CodeError {
	return errs.NewCodeError(code, msg).WithDetail(dtl)
}

var Warp = errs.Wrap
