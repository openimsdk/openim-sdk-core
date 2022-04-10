package open_im_sdk

import (
	"open_im_sdk/internal/login"
	"sync"
)

func init() {
	UserSDKRwLock.Lock()
	defer UserSDKRwLock.Unlock()
	UserRouterMap = make(map[string]*login.LoginMgr, 0)
}

var UserSDKRwLock sync.RWMutex
var UserRouterMap map[string]*login.LoginMgr

var userForSDK *login.LoginMgr

func GetUserWorker(uid string) *login.LoginMgr {
	UserSDKRwLock.Lock()
	defer UserSDKRwLock.Unlock()
	v, ok := UserRouterMap[uid]
	if ok {
		return v
	}
	UserRouterMap[uid] = new(login.LoginMgr)

	return UserRouterMap[uid]
}
