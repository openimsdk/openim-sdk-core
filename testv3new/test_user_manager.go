package testv3new

import (
	"context"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
	"time"

	authPB "github.com/OpenIMSDK/protocol/auth"
	"github.com/OpenIMSDK/protocol/group"
	"github.com/OpenIMSDK/protocol/sdkws"
	userPB "github.com/OpenIMSDK/protocol/user"
	"github.com/OpenIMSDK/tools/log"
)

type TestUserManager struct {
	userIDs []string
	Secret  string
}

const (
	UserRegistered   = "registered"
	UserUnregistered = "unregistered"
	Admin            = "openIM123456"
)

func NewRegisterManager() *TestUserManager {
	return &TestUserManager{[]string{}, SECRET}
}

func (t *TestUserManager) RegisterOne(ctx context.Context, userID string) error {
	// checkAccount need admin token
	var getAccountCheckReq userPB.AccountCheckReq
	var getAccountCheckResp userPB.AccountCheckResp
	getAccountCheckReq.CheckUserIDs = []string{userID}
	err := util.ApiPost(ctx, constant.AccountCheck, &getAccountCheckReq, &getAccountCheckResp)
	if err != nil {
		return err
	}
	if len(getAccountCheckResp.Results) == 1 && getAccountCheckResp.Results[0].AccountStatus == UserRegistered {
		t.userIDs = append(t.userIDs, userID)
		return nil
	} else if len(getAccountCheckResp.Results) == 1 && getAccountCheckResp.Results[0].AccountStatus == UserUnregistered {
		log.ZInfo(ctx, "not registered", "userID", userID, "results", getAccountCheckResp.Results)
	} else {
		log.ZError(ctx, "failed", err, "checkUserIDs", getAccountCheckReq.CheckUserIDs)
		return nil
	}

	var req userPB.UserRegisterReq
	req.Users = []*sdkws.UserInfo{{UserID: userID}}
	req.Secret = t.Secret
	for {
		err := util.ApiPost(ctx, constant.UserRegister, &req, nil)
		if err != nil {
			log.ZError(ctx, "post failed ,continue", err, REGISTERADDR)
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			log.ZInfo(ctx, "register ok", "addr", REGISTERADDR)
			t.userIDs = append(t.userIDs, userID)
			return nil
		}
	}
}

func (t *TestUserManager) RegisterBatch(ctx context.Context, userIDs []string) error {
	for _, userID := range userIDs {
		err := t.RegisterOne(ctx, userID)
		if err != nil {
			return err
		}
	}
	return nil
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
		AdminUserIDs:  []string{},
		OwnerUserID:   ownerUserID,
		GroupInfo: &sdkws.GroupInfo{
			GroupID:   groupID,
			GroupName: groupName,
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
		Reason:         "",
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

// func InitContext(uid string) (context.Context, *ccontext.GlobalConfig) {
// 	config := ccontext.GlobalConfig{
// 		UserID: uid, Token: "",
// 		IMConfig: sdk_struct.IMConfig{
// 			PlatformID:          int32(PLATFORMID),
// 			ApiAddr:             APIADDR,
// 			WsAddr:              WSADDR,
// 			LogLevel:            2,
// 			IsLogStandardOutput: true,
// 			LogFilePath:         "./",
// 		}}
// 	ctx := ccontext.WithInfo(context.Background(), &config)
// 	ctx = ccontext.WithOperationID(ctx, utils.OperationIDGenerator())
// 	return ctx, &config
// }
