package checker

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/sdk"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
)

// CheckLoginUsersFriends check login users friends
func CheckLoginUsersFriends(ctx context.Context) error {
	corrects := func() [2]int {
		// corrects[0] :super user friend num
		// corrects[1] :common user friend num

		return [2]int{vars.UserNum - 1, vars.SuperUserNum}
	}()

	c := &CounterChecker[*sdk.TestSDK, string]{
		CheckName:      "checkLoginUsersFriends",
		CheckerKeyName: "userID",
		GoroutineLimit: config.ErrGroupCommonLimit,
		GetTotalCount: func(ctx context.Context, t *sdk.TestSDK) (int, error) {
			friendList, err := t.GetAllFriends(ctx)
			if err != nil {
				return 0, err
			}
			return len(friendList), nil
		},
		CalCorrectCount: func(userID string) int {
			if utils.IsSuperUser(userID) {
				return corrects[0]
			} else {
				return corrects[1]
			}
		},
		LoopSlice: sdk.TestSDKs[:vars.LoginUserNum],
		GetKey: func(t *sdk.TestSDK) string {
			return t.UserID
		},
	}

	return c.LoopCheck(ctx)
}
