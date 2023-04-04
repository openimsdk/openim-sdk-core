//go:build js && wasm
// +build js,wasm

package wasm_wrapper

import (
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/event_listener"
	"syscall/js"
)

// ------------------------------------group---------------------------
type WrapperThird struct {
	*WrapperCommon
}

func NewWrapperThird(wrapperCommon *WrapperCommon) *WrapperThird {
	return &WrapperThird{WrapperCommon: wrapperCommon}
}
func (w *WrapperThird) UpdateFcmToken(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.UpdateFcmToken, callback, &args).AsyncCallWithCallback()
}
