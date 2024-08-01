package initialization

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/manager"
)

func GenUserIDs() {
	manager.NewUserManager(manager.NewMetaManager()).GenAllUserIDs()
}

func InitAllSDK(ctx context.Context) error {
	userMng := manager.NewUserManager(manager.NewMetaManager())
	err := userMng.InitAllSDK(ctx)
	if err != nil {
		return err
	}
	return nil
}
