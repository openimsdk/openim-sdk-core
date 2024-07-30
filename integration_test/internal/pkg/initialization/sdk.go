package initialization

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/manager"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
)

func GenUserIDs() {
	manager.NewUserManager(manager.NewMetaManager()).GenAllUserIDs()
}

func InitSDK(ctx context.Context) error {
	userMng := manager.NewUserManager(manager.NewMetaManager())
	err := userMng.InitSDK(ctx, vars.UserIDs...)
	if err != nil {
		return err
	}
	return nil
}
