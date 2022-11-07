// +build js,wasm

package wasm_wrapper

import (
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/event_listener"
	"syscall/js"
)

//------------------------------------group---------------------------
type WrapperGroup struct {
	*WrapperCommon
}

func NewWrapperGroup(wrapperCommon *WrapperCommon) *WrapperGroup {
	return &WrapperGroup{WrapperCommon: wrapperCommon}
}
func (w *WrapperGroup) GetGroupsInfo(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetGroupsInfo, callback, &args).AsyncCallWithCallback()
}
