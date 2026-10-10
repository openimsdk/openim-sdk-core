package checker

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
)

func CheckLoginByRateNum(ctx context.Context) error {
	correct := func() int {
		return vars.LoginUserNum
	}()

	c := &CounterChecker[int, string]{
		CheckName:      "checkLoginByRateNum",
		CheckerKeyName: "loginNum",
		GoroutineLimit: config.ErrGroupCommonLimit,
		GetTotalCount: func(ctx context.Context, t int) (int, error) {
			return int(vars.NowLoginNum.Load()), nil
		},
		CalCorrectCount: func(_ string) int {
			return correct
		},
		LoopSlice: []int{0},
		GetKey: func(t int) string {
			return "login"
		},
	}

	return c.LoopCheck(ctx)
}

// CheckAllLoginNum check if all user is login
func CheckAllLoginNum(ctx context.Context) error {
	correct := func() int {
		return vars.UserNum
	}()

	c := &CounterChecker[int, string]{
		CheckName:      "checkLoginByRateNum",
		CheckerKeyName: "loginNum",
		GoroutineLimit: config.ErrGroupCommonLimit,
		GetTotalCount: func(ctx context.Context, t int) (int, error) {
			return int(vars.NowLoginNum.Load()), nil
		},
		CalCorrectCount: func(_ string) int {
			return correct
		},
		LoopSlice: []int{0},
		GetKey: func(t int) string {
			return "login"
		},
	}

	return c.LoopCheck(ctx)
}
