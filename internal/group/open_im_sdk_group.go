package group

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/utils"
)

func (g *Group) SetGroupListener(callback open_im_sdk_callback.OnGroupListener) {
	if callback == nil {
		return
	}
	g.listener = callback
}

func (g *Group) CreateGroup(callback open_im_sdk_callback.Base, groupBaseInfo string, memberList string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, groupBaseInfo, memberList)
		var unmarshalCreateGroupBaseInfoParam sdk_params_callback.CreateGroupBaseInfoParam
		common.JsonUnmarshalAndArgsValidate(groupBaseInfo, &unmarshalCreateGroupBaseInfoParam, callback, operationID)
		var unmarshalCreateGroupMemberRoleParam sdk_params_callback.CreateGroupMemberRoleParam
		common.JsonUnmarshalAndArgsValidate(memberList, &unmarshalCreateGroupMemberRoleParam, callback, operationID)
		result := g.createGroup(callback, unmarshalCreateGroupBaseInfoParam, unmarshalCreateGroupMemberRoleParam, operationID)
		callback.OnSuccess(utils.StructToJsonString(result))
		log.NewInfo(operationID, fName, "callback: ", utils.StructToJsonString(result))
	}()
}

func (g *Group) JoinGroup(callback open_im_sdk_callback.Base, groupID, reqMsg string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupID, reqMsg)
		g.joinGroup(groupID, reqMsg, callback, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.JoinGroupCallback))
		log.NewInfo(operationID, fName, "callback: ", utils.StructToJsonString(sdk_params_callback.JoinGroupCallback))
	}()
}

func (g *Group) QuitGroup(callback open_im_sdk_callback.Base, groupID string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupID)
		g.quitGroup(groupID, callback, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.QuitGroupCallback))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonString(sdk_params_callback.QuitGroupCallback))
	}()
}

func (g *Group) GetJoinedGroupList(callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ")
		groupList := g.getJoinedGroupList(callback, operationID)
		log.Debug(operationID, "this is a dbd test", groupList)
		callback.OnSuccess(utils.StructToJsonStringDefault(groupList))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonStringDefault(groupList))
	}()
}

func (g *Group) GetGroupsInfo(callback open_im_sdk_callback.Base, groupIDList string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupIDList)
		var unmarshalGetGroupsInfoParam sdk_params_callback.GetGroupsInfoParam
		common.JsonUnmarshalAndArgsValidate(groupIDList, &unmarshalGetGroupsInfoParam, callback, operationID)
		groupsInfoList := g.getGroupsInfo(unmarshalGetGroupsInfoParam, callback, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(groupsInfoList))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonStringDefault(groupsInfoList))

	}()
}

func (g *Group) SetGroupInfo(callback open_im_sdk_callback.Base, groupInfo string, groupID string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupInfo, groupID)
		var unmarshalSetGroupInfoParam sdk_params_callback.SetGroupInfoParam
		common.JsonUnmarshalAndArgsValidate(groupInfo, &unmarshalSetGroupInfoParam, callback, operationID)
		g.setGroupInfo(callback, unmarshalSetGroupInfoParam, groupID, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.SetGroupInfoCallback))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonString(sdk_params_callback.SetGroupInfoCallback))
	}()
}

func (g *Group) GetGroupMemberList(callback open_im_sdk_callback.Base, groupID string, filter, offset, count int32, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupID, offset, count)
		groupMemberList := g.getGroupMemberList(callback, groupID, filter, offset, count, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(groupMemberList))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonStringDefault(groupMemberList))
	}()
}

func (g *Group) GetGroupMembersInfo(callback open_im_sdk_callback.Base, groupID string, userIDList string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupID, userIDList)
		var unmarshalParam sdk_params_callback.GetGroupMembersInfoParam
		common.JsonUnmarshalCallback(userIDList, &unmarshalParam, callback, operationID)
		groupMemberList := g.getGroupMembersInfo(callback, groupID, unmarshalParam, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(groupMemberList))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonStringDefault(groupMemberList))
	}()
}

func (g *Group) KickGroupMember(callback open_im_sdk_callback.Base, groupID string, reason string, userIDList string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupID, reason, userIDList)
		var unmarshalParam sdk_params_callback.KickGroupMemberParam
		common.JsonUnmarshalCallback(userIDList, &unmarshalParam, callback, operationID)
		result := g.kickGroupMember(callback, groupID, unmarshalParam, reason, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonStringDefault(result))
	}()
}

func (g *Group) TransferGroupOwner(callback open_im_sdk_callback.Base, groupID, newOwnerUserID string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupID, newOwnerUserID)
		g.transferGroupOwner(callback, groupID, newOwnerUserID, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.TransferGroupOwnerCallback))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonStringDefault(sdk_params_callback.TransferGroupOwnerCallback))
	}()
}

func (g *Group) InviteUserToGroup(callback open_im_sdk_callback.Base, groupID, reason string, userIDList string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupID, reason, userIDList)
		var unmarshalParam sdk_params_callback.InviteUserToGroupParam
		common.JsonUnmarshalAndArgsValidate(userIDList, &unmarshalParam, callback, operationID)
		result := g.inviteUserToGroup(callback, groupID, reason, unmarshalParam, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, fName, "callback: ", utils.StructToJsonStringDefault(result))
	}()
}

func (g *Group) GetRecvGroupApplicationList(callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ")
		result := g.getRecvGroupApplicationList(callback, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, fName, "callback: ", utils.StructToJsonString(result))
	}()
}

func (g *Group) GetSendGroupApplicationList(callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, "output results")
		result := g.getSendGroupApplicationList(callback, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, fName, "callback: ", utils.StructToJsonString(result))
	}()
}

func (g *Group) AcceptGroupApplication(callback open_im_sdk_callback.Base, groupID, fromUserID, handleMsg string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupID, fromUserID, handleMsg)
		g.processGroupApplication(callback, groupID, fromUserID, handleMsg, 1, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.AcceptGroupApplicationCallback))
		log.NewInfo(operationID, fName, "callback: ", utils.StructToJsonString(sdk_params_callback.AcceptGroupApplicationCallback))
	}()
}

func (g *Group) RefuseGroupApplication(callback open_im_sdk_callback.Base, groupID, fromUserID, handleMsg string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupID, fromUserID, handleMsg)
		g.processGroupApplication(callback, groupID, fromUserID, handleMsg, -1, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.RefuseGroupApplicationCallback))
		log.NewInfo(operationID, fName, "callback: ", utils.StructToJsonString(sdk_params_callback.RefuseGroupApplicationCallback))
	}()
}
