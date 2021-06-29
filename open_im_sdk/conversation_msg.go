package open_im_sdk

import (
	"encoding/json"
	"fmt"
)

type ConversationListener struct {
	ConversationListener OnConversationListener
	MsgListenerList      []OnAdvancedMsgListener
	ch                   chan cmd2Value
}

func (con ConversationListener) getCh() chan cmd2Value {
	return con.ch
}

func (con *ConversationListener) doMsg(c2v cmd2Value) {
	fmt.Println("msg come herrr 1111111")
	var conversationID string
	var groupID, userID string
	if con.ConversationListener == nil || con.MsgListenerList == nil {
		fmt.Println("not set conversationListener or MsgListenerList", len(con.MsgListenerList))
		return
	}
	var messages []*MsgStruct
	MsgList := c2v.Value.(ArrMsg)
	var msgReadList []*MsgData
	for _, v := range MsgList.Data {
		msg := &MsgStruct{
			SendID:      v.SendID,
			RecvID:      v.RecvID,
			SessionType: v.SessionType,
			MsgFrom:     v.MsgFrom,
			ContentType: v.ContentType,
			ServerMsgID: v.ServerMsgID,
			ClientMsgID: v.ServerMsgID,
			Content:     v.Content,
			SendTime:    v.SendTime,
			Seq:         v.Seq,
			PlatformID:  v.SenderPlatformID,
			Status:      MsgStatusSendSuccess,
			IsRead:      false,
		}
		if v.SendID == LoginUid { //seq對齊消息 Messages sent by myself
			if judgeMessageIfExists(msg) { //if  sent through  this terminal
				err := updateMessageSeq(msg)
				if err != nil {
					fmt.Println("err", err.Error())
				}
			} else { //同步消息       send through  other terminal
				_ = insertPushMessageToChatLog(msg)
				switch v.SessionType {
				case SingleChatType:
					conversationID = GetConversationIDBySessionType(v.RecvID, SingleChatType)
					userID = v.RecvID
				case GroupChatType:
					conversationID = GetConversationIDBySessionType(v.RecvID, GroupChatType)
					groupID = v.RecvID
				}
				conversation := ConversationStruct{
					ConversationID:    conversationID,
					ConversationType:  int(v.SessionType),
					UserID:            userID,
					GroupID:           groupID,
					RecvMsgOpt:        1,
					LatestMsg:         structToJsonString(msg),
					LatestMsgSendTime: msg.SendTime,
				}
				switch {
				case msg.ContentType <= AcceptFriendApplicationTip:
					messages = append(messages, msg)
					con.doUpdateConversation(cmd2Value{Value: updateConNode{conversationID, AddConOrUpLatMsg,
						conversation}})

				default:

				}

				//}

			}
		} else { //他人發的
			_ = insertPushMessageToChatLog(msg)
			switch v.SessionType {
			case SingleChatType:
				conversationID = GetConversationIDBySessionType(v.SendID, SingleChatType)
				userID = v.SendID
			case GroupChatType:
				conversationID = GetConversationIDBySessionType(v.RecvID, GroupChatType)
				groupID = v.RecvID
			}
			conversation := ConversationStruct{
				ConversationID:    conversationID,
				ConversationType:  int(v.SessionType),
				UserID:            userID,
				GroupID:           groupID,
				RecvMsgOpt:        1,
				LatestMsg:         structToJsonString(msg),
				LatestMsgSendTime: msg.SendTime,
			}
			switch {
			case msg.ContentType <= AcceptFriendApplicationTip:
				fmt.Println("tttttttttttttttttttt", msg.ContentType)
				con.doUpdateConversation(cmd2Value{Value: updateConNode{conversationID, IncrUnread, ""}})
				con.doUpdateConversation(cmd2Value{Value: updateConNode{conversationID, AddConOrUpLatMsg,
					conversation}})
				messages = append(messages, msg)

			default:

			}
			if msg.ContentType == C2CMessageAsRead {
				//update read state
				msgReadList = append(msgReadList, &v)

			}
		}
	}

	con.doMsgReadState(msgReadList)
	fmt.Println("length msgListenerList", con.MsgListenerList, "length message", len(messages), "msgListenerLen", len(con.MsgListenerList))
	for _, v := range con.MsgListenerList {
		for _, w := range messages {
			//De-analyze data
			err := msgHandleByContentType(w)
			if err != nil {
				fmt.Println("Parsing data error")
			} else {
				if v != nil {
					fmt.Println("msgListener,OnRecvNewMessage")
					v.OnRecvNewMessage(structToJsonString(w))
				} else {
					fmt.Println("set msgListener is err ")
				}

			}
		}
	}

}
func insertMessageAndSyncConversation() {

}

func (con *ConversationListener) doDeleteConversation(c2v cmd2Value) {
	if con.ConversationListener == nil {
		fmt.Println("not set conversationListener")
		return
	}
	node := c2v.Value.(deleteConNode)
	_ = deleteConversationModel(node.ConversationID)
	maxSeq, err := getLocalMaxSeqModel()
	if err != nil {
		return
	}
	_ = deleteMessageByConversationModel(node.SourceID, maxSeq)
	err, list := getAllConversationListModel()
	if err == nil {
		if list != nil {
			con.ConversationListener.OnConversationChanged(structToJsonString(list))
		} else {
			con.ConversationListener.OnConversationChanged(structToJsonString([]ConversationStruct{}))
		}
		totalUnreadCount, err := getTotalUnreadMsgCountModel()
		if err == nil {
			con.ConversationListener.OnTotalUnreadMessageCountChanged(totalUnreadCount)
		}

	}

}

func (con *ConversationListener) doMsgReadState(msgReadList []*MsgData) {
	var messageReceiptResp []MessageReceipt
	for _, rd := range msgReadList {
		var msgIdList []string
		err := json.Unmarshal([]byte(rd.Content), &msgIdList)
		if err != nil {
			sdkLog("unmarshal failed, err : ", err.Error())
			return
		}

		var msgIdListStatusOK []string
		for _, v := range msgIdList {
			err := setMessageHasReadByMsgID(v)
			if err != nil {
				continue
			}
			msgIdListStatusOK = append(msgIdListStatusOK, v)
		}

		if len(msgIdListStatusOK) > 0 {
			var msgRt MessageReceipt
			msgRt.ContentType = rd.ContentType
			msgRt.MsgFrom = rd.MsgFrom
			msgRt.ReadTime = rd.SendTime
			msgRt.UserID = rd.SendID
			msgRt.SessionType = rd.SessionType
			msgRt.MsgIdList = msgIdListStatusOK
			messageReceiptResp = append(messageReceiptResp, msgRt)
		}
	}

	if len(messageReceiptResp) > 0 {
		jsonResult, err := json.Marshal(messageReceiptResp)
		if err != nil {
			sdkLog("marshal failed, ", err.Error())
			return
		}
		for _, v := range con.MsgListenerList {
			sdkLog("OnRecvC2CReadReceipt: ", string(jsonResult))
			v.OnRecvC2CReadReceipt(string(jsonResult))
		}
	}
}

func (con *ConversationListener) doUpdateConversation(c2v cmd2Value) {
	if con.ConversationListener == nil {
		fmt.Println("not set conversationListener")
		return
	}
	node := c2v.Value.(updateConNode)
	switch node.Action {
	case ConAndUnreadChange:
		err, list := getAllConversationListModel()
		if err == nil {
			if list == nil {
				con.ConversationListener.OnConversationChanged(structToJsonString([]ConversationStruct{}))
			} else {
				con.ConversationListener.OnConversationChanged(structToJsonString(list))

			}
			totalUnreadCount, err := getTotalUnreadMsgCountModel()
			if err == nil {
				con.ConversationListener.OnTotalUnreadMessageCountChanged(totalUnreadCount)
			}
		}
	case AddConOrUpLatMsg:
		c := node.Args.(ConversationStruct)
		if judgeConversationIfExists(node.ConId) {
			err := setConversationLatestMsgModel(&c, node.ConId)
			if err != nil {
				sdkLog("err: ", err)
			} else {
				err, list := getAllConversationListModel()
				if err != nil {
					sdkLog("doUpdateConversation database err:", err.Error())
				} else {
					if list == nil {
						con.ConversationListener.OnConversationChanged(structToJsonString([]ConversationStruct{}))
					} else {
						con.ConversationListener.OnConversationChanged(structToJsonString(list))

					}
				}
			}
		} else {
			_ = addConversationOrUpdateLatestMsg(&c, node.ConId)
			var list []*ConversationStruct
			list = append(list, &c)
			con.ConversationListener.OnNewConversation(structToJsonString(list))

		}

	case UnreadCountSetZero:
		if err := setConversationUnreadCount(0, node.ConId); err != nil {
		} else {
			err, list := getAllConversationListModel()
			if err == nil {
				if list == nil {
					con.ConversationListener.OnConversationChanged(structToJsonString([]ConversationStruct{}))
				} else {
					con.ConversationListener.OnConversationChanged(structToJsonString(list))

				}
				totalUnreadCount, err := getTotalUnreadMsgCountModel()
				if err == nil {
					con.ConversationListener.OnTotalUnreadMessageCountChanged(totalUnreadCount)
				}

			}
		}
	case ConChange:
		err, list := getAllConversationListModel()
		if err != nil {
			sdkLog("doUpdateConversation database err:", err.Error())
		} else {
			if list == nil {
				con.ConversationListener.OnConversationChanged(structToJsonString([]ConversationStruct{}))
			} else {
				con.ConversationListener.OnConversationChanged(structToJsonString(list))

			}
		}
	case IncrUnread:
		err := incrConversationUnreadCount(node.ConId)
		if err != nil {
			sdkLog("incrConversationUnreadCount database err:", err.Error())
			return
		}
		totalUnreadCount, err := getTotalUnreadMsgCountModel()
		if err == nil {
			con.ConversationListener.OnTotalUnreadMessageCountChanged(totalUnreadCount)
		}

	}

}

func (con ConversationListener) work(c2v cmd2Value) {
	switch c2v.Cmd {
	case CmdDeleteConversation:
		con.doDeleteConversation(c2v)

	case CmdNewMsgCome:
		con.doMsg(c2v)
	case CmdUpdateConversation:
		con.doUpdateConversation(c2v)

	}
}
