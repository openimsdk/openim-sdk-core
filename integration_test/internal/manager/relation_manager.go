package manager

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/decorator"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/progress"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/reerrgroup"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/tools/log"
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

	gr, cctx := reerrgroup.WithContext(ctx, config.ErrGroupMiddleSmallLimit)

	var (
		total int
		now   int
	)
	total = vars.SuperUserNum
	progress.FuncNameBarPrint(cctx, gr, now, total)
	for i, userID := range vars.SuperUserIDs {
		i := i
		userID := userID
		gr.Go(func() error {
			friendIDs := vars.UserIDs[i+1:] // excluding oneself
			if len(friendIDs) == 0 {
				return nil
			}

			for i := 0; i < len(friendIDs); i += config.ApiParamLength {
				end := i + config.ApiParamLength
				if end > len(friendIDs) {
					end = len(friendIDs)
				}
				req := &relation.ImportFriendReq{
					OwnerUserID:   userID,
					FriendUserIDs: friendIDs[i:end],
				}
				ctx := m.BuildCtx(ctx)
				log.ZWarn(ctx, "ImportFriends begin", nil, "len", len(friendIDs))
				if err := api.ImportFriendList.Execute(ctx, req); err != nil {
					return err
				}
				log.ZWarn(ctx, "ImportFriends end", nil, "len", len(friendIDs))
			}
			return nil
		})
	}
	return gr.Wait()
}
