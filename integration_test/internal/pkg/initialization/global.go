package initialization

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/sdk"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"math"
	"math/rand"
	"time"
)

func InitGlobalData() {
	sdk.TestSDKs = make([]*sdk.TestSDK, vars.UserNum)
	vars.Contexts = make([]context.Context, vars.UserNum)
	vars.LoginEndUserNum = int(math.Floor(vars.LoginRate * float64(vars.UserNum)))
	vars.MsgConsuming = make(chan time.Duration, config.MaxCheckMsg)
	rand.New(rand.NewSource(time.Now().UnixNano()))
}
