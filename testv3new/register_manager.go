package testv3new

import (
	"context"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"open_im_sdk/testv3new/testcore"
	"time"

	authPB "github.com/OpenIMSDK/protocol/auth"
	"github.com/OpenIMSDK/protocol/group"
	"github.com/OpenIMSDK/protocol/sdkws"
	userPB "github.com/OpenIMSDK/protocol/user"
	"github.com/OpenIMSDK/tools/log"
)

type RegisterManager struct {
	userIDs []string
	Secret  string
}

const (
	UserRegistered   = "registered"
	UserUnregistered = "unregistered"
	Admin            = "openIM123456"
)

func NewRegisterManager() *RegisterManager {
	return &RegisterManager{[]string{}, testcore.SECRET}
}

func (r *RegisterManager) RegisterOne(userID string) error {
	ctx, config := InitContext(userID)
	// checkAccount need admin token
	config.Token, _ = r.GetToken(Admin)
	var getAccountCheckReq userPB.AccountCheckReq
	var getAccountCheckResp userPB.AccountCheckResp
	getAccountCheckReq.CheckUserIDs = []string{userID}
	err := util.ApiPost(ctx, constant.AccountCheck, &getAccountCheckReq, &getAccountCheckResp)
	if err != nil {
		log.ZError(ctx, "ApiPost failed", err)
		return err
	}
	if len(getAccountCheckResp.Results) == 1 && getAccountCheckResp.Results[0].AccountStatus == UserRegistered {
		log.ZError(ctx, "Already registered", err, "userID", userID, "results", getAccountCheckResp.Results)
		r.userIDs = append(r.userIDs, userID)
		return nil
	} else if len(getAccountCheckResp.Results) == 1 && getAccountCheckResp.Results[0].AccountStatus == UserUnregistered {
		log.ZInfo(ctx, "not registered", "userID", userID, "results", getAccountCheckResp.Results)
	} else {
		log.ZError(ctx, "failed", err, "checkUserIDs", getAccountCheckReq.CheckUserIDs)
		return nil
	}

	var req userPB.UserRegisterReq
	req.Users = []*sdkws.UserInfo{{UserID: userID}}
	req.Secret = r.Secret
	for {
		err := util.ApiPost(ctx, constant.UserRegister, &req, nil)
		if err != nil {
			log.ZError(ctx, "post failed ,continue", err, testcore.REGISTERADDR)
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			log.ZInfo(ctx, "register ok", "addr", testcore.REGISTERADDR)
			r.userIDs = append(r.userIDs, userID)
			return nil
		}
	}
}

func (r *RegisterManager) RegisterBatch(userIDs []string) error {
	for _, userID := range userIDs {
		err := r.RegisterOne(userID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RegisterManager) GetTokens(userIDs ...string) []string {
	var userTokens []string
	for i := 0; i < len(userIDs); i++ {
		userID := userIDs[i]
		token, _ := r.GetToken(userID)
		userTokens = append(userTokens, token)
	}
	return userTokens
}

func (r *RegisterManager) GetToken(userID string) (string, error) {
	ctx, config := InitContext(userID)
	req := authPB.UserTokenReq{PlatformID: config.PlatformID, UserID: userID, Secret: r.Secret}
	resp := authPB.UserTokenResp{}
	err := util.ApiPost(ctx, constant.GetUsersToken, &req, &resp)
	if err != nil {
		log.ZError(ctx, "ApiPost failed", err, "addr", testcore.TOKENADDR, "req", req)
		return "", err
	}
	config.Token = resp.Token
	log.ZInfo(ctx, "get token success", "token", resp.Token, "expireTimeSeconds", resp.ExpireTimeSeconds)
	return resp.Token, nil
}

func (p *PressureTester) CreateGroup(groupID string, ownerUserID string, userIDs []string, groupName string) error {
	token, err := p.registerManager.GetToken(ownerUserID)
	if err != nil {
		log.ZError(context.Background(), "get token error", err, "userID", ownerUserID)
		return err
	}
	ctx, _ := InitContextWithToken(ownerUserID, token)
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
	err = util.ApiPost(ctx, constant.CreateGroupRouter, &req, &resp)
	if err != nil {
		log.ZError(ctx, "ApiPost failed ", err, "addr", testcore.TOKENADDR, "req", req)
		return err
	}
	log.ZInfo(ctx, "create group success", "groupID", groupID, "ownerUserID", ownerUserID)
	return nil
}

func (p *PressureTester) InviteUserToGroup(groupID string, invitedUserIDs []string) error {
	ctx, config := InitContext(groupID)
	config.Token, _ = p.registerManager.GetToken(Admin)
	req := &group.InviteUserToGroupReq{
		GroupID:        groupID,
		Reason:         "",
		InvitedUserIDs: invitedUserIDs,
	}
	resp := &group.InviteUserToGroupResp{}
	err := util.ApiPost(ctx, constant.InviteUserToGroupRouter, &req, &resp)
	if err != nil {
		log.ZError(ctx, "ApiPost failed ", err, "addr", testcore.TOKENADDR, "req", req)
		return err
	}
	log.ZInfo(ctx, "create group success", "groupID", groupID, "invitedUserIDs", invitedUserIDs)
	return nil
}

func (p *PressureTester) GetGroupMembersInfo(groupID string, userIDs []string) (*group.GetGroupMembersInfoResp, error) {
	ctx, config := InitContext(groupID)
	config.Token, _ = p.registerManager.GetToken(Admin)
	req := &group.GetGroupMembersInfoReq{
		GroupID: groupID,
		UserIDs: userIDs,
	}
	resp := &group.GetGroupMembersInfoResp{}
	err := util.ApiPost(ctx, constant.GetGroupMembersInfoRouter, &req, &resp)
	if err != nil {
		log.ZError(ctx, "ApiPost failed ", err, "addr", testcore.TOKENADDR, "req", req)
		return nil, err
	}
	log.ZInfo(ctx, "create group success", "groupID", groupID, "userIDs", userIDs)
	return resp, nil
}

func InitContext(uid string) (context.Context, *ccontext.GlobalConfig) {
	config := ccontext.GlobalConfig{
		UserID: uid, Token: "",
		IMConfig: sdk_struct.IMConfig{
			PlatformID:          constant.AndroidPlatformID,
			ApiAddr:             testcore.APIADDR,
			WsAddr:              testcore.WSADDR,
			LogLevel:            2,
			IsLogStandardOutput: true,
			LogFilePath:         "./",
		}}
	ctx := ccontext.WithInfo(context.Background(), &config)
	ctx = ccontext.WithOperationID(ctx, utils.OperationIDGenerator())
	return ctx, &config
}

func InitContextWithToken(uid string, token string) (context.Context, *ccontext.GlobalConfig) {
	config := ccontext.GlobalConfig{
		UserID: uid, Token: token,
		IMConfig: sdk_struct.IMConfig{
			PlatformID:          constant.AndroidPlatformID,
			ApiAddr:             testcore.APIADDR,
			WsAddr:              testcore.WSADDR,
			LogLevel:            2,
			IsLogStandardOutput: true,
			LogFilePath:         "./",
		}}
	ctx := ccontext.WithInfo(context.Background(), &config)
	ctx = ccontext.WithOperationID(ctx, utils.OperationIDGenerator())
	return ctx, &config
}
