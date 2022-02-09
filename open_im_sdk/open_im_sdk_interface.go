package open_im_sdk

import (
	"encoding/json"
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
const sdkVersion = "Open-IM-SDK-Core-v2.0.0"

func SdkVersion() string {
	return sdkVersion
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
func GetUsersInfo(callback open_im_sdk_callback.Base, operationID string, userIDList string) {
	userForSDK.User().GetUsersInfo(callback, userIDList, operationID)
}

func SetSelfInfo(callback open_im_sdk_callback.Base, operationID string, userInfo string) {
	userForSDK.User().SetSelfInfo(callback, userInfo, operationID)
}

//////////////////////////group//////////////////////////////////////////
func SetGroupListener(callback open_im_sdk_callback.OnGroupListener) {
	if callback == nil {
		log.Error("callback is nil")
		return
	}
	userForSDK.SetGroupListener(callback)
}

func CreateGroup(callback open_im_sdk_callback.Base, operationID string, groupBaseInfo string, memberList string) {
	if callback == nil {
		log.Error("callback is nil")
		return
	}
	userForSDK.Group().CreateGroup(callback, groupBaseInfo, memberList, operationID)
}

func JoinGroup(callback open_im_sdk_callback.Base, operationID string, groupID, reqMsg string) {
	if callback == nil {
		log.Error("callback is nil")
		return
	}
	userForSDK.Group().JoinGroup(callback, groupID, reqMsg, operationID)
}

func QuitGroup(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	userForSDK.Group().QuitGroup(callback, groupID, operationID)
}

func GetJoinedGroupList(callback open_im_sdk_callback.Base, operationID string) {
	userForSDK.Group().GetJoinedGroupList(callback, operationID)
}

func GetGroupsInfo(callback open_im_sdk_callback.Base, operationID string, groupIDList string) {
	userForSDK.Group().GetGroupsInfo(callback, groupIDList, operationID)
}

func SetGroupInfo(callback open_im_sdk_callback.Base, operationID string, groupID string, groupInfo string) {
	userForSDK.Group().SetGroupInfo(callback, groupInfo, groupID, operationID)
}

func GetGroupMemberList(callback open_im_sdk_callback.Base, operationID string, groupID string, filter, offset, count int32) {
	userForSDK.Group().GetGroupMemberList(callback, groupID, filter, offset, count, operationID)
}

func GetGroupMembersInfo(callback open_im_sdk_callback.Base, operationID string, groupID string, userIDList string) {
	userForSDK.Group().GetGroupMembersInfo(callback, groupID, userIDList, operationID)
}

func KickGroupMember(callback open_im_sdk_callback.Base, operationID string, groupID string, reason string, userIDList string) {
	userForSDK.Group().KickGroupMember(callback, groupID, reason, userIDList, operationID)
}

func TransferGroupOwner(callback open_im_sdk_callback.Base, operationID string, groupID, newOwnerUserID string) {
	userForSDK.Group().TransferGroupOwner(callback, groupID, newOwnerUserID, operationID)
}

func InviteUserToGroup(callback open_im_sdk_callback.Base, operationID string, groupID, reason string, userIDList string) {
	userForSDK.Group().InviteUserToGroup(callback, groupID, reason, userIDList, operationID)
}

func GetRecvGroupApplicationList(callback open_im_sdk_callback.Base, operationID string) {
	userForSDK.Group().GetRecvGroupApplicationList(callback, operationID)
}

func AcceptGroupApplication(callback open_im_sdk_callback.Base, operationID string, groupID, fromUserID, handleMsg string) {
	userForSDK.Group().AcceptGroupApplication(callback, groupID, fromUserID, handleMsg, operationID)
}

func RefuseGroupApplication(callback open_im_sdk_callback.Base, operationID string, groupID, fromUserID, handleMsg string) {
	userForSDK.Group().RefuseGroupApplication(callback, groupID, fromUserID, handleMsg, operationID)
}

////////////////////////////friend/////////////////////////////////////

func GetDesignatedFriendsInfo(callback open_im_sdk_callback.Base, operationID string, userIDList string) {
	userForSDK.Friend().GetDesignatedFriendsInfo(callback, userIDList, operationID)
}

func GetFriendList(callback open_im_sdk_callback.Base, operationID string) {
	userForSDK.Friend().GetFriendList(callback, operationID)
}

func CheckFriend(callback open_im_sdk_callback.Base, operationID string, userIDList string) {
	userForSDK.Friend().CheckFriend(callback, userIDList, operationID)
}

func AddFriend(callback open_im_sdk_callback.Base, operationID string, userIDReqMsg string) {
	userForSDK.Friend().AddFriend(callback, userIDReqMsg, operationID)
}

func SetFriendRemark(callback open_im_sdk_callback.Base, operationID string, userIDRemark string) {
	userForSDK.Friend().SetFriendRemark(callback, userIDRemark, operationID)
}
func DeleteFriend(callback open_im_sdk_callback.Base, operationID string, friendUserID string) {
	userForSDK.Friend().DeleteFriend(callback, friendUserID, operationID)
}

func GetRecvFriendApplicationList(callback open_im_sdk_callback.Base, operationID string) {
	userForSDK.Friend().GetRecvFriendApplicationList(callback, operationID)
}

func GetSendFriendApplicationList(callback open_im_sdk_callback.Base, operationID string) {
	userForSDK.Friend().GetSendFriendApplicationList(callback, operationID)
}

func AcceptFriendApplication(callback open_im_sdk_callback.Base, operationID string, userIDHandleMsg string) {
	userForSDK.Friend().AcceptFriendApplication(callback, userIDHandleMsg, operationID)
}

func RefuseFriendApplication(callback open_im_sdk_callback.Base, operationID string, userIDHandleMsg string) {
	userForSDK.Friend().RefuseFriendApplication(callback, userIDHandleMsg, operationID)
}

func AddBlack(callback open_im_sdk_callback.Base, operationID string, blackUserID string) {
	userForSDK.Friend().AddBlack(callback, blackUserID, operationID)
}

func GetBlackList(callback open_im_sdk_callback.Base, operationID string) {
	userForSDK.Friend().GetBlackList(callback, operationID)
}

func RemoveBlack(callback open_im_sdk_callback.Base, operationID string, removeUserID string) {
	userForSDK.Friend().RemoveBlack(callback, removeUserID, operationID)
}

func SetFriendListener(listener open_im_sdk_callback.OnFriendshipListener) {
	userForSDK.SetFriendListener(listener)
}

///////////////////////conversation////////////////////////////////////

func GetAllConversationList(callback open_im_sdk_callback.Base, operationID string) {
	userForSDK.Conversation().GetAllConversationList(callback, operationID)
}
func GetConversationListSplit(callback open_im_sdk_callback.Base, operationID string, offset, count int) {
	userForSDK.Conversation().GetConversationListSplit(callback, offset, count, operationID)
}

func SetConversationRecvMessageOpt(callback open_im_sdk_callback.Base, operationID string, conversationIDList string, opt int) {
	userForSDK.Conversation().SetConversationRecvMessageOpt(callback, conversationIDList, opt, operationID)
}

func GetConversationRecvMessageOpt(callback open_im_sdk_callback.Base, operationID string, conversationIDList string) {
	userForSDK.Conversation().GetConversationRecvMessageOpt(callback, conversationIDList, operationID)
}
func GetOneConversation(callback open_im_sdk_callback.Base, operationID string, sessionType int, sourceID string) {
	userForSDK.Conversation().GetOneConversation(callback, int32(sessionType), sourceID, operationID)
}
func GetMultipleConversation(callback open_im_sdk_callback.Base, operationID string, conversationIDList string) {
	userForSDK.Conversation().GetMultipleConversation(callback, conversationIDList, operationID)
}
func DeleteConversation(callback open_im_sdk_callback.Base, operationID string, conversationID string) {
	userForSDK.Conversation().DeleteConversation(callback, conversationID, operationID)
}
func SetConversationDraft(callback open_im_sdk_callback.Base, operationID string, conversationID, draftText string) {
	userForSDK.Conversation().SetConversationDraft(callback, conversationID, draftText, operationID)
}
func PinConversation(callback open_im_sdk_callback.Base, operationID string, conversationID string, isPinned bool) {
	userForSDK.Conversation().PinConversation(callback, conversationID, isPinned, operationID)
}
func GetTotalUnreadMsgCount(callback open_im_sdk_callback.Base, operationID string) {
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
	userForSDK.Conversation().SendMessage(callback, message, recvID, groupID, offlinePushInfo, operationID)
}
func SendMessageNotOss(callback open_im_sdk_callback.SendMsgCallBack, operationID string, message, recvID, groupID string, offlinePushInfo string) {
	userForSDK.Conversation().SendMessageNotOss(callback, message, recvID, groupID, offlinePushInfo, operationID)
}

func GetHistoryMessageList(callback open_im_sdk_callback.Base, operationID string, getMessageOptions string) {
	userForSDK.Conversation().GetHistoryMessageList(callback, getMessageOptions, operationID)
}

func RevokeMessage(callback open_im_sdk_callback.Base, operationID string, message string) {
	userForSDK.Conversation().RevokeMessage(callback, message, operationID)
}
func TypingStatusUpdate(callback open_im_sdk_callback.Base, operationID string, recvID, msgTip string) {
	userForSDK.Conversation().TypingStatusUpdate(callback, recvID, msgTip, operationID)
}
func MarkC2CMessageAsRead(callback open_im_sdk_callback.Base, operationID string, userID string, msgIDList string) {
	userForSDK.Conversation().MarkC2CMessageAsRead(callback, userID, msgIDList, operationID)
}

func MarkGroupMessageHasRead(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	userForSDK.Conversation().MarkGroupMessageHasRead(callback, groupID, operationID)
}
func DeleteMessageFromLocalStorage(callback open_im_sdk_callback.Base, operationID string, message string) {
	userForSDK.Conversation().DeleteMessageFromLocalStorage(callback, message, operationID)
}
func ClearC2CHistoryMessage(callback open_im_sdk_callback.Base, operationID string, userID string) {
	userForSDK.Conversation().ClearC2CHistoryMessage(callback, userID, operationID)
}
func ClearGroupHistoryMessage(callback open_im_sdk_callback.Base, operationID string, groupID string) {
	userForSDK.Conversation().ClearGroupHistoryMessage(callback, groupID, operationID)
}
func InsertSingleMessageToLocalStorage(callback open_im_sdk_callback.Base, operationID string, message, recvID, sendID string) {
	userForSDK.Conversation().InsertSingleMessageToLocalStorage(callback, message, recvID, sendID, operationID)
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
	switch sessionType {
	case constant.SingleChatType:
		return "single_" + sourceID
	case constant.GroupChatType:
		return "group_" + sourceID
	}
	return ""
}
