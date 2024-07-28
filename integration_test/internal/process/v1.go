package process

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/manager"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/tools/log"
	"time"
)

var (
	mng     = manager.NewMetaManager()
	userIDs []string
)

func ProV1(ctx context.Context) error {
	var (
		userMng     = manager.NewUserManager(mng)
		relationMng = manager.NewRelationManager(mng)
		groupMng    = manager.NewGroupManager(mng)
	)

	if err := RegisterAndLogin(ctx, userMng); err != nil {
		return err
	}
	if err := relationMng.ImportFriends(ctx); err != nil {
		return err
	}
	if err := groupMng.CreateGroups(ctx); err != nil {
		return err
	}

	return nil
}

func RegisterAndLogin(ctx context.Context, userMng *manager.TestUserManager) error {
	t := time.Now()
	log.ZDebug(ctx, "registerAndLogin begin")

	userIDs = vars.UserIDs
	err := userMng.RegisterUsers(ctx, userIDs...)
	if err != nil {
		return err
	}
	err = userMng.InitSDK(ctx, userIDs...)
	if err != nil {
		return err
	}

	log.ZDebug(ctx, "registerAndLogin over", "time consuming", time.Since(t))
	return nil
}

func Login(ctx context.Context, userMng *manager.TestUserManager) error {
	t := time.Now()
	log.ZDebug(ctx, "login begin")

	userIDs = vars.UserIDs
	err := userMng.Login(ctx, userIDs...)
	if err != nil {
		return err
	}

	log.ZDebug(ctx, "login over", "time consuming", time.Since(t))
	return nil
}
