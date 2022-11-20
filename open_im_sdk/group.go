package open_im_sdk

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
)

func CreateGroup(callback open_im_sdk_callback.Base, operationID string, groupBaseInfo string, memberList string) {
	if callback == nil {
		log.Error("callback is nil")
		return
	}
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().CreateGroup(callback, groupBaseInfo, memberList, operationID)
}

func JoinGroup(callback open_im_sdk_callback.Base, operationID string, groupID, reqMsg string, joinSource int32) {
	if callback == nil {
		log.Error("callback is nil")
		return
	}
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().JoinGroup(callback, groupID, reqMsg, joinSource, operationID)
}

func QuitGroup(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().QuitGroup(callback, groupID, operationID)
}

func DismissGroup(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().DismissGroup(callback, groupID, operationID)
}

func ChangeGroupMute(callback open_im_sdk_callback.Base, operationID string, groupID string, isMute bool) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().ChangeGroupMute(callback, groupID, isMute, operationID)
}

func ChangeGroupMemberMute(callback open_im_sdk_callback.Base, operationID string, groupID, userID string, mutedSeconds int) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().ChangeGroupMemberMute(callback, groupID, userID, uint32(mutedSeconds), operationID)
}

func SetGroupMemberRoleLevel(callback open_im_sdk_callback.Base, operationID string, groupID, userID string, roleLevel int) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().SetGroupMemberRoleLevel(callback, groupID, userID, roleLevel, operationID)
}

func GetJoinedGroupList(callback open_im_sdk_callback.Base, operationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().GetJoinedGroupList(callback, operationID)
}

func GetGroupsInfo(callback open_im_sdk_callback.Base, operationID string, groupIDList string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().GetGroupsInfo(callback, groupIDList, operationID)
}
func SearchGroups(callback open_im_sdk_callback.Base, operationID string, searchParam string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().SearchGroups(callback, searchParam, operationID)
}

func SetGroupInfo(callback open_im_sdk_callback.Base, operationID string, groupID string, groupInfo string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().SetGroupInfo(callback, groupInfo, groupID, operationID)
}

func SetGroupVerification(callback open_im_sdk_callback.Base, operationID string, groupID string, verification int32) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().SetGroupVerification(callback, verification, groupID, operationID)
}
func SetGroupLookMemberInfo(callback open_im_sdk_callback.Base, operationID string, groupID string, rule int32) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().SetGroupLookMemberInfo(callback, rule, groupID, operationID)
}
func SetGroupApplyMemberFriend(callback open_im_sdk_callback.Base, operationID string, groupID string, rule int32) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().SetGroupApplyMemberFriend(callback, rule, groupID, operationID)
}
func GetGroupMemberList(callback open_im_sdk_callback.Base, operationID string, groupID string, filter, offset, count int32) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().GetGroupMemberList(callback, groupID, filter, offset, count, operationID)
}

func GetGroupMemberOwnerAndAdmin(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().GetGroupMemberOwnerAndAdmin(callback, groupID, operationID)
}

func GetGroupMemberListByJoinTimeFilter(callback open_im_sdk_callback.Base, operationID string, groupID string, offset, count int32, joinTimeBegin, joinTimeEnd int64, filterUserIDList string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().GetGroupMemberListByJoinTimeFilter(callback, groupID, offset, count, joinTimeBegin, joinTimeEnd, filterUserIDList, operationID)
}

func GetGroupMembersInfo(callback open_im_sdk_callback.Base, operationID string, groupID string, userIDList string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().GetGroupMembersInfo(callback, groupID, userIDList, operationID)
}

func KickGroupMember(callback open_im_sdk_callback.Base, operationID string, groupID string, reason string, userIDList string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().KickGroupMember(callback, groupID, reason, userIDList, operationID)
}

func TransferGroupOwner(callback open_im_sdk_callback.Base, operationID string, groupID, newOwnerUserID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().TransferGroupOwner(callback, groupID, newOwnerUserID, operationID)
}

func InviteUserToGroup(callback open_im_sdk_callback.Base, operationID string, groupID, reason string, userIDList string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().InviteUserToGroup(callback, groupID, reason, userIDList, operationID)
}

func GetRecvGroupApplicationList(callback open_im_sdk_callback.Base, operationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().GetRecvGroupApplicationList(callback, operationID)
}

func GetSendGroupApplicationList(callback open_im_sdk_callback.Base, operationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().GetSendGroupApplicationList(callback, operationID)
}

func AcceptGroupApplication(callback open_im_sdk_callback.Base, operationID string, groupID, fromUserID, handleMsg string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().AcceptGroupApplication(callback, groupID, fromUserID, handleMsg, operationID)
}

func RefuseGroupApplication(callback open_im_sdk_callback.Base, operationID string, groupID, fromUserID, handleMsg string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().RefuseGroupApplication(callback, groupID, fromUserID, handleMsg, operationID)
}

func SetGroupMemberNickname(callback open_im_sdk_callback.Base, operationID string, groupID, userID, groupMemberNickname string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().SetGroupMemberNickname(callback, groupID, userID, groupMemberNickname, operationID)
}

//SearchGroupMembers
func SearchGroupMembers(callback open_im_sdk_callback.Base, operationID string, searchParam string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().SearchGroupMembers(callback, searchParam, operationID)
}
