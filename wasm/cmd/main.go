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

//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"open_im_sdk/pkg/log"
	"open_im_sdk/wasm/wasm_wrapper"
	"runtime"
	"runtime/debug"

	"syscall/js"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("MAIN", "panic info is:", r, string(debug.Stack()))
		}
	}()
	fmt.Println("runtime env", runtime.GOARCH, runtime.GOOS)
	registerFunc()
	<-make(chan bool)
}

func registerFunc() {
	//register global listener function
	globalFuc := wasm_wrapper.NewWrapperCommon()
	js.Global().Set(wasm_wrapper.COMMONEVENTFUNC, js.FuncOf(globalFuc.CommonEventFunc))
	//register init login function
	wrapperInitLogin := wasm_wrapper.NewWrapperInitLogin(globalFuc)
	js.Global().Set("initSDK", js.FuncOf(wrapperInitLogin.InitSDK))
	js.Global().Set("login", js.FuncOf(wrapperInitLogin.Login))
	js.Global().Set("logout", js.FuncOf(wrapperInitLogin.Logout))
	js.Global().Set("getLoginStatus", js.FuncOf(wrapperInitLogin.GetLoginStatus))
	js.Global().Set("setAppBackgroundStatus", js.FuncOf(wrapperInitLogin.SetAppBackgroundStatus))
	js.Global().Set("networkStatusChanged", js.FuncOf(wrapperInitLogin.NetworkStatusChanged))
	//register conversation and message function
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
	js.Global().Set("getAtAllTag", js.FuncOf(wrapperConMsg.GetAtAllTag))
	js.Global().Set("markConversationMessageAsRead", js.FuncOf(wrapperConMsg.MarkConversationMessageAsRead))
	js.Global().Set("markMessagesAsReadByMsgID", js.FuncOf(wrapperConMsg.MarkMessagesAsReadByMsgID))
	js.Global().Set("sendMessage", js.FuncOf(wrapperConMsg.SendMessage))
	js.Global().Set("sendMessageNotOss", js.FuncOf(wrapperConMsg.SendMessageNotOss))
	js.Global().Set("sendMessageByBuffer", js.FuncOf(wrapperConMsg.SendMessageByBuffer))
	//js.Global().Set("setMessageReactionExtensions", js.FuncOf(wrapperConMsg.SetMessageReactionExtensions))
	//js.Global().Set("addMessageReactionExtensions", js.FuncOf(wrapperConMsg.AddMessageReactionExtensions))
	//js.Global().Set("deleteMessageReactionExtensions", js.FuncOf(wrapperConMsg.DeleteMessageReactionExtensions))
	//js.Global().Set("getMessageListReactionExtensions", js.FuncOf(wrapperConMsg.GetMessageListReactionExtensions))
	//js.Global().Set("getMessageListSomeReactionExtensions", js.FuncOf(wrapperConMsg.GetMessageListSomeReactionExtensions))
	js.Global().Set("getAllConversationList", js.FuncOf(wrapperConMsg.GetAllConversationList))
	js.Global().Set("getConversationListSplit", js.FuncOf(wrapperConMsg.GetConversationListSplit))
	js.Global().Set("getOneConversation", js.FuncOf(wrapperConMsg.GetOneConversation))
	js.Global().Set("deleteConversationAndDeleteAllMsg", js.FuncOf(wrapperConMsg.DeleteConversationAndDeleteAllMsg))
	js.Global().Set("getAdvancedHistoryMessageList", js.FuncOf(wrapperConMsg.GetAdvancedHistoryMessageList))
	js.Global().Set("getAdvancedHistoryMessageListReverse", js.FuncOf(wrapperConMsg.GetAdvancedHistoryMessageListReverse))
	js.Global().Set("getMultipleConversation", js.FuncOf(wrapperConMsg.GetMultipleConversation))
	js.Global().Set("setConversationPrivateChat", js.FuncOf(wrapperConMsg.SetConversationPrivateChat))
	js.Global().Set("setConversationRecvMessageOpt", js.FuncOf(wrapperConMsg.SetConversationRecvMessageOpt))
	js.Global().Set("setGlobalRecvMessageOpt", js.FuncOf(wrapperConMsg.SetGlobalRecvMessageOpt))
	js.Global().Set("deleteAllConversationFromLocal", js.FuncOf(wrapperConMsg.DeleteAllConversationFromLocal))
	js.Global().Set("hideConversation", js.FuncOf(wrapperConMsg.HideConversation))
	js.Global().Set("setConversationDraft", js.FuncOf(wrapperConMsg.SetConversationDraft))
	js.Global().Set("resetConversationGroupAtType", js.FuncOf(wrapperConMsg.ResetConversationGroupAtType))
	js.Global().Set("pinConversation", js.FuncOf(wrapperConMsg.PinConversation))
	js.Global().Set("getTotalUnreadMsgCount", js.FuncOf(wrapperConMsg.GetTotalUnreadMsgCount))
	js.Global().Set("setConversationBurnDuration", js.FuncOf(wrapperConMsg.SetConversationBurnDuration))
	js.Global().Set("findMessageList", js.FuncOf(wrapperConMsg.FindMessageList))

	js.Global().Set("revokeMessage", js.FuncOf(wrapperConMsg.RevokeMessage))
	js.Global().Set("typingStatusUpdate", js.FuncOf(wrapperConMsg.TypingStatusUpdate))
	js.Global().Set("deleteMessageFromLocalStorage", js.FuncOf(wrapperConMsg.DeleteMessageFromLocalStorage))
	js.Global().Set("deleteMessage", js.FuncOf(wrapperConMsg.DeleteMessage))
	js.Global().Set("deleteAllMsgFromLocalAndSvr", js.FuncOf(wrapperConMsg.DeleteAllMsgFromLocalAndSvr))
	js.Global().Set("deleteAllMsgFromLocal", js.FuncOf(wrapperConMsg.DeleteAllMsgFromLocal))
	js.Global().Set("clearConversationAndDeleteAllMsg", js.FuncOf(wrapperConMsg.ClearConversationAndDeleteAllMsg))
	js.Global().Set("insertSingleMessageToLocalStorage", js.FuncOf(wrapperConMsg.InsertSingleMessageToLocalStorage))
	js.Global().Set("insertGroupMessageToLocalStorage", js.FuncOf(wrapperConMsg.InsertGroupMessageToLocalStorage))
	js.Global().Set("searchLocalMessages", js.FuncOf(wrapperConMsg.SearchLocalMessages))
	js.Global().Set("setMessageLocalEx", js.FuncOf(wrapperConMsg.SetMessageLocalEx))

	//register group funcation
	wrapperGroup := wasm_wrapper.NewWrapperGroup(globalFuc)
	js.Global().Set("createGroup", js.FuncOf(wrapperGroup.CreateGroup))
	js.Global().Set("getSpecifiedGroupsInfo", js.FuncOf(wrapperGroup.GetSpecifiedGroupsInfo))
	js.Global().Set("joinGroup", js.FuncOf(wrapperGroup.JoinGroup))
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
	js.Global().Set("getSpecifiedGroupMembersInfo", js.FuncOf(wrapperGroup.GetSpecifiedGroupMembersInfo))
	js.Global().Set("kickGroupMember", js.FuncOf(wrapperGroup.KickGroupMember))
	js.Global().Set("transferGroupOwner", js.FuncOf(wrapperGroup.TransferGroupOwner))
	js.Global().Set("inviteUserToGroup", js.FuncOf(wrapperGroup.InviteUserToGroup))
	js.Global().Set("getGroupApplicationListAsRecipient", js.FuncOf(wrapperGroup.GetGroupApplicationListAsRecipient))
	js.Global().Set("getGroupApplicationListAsApplicant", js.FuncOf(wrapperGroup.GetGroupApplicationListAsApplicant))
	js.Global().Set("acceptGroupApplication", js.FuncOf(wrapperGroup.AcceptGroupApplication))
	js.Global().Set("refuseGroupApplication", js.FuncOf(wrapperGroup.RefuseGroupApplication))
	js.Global().Set("setGroupMemberNickname", js.FuncOf(wrapperGroup.SetGroupMemberNickname))
	js.Global().Set("searchGroupMembers", js.FuncOf(wrapperGroup.SearchGroupMembers))
	js.Global().Set("isJoinGroup", js.FuncOf(wrapperGroup.IsJoinGroup))

	wrapperUser := wasm_wrapper.NewWrapperUser(globalFuc)
	js.Global().Set("getSelfUserInfo", js.FuncOf(wrapperUser.GetSelfUserInfo))
	js.Global().Set("setSelfInfo", js.FuncOf(wrapperUser.SetSelfInfo))
	js.Global().Set("getUsersInfo", js.FuncOf(wrapperUser.GetUsersInfo))

	wrapperFriend := wasm_wrapper.NewWrapperFriend(globalFuc)
	js.Global().Set("getSpecifiedFriendsInfo", js.FuncOf(wrapperFriend.GetSpecifiedFriendsInfo))
	js.Global().Set("getFriendList", js.FuncOf(wrapperFriend.GetFriendList))
	js.Global().Set("searchFriends", js.FuncOf(wrapperFriend.SearchFriends))
	js.Global().Set("checkFriend", js.FuncOf(wrapperFriend.CheckFriend))
	js.Global().Set("addFriend", js.FuncOf(wrapperFriend.AddFriend))
	js.Global().Set("setFriendRemark", js.FuncOf(wrapperFriend.SetFriendRemark))
	js.Global().Set("deleteFriend", js.FuncOf(wrapperFriend.DeleteFriend))
	js.Global().Set("getFriendApplicationListAsRecipient", js.FuncOf(wrapperFriend.GetFriendApplicationListAsRecipient))
	js.Global().Set("getFriendApplicationListAsApplicant", js.FuncOf(wrapperFriend.GetFriendApplicationListAsApplicant))
	js.Global().Set("acceptFriendApplication", js.FuncOf(wrapperFriend.AcceptFriendApplication))
	js.Global().Set("refuseFriendApplication", js.FuncOf(wrapperFriend.RefuseFriendApplication))
	js.Global().Set("getBlackList", js.FuncOf(wrapperFriend.GetBlackList))
	js.Global().Set("removeBlack", js.FuncOf(wrapperFriend.RemoveBlack))
	js.Global().Set("addBlack", js.FuncOf(wrapperFriend.AddBlack))

	//wrapperSignaling := wasm_wrapper.NewWrapperSignaling(globalFuc)
	//js.Global().Set("signalingInviteInGroup", js.FuncOf(wrapperSignaling.SignalingInviteInGroup))
	//js.Global().Set("signalingInvite", js.FuncOf(wrapperSignaling.SignalingInvite))
	//js.Global().Set("signalingAccept", js.FuncOf(wrapperSignaling.SignalingAccept))
	//js.Global().Set("signalingReject", js.FuncOf(wrapperSignaling.SignalingReject))
	//js.Global().Set("signalingCancel", js.FuncOf(wrapperSignaling.SignalingCancel))
	//js.Global().Set("signalingHungUp", js.FuncOf(wrapperSignaling.SignalingHungUp))

	wrapperThird := wasm_wrapper.NewWrapperThird(globalFuc)
	js.Global().Set("updateFcmToken", js.FuncOf(wrapperThird.UpdateFcmToken))
	js.Global().Set("uploadFile", js.FuncOf(wrapperThird.UploadFile))

}
