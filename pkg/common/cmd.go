package common

import (
	"runtime/debug"
	"strings"
)

var packet string

func init() {
	build, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	packet = build.Main.Path
	if packet != "" && !strings.HasSuffix(packet, "/") {
		packet += "/"
	}
}

//type goroutine interface {
//	Work(cmd Cmd2Value)
//	GetCh() chan Cmd2Value
//}
//
//func DoListener(ctx context.Context, li goroutine) {
//	defer func() {
//		if r := recover(); r != nil {
//			err := fmt.Sprintf("panic: %+v\n%s", r, debug.Stack())
//			log.ZWarn(ctx, "DoListener panic", nil, "panic info", err)
//		}
//	}()
//
//	for {
//		select {
//		case cmd := <-li.GetCh():
//			log.ZInfo(cmd.Ctx, "recv cmd", "caller", cmd.Caller, "cmd", cmd.Cmd, "value", cmd.Value)
//			li.Work(cmd)
//			log.ZInfo(cmd.Ctx, "done cmd", "caller", cmd.Caller, "cmd", cmd.Cmd, "value", cmd.Value)
//		case <-ctx.Done():
//			log.ZInfo(ctx, "conversation done sdk logout.....")
//			return
//		}
//	}
//}
//
//func DoListenerWithEventQueue(ctx context.Context, queue *common.EventQueue, handler func(cmd Cmd2Value)) {
//	queue.ConsumeLoopWithLogger(ctx, func(ctx context.Context, event *common.Event) {
//		cmd, ok := event.Data.(Cmd2Value)
//		if !ok {
//			log.ZError(ctx, "invalid event type", nil)
//			return
//		}
//		log.ZInfo(cmd.Ctx, "recv cmd", "caller", cmd.Caller, "cmd", cmd.Cmd, "value", cmd.Value)
//		handler(cmd)
//		log.ZInfo(cmd.Ctx, "done cmd", "caller", cmd.Caller, "cmd", cmd.Cmd, "value", cmd.Value)
//	}, func(event *common.Event, err error) {
//		if err != nil && ctx.Err() != nil {
//			log.ZInfo(ctx, "conversation done sdk logout.....")
//		} else if err != nil {
//			log.ZError(ctx, "Pop error", err)
//		}
//	})
//}
