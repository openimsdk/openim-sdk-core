package wasm_wrapper

import (
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/event_listener"
	"syscall/js"
)

//------------------------------------group---------------------------
type WrapperUser struct {
	*WrapperCommon
	caller event_listener.Caller
}

func NewWrapperUser(wrapperCommon *WrapperCommon) *WrapperUser {
	return &WrapperUser{WrapperCommon: wrapperCommon, caller: &event_listener.ReflectCall{}}
}

func (w *WrapperUser) GetSelfUserInfo(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return w.caller.NewCaller(open_im_sdk.GetSelfUserInfo, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperUser) SetSelfInfo(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return w.caller.NewCaller(open_im_sdk.SetSelfInfo, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperUser) GetUsersInfo(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return w.caller.NewCaller(open_im_sdk.GetUsersInfo, callback, &args).AsyncCallWithCallback()
}

