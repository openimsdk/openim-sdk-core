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

//go:build js && wasm
// +build js,wasm

package indexdb

import (
	"context"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/indexdb/temp_struct"
)

type LocalChatLogs struct {
	loginUserID string
}

func NewLocalChatLogs(loginUserID string) *LocalChatLogs {
	return &LocalChatLogs{loginUserID: loginUserID}
}

func (i *LocalChatLogs) GetMessage(ctx context.Context, conversationID, clientMsgID string) (*model_struct.LocalChatLog, error) {
	msg, err := Exec(conversationID, clientMsgID)
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
func (i *LocalChatLogs) GetSendingMessageList(ctx context.Context) (result []*model_struct.LocalChatLog, err error) {
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
func (i *LocalChatLogs) UpdateMessage(ctx context.Context, conversationID string, c *model_struct.LocalChatLog) error {
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
	_, err := Exec(conversationID, c.ClientMsgID, utils.StructToJsonString(tempLocalChatLog))
	return err
}
func (i *LocalChatLogs) GetNormalMsgSeq(ctx context.Context) (uint32, error) {
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
func (i *LocalChatLogs) BatchInsertMessageList(ctx context.Context, conversationID string, messageList []*model_struct.LocalChatLog) error {
	_, err := Exec(conversationID, utils.StructToJsonString(messageList))
	return err
}
func (i *LocalChatLogs) InsertMessage(ctx context.Context, conversationID string, message *model_struct.LocalChatLog) error {
	_, err := Exec(conversationID, utils.StructToJsonString(message))
	return err
}
func (i *LocalChatLogs) GetAllUnDeleteMessageSeqList(ctx context.Context) (result []uint32, err error) {
	l, err := Exec()
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
func (i *LocalChatLogs) UpdateColumnsMessageList(ctx context.Context, clientMsgIDList []string, args map[string]interface{}) error {
	_, err := Exec(utils.StructToJsonString(clientMsgIDList), args)
	return err
}
func (i *LocalChatLogs) UpdateColumnsMessage(ctx context.Context, clientMsgID string, args map[string]interface{}) error {
	_, err := Exec(clientMsgID, utils.StructToJsonString(args))
	return err
}
func (i *LocalChatLogs) DeleteAllMessage(ctx context.Context) error {
	_, err := Exec()
	return err
}
func (i *LocalChatLogs) UpdateMessageStatusBySourceID(ctx context.Context, sourceID string, status, sessionType int32) error {
	_, err := Exec(sourceID, status, sessionType, i.loginUserID)
	return err
}
func (i *LocalChatLogs) UpdateMessageTimeAndStatus(ctx context.Context, conversationID, clientMsgID string, serverMsgID string, sendTime int64, status int32) error {
	_, err := Exec(clientMsgID, serverMsgID, sendTime, status)
	return err
}
func (i *LocalChatLogs) GetMessageList(ctx context.Context, conversationID string, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec(conversationID, count, startTime, isReverse, i.loginUserID)
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
func (i *LocalChatLogs) GetMessageListNoTime(ctx context.Context, conversationID string, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec(conversationID, count, isReverse, i.loginUserID)
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
func (i *LocalChatLogs) UpdateSingleMessageHasRead(ctx context.Context, sendID string, msgIDList []string) error {
	_, err := Exec(sendID, utils.StructToJsonString(msgIDList))
	return err
}

func (i *LocalChatLogs) SearchMessageByContentType(ctx context.Context, contentType []int, conversationID string, startTime, endTime int64, offset, count int) (messages []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec(conversationID, utils.StructToJsonString(contentType), startTime, endTime, offset, count)
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

func (i *LocalChatLogs) SuperGroupSearchMessageByContentType(ctx context.Context, contentType []int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (messages []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec(utils.StructToJsonString(contentType), sourceID, startTime, endTime, sessionType, offset, count)
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

func (i *LocalChatLogs) SearchMessageByContentTypeAndKeyword(ctx context.Context, contentType []int, conversationID string, keywordList []string, keywordListMatchType int, startTime, endTime int64) (result []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec(conversationID, utils.StructToJsonString(contentType), utils.StructToJsonString(keywordList), keywordListMatchType, startTime, endTime)
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

func (i *LocalChatLogs) MessageIfExists(ctx context.Context, clientMsgID string) (bool, error) {
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

func (i *LocalChatLogs) IsExistsInErrChatLogBySeq(ctx context.Context, seq int64) bool {
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

func (i *LocalChatLogs) MessageIfExistsBySeq(ctx context.Context, seq int64) (bool, error) {
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

func (i *LocalChatLogs) GetMultipleMessage(ctx context.Context, msgIDList []string) (result []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec(utils.StructToJsonString(msgIDList))
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

func (i *LocalChatLogs) GetLostMsgSeqList(ctx context.Context, minSeqInSvr uint32) (result []uint32, err error) {
	l, err := Exec(minSeqInSvr)
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

func (i *LocalChatLogs) GetTestMessage(ctx context.Context, seq uint32) (*model_struct.LocalChatLog, error) {
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

func (i *LocalChatLogs) UpdateMsgSenderNickname(ctx context.Context, sendID, nickname string, sType int) error {
	_, err := Exec(sendID, nickname, sType)
	return err
}

func (i *LocalChatLogs) UpdateMsgSenderFaceURL(ctx context.Context, sendID, faceURL string, sType int) error {
	_, err := Exec(sendID, faceURL, sType)
	return err
}

func (i *LocalChatLogs) UpdateMsgSenderFaceURLAndSenderNickname(ctx context.Context, conversationID, sendID, faceURL, nickname string) error {
	_, err := Exec(conversationID, sendID, faceURL, nickname)
	return err
}

func (i *LocalChatLogs) GetMsgSeqByClientMsgID(ctx context.Context, clientMsgID string) (uint32, error) {
	result, err := Exec(clientMsgID)
	if err != nil {
		return 0, err
	}
	if v, ok := result.(float64); ok {
		return uint32(v), nil
	}
	return 0, ErrType
}
func (i *LocalChatLogs) SearchAllMessageByContentType(ctx context.Context, contentType int) (result []*model_struct.LocalChatLog, err error) {
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

func (i *LocalChatLogs) GetMsgSeqListByGroupID(ctx context.Context, groupID string) (result []uint32, err error) {
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

func (i *LocalChatLogs) GetMsgSeqListByPeerUserID(ctx context.Context, userID string) (result []uint32, err error) {
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

func (i *LocalChatLogs) GetMsgSeqListBySelfUserID(ctx context.Context, userID string) (result []uint32, err error) {
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

func (i *LocalChatLogs) GetAbnormalMsgSeq(ctx context.Context) (int64, error) {
	result, err := Exec()
	if err != nil {
		return 0, err
	}
	if v, ok := result.(float64); ok {
		return int64(v), nil
	}
	return 0, ErrType
}

func (i *LocalChatLogs) GetAbnormalMsgSeqList(ctx context.Context) (result []int64, err error) {
	l, err := Exec()
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

func (i *LocalChatLogs) BatchInsertExceptionMsg(ctx context.Context, MessageList []*model_struct.LocalErrChatLog) error {
	_, err := Exec(utils.StructToJsonString(MessageList))
	return err
}

func (i *LocalChatLogs) UpdateGroupMessageHasRead(ctx context.Context, msgIDList []string, sessionType int32) error {
	_, err := Exec(utils.StructToJsonString(msgIDList), sessionType)
	return err
}

func (i *LocalChatLogs) SearchMessageByKeyword(ctx context.Context, contentType []int, keywordList []string, keywordListMatchType int, conversationID string, startTime, endTime int64, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	msgList, err := Exec(utils.StructToJsonString(contentType), utils.StructToJsonString(keywordList), keywordListMatchType, conversationID, startTime, endTime, offset, count)
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

func (i *LocalChatLogs) GetSuperGroupAbnormalMsgSeq(ctx context.Context, groupID string) (uint32, error) {
	isExist, err := Exec(groupID)
	if err != nil {
		return 0, err
	} else {
		if v, ok := isExist.(uint32); ok {
			return v, nil
		} else {
			return 0, ErrType
		}
	}
}
func (i *LocalChatLogs) GetAlreadyExistSeqList(ctx context.Context, conversationID string, lostSeqList []int64) (seqList []int64, err error) {
	return nil, nil
}
func (i *LocalChatLogs) GetMessageBySeq(ctx context.Context, conversationID string, seq int64) (*model_struct.LocalChatLog, error) {
	//TODO implement me
	panic("implement me")
}

func (i *LocalChatLogs) UpdateMessageBySeq(ctx context.Context, conversationID string, c *model_struct.LocalChatLog) error {
	//TODO implement me
	panic("implement me")
}

func (i *LocalChatLogs) UpdateMessageByClientMsgID(ctx context.Context, clientMsgID string, data map[string]any) error {
	//TODO implement me
	panic("implement me")
}

func (i *LocalChatLogs) MarkConversationMessageAsRead(ctx context.Context, conversationID string, msgIDs []string) (rowsAffected int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (i *LocalChatLogs) MarkConversationMessageAsReadBySeqs(ctx context.Context, conversationID string, seqs []int64) (rowsAffected int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (i *LocalChatLogs) GetUnreadMessage(ctx context.Context, conversationID string) (result []*model_struct.LocalChatLog, err error) {
	//TODO implement me
	panic("implement me")
}

func (i *LocalChatLogs) MarkConversationAllMessageAsRead(ctx context.Context, conversationID string) (rowsAffected int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (i *LocalChatLogs) GetMessagesByClientMsgIDs(ctx context.Context, conversationID string, msgIDs []string) (result []*model_struct.LocalChatLog, err error) {
	//TODO implement me
	panic("implement me")
}

func (i *LocalChatLogs) GetMessagesBySeqs(ctx context.Context, conversationID string, seqs []int64) (result []*model_struct.LocalChatLog, err error) {
	//TODO implement me
	panic("implement me")
}

func (i *LocalChatLogs) GetConversationNormalMsgSeq(ctx context.Context, conversationID string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (i *LocalChatLogs) GetConversationPeerNormalMsgSeq(ctx context.Context, conversationID string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (i *LocalChatLogs) GetConversationAbnormalMsgSeq(ctx context.Context, groupID string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (i *LocalChatLogs) DeleteConversationAllMessages(ctx context.Context, conversationID string) error {
	//TODO implement me
	panic("implement me")
}

func (i *LocalChatLogs) MarkDeleteConversationAllMessages(ctx context.Context, conversationID string) error {
	//TODO implement me
	panic("implement me")
}

func (i *LocalChatLogs) SuperGroupGetNormalMsgSeq(ctx context.Context) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (i *LocalChatLogs) SuperGroupGetNormalMinSeq(ctx context.Context, groupID string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (i *LocalChatLogs) SuperGroupGetTestMessage(ctx context.Context, seq int64) (*model_struct.LocalChatLog, error) {
	//TODO implement me
	panic("implement me")
}

func (i *LocalChatLogs) DeleteConversationMsgs(ctx context.Context, conversationID string, msgIDs []string) error {
	//TODO implement me
	panic("implement me")
}

func (i *LocalChatLogs) DeleteConversationMsgsBySeqs(ctx context.Context, conversationID string, seqs []int64) error {
	//TODO implement me
	panic("implement me")
}
