package testv3new

import (
	"context"
	authPB "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/auth"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	userPB "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
	"open_im_sdk/internal/util"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/testv3new/testcore"
	_ "open_im_sdk/testv3new/testcore"
	"time"
)

type RegisterManager struct {
}

func NewRegisterManager() *RegisterManager {
	return &RegisterManager{}
}

func (r *RegisterManager) RegisterOne(userID string) error {
	ctx, _ := testcore.InitContext(userID)
	account, err := checkUserAccount(ctx, userID)
	if err != nil {
		return err
	}
	if account != true {
		log.Error(userID, "user:["+userID+"] account register fail, maybe already registered")
		return err
	}
	registerUserAccount(ctx, userID)
	return nil
}

func (r *RegisterManager) RegisterBatch(userIDs []string) error {
	for i := 0; i < len(userIDs); i++ {
		uid := userIDs[i]
		err := r.RegisterOne(uid)
		if err != nil {
			log.Error(utils.OperationIDGenerator(), err)
			return err
		}
	}
	return nil
}

func (r *RegisterManager) GetTokens(userIDs ...string) []string {
	for i := 0; i < len(userIDs); i++ {
		uid := userIDs[i]
		ctx, config := testcore.InitContext(uid)
		config.Token = ""
		req := authPB.UserTokenReq{PlatformID: testcore.PlatformID, UserID: uid, Secret: testcore.Secret}
		resp := authPB.UserTokenResp{}
		err := util.ApiPost(ctx, "/auth/user_token", &req, &resp)
		if err != nil {
			log.Error(req.UserID, "ApiPost failed ", err.Error(), testcore.TOKENADDR, req)
		}
		config.Token = resp.Token
		log.Info(req.UserID, "get token: ", resp.Token, " expireTimeSeconds: ", resp.ExpireTimeSeconds)
	}
	return nil
}

func (p *PressureTester) CreateGroup(groupID string, ownerUserID string, userIDs []string) error {
	ctx, _ := testcore.InitContext(ownerUserID)
	req := &group.CreateGroupReq{
		MemberUserIDs: userIDs,
		AdminUserIDs:  []string{},
		OwnerUserID:   ownerUserID,
		GroupInfo: &sdkws.GroupInfo{
			GroupID:   groupID,
			GroupName: "test-",
			GroupType: 2,
		},
	}
	_, err := open_im_sdk.UserForSDK.Group().CreateGroup(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func checkUserAccount(ctx context.Context, uid string) (bool, error) {
	var getAccountCheckReq userPB.AccountCheckReq
	var getAccountCheckResp userPB.AccountCheckResp
	getAccountCheckReq.CheckUserIDs = []string{uid}
	for {
		err := util.ApiPost(ctx, "/user/account_check", &getAccountCheckReq, &getAccountCheckResp)
		if err != nil {
			return false, err
		}
		if len(getAccountCheckResp.Results) == 1 && getAccountCheckResp.Results[0].AccountStatus == "registered" {
			log.Warn(getAccountCheckReq.CheckUserIDs[0], "Already registered ", uid, getAccountCheckResp.Results)
			testcore.AddUserID(uid)
			return false, nil
		} else if len(getAccountCheckResp.Results) == 1 && getAccountCheckResp.Results[0].AccountStatus == "unregistered" {
			log.Info(getAccountCheckReq.CheckUserIDs[0], "not registered ", uid, getAccountCheckResp.Results)
			break
		} else {
			log.Error(getAccountCheckReq.CheckUserIDs[0], " failed, continue ", err, getAccountCheckReq.CheckUserIDs)
			continue
		}
	}
	return true, nil
}

func registerUserAccount(ctx context.Context, uid string) {
	var req userPB.UserRegisterReq
	req.Users = []*sdkws.UserInfo{{UserID: uid}}
	req.Secret = testcore.Secret
	for {
		err := util.ApiPost(ctx, "/user/user_register", &req, nil)
		if err != nil {
			log.Error("post failed ,continue ", err.Error(), testcore.REGISTERADDR)
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			log.Info("register ok ", testcore.REGISTERADDR)
			testcore.AddUserID(uid)
			return
		}
	}
}
