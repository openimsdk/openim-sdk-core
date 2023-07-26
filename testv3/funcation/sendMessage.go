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

package funcation

import (
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"strconv"

	"github.com/OpenIMSDK/protocol/sdkws"
)

func init() {
	AllLoginMgr = make(map[string]*CoreNode)
	SendSuccAllMsg = make(map[string]*SendRecvTime)
}

var SendSuccAllMsg map[string]*SendRecvTime // msgid->send+recv:
var SendFailedAllMsg map[string]string
var RecvAllMsg map[string]*SendRecvTime // msgid->send+recv

// 基准函数不应该做模拟，这一部分逻辑应该放在 test 文件中自行模拟
func DoTestSendMsg(index int, sendId, recvID string, groupID string, idx string) {
	m := "test msg " + sendId + ":" + recvID + ":" + idx
	operationID := utils.OperationIDGenerator()
	log.Info(operationID, "CreateTextMessage  conv: ", AllLoginMgr[strconv.Itoa(index)].Mgr.Conversation(), "index: ", index)
	s, err := AllLoginMgr[strconv.Itoa(index)].Mgr.Conversation().CreateTextMessage(ctx, m)
	if err != nil {
		log.Error(operationID, "CreateTextMessage", err)
		return
	}

	testSendMsg := TestSendMsgCallBack{
		OperationID: operationID,
		sendTime:    utils.GetCurrentTimestampByMill(),
		sendID:      sendId,
		recvID:      recvID,
		groupID:     groupID,
		msgID:       s.ClientMsgID,
	}
	o := sdkws.OfflinePushInfo{Title: "title", Desc: "desc"}

	log.Info(operationID, "SendMessage", sendId, recvID, groupID, testSendMsg.msgID, index)
	// 如果 recvID 为空 代表发送群聊消息，反之
	AllLoginMgr[strconv.Itoa(index)].Mgr.Conversation().SendMessage(ctx, s, recvID, groupID, &o)
	SendMsgMapLock.Lock()
	defer SendMsgMapLock.Unlock()
	x := SendRecvTime{SendTime: utils.GetCurrentTimestampByMill()}
	SendSuccAllMsg[testSendMsg.msgID] = &x
}
