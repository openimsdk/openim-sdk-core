package indexdb

import (
	"errors"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
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
	loginUserID string
}

func (i IndexDB) SearchMessageByKeywordController(contentType []int, keywordList []string, keywordListMatchType int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	switch sessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupSearchMessageByKeyword(contentType, keywordList, keywordListMatchType, sourceID, startTime, endTime, sessionType, offset, count)
	default:
		return i.SearchMessageByKeyword(contentType, keywordList, keywordListMatchType, sourceID, startTime, endTime, sessionType, offset, count)
	}
}

func (i IndexDB) SearchMessageByContentTypeController(contentType []int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	switch sessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupSearchMessageByContentType(contentType, sourceID, startTime, endTime, sessionType, offset, count)
	default:
		return i.SearchMessageByContentType(contentType, sourceID, startTime, endTime, sessionType, offset, count)
	}
}

func (i IndexDB) SearchMessageByContentTypeAndKeywordController(contentType []int, keywordList []string, keywordListMatchType int, startTime, endTime int64, operationID string) (result []*model_struct.LocalChatLog, err error) {
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

func (i IndexDB) GetMinSeq(ID string) (uint32, error) {
	result, err := Exec(ID)
	if err != nil {
		return 0, err
	}
	if v, ok := result.(float64); ok {
		return uint32(v), nil
	} else {
		return 0, ErrType
	}
}

func (i IndexDB) SetMinSeq(ID string, minSeq uint32) error {
	_, err := Exec(ID, minSeq)
	return err
}

func (i IndexDB) GetUserMinSeq() (uint32, error) {
	result, err := Exec()
	if err != nil {
		return 0, err
	}
	if v, ok := result.(float64); ok {
		return uint32(v), nil
	} else {
		return 0, ErrType
	}
}

func (i IndexDB) GetGroupMinSeq(groupID string) (uint32, error) {
	result, err := Exec(groupID)
	if err != nil {
		return 0, err
	}
	if v, ok := result.(float64); ok {
		return uint32(v), nil
	} else {
		return 0, ErrType
	}
}

func (i IndexDB) BatchInsertMessageListController(MessageList []*model_struct.LocalChatLog) error {
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

func (i IndexDB) InsertMessageController(message *model_struct.LocalChatLog) error {
	switch message.SessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupInsertMessage(message, message.RecvID)
	default:
		return i.InsertMessage(message)
	}
}

func (i IndexDB) BatchUpdateMessageList(MessageList []*model_struct.LocalChatLog) error {
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
		err := i.UpdateMessageController(v1)
		if err != nil {
			return utils.Wrap(err, "BatchUpdateMessageList failed")
		}

	}
	return nil
}

func (i IndexDB) BatchSpecialUpdateMessageList(MessageList []*model_struct.LocalChatLog) error {
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

func (i IndexDB) GetMessageController(msg *sdk_struct.MsgStruct) (*model_struct.LocalChatLog, error) {
	switch msg.SessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupGetMessage(msg)
	default:
		return i.GetMessage(msg.ClientMsgID)
	}
}
func (i IndexDB) UpdateColumnsMessageController(ClientMsgID string, groupID string, sessionType int32, args map[string]interface{}) error {
	switch sessionType {
	case constant.SuperGroupChatType:
		return utils.Wrap(i.SuperGroupUpdateColumnsMessage(ClientMsgID, groupID, args), "")
	default:
		return utils.Wrap(i.UpdateColumnsMessage(ClientMsgID, args), "")
	}
}
func (i IndexDB) UpdateMessageController(c *model_struct.LocalChatLog) error {
	switch c.SessionType {
	case constant.SuperGroupChatType:
		return utils.Wrap(i.SuperGroupUpdateMessage(c), "")
	default:
		return utils.Wrap(i.UpdateMessage(c), "")
	}
}

func (i IndexDB) UpdateMessageStatusBySourceIDController(sourceID string, status, sessionType int32) error {
	switch sessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupUpdateMessageStatusBySourceID(sourceID, status, sessionType)
	default:
		return i.UpdateMessageStatusBySourceID(sourceID, status, sessionType)
	}
}
func (i IndexDB) UpdateMessageTimeAndStatusController(msg *sdk_struct.MsgStruct) error {
	switch msg.SessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupUpdateMessageTimeAndStatus(msg)
	default:
		return i.UpdateMessageTimeAndStatus(msg.ClientMsgID, msg.ServerMsgID, msg.SendTime, msg.Status)
	}
}
func (i IndexDB) GetMessageListController(sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	switch sessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupGetMessageList(sourceID, sessionType, count, startTime, isReverse)
	default:
		return i.GetMessageList(sourceID, sessionType, count, startTime, isReverse)
	}
}
func (i IndexDB) GetMessageListNoTimeController(sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	switch sessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupGetMessageListNoTime(sourceID, sessionType, count, isReverse)
	default:
		return i.GetMessageListNoTime(sourceID, sessionType, count, isReverse)
	}
}

func (i IndexDB) UpdateGroupMessageHasReadController(msgIDList []string, groupID string, sessionType int32) error {
	switch sessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupUpdateGroupMessageHasRead(msgIDList, groupID)
	default:
		return i.UpdateGroupMessageHasRead(msgIDList, sessionType)
	}
}

func (i IndexDB) GetMultipleMessageController(msgIDList []string, groupID string, sessionType int32) (result []*model_struct.LocalChatLog, err error) {
	switch sessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupGetMultipleMessage(msgIDList, groupID)
	default:
		return i.GetMultipleMessage(msgIDList)
	}
}

func (i IndexDB) GetMsgSeqByClientMsgIDController(m *sdk_struct.MsgStruct) (uint32, error) {
	switch m.SessionType {
	case constant.SuperGroupChatType:
		return i.SuperGroupGetMsgSeqByClientMsgID(m.ClientMsgID, m.GroupID)
	default:
		return i.GetMsgSeqByClientMsgID(m.ClientMsgID)
	}
}

func (i IndexDB) GetSubDepartmentList(departmentID string, args ...int) ([]*model_struct.LocalDepartment, error) {
	panic("implement me")
}

func (i IndexDB) InsertDepartment(department *model_struct.LocalDepartment) error {
	panic("implement me")
}

func (i IndexDB) UpdateDepartment(department *model_struct.LocalDepartment) error {
	panic("implement me")
}

func (i IndexDB) DeleteDepartment(departmentID string) error {
	panic("implement me")
}

func (i IndexDB) GetDepartmentInfo(departmentID string) (*model_struct.LocalDepartment, error) {
	panic("implement me")
}

func (i IndexDB) GetAllDepartmentList() ([]*model_struct.LocalDepartment, error) {
	panic("implement me")
}

func (i IndexDB) GetParentDepartmentList(departmentID string) ([]*model_struct.LocalDepartment, error) {
	panic("implement me")
}

func (i IndexDB) GetDepartmentList(departmentList *[]*model_struct.LocalDepartment, departmentID string) error {
	panic("implement me")
}

func (i IndexDB) GetParentDepartment(departmentID string) (model_struct.LocalDepartment, error) {
	panic("implement me")
}

func (i IndexDB) SearchDepartmentMember(keyWord string, isSearchUserName, isSearchEmail, isSearchMobile, isSearchPosition, isSearchTelephone, isSearchUserEnglishName, isSearchUserID bool, offset, count int) ([]*model_struct.SearchDepartmentMemberResult, error) {
	panic("implement me")
}

func (i IndexDB) SearchDepartment(keyWord string, offset, count int) ([]*model_struct.LocalDepartment, error) {
	panic("implement me")
}

func (i IndexDB) GetJoinedSuperGroupIDList() ([]string, error) {
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

func (i IndexDB) GetReadDiffusionGroupIDList() ([]string, error) {
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

func (i IndexDB) BatchInsertExceptionMsgController(MessageList []*model_struct.LocalErrChatLog) error {
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

func (i IndexDB) GetDepartmentMemberListByDepartmentID(departmentID string, args ...int) ([]*model_struct.LocalDepartmentMember, error) {
	panic("implement me")
}

func (i IndexDB) GetAllDepartmentMemberList() ([]*model_struct.LocalDepartmentMember, error) {
	panic("implement me")
}

func (i IndexDB) InsertDepartmentMember(departmentMember *model_struct.LocalDepartmentMember) error {
	panic("implement me")
}

func (i IndexDB) BatchInsertDepartmentMember(departmentMemberList []*model_struct.LocalDepartmentMember) error {
	panic("implement me")
}

func (i IndexDB) UpdateDepartmentMember(departmentMember *model_struct.LocalDepartmentMember) error {
	panic("implement me")
}

func (i IndexDB) DeleteDepartmentMember(departmentID string, userID string) error {
	panic("implement me")
}

func (i IndexDB) GetDepartmentMemberListByUserID(userID string) ([]*model_struct.LocalDepartmentMember, error) {
	panic("implement me")
}

func (i IndexDB) InsertWorkMomentsNotification(jsonDetail string) error {
	panic("implement me")
}

func (i IndexDB) GetWorkMomentsNotification(offset, count int) (WorkMomentsNotifications []*db.LocalWorkMomentsNotification, err error) {
	panic("implement me")
}

func (i IndexDB) GetWorkMomentsNotificationLimit(pageNumber, showNumber int) (WorkMomentsNotifications []*db.LocalWorkMomentsNotification, err error) {
	panic("implement me")
}

func (i IndexDB) InitWorkMomentsNotificationUnreadCount() error {
	panic("implement me")
}

func (i IndexDB) IncrWorkMomentsNotificationUnreadCount() error {
	panic("implement me")
}

func (i IndexDB) MarkAllWorkMomentsNotificationAsRead() (err error) {
	panic("implement me")
}

func (i IndexDB) GetWorkMomentsUnReadCount() (workMomentsNotificationUnReadCount db.LocalWorkMomentsNotificationUnreadCount, err error) {
	panic("implement me")
}

func (i IndexDB) ClearWorkMomentsNotification() (err error) {
	panic("implement me")
}

func (i IndexDB) Close() error {
	_, err := Exec()
	return err
}

func (i IndexDB) InitDB(userID string, dataDir string) error {
	_, err := Exec(userID, dataDir)
	return err
}

//func (i IndexDB) GetBlackList() ([]*model_struct.LocalBlack, error) {
//	panic("implement me")
//}
//
//func (i IndexDB) GetBlackListUserID() (blackListUid []string, err error) {
//	panic("implement me")
//}
//
//func (i IndexDB) GetBlackInfoByBlockUserID(blockUserID string) (*model_struct.LocalBlack, error) {
//	panic("implement me")
//}
//
//func (i IndexDB) GetBlackInfoList(blockUserIDList []string) ([]*model_struct.LocalBlack, error) {
//	panic("implement me")
//}
//
//func (i IndexDB) InsertBlack(black *model_struct.LocalBlack) error {
//	panic("implement me")
//}
//
//func (i IndexDB) UpdateBlack(black *model_struct.LocalBlack) error {
//	panic("implement me")
//}
//
//func (i IndexDB) DeleteBlack(blockUserID string) error {
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

func (i IndexDB) SetChatLogFailedStatus() {
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
