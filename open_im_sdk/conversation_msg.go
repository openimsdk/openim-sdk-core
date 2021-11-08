package open_im_sdk

import (
	"database/sql"
	"encoding/json"
	"strings"
)

type ChatLog struct {
	MsgId            string
	SendID           string
	IsRead           int32
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

func (con *ConversationListener) getCh() chan cmd2Value {
	return con.ch
}

func (u *UserRelated) doMsgNew(c2v cmd2Value) {
	if u.MsgListenerList == nil {
		sdkLog("not set c MsgListenerList", len(u.MsgListenerList))
		return
	}
	var newMessages, msgReadList, msgRevokeList []*MsgStruct
	MsgList := c2v.Value.(ArrMsg)
	for _, v := range MsgList.GroupData {
		MsgList.SingleData = append(MsgList.SingleData, v)
	}
	sdkLog("do Msg come here,len:", len(MsgList.SingleData))
	for _, v := range MsgList.SingleData {
		msg := &MsgStruct{
			SendID:         v.SendID,
			SessionType:    v.SessionType,
			MsgFrom:        v.MsgFrom,
			ContentType:    v.ContentType,
			ServerMsgID:    v.ServerMsgID,
			ClientMsgID:    v.ClientMsgID,
			Content:        v.Content,
			SendTime:       v.SendTime,
			SenderFaceURL:  v.SenderFaceURL,
			SenderNickName: v.SenderNickName,
			Seq:            v.Seq,
			PlatformID:     v.SenderPlatformID,
			Status:         MsgStatusSendSuccess,
			IsRead:         false,
		}
		//De-analyze data
		err := u.msgHandleByContentType(msg)
		if err != nil {
			sdkLog("Parsing data error:", err.Error(), msg)
			return
		}
		switch v.SessionType {
		case SingleChatType:
			msg.RecvID = v.RecvID
		case GroupChatType:
			msg.RecvID = strings.Split(v.RecvID, " ")[1]
			msg.GroupID = msg.RecvID
		}
		if v.SendID == u.LoginUid { //seq對齊消息 Messages sent by myself  //if  sent through  this terminal
			m, err := u.getOneMessage(msg.ClientMsgID)
			if err == nil && m != nil {
				sdkLog("have message", *msg, msg.Seq)
				if m.Seq == 0 {
					err := u.updateMessageSeq(msg, MsgStatusSendSuccess)
					if err != nil {
						sdkLog("updateMessageSeq err", err.Error(), msg)
					}
				} else {
					err := u.setErrorMessageToErrorChatLog(msg)
					if err != nil {
						sdkLog("setErrorMessage  err", err.Error(), msg)
					}
				}
			} else { //同步消息       send through  other terminal
				err = u.insertPushMessageToChatLog(msg)
				if err != nil {
					sdkLog(" sync insertPushMessageToChatLog err", err.Error(), msg)
					err := u.setErrorMessageToErrorChatLog(msg)
					if err != nil {
						sdkLog("setErrorMessage err", err.Error(), *msg)
					}
				}
				c := ConversationStruct{
					//ConversationID:    conversationID,
					ConversationType: int(v.SessionType),
					//UserID:            userID,
					//GroupID:           groupID,
					RecvMsgOpt:        1,
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
					u.doUpdateConversation(cmd2Value{Value: updateConNode{c.ConversationID, UpdateFaceUrlAndNickName, c}})
				case GroupChatType:
					c.GroupID = strings.Split(v.RecvID, " ")[1]
					c.ConversationID = GetConversationIDBySessionType(c.GroupID, GroupChatType)
				}

				if msg.ContentType <= AcceptFriendApplicationTip && msg.ContentType != HasReadReceipt {
					newMessages = append(newMessages, msg)
					u.doUpdateConversation(cmd2Value{Value: updateConNode{c.ConversationID, AddConOrUpLatMsg,
						c}})
				}
				//}
			}
		} else { //他人發的
			if !u.judgeMessageIfExists(msg) { //去重操作
				if msg.ContentType != Typing && msg.ContentType != HasReadReceipt {
					c := ConversationStruct{
						//ConversationID:    conversationID,
						ConversationType: int(v.SessionType),
						//ShowName:          msg.SenderNickName,
						//FaceURL:           msg.SenderFaceURL,
						//UserID:            userID,
						//GroupID:           groupID,
						RecvMsgOpt:        1,
						LatestMsg:         structToJsonString(msg),
						LatestMsgSendTime: msg.SendTime,
					}
					err = u.insertPushMessageToChatLog(msg)
					if err != nil {
						sdkLog("insertPushMessageToChatLog err", err.Error(), msg)
					}
					switch v.SessionType {
					case SingleChatType:
						c.ConversationID = GetConversationIDBySessionType(v.SendID, SingleChatType)
						c.UserID = v.SendID
						c.ShowName = msg.SenderNickName
						c.FaceURL = msg.SenderFaceURL
					case GroupChatType:
						c.GroupID = strings.Split(v.RecvID, " ")[1]
						c.ConversationID = GetConversationIDBySessionType(c.GroupID, GroupChatType)
						faceUrl, name, err := u.getGroupNameAndFaceUrlByUid(c.GroupID)
						if err != nil {
							sdkLog("getGroupNameAndFaceUrlByUid err:", err)
						} else {
							c.ShowName = name
							c.FaceURL = faceUrl
						}
					}
					if msg.ContentType <= AcceptFriendApplicationTip || (msg.ContentType >= GroupTipBegin && msg.ContentType <= GroupTipEnd && msg.ContentType != SetGroupInfoTip && msg.ContentType != JoinGroupTip) {
						u.doUpdateConversation(cmd2Value{Value: updateConNode{c.ConversationID, AddConOrUpLatMsg,
							c}})
						if msg.ContentType != Revoke {
							u.doUpdateConversation(cmd2Value{Value: updateConNode{c.ConversationID, IncrUnread, ""}})
						}
						u.doUpdateConversation(cmd2Value{Value: updateConNode{c.ConversationID, UpdateFaceUrlAndNickName, c}})
						newMessages = append(newMessages, msg)
					}
					if msg.ContentType == SetGroupInfoTip || msg.ContentType == SetSelfInfoTip {
						u.doUpdateConversation(cmd2Value{Value: updateConNode{c.ConversationID, UpdateFaceUrlAndNickName, c}})
						newMessages = append(newMessages, msg)

					}
					if msg.ContentType == Revoke {
						msgRevokeList = append(msgRevokeList, msg)
					}
				} else {
					if msg.ContentType == Typing {
						newMessages = append(newMessages, msg)

					} else {
						err = u.insertPushMessageToChatLog(msg)
						if err != nil {
							sdkLog("insert HasReadReceipt err:", err)
						}
						//update read state
						msgReadList = append(msgReadList, msg)
					}
				}
			} else {
				err := u.setErrorMessageToErrorChatLog(msg)
				if err != nil {
					sdkLog("setErrorMessage  err", err.Error(), msg)
				} else {
					sdkLog("repeat message err", msg.ClientMsgID, msg.Seq, msg.ServerMsgID)
				}
			}
		}
	}
	u.doMsgReadState(msgReadList)
	u.revokeMessage(msgRevokeList)
	u.newMessage(newMessages)
	u.doUpdateConversation(cmd2Value{Value: updateConNode{"", ConChange, ""}})
	u.doUpdateConversation(cmd2Value{Value: updateConNode{"", TotalUnreadMessageChanged, ""}})
	sdkLog("length msgListenerList", u.MsgListenerList, "length message", len(newMessages), "msgListenerLen", len(u.MsgListenerList))

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
	//标记删除与此会话相关的消息
	err := u.setMessageStatusBySourceID(node.SourceID, MsgStatusHasDeleted, node.SessionType)
	if err != nil {
		sdkLog("setMessageStatusBySourceID err:", err.Error())
		return
	}
	//重置该会话信息，空会话
	err = u.ResetConversation(node.ConversationID)
	if err != nil {
		sdkLog("ResetConversation err:", err.Error())
	}
	u.doUpdateConversation(cmd2Value{
		Cmd:   CmdUpdateConversation,
		Value: updateConNode{ConId: node.ConversationID, Action: ConAndUnreadChange},
	})
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
	case ConAndUnreadChange:
		err, list := u.getAllConversationListModel()
		if err == nil {
			if list == nil {
				u.ConversationListenerx.OnConversationChanged(structToJsonString([]ConversationStruct{}))
			} else {
				u.ConversationListenerx.OnConversationChanged(structToJsonString(list))

			}
			totalUnreadCount, err := u.getTotalUnreadMsgCountModel()
			if err == nil {
				u.ConversationListenerx.OnTotalUnreadMessageCountChanged(totalUnreadCount)
			}
		}
	case AddConOrUpLatMsg:
		c := node.Args.(ConversationStruct)
		if u.judgeConversationIfExists(node.ConId) {
			_, o := u.getOneConversationModel(node.ConId)
			if c.LatestMsgSendTime > o.LatestMsgSendTime { //The session update of asynchronous messages is subject to the latest sending time
				err := u.setConversationLatestMsgModel(&c, node.ConId)
				if err != nil {
					sdkLog("setConversationLatestMsgModel err: ", err)
				}
			}
		} else {
			_ = u.addConversationOrUpdateLatestMsg(&c, node.ConId)
			var list []*ConversationStruct
			list = append(list, &c)
			u.ConversationListenerx.OnNewConversation(structToJsonString(list))
		}

	case UnreadCountSetZero:
		if err := u.setConversationUnreadCount(0, node.ConId); err != nil {
		} else {
			err, list := u.getAllConversationListModel()
			if err == nil {
				if list == nil {
					u.ConversationListenerx.OnConversationChanged(structToJsonString([]ConversationStruct{}))
				} else {
					u.ConversationListenerx.OnConversationChanged(structToJsonString(list))

				}
				totalUnreadCount, err := u.getTotalUnreadMsgCountModel()
				if err == nil {
					u.ConversationListenerx.OnTotalUnreadMessageCountChanged(totalUnreadCount)
				}

			}
		}
	case ConChange:
		err, list := u.getAllConversationListModel()
		if err != nil {
			sdkLog("getAllConversationListModel database err:", err.Error())
		} else {
			if list == nil {
				u.ConversationListenerx.OnConversationChanged(structToJsonString([]ConversationStruct{}))
			} else {
				u.ConversationListenerx.OnConversationChanged(structToJsonString(list))

			}
		}
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
	case Sound:
		err = jsonStringToStruct(msg.Content, &msg.SoundElem)
	case Video:
		err = jsonStringToStruct(msg.Content, &msg.VideoElem)
	case File:
		err = jsonStringToStruct(msg.Content, &msg.FileElem)
	case AtText:
		err = jsonStringToStruct(msg.Content, &msg.AtElem)
		if err == nil {
			if isContain(u.LoginUid, msg.AtElem.AtUserList) {
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
