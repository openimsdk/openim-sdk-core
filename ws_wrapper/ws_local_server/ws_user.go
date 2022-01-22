package ws_local_server

import (
	"open_im_sdk/open_im_sdk"
)

func (wsRouter *WsFuncRouter) GetUsersInfo(userIDList string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.User().GetUsersInfo(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, userIDList, operationID)
}

func (wsRouter *WsFuncRouter) SetSelfInfo(userInfo string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.User().SetSelfInfo(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, userInfo, operationID)
}

func (wsRouter *WsFuncRouter) GetSelfUserInfo(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.User().GetSelfUserInfo(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}
