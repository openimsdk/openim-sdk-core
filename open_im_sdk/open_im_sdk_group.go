package open_im_sdk

import (
	"encoding/json"
	"fmt"
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

func SetGroupListener(callback OnGroupListener) {
	if callback == nil {
		sdkLog("callback null")
		return
	}
	groupManager.listener = callback
	sdkLog("SetGroupListener ", callback)
}

func CreateGroup(gInfo string, memberList string, callback Base) {
	if callback == nil {
		sdkLog("callback is nil")
		return
	}
	go func() {
		var sctGroupInfo groupInfo
		err := json.Unmarshal([]byte(gInfo), &sctGroupInfo)
		if err != nil {
			sdkLog("unmarshal failed", gInfo, err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}

		var sctmemberList []createGroupMemberInfo
		err = json.Unmarshal([]byte(memberList), &sctmemberList)
		if err != nil {
			sdkLog("unmarshal failed, ", memberList, err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}

		sdkLog("unmarshal ok  ", gInfo, memberList, callback)

		resp, err := groupManager.createGroup(sctGroupInfo, sctmemberList)
		if err != nil {
			sdkLog("createGroup failed, ", err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		sdkLog("createGroup ok, callback success", sctGroupInfo, sctmemberList, resp)

		callback.OnSuccess(resp.Data.GroupId)
	}()
}

func JoinGroup(groupId, message string, callback Base) {
	if callback == nil {
		sdkLog("callback is nil")
		return
	}
	go func() {
		sdkLog(".................joinGroup begin ...............", groupId, message)
		err := groupManager.joinGroup(groupId, message)
		if err != nil {
			sdkLog("............joinGroup failed............ ", groupId, message, err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		sdkLog("...........joinGroup end, callback............... ", groupId, message)
		callback.OnSuccess(DeFaultSuccessMsg)
	}()
}

func QuitGroup(groupId string, callback Base) {
	if callback == nil {
		sdkLog("callback is nil")
		return
	}
	go func() {
		sdkLog("............quitGroup begin...............")
		err := groupManager.quitGroup(groupId)
		if err != nil {
			sdkLog(".........quitGroup failed.............", groupId, err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		sdkLog("..........quitGroup end callback...........", groupId)
		callback.OnSuccess(DeFaultSuccessMsg)
	}()
}

func GetJoinedGroupList(callback Base) {
	if callback == nil {
		sdkLog("Base callback is nil")
		return
	}
	go func() {
		groupInfoList, err := groupManager.getJoinedGroupListFromLocal()
		if err != nil {
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		sdkLog("getJoinedGroupListFromLocal, ", groupInfoList)

		for i, v := range groupInfoList {
			ownerId, err := findLocalGroupOwnerByGroupId(v.GroupId)
			if err != nil {
				sdkLog("findLocalGroupOwnerByGroupId failed,  ", err.Error(), v.GroupId)
				continue
			}
			sdkLog("findLocalGroupOwnerByGroupId ", v.GroupId, ownerId)
			v.OwnerId = ownerId
			number, err := getLocalGroupMemberNumByGroupId(v.GroupId)
			if err != nil {
				sdkLog("getLocalGroupMemberNumByGroupId failed, ", err.Error(), v.GroupId)
				continue
			}
			sdkLog("getLocalGroupMemberNumByGroupId ", v.GroupId, number)
			v.MemberCount = uint32(number)
			groupInfoList[i] = v
		}

		jsonGroupInfoList, err := json.Marshal(groupInfoList)
		if err != nil {
			sdkLog("marshal failed, ", err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		sdkLog("jsonGroupInfoList, ", string(jsonGroupInfoList))
		callback.OnSuccess(string(jsonGroupInfoList))
	}()
}

func GetGroupsInfo(groupIdList string, callback Base) {
	if callback == nil {
		sdkLog("Base callback is nil")
		return
	}
	go func() {
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
	}()

}

func SetGroupInfo(jsonGroupInfo string, callback Base) {
	if callback == nil {
		sdkLog("callback is nil")
		return
	}
	go func() {
		sdkLog("............SetGroupInfo begin...................")
		var newGroupInfo setGroupInfoReq
		err := json.Unmarshal([]byte(jsonGroupInfo), &newGroupInfo)
		if err != nil {
			sdkLog("............Unmarshal failed................", err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		err = groupManager.setGroupInfo(newGroupInfo)
		if err != nil {
			sdkLog("..........setGroupInfo failed........... ", err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		sdkLog(".........setGroupInfo end, callback...............", jsonGroupInfo)
		callback.OnSuccess(DeFaultSuccessMsg)
	}()
}

func GetGroupMemberList(groupId string, filter int32, next int32, callback Base) {
	if callback == nil {
		sdkLog("Base callback is nil")
		return
	}
	go func() {
		n, groupMemberResult, err := groupManager.getGroupMemberListFromLocal(groupId, filter, next)
		if err != nil {
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		sdkLog("getGroupMemberListFromLocal ", groupId, filter, next, groupMemberResult)

		var result getGroupMemberListResult
		if groupMemberResult == nil {
			groupMemberResult = make([]groupMemberFullInfo, 0)
		}
		result.Data = groupMemberResult
		result.NextSeq = n
		jsonGroupMemberResult, err := json.Marshal(result)
		if err != nil {
			sdkLog("marshal failed, ", err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		sdkLog("jsonGroupMemberResult ", string(jsonGroupMemberResult))
		callback.OnSuccess(string(jsonGroupMemberResult))
	}()
}

func GetGroupMembersInfo(groupId string, userList string, callback Base) {
	if callback == nil {
		sdkLog("Base callback is nil")
		return
	}
	go func() {
		var sctmemberList []string
		err := json.Unmarshal([]byte(userList), &sctmemberList)
		if err != nil {
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		sdkLog("GetGroupMembersInfo args: ", groupId, userList)
		r, err := groupManager.getGroupMembersInfoFromLocal(groupId, sctmemberList)
		if err != nil {
			sdkLog("getGroupMembersInfoFromLocal failed", groupId, err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		sdkLog("getGroupMembersInfoFromLocal, ", groupId, sctmemberList, r)

		jsonResult, err := json.Marshal(r)
		if err != nil {
			sdkLog("marshal failed, ", err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		sdkLog("jsonResult", string(jsonResult))
		callback.OnSuccess(string(jsonResult))
	}()
}

func KickGroupMember(groupId string, reason string, userList string, callback Base) {
	if callback == nil {
		sdkLog("callback null")
		return
	}
	go func() {
		var sctMemberList []string
		err := json.Unmarshal([]byte(userList), &sctMemberList)
		if err != nil {
			sdkLog("unmarshal failed, ", err.Error(), userList)
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		r, err := groupManager.kickGroupMember(groupId, sctMemberList, reason)
		if err != nil {
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		sdkLog("kickGroupMember, ", groupId, sctMemberList, reason)

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
		groupManager.syncGroupRequest()
		groupManager.syncGroupMemberByGroupId(groupId)
		callback.OnSuccess(DeFaultSuccessMsg)

		transfer := TransferGroupOwnerReq{
			GroupID:     groupId,
			OldOwner:    LoginUid,
			NewOwner:    userId,
			OperationID: operationIDGenerator(),
		}
		bTransfer, err := json.Marshal(transfer)
		if err != nil {
			sdkLog("TransferGroupOwner", err.Error())
			return
		}

		n := NotificationContent{1, TransferGroupTip, string(bTransfer)}
		autoSendMsg(createTextSystemMessage(n, TransferGroupOwnerTip), "", groupId, false, false, false)

	}()
}

func InviteUserToGroup(groupId, reason string, userList string, callback Base) {
	if callback == nil {
		sdkLog("callbak null")
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
		sdkLog("inviteUserToGroup, ", groupId, reason, sctUserList)

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
		return
	}()
}

func TsetGetGroupApplicationList(callback Base) string {
	if callback == nil {
		return ""
	}

	r, err := groupManager.getGroupApplicationList()
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

}

func AcceptGroupApplication(application, reason string, callback Base) {
	if callback == nil {
		return
	}
	go func() {
		var sctApplication GroupReqListInfo
		err := json.Unmarshal([]byte(application), &sctApplication)
		if err != nil {
			sdkLog("Unmarshal, ", err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}

		var access accessOrRefuseGroupApplicationReq
		access.OperationID = operationIDGenerator()
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

		err = groupManager.acceptGroupApplication(&access)
		if err != nil {
			sdkLog("acceptGroupApplication, ", err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		groupManager.syncGroupRequest()
		groupManager.syncGroupMemberByGroupId(sctApplication.GroupID)
		callback.OnSuccess(DeFaultSuccessMsg)

		user, err := getLoginUserInfoFromLocal()
		if err != nil {
			sdkLog("AcceptGroupApplication  getLoginUserInfoFromLocal err! ", err.Error())
			return
		}
		info := GroupApplicationInfo{
			Info:         access,
			HandUserID:   LoginUid,
			HandUserName: user.Name,
			HandUserIcon: user.Icon,
		}
		bInfo, err := json.Marshal(info)
		if err != nil {
			sdkLog("AcceptGroupApplication  json.Marshal err! ", err.Error())
			return
		}

		var name string
		if access.ToUser == "0" {
			name = access.FromUserNickName
		} else {
			name = access.ToUserNickname
		}
		defaultTip := fmt.Sprintf(AcceptGroupTip, name)
		n := NotificationContent{1, defaultTip, string(bInfo)}
		autoSendMsg(createTextSystemMessage(n, AcceptGroupApplicationTip), "", info.Info.GroupId, false, true, false)

	}()
}

func RefuseGroupApplication(application, reason string, callback Base) {
	if callback == nil {
		return
	}
	go func() {
		var sctApplication GroupReqListInfo
		err := json.Unmarshal([]byte(application), &sctApplication)
		if err != nil {
			sdkLog("Unmarshal, ", err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}

		var access accessOrRefuseGroupApplicationReq
		access.OperationID = operationIDGenerator()
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

		err = groupManager.refuseGroupApplication(&access)
		if err != nil {
			sdkLog("refuseGroupApplication, ", err.Error())
			callback.OnError(ErrCodeGroup, err.Error())
			return
		}
		groupManager.syncGroupRequest()
		callback.OnSuccess(DeFaultSuccessMsg)

		user, err := getLoginUserInfoFromLocal()
		if err != nil {
			sdkLog("RefuseGroupApplication  getLoginUserInfoFromLocal err! ", err.Error())
			return
		}
		info := GroupApplicationInfo{
			Info:         access,
			HandUserID:   LoginUid,
			HandUserName: user.Name,
			HandUserIcon: user.Icon,
		}
		bInfo, err := json.Marshal(info)
		if err != nil {
			sdkLog("RefuseGroupApplication  json.Marshal err! ", err.Error())
			return
		}

		recvID := ""
		if access.ToUser != "0" {
			recvID = access.ToUser
		} else {
			recvID = access.FromUser
		}

		n := NotificationContent{1, "", string(bInfo)}
		autoSendMsg(createTextSystemMessage(n, RefuseGroupApplicationTip), recvID, "", false, false, false)
	}()

}
