// Copyright © 2023 OpenIM. All rights reserved.
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

package group

import "errors"

func (x *CreateGroupReq) Check() error {
	if x.MemberUserIDs == nil {
		return errors.New("memberUserIDS is empty")
	}
	if x.GroupInfo == nil {
		return errors.New("groupInfo is empty")
	}
	if x.GroupInfo.GroupType > 2 || x.GroupInfo.GroupType < 0 {
		return errors.New("GroupType is invalid")
	}
	if x.OwnerUserID == "" {
		return errors.New("ownerUserID")
	}
	return nil
}

func (x *GetGroupsInfoReq) Check() error {
	if x.GroupIDs == nil {
		return errors.New("GroupIDs")
	}
	return nil
}

func (x *SetGroupInfoReq) Check() error {
	if x.GroupInfoForSet == nil {
		return errors.New("GroupInfoForSets is empty")
	}
	if x.GroupInfoForSet.GroupID == "" {
		return errors.New("GroupID is empty")
	}
	return nil
}

func (x *GetGroupApplicationListReq) Check() error {
	if x.Pagination == nil {
		return errors.New("pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errors.New("pageNumber is invalid")
	}
	if x.FromUserID == "" {
		return errors.New("fromUserID is empty")
	}
	return nil
}

func (x *GetUserReqApplicationListReq) Check() error {
	if x.UserID == "" {
		return errors.New("UserID is empty")
	}
	if x.Pagination == nil {
		return errors.New("pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errors.New("pageNumber is invalid")
	}
	return nil
}

func (x *TransferGroupOwnerReq) Check() error {
	if x.GroupID == "" {
		return errors.New("groupID is empty")
	}
	if x.OldOwnerUserID == "" {
		return errors.New("oldOwnerUserID is empty")
	}
	if x.NewOwnerUserID == "" {
		return errors.New("newOwnerUserID is empty")
	}
	return nil
}

func (x *JoinGroupReq) Check() error {
	if x.GroupID == "" {
		return errors.New("groupID is empty")
	}
	if x.JoinSource < 2 || x.JoinSource > 4 {
		return errors.New("joinSource is invalid")
	}
	if x.JoinSource == 2 {
		if x.InviterUserID == "" {
			return errors.New("inviterUserID is empty")
		}
	}
	return nil
}

func (x *GroupApplicationResponseReq) Check() error {
	if x.GroupID == "" {
		return errors.New("groupID is empty")
	}
	if x.FromUserID == "" {
		return errors.New("fromUserID is empty")
	}
	if x.HandleResult > 1 || x.HandleResult < -1 {
		return errors.New("handleResult is invalid")
	}
	return nil
}

func (x *QuitGroupReq) Check() error {
	if x.GroupID == "" {
		return errors.New("groupID is empty")
	}
	return nil
}

func (x *GetGroupMemberListReq) Check() error {
	if x.GroupID == "" {
		return errors.New("groupID is empty")
	}
	if x.Pagination == nil {
		return errors.New("pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errors.New("pageNumber is invalid")
	}
	if x.Filter < 0 || x.Filter > 5 {
		return errors.New("filter is invalid")
	}
	return nil
}

func (x *GetGroupMembersInfoReq) Check() error {
	if x.GroupID == "" {
		return errors.New("groupID is empty")
	}
	if x.UserIDs == nil {
		return errors.New("userIDs is empty")
	}
	return nil
}

func (x *KickGroupMemberReq) Check() error {
	if x.GroupID == "" {
		return errors.New("groupID is empty")
	}
	if x.KickedUserIDs == nil {
		return errors.New("kickUserIDs is empty")
	}
	return nil
}

func (x *GetJoinedGroupListReq) Check() error {
	if x.Pagination == nil {
		return errors.New("pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errors.New("pageNumber is invalid")
	}
	if x.FromUserID == "" {
		return errors.New("fromUserID is empty")
	}
	return nil
}

func (x *InviteUserToGroupReq) Check() error {
	if x.GroupID == "" {
		return errors.New("groupID is empty")
	}
	if x.InvitedUserIDs == nil {
		return errors.New("invitedUserIDs is empty")
	}
	return nil
}

func (x *GetGroupAllMemberReq) Check() error {
	if x.GroupID == "" {
		return errors.New("groupID is empty")
	}
	if x.Pagination == nil {
		return errors.New("pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errors.New("pageNumber is invalid")
	}
	return nil
}

func (x *GetGroupsReq) Check() error {
	if x.Pagination == nil {
		return errors.New("pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errors.New("pageNumber is invalid")
	}
	return nil
}

func (x *GetGroupMemberReq) Check() error {
	if x.GroupID == "" {
		return errors.New("groupID is empty")
	}
	return nil
}

func (x *GetGroupMembersCMSReq) Check() error {
	if x.GroupID == "" {
		return errors.New("groupID is empty")
	}
	if x.Pagination == nil {
		return errors.New("pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errors.New("pageNumber is invalid")
	}
	return nil
}

func (x *DismissGroupReq) Check() error {
	if x.GroupID == "" {
		return errors.New("groupID is empty")
	}
	return nil
}

func (x *MuteGroupMemberReq) Check() error {
	if x.GroupID == "" {
		return errors.New("groupID is empty")
	}
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	if x.MutedSeconds <= 0 {
		return errors.New("mutedSeconds is empty")
	}
	return nil
}

func (x *CancelMuteGroupMemberReq) Check() error {
	if x.GroupID == "" {
		return errors.New("groupID is empty")
	}
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	return nil
}

func (x *MuteGroupReq) Check() error {
	if x.GroupID == "" {
		return errors.New("groupID is empty")
	}
	return nil
}

func (x *CancelMuteGroupReq) Check() error {
	if x.GroupID == "" {
		return errors.New("groupID is empty")
	}
	return nil
}

func (x *GetJoinedSuperGroupListReq) Check() error {
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	return nil
}

func (x *GetSuperGroupsInfoReq) Check() error {
	if x.GroupIDs == nil {
		return errors.New("GroupIDs is empty")
	}
	return nil
}

func (x *SetGroupMemberInfo) Check() error {
	if x.GroupID == "" {
		return errors.New("GroupID is empty")
	}
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	return nil
}

func (x *SetGroupMemberInfoReq) Check() error {
	if x.Members == nil {
		return errors.New("Members is empty")
	}
	return nil
}

func (x *GetGroupAbstractInfoReq) Check() error {
	if x.GroupIDs == nil {
		return errors.New("GroupID is empty")
	}
	return nil
}

func (x *GetUserInGroupMembersReq) Check() error {
	if x.GroupIDs == nil {
		return errors.New("GroupID is empty")
	}
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	return nil
}

func (x *GetGroupMemberUserIDsReq) Check() error {
	if x.GroupID == "" {
		return errors.New("GroupID is empty")
	}
	return nil
}

func (x *GetGroupMemberRoleLevelReq) Check() error {
	if x.GroupID == "" {
		return errors.New("GroupID is empty")
	}
	if x.RoleLevels == nil {
		return errors.New("rolesLevel is empty")
	}
	return nil
}

func (x *GetGroupInfoCacheReq) Check() error {
	if x.GroupID == "" {
		return errors.New("GroupID is empty")
	}
	return nil
}

func (x *GetGroupMemberCacheReq) Check() error {
	if x.GroupID == "" {
		return errors.New("GroupID is empty")
	}
	if x.GroupMemberID == "" {
		return errors.New("GroupMemberID is empty")
	}
	return nil
}
