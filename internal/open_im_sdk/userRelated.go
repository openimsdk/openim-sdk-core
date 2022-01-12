package open_im_sdk

import (

	"open_im_sdk/internal/controller/init"
	"sync"
)



func init() {
	UserSDKRwLock.Lock()
	defer UserSDKRwLock.Unlock()
	UserRouterMap = make(map[string]*init.LoginMgr, 0)
}


var UserSDKRwLock sync.RWMutex
var UserRouterMap map[string]*init.LoginMgr

var userForSDK *init.LoginMgr

func GetUserWorker(uid string) *init.LoginMgr {
	UserSDKRwLock.Lock()
	defer UserSDKRwLock.Unlock()
	v, ok := UserRouterMap[uid]
	if ok {
		return v
	}
	UserRouterMap[uid] = new(init.LoginMgr)

	return UserRouterMap[uid]
}


