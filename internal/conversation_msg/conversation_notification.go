package conversation_msg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
	"open_im_sdk/internal/business"
	"open_im_sdk/internal/cache"
	"open_im_sdk/internal/friend"
	"open_im_sdk/internal/full"
	"open_im_sdk/internal/group"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/internal/signaling"
	"open_im_sdk/internal/user"
	"open_im_sdk/internal/util"
	"runtime"
	"strings"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/google/go-cmp/cmp"

	workMoments "open_im_sdk/internal/work_moments"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/syncer"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"

	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"sync"
	"time"
)

func NewNotificationConversation(ctx context.Context, db db_interface.DataBase,
	ch chan common.Cmd2Value,
	friend *friend.Friend, group *group.Group, user *user.User,
	conversationListener open_im_sdk_callback.OnConversationListener,
	msgListener open_im_sdk_callback.OnAdvancedMsgListener, signaling *signaling.LiveSignaling,
	workMoments *workMoments.WorkMoments, business *business.Business, cache *cache.Cache, full *full.Full, id2MinSeq map[string]int64) *NotificationConversation {
	n := &NotificationConversation{Ws: ws, db: db, p: p, recvCH: ch, loginUserID: loginUserID, platformID: platformID,
		DataDir: dataDir, friend: friend, group: group, user: user, ObjectStorage: objectStorage,
		signaling: signaling, workMoments: workMoments,
		full: full, id2MinSeq: id2MinSeq, encryptionKey: encryptionKey, business: business, IsExternalExtensions: isExternalExtensions}
	n.SetMsgListener(msgListener)
	n.SetConversationListener(conversationListener)
	n.initSyncer()
	n.cache = cache
	return n
}

type NotificationConversation struct {
	conversationSyncer *syncer.Syncer[*model_struct.LocalConversation, string]
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
	workMoments          *workMoments.WorkMoments
	business             *business.Business

	cache          *cache.Cache
	full           *full.Full
	tempMessageMap sync.Map
	encryptionKey  string

	id2MinSeq            map[string]int64
	IsExternalExtensions bool

	listenerForService open_im_sdk_callback.OnListenerForService
}

func (c *NotificationConversation) SetListenerForService(listener open_im_sdk_callback.OnListenerForService) {
	c.listenerForService = listener
}

func (c *NotificationConversation) MsgListener() open_im_sdk_callback.OnAdvancedMsgListener {
	return c.msgListener
}

func (c *NotificationConversation) SetSignaling(signaling *signaling.LiveSignaling) {
	c.signaling = signaling
}

func (c *NotificationConversation) SetMsgListener(msgListener open_im_sdk_callback.OnAdvancedMsgListener) {
	c.msgListener = msgListener
}
func (c *NotificationConversation) SetMsgKvListener(msgKvListener open_im_sdk_callback.OnMessageKvInfoListener) {
	c.msgKvListener = msgKvListener
}
func (c *NotificationConversation) SetBatchMsgListener(batchMsgListener open_im_sdk_callback.OnBatchMsgListener) {
	c.batchMsgListener = batchMsgListener
}

func (c *NotificationConversation) initSyncer() {

}

func (c *NotificationConversation) GetCh() chan common.Cmd2Value {
	return c.recvCH
}

func (c *NotificationConversation) getServerConversationList(ctx context.Context) ([]*model_struct.LocalConversation, error) {
	resp, err := util.CallApi[conversation.GetAllConversationsResp](ctx, constant.GetAllConversationsRouter, conversation.GetAllConversationsReq{OwnerUserID: c.loginUserID})
	if err != nil {
		return nil, err
	}
	return util.Batch(ServerConversationToLocal, resp.Conversations), nil
}

func (c *NotificationConversation) SyncConversations(ctx context.Context) error {
	ccTime := time.Now()
	conversationsOnServer, err := c.getServerConversationList(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "get server cost time", "cost time", time.Since(ccTime), "conversation on server", conversationsOnServer)
	conversationsOnLocal, err := c.db.GetAllConversationListToSync(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "get local cost time", "cost time", time.Since(ccTime), "conversation on local", conversationsOnLocal)
	for _, v := range conversationsOnServer {
		c.addFaceURLAndName(ctx, v)
	}
	log.ZDebug(ctx, "get local cost time", "cost time", time.Since(ccTime), "conversation on local", conversationsOnLocal)
	if err = c.conversationSyncer.Sync(ctx, conversationsOnServer, conversationsOnLocal, func(ctx context.Context, state int, conversation *model_struct.LocalConversation) error {
		if state == syncer.Update {
			c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: conversation.ConversationID, Action: constant.ConChange, Args: ""}})
		}
		return nil
	}, true); err != nil {
		return err
	}
	conversationsOnLocal, err = c.db.GetAllConversationListToSync(ctx)
	if err != nil {
		return err
	}
	c.cache.UpdateConversations(conversationsOnLocal)
	return nil
}

func (c *NotificationConversation) doReactionMsgModifier(ctx context.Context, msgReactionList []*sdk_struct.MsgStruct) {
	for _, msgStruct := range msgReactionList {
		var n server_api_params.ReactionMessageModifierNotification
		err := json.Unmarshal([]byte(msgStruct.Content), &n)
		if err != nil {
			log.Error("internal", "unmarshal failed err:", err.Error(), *msgStruct)
			continue
		}
		switch n.Operation {
		case constant.AddMessageExtensions:
			var reactionExtensionList []*sdkws.KeyValue
			for _, value := range n.SuccessReactionExtensionList {
				reactionExtensionList = append(reactionExtensionList, value)
			}
			if !(msgStruct.SendID == c.loginUserID && msgStruct.SenderPlatformID == c.platformID) {
				c.msgListener.OnRecvMessageExtensionsAdded(n.ClientMsgID, utils.StructToJsonString(reactionExtensionList))
			}
		case constant.SetMessageExtensions:
			err = c.db.GetAndUpdateMessageReactionExtension(ctx, n.ClientMsgID, n.SuccessReactionExtensionList)
			if err != nil {
				log.Error("internal", "GetAndUpdateMessageReactionExtension err:", err.Error())
				continue
			}
			var reactionExtensionList []*sdkws.KeyValue
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
		err2 := c.db.UpdateMessageController(ctx, &t)
		if err2 != nil {
			log.Error("internal", "unmarshal failed err:", err2.Error(), t)
			continue
		}
	}
}
func (c *NotificationConversation) doReactionMsgDeleter(ctx context.Context, msgReactionList []*sdk_struct.MsgStruct) {
	for _, msgStruct := range msgReactionList {
		var n server_api_params.ReactionMessageDeleteNotification
		err := json.Unmarshal([]byte(msgStruct.Content), &n)
		if err != nil {
			log.Error("internal", "unmarshal failed err:", err.Error(), *msgStruct)
			continue
		}
		err = c.db.DeleteAndUpdateMessageReactionExtension(ctx, n.ClientMsgID, n.SuccessReactionExtensionList)
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
func (c *NotificationConversation) QuoteMsgRevokeHandle(ctx context.Context, v *model_struct.LocalChatLog, revokeMsgIDList []*sdk_struct.MessageRevoked) {
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
	err = c.db.UpdateMessageController(ctx, v)
	if err != nil {
		log.NewError("internal", "unmarshall failed", v)
	}
}

// todo
func (c *NotificationConversation) doDeleteConversation(c2v common.Cmd2Value) {
	if c.ConversationListener == nil {
		log.Error("internal", "not set conversationListener")
		return
	}
	node := c2v.Value.(common.DeleteConNode)
	ctx := mcontext.NewCtx(utils.OperationIDGenerator())
	//Mark messages related to this conversation for deletion
	err := c.db.UpdateMessageStatusBySourceID(context.Background(), node.SourceID, constant.MsgStatusHasDeleted, int32(node.SessionType))
	if err != nil {
		log.Error("internal", "setMessageStatusBySourceID err:", err.Error())
		return
	}
	//Reset the session information, empty session
	err = c.db.ResetConversation(ctx, node.ConversationID)
	if err != nil {
		log.Error("internal", "ResetConversation err:", err.Error())
	}
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{"", constant.TotalUnreadMessageChanged, ""}})
}
func (c *NotificationConversation) doMsgReadState(ctx context.Context, msgReadList []*sdk_struct.MsgStruct) {
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
			m, err := c.db.GetMessage(ctx, v)
			if err != nil {
				log.Error("internal", "GetMessage err:", err, "ClientMsgID", v)
				continue
			}
			attachInfo := sdk_struct.AttachedInfoElem{}
			_ = utils.JsonStringToStruct(m.AttachedInfo, &attachInfo)
			attachInfo.HasReadTime = rd.SendTime
			m.AttachedInfo = utils.StructToJsonString(attachInfo)
			m.IsRead = true
			err = c.db.UpdateMessage(ctx, m)
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

func (c *NotificationConversation) getConversationTypeByGroupID(ctx context.Context, groupID string) (conversationID string, conversationType int32, err error) {
	g, err := c.full.GetGroupInfoByGroupID(ctx, groupID)
	if err != nil {
		return "", 0, utils.Wrap(err, "get group info error")
	}
	switch g.GroupType {
	case constant.NormalGroup:
		return utils.GetConversationIDBySessionType(groupID, constant.GroupChatType), constant.GroupChatType, nil
	case constant.SuperGroup, constant.WorkingGroup:
		return utils.GetConversationIDBySessionType(groupID, constant.SuperGroupChatType), constant.SuperGroupChatType, nil
	default:
		return "", 0, utils.Wrap(errors.New("err groupType"), "group type err")
	}
}

// todo
func (c *NotificationConversation) doUpdateConversation(c2v common.Cmd2Value) {
	if c.ConversationListener == nil {
		log.Error("internal", "not set conversationListener")
		return
	}
	ctx := mcontext.NewCtx(utils.OperationIDGenerator())
	node := c2v.Value.(common.UpdateConNode)
	switch node.Action {
	case constant.AddConOrUpLatMsg:
		var list []*model_struct.LocalConversation
		lc := node.Args.(model_struct.LocalConversation)
		oc, err := c.db.GetConversation(ctx, lc.ConversationID)
		if err == nil {
			log.Info("this is old conversation", *oc)
			if lc.LatestMsgSendTime >= oc.LatestMsgSendTime { //The session update of asynchronous messages is subject to the latest sending time
				err := c.db.UpdateColumnsConversation(nil, node.ConID, map[string]interface{}{"latest_msg_send_time": lc.LatestMsgSendTime, "latest_msg": lc.LatestMsg})
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
			err4 := c.db.InsertConversation(ctx, &lc)
			if err4 != nil {
				log.Error("internal", "insert new conversation err:", err4.Error())
			} else {
				list = append(list, &lc)
				c.ConversationListener.OnNewConversation(utils.StructToJsonString(list))
			}
		}

	case constant.UnreadCountSetZero:
		if err := c.db.UpdateColumnsConversation(ctx, node.ConID, map[string]interface{}{"unread_count": 0}); err != nil {
			log.Error("internal", "UpdateColumnsConversation err", err.Error(), node.ConID)
		} else {
			totalUnreadCount, err := c.db.GetTotalUnreadMsgCountDB(ctx)
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
		err := c.db.IncrConversationUnreadCount(ctx, node.ConID)
		if err != nil {
			log.Error("internal", "incrConversationUnreadCount database err:", err.Error())
			return
		}
	case constant.TotalUnreadMessageChanged:
		totalUnreadCount, err := c.db.GetTotalUnreadMsgCountDB(ctx)
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
			conversationID, conversationType, err := c.getConversationTypeByGroupID(ctx, st.SourceID)
			if err != nil {
				log.Error("internal", "getConversationTypeByGroupID database err:", err.Error())
				return
			}
			lc.GroupID = st.SourceID
			lc.ConversationID = conversationID
			lc.ConversationType = conversationType
		}
		c.addFaceURLAndName(ctx, &lc)
		err := c.db.UpdateConversation(ctx, &lc)
		if err != nil {
			log.Error("internal", "setConversationFaceUrlAndNickName database err:", err.Error())
			return
		}
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: lc.ConversationID, Action: constant.ConChange, Args: []string{lc.ConversationID}}})

	case constant.UpdateLatestMessageChange:
		conversationID := node.ConID
		var latestMsg sdk_struct.MsgStruct
		l, err := c.db.GetConversation(ctx, conversationID)
		if err != nil {
			log.Error("internal", "getConversationLatestMsgModel err", err.Error())
		} else {
			err := json.Unmarshal([]byte(l.LatestMsg), &latestMsg)
			if err != nil {
				log.Error("internal", "latestMsg,Unmarshal err :", err.Error())
			} else {
				latestMsg.IsRead = true
				newLatestMessage := utils.StructToJsonString(latestMsg)
				err = c.db.UpdateColumnsConversation(nil, node.ConID, map[string]interface{}{"latest_msg_send_time": latestMsg.SendTime, "latest_msg": newLatestMessage})
				if err != nil {
					log.Error("internal", "updateConversationLatestMsgModel err :", err.Error())
				}
			}
		}
	case constant.ConChange:
		cidList := node.Args.([]string)
		cLists, err := c.db.GetMultipleConversationDB(ctx, cidList)
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
		cLists, err := c.db.GetMultipleConversationDB(ctx, cidList)
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
			LocalConversation, err := c.db.GetConversation(ctx, conversationID)
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
				err := c.db.UpdateConversation(ctx, &lc)
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
		c.SyncConversations(ctx)
		err := c.SyncConversationUnreadCount(ctx)
		if err != nil {
			log.Error(operationID, "reconn sync conversation unread count err", err.Error())
		}
		totalUnreadCount, err := c.db.GetTotalUnreadMsgCountDB(ctx)
		if err != nil {
			log.Error("internal", "TotalUnreadMessageChanged database err:", err.Error())
		} else {
			c.ConversationListener.OnTotalUnreadMessageCountChanged(totalUnreadCount)
		}
	}
}

func (c *NotificationConversation) SyncConversationUnreadCount(ctx context.Context) error {
	var conversationChangedList []string
	fmt.Println("test", c.cache)
	allConversations := c.cache.GetAllHasUnreadMessageConversations()
	allConversations = c.cache.GetAllHasUnreadMessageConversations()
	log.ZDebug(ctx, "get unread message length", "len", len(allConversations))
	for _, conversation := range allConversations {
		if deleteRows := c.db.DeleteConversationUnreadMessageList(ctx, conversation.ConversationID, conversation.UpdateUnreadCountTime); deleteRows > 0 {
			log.ZDebug(ctx, "DeleteConversationUnreadMessageList", conversation.ConversationID, conversation.UpdateUnreadCountTime, "delete rows:", deleteRows)
			if err := c.db.DecrConversationUnreadCount(ctx, conversation.ConversationID, deleteRows); err != nil {
				log.ZDebug(ctx, "DecrConversationUnreadCount", conversation.ConversationID, conversation.UpdateUnreadCountTime, "decr unread count err:", err.Error())
			} else {
				conversationChangedList = append(conversationChangedList, conversation.ConversationID)
			}
		}
	}
	if len(conversationChangedList) > 0 {
		if err := common.TriggerCmdUpdateConversation(common.UpdateConNode{Action: constant.ConChange, Args: conversationChangedList}, c.GetCh()); err != nil {
			return err
		}
	}
	return nil
}

// todo
func (c *NotificationConversation) doUpdateMessage(c2v common.Cmd2Value) {
	if c.ConversationListener == nil {
		log.Error("internal", "not set conversationListener")
		return
	}

	node := c2v.Value.(common.UpdateMessageNode)
	ctx := mcontext.NewCtx(utils.OperationIDGenerator())
	switch node.Action {
	case constant.UpdateMsgFaceUrlAndNickName:
		args := node.Args.(common.UpdateMessageInfo)
		var conversationType int32
		if args.GroupID == "" {
			conversationType = constant.SingleChatType
		} else {
			var err error
			_, conversationType, err = c.getConversationTypeByGroupID(ctx, args.GroupID)
			if err != nil {
				log.Error("internal", "getConversationTypeByGroupID database err:", err.Error())
				return
			}
		}
		err := c.db.UpdateMsgSenderFaceURLAndSenderNicknameController(ctx, args.UserID, args.FaceURL, args.Nickname, int(conversationType), args.GroupID)
		if err != nil {
			log.Error("internal", "UpdateMsgSenderFaceURLAndSenderNickname err:", err.Error())
		}

	}

}

// todo
func (c *NotificationConversation) doSyncReactionExtensions(c2v common.Cmd2Value) {
	if c.ConversationListener == nil {
		log.Error("internal", "not set conversationListener")
		return
	}
	node := c2v.Value.(common.SyncReactionExtensionsNode)
	ctx := mcontext.NewCtx(node.OperationID)
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
		// for _, result := range apiResp {
		// 	log.Warn(node.OperationID, "api return reslut is:", result.ClientMsgID, result.ReactionExtensionList)
		// }
		onLocal := func(data []*model_struct.LocalChatLogReactionExtensions) []*server_api_params.SingleMessageExtensionResult {
			var result []*server_api_params.SingleMessageExtensionResult
			for _, v := range data {
				temp := new(server_api_params.SingleMessageExtensionResult)
				tempMap := make(map[string]*sdkws.KeyValue)
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
				err := c.db.InsertMessageReactionExtension(ctx, &temp)
				if err != nil {
					log.Error(node.OperationID, "InsertMessageReactionExtension err:", err.Error())
					continue
				}
			}
			var changedKv []*sdkws.KeyValue
			for _, value := range v.ReactionExtensionList {
				changedKv = append(changedKv, value)
			}
			if len(changedKv) > 0 {
				c.msgListener.OnRecvMessageExtensionsChanged(v.ClientMsgID, utils.StructToJsonString(changedKv))
			}
		}
		// for _, result := range sameA {
		// log.ZWarn(ctx, "result", result.ReactionExtensionList, result.ClientMsgID)
		// }
		for _, v := range sameA {
			log.Error(node.OperationID, "come sameA", v.ClientMsgID, v.ReactionExtensionList)
			tempMap := make(map[string]*sdkws.KeyValue)
			for _, extensions := range args.ExtendMessageList {
				if v.ClientMsgID == extensions.ClientMsgID {
					_ = json.Unmarshal(extensions.LocalReactionExtensions, &tempMap)
					break
				}
			}
			if len(v.ReactionExtensionList) == 0 {
				err := c.db.DeleteMessageReactionExtension(ctx, v.ClientMsgID)
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
				deleteKeyList, changedKv := func(local, server map[string]*sdkws.KeyValue) ([]string, []*sdkws.KeyValue) {
					var deleteKeyList []string
					var changedKv []*sdkws.KeyValue
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
				err = c.db.UpdateMessageReactionExtension(ctx, &extendMsg)
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
			message, err := c.db.GetMessageController(ctx, v)
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
		var messageChangedList []*notificationKvList
		for _, v := range apiResp {
			if v.ErrCode == 0 {
				var changedKv []*sdkws.KeyValue
				var prefixTypeKey []string
				extendMsg, _ := c.db.GetMessageReactionExtension(ctx, v.ClientMsgID)
				localKV := make(map[string]*sdkws.KeyValue)
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
				_ = c.db.UpdateMessageReactionExtension(ctx, extendMsg)
				if len(changedKv) > 0 {
					c.msgListener.OnRecvMessageExtensionsChanged(extendMsg.ClientMsgID, utils.StructToJsonString(changedKv))
				}
				prefixTypeKey = utils.RemoveRepeatedStringInList(prefixTypeKey)
				if len(prefixTypeKey) > 0 && c.msgKvListener != nil {
					var result []*sdk.SingleTypeKeyInfoSum
					oneMessageChanged := new(notificationKvList)
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

type notificationKvList struct {
	ClientMsgID   string                      `json:"clientMsgID"`
	ChangedKvList []*sdk.SingleTypeKeyInfoSum `json:"changedKvList"`
}

func (c *NotificationConversation) Work(c2v common.Cmd2Value) {
	log.Info("internal", "doListener work..", c2v.Cmd)
	ctx := context.Background()
	switch c2v.Cmd {
	case constant.CmdDeleteConversation:
		if c.LoginStatus() == constant.Logout {
			log.ZWarn(ctx, "m.LoginStatus() == constant.Logout, Goexit()", nil)
			runtime.Goexit()
		}
		log.Info("internal", "CmdDeleteConversation start ..", c2v.Cmd)
		c.doDeleteConversation(c2v)
		log.Info("internal", "CmdDeleteConversation end..", c2v.Cmd)
	case constant.CmdUpdateConversation:
		if c.LoginStatus() == constant.Logout {
			log.ZWarn(ctx, "m.LoginStatus() == constant.Logout, Goexit()", nil)
			runtime.Goexit()
		}
		log.Info("internal", "doUpdateConversation start ..", c2v.Cmd)
		c.doUpdateConversation(c2v)
		log.Info("internal", "doUpdateConversation end..", c2v.Cmd)
	case constant.CmdUpdateMessage:
		if c.LoginStatus() == constant.Logout {
			log.ZWarn(ctx, "m.LoginStatus() == constant.Logout, Goexit()", nil)
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

func (c *NotificationConversation) addFaceURLAndName(ctx context.Context, lc *model_struct.LocalConversation) error {
	switch lc.ConversationType {
	case constant.SingleChatType, constant.NotificationChatType:
		faceUrl, name, err := c.cache.GetUserNameAndFaceURL(ctx, lc.UserID)
		if err != nil {
			return err
		}
		lc.FaceURL = faceUrl
		lc.ShowName = name

	case constant.GroupChatType, constant.SuperGroupChatType:
		g, err := c.full.GetGroupInfoFromLocal2Svr(ctx, lc.GroupID, lc.ConversationType)
		if err != nil {
			return err
		}
		lc.ShowName = g.GroupName
		lc.FaceURL = g.FaceURL
	}
	return nil
}
