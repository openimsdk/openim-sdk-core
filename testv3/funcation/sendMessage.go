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
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	//	"open_im_sdk/internal/interaction"
	"open_im_sdk/internal/login"
)

func init() {
	//sdk_struct.SvrConf = sdk_struct.IMConfig{Platform: 1, ApiAddr: APIADDR, WsAddr: WSADDR, DataDir: "./", LogLevel: 6, ObjectStorage: "cos"}
	allLoginMgr = make(map[int]*CoreNode)

	SendSuccAllMsg = make(map[string]*SendRecvTime)

}

type CoreNode struct {
	token             string
	userID            string
	mgr               *login.LoginMgr
	sendMsgSuccessNum uint32
	sendMsgFailedNum  uint32
	idx               int
}

type TestSendMsgCallBack struct {
	msg         string
	OperationID string
	sendID      string
	recvID      string
	msgID       string
	sendTime    int64
	recvTime    int64
	groupID     string
}

type SendRecvTime struct {
	SendTime             int64
	SendSeccCallbackTime int64
	RecvTime             int64
	SendIDRecvID         string
}

var SendSuccAllMsg map[string]*SendRecvTime //msgid->send+recv:
var SendFailedAllMsg map[string]string
var RecvAllMsg map[string]*SendRecvTime //msgid->send+recv

func DoTestSendMsg(index int, sendId, recvID string, groupID string, idx string) {
	m := "test msg " + sendId + ":" + recvID + ":" + idx
	operationID := utils.OperationIDGenerator()
	log.Info(operationID, "CreateTextMessage  conv: ", allLoginMgr[index].mgr.Conversation(), "index: ", index)
	s, err := allLoginMgr[index].mgr.Conversation().CreateTextMessage(ctx, m)
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
	allLoginMgr[index].mgr.Conversation().SendMessage(ctx, s, recvID, groupID, &o)
	SendMsgMapLock.Lock()
	defer SendMsgMapLock.Unlock()
	x := SendRecvTime{SendTime: utils.GetCurrentTimestampByMill()}
	SendSuccAllMsg[testSendMsg.msgID] = &x
}
