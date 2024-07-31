package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/checker"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/manager"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/initialization"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/sdk"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/tools/log"
	"time"
)

const (
	sleepSec = 10
)

func Init(ctx context.Context) error {
	initialization.InitFlag()
	flag.Parse()
	initialization.SetFlagLimit()
	initialization.GenUserIDs()
	sdk.TestSDKs = make([]*sdk.TestSDK, vars.UserNum)
	vars.Contexts = make([]context.Context, vars.UserNum)
	vars.ReceiveMsgNum = make([]int64, vars.UserNum)
	if err := initialization.InitLog(config.GetConf()); err != nil {
		return err
	}
	if !vars.ShouldRegister {
		if err := initialization.InitSDK(ctx); err != nil {
			return err
		}
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
		m           = manager.NewMetaManager()
		userMng     = manager.NewUserManager(m)
		groupMng    = manager.NewGroupManager(m)
		relationMng = manager.NewRelationManager(m)
		msgMng      = manager.NewMsgManager(m)
	)
	if err = m.WithAdminToken(); err != nil {
		return err
	}
	ctx = m.BuildCtx(ctx)

	if vars.ShouldRegister {
		if err = userMng.RegisterUsers(ctx, vars.UserIDs...); err != nil {
			return err
		}
		if err = initialization.InitSDK(ctx); err != nil {
			return err
		}
	}

	if err = userMng.LoginByRate(ctx, vars.LoginRate); err != nil {
		return err
	}

	if vars.ShouldImportFriends {
		if err = relationMng.ImportFriends(ctx); err != nil {
			return err
		}
	}

	// sync data
	if vars.ShouldRegister {
		Sleep()
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

	if vars.ShouldCheckConversationNum {
		if err = checker.CheckConvNumAfterImpFriAndCrGro(ctx); err != nil {
			return err
		}
	}

	if vars.ShouldSendMsg {
		if err = msgMng.SendMessages(ctx); err != nil {
			return err
		}
		Sleep()
	}

	if vars.ShouldCheckMessageNum {
		if err = checker.CheckMessageNum(ctx); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	ctx := context.Background()
	if err := Init(ctx); err != nil {
		log.ZError(ctx, "init err", err, "stack", utils.FormatErrorStack(err))
		fmt.Println("init err")
		fmt.Println(utils.FormatErrorStack(err))
		panic(err)
	}
	if err := DoFlagFunc(ctx); err != nil {
		log.ZError(ctx, "do flag err", err, "stack", utils.FormatErrorStack(err))
		fmt.Println("do flag err")
		fmt.Println(utils.FormatErrorStack(err))
		panic(err)
	}

	log.ZInfo(ctx, "start success!")
	fmt.Println("start success!")
	select {}
}

func Sleep() {
	fmt.Printf("sleep %d s for sync data~\n", sleepSec)
	for i := 0; i < sleepSec; i++ {
		time.Sleep(time.Second)
		fmt.Println(i + 1)
	}
}
