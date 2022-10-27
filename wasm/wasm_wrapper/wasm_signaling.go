package wasm_wrapper

import (
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/event_listener"
	"syscall/js"
)

//------------------------------------group---------------------------
type WrapperSignaling struct {
	*WrapperCommon
	caller event_listener.Caller
}

func NewWrapperSignaling(wrapperCommon *WrapperCommon) *WrapperSignaling {
	return &WrapperSignaling{WrapperCommon: wrapperCommon, caller: &event_listener.ReflectCall{}}
}

func (w *WrapperSignaling) SignalingInviteInGroup(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return w.caller.NewCaller(open_im_sdk.SignalingInviteInGroup, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperSignaling) SignalingInvite(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return w.caller.NewCaller(open_im_sdk.SignalingInvite, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperSignaling) SignalingAccept(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return w.caller.NewCaller(open_im_sdk.SignalingAccept, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperSignaling) SignalingReject(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return w.caller.NewCaller(open_im_sdk.SignalingReject, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperSignaling) SignalingCancel(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return w.caller.NewCaller(open_im_sdk.SignalingCancel, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperSignaling) SignalingHungUp(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return w.caller.NewCaller(open_im_sdk.SignalingHungUp, callback, &args).AsyncCallWithCallback()
}
