package ws_local_server

import (
	"open_im_sdk/internal/controller/init"
)

func (wsRouter *WsFuncRouter) GetUsersInfo(input string, operationID string) {
	userWorker := init.GetUserWorker(wsRouter.uId)
	userWorker.GetUsersInfo(input, &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) SetSelfInfo(input string, operationID string) {
	userWorker := init.GetUserWorker(wsRouter.uId)
	userWorker.SetSelfInfo(input, &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}
