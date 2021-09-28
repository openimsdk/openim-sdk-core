package ws_local_server

import "open_im_sdk/open_im_sdk"

func (wsRouter *WsFuncRouter) GetUsersInfo(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.GetUsersInfo(input, &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) SetSelfInfo(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.SetSelfInfo(input, &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}
