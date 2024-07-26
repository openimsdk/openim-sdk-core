package process

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/manager"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/process/checker"
	"github.com/openimsdk/tools/errs"
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

	checkerMap, err := checker.CheckGroupNum(ctx)
	if err != nil {
		return err
	}
	if len(checkerMap) != 0 {
		for k, ck := range checkerMap {
			log.ZInfo(ctx, "group num un correct", "userID", k, "group num", ck.TotalCount, "correct num", ck.CorrectCount)
		}
		err = errs.New("check group number un correct!").Wrap()
		return err
	}
	return nil
}

func RegisterAndLogin(ctx context.Context, userMng *manager.TestUserManager) error {
	t := time.Now()
	log.ZDebug(ctx, "registerAndLogin begin")

	userIDs = userMng.GenUserIDs()
	err := userMng.RegisterUsers(ctx, userIDs...)
	if err != nil {
		return err
	}
	err = userMng.InitSDKAndLogin(ctx, userIDs...)
	if err != nil {
		return err
	}

	log.ZDebug(ctx, "registerAndLogin over", "time consuming", time.Since(t))
	return nil
}

func Login(ctx context.Context, userMng *manager.TestUserManager) error {
	t := time.Now()
	log.ZDebug(ctx, "login begin")

	userIDs = userMng.GenUserIDs()
	err := userMng.InitSDKAndLogin(ctx, userIDs...)
	if err != nil {
		return err
	}

	log.ZDebug(ctx, "login over", "time consuming", time.Since(t))
	return nil
}
