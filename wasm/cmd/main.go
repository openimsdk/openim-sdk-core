package main

import (
	"open_im_sdk/wasm/wasm_wrapper"

	"syscall/js"
)

//type BaseSuccessFailed struct {
//	funcName    string //e.g open_im_sdk/open_im_sdk.Login
//	operationID string
//	uid         string
//}
//
//func (b *BaseSuccessFailed) OnError(errCode int32, errMsg string) {
//	fmt.Println("OnError", errCode, errMsg)
//
//}
//
//func (b *BaseSuccessFailed) OnSuccess(data string) {
//	fmt.Println("OnError", data)
//}
//
//type InitCallback struct {
//	uid string
//}
//
//func (i *InitCallback) OnConnecting() {
//	fmt.Println("OnConnecting")
//}
//
//func (i *InitCallback) OnConnectSuccess() {
//	fmt.Println("OnConnecting")
//
//}
//
//func (i *InitCallback) OnConnectFailed(ErrCode int32, ErrMsg string) {
//	fmt.Println("OnConnecting", ErrCode, ErrMsg)
//
//}
//
//func (i *InitCallback) OnKickedOffline() {
//	fmt.Println("OnConnecting")
//
//}
//
//func (i *InitCallback) OnUserTokenExpired() {
//	fmt.Println("OnConnecting")
//
//}
//
//func (i *InitCallback) OnSelfInfoUpdated(userInfo string) {
//	fmt.Println("OnSelfInfoUpdated", userInfo)
//
//}
//
//var (
//	TESTIP = "43.155.69.205"
//	//TESTIP          = "121.37.25.71"
//	APIADDR = "http://" + TESTIP + ":10002"
//	WSADDR  = "ws://" + TESTIP + ":10001"
//)

func main() {
	//config := sdk_struct.IMConfig{
	//	Platform: 1,
	//	ApiAddr:  APIADDR,
	//	WsAddr:   WSADDR,
	//	DataDir:  "./",
	//	LogLevel: 6,
	//}
	//var listener InitCallback
	//var base BaseSuccessFailed
	//operationID := utils.OperationIDGenerator()
	//open_im_sdk.InitSDK(&listener, operationID, utils.StructToJsonString(config))
	//strMyUidx := "3984071717"
	//tokenx := test.RunGetToken(strMyUidx)
	//open_im_sdk.Login(&base, operationID, strMyUidx, tokenx)
	//opid := utils.OperationIDGenerator()
	//db := indexdb.NewIndexDB()
	//msg, err := db.GetMessage("client_msg_id_123")
	//if err != nil {
	//	log.Error(opid, "get message err:", err.Error())
	//} else {
	//	log.Info(opid, "get message is :", *msg, "get args is :", msg.ClientMsgID, msg.ServerMsgID)
	//}
	//var user model_struct.LocalUser
	//user.UserID = "111"
	//user.CreateTime = 1232
	//err = db.InsertLoginUser(&user)
	//if err != nil {
	//	log.Error(opid, "InsertLoginUser:", err.Error())
	//} else {
	//	log.Info(opid, "InsertLoginUser success:")
	//}
	//err = db.UpdateLoginUserByMap(&user, map[string]interface{}{"1": 3})
	//if err != nil {
	//	log.Error(opid, "UpdateLoginUserByMap:", err.Error())
	//} else {
	//	log.Info(opid, "UpdateLoginUserByMap success:")
	//}
	//seq, err := db.GetNormalMsgSeq()
	//if err != nil {
	//	log.Error(opid, "GetNormalMsgSeq:", err.Error())
	//} else {
	//	log.Info(opid, "GetNormalMsgSeq seq  success:", seq)
	//}

	registerFunc()
	<-make(chan bool)
}

func registerFunc() {
	//register global listener func
	globalFuc := wasm_wrapper.NewWrapperCommon()
	js.Global().Set(wasm_wrapper.COMMONEVENTFUNC, js.FuncOf(globalFuc.CommonEventFunc))
	//register init login func
	wrapperInitLogin := wasm_wrapper.NewWrapperInitLogin(globalFuc)
	js.Global().Set("initSDK", js.FuncOf(wrapperInitLogin.InitSDK))
	js.Global().Set("login", js.FuncOf(wrapperInitLogin.Login))
	js.Global().Set("logout", js.FuncOf(wrapperInitLogin.Logout))
	js.Global().Set("wakeUp", js.FuncOf(wrapperInitLogin.WakeUp))
	js.Global().Set("getLoginStatus", js.FuncOf(wrapperInitLogin.GetLoginStatus))
	//register conversation and message func
	wrapperConMsg := wasm_wrapper.NewWrapperConMsg(globalFuc)
	js.Global().Set("createTextMessage", js.FuncOf(wrapperConMsg.CreateTextMessage))
	js.Global().Set("createImageMessage", js.FuncOf(wrapperConMsg.CreateImageMessage))
	js.Global().Set("createImageMessageByURL", js.FuncOf(wrapperConMsg.CreateImageMessageByURL))
	js.Global().Set("createSoundMessageByURL", js.FuncOf(wrapperConMsg.CreateSoundMessageByURL))
	js.Global().Set("createVideoMessageByURL", js.FuncOf(wrapperConMsg.CreateVideoMessageByURL))
	js.Global().Set("createFileMessageByURL", js.FuncOf(wrapperConMsg.CreateFileMessageByURL))
	js.Global().Set("createCustomMessage", js.FuncOf(wrapperConMsg.CreateCustomMessage))
	js.Global().Set("createQuoteMessage", js.FuncOf(wrapperConMsg.CreateQuoteMessage))
	js.Global().Set("createAdvancedQuoteMessage", js.FuncOf(wrapperConMsg.CreateAdvancedQuoteMessage))
	js.Global().Set("createAdvancedTextMessage", js.FuncOf(wrapperConMsg.CreateAdvancedTextMessage))
	js.Global().Set("createCardMessage", js.FuncOf(wrapperConMsg.CreateCardMessage))
	js.Global().Set("createTextAtMessage", js.FuncOf(wrapperConMsg.CreateTextAtMessage))
	js.Global().Set("createVideoMessage", js.FuncOf(wrapperConMsg.CreateVideoMessage))
	js.Global().Set("createFileMessage", js.FuncOf(wrapperConMsg.CreateFileMessage))
	js.Global().Set("createMergerMessage", js.FuncOf(wrapperConMsg.CreateMergerMessage))
	js.Global().Set("createFaceMessage", js.FuncOf(wrapperConMsg.CreateFaceMessage))
	js.Global().Set("createForwardMessage", js.FuncOf(wrapperConMsg.CreateForwardMessage))
	js.Global().Set("createLocationMessage", js.FuncOf(wrapperConMsg.CreateLocationMessage))
	js.Global().Set("createVideoMessageFromFullPath", js.FuncOf(wrapperConMsg.CreateVideoMessageFromFullPath))
	js.Global().Set("createImageMessageFromFullPath", js.FuncOf(wrapperConMsg.CreateImageMessageFromFullPath))

	js.Global().Set("createSoundMessageFromFullPath", js.FuncOf(wrapperConMsg.CreateSoundMessageFromFullPath))
	js.Global().Set("createFileMessageFromFullPath", js.FuncOf(wrapperConMsg.CreateFileMessageFromFullPath))
	js.Global().Set("createSoundMessage", js.FuncOf(wrapperConMsg.CreateSoundMessage))
	js.Global().Set("createForwardMessage", js.FuncOf(wrapperConMsg.CreateForwardMessage))
	js.Global().Set("createLocationMessage", js.FuncOf(wrapperConMsg.CreateLocationMessage))
	js.Global().Set("createVideoMessageFromFullPath", js.FuncOf(wrapperConMsg.CreateVideoMessageFromFullPath))
	js.Global().Set("createImageMessageFromFullPath", js.FuncOf(wrapperConMsg.CreateImageMessageFromFullPath))
	js.Global().Set("markC2CMessageAsRead", js.FuncOf(wrapperConMsg.MarkC2CMessageAsRead))
	js.Global().Set("markMessageAsReadByConID", js.FuncOf(wrapperConMsg.MarkMessageAsReadByConID))
	js.Global().Set("sendMessage", js.FuncOf(wrapperConMsg.SendMessage))
	js.Global().Set("sendMessageNotOss", js.FuncOf(wrapperConMsg.SendMessageNotOss))
	js.Global().Set("getAllConversationList", js.FuncOf(wrapperConMsg.GetAllConversationList))
	js.Global().Set("getConversationListSplit", js.FuncOf(wrapperConMsg.GetConversationListSplit))
	js.Global().Set("getOneConversation", js.FuncOf(wrapperConMsg.GetOneConversation))
	js.Global().Set("deleteConversationFromLocalAndSvr", js.FuncOf(wrapperConMsg.DeleteConversationFromLocalAndSvr))
	js.Global().Set("getAdvancedHistoryMessageList", js.FuncOf(wrapperConMsg.GetAdvancedHistoryMessageList))
	js.Global().Set("getHistoryMessageList", js.FuncOf(wrapperConMsg.GetHistoryMessageList))
	js.Global().Set("getMultipleConversation", js.FuncOf(wrapperConMsg.GetMultipleConversation))
	js.Global().Set("setOneConversationPrivateChat", js.FuncOf(wrapperConMsg.SetOneConversationPrivateChat))
	js.Global().Set("setOneConversationRecvMessageOpt", js.FuncOf(wrapperConMsg.SetOneConversationRecvMessageOpt))
	js.Global().Set("setConversationRecvMessageOpt", js.FuncOf(wrapperConMsg.SetConversationRecvMessageOpt))
	js.Global().Set("setGlobalRecvMessageOpt", js.FuncOf(wrapperConMsg.SetGlobalRecvMessageOpt))
	js.Global().Set("deleteAllConversationFromLocal", js.FuncOf(wrapperConMsg.DeleteAllConversationFromLocal))
	js.Global().Set("setConversationDraft", js.FuncOf(wrapperConMsg.SetConversationDraft))
	js.Global().Set("resetConversationGroupAtType", js.FuncOf(wrapperConMsg.ResetConversationGroupAtType))
	js.Global().Set("pinConversation", js.FuncOf(wrapperConMsg.PinConversation))
	js.Global().Set("getTotalUnreadMsgCount", js.FuncOf(wrapperConMsg.GetTotalUnreadMsgCount))
	js.Global().Set("findMessageList", js.FuncOf(wrapperConMsg.FindMessageList))
	js.Global().Set("getHistoryMessageListReverse", js.FuncOf(wrapperConMsg.GetHistoryMessageListReverse))
	js.Global().Set("newRevokeMessage", js.FuncOf(wrapperConMsg.NewRevokeMessage))
	js.Global().Set("TypingStatusUpdate", js.FuncOf(wrapperConMsg.TypingStatusUpdate))
	js.Global().Set("markGroupMessageAsRead", js.FuncOf(wrapperConMsg.MarkGroupMessageAsRead))
	js.Global().Set("deleteMessageFromLocalStorage", js.FuncOf(wrapperConMsg.DeleteMessageFromLocalStorage))
	js.Global().Set("deleteMessageFromLocalAndSvr", js.FuncOf(wrapperConMsg.DeleteMessageFromLocalAndSvr))
	js.Global().Set("deleteAllMsgFromLocalAndSvr", js.FuncOf(wrapperConMsg.DeleteAllMsgFromLocalAndSvr))
	js.Global().Set("deleteAllMsgFromLocal", js.FuncOf(wrapperConMsg.DeleteAllMsgFromLocal))
	js.Global().Set("clearC2CHistoryMessage", js.FuncOf(wrapperConMsg.ClearC2CHistoryMessage))
	js.Global().Set("clearC2CHistoryMessageFromLocalAndSvr", js.FuncOf(wrapperConMsg.ClearC2CHistoryMessageFromLocalAndSvr))
	js.Global().Set("clearGroupHistoryMessage", js.FuncOf(wrapperConMsg.ClearGroupHistoryMessage))
	js.Global().Set("clearGroupHistoryMessageFromLocalAndSvr", js.FuncOf(wrapperConMsg.ClearGroupHistoryMessageFromLocalAndSvr))
	js.Global().Set("insertSingleMessageToLocalStorage", js.FuncOf(wrapperConMsg.InsertSingleMessageToLocalStorage))
	js.Global().Set("insertGroupMessageToLocalStorage", js.FuncOf(wrapperConMsg.InsertGroupMessageToLocalStorage))
	js.Global().Set("searchLocalMessages", js.FuncOf(wrapperConMsg.SearchLocalMessages))

	//register group func
	wrapperGroup := wasm_wrapper.NewWrapperGroup(globalFuc)
	js.Global().Set("createGroup", js.FuncOf(wrapperGroup.CreateGroup))
	js.Global().Set("getGroupsInfo", js.FuncOf(wrapperGroup.GetGroupsInfo))
	js.Global().Set("JoinGroup", js.FuncOf(wrapperGroup.JoinGroup))
	js.Global().Set("quitGroup", js.FuncOf(wrapperGroup.QuitGroup))
	js.Global().Set("dismissGroup", js.FuncOf(wrapperGroup.DismissGroup))
	js.Global().Set("changeGroupMute", js.FuncOf(wrapperGroup.ChangeGroupMute))
	js.Global().Set("changeGroupMemberMute", js.FuncOf(wrapperGroup.ChangeGroupMemberMute))
	js.Global().Set("setGroupMemberRoleLevel", js.FuncOf(wrapperGroup.SetGroupMemberRoleLevel))
	js.Global().Set("getJoinedGroupList", js.FuncOf(wrapperGroup.GetJoinedGroupList))
	js.Global().Set("searchGroups", js.FuncOf(wrapperGroup.SearchGroups))
	js.Global().Set("setGroupInfo", js.FuncOf(wrapperGroup.SetGroupInfo))
	js.Global().Set("setGroupVerification", js.FuncOf(wrapperGroup.SetGroupVerification))
	js.Global().Set("setGroupLookMemberInfo", js.FuncOf(wrapperGroup.SetGroupLookMemberInfo))
	js.Global().Set("setGroupApplyMemberFriend", js.FuncOf(wrapperGroup.SetGroupApplyMemberFriend))
	js.Global().Set("getGroupMemberList", js.FuncOf(wrapperGroup.GetGroupMemberList))
	js.Global().Set("getGroupMemberOwnerAndAdmin", js.FuncOf(wrapperGroup.GetGroupMemberOwnerAndAdmin))
	js.Global().Set("getGroupMemberListByJoinTimeFilter", js.FuncOf(wrapperGroup.GetGroupMemberListByJoinTimeFilter))
	js.Global().Set("getGroupMembersInfo", js.FuncOf(wrapperGroup.GetGroupMembersInfo))
	js.Global().Set("kickGroupMember", js.FuncOf(wrapperGroup.KickGroupMember))
	js.Global().Set("transferGroupOwner", js.FuncOf(wrapperGroup.TransferGroupOwner))
	js.Global().Set("inviteUserToGroup", js.FuncOf(wrapperGroup.InviteUserToGroup))
	js.Global().Set("getRecvGroupApplicationList", js.FuncOf(wrapperGroup.GetRecvGroupApplicationList))
	js.Global().Set("getSendGroupApplicationList", js.FuncOf(wrapperGroup.GetSendGroupApplicationList))
	js.Global().Set("acceptGroupApplication", js.FuncOf(wrapperGroup.AcceptGroupApplication))
	js.Global().Set("refuseGroupApplication", js.FuncOf(wrapperGroup.RefuseGroupApplication))
	js.Global().Set("setGroupMemberNickname", js.FuncOf(wrapperGroup.SetGroupMemberNickname))
	js.Global().Set("searchGroupMembers", js.FuncOf(wrapperGroup.SearchGroupMembers))
	js.Global().Set("getGroupsInfo", js.FuncOf(wrapperGroup.GetGroupsInfo))

	wrapperUser := wasm_wrapper.NewWrapperUser(globalFuc)
	js.Global().Set("getSelfUserInfo", js.FuncOf(wrapperUser.GetSelfUserInfo))
	js.Global().Set("setSelfInfo", js.FuncOf(wrapperUser.SetSelfInfo))
	js.Global().Set("getUsersInfo", js.FuncOf(wrapperUser.GetUsersInfo))

	wrapperFriend := wasm_wrapper.NewWrapperFriend(globalFuc)
	js.Global().Set("getDesignatedFriendsInfo", js.FuncOf(wrapperFriend.GetDesignatedFriendsInfo))
	js.Global().Set("getFriendList", js.FuncOf(wrapperFriend.GetFriendList))
	js.Global().Set("searchFriends", js.FuncOf(wrapperFriend.SearchFriends))
	js.Global().Set("checkFriend", js.FuncOf(wrapperFriend.CheckFriend))
	js.Global().Set("addFriend", js.FuncOf(wrapperFriend.AddFriend))
	js.Global().Set("setFriendRemark", js.FuncOf(wrapperFriend.SetFriendRemark))
	js.Global().Set("deleteFriend", js.FuncOf(wrapperFriend.DeleteFriend))
	js.Global().Set("getRecvFriendApplicationList", js.FuncOf(wrapperFriend.GetRecvFriendApplicationList))
	js.Global().Set("getSendFriendApplicationList", js.FuncOf(wrapperFriend.GetSendFriendApplicationList))
	js.Global().Set("acceptFriendApplication", js.FuncOf(wrapperFriend.AcceptFriendApplication))
	js.Global().Set("refuseFriendApplication", js.FuncOf(wrapperFriend.RefuseFriendApplication))
	js.Global().Set("getBlackList", js.FuncOf(wrapperFriend.GetBlackList))
	js.Global().Set("removeBlack", js.FuncOf(wrapperFriend.RemoveBlack))

	wrapperSignaling := wasm_wrapper.NewWrapperSignaling(globalFuc)
	js.Global().Set("signalingInviteInGroup", js.FuncOf(wrapperSignaling.SignalingInviteInGroup))
	js.Global().Set("signalingInvite", js.FuncOf(wrapperSignaling.SignalingInvite))
	js.Global().Set("signalingInvite", js.FuncOf(wrapperSignaling.SignalingInvite))
	js.Global().Set("signalingReject", js.FuncOf(wrapperSignaling.SignalingReject))
	js.Global().Set("signalingCancel", js.FuncOf(wrapperSignaling.SignalingCancel))
	js.Global().Set("signalingHungUp", js.FuncOf(wrapperSignaling.SignalingHungUp))
}
