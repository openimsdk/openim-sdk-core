package common

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/openimsdk/tools/log"
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

type Cmd2Value struct {
	Cmd    string
	Value  any
	Caller string
	Ctx    context.Context
}

func sendCmd(ch chan<- Cmd2Value, value Cmd2Value, timeout time.Duration) error {
	if value.Caller == "" {
		value.Caller = GetCaller(3)
	}
	if ch == nil {
		log.ZError(value.Ctx, "sendCmd chan is nil", ErrChanNil, "caller", value.Caller, "cmd", value.Cmd, "value", value.Value)
		return ErrChanNil
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case ch <- value:
		log.ZInfo(value.Ctx, "sendCmd chan success", "caller", value.Caller, "cmd", value.Cmd, "value", value.Value)
		return nil
	case <-timer.C:
		log.ZError(value.Ctx, "sendCmd chan timeout", ErrTimeout, "caller", value.Caller, "cmd", value.Cmd, "value", value.Value)
		return ErrTimeout
	}
}

func GetCaller(skip int) string {
	pc, _, line, ok := runtime.Caller(skip)
	if !ok {
		return "runtime.caller.failed"
	}
	name := runtime.FuncForPC(pc).Name()
	if packet != "" {
		name = strings.TrimPrefix(name, packet)
	}
	return fmt.Sprintf("%s:%d", name, line)
}

type goroutine interface {
	Work(cmd Cmd2Value)
	GetCh() chan Cmd2Value
}

func DoListener(ctx context.Context, li goroutine) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("panic: %+v\n%s", r, debug.Stack())
			log.ZWarn(ctx, "DoListener panic", nil, "panic info", err)
		}
	}()

	for {
		select {
		case cmd := <-li.GetCh():
			log.ZInfo(cmd.Ctx, "recv cmd", "caller", cmd.Caller, "cmd", cmd.Cmd, "value", cmd.Value)
			li.Work(cmd)
			log.ZInfo(cmd.Ctx, "done cmd", "caller", cmd.Caller, "cmd", cmd.Cmd, "value", cmd.Value)
		case <-ctx.Done():
			log.ZInfo(ctx, "conversation done sdk logout.....")
			return
		}
	}
}
