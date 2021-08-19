package open_im_sdk

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	ConversationListener OnConversationListener
	MsgListenerList      []OnAdvancedMsgListener
	ch                   chan cmd2Value
}

func (con ConversationListener) getCh() chan cmd2Value {
	return con.ch
}

func (con *ConversationListener) doMsg(c2v cmd2Value) {
	if con.ConversationListener == nil || con.MsgListenerList == nil {
		fmt.Println("not set conversationListener or MsgListenerList", len(con.MsgListenerList))
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
		err := msgHandleByContentType(msg)
		if err != nil {
			fmt.Println("Parsing data error:", err.Error(), msg)
		}
		switch v.SessionType {
		case SingleChatType:
			msg.RecvID = v.RecvID
		case GroupChatType:
			msg.RecvID = strings.Split(v.RecvID, " ")[1]
			msg.GroupID = msg.RecvID
		}
		if v.SendID == LoginUid { //seq對齊消息 Messages sent by myself
			if judgeMessageIfExists(msg) { //if  sent through  this terminal
				err := updateMessageSeq(msg)
				if err != nil {
					fmt.Println("err", err.Error())
				}
			} else { //同步消息       send through  other terminal
				_ = insertPushMessageToChatLog(msg)
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
				case GroupChatType:
					c.GroupID = strings.Split(v.RecvID, " ")[1]
					c.ConversationID = GetConversationIDBySessionType(c.GroupID, GroupChatType)
				}

				if msg.ContentType <= AcceptFriendApplicationTip {
					newMessages = append(newMessages, msg)
					con.doUpdateConversation(cmd2Value{Value: updateConNode{c.ConversationID, AddConOrUpLatMsg,
						c}})
				}
				//}
			}
		} else { //他人發的
			if !judgeMessageIfExists(msg) { //去重操作
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
					_ = insertPushMessageToChatLog(msg)
					switch v.SessionType {
					case SingleChatType:
						c.ConversationID = GetConversationIDBySessionType(v.SendID, SingleChatType)
						c.UserID = v.SendID
						c.ShowName = msg.SenderNickName
						c.FaceURL = msg.SenderFaceURL
					case GroupChatType:
						c.GroupID = strings.Split(v.RecvID, " ")[1]
						c.ConversationID = GetConversationIDBySessionType(c.GroupID, GroupChatType)
						faceUrl, name, err := getGroupNameAndFaceUrlByUid(c.GroupID)
						if err != nil {
							sdkLog("getGroupNameAndFaceUrlByUid err:", err)
						} else {
							c.ShowName = name
							c.FaceURL = faceUrl
						}
					}
					if msg.ContentType <= AcceptFriendApplicationTip || (msg.ContentType >= GroupTipBegin && msg.ContentType <= GroupTipEnd && msg.ContentType != SetGroupInfoTip && msg.ContentType != JoinGroupTip) {
						fmt.Println("tttttttttttttttttttt", msg.ContentType)
						con.doUpdateConversation(cmd2Value{Value: updateConNode{c.ConversationID, AddConOrUpLatMsg,
							c}})
						if msg.ContentType != Revoke {
							con.doUpdateConversation(cmd2Value{Value: updateConNode{c.ConversationID, IncrUnread, ""}})
						}
						con.doUpdateConversation(cmd2Value{Value: updateConNode{c.ConversationID, UpdateFaceUrlAndNickName, c}})
						newMessages = append(newMessages, msg)
					}
					if msg.ContentType == SetGroupInfoTip || msg.ContentType == SetSelfInfoTip {
						con.doUpdateConversation(cmd2Value{Value: updateConNode{c.ConversationID, UpdateFaceUrlAndNickName, c}})
						newMessages = append(newMessages, msg)

					}
					if msg.ContentType == Revoke {
						msgRevokeList = append(msgRevokeList, msg)
					}
				} else {
					if msg.ContentType == Typing {
						newMessages = append(newMessages, msg)

					} else {
						//update read state
						msgReadList = append(msgReadList, msg)
					}
				}
			}
		}
	}
	con.doMsgReadState(msgReadList)
	con.revokeMessage(msgRevokeList)
	con.newMessage(newMessages)
	con.doUpdateConversation(cmd2Value{Value: updateConNode{"", ConChange, ""}})
	con.doUpdateConversation(cmd2Value{Value: updateConNode{"", TotalUnreadMessageChanged, ""}})
	fmt.Println("length msgListenerList", con.MsgListenerList, "length message", len(newMessages), "msgListenerLen", len(con.MsgListenerList))

}

func (con *ConversationListener) revokeMessage(msgRevokeList []*MsgStruct) {
	for _, v := range con.MsgListenerList {
		for _, w := range msgRevokeList {
			if v != nil {
				err := setMessageStatus(w.Content, MsgStatusRevoked)
				if err != nil {
					sdkLog("setLocalMessageStatus revokeMessage err:", err.Error(), "msg", w)
				} else {
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
			if v != nil {
				fmt.Println("msgListener,OnRecvNewMessage")
				v.OnRecvNewMessage(structToJsonString(w))
			} else {
				fmt.Println("set msgListener is err ")
			}
		}
	}
}
func (con *ConversationListener) doDeleteConversation(c2v cmd2Value) {
	if con.ConversationListener == nil {
		fmt.Println("not set conversationListener")
		return
	}
	node := c2v.Value.(deleteConNode)
	//标记删除与此会话相关的消息
	err := setMessageStatusBySourceID(node.SourceID, MsgStatusHasDeleted, node.SessionType)
	if err != nil {
		sdkLog("setMessageStatusBySourceID err:", err.Error())
		return
	}
	//重置该会话信息，空会话
	err = ResetConversation(node.ConversationID)
	if err != nil {
		sdkLog("ResetConversation err:", err.Error())
	}
	con.doUpdateConversation(cmd2Value{
		Cmd:   CmdUpdateConversation,
		Value: updateConNode{ConId: node.ConversationID, Action: ConAndUnreadChange},
	})
}
func (con *ConversationListener) doMsgReadState(msgReadList []*MsgStruct) {
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
			err := setMessageHasReadByMsgID(v)
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
		for _, v := range con.MsgListenerList {
			sdkLog("OnRecvC2CReadReceipt: ", structToJsonString(messageReceiptResp))
			v.OnRecvC2CReadReceipt(structToJsonString(messageReceiptResp))
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
				sdkLog("setConversationLatestMsgModel err: ", err)
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
			sdkLog("getAllConversationListModel database err:", err.Error())
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
	case TotalUnreadMessageChanged:
		totalUnreadCount, err := getTotalUnreadMsgCountModel()
		if err != nil {
			sdkLog("TotalUnreadMessageChanged database err:", err.Error())
		} else {
			con.ConversationListener.OnTotalUnreadMessageCountChanged(totalUnreadCount)
		}
	case UpdateFaceUrlAndNickName:
		c := node.Args.(ConversationStruct)
		err := setConversationFaceUrlAndNickName(&c, node.ConId)
		if err != nil {
			sdkLog("setConversationFaceUrlAndNickName database err:", err.Error())
			return
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

func msgHandleByContentType(msg *MsgStruct) (err error) {
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
			if isContain(LoginUid, msg.AtElem.AtUserList) {
				msg.AtElem.IsAtSelf = true
			}
		}
	case Location:
		err = jsonStringToStruct(msg.Content, &msg.LocationElem)
	case Custom:
		err = jsonStringToStruct(msg.Content, &msg.CustomElem)
	}
	return err
}
func getGroupNameAndFaceUrlByUid(groupID string) (faceUrl, name string, err error) {
	groupInfo, err := getLocalGroupsInfoByGroupID(groupID)
	if err != nil {
		return "", "", err
	}
	if groupInfo.GroupId == "" {
		groupInfo, err := groupManager.getGroupInfoByGroupId(groupID)
		if err != nil {
			return "", "", err
		} else {
			return groupInfo.FaceUrl, groupInfo.GroupName, nil
		}
	} else {
		return groupInfo.FaceUrl, groupInfo.GroupName, nil
	}
}
