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
	corrects := func() [3]int {
		// corrects[0]: super user msg num
		// corrects[1]: common user msg num
		// corrects[2]: create more one large group largest user no + 1

		createdLargeGroupNum := vars.LargeGroupNum / vars.UserNum
		// if a user num smaller than remainder, it means this user created more one large group
		remainder := vars.LargeGroupNum % vars.UserNum

		largeGroupNum := ((vars.UserNum-1)*vars.GroupMessageNum+1)*vars.LargeGroupNum - createdLargeGroupNum
		// Formula:
		// largeGroupNum =
		//	// total send message num
		//	vars.GroupMessageNum*vars.UserNum*vars.LargeGroupNum +
		//	// total create group notification message
		//	vars.LargeGroupNum -
		//	// self send group message
		//	vars.GroupMessageNum*vars.LargeGroupNum -
		//	// self create group notification message. Complete the calculation based on user ID in CalCorrectCount.
		//	createdLargeGroupNum

		commonGroupNum := (vars.GroupMessageNum + 1) * (vars.CommonGroupNum * (vars.CommonGroupMemberNum - 1))
		// Formula:
		// commonGroupNum =
		// // total send group message
		// vars.GroupMessageNum*(vars.CommonGroupMemberNum*vars.CommonGroupNum) +
		// // total create group notification message
		//	(vars.CommonGroupMemberNum * vars.CommonGroupNum) -
		// // self send group message
		//	vars.GroupMessageNum*vars.CommonGroupNum -
		// // self create group notification message
		//	vars.CommonGroupNum

		groupMsgNum := largeGroupNum + commonGroupNum

		superUserMsgNum := (vars.UserNum - 1) * vars.SingleMessageNum // send message + become friend message(in CalCorrectCount)
		// Formula:
		// superUserMsgNum =
		//	// friend send message num
		//	(vars.UserNum-1)*vars.GroupMessageNum +
		//	// become friend notification message num. it`s number of friends applied for
		//	userNum

		commonUserMsgNum := vars.SuperUserNum * vars.SingleMessageNum

		return [3]int{superUserMsgNum + groupMsgNum, commonUserMsgNum + groupMsgNum, remainder}
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
			var res int
			useNum := utils.MustGetUserNum(userID)
			if utils.IsSuperUser(userID) {
				res = corrects[0] + vars.UserNum - 1 - useNum // become friend message
			} else {
				res = corrects[1]
			}
			if useNum < corrects[2] {
				res--
			}
			return res
		},
		LoopSlice: sdk.TestSDKs,
		GetKey: func(t *sdk.TestSDK) string {
			return t.UserID
		},
	}

	return c.LoopCheck(ctx)
}
