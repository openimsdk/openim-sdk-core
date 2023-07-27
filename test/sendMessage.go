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

package test

import (
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"

	"github.com/OpenIMSDK/protocol/sdkws"

	//	"open_im_sdk/internal/interaction"
	"open_im_sdk/internal/login"
)

func init() {
	//sdk_struct.SvrConf = sdk_struct.IMConfig{Platform: 1, ApiAddr: APIADDR, WsAddr: WSADDR, DataDir: "./", LogLevel: 6, ObjectStorage: "cos"}
	allLoginMgr = make(map[int]*CoreNode)

}

//funcation InitMgr(num int) {
//	log.Warn("", "allLoginMgr cap:  ", num)
//	allLoginMgr = make(map[int]*CoreNode, num)
//}

type CoreNode struct {
	token             string
	userID            string
	mgr               *login.LoginMgr
	sendMsgSuccessNum uint32
	sendMsgFailedNum  uint32
	idx               int
}

func addSendSuccess() {
	sendSuccessLock.Lock()
	defer sendSuccessLock.Unlock()
	sendSuccessCount++
}
func addSendFailed() {
	sendFailedLock.Lock()
	defer sendFailedLock.Unlock()
	sendFailedCount++
}

//
//funcation TestSendCostTime() {
//	GenWsConn(0)
//	sendID := allUserID[0]
//	recvID := allUserID[0]
//	for {
//		operationID := utils.OperationIDGenerator()
//		b := SendTextMessage("test", sendID, recvID, operationID, allWs[0])
//		if b {
//			log.Debug(operationID, sendID, recvID, "SendTextMessage success")
//		} else {
//			log.Error(operationID, sendID, recvID, "SendTextMessage failed")
//		}
//		time.Sleep(time.Duration(5) * time.Second)
//		log.Debug(operationID, "//////////////////////////////////")
//	}
//
//}
//funcation TestSend(idx int, text string, uidNum, intervalSleep int) {
//	for {
//		operationID := utils.OperationIDGenerator()
//		sendID := allUserID[idx]
//		recvID := allUserID[rand.Intn(uidNum)]
//		b := SendTextMessage(text, sendID, recvID, operationID, allWs[idx])
//		if b {
//			log.Debug(operationID, sendID, recvID, "SendTextMessage success")
//		} else {
//			log.Error(operationID, sendID, recvID, "SendTextMessage failed")
//		}
//		time.Sleep(time.Duration(rand.Intn(intervalSleep)) * time.Millisecond)
//	}
//}
//

//funcation sendPressMsg(idx int, text string, uidNum, intervalSleep int) {
//	for {
//		operationID := utils.OperationIDGenerator()
//		sendID := allUserID[idx]
//		recvID := allUserID[rand.Intn(uidNum)]
//		b := SendTextMessageOnlyForPress(text, sendID, recvID, operationID, allLoginMgr[idx].mgr.Ws())
//		if b {
//			log.Debug(operationID, sendID, recvID, "SendTextMessage success")
//		} else {
//			log.Error(operationID, sendID, recvID, "SendTextMessage failed ")
//		}
//		time.Sleep(time.Duration(rand.Intn(intervalSleep)) * time.Second)
//	}
//}

func sendPressMsg(index int, sendId, recvID string, groupID string, idx string) bool {

	return SendTextMessageOnlyForPress(idx, sendId, recvID, groupID, utils.OperationIDGenerator())
}
func SendTextMessageOnlyForPress(text, senderID, recvID, groupID, operationID string) bool {
	var wsMsgData sdkws.MsgData
	options := make(map[string]bool, 2)
	wsMsgData.SendID = senderID
	if groupID == "" {
		wsMsgData.RecvID = recvID
		wsMsgData.SessionType = constant.SingleChatType
	} else {
		wsMsgData.GroupID = groupID
		wsMsgData.SessionType = constant.SuperGroupChatType
	}

	wsMsgData.ClientMsgID = utils.GetMsgID(senderID)
	wsMsgData.SenderPlatformID = 1

	wsMsgData.MsgFrom = constant.UserMsgType
	wsMsgData.ContentType = constant.Text
	wsMsgData.Content = []byte(text)
	wsMsgData.CreateTime = utils.GetCurrentTimestampByMill()
	wsMsgData.Options = options
	wsMsgData.OfflinePushInfo = nil
	//timeout := 300
	log.Info(operationID, "SendReqTest begin ", wsMsgData)
	//flag := ws.SendReqTest(&wsMsgData, constant.WSSendMsg, timeout, senderID, operationID)
	//
	//if flag != true {
	//	log.Warn(operationID, "SendReqTest failed ", wsMsgData)
	//}
	return true
}
