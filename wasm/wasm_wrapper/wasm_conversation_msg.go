package wasm_wrapper

import (
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/event_listener"
	"syscall/js"
)

//------------------------------------message---------------------------
func CreateTextMessage(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewSendMessageCallback(utils.FirstLower(utils.GetSelfFuncName()), commonFunc)
	checker(callback, &args, 2)
	callback.EventData().SetOperationID(args[0].String())
	return js.ValueOf(open_im_sdk.CreateTextMessage(args[0].String(), args[1].String()))
}
func CreateImageMessage(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewSendMessageCallback(utils.FirstLower(utils.GetSelfFuncName()), commonFunc)
	checker(callback, &args, 2)
	callback.EventData().SetOperationID(args[0].String())
	return js.ValueOf(open_im_sdk.CreateImageMessage(args[0].String(), args[1].String()))
}
func CreateImageMessageByURL(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewSendMessageCallback(utils.FirstLower(utils.GetSelfFuncName()), commonFunc)
	checker(callback, &args, 4)
	callback.EventData().SetOperationID(args[0].String())
	return js.ValueOf(open_im_sdk.CreateImageMessageByURL(args[0].String(), args[1].String(), args[2].String(), args[3].String()))
}
func CreateCustomMessage(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewSendMessageCallback(utils.FirstLower(utils.GetSelfFuncName()), commonFunc)
	checker(callback, &args, 4)
	callback.EventData().SetOperationID(args[0].String())
	return js.ValueOf(open_im_sdk.CreateCustomMessage(args[0].String(), args[1].String(), args[2].String(), args[3].String()))
}
func CreateQuoteMessage(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewSendMessageCallback(utils.FirstLower(utils.GetSelfFuncName()), commonFunc)
	checker(callback, &args, 3)
	callback.EventData().SetOperationID(args[0].String())
	return js.ValueOf(open_im_sdk.CreateQuoteMessage(args[0].String(), args[1].String(), args[2].String()))
}
func CreateAdvancedQuoteMessage(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewSendMessageCallback(utils.FirstLower(utils.GetSelfFuncName()), commonFunc)
	checker(callback, &args, 4)
	callback.EventData().SetOperationID(args[0].String())
	return js.ValueOf(open_im_sdk.CreateAdvancedQuoteMessage(args[0].String(), args[1].String(), args[2].String(), args[3].String()))
}
func CreateAdvancedTextMessage(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewSendMessageCallback(utils.FirstLower(utils.GetSelfFuncName()), commonFunc)
	checker(callback, &args, 3)
	callback.EventData().SetOperationID(args[0].String())
	return js.ValueOf(open_im_sdk.CreateAdvancedTextMessage(args[0].String(), args[1].String(), args[2].String()))
}

func MarkC2CMessageAsRead(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), commonFunc)
	checker(callback, &args, 3)
	callback.EventData().SetOperationID(args[0].String())
	open_im_sdk.MarkC2CMessageAsRead(callback, args[0].String(), args[1].String(), args[2].String())
	return nil
}
func MarkMessageAsReadByConID(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), commonFunc)
	checker(callback, &args, 3)
	callback.EventData().SetOperationID(args[0].String())
	open_im_sdk.MarkMessageAsReadByConID(callback, args[0].String(), args[1].String(), args[2].String())
	return nil
}
func SendMessage(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewSendMessageCallback(utils.FirstLower(utils.GetSelfFuncName()), commonFunc)
	checker(callback, &args, 5)
	callback.EventData().SetOperationID(args[0].String())
	open_im_sdk.SendMessage(callback, args[0].String(), args[1].String(), args[2].String(), args[3].String(), args[4].String())
	return nil
}
func SendMessageNotOss(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewSendMessageCallback(utils.FirstLower(utils.GetSelfFuncName()), commonFunc)
	checker(callback, &args, 5)
	callback.EventData().SetOperationID(args[0].String())
	open_im_sdk.SendMessageNotOss(callback, args[0].String(), args[1].String(), args[2].String(), args[3].String(), args[4].String())
	return nil
}

//------------------------------------conversation---------------------------
func GetAllConversationList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), commonFunc)
	checker(callback, &args, 1)
	callback.EventData().SetOperationID(args[0].String())
	open_im_sdk.GetAllConversationList(callback, args[0].String())
	return nil
}
func GetOneConversation(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), commonFunc)
	checker(callback, &args, 3)
	callback.EventData().SetOperationID(args[0].String())
	open_im_sdk.GetOneConversation(callback, args[0].String(), args[1].Int(), args[2].String())
	return nil
}
func DeleteConversationFromLocalAndSvr(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), commonFunc)
	checker(callback, &args, 2)
	callback.EventData().SetOperationID(args[0].String())
	open_im_sdk.DeleteConversationFromLocalAndSvr(callback, args[0].String(), args[1].String())
	return nil
}
func GetAdvancedHistoryMessageList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), commonFunc)
	checker(callback, &args, 2)
	callback.EventData().SetOperationID(args[0].String())
	open_im_sdk.GetAdvancedHistoryMessageList(callback, args[0].String(), args[1].String())
	return nil
}
func GetHistoryMessageList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), commonFunc)
	checker(callback, &args, 2)
	callback.EventData().SetOperationID(args[0].String())
	open_im_sdk.GetHistoryMessageList(callback, args[0].String(), args[1].String())
	return nil
}
