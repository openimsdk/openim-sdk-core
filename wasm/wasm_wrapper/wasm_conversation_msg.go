package wasm_wrapper

import (
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/event_listener"
	"syscall/js"
)

//------------------------------------message---------------------------
type WrapperConMsg struct {
	*WrapperCommon
	caller event_listener.Caller
}

func NewWrapperConMsg(wrapperCommon *WrapperCommon) *WrapperConMsg {
	return &WrapperConMsg{WrapperCommon: wrapperCommon, caller: &event_listener.ReflectCall{}}
}

func (w *WrapperConMsg) CreateTextMessage(_ js.Value, args []js.Value) interface{} {
	return w.caller.NewCaller(open_im_sdk.CreateTextMessage, nil, &args).AsyncCallWithOutCallback()
}
func (w *WrapperConMsg) CreateImageMessage(_ js.Value, args []js.Value) interface{} {
	return w.caller.NewCaller(open_im_sdk.CreateImageMessage, nil, &args).AsyncCallWithOutCallback()
}
func (w *WrapperConMsg) CreateImageMessageByURL(_ js.Value, args []js.Value) interface{} {
	return w.caller.NewCaller(open_im_sdk.CreateImageMessageByURL, nil, &args).AsyncCallWithOutCallback()
}
func (w *WrapperConMsg) CreateCustomMessage(_ js.Value, args []js.Value) interface{} {
	return w.caller.NewCaller(open_im_sdk.CreateCustomMessage, nil, &args).AsyncCallWithOutCallback()
}
func (w *WrapperConMsg) CreateQuoteMessage(_ js.Value, args []js.Value) interface{} {
	return w.caller.NewCaller(open_im_sdk.CreateQuoteMessage, nil, &args).AsyncCallWithOutCallback()
}
func (w *WrapperConMsg) CreateAdvancedQuoteMessage(_ js.Value, args []js.Value) interface{} {
	return w.caller.NewCaller(open_im_sdk.CreateAdvancedQuoteMessage, nil, &args).AsyncCallWithOutCallback()
}
func (w *WrapperConMsg) CreateAdvancedTextMessage(_ js.Value, args []js.Value) interface{} {
	return w.caller.NewCaller(open_im_sdk.CreateAdvancedTextMessage, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) MarkC2CMessageAsRead(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	w.caller.NewCaller(open_im_sdk.MarkC2CMessageAsRead, callback, &args).AsyncCallWithCallback()
	return callback.HandlerFunc()
}
func (w *WrapperConMsg) MarkMessageAsReadByConID(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	w.caller.NewCaller(open_im_sdk.MarkMessageAsReadByConID, callback, &args).AsyncCallWithCallback()
	return callback.HandlerFunc()
}
func (w *WrapperConMsg) SendMessage(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewSendMessageCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	w.caller.NewCaller(open_im_sdk.SendMessage, callback, &args).AsyncCallWithCallback()
	return callback.HandlerFunc()
}
func (w *WrapperConMsg) SendMessageNotOss(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewSendMessageCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	w.caller.NewCaller(open_im_sdk.SendMessageNotOss, callback, &args).AsyncCallWithCallback()
	return callback.HandlerFunc()
}

//------------------------------------conversation---------------------------

func (w *WrapperConMsg) GetAllConversationList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	w.caller.NewCaller(open_im_sdk.GetAllConversationList, callback, &args).AsyncCallWithCallback()
	return callback.HandlerFunc()
}
func (w *WrapperConMsg) GetOneConversation(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	w.caller.NewCaller(open_im_sdk.GetOneConversation, callback, &args).AsyncCallWithCallback()
	return callback.HandlerFunc()
}
func (w *WrapperConMsg) DeleteConversationFromLocalAndSvr(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	w.caller.NewCaller(open_im_sdk.DeleteConversationFromLocalAndSvr, callback, &args).AsyncCallWithCallback()
	return callback.HandlerFunc()
}
func (w *WrapperConMsg) GetAdvancedHistoryMessageList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	w.caller.NewCaller(open_im_sdk.GetAdvancedHistoryMessageList, callback, &args).AsyncCallWithCallback()
	return callback.HandlerFunc()
}
func (w *WrapperConMsg) GetHistoryMessageList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	w.caller.NewCaller(open_im_sdk.GetHistoryMessageList, callback, &args).AsyncCallWithCallback()
	return callback.HandlerFunc()
}
