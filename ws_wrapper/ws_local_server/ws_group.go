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

type GroupCallback struct {
	uid string
}

func (g *GroupCallback) OnJoinedGroupAdded(groupInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", groupInfo, "0"}, g.uid)
}
func (g *GroupCallback) OnJoinedGroupDeleted(groupInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", groupInfo, "0"}, g.uid)
}

func (g *GroupCallback) OnGroupMemberAdded(groupMemberInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", groupMemberInfo, "0"}, g.uid)
}
func (g *GroupCallback) OnGroupMemberDeleted(groupMemberInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", groupMemberInfo, "0"}, g.uid)
}

func (g *GroupCallback) OnGroupApplicationAdded(groupApplication string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", groupApplication, "0"}, g.uid)
}
func (g *GroupCallback) OnGroupApplicationDeleted(groupApplication string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", groupApplication, "0"}, g.uid)
}

func (g *GroupCallback) OnGroupInfoChanged(groupInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", groupInfo, "0"}, g.uid)
}
func (g *GroupCallback) OnGroupMemberInfoChanged(groupMemberInfo string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", groupMemberInfo, "0"}, g.uid)
}

func (g *GroupCallback) OnGroupApplicationAccepted(groupApplication string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", groupApplication, "0"}, g.uid)
}
func (g *GroupCallback) OnGroupApplicationRejected(groupApplication string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", groupApplication, "0"}, g.uid)
}

func (wsRouter *WsFuncRouter) SetGroupListener() {
	var g GroupCallback
	g.uid = wsRouter.uId
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.SetGroupListener(&g)
}

func (wsRouter *WsFuncRouter) CreateGroup(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}

	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupBaseInfo", "memberList") {
		return
	}
	//callback common.Base, groupBaseInfo string, memberList string, operationID string
	userWorker.Group().CreateGroup(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupBaseInfo"].(string), m["memberList"].(string), operationID)
}

func (wsRouter *WsFuncRouter) JoinGroup(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}

	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupID", "reqMsg", "joinSource") {
		return
	}
	//callback common.Base, groupID, reqMsg string, operationID string
	userWorker.Group().JoinGroup(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupID"].(string), m["reqMsg"].(string), int32(m["joinSource"].(float64)), operationID)
}

func (wsRouter *WsFuncRouter) QuitGroup(groupID, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//callback common.Base, groupID string, operationID string
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, groupID, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Group().QuitGroup(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, groupID, operationID)
}

func (wsRouter *WsFuncRouter) DismissGroup(groupID, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//callback common.Base, groupID string, operationID string
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, groupID, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Group().DismissGroup(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, groupID, operationID)
}

func (wsRouter *WsFuncRouter) ChangeGroupMute(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//callback common.Base, groupID string, operationID string
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupID", "isMute") {
		return
	}
	userWorker.Group().ChangeGroupMute(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["groupID"].(string), m["isMute"].(bool), operationID)
}

func (wsRouter *WsFuncRouter) SetGroupMemberNickname(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//callback common.Base, groupID string, operationID string
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupID", "userID", "GroupMemberNickname") {
		return
	}
	userWorker.Group().SetGroupMemberNickname(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["groupID"].(string), m["userID"].(string), m["GroupMemberNickname"].(string), operationID)
}

func (wsRouter *WsFuncRouter) SetGroupMemberRoleLevel(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//callback common.Base, groupID string, operationID string
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupID", "userID", "roleLevel") {
		return
	}
	userWorker.Group().SetGroupMemberRoleLevel(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["groupID"].(string), m["userID"].(string), int(m["roleLevel"].(float64)), operationID)
}

func (wsRouter *WsFuncRouter) ChangeGroupMemberMute(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//callback common.Base, groupID string, operationID string
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupID", "userID", "mutedSeconds") {
		return
	}
	userWorker.Group().ChangeGroupMemberMute(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["groupID"].(string), m["userID"].(string), uint32(m["mutedSeconds"].(float64)), operationID)
}

func (wsRouter *WsFuncRouter) GetJoinedGroupList(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//(callback common.Base, operationID string)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Group().GetJoinedGroupList(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}

func (wsRouter *WsFuncRouter) GetGroupsInfo(input, operationID string) { //(groupIdList string, callback Base) {
	//m := make(map[string]interface{})
	//if err := json.Unmarshal([]byte(input), &m); err != nil {
	//	log.Info("unmarshal failed")
	//	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
	//	return
	//}
	//if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "groupIDList") {
	//	return
	//}

	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	//callback common.Base, groupIDList string, operationID string
	userWorker.Group().GetGroupsInfo(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		input, operationID)
}
func (wsRouter *WsFuncRouter) SearchGroups(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Group().SearchGroups(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		input, operationID)
}

func (wsRouter *WsFuncRouter) SetGroupInfo(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}

	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupInfo", "groupID") {
		return
	}
	//(callback common.Base, groupInfo string, groupID string, operationID string)
	userWorker.Group().SetGroupInfo(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupInfo"].(string), m["groupID"].(string), operationID)
}

func (wsRouter *WsFuncRouter) SetGroupVerification(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}

	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "verification", "groupID") {
		return
	}
	//(callback common.Base, groupInfo string, groupID string, operationID string)
	userWorker.Group().SetGroupVerification(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		int32(m["verification"].(float64)), m["groupID"].(string), operationID)
}

func (wsRouter *WsFuncRouter) GetGroupMemberList(input, operationID string) { //(groupId string, filter int32, next int32, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupID", "filter", "offset", "count") {
		return
	}
	//callback common.Base, groupID string, filter int32, next int32, operationID string
	userWorker.Group().GetGroupMemberList(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupID"].(string), int32(m["filter"].(float64)), int32(m["offset"].(float64)), int32(m["count"].(float64)), operationID)
}
func (wsRouter *WsFuncRouter) GetGroupMemberOwnerAndAdmin(input, operationID string) { //(groupId string, filter int32, next int32, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupID") {
		return
	}
	//callback common.Base, groupID string, filter int32, next int32, operationID string
	userWorker.Group().GetGroupMemberOwnerAndAdmin(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupID"].(string), operationID)
}
func (wsRouter *WsFuncRouter) GetGroupMemberListByJoinTimeFilter(input, operationID string) { //(groupId string, filter int32, next int32, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupID", "offset", "count", "joinTimeBegin", "joinTimeEnd", "filterUserIDList") {
		return
	}
	userWorker.Group().GetGroupMemberListByJoinTimeFilter(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupID"].(string), int32(m["offset"].(float64)), int32(m["count"].(float64)), int64(m["joinTimeBegin"].(float64)), int64(m["joinTimeEnd"].(float64)), m["filterUserIDList"].(string), operationID)
}

func (wsRouter *WsFuncRouter) GetGroupMembersInfo(input, operationID string) { //(groupId string, userList string, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupID", "userIDList") {
		return
	}
	//callback common.Base, groupID string, userIDList string, operationID string
	userWorker.Group().GetGroupMembersInfo(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupID"].(string), m["userIDList"].(string), operationID)
}

func (wsRouter *WsFuncRouter) KickGroupMember(input, operationID string) { //(groupId string, reason string, userList string, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupID", "reason", "userIDList") {
		return
	}
	//KickGroupMember(callback common.Base, groupID string, reason string, userIDList string, operationID string)
	userWorker.Group().KickGroupMember(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupID"].(string), m["reason"].(string), m["userIDList"].(string), operationID)
}

func (wsRouter *WsFuncRouter) TransferGroupOwner(input, operationID string) { //(groupId, userId string, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupID", "newOwnerUserID") {
		return
	}
	//callback common.Base, groupID, newOwnerUserID string, operationID string
	userWorker.Group().TransferGroupOwner(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupID"].(string), m["newOwnerUserID"].(string), operationID)
}

func (wsRouter *WsFuncRouter) InviteUserToGroup(input, operationID string) { //(groupId, reason string, userList string, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupID", "reason", "userIDList") {
		return
	}

	userWorker.Group().InviteUserToGroup(
		&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupID"].(string), m["reason"].(string), m["userIDList"].(string), operationID)
}

func (wsRouter *WsFuncRouter) GetRecvGroupApplicationList(input, operationID string) { //(callback Base) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Group().GetRecvGroupApplicationList(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}

func (wsRouter *WsFuncRouter) GetSendGroupApplicationList(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Group().GetSendGroupApplicationList(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}

func (wsRouter *WsFuncRouter) AcceptGroupApplication(input, operationID string) { //(application, reason string, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupID", "fromUserID", "handleMsg") {
		return
	}

	userWorker.Group().AcceptGroupApplication(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupID"].(string), m["fromUserID"].(string), m["handleMsg"].(string), operationID)
}

func (wsRouter *WsFuncRouter) RefuseGroupApplication(input, operationID string) { //(application, reason string, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}

	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupID", "fromUserID", "handleMsg") {
		return
	}
	userWorker.Group().RefuseGroupApplication(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupID"].(string), m["fromUserID"].(string), m["handleMsg"].(string), operationID)
}

func (wsRouter *WsFuncRouter) SearchGroupMembers(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "searchParam") {
		return
	}
	userWorker.Group().SearchGroupMembers(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["searchParam"].(string), operationID)
}

//SetGroupApplyMemberFriend
func (wsRouter *WsFuncRouter) SetGroupApplyMemberFriend(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "rule", "groupID") {
		return
	}
	userWorker.Group().SetGroupApplyMemberFriend(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		int32(m["rule"].(float64)), m["groupID"].(string), operationID)
}

//SetGroupApplyMemberFriend
func (wsRouter *WsFuncRouter) SetGroupLookMemberInfo(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "rule", "groupID") {
		return
	}
	userWorker.Group().SetGroupLookMemberInfo(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		int32(m["rule"].(float64)), m["groupID"].(string), operationID)
}
