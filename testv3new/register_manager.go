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

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	authPB "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/auth"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	userPB "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
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
	ctx, _ := InitContext(ownerUserID)
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
		log.ZError(ctx, "ApiPost failed ", err, "addr", testcore.TOKENADDR, "req", req)
		return err
	}
	log.ZInfo(ctx, "create group success", "groupID", groupID, "ownerUserID", ownerUserID)
	return nil
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
