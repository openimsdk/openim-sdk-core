package testv3new

import (
	"context"
	"fmt"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/sdk_struct"
	"time"

	authPB "github.com/OpenIMSDK/protocol/auth"
	"github.com/OpenIMSDK/protocol/group"
	"github.com/OpenIMSDK/protocol/sdkws"
	userPB "github.com/OpenIMSDK/protocol/user"
)

type TestUserManager struct {
	Secret string
}

const (
	UserRegistered   = "registered"
	UserUnregistered = "unregistered"
	Admin            = "openIM123456"
)

func NewTestUserManager(secret string) *TestUserManager {
	return &TestUserManager{Secret: secret}
}

func (t *TestUserManager) NewCtx() context.Context {
	return ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
		UserID: Admin,
		IMConfig: sdk_struct.IMConfig{
			PlatformID: int32(PLATFORMID),
			ApiAddr:    APIADDR,
			WsAddr:     WSADDR,
		}})
}

func (t *TestUserManager) GenUserIDs(num int) (userIDs []string) {
	for i := 0; i < num; i++ {
		userIDs = append(userIDs, fmt.Sprintf("testv3new_%d_%d", time.Now().UnixNano(), i))
	}
	return userIDs
}

func (t *TestUserManager) RegisterUsers(ctx context.Context, userIDs ...string) error {
	var users []*sdkws.UserInfo
	for _, userID := range userIDs {
		users = append(users, &sdkws.UserInfo{UserID: userID, Nickname: userID})
	}
	return util.ApiPost(ctx, constant.UserRegister, &userPB.UserRegisterReq{
		Secret: t.Secret,
		Users:  users,
	}, nil)
}

func (t *TestUserManager) GetTokens(ctx context.Context, userIDs ...string) []string {
	var userTokens []string
	for i := 0; i < len(userIDs); i++ {
		userID := userIDs[i]
		token, _ := t.GetToken(ctx, userID, constant.WindowsPlatformID)
		userTokens = append(userTokens, token)
	}
	return userTokens
}

func (r *TestUserManager) GetToken(ctx context.Context, userID string, platformID int32) (string, error) {
	req := authPB.UserTokenReq{PlatformID: platformID, UserID: userID, Secret: r.Secret}
	resp := authPB.UserTokenResp{}
	err := util.ApiPost(ctx, constant.GetUsersToken, &req, &resp)
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (p *TestUserManager) CreateGroup(ctx context.Context, groupID string, ownerUserID string, userIDs []string, groupName string) error {
	req := &group.CreateGroupReq{
		MemberUserIDs: userIDs,
		OwnerUserID:   ownerUserID,
		GroupInfo: &sdkws.GroupInfo{
			GroupID:   groupID,
			GroupName: groupName,
			GroupType: constant.WorkingGroup,
		},
	}
	resp := &group.CreateGroupResp{}
	err := util.ApiPost(ctx, constant.CreateGroupRouter, &req, &resp)
	if err != nil {
		return err
	}
	return nil
}

func (t *TestUserManager) InviteUserToGroup(ctx context.Context, groupID string, invitedUserIDs []string) error {
	req := &group.InviteUserToGroupReq{
		GroupID:        groupID,
		InvitedUserIDs: invitedUserIDs,
	}
	resp := &group.InviteUserToGroupResp{}
	err := util.ApiPost(ctx, constant.InviteUserToGroupRouter, &req, &resp)
	if err != nil {
		return err
	}
	return nil
}

func (t *TestUserManager) GetGroupMembersInfo(ctx context.Context, groupID string, userIDs []string) (*group.GetGroupMembersInfoResp, error) {
	req := &group.GetGroupMembersInfoReq{
		GroupID: groupID,
		UserIDs: userIDs,
	}
	resp := &group.GetGroupMembersInfoResp{}
	err := util.ApiPost(ctx, constant.GetGroupMembersInfoRouter, &req, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
