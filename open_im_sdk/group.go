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

import "github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"

func CreateGroup(callback open_im_sdk_callback.Base, operationID string, groupReqInfo string) {
	call(callback, operationID, IMUserContext.Group().CreateGroup, groupReqInfo)
}

func JoinGroup(callback open_im_sdk_callback.Base, operationID string, groupID string, reqMsg string, joinSource int32, ex string) {
	call(callback, operationID, IMUserContext.Group().JoinGroup, groupID, reqMsg, joinSource, ex)
}

func QuitGroup(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, IMUserContext.Group().QuitGroup, groupID)
}

func DismissGroup(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, IMUserContext.Group().DismissGroup, groupID)
}

func ChangeGroupMute(callback open_im_sdk_callback.Base, operationID string, groupID string, isMute bool) {
	call(callback, operationID, IMUserContext.Group().ChangeGroupMute, groupID, isMute)
}

func ChangeGroupMemberMute(callback open_im_sdk_callback.Base, operationID string, groupID string, userID string, mutedSeconds int) {
	call(callback, operationID, IMUserContext.Group().ChangeGroupMemberMute, groupID, userID, mutedSeconds)
}

func TransferGroupOwner(callback open_im_sdk_callback.Base, operationID string, groupID string, newOwnerUserID string) {
	call(callback, operationID, IMUserContext.Group().TransferGroupOwner, groupID, newOwnerUserID)
}

func KickGroupMember(callback open_im_sdk_callback.Base, operationID string, groupID string, reason string, userIDList string) {
	call(callback, operationID, IMUserContext.Group().KickGroupMember, groupID, reason, userIDList)
}

func SetGroupInfo(callback open_im_sdk_callback.Base, operationID string, groupInfo string) {
	call(callback, operationID, IMUserContext.Group().SetGroupInfo, groupInfo)
}

func SetGroupMemberInfo(callback open_im_sdk_callback.Base, operationID string, groupMemberInfo string) {
	call(callback, operationID, IMUserContext.Group().SetGroupMemberInfo, groupMemberInfo)
}

func GetJoinedGroupList(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, IMUserContext.Group().GetJoinedGroupList)
}

func GetJoinedGroupListPage(callback open_im_sdk_callback.Base, operationID string, offset, count int32) {
	call(callback, operationID, IMUserContext.Group().GetJoinedGroupListPage, offset, count)
}

func GetSpecifiedGroupsInfo(callback open_im_sdk_callback.Base, operationID string, groupIDList string) {
	call(callback, operationID, IMUserContext.Group().GetSpecifiedGroupsInfo, groupIDList)
}

func SearchGroups(callback open_im_sdk_callback.Base, operationID string, searchParam string) {
	call(callback, operationID, IMUserContext.Group().SearchGroups, searchParam)
}

func GetGroupMemberOwnerAndAdmin(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, IMUserContext.Group().GetGroupMemberOwnerAndAdmin, groupID)
}

func GetGroupMemberListByJoinTimeFilter(callback open_im_sdk_callback.Base, operationID string, groupID string, offset int32, count int32, joinTimeBegin int64, joinTimeEnd int64, filterUserIDList string) {
	call(callback, operationID, IMUserContext.Group().GetGroupMemberListByJoinTimeFilter, groupID, offset, count, joinTimeBegin, joinTimeEnd, filterUserIDList)
}

func GetSpecifiedGroupMembersInfo(callback open_im_sdk_callback.Base, operationID string, groupID string, userIDList string) {
	call(callback, operationID, IMUserContext.Group().GetSpecifiedGroupMembersInfo, groupID, userIDList)
}

func GetGroupMemberList(callback open_im_sdk_callback.Base, operationID string, groupID string, filter int32, offset int32, count int32) {
	call(callback, operationID, IMUserContext.Group().GetGroupMemberList, groupID, filter, offset, count)
}

func GetGroupApplicationListAsRecipient(callback open_im_sdk_callback.Base, operationID, req string) {
	call(callback, operationID, IMUserContext.Group().GetGroupApplicationListAsRecipient, req)
}

func GetGroupApplicationListAsApplicant(callback open_im_sdk_callback.Base, operationID, req string) {
	call(callback, operationID, IMUserContext.Group().GetGroupApplicationListAsApplicant, req)
}

func SearchGroupMembers(callback open_im_sdk_callback.Base, operationID string, searchParam string) {
	call(callback, operationID, IMUserContext.Group().SearchGroupMembers, searchParam)
}

func IsJoinGroup(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, IMUserContext.Group().IsJoinGroup, groupID)
}

func GetUsersInGroup(callback open_im_sdk_callback.Base, operationID string, groupID, userIDList string) {
	call(callback, operationID, IMUserContext.Group().GetUsersInGroup, groupID, userIDList)
}

func InviteUserToGroup(callback open_im_sdk_callback.Base, operationID string, groupID string, reason string, userIDList string) {
	call(callback, operationID, IMUserContext.Group().InviteUserToGroup, groupID, reason, userIDList)
}

func AcceptGroupApplication(callback open_im_sdk_callback.Base, operationID string, groupID string, fromUserID string, handleMsg string) {
	call(callback, operationID, IMUserContext.Group().AcceptGroupApplication, groupID, fromUserID, handleMsg)
}

func RefuseGroupApplication(callback open_im_sdk_callback.Base, operationID string, groupID string, fromUserID string, handleMsg string) {
	call(callback, operationID, IMUserContext.Group().RefuseGroupApplication, groupID, fromUserID, handleMsg)
}

func CheckLocalGroupFullSync(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, IMUserContext.Group().CheckLocalGroupFullSync)
}

func CheckGroupMemberFullSync(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, IMUserContext.Group().CheckGroupMemberFullSync, groupID)
}

func GetGroupApplicationUnhandledCount(callback open_im_sdk_callback.Base, operationID string, req string) {
	call(callback, operationID, IMUserContext.Group().GetGroupApplicationUnhandledCount, req)
}
