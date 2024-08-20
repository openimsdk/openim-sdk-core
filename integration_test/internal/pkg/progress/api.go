package progress

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/reerrgroup"
	"github.com/openimsdk/tools/utils/stringutil"
)

func FuncNameBarPrint(ctx context.Context, gr *reerrgroup.Group, now, total int) *Progress {
	return FuncBarPrint(ctx, stringutil.GetFuncName(1), gr, now, total)
}

func FuncBarPrint(ctx context.Context, name string, gr *reerrgroup.Group, now, total int) *Progress {
	bar := NewBar(name, now, total, false)
	p := Start(bar)
	gr.SetAfterTasks(func() error {
		p.IncBar(bar)
		return nil
	})

	go func() {
		select {
		case <-ctx.Done():
			p.Stop()
		case <-p.done:
			return // p is done
		}
	}()

	return p
}

func Start(bar ...*Bar) *Progress {
	return StartWithMode(AutoClose|ForbiddenWrite, bar...)
}

func StartWithMode(mode proFlag, bar ...*Bar) *Progress {
	p := NewProgress(mode, 0)
	p.Start()

	p.AddBar(bar...)
	return p
}
