package manager

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/tools/log"
	"golang.org/x/sync/errgroup"
	"time"
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
	tm := time.Now()
	log.ZDebug(ctx, "importFriends begin")
	defer func() {
		log.ZDebug(ctx, "importFriends end", "time consuming", time.Since(tm))
	}()

	gr, _ := errgroup.WithContext(ctx)
	gr.SetLimit(vars.ErrGroupCommonLimit)
	for i, userID := range vars.SuperUserIDs {
		i := i
		userID := userID
		gr.Go(func() error {
			friendIDs := append(vars.UserIDs[:i], vars.UserIDs[i+1:]...) // excluding oneself
			req := &relation.ImportFriendReq{
				OwnerUserID:   userID,
				FriendUserIDs: friendIDs,
			}
			resp := &relation.ImportFriendResp{}
			err := m.PostWithCtx(constant.ImportFriendListRouter, req, resp)
			if err != nil {
				return err
			}
			return nil
		})
	}
	return gr.Wait()
}
