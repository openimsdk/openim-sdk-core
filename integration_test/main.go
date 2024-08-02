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
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/process"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/sdk"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/tools/log"
	"math"
	"time"
)

const (
	sleepSec = 30
)

func Init(ctx context.Context) error {
	initialization.InitFlag()
	flag.Parse()
	initialization.SetFlagLimit()
	initialization.GenUserIDs()
	sdk.TestSDKs = make([]*sdk.TestSDK, vars.UserNum)
	vars.Contexts = make([]context.Context, vars.UserNum)
	vars.LoginEndUserNum = int(math.Floor(vars.LoginRate * float64(vars.UserNum)))
	if err := initialization.InitLog(config.GetConf()); err != nil {
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
		pro = process.NewProcess()

		m           = manager.NewMetaManager()
		userMng     = manager.NewUserManager(m)
		groupMng    = manager.NewGroupManager(m)
		relationMng = manager.NewRelationManager(m)
		msgMng      = manager.NewMsgManager(m)
		fileMng     = manager.NewFileManager(m)
	)
	if err = m.WithAdminToken(); err != nil {
		return err
	}
	ctx = m.BuildCtx(ctx)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pro.SetContext(ctx)

	checkTasks := []*process.Task{
		process.NewTask(vars.ShouldCheckGroupNum, checker.CheckGroupNum),
		process.NewTask(vars.ShouldCheckConversationNum, checker.CheckConvNumAfterImpFriAndCrGro),
		process.NewTask(vars.ShouldCheckMessageNum, checker.CheckMessageNum),
	}

	pro.AddTasks(
		process.NewTask(vars.ShouldRegister, userMng.RegisterAllUsers),
		process.NewTask(vars.ShouldRegister, userMng.InitAllSDK),
		process.NewTask(true, userMng.LoginByRate),

		process.NewTask(vars.ShouldImportFriends, relationMng.ImportFriends),
		process.NewTask(vars.ShouldImportFriends, Sleep),
		process.NewTask(vars.ShouldCreateGroup, groupMng.CreateGroups),
		process.NewTask(vars.ShouldSendMsg, msgMng.SendMessages),

		process.NewTask(true, userMng.LoginLastUsers),
		process.NewTask(true, Sleep),
	)

	pro.AddTasks(checkTasks...)

	// Uninstall and reinstall
	offline := func() {
		ctx = utils.CancelAndReBuildCtx(m.BuildCtx, cancel) // Offline
		log.ZInfo(ctx, "cancel ctx. Offline", err, "stack", utils.FormatErrorStack(err))
		Sleep()
		pro.SetContext(ctx)
	}
	pro.AddTasks(
		process.NewTask(true, process.AddConditions, pro, vars.ShouldCheckUninsAndReins),
		process.NewTask(true, offline),
		process.NewTask(true, fileMng.DeleteLocalDB),
		process.NewTask(true, userMng.InitAllSDK),
		process.NewTask(true, userMng.LoginAllUsers),
		process.NewTask(true, Sleep),
		process.NewTask(true, Sleep),
	)
	pro.AddTasks(checkTasks...)

	return pro.Exec()
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
