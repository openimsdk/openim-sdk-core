package checker

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/sdk"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
)

// CheckMessageNum check message num.
func CheckMessageNum(ctx context.Context) error {
	corrects := func() [2]int {
		// corrects[0] :super user conversion num
		// corrects[1] :common user conversion num
		largeGroupNum := ((vars.UserNum-1)*vars.GroupMessageNum + 1) * vars.LargeGroupNum
		commonGroupNum := (vars.GroupMessageNum)*(vars.CommonGroupNum*(vars.CommonGroupMemberNum-1)) +
			vars.CommonGroupMemberNum*vars.CommonGroupMemberNum
		groupMsgNum := largeGroupNum + commonGroupNum

		superUserMsgNum := (vars.UserNum - 1) * (vars.SingleMessageNum + 1) // send message + become friend message
		commonUserMsgNum := vars.SuperUserNum * (vars.SingleMessageNum + 1)

		return [2]int{superUserMsgNum + groupMsgNum, commonUserMsgNum + groupMsgNum}
	}()

	c := &CounterChecker[*sdk.TestSDK, string]{
		CheckName:      "checkMessageNum",
		CheckerKeyName: "userID",
		GoroutineLimit: config.ErrGroupCommonLimit,
		GetTotalCount: func(ctx context.Context, t *sdk.TestSDK) (int, error) {
			totalNum, err := t.SDK.Conversation().GetTotalUnreadMsgCount(ctx)
			if err != nil {
				return 0, err
			}
			return int(totalNum), nil
		},
		CalCorrectCount: func(userID string) int {
			if utils.IsSuperUser(userID) {
				return corrects[0]
			} else {
				return corrects[1]
			}
		},
		LoopSlice: sdk.TestSDKs,
		GetKey: func(t *sdk.TestSDK) string {
			return t.UserID
		},
	}

	return c.Check(ctx)
}
