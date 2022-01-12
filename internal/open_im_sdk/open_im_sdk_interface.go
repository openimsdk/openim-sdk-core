package open_im_sdk

import (
	"encoding/json"
	"open_im_sdk/internal/controller/conversation_msg"
	"open_im_sdk/internal/controller/friend"
	"open_im_sdk/internal/controller/group"
	"open_im_sdk/internal/controller/init"
	ws "open_im_sdk/internal/controller/interaction"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
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

func SdkVersion() string {
	return "Open-IM-SDK-Core-v1.0.6"
}

func InitSDK(config string, listener ws.ConnListener) bool {
	log.NewInfo("0", utils.GetSelfFuncName(), config)
	var sc utils.IMConfig
	if err := json.Unmarshal([]byte(config), &sc); err != nil {
		log.Error("initSDK failed, config: ", sc, err.Error())
		return false
	}
	log.NewInfo("0","InitSDK, config ", config, "version: ", SdkVersion())
	if userForSDK != nil {
		log.Error("0", "Logout first ")
		userForSDK.Logout(nil)
		userForSDK.UnInitSDK()
	}
	userForSDK = new(init.LoginMgr)



	if !userForSDK.InitSDK(config, listener){
		log.Error("0", "InitSDK failed ", config, listener)
	}
}

//1 no print
func SetSdkLog(flag int32) {
	constant.SdkLogFlag = flag
}
func (u *open_im_sdk.UserRelated) SetSdkLog(flag int32) {
	constant.SdkLogFlag = flag
}
func SetHearbeatInterval(interval int32) {
	open_im_sdk.hearbeatInterval = interval
}

func UnInitSDK() {
	if open_im_sdk.userForSDK == nil {
		utils.sdkLog("userForSDK nil")
		return
	}
	open_im_sdk.userForSDK.unInitSDK()
}

func Login(uid, tk string, callback init.Base) {
	if open_im_sdk.userForSDK == nil {
		utils.sdkLog("userForSDK nil")
		if callback != nil {
			callback.OnError(constant.ErrCodeInitLogin, "userForSDK nil, initSdk first ")
		}
		return
	}
	open_im_sdk.userForSDK.Login(uid, tk, callback)
}

func Logout(callback init.Base) {
	if open_im_sdk.userForSDK == nil {
		utils.sdkLog("userForSDK nil")
		if callback != nil {
			callback.OnError(constant.ErrCodeInitLogin, "userForSDK nil, initSdk first ")
		}
		return
	}
	open_im_sdk.userForSDK.logout(callback)
}

func GetLoginStatus() int {
	return open_im_sdk.userForSDK.getLoginStatus()
}

func GetLoginUser() string {
	return open_im_sdk.userForSDK.GetLoginUser()
}

func ForceSyncLoginUerInfo() {
	open_im_sdk.userForSDK.ForceSyncLoginUserInfo()
}

func ForceSyncMsg() bool {
	return open_im_sdk.userForSDK.ForceSyncMsg()
}

////////////////////////////////////////////////////////////////////
func SetGroupListener(callback group.OnGroupListener) {
	userForSDK.SetGroupListener(callback)
}

func CreateGroup(callback init.Base, gInfo string, memberList string, operationID string) {
	open_im_sdk.userForSDK.CreateGroup(callback, gInfo, memberList, operationID)
}

func JoinGroup(groupId, message string, callback init.Base) {
	open_im_sdk.userForSDK.JoinGroup(groupId, message, callback)
}

func QuitGroup(groupId string, callback init.Base) {
	open_im_sdk.userForSDK.QuitGroup(groupId, callback)
}

func GetJoinedGroupList(callback init.Base) {
	open_im_sdk.userForSDK.GetJoinedGroupList(callback)
}

func GetGroupsInfo(groupIdList string, callback init.Base) {
	open_im_sdk.userForSDK.GetGroupsInfo(groupIdList, callback)
}

func SetGroupInfo(jsonGroupInfo string, callback init.Base) {
	open_im_sdk.userForSDK.SetGroupInfo(jsonGroupInfo, callback)
}

func GetGroupMemberList(groupId string, filter int32, next int32, callback init.Base) {
	open_im_sdk.userForSDK.GetGroupMemberList(groupId, filter, next, callback)
}

func GetGroupMembersInfo(groupId string, userList string, callback init.Base) {
	open_im_sdk.userForSDK.GetGroupMembersInfo(groupId, userList, callback)
}

func KickGroupMember(groupId string, reason string, userList string, callback init.Base) {
	open_im_sdk.userForSDK.KickGroupMember(groupId, reason, userList, callback)
}

func TransferGroupOwner(groupId, userId string, callback init.Base) {
	open_im_sdk.userForSDK.TransferGroupOwner(groupId, userId, callback)
}

func InviteUserToGroup(groupId, reason string, userList string, callback init.Base) {
	open_im_sdk.userForSDK.InviteUserToGroup(groupId, reason, userList, callback)
}

func GetGroupApplicationList(callback init.Base) {
	open_im_sdk.userForSDK.GetGroupApplicationList(callback)
}

func AcceptGroupApplication(application, reason string, callback init.Base) {
	open_im_sdk.userForSDK.AcceptGroupApplication(application, reason, callback)
}

func RefuseGroupApplication(application, reason string, callback init.Base) {
	open_im_sdk.userForSDK.RefuseGroupApplication(application, reason, callback)
}

/////////////////////////////////////////////////////////////////

func GetDesignatedFriendsInfo(callback common.Base, userIDList, operationID string) {
	open_im_sdk.userForSDK.GetDesignatedFriendsInfo(callback, userIDList, operationID)
}

func GetFriendList(callback common.Base, operationID string) {
	open_im_sdk.userForSDK.GetFriendList(callback, operationID)
}

func CheckFriend(callback common.Base, userIDList, operationID string) {
	open_im_sdk.userForSDK.CheckFriend(callback, userIDList, operationID)
}

func AddFriend(callback common.Base, paramsReq, operationID string) {
	open_im_sdk.userForSDK.AddFriend(callback, paramsReq, operationID)
}

func SetFriendRemark(callback common.Base, params, operationID string) {
	open_im_sdk.userForSDK.SetFriendRemark(callback, params, operationID)
}
func DeleteFriend(callback common.Base, friendUserID, operationID string) {
	open_im_sdk.userForSDK.DeleteFriend(callback, friendUserID, operationID)
}

func GetRecvFriendApplicationList(callback common.Base, operationID string) {
	open_im_sdk.userForSDK.GetRecvFriendApplicationList(callback, operationID)
}

func GetSendFriendApplicationList(callback common.Base, operationID string) {
	open_im_sdk.userForSDK.GetSendFriendApplicationList(callback, operationID)
}

func AcceptFriendApplication(callback common.Base, params string, operationID string) {
	open_im_sdk.userForSDK.AcceptFriendApplication(callback, params, operationID)
}

func RefuseFriendApplication(callback common.Base, params, operationID string) {
	open_im_sdk.userForSDK.RefuseFriendApplication(callback, params, operationID)
}

func AddBlack(callback common.Base, blackUserID, operationID string) {
	open_im_sdk.userForSDK.AddBlack(callback, blackUserID, operationID)
}

func GetBlackList(callback common.Base, operationID string) {
	open_im_sdk.userForSDK.GetBlackList(callback, operationID)
}

func RemoveBlack(callback common.Base, removeUserID, operationID string) {
	open_im_sdk.userForSDK.RemoveBlack(callback, removeUserID, operationID)
}

func SetFriendListener(listener friend.OnFriendshipListener) bool {
	return open_im_sdk.userForSDK.SetFriendListener(listener)
}

///////////////////////////////////////////////////////////

func GetAllConversationList(callback init.Base) {
	open_im_sdk.userForSDK.GetAllConversationList(callback)
}
func GetConversationListSplit(callback init.Base, offset, count int) {
	open_im_sdk.userForSDK.GetConversationListSplit(callback, offset, count)
}
func SetConversationRecvMessageOpt(callback init.Base, conversationIDList string, opt int) {
	open_im_sdk.userForSDK.SetConversationRecvMessageOpt(callback, conversationIDList, opt)
}

func GetConversationRecvMessageOpt(callback init.Base, conversationIDList string) {
	open_im_sdk.userForSDK.GetConversationRecvMessageOpt(callback, conversationIDList)
}
func GetOneConversation(sourceID string, sessionType int, callback init.Base) {
	open_im_sdk.userForSDK.GetOneConversation(sourceID, sessionType, callback)
}
func GetMultipleConversation(conversationIDList string, callback init.Base) {
	open_im_sdk.userForSDK.GetMultipleConversation(conversationIDList, callback)
}
func DeleteConversation(conversationID string, callback init.Base) {
	open_im_sdk.userForSDK.DeleteConversation(conversationID, callback)
}
func SetConversationDraft(conversationID, draftText string, callback init.Base) {
	open_im_sdk.userForSDK.SetConversationDraft(conversationID, draftText, callback)
}
func PinConversation(conversationID string, isPinned bool, callback init.Base) {
	open_im_sdk.userForSDK.PinConversation(conversationID, isPinned, callback)
}
func GetTotalUnreadMsgCount(callback init.Base) {
	open_im_sdk.userForSDK.GetTotalUnreadMsgCount(callback)
}

func SetConversationListener(listener conversation_msg.OnConversationListener) {
	open_im_sdk.userForSDK.SetConversationListener(listener)
}

func AddAdvancedMsgListener(listener conversation_msg.OnAdvancedMsgListener) {
	open_im_sdk.userForSDK.AddAdvancedMsgListener(listener)
}

func CreateTextMessage(text string) string {
	return open_im_sdk.userForSDK.CreateTextMessage(text)
}
func CreateTextAtMessage(text, atUserList string) string {
	return open_im_sdk.userForSDK.CreateTextAtMessage(text, atUserList)
}
func CreateLocationMessage(description string, longitude, latitude float64) string {
	return open_im_sdk.userForSDK.CreateLocationMessage(description, longitude, latitude)
}
func CreateCustomMessage(data, extension string, description string) string {
	return open_im_sdk.userForSDK.CreateCustomMessage(data, extension, description)
}
func CreateQuoteMessage(text string, message string) string {
	return open_im_sdk.userForSDK.CreateQuoteMessage(text, message)
}
func CreateCardMessage(cardInfo string) string {
	return open_im_sdk.userForSDK.CreateCardMessage(cardInfo)

}
func CreateVideoMessageFromFullPath(videoFullPath string, videoType string, duration int64, snapshotFullPath string) string {
	return open_im_sdk.userForSDK.CreateVideoMessageFromFullPath(videoFullPath, videoType, duration, snapshotFullPath)
}
func CreateImageMessageFromFullPath(imageFullPath string) string {
	return open_im_sdk.userForSDK.CreateImageMessageFromFullPath(imageFullPath)
}
func CreateSoundMessageFromFullPath(soundPath string, duration int64) string {
	return open_im_sdk.userForSDK.CreateSoundMessageFromFullPath(soundPath, duration)
}
func CreateFileMessageFromFullPath(fileFullPath, fileName string) string {
	return open_im_sdk.userForSDK.CreateFileMessageFromFullPath(fileFullPath, fileName)
}
func CreateImageMessage(imagePath string) string {
	return open_im_sdk.userForSDK.CreateImageMessage(imagePath)
}
func CreateImageMessageByURL(sourcePicture, bigPicture, snapshotPicture string) string {
	return open_im_sdk.userForSDK.CreateImageMessageByURL(sourcePicture, bigPicture, snapshotPicture)
}
func SendMessageNotOss(callback conversation_msg.SendMsgCallBack, message, receiver, groupID string, onlineUserOnly bool, offlinePushInfo string) string {
	return open_im_sdk.userForSDK.SendMessageNotOss(callback, message, receiver, groupID, onlineUserOnly, offlinePushInfo)
}
func CreateSoundMessageByURL(soundBaseInfo string) string {
	return open_im_sdk.userForSDK.CreateSoundMessageByURL(soundBaseInfo)
}
func CreateSoundMessage(soundPath string, duration int64) string {
	return open_im_sdk.userForSDK.CreateSoundMessage(soundPath, duration)
}
func CreateVideoMessageByURL(videoBaseInfo string) string {
	return open_im_sdk.userForSDK.CreateVideoMessageByURL(videoBaseInfo)
}
func CreateVideoMessage(videoPath string, videoType string, duration int64, snapshotPath string) string {
	return open_im_sdk.userForSDK.CreateVideoMessage(videoPath, videoType, duration, snapshotPath)
}
func CreateFileMessageByURL(fileBaseInfo string) string {
	return open_im_sdk.userForSDK.CreateFileMessageByURL(fileBaseInfo)
}
func CreateFileMessage(filePath string, fileName string) string {
	return open_im_sdk.userForSDK.CreateFileMessage(filePath, fileName)
}
func CreateMergerMessage(messageList, title, summaryList string) string {
	return open_im_sdk.userForSDK.CreateMergerMessage(messageList, title, summaryList)
}

func CreateForwardMessage(m string) string {
	return open_im_sdk.userForSDK.CreateForwardMessage(m)
}

func SendMessage(callback conversation_msg.SendMsgCallBack, message, receiver, groupID string, onlineUserOnly bool, offlinePushInfo string) string {
	return open_im_sdk.userForSDK.SendMessage(callback, message, receiver, groupID, onlineUserOnly, offlinePushInfo)
}
func GetHistoryMessageList(callback init.Base, getMessageOptions string) {
	open_im_sdk.userForSDK.GetHistoryMessageList(callback, getMessageOptions)
}
func RevokeMessage(callback init.Base, message string) {
	open_im_sdk.userForSDK.RevokeMessage(callback, message)
}
func TypingStatusUpdate(receiver, msgTip string) {
	open_im_sdk.userForSDK.TypingStatusUpdate(receiver, msgTip)
}
func MarkC2CMessageAsRead(callback init.Base, receiver string, msgIDList string) {
	open_im_sdk.userForSDK.MarkC2CMessageAsRead(callback, receiver, msgIDList)
}

//Deprecated
func MarkSingleMessageHasRead(callback init.Base, userID string) {
	open_im_sdk.userForSDK.MarkSingleMessageHasRead(callback, userID)
}
func MarkGroupMessageHasRead(callback init.Base, groupID string) {
	open_im_sdk.userForSDK.MarkGroupMessageHasRead(callback, groupID)
}
func DeleteMessageFromLocalStorage(callback init.Base, message string) {
	open_im_sdk.userForSDK.DeleteMessageFromLocalStorage(callback, message)
}
func ClearC2CHistoryMessage(callback init.Base, userID string) {
	open_im_sdk.userForSDK.ClearC2CHistoryMessage(callback, userID)
}
func ClearGroupHistoryMessage(callback init.Base, groupID string) {
	open_im_sdk.userForSDK.ClearGroupHistoryMessage(callback, groupID)
}
func InsertSingleMessageToLocalStorage(callback init.Base, message, userID, sender string) string {
	return open_im_sdk.userForSDK.InsertSingleMessageToLocalStorage(callback, message, userID, sender)
}

func FindMessages(callback init.Base, messageIDList string) {
	open_im_sdk.userForSDK.FindMessages(callback, messageIDList)
}

func GetUsersInfo(uIDList string, cb init.Base) {
	open_im_sdk.userForSDK.GetUsersInfo(uIDList, cb)
}

func SetSelfInfo(info string, cb init.Base) {
	open_im_sdk.userForSDK.SetSelfInfo(info, cb)
}
