package manager

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/sdk"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/tools/log"
	"golang.org/x/sync/errgroup"
	"time"
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
	tm := time.Now()
	log.ZDebug(ctx, "createGroups begin")
	defer func() {
		log.ZDebug(ctx, "createGroups end", "time consuming", time.Since(tm))
	}()

	gr, _ := errgroup.WithContext(ctx)
	gr.SetLimit(vars.ErrGroupCommonLimit)
	m.createLargeGroups(ctx, gr)
	m.createCommonGroups(ctx, gr)
	return gr.Wait()
}

// createLargeGroups see CreateGroups
func (m *TestGroupManager) createLargeGroups(ctx context.Context, gr *errgroup.Group) {
	userNum := 0
	for i := 0; i < vars.LargeGroupNum; i++ {
		testSDK := sdk.TestSDKs[userNum]
		gr.Go(func() error {
			_, err := testSDK.CreateLargeGroup(ctx)
			if err != nil {
				return err
			}
			return nil
		})
		userNum = utils.NextNum(userNum)
	}
	return
}

// createLargeGroups see CreateGroups
func (m *TestGroupManager) createCommonGroups(ctx context.Context, gr *errgroup.Group) {
	for userNum := 0; userNum < vars.UserNum; userNum++ {
		testSDK := sdk.TestSDKs[userNum]
		gr.Go(func() error {
			for i := 0; i < vars.CommonGroupNum; i++ {
				_, err := testSDK.CreateCommonGroup(ctx, vars.CommonGroupMemberNum)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	return
}
