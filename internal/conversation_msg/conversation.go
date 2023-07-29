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
	"errors"
	_ "open_im_sdk/internal/common"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/sdkerrs"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"sort"
	"strings"
	"time"

	"github.com/jinzhu/copier"

	"github.com/OpenIMSDK/tools/log"

	"github.com/OpenIMSDK/protocol/msg"
	"github.com/OpenIMSDK/protocol/sdkws"

	pbConversation "github.com/OpenIMSDK/protocol/conversation"
)

func (c *Conversation) setConversation(ctx context.Context, apiReq *pbConversation.SetConversationsReq, localConversation *model_struct.LocalConversation) error {
	apiReq.Conversation.ConversationID = localConversation.ConversationID
	apiReq.Conversation.ConversationType = localConversation.ConversationType
	apiReq.Conversation.UserID = localConversation.UserID
	apiReq.Conversation.GroupID = localConversation.GroupID
	apiReq.UserIDs = []string{c.loginUserID}
	if err := util.ApiPost(ctx, constant.SetConversationsRouter, apiReq, nil); err != nil {
		return err
	}
	return nil
}

func (c *Conversation) getServerConversationList(ctx context.Context) ([]*model_struct.LocalConversation, error) {
	resp, err := util.CallApi[pbConversation.GetAllConversationsResp](ctx, constant.GetAllConversationsRouter, pbConversation.GetAllConversationsReq{OwnerUserID: c.loginUserID})
	if err != nil {
		return nil, err
	}
	return util.Batch(ServerConversationToLocal, resp.Conversations), nil
}

func (c *Conversation) getServerConversationsByIDs(ctx context.Context, conversations []string) ([]*model_struct.LocalConversation, error) {
	resp, err := util.CallApi[pbConversation.GetConversationsResp](ctx, constant.GetConversationsRouter, pbConversation.GetConversationsReq{OwnerUserID: c.loginUserID, ConversationIDs: conversations})
	if err != nil {
		return nil, err
	}
	return util.Batch(ServerConversationToLocal, resp.Conversations), nil
}

func (c *Conversation) getServerHasReadAndMaxSeqs(ctx context.Context) (map[string]*msg.Seqs, error) {
	resp := &msg.GetConversationsHasReadAndMaxSeqResp{}
	err := util.ApiPost(ctx, constant.GetConversationsHasReadAndMaxSeqRouter, msg.GetConversationsHasReadAndMaxSeqReq{UserID: c.loginUserID}, resp)
	if err != nil {
		log.ZError(ctx, "getServerHasReadAndMaxSeqs err", err)
		return nil, err
	}
	return resp.Seqs, nil
}

func (c *Conversation) getHistoryMessageList(ctx context.Context, req sdk.GetHistoryMessageListParams, isReverse bool) ([]*sdk_struct.MsgStruct, error) {
	// t := time.Now()
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
			m, err := c.db.GetMessage(ctx, conversationID, msg.ClientMsgID)
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
			conversationID = c.getConversationIDBySessionType(sourceID, constant.SingleChatType)
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
			m, err := c.db.GetMessage(ctx, conversationID, msg.ClientMsgID)
			if err != nil {
				return nil, err
			}
			startTime = m.SendTime
		}
	}
	// log.Debug("", "Assembly parameters cost time", time.Since(t))
	// t = time.Now()
	// log.Info("", "sourceID:", sourceID, "startTime:", startTime, "count:", req.Count, "not start_time", notStartTime)
	if notStartTime {
		list, err = c.db.GetMessageListNoTime(ctx, conversationID, req.Count, isReverse)
	} else {
		list, err = c.db.GetMessageList(ctx, conversationID, req.Count, startTime, isReverse)
	}
	// log.Debug("", "db cost time", time.Since(t))
	if err != nil {
		return nil, err
	}
	// t = time.Now()
	for _, v := range list {
		temp := sdk_struct.MsgStruct{}
		// tt := time.Now()
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
		var attachedInfo sdk_struct.AttachedInfoElem
		_ = utils.JsonStringToStruct(v.AttachedInfo, &attachedInfo)
		temp.AttachedInfoElem = &attachedInfo
		temp.Ex = v.Ex
		temp.IsReact = v.IsReact
		temp.IsExternalExtensions = v.IsExternalExtensions
		err := c.msgHandleByContentType(&temp)
		if err != nil {
			// log.Error("", "Parsing data error:", err.Error(), temp)
			continue
		}
		// log.Debug("", "internal unmarshal cost time", time.Since(tt))

		switch sessionType {
		case constant.GroupChatType:
			fallthrough
		case constant.SuperGroupChatType:
			temp.GroupID = temp.RecvID
			temp.RecvID = c.loginUserID
		}
		messageList = append(messageList, &temp)
	}
	// log.Debug("", "unmarshal cost time", time.Since(t))
	// t = time.Now()
	if !isReverse {
		sort.Sort(messageList)
	}
	// log.Debug("", "sort cost time", time.Since(t))
	return messageList, nil
}

func (c *Conversation) getAdvancedHistoryMessageList2(ctx context.Context, req sdk.GetAdvancedHistoryMessageListParams, isReverse bool) (*sdk.GetAdvancedHistoryMessageListCallback, error) {
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
			notStartTime = true
		} else {
			msg.SessionType = lc.ConversationType
			msg.ClientMsgID = req.StartClientMsgID
			m, err := c.db.GetMessage(ctx, conversationID, msg.ClientMsgID)
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
			conversationID = c.getConversationIDBySessionType(sourceID, constant.SingleChatType)
			sessionType = constant.SingleChatType
		}
		if req.StartClientMsgID == "" {
			notStartTime = true
		} else {
			msg.ClientMsgID = req.StartClientMsgID
			m, err := c.db.GetMessage(ctx, conversationID, msg.ClientMsgID)
			if err != nil {
				return nil, err
			}
			startTime = m.SendTime
		}
	}
	log.ZDebug(ctx, "Assembly conversation parameters", "cost time", time.Since(t), "conversationID",
		conversationID, "startTime:", startTime, "count:", req.Count, "not start_time", notStartTime)
	t = time.Now()
	if notStartTime {
		list, err = c.db.GetMessageListNoTime(ctx, conversationID, req.Count, isReverse)
	} else {
		list, err = c.db.GetMessageList(ctx, conversationID, req.Count, startTime, isReverse)
	}
	log.ZDebug(ctx, "db get messageList", "cost time", time.Since(t), "len", len(list), "err", err, "conversationID", conversationID)

	if err != nil {
		return nil, err
	}
	rawMessageLength := len(list)
	t = time.Now()
	if rawMessageLength < req.Count {
		maxSeq, minSeq, lostSeqListLength := c.messageBlocksInternalContinuityCheck(ctx, conversationID, notStartTime, isReverse, req.Count, startTime, &list, &messageListCallback)
		_ = c.messageBlocksBetweenContinuityCheck(ctx, req.LastMinSeq, maxSeq, conversationID, notStartTime, isReverse, req.Count, startTime, &list, &messageListCallback)
		if minSeq == 1 && lostSeqListLength == 0 {
			messageListCallback.IsEnd = true
		} else {
			c.messageBlocksEndContinuityCheck(ctx, minSeq, conversationID, notStartTime, isReverse, req.Count, startTime, &list, &messageListCallback)
		}
	} else {
		maxSeq, _, _ := c.messageBlocksInternalContinuityCheck(ctx, conversationID, notStartTime, isReverse, req.Count, startTime, &list, &messageListCallback)
		c.messageBlocksBetweenContinuityCheck(ctx, req.LastMinSeq, maxSeq, conversationID, notStartTime, isReverse, req.Count, startTime, &list, &messageListCallback)

	}
	log.ZDebug(ctx, "pull message", "pull cost time", time.Since(t))
	t = time.Now()
	var thisMinSeq int64
	for _, v := range list {
		if v.Seq != 0 && thisMinSeq == 0 {
			thisMinSeq = v.Seq
		}
		if v.Seq < thisMinSeq && v.Seq != 0 {
			thisMinSeq = v.Seq
		}
		if v.Status >= constant.MsgStatusHasDeleted {
			log.ZDebug(ctx, "this message has been deleted or exception message", "msg", v)
			continue
		}
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
		var attachedInfo sdk_struct.AttachedInfoElem
		_ = utils.JsonStringToStruct(v.AttachedInfo, &attachedInfo)
		temp.AttachedInfoElem = &attachedInfo
		temp.Ex = v.Ex
		temp.LocalEx = v.LocalEx
		err := c.msgHandleByContentType(&temp)
		if err != nil {
			log.ZError(ctx, "Parsing data error", err, "temp", temp)
			continue
		}
		switch sessionType {
		case constant.GroupChatType:
			fallthrough
		case constant.SuperGroupChatType:
			temp.GroupID = temp.RecvID
			temp.RecvID = c.loginUserID
		}
		if attachedInfo.IsPrivateChat && temp.SendTime+int64(attachedInfo.BurnDuration) < time.Now().Unix() {
			continue
		}
		messageList = append(messageList, &temp)
	}
	log.ZDebug(ctx, "message convert and unmarshal", "unmarshal cost time", time.Since(t))
	t = time.Now()
	if !isReverse {
		sort.Sort(messageList)
	}
	log.ZDebug(ctx, "sort", "sort cost time", time.Since(t))
	messageListCallback.MessageList = messageList
	if thisMinSeq == 0 {
		thisMinSeq = req.LastMinSeq
	}
	messageListCallback.LastMinSeq = thisMinSeq
	return &messageListCallback, nil

}

func (c *Conversation) typingStatusUpdate(ctx context.Context, recvID, msgTip string) error {
	if recvID == "" {
		return sdkerrs.ErrArgs
	}
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Typing)
	if err != nil {
		return err
	}
	s.RecvID = recvID
	s.SessionType = constant.SingleChatType
	typingElem := sdk_struct.TypingElem{}
	typingElem.MsgTips = msgTip
	s.Content = utils.StructToJsonString(typingElem)
	options := make(map[string]bool, 6)
	utils.SetSwitchFromOptions(options, constant.IsHistory, false)
	utils.SetSwitchFromOptions(options, constant.IsPersistent, false)
	utils.SetSwitchFromOptions(options, constant.IsSenderSync, false)
	utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsSenderConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
	utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)

	var wsMsgData sdkws.MsgData
	copier.Copy(&wsMsgData, s)
	wsMsgData.Content = []byte(s.Content)
	wsMsgData.CreateTime = s.CreateTime
	wsMsgData.Options = options
	var sendMsgResp sdkws.UserSendMsgResp
	err = c.LongConnMgr.SendReqWaitResp(ctx, &wsMsgData, constant.SendMsg, &sendMsgResp)
	if err != nil {
		log.ZError(ctx, "send msg to server failed", err, "message", s)
		return err
	}
	return nil

}

//	funcation (c *Conversation) markMessageAsReadByConID(callback open_im_sdk_callback.Base, msgIDList sdk.MarkMessageAsReadByConIDParams, conversationID, operationID string) {
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
//		conversationID := c.getConversationIDBySessionType(userID, constant.SingleChatType)
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

func (c *Conversation) insertMessageToLocalStorage(ctx context.Context, conversationID string, s *model_struct.LocalChatLog) error {
	return c.db.InsertMessage(ctx, conversationID, s)
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
	var startTime, endTime int64
	var list []*model_struct.LocalChatLog
	conversationMap := make(map[string]*sdk.SearchByConversationResult, 10)
	var err error
	var conversationID string
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
		_, err := c.db.GetConversation(ctx, searchParam.ConversationID)
		if err != nil {
			return nil, err
		}
		if len(searchParam.MessageTypeList) != 0 && len(searchParam.KeywordList) == 0 {
			list, err = c.db.SearchMessageByContentType(ctx, searchParam.MessageTypeList, searchParam.ConversationID, startTime, endTime, offset, searchParam.Count)
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
			list, err = c.db.SearchMessageByKeyword(ctx, newContentTypeList, searchParam.KeywordList, searchParam.KeywordListMatchType,
				searchParam.ConversationID, startTime, endTime, offset, searchParam.Count)
		}
	} else {
		//Comprehensive search, search all
		if len(searchParam.MessageTypeList) == 0 {
			searchParam.MessageTypeList = SearchContentType
		}
		list, err = c.messageController.SearchMessageByContentTypeAndKeyword(ctx, searchParam.MessageTypeList, searchParam.KeywordList, searchParam.KeywordListMatchType, startTime, endTime)
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
	// log.Debug("", "get raw data length is", len(list))
	log.ZDebug(ctx, "get raw data length is", len(list))
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
		var attachedInfo sdk_struct.AttachedInfoElem
		_ = utils.JsonStringToStruct(v.AttachedInfo, &attachedInfo)
		temp.AttachedInfoElem = &attachedInfo
		temp.Ex = v.Ex
		temp.LocalEx = v.LocalEx
		err := c.msgHandleByContentType(&temp)
		if err != nil {
			// log.Error("", "Parsing data error:", err.Error(), temp)
			log.ZError(ctx, "Parsing data error:", err, "msg", temp)
			continue
		}
		if temp.ContentType == constant.File && !c.judgeMultipleSubString(searchParam.KeywordList, temp.FileElem.FileName, searchParam.KeywordListMatchType) {
			continue
		}
		if temp.ContentType == constant.AtText && !c.judgeMultipleSubString(searchParam.KeywordList, temp.AtTextElem.Text, searchParam.KeywordListMatchType) {
			continue
		}
		if temp.ContentType == constant.Text && !c.judgeMultipleSubString(searchParam.KeywordList, temp.TextElem.Content, searchParam.KeywordListMatchType) {
			continue
		}

		switch temp.SessionType {
		case constant.SingleChatType:
			if temp.SendID == c.loginUserID {
				conversationID = c.getConversationIDBySessionType(temp.RecvID, constant.SingleChatType)
			} else {
				conversationID = c.getConversationIDBySessionType(temp.SendID, constant.SingleChatType)
			}
		case constant.GroupChatType:
			temp.GroupID = temp.RecvID
			temp.RecvID = c.loginUserID
			conversationID = c.getConversationIDBySessionType(temp.GroupID, constant.GroupChatType)
		case constant.SuperGroupChatType:
			temp.GroupID = temp.RecvID
			temp.RecvID = c.loginUserID
			conversationID = c.getConversationIDBySessionType(temp.GroupID, constant.SuperGroupChatType)
		}
		if oldItem, ok := conversationMap[conversationID]; !ok {
			searchResultItem := sdk.SearchByConversationResult{}
			localConversation, err := c.db.GetConversation(ctx, conversationID)
			if err != nil {
				// log.Error("", "get conversation err ", err.Error(), conversationID)
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
	// var req server_api_params.DelMsgListReq
	// req.SeqList = seqList
	// req.OperationID = utils.OperationIDGenerator()
	// req.OpUserID = c.loginUserID
	// req.UserID = c.loginUserID
	// operationID := req.OperationID

	// err := c.SendReqWaitResp(context.Background(), &req, constant.WsDelMsg, 30, c.loginUserID)
	// if err != nil {
	// 	return utils.Wrap(err, "SendReqWaitResp failed")
	// }
	// var delResp server_api_params.DelMsgListResp
	// err = proto.Unmarshal(resp.Data, &delResp)
	// if err != nil {
	// 	log.Error(operationID, "Unmarshal failed ", err.Error())
	// 	return utils.Wrap(err, "Unmarshal failed")
	// }
	return nil
}

// old WS method
//funcation (c *Conversation) deleteMessageFromSvr(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, operationID string) {
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
	return nil, nil
	//message, err := c.db.GetMessageController(ctx, s)
	//if err != nil {
	//	return nil, err
	//}
	//if message.Status != constant.MsgStatusSendSuccess {
	//	return nil, errors.New("only send success message can modify reaction extensions")
	//}
	//if message.SessionType != constant.SuperGroupChatType {
	//	return nil, errors.New("currently only support super group message")
	//
	//}
	//extendMsg, _ := c.db.GetMessageReactionExtension(ctx, message.ClientMsgID)
	//temp := make(map[string]*server_api_params.KeyValue)
	//_ = json.Unmarshal(extendMsg.LocalReactionExtensions, &temp)
	//reqTemp := make(map[string]*server_api_params.KeyValue)
	//for _, v := range req {
	//	if value, ok := temp[v.TypeKey]; ok {
	//		v.LatestUpdateTime = value.LatestUpdateTime
	//	}
	//	reqTemp[v.TypeKey] = v
	//}
	//var sourceID string
	//switch message.SessionType {
	//case constant.SingleChatType:
	//	if message.SendID == c.loginUserID {
	//		sourceID = message.RecvID
	//	} else {
	//		sourceID = message.SendID
	//	}
	//case constant.NotificationChatType:
	//	sourceID = message.RecvID
	//case constant.GroupChatType, constant.SuperGroupChatType:
	//	sourceID = message.RecvID
	//}
	//var apiReq server_api_params.SetMessageReactionExtensionsReq
	//apiReq.IsReact = message.IsReact
	//apiReq.ClientMsgID = message.ClientMsgID
	//apiReq.SourceID = sourceID
	//apiReq.SessionType = message.SessionType
	//apiReq.IsExternalExtensions = message.IsExternalExtensions
	//apiReq.ReactionExtensionList = reqTemp
	//apiReq.OperationID = ""
	//apiReq.MsgFirstModifyTime = message.MsgFirstModifyTime
	//resp, err := util.CallApi[server_api_params.ApiResult](ctx, constant.SetMessageReactionExtensionsRouter, &apiReq)
	//if err != nil {
	//	return nil, err
	//}
	//var msg model_struct.LocalChatLogReactionExtensions
	//msg.ClientMsgID = message.ClientMsgID
	//resultKeyMap := make(map[string]*sdkws.KeyValue)
	//for _, v := range resp.Result {
	//	if v.ErrCode == 0 {
	//		temp := new(sdkws.KeyValue)
	//		temp.TypeKey = v.TypeKey
	//		temp.Value = v.Value
	//		temp.LatestUpdateTime = v.LatestUpdateTime
	//		resultKeyMap[v.TypeKey] = temp
	//	}
	//}
	//err = c.db.GetAndUpdateMessageReactionExtension(ctx, message.ClientMsgID, resultKeyMap)
	//if err != nil {
	//	log.Error("", "GetAndUpdateMessageReactionExtension err:", err.Error())
	//}
	//if !message.IsReact {
	//	message.IsReact = resp.IsReact
	//	message.MsgFirstModifyTime = resp.MsgFirstModifyTime
	//	err = c.db.UpdateMessageController(ctx, message)
	//	if err != nil {
	//		log.Error("", "UpdateMessageController err:", err.Error(), message)
	//
	//	}
	//}
	//return resp.Result, nil
}

func (c *Conversation) addMessageReactionExtensions(ctx context.Context, s *sdk_struct.MsgStruct, req sdk.AddMessageReactionExtensionsParams) ([]*server_api_params.ExtensionResult, error) {
	return nil, nil
	//message, err := c.db.GetMessageController(ctx, s)
	//if err != nil {
	//	return nil, err
	//}
	//if message.Status != constant.MsgStatusSendSuccess || message.Seq == 0 {
	//	return nil, errors.New("only send success message can modify reaction extensions")
	//}
	//reqTemp := make(map[string]*server_api_params.KeyValue)
	//extendMsg, err := c.db.GetMessageReactionExtension(ctx, message.ClientMsgID)
	//if err == nil && extendMsg != nil {
	//	temp := make(map[string]*server_api_params.KeyValue)
	//	_ = json.Unmarshal(extendMsg.LocalReactionExtensions, &temp)
	//	for _, v := range req {
	//		if value, ok := temp[v.TypeKey]; ok {
	//			v.LatestUpdateTime = value.LatestUpdateTime
	//		}
	//		reqTemp[v.TypeKey] = v
	//	}
	//} else {
	//	for _, v := range req {
	//		reqTemp[v.TypeKey] = v
	//	}
	//}
	//var sourceID string
	//switch message.SessionType {
	//case constant.SingleChatType:
	//	if message.SendID == c.loginUserID {
	//		sourceID = message.RecvID
	//	} else {
	//		sourceID = message.SendID
	//	}
	//case constant.NotificationChatType:
	//	sourceID = message.RecvID
	//case constant.GroupChatType, constant.SuperGroupChatType:
	//	sourceID = message.RecvID
	//}
	//var apiReq server_api_params.AddMessageReactionExtensionsReq
	//apiReq.IsReact = message.IsReact
	//apiReq.ClientMsgID = message.ClientMsgID
	//apiReq.SourceID = sourceID
	//apiReq.SessionType = message.SessionType
	//apiReq.IsExternalExtensions = message.IsExternalExtensions
	//apiReq.ReactionExtensionList = reqTemp
	//apiReq.OperationID = ""
	//apiReq.MsgFirstModifyTime = message.MsgFirstModifyTime
	//apiReq.Seq = message.Seq
	//
	//resp, err := util.CallApi[server_api_params.ApiResult](ctx, constant.AddMessageReactionExtensionsRouter, &apiReq)
	//if err != nil {
	//	return nil, err
	//}
	//log.Debug("", "api return:", message.IsReact, resp)
	//if !message.IsReact {
	//	message.IsReact = resp.IsReact
	//	message.MsgFirstModifyTime = resp.MsgFirstModifyTime
	//	err = c.db.UpdateMessageController(ctx, message)
	//	if err != nil {
	//		log.Error("", "UpdateMessageController err:", err.Error(), message)
	//	}
	//}
	//return resp.Result, nil
}

func (c *Conversation) deleteMessageReactionExtensions(ctx context.Context, s *sdk_struct.MsgStruct, req sdk.DeleteMessageReactionExtensionsParams) ([]*server_api_params.ExtensionResult, error) {
	// message, err := c.GetMessageController(ctx, s)
	// if err != nil {
	// 	return nil, err
	// }
	// if message.Status != constant.MsgStatusSendSuccess {
	// 	return nil, errors.New("only send success message can modify reaction extensions")
	// }
	// if message.SessionType != constant.SuperGroupChatType {
	// 	return nil, errors.New("currently only support super group message")

	// }
	// extendMsg, _ := c.db.GetMessageReactionExtension(ctx, message.ClientMsgID)
	// temp := make(map[string]*server_api_params.KeyValue)
	// _ = json.Unmarshal(extendMsg.LocalReactionExtensions, &temp)
	// var reqTemp []*server_api_params.KeyValue
	// for _, v := range req {
	// 	if value, ok := temp[v]; ok {
	// 		var tt server_api_params.KeyValue
	// 		tt.LatestUpdateTime = value.LatestUpdateTime
	// 		tt.TypeKey = v
	// 		reqTemp = append(reqTemp, &tt)
	// 	}
	// }
	// var sourceID string
	// switch message.SessionType {
	// case constant.SingleChatType:
	// 	if message.SendID == c.loginUserID {
	// 		sourceID = message.RecvID
	// 	} else {
	// 		sourceID = message.SendID
	// 	}
	// case constant.NotificationChatType:
	// 	sourceID = message.RecvID
	// case constant.GroupChatType, constant.SuperGroupChatType:
	// 	sourceID = message.RecvID
	// }
	// var apiReq server_api_params.DeleteMessageReactionExtensionsReq
	// apiReq.ClientMsgID = message.ClientMsgID
	// apiReq.SourceID = sourceID
	// apiReq.SessionType = message.SessionType
	// apiReq.ReactionExtensionList = reqTemp
	// apiReq.OperationID = ""
	// apiReq.IsExternalExtensions = message.IsExternalExtensions
	// apiReq.MsgFirstModifyTime = message.MsgFirstModifyTime
	// resp, err := util.CallApi[server_api_params.ApiResult](ctx, constant.AddMessageReactionExtensionsRouter, &apiReq)
	// if err != nil {
	// 	return nil, err
	// }
	// var msg model_struct.LocalChatLogReactionExtensions
	// msg.ClientMsgID = message.ClientMsgID
	// resultKeyMap := make(map[string]*sdkws.KeyValue)
	// for _, v := range resp.Result {
	// 	if v.ErrCode == 0 {
	// 		temp := new(sdkws.KeyValue)
	// 		temp.TypeKey = v.TypeKey
	// 		resultKeyMap[v.TypeKey] = temp
	// 	}
	// }
	// err = c.db.DeleteAndUpdateMessageReactionExtension(ctx, message.ClientMsgID, resultKeyMap)
	// if err != nil {
	// 	log.Error("", "GetAndUpdateMessageReactionExtension err:", err.Error())
	// }
	// return resp.Result, nil
	return nil, nil
}

type syncReactionExtensionParams struct {
	MessageList         []*model_struct.LocalChatLog
	SessionType         int32
	SourceID            string
	IsExternalExtension bool
	ExtendMessageList   []*model_struct.LocalChatLogReactionExtensions
	TypeKeyList         []string
}

func (c *Conversation) getMessageListReactionExtensions(ctx context.Context, conversationID string, messageList []*sdk_struct.MsgStruct) ([]*server_api_params.SingleMessageExtensionResult, error) {
	if len(messageList) == 0 {
		return nil, errors.New("message list is null")
	}
	var msgIDs []string
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
		msgIDs = append(msgIDs, msgStruct.ClientMsgID)
	}
	isExternalExtension = c.IsExternalExtensions
	localMessageList, err := c.db.GetMessagesByClientMsgIDs(ctx, conversationID, msgIDs)
	if err != nil {
		return nil, err
	}
	for _, v := range localMessageList {
		if v.IsReact != true {
			return nil, errors.New("have not reaction message in message list:" + v.ClientMsgID)
		}
	}
	var result server_api_params.GetMessageListReactionExtensionsResp
	extendMessage, _ := c.db.GetMultipleMessageReactionExtension(ctx, msgIDs)
	for _, v := range extendMessage {
		var singleResult server_api_params.SingleMessageExtensionResult
		// temp := make(map[string]*sdkws.KeyValue)
		// _ = json.Unmarshal(v.LocalReactionExtensions, &temp)
		singleResult.ClientMsgID = v.ClientMsgID
		// singleResult.ReactionExtensionList = temp
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

//	funcation (c *Conversation) getMessageListSomeReactionExtensions(callback open_im_sdk_callback.Base, messageList []*sdk_struct.MsgStruct, keyList []string, operationID string) server_api_params.GetMessageListReactionExtensionsResp {
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
//	funcation (c *Conversation) setTypeKeyInfo(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, typeKey, ex string, isCanRepeat bool, operationID string) []*server_api_params.ExtensionResult {
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
//	funcation getIndexTypeKey(typeKey string, index int) string {
//		return typeKey + "$" + utils.IntToString(index)
//	}
func getPrefixTypeKey(typeKey string) string {
	list := strings.Split(typeKey, "$")
	if len(list) > 0 {
		return list[0]
	}
	return ""
}

//funcation (c *Conversation) getTypeKeyListInfo(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, keyList []string, operationID string) (result []*sdk.SingleTypeKeyInfoSum) {
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
//funcation (c *Conversation) getAllTypeKeyInfo(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, operationID string) (result []*sdk.SingleTypeKeyInfoSum) {
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
