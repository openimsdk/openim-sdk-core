package open_im_sdk

import (
	"context"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"sync/atomic"
)

type apiErrCallback struct {
	loginMgrCh         chan common.Cmd2Value
	listener           open_im_sdk_callback.OnConnListener
	tokenExpiredState  int32
	kickedOfflineState int32
}

func (c *apiErrCallback) OnError(ctx context.Context, err error) {
	if err == nil {
		return
	}
	codeErr, ok := errs.Unwrap(err).(errs.CodeError)
	if !ok {
		return
	}
	switch codeErr.Code() {
	case
		errs.TokenExpiredError,
		errs.TokenInvalidError,
		errs.TokenMalformedError,
		errs.TokenNotValidYetError,
		errs.TokenUnknownError,
		errs.TokenNotExistError:
		if atomic.CompareAndSwapInt32(&c.tokenExpiredState, 0, 1) {
			c.listener.OnUserTokenExpired()
			_ = common.TriggerCmdLogOut(ctx, c.loginMgrCh)
		}
	case errs.TokenKickedError:
		if atomic.CompareAndSwapInt32(&c.kickedOfflineState, 0, 1) {
			c.listener.OnKickedOffline()
			_ = common.TriggerCmdLogOut(ctx, c.loginMgrCh)
		}
	}
}
