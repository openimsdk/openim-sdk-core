package ws_local_server

import (
	"encoding/json"
	"open_im_sdk/open_im_sdk"
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

func (g *GroupCallback) OnReceiveJoinGroupApplicationAdded(groupApplication string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", groupApplication, "0"}, g.uid)
}
func (g *GroupCallback) OnReceiveJoinGroupApplicationDeleted(groupApplication string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", groupApplication, "0"}, g.uid)
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

//
//
//func (g *GroupCallback) OnMemberEnter(groupId string, memberList string) {
//	m := make(map[string]interface{}, 2)
//	m["groupId"] = groupId
//	m["memberList"] = memberList
//	j, _ := json.Marshal(m)
//	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", string(j), "0"}, g.uid)
//}
//
//func (g *GroupCallback) OnMemberLeave(groupId string, memberList string) {
//	m := make(map[string]interface{}, 2)
//	m["groupId"] = groupId
//	m["memberList"] = memberList
//	j, _ := json.Marshal(m)
//	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", string(j), "0"}, g.uid)
//}
//func (g *GroupCallback) OnMemberInvited(groupId string, opUser string, memberList string) {
//	m := make(map[string]interface{}, 3)
//	m["groupId"] = groupId
//	m["opUser"] = opUser
//	m["memberList"] = memberList
//	j, _ := json.Marshal(m)
//	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", string(j), "0"}, g.uid)
//}
//func (g *GroupCallback) OnMemberKicked(groupId string, opUser string, memberList string) {
//	m := make(map[string]interface{}, 3)
//	m["groupId"] = groupId
//	m["opUser"] = opUser
//	m["memberList"] = memberList
//
//	j, _ := json.Marshal(m)
//
//	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", string(j), "0"}, g.uid)
//}
//func (g *GroupCallback) OnGroupCreated(groupId string) {
//	m := make(map[string]interface{}, 1)
//	m["groupId"] = groupId
//	j, _ := json.Marshal(m)
//	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", string(j), "0"}, g.uid)
//}
//func (g *GroupCallback) OnGroupInfoChanged(groupId string, groupInfo string) {
//	m := make(map[string]interface{}, 2)
//	m["groupId"] = groupId
//	m["groupInfo"] = groupInfo
//	j, _ := json.Marshal(m)
//	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", string(j), "0"}, g.uid)
//}
//func (g *GroupCallback) OnReceiveJoinApplication(groupId string, member string, opReason string) {
//	m := make(map[string]interface{}, 3)
//	m["groupId"] = groupId
//	m["member"] = member
//	m["opReason"] = opReason
//	j, _ := json.Marshal(m)
//	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", string(j), "0"}, g.uid)
//
//}
//func (g *GroupCallback) OnApplicationProcessed(groupId string, opUser string, AgreeOrReject int32, opReason string) {
//	m := make(map[string]interface{}, 4)
//	m["groupId"] = groupId
//	m["opUser"] = opUser
//	m["AgreeOrReject"] = AgreeOrReject
//	m["opReason"] = opReason
//	j, _ := json.Marshal(m)
//	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", string(j), "0"}, g.uid)
//}
//
//func (g *GroupCallback) OnGroupApplicationAccepted(groupApplication string){
//
//}
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
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "groupBaseInfo", "memberList") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//callback common.Base, groupBaseInfo string, memberList string, operationID string
	userWorker.Group().CreateGroup(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupBaseInfo"].(string), m["memberList"].(string), operationID)
}

func (wsRouter *WsFuncRouter) JoinGroup(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "groupID", "reqMsg") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//callback common.Base, groupID, reqMsg string, operationID string
	userWorker.Group().JoinGroup(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupID"].(string), m["reqMsg"].(string), operationID)
}

func (wsRouter *WsFuncRouter) QuitGroup(groupID, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//callback common.Base, groupID string, operationID string
	userWorker.Group().QuitGroup(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, groupID, operationID)
}

func (wsRouter *WsFuncRouter) GetJoinedGroupList(input, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//(callback common.Base, operationID string)
	userWorker.Group().GetJoinedGroupList(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}

func (wsRouter *WsFuncRouter) GetGroupsInfo(input, operationID string) { //(groupIdList string, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "groupIDList") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//callback common.Base, groupIDList string, operationID string
	userWorker.Group().GetGroupsInfo(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupIDList"].(string), operationID)
}

func (wsRouter *WsFuncRouter) SetGroupInfo(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "groupInfo", "groupID") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//(callback common.Base, groupInfo string, groupID string, operationID string)
	userWorker.Group().SetGroupInfo(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupInfo"].(string), m["groupID"].(string), operationID)
}

func (wsRouter *WsFuncRouter) GetGroupMemberList(input, operationID string) { //(groupId string, filter int32, next int32, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "groupID", "filter", "next") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//callback common.Base, groupID string, filter int32, next int32, operationID string
	userWorker.Group().GetGroupMemberList(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupID"].(string), int32(m["filter"].(float64)), int32(m["next"].(float64)), operationID)
}

func (wsRouter *WsFuncRouter) GetGroupMembersInfo(input, operationID string) { //(groupId string, userList string, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "groupID", "userIDList") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//callback common.Base, groupID string, userIDList string, operationID string
	userWorker.Group().GetGroupMembersInfo(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupID"].(string), m["userIDList"].(string), operationID)
}

func (wsRouter *WsFuncRouter) KickGroupMember(input, operationID string) { //(groupId string, reason string, userList string, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "groupID", "reason", "userIDList") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//KickGroupMember(callback common.Base, groupID string, reason string, userIDList string, operationID string)
	userWorker.Group().KickGroupMember(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupID"].(string), m["reason"].(string), m["userIDList"].(string), operationID)
}

func (wsRouter *WsFuncRouter) TransferGroupOwner(input, operationID string) { //(groupId, userId string, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "groupID", "newOwnerUserID") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	//callback common.Base, groupID, newOwnerUserID string, operationID string
	userWorker.Group().TransferGroupOwner(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupID"].(string), m["newOwnerUserID"].(string), operationID)
}

func (wsRouter *WsFuncRouter) InviteUserToGroup(input, operationID string) { //(groupId, reason string, userList string, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "groupID", "reason", "userIDList") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Group().InviteUserToGroup(
		&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupID"].(string), m["reason"].(string), m["userIDList"].(string), operationID)
}

func (wsRouter *WsFuncRouter) GetGroupApplicationList(input, operationID string) { //(callback Base) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Group().GetGroupApplicationList(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, operationID)

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
	userWorker.Group().AcceptGroupApplication(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupID"].(string), m["fromUserID"].(string), m["handleMsg"].(string), operationID)
}

func (wsRouter *WsFuncRouter) RefuseGroupApplication(input, operationID string) { //(application, reason string, callback Base) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "groupID", "fromUserID", "handleMsg") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Group().RefuseGroupApplication(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupID"].(string), m["fromUserID"].(string), m["handleMsg"].(string), operationID)
}
