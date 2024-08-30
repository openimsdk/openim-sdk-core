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
		minT      = float64(-1)
		maxT      = float64(-1)
		totalCost float64
		count     int
	)
	for msg := range vars.RecvMsgConsuming {

		sec := msg.CostTime.Seconds()
		switch {
		case sec < config.ReceiveMsgTimeThresholdLow:
			low++
		case sec < config.ReceiveMsgTimeThresholdMedium:
			mid++
		case sec < config.ReceiveMsgTimeThresholdHigh:
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
receive msg in %d s count: %d
receive msg in %d s count: %d
receive msg in %d s count: %d
receive messages within %d s or more: %d
maximum time to receive messages: %.2f
minimum time to receive messages: %.2f
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
		totalCost/float64(count))

	fmt.Println(statStr)
	log.ZInfo(ctx, "stat msg consuming", "res", strings.Replace(statStr, "\n", "; ", -1))
}
