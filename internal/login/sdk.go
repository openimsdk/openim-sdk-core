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

func (u *LoginMgr) UploadImage(ctx context.Context, filePath string, token, obj string) (string, error) {
	return u.uploadImage(ctx, filePath, token, obj)
}

func (u *LoginMgr) UploadFile(ctx context.Context, filePath string) (string, error) {
	return u.uploadFile(ctx, filePath)
}
