package checker

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/sdk"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
)

// CheckConvNumAfterImpFriAndCrGro check conversation num after import friends and create groups.
func CheckConvNumAfterImpFriAndCrGro(ctx context.Context) error {
	corrects := func() [2]int {
		// corrects[0] :super user conversion num
		// corrects[1] :common user conversion num
		largeNum := vars.LargeGroupNum
		commonNum := 0 // cal by userNum
		groupNum := largeNum + commonNum

		superNum := vars.UserNum - 1 + groupNum
		commonUserNum := vars.SuperUserNum + groupNum

		return [2]int{superNum, commonUserNum}
	}()

	c := &CounterChecker[*sdk.TestSDK, string]{
		CheckName:      "checkConversationNum",
		CheckerKeyName: "userID",
		GoroutineLimit: config.ErrGroupCommonLimit,
		GetTotalCount: func(ctx context.Context, t *sdk.TestSDK) (int, error) {
			totalNum, err := t.GetTotalConversationCount(ctx)
			if err != nil {
				return 0, err
			}
			return totalNum, nil
		},
		CalCorrectCount: func(userID string) int {
			commonGroupNum := calCommonGroup(utils.MustGetUserNum(userID))
			if utils.IsSuperUser(userID) {
				return corrects[0] + commonGroupNum
			} else {
				return corrects[1] + commonGroupNum
			}
		},
		LoopSlice: sdk.TestSDKs,
		GetKey: func(t *sdk.TestSDK) string {
			return t.UserID
		},
	}

	return c.LoopCheck(ctx)
}
