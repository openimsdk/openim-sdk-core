package group

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/wrapperspb"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/sdk_params_callback"
	"time"
)

// deprecated use CreateGroupV2
func (g *Group) CreateGroup(ctx context.Context, groupBaseInfo sdk_params_callback.CreateGroupBaseInfoParam, memberList sdk_params_callback.CreateGroupMemberRoleParam) (*sdkws.GroupInfo, error) {
	req := &group.CreateGroupReq{
		GroupInfo: &sdkws.GroupInfo{
			GroupName:    groupBaseInfo.GroupName,
			Notification: groupBaseInfo.Notification,
			Introduction: groupBaseInfo.Introduction,
			FaceURL:      groupBaseInfo.FaceURL,
			Ex:           groupBaseInfo.Ex,
			GroupType:    groupBaseInfo.GroupType,
		},
	}
	if groupBaseInfo.NeedVerification != nil {
		req.GroupInfo.NeedVerification = *groupBaseInfo.NeedVerification
	}
	for _, info := range memberList {
		switch info.RoleLevel {
		case constant.GroupOrdinaryUsers:
			req.InitMembers = append(req.InitMembers, info.UserID)
		case constant.GroupOwner:
			req.OwnerUserID = info.UserID
		case constant.GroupAdmin:
			req.AdminUserIDs = append(req.AdminUserIDs, info.UserID)
		default:
			return nil, errs.ErrArgs.Wrap(fmt.Sprintf("CreateGroupV2: invalid role level %d", info.RoleLevel))
		}
	}
	return g.CreateGroupV2(ctx, req)
}

func (g *Group) CreateGroupV2(ctx context.Context, req *group.CreateGroupReq) (*sdkws.GroupInfo, error) {
	resp, err := util.CallApi[group.CreateGroupResp](ctx, constant.CreateGroupRouter, req)
	if err != nil {
		return nil, err
	}
	if err := g.SyncJoinedGroup(ctx); err != nil {
		return nil, err
	}
	if err := g.SyncGroupMember(ctx, resp.GroupInfo.GroupID); err != nil {
		return nil, err
	}
	return resp.GroupInfo, nil
}

func (g *Group) JoinGroup(ctx context.Context, groupID, reqMsg string, joinSource int32) error {
	if err := util.ApiPost(ctx, constant.JoinGroupRouter, &group.JoinGroupReq{GroupID: groupID, ReqMessage: reqMsg, JoinSource: joinSource}, nil); err != nil {
		return err
	}
	if err := g.SyncSelfGroupApplication(ctx); err != nil {
		return err
	}
	if err := g.SyncJoinedGroup(ctx); err != nil {
		return err
	}
	if err := g.SyncGroupMember(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) QuitGroup(ctx context.Context, groupID string) error {
	if err := util.ApiPost(ctx, constant.QuitGroupRouter, &group.QuitGroupReq{GroupID: groupID}, nil); err != nil {
		return err
	}
	if err := g.SyncJoinedGroup(ctx); err != nil {
		return err
	}
	if err := g.SyncGroupMember(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) DismissGroup(ctx context.Context, groupID string) error {
	if err := util.ApiPost(ctx, constant.DismissGroupRouter, &group.DismissGroupReq{GroupID: groupID}, nil); err != nil {
		return err
	}
	if err := g.SyncJoinedGroup(ctx); err != nil {
		return err
	}
	if err := g.SyncGroupMember(ctx, groupID); err != nil {
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
	if err := g.SyncJoinedGroup(ctx); err != nil {
		return err
	}
	return nil
}

func (g *Group) ChangeGroupMemberMute(ctx context.Context, groupID, userID string, mutedSeconds uint32) (err error) {
	if mutedSeconds == 0 {
		err = util.ApiPost(ctx, constant.CancelMuteGroupMemberRouter, &group.CancelMuteGroupMemberReq{GroupID: groupID, UserID: userID}, nil)
	} else {
		err = util.ApiPost(ctx, constant.MuteGroupMemberRouter, &group.MuteGroupMemberReq{GroupID: groupID, UserID: userID, MutedSeconds: mutedSeconds}, nil)
	}
	if err != nil {
		return err
	}
	if err := g.SyncJoinedGroup(ctx); err != nil {
		return err
	}
	if err := g.SyncGroupMember(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) SetGroupMemberRoleLevel(ctx context.Context, groupID, userID string, roleLevel int) error {
	return g.SetGroupMemberInfo(ctx, &group.SetGroupMemberInfo{GroupID: groupID, UserID: userID, RoleLevel: wrapperspb.Int32(int32(roleLevel))})
}

func (g *Group) SetGroupMemberNickname(ctx context.Context, groupID, userID string, groupMemberNickname string, operationID string) error {
	return g.SetGroupMemberInfo(ctx, &group.SetGroupMemberInfo{GroupID: groupID, UserID: userID, Nickname: wrapperspb.String(groupMemberNickname)})
}

func (g *Group) SetGroupMemberInfo(ctx context.Context, groupMemberInfo *group.SetGroupMemberInfo) error {
	if err := util.ApiPost(ctx, constant.SetGroupMemberInfoRouter, &group.SetGroupMemberInfoReq{Members: []*group.SetGroupMemberInfo{groupMemberInfo}}, nil); err != nil {
		return err
	}
	return g.SyncGroupMember(ctx, groupMemberInfo.GroupID)
}

func (g *Group) GetJoinedGroupList(ctx context.Context) ([]*model_struct.LocalGroup, error) {
	return g.db.GetJoinedGroupListDB(ctx)
}

func (g *Group) GetGroupsInfo(ctx context.Context, groupIDList []string) ([]*model_struct.LocalGroup, error) {
	groupList, err := g.db.GetJoinedGroupListDB(ctx)
	if err != nil {
		return nil, err
	}
	superGroupList, err := g.db.GetJoinedSuperGroupList(ctx)
	if err != nil {
		return nil, err
	}
	groupIDMap := utils.SliceSet(groupIDList)
	groups := append(groupList, superGroupList...)
	res := make([]*model_struct.LocalGroup, 0, len(groupIDList))
	for i, v := range groups {
		if _, ok := groupIDMap[v.GroupID]; ok {
			res = append(res, groups[i])
		}
	}
	return res, nil
}

func (g *Group) SearchGroups(ctx context.Context, param sdk_params_callback.SearchGroupsParam) ([]*model_struct.LocalGroup, error) {
	if len(param.KeywordList) == 0 || (!param.IsSearchGroupName && !param.IsSearchGroupID) {
		return nil, errs.NewCodeError(201, "keyword is null or search field all false")
	}
	groups, err := g.db.GetAllGroupInfoByGroupIDOrGroupName(ctx, param.KeywordList[0], param.IsSearchGroupID, param.IsSearchGroupName) // todo	param.KeywordList[0]
	if err != nil {
		return nil, err
	}
	return groups, nil
}
func (g *Group) SetGroupInfo(ctx context.Context, groupInfo *sdk_params_callback.SetGroupInfoParam, groupID string) error {
	return g.SetGroupInfoV2(ctx, &sdkws.GroupInfoForSet{
		GroupID:          groupID,
		GroupName:        groupInfo.GroupName,
		Notification:     groupInfo.Notification,
		Introduction:     groupInfo.Introduction,
		FaceURL:          groupInfo.FaceURL,
		Ex:               groupInfo.Ex,
		NeedVerification: wrapperspb.Int32Ptr(groupInfo.NeedVerification),
	})
}

func (g *Group) SetGroupVerification(ctx context.Context, verification int32, groupID string) error {
	return g.SetGroupInfoV2(ctx, &sdkws.GroupInfoForSet{GroupID: groupID, NeedVerification: wrapperspb.Int32(verification)})
}

func (g *Group) SetGroupLookMemberInfo(ctx context.Context, rule int32, groupID string) error {
	return g.SetGroupInfoV2(ctx, &sdkws.GroupInfoForSet{GroupID: groupID, LookMemberInfo: wrapperspb.Int32(rule)})
}

func (g *Group) SetGroupApplyMemberFriend(ctx context.Context, rule int32, groupID string) error {
	return g.SetGroupInfoV2(ctx, &sdkws.GroupInfoForSet{GroupID: groupID, ApplyMemberFriend: wrapperspb.Int32(rule)})
}

func (g *Group) SetGroupInfoV2(ctx context.Context, groupInfo *sdkws.GroupInfoForSet) error {
	if err := util.ApiPost(ctx, constant.SetGroupInfoRouter, &group.SetGroupInfoReq{GroupInfoForSet: groupInfo}, nil); err != nil {
		return err
	}
	return g.SyncJoinedGroup(ctx)
}

func (g *Group) GetGroupMemberList(ctx context.Context, groupID string, filter, offset, count int32) ([]*model_struct.LocalGroupMember, error) {
	return g.db.GetGroupMemberListSplit(ctx, groupID, filter, int(offset), int(count))
}

func (g *Group) GetGroupMemberOwnerAndAdmin(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error) {
	return g.db.GetGroupMemberOwnerAndAdmin(ctx, groupID)
}

func (g *Group) GetGroupMemberListByJoinTimeFilter(ctx context.Context, groupID string, offset, count int32, joinTimeBegin, joinTimeEnd int64, userIDs []string) ([]*model_struct.LocalGroupMember, error) {
	if joinTimeEnd == 0 {
		joinTimeEnd = time.Now().UnixMilli()
	}
	return g.db.GetGroupMemberListSplitByJoinTimeFilter(ctx, groupID, int(offset), int(count), joinTimeBegin, joinTimeEnd, userIDs)
}

func (g *Group) GetGroupMembersInfo(ctx context.Context, groupID string, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	return g.db.GetGroupSomeMemberInfo(ctx, groupID, userIDList)
}

func (g *Group) KickGroupMember(ctx context.Context, groupID string, reason string, userIDList []string) error {
	if err := util.ApiPost(ctx, constant.KickGroupMemberRouter, &group.KickGroupMemberReq{GroupID: groupID, KickedUserIDs: userIDList, Reason: reason}, nil); err != nil {
		return err
	}
	return g.SyncGroupMember(ctx, groupID)
}

func (g *Group) TransferGroupOwner(ctx context.Context, groupID, newOwnerUserID string) error {
	if err := util.ApiPost(ctx, constant.TransferGroupRouter, &group.TransferGroupOwnerReq{GroupID: groupID, NewOwnerUserID: newOwnerUserID}, nil); err != nil {
		return err
	}
	if err := g.SyncJoinedGroup(ctx); err != nil {
		return err
	}
	if err := g.SyncGroupMember(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) InviteUserToGroup(ctx context.Context, groupID, reason string, userIDList []string) error {
	if err := util.ApiPost(ctx, constant.InviteUserToGroupRouter, &group.InviteUserToGroupReq{GroupID: groupID, Reason: reason, InvitedUserIDs: userIDList}, nil); err != nil {
		return nil
	}
	if err := g.SyncJoinedGroup(ctx); err != nil {
		return err
	}
	if err := g.SyncGroupMember(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) GetRecvGroupApplicationList(ctx context.Context) ([]*model_struct.LocalAdminGroupRequest, error) {
	return g.db.GetAdminGroupApplication(ctx)
}

func (g *Group) GetSendGroupApplicationList(ctx context.Context) ([]*model_struct.LocalGroupRequest, error) {
	return g.db.GetSendGroupApplication(ctx)
}

func (g *Group) AcceptGroupApplication(ctx context.Context, groupID, fromUserID, handleMsg string) error {
	return g.HandlerGroupApplication(ctx, &group.GroupApplicationResponseReq{GroupID: groupID, FromUserID: fromUserID, HandledMsg: handleMsg, HandleResult: 1})
}

func (g *Group) RefuseGroupApplication(ctx context.Context, groupID, fromUserID, handleMsg string) error {
	return g.HandlerGroupApplication(ctx, &group.GroupApplicationResponseReq{GroupID: groupID, FromUserID: fromUserID, HandledMsg: handleMsg, HandleResult: 1})
}

func (g *Group) HandlerGroupApplication(ctx context.Context, req *group.GroupApplicationResponseReq) error {
	if err := util.ApiPost(ctx, constant.AcceptGroupApplicationRouter, req, nil); err != nil {
		return err
	}
	// SyncAdminGroupApplication todo
	return nil
}

func (g *Group) SearchGroupMembers(ctx context.Context, searchParam sdk_params_callback.SearchGroupMembersParam) ([]*model_struct.LocalGroupMember, error) {
	return g.db.SearchGroupMembersDB(ctx, searchParam.KeywordList[0], searchParam.GroupID, searchParam.IsSearchMemberNickname, searchParam.IsSearchUserID, searchParam.Offset, searchParam.Count)
}
