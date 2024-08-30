package statistics

import (
	"context"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/decorator"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/tools/log"
	"strings"
	"time"
)

func MsgConsuming(ctx context.Context) {
	defer decorator.FuncLog(ctx)()
	time.Sleep(time.Second * 5)
	close(vars.RecvMsgConsuming)
	var (
		low       int
		mid       int
		high      int
		outHigh   int
		minT      = int64(-1)
		maxT      = int64(-1)
		totalCost int64
		count     int
	)
	for msg := range vars.RecvMsgConsuming {

		sec := msg.CostTime.Milliseconds()
		switch {
		case sec < config.ReceiveMsgTimeThresholdLow*int64(time.Millisecond):
			low++
		case sec < config.ReceiveMsgTimeThresholdMedium*int64(time.Millisecond):
			mid++
		case sec < config.ReceiveMsgTimeThresholdHigh*int64(time.Millisecond):
			high++
		default:
			outHigh++
		}

		if minT == -1 || minT > sec {
			minT = sec
		}
		if maxT == -1 || maxT < sec {
			maxT = sec
		}

		totalCost += sec
		count++
	}

	statStr := `
statistic msg count: %d
statistic send msg count: %d
receive msg in %d ms count: %d
receive msg in %d ms count: %d
receive msg in %d ms count: %d
receive messages within %d s or more: %d ms
maximum time to receive messages: %d ms
minimum time to receive messages: %d ms
average time consuming: %.2f
`
	statStr = fmt.Sprintf(statStr,
		count,
		vars.SendMsgCount.Load(),
		config.ReceiveMsgTimeThresholdLow, low,
		config.ReceiveMsgTimeThresholdMedium, mid,
		config.ReceiveMsgTimeThresholdHigh, high,
		config.ReceiveMsgTimeThresholdHigh, outHigh,
		maxT,
		minT,
		float64(totalCost)/float64(count))

	fmt.Println(statStr)
	log.ZInfo(ctx, "stat msg consuming", "res", strings.Replace(statStr, "\n", "; ", -1))
}
