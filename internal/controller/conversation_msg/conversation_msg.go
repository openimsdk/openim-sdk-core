package conversation_msg

import (
	"database/sql"
	"encoding/json"
	"github.com/jinzhu/copier"
	"open_im_sdk/internal/controller/friend"
	"open_im_sdk/internal/controller/group"
	ws "open_im_sdk/internal/controller/interaction"
	"open_im_sdk/internal/controller/user"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
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
type Conversation struct {
	*ws.Ws
	*db.DataBase
	ConversationListener OnConversationListener
	MsgListenerList       []OnAdvancedMsgListener
	ch                    chan common.Cmd2Value
	loginUserID           string
	friend                *friend.Friend
	group                 *group.Group
	user                  *user.User
}

func NewConversation(ConversationListener OnConversationListener, msgListenerList []OnAdvancedMsgListener, ch chan common.Cmd2Value, loginUserID string, ws *ws.Ws) *Conversation {
	return &Conversation{ConversationListener: ConversationListener, MsgListenerList: msgListenerList, ch: ch, loginUserID: loginUserID, Ws: ws}
}

func (c *Conversation) getCh() chan common.Cmd2Value {
	return c.ch
}

func (c *Conversation) doMsgNew(c2v common.Cmd2Value) {
	if c.MsgListenerList == nil {
		log.Error("internal", "not set c MsgListenerList", len(c.MsgListenerList))
		return
	}
	var insertMsg []*db.LocalChatLog
	var errMsg, newMessages, msgReadList, msgRevokeList []*utils.MsgStruct
	var isUnreadCount, isConversationUpdate, isHistory bool
	var isCallbackUI bool
	conversationChangedSet := make(map[string]db.LocalConversation)
	newConversationSet := make(map[string]db.LocalConversation)
	//MsgList := c2v.Value.(ArrMsg)c
	//for _, v := range MsgList.GroupData {
	//	MsgList.SingleData = append(MsgList.SingleData, v)
	//}
	log.Info("internal", "do Msg come here")
	for _, v := range c.SeqMsg() {
		isHistory = utils.GetSwitchFromOptions(v.Options, constant.IsHistory)
		isUnreadCount = utils.GetSwitchFromOptions(v.Options, constant.IsUnreadCount)
		isConversationUpdate = utils.GetSwitchFromOptions(v.Options, constant.IsConversationUpdate)
		isCallbackUI = true
		msg := new(utils.MsgStruct)
		copier.Copy(msg, v)
		msg.Content = string(v.Content)
		msg.Status = constant.MsgStatusSendSuccess
		msg.IsRead = false
		log.Info("internal", "new msg, seq, ServerMsgID, ClientMsgID", msg.Seq, msg.ServerMsgID, msg.ClientMsgID)
		//De-analyze data
		err := c.msgHandleByContentType(msg)
		if err != nil {
			log.Error("internal", "Parsing data error:", err.Error())
			continue
		}
		switch v.SessionType {
		case constant.SingleChatType:
			if v.ContentType > constant.SingleTipBegin && v.ContentType < constant.SingleTipEnd {
				c.friend.DoFriendMsg(&v)
				log.Info("internal", "DoFriendMsg SingleChatType", v)
			} else if v.ContentType > constant.GroupTipBegin && v.ContentType < constant.GroupTipEnd {
				c.group.DoGroupMsg(&v)
				log.Info("internal", "DoGroupMsg SingleChatType", v)
			}
		case constant.GroupChatType:
			if v.ContentType > constant.GroupTipBegin && v.ContentType < constant.GroupTipEnd {
				c.group.DoGroupMsg(&v)
				log.Info("internal", "DoGroupMsg GroupChatType", v)
			}
		}
		if v.SendID == c.loginUserID { //seq  Messages sent by myself  //if  sent through  this terminal
			m, err := c.GetMessage(msg.ClientMsgID)
			if err == nil && m != nil {
				log.Info("internal", "have message", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, *msg)
				if m.Seq == 0 {
					insertMsg = append(insertMsg, c.msgStructToLocalChatLog(msg))
				} else {
					errMsg = append(errMsg, msg)

				}
			} else { //      send through  other terminal
				log.Info("internal", "sync message", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, *msg)
				lc := db.LocalConversation{
					ConversationType:  v.SessionType,
					LatestMsg:         utils.StructToJsonString(msg),
					LatestMsgSendTime: utils.UnixNanoSecondToTime(msg.SendTime),
				}
				switch v.SessionType {
				case constant.SingleChatType:
					lc.ConversationID = utils.GetConversationIDBySessionType(v.RecvID, constant.SingleChatType)
					lc.UserID = v.RecvID
					//localUserInfo,_ := c.user.GetLoginUser()
					//c.FaceURL = localUserInfo.FaceUrl
					//c.ShowName = localUserInfo.Nickname
				case constant.GroupChatType:
					lc.GroupID = v.GroupID
					lc.ConversationID = utils.GetConversationIDBySessionType(lc.GroupID, constant.GroupChatType)
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
			if b, _ := c.MessageIfExists(msg.ClientMsgID); !b { //Deduplication operation
				lc := db.LocalConversation{
					ConversationType:  v.SessionType,
					LatestMsg:         utils.StructToJsonString(msg),
					LatestMsgSendTime: utils.UnixNanoSecondToTime(msg.SendTime),
				}

				switch v.SessionType {
				case constant.SingleChatType:
					lc.ConversationID = utils.GetConversationIDBySessionType(v.SendID, constant.SingleChatType)
					lc.UserID = v.SendID
					lc.ShowName = msg.SenderNickname
					lc.FaceURL = msg.SenderFaceURL
				case constant.GroupChatType:
					lc.GroupID = v.GroupID
					lc.ConversationID = utils.GetConversationIDBySessionType(lc.GroupID, constant.GroupChatType)
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
				errMsg = append(errMsg, msg)
			}
		}
	}
	//Normal message storage
	err1 := c.BatchInsertMessageList(insertMsg)
	if err1 != nil {
		log.Error("internal", "insert normal message err  :", err1.Error())
	}
	//Exception message storage
	//err2, emsg2 := u.batchInsertErrorMessageToErrorChatLog(errMsg)
	//if err2 != nil {
	//	utils.sdkLog("insert err message err  :", err2.Error(), emsg2)
	//}
	//Changed conversation storage
	err3 := c.BatchUpdateConversationList(mapConversationToList(conversationChangedSet))
	if err3 != nil {
		log.Error("internal", "insert changed conversation err :", err3.Error())
	}
	//New conversation storage
	err4 := c.BatchInsertConversationList(mapConversationToList(newConversationSet))
	if err4 != nil {
		log.Error("internal", "insert new conversation err:", err4.Error())

	}
	//clear cache
	seqMap := make(map[int32]server_api_params.MsgData)
	c.SetSeqMsg(seqMap)
	if isCallbackUI {
		c.doMsgReadState(msgReadList)
		c.revokeMessage(msgRevokeList)
		c.newMessage(newMessages)
		//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", ConChange, ""}})
		log.Info("internal", "trigger map is :", newConversationSet, conversationChangedSet)
		//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewCon, mapKeyToStringList(newConversationSet)}})
		//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, mapKeyToStringList(conversationChangSet)}})
		c.ConversationListenerx.OnConversationChanged(utils.StructToJsonString(mapConversationToList(conversationChangedSet)))
		c.ConversationListenerx.OnNewConversation(utils.StructToJsonString(mapConversationToList(newConversationSet)))
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{"", constant.TotalUnreadMessageChanged, ""}})
	}
	//sdkLog("length msgListenerList", u.MsgListenerList, "length message", len(newMessages), "msgListenerLen", len(u.MsgListenerList))

}
func (c *Conversation) msgStructToLocalChatLog(m *utils.MsgStruct) *db.LocalChatLog {
	var lc db.LocalChatLog
	copier.Copy(&lc, m)
	lc.SendTime = utils.UnixNanoSecondToTime(m.SendTime)
	lc.CreateTime = utils.UnixNanoSecondToTime(m.CreateTime)

}

func (c *Conversation) revokeMessage(msgRevokeList []*utils.MsgStruct) {
	for _, v := range c.MsgListenerList {
		for _, w := range msgRevokeList {
			if v != nil {
				t := new(db.LocalChatLog)
				t.ClientMsgID = w.Content
				t.Status = constant.MsgStatusRevoked
				err := c.UpdateMessage(t)
				if err != nil {
					log.Error("internal", "setLocalMessageStatus revokeMessage err:", err.Error(), "msg", w)
				} else {
					log.Info("internal", "v.OnRecvMessageRevoked client_msg_id:", w.Content)
					v.OnRecvMessageRevoked(w.Content)
				}
			} else {
				log.Error("internal", "set msgListener is err:")
			}
		}
	}
}
func (c *Conversation) newMessage(newMessagesList []*utils.MsgStruct) {
	for _, v := range c.MsgListenerList {
		for _, w := range newMessagesList {
			log.Info("internal", "newMessage: ", w.ClientMsgID)
			if v != nil {
				log.Info("internal", "msgListener,OnRecvNewMessage")
				v.OnRecvNewMessage(utils.StructToJsonString(w))
			} else {
				log.Error("internal", "set msgListener is err ", len(c.MsgListenerList))
			}
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
	err := u.setMessageStatusBySourceID(node.SourceID, constant.MsgStatusHasDeleted, node.SessionType)
	if err != nil {
		utils.sdkLog("setMessageStatusBySourceID err:", err.Error())
		return
	}
	//Reset the session information, empty session
	err = u.ResetConversation(node.ConversationID)
	if err != nil {
		utils.sdkLog("ResetConversation err:", err.Error())
	}
	u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{"", constant.TotalUnreadMessageChanged, ""}})
}
func (c *Conversation) doMsgReadState(msgReadList []*utils.MsgStruct) {
	var messageReceiptResp []*utils.MessageReceipt
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
			msgRt := new(utils.MessageReceipt)
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

func (c *Conversation) doUpdateConversation(c2v common.Cmd2Value) {
	if c.ConversationListener == nil {
		log.Error("internal", "not set conversationListener")
		return
	}
	node := c2v.Value.(common.updateConNode)
	switch node.Action {
	case constant.AddConOrUpLatMsg:
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

	case constant.UnreadCountSetZero:
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
	case constant.IncrUnread:
		err := u.incrConversationUnreadCount(node.ConId)
		if err != nil {
			utils.sdkLog("incrConversationUnreadCount database err:", err.Error())
			return
		}
	case constant.TotalUnreadMessageChanged:
		totalUnreadCount, err := u.getTotalUnreadMsgCountModel()
		if err != nil {
			utils.sdkLog("TotalUnreadMessageChanged database err:", err.Error())
		} else {
			u.ConversationListenerx.OnTotalUnreadMessageCountChanged(totalUnreadCount)
		}
	case constant.UpdateFaceUrlAndNickName:
		c := node.Args.(ConversationStruct)
		if c.ShowName != "" || c.FaceURL != "" {
			err := u.setConversationFaceUrlAndNickName(&c, node.ConId)
			if err != nil {
				utils.sdkLog("setConversationFaceUrlAndNickName database err:", err.Error())
				return
			}
		}

	case constant.UpdateLatestMessageChange:
		conversationID := node.ConId
		var latestMsg utils.MsgStruct
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
	case constant.NewConChange:
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
	case constant.NewCon:
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

func (c *Conversation) work(c2v common.Cmd2Value) {

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

func (c *Conversation) msgHandleByContentType(msg *utils.MsgStruct) (err error) {
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

//func (c *Conversation) getGroupNameAndFaceUrlByUid(groupID string) (faceUrl, name string, err error) {
//	groupInfo, err := u.getLocalGroupsInfoByGroupID(groupID)
//	if err != nil {
//		return "", "", err
//	}
//	if groupInfo.GroupId == "" {
//		groupInfo, err := u.getGroupInfoByGroupId(groupID)
//		if err != nil {
//			return "", "", err
//		} else {
//			return groupInfo.FaceUrl, groupInfo.GroupName, nil
//		}
//	} else {
//		return groupInfo.FaceUrl, groupInfo.GroupName, nil
//	}
//}
func (c *Conversation) updateConversation(lc *db.LocalConversation, cc, nc map[string]db.LocalConversation) {
	b, err := c.ConversationIfExists(lc.ConversationID)
	if err != nil {
		log.Error("internal", lc, cc, nc, err.Error())
		return
	}
	if b {
		//_, o := u.getOneConversationModel(c.ConversationID)
		//if c.LatestMsgSendTime > o.LatestMsgSendTime { //The session update of asynchronous messages is subject to the latest sending time
		//	err := u.updateConversationLatestMsgModel(c.LatestMsgSendTime, c.LatestMsg, c.ConversationID)
		//	if err != nil {
		//		sdkLog("updateConversationLatestMsgModel err: ", err)
		//	} else {
		//		cc[c.ConversationID] = void{}
		//	}
		//}
		if oldC, ok := cc[lc.ConversationID]; ok {
			if oldC.LatestMsgSendTime.Before(lc.LatestMsgSendTime) {
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
			if oldC.LatestMsgSendTime.Before(lc.LatestMsgSendTime) {
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
func mapConversationToList(m map[string]db.LocalConversation) (cs []*db.LocalConversation) {
	for _, v := range m {
		cs = append(cs, &v)
	}
	return cs
}
