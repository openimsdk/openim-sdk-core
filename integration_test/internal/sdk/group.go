package sdk

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	sdkUtils "github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"time"
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
	memberUserIDs := datautil.Delete(utils.GenUserIDs(vars.LargeGroupMemberNum), utils.MustGetUserNum(s.UserID))
	resp, err := s.createGroup(ctx, memberUserIDs, vars.LargeGroup)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *TestSDK) createGroup(ctx context.Context, memberUserIds []string, groupType string) (*sdkws.GroupInfo, error) {
	initialMembers := memberUserIds
	if len(memberUserIds) > 1000 {
		initialMembers = memberUserIds[:1000]
	}

	g, err := s.SDK.Group().CreateGroup(ctx, &group.CreateGroupReq{
		MemberUserIDs: initialMembers,
		GroupInfo: &sdkws.GroupInfo{
			GroupName: utils.BuildGroupName(s.UserID, groupType),
			GroupType: constant.WorkingGroup,
		},
		AdminUserIDs: nil,
		OwnerUserID:  s.UserID,
	})
	if err != nil {
		return nil, err
	}

	if len(memberUserIds) > 1000 {
		for i := 1000; i < len(memberUserIds); i += config.ApiParamLength {
			end := i + config.ApiParamLength
			if end > len(memberUserIds) {
				end = len(memberUserIds)
			}
			t := time.Now()
			ctx = ccontext.WithOperationID(ctx, sdkUtils.OperationIDGenerator())
			log.ZWarn(ctx, "InviteUserToGroup begin", nil, "begin", i, "end", end, "groupID", g.GroupID)
			err := s.SDK.Group().InviteUserToGroup(ctx, g.GroupID, "", memberUserIds[i:end])
			if err != nil {
				return nil, err
			}
			log.ZWarn(ctx, "InviteUserToGroup end", nil, "begin", i, "end", end, "groupID", g.GroupID, "cost time", time.Since(t))
			time.Sleep(time.Second)
		}
	}

	return g, nil
}

func (s *TestSDK) GetAllJoinedGroups(ctx context.Context) ([]*sdkws.GroupInfo, int, error) {
	res, err := s.SDK.Group().GetServerJoinGroup(ctx)
	if err != nil {
		return nil, 0, err
	}
	return res, len(res), err
}
