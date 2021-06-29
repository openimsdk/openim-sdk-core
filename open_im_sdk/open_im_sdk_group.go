package open_im_sdk

import "encoding/json"

type OnGroupListener interface {
	//list->group
	OnMemberEnter(groupId string, memberList string)
	//group->one
	OnMemberLeave(groupId string, member string)
	//list->opUser->groupId
	OnMemberInvited(groupId string, opUser string, memberList string)
	//
	OnMemberKicked(groupId string, opUser string, memberList string)
	OnGroupCreated(groupId string)
	OnGroupInfoChanged(groupId string, groupInfo string)
	OnReceiveJoinApplication(groupId string, member string, opReason string)
	OnApplicationProcessed(groupId string, opUser string, AgreeOrReject int32, opReason string)

	OnTransferGroupOwner(groupId string, oldUserID string, newUserID string)
}

func SetGroupListener(callback OnGroupListener) {
	groupManager.listener = callback
}

func CreateGroup(gInfo string, memberList string, callback Base) {
	go func() {
		//to struct
		var sctGroupInfo groupInfo
		err := json.Unmarshal([]byte(gInfo), &sctGroupInfo)
		if err != nil {
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}

		var sctmemberList []createGroupMemberInfo
		err = json.Unmarshal([]byte(memberList), &sctmemberList)
		if err != nil {
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}

		err = groupManager.createGroup(sctGroupInfo, sctmemberList)
		if err != nil {
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		callback.OnSuccess(DeFaultSuccessMsg)
	}()
}

func JoinGroup(groupId, message string, callback Base) {
	var ui2ClientReq joinGroupReq
	if err := json.Unmarshal([]byte(groupId), &ui2ClientReq.GroupID); err != nil {
		callback.OnError(ErrCodeGroup, err.Error())
		return
	}

	if err := json.Unmarshal([]byte(message), &ui2ClientReq.Message); err != nil {
		callback.OnError(ErrCodeGroup, err.Error())
		return
	}

	err := groupManager.joinGroup(ui2ClientReq.GroupID, ui2ClientReq.Message)
	if err != nil {
		callback.OnError(ErrCodeGroup, err.Error())
		return
	}
	callback.OnSuccess(DeFaultSuccessMsg)
}

func QuitGroup(groupId string, callback Base) {
	var ui2ClientQuitGroupReq quitGroupReq
	if err := json.Unmarshal([]byte(groupId), &ui2ClientQuitGroupReq.GroupID); err != nil {
		callback.OnError(ErrCodeGroup, err.Error())
		return
	}

	err := groupManager.quitGroup(ui2ClientQuitGroupReq.GroupID)
	if err != nil {
		callback.OnError(ErrCodeGroup, err.Error())
		return
	}
	callback.OnSuccess(DeFaultSuccessMsg)
}

func GetJoinedGroupList(callback Base) {
	go func() {
		groupInfoList, err := groupManager.getJoinedGroupList()
		if err != nil {
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}

		jsonGroupInfoList, err := json.Marshal(groupInfoList)
		if err != nil {
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		callback.OnSuccess(string(jsonGroupInfoList))
	}()
}

func GetGroupsInfo(groupIdList string, callback Base) {
	var sctgroupIdList []string
	err := json.Unmarshal([]byte(groupIdList), &sctgroupIdList)
	if err != nil {
		callback.OnError(ErrCodeGroup, err.Error())
		return
	}

	groupsInfoList, err := groupManager.getGroupsInfo(sctgroupIdList)
	if err != nil {
		callback.OnError(ErrCodeGroup, err.Error())
		return
	}
	jsonList, err := json.Marshal(groupsInfoList)
	if err != nil {
		callback.OnError(ErrCodeGroup, err.Error())
		return
	}
	callback.OnSuccess(string(jsonList))
}

func SetGroupInfo(jsonGroupInfo string, callback Base) {

	var newGroupInfo groupInfo
	err := json.Unmarshal([]byte(jsonGroupInfo), &newGroupInfo)
	if err != nil {
		callback.OnError(ErrCodeGroup, err.Error())
		return
	}

	err = groupManager.setGroupInfo(newGroupInfo)
	if err != nil {
		callback.OnError(ErrCodeGroup, err.Error())
		return
	}
	callback.OnSuccess(DeFaultSuccessMsg)
}

func GetGroupMemberList(groupId string, filter int32, next int32, callback Base) {
	go func() {
		n, groupMemberResult, err := groupManager.getGroupMemberList(groupId, filter, next)
		if err != nil {
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		var result getGroupMemberListResult
		result.Data = groupMemberResult
		result.NextSeq = n
		jsonGroupMemberResult, err := json.Marshal(result)
		if err != nil {
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		callback.OnSuccess(string(jsonGroupMemberResult))
	}()
}

func GetGroupMembersInfo(groupId string, userList string, callback Base) {
	go func() {
		var sctmemberList []string
		err := json.Unmarshal([]byte(userList), &sctmemberList)
		if err != nil {
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		r, err := groupManager.getGroupMembersInfo(groupId, sctmemberList)
		if err != nil {
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		jsonResult, err := json.Marshal(r)
		if err != nil {
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		callback.OnSuccess(string(jsonResult))
	}()
}

func KickGroupMember(groupId string, userList string, reason string, callback Base) {
	if callback == nil {
		return
	}
	go func() {
		var sctMemberList []string
		err := json.Unmarshal([]byte(userList), &sctMemberList)
		if err != nil {
			sdkLog("unmarshal failed, ", err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		r, err := groupManager.kickGroupMember(groupId, sctMemberList, reason)
		if err != nil {
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}

		jsonResult, err := json.Marshal(r)
		if err != nil {
			sdkLog("marshal failed, ", err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		sdkLog("KickGroupMember, req: ", groupId, userList, reason, "resp: ", string(jsonResult))
		callback.OnSuccess(string(jsonResult))
	}()
}

func TransferGroupOwner(groupId, userId string, callback Base) {
	if callback == nil {
		return
	}
	go func() {
		err := groupManager.transferGroupOwner(groupId, userId)
		if err != nil {
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		callback.OnSuccess(DeFaultSuccessMsg)
	}()
}

func InviteUserToGroup(groupId, reason string, userList string, callback Base) {
	if callback == nil {
		return
	}
	go func() {
		var sctUserList []string
		err := json.Unmarshal([]byte(userList), &sctUserList)
		if err != nil {
			sdkLog("unmarshal failed, ", err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}

		result, err := groupManager.inviteUserToGroup(groupId, reason, sctUserList)
		if err != nil {
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		jsonResult, err := json.Marshal(result)
		if err != nil {
			sdkLog("marshal failed, ", err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		sdkLog("InviteUserToGroup, req: ", groupId, reason, userList, "resp: ", string(jsonResult))
		callback.OnSuccess(string(jsonResult))
	}()

}

func GetGroupApplicationList(callback Base) {
	if callback == nil {
		return
	}
	go func() {
		r, err := groupManager.getGroupApplicationList()
		if err != nil {
			sdkLog("getGroupApplicationList faild, ", err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		jsonResult, err := json.Marshal(r)
		if err != nil {
			sdkLog("getGroupApplicationList faild, ", err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		callback.OnSuccess(string(jsonResult))
	}()
}

func AcceptGroupApplication(application, reason string, callback Base) {
	if callback == nil {
		return
	}
	go func() {
		var sctApplication groupApplication
		err := json.Unmarshal([]byte(application), &sctApplication)
		if err != nil {
			sdkLog("Unmarshal, ", err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		err = groupManager.acceptGroupApplication(sctApplication, reason)
		if err != nil {
			sdkLog("acceptGroupApplication, ", err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		callback.OnSuccess(DeFaultSuccessMsg)
	}()
}

func RefuseGroupApplication(application, reason string, callback Base) {
	if callback == nil {
		return
	}
	go func() {
		var sctApplication groupApplication
		err := json.Unmarshal([]byte(application), &sctApplication)
		if err != nil {
			sdkLog("Unmarshal, ", err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		err = groupManager.refuseGroupApplication(sctApplication, reason)
		if err != nil {
			sdkLog("refuseGroupApplication, ", err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		callback.OnSuccess(DeFaultSuccessMsg)
	}()

}
