package manager

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/decorator"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/progress"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/reerrgroup"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/sdk"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	sdkUtils "github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/stringutil"
)

type TestGroupManager struct {
	*MetaManager
}

func NewGroupManager(m *MetaManager) *TestGroupManager {
	return &TestGroupManager{m}
}

// CreateGroups creates group chats. It needs to create both large group chats and regular group chats.
// The number of large group chats to be created is specified by vars.LargeGroupNum, and the group owner cycles from 0 to vars.UserNum.
// Every user creates regular group chats, and the number of regular group chats to be created is specified by vars.CommonGroupNum.
func (m *TestGroupManager) CreateGroups(ctx context.Context) error {
	defer decorator.FuncLog(ctx)()

	gr, cctx := reerrgroup.WithContext(ctx, config.ErrGroupCommonLimit)
	var (
		total int
		now   int
	)
	total = vars.LargeGroupNum + vars.LoginUserNum
	p := progress.FuncNameBarPrint(cctx, gr, now, total)

	m.createLargeGroups(ctx, gr, p)
	m.createCommonGroups(ctx, gr, p)
	return gr.Wait()
}

// createLargeGroups see CreateGroups
func (m *TestGroupManager) createLargeGroups(ctx context.Context, gr *reerrgroup.Group, p *progress.Progress) {
	userNum := 0

	bar := progress.NewRemoveBar(stringutil.GetSelfFuncName(), 0, vars.LargeGroupNum)
	p.AddBar(bar)

	for i := 0; i < vars.LargeGroupNum; i++ {

		ctx := vars.Contexts[userNum]
		testSDK := sdk.TestSDKs[userNum]
		gr.Go(func() error {
			ctx = ccontext.WithOperationID(ctx, sdkUtils.OperationIDGenerator())
			log.ZWarn(ctx, "createLargeGroups begin", nil)
			_, err := testSDK.CreateLargeGroup(ctx)
			if err != nil {
				return err
			}
			log.ZWarn(ctx, "createLargeGroups end", nil)

			p.IncBar(bar)

			return nil
		})
		userNum = utils.NextLoginNum(userNum)
	}
	return
}

// createLargeGroups see CreateGroups
func (m *TestGroupManager) createCommonGroups(ctx context.Context, gr *reerrgroup.Group, p *progress.Progress) {

	bar := progress.NewRemoveBar(stringutil.GetSelfFuncName(), 0, vars.CommonGroupNum*vars.LoginUserNum)
	p.AddBar(bar)

	for userNum := 0; userNum < vars.LoginUserNum; userNum++ {
		userNum := userNum
		ctx := vars.Contexts[userNum]
		testSDK := sdk.TestSDKs[userNum]

		gr.Go(func() error {
			ubar := progress.NewRemoveBar(utils.GetUserID(userNum), 0, vars.CommonGroupNum)
			p.AddBar(ubar)
			for i := 0; i < vars.CommonGroupNum; i++ {
				ctx = ccontext.WithOperationID(ctx, sdkUtils.OperationIDGenerator())
				log.ZWarn(ctx, "createCommonGroups begin", nil)
				_, err := testSDK.CreateCommonGroup(ctx, vars.CommonGroupMemberNum)
				if err != nil {
					return err
				}
				log.ZWarn(ctx, "createCommonGroups end", nil)

				p.IncBar(bar, ubar)
			}
			return nil
		})
	}
	return
}
