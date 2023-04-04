package open_im_sdk

import (
	"bytes"
	ws "open_im_sdk/internal/interaction"
	common2 "open_im_sdk/internal/obj_storage"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
)

func GetAllConversationList(callback open_im_sdk_callback.Base, operationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().GetAllConversationList(callback, operationID)
}

func GetConversationListSplit(callback open_im_sdk_callback.Base, operationID string, offset, count int) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().GetConversationListSplit(callback, offset, count, operationID)
}

func GetOneConversation(callback open_im_sdk_callback.Base, operationID string, sessionType int, sourceID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().GetOneConversation(callback, int32(sessionType), sourceID, operationID)
}

func GetMultipleConversation(callback open_im_sdk_callback.Base, operationID string, conversationIDList string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().GetMultipleConversation(callback, conversationIDList, operationID)
}

func SetOneConversationPrivateChat(callback open_im_sdk_callback.Base, operationID, conversationID string, isPrivate bool) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().SetOneConversationPrivateChat(callback, conversationID, isPrivate, operationID)
}

func SetOneConversationBurnDuration(callback open_im_sdk_callback.Base, operationID, conversationID string, burnDuration int32) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().SetOneConversationBurnDuration(callback, conversationID, burnDuration, operationID)
}

func SetOneConversationRecvMessageOpt(callback open_im_sdk_callback.Base, operationID, conversationID string, opt int) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().SetOneConversationRecvMessageOpt(callback, conversationID, opt, operationID)
}

func SetConversationRecvMessageOpt(callback open_im_sdk_callback.Base, operationID string, conversationIDList string, opt int) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().SetConversationRecvMessageOpt(callback, conversationIDList, opt, operationID)
}
func SetGlobalRecvMessageOpt(callback open_im_sdk_callback.Base, operationID string, opt int) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().SetGlobalRecvMessageOpt(callback, opt, operationID)
}

func HideConversation(callback open_im_sdk_callback.Base, operationID, conversationID string) {
	BaseCaller(userForSDK.Conversation().HideConversation, callback, conversationID, operationID)
}

// deprecated
func GetConversationRecvMessageOpt(callback open_im_sdk_callback.Base, operationID string, conversationIDList string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().GetConversationRecvMessageOpt(callback, conversationIDList, operationID)
}

func DeleteConversation(callback open_im_sdk_callback.Base, operationID string, conversationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().DeleteConversation(callback, conversationID, operationID)
}
func DeleteAllConversationFromLocal(callback open_im_sdk_callback.Base, operationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().DeleteAllConversationFromLocal(callback, operationID)
}
func SetConversationDraft(callback open_im_sdk_callback.Base, operationID string, conversationID, draftText string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().SetConversationDraft(callback, conversationID, draftText, operationID)
}
func ResetConversationGroupAtType(callback open_im_sdk_callback.Base, operationID string, conversationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().ResetConversationGroupAtType(callback, conversationID, operationID)
}

func PinConversation(callback open_im_sdk_callback.Base, operationID string, conversationID string, isPinned bool) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().PinConversation(callback, conversationID, isPinned, operationID)
}

func GetTotalUnreadMsgCount(callback open_im_sdk_callback.Base, operationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().GetTotalUnreadMsgCount(callback, operationID)
}

func CreateAdvancedTextMessage(operationID string, text, messageEntityList string) string {
	return userForSDK.Conversation().CreateAdvancedTextMessage(text, messageEntityList, operationID)
}
func CreateTextAtMessage(operationID string, text, atUserList, atUsersInfo, message string) string {
	return userForSDK.Conversation().CreateTextAtMessage(text, atUserList, atUsersInfo, message, operationID)
}

func CreateTextMessage(operationID string, text string) string {
	return userForSDK.Conversation().CreateTextMessage(text, operationID)
}

func CreateLocationMessage(operationID string, description string, longitude, latitude float64) string {
	return userForSDK.Conversation().CreateLocationMessage(description, longitude, latitude, operationID)
}
func CreateCustomMessage(operationID string, data, extension string, description string) string {
	return userForSDK.Conversation().CreateCustomMessage(data, extension, description, operationID)
}
func CreateQuoteMessage(operationID string, text string, message string) string {
	return userForSDK.Conversation().CreateQuoteMessage(text, message, operationID)
}
func CreateAdvancedQuoteMessage(operationID string, text string, message, messageEntityList string) string {
	return userForSDK.Conversation().CreateAdvancedQuoteMessage(text, message, messageEntityList, operationID)
}
func CreateCardMessage(operationID string, cardInfo string) string {
	return userForSDK.Conversation().CreateCardMessage(cardInfo, operationID)

}
func CreateVideoMessageFromFullPath(operationID string, videoFullPath string, videoType string, duration int64, snapshotFullPath string) string {
	return userForSDK.Conversation().CreateVideoMessageFromFullPath(videoFullPath, videoType, duration, snapshotFullPath, operationID)
}
func CreateImageMessageFromFullPath(operationID string, imageFullPath string) string {
	return userForSDK.Conversation().CreateImageMessageFromFullPath(imageFullPath, operationID)
}
func CreateSoundMessageFromFullPath(operationID string, soundPath string, duration int64) string {
	return userForSDK.Conversation().CreateSoundMessageFromFullPath(soundPath, duration, operationID)
}
func CreateFileMessageFromFullPath(operationID string, fileFullPath, fileName string) string {
	return userForSDK.Conversation().CreateFileMessageFromFullPath(fileFullPath, fileName, operationID)
}
func CreateImageMessage(operationID string, imagePath string) string {
	return userForSDK.Conversation().CreateImageMessage(imagePath, operationID)
}
func CreateImageMessageByURL(operationID string, sourcePicture, bigPicture, snapshotPicture string) string {
	return userForSDK.Conversation().CreateImageMessageByURL(sourcePicture, bigPicture, snapshotPicture, operationID)
}

func CreateSoundMessageByURL(operationID string, soundBaseInfo string) string {
	return userForSDK.Conversation().CreateSoundMessageByURL(soundBaseInfo, operationID)
}
func CreateSoundMessage(operationID string, soundPath string, duration int64) string {
	return userForSDK.Conversation().CreateSoundMessage(soundPath, duration, operationID)
}
func CreateVideoMessageByURL(operationID string, videoBaseInfo string) string {
	return userForSDK.Conversation().CreateVideoMessageByURL(videoBaseInfo, operationID)
}
func CreateVideoMessage(operationID string, videoPath string, videoType string, duration int64, snapshotPath string) string {
	return userForSDK.Conversation().CreateVideoMessage(videoPath, videoType, duration, snapshotPath, operationID)
}
func CreateFileMessageByURL(operationID string, fileBaseInfo string) string {
	return userForSDK.Conversation().CreateFileMessageByURL(fileBaseInfo, operationID)
}
func CreateFileMessage(operationID string, filePath string, fileName string) string {
	return userForSDK.Conversation().CreateFileMessage(filePath, fileName, operationID)
}
func CreateMergerMessage(operationID string, messageList, title, summaryList string) string {
	return userForSDK.Conversation().CreateMergerMessage(messageList, title, summaryList, operationID)
}
func CreateFaceMessage(operationID string, index int, data string) string {
	return userForSDK.Conversation().CreateFaceMessage(index, data, operationID)
}
func CreateForwardMessage(operationID string, m string) string {
	return userForSDK.Conversation().CreateForwardMessage(m, operationID)
}

func SendMessage(callback open_im_sdk_callback.SendMsgCallBack, operationID, message, recvID, groupID, offlinePushInfo string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().SendMessage(callback, message, recvID, groupID, offlinePushInfo, operationID)
}
func SendMessageNotOss(callback open_im_sdk_callback.SendMsgCallBack, operationID string, message, recvID, groupID string, offlinePushInfo string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().SendMessageNotOss(callback, message, recvID, groupID, offlinePushInfo, operationID)
}
func SendMessageByBuffer(callback open_im_sdk_callback.SendMsgCallBack, operationID string, message, recvID, groupID string, offlinePushInfo string, buffer1, buffer2 *bytes.Buffer) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().SendMessageByBuffer(callback, message, recvID, groupID, offlinePushInfo, operationID, buffer1, buffer2)
}
func FindMessageList(callback open_im_sdk_callback.Base, operationID string, findMessageOptions string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().FindMessageList(callback, findMessageOptions, operationID)
}
func GetHistoryMessageList(callback open_im_sdk_callback.Base, operationID string, getMessageOptions string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().GetHistoryMessageList(callback, getMessageOptions, operationID)
}
func GetAdvancedHistoryMessageList(callback open_im_sdk_callback.Base, operationID string, getMessageOptions string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().GetAdvancedHistoryMessageList(callback, getMessageOptions, operationID)
}
func GetAdvancedHistoryMessageListReverse(callback open_im_sdk_callback.Base, operationID string, getMessageOptions string) {
	BaseCaller(userForSDK.Conversation().GetAdvancedHistoryMessageListReverse, callback, getMessageOptions, operationID)
}
func GetHistoryMessageListReverse(callback open_im_sdk_callback.Base, operationID string, getMessageOptions string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().GetHistoryMessageListReverse(callback, getMessageOptions, operationID)
}

// deprecated
func RevokeMessage(callback open_im_sdk_callback.Base, operationID string, message string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().RevokeMessage(callback, message, operationID)
}

func NewRevokeMessage(callback open_im_sdk_callback.Base, operationID string, message string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().NewRevokeMessage(callback, message, operationID)
}
func TypingStatusUpdate(callback open_im_sdk_callback.Base, operationID string, recvID, msgTip string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().TypingStatusUpdate(callback, recvID, msgTip, operationID)
}
func MarkC2CMessageAsRead(callback open_im_sdk_callback.Base, operationID string, userID string, msgIDList string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().MarkC2CMessageAsRead(callback, userID, msgIDList, operationID)
}
func MarkMessageAsReadByConID(callback open_im_sdk_callback.Base, operationID string, conversationID string, msgIDList string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().MarkMessageAsReadByConID(callback, conversationID, msgIDList, operationID)
}

// deprecated
func MarkGroupMessageHasRead(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().MarkGroupMessageHasRead(callback, groupID, operationID)
}
func MarkGroupMessageAsRead(callback open_im_sdk_callback.Base, operationID string, groupID, msgIDList string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().MarkGroupMessageAsRead(callback, groupID, msgIDList, operationID)
}

func DeleteMessageFromLocalStorage(callback open_im_sdk_callback.Base, operationID string, message string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().DeleteMessageFromLocalStorage(callback, message, operationID)
}

func DeleteMessageFromLocalAndSvr(callback open_im_sdk_callback.Base, operationID string, message string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().DeleteMessageFromLocalAndSvr(callback, message, operationID)
}

func DeleteConversationFromLocalAndSvr(callback open_im_sdk_callback.Base, operationID string, conversationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().DeleteConversationFromLocalAndSvr(callback, conversationID, operationID)
}

func DeleteAllMsgFromLocalAndSvr(callback open_im_sdk_callback.Base, operationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().DeleteAllMsgFromLocalAndSvr(callback, operationID)
}

func DeleteAllMsgFromLocal(callback open_im_sdk_callback.Base, operationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().DeleteAllMsgFromLocal(callback, operationID)
}

func ClearC2CHistoryMessage(callback open_im_sdk_callback.Base, operationID string, userID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().ClearC2CHistoryMessage(callback, userID, operationID)
}
func ClearC2CHistoryMessageFromLocalAndSvr(callback open_im_sdk_callback.Base, operationID string, userID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().ClearC2CHistoryMessageFromLocalAndSvr(callback, userID, operationID)
}

func ClearGroupHistoryMessage(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().ClearGroupHistoryMessage(callback, groupID, operationID)
}
func ClearGroupHistoryMessageFromLocalAndSvr(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().ClearGroupHistoryMessageFromLocalAndSvr(callback, groupID, operationID)
}
func InsertSingleMessageToLocalStorage(callback open_im_sdk_callback.Base, operationID string, message, recvID, sendID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().InsertSingleMessageToLocalStorage(callback, message, recvID, sendID, operationID)
}
func InsertGroupMessageToLocalStorage(callback open_im_sdk_callback.Base, operationID string, message, groupID, sendID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().InsertGroupMessageToLocalStorage(callback, message, groupID, sendID, operationID)
}
func SearchLocalMessages(callback open_im_sdk_callback.Base, operationID string, searchParam string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().SearchLocalMessages(callback, searchParam, operationID)
}

func uploadImage(callback open_im_sdk_callback.Base, operationID string, filePath string, token, obj string) string {
	if obj == "cos" {
		p := ws.NewPostApi(token, userForSDK.ImConfig().ApiAddr)
		o := common2.NewCOS(p)
		url, _, err := o.UploadFile(filePath, func(progress int) {
			if progress == 100 {
				callback.OnSuccess("")
			}
		})

		if err != nil {
			callback.OnError(100, err.Error())
			return ""
		}
		return url

	} else {
		return ""
	}
}
func GetConversationIDBySessionType(sourceID string, sessionType int) string {
	return utils.GetConversationIDBySessionType(sourceID, sessionType)
}
func GetAtAllTag() string {
	return constant.AtAllString
}

type MessageType struct {
	TypeKey     string `json:"typeKey"`
	CanRepeat   bool   `json:"canRepeat,omitempty"`
	NeedToCount bool   `json:"needToCount,omitempty"`
	Counter     int32  `json:"counter,omitempty"`
}

// 修改
func SetMessageReactionExtensions(callback open_im_sdk_callback.Base, operationID, message, reactionExtensionList string) {
	BaseCaller(userForSDK.Conversation().SetMessageReactionExtensions, callback, message, reactionExtensionList, operationID)
}

func AddMessageReactionExtensions(callback open_im_sdk_callback.Base, operationID, message, reactionExtensionList string) {
	BaseCaller(userForSDK.Conversation().AddMessageReactionExtensions, callback, message, reactionExtensionList, operationID)
}

func DeleteMessageReactionExtensions(callback open_im_sdk_callback.Base, operationID, message, reactionExtensionKeyList string) {
	BaseCaller(userForSDK.Conversation().DeleteMessageReactionExtensions, callback, message, reactionExtensionKeyList, operationID)
}

func GetMessageListReactionExtensions(callback open_im_sdk_callback.Base, operationID, messageList string) {
	BaseCaller(userForSDK.Conversation().GetMessageListReactionExtensions, callback, messageList, operationID)
}
func GetMessageListSomeReactionExtensions(callback open_im_sdk_callback.Base, operationID, messageList, reactionExtensionKeyList string) {
	BaseCaller(userForSDK.Conversation().GetMessageListSomeReactionExtensions, callback, messageList, reactionExtensionKeyList, operationID)
}

func SetTypeKeyInfo(callback open_im_sdk_callback.Base, operationID, message, typeKey, ex string, isCanRepeat bool) {
	BaseCaller(userForSDK.Conversation().SetTypeKeyInfo, callback, message, typeKey, ex, isCanRepeat, operationID)
}
func GetTypeKeyListInfo(callback open_im_sdk_callback.Base, operationID, message, typeKeyList string) {
	BaseCaller(userForSDK.Conversation().GetTypeKeyListInfo, callback, message, typeKeyList, operationID)
}
func GetAllTypeKeyInfo(callback open_im_sdk_callback.Base, message, operationID string) {
	BaseCaller(userForSDK.Conversation().GetAllTypeKeyInfo, callback, message, operationID)
}
