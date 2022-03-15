package ws_local_server

import (
	"open_im_sdk/open_im_sdk"
)

func (wsRouter *WsFuncRouter) SignalingInvite(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Signaling().Invite(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, input, operationID)
}

func (wsRouter *WsFuncRouter) SignalingInviteInGroup(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Signaling().InviteInGroup(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, input, operationID)
}

func (wsRouter *WsFuncRouter) SignalingAccept(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Signaling().Accept(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, input, operationID)
}

func (wsRouter *WsFuncRouter) SignalingReject(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Signaling().Reject(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, input, operationID)
}

func (wsRouter *WsFuncRouter) SignalingCancel(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Signaling().Cancel(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, input, operationID)
}

func (wsRouter *WsFuncRouter) SignalingHungUp(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Signaling().HungUp(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, input, operationID)
}
