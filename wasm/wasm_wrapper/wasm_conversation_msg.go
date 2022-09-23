package wasm_wrapper

import (
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/event_listener"
	"syscall/js"
)

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
func DeleteConversationFromLocalAndSvr(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), commonFunc)
	checker(callback, &args, 2)
	callback.EventData().SetOperationID(args[0].String())
	open_im_sdk.DeleteConversationFromLocalAndSvr(callback, args[0].String(), args[1].String())
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
func CreateTextMessage(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewSendMessageCallback(utils.FirstLower(utils.GetSelfFuncName()), commonFunc)
	checker(callback, &args, 2)
	callback.EventData().SetOperationID(args[0].String())
	return open_im_sdk.CreateTextMessage(args[0].String(), args[1].String())
}
