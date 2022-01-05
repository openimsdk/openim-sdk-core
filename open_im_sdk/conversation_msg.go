package open_im_sdk

import (
	"database/sql"
	"encoding/json"
)

type ChatLog struct {
	MsgId            string
	SendID           string
	IsRead           int32
	IsFilter         int32
	Seq              int64
	Status           int32
	SessionType      int32
	RecvID           string
	ContentType      int32
	MsgFrom          int32
	Content          string
	Remark           sql.NullString
	SenderPlatformID int32
	SendTime         int64
	CreateTime       int64
}
type ConversationStruct struct {
	ConversationID    string `json:"conversationID"`
	ConversationType  int    `json:"conversationType"`
	UserID            string `json:"userID"`
	GroupID           string `json:"groupID"`
	ShowName          string `json:"showName"`
	FaceURL           string `json:"faceUrl"`
	RecvMsgOpt        int    `json:"recvMsgOpt"`
	UnreadCount       int    `json:"unreadCount"`
	GroupAtType       int    `json:"groupAtType"`
	LatestMsg         string `json:"latestMsg"`
	LatestMsgSendTime int64  `json:"latestMsgSendTime"`
	DraftText         string `json:"draftText"`
	DraftTimestamp    int64  `json:"draftTimestamp"`
	IsPinned          int    `json:"isPinned"`
}
type ConversationListener struct {
	ConversationListenerx OnConversationListener
	MsgListenerList       []OnAdvancedMsgListener
	ch                    chan cmd2Value
}
type InsertMsg struct {
	*MsgStruct
	isFilter bool
}

func (con *ConversationListener) getCh() chan cmd2Value {
	return con.ch
}

func (u *UserRelated) doMsgNew(c2v cmd2Value) {
	if u.MsgListenerList == nil {
		sdkLog("not set c MsgListenerList", len(u.MsgListenerList))
		return
	}
	var insertMsg []*InsertMsg
	var errMsg, newMessages, msgReadList, msgRevokeList []*MsgStruct
	var isUnreadCount, isConversationUpdate bool
	var isCallbackUI bool
	conversationChangedSet := make(map[string]ConversationStruct)
	newConversationSet := make(map[string]ConversationStruct)
	//MsgList := c2v.Value.(ArrMsg)
	//for _, v := range MsgList.GroupData {
	//	MsgList.SingleData = append(MsgList.SingleData, v)
	//}
	sdkLog("do Msg come here")
	u.seqMsgMutex.Lock()
	for _, v := range u.seqMsg {
		//isHistory = GetSwitchFromOptions(v.Options, IsHistory)
		isUnreadCount = GetSwitchFromOptions(v.Options, IsUnreadCount)
		isConversationUpdate = GetSwitchFromOptions(v.Options, IsConversationUpdate)
		isCallbackUI = true
		msg := &MsgStruct{
			SendID:           v.SendID,
			SessionType:      v.SessionType,
			MsgFrom:          v.MsgFrom,
			ContentType:      v.ContentType,
			ServerMsgID:      v.ServerMsgID,
			ClientMsgID:      v.ClientMsgID,
			Content:          string(v.Content),
			SendTime:         v.SendTime,
			CreateTime:       v.CreateTime,
			RecvID:           v.RecvID,
			SenderFaceURL:    v.SenderFaceURL,
			SenderNickname:   v.SenderNickname,
			Seq:              v.Seq,
			SenderPlatformID: v.SenderPlatformID,
			ForceList:        v.ForceList,
			GroupID:          v.GroupID,
			Status:           MsgStatusSendSuccess,
			IsRead:           false,
		}
		sdkLog("new msg, seq, ServerMsgID, ClientMsgID", msg.Seq, msg.ServerMsgID, msg.ClientMsgID)
		//De-analyze data
		err := u.msgHandleByContentType(msg)
		if err != nil {
			sdkLog("Parsing data error:", err.Error(), msg)
			continue
		}
		switch v.SessionType {
		case SingleChatType:
			if v.ContentType > SingleTipBegin && v.ContentType < SingleTipEnd {
				u.doFriendMsg(v)
				sdkLog("doFriendMsg, ", v)
			} else if v.ContentType > GroupTipBegin && v.ContentType < GroupTipEnd {
				u.doGroupMsg(v)
				sdkLog("doGroupMsg, SingleChat ", v)
			}
		case GroupChatType:
			if v.ContentType > GroupTipBegin && v.ContentType < GroupTipEnd {
				u.doGroupMsg(v)
				sdkLog("doGroupMsg, ", v)
			}
		}
		if v.SendID == u.loginUserID { //seq  Messages sent by myself  //if  sent through  this terminal
			m, err := u.getOneMessage(msg.ClientMsgID)
			if err == nil && m != nil {
				sdkLog("have message", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, *msg)
				if m.Seq == 0 {
					insertMsg = append(insertMsg, &InsertMsg{MsgStruct: msg})
				} else {
					errMsg = append(errMsg, msg)

				}
			} else { //      send through  other terminal
				sdkLog("sync message", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, *msg)
				c := ConversationStruct{
					ConversationType:  int(v.SessionType),
					LatestMsg:         structToJsonString(msg),
					LatestMsgSendTime: msg.SendTime,
				}
				switch v.SessionType {
				case SingleChatType:
					c.ConversationID = GetConversationIDBySessionType(v.RecvID, SingleChatType)
					c.UserID = v.RecvID
					faceUrl, name, _ := u.getUserNameAndFaceUrlByUid(c.UserID)
					c.FaceURL = faceUrl
					c.ShowName = name
				case GroupChatType:
					c.GroupID = v.GroupID
					c.ConversationID = GetConversationIDBySessionType(c.GroupID, GroupChatType)
					faceUrl, name, err := u.getGroupNameAndFaceUrlByUid(c.GroupID)
					if err != nil {
						sdkLog("getGroupNameAndFaceUrlByUid err:", err)
					} else {
						c.ShowName = name
						c.FaceURL = faceUrl
					}
				}
				if isUnreadCount {
					c.UnreadCount = 1
				}
				if isConversationUpdate {
					u.updateConversation(&c, conversationChangedSet, newConversationSet)
					insertMsg = append(insertMsg, &InsertMsg{MsgStruct: msg})
				} else {
					insertMsg = append(insertMsg, &InsertMsg{MsgStruct: msg, isFilter: true})
				}
				newMessages = append(newMessages, msg)

			}
		} else { //Sent by others
			if !u.judgeMessageIfExists(msg.ClientMsgID) { //Deduplication operation
				c := ConversationStruct{
					ConversationType:  int(v.SessionType),
					LatestMsg:         structToJsonString(msg),
					LatestMsgSendTime: msg.SendTime,
				}

				switch v.SessionType {
				case SingleChatType:
					c.ConversationID = GetConversationIDBySessionType(v.SendID, SingleChatType)
					c.UserID = v.SendID
					c.ShowName = msg.SenderNickname
					c.FaceURL = msg.SenderFaceURL
				case GroupChatType:
					c.GroupID = v.GroupID
					c.ConversationID = GetConversationIDBySessionType(c.GroupID, GroupChatType)
					faceUrl, name, err := u.getGroupNameAndFaceUrlByUid(c.GroupID)
					if err != nil {
						sdkLog("getGroupNameAndFaceUrlByUid err:", err)
					} else {
						c.ShowName = name
						c.FaceURL = faceUrl
					}
				}
				if isUnreadCount {
					c.UnreadCount = 1
				}
				//u.doUpdateConversation(cmd2Value{Value: updateConNode{c.ConversationID, UpdateFaceUrlAndNickName, c}})
				if isConversationUpdate {
					insertMsg = append(insertMsg, &InsertMsg{MsgStruct: msg})
					u.updateConversation(&c, conversationChangedSet, newConversationSet)
					newMessages = append(newMessages, msg)

				} else {
					insertMsg = append(insertMsg, &InsertMsg{MsgStruct: msg, isFilter: true})

				}
				if msg.ContentType == Revoke {
					msgRevokeList = append(msgRevokeList, msg)
				}
			} else {
				errMsg = append(errMsg, msg)
			}
		}
	}
	//Normal message storage
	err1, emsg1 := u.batchInsertMessageToChatLog(insertMsg)
	if err1 != nil {
		sdkLog("insert normal message err  :", err1.Error(), emsg1)
	}
	//Exception message storage
	err2, emsg2 := u.batchInsertErrorMessageToErrorChatLog(errMsg)
	if err2 != nil {
		sdkLog("insert err message err  :", err2.Error(), emsg2)
	}
	//Changed conversation storage
	err3 := u.batchUpdateConversationLatestMsgModel(mapConversationToList(conversationChangedSet))
	if err3 != nil {
		sdkLog("insert changed conversation err :", err3.Error())
	}
	//New conversation storage
	err4 := u.batchInsertConversationModel(mapConversationToList(newConversationSet))
	if err4 != nil {
		sdkLog("insert new conversation err:", err4.Error())
	}
	//clear cache
	func(m map[int32]*MsgData) {
		for k := range m {
			delete(m, k)
		}
	}(u.seqMsg)
	u.seqMsgMutex.Unlock()
	if isCallbackUI {
		u.doMsgReadState(msgReadList)
		u.revokeMessage(msgRevokeList)
		u.newMessage(newMessages)
		//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", ConChange, ""}})
		sdkLog("trigger map is :", newConversationSet, conversationChangedSet)
		//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewCon, mapKeyToStringList(newConversationSet)}})
		//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, mapKeyToStringList(conversationChangSet)}})
		u.ConversationListenerx.OnConversationChanged(structToJsonString(mapConversationToList(conversationChangedSet)))
		u.ConversationListenerx.OnNewConversation(structToJsonString(mapConversationToList(newConversationSet)))
		u.doUpdateConversation(cmd2Value{Value: updateConNode{"", TotalUnreadMessageChanged, ""}})
	}
	//sdkLog("length msgListenerList", u.MsgListenerList, "length message", len(newMessages), "msgListenerLen", len(u.MsgListenerList))

}

func (u *UserRelated) revokeMessage(msgRevokeList []*MsgStruct) {
	for _, v := range u.MsgListenerList {
		for _, w := range msgRevokeList {
			if v != nil {
				err := u.setMessageStatus(w.Content, MsgStatusRevoked)
				if err != nil {
					sdkLog("setLocalMessageStatus revokeMessage err:", err.Error(), "msg", w)
				} else {
					sdkLog("v.OnRecvMessageRevoked", w.Content)
					v.OnRecvMessageRevoked(w.Content)
				}
			} else {
				sdkLog("set msgListener is err:")
			}
		}
	}
}
func (con *ConversationListener) newMessage(newMessagesList []*MsgStruct) {
	for _, v := range con.MsgListenerList {
		for _, w := range newMessagesList {
			sdkLog("newMessage: ", w.ClientMsgID)
			if v != nil {
				sdkLog("msgListener,OnRecvNewMessage")
				v.OnRecvNewMessage(structToJsonString(w))
			} else {
				sdkLog("set msgListener is err ")
			}
		}
	}
}
func (u *UserRelated) doDeleteConversation(c2v cmd2Value) {
	if u.ConversationListenerx == nil {
		sdkLog("not set conversationListener")
		return
	}
	node := c2v.Value.(deleteConNode)
	//Mark messages related to this conversation for deletion
	err := u.setMessageStatusBySourceID(node.SourceID, MsgStatusHasDeleted, node.SessionType)
	if err != nil {
		sdkLog("setMessageStatusBySourceID err:", err.Error())
		return
	}
	//Reset the session information, empty session
	err = u.ResetConversation(node.ConversationID)
	if err != nil {
		sdkLog("ResetConversation err:", err.Error())
	}
	u.doUpdateConversation(cmd2Value{Value: updateConNode{"", TotalUnreadMessageChanged, ""}})
}
func (u *UserRelated) doMsgReadState(msgReadList []*MsgStruct) {
	var messageReceiptResp []*MessageReceipt
	var msgIdList []string
	for _, rd := range msgReadList {
		err := json.Unmarshal([]byte(rd.Content), &msgIdList)
		if err != nil {
			sdkLog("unmarshal failed, err : ", err.Error())
			return
		}
		var msgIdListStatusOK []string
		for _, v := range msgIdList {
			err := u.setMessageHasReadByMsgID(v)
			if err != nil {
				sdkLog("setMessageHasReadByMsgID err:", err, "msgID", v)
				continue
			}
			msgIdListStatusOK = append(msgIdListStatusOK, v)
		}
		if len(msgIdListStatusOK) > 0 {
			msgRt := new(MessageReceipt)
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
		for _, v := range u.MsgListenerList {
			sdkLog("OnRecvC2CReadReceipt: ", structToJsonString(messageReceiptResp))
			v.OnRecvC2CReadReceipt(structToJsonString(messageReceiptResp))
		}
	}
}

func (u *UserRelated) doUpdateConversation(c2v cmd2Value) {
	if u.ConversationListenerx == nil {
		sdkLog("not set conversationListener")
		return
	}
	node := c2v.Value.(updateConNode)
	switch node.Action {
	case AddConOrUpLatMsg:
		c := node.Args.(ConversationStruct)
		if u.judgeConversationIfExists(node.ConId) {
			_, o := u.getOneConversationModel(node.ConId)
			if c.LatestMsgSendTime > o.LatestMsgSendTime { //The session update of asynchronous messages is subject to the latest sending time
				err := u.updateConversationLatestMsgModel(c.LatestMsgSendTime, c.LatestMsg, node.ConId)
				if err != nil {
					sdkLog("updateConversationLatestMsgModel err: ", err)
				}
			}
		} else {
			_ = u.insertConOrUpdateLatestMsg(&c, node.ConId)
			var list []*ConversationStruct
			list = append(list, &c)
			u.ConversationListenerx.OnNewConversation(structToJsonString(list))
		}

	case UnreadCountSetZero:
		if err := u.setConversationUnreadCount(0, node.ConId); err != nil {
		} else {
			totalUnreadCount, err := u.getTotalUnreadMsgCountModel()
			if err == nil {
				u.ConversationListenerx.OnTotalUnreadMessageCountChanged(totalUnreadCount)
			} else {
				sdkLog("getTotalUnreadMsgCountModel err", err.Error())
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
	case IncrUnread:
		err := u.incrConversationUnreadCount(node.ConId)
		if err != nil {
			sdkLog("incrConversationUnreadCount database err:", err.Error())
			return
		}
	case TotalUnreadMessageChanged:
		totalUnreadCount, err := u.getTotalUnreadMsgCountModel()
		if err != nil {
			sdkLog("TotalUnreadMessageChanged database err:", err.Error())
		} else {
			u.ConversationListenerx.OnTotalUnreadMessageCountChanged(totalUnreadCount)
		}
	case UpdateFaceUrlAndNickName:
		c := node.Args.(ConversationStruct)
		if c.ShowName != "" || c.FaceURL != "" {
			err := u.setConversationFaceUrlAndNickName(&c, node.ConId)
			if err != nil {
				sdkLog("setConversationFaceUrlAndNickName database err:", err.Error())
				return
			}
		}

	case UpdateLatestMessageChange:
		conversationID := node.ConId
		var latestMsg MsgStruct
		err, l := u.getConversationLatestMsgModel(conversationID)
		if err != nil {
			sdkLog("getConversationLatestMsgModel err", err.Error())
		} else {
			err := json.Unmarshal([]byte(l), &latestMsg)
			if err != nil {
				sdkLog("latestMsg,Unmarshal err :", err.Error())
			} else {
				latestMsg.IsRead = true
				newLatestMessage := structToJsonString(latestMsg)
				err = u.updateConversationLatestMsgModel(latestMsg.SendTime, newLatestMessage, conversationID)
				if err != nil {
					sdkLog("updateConversationLatestMsgModel err :", err.Error())
				}
			}
		}
	case NewConChange:
		cidList := node.Args.([]string)
		err, cList := u.getMultipleConversationModel(cidList)
		if err != nil {
			sdkLog("getMultipleConversationModel err :", err.Error())
		} else {
			if cList != nil {
				sdkLog("getMultipleConversationModel success :", cList)
				u.ConversationListenerx.OnConversationChanged(structToJsonString(cList))
			}
		}
	case NewCon:
		cidList := node.Args.([]string)
		err, cList := u.getMultipleConversationModel(cidList)
		if err != nil {
			sdkLog("getMultipleConversationModel err :", err.Error())
		} else {
			if cList != nil {
				sdkLog("getMultipleConversationModel success :", cList)
				u.ConversationListenerx.OnNewConversation(structToJsonString(cList))
			}
		}
	}
}

func (u *UserRelated) work(c2v cmd2Value) {

	sdkLog("doListener work..", c2v.Cmd)

	switch c2v.Cmd {
	case CmdDeleteConversation:
		sdkLog("CmdDeleteConversation start ..", c2v.Cmd)
		u.doDeleteConversation(c2v)
		sdkLog("CmdDeleteConversation end..", c2v.Cmd)
	case CmdNewMsgCome:
		sdkLog("doMsgNew start..", c2v.Cmd)

		u.doMsgNew(c2v)
		sdkLog("doMsgNew end..", c2v.Cmd)
	case CmdUpdateConversation:
		sdkLog("doUpdateConversation start ..", c2v.Cmd)
		u.doUpdateConversation(c2v)
		sdkLog("doUpdateConversation end..", c2v.Cmd)
	}
}

func (u *UserRelated) msgHandleByContentType(msg *MsgStruct) (err error) {
	switch msg.ContentType {
	case Text:
	case Picture:
		err = jsonStringToStruct(msg.Content, &msg.PictureElem)
	case Voice:
		err = jsonStringToStruct(msg.Content, &msg.SoundElem)
	case Video:
		err = jsonStringToStruct(msg.Content, &msg.VideoElem)
	case File:
		err = jsonStringToStruct(msg.Content, &msg.FileElem)
	case AtText:
		err = jsonStringToStruct(msg.Content, &msg.AtElem)
		if err == nil {
			if isContain(u.loginUserID, msg.AtElem.AtUserList) {
				msg.AtElem.IsAtSelf = true
			}
		}
	case Location:
		err = jsonStringToStruct(msg.Content, &msg.LocationElem)
	case Custom:
		err = jsonStringToStruct(msg.Content, &msg.CustomElem)
	case Quote:
		err = jsonStringToStruct(msg.Content, &msg.QuoteElem)
	case Merger:
		err = jsonStringToStruct(msg.Content, &msg.MergeElem)
	}
	return err
}
func (u *UserRelated) getGroupNameAndFaceUrlByUid(groupID string) (faceUrl, name string, err error) {
	groupInfo, err := u.getLocalGroupsInfoByGroupID(groupID)
	if err != nil {
		return "", "", err
	}
	if groupInfo.GroupId == "" {
		groupInfo, err := u.getGroupInfoByGroupId(groupID)
		if err != nil {
			return "", "", err
		} else {
			return groupInfo.FaceUrl, groupInfo.GroupName, nil
		}
	} else {
		return groupInfo.FaceUrl, groupInfo.GroupName, nil
	}
}
func (u *UserRelated) updateConversation(c *ConversationStruct, cc, nc map[string]ConversationStruct) {
	if u.judgeConversationIfExists(c.ConversationID) {
		//_, o := u.getOneConversationModel(c.ConversationID)
		//if c.LatestMsgSendTime > o.LatestMsgSendTime { //The session update of asynchronous messages is subject to the latest sending time
		//	err := u.updateConversationLatestMsgModel(c.LatestMsgSendTime, c.LatestMsg, c.ConversationID)
		//	if err != nil {
		//		sdkLog("updateConversationLatestMsgModel err: ", err)
		//	} else {
		//		cc[c.ConversationID] = void{}
		//	}
		//}
		if oldC, ok := cc[c.ConversationID]; ok {
			if c.LatestMsgSendTime > oldC.LatestMsgSendTime {
				c.UnreadCount = c.UnreadCount + oldC.UnreadCount
				cc[c.ConversationID] = *c
			} else {
				oldC.UnreadCount = oldC.UnreadCount + c.UnreadCount
				cc[c.ConversationID] = oldC
			}
		} else {
			cc[c.ConversationID] = *c
		}

	} else {
		if oldC, ok := nc[c.ConversationID]; ok {
			if c.LatestMsgSendTime > oldC.LatestMsgSendTime {
				c.UnreadCount = c.UnreadCount + oldC.UnreadCount
				nc[c.ConversationID] = *c
			} else {
				oldC.UnreadCount = oldC.UnreadCount + c.UnreadCount
				cc[c.ConversationID] = oldC
			}
		} else {
			nc[c.ConversationID] = *c
		}
	}

	//if u.judgeConversationIfExists(c.ConversationID) {
	//	_, o := u.getOneConversationModel(c.ConversationID)
	//	if c.LatestMsgSendTime > o.LatestMsgSendTime { //The session update of asynchronous messages is subject to the latest sending time
	//		err := u.updateConversationLatestMsgModel(c.LatestMsgSendTime, c.LatestMsg, c.ConversationID)
	//		if err != nil {
	//			sdkLog("updateConversationLatestMsgModel err: ", err)
	//		} else {
	//			cc[c.ConversationID] = void{}
	//		}
	//	}
	//} else {
	//	err := u.insertConOrUpdateLatestMsg(c, c.ConversationID)
	//	if err != nil {
	//		sdkLog("insertConOrUpdateLatestMsg err: ", err.Error())
	//	} else {
	//		nc[c.ConversationID] = void{}
	//	}
	//	//var list []*ConversationStruct
	//	//list = append(list, c)
	//	//u.ConversationListenerx.OnNewConversation(structToJsonString(list))
	//}
}
func mapConversationToList(m map[string]ConversationStruct) (cs []*ConversationStruct) {
	for _, v := range m {
		cs = append(cs, &v)
	}
	return cs
}
