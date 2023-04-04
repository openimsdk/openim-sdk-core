package group

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/jinzhu/copier"
	"open_im_sdk/internal/util"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/sdk_params_callback"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
)

func (g *Group) SetGroupListener(callback open_im_sdk_callback.OnGroupListener) {
	if callback == nil {
		return
	}
	g.listener = callback
}

// deprecated use CreateGroupV2
func (g *Group) CreateGroup(ctx context.Context, groupBaseInfo sdk_params_callback.CreateGroupBaseInfoParam, memberList sdk_params_callback.CreateGroupMemberRoleParam) {
	req := &group.CreateGroupReq{
		GroupInfo: &sdkws.GroupInfo{
			GroupName:    groupBaseInfo.GroupName,
			Notification: groupBaseInfo.Notification,
			Introduction: groupBaseInfo.Introduction,
			FaceURL:      groupBaseInfo.FaceURL,
			Ex:           groupBaseInfo.Ex,
			GroupType:    groupBaseInfo.GroupType,
		},
	}
	if groupBaseInfo.NeedVerification != nil {
		req.GroupInfo.NeedVerification = *groupBaseInfo.NeedVerification
	}
	for _, info := range memberList {
		switch info.RoleLevel {
		case constant.GroupOrdinaryUsers:
		case constant.GroupOwner:
		case constant.GroupAdmin:

		}
	}

}

func (g *Group) CreateGroupV2(ctx context.Context, req *group.CreateGroupReq) (*sdkws.GroupInfo, error) {
	resp, err := util.ApiFunc[group.CreateGroupResp](ctx, constant.CreateGroupRouter, req)
	if err != nil {
		return nil, err
	}
	//g.SyncJoinedGroupList(operationID)
	//g.syncGroupMemberByGroupID(realData.GroupInfo.GroupID, operationID, false)
	return resp.GroupInfo, nil
}

func (g *Group) createGroupV2(callback open_im_sdk_callback.Base, group sdk.CreateGroupBaseInfoParam, memberList sdk.CreateGroupMemberRoleParam, operationID string) *sdk.CreateGroupCallback {
	apiReq := api.CreateGroupReq{}
	apiReq.OperationID = operationID
	apiReq.OwnerUserID = g.loginUserID
	apiReq.MemberList = memberList
	for _, v := range apiReq.MemberList {
		if v.RoleLevel == 0 {
			v.RoleLevel = 1
		}
	}

	copier.Copy(&apiReq, &group)
	realData := api.CreateGroupResp{}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "api req args: ", apiReq)
	g.p.PostFatalCallbackPenetrate(callback, constant.CreateGroupRouter, apiReq, &realData.GroupInfo, apiReq.OperationID)
	m := utils.JsonDataOne(&realData.GroupInfo)
	g.SyncJoinedGroupList(operationID)
	g.syncGroupMemberByGroupID(realData.GroupInfo.GroupID, operationID, false)
	return (*sdk.CreateGroupCallback)(&m)
}

func (g *Group) JoinGroup(callback open_im_sdk_callback.Base, groupID, reqMsg string, joinSource int32, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupID, reqMsg, joinSource)
		g.joinGroup(groupID, reqMsg, joinSource, callback, operationID)
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

func (g *Group) DismissGroup(callback open_im_sdk_callback.Base, groupID string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupID)
		g.dismissGroup(groupID, callback, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.DismissGroupCallback))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonString(sdk_params_callback.DismissGroupCallback))
	}()
}

func (g *Group) ChangeGroupMute(callback open_im_sdk_callback.Base, groupID string, isMute bool, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupID, isMute)
		g.changeGroupMute(groupID, isMute, callback, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.GroupMuteChangeCallback))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonString(sdk_params_callback.GroupMuteChangeCallback))
	}()
}

func (g *Group) ChangeGroupMemberMute(callback open_im_sdk_callback.Base, groupID, userID string, mutedSeconds uint32, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupID, userID, mutedSeconds)
		g.changeGroupMemberMute(groupID, userID, mutedSeconds, callback, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.GroupMemberMuteChangeCallback))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonString(sdk_params_callback.GroupMemberMuteChangeCallback))
	}()
}

func (g *Group) SetGroupMemberRoleLevel(callback open_im_sdk_callback.Base, groupID, userID string, roleLevel int, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupID, userID, roleLevel)
		g.setGroupMemberRoleLevel(callback, groupID, userID, roleLevel, operationID)
		callback.OnSuccess(constant.SuccessCallbackDefault)
		log.NewInfo(operationID, fName, " callback: ", constant.SuccessCallbackDefault)
	}()
}

func (g *Group) SetGroupMemberInfo(callback open_im_sdk_callback.Base, groupMemberInfo string, operationID string) {
	var unmarshalSetGroupMemberInfoParam sdk_params_callback.SetGroupMemberInfoParam
	common.JsonUnmarshalAndArgsValidate(groupMemberInfo, &unmarshalSetGroupMemberInfoParam, callback, operationID)
	g.setGroupMemberInfo(callback, unmarshalSetGroupMemberInfoParam, operationID)
	callback.OnSuccess(utils.StructToJsonStringDefault(sdk_params_callback.SetGroupMemberInfoCallback))
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

func (g *Group) SearchGroups(callback open_im_sdk_callback.Base, searchParam, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", searchParam)
		var unmarshalGetGroupsInfoParam sdk_params_callback.SearchGroupsParam
		common.JsonUnmarshalAndArgsValidate(searchParam, &unmarshalGetGroupsInfoParam, callback, operationID)
		unmarshalGetGroupsInfoParam.KeywordList = utils.TrimStringList(unmarshalGetGroupsInfoParam.KeywordList)
		groupsInfoList := g.searchGroups(callback, unmarshalGetGroupsInfoParam, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(groupsInfoList))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonStringDefault(groupsInfoList), len(groupsInfoList))

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

//func SetGroupVerification(callback open_im_sdk_callback.Base, operationID string, groupID string, verification int32)

func (g *Group) SetGroupVerification(callback open_im_sdk_callback.Base, verification int32, groupID string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", verification, groupID)
		var unmarshalSetGroupInfoParam sdk_params_callback.SetGroupInfoParam
		n := verification
		unmarshalSetGroupInfoParam.NeedVerification = &n
		g.setGroupInfo(callback, unmarshalSetGroupInfoParam, groupID, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.SetGroupInfoCallback))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonString(sdk_params_callback.SetGroupInfoCallback))
	}()
}
func (g *Group) SetGroupLookMemberInfo(callback open_im_sdk_callback.Base, rule int32, groupID string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", rule, groupID)
		apiReq := api.SetGroupInfoReq{}
		apiReq.GroupID = groupID
		apiReq.LookMemberInfo = &rule
		apiReq.OperationID = operationID
		g.modifyGroupInfo(callback, apiReq, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.SetGroupInfoCallback))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonString(sdk_params_callback.SetGroupInfoCallback))
	}()
}
func (g *Group) SetGroupApplyMemberFriend(callback open_im_sdk_callback.Base, rule int32, groupID string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", rule, groupID)
		apiReq := api.SetGroupInfoReq{}
		apiReq.GroupID = groupID
		apiReq.ApplyMemberFriend = &rule
		apiReq.OperationID = operationID
		g.modifyGroupInfo(callback, apiReq, operationID)
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
func (g *Group) GetGroupMemberOwnerAndAdmin(callback open_im_sdk_callback.Base, groupID string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupID)
		groupMemberList := g.getGroupMemberOwnerAndAdmin(callback, groupID, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(groupMemberList))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonStringDefault(groupMemberList))
	}()
}

// getGroupMemberListByJoinTimeFilter
func (g *Group) GetGroupMemberListByJoinTimeFilter(callback open_im_sdk_callback.Base, groupID string, offset, count int32, joinTimeBegin, joinTimeEnd int64, filterUserID, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupID, offset, count, filterUserID)
		var userIDList []string
		common.JsonUnmarshalAndArgsValidate(filterUserID, &userIDList, callback, operationID)
		groupMemberList := g.getGroupMemberListByJoinTimeFilter(callback, groupID, offset, count, joinTimeBegin, joinTimeEnd, userIDList, operationID)
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

func (g *Group) SetGroupMemberNickname(callback open_im_sdk_callback.Base, groupID, userID string, GroupMemberNickname string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupID, userID, GroupMemberNickname)
		g.setGroupMemberNickname(callback, groupID, userID, GroupMemberNickname, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.SetGroupMemberNicknameCallback))
		log.NewInfo(operationID, fName, "callback: ", utils.StructToJsonString(sdk_params_callback.SetGroupMemberNicknameCallback))
	}()
}

func (g *Group) SearchGroupMembers(callback open_im_sdk_callback.Base, searchParam string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", searchParam)
		var unmarshalSearchGroupMembersParam sdk_params_callback.SearchGroupMembersParam
		common.JsonUnmarshalAndArgsValidate(searchParam, &unmarshalSearchGroupMembersParam, callback, operationID)
		members := g.searchGroupMembers(callback, unmarshalSearchGroupMembersParam, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(members))
		log.NewInfo(operationID, fName, "callback: ", utils.StructToJsonStringDefault(members))
	}()
}
