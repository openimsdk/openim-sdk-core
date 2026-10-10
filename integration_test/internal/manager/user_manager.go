package manager

import (
	"context"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/decorator"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/progress"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/reerrgroup"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/sdk_user_simulator"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/sdk"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	sdkUtils "github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/protocol/sdkws"
	userPB "github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
)

type TestUserManager struct {
	*MetaManager
}

func NewUserManager(m *MetaManager) *TestUserManager {
	return &TestUserManager{m}
}

func (t *TestUserManager) GenAllUserIDs() []string {
	ids := utils.GenUserIDs(vars.UserNum)
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

	for i := 0; i < len(users); i += config.ApiParamLength {
		end := i + config.ApiParamLength
		if end > len(users) {
			end = len(users)
		}
		if err := t.PostWithCtx(api.UserRegister.Route(), &userPB.UserRegisterReq{
			Users: users[i:end],
		}, nil); err != nil {
			return err
		}
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
		total int
		now   int
	)
	total = len(userIDs)
	progress.FuncNameBarPrint(cctx, gr, now, total)

	for _, userID := range userIDs {
		userID := userID
		gr.Go(func() error {
			userNum := utils.MustGetUserNum(userID)
			mgr, err := sdk_user_simulator.InitSDK(userID, t.IMConfig)
			if err != nil {
				return err
			}
			sdk.TestSDKs[userNum] = sdk.NewTestSDK(userID, userNum, mgr) // init sdk
			ctx := mgr.Context()
			ctx = ccontext.WithOperationID(ctx, sdkUtils.OperationIDGenerator())
			ctx = mcontext.SetOpUserID(ctx, userID)
			ctx, cancel := context.WithCancel(ctx)
			vars.Contexts[userNum] = ctx // init ctx
			vars.Cancels[userNum] = cancel
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
	userIDs := vars.UserIDs[:vars.LoginUserNum]
	return t.login(ctx, userIDs...)
}

func (t *TestUserManager) LoginLastUsers(ctx context.Context) error {
	userIDs := vars.UserIDs[vars.LoginUserNum:]
	return t.login(ctx, userIDs...)
}

func (t *TestUserManager) login(ctx context.Context, userIDs ...string) error {
	defer decorator.FuncLog(ctx)()

	log.ZDebug(ctx, "login users", "len", len(userIDs), "userIDs", userIDs)

	gr, cctx := reerrgroup.WithContext(ctx, config.ErrGroupCommonLimit)

	var (
		total int
		now   int
	)
	total = len(userIDs)
	progress.FuncNameBarPrint(cctx, gr, now, total)
	for _, userID := range userIDs {
		userID := userID
		gr.Go(func() error {
			token, err := t.GetUserToken(userID, config.PlatformID)
			if err != nil {
				return err
			}
			userNum := utils.MustGetUserNum(userID)
			err = sdk.TestSDKs[userNum].SDK.Login(vars.Contexts[userNum], userID, token)
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

func (t *TestUserManager) ForceLogoutAllUsers(ctx context.Context) error {
	return t.forceLogout(ctx, vars.UserIDs...)
}

func (t *TestUserManager) forceLogout(ctx context.Context, userIDs ...string) error {
	defer decorator.FuncLog(ctx)()

	log.ZDebug(ctx, "logout users", "len", len(userIDs), "userIDs", userIDs)

	cancelBar := progress.NewBar("cancel ctx", 0, len(userIDs), false)
	pro := progress.Start(cancelBar)
	for _, userID := range userIDs {
		pro.IncBar(cancelBar)
		vars.Cancels[utils.MustGetUserNum(userID)]()
	}

	sleepTime := 30 * 2 // unit: second. pone wait * 2
	sleepBar := progress.NewBar("sleep", 0, sleepTime, false)
	pro.AddBar(sleepBar)
	for i := 0; i < sleepTime; i++ {
		time.Sleep(time.Second)
		pro.IncBar(sleepBar)
	}
	return nil
}
