package open_im_sdk

import (
	"encoding/json"
	"fmt"
	"open_im_sdk/internal/controller/conversation_msg"
	"open_im_sdk/internal/controller/friend"
	"open_im_sdk/internal/controller/group"
	ws "open_im_sdk/internal/controller/interaction"
	"open_im_sdk/internal/controller/login"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

/*
var RouteMapSDK map[string]reflect.Value

func init(){
	RouteMapSDK = make(map[string]reflect.Value, 0)
	vf := reflect.ValueOf(&userForSDK)
	vft := vf.Type()
	mNum := vf.NumMethod()
	fmt.Println("vft num ", vft, mNum)
	for i := 0; i < mNum; i++ {
		mName := vft.Method(i).Name
		RouteMapSDK[mName] = vf.Method(i)
		fmt.Println("func ", vf.Method(i))
	}

}
*/
const sdkVersion = "Open-IM-SDK-Core-v2.0.0"

func SdkVersion() string {
	return sdkVersion
}

func InitSDK(config string, operationID string, listener ws.ConnListener) bool {
	if err := json.Unmarshal([]byte(config), &sdk_struct.SvrConf); err != nil {
		return false
	}
	fmt.Println("config ", config, sdk_struct.SvrConf)
	log.NewPrivateLog("", sdk_struct.SvrConf.LogLevel)
	log.NewInfo("0", utils.GetSelfFuncName(), config, SdkVersion())
	if listener == nil || config == "" {
		log.Error("0", "listener or config is nil")
		return false
	}
	if userForSDK != nil {
		log.Warn("0", "Initialize multiple times, call logout")
		userForSDK.Logout(nil, utils.OperationIDGenerator())
	}
	userForSDK = new(login.LoginMgr)
	return userForSDK.InitSDK(sdk_struct.SvrConf, listener)
}

func Login(userID, operationID string, token string, callback common.Base) {
	if callback == nil {
		log.Error("callback is nil")
		return
	}
	if userForSDK == nil {
		callback.OnError(constant.ErrArgs.ErrCode, constant.ErrArgs.ErrMsg)
		return
	}
	userForSDK.Login(callback, userID, token, operationID)
}

func Logout(callback common.Base, operationID string) {
	if callback == nil {
		log.Error("callback is nil")
		return
	}
	if userForSDK == nil {
		callback.OnError(constant.ErrArgs.ErrCode, constant.ErrArgs.ErrMsg)
		return
	}
	userForSDK.Logout(callback, operationID)
}

func GetLoginStatus() int32 {
	return userForSDK.GetLoginStatus()
}

func GetLoginUser() string {
	return userForSDK.GetLoginUser()
}

///////////////////////user/////////////////////
func GetUsersInfo(uIDList string, operationID string, callback common.Base) {
	userForSDK.User().GetUsersInfo(callback, uIDList, operationID)
}

func SetSelfInfo(callback common.Base, operationID string, info string) {
	userForSDK.User().SetSelfInfo(callback, info, operationID)
}

//////////////////////////group//////////////////////////////////////////
func SetGroupListener(callback group.OnGroupListener) {
	if callback == nil {
		log.Error("callback is nil")
		return
	}
	userForSDK.SetGroupListener(callback)
}

func CreateGroup(callback common.Base, operationID string, gInfo string, memberList string) {
	if callback == nil {
		log.Error("callback is nil")
		return
	}
	userForSDK.Group().CreateGroup(callback, gInfo, memberList, operationID)
}

func JoinGroup(callback common.Base, operationID string, groupID, message string) {
	if callback == nil {
		log.Error("callback is nil")
		return
	}
	userForSDK.Group().JoinGroup(callback, groupID, message, operationID)
}

func QuitGroup(callback common.Base, operationID string, groupID string) {
	userForSDK.Group().QuitGroup(callback, groupID, operationID)
}

func GetJoinedGroupList(callback common.Base, operationID string) {
	userForSDK.Group().GetJoinedGroupList(callback, operationID)
}

func GetGroupsInfo(callback common.Base, operationID string, groupIDList string) {
	userForSDK.Group().GetGroupsInfo(callback, groupIDList, operationID)
}

func SetGroupInfo(callback common.Base, operationID string, groupID string, jsonGroupInfo string) {
	userForSDK.Group().SetGroupInfo(callback, jsonGroupInfo, groupID, operationID)
}

func GetGroupMemberList(callback common.Base, operationID string, groupID string, filter int32, next int32) {
	userForSDK.Group().GetGroupMemberList(callback, groupID, filter, next, operationID)
}

func GetGroupMembersInfo(callback common.Base, operationID string, groupID string, userList string) {
	userForSDK.Group().GetGroupMembersInfo(callback, groupID, userList, operationID)
}

func KickGroupMember(callback common.Base, operationID string, groupID string, reason string, userList string) {
	userForSDK.Group().KickGroupMember(callback, groupID, reason, userList, operationID)
}

func TransferGroupOwner(callback common.Base, operationID string, groupID, userId string) {
	userForSDK.Group().TransferGroupOwner(callback, groupID, userId, operationID)
}

func InviteUserToGroup(callback common.Base, operationID string, groupID, reason string, userList string) {
	userForSDK.Group().InviteUserToGroup(callback, groupID, reason, userList, operationID)
}

func GetGroupApplicationList(callback common.Base, operationID string) {
	userForSDK.Group().GetGroupApplicationList(callback, operationID)
}

func AcceptGroupApplication(callback common.Base, groupID, fromUserID, handleMsg string, operationID string) {
	userForSDK.Group().AcceptGroupApplication(callback, groupID, fromUserID, handleMsg, operationID)
}

func RefuseGroupApplication(callback common.Base, operationID string, groupID, fromUserID, handleMsg string) {
	userForSDK.Group().RefuseGroupApplication(callback, groupID, fromUserID, handleMsg, operationID)
}

////////////////////////////friend/////////////////////////////////////

func GetDesignatedFriendsInfo(callback common.Base, operationID string, userIDList string) {
	userForSDK.Friend().GetDesignatedFriendsInfo(callback, userIDList, operationID)
}

func GetFriendList(callback common.Base, operationID string) {
	userForSDK.Friend().GetFriendList(callback, operationID)
}

func CheckFriend(callback common.Base, operationID string, userIDList string) {
	userForSDK.Friend().CheckFriend(callback, userIDList, operationID)
}

func AddFriend(callback common.Base, operationID string, paramsReq string) {
	userForSDK.Friend().AddFriend(callback, paramsReq, operationID)
}

func SetFriendRemark(callback common.Base, operationID string, params string) {
	userForSDK.Friend().SetFriendRemark(callback, params, operationID)
}
func DeleteFriend(callback common.Base, operationID string, friendUserID string) {
	userForSDK.Friend().DeleteFriend(callback, friendUserID, operationID)
}

func GetRecvFriendApplicationList(callback common.Base, operationID string) {
	userForSDK.Friend().GetRecvFriendApplicationList(callback, operationID)
}

func GetSendFriendApplicationList(callback common.Base, operationID string) {
	userForSDK.Friend().GetSendFriendApplicationList(callback, operationID)
}

func AcceptFriendApplication(callback common.Base, operationID string, params string) {
	userForSDK.Friend().AcceptFriendApplication(callback, params, operationID)
}

func RefuseFriendApplication(callback common.Base, operationID string, params string) {
	userForSDK.Friend().RefuseFriendApplication(callback, params, operationID)
}

func AddBlack(callback common.Base, operationID string, blackUserID string) {
	userForSDK.Friend().AddBlack(callback, blackUserID, operationID)
}

func GetBlackList(callback common.Base, operationID string) {
	userForSDK.Friend().GetBlackList(callback, operationID)
}

func RemoveBlack(callback common.Base, operationID string, removeUserID string) {
	userForSDK.Friend().RemoveBlack(callback, removeUserID, operationID)
}

func SetFriendListener(listener friend.OnFriendshipListener) {
	userForSDK.SetFriendListener(listener)
}

///////////////////////conversation////////////////////////////////////

//func GetAllConversationList(callback common.Base, operationID string) {
//	userForSDK.Conversation().GetAllConversationList(callback)
//}
//func GetConversationListSplit(callback common.Base, operationID string, offset, count int) {
//	userForSDK.Conversation().GetConversationListSplit(callback, offset, count)
//}
//func SetConversationRecvMessageOpt(callback common.Base, operationID string, conversationIDList string, opt int) {
//	userForSDK.Conversation().SetConversationRecvMessageOpt(callback, conversationIDList, opt)
//}
//
//func GetConversationRecvMessageOpt(callback common.Base, operationID string, conversationIDList string) {
//	userForSDK.Conversation().GetConversationRecvMessageOpt(callback, conversationIDList)
//}
//func GetOneConversation(operationID string, sourceID string, sessionType int, callback common.Base) {
//	userForSDK.Conversation().GetOneConversation(sourceID, sessionType, callback)
//}
//func GetMultipleConversation(operationID string, conversationIDList string, callback common.Base) {
//	userForSDK.Conversation().GetMultipleConversation(conversationIDList, callback)
//}
//func DeleteConversation(operationID string, conversationID string, callback common.Base) {
//	userForSDK.Conversation().DeleteConversation(conversationID, callback)
//}
//func SetConversationDraft(operationID string, conversationID, draftText string, callback common.Base) {
//	userForSDK.Conversation().SetConversationDraft(conversationID, draftText, callback)
//}
//func PinConversation(operationID string, conversationID string, isPinned bool, callback common.Base) {
//	userForSDK.Conversation().PinConversation(conversationID, isPinned, callback)
//}
//func GetTotalUnreadMsgCount(callback common.Base,operationID string) {
//	userForSDK.Conversation().GetTotalUnreadMsgCount(callback)
//}
//
func SetConversationListener(listener conversation_msg.OnConversationListener) {
	userForSDK.SetConversationListener(listener)
}

//
//func AddAdvancedMsgListener(listener conversation_msg.OnAdvancedMsgListener) {
//	userForSDK.Conversation().AddAdvancedMsgListener(listener)
//}
//
//func CreateTextMessage(operationID string, text string) string {
//	return userForSDK.Conversation().CreateTextMessage(text)
//}
//func CreateTextAtMessage(operationID string, text, atUserList string) string {
//	return userForSDK.Conversation().CreateTextAtMessage(text, atUserList)
//}
//func CreateLocationMessage(operationID string, description string, longitude, latitude float64) string {
//	return userForSDK.Conversation().CreateLocationMessage(description, longitude, latitude)
//}
//func CreateCustomMessage(operationID string, data, extension string, description string) string {
//	return userForSDK.Conversation().CreateCustomMessage(data, extension, description)
//}
//func CreateQuoteMessage(operationID string, text string, message string) string {
//	return userForSDK.Conversation().CreateQuoteMessage(text, message)
//}
//func CreateCardMessage(operationID string, cardInfo string) string {
//	return userForSDK.Conversation().CreateCardMessage(cardInfo)
//
//}
//func CreateVideoMessageFromFullPath(operationID string, videoFullPath string, videoType string, duration int64, snapshotFullPath string) string {
//	return userForSDK.Conversation().CreateVideoMessageFromFullPath(videoFullPath, videoType, duration, snapshotFullPath)
//}
//func CreateImageMessageFromFullPath(operationID string, imageFullPath string) string {
//	return userForSDK.Conversation().CreateImageMessageFromFullPath(imageFullPath)
//}
//func CreateSoundMessageFromFullPath(operationID string, soundPath string, duration int64) string {
//	return userForSDK.Conversation().CreateSoundMessageFromFullPath(soundPath, duration)
//}
//func CreateFileMessageFromFullPath(operationID string, fileFullPath, fileName string) string {
//	return userForSDK.Conversation().CreateFileMessageFromFullPath(fileFullPath, fileName)
//}
//func CreateImageMessage(operationID string, imagePath string) string {
//	return userForSDK.Conversation().CreateImageMessage(imagePath)
//}
//func CreateImageMessageByURL(operationID string, sourcePicture, bigPicture, snapshotPicture string) string {
//	return userForSDK.Conversation().CreateImageMessageByURL(sourcePicture, bigPicture, snapshotPicture)
//}
//func SendMessageNotOss(callback conversation_msg.SendMsgCallBack, operationID string, message, receiver, groupID string, onlineUserOnly bool, offlinePushInfo string) string {
//	return userForSDK.Conversation().SendMessageNotOss(callback, message, receiver, groupID, onlineUserOnly, offlinePushInfo)
//}
//func CreateSoundMessageByURL(operationID string, soundBaseInfo string) string {
//	return userForSDK.Conversation().CreateSoundMessageByURL(soundBaseInfo)
//}
//func CreateSoundMessage(operationID string, soundPath string, duration int64) string {
//	return userForSDK.Conversation().CreateSoundMessage(soundPath, duration)
//}
//func CreateVideoMessageByURL(operationID string, videoBaseInfo string) string {
//	return userForSDK.Conversation().CreateVideoMessageByURL(videoBaseInfo)
//}
//func CreateVideoMessage(operationID string, videoPath string, videoType string, duration int64, snapshotPath string) string {
//	return userForSDK.Conversation().CreateVideoMessage(videoPath, videoType, duration, snapshotPath)
//}
//func CreateFileMessageByURL(operationID string, fileBaseInfo string) string {
//	return userForSDK.Conversation().CreateFileMessageByURL(fileBaseInfo)
//}
//func CreateFileMessage(operationID string, filePath string, fileName string) string {
//	return userForSDK.Conversation().CreateFileMessage(filePath, fileName)
//}
//func CreateMergerMessage(operationID string, messageList, title, summaryList string) string {
//	return userForSDK.Conversation().CreateMergerMessage(messageList, title, summaryList)
//}
//
//func CreateForwardMessage(operationID string, m string,) string {
//	return userForSDK.Conversation().CreateForwardMessage(m)
//}
//
//func SendMessage(callback conversation_msg.SendMsgCallBack, operationID string, message, receiver, groupID string, onlineUserOnly bool, offlinePushInfo string) string {
//	return userForSDK.Conversation().SendMessage(callback, message, receiver, groupID, onlineUserOnly, offlinePushInfo)
//}
//func GetHistoryMessageList(callback common.Base, operationID string, getMessageOptions string) {
//	userForSDK.Conversation().GetHistoryMessageList(callback, getMessageOptions)
//}
//func RevokeMessage(callback common.Base, operationID string, message string) {
//	userForSDK.Conversation().RevokeMessage(callback, message)
//}
//func TypingStatusUpdate(operationID string, receiver, msgTip string) {
//	userForSDK.Conversation().TypingStatusUpdate(receiver, msgTip)
//}
//func MarkC2CMessageAsRead(callback common.Base, operationID string, receiver string, msgIDList string) {
//	userForSDK.Conversation().MarkC2CMessageAsRead(callback, receiver, msgIDList)
//}
//
////Deprecated
//func MarkSingleMessageHasRead(callback common.Base, operationID string, userID string) {
//	userForSDK.Conversation().MarkSingleMessageHasRead(callback, userID)
//}
//func MarkGroupMessageHasRead(callback common.Base, operationID string, groupID string) {
//	userForSDK.Conversation().MarkGroupMessageHasRead(callback, groupID)
//}
//func DeleteMessageFromLocalStorage(callback common.Base, operationID string, message string) {
//	userForSDK.Conversation().DeleteMessageFromLocalStorage(callback, message)
//}
//func ClearC2CHistoryMessage(callback common.Base, operationID string, userID string) {
//	userForSDK.Conversation().ClearC2CHistoryMessage(callback, userID)
//}
//func ClearGroupHistoryMessage(callback common.Base, operationID string, groupID string) {
//	userForSDK.Conversation().ClearGroupHistoryMessage(callback, groupID)
//}
//func InsertSingleMessageToLocalStorage(callback common.Base, operationID string, message, userID, sender string) string {
//	return userForSDK.Conversation().InsertSingleMessageToLocalStorage(callback, message, userID, sender)
//}
//
//func FindMessages(callback common.Base, operationID string, messageIDList string) {
//	userForSDK.Conversation().FindMessages(callback, messageIDList)
//}
