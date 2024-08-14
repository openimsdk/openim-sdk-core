package statistics

import (
	"context"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/decorator"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/tools/log"
	"strings"
)

func MsgConsuming(ctx context.Context) {
	defer decorator.FuncLog(ctx)

	close(vars.MsgConsuming)
	var (
		low       int
		mid       int
		high      int
		totalCost float64
		count     int
	)

	for msg := range vars.MsgConsuming {
		sec := msg.Seconds()
		switch {
		case sec < config.ReceiveMsgTimeThresholdLow:
			low++
		case sec < config.ReceiveMsgTimeThresholdMedium:
			mid++
		case sec < config.ReceiveMsgTimeThresholdHigh:
			high++
		}
		totalCost += sec
		count++
	}

	statStr := `
statistic msg count: %d
receive msg in %d s count: %d
receive msg in %d s count: %d
receive msg in %d s count: %d
average time consuming: %.2f
`
	statStr = fmt.Sprintf(statStr,
		count,
		config.ReceiveMsgTimeThresholdLow, low,
		config.ReceiveMsgTimeThresholdMedium, mid,
		config.ReceiveMsgTimeThresholdHigh, high,
		totalCost/float64(count))

	fmt.Println(statStr)
	log.ZInfo(ctx, "stat msg consuming", "res", strings.Replace(statStr, "\n", "; ", -1))
}
