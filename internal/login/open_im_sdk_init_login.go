package login

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/log"
	"sync"
)

func (u *LoginMgr) Login(callback open_im_sdk_callback.Base, userID, token string, operationID string) {
	go func() {
		u.login(userID, token, callback, operationID)
	}()
}

func (u *LoginMgr) WakeUp(callback open_im_sdk_callback.Base, operationID string) {
	go func() {
		u.wakeUp(callback, operationID)
	}()
}

func (u *LoginMgr) Logout(callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		u.logout(callback, operationID)
		return
	}
	go func() {
		u.logout(callback, operationID)
	}()
}

func (u *LoginMgr) SetAppBackgroundStatus(callback open_im_sdk_callback.Base, isBackground bool, operationID string) {
	go func() {
		log.NewInfo(operationID, "SetAppBackgroundStatus", "isBackground", isBackground)
		u.setAppBackgroundStatus(callback, isBackground, operationID)
	}()
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

func (u *LoginMgr) UploadFile(callback open_im_sdk_callback.SendMsgCallBack, filePath, operationID string) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		u.uploadFile(callback, filePath, operationID)
		wg.Done()
	}()
	wg.Wait()
}
