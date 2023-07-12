package funcation

import (
	"open_im_sdk/internal/login"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
)

func LoginOne(uid string) bool {
	// get token
	collectToken(uid)
	// init and login
	return initAndLogin(uid, AllLoginMgr[uid].Token)
}

// 批量登录
// 返回值：成功登录和失败登录的 uidList
func LoginBatch(uidList []string) ([]string, []string) {
	var successList, failList []string
	for i, uid := range uidList {
		if LoginOne(uid) == true {
			successList[i] = uid
		} else {
			failList[i] = uid
		}
	}
	return successList, failList
}

func collectToken(uid string) {
	token, _ := getToken(uid)
	coreMgrLock.Lock()
	defer coreMgrLock.Unlock()
	AllLoginMgr[uid] = &CoreNode{Token: token, UserID: uid}
}

func initAndLogin(uid, token string) bool {
	var testinit testInitLister

	lg := new(login.LoginMgr)

	lg.InitSDK(Config, &testinit)
	log.Info(uid, "new login ", lg)
	AllLoginMgr[uid].Mgr = lg
	log.Info(uid, "InitSDK ", Config, "index mgr", uid, lg)

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

	// ctx := mcontext.NewCtx(operationID)
	ctx := ccontext.WithOperationID(lg.Context(), operationID)

	err := lg.Login(ctx, uid, token)
	lg.User().GetSelfUserInfo(ctx)
	if err != nil {
		log.Error(uid, err)
		return false
	}
	return true
}
