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

//拉取的消息都需要经过块内部连续性检测以及块和上一块之间的连续性检测不连续则补，补齐的过程中如果出现任何异常只给seq从大到小到断层
//拉取消息不满量，获取服务器中该群最大seq以及用户对于此群最小seq，本地该群的最小seq，如果本地的不为0并且小于等于服务器最小的，说明已经到底部
//如果本地的为0，可以理解为初始化的时候，数据还未同步，或者异常情况，如果服务器最大seq>0说明还未到底部，否则到底部
