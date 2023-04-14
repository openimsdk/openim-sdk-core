package open_im_sdk

import (
	"bytes"
	"open_im_sdk/open_im_sdk_callback"
)

func GetAllConversationList(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, userForSDK.Conversation().GetAllConversationList)
}

func GetConversationListSplit(callback open_im_sdk_callback.Base, operationID string, offset int, count int) {
	call(callback, operationID, userForSDK.Conversation().GetConversationListSplit, offset, count)
}

func GetOneConversation(callback open_im_sdk_callback.Base, operationID string, sessionType int, sourceID string) {
	call(callback, operationID, userForSDK.Conversation().GetOneConversation, sessionType, sourceID)
}

func GetMultipleConversation(callback open_im_sdk_callback.Base, operationID string, conversationIDList string) {
	call(callback, operationID, userForSDK.Conversation().GetMultipleConversation, conversationIDList)
}

func SetConversationRecvMessageOpt(callback open_im_sdk_callback.Base, operationID string, conversationIDList string, opt int) {
	call(callback, operationID, userForSDK.Conversation().SetConversationRecvMessageOpt, conversationIDList, opt)
}

func SetGlobalRecvMessageOpt(callback open_im_sdk_callback.Base, operationID string, opt int) {
	call(callback, operationID, userForSDK.Conversation().SetGlobalRecvMessageOpt, opt)
}

func GetConversationRecvMessageOpt(callback open_im_sdk_callback.Base, operationID string, conversationIDList string) {
	call(callback, operationID, userForSDK.Conversation().GetConversationRecvMessageOpt, conversationIDList)
}

func DeleteConversation(callback open_im_sdk_callback.Base, operationID string, conversationID string) {
	call(callback, operationID, userForSDK.Conversation().DeleteConversation, conversationID)
}

func DeleteAllConversationFromLocal(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, userForSDK.Conversation().DeleteAllConversationFromLocal)
}

func SetConversationDraft(callback open_im_sdk_callback.Base, operationID string, conversationID string, draftText string) {
	call(callback, operationID, userForSDK.Conversation().SetConversationDraft, conversationID, draftText)
}

func ResetConversationGroupAtType(callback open_im_sdk_callback.Base, operationID string, conversationID string) {
	call(callback, operationID, userForSDK.Conversation().ResetConversationGroupAtType, conversationID)
}

func PinConversation(callback open_im_sdk_callback.Base, operationID string, conversationID string, isPinned bool) {
	call(callback, operationID, userForSDK.Conversation().PinConversation, conversationID, isPinned)
}

func GetTotalUnreadMsgCount(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, userForSDK.Conversation().GetTotalUnreadMsgCount)
}
func CreateAdvancedTextMessage(operationID string, text, messageEntityList string) string {
	return syncCall(operationID, userForSDK.Conversation().CreateAdvancedTextMessage, text, messageEntityList)
}
func CreateTextAtMessage(operationID string, text, atUserList, atUsersInfo, message string) string {
	return syncCall(operationID, userForSDK.Conversation().CreateTextAtMessage, text, atUserList, atUsersInfo, message)
}
func CreateTextMessage(operationID string, text string) string {
	return syncCall(operationID, userForSDK.Conversation().CreateTextMessage, text)
}

func CreateLocationMessage(operationID string, description string, longitude, latitude float64) string {
	return syncCall(operationID, userForSDK.Conversation().CreateLocationMessage, description, longitude, latitude)
}
func CreateCustomMessage(operationID string, data, extension string, description string) string {
	return syncCall(operationID, userForSDK.Conversation().CreateCustomMessage, data, extension, description)
}
func CreateQuoteMessage(operationID string, text string, message string) string {
	return syncCall(operationID, userForSDK.Conversation().CreateQuoteMessage, text, message)
}
func CreateAdvancedQuoteMessage(operationID string, text string, message, messageEntityList string) string {
	return syncCall(operationID, userForSDK.Conversation().CreateAdvancedQuoteMessage, text, message, messageEntityList)
}
func CreateCardMessage(operationID string, cardInfo string) string {
	return syncCall(operationID, userForSDK.Conversation().CreateCardMessage, cardInfo)

}
func CreateVideoMessageFromFullPath(operationID string, videoFullPath string, videoType string, duration int64, snapshotFullPath string) string {
	return syncCall(operationID, userForSDK.Conversation().CreateVideoMessageFromFullPath, videoFullPath, videoType, duration, snapshotFullPath)
}
func CreateImageMessageFromFullPath(operationID string, imageFullPath string) string {
	return syncCall(operationID, userForSDK.Conversation().CreateImageMessageFromFullPath, imageFullPath)
}
func CreateSoundMessageFromFullPath(operationID string, soundPath string, duration int64) string {
	return syncCall(operationID, userForSDK.Conversation().CreateSoundMessageFromFullPath, soundPath, duration)
}
func CreateFileMessageFromFullPath(operationID string, fileFullPath, fileName string) string {
	return syncCall(operationID, userForSDK.Conversation().CreateFileMessageFromFullPath, fileFullPath, fileName)
}
func CreateImageMessage(operationID string, imagePath string) string {
	return syncCall(operationID, userForSDK.Conversation().CreateImageMessage, imagePath)
}
func CreateImageMessageByURL(operationID string, sourcePicture, bigPicture, snapshotPicture string) string {
	return syncCall(operationID, userForSDK.Conversation().CreateImageMessageByURL, sourcePicture, bigPicture, snapshotPicture)
}

func CreateSoundMessageByURL(operationID string, soundBaseInfo string) string {
	return syncCall(operationID, userForSDK.Conversation().CreateSoundMessageByURL, soundBaseInfo)
}
func CreateSoundMessage(operationID string, soundPath string, duration int64) string {
	return syncCall(operationID, userForSDK.Conversation().CreateSoundMessage, soundPath, duration)
}
func CreateVideoMessageByURL(operationID string, videoBaseInfo string) string {
	return syncCall(operationID, userForSDK.Conversation().CreateVideoMessageByURL, videoBaseInfo)
}
func CreateVideoMessage(operationID string, videoPath string, videoType string, duration int64, snapshotPath string) string {
	return syncCall(operationID, userForSDK.Conversation().CreateVideoMessage, videoPath, videoType, duration, snapshotPath)
}
func CreateFileMessageByURL(operationID string, fileBaseInfo string) string {
	return syncCall(operationID, userForSDK.Conversation().CreateFileMessageByURL, fileBaseInfo)
}
func CreateFileMessage(operationID string, filePath string, fileName string) string {
	return syncCall(operationID, userForSDK.Conversation().CreateFileMessage, filePath, fileName)
}
func CreateMergerMessage(operationID string, messageList, title, summaryList string) string {
	return syncCall(operationID, userForSDK.Conversation().CreateMergerMessage, messageList, title, summaryList)
}
func CreateFaceMessage(operationID string, index int, data string) string {
	return syncCall(operationID, userForSDK.Conversation().CreateFaceMessage, index, data)
}
func CreateForwardMessage(operationID string, m string) string {
	return syncCall(operationID, userForSDK.Conversation().CreateForwardMessage, m)
}
func SendMessage(callback open_im_sdk_callback.SendMsgCallBack, operationID, message, recvID, groupID, offlinePushInfo string) {
	messageCall(callback, operationID, userForSDK.Conversation().SendMessage, message, recvID, groupID, offlinePushInfo)
}
func SendMessageNotOss(callback open_im_sdk_callback.SendMsgCallBack, operationID string, message, recvID, groupID string, offlinePushInfo string) {
	messageCall(callback, operationID, userForSDK.Conversation().SendMessageNotOss, message, recvID, groupID, offlinePushInfo)
}
func SendMessageByBuffer(callback open_im_sdk_callback.SendMsgCallBack, operationID string, message, recvID, groupID string, offlinePushInfo string, buffer1, buffer2 *bytes.Buffer) {
	messageCall(callback, operationID, userForSDK.Conversation().SendMessageByBuffer, message, recvID, groupID, offlinePushInfo, buffer1, buffer2)
}

func FindMessageList(callback open_im_sdk_callback.Base, operationID string, findMessageOptions string) {
	call(callback, operationID, userForSDK.Conversation().FindMessageList, findMessageOptions)
}

func GetHistoryMessageList(callback open_im_sdk_callback.Base, operationID string, getMessageOptions string) {
	call(callback, operationID, userForSDK.Conversation().GetHistoryMessageList, getMessageOptions)
}

func GetAdvancedHistoryMessageList(callback open_im_sdk_callback.Base, operationID string, getMessageOptions string) {
	call(callback, operationID, userForSDK.Conversation().GetAdvancedHistoryMessageList, getMessageOptions)
}

func GetAdvancedHistoryMessageListReverse(callback open_im_sdk_callback.Base, operationID string, getMessageOptions string) {
	call(callback, operationID, userForSDK.Conversation().GetAdvancedHistoryMessageListReverse, getMessageOptions)
}

func GetHistoryMessageListReverse(callback open_im_sdk_callback.Base, operationID string, getMessageOptions string) {
	call(callback, operationID, userForSDK.Conversation().GetHistoryMessageListReverse, getMessageOptions)
}

func RevokeMessage(callback open_im_sdk_callback.Base, operationID string, message string) {
	call(callback, operationID, userForSDK.Conversation().RevokeMessage, message)
}

func NewRevokeMessage(callback open_im_sdk_callback.Base, operationID string, message string) {
	call(callback, operationID, userForSDK.Conversation().NewRevokeMessage, message)
}

func TypingStatusUpdate(callback open_im_sdk_callback.Base, operationID string, recvID string, msgTip string) {
	call(callback, operationID, userForSDK.Conversation().TypingStatusUpdate, recvID, msgTip)
}

func MarkC2CMessageAsRead(callback open_im_sdk_callback.Base, operationID string, userID string, msgIDList string) {
	call(callback, operationID, userForSDK.Conversation().MarkC2CMessageAsRead, userID, msgIDList)
}

func MarkMessageAsReadByConID(callback open_im_sdk_callback.Base, operationID string, conversationID string, msgIDList string) {
	call(callback, operationID, userForSDK.Conversation().MarkMessageAsReadByConID, conversationID, msgIDList)
}

func MarkGroupMessageHasRead(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, userForSDK.Conversation().MarkGroupMessageHasRead, groupID)
}

func MarkGroupMessageAsRead(callback open_im_sdk_callback.Base, operationID string, groupID string, msgIDList string) {
	call(callback, operationID, userForSDK.Conversation().MarkGroupMessageAsRead, groupID, msgIDList)
}

func DeleteMessageFromLocalStorage(callback open_im_sdk_callback.Base, operationID string, message string) {
	call(callback, operationID, userForSDK.Conversation().DeleteMessageFromLocalStorage, message)
}

func DeleteMessageFromLocalAndSvr(callback open_im_sdk_callback.Base, operationID string, message string) {
	call(callback, operationID, userForSDK.Conversation().DeleteMessageFromLocalAndSvr, message)
}

func DeleteConversationFromLocalAndSvr(callback open_im_sdk_callback.Base, operationID string, conversationID string) {
	call(callback, operationID, userForSDK.Conversation().DeleteConversationFromLocalAndSvr, conversationID)
}

func DeleteAllMsgFromLocalAndSvr(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, userForSDK.Conversation().DeleteAllMsgFromLocalAndSvr)
}

func DeleteAllMsgFromLocal(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, userForSDK.Conversation().DeleteAllMsgFromLocal)
}

func ClearC2CHistoryMessage(callback open_im_sdk_callback.Base, operationID string, userID string) {
	call(callback, operationID, userForSDK.Conversation().ClearC2CHistoryMessage, userID)
}

func ClearC2CHistoryMessageFromLocalAndSvr(callback open_im_sdk_callback.Base, operationID string, userID string) {
	call(callback, operationID, userForSDK.Conversation().ClearC2CHistoryMessageFromLocalAndSvr, userID)
}

func ClearGroupHistoryMessage(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, userForSDK.Conversation().ClearGroupHistoryMessage, groupID)
}

func ClearGroupHistoryMessageFromLocalAndSvr(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	call(callback, operationID, userForSDK.Conversation().ClearGroupHistoryMessageFromLocalAndSvr, groupID)
}

func InsertSingleMessageToLocalStorage(callback open_im_sdk_callback.Base, operationID string, message string, recvID string, sendID string) {
	call(callback, operationID, userForSDK.Conversation().InsertSingleMessageToLocalStorage, message, recvID, sendID)
}

func InsertGroupMessageToLocalStorage(callback open_im_sdk_callback.Base, operationID string, message string, groupID string, sendID string) {
	call(callback, operationID, userForSDK.Conversation().InsertGroupMessageToLocalStorage, message, groupID, sendID)
}

func SearchLocalMessages(callback open_im_sdk_callback.Base, operationID string, searchParam string) {
	call(callback, operationID, userForSDK.Conversation().SearchLocalMessages, searchParam)
}
