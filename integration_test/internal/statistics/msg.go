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
		minT      *vars.StatMsg
		maxT      *vars.StatMsg
		totalCost int64
		count     int
	)
	for msg := range vars.RecvMsgConsuming {

		sec := msg.CostTime
		switch {
		case sec < config.ReceiveMsgTimeThresholdLow*1000:
			low++
		case sec < config.ReceiveMsgTimeThresholdMedium*1000:
			mid++
		case sec < config.ReceiveMsgTimeThresholdHigh*1000:
			high++
		default:
			outHigh++
		}

		if minT == nil || minT.CostTime > sec {
			minT = msg
		}
		if maxT == nil || maxT.CostTime < sec {
			maxT = msg
		}

		totalCost += sec
		count++
	}

	if minT == nil || maxT == nil {
		return
	}
	statStr := `
statistic msg count: %d
statistic send msg count: %d
receive msg in %d s count: %d
receive msg in %d s count: %d
receive msg in %d s count: %d
receive messages within %d s or more: %d
maximum time to receive messages: %d ms, create: %d, receive: %d, msg: %v
minimum time to receive messages: %d ms, create: %d, receive: %d, msg: %v
average time consuming: %.2f ms
`
	statStr = fmt.Sprintf(statStr,
		count,
		vars.SendMsgCount.Load(),
		config.ReceiveMsgTimeThresholdLow, low,
		config.ReceiveMsgTimeThresholdMedium, mid,
		config.ReceiveMsgTimeThresholdHigh, high,
		config.ReceiveMsgTimeThresholdHigh, outHigh,
		maxT.CostTime, maxT.Msg.CreateTime, maxT.ReceiveTime, *maxT.Msg,
		minT.CostTime, minT.Msg.CreateTime, minT.ReceiveTime, *minT.Msg,
		float64(totalCost)/float64(count))

	fmt.Println(statStr)
	log.ZInfo(ctx, "stat msg consuming", "res", strings.Replace(statStr, "\n", "; ", -1))
}
