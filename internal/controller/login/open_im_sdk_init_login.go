package login

import (
	"open_im_sdk/pkg/common"
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

//func InitOnce(config *utils.IMConfig) bool {
//	constant.SvrConf = *config
//	initUserRouter()
//	open_im_sdk.initAddr()
//	return true
//}
