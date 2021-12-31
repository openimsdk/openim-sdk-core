package open_im_sdk

import (
	"encoding/json"
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

func InitSDK(config string, cb IMSDKListener) bool {
	var sc IMConfig
	if err := json.Unmarshal([]byte(config), &sc); err != nil {
		sdkLog("initSDK failed, config: ", sc, err.Error())
		return false
	}
	sdkLog("InitSDK, config ", config, "version: ", SdkVersion())
	if userForSDK != nil {
		sdkLog("Logout first ")
		userForSDK.Logout(nil)
		sdkLog("unInit first ")
		userForSDK.UnInitSDK()
	}
	userForSDK = new(UserRelated)

	InitOnce(&sc)

	return userForSDK.InitSDK(config, cb)
}

//1 no print
func SetSdkLog(flag int32) {
	SdkLogFlag = flag
}
func (u *UserRelated) SetSdkLog(flag int32) {
	SdkLogFlag = flag
}
func SetHearbeatInterval(interval int32) {
	hearbeatInterval = interval
}

func UnInitSDK() {
	if userForSDK == nil {
		sdkLog("userForSDK nil")
		return
	}
	userForSDK.unInitSDK()
}

func Login(uid, tk string, callback Base) {
	if userForSDK == nil {
		sdkLog("userForSDK nil")
		if callback != nil {
			callback.OnError(ErrCodeInitLogin, "userForSDK nil, initSdk first ")
		}
		return
	}
	userForSDK.Login(uid, tk, callback)
}

func Logout(callback Base) {
	if userForSDK == nil {
		sdkLog("userForSDK nil")
		if callback != nil {
			callback.OnError(ErrCodeInitLogin, "userForSDK nil, initSdk first ")
		}
		return
	}
	userForSDK.logout(callback)
}

func GetLoginStatus() int {
	return userForSDK.getLoginStatus()
}

func GetLoginUser() string {
	return userForSDK.GetLoginUser()
}

func ForceSyncLoginUerInfo() {
	userForSDK.ForceSyncLoginUserInfo()
}

func ForceSyncMsg() bool {
	return userForSDK.ForceSyncMsg()
}

func SetGroupListener(callback OnGroupListener) {
	userForSDK.SetGroupListener(callback)
}

func CreateGroup(gInfo string, memberList string, callback Base) {
	userForSDK.CreateGroup(gInfo, memberList, callback)
}

func JoinGroup(groupId, message string, callback Base) {
	userForSDK.JoinGroup(groupId, message, callback)
}

func QuitGroup(groupId string, callback Base) {
	userForSDK.QuitGroup(groupId, callback)
}

func GetJoinedGroupList(callback Base) {
	userForSDK.GetJoinedGroupList(callback)
}

func GetGroupsInfo(groupIdList string, callback Base) {
	userForSDK.GetGroupsInfo(groupIdList, callback)
}

func SetGroupInfo(jsonGroupInfo string, callback Base) {
	userForSDK.SetGroupInfo(jsonGroupInfo, callback)
}

func GetGroupMemberList(groupId string, filter int32, next int32, callback Base) {
	userForSDK.GetGroupMemberList(groupId, filter, next, callback)
}

func GetGroupMembersInfo(groupId string, userList string, callback Base) {
	userForSDK.GetGroupMembersInfo(groupId, userList, callback)
}

func KickGroupMember(groupId string, reason string, userList string, callback Base) {
	userForSDK.KickGroupMember(groupId, reason, userList, callback)
}

func TransferGroupOwner(groupId, userId string, callback Base) {
	userForSDK.TransferGroupOwner(groupId, userId, callback)
}

func InviteUserToGroup(groupId, reason string, userList string, callback Base) {
	userForSDK.InviteUserToGroup(groupId, reason, userList, callback)
}

func GetGroupApplicationList(callback Base) {
	userForSDK.GetGroupApplicationList(callback)
}

func AcceptGroupApplication(application, reason string, callback Base) {
	userForSDK.AcceptGroupApplication(application, reason, callback)
}

func RefuseGroupApplication(application, reason string, callback Base) {
	userForSDK.RefuseGroupApplication(application, reason, callback)
}

/////////////////////////////////////////////////////////////////

func GetFriendsInfo(callback Base, uidList string) {
	userForSDK.GetFriendsInfo(callback, uidList)
}

func AddFriend(callback Base, paramsReq string) {
	userForSDK.AddFriend(callback, paramsReq)
}

func GetFriendApplicationList(callback Base) {
	userForSDK.GetFriendApplicationList(callback)
}

func AcceptFriendApplication(callback Base, uid string) {
	userForSDK.AcceptFriendApplication(callback, uid)
}

func RefuseFriendApplication(callback Base, uid string) {
	userForSDK.RefuseFriendApplication(callback, uid)
}

func CheckFriend(callback Base, uidList string) {
	userForSDK.CheckFriend(callback, uidList)
}

func DeleteFromFriendList(deleteUid string, callback Base) {
	userForSDK.DeleteFromFriendList(deleteUid, callback)
}

func GetFriendList(callback Base) {
	userForSDK.GetFriendList(callback)
}

func SetFriendInfo(comment string, callback Base) {
	userForSDK.SetFriendInfo(comment, callback)
}

func AddToBlackList(callback Base, blackUid string) {
	userForSDK.AddToBlackList(callback, blackUid)
}

func GetBlackList(callback Base) {
	userForSDK.GetBlackList(callback)
}

func DeleteFromBlackList(callback Base, deleteUid string) {
	userForSDK.DeleteFromBlackList(callback, deleteUid)
}

func SetFriendListener(listener OnFriendshipListener) bool {
	return userForSDK.SetFriendListener(listener)
}

///////////////////////////////////////////////////////////

func GetAllConversationList(callback Base) {
	userForSDK.GetAllConversationList(callback)
}
func GetConversationListSplit(callback Base, offset, count int) {
	userForSDK.GetConversationListSplit(callback, offset, count)
}
func SetConversationRecvMessageOpt(callback Base, conversationIDList string, opt int) {
	userForSDK.SetConversationRecvMessageOpt(callback, conversationIDList, opt)
}

func GetConversationRecvMessageOpt(callback Base, conversationIDList string) {
	userForSDK.GetConversationRecvMessageOpt(callback, conversationIDList)
}
func GetOneConversation(sourceID string, sessionType int, callback Base) {
	userForSDK.GetOneConversation(sourceID, sessionType, callback)
}
func GetMultipleConversation(conversationIDList string, callback Base) {
	userForSDK.GetMultipleConversation(conversationIDList, callback)
}
func DeleteConversation(conversationID string, callback Base) {
	userForSDK.DeleteConversation(conversationID, callback)
}
func SetConversationDraft(conversationID, draftText string, callback Base) {
	userForSDK.SetConversationDraft(conversationID, draftText, callback)
}
func PinConversation(conversationID string, isPinned bool, callback Base) {
	userForSDK.PinConversation(conversationID, isPinned, callback)
}
func GetTotalUnreadMsgCount(callback Base) {
	userForSDK.GetTotalUnreadMsgCount(callback)
}

func SetConversationListener(listener OnConversationListener) {
	userForSDK.SetConversationListener(listener)
}

func AddAdvancedMsgListener(listener OnAdvancedMsgListener) {
	userForSDK.AddAdvancedMsgListener(listener)
}

func CreateTextMessage(text string) string {
	return userForSDK.CreateTextMessage(text)
}
func CreateTextAtMessage(text, atUserList string) string {
	return userForSDK.CreateTextAtMessage(text, atUserList)
}
func CreateLocationMessage(description string, longitude, latitude float64) string {
	return userForSDK.CreateLocationMessage(description, longitude, latitude)
}
func CreateCustomMessage(data, extension string, description string) string {
	return userForSDK.CreateCustomMessage(data, extension, description)
}
func CreateQuoteMessage(text string, message string) string {
	return userForSDK.CreateQuoteMessage(text, message)
}
func CreateCardMessage(cardInfo string) string {
	return userForSDK.CreateCardMessage(cardInfo)

}
func CreateVideoMessageFromFullPath(videoFullPath string, videoType string, duration int64, snapshotFullPath string) string {
	return userForSDK.CreateVideoMessageFromFullPath(videoFullPath, videoType, duration, snapshotFullPath)
}
func CreateImageMessageFromFullPath(imageFullPath string) string {
	return userForSDK.CreateImageMessageFromFullPath(imageFullPath)
}
func CreateSoundMessageFromFullPath(soundPath string, duration int64) string {
	return userForSDK.CreateSoundMessageFromFullPath(soundPath, duration)
}
func CreateFileMessageFromFullPath(fileFullPath, fileName string) string {
	return userForSDK.CreateFileMessageFromFullPath(fileFullPath, fileName)
}
func CreateImageMessage(imagePath string) string {
	return userForSDK.CreateImageMessage(imagePath)
}
func CreateImageMessageByURL(sourcePicture, bigPicture, snapshotPicture string) string {
	return userForSDK.CreateImageMessageByURL(sourcePicture, bigPicture, snapshotPicture)
}
func SendMessageNotOss(callback SendMsgCallBack, message, receiver, groupID string, onlineUserOnly bool, offlinePushInfo string) string {
	return userForSDK.SendMessageNotOss(callback, message, receiver, groupID, onlineUserOnly, offlinePushInfo)
}
func CreateSoundMessageByURL(soundBaseInfo string) string {
	return userForSDK.CreateSoundMessageByURL(soundBaseInfo)
}
func CreateSoundMessage(soundPath string, duration int64) string {
	return userForSDK.CreateSoundMessage(soundPath, duration)
}
func CreateVideoMessageByURL(videoBaseInfo string) string {
	return userForSDK.CreateVideoMessageByURL(videoBaseInfo)
}
func CreateVideoMessage(videoPath string, videoType string, duration int64, snapshotPath string) string {
	return userForSDK.CreateVideoMessage(videoPath, videoType, duration, snapshotPath)
}
func CreateFileMessageByURL(fileBaseInfo string) string {
	return userForSDK.CreateFileMessageByURL(fileBaseInfo)
}
func CreateFileMessage(filePath string, fileName string) string {
	return userForSDK.CreateFileMessage(filePath, fileName)
}
func CreateMergerMessage(messageList, title, summaryList string) string {
	return userForSDK.CreateMergerMessage(messageList, title, summaryList)
}

func CreateForwardMessage(m string) string {
	return userForSDK.CreateForwardMessage(m)
}

func SendMessage(callback SendMsgCallBack, message, receiver, groupID string, onlineUserOnly bool, offlinePushInfo string) string {
	return userForSDK.SendMessage(callback, message, receiver, groupID, onlineUserOnly, offlinePushInfo)
}
func GetHistoryMessageList(callback Base, getMessageOptions string) {
	userForSDK.GetHistoryMessageList(callback, getMessageOptions)
}
func RevokeMessage(callback Base, message string) {
	userForSDK.RevokeMessage(callback, message)
}
func TypingStatusUpdate(receiver, msgTip string) {
	userForSDK.TypingStatusUpdate(receiver, msgTip)
}
func MarkC2CMessageAsRead(callback Base, receiver string, msgIDList string) {
	userForSDK.MarkC2CMessageAsRead(callback, receiver, msgIDList)
}

//Deprecated
func MarkSingleMessageHasRead(callback Base, userID string) {
	userForSDK.MarkSingleMessageHasRead(callback, userID)
}
func MarkGroupMessageHasRead(callback Base, groupID string) {
	userForSDK.MarkGroupMessageHasRead(callback, groupID)
}
func DeleteMessageFromLocalStorage(callback Base, message string) {
	userForSDK.DeleteMessageFromLocalStorage(callback, message)
}
func ClearC2CHistoryMessage(callback Base, userID string) {
	userForSDK.ClearC2CHistoryMessage(callback, userID)
}
func ClearGroupHistoryMessage(callback Base, groupID string) {
	userForSDK.ClearGroupHistoryMessage(callback, groupID)
}
func InsertSingleMessageToLocalStorage(callback Base, message, userID, sender string) string {
	return userForSDK.InsertSingleMessageToLocalStorage(callback, message, userID, sender)
}

func FindMessages(callback Base, messageIDList string) {
	userForSDK.FindMessages(callback, messageIDList)
}

func GetUsersInfo(uIDList string, cb Base) {
	userForSDK.GetUsersInfo(uIDList, cb)
}

func SetSelfInfo(info string, cb Base) {
	userForSDK.SetSelfInfo(info, cb)
}
