package open_im_sdk

import (

	"open_im_sdk/internal/controller/init"
	"sync"
)



//func (u *UserRelated) initListenerCh() {
//	u.ch = make(chan utils.cmd2Value, 1000)
//	u.ConversationCh = u.ch
//
//	u.wsNotification = make(map[string]chan utils.GeneralWsResp, 1)
//	u.seqMsg = make(map[int32]*server_api_params.MsgData, 1000)
//
//	u.receiveMessageOpt = make(map[string]int32, 1000)
//}

func initUserRouter() {
	UserSDKRwLock.Lock()
	defer UserSDKRwLock.Unlock()
	UserRouterMap = make(map[string]*init.LoginMgr, 0)
}


var UserSDKRwLock sync.RWMutex
var UserRouterMap map[string]*init.LoginMgr

var userForSDK *init.LoginMgr

