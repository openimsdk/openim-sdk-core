package manager

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/sdk_user_simulator"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/sdk"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/protocol/sdkws"
	userPB "github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/log"
	"golang.org/x/sync/errgroup"
	"time"
)

type TestUserManager struct {
	*MetaManager
}

func NewUserManager(m *MetaManager) *TestUserManager {
	return &TestUserManager{m}
}

func (t *TestUserManager) GenUserIDs() []string {
	ids := make([]string, vars.UserNum)
	for i := 0; i < vars.UserNum; i++ {
		ids[i] = utils.GetUserID(i)
	}
	vars.UserIDs = ids
	vars.SuperUserIDs = ids[:vars.SuperUserNum]
	return ids
}

func (t *TestUserManager) RegisterUsers(ctx context.Context, userIDs ...string) error {
	tm := time.Now()
	log.ZDebug(ctx, "register begin", "len userIDs", len(userIDs))
	defer func() {
		log.ZDebug(ctx, "register end", "time consuming", time.Since(tm))
	}()

	var users []*sdkws.UserInfo
	for _, userID := range userIDs {
		users = append(users, &sdkws.UserInfo{UserID: userID, Nickname: userID})
	}
	return t.PostWithCtx(constant.UserRegister, &userPB.UserRegisterReq{
		Secret: t.GetSecret(),
		Users:  users,
	}, nil)
}

func (t *TestUserManager) InitSDKAndLogin(ctx context.Context, userIDs ...string) error {
	tm := time.Now()
	log.ZDebug(ctx, "initSDKAndLogin begin", "len userIDs", len(userIDs))
	defer func() {
		log.ZDebug(ctx, "initSDKAndLogin end", "time consuming", time.Since(tm))
	}()

	gr, ctx := errgroup.WithContext(ctx)
	gr.SetLimit(vars.ErrGroupCommonLimit)
	for _, userID := range userIDs {
		userID := userID
		gr.Go(func() error {
			token, err := t.GetToken(userID, vars.PlatformID)
			if err != nil {
				return err
			}
			mgr, err := sdk_user_simulator.InitSDKAndLogin(userID, token, t.IMConfig)
			if err != nil {
				return err
			}
			userNum := utils.MustGetUserNum(userID)
			sdk.TestSDKs[userNum] = sdk.NewTestSDK(userID, userNum, mgr) // init sdk
			return nil
		})
	}
	if err := gr.Wait(); err != nil {
		return err
	}
	return nil
}
