package indexdb

import (
	"errors"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"runtime"
	"sync"
	"syscall/js"
)

type IndexDB struct {
	LocalUsers
}

type CallbackData struct {
	ErrCode int32  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	Data    string `json:"data"`
}

func Exec(args ...interface{}) (output string, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
				log.Error("", "panic info is here", err.Error())
			case error:
				err = x
				log.Error("", "panic info is here", err.Error())
			default:
				err = errors.New("unknow panic")
			}
		}
	}()
	pc, _, _, _ := runtime.Caller(1)
	funcName := utils.CleanUpfuncName(runtime.FuncForPC(pc).Name())
	data := CallbackData{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	js.Global().Call(funcName, args...).Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		log.Debug("js", "=> (main go context) "+funcName+" with respone ", args[0].String())
		interErr := utils.JsonStringToStruct(args[0].String(), &data)
		if interErr != nil {
			err = utils.Wrap(err, "return json unmarshal err from javascript")
			wg.Done()
			return nil
		}
		wg.Done()
		return nil
	}))
	wg.Wait()
	if data.ErrCode != 0 {
		return "", errors.New(data.ErrMsg)
	}
	return data.Data, err

}
func (i IndexDB) InsertGroup(groupInfo *model_struct.LocalGroup) error {
	panic("implement me")
}

func (i IndexDB) DeleteGroup(groupID string) error {
	panic("implement me")
}

func (i IndexDB) UpdateGroup(groupInfo *model_struct.LocalGroup) error {
	panic("implement me")
}

func (i IndexDB) GetJoinedGroupList() ([]*model_struct.LocalGroup, error) {
	panic("implement me")
}

func (i IndexDB) GetGroupInfoByGroupID(groupID string) (*model_struct.LocalGroup, error) {
	panic("implement me")
}

func (i IndexDB) GetAllGroupInfoByGroupIDOrGroupName(keyword string, isSearchGroupID bool, isSearchGroupName bool) ([]*model_struct.LocalGroup, error) {
	panic("implement me")
}

func (i IndexDB) AddMemberCount(groupID string) error {
	panic("implement me")
}

func (i IndexDB) SubtractMemberCount(groupID string) error {
	panic("implement me")
}

func (i IndexDB) GetJoinedWorkingGroupIDList() ([]string, error) {
	panic("implement me")
}

func (i IndexDB) GetJoinedWorkingGroupList() ([]*model_struct.LocalGroup, error) {
	panic("implement me")
}

func (i IndexDB) GetMinSeq(ID string) (uint32, error) {
	panic("implement me")
}

func (i IndexDB) SetMinSeq(ID string, minSeq uint32) error {
	panic("implement me")
}

func (i IndexDB) GetUserMinSeq() (uint32, error) {
	panic("implement me")
}

func (i IndexDB) GetGroupMinSeq(groupID string) (uint32, error) {
	panic("implement me")
}

func (i IndexDB) InsertAdminGroupRequest(groupRequest *model_struct.LocalAdminGroupRequest) error {
	panic("implement me")
}

func (i IndexDB) DeleteAdminGroupRequest(groupID, userID string) error {
	panic("implement me")
}

func (i IndexDB) UpdateAdminGroupRequest(groupRequest *model_struct.LocalAdminGroupRequest) error {
	panic("implement me")
}

func (i IndexDB) GetAdminGroupApplication() ([]*model_struct.LocalAdminGroupRequest, error) {
	panic("implement me")
}

func (i IndexDB) BatchInsertMessageList(MessageList []*model_struct.LocalChatLog) error {
	panic("implement me")
}

func (i IndexDB) BatchInsertMessageListController(MessageList []*model_struct.LocalChatLog) error {
	panic("implement me")
}

func (i IndexDB) InsertMessage(Message *model_struct.LocalChatLog) error {
	panic("implement me")
}

func (i IndexDB) InsertMessageController(message *model_struct.LocalChatLog) error {
	panic("implement me")
}

func (i IndexDB) SearchMessageByKeyword(contentType []int, keywordList []string, keywordListMatchType int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) SearchMessageByKeywordController(contentType []int, keywordList []string, keywordListMatchType int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) SearchMessageByContentType(contentType []int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) SearchMessageByContentTypeController(contentType []int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) SearchMessageByContentTypeAndKeyword(contentType []int, keywordList []string, keywordListMatchType int, startTime, endTime int64) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) SearchMessageByContentTypeAndKeywordController(contentType []int, keywordList []string, keywordListMatchType int, startTime, endTime int64, operationID string) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) BatchUpdateMessageList(MessageList []*model_struct.LocalChatLog) error {
	panic("implement me")
}

func (i IndexDB) BatchSpecialUpdateMessageList(MessageList []*model_struct.LocalChatLog) error {
	panic("implement me")
}

func (i IndexDB) MessageIfExists(ClientMsgID string) (bool, error) {
	panic("implement me")
}

func (i IndexDB) IsExistsInErrChatLogBySeq(seq int64) bool {
	panic("implement me")
}

func (i IndexDB) MessageIfExistsBySeq(seq int64) (bool, error) {
	panic("implement me")
}

//1,使用wasm原生的方式，tinygo应用于go的嵌入式领域，支持的功能有限，甚至json序列化都无法支持
//2.函数命名遵从驼峰还是帕斯卡命名法需要确定一下
//3.提供的sql生成语句中，关于bool值需要特殊处理，create语句的设计的到bool值的我会在创建语句中单独说明，这是因为在原有的sqlite中并不支持bool，用整数1或者0替代，gorm对其做了转换。
//4.提供的sql生成语句中，字段名是下划线方式 例如：recv_id，但是接口转换的数据json tag字段的风格是recvID，类似的所有的字段需要做个map映射

func (i IndexDB) GetMessage(ClientMsgID string) (*model_struct.LocalChatLog, error) {
	msg, err := Exec(ClientMsgID)
	log.Info("test", msg)
	if err != nil {
		return nil, err
	} else {
		result := model_struct.LocalChatLog{}
		err := utils.JsonStringToStruct(msg, &result)
		if err != nil {
			return nil, err
		}
		return &result, err

	}
}

func (i IndexDB) GetMessageController(msg *sdk_struct.MsgStruct) (*model_struct.LocalChatLog, error) {
	panic("implement me")
}

func (i IndexDB) GetAllUnDeleteMessageSeqList() ([]uint32, error) {
	panic("implement me")
}

func (i IndexDB) UpdateColumnsMessageList(clientMsgIDList []string, args map[string]interface{}) error {
	panic("implement me")
}

func (i IndexDB) UpdateColumnsMessage(ClientMsgID string, args map[string]interface{}) error {
	panic("implement me")
}

func (i IndexDB) UpdateColumnsMessageController(ClientMsgID string, groupID string, sessionType int32, args map[string]interface{}) error {
	panic("implement me")
}

func (i IndexDB) UpdateMessage(c *model_struct.LocalChatLog) error {
	panic("implement me")
}

func (i IndexDB) UpdateMessageController(c *model_struct.LocalChatLog) error {
	panic("implement me")
}

func (i IndexDB) DeleteAllMessage() error {
	panic("implement me")
}

func (i IndexDB) UpdateMessageStatusBySourceID(sourceID string, status, sessionType int32) error {
	panic("implement me")
}

func (i IndexDB) UpdateMessageStatusBySourceIDController(sourceID string, status, sessionType int32) error {
	panic("implement me")
}

func (i IndexDB) UpdateMessageTimeAndStatus(clientMsgID string, serverMsgID string, sendTime int64, status int32) error {
	panic("implement me")
}

func (i IndexDB) UpdateMessageTimeAndStatusController(msg *sdk_struct.MsgStruct) error {
	panic("implement me")
}

func (i IndexDB) GetMessageList(sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) GetMessageListController(sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) GetMessageListNoTime(sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) GetMessageListNoTimeController(sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) GetSendingMessageList() (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) UpdateSingleMessageHasRead(sendID string, msgIDList []string) error {
	panic("implement me")
}

func (i IndexDB) UpdateGroupMessageHasRead(msgIDList []string, sessionType int32) error {
	panic("implement me")
}

func (i IndexDB) UpdateGroupMessageHasReadController(msgIDList []string, groupID string, sessionType int32) error {
	panic("implement me")
}

func (i IndexDB) GetMultipleMessage(msgIDList []string) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) GetMultipleMessageController(msgIDList []string, groupID string, sessionType int32) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) GetNormalMsgSeq() (uint32, error) {
	panic("implement me")
}

func (i IndexDB) GetLostMsgSeqList(minSeqInSvr uint32) ([]uint32, error) {
	panic("implement me")
}

func (i IndexDB) GetSuperGroupNormalMsgSeq(groupID string) (uint32, error) {
	panic("implement me")
}

func (i IndexDB) GetTestMessage(seq uint32) (*model_struct.LocalChatLog, error) {
	panic("implement me")
}

func (i IndexDB) UpdateMsgSenderNickname(sendID, nickname string, sType int) error {
	panic("implement me")
}

func (i IndexDB) UpdateMsgSenderFaceURL(sendID, faceURL string, sType int) error {
	panic("implement me")
}

func (i IndexDB) UpdateMsgSenderFaceURLAndSenderNickname(sendID, faceURL, nickname string, sessionType int) error {
	panic("implement me")
}

func (i IndexDB) GetMsgSeqByClientMsgID(clientMsgID string) (uint32, error) {
	panic("implement me")
}

func (i IndexDB) GetMsgSeqByClientMsgIDController(m *sdk_struct.MsgStruct) (uint32, error) {
	panic("implement me")
}

func (i IndexDB) GetMsgSeqListByGroupID(groupID string) ([]uint32, error) {
	panic("implement me")
}

func (i IndexDB) GetMsgSeqListByPeerUserID(userID string) ([]uint32, error) {
	panic("implement me")
}

func (i IndexDB) GetMsgSeqListBySelfUserID(userID string) ([]uint32, error) {
	panic("implement me")
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

func (i IndexDB) InsertGroupRequest(groupRequest *model_struct.LocalGroupRequest) error {
	panic("implement me")
}

func (i IndexDB) DeleteGroupRequest(groupID, userID string) error {
	panic("implement me")
}

func (i IndexDB) UpdateGroupRequest(groupRequest *model_struct.LocalGroupRequest) error {
	panic("implement me")
}

func (i IndexDB) GetSendGroupApplication() ([]*model_struct.LocalGroupRequest, error) {
	panic("implement me")
}

func (i IndexDB) GetConversationByUserID(userID string) (*model_struct.LocalConversation, error) {
	panic("implement me")
}

func (i IndexDB) GetAllConversationList() ([]*model_struct.LocalConversation, error) {
	panic("implement me")
}

func (i IndexDB) GetHiddenConversationList() ([]*model_struct.LocalConversation, error) {
	panic("implement me")
}

func (i IndexDB) GetAllConversationListToSync() ([]*model_struct.LocalConversation, error) {
	panic("implement me")
}

func (i IndexDB) GetConversationListSplit(offset, count int) ([]*model_struct.LocalConversation, error) {
	panic("implement me")
}

func (i IndexDB) BatchInsertConversationList(conversationList []*model_struct.LocalConversation) error {
	panic("implement me")
}

func (i IndexDB) InsertConversation(conversationList *model_struct.LocalConversation) error {
	panic("implement me")
}

func (i IndexDB) DeleteConversation(conversationID string) error {
	panic("implement me")
}

func (i IndexDB) GetConversation(conversationID string) (*model_struct.LocalConversation, error) {
	panic("implement me")
}

func (i IndexDB) UpdateConversation(c *model_struct.LocalConversation) error {
	panic("implement me")
}

func (i IndexDB) UpdateConversationForSync(c *model_struct.LocalConversation) error {
	panic("implement me")
}

func (i IndexDB) BatchUpdateConversationList(conversationList []*model_struct.LocalConversation) error {
	panic("implement me")
}

func (i IndexDB) ConversationIfExists(conversationID string) (bool, error) {
	panic("implement me")
}

func (i IndexDB) ResetConversation(conversationID string) error {
	panic("implement me")
}

func (i IndexDB) ResetAllConversation() error {
	panic("implement me")
}

func (i IndexDB) ClearConversation(conversationID string) error {
	panic("implement me")
}

func (i IndexDB) CleaAllConversation() error {
	panic("implement me")
}

func (i IndexDB) SetConversationDraft(conversationID, draftText string) error {
	panic("implement me")
}

func (i IndexDB) RemoveConversationDraft(conversationID, draftText string) error {
	panic("implement me")
}

func (i IndexDB) UnPinConversation(conversationID string, isPinned int) error {
	panic("implement me")
}

func (i IndexDB) UpdateColumnsConversation(conversationID string, args map[string]interface{}) error {
	panic("implement me")
}

func (i IndexDB) UpdateAllConversation(conversation *model_struct.LocalConversation) error {
	panic("implement me")
}

func (i IndexDB) IncrConversationUnreadCount(conversationID string) error {
	panic("implement me")
}

func (i IndexDB) GetTotalUnreadMsgCount() (totalUnreadCount int32, err error) {
	panic("implement me")
}

func (i IndexDB) SetMultipleConversationRecvMsgOpt(conversationIDList []string, opt int) (err error) {
	panic("implement me")
}

func (i IndexDB) GetMultipleConversation(conversationIDList []string) (result []*model_struct.LocalConversation, err error) {
	panic("implement me")
}

func (i IndexDB) GetJoinedSuperGroupList() ([]*model_struct.LocalGroup, error) {
	panic("implement me")
}

func (i IndexDB) GetJoinedSuperGroupIDList() ([]string, error) {
	panic("implement me")
}

func (i IndexDB) InsertSuperGroup(groupInfo *model_struct.LocalGroup) error {
	panic("implement me")
}

func (i IndexDB) DeleteAllSuperGroup() error {
	panic("implement me")
}

func (i IndexDB) GetSuperGroupInfoByGroupID(groupID string) (*model_struct.LocalGroup, error) {
	panic("implement me")
}

func (i IndexDB) UpdateSuperGroup(groupInfo *model_struct.LocalGroup) error {
	panic("implement me")
}

func (i IndexDB) DeleteSuperGroup(groupID string) error {
	panic("implement me")
}

func (i IndexDB) GetReadDiffusionGroupIDList() ([]string, error) {
	panic("implement me")
}

func (i IndexDB) GetGroupMemberInfoByGroupIDUserID(groupID, userID string) (*model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (i IndexDB) GetAllGroupMemberList() ([]model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (i IndexDB) GetAllGroupMemberUserIDList() ([]model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (i IndexDB) GetGroupMemberCount(groupID string) (uint32, error) {
	panic("implement me")
}

func (i IndexDB) GetGroupSomeMemberInfo(groupID string, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (i IndexDB) GetGroupAdminID(groupID string) ([]string, error) {
	panic("implement me")
}

func (i IndexDB) GetGroupMemberListByGroupID(groupID string) ([]*model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (i IndexDB) GetGroupMemberListSplit(groupID string, filter int32, offset, count int) ([]*model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (i IndexDB) GetGroupMemberOwnerAndAdmin(groupID string) ([]*model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (i IndexDB) GetGroupMemberOwner(groupID string) (*model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (i IndexDB) GetGroupMemberListSplitByJoinTimeFilter(groupID string, offset, count int, joinTimeBegin, joinTimeEnd int64, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (i IndexDB) GetGroupOwnerAndAdminByGroupID(groupID string) ([]*model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (i IndexDB) GetGroupMemberUIDListByGroupID(groupID string) (result []string, err error) {
	panic("implement me")
}

func (i IndexDB) InsertGroupMember(groupMember *model_struct.LocalGroupMember) error {
	panic("implement me")
}

func (i IndexDB) BatchInsertGroupMember(groupMemberList []*model_struct.LocalGroupMember) error {
	panic("implement me")
}

func (i IndexDB) DeleteGroupMember(groupID, userID string) error {
	panic("implement me")
}

func (i IndexDB) DeleteGroupAllMembers(groupID string) error {
	panic("implement me")
}

func (i IndexDB) UpdateGroupMember(groupMember *model_struct.LocalGroupMember) error {
	panic("implement me")
}

func (i IndexDB) UpdateGroupMemberField(groupID, userID string, args map[string]interface{}) error {
	panic("implement me")
}

func (i IndexDB) GetGroupMemberInfoIfOwnerOrAdmin() ([]*model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (i IndexDB) SearchGroupMembers(keyword string, groupID string, isSearchMemberNickname, isSearchUserID bool, offset, count int) (result []*model_struct.LocalGroupMember, err error) {
	panic("implement me")
}

func (i IndexDB) InitSuperLocalErrChatLog(groupID string) {
	panic("implement me")
}

func (i IndexDB) SuperBatchInsertExceptionMsg(MessageList []*model_struct.LocalErrChatLog, groupID string) error {
	panic("implement me")
}

func (i IndexDB) GetAbnormalMsgSeq() (uint32, error) {
	panic("implement me")
}

func (i IndexDB) GetAbnormalMsgSeqList() ([]uint32, error) {
	panic("implement me")
}

func (i IndexDB) BatchInsertExceptionMsg(MessageList []*model_struct.LocalErrChatLog) error {
	panic("implement me")
}

func (i IndexDB) BatchInsertExceptionMsgController(MessageList []*model_struct.LocalErrChatLog) error {
	panic("implement me")
}

func (i IndexDB) GetSuperGroupAbnormalMsgSeq(groupID string) (uint32, error) {
	panic("implement me")
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

func (i IndexDB) BatchInsertTempCacheMessageList(MessageList []*model_struct.TempCacheLocalChatLog) error {
	panic("implement me")
}

func (i IndexDB) InsertTempCacheMessage(Message *model_struct.TempCacheLocalChatLog) error {
	panic("implement me")
}

func (i IndexDB) InsertFriend(friend *model_struct.LocalFriend) error {
	panic("implement me")
}

func (i IndexDB) DeleteFriend(friendUserID string) error {
	panic("implement me")
}

func (i IndexDB) UpdateFriend(friend *model_struct.LocalFriend) error {
	panic("implement me")
}

func (i IndexDB) GetAllFriendList() ([]*model_struct.LocalFriend, error) {
	panic("implement me")
}

func (i IndexDB) SearchFriendList(keyword string, isSearchUserID, isSearchNickname, isSearchRemark bool) ([]*model_struct.LocalFriend, error) {
	panic("implement me")
}

func (i IndexDB) GetFriendInfoByFriendUserID(FriendUserID string) (*model_struct.LocalFriend, error) {
	panic("implement me")
}

func (i IndexDB) GetFriendInfoList(friendUserIDList []string) ([]*model_struct.LocalFriend, error) {
	panic("implement me")
}

func (i IndexDB) InsertFriendRequest(friendRequest *model_struct.LocalFriendRequest) error {
	panic("implement me")
}

func (i IndexDB) DeleteFriendRequestBothUserID(fromUserID, toUserID string) error {
	panic("implement me")
}

func (i IndexDB) UpdateFriendRequest(friendRequest *model_struct.LocalFriendRequest) error {
	panic("implement me")
}

func (i IndexDB) GetRecvFriendApplication() ([]*model_struct.LocalFriendRequest, error) {
	panic("implement me")
}

func (i IndexDB) GetSendFriendApplication() ([]*model_struct.LocalFriendRequest, error) {
	panic("implement me")
}

func (i IndexDB) GetFriendApplicationByBothID(fromUserID, toUserID string) (*model_struct.LocalFriendRequest, error) {
	panic("implement me")
}

func (i IndexDB) InitSuperLocalChatLog(groupID string) {
	panic("implement me")
}

func (i IndexDB) SuperGroupBatchInsertMessageList(MessageList []*model_struct.LocalChatLog, groupID string) error {
	panic("implement me")
}

func (i IndexDB) SuperGroupInsertMessage(Message *model_struct.LocalChatLog, groupID string) error {
	panic("implement me")
}

func (i IndexDB) SuperGroupDeleteAllMessage(groupID string) error {
	panic("implement me")
}

func (i IndexDB) SuperGroupSearchMessageByKeyword(contentType []int, keywordList []string, keywordListMatchType int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) SuperGroupSearchMessageByContentType(contentType []int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
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

func (i IndexDB) SuperGroupGetMessage(msg *sdk_struct.MsgStruct) (*model_struct.LocalChatLog, error) {
	panic("implement me")
}

func (i IndexDB) SuperGroupGetAllUnDeleteMessageSeqList() ([]uint32, error) {
	panic("implement me")
}

func (i IndexDB) SuperGroupUpdateColumnsMessage(ClientMsgID, groupID string, args map[string]interface{}) error {
	panic("implement me")
}

func (i IndexDB) SuperGroupUpdateMessage(c *model_struct.LocalChatLog) error {
	panic("implement me")
}

func (i IndexDB) SuperGroupUpdateMessageStatusBySourceID(sourceID string, status, sessionType int32) error {
	panic("implement me")
}

func (i IndexDB) SuperGroupUpdateMessageTimeAndStatus(msg *sdk_struct.MsgStruct) error {
	panic("implement me")
}

func (i IndexDB) SuperGroupGetMessageList(sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) SuperGroupGetMessageListNoTime(sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) SuperGroupGetSendingMessageList() (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) SuperGroupUpdateGroupMessageHasRead(msgIDList []string, groupID string) error {
	panic("implement me")
}

func (i IndexDB) SuperGroupGetMultipleMessage(conversationIDList []string, groupID string) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) SuperGroupGetNormalMsgSeq() (uint32, error) {
	panic("implement me")
}

func (i IndexDB) SuperGroupGetNormalMinSeq(groupID string) (uint32, error) {
	panic("implement me")
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

func (i IndexDB) CloseDB() error {
	panic("implement me")
}

func (i IndexDB) SetChatLogFailedStatus() {
	panic("implement me")
}

func (i IndexDB) InitDB() error {
	panic("implement me")
}

func (i IndexDB) GetBlackList() ([]*model_struct.LocalBlack, error) {
	panic("implement me")
}

func (i IndexDB) GetBlackListUserID() (blackListUid []string, err error) {
	panic("implement me")
}

func (i IndexDB) GetBlackInfoByBlockUserID(blockUserID string) (*model_struct.LocalBlack, error) {
	panic("implement me")
}

func (i IndexDB) GetBlackInfoList(blockUserIDList []string) ([]*model_struct.LocalBlack, error) {
	panic("implement me")
}

func (i IndexDB) InsertBlack(black *model_struct.LocalBlack) error {
	panic("implement me")
}

func (i IndexDB) UpdateBlack(black *model_struct.LocalBlack) error {
	panic("implement me")
}

func (i IndexDB) DeleteBlack(blockUserID string) error {
	panic("implement me")
}

func NewIndexDB() *IndexDB {
	return &IndexDB{}
}
