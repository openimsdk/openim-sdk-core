package conversation_msg

import (
	"encoding/json"
	"github.com/jinzhu/copier"
	common2 "open_im_sdk/internal/common"
	"open_im_sdk/internal/friend"
	"open_im_sdk/internal/group"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/internal/user"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

type Conversation struct {
	*ws.Ws
	db                   *db.DataBase
	p                    *ws.PostApi
	ConversationListener OnConversationListener
	msgListener          OnAdvancedMsgListener
	ch                   chan common.Cmd2Value
	loginUserID          string
	platformID           int32
	DataDir              string
	friend               *friend.Friend
	group                *group.Group
	user                 *user.User
	common2.ObjectStorage
}

func (c *Conversation) SetMsgListener(msgListener OnAdvancedMsgListener) {
	c.msgListener = msgListener
}
func NewConversation(ws *ws.Ws, db *db.DataBase, p *ws.PostApi,
	ch chan common.Cmd2Value, loginUserID string, platformID int32, dataDir string,
	friend *friend.Friend, group *group.Group, user *user.User,
	objectStorage common2.ObjectStorage) *Conversation {
	return &Conversation{Ws: ws, db: db, p: p, ch: ch, loginUserID: loginUserID, platformID: platformID, DataDir: dataDir, friend: friend, group: group, user: user, ObjectStorage: objectStorage}
}

//func NewConversation() *Conversation {
//	return &Conversation{}
//}

func (c *Conversation) Init(ws *ws.Ws, db *db.DataBase, ch chan common.Cmd2Value, loginUserID string, friend *friend.Friend, group *group.Group, user *user.User) {
	c.Ws = ws
	c.db = db
	c.ch = ch
	c.loginUserID = c.loginUserID
	c.friend = friend
	c.group = group
	c.user = user
	go common.DoListener(c)
}

func (c *Conversation) GetCh() chan common.Cmd2Value {
	return c.ch
}

func (c *Conversation) doMsgNew(c2v common.Cmd2Value) {
	allMsg := c2v.Value.([]*server_api_params.MsgData)
	operationID := utils.OperationIDGenerator()
	if c.msgListener == nil {
		log.Error(operationID, "not set c MsgListenerList")
		return
	}
	var insertMsg []*db.LocalChatLog
	var exceptionMsg []*db.LocalErrChatLog
	var newMessages, msgReadList, msgRevokeList []*sdk_struct.MsgStruct
	var isUnreadCount, isConversationUpdate, isHistory bool
	conversationChangedSet := make(map[string]db.LocalConversation)
	newConversationSet := make(map[string]db.LocalConversation)
	log.Info(operationID, "do Msg come here")
	for _, v := range allMsg {
		isHistory = utils.GetSwitchFromOptions(v.Options, constant.IsHistory)
		isUnreadCount = utils.GetSwitchFromOptions(v.Options, constant.IsUnreadCount)
		isConversationUpdate = utils.GetSwitchFromOptions(v.Options, constant.IsConversationUpdate)
		msg := new(sdk_struct.MsgStruct)
		copier.Copy(msg, v)
		msg.Content = string(v.Content)
		msg.Status = constant.MsgStatusSendSuccess
		msg.IsRead = false
		log.Info(operationID, "new msg, seq, ServerMsgID, ClientMsgID", msg.Seq, msg.ServerMsgID, msg.ClientMsgID)
		//De-analyze data
		err := c.msgHandleByContentType(msg)
		if err != nil {
			log.Error(operationID, "Parsing data error:", err.Error())
			continue
		}
		if msg.ClientMsgID == "" {
			exceptionMsg = append(exceptionMsg, c.msgStructToLocalErrChatLog(msg))
			continue
		}
		switch v.SessionType {
		case constant.SingleChatType:
			if v.ContentType > constant.FriendNotificationBegin && v.ContentType < constant.FriendNotificationEnd {
				c.friend.DoNotification(v)
				log.Info("internal", "DoFriendMsg SingleChatType", v)
			} else if v.ContentType > constant.UserNotificationBegin && v.ContentType < constant.UserNotificationEnd {
				c.user.DoNotification(v)
			}
		case constant.GroupChatType:
			if v.ContentType > constant.GroupNotificationBegin && v.ContentType < constant.GroupNotificationEnd {
				c.group.DoNotification(v)
				log.Info(operationID, "DoGroupMsg SingleChatType", v)
			}
		}
		if v.SendID == c.loginUserID { //seq  Messages sent by myself  //if  sent through  this terminal
			m, err := c.db.GetMessage(msg.ClientMsgID)
			if err == nil && m != nil {
				log.Info("internal", "have message", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, *msg)
				if m.Seq == 0 {
					insertMsg = append(insertMsg, c.msgStructToLocalChatLog(msg))
				} else {
					exceptionMsg = append(exceptionMsg, c.msgStructToLocalErrChatLog(msg))

				}
			} else { //      send through  other terminal
				log.Info(operationID, "sync message", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, *msg)
				lc := db.LocalConversation{
					ConversationType:  v.SessionType,
					LatestMsg:         utils.StructToJsonString(msg),
					LatestMsgSendTime: msg.SendTime,
				}
				switch v.SessionType {
				case constant.SingleChatType:
					lc.ConversationID = c.GetConversationIDBySessionType(v.RecvID, constant.SingleChatType)
					lc.UserID = v.RecvID
					//localUserInfo,_ := c.user.GetLoginUser()
					//c.FaceURL = localUserInfo.FaceUrl
					//c.ShowName = localUserInfo.Nickname
				case constant.GroupChatType:
					lc.GroupID = v.GroupID
					lc.ConversationID = c.GetConversationIDBySessionType(lc.GroupID, constant.GroupChatType)
					//faceUrl, name, err := u.getGroupNameAndFaceUrlByUid(c.GroupID)
					//if err != nil {
					//	utils.sdkLog("getGroupNameAndFaceUrlByUid err:", err)
					//} else {
					//	c.ShowName = name
					//	c.FaceURL = faceUrl
					//}
				}
				if isUnreadCount {
					lc.UnreadCount = 1
				}
				if isConversationUpdate {
					c.updateConversation(&lc, conversationChangedSet, newConversationSet)
				} else {
					msg.Status = constant.MsgStatusFiltered
				}
				if isHistory {
					insertMsg = append(insertMsg, c.msgStructToLocalChatLog(msg))
				}
				newMessages = append(newMessages, msg)

			}
		} else { //Sent by others
			if b, _ := c.db.MessageIfExists(msg.ClientMsgID); !b { //Deduplication operation
				lc := db.LocalConversation{
					ConversationType:  v.SessionType,
					LatestMsg:         utils.StructToJsonString(msg),
					LatestMsgSendTime: msg.SendTime,
				}

				switch v.SessionType {
				case constant.SingleChatType:
					lc.ConversationID = c.GetConversationIDBySessionType(v.SendID, constant.SingleChatType)
					lc.UserID = v.SendID
					lc.ShowName = msg.SenderNickname
					lc.FaceURL = msg.SenderFaceURL
				case constant.GroupChatType:
					lc.GroupID = v.GroupID
					lc.ConversationID = c.GetConversationIDBySessionType(lc.GroupID, constant.GroupChatType)
					//faceUrl, name, err := u.getGroupNameAndFaceUrlByUid(c.GroupID)
					//if err != nil {
					//	utils.sdkLog("getGroupNameAndFaceUrlByUid err:", err)
					//} else {
					//	c.ShowName = name
					//	c.FaceURL = faceUrl
					//}
				}
				if isUnreadCount {
					lc.UnreadCount = 1
				}

				//u.doUpdateConversation(cmd2Value{Value: updateConNode{c.ConversationID, UpdateFaceUrlAndNickName, c}})
				if isConversationUpdate {
					c.updateConversation(&lc, conversationChangedSet, newConversationSet)
					newMessages = append(newMessages, msg)
				} else {
					msg.Status = constant.MsgStatusFiltered
				}
				if isHistory {
					insertMsg = append(insertMsg, c.msgStructToLocalChatLog(msg))
				}
				if msg.ContentType == constant.Revoke {
					msgRevokeList = append(msgRevokeList, msg)
				}
			} else {
				exceptionMsg = append(exceptionMsg, c.msgStructToLocalErrChatLog(msg))
			}
		}
	}
	//Normal message storage
	err1 := c.db.BatchInsertMessageList(insertMsg)
	if err1 != nil {
		log.Error(operationID, "insert normal message err  :", err1.Error())
	}
	//Exception message storage
	if exceptionMsg != nil {
		err2 := c.db.BatchInsertExceptionMsgToErrorChatLog(exceptionMsg)
		if err2 != nil {
			log.Error(operationID, "insert err message err  :", err2.Error())

		}
	}

	//Changed conversation storage
	err3 := c.db.BatchUpdateConversationList(mapConversationToList(conversationChangedSet))
	if err3 != nil {
		log.Error(operationID, "insert changed conversation err :", err3.Error())
	}
	//New conversation storage
	err4 := c.db.BatchInsertConversationList(mapConversationToList(newConversationSet))
	if err4 != nil {
		log.Error(operationID, "insert new conversation err:", err4.Error())

	}

	c.doMsgReadState(msgReadList)
	c.revokeMessage(msgRevokeList)
	c.newMessage(newMessages)
	//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", ConChange, ""}})
	log.Info(operationID, "trigger map is :", newConversationSet, conversationChangedSet)
	//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewCon, mapKeyToStringList(newConversationSet)}})
	//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, mapKeyToStringList(conversationChangSet)}})
	if len(conversationChangedSet) != 0 {
		c.ConversationListener.OnConversationChanged(utils.StructToJsonString(mapConversationToList(conversationChangedSet)))
	}
	if len(newConversationSet) != 0 {
		c.ConversationListener.OnNewConversation(utils.StructToJsonString(mapConversationToList(newConversationSet)))
	}
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{"", constant.TotalUnreadMessageChanged, ""}})
	//sdkLog("length msgListenerList", u.MsgListenerList, "length message", len(newMessages), "msgListenerLen", len(u.MsgListenerList))

}
func (c *Conversation) msgStructToLocalChatLog(m *sdk_struct.MsgStruct) *db.LocalChatLog {
	var lc db.LocalChatLog
	copier.Copy(&lc, m)
	return &lc
}
func (c *Conversation) msgStructToLocalErrChatLog(m *sdk_struct.MsgStruct) *db.LocalErrChatLog {
	var lc db.LocalErrChatLog
	copier.Copy(&lc, m)
	return &lc
}

func (c *Conversation) revokeMessage(msgRevokeList []*sdk_struct.MsgStruct) {
	for _, w := range msgRevokeList {
		if c.msgListener != nil {
			t := new(db.LocalChatLog)
			t.ClientMsgID = w.Content
			t.Status = constant.MsgStatusRevoked
			err := c.db.UpdateMessage(t)
			if err != nil {
				log.Error("internal", "setLocalMessageStatus revokeMessage err:", err.Error(), "msg", w)
			} else {
				log.Info("internal", "v.OnRecvMessageRevoked client_msg_id:", w.Content)
				c.msgListener.OnRecvMessageRevoked(w.Content)
			}
		} else {
			log.Error("internal", "set msgListener is err:")
		}
	}

}
func (c *Conversation) newMessage(newMessagesList []*sdk_struct.MsgStruct) {
	for _, w := range newMessagesList {
		log.Info("internal", "newMessage: ", w.ClientMsgID)
		if c.msgListener != nil {
			log.Info("internal", "msgListener,OnRecvNewMessage")
			c.msgListener.OnRecvNewMessage(utils.StructToJsonString(w))
		} else {
			log.Error("internal", "set msgListener is err ")
		}
	}
}
func (c *Conversation) doDeleteConversation(c2v common.Cmd2Value) {
	if c.ConversationListener == nil {
		log.Error("internal", "not set conversationListener")
		return
	}
	node := c2v.Value.(common.DeleteConNode)
	//Mark messages related to this conversation for deletion
	err := c.db.UpdateMessageStatusBySourceID(node.SourceID, constant.MsgStatusHasDeleted, int32(node.SessionType))
	if err != nil {
		log.Error("internal", "setMessageStatusBySourceID err:", err.Error())
		return
	}
	//Reset the session information, empty session
	err = c.db.ResetConversation(node.ConversationID)
	if err != nil {
		log.Error("internal", "ResetConversation err:", err.Error())
	}
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{"", constant.TotalUnreadMessageChanged, ""}})
}
func (c *Conversation) doMsgReadState(msgReadList []*sdk_struct.MsgStruct) {
	var messageReceiptResp []*sdk_struct.MessageReceipt
	var msgIdList []string
	for _, rd := range msgReadList {
		err := json.Unmarshal([]byte(rd.Content), &msgIdList)
		if err != nil {
			log.Error("internal", "unmarshal failed, err : ", err.Error())
			return
		}
		var msgIdListStatusOK []string
		for _, v := range msgIdList {
			t := new(db.LocalChatLog)
			t.ClientMsgID = v
			t.IsRead = constant.HasRead
			err := c.db.UpdateMessage(t)
			if err != nil {
				log.Error("internal", "setMessageHasReadByMsgID err:", err, "ClientMsgID", v)
				continue
			}
			msgIdListStatusOK = append(msgIdListStatusOK, v)
		}
		if len(msgIdListStatusOK) > 0 {
			msgRt := new(sdk_struct.MessageReceipt)
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

		log.Info("internal", "OnRecvC2CReadReceipt: ", utils.StructToJsonString(messageReceiptResp))
		c.msgListener.OnRecvC2CReadReceipt(utils.StructToJsonString(messageReceiptResp))
	}
}

func (c *Conversation) doUpdateConversation(c2v common.Cmd2Value) {
	if c.ConversationListener == nil {
		log.Error("internal", "not set conversationListener")
		return
	}
	node := c2v.Value.(common.UpdateConNode)
	switch node.Action {
	case constant.AddConOrUpLatMsg:
		lc := node.Args.(db.LocalConversation)
		oc, err := c.db.GetConversation(node.ConId)
		if err == nil && oc != nil {
			if lc.LatestMsgSendTime > oc.LatestMsgSendTime { //The session update of asynchronous messages is subject to the latest sending time
				err := c.db.UpdateColumnsConversation(node.ConId, map[string]interface{}{"latest_msg_send_time": lc.LatestMsgSendTime, "latest_msg": lc.LatestMsg})
				if err != nil {
					log.Error("internal", "updateConversationLatestMsgModel err: ", err)
				}
			}
		} else {
			err4 := c.db.InsertConversation(&lc)
			if err4 != nil {
				log.Error("internal", "insert new conversation err:", err4.Error())

			}
			var list []*db.LocalConversation
			list = append(list, &lc)
			c.ConversationListener.OnNewConversation(utils.StructToJsonString(list))
		}

	case constant.UnreadCountSetZero:
		if err := c.db.UpdateColumnsConversation(node.ConId, map[string]interface{}{"unread_count": 0}); err != nil {
		} else {
			totalUnreadCount, err := c.db.GetTotalUnreadMsgCount()
			if err == nil {
				c.ConversationListener.OnTotalUnreadMessageCountChanged(totalUnreadCount)
			} else {
				log.Error("internal", "getTotalUnreadMsgCountModel err", err.Error())
			}

		}
	//case ConChange:
	//	err, list := u.getAllConversationListModel()
	//	if err != nil {
	//		sdkLog("getAllConversationListModel database err:", err.Error())
	//	} else {
	//		if list == nil {
	//			u.ConversationListenerx.OnConversationChanged(structToJsonString([]ConversationStruct{}))
	//		} else {
	//			u.ConversationListenerx.OnConversationChanged(structToJsonString(list))
	//
	//		}
	//	}
	case constant.IncrUnread:
		err := c.db.IncrConversationUnreadCount(node.ConId)
		if err != nil {
			log.Error("internal", "incrConversationUnreadCount database err:", err.Error())
			return
		}
	case constant.TotalUnreadMessageChanged:
		totalUnreadCount, err := c.db.GetTotalUnreadMsgCount()
		if err != nil {
			log.Error("internal", "TotalUnreadMessageChanged database err:", err.Error())
		} else {
			c.ConversationListener.OnTotalUnreadMessageCountChanged(totalUnreadCount)
		}
	case constant.UpdateFaceUrlAndNickName:
		lc := node.Args.(db.LocalConversation)
		if lc.ShowName != "" || lc.FaceURL != "" {

			err := c.db.UpdateConversation(&lc)
			if err != nil {
				log.Error("internal", "setConversationFaceUrlAndNickName database err:", err.Error())
				return
			}
		}

	case constant.UpdateLatestMessageChange:
		conversationID := node.ConId
		var latestMsg sdk_struct.MsgStruct
		l, err := c.db.GetConversation(conversationID)
		if err != nil {
			log.Error("internal", "getConversationLatestMsgModel err", err.Error())
		} else {
			err := json.Unmarshal([]byte(l.LatestMsg), &latestMsg)
			if err != nil {
				log.Error("internal", "latestMsg,Unmarshal err :", err.Error())
			} else {
				latestMsg.IsRead = true
				newLatestMessage := utils.StructToJsonString(latestMsg)
				err = c.db.UpdateColumnsConversation(node.ConId, map[string]interface{}{"latest_msg_send_time": latestMsg.SendTime, "latest_msg": newLatestMessage})
				if err != nil {
					log.Error("internal", "updateConversationLatestMsgModel err :", err.Error())
				}
			}
		}
	case constant.NewConChange:
		cidList := node.Args.([]string)
		cLists, err := c.db.GetMultipleConversation(cidList)
		if err != nil {
			log.Error("internal", "getMultipleConversationModel err :", err.Error())
		} else {
			if cLists != nil {
				log.Info("internal", "getMultipleConversationModel success :", cLists)
				c.ConversationListener.OnConversationChanged(utils.StructToJsonString(cLists))
			}
		}
	case constant.NewCon:
		cidList := node.Args.([]string)
		cLists, err := c.db.GetMultipleConversation(cidList)
		if err != nil {
			log.Error("internal", "getMultipleConversationModel err :", err.Error())
		} else {
			if cLists != nil {
				log.Info("internal", "getMultipleConversationModel success :", cLists)
				c.ConversationListener.OnNewConversation(utils.StructToJsonString(cLists))
			}
		}
	}
}

func (c *Conversation) Work(c2v common.Cmd2Value) {

	log.Info("internal", "doListener work..", c2v.Cmd)

	switch c2v.Cmd {
	case constant.CmdDeleteConversation:
		log.Info("internal", "CmdDeleteConversation start ..", c2v.Cmd)
		c.doDeleteConversation(c2v)
		log.Info("internal", "CmdDeleteConversation end..", c2v.Cmd)
	case constant.CmdNewMsgCome:
		log.Info("internal", "doMsgNew start..", c2v.Cmd)
		c.doMsgNew(c2v)
		log.Info("internal", "doMsgNew end..", c2v.Cmd)

	case constant.CmdUpdateConversation:
		log.Info("internal", "doUpdateConversation start ..", c2v.Cmd)
		c.doUpdateConversation(c2v)
		log.Info("internal", "doUpdateConversation end..", c2v.Cmd)
	}
}

func (c *Conversation) msgHandleByContentType(msg *sdk_struct.MsgStruct) (err error) {
	switch msg.ContentType {
	case constant.Text:
	case constant.Picture:
		err = utils.JsonStringToStruct(msg.Content, &msg.PictureElem)
	case constant.Voice:
		err = utils.JsonStringToStruct(msg.Content, &msg.SoundElem)
	case constant.Video:
		err = utils.JsonStringToStruct(msg.Content, &msg.VideoElem)
	case constant.File:
		err = utils.JsonStringToStruct(msg.Content, &msg.FileElem)
	case constant.AtText:
		err = utils.JsonStringToStruct(msg.Content, &msg.AtElem)
		if err == nil {
			if utils.IsContain(c.loginUserID, msg.AtElem.AtUserList) {
				msg.AtElem.IsAtSelf = true
			}
		}
	case constant.Location:
		err = utils.JsonStringToStruct(msg.Content, &msg.LocationElem)
	case constant.Custom:
		err = utils.JsonStringToStruct(msg.Content, &msg.CustomElem)
	case constant.Quote:
		err = utils.JsonStringToStruct(msg.Content, &msg.QuoteElem)
	case constant.Merger:
		err = utils.JsonStringToStruct(msg.Content, &msg.MergeElem)
	}
	return err
}
func (c *Conversation) updateConversation(lc *db.LocalConversation, cc, nc map[string]db.LocalConversation) {
	b, err := c.db.ConversationIfExists(lc.ConversationID)
	if err != nil {
		log.Error("internal", lc, cc, nc, err.Error())
		return
	}
	if b {
		if oldC, ok := cc[lc.ConversationID]; ok {
			if lc.LatestMsgSendTime > oldC.LatestMsgSendTime {
				lc.UnreadCount = lc.UnreadCount + oldC.UnreadCount
				cc[lc.ConversationID] = *lc
			} else {
				oldC.UnreadCount = oldC.UnreadCount + lc.UnreadCount
				cc[lc.ConversationID] = oldC
			}
		} else {
			cc[lc.ConversationID] = *lc
		}

	} else {
		if oldC, ok := nc[lc.ConversationID]; ok {
			if lc.LatestMsgSendTime > oldC.LatestMsgSendTime {
				lc.UnreadCount = lc.UnreadCount + oldC.UnreadCount
				nc[lc.ConversationID] = *lc
			} else {
				oldC.UnreadCount = oldC.UnreadCount + lc.UnreadCount
				cc[lc.ConversationID] = oldC
			}
		} else {
			nc[lc.ConversationID] = *lc
		}
	}

}
func mapConversationToList(m map[string]db.LocalConversation) (cs []*db.LocalConversation) {
	for _, v := range m {
		cs = append(cs, &v)
	}
	return cs
}
