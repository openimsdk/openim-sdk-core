package conversation_msg

import (
	"database/sql"
	"encoding/json"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/open_im_sdk/utils"
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
	ch                    chan open_im_sdk.cmd2Value
}
type InsertMsg struct {
	*open_im_sdk.MsgStruct
	isFilter bool
}

func (con *ConversationListener) getCh() chan open_im_sdk.cmd2Value {
	return con.ch
}

func (u *open_im_sdk.UserRelated) doMsgNew(c2v open_im_sdk.cmd2Value) {
	if u.MsgListenerList == nil {
		utils.sdkLog("not set c MsgListenerList", len(u.MsgListenerList))
		return
	}
	var insertMsg []*InsertMsg
	var errMsg, newMessages, msgReadList, msgRevokeList []*open_im_sdk.MsgStruct
	var isUnreadCount, isConversationUpdate bool
	var isCallbackUI bool
	conversationChangedSet := make(map[string]ConversationStruct)
	newConversationSet := make(map[string]ConversationStruct)
	//MsgList := c2v.Value.(ArrMsg)
	//for _, v := range MsgList.GroupData {
	//	MsgList.SingleData = append(MsgList.SingleData, v)
	//}
	utils.sdkLog("do Msg come here")
	u.seqMsgMutex.Lock()
	for _, v := range u.seqMsg {
		//isHistory = GetSwitchFromOptions(v.Options, IsHistory)
		isUnreadCount = utils.GetSwitchFromOptions(v.Options, open_im_sdk.IsUnreadCount)
		isConversationUpdate = utils.GetSwitchFromOptions(v.Options, open_im_sdk.IsConversationUpdate)
		isCallbackUI = true
		msg := &open_im_sdk.MsgStruct{
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
			Status:           open_im_sdk.MsgStatusSendSuccess,
			IsRead:           false,
		}
		utils.sdkLog("new msg, seq, ServerMsgID, ClientMsgID", msg.Seq, msg.ServerMsgID, msg.ClientMsgID)
		//De-analyze data
		err := u.msgHandleByContentType(msg)
		if err != nil {
			utils.sdkLog("Parsing data error:", err.Error(), msg)
			continue
		}
		switch v.SessionType {
		case open_im_sdk.SingleChatType:
			if v.ContentType > open_im_sdk.SingleTipBegin && v.ContentType < open_im_sdk.SingleTipEnd {
				u.doFriendMsg(v)
				utils.sdkLog("doFriendMsg, ", v)
			} else if v.ContentType > open_im_sdk.GroupTipBegin && v.ContentType < open_im_sdk.GroupTipEnd {
				u.doGroupMsg(v)
				utils.sdkLog("doGroupMsg, SingleChat ", v)
			}
		case open_im_sdk.GroupChatType:
			if v.ContentType > open_im_sdk.GroupTipBegin && v.ContentType < open_im_sdk.GroupTipEnd {
				u.doGroupMsg(v)
				utils.sdkLog("doGroupMsg, ", v)
			}
		}
		if v.SendID == u.loginUserID { //seq  Messages sent by myself  //if  sent through  this terminal
			m, err := u.getOneMessage(msg.ClientMsgID)
			if err == nil && m != nil {
				utils.sdkLog("have message", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, *msg)
				if m.Seq == 0 {
					insertMsg = append(insertMsg, &InsertMsg{MsgStruct: msg})
				} else {
					errMsg = append(errMsg, msg)

				}
			} else { //      send through  other terminal
				utils.sdkLog("sync message", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, *msg)
				c := ConversationStruct{
					ConversationType:  int(v.SessionType),
					LatestMsg:         utils.structToJsonString(msg),
					LatestMsgSendTime: msg.SendTime,
				}
				switch v.SessionType {
				case open_im_sdk.SingleChatType:
					c.ConversationID = utils.GetConversationIDBySessionType(v.RecvID, open_im_sdk.SingleChatType)
					c.UserID = v.RecvID
					faceUrl, name, _ := u.getUserNameAndFaceUrlByUid(c.UserID)
					c.FaceURL = faceUrl
					c.ShowName = name
				case open_im_sdk.GroupChatType:
					c.GroupID = v.GroupID
					c.ConversationID = utils.GetConversationIDBySessionType(c.GroupID, open_im_sdk.GroupChatType)
					faceUrl, name, err := u.getGroupNameAndFaceUrlByUid(c.GroupID)
					if err != nil {
						utils.sdkLog("getGroupNameAndFaceUrlByUid err:", err)
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
					LatestMsg:         utils.structToJsonString(msg),
					LatestMsgSendTime: msg.SendTime,
				}

				switch v.SessionType {
				case open_im_sdk.SingleChatType:
					c.ConversationID = utils.GetConversationIDBySessionType(v.SendID, open_im_sdk.SingleChatType)
					c.UserID = v.SendID
					c.ShowName = msg.SenderNickname
					c.FaceURL = msg.SenderFaceURL
				case open_im_sdk.GroupChatType:
					c.GroupID = v.GroupID
					c.ConversationID = utils.GetConversationIDBySessionType(c.GroupID, open_im_sdk.GroupChatType)
					faceUrl, name, err := u.getGroupNameAndFaceUrlByUid(c.GroupID)
					if err != nil {
						utils.sdkLog("getGroupNameAndFaceUrlByUid err:", err)
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
				if msg.ContentType == open_im_sdk.Revoke {
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
		utils.sdkLog("insert normal message err  :", err1.Error(), emsg1)
	}
	//Exception message storage
	err2, emsg2 := u.batchInsertErrorMessageToErrorChatLog(errMsg)
	if err2 != nil {
		utils.sdkLog("insert err message err  :", err2.Error(), emsg2)
	}
	//Changed conversation storage
	err3 := u.batchUpdateConversationLatestMsgModel(mapConversationToList(conversationChangedSet))
	if err3 != nil {
		utils.sdkLog("insert changed conversation err :", err3.Error())
	}
	//New conversation storage
	err4 := u.batchInsertConversationModel(mapConversationToList(newConversationSet))
	if err4 != nil {
		utils.sdkLog("insert new conversation err:", err4.Error())
	}
	//clear cache
	func(m map[int32]*open_im_sdk.MsgData) {
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
		utils.sdkLog("trigger map is :", newConversationSet, conversationChangedSet)
		//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewCon, mapKeyToStringList(newConversationSet)}})
		//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, mapKeyToStringList(conversationChangSet)}})
		u.ConversationListenerx.OnConversationChanged(utils.structToJsonString(mapConversationToList(conversationChangedSet)))
		u.ConversationListenerx.OnNewConversation(utils.structToJsonString(mapConversationToList(newConversationSet)))
		u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{"", open_im_sdk.TotalUnreadMessageChanged, ""}})
	}
	//sdkLog("length msgListenerList", u.MsgListenerList, "length message", len(newMessages), "msgListenerLen", len(u.MsgListenerList))

}

func (u *open_im_sdk.UserRelated) revokeMessage(msgRevokeList []*open_im_sdk.MsgStruct) {
	for _, v := range u.MsgListenerList {
		for _, w := range msgRevokeList {
			if v != nil {
				err := u.setMessageStatus(w.Content, open_im_sdk.MsgStatusRevoked)
				if err != nil {
					utils.sdkLog("setLocalMessageStatus revokeMessage err:", err.Error(), "msg", w)
				} else {
					utils.sdkLog("v.OnRecvMessageRevoked", w.Content)
					v.OnRecvMessageRevoked(w.Content)
				}
			} else {
				utils.sdkLog("set msgListener is err:")
			}
		}
	}
}
func (con *ConversationListener) newMessage(newMessagesList []*open_im_sdk.MsgStruct) {
	for _, v := range con.MsgListenerList {
		for _, w := range newMessagesList {
			utils.sdkLog("newMessage: ", w.ClientMsgID)
			if v != nil {
				utils.sdkLog("msgListener,OnRecvNewMessage")
				v.OnRecvNewMessage(utils.structToJsonString(w))
			} else {
				utils.sdkLog("set msgListener is err ")
			}
		}
	}
}
func (u *open_im_sdk.UserRelated) doDeleteConversation(c2v open_im_sdk.cmd2Value) {
	if u.ConversationListenerx == nil {
		utils.sdkLog("not set conversationListener")
		return
	}
	node := c2v.Value.(open_im_sdk.deleteConNode)
	//Mark messages related to this conversation for deletion
	err := u.setMessageStatusBySourceID(node.SourceID, open_im_sdk.MsgStatusHasDeleted, node.SessionType)
	if err != nil {
		utils.sdkLog("setMessageStatusBySourceID err:", err.Error())
		return
	}
	//Reset the session information, empty session
	err = u.ResetConversation(node.ConversationID)
	if err != nil {
		utils.sdkLog("ResetConversation err:", err.Error())
	}
	u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{"", open_im_sdk.TotalUnreadMessageChanged, ""}})
}
func (u *open_im_sdk.UserRelated) doMsgReadState(msgReadList []*open_im_sdk.MsgStruct) {
	var messageReceiptResp []*open_im_sdk.MessageReceipt
	var msgIdList []string
	for _, rd := range msgReadList {
		err := json.Unmarshal([]byte(rd.Content), &msgIdList)
		if err != nil {
			utils.sdkLog("unmarshal failed, err : ", err.Error())
			return
		}
		var msgIdListStatusOK []string
		for _, v := range msgIdList {
			err := u.setMessageHasReadByMsgID(v)
			if err != nil {
				utils.sdkLog("setMessageHasReadByMsgID err:", err, "msgID", v)
				continue
			}
			msgIdListStatusOK = append(msgIdListStatusOK, v)
		}
		if len(msgIdListStatusOK) > 0 {
			msgRt := new(open_im_sdk.MessageReceipt)
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
			utils.sdkLog("OnRecvC2CReadReceipt: ", utils.structToJsonString(messageReceiptResp))
			v.OnRecvC2CReadReceipt(utils.structToJsonString(messageReceiptResp))
		}
	}
}

func (u *open_im_sdk.UserRelated) doUpdateConversation(c2v open_im_sdk.cmd2Value) {
	if u.ConversationListenerx == nil {
		utils.sdkLog("not set conversationListener")
		return
	}
	node := c2v.Value.(open_im_sdk.updateConNode)
	switch node.Action {
	case open_im_sdk.AddConOrUpLatMsg:
		c := node.Args.(ConversationStruct)
		if u.judgeConversationIfExists(node.ConId) {
			_, o := u.getOneConversationModel(node.ConId)
			if c.LatestMsgSendTime > o.LatestMsgSendTime { //The session update of asynchronous messages is subject to the latest sending time
				err := u.updateConversationLatestMsgModel(c.LatestMsgSendTime, c.LatestMsg, node.ConId)
				if err != nil {
					utils.sdkLog("updateConversationLatestMsgModel err: ", err)
				}
			}
		} else {
			_ = u.insertConOrUpdateLatestMsg(&c, node.ConId)
			var list []*ConversationStruct
			list = append(list, &c)
			u.ConversationListenerx.OnNewConversation(utils.structToJsonString(list))
		}

	case open_im_sdk.UnreadCountSetZero:
		if err := u.setConversationUnreadCount(0, node.ConId); err != nil {
		} else {
			totalUnreadCount, err := u.getTotalUnreadMsgCountModel()
			if err == nil {
				u.ConversationListenerx.OnTotalUnreadMessageCountChanged(totalUnreadCount)
			} else {
				utils.sdkLog("getTotalUnreadMsgCountModel err", err.Error())
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
	case open_im_sdk.IncrUnread:
		err := u.incrConversationUnreadCount(node.ConId)
		if err != nil {
			utils.sdkLog("incrConversationUnreadCount database err:", err.Error())
			return
		}
	case open_im_sdk.TotalUnreadMessageChanged:
		totalUnreadCount, err := u.getTotalUnreadMsgCountModel()
		if err != nil {
			utils.sdkLog("TotalUnreadMessageChanged database err:", err.Error())
		} else {
			u.ConversationListenerx.OnTotalUnreadMessageCountChanged(totalUnreadCount)
		}
	case open_im_sdk.UpdateFaceUrlAndNickName:
		c := node.Args.(ConversationStruct)
		if c.ShowName != "" || c.FaceURL != "" {
			err := u.setConversationFaceUrlAndNickName(&c, node.ConId)
			if err != nil {
				utils.sdkLog("setConversationFaceUrlAndNickName database err:", err.Error())
				return
			}
		}

	case open_im_sdk.UpdateLatestMessageChange:
		conversationID := node.ConId
		var latestMsg open_im_sdk.MsgStruct
		err, l := u.getConversationLatestMsgModel(conversationID)
		if err != nil {
			utils.sdkLog("getConversationLatestMsgModel err", err.Error())
		} else {
			err := json.Unmarshal([]byte(l), &latestMsg)
			if err != nil {
				utils.sdkLog("latestMsg,Unmarshal err :", err.Error())
			} else {
				latestMsg.IsRead = true
				newLatestMessage := utils.structToJsonString(latestMsg)
				err = u.updateConversationLatestMsgModel(latestMsg.SendTime, newLatestMessage, conversationID)
				if err != nil {
					utils.sdkLog("updateConversationLatestMsgModel err :", err.Error())
				}
			}
		}
	case open_im_sdk.NewConChange:
		cidList := node.Args.([]string)
		err, cList := u.getMultipleConversationModel(cidList)
		if err != nil {
			utils.sdkLog("getMultipleConversationModel err :", err.Error())
		} else {
			if cList != nil {
				utils.sdkLog("getMultipleConversationModel success :", cList)
				u.ConversationListenerx.OnConversationChanged(utils.structToJsonString(cList))
			}
		}
	case open_im_sdk.NewCon:
		cidList := node.Args.([]string)
		err, cList := u.getMultipleConversationModel(cidList)
		if err != nil {
			utils.sdkLog("getMultipleConversationModel err :", err.Error())
		} else {
			if cList != nil {
				utils.sdkLog("getMultipleConversationModel success :", cList)
				u.ConversationListenerx.OnNewConversation(utils.structToJsonString(cList))
			}
		}
	}
}

func (u *open_im_sdk.UserRelated) work(c2v open_im_sdk.cmd2Value) {

	utils.sdkLog("doListener work..", c2v.Cmd)

	switch c2v.Cmd {
	case open_im_sdk.CmdDeleteConversation:
		utils.sdkLog("CmdDeleteConversation start ..", c2v.Cmd)
		u.doDeleteConversation(c2v)
		utils.sdkLog("CmdDeleteConversation end..", c2v.Cmd)
	case open_im_sdk.CmdNewMsgCome:
		utils.sdkLog("doMsgNew start..", c2v.Cmd)

		u.doMsgNew(c2v)
		utils.sdkLog("doMsgNew end..", c2v.Cmd)
	case open_im_sdk.CmdUpdateConversation:
		utils.sdkLog("doUpdateConversation start ..", c2v.Cmd)
		u.doUpdateConversation(c2v)
		utils.sdkLog("doUpdateConversation end..", c2v.Cmd)
	}
}

func (u *open_im_sdk.UserRelated) msgHandleByContentType(msg *open_im_sdk.MsgStruct) (err error) {
	switch msg.ContentType {
	case open_im_sdk.Text:
	case open_im_sdk.Picture:
		err = utils.jsonStringToStruct(msg.Content, &msg.PictureElem)
	case open_im_sdk.Voice:
		err = utils.jsonStringToStruct(msg.Content, &msg.SoundElem)
	case open_im_sdk.Video:
		err = utils.jsonStringToStruct(msg.Content, &msg.VideoElem)
	case open_im_sdk.File:
		err = utils.jsonStringToStruct(msg.Content, &msg.FileElem)
	case open_im_sdk.AtText:
		err = utils.jsonStringToStruct(msg.Content, &msg.AtElem)
		if err == nil {
			if utils.isContain(u.loginUserID, msg.AtElem.AtUserList) {
				msg.AtElem.IsAtSelf = true
			}
		}
	case open_im_sdk.Location:
		err = utils.jsonStringToStruct(msg.Content, &msg.LocationElem)
	case open_im_sdk.Custom:
		err = utils.jsonStringToStruct(msg.Content, &msg.CustomElem)
	case open_im_sdk.Quote:
		err = utils.jsonStringToStruct(msg.Content, &msg.QuoteElem)
	case open_im_sdk.Merger:
		err = utils.jsonStringToStruct(msg.Content, &msg.MergeElem)
	}
	return err
}
func (u *open_im_sdk.UserRelated) getGroupNameAndFaceUrlByUid(groupID string) (faceUrl, name string, err error) {
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
func (u *open_im_sdk.UserRelated) updateConversation(c *ConversationStruct, cc, nc map[string]ConversationStruct) {
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
