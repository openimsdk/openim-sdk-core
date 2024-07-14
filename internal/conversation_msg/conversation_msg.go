// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conversation_msg

import (
	"context"
	"encoding/json"
	"errors"
	"math"

	"github.com/openimsdk/openim-sdk-core/v3/internal/business"
	"github.com/openimsdk/openim-sdk-core/v3/internal/cache"
	"github.com/openimsdk/openim-sdk-core/v3/internal/file"
	"github.com/openimsdk/openim-sdk-core/v3/internal/friend"
	"github.com/openimsdk/openim-sdk-core/v3/internal/full"
	"github.com/openimsdk/openim-sdk-core/v3/internal/group"
	"github.com/openimsdk/openim-sdk-core/v3/internal/interaction"
	"github.com/openimsdk/openim-sdk-core/v3/internal/user"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/page"
	sdk "github.com/openimsdk/openim-sdk-core/v3/pkg/sdk_params_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/syncer"
	pbConversation "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"

	"sort"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"

	"github.com/jinzhu/copier"
)

const (
	conversationSyncLimit = math.MaxInt64
)

var SearchContentType = []int{constant.Text, constant.AtText, constant.File}

type Conversation struct {
	*interaction.LongConnMgr
	conversationSyncer   *syncer.Syncer[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string]
	db                   db_interface.DataBase
	ConversationListener func() open_im_sdk_callback.OnConversationListener
	msgListener          func() open_im_sdk_callback.OnAdvancedMsgListener
	msgKvListener        func() open_im_sdk_callback.OnMessageKvInfoListener
	batchMsgListener     func() open_im_sdk_callback.OnBatchMsgListener
	recvCH               chan common.Cmd2Value
	loginUserID          string
	platformID           int32
	DataDir              string
	friend               *friend.Friend
	group                *group.Group
	user                 *user.User
	file                 *file.File
	business             *business.Business
	messageController    *MessageController
	cache                *cache.Cache[string, *model_struct.LocalConversation]
	full                 *full.Full
	maxSeqRecorder       MaxSeqRecorder
	IsExternalExtensions bool

	startTime time.Time

	typing *typing
}

func (c *Conversation) SetMsgListener(msgListener func() open_im_sdk_callback.OnAdvancedMsgListener) {
	c.msgListener = msgListener
}

func (c *Conversation) SetMsgKvListener(msgKvListener func() open_im_sdk_callback.OnMessageKvInfoListener) {
	c.msgKvListener = msgKvListener
}

func (c *Conversation) SetBatchMsgListener(batchMsgListener func() open_im_sdk_callback.OnBatchMsgListener) {
	c.batchMsgListener = batchMsgListener
}

func NewConversation(ctx context.Context, longConnMgr *interaction.LongConnMgr, db db_interface.DataBase,
	ch chan common.Cmd2Value, friend *friend.Friend, group *group.Group, user *user.User, business *business.Business,
	full *full.Full, file *file.File) *Conversation {
	info := ccontext.Info(ctx)
	n := &Conversation{db: db,
		LongConnMgr:          longConnMgr,
		recvCH:               ch,
		loginUserID:          info.UserID(),
		platformID:           info.PlatformID(),
		DataDir:              info.DataDir(),
		friend:               friend,
		group:                group,
		user:                 user,
		full:                 full,
		business:             business,
		file:                 file,
		messageController:    NewMessageController(db, ch),
		IsExternalExtensions: info.IsExternalExtensions(),
		maxSeqRecorder:       NewMaxSeqRecorder(),
	}
	n.typing = newTyping(n)
	n.initSyncer()
	n.cache = cache.NewCache[string, *model_struct.LocalConversation]()
	return n
}

func (c *Conversation) initSyncer() {
	c.conversationSyncer = syncer.New2[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](
		syncer.WithInsert[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](func(ctx context.Context, value *model_struct.LocalConversation) error {
			if err := c.batchAddFaceURLAndName(ctx, value); err != nil {
				return err
			}
			return c.db.InsertConversation(ctx, value)
		}),
		syncer.WithDelete[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](func(ctx context.Context, value *model_struct.LocalConversation) error {
			return c.db.DeleteConversation(ctx, value.ConversationID)
		}),
		syncer.WithUpdate[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](func(ctx context.Context, serverConversation, localConversation *model_struct.LocalConversation) error {
			return c.db.UpdateColumnsConversation(ctx, serverConversation.ConversationID,
				map[string]interface{}{"recv_msg_opt": serverConversation.RecvMsgOpt,
					"is_pinned": serverConversation.IsPinned, "is_private_chat": serverConversation.IsPrivateChat, "burn_duration": serverConversation.BurnDuration,
					"is_not_in_group": serverConversation.IsNotInGroup, "group_at_type": serverConversation.GroupAtType,
					"update_unread_count_time": serverConversation.UpdateUnreadCountTime,
					"attached_info":            serverConversation.AttachedInfo, "ex": serverConversation.Ex, "msg_destruct_time": serverConversation.MsgDestructTime,
					"is_msg_destruct": serverConversation.IsMsgDestruct,
					"max_seq":         serverConversation.MaxSeq, "min_seq": serverConversation.MinSeq, "has_read_seq": serverConversation.HasReadSeq})
		}),
		syncer.WithUUID[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](func(value *model_struct.LocalConversation) string {
			return value.ConversationID
		}),
		syncer.WithEqual[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](func(server, local *model_struct.LocalConversation) bool {
			if server.RecvMsgOpt != local.RecvMsgOpt ||
				server.IsPinned != local.IsPinned ||
				server.IsPrivateChat != local.IsPrivateChat ||
				server.BurnDuration != local.BurnDuration ||
				server.IsNotInGroup != local.IsNotInGroup ||
				server.GroupAtType != local.GroupAtType ||
				server.UpdateUnreadCountTime != local.UpdateUnreadCountTime ||
				server.AttachedInfo != local.AttachedInfo ||
				server.Ex != local.Ex ||
				server.MaxSeq != local.MaxSeq ||
				server.MinSeq != local.MinSeq ||
				server.MsgDestructTime != local.MsgDestructTime ||
				server.IsMsgDestruct != local.IsMsgDestruct {
				log.ZDebug(context.Background(), "not same", "conversationID", server.ConversationID, "server", server.RecvMsgOpt, "local", local.RecvMsgOpt)
				return false
			}
			return true
		}),
		syncer.WithNotice[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](func(ctx context.Context, state int, server, local *model_struct.LocalConversation) error {
			if state == syncer.Update || state == syncer.Insert {
				c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: server.ConversationID, Action: constant.ConChange, Args: []string{server.ConversationID}}})
			}
			return nil
		}),
		syncer.WithBatchInsert[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](func(ctx context.Context, values []*model_struct.LocalConversation) error {
			if err := c.batchAddFaceURLAndName(ctx, values...); err != nil {
				return err
			}
			return c.db.BatchInsertConversationList(ctx, values)
		}),
		syncer.WithDeleteAll[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](func(ctx context.Context, _ string) error {
			return c.db.DeleteAllConversation(ctx)
		}),
		syncer.WithBatchPageReq[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](func(entityID string) page.PageReq {
			return &pbConversation.GetOwnerConversationReq{UserID: entityID,
				Pagination: &sdkws.RequestPagination{ShowNumber: 300}}
		}),
		syncer.WithBatchPageRespConvertFunc[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](func(resp *pbConversation.GetOwnerConversationResp) []*model_struct.LocalConversation {
			return datautil.Batch(ServerConversationToLocal, resp.Conversations)
		}),
		syncer.WithReqApiRouter[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](constant.GetOwnerConversationRouter),
		syncer.WithFullSyncLimit[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](conversationSyncLimit),
	)

}

func (c *Conversation) GetCh() chan common.Cmd2Value {
	return c.recvCH
}

type onlineMsgKey struct {
	ClientMsgID string
	ServerMsgID string
}

func (c *Conversation) doMsgNew(c2v common.Cmd2Value) {
	allMsg := c2v.Value.(sdk_struct.CmdNewMsgComeToConversation).Msgs
	ctx := c2v.Ctx
	var isTriggerUnReadCount bool
	insertMsg := make(map[string][]*model_struct.LocalChatLog, 10)
	updateMsg := make(map[string][]*model_struct.LocalChatLog, 10)
	var exceptionMsg []*model_struct.LocalErrChatLog
	//var unreadMessages []*model_struct.LocalConversationUnreadMessage
	var newMessages sdk_struct.NewMsgList
	// var reactionMsgModifierList, reactionMsgDeleterList sdk_struct.NewMsgList
	var isUnreadCount, isConversationUpdate, isHistory, isNotPrivate, isSenderConversationUpdate bool
	conversationChangedSet := make(map[string]*model_struct.LocalConversation)
	newConversationSet := make(map[string]*model_struct.LocalConversation)
	conversationSet := make(map[string]*model_struct.LocalConversation)
	phConversationChangedSet := make(map[string]*model_struct.LocalConversation)
	phNewConversationSet := make(map[string]*model_struct.LocalConversation)
	log.ZDebug(ctx, "message come here conversation ch", "conversation length", len(allMsg))
	b := time.Now()
	onlineMap := make(map[onlineMsgKey]struct{})
	for conversationID, msgs := range allMsg {
		log.ZDebug(ctx, "parse message in one conversation", "conversationID",
			conversationID, "message length", len(msgs.Msgs))
		var insertMessage, selfInsertMessage, othersInsertMessage []*model_struct.LocalChatLog
		var updateMessage []*model_struct.LocalChatLog
		for _, v := range msgs.Msgs {
			log.ZDebug(ctx, "parse message ", "conversationID", conversationID, "msg", v)
			isHistory = utils.GetSwitchFromOptions(v.Options, constant.IsHistory)
			isUnreadCount = utils.GetSwitchFromOptions(v.Options, constant.IsUnreadCount)
			isConversationUpdate = utils.GetSwitchFromOptions(v.Options, constant.IsConversationUpdate)
			isNotPrivate = utils.GetSwitchFromOptions(v.Options, constant.IsNotPrivate)
			isSenderConversationUpdate = utils.GetSwitchFromOptions(v.Options, constant.IsSenderConversationUpdate)
			msg := &sdk_struct.MsgStruct{}
			copier.Copy(msg, v)
			msg.Content = string(v.Content)
			var attachedInfo sdk_struct.AttachedInfoElem
			_ = utils.JsonStringToStruct(v.AttachedInfo, &attachedInfo)
			msg.AttachedInfoElem = &attachedInfo

			msg.Status = constant.MsgStatusSendSuccess
			// msg.IsRead = false
			//De-analyze data
			err := c.msgHandleByContentType(msg)
			if err != nil {
				log.ZError(ctx, "Parsing data error:", err, "type: ", msg.ContentType)
				continue
			}
			//When the message has been marked and deleted by the cloud, it is directly inserted locally without any conversation and message update.
			if msg.Status == constant.MsgStatusHasDeleted {
				insertMessage = append(insertMessage, c.msgStructToLocalChatLog(msg))
				continue
			}
			if !isNotPrivate {
				msg.AttachedInfoElem.IsPrivateChat = true
			}
			if msg.ClientMsgID == "" {
				exceptionMsg = append(exceptionMsg, c.msgStructToLocalErrChatLog(msg))
				continue
			}
			if conversationID == "" {
				log.ZError(ctx, "conversationID is empty", errors.New("conversationID is empty"), "msg", msg)
				continue
			}
			if !isHistory {
				onlineMap[onlineMsgKey{ClientMsgID: v.ClientMsgID, ServerMsgID: v.ServerMsgID}] = struct{}{}
				newMessages = append(newMessages, msg)

			}
			log.ZDebug(ctx, "decode message", "msg", msg)
			if v.SendID == c.loginUserID { //seq
				// Messages sent by myself  //if  sent through  this terminal
				m, err := c.db.GetMessage(ctx, conversationID, msg.ClientMsgID)
				if err == nil {
					log.ZInfo(ctx, "have message", "msg", msg)
					if m.Seq == 0 {
						if !isConversationUpdate {
							msg.Status = constant.MsgStatusFiltered
						}
						updateMessage = append(updateMessage, c.msgStructToLocalChatLog(msg))
					} else {
						exceptionMsg = append(exceptionMsg, c.msgStructToLocalErrChatLog(msg))
					}
				} else {
					log.ZInfo(ctx, "sync message", "msg", msg)
					lc := model_struct.LocalConversation{
						ConversationType:  v.SessionType,
						LatestMsg:         utils.StructToJsonString(msg),
						LatestMsgSendTime: msg.SendTime,
						ConversationID:    conversationID,
					}
					switch v.SessionType {
					case constant.SingleChatType:
						lc.UserID = v.RecvID
					case constant.GroupChatType, constant.SuperGroupChatType:
						lc.GroupID = v.GroupID
					}
					if isConversationUpdate {
						if isSenderConversationUpdate {
							log.ZDebug(ctx, "updateConversation msg", "message", v, "conversation", lc)
							c.updateConversation(&lc, conversationSet)
						}
						newMessages = append(newMessages, msg)
					}
					if isHistory {
						selfInsertMessage = append(selfInsertMessage, c.msgStructToLocalChatLog(msg))
					}
				}
			} else { //Sent by others
				if _, err := c.db.GetMessage(ctx, conversationID, msg.ClientMsgID); err != nil { //Deduplication operation
					lc := model_struct.LocalConversation{
						ConversationType:  v.SessionType,
						LatestMsg:         utils.StructToJsonString(msg),
						LatestMsgSendTime: msg.SendTime,
						ConversationID:    conversationID,
					}
					switch v.SessionType {
					case constant.SingleChatType:
						lc.UserID = v.SendID
						lc.ShowName = msg.SenderNickname
						lc.FaceURL = msg.SenderFaceURL
					case constant.GroupChatType, constant.SuperGroupChatType:
						lc.GroupID = v.GroupID
					case constant.NotificationChatType:
						lc.UserID = v.SendID
					}
					if isUnreadCount {
						//cacheConversation := c.cache.GetConversation(lc.ConversationID)
						if c.maxSeqRecorder.IsNewMsg(conversationID, msg.Seq) {
							isTriggerUnReadCount = true
							lc.UnreadCount = 1
							c.maxSeqRecorder.Incr(conversationID, 1)
						}
					}
					if isConversationUpdate {
						c.updateConversation(&lc, conversationSet)
						newMessages = append(newMessages, msg)
					}
					if isHistory {
						othersInsertMessage = append(othersInsertMessage, c.msgStructToLocalChatLog(msg))
					}

				} else {
					exceptionMsg = append(exceptionMsg, c.msgStructToLocalErrChatLog(msg))
					log.ZWarn(ctx, "Deduplication operation ", nil, "msg", *c.msgStructToLocalErrChatLog(msg))
					msg.Status = constant.MsgStatusFiltered
					msg.ClientMsgID = msg.ClientMsgID + utils.Int64ToString(utils.GetCurrentTimestampByNano())
					othersInsertMessage = append(othersInsertMessage, c.msgStructToLocalChatLog(msg))
				}
			}
		}
		insertMsg[conversationID] = append(insertMessage, c.faceURLAndNicknameHandle(ctx, selfInsertMessage, othersInsertMessage, conversationID)...)
		updateMsg[conversationID] = updateMessage
	}
	list, err := c.db.GetAllConversationListDB(ctx)
	if err != nil {
		log.ZError(ctx, "GetAllConversationListDB", err)
	}
	m := make(map[string]*model_struct.LocalConversation)
	listToMap(list, m)
	log.ZDebug(ctx, "listToMap: ", "local conversation", list, "generated c map", conversationSet)
	c.diff(ctx, m, conversationSet, conversationChangedSet, newConversationSet)
	log.ZInfo(ctx, "trigger map is :", "newConversations", newConversationSet, "changedConversations", conversationChangedSet)

	//seq sync message update
	if err := c.messageController.BatchUpdateMessageList(ctx, updateMsg); err != nil {
		log.ZError(ctx, "sync seq normal message err  :", err)
	}

	//Normal message storage
	_ = c.messageController.BatchInsertMessageList(ctx, insertMsg)

	hList, _ := c.db.GetHiddenConversationList(ctx)
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
			if v.UnreadCount != 0 {
				nc.UnreadCount = v.UnreadCount
			}
			nc.IsNotInGroup = v.IsNotInGroup
			nc.AttachedInfo = v.AttachedInfo
			nc.Ex = v.Ex
			nc.IsMsgDestruct = v.IsMsgDestruct
			nc.MsgDestructTime = v.MsgDestructTime
		}
	}

	for k, v := range newConversationSet {
		if _, ok := phConversationChangedSet[v.ConversationID]; !ok {
			phNewConversationSet[k] = v
		}
	}

	if err := c.db.BatchUpdateConversationList(ctx, append(mapConversationToList(conversationChangedSet), mapConversationToList(phConversationChangedSet)...)); err != nil {
		log.ZError(ctx, "insert changed conversation err :", err)
	}
	//New conversation storage

	if err := c.db.BatchInsertConversationList(ctx, mapConversationToList(phNewConversationSet)); err != nil {
		log.ZError(ctx, "insert new conversation err:", err)
	}
	if c.batchMsgListener() != nil {
		c.batchNewMessages(ctx, newMessages)
	} else {
		c.newMessage(ctx, newMessages, conversationChangedSet, newConversationSet, onlineMap)
	}
	if len(newConversationSet) > 0 {
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.NewConDirect, Args: utils.StructToJsonString(mapConversationToList(newConversationSet))}})
	}
	if len(conversationChangedSet) > 0 {
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.ConChangeDirect, Args: utils.StructToJsonString(mapConversationToList(conversationChangedSet))}})
	}

	if isTriggerUnReadCount {
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.TotalUnreadMessageChanged, Args: ""}})
	}

	for _, msgs := range allMsg {
		for _, msg := range msgs.Msgs {
			if msg.ContentType == constant.Typing {
				c.typing.onNewMsg(ctx, msg)
			}
		}
	}
	log.ZDebug(ctx, "insert msg", "cost time", time.Since(b).Seconds(), "len", len(allMsg))
}

func listToMap(list []*model_struct.LocalConversation, m map[string]*model_struct.LocalConversation) {
	for _, v := range list {
		m[v.ConversationID] = v
	}

}
func (c *Conversation) diff(ctx context.Context, local, generated, cc, nc map[string]*model_struct.LocalConversation) {
	var newConversations []*model_struct.LocalConversation
	for _, v := range generated {
		if localC, ok := local[v.ConversationID]; ok {

			if v.LatestMsgSendTime > localC.LatestMsgSendTime {
				localC.UnreadCount = localC.UnreadCount + v.UnreadCount
				localC.LatestMsg = v.LatestMsg
				localC.LatestMsgSendTime = v.LatestMsgSendTime
				cc[v.ConversationID] = localC
			} else {
				localC.UnreadCount = localC.UnreadCount + v.UnreadCount
				cc[v.ConversationID] = localC
			}

		} else {
			newConversations = append(newConversations, v)
		}
	}
	if err := c.batchAddFaceURLAndName(ctx, newConversations...); err != nil {
		log.ZError(ctx, "batchAddFaceURLAndName err", err, "conversations", newConversations)
	} else {
		for _, v := range newConversations {
			nc[v.ConversationID] = v
		}
	}
}
func (c *Conversation) genConversationGroupAtType(lc *model_struct.LocalConversation, s *sdk_struct.MsgStruct) {
	if s.ContentType == constant.AtText {
		tagMe := utils.IsContain(c.loginUserID, s.AtTextElem.AtUserList)
		tagAll := utils.IsContain(constant.AtAllString, s.AtTextElem.AtUserList)
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

func (c *Conversation) msgStructToLocalErrChatLog(m *sdk_struct.MsgStruct) *model_struct.LocalErrChatLog {
	var lc model_struct.LocalErrChatLog
	copier.Copy(&lc, m)
	if m.SessionType == constant.GroupChatType || m.SessionType == constant.SuperGroupChatType {
		lc.RecvID = m.GroupID
	}
	return &lc
}

func (c *Conversation) tempCacheChatLog(ctx context.Context, messageList []*sdk_struct.MsgStruct) {
	var newMessageList []*model_struct.TempCacheLocalChatLog
	copier.Copy(&newMessageList, &messageList)
	if err := c.db.BatchInsertTempCacheMessageList(ctx, newMessageList); err != nil {
		// log.Error("", "BatchInsertTempCacheMessageList detail err:", err.Error(), len(newMessageList))
		for _, v := range newMessageList {
			err := c.db.InsertTempCacheMessage(ctx, v)
			if err != nil {
				log.ZWarn(ctx, "InsertTempCacheMessage operation", err, "chat err log: ", *v)
			}
		}
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

// func (c *Conversation) doReactionMsgModifier(ctx context.Context, msgReactionList []*sdk_struct.MsgStruct) {
// 	for _, msgStruct := range msgReactionList {
// 		var n server_api_params.ReactionMessageModifierNotification
// 		err := json.Unmarshal([]byte(msgStruct.Content), &n)
// 		if err != nil {
// 			// log.Error("internal", "unmarshal failed err:", err.Error(), *msgStruct)
// 			continue
// 		}
// 		switch n.Operation {
// 		case constant.AddMessageExtensions:
// 			var reactionExtensionList []*sdkws.KeyValue
// 			for _, value := range n.SuccessReactionExtensionList {
// 				reactionExtensionList = append(reactionExtensionList, value)
// 			}
// 			if !(msgStruct.SendID == c.loginUserID && msgStruct.SenderPlatformID == c.platformID) {
// 				c.msgListener.OnRecvMessageExtensionsAdded(n.ClientMsgID, utils.StructToJsonString(reactionExtensionList))
// 			}
// 		case constant.SetMessageExtensions:
// 			err = c.db.GetAndUpdateMessageReactionExtension(ctx, n.ClientMsgID, n.SuccessReactionExtensionList)
// 			if err != nil {
// 				// log.Error("internal", "GetAndUpdateMessageReactionExtension err:", err.Error())
// 				continue
// 			}
// 			var reactionExtensionList []*sdkws.KeyValue
// 			for _, value := range n.SuccessReactionExtensionList {
// 				reactionExtensionList = append(reactionExtensionList, value)
// 			}
// 			if !(msgStruct.SendID == c.loginUserID && msgStruct.SenderPlatformID == c.platformID) {
// 				c.msgListener.OnRecvMessageExtensionsChanged(n.ClientMsgID, utils.StructToJsonString(reactionExtensionList))
// 			}

// 		}
// 		t := model_struct.LocalChatLog{}
// 		t.ClientMsgID = n.ClientMsgID
// 		t.SessionType = n.SessionType
// 		t.IsExternalExtensions = n.IsExternalExtensions
// 		t.IsReact = n.IsReact
// 		t.MsgFirstModifyTime = n.MsgFirstModifyTime
// 		if n.SessionType == constant.GroupChatType || n.SessionType == constant.SuperGroupChatType {
// 			t.RecvID = n.SourceID
// 		}
// 		//todo
// 		err2 := c.db.UpdateMessage(ctx, "", &t)
// 		if err2 != nil {
// 			// log.Error("internal", "unmarshal failed err:", err2.Error(), t)
// 			continue
// 		}
// 	}
// }

func (c *Conversation) doReactionMsgDeleter(ctx context.Context, msgReactionList []*sdk_struct.MsgStruct) {
	// for _, msgStruct := range msgReactionList {
	// 	var n server_api_params.ReactionMessageDeleteNotification
	// 	err := json.Unmarshal([]byte(msgStruct.Content), &n)
	// 	if err != nil {
	// 		// log.Error("internal", "unmarshal failed err:", err.Error(), *msgStruct)
	// 		continue
	// 	}
	// 	err = c.db.DeleteAndUpdateMessageReactionExtension(ctx, n.ClientMsgID, n.SuccessReactionExtensionList)
	// 	if err != nil {
	// 		// log.Error("internal", "GetAndUpdateMessageReactionExtension err:", err.Error())
	// 		continue
	// 	}
	// 	var deleteKeyList []string
	// 	for _, value := range n.SuccessReactionExtensionList {
	// 		deleteKeyList = append(deleteKeyList, value.TypeKey)
	// 	}
	// 	c.msgListener.OnRecvMessageExtensionsDeleted(n.ClientMsgID, utils.StructToJsonString(deleteKeyList))

	// }
}

func (c *Conversation) newMessage(ctx context.Context, newMessagesList sdk_struct.NewMsgList, cc, nc map[string]*model_struct.LocalConversation, onlineMsg map[onlineMsgKey]struct{}) {
	sort.Sort(newMessagesList)
	if c.GetBackground() {
		u, err := c.user.GetSelfUserInfo(ctx)
		if err != nil {
			log.ZWarn(ctx, "GetSelfUserInfo err", err)
			return
		}
		if u.GlobalRecvMsgOpt != constant.ReceiveMessage {
			return
		}
		for _, w := range newMessagesList {
			conversationID := utils.GetConversationIDByMsg(w)
			if v, ok := cc[conversationID]; ok && v.RecvMsgOpt == constant.ReceiveMessage {
				c.msgListener().OnRecvOfflineNewMessage(utils.StructToJsonString(w))
			}
			if v, ok := nc[conversationID]; ok && v.RecvMsgOpt == constant.ReceiveMessage {
				c.msgListener().OnRecvOfflineNewMessage(utils.StructToJsonString(w))
			}
		}
	} else {
		for _, w := range newMessagesList {
			if w.ContentType == constant.Typing {
				continue
			}
			if _, ok := onlineMsg[onlineMsgKey{ClientMsgID: w.ClientMsgID, ServerMsgID: w.ServerMsgID}]; ok {
				c.msgListener().OnRecvOnlineOnlyMessage(utils.StructToJsonString(w))
			} else {
				c.msgListener().OnRecvNewMessage(utils.StructToJsonString(w))
			}
		}
	}
}

func (c *Conversation) batchNewMessages(ctx context.Context, newMessagesList sdk_struct.NewMsgList) {
	sort.Sort(newMessagesList)
	if len(newMessagesList) > 0 {
		c.batchMsgListener().OnRecvNewMessages(utils.StructToJsonString(newMessagesList))
		//if c.IsBackground {
		//	c.batchMsgListener.OnRecvOfflineNewMessages(utils.StructToJsonString(newMessagesList))
		//}
	}

}

func (c *Conversation) doMsgReadState(ctx context.Context, msgReadList []*sdk_struct.MsgStruct) {
	var messageReceiptResp []*sdk_struct.MessageReceipt
	var msgIdList []string
	chrsList := make(map[string][]string)
	var conversationID string

	for _, rd := range msgReadList {
		err := json.Unmarshal([]byte(rd.Content), &msgIdList)
		if err != nil {
			// log.Error("internal", "unmarshal failed, err : ", err.Error())
			return
		}
		var msgIdListStatusOK []string
		for _, v := range msgIdList {
			//m, err := c.db.GetMessage(ctx, v)
			//if err != nil {
			//	log.Error("internal", "GetMessage err:", err, "ClientMsgID", v)
			//	continue
			//}
			//attachInfo := sdk_struct.AttachedInfoElem{}
			//_ = utils.JsonStringToStruct(m.AttachedInfo, &attachInfo)
			//attachInfo.HasReadTime = rd.SendTime
			//m.AttachedInfo = utils.StructToJsonString(attachInfo)
			//m.IsRead = true
			//err = c.db.UpdateMessage(ctx, m)
			//if err != nil {
			//	log.Error("internal", "setMessageHasReadByMsgID err:", err, "ClientMsgID", v)
			//	continue
			//}

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
			conversationID = c.getConversationIDBySessionType(rd.RecvID, constant.SingleChatType)
		} else {
			conversationID = c.getConversationIDBySessionType(rd.SendID, constant.SingleChatType)
		}
		if v, ok := chrsList[conversationID]; ok {
			chrsList[conversationID] = append(v, msgIdListStatusOK...)
		} else {
			chrsList[conversationID] = msgIdListStatusOK
		}
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.ConversationLatestMsgHasRead, Args: chrsList}})
	}
	if len(messageReceiptResp) > 0 {

		// log.Info("internal", "OnRecvC2CReadReceipt: ", utils.StructToJsonString(messageReceiptResp))
		c.msgListener().OnRecvC2CReadReceipt(utils.StructToJsonString(messageReceiptResp))
	}
}

type messageKvList struct {
	ClientMsgID   string                      `json:"clientMsgID"`
	ChangedKvList []*sdk.SingleTypeKeyInfoSum `json:"changedKvList"`
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
	switch msg.ContentType {
	case constant.Text:
		t := sdk_struct.TextElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.TextElem = &t
	case constant.Picture:
		t := sdk_struct.PictureElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.PictureElem = &t
	case constant.Sound:
		t := sdk_struct.SoundElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.SoundElem = &t
	case constant.Video:
		t := sdk_struct.VideoElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.VideoElem = &t
	case constant.File:
		t := sdk_struct.FileElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.FileElem = &t
	case constant.AdvancedText:
		t := sdk_struct.AdvancedTextElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.AdvancedTextElem = &t
	case constant.AtText:
		t := sdk_struct.AtTextElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.AtTextElem = &t
		if err == nil {
			if utils.IsContain(c.loginUserID, msg.AtTextElem.AtUserList) {
				msg.AtTextElem.IsAtSelf = true
			}
		}
	case constant.Location:
		t := sdk_struct.LocationElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.LocationElem = &t
	case constant.Custom:
		fallthrough
	case constant.CustomMsgNotTriggerConversation:
		fallthrough
	case constant.CustomMsgOnlineOnly:
		t := sdk_struct.CustomElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.CustomElem = &t
	case constant.Typing:
		t := sdk_struct.TypingElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.TypingElem = &t
	case constant.Quote:
		t := sdk_struct.QuoteElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.QuoteElem = &t
	case constant.Merger:
		t := sdk_struct.MergeElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.MergeElem = &t
	case constant.Face:
		t := sdk_struct.FaceElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.FaceElem = &t
	case constant.Card:
		t := sdk_struct.CardElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.CardElem = &t
	default:
		t := sdk_struct.NotificationElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.NotificationElem = &t
	}
	msg.Content = ""

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
func (c *Conversation) addFaceURLAndName(ctx context.Context, lc *model_struct.LocalConversation) error {
	switch lc.ConversationType {
	case constant.SingleChatType, constant.NotificationChatType:
		faceUrl, name, err := c.getUserNameAndFaceURL(ctx, lc.UserID)
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

func (c *Conversation) batchAddFaceURLAndName(ctx context.Context, conversations ...*model_struct.LocalConversation) error {
	if len(conversations) == 0 {
		return nil
	}
	var userIDs, groupIDs []string
	for _, conversation := range conversations {
		if conversation.ConversationType == constant.SingleChatType ||
			conversation.ConversationType == constant.NotificationChatType {
			userIDs = append(userIDs, conversation.UserID)
		} else if conversation.ConversationType == constant.SuperGroupChatType {
			groupIDs = append(groupIDs, conversation.GroupID)
		}
	}

	// if userIDs = nil, return nil, nil
	users, err := c.batchGetUserNameAndFaceURL(ctx, userIDs...)
	if err != nil {
		return err
	}

	groups, err := c.full.GetGroupsInfo(ctx, groupIDs...)
	if err != nil {
		return err
	}
	for _, conversation := range conversations {
		if conversation.ConversationType == constant.SingleChatType ||
			conversation.ConversationType == constant.NotificationChatType {
			if v, ok := users[conversation.UserID]; ok {
				conversation.FaceURL = v.FaceURL
				conversation.ShowName = v.Nickname
			} else {
				log.ZWarn(ctx, "user info not found", errors.New("user not found"),
					"userID", conversation.UserID)
			}
		} else if conversation.ConversationType == constant.SuperGroupChatType {
			if v, ok := groups[conversation.GroupID]; ok {
				conversation.FaceURL = v.FaceURL
				conversation.ShowName = v.GroupName
			} else {
				log.ZWarn(ctx, "group info not found", errors.New("group not found"),
					"groupID", conversation.GroupID)
			}

		}
	}
	return nil
}
func (c *Conversation) batchGetUserNameAndFaceURL(ctx context.Context, userIDs ...string) (map[string]*user.BasicInfo,
	error) {
	m := make(map[string]*user.BasicInfo)
	var notCachedUserIDs []string
	var notInFriend []string

	if len(userIDs) == 0 {
		return m, nil
	}

	friendList, err := c.friend.Db().GetFriendInfoList(ctx, userIDs)
	if err != nil {
		log.ZWarn(ctx, "BatchGetUserNameAndFaceURL", err, "userIDs", userIDs)
		notInFriend = userIDs
	} else {
		notInFriend = datautil.SliceSub(userIDs, datautil.Slice(friendList, func(e *model_struct.LocalFriend) string {
			return e.FriendUserID
		}))
	}
	for _, localFriend := range friendList {
		userInfo := &user.BasicInfo{FaceURL: localFriend.FaceURL}
		if localFriend.Remark != "" {
			userInfo.Nickname = localFriend.Remark
		} else {
			userInfo.Nickname = localFriend.Nickname
		}
		m[localFriend.FriendUserID] = userInfo
	}

	for _, userID := range notInFriend {
		if value, ok := c.user.UserBasicCache.Load(userID); ok {
			m[userID] = value
		} else {
			notCachedUserIDs = append(notCachedUserIDs, userID)
		}
	}

	if len(notCachedUserIDs) > 0 {
		users, err := c.user.GetServerUserInfo(ctx, notCachedUserIDs)
		if err != nil {
			return nil, err
		}
		for _, u := range users {
			userInfo := &user.BasicInfo{FaceURL: u.FaceURL, Nickname: u.Nickname}
			m[u.UserID] = userInfo
			c.user.UserBasicCache.Store(u.UserID, userInfo)
		}
	}
	return m, nil
}
func (c *Conversation) getUserNameAndFaceURL(ctx context.Context, userID string) (faceURL, name string, err error) {
	//find in cache
	if value, ok := c.user.UserBasicCache.Load(userID); ok {
		return value.FaceURL, value.Nickname, nil
	}
	//get from local db
	friendInfo, err := c.friend.Db().GetFriendInfoByFriendUserID(ctx, userID)
	if err == nil {
		faceURL = friendInfo.FaceURL
		if friendInfo.Remark != "" {
			name = friendInfo.Remark
		} else {
			name = friendInfo.Nickname
		}
		return faceURL, name, nil
	}
	//get from server db
	users, err := c.user.GetServerUserInfo(ctx, []string{userID})
	if err != nil {
		return "", "", err
	}
	if len(users) == 0 {
		return "", "", sdkerrs.ErrUserIDNotFound.WrapMsg(userID)
	}
	c.user.UserBasicCache.Store(userID, &user.BasicInfo{FaceURL: users[0].FaceURL, Nickname: users[0].Nickname})
	return users[0].FaceURL, users[0].Nickname, nil
}

func (c *Conversation) GetInputStates(ctx context.Context, conversationID string, userID string) ([]int32, error) {
	return c.typing.GetInputStates(conversationID, userID), nil
}

func (c *Conversation) ChangeInputStates(ctx context.Context, conversationID string, focus bool) error {
	return c.typing.ChangeInputStates(ctx, conversationID, focus)
}
