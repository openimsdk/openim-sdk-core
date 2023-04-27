// Copyright © 2023 OpenIM SDK.
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

import "context"

import (
	"errors"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"runtime"
	"syscall/js"
	"time"
)

//1,使用wasm原生的方式，tinygo应用于go的嵌入式领域，支持的功能有限，支持go语言的子集,甚至json序列化都无法支持
//2.函数命名遵从驼峰命名
//3.提供的sql生成语句中，关于bool值需要特殊处理，create语句的设计的到bool值的需要在创建语句中单独说明，这是因为在原有的sqlite中并不支持bool，用整数1或者0替代，gorm对其做了转换。
//4.提供的sql生成语句中，字段名是下划线方式 例如：recv_id，但是接口转换的数据json tag字段的风格是recvID，类似的所有的字段需要做个map映射
//5.任何涉及到gorm获取的是否需要返回错误，比如take和find都需要在文档上说明
//6.任何涉及到update的操作，一定要看gorm原型中实现，如果有select(*)则不需要用temp_struct中的结构体
//7.任何和接口重名的时候，db接口统一加上后缀DB
//8.任何map类型统一使用json字符串转换，文档说明

type IndexDB struct {
	LocalUsers
	LocalConversations
	*LocalChatLogs
	LocalSuperGroupChatLogs
	LocalSuperGroup
	LocalConversationUnreadMessages
	LocalGroups
	LocalGroupMember
	LocalGroupRequest
	LocalCacheMessage
	*FriendRequest
	*Black
	*Friend
	LocalChatLogReactionExtensions
	loginUserID string
}

func (i IndexDB) SearchMessageByKeywordController(ctx context.Context, contentType []int, keywordList []string, keywordListMatchType int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	switch sessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupSearchMessageByKeyword(contentType, keywordList, keywordListMatchType, sourceID, startTime, endTime, sessionType, offset, count)
	default:
		return i.SearchMessageByKeyword(contentType, keywordList, keywordListMatchType, sourceID, startTime, endTime, sessionType, offset, count)
	}
}

func (i IndexDB) SearchMessageByContentTypeController(ctx context.Context, contentType []int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	switch sessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupSearchMessageByContentType(contentType, sourceID, startTime, endTime, sessionType, offset, count)
	default:
		return i.SearchMessageByContentType(contentType, sourceID, startTime, endTime, sessionType, offset, count)
	}
}

func (i IndexDB) SearchMessageByContentTypeAndKeywordController(ctx context.Context, contentType []int, keywordList []string, keywordListMatchType int, startTime, endTime int64, operationID string) (result []*model_struct.LocalChatLog, err error) {
	list, err := i.SearchMessageByContentTypeAndKeyword(contentType, keywordList, keywordListMatchType, startTime, endTime)
	if err != nil {
		return nil, err
	}
	superGroupIDList, err := i.GetJoinedSuperGroupIDList()
	if err != nil {
		return nil, err
	}
	for _, v := range superGroupIDList {
		sList, err := i.SuperGroupSearchMessageByContentTypeAndKeyword(contentType, keywordList, keywordListMatchType, startTime, endTime, v)
		if err != nil {
			log.Error(operationID, "search message in group err", err.Error(), v)
			continue
		}
		list = append(list, sList...)
	}
	workingGroupIDList, err := i.GetJoinedWorkingGroupIDList()
	if err != nil {
		return nil, err
	}
	for _, v := range workingGroupIDList {
		sList, err := i.SuperGroupSearchMessageByContentTypeAndKeyword(contentType, keywordList, keywordListMatchType, startTime, endTime, v)
		if err != nil {
			log.Error(operationID, "search message in group err", err.Error(), v)
			continue
		}
		list = append(list, sList...)
	}
	return list, nil
}
func (i IndexDB) UpdateMsgSenderFaceURLAndSenderNicknameController(ctx context.Context, sendID, faceURL, nickname string, sessionType int, groupID string) error {
	switch sessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupUpdateMsgSenderFaceURLAndSenderNickname(sendID, faceURL, nickname, sessionType, groupID)
	default:
		return i.UpdateMsgSenderFaceURLAndSenderNickname(sendID, faceURL, nickname, sessionType)
	}
}

type CallbackData struct {
	ErrCode int32       `json:"errCode"`
	ErrMsg  string      `json:"errMsg"`
	Data    interface{} `json:"data"`
}

const TIMEOUT = 5

var ErrTimoutFromJavaScript = errors.New("invoke javascript timeout，maybe should check  function from javascript")
var jsErr = js.Global().Get("Error")

func Exec(args ...interface{}) (output interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = utils.Wrap(errors.New(x), "")
			case error:
				err = x
			default:
				err = utils.Wrap(errors.New("unknown panic"), "")
			}
		}
	}()
	thenChannel := make(chan []js.Value)
	defer close(thenChannel)
	catchChannel := make(chan []js.Value)
	defer close(catchChannel)
	pc, _, _, _ := runtime.Caller(1)
	funcName := utils.CleanUpfuncName(runtime.FuncForPC(pc).Name())
	data := CallbackData{}
	thenFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case string:
					err = utils.Wrap(errors.New(x), "")
				case error:
					err = x
				default:
					err = utils.Wrap(errors.New("unknown panic"), "")
				}
			}
		}()
		log.Debug("js then func", "=> (main go context) "+funcName+" with respone ", args[0].String())
		thenChannel <- args
		return nil
	})
	defer thenFunc.Release()
	catchFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case string:
					err = utils.Wrap(errors.New(x), "")
				case error:
					err = x
				default:
					err = utils.Wrap(errors.New("unknown panic"), "")
				}
			}
		}()
		log.Debug("js catch func", "=> (main go context) "+funcName+" with respone ", args[0].String())
		catchChannel <- args
		return nil
	})
	defer catchFunc.Release()
	js.Global().Call(utils.FirstLower(funcName), args...).Call("then", thenFunc).Call("catch", catchFunc)
	select {
	case result := <-thenChannel:
		interErr := utils.JsonStringToStruct(result[0].String(), &data)
		if interErr != nil {
			err = utils.Wrap(err, "return json unmarshal err from javascript")
		}
	case catch := <-catchChannel:
		if catch[0].InstanceOf(jsErr) {
			return nil, js.Error{Value: catch[0]}
		} else {
			panic("unknown javascript exception")
		}
	case <-time.After(TIMEOUT * time.Second):
		panic(ErrTimoutFromJavaScript)
	}
	if data.ErrCode != 0 {
		return "", errors.New(data.ErrMsg)
	}
	return data.Data, err
}

func (i IndexDB) BatchInsertMessageListController(ctx context.Context, MessageList []*model_struct.LocalChatLog) error {
	if len(MessageList) == 0 {
		return nil
	}
	switch MessageList[len(MessageList)-1].SessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupBatchInsertMessageList(MessageList, MessageList[len(MessageList)-1].RecvID)
	default:
		return i.BatchInsertMessageList(MessageList)
	}
}

func (i IndexDB) InsertMessageController(ctx context.Context, message *model_struct.LocalChatLog) error {
	switch message.SessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupInsertMessage(message, message.RecvID)
	default:
		return i.InsertMessage(message)
	}
}

func (i IndexDB) BatchUpdateMessageList(ctx context.Context, MessageList []*model_struct.LocalChatLog) error {
	if MessageList == nil {
		return nil
	}
	for _, v := range MessageList {
		v1 := new(model_struct.LocalChatLog)
		v1.ClientMsgID = v.ClientMsgID
		v1.Seq = v.Seq
		v1.Status = v.Status
		v1.RecvID = v.RecvID
		v1.SessionType = v.SessionType
		v1.ServerMsgID = v.ServerMsgID
		v1.Ex = v.Ex
		err := i.UpdateMessageController(v1)
		if err != nil {
			return utils.Wrap(err, "BatchUpdateMessageList failed")
		}

	}
	return nil
}

func (i IndexDB) BatchSpecialUpdateMessageList(ctx context.Context, MessageList []*model_struct.LocalChatLog) error {
	if MessageList == nil {
		return nil
	}
	for _, v := range MessageList {
		v1 := new(model_struct.LocalChatLog)
		v1.ClientMsgID = v.ClientMsgID
		v1.ServerMsgID = v.ServerMsgID
		v1.SendID = v.SendID
		v1.RecvID = v.RecvID
		v1.SenderPlatformID = v.SenderPlatformID
		v1.SenderNickname = v.SenderNickname
		v1.SenderFaceURL = v.SenderFaceURL
		v1.SessionType = v.SessionType
		v1.MsgFrom = v.MsgFrom
		v1.ContentType = v.ContentType
		v1.Content = v.Content
		v1.Seq = v.Seq
		v1.SendTime = v.SendTime
		v1.CreateTime = v.CreateTime
		v1.AttachedInfo = v.AttachedInfo
		v1.Ex = v.Ex
		err := i.UpdateMessageController(v1)
		if err != nil {
			log.Error("", "update single message failed", *v)
			return utils.Wrap(err, "BatchUpdateMessageList failed")
		}

	}
	return nil
}

func (i IndexDB) GetMessageController(ctx context.Context, msg *sdk_struct.MsgStruct) (*model_struct.LocalChatLog, error) {
	switch msg.SessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupGetMessage(msg)
	default:
		return i.GetMessage(msg.ClientMsgID)
	}
}
func (i IndexDB) UpdateColumnsMessageController(ctx context.Context, ClientMsgID string, groupID string, sessionType int32, args map[string]interface{}) error {
	switch sessionType {
	case constant.SuperGroupChatType:
		return utils.Wrap(i.SuperGroupUpdateColumnsMessage(ClientMsgID, groupID, args), "")
	default:
		return utils.Wrap(i.UpdateColumnsMessage(ClientMsgID, args), "")
	}
}
func (i IndexDB) UpdateMessageController(ctx context.Context, c *model_struct.LocalChatLog) error {
	switch c.SessionType {
	case constant.SuperGroupChatType:
		return utils.Wrap(i.SuperGroupUpdateMessage(c), "")
	default:
		return utils.Wrap(i.UpdateMessage(c), "")
	}
}

func (i IndexDB) UpdateMessageStatusBySourceIDController(ctx context.Context, sourceID string, status, sessionType int32) error {
	switch sessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupUpdateMessageStatusBySourceID(sourceID, status, sessionType)
	default:
		return i.UpdateMessageStatusBySourceID(sourceID, status, sessionType)
	}
}
func (i IndexDB) UpdateMessageTimeAndStatusController(ctx context.Context, msg *sdk_struct.MsgStruct) error {
	switch msg.SessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupUpdateMessageTimeAndStatus(msg)
	default:
		return i.UpdateMessageTimeAndStatus(msg.ClientMsgID, msg.ServerMsgID, msg.SendTime, msg.Status)
	}
}
func (i IndexDB) GetMessageListController(ctx context.Context, sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	switch sessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupGetMessageList(sourceID, sessionType, count, startTime, isReverse)
	default:
		return i.GetMessageList(sourceID, sessionType, count, startTime, isReverse)
	}
}
func (i IndexDB) GetMessageListNoTimeController(ctx context.Context, sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	switch sessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupGetMessageListNoTime(sourceID, sessionType, count, isReverse)
	default:
		return i.GetMessageListNoTime(sourceID, sessionType, count, isReverse)
	}
}

func (i IndexDB) UpdateGroupMessageHasReadController(ctx context.Context, msgIDList []string, groupID string, sessionType int32) error {
	switch sessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupUpdateGroupMessageHasRead(msgIDList, groupID)
	default:
		return i.UpdateGroupMessageHasRead(msgIDList, sessionType)
	}
}

func (i IndexDB) GetMultipleMessageController(ctx context.Context, msgIDList []string, groupID string, sessionType int32) (result []*model_struct.LocalChatLog, err error) {
	switch sessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupGetMultipleMessage(msgIDList, groupID)
	default:
		return i.GetMultipleMessage(msgIDList)
	}
}

func (i IndexDB) GetMsgSeqByClientMsgIDController(ctx context.Context, m *sdk_struct.MsgStruct) (uint32, error) {
	switch m.SessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupGetMsgSeqByClientMsgID(m.ClientMsgID, m.GroupID)
	default:
		return i.GetMsgSeqByClientMsgID(m.ClientMsgID)
	}
}

func (i IndexDB) GetSubDepartmentList(ctx context.Context, departmentID string, args ...int) ([]*model_struct.LocalDepartment, error) {
	panic("implement me")
}

func (i IndexDB) InsertDepartment(ctx context.Context, department *model_struct.LocalDepartment) error {
	panic("implement me")
}

func (i IndexDB) UpdateDepartment(ctx context.Context, department *model_struct.LocalDepartment) error {
	panic("implement me")
}

func (i IndexDB) DeleteDepartment(ctx context.Context, departmentID string) error {
	panic("implement me")
}

func (i IndexDB) GetDepartmentInfo(ctx context.Context, departmentID string) (*model_struct.LocalDepartment, error) {
	panic("implement me")
}

func (i IndexDB) GetAllDepartmentList(ctx context.Context) ([]*model_struct.LocalDepartment, error) {
	panic("implement me")
}

func (i IndexDB) GetParentDepartmentList(ctx context.Context, departmentID string) ([]*model_struct.LocalDepartment, error) {
	panic("implement me")
}

func (i IndexDB) GetDepartmentList(ctx context.Context, departmentList *[]*model_struct.LocalDepartment, departmentID string) error {
	panic("implement me")
}

func (i IndexDB) GetParentDepartment(ctx context.Context, departmentID string) (model_struct.LocalDepartment, error) {
	panic("implement me")
}

func (i IndexDB) SearchDepartmentMember(ctx context.Context, keyWord string, isSearchUserName, isSearchEmail, isSearchMobile, isSearchPosition, isSearchTelephone, isSearchUserEnglishName, isSearchUserID bool, offset, count int) ([]*model_struct.SearchDepartmentMemberResult, error) {
	panic("implement me")
}

func (i IndexDB) SearchDepartment(ctx context.Context, keyWord string, offset, count int) ([]*model_struct.LocalDepartment, error) {
	panic("implement me")
}

func (i IndexDB) GetJoinedSuperGroupIDList(ctx context.Context) ([]string, error) {
	groupIDList, err := Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := groupIDList.(string); ok {
			var result []string
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

func (i IndexDB) GetReadDiffusionGroupIDList(ctx context.Context) ([]string, error) {
	g1, err1 := i.GetJoinedSuperGroupIDList()
	g2, err2 := i.GetJoinedWorkingGroupIDList()
	var groupIDList []string
	if err1 == nil {
		groupIDList = append(groupIDList, g1...)
	}
	if err2 == nil {
		groupIDList = append(groupIDList, g2...)
	}
	var err error
	if err1 != nil {
		err = err1
	}
	if err2 != nil {
		err = err2
	}
	return groupIDList, err
}

func (i IndexDB) BatchInsertExceptionMsgController(ctx context.Context, MessageList []*model_struct.LocalErrChatLog) error {
	if len(MessageList) == 0 {
		return nil
	}
	switch MessageList[len(MessageList)-1].SessionType {
	case constant.SuperGroupChatType:
		return i.SuperBatchInsertExceptionMsg(MessageList, MessageList[len(MessageList)-1].RecvID)
	default:
		return i.BatchInsertExceptionMsg(MessageList)
	}
}

func (i IndexDB) GetDepartmentMemberListByDepartmentID(ctx context.Context, departmentID string, args ...int) ([]*model_struct.LocalDepartmentMember, error) {
	panic("implement me")
}

func (i IndexDB) GetAllDepartmentMemberList(ctx context.Context) ([]*model_struct.LocalDepartmentMember, error) {
	panic("implement me")
}

func (i IndexDB) InsertDepartmentMember(ctx context.Context, departmentMember *model_struct.LocalDepartmentMember) error {
	panic("implement me")
}

func (i IndexDB) BatchInsertDepartmentMember(ctx context.Context, departmentMemberList []*model_struct.LocalDepartmentMember) error {
	panic("implement me")
}

func (i IndexDB) UpdateDepartmentMember(ctx context.Context, departmentMember *model_struct.LocalDepartmentMember) error {
	panic("implement me")
}

func (i IndexDB) DeleteDepartmentMember(ctx context.Context, departmentID string, userID string) error {
	panic("implement me")
}

func (i IndexDB) GetDepartmentMemberListByUserID(ctx context.Context, userID string) ([]*model_struct.LocalDepartmentMember, error) {
	panic("implement me")
}

func (i IndexDB) InsertWorkMomentsNotification(ctx context.Context, jsonDetail string) error {
	panic("implement me")
}

func (i IndexDB) GetWorkMomentsNotification(ctx context.Context, offset, count int) (WorkMomentsNotifications []*model_struct.LocalWorkMomentsNotification, err error) {
	panic("implement me")
}

func (i IndexDB) GetWorkMomentsNotificationLimit(ctx context.Context, pageNumber, showNumber int) (WorkMomentsNotifications []*model_struct.LocalWorkMomentsNotification, err error) {
	panic("implement me")
}

func (i IndexDB) InitWorkMomentsNotificationUnreadCount(ctx context.Context) error {
	panic("implement me")
}

func (i IndexDB) IncrWorkMomentsNotificationUnreadCount(ctx context.Context) error {
	panic("implement me")
}

func (i IndexDB) MarkAllWorkMomentsNotificationAsRead(ctx context.Context) (err error) {
	panic("implement me")
}

func (i IndexDB) GetWorkMomentsUnReadCount(ctx context.Context) (workMomentsNotificationUnReadCount model_struct.LocalWorkMomentsNotificationUnreadCount, err error) {
	panic("implement me")
}

func (i IndexDB) ClearWorkMomentsNotification(ctx context.Context) (err error) {
	panic("implement me")
}

func (i IndexDB) Close(ctx context.Context) error {
	_, err := Exec()
	return err
}

func (i IndexDB) InitDB(ctx context.Context, userID string, dataDir string) error {
	_, err := Exec(userID, dataDir)
	return err
}

//func (i IndexDB) GetBlackList(ctx context.Context, ) ([]*model_struct.LocalBlack, error) {
//	panic("implement me")
//}
//
//func (i IndexDB) GetBlackListUserID(ctx context.Context, ) (blackListUid []string, err error) {
//	panic("implement me")
//}
//
//func (i IndexDB) GetBlackInfoByBlockUserID(ctx context.Context, blockUserID string) (*model_struct.LocalBlack, error) {
//	panic("implement me")
//}
//
//func (i IndexDB) GetBlackInfoList(ctx context.Context, blockUserIDList []string) ([]*model_struct.LocalBlack, error) {
//	panic("implement me")
//}
//
//func (i IndexDB) InsertBlack(ctx context.Context, black *model_struct.LocalBlack) error {
//	panic("implement me")
//}
//
//func (i IndexDB) UpdateBlack(ctx context.Context, black *model_struct.LocalBlack) error {
//	panic("implement me")
//}
//
//func (i IndexDB) DeleteBlack(ctx context.Context, blockUserID string) error {
//	panic("implement me")
//}

func NewIndexDB(loginUserID string) *IndexDB {
	return &IndexDB{
		LocalChatLogs: NewLocalChatLogs(loginUserID),
		FriendRequest: NewFriendRequest(loginUserID),
		Black:         NewBlack(loginUserID),
		Friend:        NewFriend(loginUserID),
		loginUserID:   loginUserID,
	}
}

func (i IndexDB) SetChatLogFailedStatus(ctx context.Context) {
	msgList, err := i.GetSendingMessageList()
	if err != nil {
		log.Error("", "GetSendingMessageList failed ", err.Error())
		return
	}
	for _, v := range msgList {
		v.Status = constant.MsgStatusSendFailed
		err := i.UpdateMessage(v)
		if err != nil {
			log.Error("", "UpdateMessage failed ", err.Error(), v)
			continue
		}
	}
	groupIDList, err := i.GetReadDiffusionGroupIDList()
	if err != nil {
		log.Error("", "GetReadDiffusionGroupIDList failed ", err.Error())
		return
	}
	for _, v := range groupIDList {
		msgList, err := i.SuperGroupGetSendingMessageList(v)
		if err != nil {
			log.Error("", "GetSendingMessageList failed ", err.Error())
			return
		}
		if len(msgList) > 0 {
			for _, v := range msgList {
				v.Status = constant.MsgStatusSendFailed
				err := i.SuperGroupUpdateMessage(v)
				if err != nil {
					log.Error("", "UpdateMessage failed ", err.Error(), v)
					continue
				}
			}
		}

	}
}
