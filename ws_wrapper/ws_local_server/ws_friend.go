package ws_local_server

import (
	"open_im_sdk/open_im_sdk"
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
	userWorker.Friend().GetDesignatedFriendsInfo(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, userIDList, operationID)
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
func (wsRouter *WsFuncRouter) CheckFriend(userIDList string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Friend().CheckFriend(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, userIDList, operationID)
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
func (wsRouter *WsFuncRouter) SetFriendRemark(remark string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Friend().SetFriendRemark(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, remark, operationID)
}

func (wsRouter *WsFuncRouter) AddBlack(blackUserID string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Friend().AddBlack(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, blackUserID, operationID)
}

func (wsRouter *WsFuncRouter) GetBlackList(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Friend().GetBlackList(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}

func (wsRouter *WsFuncRouter) RemoveBlack(removeUserID string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Friend().RemoveBlack(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, removeUserID, operationID)
}
