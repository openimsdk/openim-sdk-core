package db

import "context"

import (
	"errors"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"open_im_sdk/wasm/indexdb"
)

var ErrType = errors.New("from javascript data type err")

type IndexDB struct {
	*indexdb.LocalUsers
	*indexdb.LocalConversations
	*indexdb.LocalChatLogs
	*indexdb.LocalSuperGroupChatLogs
	*indexdb.LocalSuperGroup
	*indexdb.LocalConversationUnreadMessages
	*indexdb.LocalGroups
	*indexdb.LocalGroupMember
	*indexdb.LocalCacheMessage
	*indexdb.FriendRequest
	*indexdb.Black
	*indexdb.Friend
	*indexdb.LocalGroupRequest
	*indexdb.LocalChatLogReactionExtensions
	loginUserID string
}

func (i IndexDB) GetJoinedSuperGroupIDList(ctx context.Context) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) GetReadDiffusionGroupIDList(ctx context.Context) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) BatchInsertExceptionMsgController(ctx context.Context, MessageList []*model_struct.LocalErrChatLog) error {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) GetDepartmentMemberListByDepartmentID(ctx context.Context, departmentID string, args ...int) ([]*model_struct.LocalDepartmentMember, error) {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) GetAllDepartmentMemberList(ctx context.Context) ([]*model_struct.LocalDepartmentMember, error) {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) InsertDepartmentMember(ctx context.Context, departmentMember *model_struct.LocalDepartmentMember) error {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) BatchInsertDepartmentMember(ctx context.Context, departmentMemberList []*model_struct.LocalDepartmentMember) error {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) UpdateDepartmentMember(ctx context.Context, departmentMember *model_struct.LocalDepartmentMember) error {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) DeleteDepartmentMember(ctx context.Context, departmentID string, userID string) error {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) GetDepartmentMemberListByUserID(ctx context.Context, userID string) ([]*model_struct.LocalDepartmentMember, error) {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) SuperGroupUpdateSpecificContentTypeMessage(ctx context.Context, contentType int, groupID string, args map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) SuperGroupUpdateGroupMessageFields(ctx context.Context, msgIDList []string, groupID string, args map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) SuperGroupGetAlreadyExistSeqList(ctx context.Context, groupID string, lostSeqList []uint32) (seqList []uint32, err error) {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) InsertWorkMomentsNotification(ctx context.Context, jsonDetail string) error {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) GetWorkMomentsNotification(ctx context.Context, offset, count int) (WorkMomentsNotifications []*model_struct.LocalWorkMomentsNotification, err error) {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) GetWorkMomentsNotificationLimit(ctx context.Context, pageNumber, showNumber int) (WorkMomentsNotifications []*model_struct.LocalWorkMomentsNotification, err error) {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) InitWorkMomentsNotificationUnreadCount(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) IncrWorkMomentsNotificationUnreadCount(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) MarkAllWorkMomentsNotificationAsRead(ctx context.Context) (err error) {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) GetWorkMomentsUnReadCount(ctx context.Context) (workMomentsNotificationUnReadCount model_struct.LocalWorkMomentsNotificationUnreadCount, err error) {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) ClearWorkMomentsNotification(ctx context.Context) (err error) {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) Close(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
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

func (i IndexDB) InitDB(ctx context.Context, userID string, dataDir string) error {
	_, err := indexdb.Exec(userID, dataDir)
	return err
}

func NewDataBase(loginUserID string, dbDir string, operationID string) (*IndexDB, error) {
	return &IndexDB{
		LocalUsers:                      indexdb.NewLocalUsers(),
		LocalConversations:              indexdb.NewLocalConversations(),
		LocalChatLogs:                   indexdb.NewLocalChatLogs(loginUserID),
		LocalSuperGroupChatLogs:         indexdb.NewLocalSuperGroupChatLogs(),
		LocalSuperGroup:                 indexdb.NewLocalSuperGroup(),
		LocalConversationUnreadMessages: indexdb.NewLocalConversationUnreadMessages(),
		LocalGroups:                     indexdb.NewLocalGroups(),
		LocalGroupMember:                indexdb.NewLocalGroupMember(),
		LocalCacheMessage:               indexdb.NewLocalCacheMessage(),
		FriendRequest:                   indexdb.NewFriendRequest(loginUserID),
		Black:                           indexdb.NewBlack(loginUserID),
		Friend:                          indexdb.NewFriend(loginUserID),
		LocalGroupRequest:               indexdb.NewLocalGroupRequest(),
		LocalChatLogReactionExtensions:  indexdb.NewLocalChatLogReactionExtensions(),
		loginUserID:                     loginUserID,
	}, nil
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
