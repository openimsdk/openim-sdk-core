package ws_local_server

import "open_im_sdk/open_im_sdk"

type WorkMomentsCallback struct {
	uid string
}

func (wsRouter *WsFuncRouter) SetWorkMomentsListener() {
	var cb WorkMomentsCallback
	cb.uid = wsRouter.uId
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.SetWorkMomentsListener(&cb)
}

func (w *WorkMomentsCallback) OnRecvNewNotification() {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", "", "0"}, w.uid)
}

func (wsRouter *WsFuncRouter) GetWorkMomentsUnReadCount(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.WorkMoments().GetWorkMomentsUnReadCount(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}

func (wsRouter *WsFuncRouter) GetWorkMomentsNotification(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.WorkMoments().GetWorkMomentsNotification(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}

func (wsRouter *WsFuncRouter) ClearWorkMomentsNotification(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.WorkMoments().ClearWorkMomentsNotification(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}
