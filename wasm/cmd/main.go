package main

import (
	"open_im_sdk/pkg/log"
	"open_im_sdk/wasm/wasm_wrapper"
	"runtime/debug"

	"syscall/js"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			log.Error("MAIN", "panic info is:", r, string(debug.Stack()))
		}
	}()
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
	js.Global().Set("setAppBackgroundStatus", js.FuncOf(wrapperInitLogin.SetAppBackgroundStatus))
	//register conversation and message func
	wrapperConMsg := wasm_wrapper.NewWrapperConMsg(globalFuc)
	js.Global().Set("createTextMessage", js.FuncOf(wrapperConMsg.CreateTextMessage))
	js.Global().Set("createImageMessage", js.FuncOf(wrapperConMsg.CreateImageMessage))
	js.Global().Set("createImageMessageByURL", js.FuncOf(wrapperConMsg.CreateImageMessageByURL))
	js.Global().Set("createSoundMessageByURL", js.FuncOf(wrapperConMsg.CreateSoundMessageByURL))
	js.Global().Set("CreateVideoMessageByURL", js.FuncOf(wrapperConMsg.CreateVideoMessageByURL))
	js.Global().Set("createFileMessageByURL", js.FuncOf(wrapperConMsg.CreateFileMessageByURL))
	js.Global().Set("createCustomMessage", js.FuncOf(wrapperConMsg.CreateCustomMessage))
	js.Global().Set("createQuoteMessage", js.FuncOf(wrapperConMsg.CreateQuoteMessage))
	js.Global().Set("createAdvancedQuoteMessage", js.FuncOf(wrapperConMsg.CreateAdvancedQuoteMessage))
	js.Global().Set("createAdvancedTextMessage", js.FuncOf(wrapperConMsg.CreateAdvancedTextMessage))
	js.Global().Set("markC2CMessageAsRead", js.FuncOf(wrapperConMsg.MarkC2CMessageAsRead))
	js.Global().Set("markMessageAsReadByConID", js.FuncOf(wrapperConMsg.MarkMessageAsReadByConID))
	js.Global().Set("sendMessage", js.FuncOf(wrapperConMsg.SendMessage))
	js.Global().Set("sendMessageNotOss", js.FuncOf(wrapperConMsg.SendMessageNotOss))
	js.Global().Set("newRevokeMessage", js.FuncOf(wrapperConMsg.NewRevokeMessage))
	js.Global().Set("modifyGroupMessageReaction", js.FuncOf(wrapperConMsg.ModifyGroupMessageReaction))
	js.Global().Set("setMessageReactionExtensions", js.FuncOf(wrapperConMsg.SetMessageReactionExtensions))
	js.Global().Set("addMessageReactionExtensions", js.FuncOf(wrapperConMsg.AddMessageReactionExtensions))
	js.Global().Set("deleteMessageReactionExtensions", js.FuncOf(wrapperConMsg.DeleteMessageReactionExtensions))
	js.Global().Set("getMessageListReactionExtensions", js.FuncOf(wrapperConMsg.GetMessageListReactionExtensions))
	js.Global().Set("getMessageListSomeReactionExtensions", js.FuncOf(wrapperConMsg.GetMessageListSomeReactionExtensions))

	js.Global().Set("getAllConversationList", js.FuncOf(wrapperConMsg.GetAllConversationList))
	js.Global().Set("getOneConversation", js.FuncOf(wrapperConMsg.GetOneConversation))
	js.Global().Set("deleteConversationFromLocalAndSvr", js.FuncOf(wrapperConMsg.DeleteConversationFromLocalAndSvr))
	js.Global().Set("getAdvancedHistoryMessageList", js.FuncOf(wrapperConMsg.GetAdvancedHistoryMessageList))
	js.Global().Set("getHistoryMessageList", js.FuncOf(wrapperConMsg.GetHistoryMessageList))
	//register group func
	wrapperGroup := wasm_wrapper.NewWrapperGroup(globalFuc)
	js.Global().Set("getGroupsInfo", js.FuncOf(wrapperGroup.GetGroupsInfo))

	wrapperUser := wasm_wrapper.NewWrapperUser(globalFuc)
	js.Global().Set("getSelfUserInfo", js.FuncOf(wrapperUser.GetSelfUserInfo))

	wrapperThird := wasm_wrapper.NewWrapperThird(globalFuc)
	js.Global().Set("updateFcmToken", js.FuncOf(wrapperThird.UpdateFcmToken))

}
