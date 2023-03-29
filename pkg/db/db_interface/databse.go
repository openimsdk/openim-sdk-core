package db_interface

import (
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/sdk_struct"
)

type DataBase interface {
	InsertGroup(groupInfo *model_struct.LocalGroup) error
	DeleteGroup(groupID string) error
	UpdateGroup(groupInfo *model_struct.LocalGroup) error
	GetJoinedGroupListDB() ([]*model_struct.LocalGroup, error)
	GetGroupInfoByGroupID(groupID string) (*model_struct.LocalGroup, error)
	GetAllGroupInfoByGroupIDOrGroupName(keyword string, isSearchGroupID bool, isSearchGroupName bool) ([]*model_struct.LocalGroup, error)
	AddMemberCount(groupID string) error
	SubtractMemberCount(groupID string) error
	GetJoinedWorkingGroupIDList() ([]string, error)
	GetJoinedWorkingGroupList() ([]*model_struct.LocalGroup, error)
	GetMinSeq(ID string) (uint32, error)
	SetMinSeq(ID string, minSeq uint32) error
	GetUserMinSeq() (uint32, error)
	GetGroupMinSeq(groupID string) (uint32, error)
	InsertAdminGroupRequest(groupRequest *model_struct.LocalAdminGroupRequest) error
	DeleteAdminGroupRequest(groupID, userID string) error
	UpdateAdminGroupRequest(groupRequest *model_struct.LocalAdminGroupRequest) error
	GetAdminGroupApplication() ([]*model_struct.LocalAdminGroupRequest, error)
	BatchInsertMessageList(MessageList []*model_struct.LocalChatLog) error
	BatchInsertMessageListController(MessageList []*model_struct.LocalChatLog) error
	InsertMessage(Message *model_struct.LocalChatLog) error
	InsertMessageController(message *model_struct.LocalChatLog) error
	SearchMessageByKeyword(contentType []int, keywordList []string, keywordListMatchType int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error)
	SearchMessageByKeywordController(contentType []int, keywordList []string, keywordListMatchType int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error)
	SearchMessageByContentType(contentType []int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error)
	SearchMessageByContentTypeController(contentType []int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error)
	SearchMessageByContentTypeAndKeyword(contentType []int, keywordList []string, keywordListMatchType int, startTime, endTime int64) (result []*model_struct.LocalChatLog, err error)
	SearchMessageByContentTypeAndKeywordController(contentType []int, keywordList []string, keywordListMatchType int, startTime, endTime int64, operationID string) (result []*model_struct.LocalChatLog, err error)
	BatchUpdateMessageList(MessageList []*model_struct.LocalChatLog) error
	BatchSpecialUpdateMessageList(MessageList []*model_struct.LocalChatLog) error
	MessageIfExists(ClientMsgID string) (bool, error)
	IsExistsInErrChatLogBySeq(seq int64) bool
	MessageIfExistsBySeq(seq int64) (bool, error)
	GetMessage(ClientMsgID string) (*model_struct.LocalChatLog, error)
	GetMessageController(msg *sdk_struct.MsgStruct) (*model_struct.LocalChatLog, error)
	GetAllUnDeleteMessageSeqList() ([]uint32, error)
	UpdateColumnsMessageList(clientMsgIDList []string, args map[string]interface{}) error
	UpdateColumnsMessage(ClientMsgID string, args map[string]interface{}) error
	UpdateColumnsMessageController(ClientMsgID string, groupID string, sessionType int32, args map[string]interface{}) error
	UpdateMessage(c *model_struct.LocalChatLog) error
	UpdateMessageController(c *model_struct.LocalChatLog) error
	DeleteAllMessage() error
	UpdateMessageStatusBySourceID(sourceID string, status, sessionType int32) error
	UpdateMessageStatusBySourceIDController(sourceID string, status, sessionType int32) error
	UpdateMessageTimeAndStatus(clientMsgID string, serverMsgID string, sendTime int64, status int32) error
	UpdateMessageTimeAndStatusController(msg *sdk_struct.MsgStruct) error
	GetMessageList(sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error)
	GetMessageListController(sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error)
	GetMessageListNoTime(sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error)
	GetMessageListNoTimeController(sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error)
	GetSendingMessageList() (result []*model_struct.LocalChatLog, err error)
	UpdateSingleMessageHasRead(sendID string, msgIDList []string) error
	UpdateGroupMessageHasRead(msgIDList []string, sessionType int32) error
	UpdateGroupMessageHasReadController(msgIDList []string, groupID string, sessionType int32) error
	GetMultipleMessage(msgIDList []string) (result []*model_struct.LocalChatLog, err error)
	GetMultipleMessageController(msgIDList []string, groupID string, sessionType int32) (result []*model_struct.LocalChatLog, err error)
	GetNormalMsgSeq() (uint32, error)
	GetLostMsgSeqList(minSeqInSvr uint32) ([]uint32, error)
	GetSuperGroupNormalMsgSeq(groupID string) (uint32, error)
	GetTestMessage(seq uint32) (*model_struct.LocalChatLog, error)
	UpdateMsgSenderNickname(sendID, nickname string, sType int) error
	UpdateMsgSenderFaceURL(sendID, faceURL string, sType int) error
	UpdateMsgSenderFaceURLAndSenderNickname(sendID, faceURL, nickname string, sessionType int) error
	UpdateMsgSenderFaceURLAndSenderNicknameController(sendID, faceURL, nickname string, sessionType int, groupID string) error
	GetMsgSeqByClientMsgID(clientMsgID string) (uint32, error)
	GetMsgSeqByClientMsgIDController(m *sdk_struct.MsgStruct) (uint32, error)
	GetMsgSeqListByGroupID(groupID string) ([]uint32, error)
	GetMsgSeqListByPeerUserID(userID string) ([]uint32, error)
	GetMsgSeqListBySelfUserID(userID string) ([]uint32, error)
	GetSubDepartmentList(departmentID string, args ...int) ([]*model_struct.LocalDepartment, error)
	InsertDepartment(department *model_struct.LocalDepartment) error
	UpdateDepartment(department *model_struct.LocalDepartment) error
	DeleteDepartment(departmentID string) error
	GetDepartmentInfo(departmentID string) (*model_struct.LocalDepartment, error)
	GetAllDepartmentList() ([]*model_struct.LocalDepartment, error)
	GetParentDepartmentList(departmentID string) ([]*model_struct.LocalDepartment, error)
	GetDepartmentList(departmentList *[]*model_struct.LocalDepartment, departmentID string) error
	GetParentDepartment(departmentID string) (model_struct.LocalDepartment, error)
	SearchDepartmentMember(keyWord string, isSearchUserName, isSearchEmail, isSearchMobile, isSearchPosition, isSearchTelephone, isSearchUserEnglishName, isSearchUserID bool, offset, count int) ([]*model_struct.SearchDepartmentMemberResult, error)
	SearchDepartment(keyWord string, offset, count int) ([]*model_struct.LocalDepartment, error)
	InsertGroupRequest(groupRequest *model_struct.LocalGroupRequest) error
	DeleteGroupRequest(groupID, userID string) error
	UpdateGroupRequest(groupRequest *model_struct.LocalGroupRequest) error
	GetSendGroupApplication() ([]*model_struct.LocalGroupRequest, error)
	GetConversationByUserID(userID string) (*model_struct.LocalConversation, error)
	GetAllConversationListDB() ([]*model_struct.LocalConversation, error)
	GetHiddenConversationList() ([]*model_struct.LocalConversation, error)
	GetAllConversationListToSync() ([]*model_struct.LocalConversation, error)
	GetConversationListSplitDB(offset, count int) ([]*model_struct.LocalConversation, error)
	BatchInsertConversationList(conversationList []*model_struct.LocalConversation) error
	InsertConversation(conversationList *model_struct.LocalConversation) error
	DeleteConversation(conversationID string) error
	GetConversation(conversationID string) (*model_struct.LocalConversation, error)
	UpdateConversation(c *model_struct.LocalConversation) error
	UpdateConversationForSync(c *model_struct.LocalConversation) error
	BatchUpdateConversationList(conversationList []*model_struct.LocalConversation) error
	ConversationIfExists(conversationID string) (bool, error)
	ResetConversation(conversationID string) error
	ResetAllConversation() error
	ClearConversation(conversationID string) error
	ClearAllConversation() error
	SetConversationDraft(conversationID, draftText string) error
	RemoveConversationDraft(conversationID, draftText string) error
	UnPinConversation(conversationID string, isPinned int) error
	UpdateColumnsConversation(conversationID string, args map[string]interface{}) error
	UpdateAllConversation(conversation *model_struct.LocalConversation) error
	IncrConversationUnreadCount(conversationID string) error
	DecrConversationUnreadCount(conversationID string, count int64) (err error)
	GetTotalUnreadMsgCountDB() (totalUnreadCount int32, err error)
	SetMultipleConversationRecvMsgOpt(conversationIDList []string, opt int) (err error)
	GetMultipleConversationDB(conversationIDList []string) (result []*model_struct.LocalConversation, err error)
	GetJoinedSuperGroupList() ([]*model_struct.LocalGroup, error)
	GetJoinedSuperGroupIDList() ([]string, error)
	InsertSuperGroup(groupInfo *model_struct.LocalGroup) error
	DeleteAllSuperGroup() error
	GetSuperGroupInfoByGroupID(groupID string) (*model_struct.LocalGroup, error)
	UpdateSuperGroup(groupInfo *model_struct.LocalGroup) error
	DeleteSuperGroup(groupID string) error
	GetReadDiffusionGroupIDList() ([]string, error)
	GetGroupMemberInfoByGroupIDUserID(groupID, userID string) (*model_struct.LocalGroupMember, error)
	GetAllGroupMemberList() ([]model_struct.LocalGroupMember, error)
	GetAllGroupMemberUserIDList() ([]model_struct.LocalGroupMember, error)
	GetGroupMemberCount(groupID string) (uint32, error)
	GetGroupSomeMemberInfo(groupID string, userIDList []string) ([]*model_struct.LocalGroupMember, error)
	GetGroupAdminID(groupID string) ([]string, error)
	GetGroupMemberListByGroupID(groupID string) ([]*model_struct.LocalGroupMember, error)
	GetGroupMemberListSplit(groupID string, filter int32, offset, count int) ([]*model_struct.LocalGroupMember, error)
	GetGroupMemberOwnerAndAdmin(groupID string) ([]*model_struct.LocalGroupMember, error)
	GetGroupMemberOwner(groupID string) (*model_struct.LocalGroupMember, error)
	GetGroupMemberListSplitByJoinTimeFilter(groupID string, offset, count int, joinTimeBegin, joinTimeEnd int64, userIDList []string) ([]*model_struct.LocalGroupMember, error)
	GetGroupOwnerAndAdminByGroupID(groupID string) ([]*model_struct.LocalGroupMember, error)
	GetGroupMemberUIDListByGroupID(groupID string) (result []string, err error)
	InsertGroupMember(groupMember *model_struct.LocalGroupMember) error
	BatchInsertGroupMember(groupMemberList []*model_struct.LocalGroupMember) error
	DeleteGroupMember(groupID, userID string) error
	DeleteGroupAllMembers(groupID string) error
	UpdateGroupMember(groupMember *model_struct.LocalGroupMember) error
	UpdateGroupMemberField(groupID, userID string, args map[string]interface{}) error
	GetGroupMemberInfoIfOwnerOrAdmin() ([]*model_struct.LocalGroupMember, error)
	SearchGroupMembersDB(keyword string, groupID string, isSearchMemberNickname, isSearchUserID bool, offset, count int) (result []*model_struct.LocalGroupMember, err error)
	InitSuperLocalErrChatLog(groupID string)
	SuperBatchInsertExceptionMsg(MessageList []*model_struct.LocalErrChatLog, groupID string) error
	GetAbnormalMsgSeq() (uint32, error)
	GetAbnormalMsgSeqList() ([]uint32, error)
	BatchInsertExceptionMsg(MessageList []*model_struct.LocalErrChatLog) error
	BatchInsertExceptionMsgController(MessageList []*model_struct.LocalErrChatLog) error
	GetSuperGroupAbnormalMsgSeq(groupID string) (uint32, error)
	GetDepartmentMemberListByDepartmentID(departmentID string, args ...int) ([]*model_struct.LocalDepartmentMember, error)
	GetAllDepartmentMemberList() ([]*model_struct.LocalDepartmentMember, error)
	InsertDepartmentMember(departmentMember *model_struct.LocalDepartmentMember) error
	BatchInsertDepartmentMember(departmentMemberList []*model_struct.LocalDepartmentMember) error
	UpdateDepartmentMember(departmentMember *model_struct.LocalDepartmentMember) error
	DeleteDepartmentMember(departmentID string, userID string) error
	GetDepartmentMemberListByUserID(userID string) ([]*model_struct.LocalDepartmentMember, error)
	BatchInsertTempCacheMessageList(MessageList []*model_struct.TempCacheLocalChatLog) error
	InsertTempCacheMessage(Message *model_struct.TempCacheLocalChatLog) error
	GetLoginUser(userID string) (*model_struct.LocalUser, error)
	UpdateLoginUser(user *model_struct.LocalUser) error
	UpdateLoginUserByMap(user *model_struct.LocalUser, args map[string]interface{}) error
	InsertLoginUser(user *model_struct.LocalUser) error
	InsertFriend(friend *model_struct.LocalFriend) error
	DeleteFriendDB(friendUserID string) error
	UpdateFriend(friend *model_struct.LocalFriend) error
	GetAllFriendList() ([]*model_struct.LocalFriend, error)
	SearchFriendList(keyword string, isSearchUserID, isSearchNickname, isSearchRemark bool) ([]*model_struct.LocalFriend, error)
	GetFriendInfoByFriendUserID(FriendUserID string) (*model_struct.LocalFriend, error)
	GetFriendInfoList(friendUserIDList []string) ([]*model_struct.LocalFriend, error)
	InsertFriendRequest(friendRequest *model_struct.LocalFriendRequest) error
	DeleteFriendRequestBothUserID(fromUserID, toUserID string) error
	UpdateFriendRequest(friendRequest *model_struct.LocalFriendRequest) error
	GetRecvFriendApplication() ([]*model_struct.LocalFriendRequest, error)
	GetSendFriendApplication() ([]*model_struct.LocalFriendRequest, error)
	GetFriendApplicationByBothID(fromUserID, toUserID string) (*model_struct.LocalFriendRequest, error)
	InitSuperLocalChatLog(groupID string)
	SuperGroupBatchInsertMessageList(MessageList []*model_struct.LocalChatLog, groupID string) error
	SuperGroupInsertMessage(Message *model_struct.LocalChatLog, groupID string) error
	SuperGroupDeleteAllMessage(groupID string) error
	SuperGroupSearchMessageByKeyword(contentType []int, keywordList []string, keywordListMatchType int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error)
	SuperGroupSearchMessageByContentType(contentType []int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error)
	SuperGroupSearchMessageByContentTypeAndKeyword(contentType []int, keywordList []string, keywordListMatchType int, startTime, endTime int64, groupID string) (result []*model_struct.LocalChatLog, err error)
	SuperGroupBatchUpdateMessageList(MessageList []*model_struct.LocalChatLog) error
	SuperGroupMessageIfExists(ClientMsgID string) (bool, error)
	SuperGroupIsExistsInErrChatLogBySeq(seq int64) bool
	SuperGroupMessageIfExistsBySeq(seq int64) (bool, error)
	SuperGroupGetMessage(msg *sdk_struct.MsgStruct) (*model_struct.LocalChatLog, error)
	SuperGroupGetAllUnDeleteMessageSeqList() ([]uint32, error)
	SuperGroupUpdateColumnsMessage(ClientMsgID, groupID string, args map[string]interface{}) error
	SuperGroupUpdateMessage(c *model_struct.LocalChatLog) error
	SuperGroupUpdateSpecificContentTypeMessage(contentType int, groupID string, args map[string]interface{}) error
	SuperGroupUpdateMessageStatusBySourceID(sourceID string, status, sessionType int32) error
	SuperGroupUpdateMessageTimeAndStatus(msg *sdk_struct.MsgStruct) error
	SuperGroupGetMessageList(sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error)
	SuperGroupGetMessageListNoTime(sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error)
	SuperGroupGetSendingMessageList(groupID string) (result []*model_struct.LocalChatLog, err error)
	SuperGroupUpdateGroupMessageHasRead(msgIDList []string, groupID string) error
	SuperGroupUpdateGroupMessageFields(msgIDList []string, groupID string, args map[string]interface{}) error
	SuperGroupGetMultipleMessage(conversationIDList []string, groupID string) (result []*model_struct.LocalChatLog, err error)
	SuperGroupGetNormalMsgSeq() (uint32, error)
	SuperGroupGetNormalMinSeq(groupID string) (uint32, error)
	SuperGroupGetTestMessage(seq uint32) (*model_struct.LocalChatLog, error)
	SuperGroupUpdateMsgSenderNickname(sendID, nickname string, sType int) error
	SuperGroupUpdateMsgSenderFaceURL(sendID, faceURL string, sType int) error
	SuperGroupUpdateMsgSenderFaceURLAndSenderNickname(sendID, faceURL, nickname string, sessionType int, groupID string) error
	SuperGroupGetMsgSeqByClientMsgID(clientMsgID string, groupID string) (uint32, error)
	SuperGroupGetMsgSeqListByGroupID(groupID string) ([]uint32, error)
	SuperGroupGetMsgSeqListByPeerUserID(userID string) ([]uint32, error)
	SuperGroupGetMsgSeqListBySelfUserID(userID string) ([]uint32, error)
	SuperGroupGetAlreadyExistSeqList(groupID string, lostSeqList []uint32) (seqList []uint32, err error)
	InsertWorkMomentsNotification(jsonDetail string) error
	GetWorkMomentsNotification(offset, count int) (WorkMomentsNotifications []*db.LocalWorkMomentsNotification, err error)
	GetWorkMomentsNotificationLimit(pageNumber, showNumber int) (WorkMomentsNotifications []*db.LocalWorkMomentsNotification, err error)
	InitWorkMomentsNotificationUnreadCount() error
	IncrWorkMomentsNotificationUnreadCount() error
	MarkAllWorkMomentsNotificationAsRead() (err error)
	GetWorkMomentsUnReadCount() (workMomentsNotificationUnReadCount db.LocalWorkMomentsNotificationUnreadCount, err error)
	ClearWorkMomentsNotification() (err error)
	Close() error
	SetChatLogFailedStatus()
	InitDB(userID string, dataDir string) error
	GetBlackListDB() ([]*model_struct.LocalBlack, error)
	GetBlackListUserID() (blackListUid []string, err error)
	GetBlackInfoByBlockUserID(blockUserID string) (*model_struct.LocalBlack, error)
	GetBlackInfoList(blockUserIDList []string) ([]*model_struct.LocalBlack, error)
	InsertBlack(black *model_struct.LocalBlack) error
	UpdateBlack(black *model_struct.LocalBlack) error
	DeleteBlack(blockUserID string) error
	BatchInsertConversationUnreadMessageList(messageList []*model_struct.LocalConversationUnreadMessage) error
	DeleteConversationUnreadMessageList(conversationID string, sendTime int64) int64
	SearchAllMessageByContentType(contentType int) ([]*model_struct.LocalChatLog, error)
	SuperGroupSearchAllMessageByContentType(superGroupID string, contentType int32) ([]*model_struct.LocalChatLog, error)
	GetMessageReactionExtension(msgID string) (result *model_struct.LocalChatLogReactionExtensions, err error)
	InsertMessageReactionExtension(messageReactionExtension *model_struct.LocalChatLogReactionExtensions) error
	UpdateMessageReactionExtension(c *model_struct.LocalChatLogReactionExtensions) error
	GetAndUpdateMessageReactionExtension(msgID string, m map[string]*server_api_params.KeyValue) error
	DeleteAndUpdateMessageReactionExtension(msgID string, m map[string]*server_api_params.KeyValue) error
	GetMultipleMessageReactionExtension(msgIDList []string) (result []*model_struct.LocalChatLogReactionExtensions, err error)
	DeleteMessageReactionExtension(msgID string) error
}
