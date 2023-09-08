package module

import (
	"context"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
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
	const batch = 2000
	var memberUserIDs []string
	if len(userIDs) > batch {
		memberUserIDs = userIDs[:batch]
	} else {
		memberUserIDs = userIDs
	}
	req := &group.CreateGroupReq{
		MemberUserIDs: memberUserIDs,
		OwnerUserID:   ownerUserID,
		GroupInfo: &sdkws.GroupInfo{
			GroupID:       groupID,
			GroupName:     groupName,
			GroupType:     constant.WorkingGroup,
			CreatorUserID: ownerUserID,
		},
	}
	resp := &group.CreateGroupResp{}
	if err := t.postWithCtx(constant.CreateGroupRouter, &req, &resp); err != nil {
		return err
	}
	if len(userIDs) > batch {
		num := len(userIDs) / batch
		if len(userIDs)%batch != 0 {
			num++
		}
		for i := 1; i < num; i++ {
			start := batch * i
			end := batch*i + batch
			if len(userIDs) < end {
				end = len(userIDs)
			}
			req := map[string]any{
				"groupID":        groupID,
				"invitedUserIDs": userIDs[start:end],
				"reason":         "test",
			}
			resp := struct{}{}
			if err := t.postWithCtx(constant.RouterGroup+"/invite_user_to_group", req, &resp); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *TestGroupManager) InviteUserToGroup(ctx context.Context, groupID string, invitedUserIDs []string) error {
	req := &group.InviteUserToGroupReq{
		GroupID:        groupID,
		InvitedUserIDs: invitedUserIDs,
	}
	resp := &group.InviteUserToGroupResp{}
	return t.postWithCtx(constant.InviteUserToGroupRouter, &req, &resp)
}
