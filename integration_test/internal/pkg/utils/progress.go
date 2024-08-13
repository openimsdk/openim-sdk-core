package utils

import (
	"context"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/reerrgroup"
	"github.com/openimsdk/tools/utils/formatutil"
	"github.com/openimsdk/tools/utils/stringutil"
	"sync/atomic"
	"time"
)

func FuncProgressBarPrint(ctx context.Context, gr *reerrgroup.Group, progress, total *atomic.Int64) {
	ProgressBarPrint(ctx, stringutil.GetFuncName(1), gr, progress, total)
}
func ProgressBarPrint(ctx context.Context, title string, gr *reerrgroup.Group, progress, total *atomic.Int64) {
	gr.SetAfterTasks(func() error {
		progress.Add(1)
		return nil
	})
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Print(formatutil.ProgressBar(title, int(progress.Load()), int(total.Load())))
			}
			time.Sleep(config.ProgressWaitSec * time.Second)
		}
	}()
}
