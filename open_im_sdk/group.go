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

package open_im_sdk

import "open_im_sdk/open_im_sdk_callback"

//funcation CreateGroup(callback open_im_sdk_callback.Base, operationID string, groupBaseInfo string, memberList string) {
//	call(callback, operationID, UserForSDK.Group().CreateGroup, groupBaseInfo, memberList)
//}

func CreateGroup(callback open_im_sdk_callback.Base, operationID string, groupReqInfo string) {
	call(callback, operationID, UserForSDK.Group().CreateGroup, groupReqInfo)
}

func JoinGroup(callback open_im_sdk_callback.Base, operationID string, groupID string, reqMsg string, joinSource int32) {
	call(callback, operationID, UserForSDK.Group().JoinGroup, groupID, reqMsg, joinSource)
}

func QuitGroup(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, UserForSDK.Group().QuitGroup, groupID)
}

func DismissGroup(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, UserForSDK.Group().DismissGroup, groupID)
}

func ChangeGroupMute(callback open_im_sdk_callback.Base, operationID string, groupID string, isMute bool) {
	call(callback, operationID, UserForSDK.Group().ChangeGroupMute, groupID, isMute)
}

func ChangeGroupMemberMute(callback open_im_sdk_callback.Base, operationID string, groupID string, userID string, mutedSeconds int) {
	call(callback, operationID, UserForSDK.Group().ChangeGroupMemberMute, groupID, userID, mutedSeconds)
}

func SetGroupMemberRoleLevel(callback open_im_sdk_callback.Base, operationID string, groupID string, userID string, roleLevel int) {
	call(callback, operationID, UserForSDK.Group().SetGroupMemberRoleLevel, groupID, userID, roleLevel)
}

func SetGroupMemberInfo(callback open_im_sdk_callback.Base, operationID string, groupMemberInfo string) {
	call(callback, operationID, UserForSDK.Group().SetGroupMemberInfo, groupMemberInfo)
}

func GetJoinedGroupList(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Group().GetJoinedGroupList)
}

func GetSpecifiedGroupsInfo(callback open_im_sdk_callback.Base, operationID string, groupIDList string) {
	call(callback, operationID, UserForSDK.Group().GetSpecifiedGroupsInfo, groupIDList)
}

func SearchGroups(callback open_im_sdk_callback.Base, operationID string, searchParam string) {
	call(callback, operationID, UserForSDK.Group().SearchGroups, searchParam)
}

func SetGroupInfo(callback open_im_sdk_callback.Base, operationID string, groupInfo string) {
	call(callback, operationID, UserForSDK.Group().SetGroupInfo, groupInfo)
}

func SetGroupVerification(callback open_im_sdk_callback.Base, operationID string, groupID string, verification int32) {
	call(callback, operationID, UserForSDK.Group().SetGroupVerification, groupID, verification)
}

func SetGroupLookMemberInfo(callback open_im_sdk_callback.Base, operationID string, groupID string, rule int32) {
	call(callback, operationID, UserForSDK.Group().SetGroupLookMemberInfo, groupID, rule)
}

func SetGroupApplyMemberFriend(callback open_im_sdk_callback.Base, operationID string, groupID string, rule int32) {
	call(callback, operationID, UserForSDK.Group().SetGroupApplyMemberFriend, groupID, rule)
}

func GetGroupMemberList(callback open_im_sdk_callback.Base, operationID string, groupID string, filter int32, offset int32, count int32) {
	call(callback, operationID, UserForSDK.Group().GetGroupMemberList, groupID, filter, offset, count)
}

func GetGroupMemberOwnerAndAdmin(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, UserForSDK.Group().GetGroupMemberOwnerAndAdmin, groupID)
}

func GetGroupMemberListByJoinTimeFilter(callback open_im_sdk_callback.Base, operationID string, groupID string, offset int32, count int32, joinTimeBegin int64, joinTimeEnd int64, filterUserIDList string) {
	call(callback, operationID, UserForSDK.Group().GetGroupMemberListByJoinTimeFilter, groupID, offset, count, joinTimeBegin, joinTimeEnd, filterUserIDList)
}

func GetSpecifiedGroupMembersInfo(callback open_im_sdk_callback.Base, operationID string, groupID string, userIDList string) {
	call(callback, operationID, UserForSDK.Group().GetSpecifiedGroupMembersInfo, groupID, userIDList)
}

func KickGroupMember(callback open_im_sdk_callback.Base, operationID string, groupID string, reason string, userIDList string) {
	call(callback, operationID, UserForSDK.Group().KickGroupMember, groupID, reason, userIDList)
}

func TransferGroupOwner(callback open_im_sdk_callback.Base, operationID string, groupID string, newOwnerUserID string) {
	call(callback, operationID, UserForSDK.Group().TransferGroupOwner, groupID, newOwnerUserID)
}

func InviteUserToGroup(callback open_im_sdk_callback.Base, operationID string, groupID string, reason string, userIDList string) {
	call(callback, operationID, UserForSDK.Group().InviteUserToGroup, groupID, reason, userIDList)
}

func GetGroupApplicationListAsRecipient(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Group().GetGroupApplicationListAsRecipient)
}

func GetGroupApplicationListAsApplicant(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Group().GetGroupApplicationListAsApplicant)
}

func AcceptGroupApplication(callback open_im_sdk_callback.Base, operationID string, groupID string, fromUserID string, handleMsg string) {
	call(callback, operationID, UserForSDK.Group().AcceptGroupApplication, groupID, fromUserID, handleMsg)
}

func RefuseGroupApplication(callback open_im_sdk_callback.Base, operationID string, groupID string, fromUserID string, handleMsg string) {
	call(callback, operationID, UserForSDK.Group().RefuseGroupApplication, groupID, fromUserID, handleMsg)
}

func SetGroupMemberNickname(callback open_im_sdk_callback.Base, operationID string, groupID string, userID string, groupMemberNickname string) {
	call(callback, operationID, UserForSDK.Group().SetGroupMemberNickname, groupID, userID, groupMemberNickname)
}

func SearchGroupMembers(callback open_im_sdk_callback.Base, operationID string, searchParam string) {
	call(callback, operationID, UserForSDK.Group().SearchGroupMembers, searchParam)
}

func IsJoinGroup(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, UserForSDK.Group().IsJoinGroup, groupID)
}
