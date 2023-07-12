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
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.WorkMoments().GetWorkMomentsNotification(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, int(m["offset"].(float64)), int(m["count"].(float64)), operationID)
}

func (wsRouter *WsFuncRouter) ClearWorkMomentsNotification(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.WorkMoments().ClearWorkMomentsNotification(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}
