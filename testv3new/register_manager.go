package testv3new

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	authPB "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/auth"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	userPB "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"open_im_sdk/testv3new/testcore"
	"sync"
	"time"
)

var userLock sync.RWMutex

type RegisterManager struct {
	AllUserIDs   []string
	Secret       string
	IMConfig     sdk_struct.IMConfig
	GlobalConfig ccontext.GlobalConfig
	ctx          context.Context
}

func NewRegisterManager() *RegisterManager {
	imConfig := sdk_struct.IMConfig{
		ApiAddr:             testcore.APIADDR,
		WsAddr:              testcore.WSADDR,
		PlatformID:          int32(1),
		DataDir:             "./../",
		LogLevel:            uint32(1),
		IsLogStandardOutput: true,
	}
	globalConfig := ccontext.GlobalConfig{
		IMConfig: imConfig,
	}
	ctx := ccontext.WithInfo(context.Background(), &globalConfig)
	ctx = ccontext.WithOperationID(ctx, utils.OperationIDGenerator())
	return &RegisterManager{nil, "tuoyun", imConfig, globalConfig, ctx}
}

func (r *RegisterManager) RegisterOne(userID string) error {
	r.GlobalConfig.UserID = userID
	var getAccountCheckReq userPB.AccountCheckReq
	var getAccountCheckResp userPB.AccountCheckResp
	getAccountCheckReq.CheckUserIDs = []string{userID}
	for {
		err := util.ApiPost(r.ctx, constant.AccountCheck, &getAccountCheckReq, &getAccountCheckResp)
		if err != nil {
			return err
		}
		if len(getAccountCheckResp.Results) == 1 && getAccountCheckResp.Results[0].AccountStatus == "registered" {
			log.ZWarn(r.ctx, "Already registered", err, userID, getAccountCheckResp.Results)
			userLock.Lock()
			r.AllUserIDs = append(r.AllUserIDs, userID)
			userLock.Unlock()
		} else if len(getAccountCheckResp.Results) == 1 && getAccountCheckResp.Results[0].AccountStatus == "unregistered" {
			log.ZInfo(r.ctx, "not registered", userID, getAccountCheckResp.Results)
			break
		} else {
			log.ZError(r.ctx, " failed, continue ", err, getAccountCheckReq.CheckUserIDs)
			continue
		}
	}

	var req userPB.UserRegisterReq
	req.Users = []*sdkws.UserInfo{{UserID: userID}}
	req.Secret = r.Secret
	for {
		err := util.ApiPost(r.ctx, constant.UserRegister, &req, nil)
		if err != nil {
			log.ZError(r.ctx, "post failed ,continue", err, testcore.REGISTERADDR)
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			log.ZInfo(r.ctx, "register ok", "addr", testcore.REGISTERADDR)

			r.AllUserIDs = append(r.AllUserIDs, userID)
			return nil
		}
	}
}

func (r *RegisterManager) RegisterBatch(userIDs []string) error {
	r.GlobalConfig.UserID = userIDs[0]
	var getAccountCheckReq userPB.AccountCheckReq
	var getAccountCheckResp userPB.AccountCheckResp
	getAccountCheckReq.CheckUserIDs = userIDs
	for {
		err := util.ApiPost(r.ctx, constant.AccountCheck, &getAccountCheckReq, &getAccountCheckResp)
		if err != nil {
			return err
		}
		if len(getAccountCheckResp.Results) == 1 && getAccountCheckResp.Results[0].AccountStatus == "registered" {
			log.ZWarn(r.ctx, "Already registered", err, userIDs, getAccountCheckResp.Results)
			userLock.Lock()
			r.AllUserIDs = append(r.AllUserIDs, userIDs...)
			userLock.Unlock()
		} else if len(getAccountCheckResp.Results) == 1 && getAccountCheckResp.Results[0].AccountStatus == "unregistered" {
			log.ZInfo(r.ctx, "not registered", userIDs, getAccountCheckResp.Results)
			break
		} else {
			log.ZError(r.ctx, " failed, continue", err, getAccountCheckReq.CheckUserIDs)
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
		err := util.ApiPost(r.ctx, constant.UserRegister, &req, nil)
		if err != nil {
			log.ZError(r.ctx, "post failed ,continue", err, testcore.REGISTERADDR)
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			log.ZInfo(r.ctx, "register ok", "add", testcore.REGISTERADDR)
			r.AllUserIDs = append(r.AllUserIDs, userIDs...)
			return nil
		}
	}
}

func (r *RegisterManager) GetTokens(userIDs ...string) []string {
	for i := 0; i < len(userIDs); i++ {
		uid := userIDs[i]
		r.GlobalConfig.Token = ""
		req := authPB.UserTokenReq{PlatformID: r.IMConfig.PlatformID, UserID: uid, Secret: r.Secret}
		resp := authPB.UserTokenResp{}
		err := util.ApiPost(r.ctx, constant.GetUsersToken, &req, &resp)
		if err != nil {
			log.ZError(r.ctx, "ApiPost failed", err, "addr", testcore.TOKENADDR, "req", req)
		}
		r.GlobalConfig.Token = resp.Token
		log.ZInfo(r.ctx, "get token success", "token", resp.Token, "expireTimeSeconds", resp.ExpireTimeSeconds)
	}
	return nil
}

func (p *PressureTester) CreateGroup(groupID string, ownerUserID string, userIDs []string, groupName string) error {
	ctx := ccontext.WithOperationID(context.Background(), utils.OperationIDGenerator())
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
	}
	log.ZInfo(ctx, "create group success", "groupID", groupID, "ownerUserID", ownerUserID)
	return nil
}
