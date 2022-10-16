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
	caller *ReflectCall
}

func NewWrapperConMsg(wrapperCommon *WrapperCommon) *WrapperConMsg {
	return &WrapperConMsg{WrapperCommon: wrapperCommon, caller: &ReflectCall{}}
}

func (w *WrapperConMsg) CreateTextMessage(_ js.Value, args []js.Value) interface{} {
	return js.ValueOf(w.caller.InitData(open_im_sdk.CreateTextMessage, nil, &args).Call())
}
func (w *WrapperConMsg) CreateImageMessage(_ js.Value, args []js.Value) interface{} {
	return js.ValueOf(w.caller.InitData(open_im_sdk.CreateImageMessage, nil, &args).Call())
}
func (w *WrapperConMsg) CreateImageMessageByURL(_ js.Value, args []js.Value) interface{} {
	return js.ValueOf(w.caller.InitData(open_im_sdk.CreateImageMessageByURL, nil, &args).Call())
}
func (w *WrapperConMsg) CreateCustomMessage(_ js.Value, args []js.Value) interface{} {
	return js.ValueOf(w.caller.InitData(open_im_sdk.CreateCustomMessage, nil, &args).Call())
}
func (w *WrapperConMsg) CreateQuoteMessage(_ js.Value, args []js.Value) interface{} {
	return js.ValueOf(w.caller.InitData(open_im_sdk.CreateQuoteMessage, nil, &args).Call())
}
func (w *WrapperConMsg) CreateAdvancedQuoteMessage(_ js.Value, args []js.Value) interface{} {
	return js.ValueOf(w.caller.InitData(open_im_sdk.CreateAdvancedQuoteMessage, nil, &args).Call())
}
func (w *WrapperConMsg) CreateAdvancedTextMessage(_ js.Value, args []js.Value) interface{} {
	return js.ValueOf(w.caller.InitData(open_im_sdk.CreateAdvancedTextMessage, nil, &args).Call())
}

func (w *WrapperConMsg) MarkC2CMessageAsRead(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	w.caller.InitData(open_im_sdk.MarkC2CMessageAsRead, callback, &args).Call()
	return callback.HandlerFunc()
}
func (w *WrapperConMsg) MarkMessageAsReadByConID(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	w.caller.InitData(open_im_sdk.MarkMessageAsReadByConID, callback, &args).Call()
	return callback.HandlerFunc()
}
func (w *WrapperConMsg) SendMessage(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewSendMessageCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	w.caller.InitData(open_im_sdk.SendMessage, callback, &args).Call()
	return callback.HandlerFunc()
}
func (w *WrapperConMsg) SendMessageNotOss(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewSendMessageCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	w.caller.InitData(open_im_sdk.SendMessageNotOss, callback, &args).Call()
	return callback.HandlerFunc()
}

//------------------------------------conversation---------------------------

func (w *WrapperConMsg) GetAllConversationList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	w.caller.InitData(open_im_sdk.GetAllConversationList, callback, &args).Call()
	return callback.HandlerFunc()
}
func (w *WrapperConMsg) GetOneConversation(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	w.caller.InitData(open_im_sdk.GetOneConversation, callback, &args).Call()
	return callback.HandlerFunc()
}
func (w *WrapperConMsg) DeleteConversationFromLocalAndSvr(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	w.caller.InitData(open_im_sdk.DeleteConversationFromLocalAndSvr, callback, &args).Call()
	return callback.HandlerFunc()
}
func (w *WrapperConMsg) GetAdvancedHistoryMessageList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	w.caller.InitData(open_im_sdk.GetAdvancedHistoryMessageList, callback, &args).Call()
	return callback.HandlerFunc()
}
func (w *WrapperConMsg) GetHistoryMessageList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	w.caller.InitData(open_im_sdk.GetHistoryMessageList, callback, &args).Call()
	return callback.HandlerFunc()
}
