package common

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
)

const timeout = time.Second * 10

var (
	ErrTimeout = errors.New("send cmd timeout")
)

// Cmd2Value holds a dispatched command.
type Cmd2Value struct {
	Cmd    string
	Value  any
	Caller string
	Ctx    context.Context
}

func DispatchCmd(ctx context.Context, cmd string, val any, queue chan Cmd2Value) error {
	c2v := Cmd2Value{Cmd: cmd, Value: val, Ctx: ctx}
	return sendCmdToChan(queue, c2v, timeout)
}

func DispatchUpdateConversation(ctx context.Context, node UpdateConNode, queue chan Cmd2Value) error {
	return DispatchCmd(ctx, constant.CmdUpdateConversation, node, queue)
}

func DispatchUpdateMessage(ctx context.Context, node UpdateMessageNode, queue chan Cmd2Value) error {
	return DispatchCmd(ctx, constant.CmdUpdateMessage, node, queue)
}

func DispatchNewMessage(ctx context.Context, msg sdk_struct.CmdNewMsgComeToConversation, queue chan Cmd2Value) error {
	return DispatchCmd(ctx, constant.CmdNewMsgCome, msg, queue)
}

func DispatchMsgSyncInReinstall(ctx context.Context, msg sdk_struct.CmdMsgSyncInReinstall, queue chan Cmd2Value) error {
	return DispatchCmd(ctx, constant.CmdMsgSyncInReinstall, msg, queue)
}

func DispatchNotification(ctx context.Context, msg sdk_struct.CmdNewMsgComeToConversation, queue chan Cmd2Value) {
	_ = DispatchCmd(ctx, constant.CmdNotification, msg, queue)
}

func DispatchSyncFlag(ctx context.Context, syncFlag int, queue chan Cmd2Value) {
	_ = DispatchCmd(ctx, constant.CmdSyncFlag, sdk_struct.CmdNewMsgComeToConversation{SyncFlag: syncFlag}, queue)
}

func DispatchSyncData(ctx context.Context, queue chan Cmd2Value) {
	_ = DispatchCmd(ctx, constant.CmdSyncData, nil, queue)
}

func DispatchPushMsg(ctx context.Context, msg *sdkws.PushMessages, queue chan Cmd2Value) error {
	return DispatchCmd(ctx, constant.CmdPushMsg, msg, queue)
}

func DispatchLogout(ctx context.Context, queue chan Cmd2Value) error {
	return DispatchCmd(ctx, constant.CmdLogOut, nil, queue)
}

func DispatchConnected(ctx context.Context, queue chan Cmd2Value) error {
	return DispatchCmd(ctx, constant.CmdConnSuccesss, nil, queue)
}

func DispatchWakeUp(ctx context.Context, queue chan Cmd2Value) error {
	return DispatchCmd(ctx, constant.CmdWakeUpDataSync, nil, queue)
}

func DispatchIMSync(ctx context.Context, conversationIDs []string, queue chan Cmd2Value) error {
	return DispatchCmd(ctx, constant.CmdIMMessageSync, conversationIDs, queue)
}

type DeleteConNode struct {
	SourceID       string
	ConversationID string
	SessionType    int
}

type SyncReactionExtensionsNode struct {
	OperationID string
	Action      int
	Args        interface{}
}

type UpdateConNode struct {
	ConID  string
	Action int //1 Delete the conversation; 2 Update the latest news in the conversation or add a conversation; 3 Put a conversation on the top;
	// 4 Cancel a conversation on the top, 5 Messages are not read and set to 0, 6 New conversations
	Args interface{}
}

type UpdateMessageNode struct {
	Action int
	Args   interface{}
}

type UpdateConInfo struct {
	UserID  string
	GroupID string
}

type UpdateMessageInfo struct {
	SessionType int32
	UserID      string
	FaceURL     string
	Nickname    string
	GroupID     string
}

type SourceIDAndSessionType struct {
	SourceID    string
	SessionType int32
	FaceURL     string
	Nickname    string
}

// General send logic for channel.
func sendCmdToChan(ch chan<- Cmd2Value, value Cmd2Value, timeout time.Duration) error {
	if value.Caller == "" {
		value.Caller = GetCaller(4)
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case ch <- value:
		log.ZInfo(value.Ctx, "sendCmd channel success", "caller", value.Caller, "cmd", value.Cmd, "value", value.Value)
		return nil
	case <-timer.C:
		log.ZError(value.Ctx, "sendCmd channel timeout", ErrTimeout, "caller", value.Caller, "cmd", value.Cmd, "value", value.Value)
		return ErrTimeout
	}
}

// GetCaller returns the function and line number
func GetCaller(skip int) string {
	pc, _, line, ok := runtime.Caller(skip)
	if !ok {
		return "runtime.caller.failed"
	}
	name := runtime.FuncForPC(pc).Name()
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
