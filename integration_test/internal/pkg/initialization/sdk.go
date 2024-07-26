package initialization

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/manager"
)

func InitSDK(ctx context.Context) error {
	userMng := manager.NewUserManager(manager.NewMetaManager())
	userIDs := userMng.GenAllUserIDs()
	err := userMng.InitSDK(ctx, userIDs...)
	if err != nil {
		return err
	}
	return nil
}
