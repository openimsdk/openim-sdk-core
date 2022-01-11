package init

import (
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/utils"
)

//type IMSDKListener interface {
//	OnConnecting()
//	OnConnectSuccess()
//	OnConnectFailed(ErrCode int32, ErrMsg string)
//	OnKickedOffline()
//	OnUserTokenExpired()
//	OnSelfInfoUpdated(userInfo string)
//}

func InitOnce(config *utils.IMConfig) bool {
	constant.SvrConf = *config
	initUserRouter()
	open_im_sdk.initAddr()
	utils.sdkLog("InitOnce success, ", *config)
	return true
}

func GetUserWorker(uid string) *constant.UserRelated {
	constant.UserSDKRwLock.Lock()
	defer constant.UserSDKRwLock.Unlock()
	v, ok := constant.UserRouterMap[uid]
	if ok {
		return v
	}
	constant.UserRouterMap[uid] = new(constant.UserRelated)

	return constant.UserRouterMap[uid]
}

func (u *constant.open_im_sdk) InitSDK(config string, cb IMSDKListener) bool {
	return u.initSDK(config, cb)
}

func (u *constant.open_im_sdk) UnInitSDK() {
	u.unInitSDK()
}

func (u *constant.open_im_sdk) Login(uid, tk string, callback Base) {
	if callback == nil {
		utils.sdkLog("callback is null")
		return
	}
	//	go func() {
	u.login(uid, tk, callback)
	//	}()
}

func (u *constant.open_im_sdk) ForceReConn() {
	if u.conn != nil {
		u.conn.Close()
	}
}

func (u *constant.open_im_sdk) Logout(callback Base) {
	u.logout(callback)
}

func (u *constant.open_im_sdk) GetLoginStatus() int {
	return u.getLoginStatus()
}

func (u *constant.open_im_sdk) GetLoginUser() string {
	return u.loginUserID
}

func (u *constant.open_im_sdk) ForceSyncLoginUserInfo() {
	u.syncLoginUserInfo()
}

func initUserRouter() {
	constant.UserSDKRwLock.Lock()
	defer constant.UserSDKRwLock.Unlock()
	constant.UserRouterMap = make(map[string]*constant.UserRelated, 0)
}
