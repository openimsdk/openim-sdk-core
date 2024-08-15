package manager

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/decorator"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/reerrgroup"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/protocol/relation"
	"sync/atomic"
)

type TestRelationManager struct {
	*MetaManager
}

func NewRelationManager(m *MetaManager) *TestRelationManager {
	return &TestRelationManager{m}
}

// ImportFriends Import all users as friends of the superuser (excluding themselves),
// making the superuser have all users as friends,
// while regular users have the superuser as their only friend.
// A superuser is defined as a user who has all users as friends,
// their IDs range from 0 to vars.SuperUserNum.
func (m *TestRelationManager) ImportFriends(ctx context.Context) error {
	defer decorator.FuncLog(ctx)()

	gr, cctx := reerrgroup.WithContext(ctx, config.ErrGroupSmallLimit)

	var (
		total    atomic.Int64
		progress atomic.Int64
	)
	total.Add(int64(vars.SuperUserNum))
	utils.FuncProgressBarPrint(cctx, gr, &progress, &total)
	for i, userID := range vars.SuperUserIDs {
		i := i
		userID := userID
		gr.Go(func() error {
			friendIDs := vars.UserIDs[i+1:] // excluding oneself
			req := &relation.ImportFriendReq{
				OwnerUserID:   userID,
				FriendUserIDs: friendIDs,
			}
			_, err := util.CallApi[relation.ImportFriendResp](m.BuildCtx(ctx), constant.ImportFriendListRouter, req)
			if err != nil {
				return err
			}
			return nil
		})
	}
	return gr.Wait()
}
