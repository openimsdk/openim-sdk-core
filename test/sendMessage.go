package test

import (
	"github.com/gorilla/websocket"
	"github.com/jinzhu/copier"
	"open_im_sdk/internal/interaction"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"time"
)

func GenUserIDToken(userID string) (string, error) {
	return "", nil
}

func GenWs(userID, token string) *ws.Ws {
	return nil
}

func test(num int) {

	for i := 0; i < num; i++ {

	}

}

func SendTextMessage(text, senderID, recvID, operationID string, ws *interaction.Ws) bool {
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
