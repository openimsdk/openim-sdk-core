package manager

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/decorator"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/sdk_user_simulator"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/sdk"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/protocol/sdkws"
	userPB "github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/log"
	"golang.org/x/sync/errgroup"
	"math"
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

func (t *TestUserManager) RegisterUsers(ctx context.Context, userIDs ...string) error {
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

func (t *TestUserManager) InitSDK(ctx context.Context, userIDs ...string) error {
	defer decorator.FuncLog(ctx)()

	gr, _ := errgroup.WithContext(ctx)
	gr.SetLimit(config.ErrGroupCommonLimit)
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
			return nil
		})
	}
	if err := gr.Wait(); err != nil {
		return err
	}
	return nil
}

func (t *TestUserManager) LoginByRate(ctx context.Context, rate float64) error {
	right := int(math.Ceil(rate * float64(vars.UserNum)))
	userIDs := vars.UserIDs[:right]
	log.ZDebug(ctx, "login users", "len", len(userIDs))
	return t.Login(ctx, userIDs...)
}

func (t *TestUserManager) Login(ctx context.Context, userIDs ...string) error {
	defer decorator.FuncLog(ctx)()

	gr, _ := errgroup.WithContext(ctx)
	gr.SetLimit(config.ErrGroupCommonLimit)
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
