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

package group

import (
	"context"
	"time"

	"github.com/openimsdk/tools/errs"

	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/datafetcher"

	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdk_params_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"

	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/wrapperspb"
)

func (g *Group) CreateGroup(ctx context.Context, req *group.CreateGroupReq) (*sdkws.GroupInfo, error) {
	if req.OwnerUserID == "" {
		req.OwnerUserID = g.loginUserID
	}
	if req.GroupInfo.GroupType != constant.WorkingGroup {
		return nil, sdkerrs.ErrGroupType
	}
	req.GroupInfo.CreatorUserID = g.loginUserID
	resp, err := util.CallApi[group.CreateGroupResp](ctx, constant.CreateGroupRouter, req)
	if err != nil {
		return nil, err
	}

	g.groupSyncMutex.Lock()
	defer g.groupSyncMutex.Unlock()

	if err := g.IncrSyncJoinGroup(ctx); err != nil {
		return nil, err
	}
	if err := g.IncrSyncGroupAndMember(ctx, resp.GroupInfo.GroupID); err != nil {
		return nil, err
	}
	return resp.GroupInfo, nil
}

func (g *Group) JoinGroup(ctx context.Context, groupID, reqMsg string, joinSource int32, ex string) error {
	if err := util.ApiPost(ctx, constant.JoinGroupRouter, &group.JoinGroupReq{GroupID: groupID, ReqMessage: reqMsg, JoinSource: joinSource, InviterUserID: g.loginUserID, Ex: ex}, nil); err != nil {
		return err
	}
	if err := g.SyncSelfGroupApplications(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) QuitGroup(ctx context.Context, groupID string) error {
	if err := util.ApiPost(ctx, constant.QuitGroupRouter, &group.QuitGroupReq{GroupID: groupID}, nil); err != nil {
		return err
	}
	return nil
}

func (g *Group) DismissGroup(ctx context.Context, groupID string) error {
	if err := util.ApiPost(ctx, constant.DismissGroupRouter, &group.DismissGroupReq{GroupID: groupID}, nil); err != nil {
		return err
	}
	return nil
}

func (g *Group) SetGroupApplyMemberFriend(ctx context.Context, groupID string, rule int32) error {
	return g.SetGroupInfo(ctx, &sdkws.GroupInfoForSet{GroupID: groupID, ApplyMemberFriend: wrapperspb.Int32(rule)})
}

func (g *Group) SetGroupLookMemberInfo(ctx context.Context, groupID string, rule int32) error {
	return g.SetGroupInfo(ctx, &sdkws.GroupInfoForSet{GroupID: groupID, LookMemberInfo: wrapperspb.Int32(rule)})
}

func (g *Group) SetGroupVerification(ctx context.Context, groupID string, verification int32) error {
	return g.SetGroupInfo(ctx, &sdkws.GroupInfoForSet{GroupID: groupID, NeedVerification: wrapperspb.Int32(verification)})
}

func (g *Group) ChangeGroupMute(ctx context.Context, groupID string, isMute bool) (err error) {
	if isMute {
		err = util.ApiPost(ctx, constant.MuteGroupRouter, &group.MuteGroupReq{GroupID: groupID}, nil)
	} else {
		err = util.ApiPost(ctx, constant.CancelMuteGroupRouter, &group.CancelMuteGroupReq{GroupID: groupID}, nil)
	}
	if err != nil {
		return err
	}

	g.groupSyncMutex.Lock()
	defer g.groupSyncMutex.Unlock()

	if err := g.IncrSyncGroupAndMember(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) ChangeGroupMemberMute(ctx context.Context, groupID, userID string, mutedSeconds int) (err error) {
	if mutedSeconds == 0 {
		err = util.ApiPost(ctx, constant.CancelMuteGroupMemberRouter, &group.CancelMuteGroupMemberReq{GroupID: groupID, UserID: userID}, nil)
	} else {
		err = util.ApiPost(ctx, constant.MuteGroupMemberRouter, &group.MuteGroupMemberReq{GroupID: groupID, UserID: userID, MutedSeconds: uint32(mutedSeconds)}, nil)
	}
	if err != nil {
		return err
	}
	return nil
}

func (g *Group) TransferGroupOwner(ctx context.Context, groupID, newOwnerUserID string) error {
	if err := util.ApiPost(ctx, constant.TransferGroupRouter, &group.TransferGroupOwnerReq{GroupID: groupID, OldOwnerUserID: g.loginUserID, NewOwnerUserID: newOwnerUserID}, nil); err != nil {
		return err
	}

	g.groupSyncMutex.Lock()
	defer g.groupSyncMutex.Unlock()

	if err := g.IncrSyncGroupAndMember(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) KickGroupMember(ctx context.Context, groupID string, reason string, userIDList []string) error {
	if err := util.ApiPost(ctx, constant.KickGroupMemberRouter, &group.KickGroupMemberReq{GroupID: groupID, KickedUserIDs: userIDList, Reason: reason}, nil); err != nil {
		return err
	}

	g.groupSyncMutex.Lock()
	defer g.groupSyncMutex.Unlock()

	return g.IncrSyncGroupAndMember(ctx, groupID)
}

func (g *Group) SetGroupInfo(ctx context.Context, groupInfo *sdkws.GroupInfoForSet) error {
	if err := util.ApiPost(ctx, constant.SetGroupInfoRouter, &group.SetGroupInfoReq{GroupInfoForSet: groupInfo}, nil); err != nil {
		return err
	}

	g.groupSyncMutex.Lock()
	defer g.groupSyncMutex.Unlock()

	return g.IncrSyncJoinGroup(ctx)
}

func (g *Group) SetGroupMemberInfo(ctx context.Context, groupMemberInfo *group.SetGroupMemberInfo) error {
	if err := util.ApiPost(ctx, constant.SetGroupMemberInfoRouter, &group.SetGroupMemberInfoReq{Members: []*group.SetGroupMemberInfo{groupMemberInfo}}, nil); err != nil {
		return err
	}

	g.groupSyncMutex.Lock()
	defer g.groupSyncMutex.Unlock()

	return g.IncrSyncGroupAndMember(ctx, groupMemberInfo.GroupID)
}

func (g *Group) SetGroupMemberRoleLevel(ctx context.Context, groupID, userID string, roleLevel int) error {
	return g.SetGroupMemberInfo(ctx, &group.SetGroupMemberInfo{GroupID: groupID, UserID: userID, RoleLevel: wrapperspb.Int32(int32(roleLevel))})
}

func (g *Group) SetGroupMemberNickname(ctx context.Context, groupID, userID string, groupMemberNickname string) error {
	return g.SetGroupMemberInfo(ctx, &group.SetGroupMemberInfo{GroupID: groupID, UserID: userID, Nickname: wrapperspb.String(groupMemberNickname)})
}

func (g *Group) GetJoinedGroupList(ctx context.Context) ([]*model_struct.LocalGroup, error) {
	return g.db.GetJoinedGroupListDB(ctx)
}

func (g *Group) GetJoinedGroupListPage(ctx context.Context, offset, count int32) ([]*model_struct.LocalGroup, error) {
	dataFetcher := datafetcher.NewDataFetcher(
		g.db,
		g.groupTableName(),
		g.loginUserID,
		func(localGroup *model_struct.LocalGroup) string {
			return localGroup.GroupID
		},
		func(ctx context.Context, values []*model_struct.LocalGroup) error {
			return g.db.BatchInsertGroup(ctx, values)
		},
		func(ctx context.Context, groupIDs []string) ([]*model_struct.LocalGroup, bool, error) {
			localGroups, err := g.db.GetGroups(ctx, groupIDs)
			return localGroups, true, err
		},
		func(ctx context.Context, groupIDs []string) ([]*model_struct.LocalGroup, error) {
			serverGroupInfo, err := g.getGroupsInfoFromSvr(ctx, groupIDs)
			if err != nil {
				return nil, err
			}
			return datautil.Batch(ServerGroupToLocalGroup, serverGroupInfo), nil
		},
	)
	return dataFetcher.FetchWithPagination(ctx, int(offset), int(count))
}

func (g *Group) GetSpecifiedGroupsInfo(ctx context.Context, groupIDs []string) ([]*model_struct.LocalGroup, error) {
	dataFetcher := datafetcher.NewDataFetcher(
		g.db,
		g.groupTableName(),
		g.loginUserID,
		func(localGroup *model_struct.LocalGroup) string {
			return localGroup.GroupID
		},
		func(ctx context.Context, values []*model_struct.LocalGroup) error {
			return g.db.BatchInsertGroup(ctx, values)
		},
		func(ctx context.Context, groupIDs []string) ([]*model_struct.LocalGroup, bool, error) {
			localGroups, err := g.db.GetGroups(ctx, groupIDs)
			return localGroups, true, err
		},
		func(ctx context.Context, groupIDs []string) ([]*model_struct.LocalGroup, error) {
			serverGroupInfo, err := g.getGroupsInfoFromSvr(ctx, groupIDs)
			if err != nil {
				return nil, err
			}
			return datautil.Batch(ServerGroupToLocalGroup, serverGroupInfo), nil
		},
	)
	return dataFetcher.FetchMissingAndCombineLocal(ctx, groupIDs)
}

func (g *Group) GetJoinedGroupListPageV2(ctx context.Context, offset, count int32) (*GetGroupListV2Response, error) {
	dataFetcher := datafetcher.NewDataFetcher(
		g.db,
		g.groupTableName(),
		g.loginUserID,
		func(localGroup *model_struct.LocalGroup) string {
			return localGroup.GroupID
		},
		func(ctx context.Context, values []*model_struct.LocalGroup) error {
			return g.db.BatchInsertGroup(ctx, values)
		},
		func(ctx context.Context, groupIDs []string) ([]*model_struct.LocalGroup, bool, error) {
			localGroups, err := g.db.GetGroups(ctx, groupIDs)
			return localGroups, true, err
		},
		func(ctx context.Context, groupIDs []string) ([]*model_struct.LocalGroup, error) {
			serverGroupInfo, err := g.getGroupsInfoFromSvr(ctx, groupIDs)
			if err != nil {
				return nil, err
			}
			return datautil.Batch(ServerGroupToLocalGroup, serverGroupInfo), nil
		},
	)

	groupsList, isEnd, err := dataFetcher.FetchWithPaginationV2(ctx, int(offset), int(count))
	if err != nil {
		return nil, err
	}

	resp := &GetGroupListV2Response{
		GroupsList: groupsList,
		IsEnd:      isEnd,
	}
	return resp, nil
}

func (g *Group) SearchGroups(ctx context.Context, param sdk_params_callback.SearchGroupsParam) ([]*model_struct.LocalGroup, error) {
	if len(param.KeywordList) == 0 || (!param.IsSearchGroupName && !param.IsSearchGroupID) {
		return nil, sdkerrs.ErrArgs.WrapMsg("keyword is null or search field all false")
	}
	groups, err := g.db.GetAllGroupInfoByGroupIDOrGroupName(ctx, param.KeywordList[0], param.IsSearchGroupID, param.IsSearchGroupName) // todo	param.KeywordList[0]
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// funcation (g *Group) SetGroupInfo(ctx context.Context, groupInfo *sdk_params_callback.SetGroupInfoParam, groupID string) error {
//	return g.SetGroupInfo(ctx, &sdkws.GroupInfoForSet{
//		GroupID:          groupID,
//		GroupName:        groupInfo.GroupName,
//		Notification:     groupInfo.Notification,
//		Introduction:     groupInfo.Introduction,
//		FaceURL:          groupInfo.FaceURL,
//		Ex:               groupInfo.Ex,
//		NeedVerification: wrapperspb.Int32Ptr(groupInfo.NeedVerification),
//	})
// }

func (g *Group) GetGroupMemberOwnerAndAdmin(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error) {
	return g.db.GetGroupMemberOwnerAndAdminDB(ctx, groupID)
}

func (g *Group) GetGroupMemberListByJoinTimeFilter(ctx context.Context, groupID string, offset, count int32, joinTimeBegin, joinTimeEnd int64, userIDs []string) ([]*model_struct.LocalGroupMember, error) {
	if joinTimeEnd == 0 {
		joinTimeEnd = time.Now().UnixMilli()
	}

	dataFetcher := datafetcher.NewDataFetcher(
		g.db,
		g.groupAndMemberVersionTableName(),
		groupID,
		func(localGroupMember *model_struct.LocalGroupMember) string {
			return localGroupMember.UserID
		},
		func(ctx context.Context, values []*model_struct.LocalGroupMember) error {
			return g.db.BatchInsertGroupMember(ctx, values)
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalGroupMember, bool, error) {
			localGroupMembers, err := g.db.GetGroupMemberListSplitByJoinTimeFilter(ctx, groupID, int(offset), int(count), joinTimeBegin, joinTimeEnd, userIDs)
			return localGroupMembers, true, err
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalGroupMember, error) {
			serverGroupMember, err := g.GetDesignatedGroupMembers(ctx, groupID, userIDs)
			if err != nil {
				return nil, err
			}
			return datautil.Batch(ServerGroupMemberToLocalGroupMember, serverGroupMember), nil
		},
	)

	return dataFetcher.FetchWithPagination(ctx, int(offset), int(count))
}

func (g *Group) GetGroupMemberListByJoinTimeFilterV2(ctx context.Context, groupID string, offset, count int32, joinTimeBegin, joinTimeEnd int64, userIDs []string) (*GetGroupMemberListV2Response, error) {
	if joinTimeEnd == 0 {
		joinTimeEnd = time.Now().UnixMilli()
	}

	dataFetcher := datafetcher.NewDataFetcher(
		g.db,
		g.groupAndMemberVersionTableName(),
		groupID,
		func(localGroupMember *model_struct.LocalGroupMember) string {
			return localGroupMember.UserID
		},
		func(ctx context.Context, values []*model_struct.LocalGroupMember) error {
			return g.db.BatchInsertGroupMember(ctx, values)
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalGroupMember, bool, error) {
			localGroupMembers, err := g.db.GetGroupMemberListSplitByJoinTimeFilter(ctx, groupID, int(offset), int(count), joinTimeBegin, joinTimeEnd, userIDs)
			return localGroupMembers, true, err
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalGroupMember, error) {
			serverGroupMember, err := g.GetDesignatedGroupMembers(ctx, groupID, userIDs)
			if err != nil {
				return nil, err
			}
			return datautil.Batch(ServerGroupMemberToLocalGroupMember, serverGroupMember), nil
		},
	)

	groupMembersList, isEnd, err := dataFetcher.FetchWithPaginationV2(ctx, int(offset), int(count))
	if err != nil {
		return nil, err
	}
	resp := &GetGroupMemberListV2Response{
		GroupMembersList: groupMembersList,
		IsEnd:            isEnd,
	}
	return resp, nil
}

func (g *Group) GetSpecifiedGroupMembersInfo(ctx context.Context, groupID string, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	lvs, err := g.db.GetVersionSync(ctx, g.groupTableName(), g.loginUserID)
	if err != nil {
		return nil, err
	}
	if datautil.Contain(groupID, lvs.UIDList...) {

		_, err := g.db.GetVersionSync(ctx, g.groupAndMemberVersionTableName(), groupID)
		if err != nil {
			if errs.Unwrap(err) != errs.ErrRecordNotFound {
				return nil, err
			}

			g.groupSyncMutex.Lock()
			defer g.groupSyncMutex.Unlock()

			err := g.IncrSyncGroupAndMember(ctx, groupID)
			if err != nil {
				return nil, err
			}
		}
	} else { // If the user is no longer in the group, return nil immediately
		return nil, nil
	}
	dataFetcher := datafetcher.NewDataFetcher(
		g.db,
		g.groupAndMemberVersionTableName(),
		groupID,
		func(localGroupMember *model_struct.LocalGroupMember) string {
			return localGroupMember.UserID
		},
		func(ctx context.Context, values []*model_struct.LocalGroupMember) error {
			return g.db.BatchInsertGroupMember(ctx, values)
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalGroupMember, bool, error) {
			localGroupMembers, err := g.db.GetGroupSomeMemberInfo(ctx, groupID, userIDList)
			if err != nil {
				return nil, false, err
			}
			localGroup, err := g.db.GetGroupInfoByGroupID(ctx, groupID)
			if err != nil {
				return nil, false, err
			}
			if localGroup.MemberCount < groupMemberSyncLimit {
				return localGroupMembers, false, nil
			}
			return localGroupMembers, true, nil
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalGroupMember, error) {
			serverGroupMember, err := g.GetDesignatedGroupMembers(ctx, groupID, userIDs)
			if err != nil {
				return nil, err
			}
			if len(serverGroupMember) == 0 {
				return nil, nil
			}
			return datautil.Batch(ServerGroupMemberToLocalGroupMember, serverGroupMember), nil
		},
	)
	return dataFetcher.FetchMissingAndFillLocal(ctx, userIDList)
}

func (g *Group) GetGroupMemberList(ctx context.Context, groupID string, filter, offset, count int32) ([]*model_struct.LocalGroupMember, error) {
	lvs, err := g.db.GetVersionSync(ctx, g.groupTableName(), g.loginUserID)
	if err != nil {
		return nil, err
	}
	if datautil.Contain(groupID, lvs.UIDList...) {

		_, err := g.db.GetVersionSync(ctx, g.groupAndMemberVersionTableName(), groupID)
		if err != nil {
			if errs.Unwrap(err) != errs.ErrRecordNotFound {
				return nil, err
			}

			g.groupSyncMutex.Lock()
			defer g.groupSyncMutex.Unlock()

			err := g.IncrSyncGroupAndMember(ctx, groupID)
			if err != nil {
				return nil, err
			}
		}
	} else { // If the user is no longer in the group, return nil immediately
		return nil, nil
	}

	dataFetcher := datafetcher.NewDataFetcher(
		g.db,
		g.groupAndMemberVersionTableName(),
		groupID,
		func(localGroupMember *model_struct.LocalGroupMember) string {
			return localGroupMember.UserID
		},
		func(ctx context.Context, values []*model_struct.LocalGroupMember) error {
			return g.db.BatchInsertGroupMember(ctx, values)
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalGroupMember, bool, error) {
			localGroupMembers, err := g.db.GetGroupMemberListByUserIDs(ctx, groupID, filter, userIDs)
			if err != nil {
				return nil, false, err
			}
			switch filter {
			case constant.GroupFilterOwner:
				fallthrough
			case constant.GroupFilterAdmin:
				fallthrough
			case constant.GroupFilterOwnerAndAdmin:
				return localGroupMembers, false, nil
			case constant.GroupFilterAll:
			case constant.GroupFilterOrdinaryUsers:
			case constant.GroupFilterAdminAndOrdinaryUsers:
				return localGroupMembers, true, nil
			}
			return nil, false, sdkerrs.ErrArgs
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalGroupMember, error) {
			serverGroupMember, err := g.GetDesignatedGroupMembers(ctx, groupID, userIDs)
			if err != nil {
				return nil, err
			}
			return datautil.Batch(ServerGroupMemberToLocalGroupMember, serverGroupMember), nil
		},
	)
	switch filter {
	case constant.GroupFilterOrdinaryUsers:
		groupOwnerAndGroupMember, err := g.db.GetGroupMemberListSplit(ctx, groupID, constant.GroupFilterOwnerAndAdmin, 0, 100)
		if err != nil {
			return nil, err
		}
		offset = offset + int32(len(groupOwnerAndGroupMember))
	case constant.GroupFilterAdminAndOrdinaryUsers:
		groupOwnerAndGroupMember, err := g.db.GetGroupMemberListSplit(ctx, groupID, constant.GroupFilterOwner, 0, 100)
		if err != nil {
			return nil, err
		}
		offset = offset + int32(len(groupOwnerAndGroupMember))
	}
	return dataFetcher.FetchWithPagination(ctx, int(offset), int(count))
}

func (g *Group) GetGroupMemberListV2(ctx context.Context, groupID string, filter, offset, count int32) (*GetGroupMemberListV2Response, error) {
	dataFetcher := datafetcher.NewDataFetcher(
		g.db,
		g.groupAndMemberVersionTableName(),
		groupID,
		func(localGroupMember *model_struct.LocalGroupMember) string {
			return localGroupMember.UserID
		},
		func(ctx context.Context, values []*model_struct.LocalGroupMember) error {
			return g.db.BatchInsertGroupMember(ctx, values)
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalGroupMember, bool, error) {
			localGroupMembers, err := g.db.GetGroupMemberListByUserIDs(ctx, groupID, filter, userIDs)
			return localGroupMembers, true, err
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalGroupMember, error) {
			serverGroupMember, err := g.GetDesignatedGroupMembers(ctx, groupID, userIDs)
			if err != nil {
				return nil, err
			}
			return datautil.Batch(ServerGroupMemberToLocalGroupMember, serverGroupMember), nil
		},
	)
	groupMembersList, isEnd, err := dataFetcher.FetchWithPaginationV2(ctx, int(offset), int(count))
	if err != nil {
		return nil, err
	}
	resp := &GetGroupMemberListV2Response{
		GroupMembersList: groupMembersList,
		IsEnd:            isEnd,
	}

	return resp, nil
}

func (g *Group) GetGroupApplicationListAsRecipient(ctx context.Context) ([]*model_struct.LocalAdminGroupRequest, error) {
	return g.db.GetAdminGroupApplication(ctx)
}

func (g *Group) GetGroupApplicationListAsApplicant(ctx context.Context) ([]*model_struct.LocalGroupRequest, error) {
	return g.db.GetSendGroupApplication(ctx)
}

func (g *Group) SearchGroupMembers(ctx context.Context, searchParam *sdk_params_callback.SearchGroupMembersParam) ([]*model_struct.LocalGroupMember, error) {
	return g.db.SearchGroupMembersDB(ctx, searchParam.KeywordList[0], searchParam.GroupID, searchParam.IsSearchMemberNickname, searchParam.IsSearchUserID, searchParam.Offset, searchParam.Count)
}

func (g *Group) IsJoinGroup(ctx context.Context, groupID string) (bool, error) {
	lvs, err := g.db.GetVersionSync(ctx, g.groupTableName(), g.loginUserID)
	if err != nil {
		return false, err
	}
	if datautil.Contain(groupID, lvs.UIDList...) {
		return true, nil
	}
	return false, nil
}

func (g *Group) GetUsersInGroup(ctx context.Context, groupID string, userIDList []string) ([]string, error) {
	lvs, err := g.db.GetVersionSync(ctx, g.groupTableName(), g.loginUserID)
	if err != nil {
		return nil, err
	}
	if !datautil.Contain(groupID, lvs.UIDList...) {
		return nil, nil
	}
	lvs, err = g.db.GetVersionSync(ctx, g.groupAndMemberVersionTableName(), groupID)
	if err != nil {
		return nil, err
	}

	groupMembersMap := datautil.SliceSetAny(lvs.UIDList, func(e string) string {
		return e
	})

	var usersInGroup []string
	for _, userID := range userIDList {
		if _, exists := groupMembersMap[userID]; exists {
			usersInGroup = append(usersInGroup, userID)
		}
	}

	return usersInGroup, nil

}

func (g *Group) InviteUserToGroup(ctx context.Context, groupID, reason string, userIDList []string) error {
	if err := util.ApiPost(ctx, constant.InviteUserToGroupRouter, &group.InviteUserToGroupReq{GroupID: groupID, Reason: reason, InvitedUserIDs: userIDList}, nil); err != nil {
		return err
	}

	g.groupSyncMutex.Lock()
	defer g.groupSyncMutex.Unlock()

	if err := g.IncrSyncGroupAndMember(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) AcceptGroupApplication(ctx context.Context, groupID, fromUserID, handleMsg string) error {
	return g.HandlerGroupApplication(ctx, &group.GroupApplicationResponseReq{GroupID: groupID, FromUserID: fromUserID, HandledMsg: handleMsg, HandleResult: constant.GroupResponseAgree})
}

func (g *Group) RefuseGroupApplication(ctx context.Context, groupID, fromUserID, handleMsg string) error {
	return g.HandlerGroupApplication(ctx, &group.GroupApplicationResponseReq{GroupID: groupID, FromUserID: fromUserID, HandledMsg: handleMsg, HandleResult: constant.GroupResponseRefuse})
}

func (g *Group) HandlerGroupApplication(ctx context.Context, req *group.GroupApplicationResponseReq) error {
	if err := util.ApiPost(ctx, constant.AcceptGroupApplicationRouter, req, nil); err != nil {
		return err
	}
	// SyncAdminGroupApplication todo
	return nil
}

//func (g *Group) SearchGroupMembersV2(ctx context.Context, req *group.SearchGroupMemberReq) ([]*model_struct.LocalGroupMember, error) {
//	if err := req.Check(); err != nil {
//		return nil, err
//	}
//	info, err := g.db.GetGroupInfoByGroupID(ctx, req.GroupID)
//	if err != nil {
//		return nil, err
//	}
//	if info.MemberCount <= pconstant.MaxSyncPullNumber {
//		return g.db.SearchGroupMembersDB(ctx, req.Keyword, req.GroupID, true, false,
//			int((req.Pagination.PageNumber-1)*req.Pagination.ShowNumber), int(req.Pagination.ShowNumber))
//	}
//	resp, err := util.CallApi[group.SearchGroupMemberResp](ctx, constant.SearchGroupMember, req)
//	if err != nil {
//		return nil, err
//	}
//	return datautil.Slice(resp.Members, g.pbGroupMemberToLocal), nil
//}

func (g *Group) pbGroupMemberToLocal(pb *sdkws.GroupMemberFullInfo) *model_struct.LocalGroupMember {
	return &model_struct.LocalGroupMember{
		GroupID:        pb.GroupID,
		UserID:         pb.UserID,
		Nickname:       pb.Nickname,
		FaceURL:        pb.FaceURL,
		RoleLevel:      pb.RoleLevel,
		JoinTime:       pb.JoinTime,
		JoinSource:     pb.JoinSource,
		InviterUserID:  pb.InviterUserID,
		MuteEndTime:    pb.MuteEndTime,
		OperatorUserID: pb.OperatorUserID,
		Ex:             pb.Ex,
		// AttachedInfo:   pb.AttachedInfo,
	}
}
