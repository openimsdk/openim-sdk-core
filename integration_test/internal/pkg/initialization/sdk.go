package initialization

import (
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/manager"
)

func GenUserIDs() {
	manager.NewUserManager(manager.NewMetaManager()).GenAllUserIDs()
}
