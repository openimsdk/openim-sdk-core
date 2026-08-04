// Copyright © 2023 OpenIM SDK. All rights reserved.

//go:build js && wasm

package open_im_sdk

import (
	"context"
	"sync/atomic"

	"github.com/openimsdk/tools/log"
)

// handleKickedOffline 仅通知业务层。
// 不在此处 DispatchLogout：logout 的关连接/关库/initResources 若紧跟在
// OnKickedOffline（H5 location.replace）之后执行，Mobile Safari 仍会卡死主线程。
// 页面硬跳卸载 WASM 即可；业务侧勿再二次调用 IMSDK.logout()。
func (c *apiErrCallback) handleKickedOffline(ctx context.Context, err error) {
	if !atomic.CompareAndSwapInt32(&c.kickedOfflineState, 0, 1) {
		return
	}
	log.ZError(ctx, "OnKickedOffline callback", err)
	c.listener().OnKickedOffline()
}
