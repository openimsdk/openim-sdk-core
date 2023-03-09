package conversation_msg

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

// 检测其内部连续性，如果不连续，则向前补齐,获取这一组消息的最大最小seq，以及需要补齐的seq列表长度
func (c *Conversation) messageBlocksInternalContinuityCheck(sourceID string, notStartTime, isReverse bool, count,
	sessionType int, startTime int64, list *[]*model_struct.LocalChatLog, messageListCallback *sdk.GetAdvancedHistoryMessageListCallback,
	operationID string) (max, min uint32, length int) {
	var lostSeqListLength int
	maxSeq, minSeq, haveSeqList := c.getMaxAndMinHaveSeqList(*list)
	log.Debug(operationID, utils.GetSelfFuncName(), "getMaxAndMinHaveSeqList is:", maxSeq, minSeq, haveSeqList)
	if maxSeq != 0 && minSeq != 0 {
		successiveSeqList := func(max, min uint32) (seqList []uint32) {
			for i := min; i <= max; i++ {
				seqList = append(seqList, i)
			}
			return seqList
		}(maxSeq, minSeq)
		lostSeqList := utils.DifferenceSubset(successiveSeqList, haveSeqList)
		lostSeqListLength = len(lostSeqList)
		log.Debug(operationID, "get lost seqList is :", maxSeq, minSeq, lostSeqList, "length:", lostSeqListLength)
		if lostSeqListLength > 0 {
			var pullSeqList []uint32
			if lostSeqListLength <= constant.PullMsgNumForReadDiffusion {
				pullSeqList = lostSeqList
			} else {
				pullSeqList = lostSeqList[lostSeqListLength-constant.PullMsgNumForReadDiffusion : lostSeqListLength]
			}
			c.pullMessageAndReGetHistoryMessages(sourceID, pullSeqList, notStartTime, isReverse, count, sessionType, startTime, list, messageListCallback, operationID)
		}

	}
	return maxSeq, minSeq, lostSeqListLength
}

// 检测消息块之间的连续性，如果不连续，则向前补齐,返回块之间是否连续，bool
func (c *Conversation) messageBlocksBetweenContinuityCheck(lastMinSeq, maxSeq uint32, sourceID string, notStartTime, isReverse bool, count, sessionType int, startTime int64, list *[]*model_struct.LocalChatLog, messageListCallback *sdk.GetAdvancedHistoryMessageListCallback, operationID string) bool {
	if lastMinSeq != 0 {
		log.Debug(operationID, "get lost LastMinSeq is :", lastMinSeq, "thisMaxSeq is :", maxSeq)
		if maxSeq != 0 {
			if maxSeq+1 != lastMinSeq {
				startSeq := int64(lastMinSeq) - constant.PullMsgNumForReadDiffusion
				if startSeq <= int64(maxSeq) {
					startSeq = int64(maxSeq) + 1
				}
				successiveSeqList := func(max, min uint32) (seqList []uint32) {
					for i := min; i <= max; i++ {
						seqList = append(seqList, i)
					}
					return seqList
				}(lastMinSeq-1, uint32(startSeq))
				log.Debug(operationID, "get lost successiveSeqList is :", successiveSeqList, len(successiveSeqList))
				if len(successiveSeqList) > 0 {
					c.pullMessageAndReGetHistoryMessages(sourceID, successiveSeqList, notStartTime, isReverse, count, sessionType, startTime, list, messageListCallback, operationID)
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

// 检测其内部连续性，如果不连续，则向前补齐,获取这一组消息的最大最小seq，以及需要补齐的seq列表长度
func (c *Conversation) messageBlocksEndContinuityCheck(minSeq uint32, sourceID string, notStartTime, isReverse bool, count, sessionType int, startTime int64, list *[]*model_struct.LocalChatLog, messageListCallback *sdk.GetAdvancedHistoryMessageListCallback, operationID string) {
	var minSeqServer uint32
	var maxSeqServer uint32
	resp, err := c.SendReqWaitResp(&server_api_params.GetMaxAndMinSeqReq{UserID: c.loginUserID, GroupIDList: []string{sourceID}}, constant.WSGetNewestSeq, 1, 1, c.loginUserID, operationID)
	if err != nil {
		log.Error(operationID, "SendReqWaitResp failed ", err.Error(), constant.WSGetNewestSeq, 1, c.loginUserID)
	} else {
		var wsSeqResp server_api_params.GetMaxAndMinSeqResp
		err = proto.Unmarshal(resp.Data, &wsSeqResp)
		if err != nil {
			log.Error(operationID, "Unmarshal failed", err.Error())
		} else if wsSeqResp.ErrCode != 0 {
			log.Error(operationID, "GetMaxAndMinSeqReq failed ", wsSeqResp.ErrCode, wsSeqResp.ErrMsg)
		} else {
			if value, ok := wsSeqResp.GroupMaxAndMinSeq[sourceID]; ok {
				minSeqServer = value.MinSeq
				if value.MinSeq == 0 {
					minSeqServer = 1
				}
				maxSeqServer = value.MaxSeq
			}
		}
	}
	log.Error(operationID, "from server min seq is", minSeqServer, maxSeqServer)
	//seq, err := c.db.SuperGroupGetNormalMinSeq(sourceID)
	//if err != nil {
	//	log.Error(operationID, "SuperGroupGetNormalMinSeq err:", err.Error())
	//}
	//log.Error(operationID, sourceID+":table min seq is ", seq)
	if minSeq != 0 {
		if minSeq <= minSeqServer {
			messageListCallback.IsEnd = true
		} else {
			seqList := func(seq uint32) (seqList []uint32) {
				startSeq := int64(seq) - constant.PullMsgNumForReadDiffusion
				if startSeq <= int64(minSeqServer) {
					if minSeqServer == 0 {
						startSeq = 1
					} else {
						startSeq = int64(minSeqServer)
					}
				}
				log.Debug(operationID, "pull start is ", startSeq)
				for i := startSeq; i < int64(seq); i++ {
					seqList = append(seqList, uint32(i))
				}
				log.Debug(operationID, "pull seqList is ", seqList)
				return seqList
			}(minSeq)
			log.Debug(operationID, "pull seqList is ", seqList, len(seqList))
			if len(seqList) > 0 {
				c.pullMessageAndReGetHistoryMessages(sourceID, seqList, notStartTime, isReverse, count, sessionType, startTime, list, messageListCallback, operationID)
			}
		}
	} else {
		//local don't have messages,本地无消息，但是服务器最大消息不为0
		if int64(maxSeqServer)-int64(minSeqServer) >= 0 {
			messageListCallback.IsEnd = false
		} else {
			messageListCallback.IsEnd = true
		}

	}

}
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

// 1、保证单次拉取消息量低于sdk单次从服务器拉取量
// 2、块中连续性检测
// 3、块之间连续性检测
func (c *Conversation) pullMessageAndReGetHistoryMessages(sourceID string, seqList []uint32, notStartTime, isReverse bool, count, sessionType int, startTime int64, list *[]*model_struct.LocalChatLog, messageListCallback *sdk.GetAdvancedHistoryMessageListCallback, operationID string) {
	existedSeqList, err := c.db.SuperGroupGetAlreadyExistSeqList(sourceID, seqList)
	if err != nil {
		log.Error(operationID, "SuperGroupGetAlreadyExistSeqList err", err.Error(), sourceID, seqList)
		return
	}
	if len(existedSeqList) == len(seqList) {
		log.Debug(operationID, "do not pull message")
		return
	}
	newSeqList := utils.DifferenceSubset(seqList, existedSeqList)
	if len(newSeqList) == 0 {
		log.Debug(operationID, "do not pull message")
		return
	}
	var pullMsgReq server_api_params.PullMessageBySeqListReq
	pullMsgReq.UserID = c.loginUserID
	pullMsgReq.GroupSeqList = make(map[string]*server_api_params.SeqList, 0)
	pullMsgReq.GroupSeqList[sourceID] = &server_api_params.SeqList{SeqList: newSeqList}

	pullMsgReq.OperationID = operationID
	log.Debug(operationID, "read diffusion group pull message, req: ", pullMsgReq)
	resp, err := c.SendReqWaitResp(&pullMsgReq, constant.WSPullMsgBySeqList, 2, 1, c.loginUserID, operationID)
	if err != nil {
		errHandle(newSeqList, list, err, messageListCallback)
		log.Error(operationID, "SendReqWaitResp failed ", err.Error(), constant.WSPullMsgBySeqList, 1, 2, c.loginUserID)
	} else {
		var pullMsgResp server_api_params.PullMessageBySeqListResp
		err = proto.Unmarshal(resp.Data, &pullMsgResp)
		if err != nil {
			errHandle(newSeqList, list, err, messageListCallback)
			log.Error(operationID, "pullMsgResp Unmarshal failed ", err.Error())
		} else {
			log.Debug(operationID, "syncMsgFromServerSplit pull msg ", pullMsgReq.String(), pullMsgResp.String())
			if v, ok := pullMsgResp.GroupMsgDataList[sourceID]; ok {
				c.pullMessageIntoTable(v.MsgDataList, operationID)
			}
			if notStartTime {
				*list, err = c.db.GetMessageListNoTimeController(sourceID, sessionType, count, isReverse)
			} else {
				*list, err = c.db.GetMessageListController(sourceID, sessionType, count, startTime, isReverse)
			}
		}

	}
}
func errHandle(seqList []uint32, list *[]*model_struct.LocalChatLog, err error, messageListCallback *sdk.GetAdvancedHistoryMessageListCallback) {
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
func (c *Conversation) pullMessageIntoTable(pullMsgData []*server_api_params.MsgData, operationID string) {

	var insertMsg, specialUpdateMsg []*model_struct.LocalChatLog
	var exceptionMsg []*model_struct.LocalErrChatLog
	var msgReadList, groupMsgReadList, msgRevokeList, newMsgRevokeList sdk_struct.NewMsgList

	log.Info(operationID, "do Msg come here, len: ", len(pullMsgData))
	b := utils.GetCurrentTimestampByMill()

	for _, v := range pullMsgData {
		isConversationUpdate := utils.GetSwitchFromOptions(v.Options, constant.IsConversationUpdate)
		log.Info(operationID, "do Msg come here, msg detail ", v.RecvID, v.SendID, v.ClientMsgID, v.ServerMsgID, v.Seq, c.loginUserID)
		msg := new(sdk_struct.MsgStruct)
		copier.Copy(msg, v)
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
		if !isConversationUpdate {
			msg.Status = constant.MsgStatusFiltered
		}
		msg.IsRead = false
		//		log.Info(operationID, "new msg, seq, ServerMsgID, ClientMsgID", msg.Seq, msg.ServerMsgID, msg.ClientMsgID)
		//De-analyze data
		if msg.ClientMsgID == "" {
			exceptionMsg = append(exceptionMsg, c.msgStructToLocalErrChatLog(msg))
			continue
		}
		if v.SendID == c.loginUserID { //seq
			// Messages sent by myself  //if  sent through  this terminal
			m, err := c.db.GetMessageController(msg)
			if err == nil {
				log.Info(operationID, "have message", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, *msg)
				if m.Seq == 0 {
					specialUpdateMsg = append(specialUpdateMsg, c.msgStructToLocalChatLog(msg))

				} else {
					exceptionMsg = append(exceptionMsg, c.msgStructToLocalErrChatLog(msg))
				}
			} else { //      send through  other terminal
				log.Info(operationID, "sync message", msg.Seq, msg.ServerMsgID, msg.ClientMsgID, *msg)
				insertMsg = append(insertMsg, c.msgStructToLocalChatLog(msg))
				switch msg.ContentType {
				case constant.Revoke:
					msgRevokeList = append(msgRevokeList, msg)
				case constant.HasReadReceipt:
					msgReadList = append(msgReadList, msg)
				case constant.GroupHasReadReceipt:
					groupMsgReadList = append(groupMsgReadList, msg)
				case constant.AdvancedRevoke:
					newMsgRevokeList = append(newMsgRevokeList, msg)
				default:
				}
			}
		} else { //Sent by others
			if oldMessage, err := c.db.GetMessageController(msg); err != nil { //Deduplication operation
				log.Warn("test", "trigger msg is ", msg.SenderNickname, msg.SenderFaceURL)
				insertMsg = append(insertMsg, c.msgStructToLocalChatLog(msg))
				switch msg.ContentType {
				case constant.Revoke:
					msgRevokeList = append(msgRevokeList, msg)
				case constant.HasReadReceipt:
					msgReadList = append(msgReadList, msg)
				case constant.GroupHasReadReceipt:
					groupMsgReadList = append(groupMsgReadList, msg)
				case constant.AdvancedRevoke:
					newMsgRevokeList = append(newMsgRevokeList, msg)

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

	//update message
	err6 := c.db.BatchSpecialUpdateMessageList(specialUpdateMsg)
	if err6 != nil {
		log.Error(operationID, "sync seq normal message err  :", err6.Error())
	}
	b3 := utils.GetCurrentTimestampByMill()
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
		log.Error(operationID, "BatchInsertExceptionMsgController err message err  :", err2.Error())

	}
	b8 := utils.GetCurrentTimestampByMill()
	c.DoGroupMsgReadState(groupMsgReadList)
	b9 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "DoGroupMsgReadState  cost time : ", b9-b8, "len: ", len(groupMsgReadList))

	c.revokeMessage(msgRevokeList)
	c.newRevokeMessage(newMsgRevokeList)
	b10 := utils.GetCurrentTimestampByMill()
	log.Debug(operationID, "revokeMessage  cost time : ", b10-b9)
	log.Info(operationID, "insert msg, total cost time: ", utils.GetCurrentTimestampByMill()-b, "len:  ", len(pullMsgData))
}

//拉取的消息都需要经过块内部连续性检测以及块和上一块之间的连续性检测不连续则补，补齐的过程中如果出现任何异常只给seq从大到小到断层
//拉取消息不满量，获取服务器中该群最大seq以及用户对于此群最小seq，本地该群的最小seq，如果本地的不为0并且小于等于服务器最小的，说明已经到底部
//如果本地的为0，可以理解为初始化的时候，数据还未同步，或者异常情况，如果服务器最大seq-服务器最小seq>=0说明还未到底部，否则到底部
