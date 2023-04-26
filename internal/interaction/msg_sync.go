// Copyright © 2023 OpenIM SDK.
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

package interaction

import (
	"context"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"runtime"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
)

type SeqPair struct {
	BeginSeq uint32
	EndSeq   uint32
}

type MsgSync struct {
	db_interface.DataBase
	*Ws
	LoginUserID        string
	conversationCh     chan common.Cmd2Value
	PushMsgAndMaxSeqCh chan common.Cmd2Value // ws ch

	selfMsgSync               *SelfMsgSync
	isSyncFinished            bool
	readDiffusionGroupMsgSync *ReadDiffusionGroupMsgSync
}

func (m *MsgSync) compareSeq() {
	operationID := utils.OperationIDGenerator()
	ctx := mcontext.NewCtx(operationID)
	m.selfMsgSync.compareSeq(ctx)
	m.readDiffusionGroupMsgSync.compareSeq(operationID)
}

func (m *MsgSync) doMaxSeq(cmd common.Cmd2Value) {
	operationID := cmd.Value.(sdk_struct.CmdMaxSeqToMsgSync).OperationID
	if !m.isSyncFinished {
		m.readDiffusionGroupMsgSync.TriggerCmdNewMsgCome(nil, operationID, constant.MsgSyncBegin)
	}
	m.readDiffusionGroupMsgSync.doMaxSeq(cmd)
	m.selfMsgSync.doMaxSeq(cmd)
	if !m.isSyncFinished {
		m.readDiffusionGroupMsgSync.TriggerCmdNewMsgCome(nil, operationID, constant.MsgSyncEnd)
	}
	m.isSyncFinished = true
}

func (m *MsgSync) doPushMsg(cmd common.Cmd2Value) {
	msg := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).Msg
	switch msg.SessionType {
	case constant.SuperGroupChatType:
		m.readDiffusionGroupMsgSync.doPushMsg(cmd)
	default:
		m.selfMsgSync.doPushMsg(cmd)
	}
}

func (m *MsgSync) Work(cmd common.Cmd2Value) {
	switch cmd.Cmd {
	case constant.CmdPushMsg:
		if m.LoginStatus() == constant.Logout {
			log.Warn("", "m.LoginStatus() == constant.Logout, Goexit()")
			runtime.Goexit()
		}
		m.doPushMsg(cmd)
	case constant.CmdMaxSeq:
		if m.LoginStatus() == constant.Logout {
			log.Warn("", "m.LoginStatus() == constant.Logout, Goexit()")
			runtime.Goexit()
		}
		m.doMaxSeq(cmd)
	default:
		log.Error("", "cmd failed ", cmd.Cmd)
	}
}

func (m *MsgSync) GetCh() chan common.Cmd2Value {
	return m.PushMsgAndMaxSeqCh
}

func NewMsgSync(ctx context.Context, dataBase db_interface.DataBase, ch chan common.Cmd2Value, pushMsgAndMaxSeqCh chan common.Cmd2Value) *MsgSync {
	p := &MsgSync{DataBase: dataBase, Ws: ws, LoginUserID: loginUserID, conversationCh: ch, PushMsgAndMaxSeqCh: pushMsgAndMaxSeqCh}
	//	p.superGroupMsgSync = NewSuperGroupMsgSync(dataBase, ws, loginUserID, ch, joinedSuperGroupCh)
	p.selfMsgSync = NewSelfMsgSync(dataBase, ws, loginUserID, ch)
	p.readDiffusionGroupMsgSync = NewReadDiffusionGroupMsgSync(dataBase, ws, loginUserID, ch, joinedSuperGroupCh)
	//	p.selfMsgSync = NewSelfMsgSyncLatestModel(dataBase, ws, loginUserID, ch)
	p.compareSeq()
	go common.DoListener(p)
	return p
}
