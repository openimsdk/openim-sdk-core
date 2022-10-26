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
		ServerMsgID:      c.ServerMsgID,
		SendID:           c.SendID,
		RecvID:           c.RecvID,
		SenderPlatformID: c.SenderPlatformID,
		SenderNickname:   c.SenderNickname,
		SenderFaceURL:    c.SenderFaceURL,
		SessionType:      c.SessionType,
		MsgFrom:          c.MsgFrom,
		ContentType:      c.ContentType,
		Content:          c.Content,
		IsRead:           c.IsRead,
		Status:           c.Status,
		Seq:              c.Seq,
		SendTime:         c.SendTime,
		CreateTime:       c.CreateTime,
		AttachedInfo:     c.AttachedInfo,
		Ex:               c.Ex,
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
