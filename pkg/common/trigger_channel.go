// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
)

const timeout = time.Second * 10

var (
	ErrTimeout = errors.New("send cmd timeout")
)

// CmdPriorityMap centralized priority definition
var CmdPriorityMap = map[string]int{
	constant.CmdNewMsgCome:         5,
	constant.CmdUpdateConversation: 4,
	constant.CmdMsgSyncInReinstall: 1,
	constant.CmdNotification:       1,
	constant.CmdSyncFlag:           1,
	constant.CmdSyncData:           1,
	constant.CmdUpdateMessage:      1,

	constant.CmdLogOut:         1,
	constant.CmdPushMsg:        1,
	constant.CmdIMMessageSync:  1,
	constant.CmdConnSuccesss:   1,
	constant.CmdWakeUpDataSync: 1,
}

// Cmd2Value holds a dispatched command.
type Cmd2Value struct {
	Cmd    string
	Value  any
	Caller string
	Ctx    context.Context
}

func DispatchCmd(ctx context.Context, cmd string, val any, queue *EventQueue) error {
	c2v := Cmd2Value{Cmd: cmd, Value: val, Ctx: ctx}
	return sendCmdToQueue(ctx, queue, c2v, timeout)
}

func DispatchUpdateConversation(ctx context.Context, node UpdateConNode, queue *EventQueue) error {
	return DispatchCmd(ctx, constant.CmdUpdateConversation, node, queue)
}

func DispatchUpdateMessage(ctx context.Context, node UpdateMessageNode, queue *EventQueue) error {
	return DispatchCmd(ctx, constant.CmdUpdateMessage, node, queue)
}

func DispatchNewMessage(ctx context.Context, msg sdk_struct.CmdNewMsgComeToConversation, queue *EventQueue) error {
	return DispatchCmd(ctx, constant.CmdNewMsgCome, msg, queue)
}

func DispatchMsgSyncInReinstall(ctx context.Context, msg sdk_struct.CmdMsgSyncInReinstall, queue *EventQueue) error {
	return DispatchCmd(ctx, constant.CmdMsgSyncInReinstall, msg, queue)
}

func DispatchNotification(ctx context.Context, msg sdk_struct.CmdNewMsgComeToConversation, queue *EventQueue) {
	_ = DispatchCmd(ctx, constant.CmdNotification, msg, queue)
}

func DispatchSyncFlag(ctx context.Context, syncFlag int, queue *EventQueue) {
	_ = DispatchCmd(ctx, constant.CmdSyncFlag, sdk_struct.CmdNewMsgComeToConversation{SyncFlag: syncFlag}, queue)
}

func DispatchSyncFlagWithMeta(ctx context.Context, syncFlag int, seqs map[string]*msg.Seqs, queue *EventQueue) {
	_ = DispatchCmd(ctx, constant.CmdSyncFlag, sdk_struct.CmdNewMsgComeToConversation{Seqs: seqs, SyncFlag: syncFlag}, queue)
}

func DispatchSyncData(ctx context.Context, queue *EventQueue) {
	_ = DispatchCmd(ctx, constant.CmdSyncData, nil, queue)
}

// Legacy channel-based dispatchers

func DispatchPushMsg(ctx context.Context, msg *sdkws.PushMessages, ch chan Cmd2Value) error {
	return sendCmdToChan(ch, Cmd2Value{Cmd: constant.CmdPushMsg, Value: msg, Ctx: ctx}, timeout)
}

func DispatchLogout(ctx context.Context, ch chan Cmd2Value) error {
	return sendCmdToChan(ch, Cmd2Value{Cmd: constant.CmdLogOut, Ctx: ctx}, timeout)
}

func DispatchConnected(ctx context.Context, ch chan Cmd2Value) error {
	return sendCmdToChan(ch, Cmd2Value{Cmd: constant.CmdConnSuccesss, Ctx: ctx}, timeout)
}

func DispatchWakeUp(ctx context.Context, ch chan Cmd2Value) error {
	return sendCmdToChan(ch, Cmd2Value{Cmd: constant.CmdWakeUpDataSync, Ctx: ctx}, timeout)
}

func DispatchIMSync(ctx context.Context, conversationIDs []string, ch chan Cmd2Value) error {
	return sendCmdToChan(ch, Cmd2Value{Cmd: constant.CmdIMMessageSync, Value: conversationIDs, Ctx: ctx}, timeout)
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

// General send logic for event queue.
func sendCmdToQueue(ctx context.Context, queue *EventQueue, value Cmd2Value, timeout time.Duration) error {
	if value.Caller == "" {
		value.Caller = GetCaller(4)
	}

	priority, ok := CmdPriorityMap[value.Cmd]
	if !ok {
		priority = 1 // default
	}
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	_, err := queue.ProduceWithContext(ctxWithTimeout, value, priority)
	if err != nil {
		log.ZError(ctx, "sendCmd queue produce error", err, "caller", value.Caller, "cmd", value.Cmd, "value", value.Value)
		return err
	}
	log.ZInfo(ctx, "sendCmd queue produce success", "caller", value.Caller, "cmd", value.Cmd, "value", value.Value)
	return nil
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
