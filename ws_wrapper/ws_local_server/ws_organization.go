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

type OrganizationCallback struct {
	uid string
}

func (g *OrganizationCallback) OnOrganizationUpdated() {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", "", "0"}, g.uid)
}

func (wsRouter *WsFuncRouter) SetOrganizationListener() {
	var g OrganizationCallback
	g.uid = wsRouter.uId
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.SetOrganizationListener(&g)
}

func (wsRouter *WsFuncRouter) GetSubDepartment(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}

	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "departmentID", "offset", "count") {
		return
	}
	//callback open_im_sdk_callback.Base, departmentID string, offset, count int, operationID string
	userWorker.Organization().GetSubDepartment(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["departmentID"].(string), int(m["offset"].(float64)), int(m["count"].(float64)),
		operationID)
}

func (wsRouter *WsFuncRouter) GetDepartmentMember(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}

	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "departmentID", "offset", "count") {
		return
	}
	//callback open_im_sdk_callback.Base, departmentID string, offset, count int, operationID string
	userWorker.Organization().GetDepartmentMember(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["departmentID"].(string), int(m["offset"].(float64)), int(m["count"].(float64)),
		operationID)
}

func (wsRouter *WsFuncRouter) GetUserInDepartment(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//callback open_im_sdk_callback.Base, departmentID string, offset, count int, operationID string
	userWorker.Organization().GetUserInDepartment(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		input, operationID)
}

func (wsRouter *WsFuncRouter) GetDepartmentMemberAndSubDepartment(input, operationID string) {

	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//callback open_im_sdk_callback.Base, departmentID string, offset, count int, operationID string
	userWorker.Organization().GetDepartmentMemberAndSubDepartment(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		input, operationID)
}

func (wsRouter *WsFuncRouter) GetParentDepartmentList(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}

	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "departmentID") {
		return
	}
	//callback open_im_sdk_callback.Base, departmentID string, offset, count int, operationID string
	userWorker.Organization().GetParentDepartmentList(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["departmentID"].(string), operationID)
}

func (wsRouter *WsFuncRouter) GetDepartmentInfo(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}

	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "departmentID") {
		return
	}
	userWorker.Organization().GetDepartmentInfo(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["departmentID"].(string), operationID)
}

func (wsRouter *WsFuncRouter) SearchOrganization(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}

	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "input","offset","count") {
		return
	}
	userWorker.Organization().SearchOrganization(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["input"].(string), int(m["offset"].(float64)), int(m["count"].(float64)), operationID)
}
