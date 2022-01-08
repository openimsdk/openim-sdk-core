package interaction

import (
	"github.com/golang/protobuf/proto"
	"net/http"
	"open_im_sdk/pkg/utils"
	"os"
	"time"
)

type Ws struct {
	WsRespAsyn
	WsConn
}

func (ws *Ws) WaitResp(ch chan GeneralWsResp, timeout int) (*GeneralWsResp, error) {

}

func (ws *Ws) SendReqWaitResp(buff []byte, reqIdentifier int32, timeout int, SenderID string) (*GeneralWsResp, error) {
	var wsReq GeneralWsReq
	wsReq.ReqIdentifier = reqIdentifier
	wsReq.OperationID = utils.OperationIDGenerator()
	msgIncr, ch := ws.AddCh(SenderID)
	wsReq.SendID = SenderID
	//wsReq.Token = u.token
	wsReq.MsgIncr = msgIncr
	wsReq.Data = buff
	err, _ := ws.writeBinaryMsg(wsReq)
	if err != nil {
		return nil, err
	}
	select {
	case r := <-ch:

		if r.ErrCode != 0 {
			return nil, err
		} else {
			callback.OnSuccess("")
			callback.OnProgress(100)
			for _, v := range delFile {
				err := os.Remove(v)
				if err != nil {
					sdkLog("remove failed,", err.Error(), v)
				}
				sdkLog("remove file: ", v)
			}
			var sendMsgResp UserSendMsgResp
			err = proto.Unmarshal(r.Data, &sendMsgResp)
			if err != nil {
				sdkLog("Unmarshal failed ", err.Error())
			}
			_ = u.updateMessageTimeAndMsgIDStatus(sendMsgResp.ClientMsgID, sendMsgResp.SendTime, MsgStatusSendSuccess)

			s.ServerMsgID = sendMsgResp.ServerMsgID
			s.SendTime = sendMsgResp.SendTime
			s.Status = MsgStatusSendSuccess
			c.LatestMsg = structToJsonString(s)
			c.LatestMsgSendTime = s.SendTime
			_ = u.triggerCmdUpdateConversation(updateConNode{conversationID, AddConOrUpLatMsg,
				c})
			_ = u.triggerCmdUpdateConversation(updateConNode{conversationID, ConChange, ""})
		}
		breakFlag = 1
	case <-time.After(time.Second * time.Duration(timeout)):
		var flag bool
		sdkLog("ws ch recvMsg err: ", wsReq.OperationID)
		if connSend != u.conn {
			sdkLog("old conn != current conn  ", connSend, u.conn)
			flag = false // error
		} else {
			flag = false //error
			for tr := 0; tr < 3; tr++ {
				err = u.sendPingMsg()
				if err != nil {
					sdkLog("sendPingMsg failed ", wsReq.OperationID, err.Error(), tr)
					time.Sleep(time.Duration(30) * time.Second)
				} else {
					sdkLog("sendPingMsg ok ", wsReq.OperationID)
					flag = true //wait continue
					break
				}
			}
		}
		if flag == false {
			callback.OnError(http.StatusRequestTimeout, http.StatusText(http.StatusRequestTimeout))
			u.sendMessageFailedHandle(&s, &c, conversationID)
			sdkLog("onError callback ", wsReq.OperationID)
			breakFlag = 1
			break
		} else {
			sdkLog("wait resp continue", wsReq.OperationID)
			breakFlag = 0
			continue
		}
	}

}
