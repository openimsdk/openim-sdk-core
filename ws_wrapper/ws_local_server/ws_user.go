package ws_local_server

import (
	"open_im_sdk/internal/open_im_sdk"
)

func (wsRouter *WsFuncRouter) GetUsersInfo(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.User().GetUsersInfo(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, input, operationID)
}

func (wsRouter *WsFuncRouter) SetSelfInfo(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.User().SetSelfInfo(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, input, operationID)
}
