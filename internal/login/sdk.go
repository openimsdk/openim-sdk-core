// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package login

import (
	"context"
)

func (u *LoginMgr) Login(ctx context.Context, userID, token string) error {
	return u.login(ctx, userID, token)
}

func (u *LoginMgr) WakeUp(ctx context.Context) error {
	return u.wakeUp(ctx)
}

func (u *LoginMgr) Logout(ctx context.Context) error {
	return u.logout(ctx)
}

func (u *LoginMgr) SetAppBackgroundStatus(ctx context.Context, isBackground bool) error {
	return u.setAppBackgroundStatus(ctx, isBackground)
}
