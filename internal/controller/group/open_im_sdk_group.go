package group

import (
	"encoding/json"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/utils"
)

type OnGroupListener interface {
	OnMemberEnter(groupId string, memberList string)
	OnMemberLeave(groupId string, member string)
	OnMemberInvited(groupId string, opUser string, memberList string)
	OnMemberKicked(groupId string, opUser string, memberList string)
	OnGroupCreated(groupId string)
	OnGroupInfoChanged(groupId string, groupInfo string)
	OnReceiveJoinApplication(groupId string, member string, opReason string)
	OnApplicationProcessed(groupId string, opUser string, AgreeOrReject int32, opReason string)
}

func (u *Group) SetGroupListener(callback OnGroupListener) {
	if callback == nil {
		utils.sdkLog("callback null")
		return
	}
	u.listener = callback
	utils.sdkLog("SetGroupListener ", callback)
}

func (u *Group) CreateGroup(callback common.Base, groupBaseInfo string, memberList string, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "CreateGroup args: ", groupBaseInfo, memberList)
		var unmarshalCreateGroupBaseInfoParam sdk_params_callback.CreateGroupBaseInfoParam
		common.JsonUnmarshalAndArgsValidate(groupBaseInfo, &unmarshalCreateGroupBaseInfoParam, callback, operationID)
		var unmarshalCreateGroupMemberRoleParam sdk_params_callback.CreateGroupMemberRoleParam
		common.JsonUnmarshalAndArgsValidate(memberList, &unmarshalCreateGroupMemberRoleParam, callback, operationID)
		result := u.createGroup(unmarshalCreateGroupBaseInfoParam, unmarshalCreateGroupMemberRoleParam, operationID)
		callback.OnSuccess(utils.StructToJsonString(result))
		log.NewInfo(operationID, "CreateGroup callback: ", utils.StructToJsonString(result))
	}()
}

func (u *Group) JoinGroup(callback common.Base, groupID, reqMsg string, operationID string) {
	if callback == nil {
		return
	}
	go func() {

		log.NewInfo(operationID, utils.RunFuncName(), "args: ", groupID, reqMsg)
		err := u.joinGroup(groupID, reqMsg, callback, operationID)
		result := u.createGroup(groupID, reqMsg, callback, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.JoinGroupCallback))
		callback.OnSuccess(utils.StructToJsonString(result))
		log.NewInfo(operationID, "CreateGroup callback: ", utils.StructToJsonString(result))
		err := u.joinGroup(groupID, message, callback, operationID)
		if err != nil {
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		callback.OnSuccess(constant.DeFaultSuccessMsg)
	}()
}

func (u *Group) QuitGroup(groupId string, callback open_im_sdk.Base, operationID string) {
	if callback == nil {
		utils.sdkLog("callback is nil")
		return
	}
	go func() {
		utils.sdkLog("............quitGroup begin...............")
		err := u.quitGroup(groupId)
		if err != nil {
			utils.sdkLog(".........quitGroup failed.............", groupId, err.Error())
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		utils.sdkLog("..........quitGroup end callback...........", groupId)
		callback.OnSuccess(constant.DeFaultSuccessMsg)
	}()
}

func (u *Group) GetJoinedGroupList(callback open_im_sdk.Base, operationID string) {
	if callback == nil {
		utils.sdkLog("Base callback is nil")
		return
	}
	go func() {
		groupInfoList, err := u.getJoinedGroupListFromLocal()
		if err != nil {
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		utils.sdkLog("getJoinedGroupListFromLocal, ", groupInfoList)

		for i, v := range groupInfoList {
			g, err := u._getGroupInfoByGroupID(v.GroupId)
			if err != nil {
				utils.sdkLog("findLocalGroupOwnerByGroupId failed,  ", err.Error(), v.GroupId)
				continue
			}
			utils.sdkLog("findLocalGroupOwnerByGroupId ", v.GroupId, g.OwnerUserID)
			v.OwnerId = g.OwnerUserID
			utils.sdkLog("getLocalGroupMemberNumByGroupId ", v.GroupId, g.MemberCount)
			v.MemberCount = uint32(g.MemberCount)
			groupInfoList[i] = v
		}

		jsonGroupInfoList, err := json.Marshal(groupInfoList)
		if err != nil {
			utils.sdkLog("marshal failed, ", err.Error())
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		utils.sdkLog("jsonGroupInfoList, ", string(jsonGroupInfoList))
		callback.OnSuccess(string(jsonGroupInfoList))
	}()
}

func (u *Group) GetGroupsInfo(groupIdList string, callback open_im_sdk.Base, operationID string) {
	if callback == nil {
		utils.sdkLog("Base callback is nil")
		return
	}
	go func() {
		var sctgroupIdList []string
		err := json.Unmarshal([]byte(groupIdList), &sctgroupIdList)
		if err != nil {
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}

		groupsInfoList, err := u.getGroupsInfo(sctgroupIdList)
		if err != nil {
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		jsonList, err := json.Marshal(groupsInfoList)
		if err != nil {
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		callback.OnSuccess(string(jsonList))
	}()

}

func (u *Group) SetGroupInfo(jsonGroupInfo string, callback open_im_sdk.Base, operationID string) {
	if callback == nil {
		utils.sdkLog("callback is nil")
		return
	}
	go func() {
		utils.sdkLog("............SetGroupInfo begin...................")
		var newGroupInfo open_im_sdk.setGroupInfoReq
		err := json.Unmarshal([]byte(jsonGroupInfo), &newGroupInfo)
		if err != nil {
			utils.sdkLog("............Unmarshal failed................", err.Error())
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		err = u.setGroupInfo(newGroupInfo)
		if err != nil {
			utils.sdkLog("..........setGroupInfo failed........... ", err.Error())
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		utils.sdkLog(".........setGroupInfo end, callback...............", jsonGroupInfo)
		callback.OnSuccess(constant.DeFaultSuccessMsg)
	}()
}

func (u *Group) GetGroupMemberList(groupId string, filter int32, next int32, callback open_im_sdk.Base, operationID string) {
	if callback == nil {
		utils.sdkLog("Base callback is nil")
		return
	}
	go func() {
		n, groupMemberResult, err := u.getGroupMemberListFromLocal(groupId, filter, next)
		if err != nil {
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		utils.sdkLog("getGroupMemberListFromLocal ", groupId, filter, next, groupMemberResult)

		var result open_im_sdk.getGroupMemberListResult
		if groupMemberResult == nil {
			groupMemberResult = make([]open_im_sdk.groupMemberFullInfo, 0)
		}
		result.Data = groupMemberResult
		result.NextSeq = n
		jsonGroupMemberResult, err := json.Marshal(result)
		if err != nil {
			utils.sdkLog("marshal failed, ", err.Error())
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		utils.sdkLog("jsonGroupMemberResult ", string(jsonGroupMemberResult))
		callback.OnSuccess(string(jsonGroupMemberResult))
	}()
}

func (u *Group) GetGroupMembersInfo(groupId string, userList string, callback open_im_sdk.Base, operationID string) {
	if callback == nil {
		utils.sdkLog("Base callback is nil")
		return
	}
	go func() {
		var sctmemberList []string
		err := json.Unmarshal([]byte(userList), &sctmemberList)
		if err != nil {
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		utils.sdkLog("GetGroupMembersInfo args: ", groupId, userList)
		r, err := u.getGroupMembersInfoFromLocal(groupId, sctmemberList)
		if err != nil {
			utils.sdkLog("getGroupMembersInfoFromLocal failed", groupId, err.Error())
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		utils.sdkLog("getGroupMembersInfoFromLocal, ", groupId, sctmemberList, r)

		jsonResult, err := json.Marshal(r)
		if err != nil {
			utils.sdkLog("marshal failed, ", err.Error())
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		utils.sdkLog("jsonResult", string(jsonResult))
		callback.OnSuccess(string(jsonResult))
	}()
}

func (u *Group) KickGroupMember(groupId string, reason string, userList string, callback open_im_sdk.Base, operationID string) {
	if callback == nil {
		utils.sdkLog("callback null")
		return
	}
	go func() {
		var sctMemberList []string
		err := json.Unmarshal([]byte(userList), &sctMemberList)
		if err != nil {
			utils.sdkLog("unmarshal failed, ", err.Error(), userList)
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		r, err := u.kickGroupMember(groupId, sctMemberList, reason)
		if err != nil {
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		utils.sdkLog("kickGroupMember, ", groupId, sctMemberList, reason)

		jsonResult, err := json.Marshal(r)
		if err != nil {
			utils.sdkLog("marshal failed, ", err.Error())
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		utils.sdkLog("KickGroupMember, req: ", groupId, userList, reason, "resp: ", string(jsonResult))
		callback.OnSuccess(string(jsonResult))
	}()
}

func (u *Group) TransferGroupOwner(groupId, userId string, callback open_im_sdk.Base, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		err := u.transferGroupOwner(groupId, userId)
		if err != nil {
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		u.syncGroupRequest()
		u.syncGroupMemberByGroupId(groupId)
		callback.OnSuccess(constant.DeFaultSuccessMsg)
	}()
}

func (u *Group) InviteUserToGroup(groupId, reason string, userList string, callback open_im_sdk.Base, operationID string) {
	if callback == nil {
		utils.sdkLog("callbak null")
		return
	}
	go func() {
		var sctUserList []string
		err := json.Unmarshal([]byte(userList), &sctUserList)
		if err != nil {
			utils.sdkLog("unmarshal failed, ", err.Error())
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}

		result, err := u.inviteUserToGroup(groupId, reason, sctUserList)
		if err != nil {
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		utils.sdkLog("inviteUserToGroup, ", groupId, reason, sctUserList)

		jsonResult, err := json.Marshal(result)
		if err != nil {
			utils.sdkLog("marshal failed, ", err.Error())
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		utils.sdkLog("InviteUserToGroup, req: ", groupId, reason, userList, "resp: ", string(jsonResult))

		callback.OnSuccess(string(jsonResult))
	}()

}

func (u *Group) GetGroupApplicationList(callback open_im_sdk.Base, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		r, err := u.getGroupApplicationList()
		if err != nil {
			utils.sdkLog("getGroupApplicationList faild, ", err.Error())
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		jsonResult, err := json.Marshal(r)
		if err != nil {
			utils.sdkLog("getGroupApplicationList faild, ", err.Error())
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		callback.OnSuccess(string(jsonResult))
		return
	}()
}

/*
func (u *UserRelated) TsetGetGroupApplicationList(callback Base) string {
	if callback == nil {
		return ""
	}

	r, err := u.getGroupApplicationList()
	if err != nil {
		sdkLog("getGroupApplicationList faild, ", err.Error())
		callback.OnError(ErrCodeGroup, err.Error())
		return ""
	}
	jsonResult, err := json.Marshal(r)
	if err != nil {
		sdkLog("getGroupApplicationList faild, ", err.Error())
		callback.OnError(ErrCodeGroup, err.Error())
		return ""
	}
	callback.OnSuccess(string(jsonResult))
	return string(jsonResult)

}*/

func (u *Group) AcceptGroupApplication(application, reason string, callback open_im_sdk.Base, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		var sctApplication utils.GroupReqListInfo
		err := json.Unmarshal([]byte(application), &sctApplication)
		if err != nil {
			utils.sdkLog("Unmarshal, ", err.Error())
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}

		var access open_im_sdk.accessOrRefuseGroupApplicationReq
		access.OperationID = utils.operationIDGenerator()
		access.GroupId = sctApplication.GroupID
		access.FromUser = sctApplication.FromUserID
		access.FromUserNickName = sctApplication.FromUserNickname
		access.FromUserFaceUrl = sctApplication.FromUserFaceUrl
		access.ToUser = sctApplication.ToUserID
		access.ToUserNickname = sctApplication.ToUserNickname
		access.ToUserFaceUrl = sctApplication.ToUserFaceUrl
		access.AddTime = sctApplication.AddTime
		access.RequestMsg = sctApplication.RequestMsg
		access.HandledMsg = reason
		access.Type = sctApplication.Type
		access.HandleStatus = 2
		access.HandleResult = 1

		err = u.acceptGroupApplication(&access)
		if err != nil {
			utils.sdkLog("acceptGroupApplication, ", err.Error())
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		u.syncGroupRequest()
		u.syncGroupMemberByGroupId(sctApplication.GroupID)
		callback.OnSuccess(constant.DeFaultSuccessMsg)
	}()
}

func (u *Group) RefuseGroupApplication(application, reason string, callback open_im_sdk.Base, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		var sctApplication utils.GroupReqListInfo
		err := json.Unmarshal([]byte(application), &sctApplication)
		if err != nil {
			utils.sdkLog("Unmarshal, ", err.Error())
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}

		var access open_im_sdk.accessOrRefuseGroupApplicationReq
		access.OperationID = utils.operationIDGenerator()
		access.GroupId = sctApplication.GroupID
		access.FromUser = sctApplication.FromUserID
		access.FromUserNickName = sctApplication.FromUserNickname
		access.FromUserFaceUrl = sctApplication.FromUserFaceUrl
		access.ToUser = sctApplication.ToUserID
		access.ToUserNickname = sctApplication.ToUserNickname
		access.ToUserFaceUrl = sctApplication.ToUserFaceUrl
		access.AddTime = sctApplication.AddTime
		access.RequestMsg = sctApplication.RequestMsg
		access.HandledMsg = reason
		access.Type = sctApplication.Type
		access.HandleStatus = 2
		access.HandleResult = 0

		err = u.refuseGroupApplication(&access)
		if err != nil {
			utils.sdkLog("refuseGroupApplication, ", err.Error())
			callback.OnError(constant.ErrCodeGroup, err.Error())
			return
		}
		u.syncGroupRequest()
		callback.OnSuccess(constant.DeFaultSuccessMsg)
	}()

}
