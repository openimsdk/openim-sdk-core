package indexdb

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/indexdb/temp_struct"
)

//1,使用wasm原生的方式，tinygo应用于go的嵌入式领域，支持的功能有限，甚至json序列化都无法支持
//2.函数命名遵从驼峰命名
//3.提供的sql生成语句中，关于bool值需要特殊处理，create语句的设计的到bool值的我会在创建语句中单独说明，这是因为在原有的sqlite中并不支持bool，用整数1或者0替代，gorm对其做了转换。
//4.提供的sql生成语句中，字段名是下划线方式 例如：recv_id，但是接口转换的数据json tag字段的风格是recvID，类似的所有的字段需要做个map映射

type LocalChatLogs struct {
	loginUserID string
}

func NewLocalChatLogs(loginUserID string) *LocalChatLogs {
	return &LocalChatLogs{loginUserID: loginUserID}
}

func (i *LocalChatLogs) GetMessage(clientMsgID string) (*model_struct.LocalChatLog, error) {
	msg, err := Exec(clientMsgID)
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
func (i *LocalChatLogs) GetSendingMessageList() (result []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec()
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
func (i *LocalChatLogs) UpdateMessage(c *model_struct.LocalChatLog) error {
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
	_, err := Exec(c.ClientMsgID, utils.StructToJsonString(tempLocalChatLog))
	return err
}
func (i *LocalChatLogs) GetNormalMsgSeq() (uint32, error) {
	seq, err := Exec()
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
func (i *LocalChatLogs) BatchInsertMessageList(messageList []*model_struct.LocalChatLog) error {
	_, err := Exec(utils.StructToJsonString(messageList))
	return err
}
func (i *LocalChatLogs) InsertMessage(message *model_struct.LocalChatLog) error {
	_, err := Exec(utils.StructToJsonString(message))
	return err
}
func (i *LocalChatLogs) GetAllUnDeleteMessageSeqList() ([]uint32, error) {
	panic("implement me")
}
func (i *LocalChatLogs) UpdateColumnsMessageList(clientMsgIDList []string, args map[string]interface{}) error {
	_, err := Exec(utils.StructToJsonString(clientMsgIDList), args)
	return err
}
func (i *LocalChatLogs) UpdateColumnsMessage(clientMsgID string, args map[string]interface{}) error {
	_, err := Exec(clientMsgID, utils.StructToJsonString(args))
	return err
}
func (i *LocalChatLogs) DeleteAllMessage() error {
	panic("implement me")
}
func (i *LocalChatLogs) UpdateMessageStatusBySourceID(sourceID string, status, sessionType int32) error {
	_, err := Exec(sourceID, status, sessionType)
	return err
}
func (i *LocalChatLogs) UpdateMessageTimeAndStatus(clientMsgID string, serverMsgID string, sendTime int64, status int32) error {
	_, err := Exec(clientMsgID, serverMsgID, sendTime, status)
	return err
}
func (i *LocalChatLogs) GetMessageList(sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec(sourceID, sessionType, count, startTime, isReverse, i.loginUserID)
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
func (i *LocalChatLogs) GetMessageListNoTime(sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec(sourceID, sessionType, count, isReverse, i.loginUserID)
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
func (i *LocalChatLogs) UpdateSingleMessageHasRead(sendID string, msgIDList []string) error {
	panic("implement me")
}

func (i *LocalChatLogs) SearchMessageByContentType(contentType []int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (messages []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec(contentType, sourceID, startTime, endTime, sessionType, offset, count)
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

func (i *LocalChatLogs) SuperGroupSearchMessageByContentType(contentType []int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (messages []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec(contentType, sourceID, startTime, endTime, sessionType, offset, count)
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

func (i *LocalChatLogs) SearchMessageByContentTypeAndKeyword(contentType []int, keywordList []string, keywordListMatchType int, startTime, endTime int64) (result []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec(contentType, keywordList, keywordListMatchType, startTime, endTime)
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

func (i *LocalChatLogs) MessageIfExists(clientMsgID string) (bool, error) {
	isExist, err := Exec(clientMsgID)
	if err != nil {
		return false, err
	} else {
		if v, ok := isExist.(bool); ok {
			return v, nil
		} else {
			return false, ErrType
		}
	}
}

func (i *LocalChatLogs) IsExistsInErrChatLogBySeq(seq int64) bool {
	isExist, err := Exec(seq)
	if err != nil {
		return false
	} else {
		if v, ok := isExist.(bool); ok {
			return v
		} else {
			return false
		}
	}
}

func (i *LocalChatLogs) MessageIfExistsBySeq(seq int64) (bool, error) {
	isExist, err := Exec(seq)
	if err != nil {
		return false, err
	} else {
		if v, ok := isExist.(bool); ok {
			return v, nil
		} else {
			return false, ErrType
		}
	}
}

func (i *LocalChatLogs) UpdateGroupMessageHasRead(msgIDList []string, sessionType int32) error {
	_, err := Exec(msgIDList, sessionType)
	return err
}

func (i *LocalChatLogs) GetMultipleMessage(msgIDList []string) (result []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec(msgIDList)
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

func (i *LocalChatLogs) GetLostMsgSeqList(minSeqInSvr uint32) (result []uint32, err error) {
	l, err := Exec(minSeqInSvr)
	if err != nil {
		return nil, err
	} else {
		if v, ok := l.([]float64); ok {
			for _, v := range v {
				v1 := uint32(v)
				result = append(result, v1)
			}
			return result, err
		} else {
			return nil, ErrType
		}
	}
}

func (i *LocalChatLogs) GetTestMessage(seq uint32) (*model_struct.LocalChatLog, error) {
	msg, err := Exec(seq)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msg.(model_struct.LocalChatLog); ok {
			return &v, err
		} else {
			return nil, ErrType
		}
	}
}

func (i *LocalChatLogs) UpdateMsgSenderNickname(sendID, nickname string, sType int) error {
	_, err := Exec(sendID, nickname, sType)
	return err
}

func (i *LocalChatLogs) UpdateMsgSenderFaceURL(sendID, faceURL string, sType int) error {
	_, err := Exec(sendID, faceURL, sType)
	return err
}

func (i *LocalChatLogs) UpdateMsgSenderFaceURLAndSenderNickname(sendID, faceURL, nickname string, sessionType int) error {
	_, err := Exec(sendID, faceURL, nickname, sessionType)
	return err
}

func (i *LocalChatLogs) GetMsgSeqByClientMsgID(clientMsgID string) (uint32, error) {
	result, err := Exec(clientMsgID)
	if err != nil {
		return 0, err
	}
	if v, ok := result.(float64); ok {
		return uint32(v), nil
	}
	return 0, ErrType
}

func (i *LocalChatLogs) SearchAllMessageByContentType(contentType int) (result []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec(contentType)
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

func (i *LocalChatLogs) GetAbnormalMsgSeq() (uint32, error) {
	return 0, nil
}

func (i *LocalChatLogs) GetAbnormalMsgSeqList() ([]uint32, error) {
	panic("implement me")
}

func (i *LocalChatLogs) BatchInsertExceptionMsg(MessageList []*model_struct.LocalErrChatLog) error {
	panic("implement me")
}

func (i *LocalChatLogs) SetChatLogFailedStatus() {
	panic("implement me")
}

func (i *LocalChatLogs) GetMsgSeqListByGroupID(groupID string) (result []uint32, err error) {
	l, err := Exec(groupID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := l.(string); ok {
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return result, err
		} else {
			return nil, ErrType
		}
	}
}

func (i *LocalChatLogs) GetMsgSeqListByPeerUserID(userID string) (result []uint32, err error) {
	l, err := Exec(userID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := l.(string); ok {
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return result, err
		} else {
			return nil, ErrType
		}
	}
}

func (i *LocalChatLogs) GetMsgSeqListBySelfUserID(userID string) (result []uint32, err error) {
	l, err := Exec(userID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := l.(string); ok {
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return result, err
		} else {
			return nil, ErrType
		}
	}
}
