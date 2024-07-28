package sdk

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/utils/datautil"
)

// CreateCommonGroup create a regular group. Group members are the users with IDs
// starting from the current user's ID up to the next memberNum users.
func (s *TestSDK) CreateCommonGroup(ctx context.Context, memberNum int) (*sdkws.GroupInfo, error) {
	memberUserIds := utils.NextOffsetUserIDs(s.Num, memberNum-1) // 1 is oneself
	resp, err := s.createGroup(ctx, memberUserIds, vars.CommonGroup)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateLargeGroup create a large group. Group members are all users.
func (s *TestSDK) CreateLargeGroup(ctx context.Context) (*sdkws.GroupInfo, error) {
	memberUserIDs := datautil.Delete(datautil.CopySlice(vars.UserIDs), utils.MustGetUserNum(s.UserID))
	resp, err := s.createGroup(ctx, memberUserIDs, vars.LargeGroup)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *TestSDK) createGroup(ctx context.Context, memberUserIds []string, groupType string) (*sdkws.GroupInfo, error) {
	return s.SDK.Group().CreateGroup(ctx, &group.CreateGroupReq{
		MemberUserIDs: memberUserIds,
		GroupInfo: &sdkws.GroupInfo{
			GroupName: utils.BuildGroupName(s.UserID, groupType),
			GroupType: constant.WorkingGroup,
		},
		AdminUserIDs: nil,
		OwnerUserID:  s.UserID,
	})
}

func (s *TestSDK) GetAllJoinedGroup(ctx context.Context) ([]*sdkws.GroupInfo, int, error) {
	res, err := s.SDK.Group().GetServerJoinGroup(ctx)
	if err != nil {
		return nil, 0, err
	}
	return res, len(res), err
}
