package indexdb

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"open_im_sdk/wasm/indexdb/temp_struct"
)

type LocalSuperGroupChatLogs struct{}

func (i *LocalSuperGroupChatLogs) GetSuperGroupNormalMsgSeq(groupID string) (uint32, error) {
	seq, err := Exec(groupID)
	if err != nil {
		return 0, err
	} else {
		if v, ok := seq.(float64); ok {
			var result uint32
			result = uint32(v)
			return result, err
		} else {
			return 0, ErrType
		}
	}
}
func (i *LocalSuperGroupChatLogs) SuperGroupGetNormalMinSeq(groupID string) (uint32, error) {
	seq, err := Exec(groupID)
	if err != nil {
		return 0, err
	} else {
		if v, ok := seq.(float64); ok {
			var result uint32
			result = uint32(v)
			return result, err
		} else {
			return 0, ErrType
		}
	}
}

func (i *LocalSuperGroupChatLogs) SuperGroupGetMessage(message *sdk_struct.MsgStruct) (*model_struct.LocalChatLog, error) {
	msg, err := Exec(message.GroupID, message.ClientMsgID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msg.(string); ok {
			result := model_struct.LocalChatLog{}
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return &result, err
		} else {
			return nil, ErrType
		}
	}
}

func (i *LocalSuperGroupChatLogs) SuperGroupUpdateMessage(c *model_struct.LocalChatLog) error {
	if c.ClientMsgID == "" {
		return PrimaryKeyNull
	}
	tempLocalChatLog := temp_struct.LocalChatLog{
		ServerMsgID:          c.ServerMsgID,
		SendID:               c.SendID,
		RecvID:               c.RecvID,
		SenderPlatformID:     c.SenderPlatformID,
		SenderNickname:       c.SenderNickname,
		SenderFaceURL:        c.SenderFaceURL,
		SessionType:          c.SessionType,
		MsgFrom:              c.MsgFrom,
		ContentType:          c.ContentType,
		Content:              c.Content,
		IsRead:               c.IsRead,
		Status:               c.Status,
		Seq:                  c.Seq,
		SendTime:             c.SendTime,
		CreateTime:           c.CreateTime,
		AttachedInfo:         c.AttachedInfo,
		Ex:                   c.Ex,
		IsReact:              c.IsReact,
		IsExternalExtensions: c.IsExternalExtensions,
		MsgFirstModifyTime:   c.MsgFirstModifyTime,
	}
	_, err := Exec(c.RecvID, c.ClientMsgID, utils.StructToJsonString(tempLocalChatLog))
	return err

}
func (i *LocalSuperGroupChatLogs) SuperGroupBatchInsertMessageList(messageList []*model_struct.LocalChatLog, groupID string) error {
	_, err := Exec(utils.StructToJsonString(messageList), groupID)
	return err
}
func (i *LocalSuperGroupChatLogs) SuperGroupInsertMessage(message *model_struct.LocalChatLog, groupID string) error {
	_, err := Exec(utils.StructToJsonString(message), groupID)
	return err
}
func (i *LocalSuperGroupChatLogs) SuperGroupGetMultipleMessage(msgIDList []string, groupID string) (result []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec(msgIDList, groupID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msgList.(string); ok {
			var temp []model_struct.LocalChatLog
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, ErrType
		}
	}
}
func (i *LocalSuperGroupChatLogs) SuperGroupUpdateMessageTimeAndStatus(msg *sdk_struct.MsgStruct) error {
	_, err := Exec(msg.GroupID, msg.ClientMsgID, msg.ServerMsgID, msg.SendTime, msg.Status)
	return err
}
func (i *LocalSuperGroupChatLogs) SuperGroupGetMessageListNoTime(sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec(sourceID, sessionType, count, isReverse)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msgList.(string); ok {
			var temp []model_struct.LocalChatLog
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, ErrType
		}
	}
}
func (i *LocalSuperGroupChatLogs) SuperGroupGetMessageList(sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec(sourceID, sessionType, count, startTime, isReverse)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msgList.(string); ok {
			var temp []model_struct.LocalChatLog
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, ErrType
		}
	}
}

func (i *LocalChatLogs) SuperGroupSearchMessageByKeyword(contentType []int, keywordList []string, keywordListMatchType int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec(contentType, keywordList, keywordListMatchType, sourceID, startTime, endTime, sessionType, offset, count)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msgList.(string); ok {
			var temp []model_struct.LocalChatLog
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, ErrType
		}
	}
}

func (i *LocalChatLogs) SearchMessageByKeyword(contentType []int, keywordList []string, keywordListMatchType int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (messages []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec(contentType, keywordList, keywordListMatchType, sourceID, startTime, endTime, sessionType, offset, count)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msgList.(string); ok {
			var temp []model_struct.LocalChatLog
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				messages = append(messages, &v1)
			}
			return messages, err
		} else {
			return nil, ErrType
		}
	}
}

func (i *LocalChatLogs) SuperGroupSearchAllMessageByContentType(superGroupID string, contentType int32) (result []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec(superGroupID, contentType)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msgList.(string); ok {
			var temp []*model_struct.LocalChatLog
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, v1)
			}
			return result, err
		} else {
			return nil, ErrType
		}
	}
}

func (i *LocalChatLogs) SuperGroupUpdateColumnsMessage(clientMsgID, groupID string, args map[string]interface{}) error {
	_, err := Exec(clientMsgID, groupID, utils.StructToJsonString(args))
	return err
}
func (i *LocalChatLogs) GetSuperGroupAbnormalMsgSeq(groupID string) (uint32, error) {
	return 0, nil
}

func (i *LocalChatLogs) InitSuperLocalErrChatLog(groupID string) {
	panic("implement me")
}
func (i *LocalChatLogs) SuperBatchInsertExceptionMsg(MessageList []*model_struct.LocalErrChatLog, groupID string) error {
	return nil
}

func (i IndexDB) InitSuperLocalChatLog(groupID string) {
	panic("implement me")
}

func (i IndexDB) SuperGroupDeleteAllMessage(groupID string) error {
	panic("implement me")
}

func (i IndexDB) SuperGroupSearchMessageByContentTypeAndKeyword(contentType []int, keywordList []string, keywordListMatchType int, startTime, endTime int64, groupID string) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) SuperGroupBatchUpdateMessageList(MessageList []*model_struct.LocalChatLog) error {
	panic("implement me")
}

func (i IndexDB) SuperGroupMessageIfExists(ClientMsgID string) (bool, error) {
	panic("implement me")
}

func (i IndexDB) SuperGroupIsExistsInErrChatLogBySeq(seq int64) bool {
	panic("implement me")
}

func (i IndexDB) SuperGroupMessageIfExistsBySeq(seq int64) (bool, error) {
	panic("implement me")
}

func (i IndexDB) SuperGroupGetAllUnDeleteMessageSeqList() ([]uint32, error) {
	panic("implement me")
}

func (i IndexDB) SuperGroupUpdateMessageStatusBySourceID(sourceID string, status, sessionType int32) error {
	panic("implement me")
}

func (i IndexDB) SuperGroupGetSendingMessageList() (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) SuperGroupUpdateGroupMessageHasRead(msgIDList []string, groupID string) error {
	panic("implement me")
}

func (i IndexDB) SuperGroupGetNormalMsgSeq() (uint32, error) {
	return 0, nil
}

func (i IndexDB) SuperGroupGetTestMessage(seq uint32) (*model_struct.LocalChatLog, error) {
	panic("implement me")
}

func (i IndexDB) SuperGroupUpdateMsgSenderNickname(sendID, nickname string, sType int) error {
	panic("implement me")
}

func (i IndexDB) SuperGroupUpdateMsgSenderFaceURL(sendID, faceURL string, sType int) error {
	panic("implement me")
}

func (i IndexDB) SuperGroupUpdateMsgSenderFaceURLAndSenderNickname(sendID, faceURL, nickname string, sessionType int) error {
	panic("implement me")
}

func (i IndexDB) SuperGroupGetMsgSeqByClientMsgID(clientMsgID string, groupID string) (uint32, error) {
	panic("implement me")
}

func (i IndexDB) SuperGroupGetMsgSeqListByGroupID(groupID string) ([]uint32, error) {
	panic("implement me")
}

func (i IndexDB) SuperGroupGetMsgSeqListByPeerUserID(userID string) ([]uint32, error) {
	panic("implement me")
}

func (i IndexDB) SuperGroupGetMsgSeqListBySelfUserID(userID string) ([]uint32, error) {
	panic("implement me")
}
