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
	"bytes"
	"open_im_sdk/open_im_sdk_callback"
)

func GetAllConversationList(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Conversation().GetAllConversationList)
}

func GetConversationListSplit(callback open_im_sdk_callback.Base, operationID string, offset int, count int) {
	call(callback, operationID, UserForSDK.Conversation().GetConversationListSplit, offset, count)
}

func GetOneConversation(callback open_im_sdk_callback.Base, operationID string, sessionType int, sourceID string) {
	call(callback, operationID, UserForSDK.Conversation().GetOneConversation, sessionType, sourceID)
}

func GetMultipleConversation(callback open_im_sdk_callback.Base, operationID string, conversationIDList string) {
	call(callback, operationID, UserForSDK.Conversation().GetMultipleConversation, conversationIDList)
}

func SetConversationRecvMessageOpt(callback open_im_sdk_callback.Base, operationID string, conversationIDList string, opt int) {
	call(callback, operationID, UserForSDK.Conversation().SetConversationRecvMessageOpt, conversationIDList, opt)
}

func SetGlobalRecvMessageOpt(callback open_im_sdk_callback.Base, operationID string, opt int) {
	call(callback, operationID, UserForSDK.Conversation().SetGlobalRecvMessageOpt, opt)
}

func GetConversationRecvMessageOpt(callback open_im_sdk_callback.Base, operationID string, conversationIDList string) {
	call(callback, operationID, UserForSDK.Conversation().GetConversationRecvMessageOpt, conversationIDList)
}

func DeleteConversation(callback open_im_sdk_callback.Base, operationID string, conversationID string) {
	call(callback, operationID, UserForSDK.Conversation().DeleteConversation, conversationID)
}

func DeleteAllConversationFromLocal(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Conversation().DeleteAllConversationFromLocal)
}

func SetConversationDraft(callback open_im_sdk_callback.Base, operationID string, conversationID string, draftText string) {
	call(callback, operationID, UserForSDK.Conversation().SetConversationDraft, conversationID, draftText)
}

func ResetConversationGroupAtType(callback open_im_sdk_callback.Base, operationID string, conversationID string) {
	call(callback, operationID, UserForSDK.Conversation().ResetConversationGroupAtType, conversationID)
}

func PinConversation(callback open_im_sdk_callback.Base, operationID string, conversationID string, isPinned bool) {
	call(callback, operationID, UserForSDK.Conversation().PinConversation, conversationID, isPinned)
}

func SetOneConversationPrivateChat(callback open_im_sdk_callback.Base, operationID string, conversationID string, isPrivate bool) {
	call(callback, operationID, UserForSDK.Conversation().SetOneConversationPrivateChat, conversationID, isPrivate)
}

func SetOneConversationBurnDuration(callback open_im_sdk_callback.Base, operationID string, conversationID string, duration int) {
	call(callback, operationID, UserForSDK.Conversation().SetOneConversationBurnDuration, conversationID, duration)
}

func SetOneConversationRecvMessageOpt(callback open_im_sdk_callback.Base, operationID string, conversationID string, opt int) {
	call(callback, operationID, UserForSDK.Conversation().SetOneConversationRecvMessageOpt, conversationID, opt)
}

func GetTotalUnreadMsgCount(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Conversation().GetTotalUnreadMsgCount)
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
func CreateVideoMessageFromFullPath(operationID string, videoFullPath string, videoType string, duration int64, snapshotFullPath string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateVideoMessageFromFullPath, videoFullPath, videoType, duration, snapshotFullPath)
}
func CreateImageMessageFromFullPath(operationID string, imageFullPath string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateImageMessageFromFullPath, imageFullPath)
}
func CreateSoundMessageFromFullPath(operationID string, soundPath string, duration int64) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateSoundMessageFromFullPath, soundPath, duration)
}
func CreateFileMessageFromFullPath(operationID string, fileFullPath, fileName string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateFileMessageFromFullPath, fileFullPath, fileName)
}
func CreateImageMessage(operationID string, imagePath string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateImageMessage, imagePath)
}
func CreateImageMessageByURL(operationID string, sourcePicture, bigPicture, snapshotPicture string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateImageMessageByURL, sourcePicture, bigPicture, snapshotPicture)
}

func CreateSoundMessageByURL(operationID string, soundBaseInfo string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateSoundMessageByURL, soundBaseInfo)
}
func CreateSoundMessage(operationID string, soundPath string, duration int64) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateSoundMessage, soundPath, duration)
}
func CreateVideoMessageByURL(operationID string, videoBaseInfo string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateVideoMessageByURL, videoBaseInfo)
}
func CreateVideoMessage(operationID string, videoPath string, videoType string, duration int64, snapshotPath string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateVideoMessage, videoPath, videoType, duration, snapshotPath)
}
func CreateFileMessageByURL(operationID string, fileBaseInfo string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateFileMessageByURL, fileBaseInfo)
}
func CreateFileMessage(operationID string, filePath string, fileName string) string {
	return syncCall(operationID, UserForSDK.Conversation().CreateFileMessage, filePath, fileName)
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
func SendMessage(callback open_im_sdk_callback.SendMsgCallBack, operationID, message, recvID, groupID, offlinePushInfo string) {
	messageCall(callback, operationID, UserForSDK.Conversation().SendMessage, message, recvID, groupID, offlinePushInfo)
}
func SendMessageNotOss(callback open_im_sdk_callback.SendMsgCallBack, operationID string, message, recvID, groupID string, offlinePushInfo string) {
	messageCall(callback, operationID, UserForSDK.Conversation().SendMessageNotOss, message, recvID, groupID, offlinePushInfo)
}
func SendMessageByBuffer(callback open_im_sdk_callback.SendMsgCallBack, operationID string, message, recvID, groupID string, offlinePushInfo string, buffer1, buffer2 *bytes.Buffer) {
	messageCall(callback, operationID, UserForSDK.Conversation().SendMessageByBuffer, message, recvID, groupID, offlinePushInfo, buffer1, buffer2)
}

func FindMessageList(callback open_im_sdk_callback.Base, operationID string, findMessageOptions string) {
	call(callback, operationID, UserForSDK.Conversation().FindMessageList, findMessageOptions)
}

func GetHistoryMessageList(callback open_im_sdk_callback.Base, operationID string, getMessageOptions string) {
	call(callback, operationID, UserForSDK.Conversation().GetHistoryMessageList, getMessageOptions)
}

func GetAdvancedHistoryMessageList(callback open_im_sdk_callback.Base, operationID string, getMessageOptions string) {
	call(callback, operationID, UserForSDK.Conversation().GetAdvancedHistoryMessageList, getMessageOptions)
}

func GetAdvancedHistoryMessageListReverse(callback open_im_sdk_callback.Base, operationID string, getMessageOptions string) {
	call(callback, operationID, UserForSDK.Conversation().GetAdvancedHistoryMessageListReverse, getMessageOptions)
}

func GetHistoryMessageListReverse(callback open_im_sdk_callback.Base, operationID string, getMessageOptions string) {
	call(callback, operationID, UserForSDK.Conversation().GetHistoryMessageListReverse, getMessageOptions)
}

func RevokeMessage(callback open_im_sdk_callback.Base, operationID string, message string) {
	call(callback, operationID, UserForSDK.Conversation().RevokeMessage, message)
}

func NewRevokeMessage(callback open_im_sdk_callback.Base, operationID string, message string) {
	call(callback, operationID, UserForSDK.Conversation().NewRevokeMessage, message)
}

func TypingStatusUpdate(callback open_im_sdk_callback.Base, operationID string, recvID string, msgTip string) {
	call(callback, operationID, UserForSDK.Conversation().TypingStatusUpdate, recvID, msgTip)
}

func MarkC2CMessageAsRead(callback open_im_sdk_callback.Base, operationID string, userID string, msgIDList string) {
	call(callback, operationID, UserForSDK.Conversation().MarkC2CMessageAsRead, userID, msgIDList)
}

func MarkMessageAsReadByConID(callback open_im_sdk_callback.Base, operationID string, conversationID string, msgIDList string) {
	call(callback, operationID, UserForSDK.Conversation().MarkMessageAsReadByConID, conversationID, msgIDList)
}

func MarkGroupMessageHasRead(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, UserForSDK.Conversation().MarkGroupMessageHasRead, groupID)
}

func MarkGroupMessageAsRead(callback open_im_sdk_callback.Base, operationID string, groupID string, msgIDList string) {
	call(callback, operationID, UserForSDK.Conversation().MarkGroupMessageAsRead, groupID, msgIDList)
}

func DeleteMessageFromLocalStorage(callback open_im_sdk_callback.Base, operationID string, message string) {
	call(callback, operationID, UserForSDK.Conversation().DeleteMessageFromLocalStorage, message)
}

func DeleteMessageFromLocalAndSvr(callback open_im_sdk_callback.Base, operationID string, message string) {
	call(callback, operationID, UserForSDK.Conversation().DeleteMessageFromLocalAndSvr, message)
}

func DeleteConversationFromLocalAndSvr(callback open_im_sdk_callback.Base, operationID string, conversationID string) {
	call(callback, operationID, UserForSDK.Conversation().DeleteConversationFromLocalAndSvr, conversationID)
}

func DeleteAllMsgFromLocalAndSvr(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Conversation().DeleteAllMsgFromLocalAndSvr)
}

func DeleteAllMsgFromLocal(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Conversation().DeleteAllMsgFromLocal)
}

func ClearC2CHistoryMessage(callback open_im_sdk_callback.Base, operationID string, userID string) {
	call(callback, operationID, UserForSDK.Conversation().ClearC2CHistoryMessage, userID)
}

func ClearC2CHistoryMessageFromLocalAndSvr(callback open_im_sdk_callback.Base, operationID string, userID string) {
	call(callback, operationID, UserForSDK.Conversation().ClearC2CHistoryMessageFromLocalAndSvr, userID)
}

func ClearGroupHistoryMessage(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, UserForSDK.Conversation().ClearGroupHistoryMessage, groupID)
}

func ClearGroupHistoryMessageFromLocalAndSvr(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, UserForSDK.Conversation().ClearGroupHistoryMessageFromLocalAndSvr, groupID)
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