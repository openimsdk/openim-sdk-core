package login

import (
	"open_im_sdk/pkg/common"
	"sync"
)

func (u *LoginMgr) Login(callback common.Base, userID, token string, operationID string) {
	go func() {
		u.login(userID, token, callback, operationID)
	}()
}

func (u *LoginMgr) Logout(callback common.Base, operationID string) {
	go func() {
		u.logout(callback, operationID)
	}()
}

func (u *LoginMgr) UploadImage(callback common.Base, filePath string, token, obj string, operationID string) string {
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

//func InitOnce(config *utils.IMConfig) bool {
//	constant.SvrConf = *config
//	initUserRouter()
//	open_im_sdk.initAddr()
//	return true
//}
