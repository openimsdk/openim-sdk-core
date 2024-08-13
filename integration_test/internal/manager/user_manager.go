package manager

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/decorator"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/reerrgroup"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/sdk_user_simulator"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/sdk"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/protocol/sdkws"
	userPB "github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"sync/atomic"
)

type TestUserManager struct {
	*MetaManager
}

func NewUserManager(m *MetaManager) *TestUserManager {
	return &TestUserManager{m}
}

func (t *TestUserManager) GenAllUserIDs() []string {
	ids := make([]string, vars.UserNum)
	for i := 0; i < vars.UserNum; i++ {
		ids[i] = utils.GetUserID(i)
	}
	vars.UserIDs = ids
	vars.SuperUserIDs = ids[:vars.SuperUserNum]
	return ids
}

func (t *TestUserManager) RegisterAllUsers(ctx context.Context) error {
	return t.registerUsers(ctx, vars.UserIDs...)
}

func (t *TestUserManager) registerUsers(ctx context.Context, userIDs ...string) error {
	defer decorator.FuncLog(ctx)()

	var users []*sdkws.UserInfo
	for _, userID := range userIDs {
		users = append(users, &sdkws.UserInfo{UserID: userID, Nickname: userID})
	}
	if err := t.PostWithCtx(constant.UserRegister, &userPB.UserRegisterReq{
		Secret: t.GetSecret(),
		Users:  users,
	}, nil); err != nil {
		return err
	}
	return nil
}

func (t *TestUserManager) InitAllSDK(ctx context.Context) error {
	return t.initSDK(ctx, vars.UserIDs...)
}

func (t *TestUserManager) initSDK(ctx context.Context, userIDs ...string) error {
	defer decorator.FuncLog(ctx)()

	gr, cctx := reerrgroup.WithContext(ctx, config.ErrGroupCommonLimit)

	var (
		total    atomic.Int64
		progress atomic.Int64
	)
	total.Add(int64(len(userIDs)))
	utils.FuncProgressBarPrint(cctx, gr, &progress, &total)

	for _, userID := range userIDs {
		userID := userID
		gr.Go(func() error {
			userNum := utils.MustGetUserNum(userID)
			token, err := t.GetToken(userID, config.PlatformID)
			if err != nil {
				return err
			}
			ctx, mgr, err := sdk_user_simulator.InitSDK(ctx, userID, token, t.IMConfig)
			if err != nil {
				return err
			}
			sdk.TestSDKs[userNum] = sdk.NewTestSDK(userID, userNum, mgr) // init sdk
			vars.Contexts[userNum] = ctx                                 // init ctx
			log.ZDebug(ctx, "init sdk", "operationID", mcontext.GetOperationID(ctx), "op userID", userID)
			return nil
		})
	}
	if err := gr.Wait(); err != nil {
		return err
	}
	return nil
}

func (t *TestUserManager) LoginAllUsers(ctx context.Context) error {
	return t.login(ctx, vars.UserIDs...)
}

func (t *TestUserManager) LoginByRate(ctx context.Context) error {
	userIDs := vars.UserIDs[:vars.LoginEndUserNum]
	return t.login(ctx, userIDs...)
}

func (t *TestUserManager) LoginLastUsers(ctx context.Context) error {
	userIDs := vars.UserIDs[vars.LoginEndUserNum:]
	return t.login(ctx, userIDs...)
}

func (t *TestUserManager) login(ctx context.Context, userIDs ...string) error {
	defer decorator.FuncLog(ctx)()

	log.ZDebug(ctx, "login users", "len", len(userIDs))

	gr, cctx := reerrgroup.WithContext(ctx, config.ErrGroupCommonLimit)

	var (
		total    atomic.Int64
		progress atomic.Int64
	)
	total.Add(int64(len(userIDs)))
	utils.FuncProgressBarPrint(cctx, gr, &progress, &total)
	for _, userID := range userIDs {
		userID := userID
		gr.Go(func() error {
			token, err := t.GetToken(userID, config.PlatformID)
			userNum := utils.MustGetUserNum(userID)
			err = sdk.TestSDKs[userNum].SDK.LoginWithOutInit(vars.Contexts[userNum], userID, token)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := gr.Wait(); err != nil {
		return err
	}
	return nil
}
