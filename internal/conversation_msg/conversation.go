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
	"sort"
	"sync"
	"time"

	"github.com/jinzhu/copier"
	"golang.org/x/sync/errgroup"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	sdk "github.com/openimsdk/openim-sdk-core/v3/pkg/sdk_params_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"

	"github.com/openimsdk/tools/log"

	"github.com/openimsdk/protocol/sdkws"

	pbConversation "github.com/openimsdk/protocol/conversation"
)

func (c *Conversation) setConversation(ctx context.Context, apiReq *pbConversation.SetConversationsReq, localConversation *model_struct.LocalConversation) error {
	apiReq.Conversation.ConversationID = localConversation.ConversationID
	apiReq.Conversation.ConversationType = localConversation.ConversationType
	apiReq.Conversation.UserID = localConversation.UserID
	apiReq.Conversation.GroupID = localConversation.GroupID
	apiReq.UserIDs = []string{c.loginUserID}

	return api.SetConversations.Execute(ctx, apiReq)
}

func (c *Conversation) getAdvancedHistoryMessageList(ctx context.Context, req sdk.GetAdvancedHistoryMessageListParams, isReverse bool) (*sdk.GetAdvancedHistoryMessageListCallback, error) {
	t := time.Now()
	var messageListCallback sdk.GetAdvancedHistoryMessageListCallback
	var conversationID string
	var startTime int64
	var sessionType int
	var list []*model_struct.LocalChatLog
	var err error
	var messageList sdk_struct.NewMsgList
	var notStartTime bool
	conversationID = req.ConversationID
	lc, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	sessionType = int(lc.ConversationType)
	if req.StartClientMsgID == "" {
		notStartTime = true
	} else {
		m, err := c.db.GetMessage(ctx, conversationID, req.StartClientMsgID)
		if err != nil {
			return nil, err
		}
		startTime = m.SendTime
	}
	log.ZDebug(ctx, "Assembly conversation parameters", "cost time", time.Since(t), "conversationID",
		conversationID, "startTime:", startTime, "count:", req.Count, "not start_time", notStartTime)
	t = time.Now()
	if notStartTime {
		list, err = c.db.GetMessageListNoTime(ctx, conversationID, req.Count, isReverse)
	} else {
		list, err = c.db.GetMessageList(ctx, conversationID, req.Count, startTime, isReverse)
	}
	log.ZDebug(ctx, "db get messageList", "cost time", time.Since(t), "len", len(list), "err",
		err, "conversationID", conversationID)

	if err != nil {
		return nil, err
	}
	rawMessageLength := len(list)
	t = time.Now()
	if rawMessageLength < req.Count {
		maxSeq, minSeq, lostSeqListLength := c.messageBlocksInternalContinuityCheck(ctx,
			conversationID, notStartTime, isReverse, req.Count, startTime, &list, &messageListCallback)
		_ = c.messageBlocksBetweenContinuityCheck(ctx, req.LastMinSeq, maxSeq, conversationID,
			notStartTime, isReverse, req.Count, startTime, &list, &messageListCallback)
		if minSeq == 1 && lostSeqListLength == 0 {
			messageListCallback.IsEnd = true
		} else {
			c.messageBlocksEndContinuityCheck(ctx, minSeq, conversationID, notStartTime, isReverse,
				req.Count, startTime, &list, &messageListCallback)
		}
	} else {
		maxSeq, _, _ := c.messageBlocksInternalContinuityCheck(ctx, conversationID, notStartTime, isReverse,
			req.Count, startTime, &list, &messageListCallback)
		c.messageBlocksBetweenContinuityCheck(ctx, req.LastMinSeq, maxSeq, conversationID, notStartTime,
			isReverse, req.Count, startTime, &list, &messageListCallback)

	}
	log.ZDebug(ctx, "pull message", "pull cost time", time.Since(t))
	t = time.Now()
	//var thisMinSeq int64
	//for _, v := range list {
	//	if v.Seq != 0 && thisMinSeq == 0 {
	//		thisMinSeq = v.Seq
	//	}
	//	if v.Seq < thisMinSeq && v.Seq != 0 {
	//		thisMinSeq = v.Seq
	//	}
	//	if v.Status >= constant.MsgStatusHasDeleted {
	//		log.ZDebug(ctx, "this message has been deleted or exception message", "msg", v)
	//		continue
	//	}
	//	temp := sdk_struct.MsgStruct{}
	//	temp.ClientMsgID = v.ClientMsgID
	//	temp.ServerMsgID = v.ServerMsgID
	//	temp.CreateTime = v.CreateTime
	//	temp.SendTime = v.SendTime
	//	temp.SessionType = v.SessionType
	//	temp.SendID = v.SendID
	//	temp.RecvID = v.RecvID
	//	temp.MsgFrom = v.MsgFrom
	//	temp.ContentType = v.ContentType
	//	temp.SenderPlatformID = v.SenderPlatformID
	//	temp.SenderNickname = v.SenderNickname
	//	temp.SenderFaceURL = v.SenderFaceURL
	//	temp.Content = v.Content
	//	temp.Seq = v.Seq
	//	temp.IsRead = v.IsRead
	//	temp.Status = v.Status
	//	var attachedInfo sdk_struct.AttachedInfoElem
	//	_ = utils.JsonStringToStruct(v.AttachedInfo, &attachedInfo)
	//	temp.AttachedInfoElem = &attachedInfo
	//	temp.Ex = v.Ex
	//	temp.LocalEx = v.LocalEx
	//	err := c.msgHandleByContentType(&temp)
	//	if err != nil {
	//		log.ZError(ctx, "Parsing data error", err, "temp", temp)
	//		continue
	//	}
	//	switch sessionType {
	//	case constant.WriteGroupChatType:
	//		fallthrough
	//	case constant.ReadGroupChatType:
	//		temp.GroupID = temp.RecvID
	//		temp.RecvID = c.loginUserID
	//	}
	//	if attachedInfo.IsPrivateChat && temp.SendTime+int64(attachedInfo.BurnDuration) < time.Now().Unix() {
	//		continue
	//	}
	//	messageList = append(messageList, &temp)
	//}
	var thisMinSeq int64
	thisMinSeq, messageList = c.LocalChatLog2MsgStruct(ctx, list, sessionType)
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

func (c *Conversation) LocalChatLog2MsgStruct(ctx context.Context, list []*model_struct.LocalChatLog, sessionType int) (int64, []*sdk_struct.MsgStruct) {
	messageList := make([]*sdk_struct.MsgStruct, 0, len(list))
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
		case constant.WriteGroupChatType:
			fallthrough
		case constant.ReadGroupChatType:
			temp.GroupID = temp.RecvID
			temp.RecvID = c.loginUserID
		}
		if attachedInfo.IsPrivateChat && temp.SendTime+int64(attachedInfo.BurnDuration) < time.Now().Unix() {
			continue
		}
		messageList = append(messageList, &temp)
	}
	return thisMinSeq, messageList
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

// searchLocalMessages searches for local messages based on the given search parameters.
func (c *Conversation) searchLocalMessages(ctx context.Context, searchParam *sdk.SearchLocalMessagesParams) (*sdk.SearchLocalMessagesCallback, error) {
	var r sdk.SearchLocalMessagesCallback                                   // Initialize the result structure
	var startTime, endTime int64                                            // Variables to hold start and end times for the search
	var list []*model_struct.LocalChatLog                                   // Slice to store the search results
	conversationMap := make(map[string]*sdk.SearchByConversationResult, 10) // Map to store results grouped by conversation, with initial capacity of 10
	var err error                                                           // Variable to store any errors encountered
	var conversationID string                                               // Variable to store the current conversation ID

	// Set the end time for the search; if SearchTimePosition is 0, use the current timestamp
	if searchParam.SearchTimePosition == 0 {
		endTime = time.Now().Unix()
	} else {
		endTime = searchParam.SearchTimePosition
	}

	// Set the start time based on the specified time period
	if searchParam.SearchTimePeriod != 0 {
		startTime = endTime - searchParam.SearchTimePeriod
	}
	// Convert start and end times to milliseconds
	startTime = utils.UnixSecondToTime(startTime).UnixNano() / 1e6
	endTime = utils.UnixSecondToTime(endTime).UnixNano() / 1e6

	// Validate that either keyword list or message type list is provided
	if len(searchParam.KeywordList) == 0 && len(searchParam.MessageTypeList) == 0 {
		return nil, errors.New("keywordlist and messageTypelist all null")
	}

	// Search in a specific conversation if ConversationID is provided
	if searchParam.ConversationID != "" {
		// Validate pagination parameters
		if searchParam.PageIndex < 1 || searchParam.Count < 1 {
			return nil, errors.New("page or count is null")
		}
		offset := (searchParam.PageIndex - 1) * searchParam.Count
		_, err := c.db.GetConversation(ctx, searchParam.ConversationID)
		if err != nil {
			return nil, err
		}
		// Search by content type or keyword based on provided parameters
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
		// Comprehensive search across all conversations
		if len(searchParam.MessageTypeList) == 0 {
			searchParam.MessageTypeList = SearchContentType
		}
		list, err = c.searchMessageByContentTypeAndKeyword(ctx, searchParam.MessageTypeList, searchParam.KeywordList, searchParam.KeywordListMatchType, startTime, endTime)
	}

	// Handle any errors encountered during the search
	if err != nil {
		return nil, err
	}

	// Logging and processing each message in the search results
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
		if c.filterMsg(&temp, searchParam) {
			continue
		}
		// Determine the conversation ID based on the session type
		switch temp.SessionType {
		case constant.SingleChatType:
			if temp.SendID == c.loginUserID {
				conversationID = c.getConversationIDBySessionType(temp.RecvID, constant.SingleChatType)
			} else {
				conversationID = c.getConversationIDBySessionType(temp.SendID, constant.SingleChatType)
			}
		case constant.WriteGroupChatType:
			temp.GroupID = temp.RecvID
			temp.RecvID = c.loginUserID
			conversationID = c.getConversationIDBySessionType(temp.GroupID, constant.WriteGroupChatType)
		case constant.ReadGroupChatType:
			temp.GroupID = temp.RecvID
			temp.RecvID = c.loginUserID
			conversationID = c.getConversationIDBySessionType(temp.GroupID, constant.ReadGroupChatType)
		}
		// Populate the conversationMap with search results
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
			searchResultItem.LatestMsgSendTime = localConversation.LatestMsgSendTime
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

	// Compile the results from the conversationMap into the response structure
	for _, v := range conversationMap {
		r.SearchResultItems = append(r.SearchResultItems, v)
		r.TotalCount += v.MessageCount
	}

	// Sort the search results based on the latest message send time
	sort.Slice(r.SearchResultItems, func(i, j int) bool {
		return r.SearchResultItems[i].LatestMsgSendTime > r.SearchResultItems[j].LatestMsgSendTime
	})

	return &r, nil // Return the final search results
}

func (c *Conversation) searchMessageByContentTypeAndKeyword(ctx context.Context, contentType []int, keywordList []string,
	keywordListMatchType int, startTime, endTime int64) (result []*model_struct.LocalChatLog, err error) {
	var list []*model_struct.LocalChatLog
	conversationIDList, err := c.db.GetAllConversationIDList(ctx)
	if err != nil {
		return nil, err
	}

	var mu sync.Mutex
	g, _ := errgroup.WithContext(ctx)
	g.SetLimit(searchMessageGoroutineLimit)
	for _, v := range conversationIDList {
		conversationID := v
		g.Go(func() error {
			sList, err := c.db.SearchMessageByContentTypeAndKeyword(ctx, contentType, conversationID, keywordList, keywordListMatchType, startTime, endTime)
			if err != nil {
				log.ZWarn(ctx, "search conversation message", err, "conversationID", conversationID)
				return nil
			}

			mu.Lock()
			list = append(list, sList...)
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return list, nil
}

// true is filter, false is not filter
func (c *Conversation) filterMsg(temp *sdk_struct.MsgStruct, searchParam *sdk.SearchLocalMessagesParams) bool {
	switch temp.ContentType {
	case constant.Text:
		return !c.judgeMultipleSubString(searchParam.KeywordList, temp.TextElem.Content,
			searchParam.KeywordListMatchType)
	case constant.AtText:
		return !c.judgeMultipleSubString(searchParam.KeywordList, temp.AtTextElem.Text,
			searchParam.KeywordListMatchType)
	case constant.File:
		return !c.judgeMultipleSubString(searchParam.KeywordList, temp.FileElem.FileName,
			searchParam.KeywordListMatchType)
	case constant.Merger:
		if !c.judgeMultipleSubString(searchParam.KeywordList, temp.MergeElem.Title, searchParam.KeywordListMatchType) {
			for _, msgStruct := range temp.MergeElem.MultiMessage {
				if c.filterMsg(msgStruct, searchParam) {
					continue
				} else {
					break
				}
			}
		}
	case constant.Card:
		return !c.judgeMultipleSubString(searchParam.KeywordList, temp.CardElem.Nickname,
			searchParam.KeywordListMatchType)
	case constant.Location:
		return !c.judgeMultipleSubString(searchParam.KeywordList, temp.LocationElem.Description,
			searchParam.KeywordListMatchType)
	case constant.Custom:
		return !c.judgeMultipleSubString(searchParam.KeywordList, temp.CustomElem.Description,
			searchParam.KeywordListMatchType)
	case constant.Quote:
		if !c.judgeMultipleSubString(searchParam.KeywordList, temp.QuoteElem.Text, searchParam.KeywordListMatchType) {
			return c.filterMsg(temp.QuoteElem.QuoteMessage, searchParam)
		}
	case constant.Picture:
		fallthrough
	case constant.Video:
		if len(searchParam.KeywordList) == 0 {
			return false
		} else {
			return true
		}
	default:
		return true
	}
	return false
}
