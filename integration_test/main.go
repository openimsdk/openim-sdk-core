package main

import (
	"context"
	"flag"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/manager"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/initialization"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/process/checker"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/tools/log"
)

func Init(ctx context.Context) error {
	initialization.InitFlag()
	flag.Parse()
	initialization.SetFlagLimit()
	if err := initialization.InitSDK(ctx); err != nil {
		return err
	}
	return nil
}

func DoFlagFunc(ctx context.Context) (err error) {
	defer func() {
		// capture check err
		ch := checker.CloseAndGetCheckErrChan()
		for e := range ch {
			if err == nil {
				err = e
			}
			log.ZError(ctx, "check err", err)
		}
	}()
	var (
		m        = manager.NewMetaManager()
		userMng  = manager.NewUserManager(m)
		groupMng = manager.NewGroupManager(m)
	)

	if vars.ShouldRegister {
		if err = userMng.RegisterUsers(ctx, vars.UserIDs...); err != nil {
			return err
		}
	}

	if err = userMng.LoginByRate(ctx, vars.LoginRate); err != nil {
		return err
	}

	if vars.ShouldCreateGroup {
		if err = groupMng.CreateGroups(ctx); err != nil {
			return err
		}
	}

	if vars.ShouldCheckGroupNum {
		if err = checker.CheckGroupNum(ctx); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	ctx := context.Background()
	if err := Init(ctx); err != nil {
		panic(err)
	}

	if err := DoFlagFunc(ctx); err != nil {
		panic(err)
	}
	//if err := process.ProV1(ctx); err != nil {
	//	panic(err)
	//}
	log.ZInfo(ctx, "start success!")
}
