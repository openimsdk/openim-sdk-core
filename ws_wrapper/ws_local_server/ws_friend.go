package ws_local_server

import (
	"open_im_sdk/internal/open_im_sdk"
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
	return userWorker.Friend().SetFriendListener(&fr)
}

//1
func (wsRouter *WsFuncRouter) GetDesignatedFriendsInfo(uidList string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Friend().GetDesignatedFriendsInfo(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, uidList, operationID)
}

//1
func (wsRouter *WsFuncRouter) AddFriend(paramsReq string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Friend().AddFriend(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, paramsReq, operationID)
}

func (wsRouter *WsFuncRouter) GetRecvFriendApplicationList(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Friend().GetRecvFriendApplicationList(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}

func (wsRouter *WsFuncRouter) GetSendFriendApplicationList(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Friend().GetSendFriendApplicationList(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}

//1
func (wsRouter *WsFuncRouter) AcceptFriendApplication(params string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Friend().AcceptFriendApplication(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, params, operationID)
}

//1
func (wsRouter *WsFuncRouter) RefuseFriendApplication(params string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Friend().RefuseFriendApplication(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, params, operationID)
}

//1
func (wsRouter *WsFuncRouter) CheckFriend(uidList string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Friend().CheckFriend(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, uidList, operationID)
}

//1
func (wsRouter *WsFuncRouter) DeleteFriend(friendUserID string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Friend().DeleteFriend(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, friendUserID, operationID)
}

//1
func (wsRouter *WsFuncRouter) GetFriendList(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Friend().GetFriendList(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}

//1
func (wsRouter *WsFuncRouter) SetFriendRemark(comment string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Friend().SetFriendRemark(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, comment, operationID)
}

func (wsRouter *WsFuncRouter) AddBlack(blackUid string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Friend().AddBlack(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, blackUid, operationID)
}

func (wsRouter *WsFuncRouter) GetBlackList(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Friend().GetBlackList(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}

func (wsRouter *WsFuncRouter) RemoveBlack(removeUserID string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Friend().RemoveBlack(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, removeUserID, operationID)
}
