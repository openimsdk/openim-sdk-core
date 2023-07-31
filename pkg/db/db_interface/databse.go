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

package db_interface

import (
	"context"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/sdk_struct"
)

type GroupDatabase interface {
	InsertGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error
	DeleteGroup(ctx context.Context, groupID string) error
	UpdateGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error
	GetJoinedGroupListDB(ctx context.Context) ([]*model_struct.LocalGroup, error)
	GetGroups(ctx context.Context, groupIDs []string) ([]*model_struct.LocalGroup, error)
	GetGroupInfoByGroupID(ctx context.Context, groupID string) (*model_struct.LocalGroup, error)
	GetAllGroupInfoByGroupIDOrGroupName(ctx context.Context, keyword string, isSearchGroupID bool, isSearchGroupName bool) ([]*model_struct.LocalGroup, error)
	AddMemberCount(ctx context.Context, groupID string) error
	SubtractMemberCount(ctx context.Context, groupID string) error
	GetJoinedWorkingGroupIDList(ctx context.Context) ([]string, error)
	GetJoinedWorkingGroupList(ctx context.Context) ([]*model_struct.LocalGroup, error)
	GetUserJoinedGroupIDs(ctx context.Context, userID string) ([]string, error)

	InsertAdminGroupRequest(ctx context.Context, groupRequest *model_struct.LocalAdminGroupRequest) error
	DeleteAdminGroupRequest(ctx context.Context, groupID, userID string) error
	UpdateAdminGroupRequest(ctx context.Context, groupRequest *model_struct.LocalAdminGroupRequest) error
	GetAdminGroupApplication(ctx context.Context) ([]*model_struct.LocalAdminGroupRequest, error)
	InsertGroupRequest(ctx context.Context, groupRequest *model_struct.LocalGroupRequest) error
	DeleteGroupRequest(ctx context.Context, groupID, userID string) error
	UpdateGroupRequest(ctx context.Context, groupRequest *model_struct.LocalGroupRequest) error
	GetSendGroupApplication(ctx context.Context) ([]*model_struct.LocalGroupRequest, error)
	GetJoinedSuperGroupList(ctx context.Context) ([]*model_struct.LocalGroup, error)
	InsertSuperGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error
	DeleteAllSuperGroup(ctx context.Context) error
	GetSuperGroupInfoByGroupID(ctx context.Context, groupID string) (*model_struct.LocalGroup, error)
	UpdateSuperGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error
	DeleteSuperGroup(ctx context.Context, groupID string) error
	GetGroupMemberInfoByGroupIDUserID(ctx context.Context, groupID, userID string) (*model_struct.LocalGroupMember, error)
	GetAllGroupMemberList(ctx context.Context) ([]model_struct.LocalGroupMember, error)
	GetAllGroupMemberUserIDList(ctx context.Context) ([]model_struct.LocalGroupMember, error)
	GetGroupMemberCount(ctx context.Context, groupID string) (int32, error)
	GetGroupSomeMemberInfo(ctx context.Context, groupID string, userIDList []string) ([]*model_struct.LocalGroupMember, error)
	GetGroupAdminID(ctx context.Context, groupID string) ([]string, error)
	GetGroupMemberListByGroupID(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error)
	GetGroupMemberListSplit(ctx context.Context, groupID string, filter int32, offset, count int) ([]*model_struct.LocalGroupMember, error)
	GetGroupMemberOwnerAndAdminDB(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error)
	GetGroupMemberOwner(ctx context.Context, groupID string) (*model_struct.LocalGroupMember, error)
	GetGroupMemberListSplitByJoinTimeFilter(ctx context.Context, groupID string, offset, count int, joinTimeBegin, joinTimeEnd int64, userIDList []string) ([]*model_struct.LocalGroupMember, error)
	GetGroupOwnerAndAdminByGroupID(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error)
	GetGroupMemberUIDListByGroupID(ctx context.Context, groupID string) (result []string, err error)
	GetGroupMemberAllGroupIDs(ctx context.Context) ([]string, error)
	InsertGroupMember(ctx context.Context, groupMember *model_struct.LocalGroupMember) error
	BatchInsertGroupMember(ctx context.Context, groupMemberList []*model_struct.LocalGroupMember) error
	DeleteGroupMember(ctx context.Context, groupID, userID string) error
	DeleteGroupAllMembers(ctx context.Context, groupID string) error
	UpdateGroupMember(ctx context.Context, groupMember *model_struct.LocalGroupMember) error
	UpdateGroupMemberField(ctx context.Context, groupID, userID string, args map[string]interface{}) error
	GetGroupMemberInfoIfOwnerOrAdmin(ctx context.Context) ([]*model_struct.LocalGroupMember, error)
	SearchGroupMembersDB(ctx context.Context, keyword string, groupID string, isSearchMemberNickname, isSearchUserID bool, offset, count int) (result []*model_struct.LocalGroupMember, err error)
}

type MessageDatabase interface {
	BatchInsertMessageList(ctx context.Context, conversationID string, MessageList []*model_struct.LocalChatLog) error
	//BatchInsertMessageListController(ctx context.Context, MessageList []*model_struct.LocalChatLog) error
	InsertMessage(ctx context.Context, conversationID string, Message *model_struct.LocalChatLog) error
	//InsertMessageController(ctx context.Context, message *model_struct.LocalChatLog) error
	SearchMessageByKeyword(ctx context.Context, contentType []int, keywordList []string, keywordListMatchType int, conversationID string, startTime, endTime int64, offset, count int) (result []*model_struct.LocalChatLog, err error)
	//SearchMessageByKeywordController(ctx context.Context, contentType []int, keywordList []string, keywordListMatchType int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error)
	SearchMessageByContentType(ctx context.Context, contentType []int, conversationID string, startTime, endTime int64, offset, count int) (result []*model_struct.LocalChatLog, err error)
	//SearchMessageByContentTypeController(ctx context.Context, contentType []int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error)
	SearchMessageByContentTypeAndKeyword(ctx context.Context, contentType []int, conversationID string, keywordList []string, keywordListMatchType int, startTime, endTime int64) (result []*model_struct.LocalChatLog, err error)
	//SearchMessageByContentTypeAndKeywordController(ctx context.Context, contentType []int, keywordList []string, keywordListMatchType int, startTime, endTime int64) (result []*model_struct.LocalChatLog, err error)
	//BatchUpdateMessageList(ctx context.Context, MessageList []*model_struct.LocalChatLog) error
	//BatchSpecialUpdateMessageList(ctx context.Context, MessageList []*model_struct.LocalChatLog) error
	MessageIfExists(ctx context.Context, ClientMsgID string) (bool, error)
	IsExistsInErrChatLogBySeq(ctx context.Context, seq int64) bool
	MessageIfExistsBySeq(ctx context.Context, seq int64) (bool, error)
	GetMessage(ctx context.Context, conversationID, clientMsgID string) (*model_struct.LocalChatLog, error)
	GetMessageBySeq(ctx context.Context, conversationID string, seq int64) (*model_struct.LocalChatLog, error)
	//GetMessageController(ctx context.Context, conversationID, clientMsgID string) (*model_struct.LocalChatLog, error)
	UpdateColumnsMessageList(ctx context.Context, clientMsgIDList []string, args map[string]interface{}) error
	UpdateColumnsMessage(ctx context.Context, conversationID string, ClientMsgID string, args map[string]interface{}) error
	//UpdateColumnsMessageController(ctx context.Context, ClientMsgID string, groupID string, sessionType int32, args map[string]interface{}) error
	UpdateMessage(ctx context.Context, conversationID string, c *model_struct.LocalChatLog) error
	UpdateMessageBySeq(ctx context.Context, conversationID string, c *model_struct.LocalChatLog) error
	//UpdateMessageController(ctx context.Context, c *model_struct.LocalChatLog) error
	DeleteAllMessage(ctx context.Context) error
	UpdateMessageStatusBySourceID(ctx context.Context, sourceID string, status, sessionType int32) error
	//UpdateMessageStatusBySourceIDController(ctx context.Context, sourceID string, status, sessionType int32) error
	UpdateMessageTimeAndStatus(ctx context.Context, conversationID, clientMsgID string, serverMsgID string, sendTime int64, status int32) error
	//UpdateMessageTimeAndStatusController(ctx context.Context, msg *sdk_struct.MsgStruct) error
	GetMessageList(ctx context.Context, conversationID string, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error)
	//GetMessageListController(ctx context.Context, sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error)
	GetMessageListNoTime(ctx context.Context, conversationID string, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error)
	//GetMessageListNoTimeController(ctx context.Context, sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error)
	GetSendingMessageList(ctx context.Context) (result []*model_struct.LocalChatLog, err error)
	MarkConversationMessageAsReadDB(ctx context.Context, conversationID string, msgIDs []string) (rowsAffected int64, err error)
	MarkConversationMessageAsReadBySeqs(ctx context.Context, conversationID string, seqs []int64) (rowsAffected int64, err error)
	GetUnreadMessage(ctx context.Context, conversationID string) (result []*model_struct.LocalChatLog, err error)
	MarkConversationAllMessageAsRead(ctx context.Context, conversationID string) (rowsAffected int64, err error)
	GetMessagesByClientMsgIDs(ctx context.Context, conversationID string, msgIDs []string) (result []*model_struct.LocalChatLog, err error)
	GetMessagesBySeqs(ctx context.Context, conversationID string, seqs []int64) (result []*model_struct.LocalChatLog, err error)
	GetConversationNormalMsgSeq(ctx context.Context, conversationID string) (int64, error)
	GetConversationPeerNormalMsgSeq(ctx context.Context, conversationID string) (int64, error)

	GetTestMessage(ctx context.Context, seq uint32) (*model_struct.LocalChatLog, error)
	UpdateMsgSenderNickname(ctx context.Context, sendID, nickname string, sType int) error
	UpdateMsgSenderFaceURL(ctx context.Context, sendID, faceURL string, sType int) error
	//UpdateMsgSenderFaceURLAndSenderNicknameController(ctx context.Context, sendID, faceURL, nickname string, sessionType int) error
	UpdateMsgSenderFaceURLAndSenderNickname(ctx context.Context, conversationID, sendID, faceURL, nickname string) error
	GetMsgSeqByClientMsgID(ctx context.Context, clientMsgID string) (uint32, error)
	//GetMsgSeqByClientMsgIDController(ctx context.Context, m *sdk_struct.MsgStruct) (uint32, error)
	GetMsgSeqListByGroupID(ctx context.Context, groupID string) ([]uint32, error)
	GetMsgSeqListByPeerUserID(ctx context.Context, userID string) ([]uint32, error)
	GetMsgSeqListBySelfUserID(ctx context.Context, userID string) ([]uint32, error)
	InitSuperLocalErrChatLog(ctx context.Context, groupID string)
	SuperBatchInsertExceptionMsg(ctx context.Context, MessageList []*model_struct.LocalErrChatLog, groupID string) error
	GetAbnormalMsgSeq(ctx context.Context) (int64, error)
	GetAbnormalMsgSeqList(ctx context.Context) ([]int64, error)
	BatchInsertExceptionMsg(ctx context.Context, MessageList []*model_struct.LocalErrChatLog) error
	GetConversationAbnormalMsgSeq(ctx context.Context, groupID string) (int64, error)
	BatchInsertTempCacheMessageList(ctx context.Context, MessageList []*model_struct.TempCacheLocalChatLog) error
	InsertTempCacheMessage(ctx context.Context, Message *model_struct.TempCacheLocalChatLog) error
	InitSuperLocalChatLog(ctx context.Context, groupID string)
	SuperGroupBatchInsertMessageList(ctx context.Context, MessageList []*model_struct.LocalChatLog, groupID string) error
	SuperGroupInsertMessage(ctx context.Context, Message *model_struct.LocalChatLog, groupID string) error
	SuperGroupDeleteAllMessage(ctx context.Context, groupID string) error
	DeleteConversationAllMessages(ctx context.Context, conversationID string) error
	MarkDeleteConversationAllMessages(ctx context.Context, conversationID string) error
	SuperGroupSearchMessageByKeyword(ctx context.Context, contentType []int, keywordList []string, keywordListMatchType int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error)
	//SuperGroupSearchMessageByContentType(ctx context.Context, contentType []int, sourceID string, startTime, endTime int64, sessionType, offset, count int) (result []*model_struct.LocalChatLog, err error)
	SuperGroupSearchMessageByContentTypeAndKeyword(ctx context.Context, contentType []int, keywordList []string, keywordListMatchType int, startTime, endTime int64, groupID string) (result []*model_struct.LocalChatLog, err error)

	SuperGroupBatchUpdateMessageList(ctx context.Context, MessageList []*model_struct.LocalChatLog) error
	SuperGroupMessageIfExists(ctx context.Context, ClientMsgID string) (bool, error)
	SuperGroupIsExistsInErrChatLogBySeq(ctx context.Context, seq int64) bool
	SuperGroupMessageIfExistsBySeq(ctx context.Context, seq int64) (bool, error)
	SuperGroupGetMessage(ctx context.Context, msg *sdk_struct.MsgStruct) (*model_struct.LocalChatLog, error)
	SuperGroupGetAllUnDeleteMessageSeqList(ctx context.Context) ([]uint32, error)
	SuperGroupUpdateColumnsMessage(ctx context.Context, ClientMsgID, groupID string, args map[string]interface{}) error
	SuperGroupUpdateMessage(ctx context.Context, c *model_struct.LocalChatLog) error
	SuperGroupUpdateMessageStatusBySourceID(ctx context.Context, sourceID string, status, sessionType int32) error
	SuperGroupUpdateMessageTimeAndStatus(ctx context.Context, msg *sdk_struct.MsgStruct) error
	SuperGroupGetMessageList(ctx context.Context, sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error)
	SuperGroupGetMessageListNoTime(ctx context.Context, sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error)
	SuperGroupGetSendingMessageList(ctx context.Context, groupID string) (result []*model_struct.LocalChatLog, err error)
	SuperGroupUpdateGroupMessageHasRead(ctx context.Context, msgIDList []string, groupID string) error
	//SuperGroupUpdateGroupMessageFields(ctx context.Context, msgIDList []string, groupID string, args map[string]interface{}) error

	SuperGroupUpdateMsgSenderNickname(ctx context.Context, sendID, nickname string, sType int) error
	SuperGroupUpdateMsgSenderFaceURL(ctx context.Context, sendID, faceURL string, sType int) error
	SuperGroupUpdateMsgSenderFaceURLAndSenderNickname(ctx context.Context, sendID, faceURL, nickname string, sessionType int, groupID string) error
	SuperGroupGetMsgSeqByClientMsgID(ctx context.Context, clientMsgID string, groupID string) (uint32, error)
	SuperGroupGetMsgSeqListByGroupID(ctx context.Context, groupID string) ([]uint32, error)
	SuperGroupGetMsgSeqListByPeerUserID(ctx context.Context, userID string) ([]uint32, error)
	SuperGroupGetMsgSeqListBySelfUserID(ctx context.Context, userID string) ([]uint32, error)
	GetAlreadyExistSeqList(ctx context.Context, conversationID string, lostSeqList []int64) (seqList []int64, err error)

	BatchInsertConversationUnreadMessageList(ctx context.Context, messageList []*model_struct.LocalConversationUnreadMessage) error
	DeleteConversationUnreadMessageList(ctx context.Context, conversationID string, sendTime int64) int64
	DeleteConversationMsgs(ctx context.Context, conversationID string, msgIDs []string) error
	//DeleteConversationMsgsBySeqs(ctx context.Context, conversationID string, seqs []int64) error
	SetNotificationSeq(ctx context.Context, conversationID string, seq int64) error
	GetNotificationAllSeqs(ctx context.Context) ([]*model_struct.NotificationSeqs, error)
}

type ConversationDatabase interface {
	GetConversationByUserID(ctx context.Context, userID string) (*model_struct.LocalConversation, error)
	GetAllConversationListDB(ctx context.Context) ([]*model_struct.LocalConversation, error)
	GetHiddenConversationList(ctx context.Context) ([]*model_struct.LocalConversation, error)
	GetAllConversations(ctx context.Context) ([]*model_struct.LocalConversation, error)
	GetAllSingleConversationIDList(ctx context.Context) (result []string, err error)
	GetAllConversationIDList(ctx context.Context) (result []string, err error)
	GetConversationListSplitDB(ctx context.Context, offset, count int) ([]*model_struct.LocalConversation, error)
	BatchInsertConversationList(ctx context.Context, conversationList []*model_struct.LocalConversation) error
	UpdateOrCreateConversations(ctx context.Context, conversationList []*model_struct.LocalConversation) error
	InsertConversation(ctx context.Context, conversationList *model_struct.LocalConversation) error
	DeleteConversation(ctx context.Context, conversationID string) error
	GetConversation(ctx context.Context, conversationID string) (*model_struct.LocalConversation, error)
	UpdateConversation(ctx context.Context, c *model_struct.LocalConversation) error
	UpdateConversationForSync(ctx context.Context, c *model_struct.LocalConversation) error
	BatchUpdateConversationList(ctx context.Context, conversationList []*model_struct.LocalConversation) error
	ConversationIfExists(ctx context.Context, conversationID string) (bool, error)
	ResetConversation(ctx context.Context, conversationID string) error
	ResetAllConversation(ctx context.Context) error
	ClearConversation(ctx context.Context, conversationID string) error
	SetConversationDraftDB(ctx context.Context, conversationID, draftText string) error
	RemoveConversationDraft(ctx context.Context, conversationID, draftText string) error
	UnPinConversation(ctx context.Context, conversationID string, isPinned int) error
	UpdateColumnsConversation(ctx context.Context, conversationID string, args map[string]interface{}) error
	UpdateAllConversation(ctx context.Context, conversation *model_struct.LocalConversation) error
	IncrConversationUnreadCount(ctx context.Context, conversationID string) error
	DecrConversationUnreadCount(ctx context.Context, conversationID string, count int64) (err error)
	GetTotalUnreadMsgCountDB(ctx context.Context) (totalUnreadCount int32, err error)
	SetMultipleConversationRecvMsgOpt(ctx context.Context, conversationIDList []string, opt int) (err error)
	GetMultipleConversationDB(ctx context.Context, conversationIDList []string) (result []*model_struct.LocalConversation, err error)
	SearchAllMessageByContentType(ctx context.Context, conversationID string, contentType int) ([]*model_struct.LocalChatLog, error)
	SuperGroupSearchAllMessageByContentType(ctx context.Context, superGroupID string, contentType int32) ([]*model_struct.LocalChatLog, error)
}

type UserDatabase interface {
	GetLoginUser(ctx context.Context, userID string) (*model_struct.LocalUser, error)
	UpdateLoginUser(ctx context.Context, user *model_struct.LocalUser) error
	UpdateLoginUserByMap(ctx context.Context, user *model_struct.LocalUser, args map[string]interface{}) error
	InsertLoginUser(ctx context.Context, user *model_struct.LocalUser) error
}

type FriendDatabase interface {
	InsertFriend(ctx context.Context, friend *model_struct.LocalFriend) error
	DeleteFriendDB(ctx context.Context, friendUserID string) error
	UpdateFriend(ctx context.Context, friend *model_struct.LocalFriend) error
	GetAllFriendList(ctx context.Context) ([]*model_struct.LocalFriend, error)
	GetPageFriendList(ctx context.Context, offset, count int) ([]*model_struct.LocalFriend, error)
	SearchFriendList(ctx context.Context, keyword string, isSearchUserID, isSearchNickname, isSearchRemark bool) ([]*model_struct.LocalFriend, error)
	GetFriendInfoByFriendUserID(ctx context.Context, FriendUserID string) (*model_struct.LocalFriend, error)
	GetFriendInfoList(ctx context.Context, friendUserIDList []string) ([]*model_struct.LocalFriend, error)
	InsertFriendRequest(ctx context.Context, friendRequest *model_struct.LocalFriendRequest) error
	DeleteFriendRequestBothUserID(ctx context.Context, fromUserID, toUserID string) error
	UpdateFriendRequest(ctx context.Context, friendRequest *model_struct.LocalFriendRequest) error
	GetRecvFriendApplication(ctx context.Context) ([]*model_struct.LocalFriendRequest, error)
	GetSendFriendApplication(ctx context.Context) ([]*model_struct.LocalFriendRequest, error)
	GetFriendApplicationByBothID(ctx context.Context, fromUserID, toUserID string) (*model_struct.LocalFriendRequest, error)
	GetBothFriendReq(ctx context.Context, fromUserID, toUserID string) ([]*model_struct.LocalFriendRequest, error)

	GetBlackListDB(ctx context.Context) ([]*model_struct.LocalBlack, error)
	GetBlackListUserID(ctx context.Context) (blackListUid []string, err error)
	GetBlackInfoByBlockUserID(ctx context.Context, blockUserID string) (*model_struct.LocalBlack, error)
	GetBlackInfoList(ctx context.Context, blockUserIDList []string) ([]*model_struct.LocalBlack, error)
	InsertBlack(ctx context.Context, black *model_struct.LocalBlack) error
	UpdateBlack(ctx context.Context, black *model_struct.LocalBlack) error
	DeleteBlack(ctx context.Context, blockUserID string) error
}

type ReactionDatabase interface {
	GetMessageReactionExtension(ctx context.Context, msgID string) (result *model_struct.LocalChatLogReactionExtensions, err error)
	InsertMessageReactionExtension(ctx context.Context, messageReactionExtension *model_struct.LocalChatLogReactionExtensions) error
	UpdateMessageReactionExtension(ctx context.Context, c *model_struct.LocalChatLogReactionExtensions) error
	// GetAndUpdateMessageReactionExtension(ctx context.Context, msgID string, m map[string]*sdkws.KeyValue) error
	// DeleteAndUpdateMessageReactionExtension(ctx context.Context, msgID string, m map[string]*sdkws.KeyValue) error
	GetMultipleMessageReactionExtension(ctx context.Context, msgIDList []string) (result []*model_struct.LocalChatLogReactionExtensions, err error)
	DeleteMessageReactionExtension(ctx context.Context, msgID string) error
}

type S3Database interface {
	GetUpload(ctx context.Context, partHash string) (*model_struct.LocalUpload, error)
	InsertUpload(ctx context.Context, upload *model_struct.LocalUpload) error
	DeleteUpload(ctx context.Context, partHash string) error
	UpdateUpload(ctx context.Context, upload *model_struct.LocalUpload) error
	DeleteExpireUpload(ctx context.Context) error
}

type DataBase interface {
	Close(ctx context.Context) error
	InitDB(ctx context.Context, userID string, dataDir string) error
	GroupDatabase
	MessageDatabase
	ConversationDatabase
	UserDatabase
	FriendDatabase
	ReactionDatabase
	S3Database
}
