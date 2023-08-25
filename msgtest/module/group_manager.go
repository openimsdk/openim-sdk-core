package module

import (
	"context"
	"fmt"
	"open_im_sdk/pkg/constant"
	"time"

	"github.com/OpenIMSDK/protocol/group"
	"github.com/OpenIMSDK/protocol/sdkws"
)

type TestGroupManager struct {
	*MetaManager
}

func (t *TestGroupManager) GenGroupID(prefix string) string {
	return fmt.Sprintf("%s_test_group_id_%d", prefix, time.Now().UnixNano())
}

func (t *TestGroupManager) CreateGroup(groupID string, groupName string, ownerUserID string, userIDs []string) error {
	req := &group.CreateGroupReq{
		MemberUserIDs: userIDs,
		OwnerUserID:   ownerUserID,
		GroupInfo: &sdkws.GroupInfo{
			GroupID:       groupID,
			GroupName:     groupName,
			GroupType:     constant.WorkingGroup,
			CreatorUserID: ownerUserID,
		},
	}
	resp := &group.CreateGroupResp{}
	return t.postWithCtx(constant.CreateGroupRouter, &req, &resp)
}

func (t *TestGroupManager) InviteUserToGroup(ctx context.Context, groupID string, invitedUserIDs []string) error {
	req := &group.InviteUserToGroupReq{
		GroupID:        groupID,
		InvitedUserIDs: invitedUserIDs,
	}
	resp := &group.InviteUserToGroupResp{}
	return t.postWithCtx(constant.InviteUserToGroupRouter, &req, &resp)
}
