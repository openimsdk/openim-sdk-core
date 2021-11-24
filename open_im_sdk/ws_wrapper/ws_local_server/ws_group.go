package ws_local_server

import (
	"encoding/json"
	"open_im_sdk/open_im_sdk"
)

type GroupCallback struct {
	uid string
}

func (g *GroupCallback) OnMemberEnter(groupId string, memberList string) {
	m := make(map[string]interface{}, 2)
	m["groupId"] = groupId
	m["memberList"] = memberList
	j, _ := json.Marshal(m)
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", string(j), "0"}, g.uid)
}
func (g *GroupCallback) OnMemberLeave(groupId string, memberList string) {
	m := make(map[string]interface{}, 2)
	m["groupId"] = groupId
	m["memberList"] = memberList
	j, _ := json.Marshal(m)
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", string(j), "0"}, g.uid)
}
func (g *GroupCallback) OnMemberInvited(groupId string, opUser string, memberList string) {
	m := make(map[string]interface{}, 3)
	m["groupId"] = groupId
	m["opUser"] = opUser
	m["memberList"] = memberList
	j, _ := json.Marshal(m)
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", string(j), "0"}, g.uid)
}
func (g *GroupCallback) OnMemberKicked(groupId string, opUser string, memberList string) {
	m := make(map[string]interface{}, 3)
	m["groupId"] = groupId
	m["opUser"] = opUser
	m["memberList"] = memberList

	j, _ := json.Marshal(m)

	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", string(j), "0"}, g.uid)
}
func (g *GroupCallback) OnGroupCreated(groupId string) {
	m := make(map[string]interface{}, 1)
	m["groupId"] = groupId
	j, _ := json.Marshal(m)
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", string(j), "0"}, g.uid)
}
func (g *GroupCallback) OnGroupInfoChanged(groupId string, groupInfo string) {
	m := make(map[string]interface{}, 2)
	m["groupId"] = groupId
	m["groupInfo"] = groupInfo
	j, _ := json.Marshal(m)
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", string(j), "0"}, g.uid)
}
func (g *GroupCallback) OnReceiveJoinApplication(groupId string, member string, opReason string) {
	m := make(map[string]interface{}, 3)
	m["groupId"] = groupId
	m["member"] = member
	m["opReason"] = opReason
	j, _ := json.Marshal(m)
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", string(j), "0"}, g.uid)

}
func (g *GroupCallback) OnApplicationProcessed(groupId string, opUser string, AgreeOrReject int32, opReason string) {
	m := make(map[string]interface{}, 4)
	m["groupId"] = groupId
	m["opUser"] = opUser
	m["AgreeOrReject"] = AgreeOrReject
	m["opReason"] = opReason
	j, _ := json.Marshal(m)
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", string(j), "0"}, g.uid)
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
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "gInfo", "memberList") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.CreateGroup(m["gInfo"].(string), m["memberList"].(string), &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) JoinGroup(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "groupId", "message") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.JoinGroup(m["groupId"].(string), m["message"].(string), &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) QuitGroup(groupId, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.QuitGroup(groupId, &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) GetJoinedGroupList(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.GetJoinedGroupList(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) GetGroupsInfo(input, operationID string) { //(groupIdList string, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "groupIdList") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.GetGroupsInfo(m["groupIdList"].(string), &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) SetGroupInfo(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "groupInfo") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.SetGroupInfo(m["groupInfo"].(string), &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) GetGroupMemberList(input, operationID string) { //(groupId string, filter int32, next int32, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "groupId", "filter", "next") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.GetGroupMemberList(m["groupId"].(string), int32(m["filter"].(float64)), int32(m["next"].(float64)), &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) GetGroupMembersInfo(input, operationID string) { //(groupId string, userList string, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "groupId", "userList") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.GetGroupMembersInfo(m["groupId"].(string), m["userList"].(string), &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) KickGroupMember(input, operationID string) { //(groupId string, reason string, userList string, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "groupId", "reason", "userList") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.KickGroupMember(m["groupId"].(string), m["reason"].(string), m["userList"].(string), &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) TransferGroupOwner(input, operationID string) { //(groupId, userId string, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "groupId", "userId") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.TransferGroupOwner(m["groupId"].(string), m["userId"].(string), &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) InviteUserToGroup(input, operationID string) { //(groupId, reason string, userList string, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "groupId", "reason", "userList") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.InviteUserToGroup(m["groupId"].(string), m["reason"].(string), m["userList"].(string), &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) GetGroupApplicationList(input, operationID string) { //(callback Base) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.GetGroupApplicationList(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})

}

func (wsRouter *WsFuncRouter) AcceptGroupApplication(input, operationID string) { //(application, reason string, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "application", "reason") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.AcceptGroupApplication(m["application"].(string), m["reason"].(string), &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})

}

func (wsRouter *WsFuncRouter) RefuseGroupApplication(input, operationID string) { //(application, reason string, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "application", "reason") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.RefuseGroupApplication(m["application"].(string), m["reason"].(string), &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}
