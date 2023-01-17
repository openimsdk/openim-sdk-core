//+build !js

package indexdb

import (
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/sdk_struct"
)

type NilIndexDB struct {
}

func (n NilIndexDB) InsertGroup(groupInfo *model_struct.LocalGroup) error {
	panic("implement me")
}

func (n NilIndexDB) DeleteGroup(groupID string) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateGroup(groupInfo *model_struct.LocalGroup) error {
	panic("implement me")
}

func (n NilIndexDB) GetJoinedGroupListDB() ([]*model_struct.LocalGroup, error) {
	panic("implement me")
}

func (n NilIndexDB) GetGroupInfoByGroupID(groupID string) (*model_struct.LocalGroup, error) {
	panic("implement me")
}

func (n NilIndexDB) GetAllGroupInfoByGroupIDOrGroupName(keyword string, isSearchGroupID bool, isSearchGroupName bool) ([]*model_struct.LocalGroup, error) {
	panic("implement me")
}

func (n NilIndexDB) AddMemberCount(groupID string) error {
	panic("implement me")
}

func (n NilIndexDB) SubtractMemberCount(groupID string) error {
	panic("implement me")
}

func (n NilIndexDB) GetJoinedWorkingGroupIDList() ([]string, error) {
	panic("implement me")
}

func (n NilIndexDB) GetJoinedWorkingGroupList() ([]*model_struct.LocalGroup, error) {
	panic("implement me")
}

func (n NilIndexDB) GetMinSeq(ID string) (uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) SetMinSeq(ID string, minSeq uint32) error {
	panic("implement me")
}

func (n NilIndexDB) GetUserMinSeq() (uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) GetGroupMinSeq(groupID string) (uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) InsertAdminGroupRequest(groupRequest *model_struct.LocalAdminGroupRequest) error {
	panic("implement me")
}

func (n NilIndexDB) DeleteAdminGroupRequest(groupID, userID string) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateAdminGroupRequest(groupRequest *model_struct.LocalAdminGroupRequest) error {
	panic("implement me")
}

func (n NilIndexDB) GetAdminGroupApplication() ([]*model_struct.LocalAdminGroupRequest, error) {
	panic("implement me")
}

func (n NilIndexDB) BatchInsertMessageList(MessageList []*model_struct.LocalChatLog) error {
	panic("implement me")
}

func (n NilIndexDB) BatchInsertMessageListController(MessageList []*model_struct.LocalChatLog) error {
	panic("implement me")
}

func (n NilIndexDB) InsertMessage(Message *model_struct.LocalChatLog) error {
	panic("implement me")
}

func (n NilIndexDB) InsertMessageController(message *model_struct.LocalChatLog) error {
	panic("implement me")
}

func (n NilIndexDB) SearchMessageByKeyword(contentType []int, keywordList []string, keywordListMatchType int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (n NilIndexDB) SearchMessageByKeywordController(contentType []int, keywordList []string, keywordListMatchType int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (n NilIndexDB) SearchMessageByContentType(contentType []int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (n NilIndexDB) SearchMessageByContentTypeController(contentType []int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (n NilIndexDB) SearchMessageByContentTypeAndKeyword(contentType []int, keywordList []string, keywordListMatchType int, startTime, endTime int64) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (n NilIndexDB) SearchMessageByContentTypeAndKeywordController(contentType []int, keywordList []string, keywordListMatchType int, startTime, endTime int64, operationID string) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (n NilIndexDB) BatchUpdateMessageList(MessageList []*model_struct.LocalChatLog) error {
	panic("implement me")
}

func (n NilIndexDB) BatchSpecialUpdateMessageList(MessageList []*model_struct.LocalChatLog) error {
	panic("implement me")
}

func (n NilIndexDB) MessageIfExists(ClientMsgID string) (bool, error) {
	panic("implement me")
}

func (n NilIndexDB) IsExistsInErrChatLogBySeq(seq int64) bool {
	panic("implement me")
}

func (n NilIndexDB) MessageIfExistsBySeq(seq int64) (bool, error) {
	panic("implement me")
}

func (n NilIndexDB) GetMessage(ClientMsgID string) (*model_struct.LocalChatLog, error) {
	panic("implement me")
}

func (n NilIndexDB) GetMessageController(msg *sdk_struct.MsgStruct) (*model_struct.LocalChatLog, error) {
	panic("implement me")
}

func (n NilIndexDB) GetAllUnDeleteMessageSeqList() ([]uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) UpdateColumnsMessageList(clientMsgIDList []string, args map[string]interface{}) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateColumnsMessage(ClientMsgID string, args map[string]interface{}) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateColumnsMessageController(ClientMsgID string, groupID string, sessionType int32, args map[string]interface{}) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateMessage(c *model_struct.LocalChatLog) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateMessageController(c *model_struct.LocalChatLog) error {
	panic("implement me")
}

func (n NilIndexDB) DeleteAllMessage() error {
	panic("implement me")
}

func (n NilIndexDB) UpdateMessageStatusBySourceID(sourceID string, status, sessionType int32) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateMessageStatusBySourceIDController(sourceID string, status, sessionType int32) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateMessageTimeAndStatus(clientMsgID string, serverMsgID string, sendTime int64, status int32) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateMessageTimeAndStatusController(msg *sdk_struct.MsgStruct) error {
	panic("implement me")
}

func (n NilIndexDB) GetMessageList(sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (n NilIndexDB) GetMessageListController(sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (n NilIndexDB) GetMessageListNoTime(sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (n NilIndexDB) GetMessageListNoTimeController(sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (n NilIndexDB) GetSendingMessageList() (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (n NilIndexDB) UpdateSingleMessageHasRead(sendID string, msgIDList []string) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateGroupMessageHasRead(msgIDList []string, sessionType int32) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateGroupMessageHasReadController(msgIDList []string, groupID string, sessionType int32) error {
	panic("implement me")
}

func (n NilIndexDB) GetMultipleMessage(msgIDList []string) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (n NilIndexDB) GetMultipleMessageController(msgIDList []string, groupID string, sessionType int32) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (n NilIndexDB) GetNormalMsgSeq() (uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) GetLostMsgSeqList(minSeqInSvr uint32) ([]uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) GetSuperGroupNormalMsgSeq(groupID string) (uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) GetTestMessage(seq uint32) (*model_struct.LocalChatLog, error) {
	panic("implement me")
}

func (n NilIndexDB) UpdateMsgSenderNickname(sendID, nickname string, sType int) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateMsgSenderFaceURL(sendID, faceURL string, sType int) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateMsgSenderFaceURLAndSenderNickname(sendID, faceURL, nickname string, sessionType int) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateMsgSenderFaceURLAndSenderNicknameController(sendID, faceURL, nickname string, sessionType int, groupID string) error {
	panic("implement me")
}

func (n NilIndexDB) GetMsgSeqByClientMsgID(clientMsgID string) (uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) GetMsgSeqByClientMsgIDController(m *sdk_struct.MsgStruct) (uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) GetMsgSeqListByGroupID(groupID string) ([]uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) GetMsgSeqListByPeerUserID(userID string) ([]uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) GetMsgSeqListBySelfUserID(userID string) ([]uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) GetSubDepartmentList(departmentID string, args ...int) ([]*model_struct.LocalDepartment, error) {
	panic("implement me")
}

func (n NilIndexDB) InsertDepartment(department *model_struct.LocalDepartment) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateDepartment(department *model_struct.LocalDepartment) error {
	panic("implement me")
}

func (n NilIndexDB) DeleteDepartment(departmentID string) error {
	panic("implement me")
}

func (n NilIndexDB) GetDepartmentInfo(departmentID string) (*model_struct.LocalDepartment, error) {
	panic("implement me")
}

func (n NilIndexDB) GetAllDepartmentList() ([]*model_struct.LocalDepartment, error) {
	panic("implement me")
}

func (n NilIndexDB) GetParentDepartmentList(departmentID string) ([]*model_struct.LocalDepartment, error) {
	panic("implement me")
}

func (n NilIndexDB) GetDepartmentList(departmentList *[]*model_struct.LocalDepartment, departmentID string) error {
	panic("implement me")
}

func (n NilIndexDB) GetParentDepartment(departmentID string) (model_struct.LocalDepartment, error) {
	panic("implement me")
}

func (n NilIndexDB) SearchDepartmentMember(keyWord string, isSearchUserName, isSearchEmail, isSearchMobile, isSearchPosition, isSearchTelephone, isSearchUserEnglishName, isSearchUserID bool, offset, count int) ([]*model_struct.SearchDepartmentMemberResult, error) {
	panic("implement me")
}

func (n NilIndexDB) SearchDepartment(keyWord string, offset, count int) ([]*model_struct.LocalDepartment, error) {
	panic("implement me")
}

func (n NilIndexDB) InsertGroupRequest(groupRequest *model_struct.LocalGroupRequest) error {
	panic("implement me")
}

func (n NilIndexDB) DeleteGroupRequest(groupID, userID string) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateGroupRequest(groupRequest *model_struct.LocalGroupRequest) error {
	panic("implement me")
}

func (n NilIndexDB) GetSendGroupApplication() ([]*model_struct.LocalGroupRequest, error) {
	panic("implement me")
}

func (n NilIndexDB) GetConversationByUserID(userID string) (*model_struct.LocalConversation, error) {
	panic("implement me")
}

func (n NilIndexDB) GetAllConversationListDB() ([]*model_struct.LocalConversation, error) {
	panic("implement me")
}

func (n NilIndexDB) GetHiddenConversationList() ([]*model_struct.LocalConversation, error) {
	panic("implement me")
}

func (n NilIndexDB) GetAllConversationListToSync() ([]*model_struct.LocalConversation, error) {
	panic("implement me")
}

func (n NilIndexDB) GetConversationListSplitDB(offset, count int) ([]*model_struct.LocalConversation, error) {
	panic("implement me")
}

func (n NilIndexDB) BatchInsertConversationList(conversationList []*model_struct.LocalConversation) error {
	panic("implement me")
}

func (n NilIndexDB) InsertConversation(conversationList *model_struct.LocalConversation) error {
	panic("implement me")
}

func (n NilIndexDB) DeleteConversation(conversationID string) error {
	panic("implement me")
}

func (n NilIndexDB) GetConversation(conversationID string) (*model_struct.LocalConversation, error) {
	panic("implement me")
}

func (n NilIndexDB) UpdateConversation(c *model_struct.LocalConversation) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateConversationForSync(c *model_struct.LocalConversation) error {
	panic("implement me")
}

func (n NilIndexDB) BatchUpdateConversationList(conversationList []*model_struct.LocalConversation) error {
	panic("implement me")
}

func (n NilIndexDB) ConversationIfExists(conversationID string) (bool, error) {
	panic("implement me")
}

func (n NilIndexDB) ResetConversation(conversationID string) error {
	panic("implement me")
}

func (n NilIndexDB) ResetAllConversation() error {
	panic("implement me")
}

func (n NilIndexDB) ClearConversation(conversationID string) error {
	panic("implement me")
}

func (n NilIndexDB) ClearAllConversation() error {
	panic("implement me")
}

func (n NilIndexDB) SetConversationDraft(conversationID, draftText string) error {
	panic("implement me")
}

func (n NilIndexDB) RemoveConversationDraft(conversationID, draftText string) error {
	panic("implement me")
}

func (n NilIndexDB) UnPinConversation(conversationID string, isPinned int) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateColumnsConversation(conversationID string, args map[string]interface{}) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateAllConversation(conversation *model_struct.LocalConversation) error {
	panic("implement me")
}

func (n NilIndexDB) IncrConversationUnreadCount(conversationID string) error {
	panic("implement me")
}

func (n NilIndexDB) DecrConversationUnreadCount(conversationID string, count int64) (err error) {
	panic("implement me")
}

func (n NilIndexDB) GetTotalUnreadMsgCountDB() (totalUnreadCount int32, err error) {
	panic("implement me")
}

func (n NilIndexDB) SetMultipleConversationRecvMsgOpt(conversationIDList []string, opt int) (err error) {
	panic("implement me")
}

func (n NilIndexDB) GetMultipleConversationDB(conversationIDList []string) (result []*model_struct.LocalConversation, err error) {
	panic("implement me")
}

func (n NilIndexDB) GetJoinedSuperGroupList() ([]*model_struct.LocalGroup, error) {
	panic("implement me")
}

func (n NilIndexDB) GetJoinedSuperGroupIDList() ([]string, error) {
	panic("implement me")
}

func (n NilIndexDB) InsertSuperGroup(groupInfo *model_struct.LocalGroup) error {
	panic("implement me")
}

func (n NilIndexDB) DeleteAllSuperGroup() error {
	panic("implement me")
}

func (n NilIndexDB) GetSuperGroupInfoByGroupID(groupID string) (*model_struct.LocalGroup, error) {
	panic("implement me")
}

func (n NilIndexDB) UpdateSuperGroup(groupInfo *model_struct.LocalGroup) error {
	panic("implement me")
}

func (n NilIndexDB) DeleteSuperGroup(groupID string) error {
	panic("implement me")
}

func (n NilIndexDB) GetReadDiffusionGroupIDList() ([]string, error) {
	panic("implement me")
}

func (n NilIndexDB) GetGroupMemberInfoByGroupIDUserID(groupID, userID string) (*model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (n NilIndexDB) GetAllGroupMemberList() ([]model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (n NilIndexDB) GetAllGroupMemberUserIDList() ([]model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (n NilIndexDB) GetGroupMemberCount(groupID string) (uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) GetGroupSomeMemberInfo(groupID string, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (n NilIndexDB) GetGroupAdminID(groupID string) ([]string, error) {
	panic("implement me")
}

func (n NilIndexDB) GetGroupMemberListByGroupID(groupID string) ([]*model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (n NilIndexDB) GetGroupMemberListSplit(groupID string, filter int32, offset, count int) ([]*model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (n NilIndexDB) GetGroupMemberOwnerAndAdmin(groupID string) ([]*model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (n NilIndexDB) GetGroupMemberOwner(groupID string) (*model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (n NilIndexDB) GetGroupMemberListSplitByJoinTimeFilter(groupID string, offset, count int, joinTimeBegin, joinTimeEnd int64, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (n NilIndexDB) GetGroupOwnerAndAdminByGroupID(groupID string) ([]*model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (n NilIndexDB) GetGroupMemberUIDListByGroupID(groupID string) (result []string, err error) {
	panic("implement me")
}

func (n NilIndexDB) InsertGroupMember(groupMember *model_struct.LocalGroupMember) error {
	panic("implement me")
}

func (n NilIndexDB) BatchInsertGroupMember(groupMemberList []*model_struct.LocalGroupMember) error {
	panic("implement me")
}

func (n NilIndexDB) DeleteGroupMember(groupID, userID string) error {
	panic("implement me")
}

func (n NilIndexDB) DeleteGroupAllMembers(groupID string) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateGroupMember(groupMember *model_struct.LocalGroupMember) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateGroupMemberField(groupID, userID string, args map[string]interface{}) error {
	panic("implement me")
}

func (n NilIndexDB) GetGroupMemberInfoIfOwnerOrAdmin() ([]*model_struct.LocalGroupMember, error) {
	panic("implement me")
}

func (n NilIndexDB) SearchGroupMembersDB(keyword string, groupID string, isSearchMemberNickname, isSearchUserID bool, offset, count int) (result []*model_struct.LocalGroupMember, err error) {
	panic("implement me")
}

func (n NilIndexDB) InitSuperLocalErrChatLog(groupID string) {
	panic("implement me")
}

func (n NilIndexDB) SuperBatchInsertExceptionMsg(MessageList []*model_struct.LocalErrChatLog, groupID string) error {
	panic("implement me")
}

func (n NilIndexDB) GetAbnormalMsgSeq() (uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) GetAbnormalMsgSeqList() ([]uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) BatchInsertExceptionMsg(MessageList []*model_struct.LocalErrChatLog) error {
	panic("implement me")
}

func (n NilIndexDB) BatchInsertExceptionMsgController(MessageList []*model_struct.LocalErrChatLog) error {
	panic("implement me")
}

func (n NilIndexDB) GetSuperGroupAbnormalMsgSeq(groupID string) (uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) GetDepartmentMemberListByDepartmentID(departmentID string, args ...int) ([]*model_struct.LocalDepartmentMember, error) {
	panic("implement me")
}

func (n NilIndexDB) GetAllDepartmentMemberList() ([]*model_struct.LocalDepartmentMember, error) {
	panic("implement me")
}

func (n NilIndexDB) InsertDepartmentMember(departmentMember *model_struct.LocalDepartmentMember) error {
	panic("implement me")
}

func (n NilIndexDB) BatchInsertDepartmentMember(departmentMemberList []*model_struct.LocalDepartmentMember) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateDepartmentMember(departmentMember *model_struct.LocalDepartmentMember) error {
	panic("implement me")
}

func (n NilIndexDB) DeleteDepartmentMember(departmentID string, userID string) error {
	panic("implement me")
}

func (n NilIndexDB) GetDepartmentMemberListByUserID(userID string) ([]*model_struct.LocalDepartmentMember, error) {
	panic("implement me")
}

func (n NilIndexDB) BatchInsertTempCacheMessageList(MessageList []*model_struct.TempCacheLocalChatLog) error {
	panic("implement me")
}

func (n NilIndexDB) InsertTempCacheMessage(Message *model_struct.TempCacheLocalChatLog) error {
	panic("implement me")
}

func (n NilIndexDB) GetLoginUser(userID string) (*model_struct.LocalUser, error) {
	panic("implement me")
}

func (n NilIndexDB) UpdateLoginUser(user *model_struct.LocalUser) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateLoginUserByMap(user *model_struct.LocalUser, args map[string]interface{}) error {
	panic("implement me")
}

func (n NilIndexDB) InsertLoginUser(user *model_struct.LocalUser) error {
	panic("implement me")
}

func (n NilIndexDB) InsertFriend(friend *model_struct.LocalFriend) error {
	panic("implement me")
}

func (n NilIndexDB) DeleteFriendDB(friendUserID string) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateFriend(friend *model_struct.LocalFriend) error {
	panic("implement me")
}

func (n NilIndexDB) GetAllFriendList() ([]*model_struct.LocalFriend, error) {
	panic("implement me")
}

func (n NilIndexDB) SearchFriendList(keyword string, isSearchUserID, isSearchNickname, isSearchRemark bool) ([]*model_struct.LocalFriend, error) {
	panic("implement me")
}

func (n NilIndexDB) GetFriendInfoByFriendUserID(FriendUserID string) (*model_struct.LocalFriend, error) {
	panic("implement me")
}

func (n NilIndexDB) GetFriendInfoList(friendUserIDList []string) ([]*model_struct.LocalFriend, error) {
	panic("implement me")
}

func (n NilIndexDB) InsertFriendRequest(friendRequest *model_struct.LocalFriendRequest) error {
	panic("implement me")
}

func (n NilIndexDB) DeleteFriendRequestBothUserID(fromUserID, toUserID string) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateFriendRequest(friendRequest *model_struct.LocalFriendRequest) error {
	panic("implement me")
}

func (n NilIndexDB) GetRecvFriendApplication() ([]*model_struct.LocalFriendRequest, error) {
	panic("implement me")
}

func (n NilIndexDB) GetSendFriendApplication() ([]*model_struct.LocalFriendRequest, error) {
	panic("implement me")
}

func (n NilIndexDB) GetFriendApplicationByBothID(fromUserID, toUserID string) (*model_struct.LocalFriendRequest, error) {
	panic("implement me")
}

func (n NilIndexDB) InitSuperLocalChatLog(groupID string) {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupBatchInsertMessageList(MessageList []*model_struct.LocalChatLog, groupID string) error {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupInsertMessage(Message *model_struct.LocalChatLog, groupID string) error {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupDeleteAllMessage(groupID string) error {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupSearchMessageByKeyword(contentType []int, keywordList []string, keywordListMatchType int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupSearchMessageByContentType(contentType []int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupSearchMessageByContentTypeAndKeyword(contentType []int, keywordList []string, keywordListMatchType int, startTime, endTime int64, groupID string) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupBatchUpdateMessageList(MessageList []*model_struct.LocalChatLog) error {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupMessageIfExists(ClientMsgID string) (bool, error) {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupIsExistsInErrChatLogBySeq(seq int64) bool {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupMessageIfExistsBySeq(seq int64) (bool, error) {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupGetMessage(msg *sdk_struct.MsgStruct) (*model_struct.LocalChatLog, error) {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupGetAllUnDeleteMessageSeqList() ([]uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupUpdateColumnsMessage(ClientMsgID, groupID string, args map[string]interface{}) error {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupUpdateMessage(c *model_struct.LocalChatLog) error {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupUpdateMessageStatusBySourceID(sourceID string, status, sessionType int32) error {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupUpdateMessageTimeAndStatus(msg *sdk_struct.MsgStruct) error {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupGetMessageList(sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupGetMessageListNoTime(sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupGetSendingMessageList(groupID string) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupUpdateGroupMessageHasRead(msgIDList []string, groupID string) error {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupGetMultipleMessage(conversationIDList []string, groupID string) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupGetNormalMsgSeq() (uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupGetNormalMinSeq(groupID string) (uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupGetTestMessage(seq uint32) (*model_struct.LocalChatLog, error) {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupUpdateMsgSenderNickname(sendID, nickname string, sType int) error {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupUpdateMsgSenderFaceURL(sendID, faceURL string, sType int) error {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupUpdateMsgSenderFaceURLAndSenderNickname(sendID, faceURL, nickname string, sessionType int, groupID string) error {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupGetMsgSeqByClientMsgID(clientMsgID string, groupID string) (uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupGetMsgSeqListByGroupID(groupID string) ([]uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupGetMsgSeqListByPeerUserID(userID string) ([]uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupGetMsgSeqListBySelfUserID(userID string) ([]uint32, error) {
	panic("implement me")
}

func (n NilIndexDB) InsertWorkMomentsNotification(jsonDetail string) error {
	panic("implement me")
}

func (n NilIndexDB) GetWorkMomentsNotification(offset, count int) (WorkMomentsNotifications []*db.LocalWorkMomentsNotification, err error) {
	panic("implement me")
}

func (n NilIndexDB) GetWorkMomentsNotificationLimit(pageNumber, showNumber int) (WorkMomentsNotifications []*db.LocalWorkMomentsNotification, err error) {
	panic("implement me")
}

func (n NilIndexDB) InitWorkMomentsNotificationUnreadCount() error {
	panic("implement me")
}

func (n NilIndexDB) IncrWorkMomentsNotificationUnreadCount() error {
	panic("implement me")
}

func (n NilIndexDB) MarkAllWorkMomentsNotificationAsRead() (err error) {
	panic("implement me")
}

func (n NilIndexDB) GetWorkMomentsUnReadCount() (workMomentsNotificationUnReadCount db.LocalWorkMomentsNotificationUnreadCount, err error) {
	panic("implement me")
}

func (n NilIndexDB) ClearWorkMomentsNotification() (err error) {
	panic("implement me")
}

func (n NilIndexDB) Close() error {
	panic("implement me")
}

func (n NilIndexDB) SetChatLogFailedStatus() {
	panic("implement me")
}

func (n NilIndexDB) InitDB(userID string, dataDir string) error {
	panic("implement me")
}

func (n NilIndexDB) GetBlackListDB() ([]*model_struct.LocalBlack, error) {
	panic("implement me")
}

func (n NilIndexDB) GetBlackListUserID() (blackListUid []string, err error) {
	panic("implement me")
}

func (n NilIndexDB) GetBlackInfoByBlockUserID(blockUserID string) (*model_struct.LocalBlack, error) {
	panic("implement me")
}

func (n NilIndexDB) GetBlackInfoList(blockUserIDList []string) ([]*model_struct.LocalBlack, error) {
	panic("implement me")
}

func (n NilIndexDB) InsertBlack(black *model_struct.LocalBlack) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateBlack(black *model_struct.LocalBlack) error {
	panic("implement me")
}

func (n NilIndexDB) DeleteBlack(blockUserID string) error {
	panic("implement me")
}

func (n NilIndexDB) BatchInsertConversationUnreadMessageList(messageList []*model_struct.LocalConversationUnreadMessage) error {
	panic("implement me")
}

func (n NilIndexDB) DeleteConversationUnreadMessageList(conversationID string, sendTime int64) int64 {
	panic("implement me")
}

func (n NilIndexDB) SearchAllMessageByContentType(contentType int) ([]*model_struct.LocalChatLog, error) {
	panic("implement me")
}

func (n NilIndexDB) SuperGroupSearchAllMessageByContentType(superGroupID string, contentType int32) ([]*model_struct.LocalChatLog, error) {
	panic("implement me")
}

func (n NilIndexDB) GetMessageReactionExtension(msgID string) (result *model_struct.LocalChatLogReactionExtensions, err error) {
	panic("implement me")
}

func (n NilIndexDB) InsertMessageReactionExtension(messageReactionExtension *model_struct.LocalChatLogReactionExtensions) error {
	panic("implement me")
}

func (n NilIndexDB) UpdateMessageReactionExtension(c *model_struct.LocalChatLogReactionExtensions) error {
	panic("implement me")
}

func (n NilIndexDB) GetAndUpdateMessageReactionExtension(msgID string, m map[string]*server_api_params.KeyValue) error {
	panic("implement me")
}

func (n NilIndexDB) DeleteAndUpdateMessageReactionExtension(msgID string, m map[string]*server_api_params.KeyValue) error {
	panic("implement me")
}

func NewIndexDB(loginUserID string) *NilIndexDB {
	return &NilIndexDB{}
}
