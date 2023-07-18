// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	for _, uid := range uidList {
		if LoginOne(uid) == true {
			successList = append(successList, uid)
		} else {
			failList = append(failList, uid)
		}
	}
	return successList, failList
}

func collectToken(uid string) {
	token, _ := GetToken(uid)
	coreMgrLock.Lock()
	defer coreMgrLock.Unlock()
	AllLoginMgr[uid] = &CoreNode{Token: token, UserID: uid}
}

func initAndLogin(uid, token string) bool {
	var init initLister

	lg := new(login.LoginMgr)

	lg.InitSDK(Config, &init)
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
