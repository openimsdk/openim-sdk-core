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
	createdLargeGroupNum := vars.LargeGroupNum / vars.LoginUserNum
	corrects := func() [3]int {
		// corrects[0]: super user msg num
		// corrects[1]: common user msg num
		// corrects[2]: create more one large group largest user no + 1

		// if a user num smaller than remainder, it means this user created more one large group
		remainder := vars.LargeGroupNum % vars.LoginUserNum

		largeGroupNum :=
			// total send message num +
			vars.GroupMessageNum*vars.LoginUserNum*vars.LargeGroupNum +
				// total create group notification message -
				vars.LargeGroupNum
		// self send group message(cal by userID) -
		// self create group notification message. Complete the calculation based on user ID in CalCorrectCount.

		commonGroupNum := 0
		// self create group notification message
		// Formula:
		// commonGroupNum =
		// total send group message(cal by userID) +
		// total create group notification message(cal by userID) -
		// self send group message(cal by userID) -
		// self create group notification message(cal by userID)

		groupMsgNum := largeGroupNum + commonGroupNum

		superUserMsgNum := 0
		// Formula:
		// superUserMsgNum =
		//	friend send message num(cal by userID) +
		//	become friend notification message num(cal by userID)

		commonUserMsgNum := min(vars.LoginUserNum, vars.SuperUserNum) * vars.SingleMessageNum

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
			userNum := utils.MustGetUserNum(userID)
			if utils.IsSuperUser(userID) {
				res += corrects[0]
				res += vars.UserNum - 1 - userNum // become friend notification message num
				if userNum < vars.LoginUserNum {
					// friend send message num
					res += vars.SingleMessageNum * (vars.LoginUserNum - 1)
					// self send large group message
					res -= vars.GroupMessageNum * vars.LargeGroupNum
					res -= createdLargeGroupNum
				} else {
					// friend send message num
					res += vars.SingleMessageNum * vars.LoginUserNum
					// self send large group message
					res -= 0
				}
			} else {
				res += corrects[1]
				if userNum < vars.LoginUserNum {
					// self send large group message
					res -= vars.GroupMessageNum * vars.LargeGroupNum
					// self created large group num
					res -= createdLargeGroupNum
				} else {
					// self send large group message
					res -= 0
				}
			}

			// commonGroupNum =
			// total send group message(cal by userID) +
			// total create group notification message(cal by userID) -
			// self send group message(cal by userID) -
			// self create group notification message(cal by userID)
			comGroupNum := calCommonGroup(userNum) * (vars.GroupMessageNum + 1)
			if utils.IsNumLogin(userNum) {
				comGroupNum -= vars.CommonGroupNum * (vars.GroupMessageNum + 1)
			}
			res += comGroupNum

			// create more one large group
			if userNum < corrects[2] {
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
