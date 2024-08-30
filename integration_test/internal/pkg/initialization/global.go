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
	vars.Cancels = make([]context.CancelFunc, vars.UserNum)
	vars.LoginUserNum = int(math.Floor(vars.LoginRate * float64(vars.UserNum)))
	vars.RecvMsgConsuming = make(chan *vars.StatMsg, config.MaxCheckMsg)
	rand.New(rand.NewSource(time.Now().UnixNano()))
}
