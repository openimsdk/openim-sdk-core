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
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	sdk "github.com/openimsdk/openim-sdk-core/v3/pkg/sdk_params_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/tools/log"

	"github.com/openimsdk/protocol/sdkws"
)

// Check for internal continuity. If discontinuity is found, fill in the gaps.
// Retrieve the maximum and minimum seq of this group of messages, as well as the length of the seq list that needs to be filled in.
func (c *Conversation) messageBlocksInternalContinuityCheck(ctx context.Context, conversationID string, notStartTime, isReverse bool, count int,
	startTime int64, list *[]*model_struct.LocalChatLog, messageListCallback *sdk.GetAdvancedHistoryMessageListCallback) (max, min int64, length int) {
	var lostSeqListLength int
	maxSeq, minSeq, haveSeqList := c.getMaxAndMinHaveSeqList(*list)
	log.ZDebug(ctx, "getMaxAndMinHaveSeqList is:", "maxSeq", maxSeq, "minSeq", minSeq, "haveSeqList", haveSeqList)
	if maxSeq != 0 && minSeq != 0 {
		var lostSeqList []int64
		haveSeqSet := datautil.SliceSetAny(haveSeqList, func(e int64) int64 {
			return e
		})
		for i := minSeq; i <= maxSeq; i++ {
			if _, found := haveSeqSet[i]; !found {
				lostSeqList = append(lostSeqList, i)
			}
		}
		lostSeqListLength = len(lostSeqList)
		log.ZDebug(ctx, "get lost seqList is :", "maxSeq", maxSeq, "minSeq", minSeq, "lostSeqList", lostSeqList, "length:", lostSeqListLength)
		if lostSeqListLength > 0 {
			var pullSeqList []int64
			if lostSeqListLength <= constant.PullMsgNumForReadDiffusion {
				pullSeqList = lostSeqList
			} else {
				pullSeqList = lostSeqList[lostSeqListLength-constant.PullMsgNumForReadDiffusion : lostSeqListLength]
			}
			log.ZDebug(ctx, "messageBlocksInternalContinuityCheck", "pullSeqList", pullSeqList)
			c.pullMessageAndReGetHistoryMessages(ctx, conversationID, pullSeqList, notStartTime, isReverse, count, startTime, list, messageListCallback)
		}

	}
	return maxSeq, minSeq, lostSeqListLength
}

// Check the continuity between message blocks. If discontinuity is found, fill in the gaps forward.
// Returns whether the blocks are continuous as a boolean value.
func (c *Conversation) messageBlocksBetweenContinuityCheck(ctx context.Context, lastMinSeq, maxSeq int64, conversationID string,
	notStartTime, isReverse bool, count int, startTime int64, list *[]*model_struct.LocalChatLog, messageListCallback *sdk.GetAdvancedHistoryMessageListCallback) bool {
	if lastMinSeq != 0 {
		log.ZDebug(ctx, "get lost LastMinSeq is :", "lastMinSeq", lastMinSeq, "thisMaxSeq", maxSeq)
		if maxSeq != 0 {
			if maxSeq+1 != lastMinSeq {
				startSeq := int64(lastMinSeq) - constant.PullMsgNumForReadDiffusion
				if startSeq <= maxSeq {
					startSeq = int64(maxSeq) + 1
				}
				successiveSeqList := func(max, min int64) (seqList []int64) {
					for i := min; i <= max; i++ {
						seqList = append(seqList, i)
					}
					return seqList
				}(lastMinSeq-1, startSeq)
				log.ZDebug(ctx, "get lost successiveSeqList is :", "successiveSeqList", successiveSeqList, "length:", len(successiveSeqList))
				if len(successiveSeqList) > 0 {
					log.ZDebug(ctx, "messageBlocksBetweenContinuityCheck", "successiveSeqList", successiveSeqList)
					c.pullMessageAndReGetHistoryMessages(ctx, conversationID, successiveSeqList, notStartTime, isReverse, count, startTime, list, messageListCallback)
				}
			} else {
				return true
			}

		} else {
			return true
		}

	} else {
		return true
	}

	return false
}

func (c *Conversation) messageBlocksEndContinuityCheck(ctx context.Context, minSeq int64, conversationID string, notStartTime,
	isReverse bool, count int, startTime int64, list *[]*model_struct.LocalChatLog, messageListCallback *sdk.GetAdvancedHistoryMessageListCallback) {
	if minSeq != 0 {
		seqList := func(seq int64) (seqList []int64) {
			startSeq := seq - constant.PullMsgNumForReadDiffusion
			if startSeq <= 0 {
				startSeq = 1
			}
			log.ZDebug(ctx, "pull start is ", "start seq", startSeq)
			for i := startSeq; i < seq; i++ {
				seqList = append(seqList, i)
			}
			return seqList
		}(minSeq)
		log.ZDebug(ctx, "pull seqList is ", "seqList", seqList, "len", len(seqList))

		if len(seqList) > 0 {
			log.ZDebug(ctx, "messageBlocksEndContinuityCheck", "seqList", seqList)
			c.pullMessageAndReGetHistoryMessages(ctx, conversationID, seqList, notStartTime, isReverse, count, startTime, list, messageListCallback)
		}

	} else {
		log.ZDebug(ctx, "messageBlocksEndContinuityCheck", "minSeq", minSeq, "conversationID", conversationID)
		//local don't have messages,本地无消息，但是服务器最大消息不为0
		seqList := []int64{0, 0}
		c.pullMessageAndReGetHistoryMessages(ctx, conversationID, seqList, notStartTime, isReverse, count, startTime, list, messageListCallback)

	}

}
func (c *Conversation) getMaxAndMinHaveSeqList(messages []*model_struct.LocalChatLog) (max, min int64, seqList []int64) {
	for i := 0; i < len(messages); i++ {
		if messages[i].Seq != 0 {
			seqList = append(seqList, messages[i].Seq)
		}
		if messages[i].Seq != 0 && min == 0 && max == 0 {
			min = messages[i].Seq
			max = messages[i].Seq
		}
		if messages[i].Seq < min && messages[i].Seq != 0 {
			min = messages[i].Seq
		}
		if messages[i].Seq > max {
			max = messages[i].Seq

		}
	}
	return max, min, seqList
}

// 1、保证单次拉取消息量低于sdk单次从服务器拉取量
// 2、块中连续性检测
// 3、块之间连续性检测
func (c *Conversation) pullMessageAndReGetHistoryMessages(ctx context.Context, conversationID string, seqList []int64,
	notStartTime, isReverse bool, count int, startTime int64, list *[]*model_struct.LocalChatLog,
	messageListCallback *sdk.GetAdvancedHistoryMessageListCallback) {
	existedSeqList, err := c.db.GetAlreadyExistSeqList(ctx, conversationID, seqList)
	if err != nil {
		log.ZError(ctx, "GetAlreadyExistSeqList err", err, "conversationID", conversationID,
			"seqList", seqList)
		return
	}
	if len(existedSeqList) > 0 {
		log.ZWarn(ctx, "GetAlreadyExistSeqList", nil, "conversationID", conversationID, "seqList", seqList, "existedSeqList", existedSeqList)
	}
	if len(existedSeqList) == len(seqList) {
		log.ZDebug(ctx, "do not pull message", "seqList", seqList, "existedSeqList", existedSeqList)
		return
	}
	newSeqList := utils.DifferenceSubset(seqList, existedSeqList)
	if len(newSeqList) == 0 {
		log.ZDebug(ctx, "do not pull message", "seqList", seqList, "existedSeqList", existedSeqList,
			"newSeqList", newSeqList)
		return
	}
	var getSeqMessageResp msg.GetSeqMessageResp
	var getSeqMessageReq msg.GetSeqMessageReq
	getSeqMessageReq.UserID = c.loginUserID
	var conversationSeqs msg.ConversationSeqs
	conversationSeqs.ConversationID = conversationID
	conversationSeqs.Seqs = newSeqList
	log.ZDebug(ctx, "conversation pull message,  ", "req", getSeqMessageReq)
	if notStartTime && !c.LongConnMgr.IsConnected() {
		return
	}
	err = c.SendReqWaitResp(ctx, &getSeqMessageReq, constant.PullMsgBySeqList, &getSeqMessageResp)
	if err != nil {
		errHandle(newSeqList, list, err, messageListCallback)
		log.ZDebug(ctx, "pull SendReqWaitResp failed", err, "req")
	} else {
		log.ZDebug(ctx, "syncMsgFromServerSplit pull msg", "resp", getSeqMessageResp)
		if getSeqMessageResp.Msgs == nil {
			log.ZWarn(ctx, "syncMsgFromServerSplit pull msg is null", errors.New("pull message is null"),
				"req", getSeqMessageResp)
			return
		}
		if v, ok := getSeqMessageResp.Msgs[conversationID]; ok {
			c.pullMessageIntoTable(ctx, getSeqMessageResp.Msgs)
			messageListCallback.IsEnd = v.IsEnd

			if notStartTime {
				*list, err = c.db.GetMessageListNoTime(ctx, conversationID, count, isReverse)
			} else {
				*list, err = c.db.GetMessageList(ctx, conversationID, count, startTime, isReverse)
			}
		}

	}
}
func errHandle(seqList []int64, list *[]*model_struct.LocalChatLog, err error, messageListCallback *sdk.GetAdvancedHistoryMessageListCallback) {
	messageListCallback.ErrCode = 100
	messageListCallback.ErrMsg = err.Error()
	var result []*model_struct.LocalChatLog
	needPullMaxSeq := seqList[len(seqList)-1]
	for _, chatLog := range *list {
		if chatLog.Seq == 0 || chatLog.Seq > needPullMaxSeq {
			temp := chatLog
			result = append(result, temp)
		} else {
			if chatLog.Seq <= needPullMaxSeq {
				break
			}
		}
	}
	*list = result
}
func (c *Conversation) pullMessageIntoTable(ctx context.Context, pullMsgData map[string]*sdkws.PullMsgs) {
	insertMsg := make(map[string][]*model_struct.LocalChatLog, 20)
	updateMsg := make(map[string][]*model_struct.LocalChatLog, 30)
	var insertMessage, selfInsertMessage, othersInsertMessage []*model_struct.LocalChatLog
	var updateMessage []*model_struct.LocalChatLog
	var exceptionMsg []*model_struct.LocalErrChatLog

	log.ZDebug(ctx, "do Msg come here, len: ", "msg length", len(pullMsgData))
	for conversationID, msgs := range pullMsgData {
		msgIDs := datautil.Slice(msgs.Msgs, func(msg *sdkws.MsgData) string {
			return msg.ClientMsgID
		})
		localMessages, err := c.db.GetMessagesByClientMsgIDs(ctx, conversationID, msgIDs)
		if err != nil {
			log.ZWarn(ctx, "Failed to get messages by ClientMsgIDs", err)
		}
		localMessagesMap := datautil.SliceToMap(localMessages, func(msg *model_struct.LocalChatLog) string { return msg.ClientMsgID })
		for _, v := range msgs.Msgs {
			log.ZDebug(ctx, "msg detail", "msg", v, "conversationID", conversationID)
			msg := c.msgDataToLocalChatLog(v)
			//When the message has been marked and deleted by the cloud, it is directly inserted locally
			//without any conversation and message update.
			if msg.Status == constant.MsgStatusHasDeleted {
				insertMessage = append(insertMessage, msg)
				continue
			}
			msg.Status = constant.MsgStatusSendSuccess
			// The message might be a filler provided by the server due to a gap in the sequence.
			if msg.ClientMsgID == "" {
				msg.ClientMsgID = utils.GetMsgID(c.loginUserID) + utils.Int64ToString(msg.Seq)
				exceptionMsg = append(exceptionMsg, c.msgDataToLocalErrChatLog(msg))
				insertMessage = append(insertMessage, msg)
				continue
			}
			existingMsg, exists := localMessagesMap[msg.ClientMsgID]
			if v.SendID == c.loginUserID { //seq
				// Messages sent by myself  //if  sent through  this terminal
				if exists {
					log.ZDebug(ctx, "have message", "msg", msg)
					if existingMsg.Seq == 0 {
						updateMessage = append(updateMessage, msg)

					} else {
						// The message you sent is duplicated, possibly due to a resend or the server consuming
						// the message multiple times.
						msg.ClientMsgID = msg.ClientMsgID + utils.Int64ToString(msg.Seq)
						exceptionMsg = append(exceptionMsg, c.msgDataToLocalErrChatLog(msg))
						insertMessage = append(insertMessage, msg)
					}
				} else { //      send through  other terminal
					log.ZDebug(ctx, "sync message", "msg", msg)
					selfInsertMessage = append(selfInsertMessage, msg)
				}
			} else { //Sent by others
				if !exists {
					othersInsertMessage = append(othersInsertMessage, msg)

				} else {
					// The message sent by others is duplicated, possibly due to a resend or the server consuming
					// the message multiple times.
					msg.ClientMsgID = msg.ClientMsgID + utils.Int64ToString(msg.Seq)
					exceptionMsg = append(exceptionMsg, c.msgDataToLocalErrChatLog(msg))
					insertMessage = append(insertMessage, msg)
				}
			}

		}
		timeNow := time.Now()
		insertMsg[conversationID] = append(insertMessage, c.faceURLAndNicknameHandle(ctx, selfInsertMessage, othersInsertMessage, conversationID)...)
		updateMsg[conversationID] = updateMessage
		log.ZDebug(ctx, "faceURLAndNicknameHandle, ", "cost time", time.Since(timeNow).Milliseconds(),
			"updateMsg", updateMessage, "insertMsg", insertMessage, "selfInsertMessage", selfInsertMessage, "othersInsertMessage", othersInsertMessage)

		//update message
		if err6 := c.batchUpdateMessageList(ctx, updateMsg); err6 != nil {
			log.ZError(ctx, "sync seq normal message err  :", err6)
		}
		timeNow = time.Now()
		//Normal message storage
		_ = c.batchInsertMessageList(ctx, insertMsg)
		log.ZDebug(ctx, "BatchInsertMessageListController, ", "cost time", time.Since(timeNow).Milliseconds())

		//Exception message storage
		for _, v := range exceptionMsg {
			log.ZWarn(ctx, "exceptionMsg show: ", nil, "msg", *v)
		}

	}
}

// 拉取的消息都需要经过块内部连续性检测以及块和上一块之间的连续性检测不连续则补，补齐的过程中如果出现任何异常只给seq从大到小到断层
// 拉取消息不满量，获取服务器中该群最大seq以及用户对于此群最小seq，本地该群的最小seq，如果本地的不为0并且小于等于服务器最小的，说明已经到底部
// 如果本地的为0，可以理解为初始化的时候，数据还未同步，或者异常情况，如果服务器最大seq-服务器最小seq>=0说明还未到底部，否则到底部

// faceURLAndNicknameHandle handles the assignment of face URLs and nicknames for chat logs
// based on the conversation type (single chat or group chat).
// It first retrieves the conversation information using the provided conversationID.
// Depending on the conversation type, it delegates the handling to either singleHandle (for single chats)
// or groupHandle (for group chats). If conversation information retrieval fails, it returns the merged chat logs.
func (c *Conversation) faceURLAndNicknameHandle(ctx context.Context, self, others []*model_struct.LocalChatLog, conversationID string) []*model_struct.LocalChatLog {
	lc, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return append(self, others...)
	}
	switch lc.ConversationType {
	case constant.SingleChatType:
		c.singleHandle(ctx, self, others, lc)
	case constant.SuperGroupChatType:
		c.groupHandle(ctx, self, others, lc)
	}
	return append(self, others...)
}

// singleHandle processes chat logs for single chat conversations.
// It updates the SenderFaceURL and SenderNickname fields for messages in the `self` list
// using the logged-in user's information, and for messages in the `others` list
// using the other party's information if available in the conversation.
func (c *Conversation) singleHandle(ctx context.Context, self, others []*model_struct.LocalChatLog, lc *model_struct.LocalConversation) {
	if len(self) > 0 {
		userInfo, err := c.db.GetLoginUser(ctx, c.loginUserID)
		if err == nil {
			for _, chatLog := range self {
				chatLog.SenderFaceURL = userInfo.FaceURL
				chatLog.SenderNickname = userInfo.Nickname
			}
		}
	}

	if lc.FaceURL != "" && lc.ShowName != "" {
		for _, chatLog := range others {
			chatLog.SenderFaceURL = lc.FaceURL
			chatLog.SenderNickname = lc.ShowName
		}
	}
}

// groupHandle processes chat logs for group chat conversations.
// It merges the `self` and `others` chat logs and updates the SenderFaceURL and SenderNickname fields
// using the group members' information. If group member information is not available,
// it attempts to retrieve the sender's information from a local cache.
func (c *Conversation) groupHandle(ctx context.Context, self, others []*model_struct.LocalChatLog, lc *model_struct.LocalConversation) {
	allMessage := append(self, others...)

	allSenders := datautil.Slice(allMessage, func(e *model_struct.LocalChatLog) string {
		return e.SendID
	})
	localGroupMemberInfo, err := c.group.GetSpecifiedGroupMembersInfo(ctx, lc.GroupID, datautil.Distinct(allSenders))
	if err != nil {
		log.ZError(ctx, "get group member info err", err)
		return
	}
	groupMap := datautil.SliceToMap(localGroupMemberInfo, func(e *model_struct.LocalGroupMember) string {
		return e.UserID
	})
	for _, chatLog := range allMessage {
		if g, ok := groupMap[chatLog.SendID]; ok { // If group member info is successfully retrieved
			if g.FaceURL != "" && g.Nickname != "" {
				chatLog.SenderFaceURL = g.FaceURL
				chatLog.SenderNickname = g.Nickname
			}
		} else { // Otherwise, retrieve from local temporary cache
			faceURL, name, err := c.getUserNameAndFaceURL(ctx, chatLog.SendID)
			if err != nil {
				log.ZWarn(ctx, "getUserNameAndFaceURL error", err, "senderID", chatLog.SendID)
			} else if faceURL != "" && name != "" {
				chatLog.SenderFaceURL = faceURL
				chatLog.SenderNickname = name
			}
		}
	}
}
