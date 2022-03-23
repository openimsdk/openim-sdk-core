package open_im_sdk

import (
	"encoding/json"
	"errors"
	common2 "open_im_sdk/internal/common"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/internal/login"
	"open_im_sdk/open_im_sdk_callback"
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

func SdkVersion() string {
	return constant.SdkVersion + constant.BigVersion + constant.UpdateVersion
}

func InitSDK(listener open_im_sdk_callback.OnConnListener, operationID string, config string) bool {
	if err := json.Unmarshal([]byte(config), &sdk_struct.SvrConf); err != nil {
		log.Error(operationID, "Unmarshal failed ", err.Error(), config)
		return false
	}
	log.Info(operationID, "config ", config, sdk_struct.SvrConf)
	log.NewPrivateLog("", sdk_struct.SvrConf.LogLevel)
	log.NewInfo(operationID, utils.GetSelfFuncName(), config, SdkVersion())
	if listener == nil || config == "" {
		log.Error(operationID, "listener or config is nil")
		return false
	}
	if userForSDK != nil {
		log.Warn(operationID, "Initialize multiple times, call logout")
		userForSDK.Logout(nil, utils.OperationIDGenerator())
	}
	userForSDK = new(login.LoginMgr)
	return userForSDK.InitSDK(sdk_struct.SvrConf, listener, operationID)
}

func Login(callback open_im_sdk_callback.Base, operationID string, userID, token string) {
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

func UploadImage(callback open_im_sdk_callback.Base, operationID string, filePath string, token, obj string) string {
	return userForSDK.UploadImage(callback, filePath, token, obj, operationID)
}

func Logout(callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		log.Error("callback is nil")
		return
	}
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
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
func GetUsersInfo(callback open_im_sdk_callback.Base, operationID string, userIDList string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Full().GetUsersInfo(callback, userIDList, operationID)
}

func SetSelfInfo(callback open_im_sdk_callback.Base, operationID string, userInfo string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.User().SetSelfInfo(callback, userInfo, operationID)
}

func GetSelfUserInfo(callback open_im_sdk_callback.Base, operationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.User().GetSelfUserInfo(callback, operationID)
}

//////////////////////////group//////////////////////////////////////////
func SetGroupListener(callback open_im_sdk_callback.OnGroupListener) {
	if callback == nil || userForSDK == nil {
		log.Error("callback or userForSDK is nil")
		return
	}

	userForSDK.SetGroupListener(callback)
}

func CreateGroup(callback open_im_sdk_callback.Base, operationID string, groupBaseInfo string, memberList string) {
	if callback == nil {
		log.Error("callback is nil")
		return
	}
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().CreateGroup(callback, groupBaseInfo, memberList, operationID)
}

func JoinGroup(callback open_im_sdk_callback.Base, operationID string, groupID, reqMsg string) {
	if callback == nil {
		log.Error("callback is nil")
		return
	}
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().JoinGroup(callback, groupID, reqMsg, operationID)
}

func QuitGroup(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().QuitGroup(callback, groupID, operationID)
}

func GetJoinedGroupList(callback open_im_sdk_callback.Base, operationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().GetJoinedGroupList(callback, operationID)
}

func GetGroupsInfo(callback open_im_sdk_callback.Base, operationID string, groupIDList string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().GetGroupsInfo(callback, groupIDList, operationID)
}

func SetGroupInfo(callback open_im_sdk_callback.Base, operationID string, groupID string, groupInfo string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().SetGroupInfo(callback, groupInfo, groupID, operationID)
}

func GetGroupMemberList(callback open_im_sdk_callback.Base, operationID string, groupID string, filter, offset, count int32) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().GetGroupMemberList(callback, groupID, filter, offset, count, operationID)
}

func GetGroupMembersInfo(callback open_im_sdk_callback.Base, operationID string, groupID string, userIDList string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().GetGroupMembersInfo(callback, groupID, userIDList, operationID)
}

func KickGroupMember(callback open_im_sdk_callback.Base, operationID string, groupID string, reason string, userIDList string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().KickGroupMember(callback, groupID, reason, userIDList, operationID)
}

func TransferGroupOwner(callback open_im_sdk_callback.Base, operationID string, groupID, newOwnerUserID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().TransferGroupOwner(callback, groupID, newOwnerUserID, operationID)
}

func InviteUserToGroup(callback open_im_sdk_callback.Base, operationID string, groupID, reason string, userIDList string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().InviteUserToGroup(callback, groupID, reason, userIDList, operationID)
}

func GetRecvGroupApplicationList(callback open_im_sdk_callback.Base, operationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().GetRecvGroupApplicationList(callback, operationID)
}

func GetSendGroupApplicationList(callback open_im_sdk_callback.Base, operationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().GetSendGroupApplicationList(callback, operationID)
}

func AcceptGroupApplication(callback open_im_sdk_callback.Base, operationID string, groupID, fromUserID, handleMsg string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().AcceptGroupApplication(callback, groupID, fromUserID, handleMsg, operationID)
}

func RefuseGroupApplication(callback open_im_sdk_callback.Base, operationID string, groupID, fromUserID, handleMsg string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Group().RefuseGroupApplication(callback, groupID, fromUserID, handleMsg, operationID)
}

////////////////////////////friend/////////////////////////////////////

func GetDesignatedFriendsInfo(callback open_im_sdk_callback.Base, operationID string, userIDList string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Friend().GetDesignatedFriendsInfo(callback, userIDList, operationID)
}

func GetFriendList(callback open_im_sdk_callback.Base, operationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Friend().GetFriendList(callback, operationID)
}

func CheckFriend(callback open_im_sdk_callback.Base, operationID string, userIDList string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Friend().CheckFriend(callback, userIDList, operationID)
}

func AddFriend(callback open_im_sdk_callback.Base, operationID string, userIDReqMsg string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Friend().AddFriend(callback, userIDReqMsg, operationID)
}

func SetFriendRemark(callback open_im_sdk_callback.Base, operationID string, userIDRemark string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Friend().SetFriendRemark(callback, userIDRemark, operationID)
}
func DeleteFriend(callback open_im_sdk_callback.Base, operationID string, friendUserID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Friend().DeleteFriend(callback, friendUserID, operationID)
}

func GetRecvFriendApplicationList(callback open_im_sdk_callback.Base, operationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Friend().GetRecvFriendApplicationList(callback, operationID)
}

func GetSendFriendApplicationList(callback open_im_sdk_callback.Base, operationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Friend().GetSendFriendApplicationList(callback, operationID)
}

func AcceptFriendApplication(callback open_im_sdk_callback.Base, operationID string, userIDHandleMsg string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Friend().AcceptFriendApplication(callback, userIDHandleMsg, operationID)
}

func RefuseFriendApplication(callback open_im_sdk_callback.Base, operationID string, userIDHandleMsg string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Friend().RefuseFriendApplication(callback, userIDHandleMsg, operationID)
}

func AddBlack(callback open_im_sdk_callback.Base, operationID string, blackUserID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Friend().AddBlack(callback, blackUserID, operationID)
}

func GetBlackList(callback open_im_sdk_callback.Base, operationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Friend().GetBlackList(callback, operationID)
}

func RemoveBlack(callback open_im_sdk_callback.Base, operationID string, removeUserID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Friend().RemoveBlack(callback, removeUserID, operationID)
}

func SetFriendListener(listener open_im_sdk_callback.OnFriendshipListener) {
	userForSDK.SetFriendListener(listener)
}

///////////////////////conversation////////////////////////////////////

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

func SetConversationRecvMessageOpt(callback open_im_sdk_callback.Base, operationID string, conversationIDList string, opt int) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().SetConversationRecvMessageOpt(callback, conversationIDList, opt, operationID)
}

func GetConversationRecvMessageOpt(callback open_im_sdk_callback.Base, operationID string, conversationIDList string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().GetConversationRecvMessageOpt(callback, conversationIDList, operationID)
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
func DeleteConversation(callback open_im_sdk_callback.Base, operationID string, conversationID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().DeleteConversation(callback, conversationID, operationID)
}
func SetConversationDraft(callback open_im_sdk_callback.Base, operationID string, conversationID, draftText string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().SetConversationDraft(callback, conversationID, draftText, operationID)
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

//
func SetConversationListener(listener open_im_sdk_callback.OnConversationListener) {
	userForSDK.SetConversationListener(listener)
}
func SetAdvancedMsgListener(listener open_im_sdk_callback.OnAdvancedMsgListener) {
	userForSDK.SetAdvancedMsgListener(listener)
}

func SetUserListener(listener open_im_sdk_callback.OnUserListener) {
	userForSDK.SetUserListener(listener)
}

func CreateTextAtMessage(operationID string, text, atUserList string) string {
	return userForSDK.Conversation().CreateTextAtMessage(text, atUserList, operationID)
}

//
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

func GetHistoryMessageList(callback open_im_sdk_callback.Base, operationID string, getMessageOptions string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().GetHistoryMessageList(callback, getMessageOptions, operationID)
}

func RevokeMessage(callback open_im_sdk_callback.Base, operationID string, message string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().RevokeMessage(callback, message, operationID)
}
func UpdateMessage(callback open_im_sdk_callback.Base, operationID string, message string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().UpdateMessage(callback, message, operationID)
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

func MarkGroupMessageHasRead(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().MarkGroupMessageHasRead(callback, groupID, operationID)
}
func DeleteMessageFromLocalStorage(callback open_im_sdk_callback.Base, operationID string, message string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().DeleteMessageFromLocalStorage(callback, message, operationID)
}
func ClearC2CHistoryMessage(callback open_im_sdk_callback.Base, operationID string, userID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().ClearC2CHistoryMessage(callback, userID, operationID)
}
func ClearGroupHistoryMessage(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().ClearGroupHistoryMessage(callback, groupID, operationID)
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

//func FindMessages(callback common.Base, operationID string, messageIDList string) {
//	userForSDK.Conversation().FindMessages(callback, messageIDList)
//}

func InitOnce(config *sdk_struct.IMConfig) bool {
	sdk_struct.SvrConf = *config
	return true
}

func CheckToken(userID, token string) error {
	return login.CheckToken(userID, token, "")
}

func CheckResourceLoad(uSDK *login.LoginMgr) error {
	if uSDK == nil {
		//	callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return utils.Wrap(errors.New("CheckResourceLoad failed uSDK == nil "), "")
	}
	if uSDK.Friend() == nil || uSDK.User() == nil || uSDK.Group() == nil || uSDK.Conversation() == nil ||
		uSDK.Full() == nil {
		return utils.Wrap(errors.New("CheckResourceLoad failed, resource nil "), "")
	}
	return nil
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
