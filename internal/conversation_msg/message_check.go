package conversation_msg

import "open_im_sdk/pkg/db/model_struct"

//
//import (
//	"open_im_sdk/pkg/constant"
//	"open_im_sdk/pkg/db/model_struct"
//	"open_im_sdk/pkg/log"
//	sdk "open_im_sdk/pkg/sdk_params_callback"
//	"open_im_sdk/pkg/utils"
//)
//
//func (c *Conversation) messageBlocksInternalContinuityCheck(sourceID string, seqList []uint32, notStartTime, isReverse bool, count, sessionType int, startTime int64, list *[]*model_struct.LocalChatLog, messageListCallback *sdk.GetAdvancedHistoryMessageListCallback, operationID string)  {
//	maxSeq, minSeq, haveSeqList := func(messages []*model_struct.LocalChatLog) (max, min uint32, seqList []uint32) {
//		for i := 0; i < len(messages); i++ {
//			if messages[i].Seq != 0 {
//				seqList = append(seqList, messages[i].Seq)
//			}
//			if messages[i].Seq != 0 && min == 0 &&max == 0{
//				min = messages[i].Seq
//				max = messages[i].Seq
//			}
//			if messages[i].Seq < min {
//				min = messages[i].Seq
//			}
//			if messages[i].Seq > max {
//				max = messages[i].Seq
//
//			}
//		}
//		return max, min, seqList
//	}(*list)
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
//			c.pullMessageAndReGetHistoryMessages(sourceID, pullSeqList, notStartTime, isReverse, count, sessionType, startTime, list, messageListCallback, operationID)
//		}
//	}
//}
//func (c *Conversation) messageBlocksBetweenContinuityCheck(sourceID string, seqList []uint32, notStartTime, isReverse bool, count, sessionType int, startTime int64, list *[]*model_struct.LocalChatLog, messageListCallback *sdk.GetAdvancedHistoryMessageListCallback, operationID string)  {
//	maxSeq, minSeq, haveSeqList := func(messages []*model_struct.LocalChatLog) (max, min uint32, seqList []uint32) {
//		for _, message := range messages {
//			if message.Seq != 0 {
//				max = message.Seq
//				min = message.Seq
//				break
//			}
//		}
//		for i := 0; i < len(messages); i++ {
//			if messages[i].Seq != 0 {
//				seqList = append(seqList, messages[i].Seq)
//			}
//			if messages[i].Seq > max {
//				max = messages[i].Seq
//
//			}
//			if messages[i].Seq < min {
//				min = messages[i].Seq
//			}
//		}
//		return max, min, seqList
//	}(list)
//}
//func (c *Conversation) messageBlocksEndContinuityCheck(sourceID string, seqList []uint32, notStartTime, isReverse bool, count, sessionType int, startTime int64, list *[]*model_struct.LocalChatLog, messageListCallback *sdk.GetAdvancedHistoryMessageListCallback, operationID string)  {
//	maxSeq, minSeq, haveSeqList := func(messages []*model_struct.LocalChatLog) (max, min uint32, seqList []uint32) {
//		for _, message := range messages {
//			if message.Seq != 0 {
//				max = message.Seq
//				min = message.Seq
//				break
//			}
//		}
//		for i := 0; i < len(messages); i++ {
//			if messages[i].Seq != 0 {
//				seqList = append(seqList, messages[i].Seq)
//			}
//			if messages[i].Seq > max {
//				max = messages[i].Seq
//
//			}
//			if messages[i].Seq < min {
//				min = messages[i].Seq
//			}
//		}
//		return max, min, seqList
//	}(list)
//}
func (c *Conversation) getMaxAndMinHaveSeqList(messages []*model_struct.LocalChatLog) (max, min uint32, seqList []uint32) {
	for i := 0; i < len(messages); i++ {
		if messages[i].Seq != 0 {
			seqList = append(seqList, messages[i].Seq)
		}
		if messages[i].Seq != 0 && min == 0 && max == 0 {
			min = messages[i].Seq
			max = messages[i].Seq
		}
		if messages[i].Seq < min {
			min = messages[i].Seq
		}
		if messages[i].Seq > max {
			max = messages[i].Seq

		}
	}
	return max, min, seqList
}
