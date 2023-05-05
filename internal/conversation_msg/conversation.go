package conversation_msg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	_ "open_im_sdk/internal/common"
	"open_im_sdk/internal/util"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/syncer"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"sort"
	"strings"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"

	pbConversation "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
	pbMsg "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
)

func (c *Conversation) setConversation(ctx context.Context, apiReq *pbConversation.ModifyConversationFieldReq, localConversation *model_struct.LocalConversation) error {
	apiReq.Conversation.OwnerUserID = c.loginUserID
	apiReq.Conversation.ConversationID = localConversation.ConversationID
	apiReq.Conversation.ConversationType = localConversation.ConversationType
	apiReq.Conversation.UserID = localConversation.UserID
	apiReq.Conversation.GroupID = localConversation.GroupID
	apiReq.UserIDList = []string{c.loginUserID}
	if err := util.ApiPost(ctx, constant.ModifyConversationFieldRouter, apiReq, nil); err != nil {
		return err
	}
	return nil
}

func (c *Conversation) setOneConversationUnread(ctx context.Context, conversationID string, unreadCount int) error {
	apiReq := &pbConversation.ModifyConversationFieldReq{}
	localConversation, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	if localConversation.UnreadCount == 0 {
		return nil
	}
	apiReq.Conversation.UpdateUnreadCountTime = localConversation.LatestMsgSendTime
	apiReq.Conversation.UnreadCount = int32(unreadCount)
	apiReq.FieldType = constant.FieldUnread
	err = c.setConversation(ctx, apiReq, localConversation)
	if err != nil {
		return err
	}
	deleteRows := c.db.DeleteConversationUnreadMessageList(ctx, localConversation.ConversationID, localConversation.LatestMsgSendTime)
	if deleteRows == 0 {
		log.ZError(ctx, "DeleteConversationUnreadMessageList err", nil, "conversationID", localConversation.ConversationID, "latestMsgSendTime", localConversation.LatestMsgSendTime)
	}
	return nil
}

func (c *Conversation) deleteConversation(ctx context.Context, conversationID string) error {
	lc, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	var sourceID string
	switch lc.ConversationType {
	case constant.SingleChatType, constant.NotificationChatType:
		sourceID = lc.UserID
	case constant.GroupChatType, constant.SuperGroupChatType:
		sourceID = lc.GroupID
	}
	if lc.ConversationType == constant.SuperGroupChatType {
		err = c.db.SuperGroupDeleteAllMessage(ctx, lc.GroupID)
		if err != nil {
			return err
		}
	} else {
		//Mark messages related to this conversation for deletion
		err = c.db.UpdateMessageStatusBySourceIDController(ctx, sourceID, constant.MsgStatusHasDeleted, lc.ConversationType)
		if err != nil {
			return err
		}
	}
	//Reset the session information, empty session
	err = c.db.ResetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: "", Action: constant.TotalUnreadMessageChanged, Args: ""}})
	return nil
}

func (c *Conversation) getServerConversationList(ctx context.Context) ([]*model_struct.LocalConversation, error) {
	resp, err := util.CallApi[conversation.GetAllConversationsResp](ctx, constant.GetAllConversationsRouter, conversation.GetAllConversationsReq{OwnerUserID: c.loginUserID})
	if err != nil {
		return nil, err
	}
	return util.Batch(ServerConversationToLocal, resp.Conversations), nil
}

func (c *Conversation) SyncConversations(ctx context.Context) error {
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
func (c *Conversation) SyncConversationUnreadCount(ctx context.Context) error {
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
func (c *Conversation) FixVersionData(ctx context.Context) {
	switch constant.SdkVersion + constant.BigVersion + constant.UpdateVersion {
	case "v2.0.0":
		t := time.Now()
		groupIDList, err := c.db.GetReadDiffusionGroupIDList(ctx)
		if err != nil {
			log.Error("", "GetReadDiffusionGroupIDList failed ", err.Error())
			return
		}
		log.Info("", "fix version data start", groupIDList)
		for _, v := range groupIDList {
			err := c.db.SuperGroupUpdateSpecificContentTypeMessage(ctx, constant.ReactionMessageModifier, v, map[string]interface{}{"status": constant.MsgStatusFiltered})
			if err != nil {
				log.Error("", "SuperGroupUpdateSpecificContentTypeMessage failed ", err.Error())
				continue
			}
			msgList, err := c.db.SuperGroupSearchAllMessageByContentType(ctx, v, constant.ReactionMessageModifier)
			if err != nil {
				log.NewError("internal", "SuperGroupSearchMessageByContentTypeNotOffset failed", v, err.Error())
				continue
			}
			var reactionMsgIDList []string
			for _, value := range msgList {
				var n server_api_params.ReactionMessageModifierNotification
				err := json.Unmarshal([]byte(value.Content), &n)
				if err != nil {
					log.Error("internal", "unmarshal failed err:", err.Error(), *value)
					continue
				}
				reactionMsgIDList = append(reactionMsgIDList, n.ClientMsgID)
			}
			if len(reactionMsgIDList) > 0 {
				err := c.db.SuperGroupUpdateGroupMessageFields(ctx, reactionMsgIDList, v, map[string]interface{}{"is_react": true})
				if err != nil {
					log.Error("internal", "unmarshal failed err:", err.Error(), reactionMsgIDList, v)
					continue
				}
			}

		}
		log.Info("", "fix version data end", groupIDList, "cost time:", time.Since(t))
	}
}

func (c *Conversation) getHistoryMessageList(ctx context.Context, req sdk.GetHistoryMessageListParams, isReverse bool) ([]*sdk_struct.MsgStruct, error) {
	t := time.Now()
	var sourceID string
	var conversationID string
	var startTime int64
	var sessionType int
	var list []*model_struct.LocalChatLog
	var err error
	var messageList sdk_struct.NewMsgList
	var msg sdk_struct.MsgStruct
	var notStartTime bool
	if req.ConversationID != "" {
		conversationID = req.ConversationID
		lc, err := c.db.GetConversation(ctx, conversationID)
		if err != nil {
			return nil, err
		}
		switch lc.ConversationType {
		case constant.SingleChatType, constant.NotificationChatType:
			sourceID = lc.UserID
		case constant.GroupChatType, constant.SuperGroupChatType:
			sourceID = lc.GroupID
			msg.GroupID = lc.GroupID
		}
		sessionType = int(lc.ConversationType)
		if req.StartClientMsgID == "" {
			//startTime = lc.LatestMsgSendTime + TimeOffset
			////startTime = utils.GetCurrentTimestampByMill()
			notStartTime = true
		} else {
			msg.SessionType = lc.ConversationType
			msg.ClientMsgID = req.StartClientMsgID
			m, err := c.db.GetMessageController(ctx, &msg)
			if err != nil {
				return nil, err
			}
			startTime = m.SendTime
		}
	} else {
		if req.UserID == "" {
			newConversationID, newSessionType, err := c.getConversationTypeByGroupID(ctx, req.GroupID)
			if err != nil {
				return nil, err
			}
			sourceID = req.GroupID
			sessionType = int(newSessionType)
			conversationID = newConversationID
			msg.GroupID = req.GroupID
			msg.SessionType = newSessionType
		} else {
			sourceID = req.UserID
			conversationID = utils.GetConversationIDBySessionType(sourceID, constant.SingleChatType)
			sessionType = constant.SingleChatType
		}
		if req.StartClientMsgID == "" {
			//lc, err := c.db.GetConversation(conversationID)
			//if err != nil {
			//	return nil
			//}
			//startTime = lc.LatestMsgSendTime + TimeOffset
			//startTime = utils.GetCurrentTimestampByMill()
			notStartTime = true
		} else {
			msg.ClientMsgID = req.StartClientMsgID
			m, err := c.db.GetMessageController(ctx, &msg)
			if err != nil {
				return nil, err
			}
			startTime = m.SendTime
		}
	}
	log.Debug("", "Assembly parameters cost time", time.Since(t))
	t = time.Now()
	log.Info("", "sourceID:", sourceID, "startTime:", startTime, "count:", req.Count, "not start_time", notStartTime)
	if notStartTime {
		list, err = c.db.GetMessageListNoTimeController(ctx, sourceID, sessionType, req.Count, isReverse)
	} else {
		list, err = c.db.GetMessageListController(ctx, sourceID, sessionType, req.Count, startTime, isReverse)
	}
	log.Debug("", "db cost time", time.Since(t))
	if err != nil {
		return nil, err
	}
	t = time.Now()
	for _, v := range list {
		temp := sdk_struct.MsgStruct{}
		tt := time.Now()
		temp.ClientMsgID = v.ClientMsgID
		temp.ServerMsgID = v.ServerMsgID
		temp.CreateTime = v.CreateTime
		temp.SendTime = v.SendTime
		temp.SessionType = v.SessionType
		temp.SendID = v.SendID
		temp.RecvID = v.RecvID
		temp.MsgFrom = v.MsgFrom
		temp.ContentType = v.ContentType
		temp.SenderPlatformID = v.SenderPlatformID
		temp.SenderNickname = v.SenderNickname
		temp.SenderFaceURL = v.SenderFaceURL
		temp.Content = v.Content
		temp.Seq = v.Seq
		temp.IsRead = v.IsRead
		temp.Status = v.Status
		temp.AttachedInfo = v.AttachedInfo
		temp.Ex = v.Ex
		temp.IsReact = v.IsReact
		temp.IsExternalExtensions = v.IsExternalExtensions
		err := c.msgHandleByContentType(&temp)
		if err != nil {
			log.Error("", "Parsing data error:", err.Error(), temp)
			continue
		}
		log.Debug("", "internal unmarshal cost time", time.Since(tt))

		switch sessionType {
		case constant.GroupChatType:
			fallthrough
		case constant.SuperGroupChatType:
			temp.GroupID = temp.RecvID
			temp.RecvID = c.loginUserID
		}
		messageList = append(messageList, &temp)
	}
	log.Debug("", "unmarshal cost time", time.Since(t))
	t = time.Now()
	if !isReverse {
		sort.Sort(messageList)
	}
	log.Debug("", "sort cost time", time.Since(t))
	return messageList, nil
}
func (c *Conversation) getAdvancedHistoryMessageList(ctx context.Context, req sdk.GetAdvancedHistoryMessageListParams, isReverse bool) (*sdk.GetAdvancedHistoryMessageListCallback, error) {
	t := time.Now()
	var messageListCallback sdk.GetAdvancedHistoryMessageListCallback
	var sourceID string
	var conversationID string
	var startTime int64

	var sessionType int
	var list []*model_struct.LocalChatLog
	var err error
	var messageList sdk_struct.NewMsgList
	var msg sdk_struct.MsgStruct
	var notStartTime bool
	if req.ConversationID != "" {
		conversationID = req.ConversationID
		lc, err := c.db.GetConversation(ctx, conversationID)
		if err != nil {
			//messageListCallback.ErrCode = 100
			//messageListCallback.ErrMsg = "conversation get err"
			return nil, err
		}
		switch lc.ConversationType {
		case constant.SingleChatType, constant.NotificationChatType:
			sourceID = lc.UserID
		case constant.GroupChatType, constant.SuperGroupChatType:
			sourceID = lc.GroupID
			msg.GroupID = lc.GroupID
		}
		sessionType = int(lc.ConversationType)
		if req.StartClientMsgID == "" {
			//startTime = lc.LatestMsgSendTime + TimeOffset
			////startTime = utils.GetCurrentTimestampByMill()
			notStartTime = true
		} else {
			msg.SessionType = lc.ConversationType
			msg.ClientMsgID = req.StartClientMsgID
			m, err := c.db.GetMessageController(ctx, &msg)
			if err != nil {
				return nil, err
			}
			startTime = m.SendTime
		}
	} else {
		if req.UserID == "" {
			newConversationID, newSessionType, err := c.getConversationTypeByGroupID(ctx, req.GroupID)
			if err != nil {
				return nil, err
			}
			sourceID = req.GroupID
			sessionType = int(newSessionType)
			conversationID = newConversationID
			msg.GroupID = req.GroupID
			msg.SessionType = newSessionType
		} else {
			sourceID = req.UserID
			conversationID = utils.GetConversationIDBySessionType(sourceID, constant.SingleChatType)
			sessionType = constant.SingleChatType
		}
		if req.StartClientMsgID == "" {
			//lc, err := c.db.GetConversation(conversationID)
			//if err != nil {
			//	return nil
			//}
			//startTime = lc.LatestMsgSendTime + TimeOffset
			//startTime = utils.GetCurrentTimestampByMill()
			notStartTime = true
		} else {
			msg.ClientMsgID = req.StartClientMsgID
			m, err := c.db.GetMessageController(ctx, &msg)
			if err != nil {
				return nil, err
			}
			startTime = m.SendTime
		}
	}
	log.Debug("", "Assembly parameters cost time", time.Since(t))
	t = time.Now()
	log.Info("", "sourceID:", sourceID, "startTime:", startTime, "count:", req.Count, "not start_time", notStartTime)
	if notStartTime {
		list, err = c.db.GetMessageListNoTimeController(ctx, sourceID, sessionType, req.Count, isReverse)
	} else {
		list, err = c.db.GetMessageListController(ctx, sourceID, sessionType, req.Count, startTime, isReverse)
	}
	log.Error("", "db cost time", time.Since(t), len(list), err, sourceID)
	t = time.Now()
	if err != nil {
		return nil, err
	}
	if isReverse {
		if len(list) < req.Count {
			messageListCallback.IsEnd = true
		}
	} else {
		switch sessionType {
		case constant.SuperGroupChatType:
			if len(list) < req.Count {
				var minSeq int64
				var maxSeq int64
				resp, err := c.SendReqWaitResp(ctx, &server_api_params.GetMaxAndMinSeqReq{UserID: c.loginUserID, GroupIDList: []string{sourceID}}, constant.WSGetNewestSeq, 1, 2, c.loginUserID)
				if err != nil {
					//log.Error(operationID, "SendReqWaitResp failed ", err.Error(), constant.WSGetNewestSeq, 30, c.loginUserID)
				} else {
					var wsSeqResp sdkws.GetMaxAndMinSeqResp
					err = proto.Unmarshal(resp.Data, &wsSeqResp)
					if err != nil {
						//log.Error(operationID, "Unmarshal failed", err.Error())
					} else {
						if value, ok := wsSeqResp.GroupMaxAndMinSeq[sourceID]; ok {
							minSeq = value.MinSeq
							if value.MinSeq == 0 {
								minSeq = 1
							}
							maxSeq = value.MaxSeq
						}
					}
				}
				log.Error("", "from server min seq is", minSeq, maxSeq)
				seq, err := c.db.SuperGroupGetNormalMinSeq(ctx, sourceID)
				if err != nil {
					log.Error("", "SuperGroupGetNormalMinSeq err:", err.Error())
				}
				log.Error("", sourceID+":table min seq is ", seq)
				if seq != 0 {
					if seq <= minSeq {
						messageListCallback.IsEnd = true
					} else {
						seqList := func(seq int64) (seqList []int64) {
							startSeq := int64(seq) - constant.PullMsgNumForReadDiffusion
							if startSeq <= 0 {
								startSeq = 1
							}
							log.Debug("", "pull start is ", startSeq)
							if startSeq < int64(minSeq) {
								startSeq = int64(minSeq)
							}
							for i := startSeq; i < int64(seq); i++ {
								seqList = append(seqList, i)
							}
							log.Debug("", "pull seqList is ", seqList)
							return seqList
						}(seq)
						log.Debug("", "pull seqList is ", seqList, len(seqList))
						if len(seqList) > 0 {
							c.pullMessageAndReGetHistoryMessages(ctx, sourceID, seqList, notStartTime, isReverse, req.Count, sessionType, startTime, &list, &messageListCallback)
						}
					}
				} else {
					//local don't have messages,本地无消息，但是服务器最大消息不为0
					if int64(maxSeq)-int64(minSeq) > 0 {
						messageListCallback.IsEnd = false
					} else {
						messageListCallback.IsEnd = true
					}

				}
			} else if len(list) == req.Count {
				maxSeq, minSeq, haveSeqList := func(messages []*model_struct.LocalChatLog) (max, min int64, seqList []int64) {
					for _, message := range messages {
						if message.Seq != 0 {
							max = message.Seq
							min = message.Seq
							break
						}
					}
					for i := 0; i < len(messages); i++ {
						if messages[i].Seq != 0 {
							seqList = append(seqList, messages[i].Seq)
						}
						if messages[i].Seq > max {
							max = messages[i].Seq

						}
						if messages[i].Seq < min {
							min = messages[i].Seq
						}
					}
					return max, min, seqList
				}(list)
				log.Debug("", "get message from local db max seq:", maxSeq, "minSeq:", minSeq, "haveSeqList:", haveSeqList, "length:", len(haveSeqList))
				if maxSeq != 0 && minSeq != 0 {
					successiveSeqList := func(max, min int64) (seqList []int64) {
						for i := min; i <= max; i++ {
							seqList = append(seqList, i)
						}
						return seqList
					}(maxSeq, minSeq)
					lostSeqList := utils.DifferenceSubset(successiveSeqList, haveSeqList)
					lostSeqListLength := len(lostSeqList)
					log.Debug("", "get lost seqList is :", lostSeqList, "length:", lostSeqListLength)
					if lostSeqListLength > 0 {
						var pullSeqList []int64
						if lostSeqListLength <= constant.PullMsgNumForReadDiffusion {
							pullSeqList = lostSeqList
						} else {
							pullSeqList = lostSeqList[lostSeqListLength-constant.PullMsgNumForReadDiffusion : lostSeqListLength]
						}
						c.pullMessageAndReGetHistoryMessages(ctx, sourceID, pullSeqList, notStartTime, isReverse, req.Count, sessionType, startTime, &list, &messageListCallback)
					} else {
						if req.LastMinSeq != 0 {
							var thisMaxSeq int64
							for i := 0; i < len(list); i++ {
								if list[i].Seq != 0 && thisMaxSeq == 0 {
									thisMaxSeq = list[i].Seq
								}
								if list[i].Seq > thisMaxSeq {
									thisMaxSeq = list[i].Seq
								}
							}
							log.Debug("", "get lost LastMinSeq is :", req.LastMinSeq, "thisMaxSeq is :", thisMaxSeq)
							if thisMaxSeq != 0 {
								if thisMaxSeq+1 != req.LastMinSeq {
									startSeq := int64(req.LastMinSeq) - constant.PullMsgNumForReadDiffusion
									if startSeq <= int64(thisMaxSeq) {
										startSeq = int64(thisMaxSeq) + 1
									}
									successiveSeqList := func(max, min int64) (seqList []int64) {
										for i := min; i <= max; i++ {
											seqList = append(seqList, i)
										}
										return seqList
									}(req.LastMinSeq-1, startSeq)
									log.Debug("", "get lost successiveSeqList is :", successiveSeqList, len(successiveSeqList))
									if len(successiveSeqList) > 0 {
										c.pullMessageAndReGetHistoryMessages(ctx, sourceID, successiveSeqList, notStartTime, isReverse, req.Count, sessionType, startTime, &list, &messageListCallback)
									}
								}

							}

						}
					}

				}
			}
		default:
			if len(list) < req.Count {
				messageListCallback.IsEnd = true
			}
		}

	}

	log.Debug("", "pull cost time", time.Since(t))
	t = time.Now()
	var thisMinSeq int64
	for _, v := range list {
		if v.Seq != 0 && thisMinSeq == 0 {
			thisMinSeq = v.Seq
		}
		if v.Seq < thisMinSeq {
			thisMinSeq = v.Seq
		}
		temp := sdk_struct.MsgStruct{}
		tt := time.Now()
		temp.ClientMsgID = v.ClientMsgID
		temp.ServerMsgID = v.ServerMsgID
		temp.CreateTime = v.CreateTime
		temp.SendTime = v.SendTime
		temp.SessionType = v.SessionType
		temp.SendID = v.SendID
		temp.RecvID = v.RecvID
		temp.MsgFrom = v.MsgFrom
		temp.ContentType = v.ContentType
		temp.SenderPlatformID = v.SenderPlatformID
		temp.SenderNickname = v.SenderNickname
		temp.SenderFaceURL = v.SenderFaceURL
		temp.Content = v.Content
		temp.Seq = v.Seq
		temp.IsRead = v.IsRead
		temp.Status = v.Status
		temp.AttachedInfo = v.AttachedInfo
		temp.Ex = v.Ex
		temp.IsReact = v.IsReact
		temp.IsExternalExtensions = v.IsExternalExtensions
		err := c.msgHandleByContentType(&temp)
		if err != nil {
			log.Error("", "Parsing data error:", err.Error(), temp)
			continue
		}
		log.Debug("", "internal unmarshal cost time", time.Since(tt))

		switch sessionType {
		case constant.GroupChatType:
			fallthrough
		case constant.SuperGroupChatType:
			temp.GroupID = temp.RecvID
			temp.RecvID = c.loginUserID
		}
		messageList = append(messageList, &temp)
	}
	log.Debug("", "unmarshal cost time", time.Since(t))
	t = time.Now()
	if !isReverse {
		sort.Sort(messageList)
	}
	log.Debug("", "sort cost time", time.Since(t))
	messageListCallback.MessageList = messageList
	messageListCallback.LastMinSeq = thisMinSeq
	return &messageListCallback, nil
}

func (c *Conversation) getAdvancedHistoryMessageList2(callback open_im_sdk_callback.Base, req sdk.GetAdvancedHistoryMessageListParams, operationID string, isReverse bool) sdk.GetAdvancedHistoryMessageListCallback {
	t := time.Now()
	var messageListCallback sdk.GetAdvancedHistoryMessageListCallback
	var sourceID string
	var conversationID string
	var startTime int64

	var sessionType int
	var list []*model_struct.LocalChatLog
	var err error
	var messageList sdk_struct.NewMsgList
	var msg sdk_struct.MsgStruct
	var notStartTime bool
	if req.ConversationID != "" {
		conversationID = req.ConversationID
		lc, err := c.db.GetConversation(context.Background(), conversationID)
		if err != nil {
			messageListCallback.ErrCode = 100
			messageListCallback.ErrMsg = "conversation get err"
			return messageListCallback
		}
		switch lc.ConversationType {
		case constant.SingleChatType, constant.NotificationChatType:
			sourceID = lc.UserID
		case constant.GroupChatType, constant.SuperGroupChatType:
			sourceID = lc.GroupID
			msg.GroupID = lc.GroupID
		}
		sessionType = int(lc.ConversationType)
		if req.StartClientMsgID == "" {
			//startTime = lc.LatestMsgSendTime + TimeOffset
			////startTime = utils.GetCurrentTimestampByMill()
			notStartTime = true
		} else {
			msg.SessionType = lc.ConversationType
			msg.ClientMsgID = req.StartClientMsgID
			m, err := c.db.GetMessageController(context.Background(), &msg)
			common.CheckDBErrCallback(callback, err, operationID)
			startTime = m.SendTime
		}
	} else {
		if req.UserID == "" {
			newConversationID, newSessionType, err := c.getConversationTypeByGroupID(context.Background(), req.GroupID)
			common.CheckDBErrCallback(callback, err, operationID)
			sourceID = req.GroupID
			sessionType = int(newSessionType)
			conversationID = newConversationID
			msg.GroupID = req.GroupID
			msg.SessionType = newSessionType
		} else {
			sourceID = req.UserID
			conversationID = utils.GetConversationIDBySessionType(sourceID, constant.SingleChatType)
			sessionType = constant.SingleChatType
		}
		if req.StartClientMsgID == "" {
			//lc, err := c.db.GetConversation(conversationID)
			//if err != nil {
			//	return nil
			//}
			//startTime = lc.LatestMsgSendTime + TimeOffset
			//startTime = utils.GetCurrentTimestampByMill()
			notStartTime = true
		} else {
			msg.ClientMsgID = req.StartClientMsgID
			m, err := c.db.GetMessageController(context.Background(), &msg)
			common.CheckDBErrCallback(callback, err, operationID)
			startTime = m.SendTime
		}
	}
	log.Debug(operationID, "Assembly parameters cost time", time.Since(t))
	t = time.Now()
	log.Info(operationID, "sourceID:", sourceID, "startTime:", startTime, "count:", req.Count, "not start_time", notStartTime)
	if notStartTime {
		list, err = c.db.GetMessageListNoTimeController(context.Background(), sourceID, sessionType, req.Count, isReverse)
	} else {
		list, err = c.db.GetMessageListController(context.Background(), sourceID, sessionType, req.Count, startTime, isReverse)
	}
	log.Error(operationID, "db cost time", time.Since(t), len(list), err, sourceID)
	t = time.Now()
	common.CheckDBErrCallback(callback, err, operationID)
	if sessionType == constant.SuperGroupChatType {
		rawMessageLength := len(list)
		if rawMessageLength < req.Count {
			maxSeq, minSeq, lostSeqListLength := c.messageBlocksInternalContinuityCheck(sourceID, notStartTime, isReverse, req.Count, sessionType, startTime, &list, &messageListCallback, operationID)
			_ = c.messageBlocksBetweenContinuityCheck(req.LastMinSeq, maxSeq, sourceID, notStartTime, isReverse, req.Count, sessionType, startTime, &list, &messageListCallback, operationID)
			if minSeq == 1 && lostSeqListLength == 0 {
				messageListCallback.IsEnd = true
			} else {
				c.messageBlocksEndContinuityCheck(nil, minSeq, sourceID, notStartTime, isReverse, req.Count, sessionType, startTime, &list, &messageListCallback)
			}
		} else {
			maxSeq, _, _ := c.messageBlocksInternalContinuityCheck(sourceID, notStartTime, isReverse, req.Count, sessionType, startTime, &list, &messageListCallback, operationID)
			c.messageBlocksBetweenContinuityCheck(req.LastMinSeq, maxSeq, sourceID, notStartTime, isReverse, req.Count, sessionType, startTime, &list, &messageListCallback, operationID)

		}

	}
	//if len(list) < req.Count && sessionType == constant.SuperGroupChatType {
	//
	//} else if len(list) == req.Count && sessionType == constant.SuperGroupChatType {
	//	if maxSeq != 0 && minSeq != 0 {
	//		successiveSeqList := func(max, min uint32) (seqList []uint32) {
	//			for i := min; i <= max; i++ {
	//				seqList = append(seqList, i)
	//			}
	//			return seqList
	//		}(maxSeq, minSeq)
	//		lostSeqList := utils.DifferenceSubset(successiveSeqList, haveSeqList)
	//		lostSeqListLength := len(lostSeqList)
	//		log.Debug(operationID, "get lost seqList is :", lostSeqList, "length:", lostSeqListLength)
	//		if lostSeqListLength > 0 {
	//			var pullSeqList []uint32
	//			if lostSeqListLength <= constant.PullMsgNumForReadDiffusion {
	//				pullSeqList = lostSeqList
	//			} else {
	//				pullSeqList = lostSeqList[lostSeqListLength-constant.PullMsgNumForReadDiffusion : lostSeqListLength]
	//			}
	//			c.pullMessageAndReGetHistoryMessages(sourceID, pullSeqList, notStartTime, isReverse, req.Count, sessionType, startTime, &list, &messageListCallback, operationID)
	//		} else {
	//			if req.LastMinSeq != 0 {
	//				var thisMaxSeq uint32
	//				for i := 0; i < len(list); i++ {
	//					if list[i].Seq != 0 && thisMaxSeq == 0 {
	//						thisMaxSeq = list[i].Seq
	//					}
	//					if list[i].Seq > thisMaxSeq {
	//						thisMaxSeq = list[i].Seq
	//					}
	//				}
	//				log.Debug(operationID, "get lost LastMinSeq is :", req.LastMinSeq, "thisMaxSeq is :", thisMaxSeq)
	//				if thisMaxSeq != 0 {
	//					if thisMaxSeq+1 != req.LastMinSeq {
	//						startSeq := int64(req.LastMinSeq) - constant.PullMsgNumForReadDiffusion
	//						if startSeq <= int64(thisMaxSeq) {
	//							startSeq = int64(thisMaxSeq) + 1
	//						}
	//						successiveSeqList := func(max, min uint32) (seqList []uint32) {
	//							for i := min; i <= max; i++ {
	//								seqList = append(seqList, i)
	//							}
	//							return seqList
	//						}(req.LastMinSeq-1, uint32(startSeq))
	//						log.Debug(operationID, "get lost successiveSeqList is :", successiveSeqList, len(successiveSeqList))
	//						if len(successiveSeqList) > 0 {
	//							c.pullMessageAndReGetHistoryMessages(sourceID, successiveSeqList, notStartTime, isReverse, req.Count, sessionType, startTime, &list, &messageListCallback, operationID)
	//						}
	//					}
	//
	//				}
	//
	//			}
	//		}
	//
	//	}
	//}
	log.Debug(operationID, "pull cost time", time.Since(t))
	t = time.Now()
	var thisMinSeq int64
	for _, v := range list {
		if v.Seq != 0 && thisMinSeq == 0 {
			thisMinSeq = v.Seq
		}
		if v.Seq < thisMinSeq {
			thisMinSeq = v.Seq
		}
		temp := sdk_struct.MsgStruct{}
		tt := time.Now()
		temp.ClientMsgID = v.ClientMsgID
		temp.ServerMsgID = v.ServerMsgID
		temp.CreateTime = v.CreateTime
		temp.SendTime = v.SendTime
		temp.SessionType = v.SessionType
		temp.SendID = v.SendID
		temp.RecvID = v.RecvID
		temp.MsgFrom = v.MsgFrom
		temp.ContentType = v.ContentType
		temp.SenderPlatformID = v.SenderPlatformID
		temp.SenderNickname = v.SenderNickname
		temp.SenderFaceURL = v.SenderFaceURL
		temp.Content = v.Content
		temp.Seq = v.Seq
		temp.IsRead = v.IsRead
		temp.Status = v.Status
		temp.AttachedInfo = v.AttachedInfo
		temp.Ex = v.Ex
		err := c.msgHandleByContentType(&temp)
		if err != nil {
			log.Error(operationID, "Parsing data error:", err.Error(), temp)
			continue
		}
		log.Debug(operationID, "internal unmarshal cost time", time.Since(tt))

		switch sessionType {
		case constant.GroupChatType:
			fallthrough
		case constant.SuperGroupChatType:
			temp.GroupID = temp.RecvID
			temp.RecvID = c.loginUserID
		}
		messageList = append(messageList, &temp)
	}
	log.Debug(operationID, "unmarshal cost time", time.Since(t))
	t = time.Now()
	if !isReverse {
		sort.Sort(messageList)
	}
	log.Debug(operationID, "sort cost time", time.Since(t))
	messageListCallback.MessageList = messageList
	if thisMinSeq == 0 {
		thisMinSeq = req.LastMinSeq
	}
	messageListCallback.LastMinSeq = thisMinSeq
	return messageListCallback
}

func (c *Conversation) revokeOneMessage(ctx context.Context, req *sdk_struct.MsgStruct) error {
	var recvID, groupID string
	var localMessage model_struct.LocalChatLog
	var lc model_struct.LocalConversation
	var conversationID string
	message, err := c.db.GetMessageController(ctx, req)
	if err != nil {
		return err
	}
	if message.Status != constant.MsgStatusSendSuccess {
		return errors.New("only send success message can be revoked")
	}
	if message.SendID != c.loginUserID {
		return errors.New("only you send message can be revoked")
	}
	//Send message internally
	switch req.SessionType {
	case constant.SingleChatType:
		recvID = req.RecvID
		conversationID = utils.GetConversationIDBySessionType(recvID, constant.SingleChatType)
	case constant.GroupChatType:
		groupID = req.GroupID
		conversationID = utils.GetConversationIDBySessionType(groupID, constant.GroupChatType)
	case constant.SuperGroupChatType:
		groupID = req.GroupID
		conversationID = utils.GetConversationIDBySessionType(groupID, constant.SuperGroupChatType)
	default:
		return errors.New("SessionType err")
	}
	req.Content = message.ClientMsgID
	req.ClientMsgID = utils.GetMsgID(message.SendID)
	req.ContentType = constant.Revoke
	req.SendTime = 0
	req.CreateTime = utils.GetCurrentTimestampByMill()
	options := make(map[string]bool, 5)
	utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
	utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
	resp, err := c.InternalSendMessage(ctx, req, recvID, groupID, &server_api_params.OfflinePushInfo{}, false, options)
	if err != nil {
		return err
	}
	req.ServerMsgID = resp.ServerMsgID
	req.SendTime = resp.SendTime
	req.Status = constant.MsgStatusSendSuccess
	msgStructToLocalChatLog(&localMessage, req)
	err = c.db.InsertMessageController(ctx, &localMessage)
	if err != nil {
		log.Error("", "inset into chat log err", localMessage, req)
	}
	err = c.db.UpdateColumnsMessageController(ctx, req.Content, groupID, req.SessionType, map[string]interface{}{"status": constant.MsgStatusRevoked})
	if err != nil {
		log.Error("", "update revoke message err", localMessage, req)
	}
	lc.LatestMsg = utils.StructToJsonString(req)
	lc.LatestMsgSendTime = req.SendTime
	lc.ConversationID = conversationID
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: lc.ConversationID, Action: constant.AddConOrUpLatMsg, Args: lc}, c.GetCh())
	return nil
}
func (c *Conversation) newRevokeOneMessage(ctx context.Context, req *sdk_struct.MsgStruct) error {
	var recvID, groupID string
	var localMessage model_struct.LocalChatLog
	var revokeMessage sdk_struct.MessageRevoked
	var lc model_struct.LocalConversation
	var conversationID string
	message, err := c.db.GetMessageController(ctx, req)
	if err != nil {
		return err
	}
	if message.Status != constant.MsgStatusSendSuccess {
		return errors.New("only send success message can be revoked")
	}
	s := sdk_struct.MsgStruct{}
	err = c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.AdvancedRevoke)
	if err != nil {
		return err
	}
	revokeMessage.ClientMsgID = message.ClientMsgID
	revokeMessage.RevokerID = c.loginUserID
	revokeMessage.RevokeTime = utils.GetCurrentTimestampBySecond()
	revokeMessage.RevokerNickname = s.SenderNickname
	revokeMessage.SourceMessageSendTime = message.SendTime
	revokeMessage.SessionType = message.SessionType
	revokeMessage.SourceMessageSendID = message.SendID
	revokeMessage.SourceMessageSenderNickname = message.SenderNickname
	revokeMessage.Seq = message.Seq
	revokeMessage.Ex = message.Ex
	//Send message internally
	switch message.SessionType {
	case constant.SingleChatType:
		if message.SendID != c.loginUserID {
			return errors.New("only you send message can be revoked")
		}
		recvID = message.RecvID
		conversationID = utils.GetConversationIDBySessionType(recvID, constant.SingleChatType)
	case constant.GroupChatType:
		if message.SendID != c.loginUserID {
			ownerID, adminIDList, err := c.group.GetGroupOwnerIDAndAdminIDList(ctx, message.RecvID)
			if err != nil {
				return err
			}
			if c.loginUserID == ownerID {
				revokeMessage.RevokerRole = constant.GroupOwner
			} else if utils.IsContain(c.loginUserID, adminIDList) {
				if utils.IsContain(message.SendID, adminIDList) || message.SendID == ownerID {
					return errors.New("you do not have this permission")
				} else {
					revokeMessage.RevokerRole = constant.GroupAdmin
				}

			} else {
				return errors.New("you do not have this permission")
			}
		}
		groupID = message.RecvID
		conversationID = utils.GetConversationIDBySessionType(groupID, constant.GroupChatType)
	case constant.SuperGroupChatType:
		if message.SendID != c.loginUserID {
			ownerID, adminIDList, err := c.group.GetGroupOwnerIDAndAdminIDList(ctx, message.RecvID)
			if err != nil {
				return err
			}
			if c.loginUserID == ownerID {
				revokeMessage.RevokerRole = constant.GroupOwner
			} else if utils.IsContain(c.loginUserID, adminIDList) {
				if utils.IsContain(message.SendID, adminIDList) || message.SendID == ownerID {
					return errors.New("you do not have this permission")
				} else {
					revokeMessage.RevokerRole = constant.GroupAdmin
				}

			} else {
				return errors.New("you do not have this permission")
			}
		}
		groupID = message.RecvID
		conversationID = utils.GetConversationIDBySessionType(groupID, constant.SuperGroupChatType)
	default:
		return errors.New("SessionType err")
	}
	s.Content = utils.StructToJsonString(revokeMessage)
	options := make(map[string]bool, 5)
	utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
	utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
	resp, err := c.InternalSendMessage(ctx, &s, recvID, groupID, &server_api_params.OfflinePushInfo{}, false, options)
	if err != nil {
		return err
	}
	s.ServerMsgID = resp.ServerMsgID
	s.SendTime = message.SendTime //New message takes the old place
	s.Status = constant.MsgStatusSendSuccess
	msgStructToLocalChatLog(&localMessage, &s)
	err = c.db.InsertMessageController(ctx, &localMessage)
	if err != nil {
		log.Error("", "inset into chat log err", localMessage, s)
	}
	err = c.db.UpdateColumnsMessageController(ctx, message.ClientMsgID, groupID, message.SessionType, map[string]interface{}{"status": constant.MsgStatusRevoked})
	if err != nil {
		log.Error("", "update revoke message err", localMessage, message, err.Error())
	}
	s.SendTime = resp.SendTime
	lc.LatestMsg = utils.StructToJsonString(s)
	lc.LatestMsgSendTime = s.SendTime
	lc.ConversationID = conversationID
	s.GroupID = groupID
	s.RecvID = recvID
	c.newRevokeMessage(ctx, []*sdk_struct.MsgStruct{&s})
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: lc.ConversationID, Action: constant.AddConOrUpLatMsg, Args: lc}, c.GetCh())
	return nil
}

func (c *Conversation) typingStatusUpdate(ctx context.Context, recvID, msgTip string) error {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Typing)
	if err != nil {
		return err
	}
	s.Content = msgTip
	options := make(map[string]bool, 6)
	utils.SetSwitchFromOptions(options, constant.IsHistory, false)
	utils.SetSwitchFromOptions(options, constant.IsPersistent, false)
	utils.SetSwitchFromOptions(options, constant.IsSenderSync, false)
	utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsSenderConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
	utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
	_, err = c.InternalSendMessage(ctx, &s, recvID, "", &server_api_params.OfflinePushInfo{}, true, options)
	return err

}

func (c *Conversation) markC2CMessageAsRead(ctx context.Context, msgIDList []string, userID string) error {
	var localMessage model_struct.LocalChatLog
	var newMessageIDList []string
	messages, err := c.db.GetMultipleMessage(ctx, msgIDList)
	if err != nil {
		return err
	}
	for _, v := range messages {
		if v.IsRead == false && v.ContentType < constant.NotificationBegin && v.SendID != c.loginUserID {
			newMessageIDList = append(newMessageIDList, v.ClientMsgID)
		}
	}
	if len(newMessageIDList) == 0 {
		return errors.New("message has been marked read or sender is yourself or notification message not support")
	}
	conversationID := utils.GetConversationIDBySessionType(userID, constant.SingleChatType)
	s := sdk_struct.MsgStruct{}
	err = c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.HasReadReceipt)
	if err != nil {
		return err
	}
	s.Content = utils.StructToJsonString(newMessageIDList)
	options := make(map[string]bool, 5)
	utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsSenderConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
	utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
	//If there is an error, the coroutine ends, so judgment is not  required
	resp, err := c.InternalSendMessage(ctx, &s, userID, "", &server_api_params.OfflinePushInfo{}, false, options)
	if err != nil {
		return err
	}
	s.ServerMsgID = resp.ServerMsgID
	s.SendTime = resp.SendTime
	s.Status = constant.MsgStatusFiltered
	msgStructToLocalChatLog(&localMessage, &s)
	err = c.db.InsertMessage(ctx, &localMessage)
	if err != nil {
		log.Error("", "inset into chat log err", localMessage, s, err.Error())
	}

	err2 := c.db.UpdateSingleMessageHasRead(ctx, userID, newMessageIDList)
	if err2 != nil {
		log.Error("", "update message has read error", newMessageIDList, userID, err2.Error())
	}
	newMessages, err3 := c.db.GetMultipleMessage(ctx, newMessageIDList)
	if err3 != nil {
		log.Error("", "get messages error", newMessageIDList, userID, err3.Error())
	}
	for _, v := range newMessages {
		attachInfo := sdk_struct.AttachedInfoElem{}
		_ = utils.JsonStringToStruct(v.AttachedInfo, &attachInfo)
		attachInfo.HasReadTime = s.SendTime
		v.AttachedInfo = utils.StructToJsonString(attachInfo)
		err = c.db.UpdateMessage(ctx, v)
		if err != nil {
			log.Error("internal", "setMessageHasReadByMsgID err:", err, "ClientMsgID", v)
			continue
		}
	}
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.UpdateLatestMessageChange}, c.GetCh())
	return nil
	//_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.ch)
}
func (c *Conversation) markGroupMessageAsRead(ctx context.Context, msgIDList []string, groupID string) error {
	conversationID, conversationType, err := c.getConversationTypeByGroupID(ctx, groupID)
	if err != nil {
		return err
	}
	if len(msgIDList) == 0 {
		_ = c.setOneConversationUnread(ctx, conversationID, 0)
		_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.UnreadCountSetZero}, c.GetCh())
		_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())
		return nil
	}
	var localMessage model_struct.LocalChatLog
	allUserMessage := make(map[string][]string, 3)
	messages, err := c.db.GetMultipleMessageController(ctx, msgIDList, groupID, conversationType)
	if err != nil {
		return err
	}
	for _, v := range messages {
		log.Debug("", "get group info is test2", v.ClientMsgID, v.SessionType)
		if v.IsRead == false && v.ContentType < constant.NotificationBegin && v.SendID != c.loginUserID {
			if msgIDList, ok := allUserMessage[v.SendID]; ok {
				msgIDList = append(msgIDList, v.ClientMsgID)
				allUserMessage[v.SendID] = msgIDList
			} else {
				allUserMessage[v.SendID] = []string{v.ClientMsgID}
			}
		}
	}
	if len(allUserMessage) == 0 {
		return errors.New("message has been marked read or sender is yourself or notification message not support")
	}

	for userID, list := range allUserMessage {
		s := sdk_struct.MsgStruct{}
		s.GroupID = groupID
		err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.GroupHasReadReceipt)
		if err != nil {
			return err
		}
		s.Content = utils.StructToJsonString(list)
		options := make(map[string]bool, 5)
		utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
		utils.SetSwitchFromOptions(options, constant.IsSenderConversationUpdate, false)
		utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
		utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
		//If there is an error, the coroutine ends, so judgment is not  required
		resp, err := c.InternalSendMessage(ctx, &s, userID, "", &server_api_params.OfflinePushInfo{}, false, options)
		if err != nil {
			return err
		}
		s.ServerMsgID = resp.ServerMsgID
		s.SendTime = resp.SendTime
		s.Status = constant.MsgStatusFiltered
		msgStructToLocalChatLog(&localMessage, &s)
		err = c.db.InsertMessageController(ctx, &localMessage)
		if err != nil {
			log.Error(
				"", "inset into chat log err", localMessage, s, err.Error())
		}
		log.Debug("", "get group info is test3", list, conversationType)
		err2 := c.db.UpdateGroupMessageHasReadController(ctx, list, groupID, conversationType)
		if err2 != nil {
			log.Error("", "update message has read err", list, userID, err2.Error())
		}
	}
	return nil
}

//	func (c *Conversation) markMessageAsReadByConID(callback open_im_sdk_callback.Base, msgIDList sdk.MarkMessageAsReadByConIDParams, conversationID, operationID string) {
//		var localMessage db.LocalChatLog
//		var newMessageIDList []string
//		messages, err := c.db.GetMultipleMessage(msgIDList)
//		common.CheckDBErrCallback(callback, err, operationID)
//		for _, v := range messages {
//			if v.IsRead == false && v.ContentType < constant.NotificationBegin && v.SendID != c.loginUserID {
//				newMessageIDList = append(newMessageIDList, v.ClientMsgID)
//			}
//		}
//		if len(newMessageIDList) == 0 {
//			common.CheckAnyErrCallback(callback, 201, errors.New("message has been marked read or sender is yourself"), operationID)
//		}
//		conversationID := utils.GetConversationIDBySessionType(userID, constant.SingleChatType)
//		s := sdk_struct.MsgStruct{}
//		c.initBasicInfo(&s, constant.UserMsgType, constant.HasReadReceipt, operationID)
//		s.Content = utils.StructToJsonString(newMessageIDList)
//		options := make(map[string]bool, 5)
//		utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
//		utils.SetSwitchFromOptions(options, constant.IsSenderConversationUpdate, false)
//		utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
//		utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
//		//If there is an error, the coroutine ends, so judgment is not  required
//		resp, _ := c.InternalSendMessage(callback, &s, userID, "", operationID, &server_api_params.OfflinePushInfo{}, false, options)
//		s.ServerMsgID = resp.ServerMsgID
//		s.SendTime = resp.SendTime
//		s.Status = constant.MsgStatusFiltered
//		msgStructToLocalChatLog(&localMessage, &s)
//		err = c.db.InsertMessage(&localMessage)
//		if err != nil {
//			log.Error(operationID, "inset into chat log err", localMessage, s, err.Error())
//		}
//		err2 := c.db.UpdateMessageHasRead(userID, newMessageIDList, constant.SingleChatType)
//		if err2 != nil {
//			log.Error(operationID, "update message has read error", newMessageIDList, userID, err2.Error())
//		}
//		_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.UpdateLatestMessageChange}, c.ch)
//		//_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.ch)
//	}
func (c *Conversation) insertMessageToLocalStorage(ctx context.Context, s *model_struct.LocalChatLog) error {
	return c.db.InsertMessageController(ctx, s)

}

func (c *Conversation) clearGroupHistoryMessage(ctx context.Context, groupID string) error {
	_, sessionType, err := c.getConversationTypeByGroupID(ctx, groupID)
	if err != nil {
		return err
	}
	conversationID := utils.GetConversationIDBySessionType(groupID, int(sessionType))
	switch sessionType {
	case constant.SuperGroupChatType:
		err = c.db.SuperGroupDeleteAllMessage(ctx, groupID)
		if err != nil {
			return err
		}
	default:
		err = c.db.UpdateMessageStatusBySourceIDController(ctx, groupID, constant.MsgStatusHasDeleted, sessionType)
		if err != nil {
			return err
		}
	}

	err = c.db.ClearConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())
	return nil

}

func (c *Conversation) clearC2CHistoryMessage(ctx context.Context, userID string) error {
	conversationID := utils.GetConversationIDBySessionType(userID, constant.SingleChatType)
	err := c.db.UpdateMessageStatusBySourceID(ctx, userID, constant.MsgStatusHasDeleted, constant.SingleChatType)
	if err != nil {
		return err
	}
	err = c.db.ClearConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())
	return nil
}

func (c *Conversation) deleteMessageFromSvr(ctx context.Context, s *sdk_struct.MsgStruct) error {
	seq, err := c.db.GetMsgSeqByClientMsgIDController(ctx, s)
	if err != nil {
		return err
	}
	switch s.SessionType {
	case constant.SingleChatType, constant.GroupChatType:
		var apiReq pbMsg.DelMsgsReq
		apiReq.Seqs = utils.Uint32ListConvert([]uint32{seq})
		apiReq.UserID = c.loginUserID
		return util.ApiPost(ctx, constant.DeleteMsgRouter, &apiReq, nil)
	case constant.SuperGroupChatType:
		var apiReq pbMsg.DelSuperGroupMsgReq
		apiReq.UserID = c.loginUserID
		apiReq.GroupID = s.GroupID
		return util.ApiPost(ctx, constant.DeleteSuperGroupMsgRouter, &apiReq, nil)

	}
	return errors.New("session type error")

}

func (c *Conversation) clearMessageFromSvr(ctx context.Context) error {
	var apiReq pbMsg.ClearMsgReq
	apiReq.UserID = c.loginUserID
	err := util.ApiPost(ctx, constant.ClearMsgRouter, &apiReq, nil)
	if err != nil {
		return err
	}
	groupIDList, err := c.full.GetReadDiffusionGroupIDList(ctx)
	if err != nil {
		return err
	}
	var superGroupApiReq pbMsg.DelSuperGroupMsgReq
	superGroupApiReq.UserID = c.loginUserID
	for _, v := range groupIDList {
		superGroupApiReq.GroupID = v
		err := util.ApiPost(ctx, constant.DeleteSuperGroupMsgRouter, &superGroupApiReq, nil)
		if err != nil {
			//log.
		}

	}
	return nil
}

func (c *Conversation) deleteMessageFromLocalStorage(ctx context.Context, s *sdk_struct.MsgStruct) error {
	var conversation model_struct.LocalConversation
	var latestMsg sdk_struct.MsgStruct
	var conversationID string
	var sourceID string
	chatLog := model_struct.LocalChatLog{ClientMsgID: s.ClientMsgID, Status: constant.MsgStatusHasDeleted, SessionType: s.SessionType}

	switch s.SessionType {
	case constant.GroupChatType:
		conversationID = utils.GetConversationIDBySessionType(s.GroupID, constant.GroupChatType)
		sourceID = s.GroupID
	case constant.SingleChatType:
		if s.SendID != c.loginUserID {
			conversationID = utils.GetConversationIDBySessionType(s.SendID, constant.SingleChatType)
			sourceID = s.SendID
		} else {
			conversationID = utils.GetConversationIDBySessionType(s.RecvID, constant.SingleChatType)
			sourceID = s.RecvID
		}
	case constant.SuperGroupChatType:
		conversationID = utils.GetConversationIDBySessionType(s.GroupID, constant.SuperGroupChatType)
		sourceID = s.GroupID
		chatLog.RecvID = s.GroupID
	}
	err := c.db.UpdateMessageController(ctx, &chatLog)
	if err != nil {
		return err
	}
	LocalConversation, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	err = utils.JsonStringToStruct(LocalConversation.LatestMsg, &latestMsg)
	if err != nil {
		return err
	}

	if s.ClientMsgID == latestMsg.ClientMsgID { //If the deleted message is the latest message of the conversation, update the latest message of the conversation
		list, err := c.db.GetMessageListNoTimeController(ctx, sourceID, int(s.SessionType), 1, false)
		if err != nil {
			return err
		}
		conversation.ConversationID = conversationID
		if list == nil {
			conversation.LatestMsg = ""
			conversation.LatestMsgSendTime = s.SendTime
		} else {
			copier.Copy(&latestMsg, list[0])
			err := c.msgConvert(&latestMsg)
			if err != nil {
				log.Error("", "Parsing data error:", err.Error(), latestMsg)
			}
			conversation.LatestMsg = utils.StructToJsonString(latestMsg)
			conversation.LatestMsgSendTime = latestMsg.SendTime
		}
		err = c.db.UpdateColumnsConversation(ctx, conversation.ConversationID, map[string]interface{}{"latest_msg_send_time": conversation.LatestMsgSendTime, "latest_msg": conversation.LatestMsg})
		if err != nil {
			log.Error("internal", "updateConversationLatestMsgModel err: ", err)
		} else {
			_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())
		}
	}
	return nil
}
func (c *Conversation) judgeMultipleSubString(keywordList []string, main string, keywordListMatchType int) bool {
	if len(keywordList) == 0 {
		return true
	}
	if keywordListMatchType == constant.KeywordMatchOr {
		for _, v := range keywordList {
			if utils.KMP(main, v) {
				return true
			}
		}
		return false
	} else {
		for _, v := range keywordList {
			if !utils.KMP(main, v) {
				return false
			}
		}
	}
	return true
}

func (c *Conversation) searchLocalMessages(ctx context.Context, searchParam *sdk.SearchLocalMessagesParams) (*sdk.SearchLocalMessagesCallback, error) {
	var r sdk.SearchLocalMessagesCallback
	var conversationID, sourceID string
	var startTime, endTime int64
	var list []*model_struct.LocalChatLog
	conversationMap := make(map[string]*sdk.SearchByConversationResult, 10)
	var err error

	if searchParam.SearchTimePosition == 0 {
		endTime = utils.GetCurrentTimestampBySecond()
	} else {
		endTime = searchParam.SearchTimePosition
	}
	if searchParam.SearchTimePeriod != 0 {
		startTime = endTime - searchParam.SearchTimePeriod
	}
	startTime = utils.UnixSecondToTime(startTime).UnixNano() / 1e6
	endTime = utils.UnixSecondToTime(endTime).UnixNano() / 1e6
	if len(searchParam.KeywordList) == 0 && len(searchParam.MessageTypeList) == 0 {
		return nil, errors.New("keywordlist and messageTypelist all null")
	}
	if searchParam.ConversationID != "" {
		if searchParam.PageIndex < 1 || searchParam.Count < 1 {
			return nil, errors.New("page or count is null")
		}
		offset := (searchParam.PageIndex - 1) * searchParam.Count
		localConversation, err := c.db.GetConversation(ctx, searchParam.ConversationID)
		if err != nil {
			return nil, err
		}
		switch localConversation.ConversationType {
		case constant.SingleChatType:
			sourceID = localConversation.UserID
		case constant.GroupChatType:
			sourceID = localConversation.GroupID
		case constant.SuperGroupChatType:
			sourceID = localConversation.GroupID
		}
		if len(searchParam.MessageTypeList) != 0 && len(searchParam.KeywordList) == 0 {
			list, err = c.db.SearchMessageByContentTypeController(ctx, searchParam.MessageTypeList, sourceID, startTime, endTime, int(localConversation.ConversationType), offset, searchParam.Count)
		} else {
			newContentTypeList := func(list []int) (result []int) {
				for _, v := range list {
					if utils.IsContainInt(v, SearchContentType) {
						result = append(result, v)
					}
				}
				return result
			}(searchParam.MessageTypeList)
			if len(newContentTypeList) == 0 {
				newContentTypeList = SearchContentType
			}
			list, err = c.db.SearchMessageByKeywordController(ctx, newContentTypeList, searchParam.KeywordList, searchParam.KeywordListMatchType, sourceID, startTime, endTime, int(localConversation.ConversationType), offset, searchParam.Count)
		}
	} else {
		//Comprehensive search, search all
		if len(searchParam.MessageTypeList) == 0 {
			searchParam.MessageTypeList = SearchContentType
		}
		list, err = c.db.SearchMessageByContentTypeAndKeywordController(ctx, searchParam.MessageTypeList, searchParam.KeywordList, searchParam.KeywordListMatchType, startTime, endTime)
	}
	if err != nil {
		return nil, err
	}
	//localChatLogToMsgStruct(&messageList, list)

	//log.Debug("hahh",utils.KMP("SSSsdf3434","s"))
	//log.Debug("hahh",utils.KMP("SSSsdf3434","g"))
	//log.Debug("hahh",utils.KMP("SSSsdf3434","3434"))
	//log.Debug("hahh",utils.KMP("SSSsdf3434","F3434"))
	//log.Debug("hahh",utils.KMP("SSSsdf3434","SDF3"))
	log.Debug("", "get raw data length is", len(list))
	for _, v := range list {
		temp := sdk_struct.MsgStruct{}
		temp.ClientMsgID = v.ClientMsgID
		temp.ServerMsgID = v.ServerMsgID
		temp.CreateTime = v.CreateTime
		temp.SendTime = v.SendTime
		temp.SessionType = v.SessionType
		temp.SendID = v.SendID
		temp.RecvID = v.RecvID
		temp.MsgFrom = v.MsgFrom
		temp.ContentType = v.ContentType
		temp.SenderPlatformID = v.SenderPlatformID
		temp.SenderNickname = v.SenderNickname
		temp.SenderFaceURL = v.SenderFaceURL
		temp.Content = v.Content
		temp.Seq = v.Seq
		temp.IsRead = v.IsRead
		temp.Status = v.Status
		temp.AttachedInfo = v.AttachedInfo
		temp.Ex = v.Ex
		err := c.msgHandleByContentType(&temp)
		if err != nil {
			log.Error("", "Parsing data error:", err.Error(), temp)
			continue
		}
		if temp.ContentType == constant.File && !c.judgeMultipleSubString(searchParam.KeywordList, temp.FileElem.FileName, searchParam.KeywordListMatchType) {
			continue
		}
		if temp.ContentType == constant.AtText && !c.judgeMultipleSubString(searchParam.KeywordList, temp.AtElem.Text, searchParam.KeywordListMatchType) {
			continue
		}
		switch temp.SessionType {
		case constant.SingleChatType:
			if temp.SendID == c.loginUserID {
				conversationID = utils.GetConversationIDBySessionType(temp.RecvID, constant.SingleChatType)
			} else {
				conversationID = utils.GetConversationIDBySessionType(temp.SendID, constant.SingleChatType)
			}
		case constant.GroupChatType:
			temp.GroupID = temp.RecvID
			temp.RecvID = c.loginUserID
			conversationID = utils.GetConversationIDBySessionType(temp.GroupID, constant.GroupChatType)
		case constant.SuperGroupChatType:
			temp.GroupID = temp.RecvID
			temp.RecvID = c.loginUserID
			conversationID = utils.GetConversationIDBySessionType(temp.GroupID, constant.SuperGroupChatType)
		}
		if oldItem, ok := conversationMap[conversationID]; !ok {
			searchResultItem := sdk.SearchByConversationResult{}
			localConversation, err := c.db.GetConversation(ctx, conversationID)
			if err != nil {
				log.Error("", "get conversation err ", err.Error(), conversationID)
				continue
			}
			searchResultItem.ConversationID = conversationID
			searchResultItem.FaceURL = localConversation.FaceURL
			searchResultItem.ShowName = localConversation.ShowName
			searchResultItem.ConversationType = localConversation.ConversationType
			searchResultItem.MessageList = append(searchResultItem.MessageList, &temp)
			searchResultItem.MessageCount++
			conversationMap[conversationID] = &searchResultItem
		} else {
			oldItem.MessageCount++
			oldItem.MessageList = append(oldItem.MessageList, &temp)
			conversationMap[conversationID] = oldItem
		}
	}
	for _, v := range conversationMap {
		r.SearchResultItems = append(r.SearchResultItems, v)
		r.TotalCount += v.MessageCount

	}
	return &r, nil
}

func (c *Conversation) delMsgBySeq(seqList []uint32) error {
	var SPLIT = 1000
	for i := 0; i < len(seqList)/SPLIT; i++ {
		if err := c.delMsgBySeqSplit(seqList[i*SPLIT : (i+1)*SPLIT]); err != nil {
			return utils.Wrap(err, "")
		}
	}
	return nil
}

func (c *Conversation) delMsgBySeqSplit(seqList []uint32) error {
	var req server_api_params.DelMsgListReq
	req.SeqList = seqList
	req.OperationID = utils.OperationIDGenerator()
	req.OpUserID = c.loginUserID
	req.UserID = c.loginUserID
	operationID := req.OperationID

	resp, err := c.Ws.SendReqWaitResp(context.Background(), &req, constant.WsDelMsg, 30, 5, c.loginUserID)
	if err != nil {
		return utils.Wrap(err, "SendReqWaitResp failed")
	}
	var delResp server_api_params.DelMsgListResp
	err = proto.Unmarshal(resp.Data, &delResp)
	if err != nil {
		log.Error(operationID, "Unmarshal failed ", err.Error())
		return utils.Wrap(err, "Unmarshal failed")
	}
	return nil
}

// old WS method
//func (c *Conversation) deleteMessageFromSvr(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, operationID string) {
//	seq, err := c.db.GetMsgSeqByClientMsgID(s.ClientMsgID)
//	common.CheckDBErrCallback(callback, err, operationID)
//	if seq == 0 {
//		err = errors.New("seq == 0 ")
//		common.CheckArgsErrCallback(callback, err, operationID)
//	}
//	seqList := []uint32{seq}
//	err = c.delMsgBySeq(seqList)
//	common.CheckArgsErrCallback(callback, err, operationID)
//}

func (c *Conversation) deleteConversationAndMsgFromSvr(ctx context.Context, conversationID string) error {
	local, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	var seqList []uint32
	switch local.ConversationType {
	case constant.SingleChatType, constant.NotificationChatType:
		peerUserID := local.UserID
		if peerUserID != c.loginUserID {
			seqList, err = c.db.GetMsgSeqListByPeerUserID(ctx, peerUserID)
		} else {
			seqList, err = c.db.GetMsgSeqListBySelfUserID(ctx, c.loginUserID)
		}
		log.ZDebug(ctx, utils.GetSelfFuncName(), "seqList: ", seqList)
		if err != nil {
			return err
		}
	case constant.GroupChatType:
		groupID := local.GroupID
		seqList, err = c.db.GetMsgSeqListByGroupID(ctx, groupID)
		log.ZDebug(ctx, utils.GetSelfFuncName(), "seqList: ", seqList)
		if err != nil {
			return err
		}
	case constant.SuperGroupChatType:
		var apiReq pbMsg.DelSuperGroupMsgReq
		apiReq.UserID = c.loginUserID
		apiReq.GroupID = local.GroupID
		if err := util.ApiPost(ctx, constant.DeleteSuperGroupMsgRouter, &apiReq, nil); err != nil {
			return err
		}

	}
	var apiReq pbMsg.DelMsgsReq
	apiReq.UserID = c.loginUserID
	apiReq.Seqs = utils.Uint32ListConvert(seqList)
	return util.ApiPost(ctx, constant.DeleteMsgRouter, &apiReq, nil)

}

func (c *Conversation) deleteAllMsgFromLocal(ctx context.Context) error {
	//log.NewInfo(operationID, utils.GetSelfFuncName())
	err := c.db.DeleteAllMessage(ctx)
	if err != nil {
		return err
	}
	groupIDList, err := c.full.GetReadDiffusionGroupIDList(ctx)
	if err != nil {
		return err
	}
	for _, v := range groupIDList {
		err = c.db.SuperGroupDeleteAllMessage(ctx, v)
		if err != nil {
			//log.Error(operationID, "SuperGroupDeleteAllMessage err", err.Error())
			continue
		}
	}
	err = c.db.ClearAllConversation(ctx)
	if err != nil {
		return err
	}
	conversationList, err := c.db.GetAllConversationListDB(ctx)
	if err != nil {
		return err
	}
	var cidList []string
	for _, conversation := range conversationList {
		cidList = append(cidList, conversation.ConversationID)
	}
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{Action: constant.ConChange, Args: cidList}, c.GetCh())
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{"", constant.TotalUnreadMessageChanged, ""}})
	return nil

}

func isContainMessageReaction(reactionType int, list []*sdk_struct.ReactionElem) (bool, *sdk_struct.ReactionElem) {
	for _, v := range list {
		if v.Type == reactionType {
			return true, v
		}
	}
	return false, nil
}
func isContainUserReactionElem(useID string, list []*sdk_struct.UserReactionElem) (bool, *sdk_struct.UserReactionElem) {
	for _, v := range list {
		if v.UserID == useID {
			return true, v
		}
	}
	return false, nil
}

func DeleteUserReactionElem(a []*sdk_struct.UserReactionElem, userID string) []*sdk_struct.UserReactionElem {
	j := 0
	for _, v := range a {
		if v.UserID != userID {
			a[j] = v
			j++
		}
	}
	return a[:j]
}
func (c *Conversation) setMessageReactionExtensions(ctx context.Context, s *sdk_struct.MsgStruct, req sdk.SetMessageReactionExtensionsParams) ([]*server_api_params.ExtensionResult, error) {
	message, err := c.db.GetMessageController(ctx, s)
	if err != nil {
		return nil, err
	}
	if message.Status != constant.MsgStatusSendSuccess {
		return nil, errors.New("only send success message can modify reaction extensions")
	}
	if message.SessionType != constant.SuperGroupChatType {
		return nil, errors.New("currently only support super group message")

	}
	extendMsg, _ := c.db.GetMessageReactionExtension(ctx, message.ClientMsgID)
	temp := make(map[string]*server_api_params.KeyValue)
	_ = json.Unmarshal(extendMsg.LocalReactionExtensions, &temp)
	reqTemp := make(map[string]*server_api_params.KeyValue)
	for _, v := range req {
		if value, ok := temp[v.TypeKey]; ok {
			v.LatestUpdateTime = value.LatestUpdateTime
		}
		reqTemp[v.TypeKey] = v
	}
	var sourceID string
	switch message.SessionType {
	case constant.SingleChatType:
		if message.SendID == c.loginUserID {
			sourceID = message.RecvID
		} else {
			sourceID = message.SendID
		}
	case constant.NotificationChatType:
		sourceID = message.RecvID
	case constant.GroupChatType, constant.SuperGroupChatType:
		sourceID = message.RecvID
	}
	var apiReq server_api_params.SetMessageReactionExtensionsReq
	apiReq.IsReact = message.IsReact
	apiReq.ClientMsgID = message.ClientMsgID
	apiReq.SourceID = sourceID
	apiReq.SessionType = message.SessionType
	apiReq.IsExternalExtensions = message.IsExternalExtensions
	apiReq.ReactionExtensionList = reqTemp
	apiReq.OperationID = ""
	apiReq.MsgFirstModifyTime = message.MsgFirstModifyTime
	resp, err := util.CallApi[server_api_params.ApiResult](ctx, constant.SetMessageReactionExtensionsRouter, &apiReq)
	if err != nil {
		return nil, err
	}
	var msg model_struct.LocalChatLogReactionExtensions
	msg.ClientMsgID = message.ClientMsgID
	resultKeyMap := make(map[string]*sdkws.KeyValue)
	for _, v := range resp.Result {
		if v.ErrCode == 0 {
			temp := new(sdkws.KeyValue)
			temp.TypeKey = v.TypeKey
			temp.Value = v.Value
			temp.LatestUpdateTime = v.LatestUpdateTime
			resultKeyMap[v.TypeKey] = temp
		}
	}
	err = c.db.GetAndUpdateMessageReactionExtension(ctx, message.ClientMsgID, resultKeyMap)
	if err != nil {
		log.Error("", "GetAndUpdateMessageReactionExtension err:", err.Error())
	}
	if !message.IsReact {
		message.IsReact = resp.IsReact
		message.MsgFirstModifyTime = resp.MsgFirstModifyTime
		err = c.db.UpdateMessageController(ctx, message)
		if err != nil {
			log.Error("", "UpdateMessageController err:", err.Error(), message)

		}
	}
	return resp.Result, nil
}
func (c *Conversation) addMessageReactionExtensions(ctx context.Context, s *sdk_struct.MsgStruct, req sdk.AddMessageReactionExtensionsParams) ([]*server_api_params.ExtensionResult, error) {
	message, err := c.db.GetMessageController(ctx, s)
	if err != nil {
		return nil, err
	}
	if message.Status != constant.MsgStatusSendSuccess || message.Seq == 0 {
		return nil, errors.New("only send success message can modify reaction extensions")
	}
	reqTemp := make(map[string]*server_api_params.KeyValue)
	extendMsg, err := c.db.GetMessageReactionExtension(ctx, message.ClientMsgID)
	if err == nil && extendMsg != nil {
		temp := make(map[string]*server_api_params.KeyValue)
		_ = json.Unmarshal(extendMsg.LocalReactionExtensions, &temp)
		for _, v := range req {
			if value, ok := temp[v.TypeKey]; ok {
				v.LatestUpdateTime = value.LatestUpdateTime
			}
			reqTemp[v.TypeKey] = v
		}
	} else {
		for _, v := range req {
			reqTemp[v.TypeKey] = v
		}
	}
	var sourceID string
	switch message.SessionType {
	case constant.SingleChatType:
		if message.SendID == c.loginUserID {
			sourceID = message.RecvID
		} else {
			sourceID = message.SendID
		}
	case constant.NotificationChatType:
		sourceID = message.RecvID
	case constant.GroupChatType, constant.SuperGroupChatType:
		sourceID = message.RecvID
	}
	var apiReq server_api_params.AddMessageReactionExtensionsReq
	apiReq.IsReact = message.IsReact
	apiReq.ClientMsgID = message.ClientMsgID
	apiReq.SourceID = sourceID
	apiReq.SessionType = message.SessionType
	apiReq.IsExternalExtensions = message.IsExternalExtensions
	apiReq.ReactionExtensionList = reqTemp
	apiReq.OperationID = ""
	apiReq.MsgFirstModifyTime = message.MsgFirstModifyTime
	apiReq.Seq = message.Seq

	resp, err := util.CallApi[server_api_params.ApiResult](ctx, constant.AddMessageReactionExtensionsRouter, &apiReq)
	if err != nil {
		return nil, err
	}
	log.Debug("", "api return:", message.IsReact, resp)
	if !message.IsReact {
		message.IsReact = resp.IsReact
		message.MsgFirstModifyTime = resp.MsgFirstModifyTime
		err = c.db.UpdateMessageController(ctx, message)
		if err != nil {
			log.Error("", "UpdateMessageController err:", err.Error(), message)
		}
	}
	return resp.Result, nil
}

func (c *Conversation) deleteMessageReactionExtensions(ctx context.Context, s *sdk_struct.MsgStruct, req sdk.DeleteMessageReactionExtensionsParams) ([]*server_api_params.ExtensionResult, error) {
	message, err := c.db.GetMessageController(ctx, s)
	if err != nil {
		return nil, err
	}
	if message.Status != constant.MsgStatusSendSuccess {
		return nil, errors.New("only send success message can modify reaction extensions")
	}
	if message.SessionType != constant.SuperGroupChatType {
		return nil, errors.New("currently only support super group message")

	}
	extendMsg, _ := c.db.GetMessageReactionExtension(ctx, message.ClientMsgID)
	temp := make(map[string]*server_api_params.KeyValue)
	_ = json.Unmarshal(extendMsg.LocalReactionExtensions, &temp)
	var reqTemp []*server_api_params.KeyValue
	for _, v := range req {
		if value, ok := temp[v]; ok {
			var tt server_api_params.KeyValue
			tt.LatestUpdateTime = value.LatestUpdateTime
			tt.TypeKey = v
			reqTemp = append(reqTemp, &tt)
		}
	}
	var sourceID string
	switch message.SessionType {
	case constant.SingleChatType:
		if message.SendID == c.loginUserID {
			sourceID = message.RecvID
		} else {
			sourceID = message.SendID
		}
	case constant.NotificationChatType:
		sourceID = message.RecvID
	case constant.GroupChatType, constant.SuperGroupChatType:
		sourceID = message.RecvID
	}
	var apiReq server_api_params.DeleteMessageReactionExtensionsReq
	apiReq.ClientMsgID = message.ClientMsgID
	apiReq.SourceID = sourceID
	apiReq.SessionType = message.SessionType
	apiReq.ReactionExtensionList = reqTemp
	apiReq.OperationID = ""
	apiReq.IsExternalExtensions = message.IsExternalExtensions
	apiReq.MsgFirstModifyTime = message.MsgFirstModifyTime
	resp, err := util.CallApi[server_api_params.ApiResult](ctx, constant.AddMessageReactionExtensionsRouter, &apiReq)
	if err != nil {
		return nil, err
	}
	var msg model_struct.LocalChatLogReactionExtensions
	msg.ClientMsgID = message.ClientMsgID
	resultKeyMap := make(map[string]*sdkws.KeyValue)
	for _, v := range resp.Result {
		if v.ErrCode == 0 {
			temp := new(sdkws.KeyValue)
			temp.TypeKey = v.TypeKey
			resultKeyMap[v.TypeKey] = temp
		}
	}
	err = c.db.DeleteAndUpdateMessageReactionExtension(ctx, message.ClientMsgID, resultKeyMap)
	if err != nil {
		log.Error("", "GetAndUpdateMessageReactionExtension err:", err.Error())
	}
	return resp.Result, nil
}

type syncReactionExtensionParams struct {
	MessageList         []*model_struct.LocalChatLog
	SessionType         int32
	SourceID            string
	IsExternalExtension bool
	ExtendMessageList   []*model_struct.LocalChatLogReactionExtensions
	TypeKeyList         []string
}

func (c *Conversation) getMessageListReactionExtensions(ctx context.Context, messageList []*sdk_struct.MsgStruct) ([]*server_api_params.SingleMessageExtensionResult, error) {
	if len(messageList) == 0 {
		return nil, errors.New("message list is null")
	}
	var msgIDList []string
	var sourceID string
	var sessionType int32
	var isExternalExtension bool
	for _, msgStruct := range messageList {
		switch msgStruct.SessionType {
		case constant.SingleChatType:
			if msgStruct.SendID == c.loginUserID {
				sourceID = msgStruct.RecvID
			} else {
				sourceID = msgStruct.SendID
			}
		case constant.NotificationChatType:
			sourceID = msgStruct.RecvID
		case constant.GroupChatType, constant.SuperGroupChatType:
			sourceID = msgStruct.GroupID
		}
		sessionType = msgStruct.SessionType
		msgIDList = append(msgIDList, msgStruct.ClientMsgID)
	}
	isExternalExtension = c.IsExternalExtensions
	localMessageList, err := c.db.GetMultipleMessageController(ctx, msgIDList, sourceID, sessionType)
	if err != nil {
		return nil, err
	}
	for _, v := range localMessageList {
		if v.IsReact != true {
			return nil, errors.New("have not reaction message in message list:" + v.ClientMsgID)
		}
	}
	var result server_api_params.GetMessageListReactionExtensionsResp
	extendMessage, _ := c.db.GetMultipleMessageReactionExtension(ctx, msgIDList)
	for _, v := range extendMessage {
		var singleResult server_api_params.SingleMessageExtensionResult
		temp := make(map[string]*sdkws.KeyValue)
		_ = json.Unmarshal(v.LocalReactionExtensions, &temp)
		singleResult.ClientMsgID = v.ClientMsgID
		singleResult.ReactionExtensionList = temp
		result = append(result, &singleResult)
	}
	args := syncReactionExtensionParams{}
	args.MessageList = localMessageList
	args.SourceID = sourceID
	args.SessionType = sessionType
	args.ExtendMessageList = extendMessage
	args.IsExternalExtension = isExternalExtension
	_ = common.TriggerCmdSyncReactionExtensions(common.SyncReactionExtensionsNode{
		OperationID: "",
		Action:      constant.SyncMessageListReactionExtensions,
		Args:        args,
	}, c.GetCh())
	return result, nil

}

//	func (c *Conversation) getMessageListSomeReactionExtensions(callback open_im_sdk_callback.Base, messageList []*sdk_struct.MsgStruct, keyList []string, operationID string) server_api_params.GetMessageListReactionExtensionsResp {
//		if len(messageList) == 0 {
//			common.CheckAnyErrCallback(callback, 201, errors.New("message list is null"), operationID)
//		}
//		var msgIDList []string
//		var sourceID string
//		var sessionType int32
//		var isExternalExtension bool
//		for _, msgStruct := range messageList {
//			switch msgStruct.SessionType {
//			case constant.SingleChatType:
//				if msgStruct.SendID == c.loginUserID {
//					sourceID = msgStruct.RecvID
//				} else {
//					sourceID = msgStruct.SendID
//				}
//			case constant.NotificationChatType:
//				sourceID = msgStruct.RecvID
//			case constant.GroupChatType, constant.SuperGroupChatType:
//				sourceID = msgStruct.GroupID
//			}
//			sessionType = msgStruct.SessionType
//			isExternalExtension = msgStruct.IsExternalExtensions
//			msgIDList = append(msgIDList, msgStruct.ClientMsgID)
//		}
//		localMessageList, err := c.db.GetMultipleMessageController(msgIDList, sourceID, sessionType)
//		common.CheckDBErrCallback(callback, err, operationID)
//		var result server_api_params.GetMessageListReactionExtensionsResp
//		extendMsgs, _ := c.db.GetMultipleMessageReactionExtension(msgIDList)
//		for _, v := range extendMsgs {
//			var singleResult server_api_params.SingleMessageExtensionResult
//			temp := make(map[string]*server_api_params.KeyValue)
//			_ = json.Unmarshal(v.LocalReactionExtensions, &temp)
//			for s, _ := range temp {
//				if !utils.IsContain(s, keyList) {
//					delete(temp, s)
//				}
//			}
//			singleResult.ClientMsgID = v.ClientMsgID
//			singleResult.ReactionExtensionList = temp
//			result = append(result, &singleResult)
//		}
//		args := syncReactionExtensionParams{}
//		args.MessageList = localMessageList
//		args.SourceID = sourceID
//		args.TypeKeyList = keyList
//		args.SessionType = sessionType
//		args.ExtendMessageList = extendMsgs
//		args.IsExternalExtension = isExternalExtension
//		_ = common.TriggerCmdSyncReactionExtensions(common.SyncReactionExtensionsNode{
//			OperationID: operationID,
//			Action:      constant.SyncMessageListReactionExtensions,
//			Args:        args,
//		}, c.GetCh())
//		return result
//	}
//
//	func (c *Conversation) setTypeKeyInfo(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, typeKey, ex string, isCanRepeat bool, operationID string) []*server_api_params.ExtensionResult {
//		message, err := c.db.GetMessageController(s)
//		common.CheckDBErrCallback(callback, err, operationID)
//		if message.Status != constant.MsgStatusSendSuccess {
//			common.CheckAnyErrCallback(callback, 201, errors.New("only send success message can modify reaction extensions"), operationID)
//		}
//		extendMsg, _ := c.db.GetMessageReactionExtension(message.ClientMsgID)
//		temp := make(map[string]*server_api_params.KeyValue)
//		_ = json.Unmarshal(extendMsg.LocalReactionExtensions, &temp)
//		var flag bool
//		var isContainSelfK string
//		var dbIsCanRepeat bool
//		var deletedKeyValue server_api_params.KeyValue
//		var maxTypeKey string
//		var maxTypeKeyValue server_api_params.KeyValue
//		reqTemp := make(map[string]*server_api_params.KeyValue)
//		for k, v := range temp {
//			if strings.HasPrefix(k, typeKey) {
//				flag = true
//				singleTypeKeyInfo := new(sdk.SingleTypeKeyInfo)
//				_ = json.Unmarshal([]byte(v.Value), singleTypeKeyInfo)
//				if _, ok := singleTypeKeyInfo.InfoList[c.loginUserID]; ok {
//					isContainSelfK = k
//					dbIsCanRepeat = singleTypeKeyInfo.IsCanRepeat
//					delete(singleTypeKeyInfo.InfoList, c.loginUserID)
//					singleTypeKeyInfo.Counter--
//					deletedKeyValue.TypeKey = v.TypeKey
//					deletedKeyValue.Value = utils.StructToJsonString(singleTypeKeyInfo)
//					deletedKeyValue.LatestUpdateTime = v.LatestUpdateTime
//				}
//				if k > maxTypeKey {
//					maxTypeKey = k
//					maxTypeKeyValue = *v
//				}
//			}
//		}
//		if !flag {
//			if len(temp) >= 300 {
//				common.CheckAnyErrCallback(callback, 202, errors.New("number of keys only can support 300"), operationID)
//			}
//			singleTypeKeyInfo := new(sdk.SingleTypeKeyInfo)
//			singleTypeKeyInfo.TypeKey = getIndexTypeKey(typeKey, 0)
//			singleTypeKeyInfo.Counter = 1
//			singleTypeKeyInfo.IsCanRepeat = isCanRepeat
//			singleTypeKeyInfo.Index = 0
//			userInfo := new(sdk.Info)
//			userInfo.UserID = c.loginUserID
//			userInfo.Ex = ex
//			singleTypeKeyInfo.InfoList[c.loginUserID] = userInfo
//			keyValue := new(server_api_params.KeyValue)
//			keyValue.TypeKey = singleTypeKeyInfo.TypeKey
//			keyValue.Value = utils.StructToJsonString(singleTypeKeyInfo)
//			reqTemp[singleTypeKeyInfo.TypeKey] = keyValue
//		} else {
//			if isContainSelfK != "" && !dbIsCanRepeat {
//				//删除操作
//				reqTemp[isContainSelfK] = &deletedKeyValue
//			} else {
//				singleTypeKeyInfo := new(sdk.SingleTypeKeyInfo)
//				_ = json.Unmarshal([]byte(maxTypeKeyValue.Value), singleTypeKeyInfo)
//				userInfo := new(sdk.Info)
//				userInfo.UserID = c.loginUserID
//				userInfo.Ex = ex
//				singleTypeKeyInfo.Counter++
//				singleTypeKeyInfo.InfoList[c.loginUserID] = userInfo
//				maxTypeKeyValue.Value = utils.StructToJsonString(singleTypeKeyInfo)
//				data, _ := json.Marshal(maxTypeKeyValue)
//				if len(data) > 1000 { //单key超过了1kb
//					if len(temp) >= 300 {
//						common.CheckAnyErrCallback(callback, 202, errors.New("number of keys only can support 300"), operationID)
//					}
//					newSingleTypeKeyInfo := new(sdk.SingleTypeKeyInfo)
//					newSingleTypeKeyInfo.TypeKey = getIndexTypeKey(typeKey, singleTypeKeyInfo.Index+1)
//					newSingleTypeKeyInfo.Counter = 1
//					newSingleTypeKeyInfo.IsCanRepeat = singleTypeKeyInfo.IsCanRepeat
//					newSingleTypeKeyInfo.Index = singleTypeKeyInfo.Index + 1
//					userInfo := new(sdk.Info)
//					userInfo.UserID = c.loginUserID
//					userInfo.Ex = ex
//					newSingleTypeKeyInfo.InfoList[c.loginUserID] = userInfo
//					keyValue := new(server_api_params.KeyValue)
//					keyValue.TypeKey = newSingleTypeKeyInfo.TypeKey
//					keyValue.Value = utils.StructToJsonString(newSingleTypeKeyInfo)
//					reqTemp[singleTypeKeyInfo.TypeKey] = keyValue
//				} else {
//					reqTemp[maxTypeKey] = &maxTypeKeyValue
//				}
//
//			}
//		}
//		var sourceID string
//		switch message.SessionType {
//		case constant.SingleChatType:
//			sourceID = message.SendID + message.RecvID
//		case constant.NotificationChatType:
//			sourceID = message.RecvID
//		case constant.GroupChatType, constant.SuperGroupChatType:
//			sourceID = message.RecvID
//		}
//		var apiReq server_api_params.SetMessageReactionExtensionsReq
//		apiReq.IsReact = message.IsReact
//		apiReq.ClientMsgID = message.ClientMsgID
//		apiReq.SourceID = sourceID
//		apiReq.SessionType = message.SessionType
//		apiReq.IsExternalExtensions = message.IsExternalExtensions
//		apiReq.ReactionExtensionList = reqTemp
//		apiReq.OperationID = operationID
//		apiReq.MsgFirstModifyTime = message.MsgFirstModifyTime
//		var apiResp server_api_params.SetMessageReactionExtensionsResp
//		c.p.PostFatalCallback(callback, constant.SetMessageReactionExtensionsRouter, apiReq, &apiResp.ApiResult, apiReq.OperationID)
//		var msg model_struct.LocalChatLogReactionExtensions
//		msg.ClientMsgID = message.ClientMsgID
//		resultKeyMap := make(map[string]*server_api_params.KeyValue)
//		for _, v := range apiResp.ApiResult.Result {
//			if v.ErrCode == 0 {
//				temp := new(server_api_params.KeyValue)
//				temp.TypeKey = v.TypeKey
//				temp.Value = v.Value
//				temp.LatestUpdateTime = v.LatestUpdateTime
//				resultKeyMap[v.TypeKey] = temp
//			}
//		}
//		err = c.db.GetAndUpdateMessageReactionExtension(message.ClientMsgID, resultKeyMap)
//		if err != nil {
//			log.Error(operationID, "GetAndUpdateMessageReactionExtension err:", err.Error())
//		}
//		if !message.IsReact {
//			message.IsReact = apiResp.ApiResult.IsReact
//			message.MsgFirstModifyTime = apiResp.ApiResult.MsgFirstModifyTime
//			err = c.db.UpdateMessageController(message)
//			if err != nil {
//				log.Error(operationID, "UpdateMessageController err:", err.Error(), message)
//
//			}
//		}
//		return apiResp.ApiResult.Result
//	}
//
//	func getIndexTypeKey(typeKey string, index int) string {
//		return typeKey + "$" + utils.IntToString(index)
//	}
func getPrefixTypeKey(typeKey string) string {
	list := strings.Split(typeKey, "$")
	if len(list) > 0 {
		return list[0]
	}
	return ""
}

//func (c *Conversation) getTypeKeyListInfo(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, keyList []string, operationID string) (result []*sdk.SingleTypeKeyInfoSum) {
//	message, err := c.db.GetMessageController(s)
//	common.CheckDBErrCallback(callback, err, operationID)
//	if message.Status != constant.MsgStatusSendSuccess {
//		common.CheckAnyErrCallback(callback, 201, errors.New("only send success message can modify reaction extensions"), operationID)
//	}
//	if !message.IsReact {
//		common.CheckAnyErrCallback(callback, 202, errors.New("can get message reaction ex"), operationID)
//	}
//	extendMsg, _ := c.db.GetMessageReactionExtension(message.ClientMsgID)
//	temp := make(map[string]*server_api_params.KeyValue)
//	_ = json.Unmarshal(extendMsg.LocalReactionExtensions, &temp)
//	for _, v := range keyList {
//		singleResult := new(sdk.SingleTypeKeyInfoSum)
//		singleResult.TypeKey = v
//		for typeKey, value := range temp {
//			if strings.HasPrefix(typeKey, v) {
//				singleTypeKeyInfo := new(sdk.SingleTypeKeyInfo)
//				_ = json.Unmarshal([]byte(value.Value), singleTypeKeyInfo)
//				if _, ok := singleTypeKeyInfo.InfoList[c.loginUserID]; ok {
//					singleResult.IsContainSelf = true
//				}
//				for _, info := range singleTypeKeyInfo.InfoList {
//					v := *info
//					singleResult.InfoList = append(singleResult.InfoList, &v)
//				}
//				singleResult.Counter += singleTypeKeyInfo.Counter
//			}
//		}
//		result = append(result, singleResult)
//	}
//	messageList := []*sdk_struct.MsgStruct{s}
//	_ = common.TriggerCmdSyncReactionExtensions(common.SyncReactionExtensionsNode{
//		OperationID: operationID,
//		Action:      constant.SyncMessageListTypeKeyInfo,
//		Args:        messageList,
//	}, c.GetCh())
//
//	return result
//}
//
//func (c *Conversation) getAllTypeKeyInfo(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, operationID string) (result []*sdk.SingleTypeKeyInfoSum) {
//	message, err := c.db.GetMessageController(s)
//	common.CheckDBErrCallback(callback, err, operationID)
//	if message.Status != constant.MsgStatusSendSuccess {
//		common.CheckAnyErrCallback(callback, 201, errors.New("only send success message can modify reaction extensions"), operationID)
//	}
//	if !message.IsReact {
//		common.CheckAnyErrCallback(callback, 202, errors.New("can get message reaction ex"), operationID)
//	}
//	extendMsg, _ := c.db.GetMessageReactionExtension(message.ClientMsgID)
//	temp := make(map[string]*server_api_params.KeyValue)
//	_ = json.Unmarshal(extendMsg.LocalReactionExtensions, &temp)
//	mapResult := make(map[string]*sdk.SingleTypeKeyInfoSum)
//	for typeKey, value := range temp {
//		singleTypeKeyInfo := new(sdk.SingleTypeKeyInfo)
//		err := json.Unmarshal([]byte(value.Value), singleTypeKeyInfo)
//		if err != nil {
//			log.Warn(operationID, "not this type ", value.Value)
//			continue
//		}
//		prefixKey := getPrefixTypeKey(typeKey)
//		if v, ok := mapResult[prefixKey]; ok {
//			for _, info := range singleTypeKeyInfo.InfoList {
//				t := *info
//				v.InfoList = append(v.InfoList, &t)
//			}
//			if _, ok := singleTypeKeyInfo.InfoList[c.loginUserID]; ok {
//				v.IsContainSelf = true
//			}
//			v.Counter += singleTypeKeyInfo.Counter
//		} else {
//			v := new(sdk.SingleTypeKeyInfoSum)
//			v.TypeKey = prefixKey
//			v.Counter = singleTypeKeyInfo.Counter
//			for _, info := range singleTypeKeyInfo.InfoList {
//				t := *info
//				v.InfoList = append(v.InfoList, &t)
//			}
//			if _, ok := singleTypeKeyInfo.InfoList[c.loginUserID]; ok {
//				v.IsContainSelf = true
//			}
//			mapResult[prefixKey] = v
//		}
//	}
//	for _, v := range mapResult {
//		result = append(result, v)
//
//	}
//	return result
//}
