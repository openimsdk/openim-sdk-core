package ws_local_server

import (
	"open_im_sdk/open_im_sdk"
)

type FriendCallback struct {
	uid string
}

func (f *FriendCallback) OnFriendApplicationListAdded(applyUserInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", applyUserInfo, "0"}, f.uid)
}
func (f *FriendCallback) OnFriendApplicationListDeleted(applyUserInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", applyUserInfo, "0"}, f.uid)
}
func (f *FriendCallback) OnFriendApplicationListAccept(applyUserInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", applyUserInfo, "0"}, f.uid)
}
func (f *FriendCallback) OnFriendApplicationListReject(applyUserInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", applyUserInfo, "0"}, f.uid)
}
func (f *FriendCallback) OnFriendListAdded(friendInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", friendInfo, "0"}, f.uid)
}
func (f *FriendCallback) OnFriendListDeleted(friendInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", friendInfo, "0"}, f.uid)
}
func (f *FriendCallback) OnBlackListAdd(userInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", userInfo, "0"}, f.uid)
}
func (f *FriendCallback) OnBlackListDeleted(userInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", userInfo, "0"}, f.uid)
}
func (f *FriendCallback) OnFriendInfoChanged(friendInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", friendInfo, "0"}, f.uid)
}

func (wsRouter *WsFuncRouter) SetFriendListener() bool {
	var fr FriendCallback
	fr.uid = wsRouter.uId
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	return userWorker.SetFriendListener(&fr)
}

//1
func (wsRouter *WsFuncRouter) GetFriendsInfo(uidList string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.GetFriendsInfo(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, uidList)
}

//1
func (wsRouter *WsFuncRouter) AddFriend(paramsReq string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.AddFriend(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, paramsReq)
}

//1
func (wsRouter *WsFuncRouter) GetFriendApplicationList(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.GetFriendApplicationList(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

//1
func (wsRouter *WsFuncRouter) AcceptFriendApplication(uid string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.AcceptFriendApplication(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, uid)
}

//1
func (wsRouter *WsFuncRouter) RefuseFriendApplication(uid string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.RefuseFriendApplication(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, uid)
}

//1
func (wsRouter *WsFuncRouter) CheckFriend(uidList string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.CheckFriend(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, uidList)
}

//1
func (wsRouter *WsFuncRouter) DeleteFromFriendList(deleteUid string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.DeleteFromFriendList(deleteUid, &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

//1
func (wsRouter *WsFuncRouter) GetFriendList(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.GetFriendList(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

//1
func (wsRouter *WsFuncRouter) SetFriendInfo(comment string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.SetFriendInfo(comment, &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) AddToBlackList(blackUid string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.AddToBlackList(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, blackUid)
}

func (wsRouter *WsFuncRouter) GetBlackList(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.GetBlackList(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) DeleteFromBlackList(deleteUid string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.DeleteFromBlackList(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, deleteUid)
}
