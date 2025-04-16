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
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"

	"github.com/openimsdk/protocol/sdkws"
)

const timeout = time.Millisecond * 10000

var (
	ErrChanNil = errs.New("channal == nil")
	ErrTimeout = errors.New("send cmd timeout")
)

func TriggerCmdNewMsgCome(ctx context.Context, msg sdk_struct.CmdNewMsgComeToConversation, queue *EventQueue) error {

	c2v := Cmd2Value{Cmd: constant.CmdNewMsgCome, Value: msg, Ctx: ctx}
	return sendCmd2(ctx, queue, c2v, 5, timeout)
}

func TriggerCmdMsgSyncInReinstall(ctx context.Context, msg sdk_struct.CmdMsgSyncInReinstall, queue *EventQueue) error {

	c2v := Cmd2Value{Cmd: constant.CmdMsgSyncInReinstall, Value: msg, Ctx: ctx}
	return sendCmd2(ctx, queue, c2v, 5, timeout)
}

func TriggerCmdNotification(ctx context.Context, msg sdk_struct.CmdNewMsgComeToConversation, queue *EventQueue) {
	c2v := Cmd2Value{Cmd: constant.CmdNotification, Value: msg, Ctx: ctx}
	err := sendCmd2(ctx, queue, c2v, 5, timeout)
	if err != nil {
		log.ZWarn(ctx, "TriggerCmdNotification error", err, "msg", msg)
	}
}

func TriggerCmdSyncFlag(ctx context.Context, syncFlag int, queue *EventQueue) {
	c2v := Cmd2Value{Cmd: constant.CmdSyncFlag, Value: sdk_struct.CmdNewMsgComeToConversation{SyncFlag: syncFlag}, Ctx: ctx}
	err := sendCmd2(ctx, queue, c2v, 5, timeout)
	if err != nil {
		log.ZWarn(ctx, "TriggerCmdNotification error", err, "syncFlag", syncFlag)
	}
}

func TriggerCmdSyncFlagAndConversationMetaData(ctx context.Context, syncFlag int, seqs map[string]*msg.Seqs, queue *EventQueue) {
	c2v := Cmd2Value{Cmd: constant.CmdSyncFlag, Value: sdk_struct.CmdNewMsgComeToConversation{
		Seqs:     seqs,
		SyncFlag: syncFlag,
	}, Ctx: ctx}
	err := sendCmd2(ctx, queue, c2v, 5, timeout)
	if err != nil {
		log.ZWarn(ctx, "TriggerCmdNotification error", err, "syncFlag", syncFlag)
	}
}

func TriggerCmdWakeUpDataSync(ctx context.Context, ch chan Cmd2Value) error {
	c2v := Cmd2Value{Cmd: constant.CmdWakeUpDataSync, Value: nil, Ctx: ctx}
	return sendCmd(ch, c2v, timeout)
}

func TriggerCmdIMMessageSync(ctx context.Context, conversationIDs []string, ch chan Cmd2Value) error {
	c2v := Cmd2Value{Cmd: constant.CmdIMMessageSync, Value: conversationIDs, Ctx: ctx}
	return sendCmd(ch, c2v, timeout)
}

func TriggerCmdSyncData(ctx context.Context, queue *EventQueue) {
	c2v := Cmd2Value{Cmd: constant.CmdSyncData, Value: nil, Ctx: ctx}
	err := sendCmd2(ctx, queue, c2v, 5, timeout)
	if err != nil {
		log.ZWarn(ctx, "TriggerCmdSyncData error", err)
	}
}

//func TriggerCmdUpdateConversation(ctx context.Context, node UpdateConNode, conversationCh chan<- Cmd2Value) error {
//	c2v := Cmd2Value{
//		Cmd:   constant.CmdUpdateConversation,
//		Value: node,
//		Ctx:   ctx,
//	}
//	err := sendCmd2(ctx, queue, c2v, 5, timeout)
//	if err != nil {
//		_, _ = len(conversationCh), cap(conversationCh)
//	}
//	return err
//}

func TriggerCmdUpdateConversation(ctx context.Context, node UpdateConNode, queue *EventQueue) error {
	c2v := Cmd2Value{
		Cmd:   constant.CmdUpdateConversation,
		Value: node,
		Ctx:   ctx,
	}
	return sendCmd2(ctx, queue, c2v, 5, timeout)
}

func TriggerCmdUpdateMessage(ctx context.Context, node UpdateMessageNode, queue *EventQueue) error {
	c2v := Cmd2Value{
		Cmd:   constant.CmdUpdateMessage,
		Value: node,
		Ctx:   ctx,
	}
	return sendCmd2(ctx, queue, c2v, 5, timeout)
}

// TriggerCmdPushMsg Push message, msg for msgData slice
func TriggerCmdPushMsg(ctx context.Context, msg *sdkws.PushMessages, ch chan Cmd2Value) error {
	c2v := Cmd2Value{Cmd: constant.CmdPushMsg, Value: msg, Ctx: ctx}
	return sendCmd(ch, c2v, timeout)
}

func TriggerCmdLogOut(ctx context.Context, ch chan Cmd2Value) error {
	c2v := Cmd2Value{Cmd: constant.CmdLogOut, Ctx: ctx}
	return sendCmd(ch, c2v, timeout)
}

// TriggerCmdConnected Connection success trigger
func TriggerCmdConnected(ctx context.Context, ch chan Cmd2Value) error {
	c2v := Cmd2Value{Cmd: constant.CmdConnSuccesss, Value: nil, Ctx: ctx}
	return sendCmd(ch, c2v, timeout)
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
