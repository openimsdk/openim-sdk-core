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

package open_im_sdk

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
)

func GetAllConversationList(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Conversation().GetAllConversationList)
}

func GetConversationListSplit(callback open_im_sdk_callback.Base, operationID string, offset int, count int) {
	call(callback, operationID, UserForSDK.Conversation().GetConversationListSplit, offset, count)
}

func GetOneConversation(callback open_im_sdk_callback.Base, operationID string, sessionType int32, sourceID string) {
	call(callback, operationID, UserForSDK.Conversation().GetOneConversation, sessionType, sourceID)
}

func GetMultipleConversation(callback open_im_sdk_callback.Base, operationID string, conversationIDList string) {
	call(callback, operationID, UserForSDK.Conversation().GetMultipleConversation, conversationIDList)
}

func SetConversation(callback open_im_sdk_callback.Base, operationID string, conversationID string, req string) {
	call(callback, operationID, UserForSDK.Conversation().SetConversation, conversationID, req)
}

func HideConversation(callback open_im_sdk_callback.Base, operationID string, conversationID string) {
	call(callback, operationID, UserForSDK.Conversation().HideConversation, conversationID)
}

func SetConversationDraft(callback open_im_sdk_callback.Base, operationID string, conversationID string, draftText string) {
	call(callback, operationID, UserForSDK.Conversation().SetConversationDraft, conversationID, draftText)
}

func GetTotalUnreadMsgCount(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Conversation().GetTotalUnreadMsgCount)
}
func GetAtAllTag(operationID string) string {
	return syncCall(operationID, UserForSDK.Conversation().GetAtAllTag)

}
func CreateAdvancedTextMessage(operationID string, text, messageEntityList string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateAdvancedTextMessage, text, messageEntityList)
}
func CreateTextAtMessage(operationID string, text, atUserList, atUsersInfo, message string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateTextAtMessage, text, atUserList, atUsersInfo, message)
}
func CreateTextMessage(operationID string, text string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateTextMessage, text)
}

func CreateLocationMessage(operationID string, description string, longitude, latitude float64) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateLocationMessage, description, longitude, latitude)
}
func CreateCustomMessage(operationID string, data, extension string, description string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateCustomMessage, data, extension, description)
}
func CreateQuoteMessage(operationID string, text string, message string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateQuoteMessage, text, message)
}
func CreateAdvancedQuoteMessage(operationID string, text string, message, messageEntityList string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateAdvancedQuoteMessage, text, message, messageEntityList)
}
func CreateCardMessage(operationID string, cardInfo string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateCardMessage, cardInfo)

}
func CreateImageMessage(operationID string, imageSourcePath string, sourcePicture, bigPicture, snapshotPicture string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateImageMessage, imageSourcePath, sourcePicture, bigPicture, snapshotPicture)
}
func CreateSoundMessage(operationID string, soundPath string, duration int64, soundBaseInfo string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateSoundMessage, soundPath, duration, soundBaseInfo)
}
func CreateVideoMessage(operationID string, videoSourcePath string, videoType string, duration int64, snapshotSourcePath string, videoBaseInfo string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateVideoMessage, videoSourcePath, videoType, duration, snapshotSourcePath, videoBaseInfo)
}
func CreateFileMessage(operationID string, fileSourcePath string, fileName string, fileBaseInfo string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateFileMessage, fileSourcePath, fileName, fileBaseInfo)
}
func CreateMergerMessage(operationID string, messageList, title, summaryList string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateMergerMessage, messageList, title, summaryList)
}
func CreateFaceMessage(operationID string, index int, data string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateFaceMessage, index, data)
}
func CreateForwardMessage(operationID string, m string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateForwardMessage, m)
}
func GetConversationIDBySessionType(operationID string, sourceID string, sessionType int) string {
	return UserForSDK.Conversation().GetConversationIDBySessionType(context.Background(), sourceID, sessionType)
}
func SendMessage(callback open_im_sdk_callback.SendMsgCallBack, operationID, message, recvID, groupID, offlinePushInfo string, isOnlineOnly bool) {
	messageCall(callback, operationID, UserForSDK.Conversation().SendMessage, message, recvID, groupID, offlinePushInfo, isOnlineOnly)
}

func FindMessageList(callback open_im_sdk_callback.Base, operationID string, findMessageOptions string) {
	call(callback, operationID, UserForSDK.Conversation().FindMessageList, findMessageOptions)
}

func GetAdvancedHistoryMessageList(callback open_im_sdk_callback.Base, operationID string, getMessageOptions string) {
	call(callback, operationID, UserForSDK.Conversation().GetAdvancedHistoryMessageList, getMessageOptions)
}

func GetAdvancedHistoryMessageListReverse(callback open_im_sdk_callback.Base, operationID string, getMessageOptions string) {
	call(callback, operationID, UserForSDK.Conversation().GetAdvancedHistoryMessageListReverse, getMessageOptions)
}

func RevokeMessage(callback open_im_sdk_callback.Base, operationID string, conversationID, clientMsgID string) {
	call(callback, operationID, UserForSDK.Conversation().RevokeMessage, conversationID, clientMsgID)
}

func TypingStatusUpdate(callback open_im_sdk_callback.Base, operationID string, recvID string, msgTip string) {
	call(callback, operationID, UserForSDK.Conversation().TypingStatusUpdate, recvID, msgTip)
}

// mark as read
func MarkConversationMessageAsRead(callback open_im_sdk_callback.Base, operationID string, conversationID string) {
	call(callback, operationID, UserForSDK.Conversation().MarkConversationMessageAsRead, conversationID)
}

func MarkAllConversationMessageAsRead(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Conversation().MarkAllConversationMessageAsRead)
}

func MarkMessagesAsReadByMsgID(callback open_im_sdk_callback.Base, operationID string, conversationID string, clientMsgIDs string) {
	call(callback, operationID, UserForSDK.Conversation().MarkMessagesAsReadByMsgID, conversationID, clientMsgIDs)
}

func DeleteMessageFromLocalStorage(callback open_im_sdk_callback.Base, operationID string, conversationID string, clientMsgID string) {
	call(callback, operationID, UserForSDK.Conversation().DeleteMessageFromLocalStorage, conversationID, clientMsgID)
}

func DeleteMessage(callback open_im_sdk_callback.Base, operationID string, conversationID string, clientMsgID string) {
	call(callback, operationID, UserForSDK.Conversation().DeleteMessage, conversationID, clientMsgID)
}

func HideAllConversations(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Conversation().HideAllConversations)
}

func DeleteAllMsgFromLocalAndSvr(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Conversation().DeleteAllMsgFromLocalAndServer)
}

func DeleteAllMsgFromLocal(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Conversation().DeleteAllMessageFromLocalStorage)
}

func ClearConversationAndDeleteAllMsg(callback open_im_sdk_callback.Base, operationID string, conversationID string) {
	call(callback, operationID, UserForSDK.Conversation().ClearConversationAndDeleteAllMsg, conversationID)
}

func DeleteConversationAndDeleteAllMsg(callback open_im_sdk_callback.Base, operationID string, conversationID string) {
	call(callback, operationID, UserForSDK.Conversation().DeleteConversationAndDeleteAllMsg, conversationID)
}

func InsertSingleMessageToLocalStorage(callback open_im_sdk_callback.Base, operationID string, message string, recvID string, sendID string) {
	call(callback, operationID, UserForSDK.Conversation().InsertSingleMessageToLocalStorage, message, recvID, sendID)
}

func InsertGroupMessageToLocalStorage(callback open_im_sdk_callback.Base, operationID string, message string, groupID string, sendID string) {
	call(callback, operationID, UserForSDK.Conversation().InsertGroupMessageToLocalStorage, message, groupID, sendID)
}

func SearchLocalMessages(callback open_im_sdk_callback.Base, operationID string, searchParam string) {
	call(callback, operationID, UserForSDK.Conversation().SearchLocalMessages, searchParam)
}
func SetMessageLocalEx(callback open_im_sdk_callback.Base, operationID string, conversationID, clientMsgID, localEx string) {
	call(callback, operationID, UserForSDK.Conversation().SetMessageLocalEx, conversationID, clientMsgID, localEx)
}

func SearchConversation(callback open_im_sdk_callback.Base, operationID string, searchParam string) {
	call(callback, operationID, UserForSDK.Conversation().SearchConversation, searchParam)
}

func ChangeInputStates(callback open_im_sdk_callback.Base, operationID string, conversationID string, focus bool) {
	call(callback, operationID, UserForSDK.Conversation().ChangeInputStates, conversationID, focus)
}

func GetInputStates(callback open_im_sdk_callback.Base, operationID string, conversationID string, userID string) {
	call(callback, operationID, UserForSDK.Conversation().GetInputStates, conversationID, userID)
}
