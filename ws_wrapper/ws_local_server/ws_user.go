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

package ws_local_server

import (
	"open_im_sdk/open_im_sdk"
)

func (wsRouter *WsFuncRouter) GetUsersInfo(userIDList string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, userIDList, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Full().GetUsersInfo(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, userIDList, operationID)
}

func (wsRouter *WsFuncRouter) SetSelfInfo(userInfo string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, userInfo, operationID, runFuncName(), nil) {
		return
	}
	userWorker.User().SetSelfInfo(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, userInfo, operationID)
}

func (wsRouter *WsFuncRouter) GetSelfUserInfo(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.User().GetSelfUserInfo(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}

type UserCallback struct {
	uid string
}

func (u *UserCallback) OnSelfInfoUpdated(userInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", userInfo, "0"}, u.uid)
}

func (wsRouter *WsFuncRouter) SetUserListener() {
	var u UserCallback
	u.uid = wsRouter.uId
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.SetUserListener(&u)
}
