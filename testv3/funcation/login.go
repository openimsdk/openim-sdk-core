package funcation

import (
	"open_im_sdk/internal/login"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"sync"
)

func LoginOne(idx int, uid string) bool {
	// get token
	collectToken(idx, uid)
	// init and login
	return initAndLogin(idx, uid, AllLoginMgr[idx].Token)
}

// 批量登录
// 返回值：成功登录和失败登录的 uidList
func LoginBatch(uidList []string) ([]string, []string) {
	var successList, failList []string
	var wg sync.WaitGroup
	wg.Add(len(uidList))
	for i, uid := range uidList {
		uid := uid
		i := i
		go func(idx int) {
			if LoginOne(idx, uid) == true {
				successList[i] = uid
			} else {
				failList[i] = uid
			}
		}(i)
	}
	wg.Wait()
	return successList, failList
}

func collectToken(idx int, uid string) {
	token, _ := getToken(uid)
	coreMgrLock.Lock()
	defer coreMgrLock.Unlock()
	AllLoginMgr[idx] = &CoreNode{Token: token, UserID: uid}
}

func initAndLogin(idx int, uid, token string) bool {
	var testinit testInitLister

	lg := new(login.LoginMgr)

	lg.InitSDK(cf, &testinit)
	log.Info(uid, "new login ", lg)
	AllLoginMgr[idx].mgr = lg
	log.Info(uid, "InitSDK ", cf, "index mgr", idx, lg)

	lg.SetConversationListener(&testConversation)

	var testUser userCallback
	lg.SetUserListener(testUser)

	var msgCallBack MsgListenerCallBak
	lg.SetAdvancedMsgListener(&msgCallBack)

	var friendListener testFriendListener
	lg.SetFriendListener(friendListener)

	var groupListener testGroupListener
	lg.SetGroupListener(groupListener)

	var callback BaseSuccessFailed
	callback.funcName = utils.GetSelfFuncName()

	operationID := utils.OperationIDGenerator()

	//ctx := mcontext.NewCtx(operationID)
	ctx := ccontext.WithOperationID(lg.Context(), operationID)

	err := lg.Login(ctx, uid, token)
	//lg.User().GetSelfUserInfo(ctx)
	if err != nil {
		log.Error(uid, err)
		return false
	}
	return true
}
