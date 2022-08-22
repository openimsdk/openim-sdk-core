package test

import (
	"open_im_sdk/internal/interaction"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	//	"open_im_sdk/internal/interaction"
	"open_im_sdk/internal/login"
	"open_im_sdk/sdk_struct"
)

func init() {
	sdk_struct.SvrConf = sdk_struct.IMConfig{Platform: 1, ApiAddr: APIADDR, WsAddr: WSADDR, DataDir: "./", LogLevel: 6, ObjectStorage: "cos"}

}

func InitMgr(num int) {
	allLoginMgr = make(map[int]*CoreNode, num)
}

type CoreNode struct {
	token             string
	userID            string
	mgr               *login.LoginMgr
	sendMsgSuccessNum uint32
	sendMsgFailedNum  uint32
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
//func TestSendCostTime() {
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
//func TestSend(idx int, text string, uidNum, intervalSleep int) {
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

//func sendPressMsg(idx int, text string, uidNum, intervalSleep int) {
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

func sendPressMsg(index int, sendId, recvID string, idx string) bool {
	return SendTextMessageOnlyForPress(idx, sendId, recvID, "opid", allLoginMgr[index].mgr.Ws())
}
func SendTextMessageOnlyForPress(text, senderID, recvID, operationID string, ws *interaction.Ws) bool {
	var wsMsgData server_api_params.MsgData
	options := make(map[string]bool, 2)
	wsMsgData.SendID = senderID
	wsMsgData.RecvID = recvID
	wsMsgData.ClientMsgID = utils.GetMsgID(senderID)
	wsMsgData.SenderPlatformID = 1
	wsMsgData.SessionType = constant.SingleChatType
	wsMsgData.MsgFrom = constant.UserMsgType
	wsMsgData.ContentType = constant.Text
	wsMsgData.Content = []byte(text)
	wsMsgData.CreateTime = utils.GetCurrentTimestampByMill()
	wsMsgData.Options = options
	wsMsgData.OfflinePushInfo = nil
	timeout := 300
	return ws.SendReqTest(&wsMsgData, constant.WSSendMsg, timeout, senderID, operationID)
}
