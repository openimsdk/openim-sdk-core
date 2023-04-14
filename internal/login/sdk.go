package login

import (
	"context"
	"open_im_sdk/open_im_sdk_callback"
	"sync"
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

func (u *LoginMgr) UploadImage(callback open_im_sdk_callback.Base, filePath string, token, obj string, operationID string) string {
	var wg sync.WaitGroup
	wg.Add(1)
	url := ""
	go func() {
		url = u.uploadImage(callback, filePath, token, obj, operationID)
		wg.Done()
	}()

	wg.Wait()
	return url
}

func (u *LoginMgr) UploadFile(ctx context.Context, filePath string) error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		u.uploadFile(callback, filePath, operationID)
		wg.Done()
	}()
	wg.Wait()
}
