package open_im_sdk

import "open_im_sdk/open_im_sdk_callback"

func CreateGroup(callback open_im_sdk_callback.Base, operationID string, groupBaseInfo string, memberList string) {
	call(callback, operationID, userForSDK.Group().CreateGroup, groupBaseInfo, memberList)
}

func CreateGroupV2(callback open_im_sdk_callback.Base, operationID string, group string) {
	call(callback, operationID, userForSDK.Group().CreateGroupV2, group)
}

func JoinGroup(callback open_im_sdk_callback.Base, operationID string, groupID string, reqMsg string, joinSource int32) {
	call(callback, operationID, userForSDK.Group().JoinGroup, groupID, reqMsg, joinSource)
}

func QuitGroup(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, userForSDK.Group().QuitGroup, groupID)
}

func DismissGroup(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, userForSDK.Group().DismissGroup, groupID)
}

func ChangeGroupMute(callback open_im_sdk_callback.Base, operationID string, groupID string, isMute bool) {
	call(callback, operationID, userForSDK.Group().ChangeGroupMute, groupID, isMute)
}

func ChangeGroupMemberMute(callback open_im_sdk_callback.Base, operationID string, groupID string, userID string, mutedSeconds int) {
	call(callback, operationID, userForSDK.Group().ChangeGroupMemberMute, groupID, userID, mutedSeconds)
}

func SetGroupMemberRoleLevel(callback open_im_sdk_callback.Base, operationID string, groupID string, userID string, roleLevel int) {
	call(callback, operationID, userForSDK.Group().SetGroupMemberRoleLevel, groupID, userID, roleLevel)
}

func SetGroupMemberInfo(callback open_im_sdk_callback.Base, operationID string, groupMemberInfo string) {
	call(callback, operationID, userForSDK.Group().SetGroupMemberInfo, groupMemberInfo)
}

func GetJoinedGroupList(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, userForSDK.Group().GetJoinedGroupList)
}

func GetGroupsInfo(callback open_im_sdk_callback.Base, operationID string, groupIDList string) {
	call(callback, operationID, userForSDK.Group().GetGroupsInfo, groupIDList)
}

func SearchGroups(callback open_im_sdk_callback.Base, operationID string, searchParam string) {
	call(callback, operationID, userForSDK.Group().SearchGroups, searchParam)
}

func SetGroupInfo(callback open_im_sdk_callback.Base, operationID string, groupID string, groupInfo string) {
	call(callback, operationID, userForSDK.Group().SetGroupInfo, groupID, groupInfo)
}

func SetGroupVerification(callback open_im_sdk_callback.Base, operationID string, groupID string, verification int32) {
	call(callback, operationID, userForSDK.Group().SetGroupVerification, groupID, verification)
}

func SetGroupLookMemberInfo(callback open_im_sdk_callback.Base, operationID string, groupID string, rule int32) {
	call(callback, operationID, userForSDK.Group().SetGroupLookMemberInfo, groupID, rule)
}

func SetGroupApplyMemberFriend(callback open_im_sdk_callback.Base, operationID string, groupID string, rule int32) {
	call(callback, operationID, userForSDK.Group().SetGroupApplyMemberFriend, groupID, rule)
}

func GetGroupMemberList(callback open_im_sdk_callback.Base, operationID string, groupID string, filter int32, offset int32, count int32) {
	call(callback, operationID, userForSDK.Group().GetGroupMemberList, groupID, filter, offset, count)
}

func GetGroupMemberOwnerAndAdmin(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, userForSDK.Group().GetGroupMemberOwnerAndAdmin, groupID)
}

func GetGroupMemberListByJoinTimeFilter(callback open_im_sdk_callback.Base, operationID string, groupID string, offset int32, count int32, joinTimeBegin int64, joinTimeEnd int64, filterUserIDList string) {
	call(callback, operationID, userForSDK.Group().GetGroupMemberListByJoinTimeFilter, groupID, offset, count, joinTimeBegin, joinTimeEnd, filterUserIDList)
}

func GetGroupMembersInfo(callback open_im_sdk_callback.Base, operationID string, groupID string, userIDList string) {
	call(callback, operationID, userForSDK.Group().GetGroupMembersInfo, groupID, userIDList)
}

func KickGroupMember(callback open_im_sdk_callback.Base, operationID string, groupID string, reason string, userIDList string) {
	call(callback, operationID, userForSDK.Group().KickGroupMember, groupID, reason, userIDList)
}

func TransferGroupOwner(callback open_im_sdk_callback.Base, operationID string, groupID string, newOwnerUserID string) {
	call(callback, operationID, userForSDK.Group().TransferGroupOwner, groupID, newOwnerUserID)
}

func InviteUserToGroup(callback open_im_sdk_callback.Base, operationID string, groupID string, reason string, userIDList string) {
	call(callback, operationID, userForSDK.Group().InviteUserToGroup, groupID, reason, userIDList)
}

func GetRecvGroupApplicationList(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, userForSDK.Group().GetRecvGroupApplicationList)
}

func GetSendGroupApplicationList(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, userForSDK.Group().GetSendGroupApplicationList)
}

func AcceptGroupApplication(callback open_im_sdk_callback.Base, operationID string, groupID string, fromUserID string, handleMsg string) {
	call(callback, operationID, userForSDK.Group().AcceptGroupApplication, groupID, fromUserID, handleMsg)
}

func RefuseGroupApplication(callback open_im_sdk_callback.Base, operationID string, groupID string, fromUserID string, handleMsg string) {
	call(callback, operationID, userForSDK.Group().RefuseGroupApplication, groupID, fromUserID, handleMsg)
}

func SetGroupMemberNickname(callback open_im_sdk_callback.Base, operationID string, groupID string, userID string, groupMemberNickname string) {
	call(callback, operationID, userForSDK.Group().SetGroupMemberNickname, groupID, userID, groupMemberNickname)
}

func SearchGroupMembers(callback open_im_sdk_callback.Base, operationID string, searchParam string) {
	call(callback, operationID, userForSDK.Group().SearchGroupMembers, searchParam)
}
