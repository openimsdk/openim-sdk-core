package conversation_msg

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"open_im_sdk/internal/business"
	"open_im_sdk/internal/cache"
	common2 "open_im_sdk/internal/common"
	"open_im_sdk/internal/friend"
	"open_im_sdk/internal/full"
	"open_im_sdk/internal/group"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/internal/organization"
	"open_im_sdk/internal/signaling"
	"open_im_sdk/internal/user"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	"runtime"
	"strings"

	workMoments "open_im_sdk/internal/work_moments"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"sort"
	"sync"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
)

var SearchContentType = []int{constant.Text, constant.AtText, constant.File}

type Conversation struct {
	*ws.Ws
	db                   db_interface.DataBase
	p                    *ws.PostApi
	ConversationListener open_im_sdk_callback.OnConversationListener
	msgListener          open_im_sdk_callback.OnAdvancedMsgListener
	msgKvListener        open_im_sdk_callback.OnMessageKvInfoListener
	batchMsgListener     open_im_sdk_callback.OnBatchMsgListener
	recvCH               chan common.Cmd2Value
	loginUserID          string
	platformID           int32
	DataDir              string
	friend               *friend.Friend
	group                *group.Group
	user                 *user.User
	signaling            *signaling.LiveSignaling
	//advancedFunction     advanced_interface.AdvancedFunction
	organization *organization.Organization
	workMoments  *workMoments.WorkMoments
	business     *business.Business
	common2.ObjectStorage

	cache          *cache.Cache
	full           *full.Full
	tempMessageMap sync.Map
	encryptionKey  string

	id2MinSeq            map[string]uint32
	IsExternalExtensions bool

	listenerForService open_im_sdk_callback.OnListenerForService
}

//func (c *Conversation) SetAdvancedFunction(advancedFunction advanced_interface.AdvancedFunction) {
//	c.advancedFunction = advancedFunction
//}

func (c *Conversation) SetListenerForService(listener open_im_sdk_callback.OnListenerForService) {
	c.listenerForService = listener
}

func (c *Conversation) MsgListener() open_im_sdk_callback.OnAdvancedMsgListener {
	return c.msgListener
}

func (c *Conversation) SetSignaling(signaling *signaling.LiveSignaling) {
	c.signaling = signaling
}

func (c *Conversation) SetMsgListener(msgListener open_im_sdk_callback.OnAdvancedMsgListener) {
	c.msgListener = msgListener
}
func (c *Conversation) SetMsgKvListener(msgKvListener open_im_sdk_callback.OnMessageKvInfoListener) {
	c.msgKvListener = msgKvListener
}
func (c *Conversation) SetBatchMsgListener(batchMsgListener open_im_sdk_callback.OnBatchMsgListener) {
	c.batchMsgListener = batchMsgListener
}

func NewConversation(ws *ws.Ws, db db_interface.DataBase, p *ws.PostApi,
	ch chan common.Cmd2Value, loginUserID string, platformID int32, dataDir, encryptionKey string,
	friend *friend.Friend, group *group.Group, user *user.User,
	objectStorage common2.ObjectStorage, conversationListener open_im_sdk_callback.OnConversationListener,
	msgListener open_im_sdk_callback.OnAdvancedMsgListener, organization *organization.Organization, signaling *signaling.LiveSignaling,
	workMoments *workMoments.WorkMoments, business *business.Business, cache *cache.Cache, full *full.Full, id2MinSeq map[string]uint32, isExternalExtensions bool) *Conversation {
	n := &Conversation{Ws: ws, db: db, p: p, recvCH: ch, loginUserID: loginUserID, platformID: platformID,
		DataDir: dataDir, friend: friend, group: group, user: user, ObjectStorage: objectStorage,
		signaling: signaling, organization: organization, workMoments: workMoments,
		full: full, id2MinSeq: id2MinSeq, encryptionKey: encryptionKey, business: business, IsExternalExtensions: isExternalExtensions}
	n.SetMsgListener(msgListener)
	n.SetConversationListener(conversationListener)
	n.cache = cache
	return n
}

func (c *Conversation) GetCh() chan common.Cmd2Value {
	return c.recvCH
}
func (c *Conversation) doMsgNew(c2v common.Cmd2Value) {
	operationID := c2v.Value.(sdk_struct.CmdNewMsgComeToConversation).OperationID
	allMsg := c2v.Value.(sdk_struct.CmdNewMsgComeToConversation).MsgList
	syncFlag := c2v.Value.(sdk_struct.CmdNewMsgComeToConversation).SyncFlag
	if c.msgListener == nil || c.ConversationListener == nil {
		for _, v := range allMsg {
			if v.ContentType > constant.SignalingNotificationBegin && v.ContentType < constant.SignalingNotificationEnd {
				log.Info(operationID, "signaling DoNotification ", v, "signaling:", c.signaling)
				c.signaling.DoNotification(v, c.GetCh(), operationID)
			} else {
				log.Info(operationID, "listener is nil, do nothing ", v)
			}
		}
	}
	if c.msgListener == nil {
		log.Error(operationID, "not set c MsgListenerList")
		return
	}
	if c.ConversationListener == nil {
		log.Error(operationID, "not set c ConversationListener")
		return
	}
	if syncFlag == constant.MsgSyncBegin {
		c.ConversationListener.OnSyncServerStart()
	}
	if syncFlag == constant.MsgSyncFailed {
		c.ConversationListener.OnSyncServerFailed()
	}

	var isTriggerUnReadCount bool
	var insertMsg, updateMsg []*model_struct.LocalChatLog
	var exceptionMsg []*model_struct.LocalErrChatLog
	var unreadMessages []*model_struct.LocalConversationUnreadMessage
	var newMessages, msgReadList, groupMsgReadList, msgRevokeList, newMsgRevokeList, reactionMsgModifierList, reactionMsgDeleterList sdk_struct.NewMsgList
	var isUnreadCount, isConversationUpdate, isHistory, isNotPrivate, isSenderConversationUpdate, isSenderNotificationPush bool
	conversationChangedSet := make(map[string]*model_struct.LocalConversation)
	newConversationSet := make(map[string]*model_struct.LocalConversation)
	conversationSet := make(map[string]*model_struct.LocalConversation)
	phConversationChangedSet := make(map[string]*model_struct.LocalConversation)
	phNewConversationSet := make(map[string]*model_struct.LocalConversation)
	log.Info(operationID, "do Msg come here, len: ", len(allMsg), len(c.GetCh()))
	b := time.Now()

	for _, v := range allMsg {
		log.Info(operationID, "do Msg come here, msg detail ", v.RecvID, v.SendID, v.ClientMsgID, v.ServerMsgID, v.Seq, c.loginUserID)
		isHistory = utils.GetSwitchFromOptions(v.Options, constant.IsHistory)
		isUnreadCount = utils.GetSwitchFromOptions(v.Options, constant.IsUnreadCount)
		isConversationUpdate = utils.GetSwitchFromOptions(v.Options, constant.IsConversationUpdate)
		isNotPrivate = utils.GetSwitchFromOptions(v.Options, constant.IsNotPrivate)
		isSenderConversationUpdate = utils.GetSwitchFromOptions(v.Options, constant.IsSenderConversationUpdate)
		isSenderNotificationPush = utils.GetSwitchFromOptions(v.Options, constant.IsSenderNotificationPush)
		msg := new(sdk_struct.MsgStruct)
		copier.Copy(msg, v)
		if v.OfflinePushInfo != nil {
			msg.OfflinePush = *v.OfflinePushInfo
		}
		log.Info(operationID, "after copy msg result is", msg.OfflinePush)
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
			log.Error(operationID, "Parsing data error:", err.Error(), *msg, "type: ", msg.ContentType)
			continue
		}
		if !isSenderNotificationPush {
			msg.AttachedInfoElem.NotSenderNotificationPush = true
			msg.AttachedInfo = utils.StructToJsonString(msg.AttachedInfoElem)
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
		case v.ContentType == constant.MsgDeleteNotification:
			c.full.SuperGroup.DoNotification(v, c.GetCh())
		case v.ContentType == constant.SuperGroupUpdateNotification:
			c.full.SuperGroup.DoNotification(v, c.GetCh())
			continue
		case v.ContentType == constant.ConversationUnreadNotification:
			var unreadArgs server_api_params.ConversationUpdateTips
			_ = proto.Unmarshal(tips.Detail, &unreadArgs)
			log.Debug(operationID, "ConversationUnreadNotification come here", unreadArgs.String())
			for _, v := range unreadArgs.ConversationIDList {
				c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: v, Action: constant.UnreadCountSetZero}})
				c.db.DeleteConversationUnreadMessageList(v, unreadArgs.UpdateUnreadCountTime)
			}
			c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.ConChange, Args: unreadArgs.ConversationIDList}})
			continue
		case v.ContentType == constant.BusinessNotification:
			log.NewInfo(operationID, utils.GetSelfFuncName(), "recv businessNotification", tips.JsonDetail)
			c.business.DoNotification(tips.JsonDetail, operationID)
			continue
		}

		switch v.SessionType {
		case constant.SingleChatType:
			if v.ContentType > constant.FriendNotificationBegin && v.ContentType < constant.FriendNotificationEnd {
				c.friend.DoNotification(v, c.GetCh())
				log.Info(operationID, "DoFriendMsg SingleChatType", v)
			} else if v.ContentType > constant.UserNotificationBegin && v.ContentType < constant.UserNotificationEnd {
				log.Info(operationID, "DoFriendMsg  DoUserMsg SingleChatType", v)
				c.user.DoNotification(v)
				//	c.friend.DoNotification(v, c.GetCh())
			} else if v.ContentType == constant.GroupApplicationRejectedNotification ||
				v.ContentType == constant.GroupApplicationAcceptedNotification ||
				v.ContentType == constant.JoinGroupApplicationNotification {
				log.Info(operationID, "DoGroupMsg SingleChatType", v)
				c.group.DoNotification(v, c.GetCh())
			} else if v.ContentType > constant.SignalingNotificationBegin && v.ContentType < constant.SignalingNotificationEnd {
				log.Info(operationID, "signaling DoNotification ", v)
				c.signaling.DoNotification(v, c.GetCh(), operationID)
				continue
			} else if v.ContentType == constant.OrganizationChangedNotification {
				log.Info(operationID, "Organization Changed Notification ")
				c.organization.DoNotification(v, c.GetCh(), operationID)
			} else if v.ContentType == constant.WorkMomentNotification {
				log.Info(operationID, "WorkMoment New Notification")
				c.workMoments.DoNotification(tips.JsonDetail, operationID)
			}
		case constant.GroupChatType, constant.SuperGroupChatType:
			if v.ContentType > constant.GroupNotificationBegin && v.ContentType < constant.GroupNotificationEnd {
				c.group.DoNotification(v, c.GetCh())
				log.Info(operationID, "DoGroupMsg SingleChatType", v)
			} else if v.ContentType > constant.SignalingNotificationBegin && v.ContentType < constant.SignalingNotificationEnd {
				log.Info(operationID, "signaling DoNotification ", v)
				c.signaling.DoNotification(v, c.GetCh(), operationID)
				continue
			}
		}
		if v.SendID == c.loginUserID { //seq
			// Messages sent by myself  //if  sent through  this terminal
			m, err := c.db.GetMessageController(msg)
			if err == nil {
				log.Info(operationID, "have message", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, *msg)
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
				lc := model_struct.LocalConversation{
					ConversationType:  v.SessionType,
					LatestMsg:         utils.StructToJsonString(msg),
					LatestMsgSendTime: msg.SendTime,
				}
				switch v.SessionType {
				case constant.SingleChatType:
					lc.ConversationID = utils.GetConversationIDBySessionType(v.RecvID, constant.SingleChatType)
					lc.UserID = v.RecvID
				case constant.GroupChatType:
					lc.GroupID = v.GroupID
					lc.ConversationID = utils.GetConversationIDBySessionType(lc.GroupID, constant.GroupChatType)
				case constant.SuperGroupChatType:
					lc.GroupID = v.GroupID
					lc.ConversationID = utils.GetConversationIDBySessionType(lc.GroupID, constant.SuperGroupChatType)
				}
				if isConversationUpdate {
					if isSenderConversationUpdate {
						log.Debug(operationID, "updateConversation msg", v, lc)
						c.updateConversation(&lc, conversationSet)
					}
					newMessages = append(newMessages, msg)
				} else {
					msg.Status = constant.MsgStatusFiltered
				}
				if isHistory {
					insertMsg = append(insertMsg, c.msgStructToLocalChatLog(msg))
				}
				switch msg.ContentType {
				case constant.Revoke:
					msgRevokeList = append(msgRevokeList, msg)
				case constant.HasReadReceipt:
					msgReadList = append(msgReadList, msg)
				case constant.GroupHasReadReceipt:
					groupMsgReadList = append(groupMsgReadList, msg)
				case constant.AdvancedRevoke:
					newMsgRevokeList = append(newMsgRevokeList, msg)
					newMessages = removeElementInList(newMessages, msg)
				case constant.ReactionMessageModifier:
					reactionMsgModifierList = append(reactionMsgModifierList, msg)
				case constant.ReactionMessageDeleter:
					reactionMsgDeleterList = append(reactionMsgDeleterList, msg)
				default:
				}
			}
		} else { //Sent by others
			log.NewDebug("internal4", v)
			if _, err := c.db.GetMessageController(msg); err != nil { //Deduplication operation
				lc := model_struct.LocalConversation{
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
				case constant.SuperGroupChatType:
					lc.GroupID = v.GroupID
					lc.ConversationID = utils.GetConversationIDBySessionType(lc.GroupID, constant.SuperGroupChatType)
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
					cacheConversation := c.cache.GetConversation(lc.ConversationID)
					if msg.SendTime > cacheConversation.UpdateUnreadCountTime {
						isTriggerUnReadCount = true
						lc.UnreadCount = 1
						tempUnreadMessages := model_struct.LocalConversationUnreadMessage{ConversationID: lc.ConversationID, ClientMsgID: msg.ClientMsgID, SendTime: msg.SendTime}
						unreadMessages = append(unreadMessages, &tempUnreadMessages)
					}
				}
				if isConversationUpdate {
					c.updateConversation(&lc, conversationSet)
					newMessages = append(newMessages, msg)
				} else {
					msg.Status = constant.MsgStatusFiltered
				}
				if isHistory {
					log.Debug(operationID, "trigger msg is ", msg.SenderNickname, msg.SenderFaceURL)
					insertMsg = append(insertMsg, c.msgStructToLocalChatLog(msg))
				}
				switch msg.ContentType {
				case constant.Revoke:
					msgRevokeList = append(msgRevokeList, msg)
				case constant.HasReadReceipt:
					msgReadList = append(msgReadList, msg)
				case constant.GroupHasReadReceipt:
					groupMsgReadList = append(groupMsgReadList, msg)
				case constant.Typing:
					newMessages = append(newMessages, msg)
				case constant.CustomMsgOnlineOnly:
					newMessages = append(newMessages, msg)
				case constant.CustomMsgNotTriggerConversation:
					newMessages = append(newMessages, msg)
				case constant.OANotification:
					if !isConversationUpdate {
						newMessages = append(newMessages, msg)
					}
				case constant.AdvancedRevoke:
					newMsgRevokeList = append(newMsgRevokeList, msg)
					newMessages = removeElementInList(newMessages, msg)
				case constant.ReactionMessageModifier:
					reactionMsgModifierList = append(reactionMsgModifierList, msg)
				case constant.ReactionMessageDeleter:
					reactionMsgDeleterList = append(reactionMsgDeleterList, msg)
				default:
				}

			} else {
				exceptionMsg = append(exceptionMsg, c.msgStructToLocalErrChatLog(msg))
				log.Warn(operationID, "Deduplication operation ", *c.msgStructToLocalErrChatLog(msg))
			}
		}
	}
	log.NewDebug("internal", newMsgRevokeList)

	b1 := utils.GetCurrentTimestampByMill()
	log.Info(operationID, "generate conversation map is :", conversationSet)
	log.Debug(operationID, "before insert msg cost time : ", time.Since(b))

	list, err := c.db.GetAllConversationListDB()
	if err != nil {
		log.Error(operationID, "GetAllConversationListDB", "error", err.Error())
	}
	m := make(map[string]*model_struct.LocalConversation)
	listToMap(list, m)
	log.Debug(operationID, "listToMap: ", list, conversationSet)
	c.diff(m, conversationSet, conversationChangedSet, newConversationSet)
	log.Info(operationID, "trigger map is :", "newConversations", newConversationSet, "changedConversations", conversationChangedSet)
	b2 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "listToMap diff, cost time : ", b2-b1)

	//seq sync message update
	err5 := c.db.BatchUpdateMessageList(updateMsg)
	if err5 != nil {
		log.Error(operationID, "sync seq normal message err  :", err5.Error())
	}
	b3 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "BatchUpdateMessageList, cost time : ", b3-b2)

	//Normal message storage
	err1 := c.db.BatchInsertMessageListController(insertMsg)
	if err1 != nil {
		log.Error(operationID, "insert GetMessage detail err:", err1.Error(), len(insertMsg))
		for _, v := range insertMsg {
			e := c.db.InsertMessageController(v)
			if e != nil {
				errChatLog := &model_struct.LocalErrChatLog{}
				copier.Copy(errChatLog, v)
				exceptionMsg = append(exceptionMsg, errChatLog)
				log.Warn(operationID, "InsertMessage operation ", "chat err log: ", errChatLog, "chat log: ", v, e.Error())
			}
		}
	}
	b4 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "BatchInsertMessageListController, cost time : ", b4-b3)

	//Exception message storage
	for _, v := range exceptionMsg {
		log.Warn(operationID, "exceptionMsg show: ", *v)
	}

	err2 := c.db.BatchInsertExceptionMsgController(exceptionMsg)
	if err2 != nil {
		log.Error(operationID, "insert err message err  :", err2.Error())

	}
	hList, _ := c.db.GetHiddenConversationList()
	for _, v := range hList {
		if nc, ok := newConversationSet[v.ConversationID]; ok {
			phConversationChangedSet[v.ConversationID] = nc
			nc.RecvMsgOpt = v.RecvMsgOpt
			nc.GroupAtType = v.GroupAtType
			nc.IsPinned = v.IsPinned
			nc.IsPrivateChat = v.IsPrivateChat
			if nc.IsPrivateChat {
				nc.BurnDuration = v.BurnDuration
			}
			nc.IsNotInGroup = v.IsNotInGroup
			nc.AttachedInfo = v.AttachedInfo
			nc.Ex = v.Ex
		}
	}
	b5 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "GetHiddenConversationList, cost time : ", b5-b4)

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
	b6 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "BatchUpdateConversationList, cost time : ", b6-b5)
	//New conversation storage
	err4 := c.db.BatchInsertConversationList(mapConversationToList(phNewConversationSet))
	if err4 != nil {
		log.Error(operationID, "insert new conversation err:", err4.Error())
	}
	b7 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "BatchInsertConversationList, cost time : ", b7-b6)
	unreadMessageErr := c.db.BatchInsertConversationUnreadMessageList(unreadMessages)
	if unreadMessageErr != nil {
		log.Error(operationID, "insert BatchInsertConversationUnreadMessageList err:", unreadMessageErr.Error())
	}
	c.doMsgReadState(msgReadList)
	b8 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "doMsgReadState  cost time : ", b8-b7)

	c.DoGroupMsgReadState(groupMsgReadList)
	b9 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "DoGroupMsgReadState  cost time : ", b9-b8, "len: ", len(groupMsgReadList))

	c.revokeMessage(msgRevokeList)
	b10 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "revokeMessage  cost time : ", b10-b9)
	if c.batchMsgListener != nil {
		c.batchNewMessages(newMessages)
		b11 := utils.GetCurrentTimestampByMill()
		log.Debug(operationID, "batchNewMessages  cost time : ", b11-b10)
	} else {
		c.newMessage(newMessages)
		b12 := utils.GetCurrentTimestampByMill()
		log.Debug(operationID, "newMessage  cost time : ", b12-b10)
	}
	c.newRevokeMessage(newMsgRevokeList)
	c.doReactionMsgModifier(reactionMsgModifierList)
	c.doReactionMsgDeleter(reactionMsgDeleterList)
	//log.Info(operationID, "trigger map is :", newConversationSet, conversationChangedSet)
	if len(newConversationSet) > 0 {
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.NewConDirect, Args: utils.StructToJsonString(mapConversationToList(newConversationSet))}})

	}
	if len(conversationChangedSet) > 0 {
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.ConChangeDirect, Args: utils.StructToJsonString(mapConversationToList(conversationChangedSet))}})
	}

	if isTriggerUnReadCount {
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.TotalUnreadMessageChanged, Args: ""}})
	}
	if syncFlag == constant.MsgSyncEnd {
		c.ConversationListener.OnSyncServerFinish()
	}
	log.Debug(operationID, "insert msg, total cost time: ", time.Since(b), "len:  ", len(allMsg))
}
func (c *Conversation) doSuperGroupMsgNew(c2v common.Cmd2Value) {
	operationID := c2v.Value.(sdk_struct.CmdNewMsgComeToConversation).OperationID
	allMsg := c2v.Value.(sdk_struct.CmdNewMsgComeToConversation).MsgList
	syncFlag := c2v.Value.(sdk_struct.CmdNewMsgComeToConversation).SyncFlag
	log.Info(operationID, utils.GetSelfFuncName(), " args ", "syncFlag ", syncFlag)
	if c.msgListener == nil {
		log.Error(operationID, "not set c MsgListenerList")
		return
	}
	if c.ConversationListener == nil {
		log.Error(operationID, "not set c ConversationListener")
		return
	}
	if syncFlag == constant.MsgSyncBegin {
		log.Info(operationID, "OnSyncServerStart() ")
		c.ConversationListener.OnSyncServerStart()
	}
	if syncFlag == constant.MsgSyncFailed {
		c.ConversationListener.OnSyncServerFailed()
	}
	var isTriggerUnReadCount bool
	var insertMsg, updateMsg, specialUpdateMsg []*model_struct.LocalChatLog
	var exceptionMsg []*model_struct.LocalErrChatLog
	var unreadMessages []*model_struct.LocalConversationUnreadMessage
	var newMessages, msgReadList, groupMsgReadList, msgRevokeList, newMsgRevokeList, reactionMsgModifierList, reactionMsgDeleterList sdk_struct.NewMsgList
	var isUnreadCount, isConversationUpdate, isHistory, isNotPrivate, isSenderConversationUpdate, isSenderNotificationPush bool
	conversationChangedSet := make(map[string]*model_struct.LocalConversation)
	newConversationSet := make(map[string]*model_struct.LocalConversation)
	conversationSet := make(map[string]*model_struct.LocalConversation)
	phConversationChangedSet := make(map[string]*model_struct.LocalConversation)
	phNewConversationSet := make(map[string]*model_struct.LocalConversation)
	log.Info(operationID, "do Msg come here, len: ", len(allMsg))
	b := utils.GetCurrentTimestampByMill()

	for _, v := range allMsg {
		log.Info(operationID, "do Msg come here, msg detail ", *v, c.loginUserID)
		isHistory = utils.GetSwitchFromOptions(v.Options, constant.IsHistory)
		isUnreadCount = utils.GetSwitchFromOptions(v.Options, constant.IsUnreadCount)
		isConversationUpdate = utils.GetSwitchFromOptions(v.Options, constant.IsConversationUpdate)
		isNotPrivate = utils.GetSwitchFromOptions(v.Options, constant.IsNotPrivate)
		isSenderConversationUpdate = utils.GetSwitchFromOptions(v.Options, constant.IsSenderConversationUpdate)
		isSenderNotificationPush = utils.GetSwitchFromOptions(v.Options, constant.IsSenderNotificationPush)
		msg := new(sdk_struct.MsgStruct)
		copier.Copy(msg, v)
		if v.OfflinePushInfo != nil {
			msg.OfflinePush = *v.OfflinePushInfo
		}
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
			log.Error(operationID, "Parsing data error:", err.Error(), *msg, "type: ", msg.ContentType)
			continue
		}
		if !isSenderNotificationPush {
			msg.AttachedInfoElem.NotSenderNotificationPush = true
			msg.AttachedInfo = utils.StructToJsonString(msg.AttachedInfoElem)
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
		case v.ContentType == constant.SuperGroupUpdateNotification:
			c.full.SuperGroup.DoNotification(v, c.GetCh())
		case v.ContentType == constant.ConversationUnreadNotification:
			var unreadArgs server_api_params.ConversationUpdateTips
			_ = proto.Unmarshal(tips.Detail, &unreadArgs)
			for _, v := range unreadArgs.ConversationIDList {
				c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: v, Action: constant.UnreadCountSetZero}})
				c.db.DeleteConversationUnreadMessageList(v, unreadArgs.UpdateUnreadCountTime)
			}
			c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.ConChange, Args: unreadArgs.ConversationIDList}})
			continue
		case v.ContentType == constant.BusinessNotification:
			log.NewInfo(operationID, utils.GetSelfFuncName(), "recv businessNotification", tips.JsonDetail)
			c.business.DoNotification(tips.JsonDetail, operationID)
			continue
		}
		switch v.SessionType {
		case constant.SingleChatType:
			if v.ContentType > constant.FriendNotificationBegin && v.ContentType < constant.FriendNotificationEnd {
				c.friend.DoNotification(v, c.GetCh())
				log.Info(operationID, "DoFriendMsg SingleChatType", v)
			} else if v.ContentType > constant.UserNotificationBegin && v.ContentType < constant.UserNotificationEnd {
				log.Info(operationID, "DoFriendMsg  DoUserMsg SingleChatType", v)
				c.user.DoNotification(v)
				//c.friend.DoNotification(v, c.GetCh())
			} else if v.ContentType == constant.GroupApplicationRejectedNotification ||
				v.ContentType == constant.GroupApplicationAcceptedNotification ||
				v.ContentType == constant.JoinGroupApplicationNotification {
				log.Info(operationID, "DoGroupMsg SingleChatType", v)
				c.group.DoNotification(v, c.GetCh())
			} else if v.ContentType > constant.SignalingNotificationBegin && v.ContentType < constant.SignalingNotificationEnd {
				log.Info(operationID, "signaling DoNotification ", v)
				c.signaling.DoNotification(v, c.GetCh(), operationID)
				continue
			} else if v.ContentType == constant.OrganizationChangedNotification {
				log.Info(operationID, "Organization Changed Notification ")
				c.organization.DoNotification(v, c.GetCh(), operationID)
			} else if v.ContentType == constant.WorkMomentNotification {
				log.Info(operationID, "WorkMoment New Notification")
				c.workMoments.DoNotification(tips.JsonDetail, operationID)
			}
		case constant.GroupChatType, constant.SuperGroupChatType:
			if v.ContentType > constant.GroupNotificationBegin && v.ContentType < constant.GroupNotificationEnd {
				c.group.DoNotification(v, c.GetCh())
				log.Info(operationID, "DoGroupMsg SingleChatType", v)
			} else if v.ContentType > constant.SignalingNotificationBegin && v.ContentType < constant.SignalingNotificationEnd {
				log.Info(operationID, "signaling DoNotification ", v)
				c.signaling.DoNotification(v, c.GetCh(), operationID)
				continue
			}
		}
		if v.SendID == c.loginUserID { //seq
			// Messages sent by myself  //if  sent through  this terminal
			m, err := c.db.GetMessageController(msg)
			if err == nil {
				log.Info(operationID, "have message", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, *msg)
				if m.Seq == 0 {
					if m.CreateTime == 0 {
						specialUpdateMsg = append(specialUpdateMsg, c.msgStructToLocalChatLog(msg))
					} else {
						if !isConversationUpdate {
							msg.Status = constant.MsgStatusFiltered
						}
						updateMsg = append(updateMsg, c.msgStructToLocalChatLog(msg))
					}
				} else {
					exceptionMsg = append(exceptionMsg, c.msgStructToLocalErrChatLog(msg))
				}
			} else { //      send through  other terminal
				log.Info(operationID, "sync message", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, *msg)
				lc := model_struct.LocalConversation{
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
				case constant.SuperGroupChatType:
					lc.GroupID = v.GroupID
					lc.ConversationID = utils.GetConversationIDBySessionType(lc.GroupID, constant.SuperGroupChatType)
					//faceUrl, name, err := u.getGroupNameAndFaceUrlByUid(c.GroupID)
					//if err != nil {
					//	utils.sdkLog("getGroupNameAndFaceUrlByUid err:", err)
					//} else {
					//	c.ShowName = name
					//	c.FaceURL = faceUrl
					//}

				}
				if isConversationUpdate {
					if isSenderConversationUpdate {
						log.Debug(operationID, "updateConversation msg", v, lc)
						c.updateConversation(&lc, conversationSet)
					}
					newMessages = append(newMessages, msg)
				} else {
					msg.Status = constant.MsgStatusFiltered
				}
				if isHistory {
					insertMsg = append(insertMsg, c.msgStructToLocalChatLog(msg))
				}
				switch msg.ContentType {
				case constant.Revoke:
					msgRevokeList = append(msgRevokeList, msg)
				case constant.HasReadReceipt:
					msgReadList = append(msgReadList, msg)
				case constant.GroupHasReadReceipt:
					groupMsgReadList = append(groupMsgReadList, msg)
				case constant.AdvancedRevoke:
					newMsgRevokeList = append(newMsgRevokeList, msg)
					newMessages = removeElementInList(newMessages, msg)
				case constant.ReactionMessageModifier:
					reactionMsgModifierList = append(reactionMsgModifierList, msg)
				case constant.ReactionMessageDeleter:
					reactionMsgDeleterList = append(reactionMsgDeleterList, msg)
				default:
				}
			}
		} else { //Sent by others
			if oldMessage, err := c.db.GetMessageController(msg); err != nil { //Deduplication operation
				lc := model_struct.LocalConversation{
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
				case constant.SuperGroupChatType:
					lc.GroupID = v.GroupID
					lc.ConversationID = utils.GetConversationIDBySessionType(lc.GroupID, constant.SuperGroupChatType)
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
					cacheConversation := c.cache.GetConversation(lc.ConversationID)
					if msg.SendTime > cacheConversation.UpdateUnreadCountTime {
						isTriggerUnReadCount = true
						lc.UnreadCount = 1
						tempUnreadMessages := model_struct.LocalConversationUnreadMessage{ConversationID: lc.ConversationID, ClientMsgID: msg.ClientMsgID, SendTime: msg.SendTime}
						unreadMessages = append(unreadMessages, &tempUnreadMessages)
					}
				}
				if isConversationUpdate {
					c.updateConversation(&lc, conversationSet)
					newMessages = append(newMessages, msg)
				} else {
					msg.Status = constant.MsgStatusFiltered
				}
				if isHistory {
					log.Debug(operationID, "trigger msg is ", msg.SenderNickname, msg.SenderFaceURL)
					insertMsg = append(insertMsg, c.msgStructToLocalChatLog(msg))
				}
				switch msg.ContentType {
				case constant.Revoke:
					msgRevokeList = append(msgRevokeList, msg)
				case constant.HasReadReceipt:
					msgReadList = append(msgReadList, msg)
				case constant.GroupHasReadReceipt:
					groupMsgReadList = append(groupMsgReadList, msg)
				case constant.CustomMsgOnlineOnly:
					newMessages = append(newMessages, msg)
				case constant.CustomMsgNotTriggerConversation:
					newMessages = append(newMessages, msg)
				case constant.OANotification:
					if !isConversationUpdate {
						newMessages = append(newMessages, msg)
					}
				case constant.Typing:
					newMessages = append(newMessages, msg)
				case constant.AdvancedRevoke:
					newMsgRevokeList = append(newMsgRevokeList, msg)
					newMessages = removeElementInList(newMessages, msg)
				case constant.ReactionMessageModifier:
					reactionMsgModifierList = append(reactionMsgModifierList, msg)
				case constant.ReactionMessageDeleter:
					reactionMsgDeleterList = append(reactionMsgDeleterList, msg)
				default:
				}

			} else {
				if oldMessage.Seq == 0 {
					specialUpdateMsg = append(specialUpdateMsg, c.msgStructToLocalChatLog(msg))
				} else {
					exceptionMsg = append(exceptionMsg, c.msgStructToLocalErrChatLog(msg))
					log.Warn(operationID, "Deduplication operation ", *c.msgStructToLocalErrChatLog(msg))
				}
			}
		}
	}
	b1 := utils.GetCurrentTimestampByMill()
	log.Info(operationID, "generate conversation map is :", conversationSet)
	log.Debug(operationID, "before insert msg cost time : ", b1-b)

	list, err := c.db.GetAllConversationListDB()
	if err != nil {
		log.Error(operationID, "GetAllConversationListDB", "error", err.Error())
	}
	m := make(map[string]*model_struct.LocalConversation)
	listToMap(list, m)
	log.Debug(operationID, "listToMap: ", list, conversationSet)
	c.diff(m, conversationSet, conversationChangedSet, newConversationSet)
	log.Info(operationID, "trigger map is :", "newConversations", newConversationSet, "changedConversations", conversationChangedSet)
	b2 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "listToMap diff, cost time : ", b2-b1)
	//update message
	err6 := c.db.BatchSpecialUpdateMessageList(specialUpdateMsg)
	if err6 != nil {
		log.Error(operationID, "sync seq normal message err  :", err6.Error())
	}
	//seq sync message update
	err5 := c.db.BatchUpdateMessageList(updateMsg)
	if err5 != nil {
		log.Error(operationID, "sync seq normal message err  :", err5.Error())
	}
	b3 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "BatchUpdateMessageList, cost time : ", b3-b2)

	//Normal message storage
	err1 := c.db.BatchInsertMessageListController(insertMsg)
	if err1 != nil {
		log.Error(operationID, "insert GetMessage detail err:", err1.Error(), len(insertMsg))
		for _, v := range insertMsg {
			e := c.db.InsertMessageController(v)
			if e != nil {
				errChatLog := &model_struct.LocalErrChatLog{}
				copier.Copy(errChatLog, v)
				exceptionMsg = append(exceptionMsg, errChatLog)
				log.Warn(operationID, "InsertMessage operation ", "chat err log: ", errChatLog, "chat log: ", v, e.Error())
			}
		}
	}
	b4 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "BatchInsertMessageListController, cost time : ", b4-b3)

	//Exception message storage
	for _, v := range exceptionMsg {
		log.Warn(operationID, "exceptionMsg show: ", *v)
	}

	err2 := c.db.BatchInsertExceptionMsgController(exceptionMsg)
	if err2 != nil {
		log.Error(operationID, "insert err message err  :", err2.Error())

	}
	hList, _ := c.db.GetHiddenConversationList()
	for _, v := range hList {
		if nc, ok := newConversationSet[v.ConversationID]; ok {
			phConversationChangedSet[v.ConversationID] = nc
			nc.RecvMsgOpt = v.RecvMsgOpt
			nc.GroupAtType = v.GroupAtType
			nc.IsPinned = v.IsPinned
			nc.IsPrivateChat = v.IsPrivateChat
			if nc.IsPrivateChat {
				nc.BurnDuration = v.BurnDuration
			}
			nc.IsNotInGroup = v.IsNotInGroup
			nc.AttachedInfo = v.AttachedInfo
			nc.Ex = v.Ex
		}
	}
	b5 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "GetHiddenConversationList, cost time : ", b5-b4)

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
	b6 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "BatchUpdateConversationList, cost time : ", b6-b5)
	//New conversation storage
	err4 := c.db.BatchInsertConversationList(mapConversationToList(phNewConversationSet))
	if err4 != nil {
		log.Error(operationID, "insert new conversation err:", err4.Error())
	}
	b7 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "BatchInsertConversationList, cost time : ", b7-b6)
	unreadMessageErr := c.db.BatchInsertConversationUnreadMessageList(unreadMessages)
	if unreadMessageErr != nil {
		log.Error(operationID, "insert BatchInsertConversationUnreadMessageList err:", unreadMessageErr.Error())
	}
	c.doMsgReadState(msgReadList)
	b8 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "doMsgReadState  cost time : ", b8-b7)

	c.DoGroupMsgReadState(groupMsgReadList)
	b9 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "DoGroupMsgReadState  cost time : ", b9-b8, "len: ", len(groupMsgReadList))

	c.revokeMessage(msgRevokeList)
	b10 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "revokeMessage  cost time : ", b10-b9)
	if c.batchMsgListener != nil {
		c.batchNewMessages(newMessages)
		b11 := utils.GetCurrentTimestampByMill()
		log.Debug(operationID, "batchNewMessages  cost time : ", b11-b10)
	} else {
		c.newMessage(newMessages)
		b12 := utils.GetCurrentTimestampByMill()
		log.Debug(operationID, "newMessage  cost time : ", b12-b10)
	}
	c.newRevokeMessage(newMsgRevokeList)
	c.doReactionMsgModifier(reactionMsgModifierList)
	c.doReactionMsgDeleter(reactionMsgDeleterList)
	//log.Info(operationID, "trigger map is :", newConversationSet, conversationChangedSet)
	if len(newConversationSet) > 0 {
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{"", constant.NewConDirect, utils.StructToJsonString(mapConversationToList(newConversationSet))}})

	}
	if len(conversationChangedSet) > 0 {
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{"", constant.ConChangeDirect, utils.StructToJsonString(mapConversationToList(conversationChangedSet))}})
	}

	if isTriggerUnReadCount {
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{"", constant.TotalUnreadMessageChanged, ""}})
	}
	if syncFlag == constant.MsgSyncEnd {
		log.Info(operationID, "OnSyncServerFinish() ")
		c.ConversationListener.OnSyncServerFinish()
	}
	log.Info(operationID, "insert msg, total cost time: ", utils.GetCurrentTimestampByMill()-b, "len:  ", len(allMsg))
}
func listToMap(list []*model_struct.LocalConversation, m map[string]*model_struct.LocalConversation) {
	for _, v := range list {
		m[v.ConversationID] = v
	}

}
func removeElementInList(a sdk_struct.NewMsgList, e *sdk_struct.MsgStruct) (b sdk_struct.NewMsgList) {
	for i := 0; i < len(a); i++ {
		if a[i] != e {
			b = append(b, a[i])
		}
	}
	return b
}
func (c *Conversation) diff(local, generated, cc, nc map[string]*model_struct.LocalConversation) {
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
func (c *Conversation) genConversationGroupAtType(lc *model_struct.LocalConversation, s *sdk_struct.MsgStruct) {
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
func (c *Conversation) msgStructToLocalChatLog(m *sdk_struct.MsgStruct) *model_struct.LocalChatLog {
	var lc model_struct.LocalChatLog
	copier.Copy(&lc, m)
	if m.SessionType == constant.GroupChatType || m.SessionType == constant.SuperGroupChatType {
		lc.RecvID = m.GroupID
	}
	return &lc
}
func (c *Conversation) msgStructToLocalErrChatLog(m *sdk_struct.MsgStruct) *model_struct.LocalErrChatLog {
	var lc model_struct.LocalErrChatLog
	copier.Copy(&lc, m)
	if m.SessionType == constant.GroupChatType || m.SessionType == constant.SuperGroupChatType {
		lc.RecvID = m.GroupID
	}
	return &lc
}

// deprecated
func (c *Conversation) revokeMessage(msgRevokeList []*sdk_struct.MsgStruct) {
	for _, w := range msgRevokeList {
		if c.msgListener != nil {
			t := new(model_struct.LocalChatLog)
			t.ClientMsgID = w.Content
			t.Status = constant.MsgStatusRevoked
			t.SessionType = w.SessionType
			t.RecvID = w.GroupID
			err := c.db.UpdateMessageController(t)
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

func (c *Conversation) tempCacheChatLog(messageList []*sdk_struct.MsgStruct) {
	var newMessageList []*model_struct.TempCacheLocalChatLog
	copier.Copy(&newMessageList, &messageList)
	err1 := c.db.BatchInsertTempCacheMessageList(newMessageList)
	if err1 != nil {
		log.Error("", "BatchInsertTempCacheMessageList detail err:", err1.Error(), len(newMessageList))
		for _, v := range newMessageList {
			e := c.db.InsertTempCacheMessage(v)
			if e != nil {
				log.Warn("", "InsertTempCacheMessage operation ", "chat err log: ", *v, e.Error())
			}
		}
	}
}
func (c *Conversation) newRevokeMessage(msgRevokeList []*sdk_struct.MsgStruct) {
	var failedRevokeMessageList []*sdk_struct.MsgStruct
	var superGroupIDList []string
	var revokeMessageRevoked []*sdk_struct.MessageRevoked
	var superGroupRevokeMessageRevoked []*sdk_struct.MessageRevoked
	log.NewDebug("revoke msg", msgRevokeList)
	for _, w := range msgRevokeList {
		log.NewDebug("msg revoke", w)
		var msg sdk_struct.MessageRevoked
		err := json.Unmarshal([]byte(w.Content), &msg)
		if err != nil {
			log.Error("internal", "unmarshal failed, err : ", err.Error(), *w)
			continue
		}
		t := new(model_struct.LocalChatLog)
		t.ClientMsgID = msg.ClientMsgID
		t.Status = constant.MsgStatusRevoked
		t.SessionType = msg.SessionType
		t.RecvID = w.GroupID
		err = c.db.UpdateMessageController(t)
		if err != nil {
			log.Error("internal", "setLocalMessageStatus revokeMessage err:", err.Error(), "msg", w)
			failedRevokeMessageList = append(failedRevokeMessageList, w)
			switch msg.SessionType {
			case constant.SuperGroupChatType:
				err := c.db.InsertMessageController(t)
				if err != nil {
					log.Error("internal", "InsertMessageController preDefine Message err:", err.Error(), "msg", *t)
				}
			}
		} else {
			t := new(model_struct.LocalChatLog)
			t.ClientMsgID = w.ClientMsgID
			t.SendTime = msg.SourceMessageSendTime
			t.SessionType = w.SessionType
			t.RecvID = w.GroupID
			err = c.db.UpdateMessageController(t)
			if err != nil {
				log.Error("internal", "setLocalMessageStatus revokeMessage err:", err.Error(), "msg", w)
			}
			log.Info("internal", "v.OnNewRecvMessageRevoked client_msg_id:", msg.ClientMsgID)
			if c.msgListener != nil {
				c.msgListener.OnNewRecvMessageRevoked(w.Content)
			} else {
				log.Error("internal", "set msgListener is err:")
			}
			if msg.SessionType != constant.SuperGroupChatType {
				revokeMessageRevoked = append(revokeMessageRevoked, &msg)
			} else {
				if !utils.IsContain(w.RecvID, superGroupIDList) {
					superGroupIDList = append(superGroupIDList, w.GroupID)
				}
				superGroupRevokeMessageRevoked = append(superGroupRevokeMessageRevoked, &msg)
			}
		}
	}
	log.NewDebug("internal, quoteRevoke Info", superGroupIDList, len(revokeMessageRevoked), len(superGroupRevokeMessageRevoked))
	if len(revokeMessageRevoked) > 0 {
		msgList, err := c.db.SearchAllMessageByContentType(constant.Quote)
		if err != nil {
			log.NewError("internal", "SearchMessageIDsByContentType failed", err.Error())
		}
		for _, v := range msgList {
			c.QuoteMsgRevokeHandle(v, revokeMessageRevoked)
		}
	}
	for _, superGroupID := range superGroupIDList {
		msgList, err := c.db.SuperGroupSearchAllMessageByContentType(superGroupID, constant.Quote)
		if err != nil {
			log.NewError("internal", "SuperGroupSearchMessageByContentTypeNotOffset failed", superGroupID, err.Error())
		}
		for _, v := range msgList {
			c.QuoteMsgRevokeHandle(v, superGroupRevokeMessageRevoked)
		}
	}
	if len(failedRevokeMessageList) > 0 {
		//c.tempCacheChatLog(failedRevokeMessageList)
	}
}
func (c *Conversation) DoMsgReaction(msgReactionList []*sdk_struct.MsgStruct) {

	//for _, v := range msgReactionList {
	//	var msg sdk_struct.MessageReaction
	//	err := json.Unmarshal([]byte(v.Content), &msg)
	//	if err != nil {
	//		log.Error("internal", "unmarshal failed, err : ", err.Error(), *v)
	//		continue
	//	}
	//	s := sdk_struct.MsgStruct{GroupID: msg.GroupID, ClientMsgID: msg.ClientMsgID, SessionType: msg.SessionType}
	//	message, err := c.db.GetMessageController(&s)
	//	if err != nil {
	//		log.Error("internal", "GetMessageController, err : ", err.Error(), s)
	//		continue
	//	}
	//	t := new(model_struct.LocalChatLog)
	//  attachInfo := sdk_struct.AttachedInfoElem{}
	//	_ = utils.JsonStringToStruct(message.AttachedInfo, &attachInfo)
	//
	//	contain, v := isContainMessageReaction(msg.ReactionType, attachInfo.MessageReactionElem)
	//	if contain {
	//		userContain, userReaction := isContainUserReactionElem(msg.UserID, v.UserReactionList)
	//		if userContain {
	//			if !v.CanRepeat && userReaction.Counter > 0 {
	//				// to do nothing
	//				continue
	//			} else {
	//				userReaction.Counter += msg.Counter
	//				v.Counter += msg.Counter
	//				if v.Counter < 0 {
	//					log.Debug("internal", "after operate all counter  < 0", v.Type, v.Counter, msg.Counter)
	//					v.Counter = 0
	//				}
	//				if userReaction.Counter <= 0 {
	//					log.Debug("internal", "after operate userReaction counter < 0", v.Type, userReaction.Counter, msg.Counter)
	//					v.UserReactionList = DeleteUserReactionElem(v.UserReactionList, c.loginUserID)
	//				}
	//			}
	//		} else {
	//			log.Debug("internal", "attachInfo.MessageReactionElem is nil", msg)
	//			u := new(sdk_struct.UserReactionElem)
	//			u.UserID = msg.UserID
	//			u.Counter = msg.Counter
	//			v.Counter += msg.Counter
	//			if v.Counter < 0 {
	//				log.Debug("internal", "after operate all counter  < 0", v.Type, v.Counter, msg.Counter)
	//				v.Counter = 0
	//			}
	//			if u.Counter <= 0 {
	//				log.Debug("internal", "after operate userReaction counter < 0", v.Type, u.Counter, msg.Counter)
	//				v.UserReactionList = DeleteUserReactionElem(v.UserReactionList, msg.UserID)
	//			}
	//			v.UserReactionList = append(v.UserReactionList, u)
	//
	//		}
	//
	//	} else {
	//		log.Debug("internal", "attachInfo.MessageReactionElem is nil", msg)
	//		t := new(sdk_struct.ReactionElem)
	//		t.Counter = msg.Counter
	//		t.Type = msg.ReactionType
	//		u := new(sdk_struct.UserReactionElem)
	//		u.UserID = msg.UserID
	//		u.Counter = msg.Counter
	//		t.UserReactionList = append(t.UserReactionList, u)
	//		attachInfo.MessageReactionElem = append(attachInfo.MessageReactionElem, t)
	//
	//	}
	//
	//	t.AttachedInfo = utils.StructToJsonString(attachInfo)
	//	t.ClientMsgID = message.ClientMsgID
	//
	//	t.SessionType = message.SessionType
	//	t.RecvID = message.RecvID
	//	err1 := c.db.UpdateMessageController(t)
	//	if err1 != nil {
	//		log.Error("internal", "UpdateMessageController err:", err1, "ClientMsgID", *t, message)
	//	}
	//	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{"", constant.MessageChange, &s}})
	//
	//}
}

func (c *Conversation) doReactionMsgModifier(msgReactionList []*sdk_struct.MsgStruct) {
	for _, msgStruct := range msgReactionList {
		var n server_api_params.ReactionMessageModifierNotification
		err := json.Unmarshal([]byte(msgStruct.Content), &n)
		if err != nil {
			log.Error("internal", "unmarshal failed err:", err.Error(), *msgStruct)
			continue
		}
		switch n.Operation {
		case constant.AddMessageExtensions:
			var reactionExtensionList []*server_api_params.KeyValue
			for _, value := range n.SuccessReactionExtensionList {
				reactionExtensionList = append(reactionExtensionList, value)
			}
			if !(msgStruct.SendID == c.loginUserID && msgStruct.SenderPlatformID == c.platformID) {
				c.msgListener.OnRecvMessageExtensionsAdded(n.ClientMsgID, utils.StructToJsonString(reactionExtensionList))
			}
		case constant.SetMessageExtensions:
			err = c.db.GetAndUpdateMessageReactionExtension(n.ClientMsgID, n.SuccessReactionExtensionList)
			if err != nil {
				log.Error("internal", "GetAndUpdateMessageReactionExtension err:", err.Error())
				continue
			}
			var reactionExtensionList []*server_api_params.KeyValue
			for _, value := range n.SuccessReactionExtensionList {
				reactionExtensionList = append(reactionExtensionList, value)
			}
			if !(msgStruct.SendID == c.loginUserID && msgStruct.SenderPlatformID == c.platformID) {
				c.msgListener.OnRecvMessageExtensionsChanged(n.ClientMsgID, utils.StructToJsonString(reactionExtensionList))
			}

		}
		t := model_struct.LocalChatLog{}
		t.ClientMsgID = n.ClientMsgID
		t.SessionType = n.SessionType
		t.IsExternalExtensions = n.IsExternalExtensions
		t.IsReact = n.IsReact
		t.MsgFirstModifyTime = n.MsgFirstModifyTime
		if n.SessionType == constant.GroupChatType || n.SessionType == constant.SuperGroupChatType {
			t.RecvID = n.SourceID
		}
		err2 := c.db.UpdateMessageController(&t)
		if err2 != nil {
			log.Error("internal", "unmarshal failed err:", err2.Error(), t)
			continue
		}

	}

}
func (c *Conversation) doReactionMsgDeleter(msgReactionList []*sdk_struct.MsgStruct) {
	for _, msgStruct := range msgReactionList {
		var n server_api_params.ReactionMessageDeleteNotification
		err := json.Unmarshal([]byte(msgStruct.Content), &n)
		if err != nil {
			log.Error("internal", "unmarshal failed err:", err.Error(), *msgStruct)
			continue
		}
		err = c.db.DeleteAndUpdateMessageReactionExtension(n.ClientMsgID, n.SuccessReactionExtensionList)
		if err != nil {
			log.Error("internal", "GetAndUpdateMessageReactionExtension err:", err.Error())
			continue
		}
		var deleteKeyList []string
		for _, value := range n.SuccessReactionExtensionList {
			deleteKeyList = append(deleteKeyList, value.TypeKey)
		}
		c.msgListener.OnRecvMessageExtensionsDeleted(n.ClientMsgID, utils.StructToJsonString(deleteKeyList))

	}

}
func (c *Conversation) QuoteMsgRevokeHandle(v *model_struct.LocalChatLog, revokeMsgIDList []*sdk_struct.MessageRevoked) {
	s := sdk_struct.MsgStruct{}
	err := utils.JsonStringToStruct(v.Content, &s.QuoteElem)
	if err != nil {
		log.NewError("internal", "unmarshall failed", s.Content)
	}
	if s.QuoteElem.QuoteMessage == nil {
		return
	}
	ok, revokeMessage := isContainRevokedList(s.QuoteElem.QuoteMessage.ClientMsgID, revokeMsgIDList)
	if !ok {
		return
	}
	s.QuoteElem.QuoteMessage.Content = utils.StructToJsonString(revokeMessage)
	s.QuoteElem.QuoteMessage.ContentType = constant.AdvancedRevoke
	v.Content = utils.StructToJsonString(s.QuoteElem)
	err = c.db.UpdateMessageController(v)
	if err != nil {
		log.NewError("internal", "unmarshall failed", v)
	}
}
func isContainRevokedList(target string, List []*sdk_struct.MessageRevoked) (bool, *sdk_struct.MessageRevoked) {
	for _, element := range List {
		if target == element.ClientMsgID {
			return true, element
		}
	}
	return false, nil
}

func (c *Conversation) DoGroupMsgReadState(groupMsgReadList []*sdk_struct.MsgStruct) {
	var groupMessageReceiptResp []*sdk_struct.MessageReceipt
	var failedMessageList []*sdk_struct.MsgStruct
	userMsgMap := make(map[string]map[string][]string)
	//var temp []*sdk_struct.MessageReceipt
	for _, rd := range groupMsgReadList {
		var list []string
		err := json.Unmarshal([]byte(rd.Content), &list)
		if err != nil {
			log.Error("internal", "unmarshal failed, err : ", err.Error(), rd)
			continue
		}
		if groupMap, ok := userMsgMap[rd.SendID]; ok {
			if oldMsgIDList, ok := groupMap[rd.GroupID]; ok {
				oldMsgIDList = append(oldMsgIDList, list...)
				groupMap[rd.GroupID] = oldMsgIDList
			} else {
				groupMap[rd.GroupID] = list
			}
		} else {
			g := make(map[string][]string)
			g[rd.GroupID] = list
			userMsgMap[rd.SendID] = g
		}

	}
	for userID, m := range userMsgMap {
		for groupID, msgIDList := range m {
			var successMsgIDlist []string
			var failedMsgIDList []string
			newMsgID := utils.RemoveRepeatedStringInList(msgIDList)
			_, sessionType, err := c.getConversationTypeByGroupID(groupID)
			if err != nil {
				log.Error("internal", "GetGroupInfoByGroupID err:", err.Error(), "groupID", groupID)
				continue
			}
			messages, err := c.db.GetMultipleMessageController(newMsgID, groupID, sessionType)
			if err != nil {
				log.Error("internal", "GetMessage err:", err.Error(), "ClientMsgID", newMsgID)
				continue
			}
			msgRt := new(sdk_struct.MessageReceipt)
			msgRt.UserID = userID
			msgRt.GroupID = groupID
			msgRt.SessionType = sessionType
			msgRt.ContentType = constant.GroupHasReadReceipt

			for _, message := range messages {
				t := new(model_struct.LocalChatLog)
				if userID != c.loginUserID {
					attachInfo := sdk_struct.AttachedInfoElem{}
					_ = utils.JsonStringToStruct(message.AttachedInfo, &attachInfo)
					attachInfo.GroupHasReadInfo.HasReadUserIDList = utils.RemoveRepeatedStringInList(append(attachInfo.GroupHasReadInfo.HasReadUserIDList, userID))
					attachInfo.GroupHasReadInfo.HasReadCount = int32(len(attachInfo.GroupHasReadInfo.HasReadUserIDList))
					t.AttachedInfo = utils.StructToJsonString(attachInfo)
				}
				t.ClientMsgID = message.ClientMsgID
				t.IsRead = true
				t.SessionType = message.SessionType
				t.RecvID = message.RecvID
				err1 := c.db.UpdateMessageController(t)
				if err1 != nil {
					log.Error("internal", "setGroupMessageHasReadByMsgID err:", err1, "ClientMsgID", t, message)
					continue
				}
				successMsgIDlist = append(successMsgIDlist, message.ClientMsgID)
			}
			failedMsgIDList = utils.DifferenceSubsetString(newMsgID, successMsgIDlist)
			if len(successMsgIDlist) != 0 {
				msgRt.MsgIDList = successMsgIDlist
				groupMessageReceiptResp = append(groupMessageReceiptResp, msgRt)
			}
			if len(failedMsgIDList) != 0 {
				m := new(sdk_struct.MsgStruct)
				m.ClientMsgID = utils.GetMsgID(userID)
				m.SendID = userID
				m.GroupID = groupID
				m.ContentType = constant.GroupHasReadReceipt
				m.Content = utils.StructToJsonString(failedMsgIDList)
				m.Status = constant.MsgStatusFiltered
				failedMessageList = append(failedMessageList, m)
			}
		}
	}
	if len(groupMessageReceiptResp) > 0 {
		log.Info("internal", "OnRecvGroupReadReceipt: ", utils.StructToJsonString(groupMessageReceiptResp))
		c.msgListener.OnRecvGroupReadReceipt(utils.StructToJsonString(groupMessageReceiptResp))
	}
	if len(failedMessageList) > 0 {
		//c.tempCacheChatLog(failedMessageList)
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
		if c.listenerForService != nil {
			log.Info("internal", "msgListener,OnRecvNewMessage")
			c.listenerForService.OnRecvNewMessage(utils.StructToJsonString(w))
		}
	}
}
func (c *Conversation) batchNewMessages(newMessagesList sdk_struct.NewMsgList) {
	sort.Sort(newMessagesList)
	if c.batchMsgListener != nil {
		log.Info("internal", "batchMsgListener,OnRecvNewMessage")
		if len(newMessagesList) > 0 {
			for _, v := range newMessagesList {
				log.Info("internal", "trigger msg is ", v.SenderFaceURL, v.SenderNickname)
			}
			c.batchMsgListener.OnRecvNewMessages(utils.StructToJsonString(newMessagesList))
		}
	} else {
		log.Warn("internal", "not set batchMsgListener ")
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
	chrsList := make(map[string][]string)
	var conversationID string

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
			msgRt.MsgIDList = msgIdListStatusOK
			messageReceiptResp = append(messageReceiptResp, msgRt)
		}
		if rd.SendID == c.loginUserID {
			conversationID = utils.GetConversationIDBySessionType(rd.RecvID, constant.SingleChatType)
		} else {
			conversationID = utils.GetConversationIDBySessionType(rd.SendID, constant.SingleChatType)
		}
		if v, ok := chrsList[conversationID]; ok {
			chrsList[conversationID] = append(v, msgIdListStatusOK...)
		} else {
			chrsList[conversationID] = msgIdListStatusOK
		}
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.ConversationLatestMsgHasRead, Args: chrsList}})
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
		var list []*model_struct.LocalConversation
		lc := node.Args.(model_struct.LocalConversation)
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
			log.Error("internal", "UpdateColumnsConversation err", err.Error(), node.ConID)
		} else {
			totalUnreadCount, err := c.db.GetTotalUnreadMsgCountDB()
			if err == nil {
				c.ConversationListener.OnTotalUnreadMessageCountChanged(totalUnreadCount)
			} else {
				log.Error("internal", "getTotalUnreadMsgCountModel err", err.Error(), node.ConID)
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
		totalUnreadCount, err := c.db.GetTotalUnreadMsgCountDB()
		if err != nil {
			log.Error("internal", "TotalUnreadMessageChanged database err:", err.Error())
		} else {
			c.ConversationListener.OnTotalUnreadMessageCountChanged(totalUnreadCount)
		}
	case constant.UpdateConFaceUrlAndNickName:
		var lc model_struct.LocalConversation
		st := node.Args.(common.SourceIDAndSessionType)
		switch st.SessionType {
		case constant.SingleChatType:
			lc.UserID = st.SourceID
			lc.ConversationID = utils.GetConversationIDBySessionType(st.SourceID, constant.SingleChatType)
			lc.ConversationType = constant.SingleChatType
		case constant.GroupChatType:
			conversationID, conversationType, err := c.getConversationTypeByGroupID(st.SourceID)
			if err != nil {
				log.Error("internal", "getConversationTypeByGroupID database err:", err.Error())
				return
			}
			lc.GroupID = st.SourceID
			lc.ConversationID = conversationID
			lc.ConversationType = conversationType
		}
		c.addFaceURLAndName(&lc)
		err := c.db.UpdateConversation(&lc)
		if err != nil {
			log.Error("internal", "setConversationFaceUrlAndNickName database err:", err.Error())
			return
		}
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: lc.ConversationID, Action: constant.ConChange, Args: []string{lc.ConversationID}}})

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
		cLists, err := c.db.GetMultipleConversationDB(cidList)
		if err != nil {
			log.Error("internal", "getMultipleConversationModel err :", err.Error())
		} else {
			var newCList []*model_struct.LocalConversation
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
		cLists, err := c.db.GetMultipleConversationDB(cidList)
		if err != nil {
			log.Error("internal", "getMultipleConversationModel err :", err.Error())
		} else {
			if cLists != nil {
				log.Info("internal", "getMultipleConversationModel success :", cLists)
				c.ConversationListener.OnNewConversation(utils.StructToJsonString(cLists))
			}
		}
	case constant.ConChangeDirect:
		cidList := node.Args.(string)
		c.ConversationListener.OnConversationChanged(cidList)

	case constant.NewConDirect:
		cidList := node.Args.(string)
		log.Debug("internal", "NewConversation", cidList)
		c.ConversationListener.OnNewConversation(cidList)

	case constant.ConversationLatestMsgHasRead:
		hasReadMsgList := node.Args.(map[string][]string)
		var result []*model_struct.LocalConversation
		var latestMsg sdk_struct.MsgStruct
		var lc model_struct.LocalConversation
		for conversationID, msgIDList := range hasReadMsgList {
			LocalConversation, err := c.db.GetConversation(conversationID)
			if err != nil {
				log.Error("internal", "get conversation err", err.Error(), conversationID)
				continue
			}
			err = utils.JsonStringToStruct(LocalConversation.LatestMsg, &latestMsg)
			if err != nil {
				log.Error("internal", "JsonStringToStruct err", err.Error(), conversationID)
				continue
			}
			if utils.IsContain(latestMsg.ClientMsgID, msgIDList) {
				latestMsg.IsRead = true
				lc.ConversationID = conversationID
				lc.LatestMsg = utils.StructToJsonString(latestMsg)
				LocalConversation.LatestMsg = utils.StructToJsonString(latestMsg)
				err := c.db.UpdateConversation(&lc)
				if err != nil {
					log.Error("internal", "UpdateConversation database err:", err.Error())
					continue
				} else {
					result = append(result, LocalConversation)
				}
			}
		}
		if result != nil {
			log.Info("internal", "getMultipleConversationModel success :", result)
			c.ConversationListener.OnNewConversation(utils.StructToJsonString(result))
		}
	case constant.SyncConversation:
		operationID := node.Args.(string)
		log.Debug(operationID, "reconn sync conversation start")
		c.SyncConversations(operationID, 0)
		c.SyncConversationUnreadCount(operationID)
		totalUnreadCount, err := c.db.GetTotalUnreadMsgCountDB()
		if err != nil {
			log.Error("internal", "TotalUnreadMessageChanged database err:", err.Error())
		} else {
			c.ConversationListener.OnTotalUnreadMessageCountChanged(totalUnreadCount)
		}

	}
}
func (c *Conversation) doUpdateMessage(c2v common.Cmd2Value) {
	if c.ConversationListener == nil {
		log.Error("internal", "not set conversationListener")
		return
	}

	node := c2v.Value.(common.UpdateMessageNode)
	switch node.Action {
	case constant.UpdateMsgFaceUrlAndNickName:
		args := node.Args.(common.UpdateMessageInfo)
		var conversationType int32
		if args.GroupID == "" {
			conversationType = constant.SingleChatType
		} else {
			var err error
			_, conversationType, err = c.getConversationTypeByGroupID(args.GroupID)
			if err != nil {
				log.Error("internal", "getConversationTypeByGroupID database err:", err.Error())
				return
			}
		}
		err := c.db.UpdateMsgSenderFaceURLAndSenderNicknameController(args.UserID, args.FaceURL, args.Nickname, int(conversationType), args.GroupID)
		if err != nil {
			log.Error("internal", "UpdateMsgSenderFaceURLAndSenderNickname err:", err.Error())
		}

	}

}
func (c *Conversation) doSyncReactionExtensions(c2v common.Cmd2Value) {
	if c.ConversationListener == nil {
		log.Error("internal", "not set conversationListener")
		return
	}
	node := c2v.Value.(common.SyncReactionExtensionsNode)
	switch node.Action {
	case constant.SyncMessageListReactionExtensions:
		args := node.Args.(syncReactionExtensionParams)
		log.Error(node.OperationID, "come SyncMessageListReactionExtensions", args)
		var reqList []server_api_params.OperateMessageListReactionExtensionsReq
		for _, v := range args.MessageList {
			var temp server_api_params.OperateMessageListReactionExtensionsReq
			temp.ClientMsgID = v.ClientMsgID
			temp.MsgFirstModifyTime = v.MsgFirstModifyTime
			reqList = append(reqList, temp)
		}
		var apiReq server_api_params.GetMessageListReactionExtensionsReq
		apiReq.SourceID = args.SourceID
		apiReq.TypeKeyList = args.TypeKeyList
		apiReq.SessionType = args.SessionType
		apiReq.MessageReactionKeyList = reqList
		apiReq.IsExternalExtensions = args.IsExternalExtension
		apiReq.OperationID = node.OperationID
		var apiResp server_api_params.GetMessageListReactionExtensionsResp
		err := c.p.PostReturn(constant.GetMessageListReactionExtensionsRouter, apiReq, &apiResp)
		if err != nil {
			log.NewError(node.OperationID, utils.GetSelfFuncName(), "getMessageListReactionExtensions err:", err.Error())
			return
		}
		for _, result := range apiResp {
			log.Warn(node.OperationID, "api return reslut is:", result.ClientMsgID, result.ReactionExtensionList)

		}
		onLocal := func(data []*model_struct.LocalChatLogReactionExtensions) []*server_api_params.SingleMessageExtensionResult {
			var result []*server_api_params.SingleMessageExtensionResult
			for _, v := range data {
				temp := new(server_api_params.SingleMessageExtensionResult)
				tempMap := make(map[string]*server_api_params.KeyValue)
				_ = json.Unmarshal(v.LocalReactionExtensions, &tempMap)
				if len(args.TypeKeyList) != 0 {
					for s, _ := range tempMap {
						if !utils.IsContain(s, args.TypeKeyList) {
							delete(tempMap, s)
						}
					}
				}

				temp.ReactionExtensionList = tempMap
				temp.ClientMsgID = v.ClientMsgID
				result = append(result, temp)
			}
			return result
		}(args.ExtendMessageList)
		var onServer []*server_api_params.SingleMessageExtensionResult
		for _, v := range apiResp {
			if v.ErrCode == 0 {
				onServer = append(onServer, v)
			}
		}
		aInBNot, _, sameA, _ := common.CheckReactionExtensionsDiff(onServer, onLocal)
		for _, v := range aInBNot {
			log.Error(node.OperationID, "come InsertMessageReactionExtension", args, v.ClientMsgID)
			if len(v.ReactionExtensionList) > 0 {
				temp := model_struct.LocalChatLogReactionExtensions{ClientMsgID: v.ClientMsgID, LocalReactionExtensions: []byte(utils.StructToJsonString(v.ReactionExtensionList))}
				err := c.db.InsertMessageReactionExtension(&temp)
				if err != nil {
					log.Error(node.OperationID, "InsertMessageReactionExtension err:", err.Error())
					continue
				}
			}
			var changedKv []*server_api_params.KeyValue
			for _, value := range v.ReactionExtensionList {
				changedKv = append(changedKv, value)
			}
			if len(changedKv) > 0 {
				c.msgListener.OnRecvMessageExtensionsChanged(v.ClientMsgID, utils.StructToJsonString(changedKv))
			}
		}
		for _, result := range sameA {
			log.Warn(node.OperationID, "result is ", result.ReactionExtensionList, result.ClientMsgID)
		}
		for _, v := range sameA {
			log.Error(node.OperationID, "come sameA", v.ClientMsgID, v.ReactionExtensionList)
			tempMap := make(map[string]*server_api_params.KeyValue)
			for _, extensions := range args.ExtendMessageList {
				if v.ClientMsgID == extensions.ClientMsgID {
					_ = json.Unmarshal(extensions.LocalReactionExtensions, &tempMap)
					break
				}
			}
			if len(v.ReactionExtensionList) == 0 {
				err := c.db.DeleteMessageReactionExtension(v.ClientMsgID)
				if err != nil {
					log.Error(node.OperationID, "DeleteMessageReactionExtension err:", err.Error())
					continue
				}
				var deleteKeyList []string
				for key, _ := range tempMap {
					deleteKeyList = append(deleteKeyList, key)
				}
				if len(deleteKeyList) > 0 {
					c.msgListener.OnRecvMessageExtensionsDeleted(v.ClientMsgID, utils.StructToJsonString(deleteKeyList))
				}
			} else {
				deleteKeyList, changedKv := func(local, server map[string]*server_api_params.KeyValue) ([]string, []*server_api_params.KeyValue) {
					var deleteKeyList []string
					var changedKv []*server_api_params.KeyValue
					for k, v := range local {
						ia, ok := server[k]
						if ok {
							//服务器不同的kv
							if ia.Value != v.Value {
								changedKv = append(changedKv, ia)
							}
						} else {
							//服务器已经没有kv
							deleteKeyList = append(deleteKeyList, k)
						}
					}
					//从服务器新增的kv
					for k, v := range server {
						_, ok := local[k]
						if !ok {
							changedKv = append(changedKv, v)

						}
					}
					return deleteKeyList, changedKv
				}(tempMap, v.ReactionExtensionList)
				extendMsg := model_struct.LocalChatLogReactionExtensions{ClientMsgID: v.ClientMsgID, LocalReactionExtensions: []byte(utils.StructToJsonString(v.ReactionExtensionList))}
				err = c.db.UpdateMessageReactionExtension(&extendMsg)
				if err != nil {
					log.Error(node.OperationID, "UpdateMessageReactionExtension err:", err.Error())
					continue
				}
				if len(deleteKeyList) > 0 {
					c.msgListener.OnRecvMessageExtensionsDeleted(v.ClientMsgID, utils.StructToJsonString(deleteKeyList))
				}
				if len(changedKv) > 0 {
					c.msgListener.OnRecvMessageExtensionsChanged(v.ClientMsgID, utils.StructToJsonString(changedKv))
				}
			}
			//err := c.db.GetAndUpdateMessageReactionExtension(v.ClientMsgID, v.ReactionExtensionList)
			//if err != nil {
			//	log.Error(node.OperationID, "GetAndUpdateMessageReactionExtension err:", err.Error())
			//	continue
			//}
			//var changedKv []*server_api_params.KeyValue
			//for _, value := range v.ReactionExtensionList {
			//	changedKv = append(changedKv, value)
			//}
			//if len(changedKv) > 0 {
			//	c.msgListener.OnRecvMessageExtensionsChanged(v.ClientMsgID, utils.StructToJsonString(changedKv))
			//}
		}
	case constant.SyncMessageListTypeKeyInfo:
		messageList := node.Args.([]*sdk_struct.MsgStruct)
		var sourceID string
		var sessionType int32
		var reqList []server_api_params.OperateMessageListReactionExtensionsReq
		var temp server_api_params.OperateMessageListReactionExtensionsReq
		for _, v := range messageList {
			message, err := c.db.GetMessageController(v)
			if err != nil {
				log.Error(node.OperationID, "GetMessageController err:", err.Error(), *v)
				continue
			}
			temp.ClientMsgID = message.ClientMsgID
			temp.MsgFirstModifyTime = message.MsgFirstModifyTime
			reqList = append(reqList, temp)
			switch message.SessionType {
			case constant.SingleChatType:
				sourceID = message.SendID + message.RecvID
			case constant.NotificationChatType:
				sourceID = message.RecvID
			case constant.GroupChatType, constant.SuperGroupChatType:
				sourceID = message.RecvID
			}
			sessionType = message.SessionType
		}
		var apiReq server_api_params.GetMessageListReactionExtensionsReq
		apiReq.SourceID = sourceID
		apiReq.SessionType = sessionType
		apiReq.MessageReactionKeyList = reqList
		apiReq.OperationID = node.OperationID
		var apiResp server_api_params.GetMessageListReactionExtensionsResp
		err := c.p.PostReturnWithTimeOut(constant.GetMessageListReactionExtensionsRouter, apiReq, &apiResp, time.Second*2)
		if err != nil {
			log.Error(node.OperationID, "GetMessageListReactionExtensions from server err:", err.Error(), apiReq)
			return
		}
		var messageChangedList []*messageKvList
		for _, v := range apiResp {
			if v.ErrCode == 0 {
				var changedKv []*server_api_params.KeyValue
				var prefixTypeKey []string
				extendMsg, _ := c.db.GetMessageReactionExtension(v.ClientMsgID)
				localKV := make(map[string]*server_api_params.KeyValue)
				_ = json.Unmarshal(extendMsg.LocalReactionExtensions, &localKV)
				for typeKey, value := range v.ReactionExtensionList {
					oldValue, ok := localKV[typeKey]
					if ok {
						if !cmp.Equal(value, oldValue) {
							localKV[typeKey] = value
							prefixTypeKey = append(prefixTypeKey, getPrefixTypeKey(typeKey))
							changedKv = append(changedKv, value)
						}
					} else {
						localKV[typeKey] = value
						prefixTypeKey = append(prefixTypeKey, getPrefixTypeKey(typeKey))
						changedKv = append(changedKv, value)

					}

				}
				extendMsg.LocalReactionExtensions = []byte(utils.StructToJsonString(localKV))
				_ = c.db.UpdateMessageReactionExtension(extendMsg)
				if len(changedKv) > 0 {
					c.msgListener.OnRecvMessageExtensionsChanged(extendMsg.ClientMsgID, utils.StructToJsonString(changedKv))
				}
				prefixTypeKey = utils.RemoveRepeatedStringInList(prefixTypeKey)
				if len(prefixTypeKey) > 0 && c.msgKvListener != nil {
					var result []*sdk.SingleTypeKeyInfoSum
					oneMessageChanged := new(messageKvList)
					oneMessageChanged.ClientMsgID = extendMsg.ClientMsgID
					for _, v := range prefixTypeKey {
						singleResult := new(sdk.SingleTypeKeyInfoSum)
						singleResult.TypeKey = v
						for typeKey, value := range localKV {
							if strings.HasPrefix(typeKey, v) {
								singleTypeKeyInfo := new(sdk.SingleTypeKeyInfo)
								err := json.Unmarshal([]byte(value.Value), singleTypeKeyInfo)
								if err != nil {
									continue
								}
								if _, ok := singleTypeKeyInfo.InfoList[c.loginUserID]; ok {
									singleResult.IsContainSelf = true
								}
								for _, info := range singleTypeKeyInfo.InfoList {
									v := *info
									singleResult.InfoList = append(singleResult.InfoList, &v)
								}
								singleResult.Counter += singleTypeKeyInfo.Counter
							}
						}
						result = append(result, singleResult)
					}
					oneMessageChanged.ChangedKvList = result
					messageChangedList = append(messageChangedList, oneMessageChanged)
				}
			}
		}
		if len(messageChangedList) > 0 && c.msgKvListener != nil {
			c.msgKvListener.OnMessageKvInfoChanged(utils.StructToJsonString(messageChangedList))
		}

	}

}

type messageKvList struct {
	ClientMsgID   string                      `json:"clientMsgID"`
	ChangedKvList []*sdk.SingleTypeKeyInfoSum `json:"changedKvList"`
}

func (c *Conversation) Work(c2v common.Cmd2Value) {

	log.Info("internal", "doListener work..", c2v.Cmd)

	switch c2v.Cmd {
	case constant.CmdDeleteConversation:
		if c.LoginStatus() == constant.Logout {
			log.Warn("", "m.LoginStatus() == constant.Logout, Goexit()")
			runtime.Goexit()
		}
		log.Info("internal", "CmdDeleteConversation start ..", c2v.Cmd)
		c.doDeleteConversation(c2v)
		log.Info("internal", "CmdDeleteConversation end..", c2v.Cmd)
	case constant.CmdNewMsgCome:
		if c.LoginStatus() == constant.Logout {
			log.Warn("", "m.LoginStatus() == constant.Logout, Goexit()")
			runtime.Goexit()
		}
		log.Info("internal", "doMsgNew start..", c2v.Cmd)
		c.doMsgNew(c2v)
		log.Info("internal", "doMsgNew end..", c2v.Cmd)
	case constant.CmdSuperGroupMsgCome:
		if c.LoginStatus() == constant.Logout {
			log.Warn("", "m.LoginStatus() == constant.Logout, Goexit()")
			runtime.Goexit()
		}
		log.Info("internal", "doSuperGroupMsgNew start..", c2v.Cmd)
		c.doSuperGroupMsgNew(c2v)
		log.Info("internal", "doSuperGroupMsgNew end..", c2v.Cmd)
	case constant.CmdUpdateConversation:
		if c.LoginStatus() == constant.Logout {
			log.Warn("", "m.LoginStatus() == constant.Logout, Goexit()")
			runtime.Goexit()
		}
		log.Info("internal", "doUpdateConversation start ..", c2v.Cmd)
		c.doUpdateConversation(c2v)
		log.Info("internal", "doUpdateConversation end..", c2v.Cmd)
	case constant.CmdUpdateMessage:
		if c.LoginStatus() == constant.Logout {
			log.Warn("", "m.LoginStatus() == constant.Logout, Goexit()")
			runtime.Goexit()
		}
		log.Info("internal", "doUpdateMessage start ..", c2v.Cmd)
		c.doUpdateMessage(c2v)
		log.Info("internal", "doUpdateMessage end..", c2v.Cmd)
	case constant.CmSyncReactionExtensions:
		log.Info("internal", "doSyncReactionExtensions start ..", c2v.Cmd)
		c.doSyncReactionExtensions(c2v)
		log.Info("internal", "doSyncReactionExtensions end..", c2v.Cmd)

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
			if msg.AttachedInfoElem.IsEncryption && c.encryptionKey != "" && msg.AttachedInfoElem.InEncryptStatus {
				var newContent []byte
				log.NewDebug("", utils.GetSelfFuncName(), "org content, key", msg.Content, c.encryptionKey, []byte(msg.Content), msg.CreateTime, msg.AttachedInfoElem, msg.AttachedInfo)
				newContent, err = utils.AesDecrypt([]byte(msg.Content), []byte(c.encryptionKey))
				msg.Content = string(newContent)
				msg.AttachedInfoElem.InEncryptStatus = false
				msg.AttachedInfo = utils.StructToJsonString(msg.AttachedInfoElem)
			}
		case constant.Picture:
			err = utils.JsonStringToStruct(msg.Content, &msg.PictureElem)
		case constant.Voice:
			err = utils.JsonStringToStruct(msg.Content, &msg.SoundElem)
		case constant.Video:
			err = utils.JsonStringToStruct(msg.Content, &msg.VideoElem)
		case constant.File:
			err = utils.JsonStringToStruct(msg.Content, &msg.FileElem)
		case constant.AdvancedText:
			err = utils.JsonStringToStruct(msg.Content, &msg.MessageEntityElem)
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
		case constant.CustomMsgNotTriggerConversation:
			err = utils.JsonStringToStruct(msg.Content, &msg.CustomElem)
		case constant.CustomMsgOnlineOnly:
			err = utils.JsonStringToStruct(msg.Content, &msg.CustomElem)
		}
	}

	return utils.Wrap(err, "")
}
func (c *Conversation) updateConversation(lc *model_struct.LocalConversation, cs map[string]*model_struct.LocalConversation) {
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
func mapConversationToList(m map[string]*model_struct.LocalConversation) (cs []*model_struct.LocalConversation) {
	for _, v := range m {
		cs = append(cs, v)
	}
	return cs
}
func (c *Conversation) addFaceURLAndName(lc *model_struct.LocalConversation) {
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

	case constant.GroupChatType, constant.SuperGroupChatType:
		g, err := c.full.GetGroupInfoFromLocal2Svr(lc.GroupID, lc.ConversationType)
		if err != nil {
			log.Error(operationID, "GetGroupInfoByGroupID err", err.Error(), lc.GroupID, lc.ConversationType)
			return
		}
		lc.ShowName = g.GroupName
		lc.FaceURL = g.FaceURL

	}
}
