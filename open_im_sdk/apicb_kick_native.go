// Copyright © 2023 OpenIM SDK. All rights reserved.

//go:build !(js && wasm)

package open_im_sdk

import (
	"context"
	"sync/atomic"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/tools/log"
)

func (c *apiErrCallback) handleKickedOffline(ctx context.Context, err error) {
	if !atomic.CompareAndSwapInt32(&c.kickedOfflineState, 0, 1) {
		return
	}
	log.ZError(ctx, "OnKickedOffline callback", err)
	c.listener().OnKickedOffline()
	_ = common.DispatchLogout(ctx, c.loginMgrCh)
}
