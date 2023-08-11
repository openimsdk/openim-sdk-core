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
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/utils"

	"github.com/OpenIMSDK/tools/log"

	"github.com/OpenIMSDK/protocol/sdkws"
	utils2 "github.com/OpenIMSDK/tools/utils"
)

// 检测其内部连续性，如果不连续，则向前补齐,获取这一组消息的最大最小seq，以及需要补齐的seq列表长度
func (c *Conversation) messageBlocksInternalContinuityCheck(ctx context.Context, conversationID string, notStartTime, isReverse bool, count int,
	startTime int64, list *[]*model_struct.LocalChatLog, messageListCallback *sdk.GetAdvancedHistoryMessageListCallback) (max, min int64, length int) {
	var lostSeqListLength int
	maxSeq, minSeq, haveSeqList := c.getMaxAndMinHaveSeqList(*list)
	log.ZDebug(ctx, "getMaxAndMinHaveSeqList is:", "maxSeq", maxSeq, "minSeq", minSeq, "haveSeqList", haveSeqList)
	if maxSeq != 0 && minSeq != 0 {
		successiveSeqList := func(max, min int64) (seqList []int64) {
			for i := min; i <= max; i++ {
				seqList = append(seqList, i)
			}
			return seqList
		}(maxSeq, minSeq)
		lostSeqList := utils.DifferenceSubset(successiveSeqList, haveSeqList)
		lostSeqListLength = len(lostSeqList)
		log.ZDebug(ctx, "get lost seqList is :", "maxSeq", maxSeq, "minSeq", minSeq, "lostSeqList", lostSeqList, "length:", lostSeqListLength)
		if lostSeqListLength > 0 {
			var pullSeqList []int64
			if lostSeqListLength <= constant.PullMsgNumForReadDiffusion {
				pullSeqList = lostSeqList
			} else {
				pullSeqList = lostSeqList[lostSeqListLength-constant.PullMsgNumForReadDiffusion : lostSeqListLength]
			}
			c.pullMessageAndReGetHistoryMessages(ctx, conversationID, pullSeqList, notStartTime, isReverse, count, startTime, list, messageListCallback)
		}

	}
	return maxSeq, minSeq, lostSeqListLength
}

// 检测消息块之间的连续性，如果不连续，则向前补齐,返回块之间是否连续，bool
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

// 根据最小seq向前补齐消息，由服务器告诉拉取消息结果是否到底，如果网络，则向前补齐,获取这一组消息的最大最小seq，以及需要补齐的seq列表长度
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
			c.pullMessageAndReGetHistoryMessages(ctx, conversationID, seqList, notStartTime, isReverse, count, startTime, list, messageListCallback)
		}

	} else {
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
func (c *Conversation) pullMessageAndReGetHistoryMessages(ctx context.Context, conversationID string, seqList []int64, notStartTime,
	isReverse bool, count int, startTime int64, list *[]*model_struct.LocalChatLog, messageListCallback *sdk.GetAdvancedHistoryMessageListCallback) {
	existedSeqList, err := c.db.GetAlreadyExistSeqList(ctx, conversationID, seqList)
	if err != nil {
		log.ZError(ctx, "GetAlreadyExistSeqList err", err, "conversationID", conversationID, "seqList", seqList)
		return
	}
	if len(existedSeqList) == len(seqList) {
		log.ZDebug(ctx, "do not pull message", "seqList", seqList, "existedSeqList", existedSeqList)
		return
	}
	newSeqList := utils.DifferenceSubset(seqList, existedSeqList)
	if len(newSeqList) == 0 {
		log.ZDebug(ctx, "do not pull message", "seqList", seqList, "existedSeqList", existedSeqList, "newSeqList", newSeqList)
		return
	}
	var pullMsgResp sdkws.PullMessageBySeqsResp
	var pullMsgReq sdkws.PullMessageBySeqsReq
	pullMsgReq.UserID = c.loginUserID
	var seqRange sdkws.SeqRange
	seqRange.ConversationID = conversationID
	seqRange.Begin = newSeqList[0]
	seqRange.End = newSeqList[len(newSeqList)-1]
	seqRange.Num = int64(len(newSeqList))
	pullMsgReq.SeqRanges = append(pullMsgReq.SeqRanges, &seqRange)
	log.ZDebug(ctx, "conversation pull message,  ", "req", pullMsgReq)
	if notStartTime && !c.LongConnMgr.IsConnected() {
		return
	}
	err = c.SendReqWaitResp(ctx, &pullMsgReq, constant.PullMsgBySeqList, &pullMsgResp)
	if err != nil {
		errHandle(newSeqList, list, err, messageListCallback)
		log.ZDebug(ctx, "pullmsg SendReqWaitResp failed", err, "req")
	} else {
		log.ZDebug(ctx, "syncMsgFromServerSplit pull msg", "resp", pullMsgResp)
		if v, ok := pullMsgResp.Msgs[conversationID]; ok {
			c.pullMessageIntoTable(ctx, pullMsgResp.Msgs, conversationID)
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
func (c *Conversation) pullMessageIntoTable(ctx context.Context, pullMsgData map[string]*sdkws.PullMsgs, conversationID string) {
	insertMsg := make(map[string][]*model_struct.LocalChatLog, 20)
	updateMsg := make(map[string][]*model_struct.LocalChatLog, 30)
	var insertMessage, selfInsertMessage, othersInsertMessage []*model_struct.LocalChatLog
	var updateMessage []*model_struct.LocalChatLog
	var exceptionMsg []*model_struct.LocalErrChatLog

	log.ZDebug(ctx, "do Msg come here, len: ", "msg length", len(pullMsgData))
	for conversationID, msgs := range pullMsgData {
		for _, v := range msgs.Msgs {
			log.ZDebug(ctx, "msg detail", "msg", v, "conversationID", conversationID)
			msg := c.msgDataToLocalChatLog(v)
			//When the message has been marked and deleted by the cloud, it is directly inserted locally without any conversation and message update.
			if msg.Status == constant.MsgStatusHasDeleted {
				insertMessage = append(insertMessage, msg)
				continue
			}
			msg.Status = constant.MsgStatusSendSuccess
			//		log.Info(operationID, "new msg, seq, ServerMsgID, ClientMsgID", msg.Seq, msg.ServerMsgID, msg.ClientMsgID)
			//De-analyze data
			if msg.ClientMsgID == "" {
				exceptionMsg = append(exceptionMsg, c.msgDataToLocalErrChatLog(msg))
				continue
			}
			if v.SendID == c.loginUserID { //seq
				// Messages sent by myself  //if  sent through  this terminal
				m, err := c.db.GetMessage(ctx, conversationID, msg.ClientMsgID)
				if err == nil {
					log.ZInfo(ctx, "have message", "msg", msg)
					if m.Seq == 0 {
						updateMessage = append(updateMessage, msg)

					} else {
						exceptionMsg = append(exceptionMsg, c.msgDataToLocalErrChatLog(msg))
					}
				} else { //      send through  other terminal
					log.ZInfo(ctx, "sync message", "msg", msg)
					selfInsertMessage = append(selfInsertMessage, msg)
				}
			} else { //Sent by others
				if oldMessage, err := c.db.GetMessage(ctx, conversationID, msg.ClientMsgID); err != nil { //Deduplication operation
					othersInsertMessage = append(othersInsertMessage, msg)

				} else {
					if oldMessage.Seq == 0 {
						updateMessage = append(updateMessage, msg)
					}
				}
			}

			insertMsg[conversationID] = append(insertMessage, c.faceURLAndNicknameHandle(ctx, selfInsertMessage, othersInsertMessage, conversationID)...)
			updateMsg[conversationID] = updateMessage
		}

		//update message
		if err6 := c.messageController.BatchUpdateMessageList(ctx, updateMsg); err6 != nil {
			log.ZError(ctx, "sync seq normal message err  :", err6)
		}
		b3 := utils.GetCurrentTimestampByMill()
		//Normal message storage
		_ = c.messageController.BatchInsertMessageList(ctx, insertMsg)
		b4 := utils.GetCurrentTimestampByMill()
		log.ZDebug(ctx, "BatchInsertMessageListController, ", "cost time", b4-b3)

		//Exception message storage
		for _, v := range exceptionMsg {
			log.ZWarn(ctx, "exceptionMsg show: ", nil, "msg", *v)
		}

	}
}

// 拉取的消息都需要经过块内部连续性检测以及块和上一块之间的连续性检测不连续则补，补齐的过程中如果出现任何异常只给seq从大到小到断层
// 拉取消息不满量，获取服务器中该群最大seq以及用户对于此群最小seq，本地该群的最小seq，如果本地的不为0并且小于等于服务器最小的，说明已经到底部
// 如果本地的为0，可以理解为初始化的时候，数据还未同步，或者异常情况，如果服务器最大seq-服务器最小seq>=0说明还未到底部，否则到底部

func (c *Conversation) faceURLAndNicknameHandle(ctx context.Context, self, others []*model_struct.LocalChatLog, conversationID string) []*model_struct.LocalChatLog {
	lc, _ := c.db.GetConversation(ctx, conversationID)
	switch lc.ConversationType {
	case constant.SingleChatType:
		c.singleHandle(ctx, self, others, lc)
	case constant.SuperGroupChatType:
		c.groupHandle(ctx, self, others, lc)
	}
	return append(self, others...)
}
func (c *Conversation) singleHandle(ctx context.Context, self, others []*model_struct.LocalChatLog, lc *model_struct.LocalConversation) {
	userInfo, err := c.db.GetLoginUser(ctx, c.loginUserID)
	if err == nil {
		for _, chatLog := range self {
			chatLog.SenderFaceURL = userInfo.FaceURL
			chatLog.SenderNickname = userInfo.Nickname
		}
	}
	for _, chatLog := range others {
		chatLog.SenderFaceURL = lc.FaceURL
		chatLog.SenderNickname = lc.ShowName
	}

}
func (c *Conversation) groupHandle(ctx context.Context, self, others []*model_struct.LocalChatLog, lc *model_struct.LocalConversation) {
	allMessage := append(self, others...)
	localGroupMemberInfo, err := c.group.GetSpecifiedGroupMembersInfo(ctx, lc.GroupID, utils2.Slice(allMessage, func(e *model_struct.LocalChatLog) string {
		return e.SendID
	}))
	if err != nil {
		log.ZError(ctx, "get group member info err", err)
		return
	}
	groupMap := utils2.SliceToMap(localGroupMemberInfo, func(e *model_struct.LocalGroupMember) string {
		return e.UserID
	})
	for _, chatLog := range allMessage {
		if g, ok := groupMap[chatLog.SendID]; ok {
			chatLog.SenderFaceURL = g.FaceURL
			chatLog.SenderNickname = g.Nickname
		}
	}
}
