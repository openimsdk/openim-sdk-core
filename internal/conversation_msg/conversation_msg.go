package conversation_msg

import (
	"encoding/json"
	"open_im_sdk/internal/advanced_interface"
	"open_im_sdk/internal/cache"
	common2 "open_im_sdk/internal/common"
	"open_im_sdk/internal/friend"
	"open_im_sdk/internal/group"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/internal/organization"
	"open_im_sdk/internal/user"
	"open_im_sdk/internal/work_moments"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"sort"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
)

var SearchContentType = []int{constant.Text, constant.AtText, constant.File}

type Conversation struct {
	*ws.Ws
	db                   *db.DataBase
	p                    *ws.PostApi
	ConversationListener open_im_sdk_callback.OnConversationListener
	msgListener          open_im_sdk_callback.OnAdvancedMsgListener
	ch                   chan common.Cmd2Value
	loginUserID          string
	platformID           int32
	DataDir              string
	friend               *friend.Friend
	group                *group.Group
	user                 *user.User
	signaling            advanced_interface.Signaling
	advancedFunction     advanced_interface.AdvancedFunction
	organization         *organization.Organization
	workMoments          *workMoments.WorkMoments
	common2.ObjectStorage

	cache *cache.Cache
}

func (c *Conversation) SetAdvancedFunction(advancedFunction advanced_interface.AdvancedFunction) {
	c.advancedFunction = advancedFunction
}

func (c *Conversation) MsgListener() open_im_sdk_callback.OnAdvancedMsgListener {
	return c.msgListener
}

func (c *Conversation) SetSignaling(signaling advanced_interface.Signaling) {
	c.signaling = signaling
}

func (c *Conversation) SetMsgListener(msgListener open_im_sdk_callback.OnAdvancedMsgListener) {
	c.msgListener = msgListener
}
func NewConversation(ws *ws.Ws, db *db.DataBase, p *ws.PostApi,
	ch chan common.Cmd2Value, loginUserID string, platformID int32, dataDir string,
	friend *friend.Friend, group *group.Group, user *user.User,
	objectStorage common2.ObjectStorage, conversationListener open_im_sdk_callback.OnConversationListener,
	msgListener open_im_sdk_callback.OnAdvancedMsgListener, signaling advanced_interface.Signaling,
	advancedFunction advanced_interface.AdvancedFunction, organization *organization.Organization, workMoments *workMoments.WorkMoments, cache *cache.Cache) *Conversation {
	n := &Conversation{Ws: ws, db: db, p: p, ch: ch, loginUserID: loginUserID, platformID: platformID,
		DataDir: dataDir, friend: friend, group: group, user: user, ObjectStorage: objectStorage, signaling: signaling,
		advancedFunction: advancedFunction, organization: organization, workMoments: workMoments}
	n.SetMsgListener(msgListener)
	n.SetConversationListener(conversationListener)
	n.cache = cache
	return n
}

func (c *Conversation) GetCh() chan common.Cmd2Value {
	return c.ch
}

func (c *Conversation) doMsgNew(c2v common.Cmd2Value) {
	operationID := c2v.Value.(sdk_struct.CmdNewMsgComeToConversation).OperationID
	allMsg := c2v.Value.(sdk_struct.CmdNewMsgComeToConversation).MsgList
	if c.msgListener == nil {
		log.Error(operationID, "not set c MsgListenerList")
		return
	}
	var isTriggerUnReadCount bool
	var insertMsg, updateMsg []*db.LocalChatLog
	var exceptionMsg []*db.LocalErrChatLog
	var newMessages, msgReadList, groupMsgReadList, msgRevokeList sdk_struct.NewMsgList
	var isUnreadCount, isConversationUpdate, isHistory, isNotPrivate, isSenderConversationUpdate bool
	conversationChangedSet := make(map[string]*db.LocalConversation)
	newConversationSet := make(map[string]*db.LocalConversation)
	conversationSet := make(map[string]*db.LocalConversation)
	phConversationChangedSet := make(map[string]*db.LocalConversation)
	phNewConversationSet := make(map[string]*db.LocalConversation)
	log.Info(operationID, "do Msg come here")
	for _, v := range allMsg {
		isHistory = utils.GetSwitchFromOptions(v.Options, constant.IsHistory)
		isUnreadCount = utils.GetSwitchFromOptions(v.Options, constant.IsUnreadCount)
		isConversationUpdate = utils.GetSwitchFromOptions(v.Options, constant.IsConversationUpdate)
		isNotPrivate = utils.GetSwitchFromOptions(v.Options, constant.IsNotPrivate)
		isSenderConversationUpdate = utils.GetSwitchFromOptions(v.Options, constant.IsSenderConversationUpdate)
		msg := new(sdk_struct.MsgStruct)
		copier.Copy(msg, v)
		var tips server_api_params.TipsComm
		if v.ContentType >= constant.NotificationBegin && v.ContentType <= constant.NotificationEnd {
			_ = proto.Unmarshal(v.Content, &tips)
			marshaler := jsonpb.Marshaler{
				OrigName:     true,
				EnumsAsInts:  false,
				EmitDefaults: false,
			}
			msg.Content, _ = marshaler.MarshalToString(&tips)
		} else {
			msg.Content = string(v.Content)
		}
		//When the message has been marked and deleted by the cloud, it is directly inserted locally without any conversation and message update.
		if msg.Status == constant.MsgStatusHasDeleted {
			insertMsg = append(insertMsg, c.msgStructToLocalChatLog(msg))
			continue
		}
		msg.Status = constant.MsgStatusSendSuccess
		msg.IsRead = false
		//		log.Info(operationID, "new msg, seq, ServerMsgID, ClientMsgID", msg.Seq, msg.ServerMsgID, msg.ClientMsgID)
		//De-analyze data
		err := c.msgHandleByContentType(msg)
		if err != nil {
			log.Error(operationID, "Parsing data error:", err.Error())
			continue
		}
		if !isNotPrivate {
			msg.AttachedInfoElem.IsPrivateChat = true
			msg.AttachedInfo = utils.StructToJsonString(msg.AttachedInfoElem)
		}
		if msg.ClientMsgID == "" {
			exceptionMsg = append(exceptionMsg, c.msgStructToLocalErrChatLog(msg))
			continue
		}
		switch {
		case v.ContentType == constant.ConversationChangeNotification || v.ContentType == constant.ConversationPrivateChatNotification:
			log.Info(operationID, utils.GetSelfFuncName(), v)
			c.DoNotification(v)
		}
		switch v.SessionType {
		case constant.SingleChatType:
			if v.ContentType > constant.FriendNotificationBegin && v.ContentType < constant.FriendNotificationEnd {
				c.friend.DoNotification(v, c.ch)
				log.Info("internal", "DoFriendMsg SingleChatType", v)
			} else if v.ContentType > constant.UserNotificationBegin && v.ContentType < constant.UserNotificationEnd {
				log.Info("internal", "DoFriendMsg  DoUserMsg SingleChatType", v)
				c.user.DoNotification(v)
				c.friend.DoNotification(v, c.ch)
			} else if v.ContentType == constant.GroupApplicationRejectedNotification ||
				v.ContentType == constant.GroupApplicationAcceptedNotification ||
				v.ContentType == constant.JoinGroupApplicationNotification {
				log.Info("internal", "DoGroupMsg SingleChatType", v)
				c.group.DoNotification(v, c.ch)
			} else if v.ContentType > constant.SignalingNotificationBegin && v.ContentType < constant.SignalingNotificationEnd {
				log.Info(operationID, "signaling DoNotification ", v)
				c.signaling.DoNotification(v, c.ch, operationID)
				continue
			} else if v.ContentType == constant.OrganizationChangedNotification {
				log.Info(operationID, "Organization Changed Notification ")
				c.organization.DoNotification(v, c.ch, operationID)
			} else if v.ContentType == constant.WorkMomentNotification {
				log.Info(operationID, "WorkMoment New Notification")
				c.workMoments.DoNotification(tips.JsonDetail, operationID)
			}
		case constant.GroupChatType:
			if v.ContentType > constant.GroupNotificationBegin && v.ContentType < constant.GroupNotificationEnd {
				c.group.DoNotification(v, c.ch)
				log.Info(operationID, "DoGroupMsg SingleChatType", v)
			} else if v.ContentType > constant.SignalingNotificationBegin && v.ContentType < constant.SignalingNotificationEnd {
				log.Info(operationID, "signaling DoNotification ", v)
				c.signaling.DoNotification(v, c.ch, operationID)
				continue
			}
		}
		if v.SendID == c.loginUserID { //seq
			// Messages sent by myself  //if  sent through  this terminal
			m, err := c.db.GetMessage(msg.ClientMsgID)
			if err == nil {
				log.Info("internal", "have message", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, *msg)
				if m.Seq == 0 {
					if !isConversationUpdate {
						msg.Status = constant.MsgStatusFiltered
					}
					updateMsg = append(updateMsg, c.msgStructToLocalChatLog(msg))
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
				if msg.ContentType == constant.HasReadReceipt {
					msgReadList = append(msgReadList, msg)
				}
				if isConversationUpdate {
					if isSenderConversationUpdate {
						log.Debug(operationID, "updateConversation msg", v, lc)
						c.updateConversation(&lc, conversationSet)
					} else {
						//special fix
						c.updateConversation(&lc, conversationSet)
					}
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

				if msg.ContentType == constant.GroupHasReadReceipt {
					groupMsgReadList = append(groupMsgReadList, msg)

				}
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
				case constant.NotificationChatType:
					lc.ConversationID = utils.GetConversationIDBySessionType(v.SendID, constant.NotificationChatType)
					lc.UserID = v.SendID
				}
				if isUnreadCount {
					isTriggerUnReadCount = true
					lc.UnreadCount = 1
				}
				if isConversationUpdate {
					c.updateConversation(&lc, conversationSet)
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
				if msg.ContentType == constant.HasReadReceipt {
					msgReadList = append(msgReadList, msg)
				}
				if msg.ContentType == constant.GroupHasReadReceipt {
					groupMsgReadList = append(groupMsgReadList, msg)

				}
				if msg.ContentType == constant.Typing {
					newMessages = append(newMessages, msg)
				}
			} else {
				exceptionMsg = append(exceptionMsg, c.msgStructToLocalErrChatLog(msg))
				log.Warn(operationID, "Deduplication operation ", *c.msgStructToLocalErrChatLog(msg))
			}
		}
	}
	log.Info(operationID, "generate conversation map is :", conversationSet)
	list, err := c.db.GetAllConversationList()
	if err != nil {
		log.Error(operationID, "GetAllConversationList", "error", err.Error())
	}
	m := make(map[string]*db.LocalConversation)
	listToMap(list, m)
	log.Debug(operationID, "listToMap: ", list, conversationSet)
	c.diff(m, conversationSet, conversationChangedSet, newConversationSet)
	log.Info(operationID, "trigger map is :", "newConversations", newConversationSet, "changedConversations", conversationChangedSet)
	//seq sync message update
	err5 := c.db.BatchUpdateMessageList(updateMsg)
	if err5 != nil {
		log.Error(operationID, "sync seq normal message err  :", err5.Error())
	}
	//Normal message storage
	err1 := c.db.BatchInsertMessageList(insertMsg)
	if err1 != nil {
		log.Error(operationID, "insert GetMessage detail err:", err1.Error(), len(insertMsg))
		for _, v := range insertMsg {
			e := c.db.InsertMessage(v)
			if e != nil {
				errChatLog := &db.LocalErrChatLog{}
				copier.Copy(errChatLog, v)
				exceptionMsg = append(exceptionMsg, errChatLog)
				log.Warn(operationID, "InsertMessage operation ", "chat err log: ", errChatLog, "chat log: ", v, e.Error())
			}
		}
	}
	//Exception message storage
	for _, v := range exceptionMsg {
		log.Warn(operationID, "exceptionMsg show: ", *v)
	}

	err2 := c.db.BatchInsertExceptionMsgToErrorChatLog(exceptionMsg)
	if err2 != nil {
		log.Error(operationID, "insert err message err  :", err2.Error())

	}
	hList, _ := c.db.GetHiddenConversationList()
	for _, v := range hList {
		if nc, ok := newConversationSet[v.ConversationID]; ok {
			phConversationChangedSet[v.ConversationID] = nc
		}
	}
	for k, v := range newConversationSet {
		if _, ok := phConversationChangedSet[v.ConversationID]; !ok {
			phNewConversationSet[k] = v
		}
	}
	//Changed conversation storage
	err3 := c.db.BatchUpdateConversationList(append(mapConversationToList(conversationChangedSet), mapConversationToList(phConversationChangedSet)...))
	if err3 != nil {
		log.Error(operationID, "insert changed conversation err :", err3.Error())
	}
	//New conversation storage
	err4 := c.db.BatchInsertConversationList(mapConversationToList(phNewConversationSet))
	if err4 != nil {
		log.Error(operationID, "insert new conversation err:", err4.Error())
	}
	c.doMsgReadState(msgReadList)
	c.advancedFunction.DoGroupMsgReadState(groupMsgReadList)
	c.revokeMessage(msgRevokeList)
	c.newMessage(newMessages)
	//log.Info(operationID, "trigger map is :", newConversationSet, conversationChangedSet)
	if len(newConversationSet) != 0 {
		c.ConversationListener.OnNewConversation(utils.StructToJsonString(mapConversationToList(newConversationSet)))
	}
	if len(conversationChangedSet) != 0 {
		c.ConversationListener.OnConversationChanged(utils.StructToJsonString(mapConversationToList(conversationChangedSet)))
	}

	if isTriggerUnReadCount {
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{"", constant.TotalUnreadMessageChanged, ""}})
	}
}
func listToMap(list []*db.LocalConversation, m map[string]*db.LocalConversation) {
	for _, v := range list {
		m[v.ConversationID] = v
	}

}
func (c *Conversation) diff(local, generated, cc, nc map[string]*db.LocalConversation) {
	for _, v := range generated {
		log.Debug("node diff", *v)
		if localC, ok := local[v.ConversationID]; ok {
			if v.LatestMsgSendTime > localC.LatestMsgSendTime {
				localC.UnreadCount = localC.UnreadCount + v.UnreadCount
				localC.LatestMsg = v.LatestMsg
				localC.LatestMsgSendTime = v.LatestMsgSendTime
				cc[v.ConversationID] = localC
				log.Debug("", "diff1 ", *localC, *v)
			} else {
				localC.UnreadCount = localC.UnreadCount + v.UnreadCount
				cc[v.ConversationID] = localC
				log.Debug("", "diff2 ", *localC, *v)
			}

		} else {
			c.addFaceURLAndName(v)
			nc[v.ConversationID] = v
			log.Debug("", "diff3 ", *v)
		}
	}

}
func (c *Conversation) genConversationGroupAtType(lc *db.LocalConversation, s *sdk_struct.MsgStruct) {
	if s.ContentType == constant.AtText {
		tagMe := utils.IsContain(c.loginUserID, s.AtElem.AtUserList)
		tagAll := utils.IsContain(constant.AtAllString, s.AtElem.AtUserList)
		if tagAll {
			if tagMe {
				lc.GroupAtType = constant.AtAllAtMe
				return
			}
			lc.GroupAtType = constant.AtAll
			return
		}
		if tagMe {
			lc.GroupAtType = constant.AtMe
		}
	}

}
func (c *Conversation) msgStructToLocalChatLog(m *sdk_struct.MsgStruct) *db.LocalChatLog {
	var lc db.LocalChatLog
	copier.Copy(&lc, m)
	if m.SessionType == constant.GroupChatType {
		lc.RecvID = m.GroupID
	}
	return &lc
}
func (c *Conversation) msgStructToLocalErrChatLog(m *sdk_struct.MsgStruct) *db.LocalErrChatLog {
	var lc db.LocalErrChatLog
	copier.Copy(&lc, m)
	if m.SessionType == constant.GroupChatType {
		lc.RecvID = m.GroupID
	}
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
func (c *Conversation) newMessage(newMessagesList sdk_struct.NewMsgList) {
	sort.Sort(newMessagesList)
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
			m, err := c.db.GetMessage(v)
			if err != nil {
				log.Error("internal", "GetMessage err:", err, "ClientMsgID", v)
				continue
			}
			attachInfo := sdk_struct.AttachedInfoElem{}
			_ = utils.JsonStringToStruct(m.AttachedInfo, &attachInfo)
			attachInfo.HasReadTime = rd.SendTime
			m.AttachedInfo = utils.StructToJsonString(attachInfo)
			m.IsRead = true
			err = c.db.UpdateMessage(m)
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
			msgRt.SessionType = constant.SingleChatType
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
		var list []*db.LocalConversation
		lc := node.Args.(db.LocalConversation)
		oc, err := c.db.GetConversation(lc.ConversationID)
		if err == nil {
			log.Info("this is old conversation", *oc)
			if lc.LatestMsgSendTime >= oc.LatestMsgSendTime { //The session update of asynchronous messages is subject to the latest sending time
				err := c.db.UpdateColumnsConversation(node.ConID, map[string]interface{}{"latest_msg_send_time": lc.LatestMsgSendTime, "latest_msg": lc.LatestMsg})
				if err != nil {
					log.Error("internal", "updateConversationLatestMsgModel err: ", err)
				} else {
					oc.LatestMsgSendTime = lc.LatestMsgSendTime
					oc.LatestMsg = lc.LatestMsg
					list = append(list, oc)
					c.ConversationListener.OnConversationChanged(utils.StructToJsonString(list))
				}
			}
		} else {
			log.Info("this is new conversation", lc)
			err4 := c.db.InsertConversation(&lc)
			if err4 != nil {
				log.Error("internal", "insert new conversation err:", err4.Error())
			} else {
				list = append(list, &lc)
				c.ConversationListener.OnNewConversation(utils.StructToJsonString(list))
			}
		}

	case constant.UnreadCountSetZero:
		if err := c.db.UpdateColumnsConversation(node.ConID, map[string]interface{}{"unread_count": 0}); err != nil {
			log.Error("internal", "UpdateColumnsConversation err", err.Error())
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
		err := c.db.IncrConversationUnreadCount(node.ConID)
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
		var lc db.LocalConversation
		st := node.Args.(common.SourceIDAndSessionType)
		lc.ConversationID = node.ConID
		lc.ConversationType = st.SessionType
		if st.SessionType == constant.SingleChatType {
			lc.UserID = st.SourceID
		} else {
			lc.GroupID = st.SourceID
		}
		c.addFaceURLAndName(&lc)
		err := c.db.UpdateConversation(&lc)
		if err != nil {
			log.Error("internal", "setConversationFaceUrlAndNickName database err:", err.Error())
			return
		}

	case constant.UpdateLatestMessageChange:
		conversationID := node.ConID
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
				err = c.db.UpdateColumnsConversation(node.ConID, map[string]interface{}{"latest_msg_send_time": latestMsg.SendTime, "latest_msg": newLatestMessage})
				if err != nil {
					log.Error("internal", "updateConversationLatestMsgModel err :", err.Error())
				}
			}
		}
	case constant.ConChange:
		cidList := node.Args.([]string)
		cLists, err := c.db.GetMultipleConversation(cidList)
		if err != nil {
			log.Error("internal", "getMultipleConversationModel err :", err.Error())
		} else {
			var newCList []*db.LocalConversation
			for _, v := range cLists {
				if v.LatestMsgSendTime != 0 {
					newCList = append(newCList, v)
				}
			}
			log.Info("internal", "getMultipleConversationModel success :", newCList)

			c.ConversationListener.OnConversationChanged(utils.StructToJsonStringDefault(newCList))
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
func (c *Conversation) msgConvert(msg *sdk_struct.MsgStruct) (err error) {
	err = c.msgHandleByContentType(msg)
	if err != nil {
		return err
	} else {
		if msg.SessionType == constant.GroupChatType {
			msg.GroupID = msg.RecvID
			msg.RecvID = c.loginUserID
		}
		return nil
	}
}

func (c *Conversation) msgHandleByContentType(msg *sdk_struct.MsgStruct) (err error) {
	_ = utils.JsonStringToStruct(msg.AttachedInfo, &msg.AttachedInfoElem)
	if msg.ContentType >= constant.NotificationBegin && msg.ContentType <= constant.NotificationEnd {
		var tips server_api_params.TipsComm
		err = utils.JsonStringToStruct(msg.Content, &tips)
		msg.NotificationElem.Detail = tips.JsonDetail
		msg.NotificationElem.DefaultTips = tips.DefaultTips

	} else {
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
		case constant.Face:
			err = utils.JsonStringToStruct(msg.Content, &msg.FaceElem)

		}
	}

	return err
}
func (c *Conversation) updateConversation(lc *db.LocalConversation, cs map[string]*db.LocalConversation) {
	if oldC, ok := cs[lc.ConversationID]; !ok {
		cs[lc.ConversationID] = lc
	} else {
		if lc.LatestMsgSendTime > oldC.LatestMsgSendTime {
			oldC.UnreadCount = oldC.UnreadCount + lc.UnreadCount
			oldC.LatestMsg = lc.LatestMsg
			oldC.LatestMsgSendTime = lc.LatestMsgSendTime
			cs[lc.ConversationID] = oldC
		} else {
			oldC.UnreadCount = oldC.UnreadCount + lc.UnreadCount
			cs[lc.ConversationID] = oldC
		}
	}
	//if oldC, ok := cc[lc.ConversationID]; !ok {
	//	oc, err := c.db.GetConversation(lc.ConversationID)
	//	if err == nil && oc.ConversationID != "" {//如果会话已经存在
	//		if lc.LatestMsgSendTime > oc.LatestMsgSendTime {
	//			oc.UnreadCount = oc.UnreadCount + lc.UnreadCount
	//			oc.LatestMsg = lc.LatestMsg
	//			oc.LatestMsgSendTime = lc.LatestMsgSendTime
	//			cc[lc.ConversationID] = *oc
	//		} else {
	//			oc.UnreadCount = oc.UnreadCount + lc.UnreadCount
	//			cc[lc.ConversationID] = *oc
	//		}
	//	} else {
	//		if oldC, ok := nc[lc.ConversationID]; !ok {
	//			c.addFaceURLAndName(lc)
	//			nc[lc.ConversationID] = *lc
	//		} else {
	//			if lc.LatestMsgSendTime > oldC.LatestMsgSendTime {
	//				oldC.UnreadCount = oldC.UnreadCount + lc.UnreadCount
	//				oldC.LatestMsg = lc.LatestMsg
	//				oldC.LatestMsgSendTime = lc.LatestMsgSendTime
	//				nc[lc.ConversationID] = oldC
	//			} else {
	//				oldC.UnreadCount = oldC.UnreadCount + lc.UnreadCount
	//				nc[lc.ConversationID] = oldC
	//			}
	//		}
	//	}
	//} else {
	//	if lc.LatestMsgSendTime > oldC.LatestMsgSendTime {
	//		oldC.UnreadCount = oldC.UnreadCount + lc.UnreadCount
	//		oldC.LatestMsg = lc.LatestMsg
	//		oldC.LatestMsgSendTime = lc.LatestMsgSendTime
	//		cc[lc.ConversationID] = oldC
	//	} else {
	//		oldC.UnreadCount = oldC.UnreadCount + lc.UnreadCount
	//		cc[lc.ConversationID] = oldC
	//	}
	//
	//}

}
func mapConversationToList(m map[string]*db.LocalConversation) (cs []*db.LocalConversation) {
	for _, v := range m {
		cs = append(cs, v)
	}
	return cs
}
func (c *Conversation) addFaceURLAndName(lc *db.LocalConversation) {
	operationID := utils.OperationIDGenerator()
	switch lc.ConversationType {
	case constant.SingleChatType, constant.NotificationChatType:
		faceUrl, name, err, isFromSvr := c.friend.GetUserNameAndFaceUrlByUid(lc.UserID, operationID)
		if err != nil {
			log.Error(operationID, "getUserNameAndFaceUrlByUid err", err.Error(), lc.UserID)
			return
		}
		lc.FaceURL = faceUrl
		lc.ShowName = name
		if isFromSvr {
			c.cache.Update(lc.UserID, faceUrl, name)
		}

	case constant.GroupChatType:
		g, err := c.group.GetGroupInfoFromLocal2Svr(lc.GroupID)
		if err != nil {
			log.Error(operationID, "GetGroupInfoByGroupID err", err.Error(), lc.GroupID)
			return
		}
		lc.ShowName = g.GroupName
		lc.FaceURL = g.FaceURL

	}
}
