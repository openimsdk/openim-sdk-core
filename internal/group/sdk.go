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
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/sdkerrs"
	"time"

	"github.com/OpenIMSDK/tools/log"

	"github.com/OpenIMSDK/protocol/group"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/protocol/wrapperspb"
	"github.com/OpenIMSDK/tools/utils"
)

// // deprecated use CreateGroup
// funcation (g *Group) CreateGroup(ctx context.Context, groupBaseInfo sdk_params_callback.CreateGroupBaseInfoParam, memberList sdk_params_callback.CreateGroupMemberRoleParam) (*sdkws.GroupInfo, error) {
//	req := &group.CreateGroupReq{
//		GroupInfo: &sdkws.GroupInfo{
//			GroupName:    groupBaseInfo.GroupName,
//			Notification: groupBaseInfo.Notification,
//			Introduction: groupBaseInfo.Introduction,
//			FaceURL:      groupBaseInfo.FaceURL,
//			Ex:           groupBaseInfo.Ex,
//			GroupType:    groupBaseInfo.GroupType,
//		},
//	}
//	if groupBaseInfo.NeedVerification != nil {
//		req.GroupInfo.NeedVerification = *groupBaseInfo.NeedVerification
//	}
//	for _, info := range memberList {
//		switch info.RoleLevel {
//		case constant.GroupOrdinaryUsers:
//			req.InitMembers = append(req.InitMembers, info.UserID)
//		case constant.GroupOwner:
//			req.OwnerUserID = info.UserID
//		case constant.GroupAdmin:
//			req.AdminUserIDs = append(req.AdminUserIDs, info.UserID)
//		default:
//			return nil, sdkerrs.ErrArgs.Wrap(fmt.Sprintf("CreateGroup: invalid role level %d", info.RoleLevel))
//		}
//	}
//	return g.CreateGroup(ctx, req)
// }

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
	if err := g.SyncGroups(ctx, resp.GroupInfo.GroupID); err != nil {
		return nil, err
	}
	if err := g.SyncAllGroupMember(ctx, resp.GroupInfo.GroupID); err != nil {
		return nil, err
	}
	return resp.GroupInfo, nil
}

func (g *Group) JoinGroup(ctx context.Context, groupID, reqMsg string, joinSource int32) error {
	if err := util.ApiPost(ctx, constant.JoinGroupRouter, &group.JoinGroupReq{GroupID: groupID, ReqMessage: reqMsg, JoinSource: joinSource, InviterUserID: g.loginUserID}, nil); err != nil {
		return err
	}
	if err := g.SyncSelfGroupApplications(ctx, groupID); err != nil {
		return err
	}
	// if err := g.SyncJoinedGroup(ctx); err != nil {
	// 	return err
	// }
	// if err := g.SyncGroupMember(ctx, groupID); err != nil {
	// 	return err
	// }
	return nil
}

func (g *Group) QuitGroup(ctx context.Context, groupID string) error {
	if err := util.ApiPost(ctx, constant.QuitGroupRouter, &group.QuitGroupReq{GroupID: groupID}, nil); err != nil {
		return err
	}
	if err := g.db.DeleteGroupAllMembers(ctx, groupID); err != nil {
		return err
	}
	if err := g.deleteGroup(ctx, groupID); err != nil {
		return err
	}
	// if err := g.SyncGroupMember(ctx, groupID); err != nil {
	//	return err
	// }
	return nil
}

func (g *Group) DismissGroup(ctx context.Context, groupID string) error {
	if err := util.ApiPost(ctx, constant.DismissGroupRouter, &group.DismissGroupReq{GroupID: groupID}, nil); err != nil {
		return err
	}
	if err := g.deleteGroup(ctx, groupID); err != nil {
		return err
	}
	if err := g.db.DeleteGroupAllMembers(ctx, groupID); err != nil {
		return err
	}
	return nil
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
	if err := g.SyncGroups(ctx, groupID); err != nil {
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
	if err := g.SyncGroups(ctx, groupID); err != nil {
		return err
	}
	if err := g.SyncGroupMembers(ctx, groupID, userID); err != nil {
		return err
	}
	return nil
}

func (g *Group) SetGroupMemberRoleLevel(ctx context.Context, groupID, userID string, roleLevel int) error {
	return g.SetGroupMemberInfo(ctx, &group.SetGroupMemberInfo{GroupID: groupID, UserID: userID, RoleLevel: wrapperspb.Int32(int32(roleLevel))})
}

func (g *Group) SetGroupMemberNickname(ctx context.Context, groupID, userID string, groupMemberNickname string) error {
	return g.SetGroupMemberInfo(ctx, &group.SetGroupMemberInfo{GroupID: groupID, UserID: userID, Nickname: wrapperspb.String(groupMemberNickname)})
}

func (g *Group) SetGroupMemberInfo(ctx context.Context, groupMemberInfo *group.SetGroupMemberInfo) error {
	if err := util.ApiPost(ctx, constant.SetGroupMemberInfoRouter, &group.SetGroupMemberInfoReq{Members: []*group.SetGroupMemberInfo{groupMemberInfo}}, nil); err != nil {
		return err
	}
	return g.SyncGroupMembers(ctx, groupMemberInfo.GroupID, groupMemberInfo.UserID)
}

func (g *Group) GetJoinedGroupList(ctx context.Context) ([]*model_struct.LocalGroup, error) {
	return g.db.GetJoinedGroupListDB(ctx)
}

func (g *Group) GetSpecifiedGroupsInfo(ctx context.Context, groupIDs []string) ([]*model_struct.LocalGroup, error) {
	groupList, err := g.db.GetJoinedGroupListDB(ctx)
	if err != nil {
		return nil, err
	}
	superGroupList, err := g.db.GetJoinedSuperGroupList(ctx)
	if err != nil {
		return nil, err
	}
	groupIDMap := utils.SliceSet(groupIDs)
	groups := append(groupList, superGroupList...)
	res := make([]*model_struct.LocalGroup, 0, len(groupIDs))
	for i, v := range groups {
		if _, ok := groupIDMap[v.GroupID]; ok {
			delete(groupIDMap, v.GroupID)
			res = append(res, groups[i])
		}
	}
	if len(groupIDMap) > 0 {
		groups, err := util.CallApi[group.GetGroupsInfoResp](ctx, constant.GetGroupsInfoRouter, &group.GetGroupsInfoReq{GroupIDs: utils.Keys(groupIDMap)})
		if err != nil {
			log.ZError(ctx, "Call GetGroupsInfoRouter", err)
		}
		if groups != nil && len(groups.GroupInfos) > 0 {
			for i := range groups.GroupInfos {
				groups.GroupInfos[i].MemberCount = 0
			}
			res = append(res, util.Batch(ServerGroupToLocalGroup, groups.GroupInfos)...)
		}
	}
	return res, nil
}

func (g *Group) SearchGroups(ctx context.Context, param sdk_params_callback.SearchGroupsParam) ([]*model_struct.LocalGroup, error) {
	if len(param.KeywordList) == 0 || (!param.IsSearchGroupName && !param.IsSearchGroupID) {
		return nil, sdkerrs.ErrArgs.Wrap("keyword is null or search field all false")
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

func (g *Group) SetGroupVerification(ctx context.Context, groupID string, verification int32) error {
	return g.SetGroupInfo(ctx, &sdkws.GroupInfoForSet{GroupID: groupID, NeedVerification: wrapperspb.Int32(verification)})
}

func (g *Group) SetGroupLookMemberInfo(ctx context.Context, groupID string, rule int32) error {
	return g.SetGroupInfo(ctx, &sdkws.GroupInfoForSet{GroupID: groupID, LookMemberInfo: wrapperspb.Int32(rule)})
}

func (g *Group) SetGroupApplyMemberFriend(ctx context.Context, groupID string, rule int32) error {
	return g.SetGroupInfo(ctx, &sdkws.GroupInfoForSet{GroupID: groupID, ApplyMemberFriend: wrapperspb.Int32(rule)})
}

func (g *Group) SetGroupInfo(ctx context.Context, groupInfo *sdkws.GroupInfoForSet) error {
	if err := util.ApiPost(ctx, constant.SetGroupInfoRouter, &group.SetGroupInfoReq{GroupInfoForSet: groupInfo}, nil); err != nil {
		return err
	}
	return g.SyncGroups(ctx, groupInfo.GroupID)
}

func (g *Group) GetGroupMemberList(ctx context.Context, groupID string, filter, offset, count int32) ([]*model_struct.LocalGroupMember, error) {
	return g.db.GetGroupMemberListSplit(ctx, groupID, filter, int(offset), int(count))
}

func (g *Group) GetGroupMemberOwnerAndAdmin(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error) {
	return g.db.GetGroupMemberOwnerAndAdminDB(ctx, groupID)
}

func (g *Group) GetGroupMemberListByJoinTimeFilter(ctx context.Context, groupID string, offset, count int32, joinTimeBegin, joinTimeEnd int64, userIDs []string) ([]*model_struct.LocalGroupMember, error) {
	if joinTimeEnd == 0 {
		joinTimeEnd = time.Now().UnixMilli()
	}
	return g.db.GetGroupMemberListSplitByJoinTimeFilter(ctx, groupID, int(offset), int(count), joinTimeBegin, joinTimeEnd, userIDs)
}

func (g *Group) GetSpecifiedGroupMembersInfo(ctx context.Context, groupID string, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	return g.db.GetGroupSomeMemberInfo(ctx, groupID, userIDList)
}

func (g *Group) KickGroupMember(ctx context.Context, groupID string, reason string, userIDList []string) error {
	if err := util.ApiPost(ctx, constant.KickGroupMemberRouter, &group.KickGroupMemberReq{GroupID: groupID, KickedUserIDs: userIDList, Reason: reason}, nil); err != nil {
		return err
	}
	return g.SyncGroupMembers(ctx, groupID, userIDList...)
}

func (g *Group) TransferGroupOwner(ctx context.Context, groupID, newOwnerUserID string) error {
	oldOwner, err := g.db.GetGroupMemberOwner(ctx, groupID)
	if err != nil {
		return err
	}
	if err := util.ApiPost(ctx, constant.TransferGroupRouter, &group.TransferGroupOwnerReq{GroupID: groupID, OldOwnerUserID: g.loginUserID, NewOwnerUserID: newOwnerUserID}, nil); err != nil {
		return err
	}
	if err := g.SyncGroups(ctx, groupID); err != nil {
		return err
	}
	if err := g.SyncGroupMembers(ctx, groupID, newOwnerUserID, oldOwner.UserID); err != nil {
		return err
	}
	return nil
}

func (g *Group) InviteUserToGroup(ctx context.Context, groupID, reason string, userIDList []string) error {
	if err := util.ApiPost(ctx, constant.InviteUserToGroupRouter, &group.InviteUserToGroupReq{GroupID: groupID, Reason: reason, InvitedUserIDs: userIDList}, nil); err != nil {
		return err
	}
	if err := g.SyncGroups(ctx, groupID); err != nil {
		return err
	}
	if err := g.SyncGroupMembers(ctx, groupID, userIDList...); err != nil {
		return err
	}
	return nil
}

func (g *Group) GetGroupApplicationListAsRecipient(ctx context.Context) ([]*model_struct.LocalAdminGroupRequest, error) {
	return g.db.GetAdminGroupApplication(ctx)
}

func (g *Group) GetGroupApplicationListAsApplicant(ctx context.Context) ([]*model_struct.LocalGroupRequest, error) {
	return g.db.GetSendGroupApplication(ctx)
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

func (g *Group) SearchGroupMembers(ctx context.Context, searchParam *sdk_params_callback.SearchGroupMembersParam) ([]*model_struct.LocalGroupMember, error) {
	return g.db.SearchGroupMembersDB(ctx, searchParam.KeywordList[0], searchParam.GroupID, searchParam.IsSearchMemberNickname, searchParam.IsSearchUserID, searchParam.Offset, searchParam.Count)
}

func (g *Group) IsJoinGroup(ctx context.Context, groupID string) (bool, error) {
	groupList, err := g.db.GetJoinedGroupListDB(ctx)
	if err != nil {
		return false, err
	}
	for _, localGroup := range groupList {
		if localGroup.GroupID == groupID {
			return true, nil
		}
	}
	superGroupList, err := g.db.GetJoinedSuperGroupList(ctx)
	if err != nil {
		return false, err
	}
	for _, localGroup := range superGroupList {
		if localGroup.GroupID == groupID {
			return true, nil
		}
	}
	return false, nil
}
