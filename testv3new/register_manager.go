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
	AllUserIDs []string
	Secret     string
	IMConfig   sdk_struct.IMConfig
}

const (
	UserRegistered   = "registered"
	UserUnregistered = "unregistered"
	AdminToken       = ""
)

func NewRegisterManager() *RegisterManager {
	imConfig := sdk_struct.IMConfig{
		ApiAddr:    testcore.APIADDR,
		PlatformID: int32(1),
		WsAddr:     testcore.WSADDR,
		DataDir:    "./../",
		LogLevel:   uint32(0),
	}
	return &RegisterManager{nil, testcore.SECRET, imConfig}
}

func (r *RegisterManager) RegisterOne(userID string) error {
	ctx, config := InitContext(userID)
	config.IMConfig = r.IMConfig
	var getAccountCheckReq userPB.AccountCheckReq
	var getAccountCheckResp userPB.AccountCheckResp
	getAccountCheckReq.CheckUserIDs = []string{userID}
	for {
		err := util.ApiPost(ctx, constant.AccountCheck, &getAccountCheckReq, &getAccountCheckResp)
		if err != nil {
			return err
		}
		if len(getAccountCheckResp.Results) == 1 && getAccountCheckResp.Results[0].AccountStatus == UserRegistered {
			log.ZWarn(ctx, "Already registered", err, "userID", userID, "results", getAccountCheckResp.Results)
			r.AllUserIDs = append(r.AllUserIDs, userID)
		} else if len(getAccountCheckResp.Results) == 1 && getAccountCheckResp.Results[0].AccountStatus == UserUnregistered {
			log.ZInfo(ctx, "not registered", "userID", userID, "results", getAccountCheckResp.Results)
			break
		} else {
			log.ZError(ctx, " failed, continue ", err, "checkUserIDs", getAccountCheckReq.CheckUserIDs)
			continue
		}
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
			r.AllUserIDs = append(r.AllUserIDs, userID)
			return nil
		}
	}
}

func (r *RegisterManager) RegisterBatch(userIDs []string) error {
	ctx, config := InitContext(userIDs[0])
	config.IMConfig = r.IMConfig
	var getAccountCheckReq userPB.AccountCheckReq
	var getAccountCheckResp userPB.AccountCheckResp
	getAccountCheckReq.CheckUserIDs = userIDs
	for {
		err := util.ApiPost(ctx, constant.AccountCheck, &getAccountCheckReq, &getAccountCheckResp)
		if err != nil {
			log.ZError(ctx, "ApiPost failed", err)
			return err
		}
		if len(getAccountCheckResp.Results) == 1 && getAccountCheckResp.Results[0].AccountStatus == "registered" {
			log.ZWarn(ctx, "Already registered", err, "userIDs", userIDs, "resp", getAccountCheckResp.Results)
			r.AllUserIDs = append(r.AllUserIDs, userIDs...)
		} else if len(getAccountCheckResp.Results) == 1 && getAccountCheckResp.Results[0].AccountStatus == "unregistered" {
			log.ZInfo(ctx, "not registered", "userIDs", userIDs, "resp", getAccountCheckResp.Results)
			break
		} else {
			log.ZError(ctx, " failed, continue", err, "resp checkUserIDs", getAccountCheckReq.CheckUserIDs)
			continue
		}
	}
	var req userPB.UserRegisterReq
	req.Users = []*sdkws.UserInfo{}
	for _, userID := range userIDs {
		userInfo := &sdkws.UserInfo{
			UserID: userID,
		}
		req.Users = append(req.Users, userInfo)
	}

	req.Secret = r.Secret
	for {
		err := util.ApiPost(ctx, constant.UserRegister, &req, nil)
		if err != nil {
			log.ZError(ctx, "post failed ,continue", err, testcore.REGISTERADDR)
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			log.ZInfo(ctx, "register ok", "add", testcore.REGISTERADDR)
			r.AllUserIDs = append(r.AllUserIDs, userIDs...)
			return nil
		}
	}
}

func (r *RegisterManager) GetTokens(userIDs ...string) []string {
	for i := 0; i < len(userIDs); i++ {
		uid := userIDs[i]
		ctx, config := InitContext(uid)
		config.IMConfig = r.IMConfig
		config.Token = ""
		req := authPB.UserTokenReq{PlatformID: r.IMConfig.PlatformID, UserID: uid, Secret: r.Secret}
		resp := authPB.UserTokenResp{}
		err := util.ApiPost(ctx, constant.GetUsersToken, &req, &resp)
		if err != nil {
			log.ZError(ctx, "ApiPost failed", err, "addr", testcore.TOKENADDR, "req", req)
		}
		config.Token = resp.Token
		log.ZInfo(ctx, "get token success", "token", resp.Token, "expireTimeSeconds", resp.ExpireTimeSeconds)
	}
	return nil
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

func InitContext(uid string) (context.Context, ccontext.GlobalConfig) {
	config := ccontext.GlobalConfig{
		UserID: uid,
		Token:  AdminToken,
	}
	ctx := ccontext.WithInfo(context.Background(), &config)
	ctx = ccontext.WithOperationID(ctx, utils.OperationIDGenerator())
	return ctx, config
}
