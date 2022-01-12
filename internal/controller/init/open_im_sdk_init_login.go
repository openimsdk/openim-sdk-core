package init

import (
	"open_im_sdk/pkg/common"
)



func (u *LoginMgr) Login(callback common.Base, userID, token string) {
	go func() {
		u.login(userID, token, callback)
	}()
}


func (u *LoginMgr) Logout(callback common.Base){
	go func(){
		u.logout(callback)
	}()
}





//func InitOnce(config *utils.IMConfig) bool {
//	constant.SvrConf = *config
//	initUserRouter()
//	open_im_sdk.initAddr()
//	return true
//}

