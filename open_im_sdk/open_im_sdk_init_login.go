package open_im_sdk

type IMSDKListener interface {
	OnConnecting()
	OnConnectSuccess()
	OnConnectFailed(ErrCode int, ErrMsg string)
	OnKickedOffline()
	OnUserTokenExpired()
	OnSelfInfoUpdated(userInfo string)
}

func InitOnce(config *IMConfig) bool {
	SvrConf = *config
	initUserRouter()
	initAddr()
	sdkLog("InitOnce success, ", *config)
	return true
}

func GetUserWorker(uid string) *UserRelated {
	UserSDKRwLock.Lock()
	defer UserSDKRwLock.Unlock()
	v, ok := UserRouterMap[uid]
	if ok {
		return v
	}
	UserRouterMap[uid] = new(UserRelated)

	return UserRouterMap[uid]
}

func (u *UserRelated) InitSDK(config string, cb IMSDKListener) bool {
	return u.initSDK(config, cb)
}

func (u *UserRelated) UnInitSDK() {
	u.unInitSDK()
}

func (u *UserRelated) Login(uid, tk string, callback Base) {
	if callback == nil {
		sdkLog("callback is null")
		return
	}
	//	go func() {
	u.login(uid, tk, callback)
	//	}()
}

func (u *UserRelated) ForceReConn() {
	if u.conn != nil {
		u.conn.Close()
	}
}

func (u *UserRelated) Logout(callback Base) {
	u.logout(callback)
}

func (u *UserRelated) GetLoginStatus() int {
	return u.getLoginStatus()
}

func (u *UserRelated) GetLoginUser() string {
	return u.LoginUid
}

func (u *UserRelated) ForceSyncLoginUserInfo() {
	u.syncLoginUserInfo()
}

type Base interface {
	OnError(errCode int32, errMsg string)
	OnSuccess(data string)
}

func initUserRouter() {
	UserSDKRwLock.Lock()
	defer UserSDKRwLock.Unlock()
	UserRouterMap = make(map[string]*UserRelated, 0)
}
