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
	"encoding/json"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
)

type FriendCallback struct {
	uid string
}

func (f *FriendCallback) OnFriendApplicationAdded(applyUserInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", applyUserInfo, "0"}, f.uid)
}
func (f *FriendCallback) OnFriendApplicationDeleted(applyUserInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", applyUserInfo, "0"}, f.uid)
}
func (f *FriendCallback) OnFriendApplicationAccepted(applyUserInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", applyUserInfo, "0"}, f.uid)
}
func (f *FriendCallback) OnFriendApplicationRejected(applyUserInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", applyUserInfo, "0"}, f.uid)
}
func (f *FriendCallback) OnFriendAdded(friendInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", friendInfo, "0"}, f.uid)
}
func (f *FriendCallback) OnFriendDeleted(friendInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", friendInfo, "0"}, f.uid)
}
func (f *FriendCallback) OnFriendInfoChanged(userInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", userInfo, "0"}, f.uid)
}
func (f *FriendCallback) OnBlackAdded(userInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", userInfo, "0"}, f.uid)
}
func (f *FriendCallback) OnBlackDeleted(friendInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", friendInfo, "0"}, f.uid)
}

func (wsRouter *WsFuncRouter) SetFriendListener() {
	var fr FriendCallback
	fr.uid = wsRouter.uId
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.SetFriendListener(&fr)
}

//1
func (wsRouter *WsFuncRouter) GetDesignatedFriendsInfo(userIDList string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, userIDList, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Friend().GetDesignatedFriendsInfo(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, userIDList, operationID)
}

//1
func (wsRouter *WsFuncRouter) AddFriend(paramsReq string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, paramsReq, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Friend().AddFriend(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, paramsReq, operationID)
}

func (wsRouter *WsFuncRouter) GetRecvFriendApplicationList(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Friend().GetRecvFriendApplicationList(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}

func (wsRouter *WsFuncRouter) GetSendFriendApplicationList(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Friend().GetSendFriendApplicationList(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}

//1
func (wsRouter *WsFuncRouter) AcceptFriendApplication(params string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, params, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Friend().AcceptFriendApplication(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, params, operationID)
}

//1
func (wsRouter *WsFuncRouter) RefuseFriendApplication(params string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, params, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Friend().RefuseFriendApplication(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, params, operationID)
}

//1
func (wsRouter *WsFuncRouter) CheckFriend(userIDList string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, userIDList, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Friend().CheckFriend(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, userIDList, operationID)
}

//1
func (wsRouter *WsFuncRouter) DeleteFriend(friendUserID string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, friendUserID, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Friend().DeleteFriend(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, friendUserID, operationID)
}

//1
func (wsRouter *WsFuncRouter) GetFriendList(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Friend().GetFriendList(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}
func (wsRouter *WsFuncRouter) SearchFriends(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Friend().SearchFriends(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		input, operationID)
}

//1
func (wsRouter *WsFuncRouter) SetFriendRemark(remark string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, remark, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Friend().SetFriendRemark(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, remark, operationID)
}

func (wsRouter *WsFuncRouter) AddBlack(blackUserID string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, blackUserID, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Friend().AddBlack(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, blackUserID, operationID)
}

func (wsRouter *WsFuncRouter) GetBlackList(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Friend().GetBlackList(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}

func (wsRouter *WsFuncRouter) RemoveBlack(removeUserID string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, removeUserID, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Friend().RemoveBlack(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, removeUserID, operationID)
}
