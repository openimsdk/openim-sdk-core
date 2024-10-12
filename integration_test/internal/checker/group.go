package checker

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/sdk"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
)

func CheckGroupNum(ctx context.Context) error {
	correct := func() int {
		largeNum := vars.LargeGroupNum
		commonNum := 0 // cal by userNum
		return largeNum + commonNum
	}()

	c := &CounterChecker[*sdk.TestSDK, string]{
		CheckName:      "checkGroupNum",
		CheckerKeyName: "userID",
		GoroutineLimit: config.ErrGroupCommonLimit,
		GetTotalCount: func(ctx context.Context, t *sdk.TestSDK) (int, error) {
			_, groupNum, err := t.GetAllJoinedGroups(ctx)
			if err != nil {
				return 0, err
			}
			return groupNum, nil
		},
		CalCorrectCount: func(userID string) int {
			return correct + calCommonGroup(utils.MustGetUserNum(userID))
		},
		LoopSlice: sdk.TestSDKs,
		GetKey: func(t *sdk.TestSDK) string {
			return t.UserID
		},
	}

	return c.LoopCheck(ctx)
}

func calCommonGroup(userNum int) int {

	preNum := utils.NextOffsetNum(userNum, -(vars.CommonGroupMemberNum - 1))

	createNum := 0
	for i := 0; i < vars.CommonGroupMemberNum; i++ {
		if utils.IsNumLogin(preNum) {
			createNum++
		}
		preNum = utils.NextNum(preNum)
	}
	return createNum * vars.CommonGroupNum
}
